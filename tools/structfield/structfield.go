// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package structfield is a custom linter to be used by
// golangci-lint to find instances where the Go field name
// of a struct does not match the JSON or URL tag name.
// It honors idiomatic Go initialisms and handles the
// special case of `Github` vs `GitHub` as agreed upon
// by the original author of the repo.
// It also checks that fields with "omitempty" tags are reference types
// for `json` struct tags and value types for `url` struct tags.
package structfield

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("structfield", New)
}

// StructFieldPlugin is a custom linter plugin for golangci-lint.
type StructFieldPlugin struct {
	allowedTagNames map[string]bool
	allowedTagTypes map[string]bool
}

// Settings is the configuration for the structfield linter.
type Settings struct {
	AllowedTagNames []string `json:"allowed-tag-names" yaml:"allowed-tag-names"`
	AllowedTagTypes []string `json:"allowed-tag-types" yaml:"allowed-tag-types"`
}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(cfg any) (register.LinterPlugin, error) {
	allowedTagNames := map[string]bool{}
	allowedTagTypes := map[string]bool{}

	if cfg != nil {
		if settingsMap, ok := cfg.(map[string]any); ok {
			if exceptionsRaw, ok := settingsMap["allowed-tag-names"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							allowedTagNames[exception] = true
						}
					}
				}
			}

			if exceptionsRaw, ok := settingsMap["allowed-tag-types"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							allowedTagTypes[exception] = true
						}
					}
				}
			}
		}
	}

	return &StructFieldPlugin{
		allowedTagNames: allowedTagNames,
		allowedTagTypes: allowedTagTypes,
	}, nil
}

// BuildAnalyzers builds the analyzers for the StructFieldPlugin.
func (f *StructFieldPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "structfield",
			Doc: `Reports mismatches between Go field and JSON or URL tag names and types.
Note that the JSON or URL tag name is the source-of-truth and the Go field name needs to match it.
If the "json" tag contains "omitempty", then the Go field must be a reference type.
If the "url" tag contains "omitempty", then the Go field must be a value type (unless it is a timestamp).`,
			Run: func(pass *analysis.Pass) (any, error) {
				return run(pass, f.allowedTagNames, f.allowedTagTypes)
			},
		},
	}, nil
}

// GetLoadMode returns the load mode for the StructFieldPlugin.
func (f *StructFieldPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass, allowedTagNames, allowedTagTypes map[string]bool) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			t, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			structType, ok := t.Type.(*ast.StructType)
			if !ok {
				return true
			}

			// Check only exported
			if !ast.IsExported(t.Name.Name) {
				return true
			}

			for _, field := range structType.Fields.List {
				if field.Tag == nil || len(field.Names) == 0 {
					continue
				}

				processStructField(t.Name.Name, field, pass, allowedTagNames, allowedTagTypes)
			}

			return true
		})
	}
	return nil, nil
}

func processStructField(structName string, field *ast.Field, pass *analysis.Pass, allowedTagNames, allowedTagTypes map[string]bool) {
	goField := field.Names[0]
	tagValue := strings.Trim(field.Tag.Value, "`")
	structTag := reflect.StructTag(tagValue)

	processTag(structName, goField, field, structTag, "json", pass, allowedTagNames, allowedTagTypes)
	processTag(structName, goField, field, structTag, "url", pass, allowedTagNames, allowedTagTypes)
}

func processTag(structName string, goField *ast.Ident, field *ast.Field, structTag reflect.StructTag, tagType string, pass *analysis.Pass, allowedTagNames, allowedTagTypes map[string]bool) {
	tagName, ok := structTag.Lookup(tagType)
	if !ok || tagName == "-" {
		return
	}

	hasOmitEmpty := strings.Contains(tagName, ",omitempty")
	hasOmitZero := strings.Contains(tagName, ",omitzero")

	if hasOmitEmpty || hasOmitZero {
		checkGoFieldType(structName, goField.Name, tagType, field, field.Pos(), pass, allowedTagTypes, hasOmitEmpty, hasOmitZero)
		tagName = strings.ReplaceAll(tagName, ",omitzero", "")
		tagName = strings.ReplaceAll(tagName, ",omitempty", "")
	}
	if tagType == "url" {
		tagName = strings.ReplaceAll(tagName, ",comma", "")
	}
	checkGoFieldName(structName, goField.Name, tagType, tagName, goField.Pos(), pass, allowedTagNames)
}

