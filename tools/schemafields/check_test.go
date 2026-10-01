// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	fixtureRepo         = "testdata/repo"
	fixtureDescriptions = "testdata/descriptions/api.github.com.json,testdata/descriptions/ghec.json"
)

// wantFindings is every finding that the fixture produces, in report order.
const wantFindings = `WARN  github/enterprise.go:17: CreateEnterpriseRunnerGroupRequest.SelectedWorkflows: "selected_workflows" is optional in the schema, but the Go field has no omitzero, so an unset value is indistinguishable from a zero one
WARN  github/issues.go:17: CreateCommentRequest.Extra: Go field is not in the request body schema of any of its 1 operation(s)
ERROR github/runner_groups.go:14: CreateRunnerGroupRequest.Name: schema REQUIRES "name" in all 1 of its operation(s), but omitempty makes it omittable, so it can be sent as absent
WARN  github/runner_groups.go:20: CreateRunnerGroupRequest.SelectedRepositoryIDs: "selected_repository_ids" is optional in the schema, but the Go field has no omitzero, so an unset value is indistinguishable from a zero one
ERROR github/runner_groups.go:29: UpdateRunnerGroupRequest.Name: schema REQUIRES "name" and does not allow null in all 1 of its operation(s), but the Go field is a pointer, so a nil value is sent as null
WARN  github/runner_groups.go:41: ValueTypeRequest.Count: "count" is optional in the schema, but the Go field is a value type, so it is always sent: make it a pointer with omitempty
WARN  github/runner_groups.go:47: StructTypeRequest.Inner: "inner" is optional in the schema, but the Go field is a value type, so it is always sent: make it a pointer with omitempty
ERROR github/runner_groups.go:71: OneOfRequest.KindTag: schema REQUIRES "kind_tag" in all 1 of its operation(s), but omitempty makes it omittable, so it can be sent as absent
ERROR github/runner_groups.go:78: MissingRequest.required_thing: schema REQUIRES property "required_thing" but the Go struct has no field with that JSON name
ERROR github/runner_groups.go:88: SharedRequest.Name: schema REQUIRES "name" in 1 of the 2 operations that send this struct, which do not all agree, but omitempty makes it omittable, so it can be sent as absent
ERROR github/runner_groups.go:92: SharedRequest.Code: schema REQUIRES "code" in 1 of the 2 operations that send this struct, which do not all agree, but omitempty makes it omittable, so it can be sent as absent
`

// fixtureGoMod makes a copy of the fixture a buildable module, which is what -fix compiles to
// find the call sites a changed field type leaves unable to compile. The fixture cannot carry a
// go.mod of its own: script/lint.sh and test-all.sh read the modules to lint and test from
// "git ls-files '*go.mod'", so a tracked go.mod under testdata would be linted as a real module.
const fixtureGoMod = "module github.com/google/go-github/v92\n\ngo 1.26.0\n"

// copyFixture copies the fixture checkout into dir and makes the copy a module.
func copyFixture(t *testing.T, dir string) {
	t.Helper()
	copyTree(t, fixtureRepo, dir)
	mustWriteFile(t, filepath.Join(dir, "go.mod"), fixtureGoMod)
}

// runOnFixture copies the fixture checkout into a temporary directory and runs the tool
// over the copy, so that -fix cannot touch the testdata.
func runOnFixture(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	dir := t.TempDir()
	copyFixture(t, dir)
	return runIn(t, dir, args...)
}

// runIn runs the tool over the checkout at dir. The extra arguments are appended after the
// ones that point the tool at the fixture's descriptions.
func runIn(t *testing.T, dir string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	all := append([]string{"-repo", dir, "-descriptions", fixtureDescriptions}, args...)
	var out, errOut bytes.Buffer
	err = run(&out, &errOut, all)
	return out.String(), errOut.String(), err
}

func TestCheck(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := runOnFixture(t)
	if err == nil {
		t.Fatal("expected an error for the fixture findings")
	}
	assertContains(t, err.Error(), "6 schema field issue(s) found")
	assertEqual(t, wantFindings, stdout)
	assertContains(t, stderr, "checked 17 of 18 methods that take a struct request body (13 body structs, 22 fields)")
	assertContains(t, stderr, "not checked: 1 method whose body is passed by pointer (run paramcheck to convert it)")
	assertContains(t, stderr, "not checked: 1 operation use with no JSON request body")
	assertContains(t, stderr, "11 findings (5 warn, 6 error)")
	assertContains(t, stderr, "repairable by -fix: 3 new")
	assertContains(t, stderr, "3 of the errors are not repairable by -fix")
	// The fixture has no exceptions file, so nothing of it is grandfathered.
	assertNotContains(t, stderr, "grandfathered")
}

