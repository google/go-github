// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"
)

// severity is how much a finding matters.
type severity int

const (
	sevWarn severity = iota
	sevError
)

func (s severity) String() string {
	if s == sevError {
		return "error"
	}
	return "warn"
}

// level returns the GitHub Actions annotation level for the severity. GitHub does not
// understand "warn".
func (s severity) level() string {
	if s == sevError {
		return "error"
	}
	return "warning"
}

// Stable rule identifiers, so the summary groups findings by cause rather than by message
// text, and so that exceptions stay meaningful as messages change.
const (
	ruleRequiredOmit        = "required-but-omittable"
	ruleRequiredPointer     = "required-but-nullable-pointer"
	ruleOptionalPointerNil  = "optional-pointer-without-omitempty"
	ruleOptionalValueType   = "optional-value-type"
	ruleOptionalNoOmit      = "optional-without-omit"
	ruleNotInSchema         = "not-in-request-schema"
	ruleMissingRequiredProp = "missing-required-property"
)

// ruleDescriptions explains each rule in the summary.
var ruleDescriptions = map[string]string{
	ruleRequiredOmit:        "the schema REQUIRES the property, but the tag lets it be omitted",
	ruleRequiredPointer:     "the schema REQUIRES the property and does not allow null, but the field is a pointer",
	ruleOptionalPointerNil:  "the property is optional, but a nil pointer is sent as null",
	ruleOptionalValueType:   "the property is optional, but a value type cannot be omitted",
	ruleOptionalNoOmit:      "the property is optional, but the field cannot be omitted",
	ruleNotInSchema:         "the Go field is not in the request body schema",
	ruleMissingRequiredProp: "the schema requires a property the Go struct does not have",
}

// diagnostic is one finding.
type diagnostic struct {
	sev     severity
	rule    string
	file    string
	line    int
	owner   string     // Go struct name
	field   string     // Go field name, empty for struct-level findings
	info    *fieldInfo // The field itself, when the finding is about one.
	message string
	action  *fixAction // nil when -fix cannot repair the finding
}

// key returns the "Struct.Field" name used by the exceptions file.
func (d *diagnostic) key() string {
	if d.field == "" {
		return d.owner
	}
	return d.owner + "." + d.field
}

// structUse is one operation that sends a struct as its request body.
type structUse struct {
	op     *opRef
	schema *flat
}

type stats struct {
	files             int
	methods           int
	methodsWithOps    int
	bodyValue         int
	bodyPointer       int
	structsChecked    int
	usesResolved      int
	usesNotInPlan     int
	usesNoJSONBody    int
	usesUnreadable    int
	fieldsChecked     int
	fieldsConditional int
}

type checker struct {
	repo        string
	d           *descriptions
	info        *repoInfo
	diags       []*diagnostic
	byRule      map[string]int
	unannotated []*methodInfo
	// unstructured are the methods whose by-value body is not a struct this repository
	// declares, so there are no fields to check. It is only listed by -verbose.
	unstructured []*methodInfo
	// usesByStruct maps each request body struct to the operations that send it, in the
	// order they were found. A struct can be the body of more than one operation, so a
	// finding is only safe to repair when those operations agree about the field.
	usesByStruct map[string][]*structUse
	stats        stats
}

func (c *checker) add(d *diagnostic) {
	c.diags = append(c.diags, d)
	c.byRule[d.rule]++
}

