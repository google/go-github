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
// It also checks that fields with "omitempty" tags are reference types.
package structfield

import (
	"fmt"
	"go/ast"
	"go/token"
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
If the tag contains "omitempty", then the Go field must be a reference type.`,
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

func run(pass *analysis.Pass, allowedTagNameExceptions, allowedTagTypeExceptions map[string]bool) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			switch t := n.(type) {
			case *ast.TypeSpec:
				structType, ok := t.Type.(*ast.StructType)
				if !ok {
					return true
				}
				structName := t.Name.Name

				// Only check exported structs.
				if !ast.IsExported(structName) {
					return true
				}

				for _, field := range structType.Fields.List {
					if field.Tag == nil || len(field.Names) == 0 {
						continue
					}

					goField := field.Names[0]

					tagValue := strings.Trim(field.Tag.Value, "`")
					structTag := reflect.StructTag(tagValue)

					jsonTagName, ok := structTag.Lookup("json")
					if ok && jsonTagName != "-" {
						if strings.Contains(jsonTagName, ",omitempty") {
							checkGoFieldType(structName, goField.Name, field, field.Type.Pos(), pass, allowedTagTypeExceptions)
							jsonTagName = strings.ReplaceAll(jsonTagName, ",omitempty", "")
						}
						checkGoFieldName(structName, goField.Name, jsonTagName, goField.Pos(), pass, allowedTagNameExceptions)
					}

					urlTagName, ok := structTag.Lookup("url")
					if ok && urlTagName != "-" {
						if strings.Contains(urlTagName, ",omitempty") {
							checkGoFieldType(structName, goField.Name, field, field.Type.Pos(), pass, allowedTagTypeExceptions)
							urlTagName = strings.ReplaceAll(urlTagName, ",omitempty", "")
						}
						urlTagName = strings.ReplaceAll(urlTagName, ",comma", "")
						checkGoFieldName(structName, goField.Name, urlTagName, goField.Pos(), pass, allowedTagNameExceptions)
					}
				}
			}

			return true
		})
	}
	return nil, nil
}

func checkGoFieldName(structName, goFieldName, tagName string, tokenPos token.Pos, pass *analysis.Pass, allowedExceptions map[string]bool) {
	fullName := structName + "." + goFieldName
	if allowedExceptions[fullName] {
		return
	}

	want, alternate := tagNameToPascal(tagName)
	if goFieldName != want && goFieldName != alternate {
		const msg = "change Go field name %q to %q for tag %q in struct %q"
		pass.Reportf(tokenPos, msg, goFieldName, want, tagName, structName)
	}
}

func checkGoFieldType(structName, goFieldName string, field *ast.Field, tokenPos token.Pos, pass *analysis.Pass, allowedExceptions map[string]bool) {
	fullName := structName + "." + goFieldName
	if allowedExceptions[fullName] {
		return
	}

	var skip bool
	switch fieldType := field.Type.(type) {
	case *ast.StarExpr, *ast.ArrayType, *ast.MapType:
		skip = true
	case *ast.SelectorExpr:
		// check if type is json.RawMessage
		if ident, ok := fieldType.X.(*ast.Ident); ok && ident.Name == "json" && fieldType.Sel.Name == "RawMessage" {
			skip = true
		}
	case *ast.Ident:
		// check if type is `any`
		if fieldType.Name == "any" {
			skip = true
		}
	}
	if !skip {
		const msg = `change the %q field type to %q in the struct %q because its tag uses "omitempty"`
		pass.Reportf(tokenPos, msg, goFieldName, "*"+exprToString(field.Type), structName)
	}
}

func exprToString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
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
	"SAML": true, "SBOM": true, "SCIM": true,
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
