// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// maxCompilePasses bounds the repair loop. A repaired literal can expose another one, but a
// pass that repairs nothing ends the loop, so this only stops the pathological case.
const maxCompilePasses = 4

// checkTimeout bounds one compiler run, so that a checkout that hangs does not hang -fix.
const checkTimeout = 10 * time.Minute

// structLitErrRe matches the error that the compiler reports at a struct literal that still
// sets a field with the type the field had before -fix changed it, for example:
//
//	github/actions_oidc_test.go:125:21: cannot use new(false) (value of type *bool) as bool value in struct literal
var structLitErrRe = regexp.MustCompile(`^(.+?\.go):(\d+):(\d+): cannot use .*? as (.+?) value in struct literal$`)

// typeChange is one field whose Go type a repair changed. The call sites that the change
// invalidates are the struct literals that still use the old type, and the compiler names
// them exactly, which is why -fix asks it rather than guessing from the syntax.
type typeChange struct {
	file       string // the Go file that holds the declaration, relative to the repo
	structName string
	fieldName  string
	newType    string // the Go type the field has after the change, as the compiler names it
	toPointer  bool   // true when the change made the field a pointer
}

// typeChanges returns the type changes that the given repairs make. It reads the field types
// from the sources as they are now, so it must run before the repairs are written.
func typeChanges(repo string, diags []*diagnostic) []*typeChange {
	sources := map[string][]byte{}
	var out []*typeChange
	seen := map[string]bool{}
	for _, d := range diags {
		f := d.info
		if f == nil || d.action == nil || (!d.action.unwrap && !d.action.makePointer) {
			continue
		}
		src, ok := sources[f.file]
		if !ok {
			data, err := os.ReadFile(filepath.Join(repo, f.file))
			if err != nil {
				continue
			}
			src = data
			sources[f.file] = src
		}
		if f.typOff > f.typEnd || f.typEnd > len(src) {
			continue
		}
		change := &typeChange{
			file:       f.file,
			structName: d.owner,
			fieldName:  f.goName,
			toPointer:  d.action.makePointer,
		}
		switch {
		case d.action.unwrap:
			if f.starOff < 0 || f.starOff+1 > f.typEnd {
				continue
			}
			change.newType = string(src[f.starOff+1 : f.typEnd])
		case d.action.makePointer:
			change.newType = "*" + string(src[f.typOff:f.typEnd])
		}
		if change.newType == "" {
			continue
		}
		key := fmt.Sprintf("%v.%v=%v", change.structName, change.fieldName, change.newType)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, change)
	}
	return out
}

// compileErr is one struct-literal error from the compiler.
type compileErr struct {
	dir      string // the module directory the compiler ran in
	file     string // the file as the compiler named it, relative to dir
	line     int
	offset   int // byte offset of the reported expression within the file
	destType string
}

// path returns the path of the failing file relative to the repo.
func (e *compileErr) path(repo string) string {
	rel, err := filepath.Rel(repo, filepath.Join(e.dir, filepath.FromSlash(e.file)))
	if err != nil {
		return e.file
	}
	return filepath.ToSlash(rel)
}

// compileStructLitErrors builds each module, including its test files, and returns the
// struct-literal errors the compiler reports, in file order. It includes test files, because
// that is where most call sites are.
func compileStructLitErrors(ctx context.Context, repo string, dirs []string) ([]*compileErr, error) {
	var out []*compileErr
	for _, dir := range dirs {
		output, err := checkModule(ctx, dir)
		if err != nil && output == "" {
			return nil, fmt.Errorf("go test %v: %w", dir, err)
		}
		for line := range strings.Lines(output) {
			if e := parseStructLitErr(dir, line); e != nil {
				out = append(out, e)
			}
		}
	}
	// The compiler reports a module at a time, so sort by file to keep the plan stable
	// however the modules were ordered.
	slices.SortStableFunc(out, func(a, b *compileErr) int {
		return strings.Compare(a.path(repo), b.path(repo))
	})
	return out, nil
}

// checkModule compiles one module, including its test files, and returns the combined output
// even on failure, because a failed build is how the compiler reports the errors that -fix
// repairs.
//
// "go test -run=NONE" builds a test binary without running a test in it, which is the only way
// to compile test files, and it is the compiler rather than go vet that does the work, because
// go vet stops at the first type error in a package: -gcflags=-e asks the compiler for every
// error instead, so that one build reports every call site the change broke.
func checkModule(ctx context.Context, dir string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-gcflags=-e", "-vet=off", "-run=NONE", "./...")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return string(out), ctx.Err()
	}
	return string(out), err
}

