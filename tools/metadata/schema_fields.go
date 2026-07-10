// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmp"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"go.yaml.in/yaml/v3"
)

type schemaFieldCheckOptions struct {
	descriptions     []*openapiFile
	githubDir        string
	schemaNames      map[string]bool
	includeResponses bool
	// exceptions holds "Struct.Field" entries whose diagnostics are suppressed. It is
	// loaded from schema_field_exceptions.yaml by the command; see loadSchemaFieldExceptions.
	exceptions map[string]bool
}

type schemaFieldCheckResult struct {
	Summary     schemaFieldCheckSummary  `json:"summary"`
	Checked     []*schemaFieldChecked    `json:"checked"`
	Skipped     []*schemaFieldSkipped    `json:"skipped,omitempty"`
	Diagnostics []*schemaFieldDiagnostic `json:"diagnostics"`
}

type schemaFieldCheckSummary struct {
	OpenAPISchemas int `json:"openapi_schemas"`
	GoStructs      int `json:"go_structs"`
	Checked        int `json:"checked"`
	Skipped        int `json:"skipped"`
	Diagnostics    int `json:"diagnostics"`
}

type schemaFieldChecked struct {
	OpenAPISchema string `json:"openapi_schema"`
	GoStruct      string `json:"go_struct"`
	OpenAPIFile   string `json:"openapi_file,omitempty"`
	MatchReason   string `json:"match_reason"`
}

type schemaFieldSkipped struct {
	OpenAPISchema string `json:"openapi_schema"`
	OpenAPIFile   string `json:"openapi_file,omitempty"`
	Reason        string `json:"reason"`
}

type schemaFieldDiagnostic struct {
	OpenAPISchema string `json:"openapi_schema"`
	GoStruct      string `json:"go_struct"`
	Field         string `json:"field"`
	JSONName      string `json:"json_name"`
	Message       string `json:"message"`
	Filename      string `json:"filename,omitempty"`
	Line          int    `json:"line,omitempty"`
	OpenAPIFile   string `json:"openapi_file,omitempty"`
}

func (d schemaFieldDiagnostic) String() string {
	loc := diagLocation(d.Filename, d.Line)
	if loc != "" {
		loc += ": "
	}
	source := ""
	if d.OpenAPIFile != "" {
		source = fmt.Sprintf(" [%v]", d.OpenAPIFile)
	}
	return fmt.Sprintf("%v%v.%v (%v from %v): %v%v", loc, d.GoStruct, d.Field, d.JSONName, d.OpenAPISchema, d.Message, source)
}

func diagLocation(filename string, line int) string {
	if filename == "" {
		return ""
	}
	if line == 0 {
		return filename
	}
	return fmt.Sprintf("%v:%v", filename, line)
}

type goStructInfo struct {
	name     string
	filename string
	line     int
	fields   map[string]goStructField
}

type goStructField struct {
	goStruct      string
	field         string
	jsonName      string
	hasOmitOption bool
	isPointer     bool
	canBeOmitted  bool
	filename      string
	line          int
}

type openapiSchemaFields struct {
	openapiSchema string
	openapiFile   string
	required      map[string]bool
	properties    map[string]openapiSchemaProperty
}

type openapiSchemaProperty struct {
	nullable  bool
	readOnly  bool
	writeOnly bool
}

type schemaFieldMatch struct {
	schema      *openapiSchemaFields
	goStruct    *goStructInfo
	matchReason string
}

// schemaFieldExceptionsFile is the on-disk format of schema_field_exceptions.yaml: a list
// of "Struct.Field" entries whose JSON field optionality intentionally deviates from the
// OpenAPI schema, so their diagnostics are suppressed. Each entry is a known deviation
// awaiting cleanup.
type schemaFieldExceptionsFile struct {
	Exceptions []string `yaml:"exceptions"`
}

