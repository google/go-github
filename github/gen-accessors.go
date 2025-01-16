// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// gen-accessors generates accessor methods for structs with pointer fields.
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
	fileSuffix = "-accessors.go"
)

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")

	sourceTmpl = template.Must(template.New("source").Parse(source))
	testTmpl   = template.Must(template.New("test").Parse(test))

	// skipStructMethods lists "struct.method" combos to skip.
	skipStructMethods = map[string]bool{
		"RepositoryContent.GetContent":    true,
		"Client.GetBaseURL":               true,
		"Client.GetUploadURL":             true,
		"ErrorResponse.GetResponse":       true,
		"RateLimitError.GetResponse":      true,
		"AbuseRateLimitError.GetResponse": true,
	}
	// skipStructs lists structs to skip.
	skipStructs = map[string]bool{
		"Client": true,
	}

	// whitelistSliceGetters lists "struct.field" to add getter method
	whitelistSliceGetters = map[string]bool{
		"PushEvent.Commits": true,
	}
)

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

func main() {
	flag.Parse()
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	for pkgName, pkg := range pkgs {
		t := &templateData{
			filename: pkgName + fileSuffix,
			Year:     2017,
			Package:  pkgName,
			Imports:  map[string]string{},
		}
		for filename, f := range pkg.Files {
			logf("Processing %v...", filename)
			if err := t.processAST(f); err != nil {
				log.Fatal(err)
			}
		}
		if err := t.dump(); err != nil {
			log.Fatal(err)
		}
	}
	logf("Done.")
}

func (t *templateData) processAST(f *ast.File) error {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// Skip unexported identifiers.
			if !ts.Name.IsExported() {
				logf("Struct %v is unexported; skipping.", ts.Name)
				continue
			}
			// Check if the struct should be skipped.
			if skipStructs[ts.Name.Name] {
				logf("Struct %v is in skip list; skipping.", ts.Name)
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range st.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				fieldName := field.Names[0]
				// Skip unexported identifiers.
				if !fieldName.IsExported() {
					logf("Field %v is unexported; skipping.", fieldName)
					continue
				}
				// Check if "struct.method" should be skipped.
				if key := fmt.Sprintf("%v.Get%v", ts.Name, fieldName); skipStructMethods[key] {
					logf("Method %v is skip list; skipping.", key)
					continue
				}

				se, ok := field.Type.(*ast.StarExpr)
				if !ok {
					switch x := field.Type.(type) {
					case *ast.MapType:
						t.addMapType(x, ts.Name.String(), fieldName.String(), false)
						continue
					case *ast.ArrayType:
						if key := fmt.Sprintf("%v.%v", ts.Name, fieldName); whitelistSliceGetters[key] {
							logf("Method %v is whitelist; adding getter method.", key)
							t.addArrayType(x, ts.Name.String(), fieldName.String(), false)
							continue
						}
					}

					logf("Skipping field type %T, fieldName=%v", field.Type, fieldName)
					continue
				}

				switch x := se.X.(type) {
				case *ast.ArrayType:
					t.addArrayType(x, ts.Name.String(), fieldName.String(), true)
				case *ast.Ident:
					t.addIdent(x, ts.Name.String(), fieldName.String())
				case *ast.MapType:
					t.addMapType(x, ts.Name.String(), fieldName.String(), true)
				case *ast.SelectorExpr:
					t.addSelectorExpr(x, ts.Name.String(), fieldName.String())
				default:
					logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
				}
			}
		}
	}
	return nil
}

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix)
}

func (t *templateData) dump() error {
	if len(t.Getters) == 0 {
		logf("No getters for %v; skipping.", t.filename)
		return nil
	}

	// Sort getters by ReceiverType.FieldName.
	slices.SortStableFunc(t.Getters, func(a, b *getter) int {
		return strings.Compare(a.sortVal, b.sortVal)
	})

	processTemplate := func(tmpl *template.Template, filename string) error {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, t); err != nil {
			return err
		}
		clean, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("format.Source:\n%v\n%v", buf.String(), err)
		}

		logf("Writing %v...", filename)
		if err := os.Chmod(filename, 0644); err != nil {
			return fmt.Errorf("os.Chmod(%q, 0644): %v", filename, err)
		}

		if err := os.WriteFile(filename, clean, 0444); err != nil {
			return err
		}

		if err := os.Chmod(filename, 0444); err != nil {
			return fmt.Errorf("os.Chmod(%q, 0444): %v", filename, err)
		}

		return nil
	}

	if err := processTemplate(sourceTmpl, t.filename); err != nil {
		return err
	}
	return processTemplate(testTmpl, strings.ReplaceAll(t.filename, ".go", "_test.go"))
}

