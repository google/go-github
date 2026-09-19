// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command schemafields checks Go request-body struct field optionality against GitHub's
// official OpenAPI request body schemas.
//
// The struct-to-schema mapping is derived automatically, with no new annotations:
//
//	//meta:operation POST /repos/{owner}/{repo}/issues/{issue_number}/comments
//	func (s *IssuesService) CreateComment(ctx context.Context, owner, repo string, number int, body IssueCommentRequest) (*IssueComment, *Response, error) {
//
// The //meta:operation annotation already names the operation, and paramcheck already
// requires the body parameter to be named "body" and passed by value. So every by-value
// body parameter is checked against the request body schema of the operation its method is
// annotated with. Coverage therefore grows automatically as types are converted, and a
// contributor adding an endpoint is checked without having to do anything extra.
//
// CONTRIBUTING.md requires that required fields be non-pointer types without an omit
// option, and that optional fields be pointer types with "omitempty" (or "omitzero" for
// structs, slices and maps). This tool reports every request body field that disagrees with
// the schema, and -fix repairs the ones that can be repaired mechanically.
//
// A repaired field whose type changes from a value to a pointer, or the other way around,
// can require a caller or a test to be updated too, and the generated accessors must be
// regenerated with script/generate.sh. Run that before checking the result.
package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// options are the command line flags.
type options struct {
	repo            string
	descriptions    string
	cacheDir        string
	exceptions      string
	format          string
	minSeverity     string
	fix             bool
	writeExceptions bool
	verbose         bool
}

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "schemafields: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags parses the command line, writing usage to stderr.
func parseFlags(stderr io.Writer, args []string) (*options, error) {
	o := &options{}
	fs := flag.NewFlagSet("schemafields", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&o.repo, "repo", ".", "path to the go-github checkout")
	fs.StringVar(&o.descriptions, "descriptions", "",
		"comma-separated OpenAPI description files, tried in order; the default is to download the ones pinned in openapi_operations.yaml")
	fs.StringVar(&o.cacheDir, "cache-dir", "", "directory in which to cache downloaded descriptions")
	fs.StringVar(&o.exceptions, "exceptions", "",
		"path to the exceptions file of findings to ignore; the default is tools/schemafields/exceptions.txt in the checkout")
	fs.StringVar(&o.format, "format", "text", "output format: text or github")
	fs.StringVar(&o.minSeverity, "min-severity", "warn", "report findings of at least this severity: warn or error")
	fs.BoolVar(&o.fix, "fix", false,
		"rewrite the Go sources to repair the findings that can be repaired; run script/generate.sh afterward to regenerate the accessors")
	fs.BoolVar(&o.writeExceptions, "write-exceptions", false, "rewrite the exceptions file from the current findings")
	fs.BoolVar(&o.verbose, "verbose", false, "list request bodies that cannot be mapped to an operation")
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: schemafields [flags]\n\nChecks request body struct fields against GitHub's OpenAPI schemas.\n\nFlags:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if o.format != "text" && o.format != "github" {
		return nil, fmt.Errorf("unknown -format %q", o.format)
	}
	if o.minSeverity != "warn" && o.minSeverity != "error" {
		return nil, fmt.Errorf("unknown -min-severity %q", o.minSeverity)
	}
	if fs.NArg() > 0 {
		return nil, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	return o, nil
}

// exceptionsPath returns the path of the exceptions file.
func (o *options) exceptionsPath() string {
	if o.exceptions != "" {
		return o.exceptions
	}
	return filepath.Join(o.repo, "tools", "schemafields", "exceptions.txt")
}