// loadSchemaFieldExceptions reads the "Struct.Field" exception entries from filename and
// returns them as a set. A missing file yields an empty set and no error when optional is
// true, so callers relying on the default path do not need the file to exist; an explicitly
// requested file that is missing is an error.
func loadSchemaFieldExceptions(filename string, optional bool) (map[string]bool, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		if optional && errors.Is(err, fs.ErrNotExist) {
			return map[string]bool{}, nil
		}
		return nil, err
	}
	var exceptionsFile schemaFieldExceptionsFile
	if err := yaml.Unmarshal(b, &exceptionsFile); err != nil {
		return nil, fmt.Errorf("%v: %w", filename, err)
	}
	return sliceSet(exceptionsFile.Exceptions), nil
}

// filterAllowedSchemaFieldDiagnostics removes diagnostics whose "Struct.Field" is listed in
// the exceptions set.
func filterAllowedSchemaFieldDiagnostics(diagnostics []*schemaFieldDiagnostic, exceptions map[string]bool) []*schemaFieldDiagnostic {
	var kept []*schemaFieldDiagnostic
	for _, diag := range diagnostics {
		if exceptions[diag.GoStruct+"."+diag.Field] {
			continue
		}
		kept = append(kept, diag)
	}
	return kept
}

func checkSchemaFields(opts schemaFieldCheckOptions) (schemaFieldCheckResult, error) {
	if len(opts.descriptions) == 0 {
		return schemaFieldCheckResult{}, errors.New("no OpenAPI descriptions loaded")
	}

	goStructs, requestStructs, err := collectGoStructs(opts.githubDir)
	if err != nil {
		return schemaFieldCheckResult{}, err
	}

	schemas, skipped, err := collectOpenAPISchemaFields(opts.descriptions, opts.schemaNames)
	if err != nil {
		return schemaFieldCheckResult{}, err
	}

	includeResponses := opts.includeResponses || len(opts.schemaNames) > 0
	matches, matchSkipped := matchOpenAPISchemasToGoStructs(schemas, goStructs, requestStructs, len(opts.schemaNames) > 0, includeResponses)
	result := schemaFieldCheckResult{
		Skipped: append(skipped, matchSkipped...),
		Summary: schemaFieldCheckSummary{
			OpenAPISchemas: len(schemas),
			GoStructs:      len(goStructs),
		},
	}

	for _, match := range matches {
		result.Checked = append(result.Checked, &schemaFieldChecked{
			OpenAPISchema: match.schema.openapiSchema,
			GoStruct:      match.goStruct.name,
			OpenAPIFile:   match.schema.openapiFile,
			MatchReason:   match.matchReason,
		})
		result.Diagnostics = append(result.Diagnostics, compareSchemaFields(match)...)
	}

	result.Diagnostics = filterAllowedSchemaFieldDiagnostics(result.Diagnostics, opts.exceptions)
	sortSchemaFieldResult(&result)
	result.Summary.Checked = len(result.Checked)
	result.Summary.Skipped = len(result.Skipped)
	result.Summary.Diagnostics = len(result.Diagnostics)
	return result, nil
}

func sortSchemaFieldResult(result *schemaFieldCheckResult) {
	slices.SortFunc(result.Diagnostics, func(a, b *schemaFieldDiagnostic) int {
		return cmp.Or(
			cmp.Compare(a.GoStruct, b.GoStruct),
			cmp.Compare(a.JSONName, b.JSONName),
			cmp.Compare(a.Field, b.Field),
			cmp.Compare(a.OpenAPISchema, b.OpenAPISchema),
			cmp.Compare(a.OpenAPIFile, b.OpenAPIFile),
			cmp.Compare(a.Message, b.Message),
		)
	})
	slices.SortFunc(result.Checked, func(a, b *schemaFieldChecked) int {
		return cmp.Or(
			cmp.Compare(a.OpenAPISchema, b.OpenAPISchema),
			cmp.Compare(a.GoStruct, b.GoStruct),
			cmp.Compare(a.OpenAPIFile, b.OpenAPIFile),
			cmp.Compare(a.MatchReason, b.MatchReason),
		)
	})
	slices.SortFunc(result.Skipped, func(a, b *schemaFieldSkipped) int {
		return cmp.Or(
			cmp.Compare(a.OpenAPISchema, b.OpenAPISchema),
			cmp.Compare(a.OpenAPIFile, b.OpenAPIFile),
			cmp.Compare(a.Reason, b.Reason),
		)
	})
}