func newGetter(receiverType, fieldName, fieldType, zeroValue string, namedStruct bool) *getter {
	return &getter{
		sortVal:      strings.ToLower(receiverType) + "." + strings.ToLower(fieldName),
		ReceiverVar:  strings.ToLower(receiverType[:1]),
		ReceiverType: receiverType,
		FieldName:    fieldName,
		FieldType:    fieldType,
		ZeroValue:    zeroValue,
		NamedStruct:  namedStruct,
	}
}

func (t *templateData) addArrayType(x *ast.ArrayType, receiverType, fieldName string, isAPointer bool) {
	var eltType string
	var ng *getter
	switch elt := x.Elt.(type) {
	case *ast.Ident:
		eltType = elt.String()
		ng = newGetter(receiverType, fieldName, "[]"+eltType, "nil", false)
	case *ast.StarExpr:
		ident, ok := elt.X.(*ast.Ident)
		if !ok {
			return
		}
		ng = newGetter(receiverType, fieldName, "[]*"+ident.String(), "nil", false)
	default:
		logf("addArrayType: type %q, field %q: unknown elt type: %T %+v; skipping.", receiverType, fieldName, elt, elt)
		return
	}

	ng.ArrayType = !isAPointer
	t.Getters = append(t.Getters, ng)
}

func (t *templateData) addIdent(x *ast.Ident, receiverType, fieldName string) {
	var zeroValue string
	var namedStruct = false
	switch x.String() {
	case "int", "int64":
		zeroValue = "0"
	case "string":
		zeroValue = `""`
	case "bool":
		zeroValue = "false"
	case "Timestamp":
		zeroValue = "Timestamp{}"
	default:
		zeroValue = "nil"
		namedStruct = true
	}

	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, x.String(), zeroValue, namedStruct))
}

func (t *templateData) addMapType(x *ast.MapType, receiverType, fieldName string, isAPointer bool) {
	var keyType string
	switch key := x.Key.(type) {
	case *ast.Ident:
		keyType = key.String()
	default:
		logf("addMapType: type %q, field %q: unknown key type: %T %+v; skipping.", receiverType, fieldName, key, key)
		return
	}

	var valueType string
	switch value := x.Value.(type) {
	case *ast.Ident:
		valueType = value.String()
	default:
		logf("addMapType: type %q, field %q: unknown value type: %T %+v; skipping.", receiverType, fieldName, value, value)
		return
	}

	fieldType := fmt.Sprintf("map[%v]%v", keyType, valueType)
	zeroValue := fmt.Sprintf("map[%v]%v{}", keyType, valueType)
	ng := newGetter(receiverType, fieldName, fieldType, zeroValue, false)
	ng.MapType = !isAPointer
	t.Getters = append(t.Getters, ng)
}

func (t *templateData) addSelectorExpr(x *ast.SelectorExpr, receiverType, fieldName string) {
	if strings.ToLower(fieldName[:1]) == fieldName[:1] { // Non-exported field.
		return
	}

	var xX string
	if xx, ok := x.X.(*ast.Ident); ok {
		xX = xx.String()
	}

	switch xX {
	case "time", "json":
		if xX == "json" {
			t.Imports["encoding/json"] = "encoding/json"
		} else {
			t.Imports[xX] = xX
		}
		fieldType := fmt.Sprintf("%v.%v", xX, x.Sel.Name)
		zeroValue := fmt.Sprintf("%v.%v{}", xX, x.Sel.Name)
		if xX == "time" && x.Sel.Name == "Duration" {
			zeroValue = "0"
		}
		t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, zeroValue, false))
	default:
		logf("addSelectorExpr: xX %q, type %q, field %q: unknown x=%+v; skipping.", xX, receiverType, fieldName, x)
	}
}

