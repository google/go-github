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
`

// runOnFixture copies the fixture checkout into a temporary directory and runs the tool
// over the copy, so that -fix cannot touch the testdata.
func runOnFixture(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	dir := t.TempDir()
	copyTree(t, fixtureRepo, dir)
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
	assertContains(t, err.Error(), "4 schema field issue(s) found")
	assertEqual(t, wantFindings, stdout)
	assertContains(t, stderr, "checked 12 body structs: 13 resolved operation uses, 1 uses with no JSON request body")
	assertContains(t, stderr, "9 findings (9 shown, 4 errors, 7 repairable with -fix)")
	assertContains(t, stderr, "1 conditionally required, left alone")
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
	assertContains(t, stderr, "9 findings (4 shown, 4 errors")
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
	copyTree(t, fixtureRepo, dir)
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
# reports problems that a change introduces.
#
# Run "script/check-schema-fields.sh -fix" to repair what can be repaired
# automatically, and delete the lines that become obsolete.
`

// wantExceptionFile is the exceptions file that -write-exceptions writes for the fixture:
// the default header, then every finding's key in sorted order.
const wantExceptionFile = exceptionHeader + `CreateCommentRequest.Extra
CreateEnterpriseRunnerGroupRequest.SelectedWorkflows
CreateRunnerGroupRequest.Name
CreateRunnerGroupRequest.SelectedRepositoryIDs
MissingRequest.required_thing
OneOfRequest.KindTag
StructTypeRequest.Inner
UpdateRunnerGroupRequest.Name
ValueTypeRequest.Count
`

func TestWriteExceptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyTree(t, fixtureRepo, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")

	// Every finding becomes an exception, so the run that regenerated the file is
	// clean and can be used to set a baseline.
	_, stderr, err := runIn(t, dir, "-exceptions", exceptions, "-write-exceptions")
	assertNilError(t, err)
	assertContains(t, stderr, "wrote 9 exception(s)")

	data, readErr := os.ReadFile(exceptions)
	assertNilError(t, readErr)
	assertEqual(t, wantExceptionFile, string(data))
}

// TestFix repairs the fixture and compares each rewritten file with its golden copy.
func TestFix(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyTree(t, fixtureRepo, dir)

	_, stderr, err := runIn(t, dir, "-fix")
	if err == nil {
		t.Fatal("expected an error for the finding that cannot be repaired")
	}
	assertContains(t, stderr, "repaired 7 finding(s) in 2 Go file(s)")
	assertContains(t, stderr, "2 finding(s) remain, 0 of them repairable")

	for _, name := range []string{"github/enterprise.go", "github/issues.go", "github/runner_groups.go"} {
		got, readErr := os.ReadFile(filepath.Join(dir, name))
		assertNilError(t, readErr)
		checkGolden(t, name, string(got))
	}
}

// TestFixDropsObsoleteExceptions checks that a repair drops the exception that grandfathered
// the finding it repaired, and that the findings it cannot repair keep theirs.
func TestFixDropsObsoleteExceptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	copyTree(t, fixtureRepo, dir)
	exceptions := filepath.Join(t.TempDir(), "exceptions.txt")
	// Grandfather every finding, so that the run starts out clean.
	mustWriteFile(t, exceptions, wantExceptionFile)

	_, stderr, err := runIn(t, dir, "-exceptions", exceptions, "-fix")
	// The two findings that -fix cannot repair keep their exceptions, so the run
	// that repaired the other seven ends up clean.
	assertNilError(t, err)
	assertContains(t, stderr, "repaired 7 finding(s) in 2 Go file(s)")
	assertContains(t, stderr, "dropped 7 obsolete exception(s) from "+exceptions)

	data, readErr := os.ReadFile(exceptions)
	assertNilError(t, readErr)
	assertEqual(t, exceptionHeader+"CreateCommentRequest.Extra\nMissingRequest.required_thing\n", string(data))
}

// checkGolden compares got with the golden file named by the current test. A missing golden
// file is written, so that a new case regenerates itself; run UPDATE_GOLDEN=1 script/test.sh
// to refresh every golden file at once.
func checkGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", t.Name(), filepath.FromSlash(name))
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

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