func collectOpenAPISchemaFields(descriptions []*openapiFile, schemaNames map[string]bool) ([]*openapiSchemaFields, []*schemaFieldSkipped, error) {
	var schemas []*openapiSchemaFields
	var skipped []*schemaFieldSkipped
	seen := map[string]bool{}

	for _, desc := range descriptions {
		if desc.description == nil || desc.description.Components == nil || desc.description.Components.Schemas == nil {
			continue
		}

		names := make([]string, 0, len(desc.description.Components.Schemas))
		for name := range desc.description.Components.Schemas {
			if len(schemaNames) > 0 && !schemaNames[name] {
				continue
			}
			if seen[name] {
				continue
			}
			names = append(names, name)
		}
		slices.Sort(names)

		for _, name := range names {
			seen[name] = true
			schemaRef := desc.description.Components.Schemas[name]
			if schemaRef == nil || schemaRef.Value == nil {
				skipped = append(skipped, newSchemaFieldSkipped(name, desc.filename, "schema reference is unresolved"))
				continue
			}

			schema, reason, err := flattenObjectSchema(schemaRef.Value)
			if err != nil {
				return nil, nil, fmt.Errorf("%v %v: %w", desc.filename, name, err)
			}
			if reason != "" {
				skipped = append(skipped, newSchemaFieldSkipped(name, desc.filename, reason))
				continue
			}
			if len(schema.Properties) == 0 {
				skipped = append(skipped, newSchemaFieldSkipped(name, desc.filename, "schema has no object properties"))
				continue
			}

			schemas = append(schemas, &openapiSchemaFields{
				openapiSchema: name,
				openapiFile:   desc.filename,
				required:      sliceSet(schema.Required),
				properties:    schemaProperties(schema.Properties),
			})
		}
	}

	for name := range schemaNames {
		if !seen[name] {
			skipped = append(skipped, newSchemaFieldSkipped(name, "", "schema filter did not match an OpenAPI schema"))
		}
	}

	return schemas, skipped, nil
}

func newSchemaFieldSkipped(schemaName, filename, reason string) *schemaFieldSkipped {
	return &schemaFieldSkipped{
		OpenAPISchema: schemaName,
		OpenAPIFile:   filename,
		Reason:        reason,
	}
}

func flattenObjectSchema(schema *openapi3.Schema) (*openapi3.Schema, string, error) {
	if schema == nil {
		return nil, "", errors.New("schema is nil")
	}
	if hasUnsupportedComposition(schema) {
		return nil, "schema uses oneOf, anyOf, or not", nil
	}
	if len(schema.AllOf) == 0 {
		return schema, "", nil
	}

	merged := &openapi3.Schema{
		Required:   append([]string{}, schema.Required...),
		Properties: openapi3.Schemas{},
	}
	maps.Copy(merged.Properties, schema.Properties)

	for _, ref := range schema.AllOf {
		if ref == nil || ref.Value == nil {
			return nil, "schema contains an unresolved allOf reference", nil
		}
		part, reason, err := flattenObjectSchema(ref.Value)
		if err != nil || reason != "" {
			return nil, reason, err
		}
		merged.Required = append(merged.Required, part.Required...)
		maps.Copy(merged.Properties, part.Properties)
	}

	return merged, "", nil
}

func hasUnsupportedComposition(schema *openapi3.Schema) bool {
	return len(schema.OneOf) > 0 || len(schema.AnyOf) > 0 || schema.Not != nil
}