// run checks the checkout named by the flags. It writes findings to stdout and a summary to
// stderr, and returns an error when it finds a problem that is not grandfathered.
func run(stdout, stderr io.Writer, args []string) error {
	o, err := parseFlags(stderr, args)
	if err != nil {
		// -h is a request for the usage text, which parseFlags has written, and
		// not a failure.
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	info, err := scanRepo(o.repo)
	if err != nil {
		return err
	}

	files := splitList(o.descriptions)
	cacheDir := o.cacheDir
	if cacheDir == "" {
		cacheDir = defaultCacheDir()
	}
	ds, err := loadDescriptions(context.Background(), o.repo, files, cacheDir)
	if err != nil {
		return err
	}

	c := &checker{repo: o.repo, d: ds, info: info}
	c.run()

	exc, err := loadExceptions(o.exceptionsPath())
	if err != nil {
		return err
	}

	if o.fix {
		if err := fixAll(c, exc, o, stderr); err != nil {
			return err
		}
	}

	if o.writeExceptions {
		keys := make([]string, 0, len(c.diags))
		for _, d := range c.diags {
			keys = append(keys, d.key())
		}
		slices.Sort(keys)
		exc.setEntries(keys) // Dropping the duplicates that several findings share.
		if err := exc.write(exc.entries); err != nil {
			return err
		}
		// The file now holds every finding, so the run that regenerated it is clean
		// and can be used to set the baseline for a repository that has none yet.
		fmt.Fprintf(stderr, "wrote %v exception(s) to %v\n", len(exc.entries), exc.path)
	}

	threshold := sevWarn
	if o.minSeverity == "error" {
		threshold = sevError
	}
	shown, errored := 0, 0
	for _, d := range c.diags {
		if d.sev < threshold || exc.suppress(d.key()) {
			continue
		}
		shown++
		if d.sev == sevError {
			errored++
		}
		printFinding(stdout, o.format, d)
	}

	printSummary(stderr, c, o, shown, errored)
	if o.verbose {
		for _, m := range c.unannotated {
			fmt.Fprintf(stderr, "no //meta:operation annotation: %v:%v: %v (%v)\n", m.file, m.line, m.funcName, m.bodyType)
		}
	}

	obsolete := exc.obsolete(c.diags)
	if len(obsolete) > 0 {
		fmt.Fprint(stderr, "\nobsolete exceptions, which no finding needs any more (run with -fix to drop them):\n")
		for _, entry := range obsolete {
			fmt.Fprintf(stderr, "  %v\n", entry)
		}
	}

	switch {
	case errored > 0:
		return fmt.Errorf("%v schema field issue(s) found", errored)
	case len(obsolete) > 0:
		return fmt.Errorf("%v obsolete exception(s) to remove", len(obsolete))
	}
	return nil
}

// fixAll repairs what it can, re-checks the result, and drops the exceptions that the
// repairs made obsolete.
func fixAll(c *checker, exc *exceptions, o *options, stderr io.Writer) error {
	fixed, written, notes, err := applyFixes(o.repo, c.diags)
	if err != nil {
		return err
	}
	for _, note := range notes {
		fmt.Fprintf(stderr, "not repaired: %v\n", note)
	}
	// Re-read the sources, so that the summary describes the repaired tree.
	info, err := scanRepo(o.repo)
	if err != nil {
		return err
	}
	c.info = info
	c.run()
	fmt.Fprintf(stderr, "repaired %v finding(s) in %v Go file(s); %v finding(s) remain, %v of them repairable\n",
		fixed, written, len(c.diags), len(c.diags)-unrepaired(c.diags))

	obsolete := exc.obsolete(c.diags)
	if len(obsolete) > 0 {
		kept := slices.DeleteFunc(slices.Clone(exc.entries), func(entry string) bool {
			return slices.Contains(obsolete, entry)
		})
		if err := exc.write(kept); err != nil {
			return err
		}
		exc.setEntries(kept)
		fmt.Fprintf(stderr, "dropped %v obsolete exception(s) from %v\n", len(obsolete), exc.path)
	}
	return nil
}

// splitList splits a comma-separated flag value.
func splitList(value string) []string {
	var out []string
	for item := range strings.SplitSeq(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}

// printFinding writes one finding, as plain text or as a GitHub Actions annotation.
func printFinding(w io.Writer, format string, d *diagnostic) {
	if format == "github" {
		fmt.Fprintf(w, "::%v file=%v,line=%v,title=%v::%v\n",
			d.sev.level(), escapeAnnotation(d.file), d.line,
			escapeAnnotation(d.key()), escapeAnnotation(d.message))
		return
	}
	fmt.Fprintf(w, "%-5v %v:%v: %v: %v\n", strings.ToUpper(d.sev.String()), d.file, d.line, d.key(), d.message)
}

// escapeAnnotation escapes a GitHub Actions annotation argument.
func escapeAnnotation(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A", ":", "%3A", ",", "%2C").Replace(s)
}

// printSummary writes the run statistics.
func printSummary(w io.Writer, c *checker, o *options, shown, errored int) {
	s := c.stats
	repairable := 0
	for _, d := range c.diags {
		if d.action != nil {
			repairable++
		}
	}
	fmt.Fprintf(w, "\nscanned %v files, %v methods (%v with //meta:operation)\n", s.files, s.methods, s.methodsWithOps)
	fmt.Fprintf(w, "body params: %v by value, %v by pointer (skipped; run paramcheck to convert them)\n", s.bodyValue, s.bodyPointer)
	fmt.Fprintf(w, "checked %v body structs: %v resolved operation uses, %v uses with no JSON request body\n",
		s.structsChecked, s.usesResolved, s.usesNoSchema)
	fmt.Fprintf(w, "fields checked: %v (%v conditionally required, left alone)\n", s.fieldsChecked, s.fieldsConditional)
	fmt.Fprintf(w, "%v findings (%v shown, %v errors, %v repairable with -fix)\n", len(c.diags), shown, errored, repairable)

	rules := make([]string, 0, len(c.byRule))
	for rule := range c.byRule {
		rules = append(rules, rule)
	}
	slices.SortFunc(rules, func(a, b string) int {
		return cmp.Or(cmp.Compare(c.byRule[b], c.byRule[a]), strings.Compare(a, b))
	})
	for _, rule := range rules {
		fmt.Fprintf(w, "  %4d  %v: %v\n", c.byRule[rule], rule, ruleDescriptions[rule])
	}
	if o.minSeverity == "error" {
		fmt.Fprint(w, "warnings are hidden; remove -min-severity=error to see them\n")
	}
}
