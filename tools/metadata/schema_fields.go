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
	"regexp"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.yaml.in/yaml/v3"
)

type schemaFieldCheckOptions struct {
	descriptions []*openapiFile
	githubDir    string
	// exceptions holds "Struct.Field" entries whose diagnostics are suppressed. It is loaded from
	// schema_field_exceptions.yaml by the command; see loadSchemaFieldExceptions.
	exceptions []string
}

type schemaFieldCheckResult struct {
	Summary     schemaFieldCheckSummary
	Checked     []*schemaFieldChecked
	Diagnostics []*schemaFieldDiagnostic
}

type schemaFieldCheckSummary struct {
	GoStructs        int
	AnnotatedStructs int
	Checked          int
	Diagnostics      int
}

type schemaFieldChecked struct {
	Annotation  string
	GoStruct    string
	OpenAPIFile string
}

type schemaFieldDiagnostic struct {
	Annotation  string
	GoStruct    string
	Field       string
	JSONName    string
	Message     string
	Filename    string
	Line        int
	OpenAPIFile string
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
	subject := d.GoStruct
	if d.Field != "" {
		subject += "." + d.Field
	}
	if d.JSONName != "" && d.Annotation != "" {
		subject += fmt.Sprintf(" (%v from %v)", d.JSONName, d.Annotation)
	}
	return fmt.Sprintf("%v%v: %v%v", loc, subject, d.Message, source)
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

// schemaAnnotation is one "//meta:schema <role> <METHOD> <path>" line from a struct doc comment. It names
// the operation whose request or response body schema the annotated struct must match.
type schemaAnnotation struct {
	role     string // "request" or "response"
	method   string
	path     string
	filename string
	line     int
}

func (a schemaAnnotation) String() string {
	return fmt.Sprintf("%v %v %v", a.role, a.method, a.path)
}

// schemaAnnotationProblem is a malformed "//meta:schema" line that could not be parsed into a
// schemaAnnotation.
type schemaAnnotationProblem struct {
	text     string
	message  string
	filename string
	line     int
}

type goStructInfo struct {
	name               string
	filename           string
	line               int
	fields             map[string]goStructField
	annotations        []*schemaAnnotation
	annotationProblems []*schemaAnnotationProblem
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
	annotation  string
	openapiFile string
	required    []string
	properties  map[string]openapiSchemaProperty
}

type openapiSchemaProperty struct {
	nullable  bool
	readOnly  bool
	writeOnly bool
}

type schemaFieldMatch struct {
	schema   *openapiSchemaFields
	goStruct *goStructInfo
}

// schemaFieldExceptionsFile is the on-disk format of schema_field_exceptions.yaml: a list of "Struct.Field"
// entries whose JSON field optionality intentionally deviates from the OpenAPI schema, so their diagnostics
// are suppressed. Each entry is a known deviation awaiting cleanup.
type schemaFieldExceptionsFile struct {
	Exceptions []string `yaml:"exceptions"`
}

// loadSchemaFieldExceptions reads the "Struct.Field" exception entries from the exceptions file under
// workingDir and returns them. A missing file yields no exceptions and no error, so callers do not need
// the file to exist.
func loadSchemaFieldExceptions(workingDir string) ([]string, error) {
	filename := filepath.Join(workingDir, "tools/metadata/schema_field_exceptions.yaml")
	b, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var exceptionsFile schemaFieldExceptionsFile
	if err := yaml.Unmarshal(b, &exceptionsFile); err != nil {
		return nil, fmt.Errorf("%v: %w", filename, err)
	}
	return exceptionsFile.Exceptions, nil
}

// filterAllowedSchemaFieldDiagnostics removes diagnostics whose "Struct.Field" is listed in the exceptions.
func filterAllowedSchemaFieldDiagnostics(diagnostics []*schemaFieldDiagnostic, exceptions []string) []*schemaFieldDiagnostic {
	var kept []*schemaFieldDiagnostic
	for _, diag := range diagnostics {
		if slices.Contains(exceptions, diag.GoStruct+"."+diag.Field) {
			continue
		}
		kept = append(kept, diag)
	}
	return kept
}

// checkSchemaFields validates every Go struct that carries at least one "//meta:schema" annotation against
// the OpenAPI schema of the annotated operation. Structs without annotations are not checked.
func checkSchemaFields(opts schemaFieldCheckOptions) (schemaFieldCheckResult, error) {
	if len(opts.descriptions) == 0 {
		return schemaFieldCheckResult{}, errors.New("no OpenAPI descriptions loaded")
	}

	goStructs, err := collectGoStructs(opts.githubDir)
	if err != nil {
		return schemaFieldCheckResult{}, err
	}

	var result schemaFieldCheckResult
	for _, name := range slices.Sorted(maps.Keys(goStructs)) {
		goStruct := goStructs[name]
		if len(goStruct.annotations) == 0 && len(goStruct.annotationProblems) == 0 {
			continue
		}
		result.Summary.AnnotatedStructs++

		for _, problem := range goStruct.annotationProblems {
			result.Diagnostics = append(result.Diagnostics, &schemaFieldDiagnostic{
				GoStruct: goStruct.name,
				Message:  fmt.Sprintf("invalid annotation %q: %v", problem.text, problem.message),
				Filename: problem.filename,
				Line:     problem.line,
			})
		}

		for _, ann := range goStruct.annotations {
			schema, openapiFilename, problem := resolveSchemaAnnotation(opts.descriptions, ann)
			if problem != "" {
				result.Diagnostics = append(result.Diagnostics, &schemaFieldDiagnostic{
					Annotation:  ann.String(),
					GoStruct:    goStruct.name,
					Message:     problem,
					Filename:    ann.filename,
					Line:        ann.line,
					OpenAPIFile: openapiFilename,
				})
				continue
			}

			flat, reason, err := flattenObjectSchema(schema)
			if err != nil {
				return schemaFieldCheckResult{}, fmt.Errorf("%v %v: %w", goStruct.name, ann, err)
			}
			if reason != "" {
				result.Diagnostics = append(result.Diagnostics, &schemaFieldDiagnostic{
					Annotation:  ann.String(),
					GoStruct:    goStruct.name,
					Message:     "annotated schema cannot be checked: " + reason,
					Filename:    ann.filename,
					Line:        ann.line,
					OpenAPIFile: openapiFilename,
				})
				continue
			}

			match := &schemaFieldMatch{
				schema: &openapiSchemaFields{
					annotation:  ann.String(),
					openapiFile: openapiFilename,
					required:    flat.Required,
					properties:  schemaProperties(flat.Properties),
				},
				goStruct: goStruct,
			}
			result.Checked = append(result.Checked, &schemaFieldChecked{
				Annotation:  ann.String(),
				GoStruct:    goStruct.name,
				OpenAPIFile: openapiFilename,
			})
			result.Diagnostics = append(result.Diagnostics, compareSchemaFields(match)...)
		}
	}

	result.Diagnostics = filterAllowedSchemaFieldDiagnostics(result.Diagnostics, opts.exceptions)
	sortSchemaFieldResult(&result)
	result.Summary.GoStructs = len(goStructs)
	result.Summary.Checked = len(result.Checked)
	result.Summary.Diagnostics = len(result.Diagnostics)
	return result, nil
}

func sortSchemaFieldResult(result *schemaFieldCheckResult) {
	slices.SortFunc(result.Diagnostics, func(a, b *schemaFieldDiagnostic) int {
		return cmp.Or(
			cmp.Compare(a.GoStruct, b.GoStruct),
			cmp.Compare(a.JSONName, b.JSONName),
			cmp.Compare(a.Field, b.Field),
			cmp.Compare(a.Annotation, b.Annotation),
			cmp.Compare(a.OpenAPIFile, b.OpenAPIFile),
			cmp.Compare(a.Message, b.Message),
		)
	})
	slices.SortFunc(result.Checked, func(a, b *schemaFieldChecked) int {
		return cmp.Or(
			cmp.Compare(a.GoStruct, b.GoStruct),
			cmp.Compare(a.Annotation, b.Annotation),
			cmp.Compare(a.OpenAPIFile, b.OpenAPIFile),
		)
	})
}

// resolveSchemaAnnotation finds the schema named by ann in the first OpenAPI description that documents the
// annotated operation, searching the descriptions in their load order (api.github.com first, then ghec,
// then ghes). It returns a non-empty problem string when the operation cannot be found or has no matching
// JSON schema, mirroring how unknown "//meta:operation" names are reported.
func resolveSchemaAnnotation(descriptions []*openapiFile, ann *schemaAnnotation) (schema *openapi3.Schema, openapiFilename, problem string) {
	normPath := normalizeOpPath(ann.path)
	for _, desc := range descriptions {
		if desc.description == nil {
			continue
		}
		for path, pathItem := range desc.description.Paths.Map() {
			if pathItem == nil || normalizeOpPath(path) != normPath {
				continue
			}
			op := pathItem.Operations()[ann.method]
			if op == nil {
				continue
			}
			schema, problem = annotationSchema(op, ann.role)
			return schema, desc.filename, problem
		}
	}
	return nil, "", fmt.Sprintf("could not find operation %v %v in any OpenAPI description", ann.method, ann.path)
}

// annotationSchema extracts the request or response JSON schema from op according to role.
func annotationSchema(op *openapi3.Operation, role string) (*openapi3.Schema, string) {
	switch role {
	case "request":
		if op.RequestBody == nil || op.RequestBody.Value == nil {
			return nil, "operation has no request body"
		}
		return jsonContentSchema(op.RequestBody.Value.Content, "operation request body")
	default: // "response"; parseSchemaAnnotations rejects other roles
		if op.Responses == nil {
			return nil, "operation has no responses"
		}
		responses := op.Responses.Map()
		for _, code := range slices.Sorted(maps.Keys(responses)) {
			if !strings.HasPrefix(code, "2") || responses[code] == nil || responses[code].Value == nil {
				continue
			}
			if schema, problem := jsonContentSchema(responses[code].Value.Content, ""); problem == "" {
				return schema, ""
			}
		}
		return nil, "operation has no 2xx response with an application/json schema"
	}
}

// jsonContentSchema returns the application/json schema from content, or a problem string naming what is
// missing.
func jsonContentSchema(content openapi3.Content, what string) (*openapi3.Schema, string) {
	mediaType := content.Get("application/json")
	if mediaType == nil || mediaType.Schema == nil || mediaType.Schema.Value == nil {
		return nil, what + " has no application/json schema"
	}
	return mediaType.Schema.Value, ""
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
		Required:   slices.Clone(schema.Required),
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

		required := slices.Contains(match.schema.required, jsonName)
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
			Annotation:  match.schema.annotation,
			GoStruct:    match.goStruct.name,
			JSONName:    propName,
			Field:       propName,
			Message:     "OpenAPI schema property is missing from the Go struct",
			OpenAPIFile: match.schema.openapiFile,
		})
	}
	return diagnostics
}