func schemaProperties(properties openapi3.Schemas) map[string]openapiSchemaProperty {
	result := make(map[string]openapiSchemaProperty, len(properties))
	for name, propRef := range properties {
		prop := openapiSchemaProperty{}
		if propRef != nil && propRef.Value != nil {
			prop.nullable = propRef.Value.Nullable
			prop.readOnly = propRef.Value.ReadOnly
			prop.writeOnly = propRef.Value.WriteOnly
		}
		result[name] = prop
	}
	return result
}

func sliceSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func matchOpenAPISchemasToGoStructs(schemas []*openapiSchemaFields, goStructs map[string]*goStructInfo, requestStructs map[string]bool, allowSchemaNameMatch, includeResponses bool) ([]*schemaFieldMatch, []*schemaFieldSkipped) {
	var matches []*schemaFieldMatch
	var skipped []*schemaFieldSkipped

	for _, schema := range schemas {
		if allowSchemaNameMatch {
			if match, ok, reason := matchBySchemaName(schema, goStructs, requestStructs, includeResponses); ok {
				matches = append(matches, match)
				continue
			} else if reason != "" {
				skipped = append(skipped, newSchemaFieldSkipped(schema.openapiSchema, schema.openapiFile, reason))
				continue
			}
		}
		if match, ok, reason := matchByExactFieldSet(schema, goStructs, requestStructs, includeResponses); ok {
			matches = append(matches, match)
		} else {
			skipped = append(skipped, newSchemaFieldSkipped(schema.openapiSchema, schema.openapiFile, reason))
		}
	}

	return dropAmbiguousFieldSetMatches(matches, skipped)
}

// dropAmbiguousFieldSetMatches removes exact-field-set matches for a Go struct that matched more
// than one OpenAPI schema. A field set that coincidentally equals several unrelated schemas (for
// example a generic {id, type}) is not a reliable match, so it is skipped rather than reported.
func dropAmbiguousFieldSetMatches(matches []*schemaFieldMatch, skipped []*schemaFieldSkipped) ([]*schemaFieldMatch, []*schemaFieldSkipped) {
	fieldSetMatchCount := map[string]int{}
	for _, match := range matches {
		if match.matchReason == "exact JSON field set" {
			fieldSetMatchCount[match.goStruct.name]++
		}
	}

	var kept []*schemaFieldMatch
	for _, match := range matches {
		if match.matchReason == "exact JSON field set" && fieldSetMatchCount[match.goStruct.name] > 1 {
			skipped = append(skipped, newSchemaFieldSkipped(match.schema.openapiSchema, match.schema.openapiFile,
				"Go struct "+match.goStruct.name+" matches multiple schemas by field set"))
			continue
		}
		kept = append(kept, match)
	}
	return kept, skipped
}

func matchBySchemaName(schema *openapiSchemaFields, goStructs map[string]*goStructInfo, requestStructs map[string]bool, includeResponses bool) (*schemaFieldMatch, bool, string) {
	var matches []*goStructInfo
	for _, name := range goNameCandidates(schema.openapiSchema) {
		goStruct, ok := goStructs[name]
		if !ok {
			continue
		}
		if !canCheckGoStruct(goStruct, requestStructs, includeResponses) {
			continue
		}
		if !hasEnoughSharedFields(schema, goStruct) {
			continue
		}
		matches = appendUniqueGoStruct(matches, goStruct)
	}

	switch len(matches) {
	case 0:
		return nil, false, ""
	case 1:
		return &schemaFieldMatch{
			schema:      schema,
			goStruct:    matches[0],
			matchReason: "schema name",
		}, true, ""
	default:
		return nil, false, "ambiguous Go struct name match: " + joinGoStructNames(matches)
	}
}

