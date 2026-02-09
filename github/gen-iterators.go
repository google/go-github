// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// gen-iterators generates iterator methods for List methods.
//
// It is meant to be used by go-github contributors in conjunction with the
// go generate tool before sending a PR to GitHub.
// Please see the CONTRIBUTING.md file for more information.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"
)

const (
	fileSuffix = "-iterators.go"
)

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")

	sourceTmpl = template.Must(template.New("source").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
	}).Parse(source))

	testTmpl = template.Must(template.New("test").Parse(test))
)

func logf(fmt string, args ...any) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

func main() {
	flag.Parse()
	fset := token.NewFileSet()

	// Parse the current directory
	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	for pkgName, pkg := range pkgs {
		t := &templateData{
			filename: pkgName + fileSuffix,
			Package:  pkgName,
			Methods:  []*method{},
			Structs:  make(map[string]*structDef),
		}

		for _, f := range pkg.Files {
			t.processStructs(f)
		}

		for _, f := range pkg.Files {
			if err := t.processMethods(f); err != nil {
				log.Fatal(err)
			}
		}

		if err := t.dump(); err != nil {
			log.Fatal(err)
		}
	}
	logf("Done.")
}

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix) && !strings.HasPrefix(fi.Name(), "gen-")
}

type templateData struct {
	filename string
	Package  string
	Methods  []*method
	Structs  map[string]*structDef
}

type structDef struct {
	Name   string
	Fields map[string]string
	Embeds []string
}

type method struct {
	RecvType             string
	RecvVar              string
	ClientField          string
	MethodName           string
	IterMethod           string
	Args                 string
	CallArgs             string
	TestCallArgs         string
	ZeroArgs             string
	ReturnType           string
	OptsType             string
	OptsName             string
	OptsIsPtr            bool
	UseListCursorOptions bool
	UseListOptions       bool
	UsePage              bool
	TestJSON1            string
	TestJSON2            string
	TestJSON3            string
}

// customTestJSON maps method names to the JSON response they expect in tests.
// This is needed for methods that internally unmarshal a wrapper struct
// even though they return a slice.
var customTestJSON = map[string]string{
	// Uncomment in #3978: "ListAllTopics":         `{"names": []}`,
	"ListUserInstallations": `{"installations": []}`,
}

func (t *templateData) processStructs(f *ast.File) {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			sd := &structDef{
				Name:   ts.Name.Name,
				Fields: make(map[string]string),
			}

			for _, field := range st.Fields.List {
				typeStr := typeToString(field.Type)
				if len(field.Names) == 0 {
					sd.Embeds = append(sd.Embeds, strings.TrimPrefix(typeStr, "*"))
				} else {
					for _, name := range field.Names {
						sd.Fields[name.Name] = typeStr
					}
				}
			}
			t.Structs[sd.Name] = sd
		}
	}
}

func (t *templateData) hasListCursorOptions(structName string) bool {
	return t.hasOptions(structName, "ListCursorOptions")
}

func (t *templateData) hasListOptions(structName string) bool {
	return t.hasOptions(structName, "ListOptions")
}

func (t *templateData) hasOptions(structName, optionsType string) bool {
	sd, ok := t.Structs[structName]
	if !ok {
		return false
	}
	for _, embed := range sd.Embeds {
		if embed == optionsType {
			return true
		}
		if t.hasOptions(embed, optionsType) {
			return true
		}
	}
	return false
}

func (t *templateData) hasIntPage(structName string) bool {
	sd, ok := t.Structs[structName]
	if !ok {
		return false
	}
	if typeStr, ok := sd.Fields["Page"]; ok {
		return typeStr == "int"
	}
	for _, embed := range sd.Embeds {
		if t.hasIntPage(embed) {
			return true
		}
	}
	return false
}

func getZeroValue(typeStr string) string {
	switch typeStr {
	case "int", "int64", "int32":
		return "0"
	case "string":
		return `""`
	case "bool":
		return "false"
	case "context.Context":
		return "t.Context()"
	default:
		return "nil"
	}
}

