// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"maps"
	"slices"
)

// checkProblems returns the entries of openapi_operations.yaml that have stopped doing
// anything, in sorted order.
//
// The generated "openapi_operations" section cannot go stale: update-openapi replaces it
// wholesale from the descriptions at openapi_commit. The two hand-written sections are never
// pruned, so an entry that stops doing anything stays until someone notices, and the file then
// reads as though every entry were load-bearing. Nothing else reports them:
//
//   - an override whose name matches no operation is ignored by resolve() in silence, which is
//     what a renamed endpoint leaves behind;
//   - an override that sets what its operation already has changes no generated doc link, so
//     deleting it changes nothing;
//   - a name that both the "operations" section and the generated section list is resolved to
//     the hand-written entry, which hides the generated one and every change to it;
//   - a name that a section lists twice leaves the file ambiguous, because the last entry is
//     the one that decides each field.
//
// Each problem names the section to edit and the entry to look for, and says what the entry
// fails to do, so that the report is the whole instruction.
func checkProblems(m *operationsFile) []string {
	var problems []string
	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	for _, name := range duplicateNames(m.ManualOps) {
		add("operations: %q is listed more than once, so only the last entry is used", name)
	}
	for _, name := range duplicateNames(m.OverrideOps) {
		add("operation_overrides: %q is listed more than once, so the entries overlap and should be merged into one", name)
	}

	manual := make(map[string]bool, len(m.ManualOps))
	for _, op := range m.ManualOps {
		manual[op.Name] = true
	}
	generated := make(map[string]bool, len(m.OpenapiOps))
	for _, op := range m.OpenapiOps {
		generated[op.Name] = true
	}
	// Sorted, so that the report does not depend on the order of the generated section.
	for _, name := range slices.Sorted(maps.Keys(generated)) {
		if manual[name] {
			add("operations: %q is also in openapi_operations, so the hand-written entry replaces the generated one", name)
		}
	}

	for _, override := range m.OverrideOps {
		base := baseOp(m, override.Name)
		switch {
		case base == nil:
			add("operation_overrides: %q overrides no operation in openapi_operations or operations, so it is ignored", override.Name)
		case !changesOperation(base, override):
			add("operation_overrides: %q sets what its operation already has, so it can be deleted", override.Name)
		}
	}

	slices.Sort(problems)
	return problems
}

// duplicateNames returns, in sorted order, the names that ops lists more than once.
func duplicateNames(ops []*operation) []string {
	counts := make(map[string]int, len(ops))
	for _, op := range ops {
		counts[op.Name]++
	}
	var names []string
	for name, count := range counts {
		if count > 1 {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

// baseOp returns the operation that an override patches: the hand-written "operations" entry
// when there is one, because resolve() puts that entry in place of the generated one.
func baseOp(m *operationsFile, name string) *operation {
	for _, ops := range [][]*operation{m.ManualOps, m.OpenapiOps} {
		for _, op := range ops {
			if op.Name == name {
				return op
			}
		}
	}
	return nil
}

// changesOperation reports whether applying override to base changes anything that a generated
// documentation link or the schema check sees, which is the documentation URL as the generated
// comments render it and the list of description files. resolve() applies only those two
// fields, so an override that changes neither has no effect.
func changesOperation(base, override *operation) bool {
	applied := applyOverride(base, override)
	return normalizeDocURL(applied.DocumentationURL) != normalizeDocURL(base.DocumentationURL) ||
		!slices.Equal(applied.OpenAPIFiles, base.OpenAPIFiles)
}

// applyOverride returns a copy of op with the fields that an override may patch applied.
// resolve() and checkProblems() share it, so that the two cannot disagree about what an
// override does.
func applyOverride(op, override *operation) *operation {
	applied := op.clone()
	if override.DocumentationURL != "" {
		applied.DocumentationURL = override.DocumentationURL
	}
	if len(override.OpenAPIFiles) > 0 {
		applied.OpenAPIFiles = append([]string{}, override.OpenAPIFiles...)
	}
	return applied
}