// parseStructLitErr parses one line of compiler output. It returns nil for anything else, so
// that the package headers, the FAIL lines and any other diagnostic are ignored.
func parseStructLitErr(dir, line string) *compileErr {
	m := structLitErrRe.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return nil
	}
	lineNo, err := strconv.Atoi(m[2])
	if err != nil {
		return nil
	}
	col, err := strconv.Atoi(m[3])
	if err != nil {
		return nil
	}
	src, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(m[1])))
	if err != nil {
		return nil
	}
	offset := offsetOf(src, lineNo, col)
	if offset < 0 {
		return nil
	}
	return &compileErr{
		dir:      dir,
		file:     m[1],
		line:     lineNo,
		offset:   offset,
		destType: m[4],
	}
}

// offsetOf returns the byte offset of the 1-based line and column, or -1 when they are out of
// range.
func offsetOf(src []byte, line, col int) int {
	if line < 1 || col < 1 {
		return -1
	}
	offset := 0
	for i := 1; i < line; i++ {
		next := bytes.IndexByte(src[offset:], '\n')
		if next < 0 {
			return -1
		}
		offset += next + 1
	}
	offset += col - 1
	if offset >= len(src) {
		return -1
	}
	return offset
}

// callSiteFix is one literal repair, ready to apply and to print in the plan.
type callSiteFix struct {
	path   string // relative to the repo
	line   int
	field  string
	before string
	after  string
	edit   *edit
}

// fixCallSites rewrites the call sites that the type changes invalidate. It returns the number
// of literals repaired and the number of files rewritten.
func fixCallSites(ctx context.Context, repo string, changes []*typeChange, dirs []string, stderr io.Writer) (fixed, written int, err error) {
	if len(changes) == 0 {
		return 0, 0, nil
	}
	for range maxCompilePasses {
		errs, err := compileStructLitErrors(ctx, repo, dirs)
		if err != nil {
			return fixed, written, err
		}
		planned := planCallSites(repo, errs, changes)
		if len(planned) == 0 {
			reportUnrepaired(stderr, repo, errs, changes)
			return fixed, written, nil
		}
		fmt.Fprintf(stderr, "planned call site repairs (%v):\n", len(planned))
		for _, p := range planned {
			fmt.Fprintf(stderr, "  %v:%v: %v: replace %v with %v\n", p.path, p.line, p.field, p.before, p.after)
		}
		n, w, err := applyCallSites(repo, planned)
		if err != nil {
			return fixed, written, err
		}
		if n == 0 {
			return fixed, written, nil
		}
		fixed += n
		written += w
	}
	// Out of passes, which the loop above needs only when a repair exposes another one. Look
	// once more, so that what is still broken is reported rather than left for the caller to
	// find in a tree that does not compile.
	if errs, err := compileStructLitErrors(ctx, repo, dirs); err == nil {
		reportUnrepaired(stderr, repo, errs, changes)
	}
	return fixed, written, nil
}

// reportUnrepaired writes a note for each error that a type change explains and that the
// literal's shape leaves no mechanical repair for.
func reportUnrepaired(w io.Writer, repo string, errs []*compileErr, changes []*typeChange) {
	parsed := map[string]*parsedFile{}
	for _, e := range errs {
		pf := parseFor(repo, parsed, e)
		if pf == nil {
			continue
		}
		lit, err := pf.locate(e.offset)
		if err != nil {
			continue
		}
		change := matchChange(lit, e.destType, changes)
		if change == nil {
			continue
		}
		_, note := repairedExpr(pf, lit.expr, change.toPointer)
		if note == nil {
			continue // The next pass repairs it; the loop ended for another reason.
		}
		fmt.Fprintf(w, "not repaired: %v:%v: %v.%v is now %v, but %v\n",
			e.path(repo), e.line, change.structName, change.fieldName, change.newType, note)
	}
}