func checkAndReportInvalidTypesForOmitzero(structName, tagType, goFieldName string, fieldType ast.Expr, tokenPos token.Pos, pass *analysis.Pass) bool {
	switch ft := fieldType.(type) {
	case *ast.StarExpr:
		// Check for *[]T where T is builtin - should be []T
		if arrType, ok := ft.X.(*ast.ArrayType); ok {
			if ident, ok := arrType.Elt.(*ast.Ident); ok && isBuiltinType(ident.Name) {
				const msg = "change the %q field type to %q in the struct %q"
				pass.Reportf(tokenPos, msg, goFieldName, "[]"+ident.Name, structName)
			} else if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
				// Check for *[]*T - should be []*T
				if ident, ok := starExpr.X.(*ast.Ident); ok {
					const msg = "change the %q field type to %q in the struct %q"
					pass.Reportf(tokenPos, msg, goFieldName, "[]*"+ident.Name, structName)
				}
			} else {
				checkStructArrayType(structName, goFieldName, arrType, tokenPos, pass)
			}
			return true
		}
		// Check for *int and other pointers to builtin types
		if ident, ok := ft.X.(*ast.Ident); ok {
			if isBuiltinType(ident.Name) {
				const msg = `the %q field in struct %q uses "omitzero" with a primitive type; remove "omitzero" and use only "omitempty" for pointer primitive types"`
				pass.Reportf(tokenPos, msg, goFieldName, structName)
				return true
			}
		}
		// Check for *map - should be map
		if _, ok := ft.X.(*ast.MapType); ok {
			const msg = "change the %q field type to %q in the struct %q"
			pass.Reportf(tokenPos, msg, goFieldName, exprToString(ft.X), structName)
			return true
		}
		return true
	case *ast.MapType:
		return true
	case *ast.ArrayType:
		if tagType == "url" {
			const msg = "the %q field in struct %q uses unsupported omitzero tag for URL tags"
			pass.Reportf(tokenPos, msg, goFieldName, structName)
			return true
		}
		checkStructArrayType(structName, goFieldName, ft, tokenPos, pass)
		return true
	case *ast.Ident:
		if obj := pass.TypesInfo.ObjectOf(ft); obj != nil {
			switch obj.Type().Underlying().(type) {
			case *types.Struct:
				// For Struct - should be *Struct
				const msg = "change the %q field type to %q in the struct %q"
				pass.Reportf(tokenPos, msg, goFieldName, "*"+ft.Name, structName)
				return true
			case *types.Basic:
				if tagType == "url" {
					const msg = "the %q field in struct %q uses unsupported omitzero tag for URL tags"
					pass.Reportf(tokenPos, msg, goFieldName, structName)
					return true
				}
				// For Builtin - should not to be used with omitzero
				const msg = `the %q field in struct %q uses "omitzero" with a primitive type; remove "omitzero", as it is only allowed with structs, maps, and slices`
				pass.Reportf(tokenPos, msg, goFieldName, structName)
				return true
			}
		}
	case *ast.SelectorExpr:
		return true
	default:
		log.Fatalf("unhandled type: %T", ft)
	}
	return false
}

func checkGoFieldName(structName, goFieldName, tagType, tagName string, tokenPos token.Pos, pass *analysis.Pass, allowedNames map[string]bool) {
	fullName := structName + "." + goFieldName
	if allowedNames[fullName] {
		return
	}

	want, alternate := tagNameToPascal(tagName)
	if goFieldName != want && goFieldName != alternate {
		const msg = "change Go field name %q to %q for %v tag %q in struct %q"
		pass.Reportf(tokenPos, msg, goFieldName, want, tagType, tagName, structName)
	}
}