// run maps every request body struct to the operations that send it, then checks each
// struct against the schemas of those operations.
func (c *checker) run() {
	c.diags = nil
	c.unannotated = nil
	c.unstructured = nil
	c.byRule = map[string]int{}
	c.stats = stats{files: c.info.files}

	c.usesByStruct = map[string][]*structUse{}
	seenUse := map[string]bool{}

	for _, m := range c.info.methods {
		c.stats.methods++
		if len(m.ops) > 0 {
			c.stats.methodsWithOps++
		}
		if m.bodyType == "" {
			continue
		}
		if m.bodyPtr {
			// paramcheck owns the conversion of pointer bodies to by-value ones,
			// and an unconverted type is checked as soon as it is converted.
			c.stats.bodyPointer++
			continue
		}
		si, ok := c.info.structs[m.bodyType]
		if !ok {
			// The body is not a struct this repository declares, so there are no fields to
			// check, and nothing else can check it either. -verbose names these, and the
			// summary leaves them out of the methods it counts.
			c.unstructured = append(c.unstructured, m)
			continue
		}
		c.stats.bodyValue++
		if len(m.ops) == 0 {
			// Without //meta:operation there is no way to find the schema. The
			// metadata tooling requires the annotation anyway, so this is rare.
			c.unannotated = append(c.unannotated, m)
			continue
		}
		for _, op := range m.ops {
			schema, reason, err := c.d.requestSchema(op)
			if err != nil {
				// A description that cannot be read leaves this use unchecked, and the
				// count of unchecked uses says so.
				c.stats.usesUnreadable++
				continue
			}
			switch reason {
			case bodyNotInPlan:
				c.stats.usesNotInPlan++
				continue
			case bodyNoJSONBody:
				c.stats.usesNoJSONBody++
				continue
			}
			key := si.name + "|" + op.String()
			if seenUse[key] {
				continue
			}
			seenUse[key] = true
			c.stats.usesResolved++
			c.usesByStruct[si.name] = append(c.usesByStruct[si.name], &structUse{op: op, schema: schema})
		}
	}

	names := make([]string, 0, len(c.usesByStruct))
	for n := range c.usesByStruct {
		names = append(names, n)
	}
	slices.Sort(names)
	for _, name := range names {
		c.stats.structsChecked++
		c.checkStruct(c.info.structs[name], c.usesByStruct[name])
	}

	slices.SortStableFunc(c.diags, func(a, b *diagnostic) int {
		return cmp.Or(
			strings.Compare(a.file, b.file),
			cmp.Compare(a.line, b.line),
			strings.Compare(a.rule, b.rule),
		)
	})
}

func (c *checker) checkStruct(si *structInfo, uses []*structUse) {
	for _, f := range si.fields {
		c.stats.fieldsChecked++

		// One pass over the operations, counting what they agree on.
		inSchema := false
		checkable, required, nullable := 0, 0, false
		for _, u := range uses {
			meta, ok := u.schema.props[f.jsonName]
			if !ok {
				continue
			}
			inSchema = true
			if meta.readOnly {
				// A readOnly property is never sent, so optionality does not
				// apply to it.
				continue
			}
			checkable++
			if u.schema.required[f.jsonName] {
				required++
			}
			nullable = nullable || meta.nullable
		}

		if !inSchema {
			c.add(&diagnostic{
				sev: sevWarn, rule: ruleNotInSchema, file: f.file, line: f.line,
				owner: si.name, field: f.goName, info: f,
				message: fmt.Sprintf("Go field is not in the request body schema of any of its %v operation(s)", len(uses)),
			})
			continue
		}
		if checkable == 0 {
			continue // Present in the schema, but readOnly in every operation.
		}

		switch required {
		case checkable:
			c.checkRequired(f, si, checkable, len(uses), nullable)
		case 0:
			c.checkOptional(f, si)
		default:
			// Required by some operations and optional in others: conditional
			// requiredness that one Go type cannot express. Leave it to
			// author judgement.
			c.stats.fieldsConditional++
		}
	}

	// Required properties that the struct cannot supply at all.
	for _, name := range missingRequiredProps(si, uses) {
		c.add(&diagnostic{
			sev: sevError, rule: ruleMissingRequiredProp, file: si.file, line: si.line,
			owner: si.name, field: name,
			message: fmt.Sprintf("schema REQUIRES property %q but the Go struct has no field with that JSON name", name),
		})
	}
}

