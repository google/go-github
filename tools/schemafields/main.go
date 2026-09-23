// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command schemafields checks Go request-body struct field optionality against GitHub's
// official OpenAPI request body schemas.
//
// The schemas are the source of truth, and they are not in this repository. GitHub publishes
// them in github/rest-api-description, and the root openapi_operations.yaml pins the revision
// to check against in its openapi_commit field. A run downloads the three descriptions that
// revision names (api.github.com, ghec, and the newest ghes-3.x) and caches them in the user
// cache directory, so a repeated run is offline; -descriptions checks local files instead.
// Maintainers advance the pin with script/metadata.sh update-openapi, which regenerates the
// openapi_operations section from the descriptions at the revision it sees, and the linter
// workflow's update-openapi --validate checks that the pinned revision still matches. A
// finding that is then new is repaired with -fix, or recorded in
// tools/schemafields/exceptions.txt to leave it alone. CONTRIBUTING.md has the field rules
// and the full account of what -fix does.
//
// The struct-to-schema mapping needs no new annotations. //meta:operation already names the
// operation a method calls, and paramcheck already requires the body parameter to be named
// "body" and passed by value, so every by-value body parameter is checked against the request
// body schema of its operation. Coverage grows as pointer bodies are converted, and a
// contributor adding an endpoint is checked without having to do anything extra.
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
		"path to the exceptions file of findings to ignore; the default is tools/schemafields/exceptions.txt in the checkout, and an empty file, such as /dev/null, treats every finding as new")
	fs.StringVar(&o.format, "format", "text", "output format: text or github")
	fs.StringVar(&o.minSeverity, "min-severity", "warn", "report findings of at least this severity: warn or error")
	fs.BoolVar(&o.fix, "fix", false,
		"rewrite the Go sources to repair the error-severity findings that the exceptions file does not grandfather, and the call sites that a changed field type leaves unable to compile; a changed field type is generated into other files as well, so -fix runs the checkout's script/generate.sh, and fails the run rather than report success when it cannot write a repair it planned")
	fs.BoolVar(&o.writeExceptions, "write-exceptions", false, "rewrite the exceptions file from the current findings")
	fs.BoolVar(&o.verbose, "verbose", false,
		"add the full breakdown to the summary: the scan counts, the fields a rule leaves alone, the per-rule tally, and every request body that was left unchecked, with the reason")
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: schemafields [flags]\n\nChecks request body struct fields against GitHub's OpenAPI schemas.\n\nFlags:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	// The checkout path is compared with the paths of the files under it by prefix, and
	// filepath.Join drops the leading "./" of a checkout named ".", so the two spellings would
	// never match and -fix would find no module to compile for the call sites of a changed field
	// type. Making the path absolute once, here, keeps them equal.
	abs, err := filepath.Abs(o.repo)
	if err != nil {
		return nil, err
	}
	o.repo = abs
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

	unrepaired := 0
	if o.fix {
		unrepaired, err = fixAll(c, exc, o, stderr)
		if err != nil {
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
	errored := 0
	for _, d := range c.diags {
		if d.sev < threshold || exc.suppress(d.key()) {
			continue
		}
		if d.sev == sevError {
			errored++
		}
		printFinding(stdout, o.format, d)
	}

	printSummary(stderr, c, o, exc)
	if o.verbose {
		for _, m := range c.unannotated {
			fmt.Fprintf(stderr, "no //meta:operation annotation: %v:%v: %v (%v)\n", m.file, m.line, m.funcName, m.bodyType)
		}
		for _, m := range c.unstructured {
			fmt.Fprintf(stderr, "the body is not a struct this repository declares: %v:%v: %v (%v)\n", m.file, m.line, m.funcName, m.bodyType)
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
	case unrepaired > 0:
		// -fix was asked to repair these and could not, so the run must not report success:
		// whatever it left behind, the checkout no longer builds as it stands.
		return fmt.Errorf("%v repair(s) could not be applied, so the checkout may not build", unrepaired)
	case errored > 0:
		return fmt.Errorf("%v schema field issue(s) found", errored)
	case len(obsolete) > 0:
		return fmt.Errorf("%v obsolete exception(s) to remove", len(obsolete))
	}
	return nil
}

// fixCandidates returns the findings that -fix rewrites: the error-severity findings that the
// exceptions file does not grandfather.
func (c *checker) fixCandidates(exc *exceptions) []*diagnostic {
	var out []*diagnostic
	for _, d := range c.repairableFindings() {
		if !exc.suppress(d.key()) {
			out = append(out, d)
		}
	}
	return out
}

// repairableFindings returns every finding that -fix is built to rewrite, whether or not the
// exceptions file grandfathers it. It is what tells a maintainer what the baseline is holding
// back, so it does not consult the exceptions file.
func (c *checker) repairableFindings() []*diagnostic {
	var out []*diagnostic
	for _, d := range c.diags {
		switch {
		case d.sev != sevError:
			// A warning is advisory, so it is reported for a human to decide about.
		case d.action == nil || d.info == nil:
			// Nothing can be done mechanically; the summary counts these.
		case !c.agrees(d):
			// A struct can be the body of more than one operation, and their
			// schemas can disagree about the field. A repair that suits the
			// operation the finding came from would then contradict another, so
			// this one is reported for a human to decide about instead.
		default:
			out = append(out, d)
		}
	}
	return out
}

// fixAll repairs the findings that -fix is allowed to repair, repairs the call sites that the
// repairs invalidate, re-checks the result, and drops the exceptions that the repairs made
// obsolete. It returns the number of repairs that are still needed and that it could not write,
// which is what makes the run fail: a repair left unwritten leaves the checkout unable to build.
func fixAll(c *checker, exc *exceptions, o *options, stderr io.Writer) (unrepaired int, err error) {
	planned := c.fixCandidates(exc)
	// The type changes have to be read from the sources before the repairs are written.
	changes := typeChanges(o.repo, planned)
	if len(planned) == 0 {
		fmt.Fprint(stderr, "nothing to repair: -fix rewrites only the error-severity findings that the exceptions file does not grandfather\n")
	} else {
		fmt.Fprintf(stderr, "planned repairs (%v):\n", len(planned))
		for _, d := range planned {
			fmt.Fprintf(stderr, "  %v:%v: %v: %v\n", d.file, d.line, d.key(), d.action)
		}
	}

	fixed, written, notes := applyFixes(o.repo, planned)
	// A field that changed between a value type and a pointer leaves the call sites that
	// still use the old type unable to compile, so repair them from what the compiler says
	// before re-checking the tree. The count is reported on its own, because a call site is
	// not a finding.
	_, callNotes, err := repairCallSites(changes, o, stderr)
	if err != nil {
		return len(notes), err
	}
	for _, note := range append(notes, callNotes...) {
		fmt.Fprintf(stderr, "not repaired: %v\n", note)
	}
	unrepaired = len(notes) + len(callNotes)

	// Re-read the sources, so that the summary describes the repaired tree.
	info, err := scanRepo(o.repo)
	if err != nil {
		return unrepaired, err
	}
	c.info = info
	c.run()
	if len(planned) > 0 {
		fmt.Fprintf(stderr, "repaired %v finding(s) in %v Go file(s); %v finding(s) remain, %v of them repairable by -fix\n",
			fixed, written, len(c.diags), len(c.fixCandidates(exc)))
	}

	obsolete := exc.obsolete(c.diags)
	if len(obsolete) > 0 {
		kept := slices.DeleteFunc(slices.Clone(exc.entries), func(entry string) bool {
			return slices.Contains(obsolete, entry)
		})
		if err := exc.write(kept); err != nil {
			return unrepaired, err
		}
		exc.setEntries(kept)
		fmt.Fprintf(stderr, "dropped %v obsolete exception(s) from %v\n", len(obsolete), exc.path)
	}
	return unrepaired, nil
}

// repairCallSites repairs the struct literals that the type changes of the planned repairs
// invalidate, and reports how many it repaired. It returns a note for every repair that is still
// needed and that it could not write.
func repairCallSites(changes []*typeChange, o *options, stderr io.Writer) (fixed int, notes []string, err error) {
	if len(changes) == 0 {
		return 0, nil, nil
	}
	changed := make([]string, 0, len(changes))
	for _, change := range changes {
		if !slices.Contains(changed, change.file) {
			changed = append(changed, change.file)
		}
	}
	dirs := modulesToCheck(o.repo, changed)
	if len(dirs) == 0 {
		// Nothing to compile, so no call site can be found. Say so rather than report the
		// finding as repaired: the repairs are written either way, and the tree they were
		// written to may no longer build.
		return 0, []string{fmt.Sprintf("no module found for %v, so the call sites that the changed field types break were not repaired",
			strings.Join(changed, ", "))}, nil
	}
	fixed, written, notes, err := fixCallSites(context.Background(), o.repo, changes, dirs, stderr)
	if err != nil {
		return fixed, notes, err
	}
	if fixed > 0 {
		fmt.Fprintf(stderr, "repaired %v call site(s) in %v Go file(s)\n", fixed, written)
	}
	return fixed, notes, nil
}

// plural picks the word that agrees with a count of n, so that a summary line reads as a
// sentence whatever the numbers are.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// repoRelative returns a path as the summary shows it: relative to the checkout when it is
// inside it, so that it is named the way the findings name their files.
func repoRelative(repo, path string) string {
	rel, err := filepath.Rel(repo, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return path
	}
	return rel
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

// printSummary writes the run statistics: how much of the API was checked, what could not be
// checked and why, and what the findings are, counting the ones that the exceptions file
// grandfathers alongside the new ones. It reports the state of the tree, which is what a
// maintainer acts on, and -verbose adds the full breakdown.
func printSummary(w io.Writer, c *checker, o *options, exc *exceptions) {
	s := c.stats
	fmt.Fprintln(w)
	if o.verbose {
		fmt.Fprintf(w, "scanned %v files, %v methods (%v with //meta:operation)\n",
			s.files, s.methods, s.methodsWithOps)
	}
	// Every method whose request body is a struct, and how many of them this tool can judge. A
	// body that is not a struct has no fields for a schema check to apply to, so it is left out
	// of the population rather than reported: -verbose names those bodies. What this tool cannot
	// judge is the coverage it will not have until a body is converted or the pin moves, so each
	// kind of gap is named rather than folded into a total.
	bodyMethods := s.bodyValue + s.bodyPointer
	fmt.Fprintf(w, "checked %v of %v %v that %v a struct request body (%v body structs, %v fields)\n",
		s.bodyValue, bodyMethods, plural(bodyMethods, "method", "methods"), plural(bodyMethods, "takes", "take"),
		s.structsChecked, s.fieldsChecked)
	if s.bodyPointer > 0 {
		fmt.Fprintf(w, "  not checked: %v %v whose body is passed by pointer (run paramcheck to convert %v)\n",
			s.bodyPointer, plural(s.bodyPointer, "method", "methods"), plural(s.bodyPointer, "it", "them"))
	}
	if s.usesNoJSONBody > 0 {
		fmt.Fprintf(w, "  not checked: %v operation %v with no JSON request body\n",
			s.usesNoJSONBody, plural(s.usesNoJSONBody, "use", "uses"))
	}
	if s.usesNotInPlan > 0 {
		fmt.Fprintf(w, "  not checked: %v operation %v that the pinned revision does not document\n",
			s.usesNotInPlan, plural(s.usesNotInPlan, "use", "uses"))
		fmt.Fprintf(w, "    a newer openapi_commit in %v would check %v\n",
			operationsFile, plural(s.usesNotInPlan, "it", "them"))
	}
	if s.usesUnreadable > 0 {
		fmt.Fprintf(w, "  not checked: %v operation %v whose description could not be read\n",
			s.usesUnreadable, plural(s.usesUnreadable, "use", "uses"))
	}
	if o.verbose {
		fmt.Fprintf(w, "  resolved operation uses: %v\n", s.usesResolved)
		fmt.Fprintf(w, "  conditionally required, and left alone: %v of the %v fields\n",
			s.fieldsConditional, s.fieldsChecked)
	}

	// A finding the exceptions file does not name is the only kind a change can be told about,
	// but the findings it does name are still the state of the tree, so both are reported: a
	// summary that counts only the new ones is how a baseline hides an error behind "0 errors".
	warn, errs, fresh := 0, 0, 0
	for _, d := range c.diags {
		if d.sev == sevError {
			errs++
		} else {
			warn++
		}
		if !exc.suppress(d.key()) {
			fresh++
		}
	}
	switch {
	case len(c.diags) == 0:
		fmt.Fprint(w, "no findings\n")
	case fresh == 0:
		fmt.Fprintf(w, "%v %v (%v warn, %v error), all grandfathered by %v\n",
			len(c.diags), plural(len(c.diags), "finding", "findings"), warn, errs, repoRelative(c.repo, exc.path))
	case fresh == len(c.diags):
		fmt.Fprintf(w, "%v %v (%v warn, %v error)\n",
			len(c.diags), plural(len(c.diags), "finding", "findings"), warn, errs)
	default:
		fmt.Fprintf(w, "%v %v (%v warn, %v error), %v new, %v grandfathered by %v\n",
			len(c.diags), plural(len(c.diags), "finding", "findings"), warn, errs, fresh, len(c.diags)-fresh,
			repoRelative(c.repo, exc.path))
	}
	// What -fix is built to rewrite, which the exceptions file can hold back: an entry that
	// grandfathered a repairable finding is a decision that costs a repair. An error that -fix
	// will not rewrite is worth the same line, because it is the part of the baseline that only
	// a person can clear.
	fixable, freshFixable := len(c.repairableFindings()), len(c.fixCandidates(exc))
	if fixable > 0 {
		fmt.Fprintf(w, "  repairable by -fix: %v new", freshFixable)
		if held := fixable - freshFixable; held > 0 {
			fmt.Fprintf(w, ", %v grandfathered (drop the entries for -fix to repair)", held)
		}
		fmt.Fprint(w, "\n")
	}
	if manual := errs - fixable; manual > 0 {
		fmt.Fprintf(w, "  %v of the errors %v not repairable by -fix\n",
			manual, plural(manual, "is", "are"))
	}

	if o.verbose {
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
	}
	if o.minSeverity == "error" {
		fmt.Fprint(w, "warnings are hidden; remove -min-severity=error to see them\n")
	}
}