func checkGoFieldType(structName, goFieldName, tagType string, field *ast.Field, tokenPos token.Pos, pass *analysis.Pass, allowedTypes map[string]bool, omitempty, omitzero bool) {
	if allowedTypes[structName+"."+goFieldName] {
		return
	}
	switch {
	case omitzero:
		skipOmitzero := checkAndReportInvalidTypesForOmitzero(structName, tagType, goFieldName, field.Type, tokenPos, pass)
		if !skipOmitzero {
			const msg = `the %q field in struct %q uses "omitzero"; remove "omitzero", as it is only allowed with structs, maps, and slices`
			pass.Reportf(tokenPos, msg, goFieldName, structName)
		}

	case omitempty:
		if newFieldType, ok := checkAndReportInvalidTypes(structName, tagType, goFieldName, field.Type, tokenPos, pass); !ok {
			const msg = `change the %q field type to %q in the struct %q because its %v tag uses "omitempty"`
			pass.Reportf(tokenPos, msg, goFieldName, newFieldType, structName, tagType)
		}
	}
}

func checkAndReportInvalidTypes(structName, tagType, goFieldName string, fieldType ast.Expr, tokenPos token.Pos, pass *analysis.Pass) (newFieldType string, ok bool) {
	switch ft := fieldType.(type) {
	case *ast.StarExpr:
		// Check for *[]T where T is builtin - should be []T
		if arrType, ok := ft.X.(*ast.ArrayType); ok {
			if ident, ok := arrType.Elt.(*ast.Ident); ok && isBuiltinType(ident.Name) {
				const msg = "change the %q field type to %q in the struct %q"
				pass.Reportf(tokenPos, msg, goFieldName, "[]"+ident.Name, structName)
			} else if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
				// Check for *[]*T - should be []*T
				if ident, ok := starExpr.X.(*ast.Ident); ok {
					const msg = "change the %q field type to %q in the struct %q"
					pass.Reportf(tokenPos, msg, goFieldName, "[]*"+ident.Name, structName)
				}
			} else {
				checkStructArrayType(structName, goFieldName, arrType, tokenPos, pass)
			}
		}
		// Check for *map - should be map
		if _, ok := ft.X.(*ast.MapType); ok {
			const msg = "change the %q field type to %q in the struct %q"
			pass.Reportf(tokenPos, msg, goFieldName, exprToString(ft.X), structName)
		}
		if ident, ok := ft.X.(*ast.Ident); ok && tagType == "url" && isBuiltinType(ident.Name) {
			// Remove the pointer for primitives in url tags
			return ident.Name, false
		}
		return "", true
	case *ast.MapType:
		return "", true
	case *ast.ArrayType:
		checkStructArrayType(structName, goFieldName, ft, tokenPos, pass)
		return "", true
	case *ast.SelectorExpr:
		// Check for json.RawMessage
		if ident, ok := ft.X.(*ast.Ident); ok && ident.Name == "json" && ft.Sel.Name == "RawMessage" {
			return "", true
		}
	case *ast.Ident:
		// Check for `any` type
		if ft.Name == "any" {
			return "", true
		}
		if tagType == "url" && isBuiltinType(ft.Name) {
			return "", true
		}
	default:
		log.Fatalf("unhandled type: %T", ft)
	}

	newFieldType = "*" + exprToString(fieldType)
	return newFieldType, false
}

func checkStructArrayType(structName, goFieldName string, arrType *ast.ArrayType, tokenPos token.Pos, pass *analysis.Pass) {
	if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
		if ident, ok := starExpr.X.(*ast.Ident); ok && isBuiltinType(ident.Name) {
			const msg = "change the %q field type to %q in the struct %q"
			pass.Reportf(tokenPos, msg, goFieldName, "[]"+ident.Name, structName)
		}
		return
	}

	if ident, ok := arrType.Elt.(*ast.Ident); ok && ident.Obj != nil {
		if _, ok := ident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
			const msg = "change the %q field type to %q in the struct %q"
			pass.Reportf(tokenPos, msg, goFieldName, "[]*"+ident.Name, structName)
		}
	}
}

