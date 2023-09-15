// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type OperationDesc struct {
	Method           string `yaml:"method,omitempty" json:"method,omitempty"`
	EndpointURL      string `yaml:"endpoint_url,omitempty" json:"endpoint_url,omitempty"`
	DocumentationURL string `yaml:"documentation_url,omitempty" json:"documentation_url,omitempty"`
	Summary          string `yaml:"summary,omitempty" json:"summary,omitempty"`
}

type Operation struct {
	OpenAPI      OperationDesc `yaml:"openapi,omitempty" json:"openapi,omitempty"`
	Override     OperationDesc `yaml:"override,omitempty" json:"override,omitempty"`
	OpenAPIFiles []string      `yaml:"openapi_files,omitempty" json:"openapi_files,omitempty"`
	GoMethods    []string      `yaml:"go_methods,omitempty" json:"go_methods,omitempty"`
}

type operationJSON struct {
	Method      string   `json:"method,omitempty"`
	EndpointURL string   `json:"endpoint_url,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	DocumentURL string   `json:"documentation_url,omitempty"`
	Plans       []string `json:"plans,omitempty"`
	GoMethods   []string `json:"go_methods,omitempty"`
}

func (o *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&operationJSON{
		Method:      o.Method(),
		EndpointURL: o.EndpointURL(),
		Summary:     o.Summary(),
		Plans:       o.Plans(),
		DocumentURL: o.DocumentationURL(),
		GoMethods:   o.GoMethods,
	})
}

func (o *Operation) Plans() []string {
	var plans []string
	if slices.ContainsFunc(o.OpenAPIFiles, func(s string) bool {
		return strings.HasSuffix(s, "api.github.com.json")
	}) {
		plans = append(plans, "public")
	}
	if slices.ContainsFunc(o.OpenAPIFiles, func(s string) bool {
		return strings.HasSuffix(s, "ghec.json")
	}) {
		plans = append(plans, "ghec")
	}
	if slices.ContainsFunc(o.OpenAPIFiles, func(s string) bool {
		return strings.Contains(s, "/ghes")
	}) {
		plans = append(plans, "ghes")
	}
	return plans
}

func (o *Operation) Method() string {
	if o.Override.Method != "" {
		return o.Override.Method
	}
	return o.OpenAPI.Method
}

func (o *Operation) EndpointURL() string {
	if o.Override.EndpointURL != "" {
		return o.Override.EndpointURL
	}
	return o.OpenAPI.EndpointURL
}

func (o *Operation) DocumentationURL() string {
	if o.Override.DocumentationURL != "" {
		return o.Override.DocumentationURL
	}
	return o.OpenAPI.DocumentationURL
}

func (o *Operation) Summary() string {
	if o.Override.Summary != "" {
		return o.Override.Summary
	}
	return o.OpenAPI.Summary
}

func (o *Operation) Less(other *Operation) bool {
	if o.EndpointURL() != other.EndpointURL() {
		return o.EndpointURL() < other.EndpointURL()
	}
	return o.Method() < other.Method()
}

// matchesOpenAPIDesc returns true if this is describing the same operation as desc
// based on endpoint and method.
func (o *Operation) matchesOpenAPIDesc(desc OperationDesc) bool {
	if o.Method() != desc.Method {
		return false
	}
	return normalizedURL(o.EndpointURL()) == normalizedURL(desc.EndpointURL)
}

var normalizedURLs = map[string]string{}
var normalizedURLsMu sync.Mutex

// normalizedURL returns an endpoint with all templated path parameters replaced with *.
func normalizedURL(u string) string {
	normalizedURLsMu.Lock()
	defer normalizedURLsMu.Unlock()
	n, ok := normalizedURLs[u]
	if ok {
		return n
	}
	parts := strings.Split(u, "/")
	for i, p := range parts {
		if len(p) > 0 && p[0] == '{' {
			parts[i] = "*"
		}
	}
	n = strings.Join(parts, "/")
	normalizedURLs[u] = n
	return n
}

type Metadata struct {
	UndocumentedMethods []string     `yaml:"undocumented_methods,omitempty"`
	Operations          []*Operation `yaml:"operations,omitempty"`
}

func LoadMetadataFile(filename string, opFile *Metadata) (errOut error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	return yaml.NewDecoder(f).Decode(opFile)
}

func (m *Metadata) SaveFile(filename string) (errOut error) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		e := f.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	return enc.Encode(m)
}

func (m *Metadata) addOperation(filename string, desc OperationDesc) {
	for _, op := range m.Operations {
		if !op.matchesOpenAPIDesc(desc) {
			continue
		}
		if len(op.OpenAPIFiles) == 0 {
			op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
			op.OpenAPI = desc
			return
		}
		// just append to files, but only add the first ghes file
		if !strings.Contains(filename, "/ghes") {
			op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
			return
		}
		for _, f := range op.OpenAPIFiles {
			if strings.Contains(f, "/ghes") {
				return
			}
		}
		op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
		return
	}
	m.Operations = append(m.Operations, &Operation{
		OpenAPIFiles: []string{filename},
		OpenAPI:      desc,
	})
}

func (m *Metadata) operationsForMethod(method string) []*Operation {
	var operations []*Operation
	for _, op := range m.Operations {
		if !slices.Contains(op.GoMethods, method) {
			continue
		}
		operations = append(operations, op)
	}
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Less(operations[j])
	})
	return operations
}

func (m *Metadata) UpdateFromGithub(ctx context.Context, client contentsClient, ref string) error {
	descs, err := getDescriptions(ctx, client, ref)
	if err != nil {
		return err
	}
	for _, op := range m.Operations {
		op.OpenAPIFiles = op.OpenAPIFiles[:0]
	}
	for _, desc := range descs {
		for p, pathItem := range desc.description.Paths {
			for method, op := range pathItem.Operations() {
				docURL := ""
				if op.ExternalDocs != nil {
					docURL = op.ExternalDocs.URL
				}
				m.addOperation(desc.filename, OperationDesc{
					Method:           method,
					EndpointURL:      p,
					DocumentationURL: docURL,
					Summary:          op.Summary,
				})
			}
		}
	}
	sort.Slice(m.Operations, func(i, j int) bool {
		return m.Operations[i].Less(m.Operations[j])
	})
	return nil
}

// UpdateDocLinks updates the code comments in dir with doc urls from metadata.
func UpdateDocLinks(meta *Metadata, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() ||
			!strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updatedContent, err := updateDocsLinksInFile(meta, content)
		if err != nil {
			return err
		}
		if bytes.Equal(content, updatedContent) {
			return nil
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		_, err = f.Write(updatedContent)
		if err != nil {
			return err
		}
		return f.Close()
	})
}

// updateDocsLinksInFile updates in the code comments in content with doc urls from metadata.
func updateDocsLinksInFile(metadata *Metadata, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// ignore files where package is not github
	if node.Name.Name != "github" {
		return content, nil
	}

	ast.Inspect(node, func(n ast.Node) bool {
		return updateDocsLinksForNode(metadata, n)
	})

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, node)
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

var (
	docLineRE   = regexp.MustCompile(`(?i)\s*(//\s*)?GitHub\s+API\s+docs:\s*(https?://\S+)`)
	emptyLineRE = regexp.MustCompile(`^\s*(//\s*)$`)
)

func updateDocsLinksForNode(metadata *Metadata, n ast.Node) bool {
	fn, ok := n.(*ast.FuncDecl)
	if !ok || !fn.Name.IsExported() {
		return true
	}
	methodName := fn.Name.Name

	// Get the method's receiver. It can be either an identifier or a pointer to an identifier.
	// This assumes all receivers are named and we don't have something like: `func (Client) Foo()`.
	receiverType := ""
	if fn.Recv != nil {
		switch x := fn.Recv.List[0].Type.(type) {
		case *ast.Ident:
			receiverType = x.Name
		case *ast.StarExpr:
			receiverType = x.X.(*ast.Ident).Name
		}
	}
	if !ast.IsExported(receiverType) {
		return true
	}

	linksMap := map[string]struct{}{}
	ops := metadata.operationsForMethod(strings.Join([]string{receiverType, methodName}, "."))
	for _, op := range ops {
		linksMap[op.DocumentationURL()] = struct{}{}
	}

	// create copy of comment group with non-matching doc links removed
	if fn.Doc == nil {
		fn.Doc = &ast.CommentGroup{}
	}
	fnComments := make([]*ast.Comment, 0, len(fn.Doc.List))
	for _, comment := range fn.Doc.List {
		match := docLineRE.FindStringSubmatch(comment.Text)
		if match == nil {
			fnComments = append(fnComments, comment)
			continue
		}
		matchesLink := false
		for link := range linksMap {
			if sameDocLink(match[2], link) {
				matchesLink = true
				delete(linksMap, link)
				break
			}
		}
		if matchesLink {
			fnComments = append(fnComments, comment)
		}
	}

	// remove trailing empty lines
	for len(fnComments) > 0 {
		if !emptyLineRE.MatchString(fnComments[len(fnComments)-1].Text) {
			break
		}
		fnComments = fnComments[:len(fnComments)-1]
	}

	// add an empty line before adding doc links
	if len(linksMap) > 0 &&
		len(fnComments) > 0 &&
		!emptyLineRE.MatchString(fnComments[len(fnComments)-1].Text) {
		fnComments = append(fnComments, &ast.Comment{Text: "//"})
	}

	docLinks := maps.Keys(linksMap)
	sort.Strings(docLinks)

	for _, dl := range docLinks {
		fnComments = append(
			fnComments,
			&ast.Comment{
				Text: "// GitHub API docs: " + normalizeDocURLPath(dl),
			},
		)
	}
	fn.Doc.List = fnComments
	return true
}

const docURLPrefix = "https://docs.github.com/rest/"

var docURLPrefixRE = regexp.MustCompile(`^https://docs\.github\.com.*/rest/`)

func normalizeDocURLPath(u string) string {
	u = strings.Replace(u, "/en/", "/", 1)
	pre := docURLPrefixRE.FindString(u)
	if pre == "" {
		return u
	}
	if strings.Contains(u, "docs.github.com/enterprise-cloud@latest/") {
		// remove unsightly double slash
		// https://docs.github.com/enterprise-cloud@latest/
		return strings.ReplaceAll(
			u,
			"docs.github.com/enterprise-cloud@latest//",
			"docs.github.com/enterprise-cloud@latest/",
		)
	}
	if strings.Contains(u, "docs.github.com/enterprise-server") {
		return u
	}
	return docURLPrefix + strings.TrimPrefix(u, pre)
}

// sameDocLink returns true if the two doc links are going to end up rendering the same page pointed
// to the same section.
//
// If a url path starts with *./rest/ it ignores query parameters and everything before /rest/ when
// making the comparison.
func sameDocLink(left, right string) bool {
	if !docURLPrefixRE.MatchString(left) ||
		!docURLPrefixRE.MatchString(right) {
		return left == right
	}
	left = stripURLQuery(normalizeDocURLPath(left))
	right = stripURLQuery(normalizeDocURLPath(right))
	return left == right
}

func stripURLQuery(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return u
	}
	p.RawQuery = ""
	return p.String()
}