// planCallSites matches the compiler errors against the type changes and returns the repairs
// that -fix can make. An error is repaired only when the literal sets a field whose type the
// repair changed to exactly the type the compiler reports, so an unrelated error in the same
// tree is left alone.
func planCallSites(repo string, errs []*compileErr, changes []*typeChange) []*callSiteFix {
	var out []*callSiteFix
	parsed := map[string]*parsedFile{}
	for _, e := range errs {
		pf := parseFor(repo, parsed, e)
		if pf == nil {
			continue
		}
		lit, err := pf.locate(e.offset)
		if err != nil {
			continue
		}
		change := matchChange(lit, e.destType, changes)
		if change == nil {
			continue
		}
		after, err := repairedExpr(pf, lit.expr, change.toPointer)
		if err != nil {
			continue
		}
		start, end := pf.offsetOf(lit.expr.Pos()), pf.offsetOf(lit.expr.End())
		if start < 0 || end < start {
			continue
		}
		out = append(out, &callSiteFix{
			path:   e.path(repo),
			line:   e.line,
			field:  change.structName + "." + change.fieldName,
			before: string(pf.src[start:end]),
			after:  after,
			edit:   &edit{start: start, end: end, text: after},
		})
	}
	return out
}

// matchChange returns the type change that explains the compiler error at a literal: the field
// name and the type the compiler reports must both be the change's. A literal that names its
// type in full picks the change to that struct, and one that takes its type from the literal
// that holds it accepts the first change that fits.
func matchChange(lit *litAt, destType string, changes []*typeChange) *typeChange {
	var fallback *typeChange
	for _, c := range changes {
		if c.fieldName != lit.field || c.newType != destType {
			continue
		}
		if lit.litType != "" && lit.litType == c.structName {
			return c
		}
		if fallback == nil {
			fallback = c
		}
	}
	return fallback
}

// parsedFile is a Go file parsed for the repairs, kept so that several errors in one file cost
// one parse.
type parsedFile struct {
	src  []byte
	fset *token.FileSet
	file *ast.File
}

// parseFor returns the parsed file that holds the error, parsing it on first use. It returns
// nil for a file that cannot be read or parsed.
func parseFor(repo string, cache map[string]*parsedFile, e *compileErr) *parsedFile {
	path := filepath.Join(repo, filepath.FromSlash(e.path(repo)))
	if pf, ok := cache[path]; ok {
		return pf
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	pf := &parsedFile{src: src, fset: fset, file: file}
	cache[path] = pf
	return pf
}

// offsetOf returns the byte offset of a position within the parsed file.
func (p *parsedFile) offsetOf(pos token.Pos) int {
	if !pos.IsValid() {
		return -1
	}
	return p.fset.Position(pos).Offset
}

// litAt is the struct literal context of one expression: the expression itself, the field name
// of the key-value pair that holds it, and the name of the type of the composite literal,
// which is empty when the literal takes its type from the literal that holds it.
type litAt struct {
	expr    ast.Expr
	field   string
	litType string
}

// locate returns the context of the expression that starts at the given offset. Several nodes
// start at the same offset, so it keeps the widest one, which is the whole expression.
func (p *parsedFile) locate(offset int) (*litAt, error) {
	var (
		out   *litAt
		stack []ast.Node
	)
	ast.Inspect(p.file, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, n)
		expr, ok := n.(ast.Expr)
		if !ok || p.offsetOf(n.Pos()) != offset {
			return true
		}
		if out != nil && out.expr.End() >= expr.End() {
			return true
		}
		lit := &litAt{expr: expr}
		for i := len(stack) - 2; i >= 0; i-- {
			switch ancestor := stack[i].(type) {
			case *ast.KeyValueExpr:
				if lit.field == "" {
					if id, ok := ancestor.Key.(*ast.Ident); ok {
						lit.field = id.Name
					}
				}
			case *ast.CompositeLit:
				if lit.litType == "" && ancestor.Type != nil {
					lit.litType = baseTypeName(ancestor.Type)
				}
			}
		}
		out = lit
		return true
	})
	if out == nil {
		return nil, fmt.Errorf("no expression at offset %v", offset)
	}
	return out, nil
}

// baseTypeName returns the name of the type a composite literal is declared with, looking
// through the pointers, slices and maps that may hold it: a literal of []*Foo names Foo.
func baseTypeName(e ast.Expr) string {
	for {
		switch t := e.(type) {
		case *ast.StarExpr:
			e = t.X
		case *ast.ArrayType:
			e = t.Elt
		case *ast.MapType:
			e = t.Value
		case *ast.ParenExpr:
			e = t.X
		case *ast.Ident:
			return t.Name
		case *ast.SelectorExpr:
			return t.Sel.Name
		default:
			return ""
		}
	}
}