// TestSummary pins what the summary says, which is the whole of what a maintainer reads: how
// much of the checkout was checked, what was not checked and why, what the findings are
// including the ones an exceptions file grandfathers, and what -fix can and cannot do.
func TestSummary(t *testing.T) {
	t.Parallel()
	// A checker shaped like a large checkout with a baseline: 40 findings all accepted, one of
	// them repairable by -fix, and a tenth of the body structs unreachable because the pinned
	// revision predates them.
	c := &checker{
		repo: "/checkout",
		stats: stats{
			files: 209, methods: 7981, methodsWithOps: 1301,
			bodyValue: 206, bodyPointer: 62,
			structsChecked: 141, fieldsChecked: 618, fieldsConditional: 13,
			usesResolved: 193, usesNotInPlan: 10,
		},
		byRule: map[string]int{ruleNotInSchema: 1, ruleRequiredOmit: 1, ruleMissingRequiredProp: 1},
	}
	c.diags = []*diagnostic{
		{sev: sevWarn, rule: ruleNotInSchema, file: "github/x.go", line: 1, owner: "Req", field: "Extra"},
		{
			sev: sevError, rule: ruleRequiredOmit, file: "github/x.go", line: 2, owner: "Req", field: "Name",
			info: &fieldInfo{jsonName: "name"}, action: &fixAction{unomit: true},
		},
		// An error that no repair expresses: the schema requires a property the struct has
		// no field for, so only a person can settle it.
		{sev: sevError, rule: ruleMissingRequiredProp, file: "github/x.go", line: 3, owner: "Req", field: "Gone"},
	}
	exc := &exceptions{path: "/checkout/tools/schemafields/exceptions.txt"}
	exc.setEntries([]string{"Req.Extra", "Req.Name", "Req.Gone"})

	// The coverage lines are the same whatever the findings are.
	const coverage = `checked 206 of 268 methods that take a struct request body (141 body structs, 618 fields)
  not checked: 62 methods whose body is passed by pointer (run paramcheck to convert them)
  not checked: 10 operation uses that the pinned revision does not document
    a newer openapi_commit in openapi_operations.yaml would check them
`

	var buf bytes.Buffer
	printSummary(&buf, c, &options{}, exc)
	assertEqual(t, "\n"+coverage+
		"3 findings (1 warn, 2 error), all grandfathered by tools/schemafields/exceptions.txt\n"+
		"  repairable by -fix: 0 new, 1 grandfathered (drop the entries for -fix to repair)\n"+
		"  1 of the errors is not repairable by -fix\n", buf.String())

	// The same tree with no baseline: the findings are new rather than grandfathered, and the
	// repair that -fix would make is counted as new.
	buf.Reset()
	printSummary(&buf, c, &options{}, &exceptions{path: "/checkout/tools/schemafields/exceptions.txt"})
	assertEqual(t, "\n"+coverage+
		"3 findings (1 warn, 2 error)\n"+
		"  repairable by -fix: 1 new\n"+
		"  1 of the errors is not repairable by -fix\n", buf.String())

	// -verbose adds the scan counts, the operation uses that resolved, the fields that a rule
	// leaves alone, and the per-rule tally.
	buf.Reset()
	printSummary(&buf, c, &options{verbose: true}, exc)
	got := buf.String()
	assertContains(t, got, "scanned 209 files, 7981 methods (1301 with //meta:operation)")
	assertContains(t, got, "resolved operation uses: 193")
	assertContains(t, got, "conditionally required, and left alone: 13 of the 618 fields")
	assertContains(t, got, "1  not-in-request-schema: the Go field is not in the request body schema")

	// A checkout with nothing to report says so, rather than counting a baseline of nothing.
	buf.Reset()
	c.diags = nil
	c.byRule = map[string]int{}
	printSummary(&buf, c, &options{}, exc)
	assertContains(t, buf.String(), "no findings")
}