// checkRequired reports a required property that the Go field can leave out. required counts
// the operations whose schema requires the property, and total the operations that send the
// struct: they differ when the other operations' schemas do not have the property, or hold it
// readOnly, so no one Go field can suit them all and -fix leaves the field alone.
func (c *checker) checkRequired(f *fieldInfo, si *structInfo, required, total int, nullable bool) {
	in := fmt.Sprintf("in all %v of its operation(s)", total)
	if required != total {
		in = fmt.Sprintf("in %v of the %v operations that send this struct, which do not all agree", required, total)
	}
	switch {
	case f.hasOmit:
		action := &fixAction{unomit: true, unwrap: f.isPointer && !nullable}
		c.add(&diagnostic{
			sev: sevError, rule: ruleRequiredOmit, file: f.file, line: f.line,
			owner: si.name, field: f.goName, info: f, action: action,
			message: fmt.Sprintf("schema REQUIRES %q %v, but %v makes it omittable, so it can be sent as absent",
				f.jsonName, in, f.omitOption()),
		})
	case f.isPointer && !nullable:
		c.add(&diagnostic{
			sev: sevError, rule: ruleRequiredPointer, file: f.file, line: f.line,
			owner: si.name, field: f.goName, info: f, action: &fixAction{unwrap: true},
			message: fmt.Sprintf("schema REQUIRES %q and does not allow null %v, but the Go field is a pointer, so a nil value is sent as null",
				f.jsonName, in),
		})
	}
}

// agrees reports whether every operation that sends the struct of d agrees about the field, so
// that one Go field can suit all of them. An operation disagrees when its schema does not have
// the property, or holds it readOnly, which is never sent: making the Go field mandatory would
// start sending a property that operation does not accept, so the finding is reported for a
// human to settle and -fix leaves it alone.
//
// Only this kind of disagreement leaves a repair free to contradict another operation. The
// checker settles the other two before it reports a finding: required by some operations and
// optional in others is left to author judgement, and a property that any operation allows to
// be null keeps its pointer, so the repair never unwraps one.
func (c *checker) agrees(d *diagnostic) bool {
	uses := c.usesByStruct[d.owner]
	if d.info == nil || len(uses) < 2 {
		return true // One operation cannot disagree with itself.
	}
	for _, u := range uses {
		if meta, ok := u.schema.props[d.info.jsonName]; !ok || meta.readOnly {
			return false
		}
	}
	return true
}

// checkOptional reports an optional property that the Go field cannot leave out, or that it
// sends as null.
func (c *checker) checkOptional(f *fieldInfo, si *structInfo) {
	switch {
	case f.omits():
		// The tag leaves an unset value out of the body, so the field is correct.
	case f.isPointer:
		c.add(&diagnostic{
			sev: sevWarn, rule: ruleOptionalPointerNil, file: f.file, line: f.line,
			owner: si.name, field: f.goName, info: f, action: &fixAction{addOmit: "omitempty"},
			message: fmt.Sprintf("%q is optional in the schema, but the pointer has no omitempty, so a nil value is sent as null", f.jsonName),
		})
	case !f.omittable:
		// CONTRIBUTING.md asks for a pointer on a field the schema does not require,
		// and a struct value has no alternative: omitempty cannot leave one out, because
		// a struct is never empty, and the structfield linter rejects omitzero on one.
		c.add(&diagnostic{
			sev: sevWarn, rule: ruleOptionalValueType, file: f.file, line: f.line,
			owner: si.name, field: f.goName, info: f, action: &fixAction{makePointer: true, addOmit: "omitempty"},
			message: fmt.Sprintf("%q is optional in the schema, but the Go field is a value type, so it is always sent: make it a pointer with omitempty", f.jsonName),
		})
	default:
		// Slices, maps and types from other packages: omitempty would drop an empty
		// but non-nil value, so CONTRIBUTING.md asks for omitzero.
		c.add(&diagnostic{
			sev: sevWarn, rule: ruleOptionalNoOmit, file: f.file, line: f.line,
			owner: si.name, field: f.goName, info: f, action: &fixAction{addOmit: "omitzero"},
			message: fmt.Sprintf("%q is optional in the schema, but the Go field has no omitzero, so an unset value is indistinguishable from a zero one", f.jsonName),
		})
	}
}