// repairedExpr returns the text that replaces a literal expression whose field type changed.
// A field that became a pointer wraps the value, and one that became a value unwraps the
// pointer it was given.
func repairedExpr(p *parsedFile, expr ast.Expr, toPointer bool) (string, error) {
	text := p.text(expr)
	if toPointer {
		return "new(" + text + ")", nil
	}
	switch e := expr.(type) {
	case *ast.CallExpr:
		if len(e.Args) == 1 {
			if id, ok := e.Fun.(*ast.Ident); ok && (id.Name == "new" || id.Name == "Ptr") {
				return p.text(e.Args[0]), nil
			}
		}
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			return p.text(e.X), nil
		}
	}
	return "", fmt.Errorf("%v is not new(...), Ptr(...) or &..., so the repair is not mechanical", text)
}

// text returns the source text of an expression.
func (p *parsedFile) text(e ast.Expr) string {
	return string(p.src[p.offsetOf(e.Pos()):p.offsetOf(e.End())])
}

// applyCallSites rewrites the files that hold the planned repairs.
func applyCallSites(repo string, planned []*callSiteFix) (fixed, written int, err error) {
	byFile := map[string][]*callSiteFix{}
	for _, p := range planned {
		byFile[p.path] = append(byFile[p.path], p)
	}
	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	slices.Sort(files)

	for _, file := range files {
		path := filepath.Join(repo, filepath.FromSlash(file))
		src, err := os.ReadFile(path)
		if err != nil {
			return fixed, written, err
		}
		edits := make([]*edit, 0, len(byFile[file]))
		for _, p := range byFile[file] {
			edits = append(edits, p.edit)
		}
		edited, err := applyEdits(src, edits)
		if err != nil {
			return fixed, written, fmt.Errorf("%v: %w", file, err)
		}
		formatted, err := format.Source(edited)
		if err != nil {
			return fixed, written, fmt.Errorf("%v: %w", file, err)
		}
		if err := os.WriteFile(path, formatted, 0o600); err != nil {
			return fixed, written, err
		}
		fixed += len(edits)
		written++
	}
	return fixed, written, nil
}

// modulesToCheck returns the module directories that a type change can break: the module that
// holds each changed file, and every module in the checkout that requires that module, which
// is where a caller of the changed package can live.
func modulesToCheck(repo string, changed []string) []string {
	modules := map[string]bool{}
	// The module path of each changed file's module, which is what another module has to
	// require for the change to reach it.
	paths := map[string]string{}
	for _, file := range changed {
		dir := moduleOf(repo, filepath.Join(repo, filepath.FromSlash(file)))
		if dir == "" {
			continue
		}
		modules[dir] = true
		if path := modulePath(dir); path != "" {
			paths[dir] = path
		}
	}
	if len(paths) > 0 {
		for _, dir := range moduleDirsUnder(repo) {
			if modules[dir] {
				continue
			}
			for _, path := range paths {
				if requiresModule(dir, path) {
					modules[dir] = true
					break
				}
			}
		}
	}
	out := make([]string, 0, len(modules))
	for dir := range modules {
		out = append(out, dir)
	}
	slices.Sort(out)
	return out
}

// moduleOf returns the directory of the module that holds the file, by walking up to the
// nearest go.mod. It returns an empty string when the file is not in a module, including when
// the walk leaves the repository.
func moduleOf(repo, path string) string {
	for dir := filepath.Dir(path); strings.HasPrefix(dir, repo); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return ""
}

// moduleDirsUnder returns the directory of every module in the checkout, skipping the
// directories that never hold one worth checking.
func moduleDirsUnder(repo string) []string {
	var dirs []string
	_ = filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch {
			case path == repo:
			case d.Name() == ".git", d.Name() == "testdata", d.Name() == "vendor", strings.HasPrefix(d.Name(), "."):
				return fs.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	slices.Sort(dirs)
	return dirs
}

// modulePath returns the module path declared by the go.mod in dir, or "".
func modulePath(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	for line := range strings.Lines(string(data)) {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

// requiresModule reports whether the module in dir requires the given module path, which is how
// a module other than the one that changed can still be broken by the change. The module path
// the file declares itself is not a requirement, so a module whose own path merely starts with
// the one that changed is left alone.
func requiresModule(dir, path string) bool {
	if path == "" {
		return false
	}
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "module ") {
			continue
		}
		if strings.Contains(line, path) {
			return true
		}
	}
	return false
}