func TestCheckMinSeverity(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := runOnFixture(t, "-min-severity", "error")
	if err == nil {
		t.Fatal("expected an error for the fixture findings")
	}
	if strings.Contains(stdout, "WARN") {
		t.Errorf("warnings were reported at -min-severity=error:\n%v", stdout)
	}
	assertContains(t, stdout, "CreateRunnerGroupRequest.Name")
	// The summary still describes the tree rather than the threshold, and says that the
	// warnings it counts are not the ones it printed.
	assertContains(t, stderr, "11 findings (5 warn, 6 error)")
	assertContains(t, stderr, "warnings are hidden; remove -min-severity=error to see them")
}

func TestCheckGithubFormat(t *testing.T) {
	t.Parallel()
	stdout, _, err := runOnFixture(t, "-format", "github", "-min-severity", "error")
	if err == nil {
		t.Fatal("expected an error for the fixture findings")
	}
	// GitHub Actions annotations escape commas and colons in the message.
	want := `::error file=github/runner_groups.go,line=14,title=CreateRunnerGroupRequest.Name::schema REQUIRES "name" in all 1 of its operation(s)%2C but omitempty makes it omittable%2C so it can be sent as absent`
	assertContains(t, stdout, want)
}

func TestCheckVerbose(t *testing.T) {
	t.Parallel()
	_, stderr, _ := runOnFixture(t, "-verbose")
	assertContains(t, stderr, "no //meta:operation annotation: github/issues.go:54: Untracked (UntrackedRequest)")
	// -verbose is the full breakdown of the summary, which the default output leaves out.
	assertContains(t, stderr, "scanned 4 files, 18 methods (16 with //meta:operation)")
	assertContains(t, stderr, "resolved operation uses: 15")
	assertContains(t, stderr, "conditionally required, and left alone: 1 of the 22 fields")
	assertContains(t, stderr, "4  required-but-omittable: the schema REQUIRES the property, but the tag lets it be omitted")
	// The default output leaves the breakdown out, so that it reads as a verdict.
	if _, plain, _ := runOnFixture(t); strings.Contains(plain, "scanned 4 files") {
		t.Errorf("-verbose output was printed without -verbose:\n%v", plain)
	}
}

func TestCheckUnknownFlagValue(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"-format", "junit"},
		{"-min-severity", "fatal"},
		{"unexpected"},
	} {
		var stdout, stderr bytes.Buffer
		if err := run(&stdout, &stderr, args); err == nil {
			t.Errorf("run(%v) accepted an invalid argument", args)
		}
	}
}

// TestCheckHelp checks that -h prints the usage and succeeds, rather than reporting that the
// usage it just printed is an error.
func TestCheckHelp(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	assertNilError(t, run(&stdout, &stderr, []string{"-h"}))
	assertContains(t, stderr.String(), "Usage: schemafields")
	assertEqual(t, "", stdout.String())
}

func TestCheckMissingRepo(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	err := run(&stdout, &stderr, []string{"-repo", filepath.Join(t.TempDir(), "absent")})
	if err == nil {
		t.Fatal("expected an error for a -repo that does not exist")
	}
	assertContains(t, err.Error(), "absent")
}

func TestCheckMissingDescriptions(t *testing.T) {
	t.Parallel()
	// A checkout with sources but no openapi_operations.yaml, so that the run gets past
	// the scan and fails while looking for the pinned descriptions.
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "github", "doc.go"), "package github\n")
	var stdout, stderr bytes.Buffer
	err := run(&stdout, &stderr, []string{"-repo", dir})
	if err == nil {
		t.Fatal("expected an error when the checkout has no openapi_operations.yaml")
	}
	assertContains(t, err.Error(), "openapi_operations.yaml")
}

func TestExceptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")
	// Every current finding, so that the run itself is clean, plus one that nothing needs.
	mustWriteFile(t, exceptions, wantExceptionFile+"Gone.Field\n")

	stdout, stderr, err := runIn(t, dir, "-exceptions", exceptions)
	if err == nil {
		t.Fatal("expected an error for the obsolete exception")
	}
	assertContains(t, err.Error(), "1 obsolete exception(s) to remove")
	assertEqual(t, "", stdout)
	assertContains(t, stderr, "obsolete exceptions")
	assertContains(t, stderr, "  Gone.Field")
}