func matchByExactFieldSet(schema *openapiSchemaFields, goStructs map[string]*goStructInfo, requestStructs map[string]bool, includeResponses bool) (*schemaFieldMatch, bool, string) {
	if len(schema.properties) < 2 {
		return nil, false, "no unambiguous Go struct match"
	}

	var matches []*goStructInfo
	for _, goStruct := range goStructs {
		if !canCheckGoStruct(goStruct, requestStructs, includeResponses) {
			continue
		}
		if sameJSONFieldSet(schema, goStruct) {
			matches = append(matches, goStruct)
		}
	}

	switch len(matches) {
	case 0:
		return nil, false, "no unambiguous Go struct match"
	case 1:
		return &schemaFieldMatch{
			schema:      schema,
			goStruct:    matches[0],
			matchReason: "exact JSON field set",
		}, true, ""
	default:
		slices.SortFunc(matches, func(a, b *goStructInfo) int {
			return cmp.Compare(a.name, b.name)
		})
		return nil, false, "ambiguous Go struct field-set match: " + joinGoStructNames(matches)
	}
}

// canCheckGoStruct reports whether goStruct should be compared against an OpenAPI schema.
// By default only request body structs are checked; requestStructs holds the names of
// structs used as the body argument of a mutating client.NewRequest call. Response and
// other structs are only checked when includeResponses is set.
func canCheckGoStruct(goStruct *goStructInfo, requestStructs map[string]bool, includeResponses bool) bool {
	return includeResponses || requestStructs[goStruct.name]
}

func appendUniqueGoStruct(matches []*goStructInfo, goStruct *goStructInfo) []*goStructInfo {
	for _, existing := range matches {
		if existing.name == goStruct.name {
			return matches
		}
	}
	return append(matches, goStruct)
}

func joinGoStructNames(goStructs []*goStructInfo) string {
	names := make([]string, 0, len(goStructs))
	for _, goStruct := range goStructs {
		names = append(names, goStruct.name)
	}
	slices.Sort(names)
	return strings.Join(names, ", ")
}

// Thresholds for hasEnoughSharedFields. A schema-name match already agrees on the
// Go type name, so these guard only against a name that coincidentally collides with
// an unrelated struct. They are deliberately lenient: a wrong match is dropped later
// as ambiguous, but a missed match silently skips a real check.
const (
	// minSharedFieldsForMatch accepts a match once this many JSON fields overlap,
	// regardless of struct size: three shared, correctly named fields are unlikely
	// to line up by chance.
	minSharedFieldsForMatch = 3
	// minSharedFieldPercent is the fallback for structs smaller than
	// minSharedFieldsForMatch: the overlap must cover at least this share of the
	// smaller field set. Compared as shared*100 >= smallest*minSharedFieldPercent
	// to stay in integer math.
	minSharedFieldPercent = 60
)

// hasEnoughSharedFields reports whether schema and goStruct share enough JSON fields to
// treat a schema-name match as genuine rather than a coincidental name collision.
func hasEnoughSharedFields(schema *openapiSchemaFields, goStruct *goStructInfo) bool {
	shared := sharedFieldCount(schema, goStruct)
	if shared == 0 {
		return false
	}
	smallest := min(len(schema.properties), len(goStruct.fields))
	// A type with one or two fields has too little signal for a percentage test, so
	// require every field to line up before trusting the match.
	if smallest <= 2 {
		return shared == smallest
	}
	return shared >= minSharedFieldsForMatch || shared*100 >= smallest*minSharedFieldPercent
}

func sharedFieldCount(schema *openapiSchemaFields, goStruct *goStructInfo) int {
	var shared int
	for name := range schema.properties {
		if _, ok := goStruct.fields[name]; ok {
			shared++
		}
	}
	return shared
}

func sameJSONFieldSet(schema *openapiSchemaFields, goStruct *goStructInfo) bool {
	if len(schema.properties) != len(goStruct.fields) {
		return false
	}
	for name := range schema.properties {
		if _, ok := goStruct.fields[name]; !ok {
			return false
		}
	}
	return true
}