func (p openapiSchemaProperty) canCheckOptionality() bool {
	return !p.readOnly && !p.writeOnly
}

func newSchemaFieldDiagnostic(match *schemaFieldMatch, field goStructField, message string) *schemaFieldDiagnostic {
	return &schemaFieldDiagnostic{
		Annotation:  match.schema.annotation,
		GoStruct:    match.goStruct.name,
		Field:       field.field,
		JSONName:    field.jsonName,
		Message:     message,
		Filename:    field.filename,
		Line:        field.line,
		OpenAPIFile: match.schema.openapiFile,
	}
}

// metaSchemaLineRe recognizes a "//meta:schema ..." doc comment line; the arguments are validated by
// parseSchemaAnnotations.
var metaSchemaLineRe = regexp.MustCompile(`(?i)^\s*//\s*meta:schema\b(.*)$`)

// parseSchemaAnnotations extracts every "//meta:schema <role> <METHOD> <path>" line from doc. Malformed
// lines are returned as problems so they surface as diagnostics instead of being silently ignored.
func parseSchemaAnnotations(fset *token.FileSet, doc *ast.CommentGroup, filename string) ([]*schemaAnnotation, []*schemaAnnotationProblem) {
	if doc == nil {
		return nil, nil
	}
	var annotations []*schemaAnnotation
	var problems []*schemaAnnotationProblem
	for _, comment := range doc.List {
		m := metaSchemaLineRe.FindStringSubmatch(comment.Text)
		if m == nil {
			continue
		}
		line := fset.Position(comment.Pos()).Line
		args := strings.Fields(m[1])
		if len(args) != 3 {
			problems = append(problems, &schemaAnnotationProblem{
				text:     strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")),
				message:  "want //meta:schema <request|response> <METHOD> <path>",
				filename: filename,
				line:     line,
			})
			continue
		}
		role := strings.ToLower(args[0])
		if role != "request" && role != "response" {
			problems = append(problems, &schemaAnnotationProblem{
				text:     strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")),
				message:  fmt.Sprintf("unknown role %q; want request or response", args[0]),
				filename: filename,
				line:     line,
			})
			continue
		}
		annotations = append(annotations, &schemaAnnotation{
			role:     role,
			method:   strings.ToUpper(args[1]),
			path:     args[2],
			filename: filename,
			line:     line,
		})
	}
	return annotations, problems
}

// collectGoStructs parses the Go source files in dir and returns every exported struct by name, along with
// any "//meta:schema" annotations found in the struct doc comments.
func collectGoStructs(dir string) (map[string]*goStructInfo, error) {
	structs := map[string]*goStructInfo{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(filename string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(filename, ".go") || strings.HasSuffix(filename, "_test.go") {
			return nil
		}

		fileNode, err := parser.ParseFile(fset, filename, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range fileNode.Decls {
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
				doc := typeSpec.Doc
				if doc == nil && len(gen.Specs) == 1 {
					doc = gen.Doc
				}
				annotations, problems := parseSchemaAnnotations(fset, doc, filename)
				structs[typeSpec.Name.Name] = &goStructInfo{
					name:               typeSpec.Name.Name,
					filename:           filename,
					line:               fset.Position(typeSpec.Name.Pos()).Line,
					fields:             collectFieldsForStruct(fset, filename, typeSpec.Name.Name, structType),
					annotations:        annotations,
					annotationProblems: problems,
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return structs, nil
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
				jsonName:     name.Name,
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