func isBuiltinType(typeName string) bool {
	return types.Universe.Lookup(typeName) != nil
}

func exprToString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	default:
		return fmt.Sprintf("%T", e)
	}
}

func splitTag(jsonTagName string) []string {
	jsonTagName = strings.TrimPrefix(jsonTagName, "$")

	if strings.Contains(jsonTagName, "_") {
		return strings.Split(jsonTagName, "_")
	}

	if strings.Contains(jsonTagName, "-") {
		return strings.Split(jsonTagName, "-")
	}

	if strings.ToLower(jsonTagName) == jsonTagName { // single word
		return []string{jsonTagName}
	}

	s := camelCaseRE.ReplaceAllString(jsonTagName, "$1 $2")
	parts := strings.Fields(s)
	for i, part := range parts {
		parts[i] = strings.ToLower(part)
	}

	return parts
}

var camelCaseRE = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func tagNameToPascal(tagName string) (want, alternate string) {
	parts := splitTag(tagName)
	alt := make([]string, len(parts))
	for i, part := range parts {
		alt[i] = part
		if part == "" {
			continue
		}
		upper := strings.ToUpper(part)
		if initialisms[upper] {
			parts[i] = upper
			alt[i] = upper
		} else if specialCase, ok := specialCases[upper]; ok {
			parts[i] = specialCase
			alt[i] = specialCase
		} else if possibleAlternate, ok := possibleAlternates[upper]; ok {
			parts[i] = possibleAlternate
			alt[i] = strings.ToUpper(part[:1]) + part[1:]
		} else {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
			alt[i] = parts[i]
		}
	}
	return strings.Join(parts, ""), strings.Join(alt, "")
}

// Common Go initialisms that should be all caps.
var initialisms = map[string]bool{
	"API": true, "ASCII": true,
	"CAA": true, "CAS": true, "CNAME": true, "CPU": true,
	"CSS": true, "CWE": true, "CVE": true, "CVSS": true,
	"DN": true, "DNS": true,
	"EOF": true, "EPSS": true,
	"GB": true, "GHSA": true, "GPG": true, "GUID": true,
	"HTML": true, "HTTP": true, "HTTPS": true,
	"ID": true, "IDE": true, "IDP": true, "IP": true, "JIT": true,
	"JSON": true,
	"LDAP": true, "LFS": true, "LHS": true,
	"MD5": true, "MS": true, "MX": true,
	"NPM": true, "NTP": true, "NVD": true,
	"OID": true, "OS": true,
	"PEM": true, "PR": true, "QPS": true,
	"RAM": true, "RHS": true, "RPC": true,
	"SAML": true, "SAS": true, "SBOM": true, "SCIM": true,
	"SHA": true, "SHA1": true, "SHA256": true,
	"SKU": true, "SLA": true, "SMTP": true, "SNMP": true,
	"SPDX": true, "SPDXID": true, "SQL": true, "SSH": true,
	"SSL": true, "SSO": true, "SVN": true,
	"TCP": true, "TFVC": true, "TLS": true, "TTL": true,
	"UDP": true, "UI": true, "UID": true, "UUID": true,
	"URI": true, "URL": true, "UTF8": true,
	"VCF": true, "VCS": true, "VM": true,
	"XML": true, "XMPP": true, "XSRF": true, "XSS": true,
}

var specialCases = map[string]string{
	"CPUS":    "CPUs",
	"CWES":    "CWEs",
	"GRAPHQL": "GraphQL",
	"HREF":    "HRef",
	"IDS":     "IDs",
	"IPS":     "IPs",
	"OAUTH":   "OAuth",
	"OPENAPI": "OpenAPI",
	"URLS":    "URLs",
}

var possibleAlternates = map[string]string{
	"ORGANIZATION":  "Org",
	"ORGANIZATIONS": "Orgs",
	"REPOSITORY":    "Repo",
	"REPOSITORIES":  "Repos",
}
