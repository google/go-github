// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// metaOpRe matches the //meta:operation comment that ties a service method to the OpenAPI
// operation it implements.
var metaOpRe = regexp.MustCompile(`(?i)^\s*//\s*meta:operation\s+(\S+)\s+(\S+)\s*$`)

// fieldInfo describes one exported field of a Go struct.
//
// The offsets are byte offsets within the file named by file. The -fix flag uses them to
// rewrite the field declaration in place.
type fieldInfo struct {
	goName    string
	jsonName  string
	typeName  string // the named type of the field, when it is a plain identifier
	isPointer bool
	hasOmit   bool
	omitZero  bool // the omit option in the tag is omitzero rather than omitempty
	omittable bool // the type can represent "absent": pointer, slice, map, interface, selector
	isStruct  bool // the type is a struct declared in this package
	shared    bool // the declaration lists several names, which share one type and tag
	file      string
	line      int
	typOff    int // start of the field type
	typEnd    int // end of the field type
	starOff   int // offset of the "*" when isPointer, else -1
	tagOff    int // offset of the struct tag, else -1
	tagValEnd int // offset of the end of the json tag value, else -1
	omitOff   int // offset of the omit option within the tag, else -1
	omitEnd   int
}

// omitOption returns the omit option in use, or "" when the field has none.
func (f *fieldInfo) omitOption() string {
	switch {
	case !f.hasOmit:
		return ""
	case f.omitZero:
		return "omitzero"
	default:
		return "omitempty"
	}
}

// omits reports whether the struct tag can actually leave the field out of the JSON body.
// omitzero omits the zero value of any type, but omitempty cannot omit a struct, because a
// struct value is never empty. That is why CONTRIBUTING.md asks for omitzero on structs.
func (f *fieldInfo) omits() bool {
	if !f.hasOmit {
		return false
	}
	return f.omitZero || !f.isStruct
}

// canBeAbsent reports whether a field of this type can be left out of the JSON body. A
// pointer, slice, map or interface can be nil, and a type from another package can hold a
// zero value that marshals to nothing.
func canBeAbsent(e ast.Expr) bool {
	switch t := e.(type) {
	case *ast.StarExpr, *ast.ArrayType, *ast.MapType, *ast.InterfaceType, *ast.SelectorExpr:
		return true
	case *ast.Ident:
		return t.Name == "any"
	}
	return false
}

func isPointer(e ast.Expr) bool {
	_, ok := e.(*ast.StarExpr)
	return ok
}

// structInfo describes a struct declared in the scanned package.
type structInfo struct {
	name   string
	file   string
	line   int
	fields []*fieldInfo
}

// field returns the field with the given JSON name, if the struct has one.
func (s *structInfo) field(jsonName string) *fieldInfo {
	for _, f := range s.fields {
		if f.jsonName == jsonName {
			return f
		}
	}
	return nil
}

// opRef names an OpenAPI operation.
type opRef struct {
	method string
	path   string
}

func (o opRef) String() string { return o.method + " " + o.path }

// methodInfo describes a method that takes a request body.
type methodInfo struct {
	funcName string
	file     string
	line     int
	ops      []*opRef
	bodyType string // "" when the method takes no struct body
	bodyPtr  bool
}

// repoInfo is the result of scanning a go-github checkout.
type repoInfo struct {
	structs map[string]*structInfo
	methods []*methodInfo
	files   int
}

// scanRepo parses the Go sources under <repo>/github and returns the structs and the
// methods that take a body.
func scanRepo(repo string) (*repoInfo, error) {
	info := &repoInfo{structs: map[string]*structInfo{}}
	fset := token.NewFileSet()
	dir := filepath.Join(repo, "github")

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		info.files++
		rel, err := filepath.Rel(repo, path)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					info.structs[ts.Name.Name] = collectStruct(fset, rel, ts, st)
				}
			case *ast.FuncDecl:
				if d.Recv == nil { // Only methods are annotated.
					continue
				}
				m := &methodInfo{
					funcName: d.Name.Name,
					file:     rel,
					line:     fset.Position(d.Pos()).Line,
					ops:      parseMetaOps(d.Doc),
				}
				if body, ptr := findBodyParam(d); body != "" {
					m.bodyType, m.bodyPtr = body, ptr
				}
				info.methods = append(info.methods, m)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// A field's type is only known to be a struct once every file has been read.
	for _, si := range info.structs {
		for _, f := range si.fields {
			if f.typeName != "" {
				_, f.isStruct = info.structs[f.typeName]
			}
		}
	}
	return info, nil
}