// exceptionHeader is the comment block that -write-exceptions writes when the file has none.
const exceptionHeader = `# Schema field exceptions for tools/schemafields.
#
# Each line names a Go struct field whose request-body optionality disagrees with
# GitHub's OpenAPI descriptions. The findings are grandfathered so that CI only
# reports problems that a change introduces, and so that a disagreement can be
# accepted deliberately.
#
# A line here is a decision to leave the field alone, so -fix never rewrites a
# field that this file names. Delete the line to have a later -fix repair it.
#
# Run "script/check-schema-fields.sh -fix" to repair the error-severity findings
# that no line here grandfathers. It drops the lines that no finding needs any
# more, so the file can only shrink.
`

// wantExceptionFile is the exceptions file that -write-exceptions writes for the fixture:
// the default header, then every finding's key in sorted order.
const wantExceptionFile = exceptionHeader + `CreateCommentRequest.Extra
CreateEnterpriseRunnerGroupRequest.SelectedWorkflows
CreateRunnerGroupRequest.Name
CreateRunnerGroupRequest.SelectedRepositoryIDs
MissingRequest.required_thing
OneOfRequest.KindTag
SharedRequest.Code
SharedRequest.Name
StructTypeRequest.Inner
UpdateRunnerGroupRequest.Name
ValueTypeRequest.Count
`

func TestWriteExceptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")

	// Every finding becomes an exception, so the run that regenerated the file is
	// clean and can be used to set a baseline.
	_, stderr, err := runIn(t, dir, "-exceptions", exceptions, "-write-exceptions")
	assertNilError(t, err)
	assertContains(t, stderr, "wrote 11 exception(s)")

	data, readErr := os.ReadFile(exceptions)
	assertNilError(t, readErr)
	assertEqual(t, wantExceptionFile, string(data))
}

// TestFix repairs the fixture and compares each rewritten file with its golden copy.
//
// -fix repairs the findings that fail the check, so it rewrites only the three errors that
// can be repaired mechanically, and leaves the five warnings, the error that no repair can
// express, and the two errors on a struct whose operations disagree alone. Two of the repairs
// change a field type, so it repairs the call sites that the compiler then reports as well.
func TestFix(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the finding that cannot be repaired")
	}
	assertContains(t, stderr, "planned repairs (3):")
	assertContains(t, stderr, "planned call site repairs (3):")
	assertContains(t, stderr, `repaired 3 call site(s) in 1 Go file(s)`)
	assertContains(t, stderr, "repaired 3 finding(s) in 1 Go file(s)")
	assertContains(t, stderr, "8 finding(s) remain, 0 of them repairable by -fix")
	// The fields of the struct that two operations send with schemas that disagree stay as
	// they are, even though nothing grandfathers them: -fix neither plans nor notes a repair
	// for them.
	assertNotContains(t, stderr, "SharedRequest")

	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups.go")
	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")

	// A file that holds only warnings and an unrepairable error comes back unchanged.
	assertSameAsFixture(t, dir, "github/enterprise.go", "github/issues.go")
}

// TestFixWithDefaultRepo checks that -fix finds the call sites when -repo is left at its
// default, ".", which is how script/check-schema-fields.sh runs it. A checkout named "." has to
// reach the module walk that decides which modules to compile: when it does not, no call site is
// repaired, and the run reports the finding as repaired all the same, leaving a tree that no
// longer builds.
//
// The test changes the working directory, which the testing package allows only when no other
// test runs beside it, and the others resolve the fixture's descriptions relative to it.
//
//nolint:paralleltest // cannot use t.Parallel() when the test calls t.Chdir
func TestFixWithDefaultRepo(t *testing.T) {
	dir := t.TempDir()
	copyFixture(t, dir)
	// The descriptions are resolved against the working directory, which the test changes.
	descriptions := absolutePaths(t, fixtureDescriptions)
	t.Chdir(dir)

	var stdout, stderr bytes.Buffer
	err := run(&stdout, &stderr, []string{"-descriptions", descriptions, "-fix"})
	if err == nil {
		t.Fatal("expected an error for the finding that cannot be repaired")
	}
	assertContains(t, stderr.String(), "planned call site repairs (3):")
	assertContains(t, stderr.String(), `repaired 3 call site(s) in 1 Go file(s)`)
	assertNotContains(t, stderr.String(), "no module found")

	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups.go")
	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")
}