func (t *templateData) processMethods(f *ast.File) error {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil {
			continue
		}

		if !fd.Name.IsExported() || !strings.HasPrefix(fd.Name.Name, "List") {
			continue
		}

		if strings.Contains(fd.Name.Name, "MatchingRefs") {
			continue
		}

		if fd.Type.Results == nil || len(fd.Type.Results.List) != 3 {
			continue
		}

		sliceRet, ok := fd.Type.Results.List[0].Type.(*ast.ArrayType)
		if !ok {
			continue
		}
		eltType := typeToString(sliceRet.Elt)

		if typeToString(fd.Type.Results.List[1].Type) != "*Response" {
			continue
		}
		if typeToString(fd.Type.Results.List[2].Type) != "error" {
			continue
		}

		recvType := typeToString(fd.Recv.List[0].Type)
		if !strings.HasPrefix(recvType, "*") || !strings.HasSuffix(recvType, "Service") {
			continue
		}
		recvVar := ""
		if len(fd.Recv.List[0].Names) > 0 {
			recvVar = fd.Recv.List[0].Names[0].Name
		}

		args := []string{}
		callArgs := []string{}
		testCallArgs := []string{}
		zeroArgs := []string{}
		var optsType string
		var optsName string
		hasOpts := false
		optsIsPtr := false

		for _, field := range fd.Type.Params.List {
			typeStr := typeToString(field.Type)
			zeroArg := getZeroValue(typeStr)
			for _, name := range field.Names {
				args = append(args, fmt.Sprintf("%s %s", name.Name, typeStr))
				callArgs = append(callArgs, name.Name)
				zeroArgs = append(zeroArgs, zeroArg)

				if strings.HasSuffix(typeStr, "Options") {
					optsType = strings.TrimPrefix(typeStr, "*")
					optsName = name.Name
					hasOpts = true
					optsIsPtr = strings.HasPrefix(typeStr, "*")
				}
			}
			// second pass: generate testCallArgs after optsName is identified
			for _, name := range field.Names {
				if name.Name == optsName {
					testCallArgs = append(testCallArgs, name.Name)
				} else {
					testCallArgs = append(testCallArgs, zeroArg)
				}
			}
		}

		if !hasOpts {
			continue
		}

		useListCursorOptions := t.hasListCursorOptions(optsType)
		useListOptions := t.hasListOptions(optsType)
		usePage := t.hasIntPage(optsType)

		if !useListCursorOptions && !useListOptions && !usePage {
			logf("Skipping %s.%s: opts %s does not have ListCursorOptions, ListOptions, or Page int", recvType, fd.Name.Name, optsType)
			continue
		}

		recType := strings.TrimPrefix(recvType, "*")
		clientField := strings.TrimSuffix(recType, "Service")
		if clientField == "Migration" {
			clientField = "Migrations"
		}
		if clientField == "s" {
			logf("WARNING: clientField is 's' for %s.%s (recvType=%s)", recvType, fd.Name.Name, recType)
		}

		testJSON, emptyReturnValue := "[]", "{}"
		if val, ok := customTestJSON[fd.Name.Name]; ok {
			testJSON = val
		}
		if eltType == "string" {
			emptyReturnValue = `""`
		}
		testJSON1 := strings.ReplaceAll(testJSON, "[]", fmt.Sprintf("[%v,%[1]v,%[1]v]", emptyReturnValue))       // Call 1 - return 3 items
		testJSON2 := strings.ReplaceAll(testJSON, "[]", fmt.Sprintf("[%v,%[1]v,%[1]v,%[1]v]", emptyReturnValue)) // Call 1 part 2 - return 4 items
		testJSON3 := strings.ReplaceAll(testJSON, "[]", fmt.Sprintf("[%v,%[1]v]", emptyReturnValue))             // Call 2 - return 2 items

		m := &method{
			RecvType:             recType,
			RecvVar:              recvVar,
			ClientField:          clientField,
			MethodName:           fd.Name.Name,
			IterMethod:           fd.Name.Name + "Iter",
			Args:                 strings.Join(args, ", "),
			CallArgs:             strings.Join(callArgs, ", "),
			TestCallArgs:         strings.Join(testCallArgs, ", "),
			ZeroArgs:             strings.Join(zeroArgs, ", "),
			ReturnType:           eltType,
			OptsType:             optsType,
			OptsName:             optsName,
			OptsIsPtr:            optsIsPtr,
			UseListCursorOptions: useListCursorOptions,
			UseListOptions:       useListOptions,
			UsePage:              usePage,
			TestJSON1:            testJSON1,
			TestJSON2:            testJSON2,
			TestJSON3:            testJSON3,
		}
		t.Methods = append(t.Methods, m)
	}
	return nil
}

func typeToString(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + typeToString(x.X)
	case *ast.SelectorExpr:
		return typeToString(x.X) + "." + x.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(x.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", typeToString(x.Key), typeToString(x.Value))
	default:
		return ""
	}
}

func (t *templateData) dump() error {
	if len(t.Methods) == 0 {
		return nil
	}

	slices.SortStableFunc(t.Methods, func(a, b *method) int {
		if a.RecvType != b.RecvType {
			return strings.Compare(a.RecvType, b.RecvType)
		}
		return strings.Compare(a.MethodName, b.MethodName)
	})

	processTemplate := func(tmpl *template.Template, filename string) error {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, t); err != nil {
			return err
		}
		clean, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("format.Source: %v\n%s", err, buf.String())
		}
		logf("Writing %v...", filename)
		return os.WriteFile(filename, clean, 0o644)
	}

	if err := processTemplate(sourceTmpl, t.filename); err != nil {
		return err
	}
	return processTemplate(testTmpl, strings.ReplaceAll(t.filename, ".go", "_test.go"))
}