// parseMetaOps returns the operations named by the //meta:operation comments of doc.
func parseMetaOps(doc *ast.CommentGroup) []*opRef {
	if doc == nil {
		return nil
	}
	var ops []*opRef
	for _, c := range doc.List {
		m := metaOpRe.FindStringSubmatch(c.Text)
		if m == nil {
			continue
		}
		ops = append(ops, &opRef{method: strings.ToUpper(m[1]), path: m[2]})
	}
	return ops
}

// findBodyParam returns the name and pointer-ness of the parameter called "body" when its
// type is a struct declared in this package. paramcheck already requires that name, and
// requires it to be passed by value for converted types.
func findBodyParam(fn *ast.FuncDecl) (name string, pointer bool) {
	if fn.Type.Params == nil {
		return "", false
	}
	for _, p := range fn.Type.Params.List {
		for _, n := range p.Names {
			if n.Name != "body" {
				continue
			}
			t := p.Type
			ptr := false
			if star, ok := t.(*ast.StarExpr); ok {
				t, ptr = star.X, true
			}
			id, ok := t.(*ast.Ident)
			if !ok { // io.Reader or []byte: not a struct body.
				return "", false
			}
			return id.Name, ptr
		}
	}
	return "", false
}

// collectStruct returns the exported fields of a struct declaration.
func collectStruct(fset *token.FileSet, file string, ts *ast.TypeSpec, st *ast.StructType) *structInfo {
	si := &structInfo{
		name: ts.Name.Name,
		file: file,
		line: fset.Position(ts.Pos()).Line,
	}
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 { // Embedded field.
			continue
		}
		for _, n := range f.Names {
			if !n.IsExported() {
				continue
			}
			fi := &fieldInfo{
				goName:    n.Name,
				jsonName:  n.Name,
				isPointer: isPointer(f.Type),
				omittable: canBeAbsent(f.Type),
				shared:    len(f.Names) > 1,
				file:      file,
				line:      fset.Position(n.Pos()).Line,
				typOff:    fset.Position(f.Type.Pos()).Offset,
				typEnd:    fset.Position(f.Type.End()).Offset,
				starOff:   -1,
				tagOff:    -1,
				tagValEnd: -1,
				omitOff:   -1,
			}
			switch t := f.Type.(type) {
			case *ast.Ident:
				fi.typeName = t.Name
			case *ast.StarExpr:
				fi.starOff = fset.Position(t.Star).Offset
				if id, ok := t.X.(*ast.Ident); ok {
					fi.typeName = id.Name
				}
			}
			if f.Tag != nil {
				fi.tagOff = fset.Position(f.Tag.Pos()).Offset
				name, omit, ignored, valStart, valEnd := parseStructTag(f.Tag.Value)
				if ignored {
					continue
				}
				if name != "" {
					fi.jsonName = name
				}
				if valEnd > valStart {
					fi.tagValEnd = fi.tagOff + valEnd
				}
				if omit != "" {
					fi.hasOmit = true
					fi.omitZero = omit == "omitzero"
					if i := strings.Index(f.Tag.Value[valStart:valEnd], ","+omit); i >= 0 {
						fi.omitOff = fi.tagOff + valStart + i
						fi.omitEnd = fi.omitOff + len(omit) + 1
					}
				}
			}
			si.fields = append(si.fields, fi)
		}
	}
	return si
}

// parseStructTag reads the json tag of a raw struct tag literal (backticks included). It
// returns the tag name, the omit option in use, whether the field is ignored with
// json:"-", and the offsets of the tag value within the literal.
func parseStructTag(lit string) (name, omit string, ignored bool, valStart, valEnd int) {
	valStart, valEnd = -1, -1
	i := 1 // Skip the opening backtick.
	for i < len(lit) {
		for i < len(lit) && lit[i] == ' ' {
			i++
		}
		keyStart := i
		for i < len(lit) && lit[i] != ':' && lit[i] != '"' && lit[i] != '`' {
			i++
		}
		if i >= len(lit)-1 || lit[i] != ':' {
			return "", "", false, valStart, valEnd
		}
		key := lit[keyStart:i]
		i++ // Skip the colon.
		if lit[i] != '"' {
			return "", "", false, valStart, valEnd
		}
		i++ // Skip the opening quote.
		start := i
		for i < len(lit) && lit[i] != '"' {
			if lit[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(lit) {
			return "", "", false, valStart, valEnd
		}
		value := lit[start:i]
		i++ // Skip the closing quote.
		if key != "json" {
			continue
		}
		valStart, valEnd = start, i-1
		parts := strings.Split(value, ",")
		if parts[0] == "-" {
			return "", "", true, valStart, valEnd
		}
		if parts[0] != "" {
			name = parts[0]
		}
		for _, o := range parts[1:] {
			if o == "omitempty" || o == "omitzero" {
				omit = o
			}
		}
		return name, omit, false, valStart, valEnd
	}
	return "", "", false, valStart, valEnd
}