// absolutePaths returns a comma-separated flag value with each path made absolute, for a test
// that changes the working directory and so can no longer resolve them relative to it.
func absolutePaths(t *testing.T, value string) string {
	t.Helper()
	wd, err := os.Getwd()
	assertNilError(t, err)
	paths := strings.Split(value, ",")
	for i, path := range paths {
		paths[i] = filepath.Join(wd, filepath.FromSlash(path))
	}
	return strings.Join(paths, ",")
}

// TestFixReportsCallSitesItCannotRepair checks that a call site which the change breaks and
// which no mechanical repair can express is reported, rather than left in a tree that no
// longer compiles.
func TestFixReportsCallSitesItCannotRepair(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	// A literal that passes a variable, which -fix cannot unwrap the way it unwraps new(...).
	mustWriteFile(t, filepath.Join(dir, "github", "callers_extra.go"), `package github

// ptrName is a variable, so the repair of a literal that uses it is not mechanical.
var ptrName = new("d")

// PtrCaller passes the variable to the field whose type the repair changes.
var PtrCaller = UpdateRunnerGroupRequest{Name: ptrName}
`)
	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the call site that cannot be repaired")
	}
	assertContains(t, err.Error(), "1 repair(s) could not be applied")
	assertContains(t, stderr, "not repaired: github/callers_extra.go:")
	assertContains(t, stderr, "UpdateRunnerGroupRequest.Name is now string")
	assertContains(t, stderr, "is not new(...), Ptr(...) or &...")
}

// TestFixLeavesGeneratedCallSitesAlone checks that -fix neither repairs a call site in a
// generated file nor abandons the plan because of it. The generators that own such a file
// rewrite it, and chmod their output read-only, so a repair written into one would be refused
// and would be discarded by the generator's next run. The fixture has no script/generate.sh for
// -fix to run, so the file is reported and the run fails, rather than the tree being left broken
// quietly.
func TestFixLeavesGeneratedCallSitesAlone(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	// A call site in a file that a generator owns, marked as go/ast marks one, and read-only
	// the way the repository's own generators leave their output.
	const generated = `// Code generated by gen-accessors; DO NOT EDIT.

package github

// GeneratedCaller builds a request the way a generated test does.
var GeneratedCaller = UpdateRunnerGroupRequest{Name: new("g")}
`
	path := filepath.Join(dir, "github", "generated_callers.go")
	mustWriteFile(t, path, generated)
	assertNilError(t, os.Chmod(path, 0o444))

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the generated file that -fix left alone")
	}
	assertContains(t, err.Error(), "1 repair(s) could not be applied")
	assertContains(t, stderr, "not repaired: github/generated_callers.go is generated, and this checkout has no script/generate.sh to rewrite it")
	// The file is named once, though every pass of the repair loop compiles it again.
	assertEqual(t, 1, strings.Count(stderr, "is generated, and this checkout has no script/generate.sh"))
	// The plan holds the three literals of the hand-written callers and no other, so the
	// generated file is left out of it rather than tried and refused.
	assertContains(t, stderr, "planned call site repairs (3):")
	assertContains(t, stderr, `repaired 3 call site(s) in 1 Go file(s)`)

	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")
	got, readErr := os.ReadFile(path)
	assertNilError(t, readErr)
	assertEqual(t, generated, string(got))
}

// TestFixRunsTheCheckoutsGenerator checks that -fix runs the checkout's generator when a repair
// changes a field's type, rather than leaving the files the type is generated into for the user
// to regenerate. The generator's output is what makes the tree compile, so -fix reports no
// unwritten repair and the run succeeds as far as the call sites are concerned.
func TestFixRunsTheCheckoutsGenerator(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	// A generator that owns the generated file, as script/generate.sh owns the accessors.
	const generated = `// Code generated by gen-accessors; DO NOT EDIT.

package github

// GeneratedCaller is rewritten by the generator.
var GeneratedCaller = UpdateRunnerGroupRequest{Name: new("g")}
`
	path := filepath.Join(dir, "github", "generated_callers.go")
	mustWriteFile(t, path, generated)
	mustWriteFile(t, filepath.Join(dir, "script", "generate.sh"), `#!/bin/sh
# Stands in for the repository's generators, which rewrite the files they own from the
# declarations: after the repair the field is a string, so the literal in the generated file is.
cat > github/generated_callers.go <<'EOF'
// Code generated by gen-accessors; DO NOT EDIT.

package github

// GeneratedCaller is rewritten by the generator.
var GeneratedCaller = UpdateRunnerGroupRequest{Name: "g"}
EOF
`)

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the findings that cannot be repaired")
	}
	assertContains(t, stderr, "ran script/generate.sh")
	assertNotContains(t, stderr, "is generated")
	assertNotContains(t, err.Error(), "repair(s) could not be applied")
	// The three hand-written literals are repaired, and the generated file compiles as the
	// generator rewrote it.
	assertContains(t, stderr, "planned call site repairs (3):")
	assertContains(t, stderr, `repaired 3 call site(s) in 1 Go file(s)`)
	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")

	got, readErr := os.ReadFile(path)
	assertNilError(t, readErr)
	assertContains(t, string(got), `UpdateRunnerGroupRequest{Name: "g"}`)
}