// goInitialisms maps a lowercase OpenAPI name token to its idiomatic Go casing when
// building candidate struct names in goName. The structfield linter keeps a similar
// list (its `initialisms`/`specialCases`), but it lives in the separate
// github.com/google/go-github/v89/tools/structfield module and is keyed and shaped
// differently (uppercase-keyed set for a different purpose), so the two are
// intentionally not shared: unifying them would require a new shared module that both
// tools depend on. This list is deliberately small and only needs the initialisms that
// actually appear in OpenAPI schema names.
var goInitialisms = map[string]string{
	"api":     "API",
	"apis":    "APIs",
	"gpg":     "GPG",
	"html":    "HTML",
	"http":    "HTTP",
	"https":   "HTTPS",
	"id":      "ID",
	"ids":     "IDs",
	"ip":      "IP",
	"ips":     "IPs",
	"oauth":   "OAuth",
	"oidc":    "OIDC",
	"scim":    "SCIM",
	"sms":     "SMS",
	"sso":     "SSO",
	"ssh":     "SSH",
	"totp":    "TOTP",
	"url":     "URL",
	"urls":    "URLs",
	"webhook": "Webhook",
}

func goNameCandidates(openapiName string) []string {
	tokens := splitOpenAPIName(openapiName)
	if len(tokens) == 0 {
		return nil
	}

	variants := [][]string{tokens}
	allSingular := make([]string, len(tokens))
	var allChanged bool
	for i, token := range tokens {
		singular := singularize(token)
		allSingular[i] = singular
		if singular != token {
			allChanged = true
			variant := append([]string{}, tokens...)
			variant[i] = singular
			variants = append(variants, variant)
		}
	}
	if allChanged {
		variants = append(variants, allSingular)
	}

	var names []string
	seen := map[string]bool{}
	for _, variant := range variants {
		name := goName(variant)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return names
}

func splitOpenAPIName(name string) []string {
	return strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func singularize(token string) string {
	lower := strings.ToLower(token)
	switch {
	case strings.HasSuffix(lower, "ies") && len(token) > 3:
		return token[:len(token)-3] + "y"
	case strings.HasSuffix(lower, "statuses"):
		return token[:len(token)-2]
	case strings.HasSuffix(lower, "ches") || strings.HasSuffix(lower, "shes") || strings.HasSuffix(lower, "xes") || strings.HasSuffix(lower, "ses"):
		return token[:len(token)-2]
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") && len(token) > 1:
		return token[:len(token)-1]
	default:
		return token
	}
}

func goName(tokens []string) string {
	var b strings.Builder
	for _, token := range tokens {
		if token == "" {
			continue
		}
		lower := strings.ToLower(token)
		if initialism, ok := goInitialisms[lower]; ok {
			b.WriteString(initialism)
			continue
		}
		if isVersionToken(lower) {
			b.WriteString(strings.ToUpper(lower[:1]))
			b.WriteString(lower[1:])
			continue
		}
		b.WriteString(strings.ToUpper(token[:1]))
		if len(token) > 1 {
			b.WriteString(strings.ToLower(token[1:]))
		}
	}
	return b.String()
}

func isVersionToken(token string) bool {
	if len(token) < 2 || token[0] != 'v' {
		return false
	}
	for _, r := range token[1:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func compareSchemaFields(match *schemaFieldMatch) []*schemaFieldDiagnostic {
	var diagnostics []*schemaFieldDiagnostic
	for jsonName, field := range match.goStruct.fields {
		prop, inSchema := match.schema.properties[jsonName]
		if !inSchema {
			diagnostics = append(diagnostics, newSchemaFieldDiagnostic(match, field, "field is not present in the OpenAPI schema properties"))
			continue
		}
		if !prop.canCheckOptionality() {
			continue
		}

		required := match.schema.required[jsonName]
		switch {
		case required && !prop.nullable && field.isPointer:
			diagnostics = append(diagnostics, newSchemaFieldDiagnostic(match, field, "field is required and non-nullable in the OpenAPI schema but is a pointer"))
		case required && field.hasOmitOption:
			diagnostics = append(diagnostics, newSchemaFieldDiagnostic(match, field, "field is required by the OpenAPI schema but has an omit option"))
		case !required && !field.canBeOmitted:
			diagnostics = append(diagnostics, newSchemaFieldDiagnostic(match, field, "field is optional in the OpenAPI schema but is not a pointer, slice, map, interface, or selector type"))
		case !required && !field.hasOmitOption:
			diagnostics = append(diagnostics, newSchemaFieldDiagnostic(match, field, `field is optional in the OpenAPI schema but is missing "omitempty" or "omitzero"`))
		}
	}

	for propName, prop := range match.schema.properties {
		if !prop.canCheckOptionality() {
			continue
		}
		if _, ok := match.goStruct.fields[propName]; ok {
			continue
		}
		diagnostics = append(diagnostics, &schemaFieldDiagnostic{
			OpenAPISchema: match.schema.openapiSchema,
			GoStruct:      match.goStruct.name,
			JSONName:      propName,
			Field:         propName,
			Message:       "OpenAPI schema property is missing from the Go struct",
			OpenAPIFile:   match.schema.openapiFile,
		})
	}
	return diagnostics
}

func (p openapiSchemaProperty) canCheckOptionality() bool {
	return !p.readOnly && !p.writeOnly
}

func newSchemaFieldDiagnostic(match *schemaFieldMatch, field goStructField, message string) *schemaFieldDiagnostic {
	return &schemaFieldDiagnostic{
		OpenAPISchema: match.schema.openapiSchema,
		GoStruct:      match.goStruct.name,
		Field:         field.field,
		JSONName:      field.jsonName,
		Message:       message,
		Filename:      field.filename,
		Line:          field.line,
		OpenAPIFile:   match.schema.openapiFile,
	}
}

// collectGoStructs parses the Go source files in dir and returns every exported struct by
// name along with the set of struct types used exclusively as request bodies. A request body
// is the type of the body argument passed to a mutating (POST, PUT, or PATCH) client.NewRequest
// call; types that are also returned as a response (for example shared types like Label) are
// excluded because they follow the all-pointer response convention rather than the request-body
// convention.
func collectGoStructs(dir string) (map[string]*goStructInfo, map[string]bool, error) {
	structs := map[string]*goStructInfo{}
	requestStructs := map[string]bool{}
	responseStructs := map[string]bool{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(filename string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(filename, ".go") || strings.HasSuffix(filename, "_test.go") {
			return nil
		}

		fileNode, err := parser.ParseFile(fset, filename, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range fileNode.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				collectRequestStructNames(fn, requestStructs)
				collectResponseStructNames(fn, responseStructs)
				continue
			}
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !typeSpec.Name.IsExported() {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				structs[typeSpec.Name.Name] = &goStructInfo{
					name:     typeSpec.Name.Name,
					filename: filename,
					line:     fset.Position(typeSpec.Name.Pos()).Line,
					fields:   collectFieldsForStruct(fset, filename, typeSpec.Name.Name, structType),
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	// Exclude shared types that are also returned as a response; they follow the
	// all-pointer response convention rather than the request-body convention.
	for name := range responseStructs {
		delete(requestStructs, name)
	}
	return structs, requestStructs, nil
}

// collectRequestStructNames adds to requestStructs the name of the struct type passed as the
// body argument of every mutating client.NewRequest call in fn.
func collectRequestStructNames(fn *ast.FuncDecl, requestStructs map[string]bool) {
	if fn.Body == nil {
		return
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if isClientNewRequest(call) && isMutatingNewRequest(call) {
			if name := requestBodyStructName(fn, call.Args[3]); name != "" {
				requestStructs[name] = true
			}
		}
		return true
	})
}

// collectResponseStructNames adds to responseStructs the name of every struct type returned
// as a pointer (*T) or pointer slice ([]*T) by fn, which marks it as a response type.
func collectResponseStructNames(fn *ast.FuncDecl, responseStructs map[string]bool) {
	if fn.Type.Results == nil {
		return
	}
	for _, field := range fn.Type.Results.List {
		if name := responseStructName(field.Type); name != "" {
			responseStructs[name] = true
		}
	}
}

// responseStructName returns the struct name of a *T or []*T result type, or "" otherwise.
func responseStructName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.ArrayType:
		return responseStructName(t.Elt)
	}
	return ""
}

// isClientNewRequest reports whether call is of the form x.client.NewRequest(...) or client.NewRequest(...).
func isClientNewRequest(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "NewRequest" {
		return false
	}
	switch x := sel.X.(type) {
	case *ast.SelectorExpr:
		return x.Sel.Name == "client"
	case *ast.Ident:
		return x.Name == "client"
	default:
		return false
	}
}

// isMutatingNewRequest reports whether call's method argument is "PATCH", "POST", or "PUT" and a body argument is present.
func isMutatingNewRequest(call *ast.CallExpr) bool {
	if len(call.Args) < 4 {
		return false
	}
	lit, ok := call.Args[1].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	switch lit.Value {
	case `"PATCH"`, `"POST"`, `"PUT"`:
		return true
	default:
		return false
	}
}

// requestBodyStructName returns the struct type name of a client.NewRequest body argument,
// resolving a function parameter to its declared type or a composite literal to its type.
func requestBodyStructName(fn *ast.FuncDecl, arg ast.Expr) string {
	switch a := arg.(type) {
	case *ast.Ident:
		if field := findFuncParam(fn, a.Name); field != nil {
			return exprTypeName(field.Type)
		}
	case *ast.UnaryExpr:
		if a.Op == token.AND {
			return requestBodyStructName(fn, a.X)
		}
	case *ast.CompositeLit:
		return exprTypeName(a.Type)
	}
	return ""
}

func findFuncParam(fn *ast.FuncDecl, name string) *ast.Field {
	if fn.Type.Params == nil {
		return nil
	}
	for _, field := range fn.Type.Params.List {
		for _, ident := range field.Names {
			if ident.Name == name {
				return field
			}
		}
	}
	return nil
}

// exprTypeName returns the base type name of expr, unwrapping a pointer and resolving a qualified (pkg.Type) selector.
func exprTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return exprTypeName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func collectFieldsForStruct(fset *token.FileSet, filename, structName string, structType *ast.StructType) map[string]goStructField {
	fields := map[string]goStructField{}
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		for _, name := range field.Names {
			if !name.IsExported() {
				continue
			}
			info := goStructField{
				goStruct:     structName,
				field:        name.Name,
				jsonName:     defaultJSONName(name.Name),
				isPointer:    isPointerType(field.Type),
				canBeOmitted: canBeOmitted(field.Type),
				filename:     filename,
				line:         fset.Position(name.Pos()).Line,
			}
			if field.Tag != nil {
				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
				jsonName, hasOmitOption, ignored := parseJSONTag(tag.Get("json"))
				if ignored {
					continue
				}
				if jsonName != "" {
					info.jsonName = jsonName
				}
				info.hasOmitOption = hasOmitOption
			}
			if info.jsonName == "" {
				continue
			}
			fields[info.jsonName] = info
		}
	}
	return fields
}

func parseJSONTag(tag string) (name string, hasOmitOption, ignored bool) {
	if tag == "" {
		return "", false, false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	if name == "-" {
		return "", false, true
	}
	for _, opt := range parts[1:] {
		if opt == "omitempty" || opt == "omitzero" {
			hasOmitOption = true
		}
	}
	return name, hasOmitOption, false
}

func defaultJSONName(name string) string {
	return name
}

func isPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

func canBeOmitted(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.StarExpr, *ast.ArrayType, *ast.MapType, *ast.InterfaceType, *ast.SelectorExpr:
		return true
	}
	if ident, ok := expr.(*ast.Ident); ok && ident.Name == "any" {
		return true
	}
	return false
}
