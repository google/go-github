// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command list-return-structs scans the Go sources of the `github` package and
// reports the struct types that are used strictly as return values — that is,
// structs that are only ever "unmarshaled" from the GitHub API server and are
// never "marshaled" and sent back to it.
//
// A struct is considered a "return" struct if it appears in the result list of
// any function or method, or if it is (recursively) referenced as a field of
// another return struct. A struct is considered an "input" struct if it appears
// as the type of a parameter named `opts` or `body` in any function or method,
// or if it is (recursively) referenced as a field of another input struct. The
// reported set is the return structs minus the input structs.
//
// Usage:
//
//	go run tools/list-return-structs/main.go [-omitempty] [dir]
//
// The optional `dir` argument defaults to `github` (the package directory,
// relative to the current working directory). With -omitempty (or --omitempty)
// the report is reduced to structs that contain at least one field whose `json`
// struct tag includes "omitempty" or "omitzero".
//
// Structs are listed one per line in a unix-standard grep-like format, e.g.:
//
//	github/rate_limit.go:14:type Rate struct {
package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	log.SetFlags(0)

	omitempty := flag.Bool("omitempty", false, "reduce the report to structs that have at least one field tagged with \"omitempty\" or \"omitzero\"")
	flag.Usage = func() {
		log.Print("Usage: list-return-structs [-omitempty] [dir]")
		flag.PrintDefaults()
	}
	flag.Parse()

	dir := "github"
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	structs, err := analyze(dir, *omitempty)
	if err != nil {
		log.Fatalf("analyzing %s: %v", dir, err)
	}

	for _, s := range structs {
		fmt.Printf("%v:%v:type %v struct {\n", s.relPath, s.line, s.name)
	}
}

// structInfo describes a struct type declaration found while scanning the
// package.
type structInfo struct {
	name       string          // the struct's type name
	relPath    string          // display path, e.g. "github/rate_limit.go"
	line       int             // line of the "type X struct {" declaration
	structType *ast.StructType // the parsed struct type, or nil for non-struct specs
}

// parsePackage parses every non-test .go file in dir and returns the file set
// and parsed files. baseDir is the directory whose name should be preserved as
// the leading path component in structInfo.relPath (e.g. "github").
func parsePackage(dir string) (*token.FileSet, []*ast.File, string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, nil, "", fmt.Errorf("resolving dir: %w", err)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, nil, "", fmt.Errorf("reading dir: %w", err)
	}

	// relPath is computed relative to the parent of absDir so that the path
	// keeps its package directory component (e.g. "github/rate_limit.go")
	// regardless of how dir was specified.
	baseDir := filepath.Dir(absDir)

	fset := token.NewFileSet()
	var files []*ast.File
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(absDir, name)
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, "", fmt.Errorf("parsing %v: %w", path, err)
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, nil, "", fmt.Errorf("no non-test .go files found in %v", absDir)
	}

	return fset, files, baseDir, nil
}

// collectStructs builds a map of struct name -> structInfo for every struct type
// declared in the parsed files.
func collectStructs(fset *token.FileSet, files []*ast.File, baseDir string) map[string]*structInfo {
	structs := map[string]*structInfo{}

	for _, file := range files {
		absPath := fset.Position(file.Package).Filename
		relPath, err := filepath.Rel(baseDir, absPath)
		if err != nil {
			relPath = filepath.Base(absPath)
		}
		// Always report paths with forward slashes so the output matches the
		// unix-standard format on every platform (filepath.Rel uses the OS
		// separator, which is '\' on Windows).
		relPath = filepath.ToSlash(relPath)

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				structs[typeSpec.Name.Name] = &structInfo{
					name:       typeSpec.Name.Name,
					relPath:    relPath,
					line:       fset.Position(typeSpec.Pos()).Line,
					structType: structType,
				}
			}
		}
	}

	return structs
}

// collectDirectUses inspects every function and method declaration and returns
// two sets of struct names:
//   - returns: struct names appearing in a result position
//   - inputs: struct names appearing as the type of a parameter named `opts` or
//     `body`
//
// Only names that refer to package-local structs are of interest, but this
// function returns raw names; callers should intersect with the known structs
// map (closure expansion already does so).
func collectDirectUses(files []*ast.File) (returns, inputs map[string]bool) {
	returns = map[string]bool{}
	inputs = map[string]bool{}

	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type == nil {
				continue
			}

			// Input parameters named `opts` or `body`.
			if fn.Type.Params != nil {
				for _, field := range fn.Type.Params.List {
					for _, name := range field.Names {
						if name.Name != "opts" && name.Name != "body" {
							continue
						}
						if base := baseStructName(field.Type); base != "" {
							inputs[base] = true
						}
					}
				}
			}

			// Result types — every named struct returned counts as a return use.
			if fn.Type.Results != nil {
				for _, field := range fn.Type.Results.List {
					if base := baseStructName(field.Type); base != "" {
						returns[base] = true
					}
				}
			}
		}
	}

	return returns, inputs
}

