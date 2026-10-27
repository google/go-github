// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmp"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// fixAction says how a finding is repaired. The checker decides, because only it knows
// which operations agreed on the finding; the fixer turns the action into edits.
type fixAction struct {
	unomit      bool   // Remove the omit option from the json tag.
	unwrap      bool   // Replace a pointer type with the type it points at.
	addOmit     string // Add "omitempty" or "omitzero" to the json tag.
	makePointer bool   // Replace a value type with a pointer to it.
}

// String describes the repair, for the plan that -fix prints before it writes anything.
func (a *fixAction) String() string {
	var parts []string
	if a.unomit {
		parts = append(parts, "remove the omit option from the json tag")
	}
	if a.unwrap {
		parts = append(parts, "replace the pointer type with the type it points at")
	}
	if a.makePointer {
		parts = append(parts, "make the field a pointer")
	}
	if a.addOmit != "" {
		parts = append(parts, "add "+a.addOmit+" to the json tag")
	}
	return strings.Join(parts, "; ")
}

// edit is a byte-range replacement within a Go source file.
type edit struct {
	start, end int
	text       string
}

// insert returns an edit that adds text at the given offset.
func insert(at int, text string) *edit {
	return &edit{start: at, end: at, text: text}
}

// fixEdits returns the edits that repair one finding.
func fixEdits(d *diagnostic) ([]*edit, error) {
	f := d.info
	if d.action == nil || f == nil {
		return nil, errors.New("not repairable automatically")
	}
	if f.shared {
		return nil, errors.New("the declaration lists several fields, which share one type and tag")
	}
	var edits []*edit
	if d.action.unomit {
		if f.omitOff < 0 {
			return nil, errors.New("cannot locate the omit option")
		}
		edits = append(edits, &edit{start: f.omitOff, end: f.omitEnd})
	}
	if d.action.unwrap {
		if f.starOff < 0 {
			return nil, errors.New("cannot locate the pointer")
		}
		edits = append(edits, &edit{start: f.starOff, end: f.starOff + 1})
	}
	if d.action.makePointer {
		if f.typOff < 0 {
			return nil, errors.New("cannot locate the field type")
		}
		edits = append(edits, insert(f.typOff, "*"))
	}
	if d.action.addOmit != "" {
		edits = append(edits, addOmitTag(f, d.action.addOmit))
	}
	if len(edits) == 0 {
		return nil, errors.New("no repair is needed")
	}
	return edits, nil
}

// addOmitTag returns the edit that adds an omit option to the json tag of a field, creating
// the tag when the field has none.
func addOmitTag(f *fieldInfo, omit string) *edit {
	switch {
	case f.tagOff < 0: // No struct tag at all.
		return insert(f.typEnd, fmt.Sprintf(" `json:%q`", f.jsonName+","+omit))
	case f.tagValEnd > f.tagOff: // A json tag: add the option to its value.
		return insert(f.tagValEnd, ","+omit)
	default: // A tag without a json key.
		return insert(f.tagOff+1, fmt.Sprintf("json:%q ", f.jsonName+","+omit))
	}
}

// applyEdits rewrites src. The edits must not overlap.
func applyEdits(src []byte, edits []*edit) ([]byte, error) {
	slices.SortFunc(edits, func(a, b *edit) int {
		return cmp.Or(cmp.Compare(b.start, a.start), cmp.Compare(b.end, a.end))
	})
	out := src
	for _, e := range edits {
		if e.start < 0 || e.end > len(out) || e.start > e.end {
			return nil, fmt.Errorf("invalid edit at %v:%v", e.start, e.end)
		}
		out = slices.Concat(out[:e.start], []byte(e.text), out[e.end:])
	}
	return out, nil
}

// writeEdits writes the edits to the file at path, formatted as go/format formats it. The edits
// must not overlap.
func writeEdits(path string, src []byte, edits []*edit) error {
	edited, err := applyEdits(src, edits)
	if err != nil {
		return err
	}
	formatted, err := format.Source(edited)
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o600)
}

// fileError names a file and why -fix could not rewrite it. The error of the call that failed
// holds the absolute path it was given, which the note does not need: the file is already named
// relative to the checkout.
func fileError(file string, err error) string {
	if pe, ok := errors.AsType[*os.PathError](err); ok {
		return fmt.Sprintf("%v: %v", file, pe.Err)
	}
	return fmt.Sprintf("%v: %v", file, err)
}

// applyFixes rewrites the Go sources of the findings that can be repaired automatically.
// It returns the number of findings repaired, the number of files rewritten, and a note for
// each finding, or file, it could not repair. A file that cannot be rewritten is noted rather
// than fatal, so that it does not leave the findings in every file after it unrepaired.
func applyFixes(repo string, diags []*diagnostic) (fixed, written int, notes []string) {
	byFile := map[string][]*diagnostic{}
	for _, d := range diags {
		if d.action == nil || d.info == nil {
			continue
		}
		byFile[d.file] = append(byFile[d.file], d)
	}

	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	slices.Sort(files)

	for _, file := range files {
		path := filepath.Join(repo, file)
		src, err := os.ReadFile(path)
		if err != nil {
			notes = append(notes, fileError(file, err))
			continue
		}
		var edits []*edit
		var repaired []*diagnostic
		for _, d := range byFile[file] {
			fieldEdits, err := fixEdits(d)
			if err != nil {
				notes = append(notes, fmt.Sprintf("%v: %v", d.key(), err))
				continue
			}
			edits = append(edits, fieldEdits...)
			repaired = append(repaired, d)
		}
		if len(edits) == 0 {
			continue
		}
		if err := writeEdits(path, src, edits); err != nil {
			notes = append(notes, fileError(file, err))
			continue
		}
		fixed += len(repaired)
		written++
	}
	slices.Sort(notes)
	return fixed, written, notes
}