// missingRequiredProps returns, in a stable order, the properties that every one of the
// operations requires and that the struct has no field for.
func missingRequiredProps(si *structInfo, uses []*structUse) []string {
	counts := map[string]int{}
	for _, u := range uses {
		for name := range u.schema.required {
			if u.schema.props[name].readOnly { // Never sent, so it need not be settable.
				continue
			}
			counts[name]++
		}
	}
	var missing []string
	for name, count := range counts {
		if count < len(uses) || si.field(name) != nil {
			continue
		}
		missing = append(missing, name)
	}
	slices.Sort(missing)
	return missing
}

// ------------------------------------------------------------------ exceptions

// exceptions grandfathers the findings that the repository already has, so that a pull
// request is only told about the ones it introduces. An entry is a decision to leave a
// disagreement alone: -fix never rewrites a field that the file names. Entries that are no
// longer needed are reported, so that the file can only shrink.
type exceptions struct {
	path     string
	comments []string
	entries  []string        // In file order, so that a rewrite keeps the file stable.
	known    map[string]bool // The same keys, so that a lookup is not a scan.
}

// setEntries replaces the entries. The ordered slice and the lookup map are built together,
// so that they cannot disagree, and duplicates are dropped.
func (e *exceptions) setEntries(entries []string) {
	e.entries = nil
	e.known = make(map[string]bool, len(entries))
	for _, entry := range entries {
		if e.known[entry] {
			continue
		}
		e.known[entry] = true
		e.entries = append(e.entries, entry)
	}
}

// loadExceptions reads the exceptions file. A missing file is not an error.
func loadExceptions(path string) (*exceptions, error) {
	e := &exceptions{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return e, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []string
	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case strings.HasPrefix(line, "#"):
			e.comments = append(e.comments, line)
		default:
			entries = append(entries, line)
		}
	}
	e.setEntries(entries)
	return e, nil
}

// suppress reports whether key is grandfathered.
func (e *exceptions) suppress(key string) bool {
	return e != nil && e.known[key]
}

// obsolete returns the entries that no current finding needs.
func (e *exceptions) obsolete(diags []*diagnostic) []string {
	if e == nil {
		return nil
	}
	needed := make(map[string]bool, len(diags))
	for _, d := range diags {
		needed[d.key()] = true
	}
	var out []string
	for _, entry := range e.entries {
		if !needed[entry] {
			out = append(out, entry)
		}
	}
	return out
}

// write replaces the exceptions file with the given entries.
func (e *exceptions) write(entries []string) error {
	comments := e.comments
	if len(comments) == 0 {
		comments = []string{
			"# Schema field exceptions for tools/schemafields.",
			"#",
			"# Each line names a Go struct field whose request-body optionality disagrees with",
			"# GitHub's OpenAPI descriptions. The findings are grandfathered so that CI only",
			"# reports problems that a change introduces, and so that a disagreement can be",
			"# accepted deliberately.",
			"#",
			"# A line here is a decision to leave the field alone, so -fix never rewrites a",
			"# field that this file names. Delete the line to have a later -fix repair it.",
			"#",
			"# Run \"script/check-schema-fields.sh -fix\" to repair the error-severity findings",
			"# that no line here grandfathers. It drops the lines that no finding needs any",
			"# more, so the file can only shrink.",
		}
	}
	var b strings.Builder
	for _, comment := range comments {
		b.WriteString(comment)
		b.WriteByte('\n')
	}
	for _, entry := range entries {
		b.WriteString(entry)
		b.WriteByte('\n')
	}
	return os.WriteFile(e.path, []byte(b.String()), 0o600)
}