// transitiveClosure expands a starting set of struct names to include every
// package-local struct reachable as a field type (directly, or via pointers,
// slices, arrays, maps, or generic instantiations), recursively.
func transitiveClosure(start map[string]bool, structs map[string]*structInfo) map[string]bool {
	closure := map[string]bool{}
	var work []string
	for name := range start {
		closure[name] = true
		work = append(work, name)
	}

	for len(work) > 0 {
		name := work[0]
		work = work[1:]
		info, ok := structs[name]
		if !ok || info.structType == nil || info.structType.Fields == nil {
			continue
		}
		for _, field := range info.structType.Fields.List {
			if base := baseStructName(field.Type); base != "" {
				if structs[base] != nil && !closure[base] {
					closure[base] = true
					work = append(work, base)
				}
			}
		}
	}

	return closure
}

// baseStructName reduces a type expression to the name of its underlying named
// type, unwrapping pointers, slices, arrays, maps, ellipses, and generic
// instantiations. It returns "" for types that are not (or cannot be resolved
// as) a simple package-local identifier, such as qualified names (time.Time) or
// anonymous types.
func baseStructName(typ ast.Expr) string {
	switch t := typ.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return baseStructName(t.X)
	case *ast.ArrayType:
		return baseStructName(t.Elt)
	case *ast.Ellipsis:
		return baseStructName(t.Elt)
	case *ast.MapType:
		return baseStructName(t.Value)
	case *ast.IndexExpr:
		return baseStructName(t.X)
	case *ast.IndexListExpr:
		return baseStructName(t.X)
	case *ast.ParenExpr:
		return baseStructName(t.X)
	}
	return ""
}

// hasJSONOmitempty reports whether the struct has at least one field whose `json`
// struct tag includes the "omitempty" or "omitzero" option.
func hasJSONOmitempty(st *ast.StructType) bool {
	if st == nil || st.Fields == nil {
		return false
	}
	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}
		tag := strings.Trim(field.Tag.Value, "`")
		jsonTag, ok := structTagLookup(tag, "json")
		if !ok {
			continue
		}
		// A json tag looks like `name,omitempty` or `name,omitempty,omitzero`.
		for opt := range strings.SplitSeq(jsonTag, ",") {
			if opt == "omitempty" || opt == "omitzero" {
				return true
			}
		}
	}
	return false
}

// structTagLookup returns the value associated with the given key in a
// struct tag string, following the reflect.StructTag lookup rules.
func structTagLookup(tag, key string) (string, bool) {
	for tag != "" {
		// Skip leading whitespace.
		tag = strings.TrimLeft(tag, " \t")
		if tag == "" {
			break
		}
		// Scan the key up to ':'.
		i := strings.Index(tag, ":")
		if i < 0 {
			break
		}
		k := tag[:i]
		rest := tag[i+1:]
		if rest == "" || rest[0] != '"' {
			break
		}
		// Scan the quoted value.
		v, err := unquoteStructTag(rest)
		if err != nil {
			break
		}
		tag = rest[len(v)+2:]
		if k == key {
			return v, true
		}
	}
	return "", false
}

// unquoteStructTag unquotes a single double-quoted value at the start of s,
// returning the unquoted contents.
func unquoteStructTag(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' {
		return "", errors.New("bad quote")
	}
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++ // skip escaped char
		case '"':
			return s[1:i], nil
		}
	}
	return "", errors.New("unterminated quote")
}

// analyze ties the steps together: parse the package, collect structs, compute
// the return and input closures, and return the structs that are return-only,
// sorted by file path then line. When omitempty is true, only structs with at
// least one omitempty/omitzero json-tagged field are returned.
func analyze(dir string, omitempty bool) ([]*structInfo, error) {
	fset, files, baseDir, err := parsePackage(dir)
	if err != nil {
		return nil, err
	}

	structs := collectStructs(fset, files, baseDir)
	returns, inputs := collectDirectUses(files)

	returnClosure := transitiveClosure(returns, structs)
	inputClosure := transitiveClosure(inputs, structs)

	var result []*structInfo
	for name, info := range structs {
		if !returnClosure[name] {
			continue
		}
		if inputClosure[name] {
			continue
		}
		if omitempty && !hasJSONOmitempty(info.structType) {
			continue
		}
		result = append(result, info)
	}

	slices.SortFunc(result, func(a, b *structInfo) int {
		if c := cmp.Compare(a.relPath, b.relPath); c != 0 {
			return c
		}
		return cmp.Compare(a.line, b.line)
	})

	return result, nil
}
