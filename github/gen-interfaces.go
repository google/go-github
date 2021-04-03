// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-interfaces generates interfaces for each GitHub service to ease mocking.
//
// Embedding an interface into a struct makes it super-easy to only implement
// the methods that you need when mocking a service, so this auto-generates
// all service interfaces using `go generate ./...`.
//
// It is meant to be used by go-github contributors in conjunction with the
// go generate tool before sending a PR to GitHub.
// Please see the CONTRIBUTING.md file for more information.
package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

const (
	ignoreFilePrefix1 = "gen-"
	ignoreFilePrefix2 = "github-"
	outputFileSuffix  = "-interfaces.go"
)

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")

	fset *token.FileSet

	funcMap = template.FuncMap{
		"render": func(fn *ast.FuncDecl) string {
			fn.Recv = nil
			fn.Body = nil
			var buf bytes.Buffer
			printer.Fprint(&buf, fset, fn)
			return strings.ReplaceAll(buf.String(), "func ", "")
		},
	}

	sourceTmpl = template.Must(template.New("source").Funcs(funcMap).Parse(source))
)

func main() {
	flag.Parse()
	fset = token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	for pkgName, pkg := range pkgs {
		t := &templateData{
			filename: pkgName + outputFileSuffix,
			Year:     2021, // No need to change this once set (even in following years).
			Package:  pkgName,
			Imports: map[string]string{
				"context":  "context",
				"io":       "io",
				"net/http": "net/http",
				"net/url":  "net/url",
				"os":       "os",
				"time":     "time",
			},
			services: map[string]*service{},
		}
		for filename, f := range pkg.Files {
			logf("Processing %v...", filename)
			if err := t.processAST(f); err != nil {
				log.Fatal(err)
			}
		}

		if len(t.services) == 0 {
			log.Printf("No services found in package %q.", pkgName)
			continue
		}

		t.sortServicesAndMethods()

		if err := t.dump(); err != nil {
			log.Fatal(err)
		}
	}
	logf("Done.")
}

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") &&
		!strings.HasPrefix(fi.Name(), ignoreFilePrefix1) &&
		!strings.HasPrefix(fi.Name(), ignoreFilePrefix2)
}

type templateData struct {
	filename string
	Year     int
	Package  string
	Imports  map[string]string

	services       map[string]*service
	SortedServices []*service
}

type service struct {
	Name    string
	GenDecl *ast.GenDecl
	Methods []*method
}

type method struct {
	Name     string
	FuncDecl *ast.FuncDecl
}

func (t *templateData) sortServicesAndMethods() {
	t.SortedServices = make([]*service, 0, len(t.services))
	for _, svc := range t.services {
		t.SortedServices = append(t.SortedServices, svc)

		sort.Slice(svc.Methods, func(a, b int) bool { return svc.Methods[a].Name < svc.Methods[b].Name })
	}

	sort.Slice(t.SortedServices, func(a, b int) bool { return t.SortedServices[a].Name < t.SortedServices[b].Name })

}

func (t *templateData) processAST(f *ast.File) error {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok {
			// Skip unexported funcDecl.
			if !fn.Name.IsExported() {
				continue
			}

			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if _, ok := fn.Recv.List[0].Type.(*ast.Ident); ok && fn.Name.Name == "String" {
					// logf("Ignoring FuncDecl: Name=%q", fn.Name.Name)
				} else {
					logf("Found FuncDecl with receiver: Name=%q, Type=%T", fn.Name.Name, fn.Recv.List[0].Type)
					t.processFuncDeclWithRecv(fn)
				}
			} else {
				logf("Ignoring FuncDecl without receiver: Name=%q", fn.Name.Name)
			}
			continue
		}

		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			logf("Ignoring AST decl type %T", decl)
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// Skip unexported identifiers.
			if !ts.Name.IsExported() {
				continue
			}

			// Skip if this is not a service.
			serviceName := ts.Name.String()
			if !strings.HasSuffix(serviceName, "Service") {
				logf("Ignoring GenDecl with ts.Type=%T ts.Name=%q", ts.Type, serviceName)
				continue
			}

			v, ok := t.services[serviceName]
			if !ok {
				v = &service{Name: serviceName}
				t.services[serviceName] = v
			}
			v.GenDecl = gd
			logf("Found service %q", serviceName)
		}
	}
	return nil
}