const doNotEditHeader = `// Code generated by gen-iterators; DO NOT EDIT.
// Instead, please run "go generate ./..." as described here:
// https://github.com/google/go-github/blob/master/CONTRIBUTING.md#submitting-a-patch

// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`

const source = doNotEditHeader + `
package {{.Package}}

import (
	"context"
	"iter"
)

{{range .Methods}}
// {{.IterMethod}} returns an iterator that paginates through all results of {{.MethodName}}.
func ({{.RecvVar}} *{{.RecvType}}) {{.IterMethod}}({{.Args}}) iter.Seq2[{{.ReturnType}}, error] {
	return func(yield func({{.ReturnType}}, error) bool) {
		{{if .OptsIsPtr -}}
		// Create a copy of opts to avoid mutating the caller's struct
		if {{.OptsName}} == nil {
			{{.OptsName}} = &{{.OptsType}}{}
		} else {
			{{.OptsName}} = Ptr(*{{.OptsName}})
		}

		{{end}}
		for {
			items, resp, err := {{.RecvVar}}.{{.MethodName}}({{.CallArgs}})
			if err != nil {
				yield({{if hasPrefix .ReturnType "*"}}nil{{else}}*new({{.ReturnType}}){{end}}, err)
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			{{if and .UseListCursorOptions .UseListOptions}}
			if resp.Cursor == "" && resp.NextPage == 0 {
				break
			}
			{{.OptsName}}.ListCursorOptions.Cursor = resp.Cursor
			{{.OptsName}}.ListOptions.Page = resp.NextPage
			{{else if .UseListCursorOptions}}
			if resp.Cursor == "" {
				break
			}
			{{.OptsName}}.ListCursorOptions.Cursor = resp.Cursor
			{{else if .UseListOptions}}
			if resp.NextPage == 0 {
				break
			}
			{{.OptsName}}.ListOptions.Page = resp.NextPage
			{{else}}
			if resp.NextPage == 0 {
				break
			}
			{{.OptsName}}.Page = resp.NextPage
			{{end -}}
		}
	}
}
{{end}}
`

const test = doNotEditHeader + `
package {{.Package}}

import (
	"fmt"
	"net/http"
	"testing"
)

{{range .Methods}}
func Test{{.RecvType}}_{{.IterMethod}}(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	var callNum int
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		callNum++
		switch callNum {
		case 1:
			{{- if .UseListCursorOptions}}
			w.Header().Set("Link", ` + "`" + `<https://api.github.com/?cursor=yo>; rel="next"` + "`" + `)
			{{else}}
			w.Header().Set("Link", ` + "`" + `<https://api.github.com/?page=1>; rel="next"` + "`" + `)
			{{end -}}
			fmt.Fprint(w, ` + "`" + `{{.TestJSON1}}` + "`" + `)
		case 2:
			fmt.Fprint(w, ` + "`" + `{{.TestJSON2}}` + "`" + `)
		case 3:
			fmt.Fprint(w, ` + "`" + `{{.TestJSON3}}` + "`" + `)
		case 4:
			w.WriteHeader(http.StatusNotFound)
		case 5:
			fmt.Fprint(w, ` + "`" + `{{.TestJSON3}}` + "`" + `)
		}
	})

	iter := client.{{.ClientField}}.{{.IterMethod}}({{.ZeroArgs}})
	var gotItems int
	for _, err := range iter {
		gotItems++
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	if want := 7; gotItems != want {
		t.Errorf("client.{{.ClientField}}.{{.IterMethod}} call 1 got %v items; want %v", gotItems, want)
	}

	{{.OptsName}} := &{{.OptsType}}{}
	iter = client.{{.ClientField}}.{{.IterMethod}}({{.TestCallArgs}})
	gotItems = 0
	for _, err := range iter {
		gotItems++
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	if want := 2; gotItems != want {
		t.Errorf("client.{{.ClientField}}.{{.IterMethod}} call 2 got %v items; want %v", gotItems, want)
	}

	iter = client.{{.ClientField}}.{{.IterMethod}}({{.ZeroArgs}})
	gotItems = 0
	for _, err := range iter {
		gotItems++
		if err == nil {
			t.Error("expected error; got nil")
		}
	}
	if gotItems != 1 {
		t.Errorf("client.{{.ClientField}}.{{.IterMethod}} call 3 got %v items; want 1 (an error)", gotItems)
	}

	iter = client.{{.ClientField}}.{{.IterMethod}}({{.ZeroArgs}})
	gotItems = 0
	iter(func(item {{.ReturnType}}, err error) bool {
		gotItems++
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		return false
	})
	if gotItems != 1 {
		t.Errorf("client.{{.ClientField}}.{{.IterMethod}} call 4 got %v items; want 1 (an error)", gotItems)
	}
}
{{end}}
`
