// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonfieldname is a custom linter to be used by
// golangci-lint to find instances where the Go field name
// of a struct does not match the JSON tag name.
// It honors idiomatic Go initialisms and handles the
// special case of `Github` vs `GitHub` as agreed upon
// by the original author of the repo.
package jsonfieldname

import (
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("jsonfieldname", New)
}

// JSONFieldNamePlugin is a custom linter plugin for golangci-lint.
type JSONFieldNamePlugin struct {
	allowedExceptions map[string]bool
}

// Settings is the configuration for the jsonfieldname linter.
type Settings struct {
	AllowedExceptions []string `json:"allowed-exceptions" yaml:"allowed-exceptions"`
}

// New returns an analysis.Analyzer to use with golangci-lint.
// It parses the "allowed-exceptions" section to determine which warnings to skip.
func New(cfg any) (register.LinterPlugin, error) {
	allowedExceptions := map[string]bool{}

	if cfg != nil {
		if settingsMap, ok := cfg.(map[string]any); ok {
			if exceptionsRaw, ok := settingsMap["allowed-exceptions"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							allowedExceptions[exception] = true
						}
					}
				}
			}
		}
	}

	return &JSONFieldNamePlugin{allowedExceptions: allowedExceptions}, nil
}

// BuildAnalyzers builds the analyzers for the JSONFieldNamePlugin.
func (f *JSONFieldNamePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "jsonfieldname",
			Doc:  "Reports mismatches between Go field and JSON tag names. Note that the JSON tag name is the source-of-truth and the Go field name needs to match it.",
			Run: func(pass *analysis.Pass) (any, error) {
				return run(pass, f.allowedExceptions)
			},
		},
	}, nil
}

// GetLoadMode returns the load mode for the JSONFieldNamePlugin.
func (f *JSONFieldNamePlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass, allowedExceptions map[string]bool) (any, error) {
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
					if !ok || jsonTagName == "-" {
						continue
					}
					jsonTagName = strings.TrimSuffix(jsonTagName, ",omitempty")

					checkGoFieldName(structName, goField.Name, jsonTagName, goField.Pos(), pass, allowedExceptions)
				}
			}

			return true
		})
	}
	return nil, nil
}

func checkGoFieldName(structName, goFieldName, jsonTagName string, tokenPos token.Pos, pass *analysis.Pass, allowedExceptions map[string]bool) {
	fullName := structName + "." + goFieldName
	if allowedExceptions[fullName] {
		return
	}

	want, alternate := jsonTagToPascal(jsonTagName)
	if goFieldName != want && goFieldName != alternate {
		const msg = "change Go field name %q to %q for JSON tag %q in struct %q"
		pass.Reportf(tokenPos, msg, goFieldName, want, jsonTagName, structName)
	}
}

func splitJSONTag(jsonTagName string) []string {
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

func jsonTagToPascal(jsonTagName string) (want, alternate string) {
	parts := splitJSONTag(jsonTagName)
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