func (t *templateData) processFuncDeclWithRecv(fn *ast.FuncDecl) {
	serviceName := t.recvServiceName(fn)
	if serviceName == "" {
		return
	}

	svc, ok := t.services[serviceName]
	if !ok {
		svc = &service{Name: serviceName}
		t.services[serviceName] = svc
	}

	svc.Methods = append(svc.Methods, &method{Name: fn.Name.Name, FuncDecl: fn})
}

func (t *templateData) recvServiceName(fn *ast.FuncDecl) string {
	starExpr, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		logf("Ignoring FuncDecl where Type=%T (want *ast.StarExpr): Name=%q", fn.Recv.List[0].Type, fn.Name.Name)
		return ""
	}

	xIdent, ok := starExpr.X.(*ast.Ident)
	if !ok {
		logf("Ignoring FuncDecl where X=%T (want *ast.Ident): Name=%q", starExpr.X, fn.Name.Name)
		return ""
	}

	if xIdent.Obj == nil {
		if strings.HasSuffix(xIdent.Name, "Service") {
			return xIdent.Name
		}
		return ""
	}

	typeSpec, ok := xIdent.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		logf("Ignoring FuncDecl where Decl=%T (want *ast.TypeSpec): Name=%q", xIdent.Obj.Decl, fn.Name.Name)
		return ""
	}

	typeIdent, ok := typeSpec.Type.(*ast.Ident)
	if !ok {
		if strings.HasSuffix(xIdent.Name, "Service") {
			return xIdent.Name
		}
		logf("Ignoring FuncDecl where Type=%T (want *ast.Ident): Name=%q", typeSpec.Type, fn.Name.Name)
		return ""
	}

	recvType := typeIdent.Name
	if recvType != "service" {
		logf("Ignoring FuncDecl where recvType=%q (want service): Name=%q", recvType, fn.Name.Name)
		return ""
	}

	return xIdent.Name
}

func (t *templateData) dump() error {
	if len(t.services) == 0 {
		logf("No services for %v; skipping.", t.filename)
		return nil
	}

	var buf bytes.Buffer
	if err := sourceTmpl.Execute(&buf, t); err != nil {
		return err
	}
	clean, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("failed-to-format source:\n%v", buf.String())
		return err
	}

	logf("Writing %v...", t.filename)
	return ioutil.WriteFile(t.filename, clean, 0644)
}

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

const source = `// Copyright {{.Year}} The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-interfaces; DO NOT EDIT.
// Instead, please run "go generate ./...". See CONTRIBUTING.md for details.

package {{ $package := .Package }}{{ $package }}
{{ with .Imports }}
import (
  {{- range . -}}
  "{{.}}"
  {{ end -}}
)
{{ end }}
{{ range $index, $svc := .SortedServices }}
// {{ $svc.Name }}Interface defines the interface for the {{ $svc.Name }} for easy mocking.
{{ if and $svc.GenDecl $svc.GenDecl.Doc }}
{{- range $i2, $line := $svc.GenDecl.Doc.List  }}
{{- $line.Text }}
{{ end -}}
{{ end -}}
type {{ $svc.Name }}Interface interface {
{{ range $i3, $mthd := $svc.Methods }}
{{ $mthd.FuncDecl | render }}
{{ end }}
}

// {{ $svc.Name }} implements the {{ $svc.Name }}Interface.
var _ {{ $svc.Name }}Interface = &{{ $svc.Name }}{}
{{ end }}
`
