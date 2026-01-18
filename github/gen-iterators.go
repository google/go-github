// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// gen-iterators generates iterator methods for List methods.
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
	fileSuffix = "iterators.go"
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
			Package: pkgName,
			Methods: []*method{},
			Structs: make(map[string]*structDef),
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
	Package string
	Methods []*method
	Structs map[string]*structDef
}

type structDef struct {
	Name   string
	Fields map[string]string
	Embeds []string
}

type method struct {
	RecvType       string
	RecvVar        string
	ClientField    string
	MethodName     string
	IterMethod     string
	Args           string
	CallArgs       string
	ZeroArgs       string
	ReturnType     string
	OptsType       string
	OptsName       string
	OptsIsPtr      bool
	UseListOptions bool
	UsePage        bool
	TestJSON       string
}

// customTestJSON maps method names to the JSON response they expect in tests.
// This is needed for methods that internally unmarshal a wrapper struct
// even though they return a slice.
var customTestJSON = map[string]string{
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

func (t *templateData) hasListOptions(structName string) bool {
	sd, ok := t.Structs[structName]
	if !ok {
		return false
	}
	for _, embed := range sd.Embeds {
		if embed == "ListOptions" {
			return true
		}
		if t.hasListOptions(embed) {
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
		return "context.Background()"
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
		zeroArgs := []string{}
		var optsType string
		var optsName string
		hasOpts := false
		optsIsPtr := false

		for _, field := range fd.Type.Params.List {
			typeStr := typeToString(field.Type)
			for _, name := range field.Names {
				args = append(args, fmt.Sprintf("%s %s", name.Name, typeStr))
				callArgs = append(callArgs, name.Name)
				zeroArgs = append(zeroArgs, getZeroValue(typeStr))

				if strings.HasSuffix(typeStr, "Options") {
					optsType = strings.TrimPrefix(typeStr, "*")
					optsName = name.Name
					hasOpts = true
					optsIsPtr = strings.HasPrefix(typeStr, "*")
				}
			}
		}

		if !hasOpts {
			continue
		}

		useListOptions := t.hasListOptions(optsType)
		usePage := t.hasIntPage(optsType)

		if !useListOptions && !usePage {
			logf("Skipping %s.%s: opts %s does not have ListOptions or Page int", recvType, fd.Name.Name, optsType)
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

		testJSON := "[]"
		if val, ok := customTestJSON[fd.Name.Name]; ok {
			testJSON = val
		}

		m := &method{
			RecvType:       recType,
			RecvVar:        recvVar,
			ClientField:    clientField,
			MethodName:     fd.Name.Name,
			IterMethod:     fd.Name.Name + "Iter",
			Args:           strings.Join(args, ", "),
			CallArgs:       strings.Join(callArgs, ", "),
			ZeroArgs:       strings.Join(zeroArgs, ", "),
			ReturnType:     eltType,
			OptsType:       optsType,
			OptsName:       optsName,
			OptsIsPtr:      optsIsPtr,
			UseListOptions: useListOptions,
			UsePage:        usePage,
			TestJSON:       testJSON,
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
		return os.WriteFile(filename, clean, 0644)
	}

	if err := processTemplate(sourceTmpl, "iterators.go"); err != nil {
		return err
	}
	return processTemplate(testTmpl, "iterators_gen_test.go")
}

const source = `// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-iterators; DO NOT EDIT.

package {{.Package}}

import (
	"context"
	"iter"
)

{{range .Methods}}
// {{.IterMethod}} returns an iterator that paginates through all results of {{.MethodName}}.
func ({{.RecvVar}} *{{.RecvType}}) {{.IterMethod}}({{.Args}}) iter.Seq2[{{.ReturnType}}, error] {
	return func(yield func({{.ReturnType}}, error) bool) {
		{{if .OptsIsPtr}}
		// Create a copy of opts to avoid mutating the caller's struct
		if {{.OptsName}} == nil {
			{{.OptsName}} = &{{.OptsType}}{}
		} else {
			optsCopy := *{{.OptsName}}
			{{.OptsName}} = &optsCopy
		}
		{{else}}
		// Opts is value type, already a copy
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

			if resp.NextPage == 0 {
				break
			}
			{{if .UseListOptions}}
			{{.OptsName}}.ListOptions.Page = resp.NextPage
			{{else}}
			{{.OptsName}}.Page = resp.NextPage
			{{end}}
		}
	}
}
{{end}}
`

const test = `// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-iterators; DO NOT EDIT.

package {{.Package}}

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

{{range .Methods}}
func Test{{.RecvType}}_{{.IterMethod}}(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ` + "`" + `{{.TestJSON}}` + "`" + `)
	})

	ctx := context.Background()
	_ = ctx // avoid unused

	// Call iterator with zero values
	iter := client.{{.ClientField}}.{{.IterMethod}}({{.ZeroArgs}})
	for _, err := range iter {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
{{end}}
`