// TestFixReportsACheckoutThatDoesNotCompile checks that -fix does not report success when the
// checkout does not build for a reason of its own: an unresolved import leaves the package
// unbuilt, so the compiler reports no call site at all, and the run says so rather than
// repairing nothing and calling it done.
func TestFixReportsACheckoutThatDoesNotCompile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	mustWriteFile(t, filepath.Join(dir, "github", "broken.go"), `package github

import _ "example.com/nonexistent"
`)

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for a checkout that does not compile")
	}
	assertContains(t, stderr, "not repaired: github/broken.go:3:8: no required module provides package")
	assertContains(t, err.Error(), "1 repair(s) could not be applied")
	// No call site could be checked, so none was repaired, and the callers still hold the
	// type the repair changed.
	assertNotContains(t, stderr, "planned call site repairs")
	assertSameAsFixture(t, dir, "github/runner_groups_callers.go")
}

// TestFixReportsAFileItCannotWrite checks that a file that -fix cannot rewrite is reported and
// the rest of the plan is still applied: one read-only file must not leave the call sites in
// the files that come after it unrepaired.
func TestFixReportsAFileItCannotWrite(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	// A hand-written call site in a read-only file that sorts before the other callers.
	const readOnly = `package github

// ReadOnlyCaller is a call site in a file that -fix cannot rewrite.
var ReadOnlyCaller = UpdateRunnerGroupRequest{Name: new("d")}
`
	path := filepath.Join(dir, "github", "callers_readonly.go")
	mustWriteFile(t, path, readOnly)
	assertNilError(t, os.Chmod(path, 0o444))

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the file that cannot be rewritten")
	}
	// The reason is the operating system's word for why the file cannot be written, which
	// differs between them, so the file alone is asserted.
	assertContains(t, stderr, "not repaired: github/callers_readonly.go: ")
	assertEqual(t, 1, strings.Count(stderr, "not repaired: github/callers_readonly.go:"))
	// A repair that was planned and not written fails the run, because the file it was for
	// still holds the type the repair changed.
	assertContains(t, err.Error(), "1 repair(s) could not be applied")
	// The four literals of the plan include the one that cannot be written, and the three
	// that can are repaired all the same.
	assertContains(t, stderr, "planned call site repairs (4):")
	assertContains(t, stderr, `repaired 3 call site(s) in 1 Go file(s)`)
	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")

	got, readErr := os.ReadFile(path)
	assertNilError(t, readErr)
	assertEqual(t, readOnly, string(got))
}

// TestFixLeavesGrandfatheredFindings checks that -fix never rewrites a field that the
// exceptions file names, however the finding is worded, and that it still drops an entry
// that no finding needs any more.
func TestFixLeavesGrandfatheredFindings(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")
	// Grandfather every finding, which is the state of the repository, plus one entry
	// that nothing needs.
	mustWriteFile(t, exceptions, wantExceptionFile+"Gone.Field\n")

	_, stderr, err := runIn(t, dir, "-exceptions", exceptions, "-fix")
	assertNilError(t, err)
	assertContains(t, stderr, "nothing to repair")
	assertContains(t, stderr, "dropped 1 obsolete exception(s) from "+exceptions)

	data, readErr := os.ReadFile(exceptions)
	assertNilError(t, readErr)
	assertEqual(t, wantExceptionFile, string(data))
	// Nothing was repaired, so no field type changed and no call site can be broken.
	assertSameAsFixture(t, dir, "github/enterprise.go", "github/issues.go",
		"github/runner_groups.go", "github/runner_groups_callers.go")
}