type templateData struct {
	filename string
	Year     int
	Package  string
	Imports  map[string]string
	Getters  []*getter
}

type getter struct {
	sortVal      string // Lower-case version of "ReceiverType.FieldName".
	ReceiverVar  string // The one-letter variable name to match the ReceiverType.
	ReceiverType string
	FieldName    string
	FieldType    string
	ZeroValue    string
	NamedStruct  bool // Getter for named struct.
	MapType      bool
	ArrayType    bool
}

const source = `// Copyright {{.Year}} The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-accessors; DO NOT EDIT.
// Instead, please run "go generate ./..." as described here:
// https://github.com/google/go-github/blob/master/CONTRIBUTING.md#submitting-a-patch

package {{.Package}}
{{with .Imports}}
import (
  {{- range . -}}
  "{{.}}"
  {{end -}}
)
{{end}}
{{range .Getters}}
{{if .NamedStruct}}
// Get{{.FieldName}} returns the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() *{{.FieldType}} {
  if {{.ReceiverVar}} == nil {
    return {{.ZeroValue}}
  }
  return {{.ReceiverVar}}.{{.FieldName}}
}
{{else if or .MapType .ArrayType }}
// Get{{.FieldName}} returns the {{.FieldName}} {{if .MapType}}map{{else if .ArrayType }}slice{{end}} if it's non-nil, {{if .MapType}}an empty map{{else if .ArrayType }}nil{{end}} otherwise.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() {{.FieldType}} {
  if {{.ReceiverVar}} == nil || {{.ReceiverVar}}.{{.FieldName}} == nil {
    return {{.ZeroValue}}
  }
  return {{.ReceiverVar}}.{{.FieldName}}
}
{{else}}
// Get{{.FieldName}} returns the {{.FieldName}} field if it's non-nil, zero value otherwise.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() {{.FieldType}} {
  if {{.ReceiverVar}} == nil || {{.ReceiverVar}}.{{.FieldName}} == nil {
    return {{.ZeroValue}}
  }
  return *{{.ReceiverVar}}.{{.FieldName}}
}
{{end}}
{{end}}
`

const test = `// Copyright {{.Year}} The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-accessors; DO NOT EDIT.
// Instead, please run "go generate ./..." as described here:
// https://github.com/google/go-github/blob/master/CONTRIBUTING.md#submitting-a-patch

package {{.Package}}
{{with .Imports}}
import (
  "testing"
  {{range . -}}
  "{{.}}"
  {{end -}}
)
{{end}}
{{range .Getters}}
{{if .NamedStruct}}
func Test{{.ReceiverType}}_Get{{.FieldName}}(tt *testing.T) {
  tt.Parallel()
  {{.ReceiverVar}} := &{{.ReceiverType}}{}
  {{.ReceiverVar}}.Get{{.FieldName}}()
  {{.ReceiverVar}} = nil
  {{.ReceiverVar}}.Get{{.FieldName}}()
}
{{else if or .MapType .ArrayType}}
func Test{{.ReceiverType}}_Get{{.FieldName}}(tt *testing.T) {
  tt.Parallel()
  zeroValue := {{.FieldType}}{}
  {{.ReceiverVar}} := &{{.ReceiverType}}{ {{.FieldName}}: zeroValue }
  {{.ReceiverVar}}.Get{{.FieldName}}()
  {{.ReceiverVar}} = &{{.ReceiverType}}{}
  {{.ReceiverVar}}.Get{{.FieldName}}()
  {{.ReceiverVar}} = nil
  {{.ReceiverVar}}.Get{{.FieldName}}()
}
{{else}}
func Test{{.ReceiverType}}_Get{{.FieldName}}(tt *testing.T) {
  tt.Parallel()
  var zeroValue {{.FieldType}}
  {{.ReceiverVar}} := &{{.ReceiverType}}{ {{.FieldName}}: &zeroValue }
  {{.ReceiverVar}}.Get{{.FieldName}}()
  {{.ReceiverVar}} = &{{.ReceiverType}}{}
  {{.ReceiverVar}}.Get{{.FieldName}}()
  {{.ReceiverVar}} = nil
  {{.ReceiverVar}}.Get{{.FieldName}}()
}
{{end}}
{{end}}
`
