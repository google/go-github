// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"fmt"
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

	"github.com/google/go-github/v56/github"
	"gopkg.in/yaml.v3"
)

type operation struct {
	Name             string   `yaml:"name,omitempty" json:"name,omitempty"`
	DocumentationURL string   `yaml:"documentation_url,omitempty" json:"documentation_url,omitempty"`
	OpenAPIFiles     []string `yaml:"openapi_files,omitempty" json:"openapi_files,omitempty"`
}

func (o *operation) equal(other *operation) bool {
	if o.Name != other.Name || o.DocumentationURL != other.DocumentationURL {
		return false
	}
	if len(o.OpenAPIFiles) != len(other.OpenAPIFiles) {
		return false
	}
	for i := range o.OpenAPIFiles {
		if o.OpenAPIFiles[i] != other.OpenAPIFiles[i] {
			return false
		}
	}
	return true
}

func (o *operation) clone() *operation {
	return &operation{
		Name:             o.Name,
		DocumentationURL: o.DocumentationURL,
		OpenAPIFiles:     append([]string{}, o.OpenAPIFiles...),
	}
}

func operationsEqual(a, b []*operation) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].equal(b[i]) {
			return false
		}
	}
	return true
}

func sortOperations(ops []*operation) {
	sort.Slice(ops, func(i, j int) bool {
		leftVerb, leftURL := parseOpName(ops[i].Name)
		rightVerb, rightURL := parseOpName(ops[j].Name)
		if leftURL != rightURL {
			return leftURL < rightURL
		}
		return leftVerb < rightVerb
	})
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

func normalizedOpName(name string) string {
	verb, u := parseOpName(name)
	return verb + " " + normalizedURL(u)
}

func parseOpName(id string) (verb, url string) {
	verb, url, _ = strings.Cut(id, " ")
	return verb, url
}

type method struct {
	Name    string   `yaml:"name" json:"name"`
	OpNames []string `yaml:"operations,omitempty" json:"operations,omitempty"`
}

type metadata struct {
	Methods     []*method    `yaml:"methods,omitempty"`
	ManualOps   []*operation `yaml:"operations,omitempty"`
	OverrideOps []*operation `yaml:"operation_overrides,omitempty"`
	GitCommit   string       `yaml:"openapi_commit,omitempty"`
	OpenapiOps  []*operation `yaml:"openapi_operations,omitempty"`

	mu          sync.Mutex
	resolvedOps map[string]*operation
}

func (m *metadata) resolve() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.resolvedOps != nil {
		return
	}
	m.resolvedOps = map[string]*operation{}
	for _, op := range m.OpenapiOps {
		m.resolvedOps[op.Name] = op.clone()
	}
	for _, op := range m.ManualOps {
		m.resolvedOps[op.Name] = op.clone()
	}
	for _, override := range m.OverrideOps {
		_, ok := m.resolvedOps[override.Name]
		if !ok {
			continue
		}
		override = override.clone()
		if override.DocumentationURL != "" {
			m.resolvedOps[override.Name].DocumentationURL = override.DocumentationURL
		}
		if len(override.OpenAPIFiles) > 0 {
			m.resolvedOps[override.Name].OpenAPIFiles = override.OpenAPIFiles
		}
	}
}

func (m *metadata) operations() []*operation {
	m.resolve()
	ops := make([]*operation, 0, len(m.resolvedOps))
	for _, op := range m.resolvedOps {
		ops = append(ops, op)
	}
	sortOperations(ops)
	return ops
}