// TestFixRepairsNewErrors checks the case a pull request creates: the findings that the
// exceptions file does not name are repaired, and the grandfathered ones keep their lines.
func TestFixRepairsNewErrors(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyFixture(t, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")
	// Grandfather the warnings only, so that the errors are the findings a change adds.
	grandfathered := exceptionHeader + `CreateCommentRequest.Extra
CreateEnterpriseRunnerGroupRequest.SelectedWorkflows
CreateRunnerGroupRequest.SelectedRepositoryIDs
StructTypeRequest.Inner
ValueTypeRequest.Count
`
	mustWriteFile(t, exceptions, grandfathered)

	_, stderr, err := runIn(t, dir, "-exceptions", exceptions, "-fix")
	// The three errors that no repair can express are not grandfathered, so the run still
	// fails: the required property the struct cannot supply, and the two fields of the
	// struct whose operations disagree about them.
	if err == nil {
		t.Fatal("expected an error for the errors that cannot be repaired")
	}
	assertContains(t, err.Error(), "3 schema field issue(s) found")
	assertContains(t, stderr, "planned repairs (3):")
	assertContains(t, stderr, "planned call site repairs (3):")
	assertContains(t, stderr, `repaired 3 call site(s) in 1 Go file(s)`)
	assertContains(t, stderr, "repaired 3 finding(s) in 1 Go file(s)")
	assertContains(t, stderr, "8 finding(s) remain, 0 of them repairable by -fix")

	// The repaired errors had no entries, and the warnings still have findings, so the
	// exceptions file keeps every line it had.
	data, readErr := os.ReadFile(exceptions)
	assertNilError(t, readErr)
	assertEqual(t, grandfathered, string(data))
	assertSameAsFixture(t, dir, "github/enterprise.go", "github/issues.go")

	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups.go")
	checkFileGolden(t, fixedFixtureGolden, dir, "github/runner_groups_callers.go")
}

// fixedFixtureGolden names the test whose golden files hold the fixture as -fix leaves it.
// Every test that repairs the fixture repairs it the same way, so they compare against these
// files rather than each keeping a copy of the same output.
const fixedFixtureGolden = "TestFix"

// checkFileGolden compares a file of the repaired checkout with the golden file that testName
// holds under it.
func checkFileGolden(t *testing.T, testName, dir, name string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
	assertNilError(t, err)
	checkGolden(t, testName, name, string(got))
}

// checkGolden compares got with the golden file of the named test and file. A missing golden
// file is written, so that a new case regenerates itself; run UPDATE_GOLDEN=1 script/test.sh
// to refresh every golden file at once.
func checkGolden(t *testing.T, testName, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", testName, filepath.FromSlash(name))
	if _, err := os.Stat(path); os.IsNotExist(err) || os.Getenv("UPDATE_GOLDEN") != "" {
		mustWriteFile(t, path, got)
		t.Logf("wrote golden file %v", path)
		return
	}
	want, err := os.ReadFile(path)
	assertNilError(t, err)
	// The tool writes LF, but a Windows checkout converts the golden file to CRLF,
	// so compare with the endings normalized. tools/metadata does the same.
	assertEqual(t, normalizeEOL(string(want)), normalizeEOL(got))
}

// normalizeEOL removes the carriage returns of CRLF line endings, so that a file compares
// equal to one that uses LF.
func normalizeEOL(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// assertSameAsFixture fails unless the checkout at dir holds the same content as the fixture
// for each named file, which is how a test asserts that a file was not rewritten.
func assertSameAsFixture(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, name := range names {
		want, err := os.ReadFile(filepath.Join(fixtureRepo, name))
		assertNilError(t, err)
		got, err := os.ReadFile(filepath.Join(dir, name))
		assertNilError(t, err)
		assertEqual(t, string(want), string(got))
	}
}

// copyTree copies every file of src into dst, preserving relative paths.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	assertNilError(t, filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	}))
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	assertNilError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	assertNilError(t, os.WriteFile(path, []byte(content), 0o644))
}

// assertEqual fails the test unless got equals want. It takes any, so that the callers can
// compare slices and maps as well as scalars, and accepts the cmp options that comparing a
// type with unexported fields needs.
func assertEqual(t *testing.T, want, got any, opts ...cmp.Option) {
	t.Helper()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("missing %q in:\n%v", needle, haystack)
	}
}

func assertNotContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Errorf("%q unexpectedly found in:\n%v", needle, haystack)
	}
}

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