func loadMetadataFile(filename string) (*metadata, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var meta metadata
	err = yaml.Unmarshal(b, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (m *metadata) saveFile(filename string) (errOut error) {
	sortOperations(m.ManualOps)
	sortOperations(m.OverrideOps)
	sortOperations(m.OpenapiOps)
	sort.Slice(m.Methods, func(i, j int) bool {
		return m.Methods[i].Name < m.Methods[j].Name
	})
	for _, method := range m.Methods {
		sort.Strings(method.OpNames)
	}
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

func addOperation(ops []*operation, filename, opName, docURL string) []*operation {
	for _, op := range ops {
		if opName != op.Name {
			continue
		}
		if len(op.OpenAPIFiles) == 0 {
			op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
			op.DocumentationURL = docURL
			return ops
		}
		// just append to files, but only add the first ghes file
		if !strings.Contains(filename, "/ghes") {
			op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
			return ops
		}
		for _, f := range op.OpenAPIFiles {
			if strings.Contains(f, "/ghes") {
				return ops
			}
		}
		op.OpenAPIFiles = append(op.OpenAPIFiles, filename)
		return ops
	}
	return append(ops, &operation{
		Name:             opName,
		OpenAPIFiles:     []string{filename},
		DocumentationURL: docURL,
	})
}

// operationMethods returns a list methods that are mapped to the given operation id.
func (m *metadata) operationMethods(opName string) []string {
	var methods []string
	for _, method := range m.Methods {
		for _, methodOpName := range method.OpNames {
			if methodOpName == opName {
				methods = append(methods, method.Name)
			}
		}
	}
	return methods
}

func (m *metadata) getOperation(name string) *operation {
	m.resolve()
	return m.resolvedOps[name]
}

func (m *metadata) getOperationsWithNormalizedName(name string) []*operation {
	m.resolve()
	var result []*operation
	norm := normalizedOpName(name)
	for n := range m.resolvedOps {
		if normalizedOpName(n) == norm {
			result = append(result, m.resolvedOps[n])
		}
	}
	sortOperations(result)
	return result
}

func (m *metadata) getMethod(name string) *method {
	for _, method := range m.Methods {
		if method.Name == name {
			return method
		}
	}
	return nil
}

func (m *metadata) operationsForMethod(methodName string) []*operation {
	method := m.getMethod(methodName)
	if method == nil {
		return nil
	}
	var operations []*operation
	for _, name := range method.OpNames {
		op := m.getOperation(name)
		if op != nil {
			operations = append(operations, op)
		}
	}
	sortOperations(operations)
	return operations
}

func (m *metadata) canonizeMethodOperations() error {
	for _, method := range m.Methods {
		for i := range method.OpNames {
			opName := method.OpNames[i]
			if m.getOperation(opName) != nil {
				continue
			}
			ops := m.getOperationsWithNormalizedName(opName)
			switch len(ops) {
			case 0:
				return fmt.Errorf("method %q has an operation that can not be canonized to any defined name: %s", method.Name, opName)
			case 1:
				method.OpNames[i] = ops[0].Name
			default:
				candidateList := ""
				for _, op := range ops {
					candidateList += "\n    " + op.Name
				}
				return fmt.Errorf("method %q has an operation that can be canonized to multiple defined names:\n  operation: %s\n  matches: %s", method.Name, opName, candidateList)
			}
		}
	}
	return nil
}

func (m *metadata) updateFromGithub(ctx context.Context, client *github.Client, ref string) error {
	commit, resp, err := client.Repositories.GetCommit(ctx, descriptionsOwnerName, descriptionsRepoName, ref, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	ops, err := getOpsFromGithub(ctx, client, ref)
	if err != nil {
		return err
	}
	if !operationsEqual(m.OpenapiOps, ops) {
		m.OpenapiOps = ops
		m.GitCommit = commit.GetSHA()
	}
	return nil
}

// updateDocLinks updates the code comments in dir with doc urls from metadata.
func updateDocLinks(meta *metadata, dir string) error {
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
		content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
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
func updateDocsLinksInFile(metadata *metadata, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	cmap := ast.NewCommentMap(fset, file, file.Comments)

	// ignore files where package is not github
	if file.Name.Name != "github" {
		return content, nil
	}

	ast.Inspect(file, func(n ast.Node) bool {
		return updateDocsLinksForNode(metadata, cmap, n)
	})

	file.Comments = cmap.Filter(file).Comments()

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

var (
	docLineRE = regexp.MustCompile(`(?i)\s*(//\s*)?GitHub\s+API\s+docs:\s*(https?://\S+)`)
)

func updateDocsLinksForNode(metadata *metadata, cmap ast.CommentMap, n ast.Node) bool {
	fn, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}
	sm := serviceMethodFromNode(n)
	if sm == "" {
		return true
	}

	linksMap := map[string]struct{}{}
	undocMap := map[string]bool{}
	ops := metadata.operationsForMethod(sm)
	for _, op := range ops {
		if op.DocumentationURL == "" {
			undocMap[op.Name] = true
			continue
		}
		linksMap[op.DocumentationURL] = struct{}{}
	}
	var undocumentedOps []string
	for op := range undocMap {
		undocumentedOps = append(undocumentedOps, op)
	}
	sort.Strings(undocumentedOps)

	// Find the group that comes before the function
	var group *ast.CommentGroup
	for _, g := range cmap[fn] {
		if g.End() == fn.Pos()-1 {
			group = g
		}
	}

	// If there is no group, create one
	if group == nil {
		group = &ast.CommentGroup{
			List: []*ast.Comment{{Text: "//", Slash: fn.Pos() - 1}},
		}
		cmap[fn] = append(cmap[fn], group)
	}

	skipSpacer := false
	var newList []*ast.Comment
	for _, comment := range group.List {
		if strings.Contains(comment.Text, "uses the undocumented GitHub API endpoint") {
			skipSpacer = true
			continue
		}
		match := docLineRE.FindStringSubmatch(comment.Text)
		if match == nil {
			newList = append(newList, comment)
			continue
		}
		matchesLink := false
		for link := range linksMap {
			if sameDocLink(match[2], link) {
				matchesLink = true
				skipSpacer = true
				delete(linksMap, link)
				break
			}
		}
		if matchesLink {
			newList = append(newList, comment)
		}
	}
	group.List = newList

	// add an empty line before adding doc links
	if !skipSpacer {
		group.List = append(group.List, &ast.Comment{Text: "//"})
	}

	var docLinks []string
	for link := range linksMap {
		docLinks = append(docLinks, link)
	}
	sort.Strings(docLinks)

	for _, dl := range docLinks {
		group.List = append(
			group.List,
			&ast.Comment{
				Text: "// GitHub API docs: " + normalizeDocURLPath(dl),
			},
		)
	}
	_, methodName, _ := strings.Cut(sm, ".")
	for _, opName := range undocumentedOps {
		line := fmt.Sprintf("// Note: %s uses the undocumented GitHub API endpoint %q.", methodName, opName)
		group.List = append(group.List, &ast.Comment{Text: line})
	}
	group.List[0].Slash = fn.Pos() - 1
	for i := 1; i < len(group.List); i++ {
		group.List[i].Slash = token.NoPos
	}
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

func serviceMethodFromNode(node ast.Node) string {
	decl, ok := node.(*ast.FuncDecl)
	if !ok || decl.Recv == nil || len(decl.Recv.List) != 1 {
		return ""
	}
	recv := decl.Recv.List[0]
	se, ok := recv.Type.(*ast.StarExpr)
	if !ok {
		return ""
	}
	id, ok := se.X.(*ast.Ident)
	if !ok {
		return ""
	}

	// We only want exported methods on exported types where the type name ends in "Service".
	if !id.IsExported() || !decl.Name.IsExported() || !strings.HasSuffix(id.Name, "Service") {
		return ""
	}

	return id.Name + "." + decl.Name.Name
}
