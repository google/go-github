// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"maps"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/google/go-github/v84/github"
	"go.yaml.in/yaml/v3"
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
	slices.SortFunc(ops, func(a, b *operation) int {
		leftVerb, leftURL := parseOpName(a.Name)
		rightVerb, rightURL := parseOpName(b.Name)
		return cmp.Or(cmp.Compare(leftURL, rightURL), cmp.Compare(leftVerb, rightVerb))
	})
}

// normalizeOpPath returns an endpoint with all templated path parameters replaced with *.
func normalizeOpPath(opPath string) string {
	if !strings.ContainsAny(opPath, "{%") {
		return opPath
	}
	segments := strings.Split(opPath, "/")
	for i, segment := range segments {
		if len(segment) == 0 {
			continue
		}
		if segment[0] == '{' || segment[0] == '%' {
			segments[i] = "*"
		}
	}
	return strings.Join(segments, "/")
}

func normalizedOpName(name string) string {
	verb, u := parseOpName(name)
	return strings.TrimSpace(verb + " " + normalizeOpPath(u))
}

// matches something like "GET /some/path".
var opNameRe = regexp.MustCompile(`(?i)(\S+)(?:\s+(\S.*))?`)

func parseOpName(id string) (verb, url string) {
	match := opNameRe.FindStringSubmatch(id)
	if match == nil {
		return "", ""
	}
	u := strings.TrimSpace(match[2])
	if !strings.HasPrefix(u, "/") {
		u = "/" + u
	}
	return strings.ToUpper(match[1]), u
}

type operationsFile struct {
	ManualOps   []*operation `yaml:"operations,omitempty"`
	OverrideOps []*operation `yaml:"operation_overrides,omitempty"`
	GitCommit   string       `yaml:"openapi_commit,omitempty"`
	OpenapiOps  []*operation `yaml:"openapi_operations,omitempty"`

	mu          sync.Mutex
	resolvedOps map[string]*operation
}

func (m *operationsFile) resolve() {
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

func (m *operationsFile) saveFile(filename string) (errOut error) {
	sortOperations(m.ManualOps)
	sortOperations(m.OverrideOps)
	sortOperations(m.OpenapiOps)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	defer func() {
		e := enc.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	return enc.Encode(m)
}

func (m *operationsFile) updateFromGithub(ctx context.Context, client *github.Client, ref string) error {
	commit, resp, err := client.Repositories.GetCommit(ctx, descriptionsOwnerName, descriptionsRepoName, ref, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.Status)
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

func loadOperationsFile(filename string) (*operationsFile, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var opsFile operationsFile
	err = yaml.Unmarshal(b, &opsFile)
	if err != nil {
		return nil, err
	}
	return &opsFile, nil
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

func unusedOps(opsFile *operationsFile, dir string) ([]*operation, error) {
	var usedOps map[string]bool
	err := visitServiceMethods(dir, false, func(_ string, fn *ast.FuncDecl, cmap ast.CommentMap) error {
		ops, err := methodOps(opsFile, cmap, fn)
		if err != nil {
			return err
		}
		for _, op := range ops {
			if usedOps == nil {
				usedOps = map[string]bool{}
			}
			usedOps[op.Name] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var result []*operation
	opsFile.resolve()
	for opName, op := range opsFile.resolvedOps {
		if !usedOps[opName] {
			result = append(result, op)
		}
	}
	sortOperations(result)
	return result, nil
}

func updateDocsVisitor(opsFile *operationsFile) nodeVisitor {
	return func(serviceMethod string, fn *ast.FuncDecl, cmap ast.CommentMap) error {
		linksMap := map[string]struct{}{}
		undocMap := map[string]bool{}

		ops, err := methodOps(opsFile, cmap, fn)
		if err != nil {
			return err
		}
		if len(ops) == 0 {
			return fmt.Errorf("no operations defined for %v", serviceMethod)
		}

		for _, op := range ops {
			if op.DocumentationURL == "" {
				undocMap[op.Name] = true
				continue
			}
			linksMap[op.DocumentationURL] = struct{}{}
		}
		undocumentedOps := slices.Sorted(maps.Keys(undocMap))

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

		origList := group.List
		group.List = nil
		for _, comment := range origList {
			if metaOpRe.MatchString(comment.Text) ||
				docLineRE.MatchString(comment.Text) ||
				undocRE.MatchString(comment.Text) {
				continue
			}
			group.List = append(group.List, comment)
		}

		// add an empty line before adding doc links
		group.List = append(group.List, &ast.Comment{Text: "//"})

		docLinks := slices.Sorted(maps.Keys(linksMap))
		for i, dl := range docLinks {
			group.List = append(
				group.List,
				&ast.Comment{
					Text: "// GitHub API docs: " + cleanURLPath(dl),
				},
			)
			if i < len(docLinks)-1 {
				// add empty line between doc links
				group.List = append(group.List, &ast.Comment{Text: "//"})
			}
		}
		_, methodName, _ := strings.Cut(serviceMethod, ".")
		for _, opName := range undocumentedOps {
			line := fmt.Sprintf("// Note: %v uses the undocumented GitHub API endpoint %q.", methodName, opName)
			group.List = append(group.List, &ast.Comment{Text: line})
		}
		for _, op := range ops {
			group.List = append(group.List, &ast.Comment{
				Text: fmt.Sprintf("//meta:operation %v", op.Name),
			})
		}
		group.List[0].Slash = fn.Pos() - 1
		for i := 1; i < len(group.List); i++ {
			group.List[i].Slash = token.NoPos
		}
		return nil
	}
}

// updateDocs updates the code comments in dir with doc urls from metadata.
func updateDocs(opsFile *operationsFile, dir string) error {
	return visitServiceMethods(dir, true, updateDocsVisitor(opsFile))
}

type nodeVisitor func(serviceMethod string, fn *ast.FuncDecl, cmap ast.CommentMap) error

// visitServiceMethods runs visit on the ast.Node of every service method in dir. When writeFiles is true it will
// save any changes back to the original file.
func visitServiceMethods(dir string, writeFiles bool, visit nodeVisitor) error {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		filename := filepath.Join(dir, dirEntry.Name())
		if dirEntry.IsDir() ||
			!strings.HasSuffix(filename, ".go") ||
			strings.HasSuffix(filename, "_test.go") {
			continue
		}
		err = errors.Join(err, visitFileMethods(writeFiles, filename, visit))
	}
	return err
}

func visitFileMethods(updateFile bool, filename string, visit nodeVisitor) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return err
	}
	cmap := ast.NewCommentMap(fset, fileNode, fileNode.Comments)

	ast.Inspect(fileNode, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		serviceMethod := nodeServiceMethod(fn)
		if serviceMethod == "" {
			return true
		}
		e := visit(serviceMethod, fn, cmap)
		err = errors.Join(err, e)
		return true
	})
	if err != nil {
		return err
	}
	if !updateFile {
		return nil
	}
	fileNode.Comments = cmap.Filter(fileNode).Comments()
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, fileNode)
	if err != nil {
		return err
	}
	updatedContent, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	if bytes.Equal(content, updatedContent) {
		return nil
	}
	return os.WriteFile(filename, updatedContent, 0o600)
}

var (
	metaOpRe  = regexp.MustCompile(`(?i)\s*//\s*meta:operation\s+(\S.+)`)
	undocRE   = regexp.MustCompile(`(?i)\s*//\s*Note:\s+\S.+ uses the undocumented GitHub API endpoint`)
	docLineRE = regexp.MustCompile(`(?i)\s*//\s*GitHub\s+API\s+docs:`)
)

// methodOps parses a method's comments for //meta:operation lines and returns the corresponding operations.
func methodOps(opsFile *operationsFile, cmap ast.CommentMap, fn *ast.FuncDecl) ([]*operation, error) {
	var ops []*operation
	var err error
	seen := map[string]bool{}
	for _, g := range cmap[fn] {
		for _, c := range g.List {
			match := metaOpRe.FindStringSubmatch(c.Text)
			if match == nil {
				continue
			}
			opName := strings.TrimSpace(match[1])
			opsFile.resolve()
			var found []*operation
			norm := normalizedOpName(opName)
			for n := range opsFile.resolvedOps {
				if normalizedOpName(n) == norm {
					found = append(found, opsFile.resolvedOps[n])
				}
			}
			switch len(found) {
			case 0:
				err = errors.Join(err, fmt.Errorf("could not find operation %q in openapi_operations.yaml", opName))
			case 1:
				name := found[0].Name
				if seen[name] {
					err = errors.Join(err, fmt.Errorf("duplicate operation: %v", name))
				}
				seen[name] = true
				ops = append(ops, found[0])
			default:
				var foundNames []string
				for _, op := range found {
					foundNames = append(foundNames, op.Name)
				}
				slices.Sort(foundNames)
				err = errors.Join(err, fmt.Errorf("ambiguous operation %q could match any of: %v", opName, foundNames))
			}
		}
	}
	sortOperations(ops)
	return ops, err
}

// cleanURLPath runs path.Clean on the url path. This is to remove the unsightly double slashes from some
// of the urls in github's openapi descriptions.
func cleanURLPath(docURL string) string {
	u, err := url.Parse(docURL)
	if err != nil {
		return docURL
	}
	u.Path = path.Clean(u.Path)
	return u.String()
}

// nodeServiceMethod returns the name of the service method represented by fn, or "" if fn is not a service method.
// Name is in the form of "Receiver.Function", for example "IssuesService.Create".
func nodeServiceMethod(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) != 1 {
		return ""
	}
	recv := fn.Recv.List[0]
	se, ok := recv.Type.(*ast.StarExpr)
	if !ok {
		return ""
	}
	id, ok := se.X.(*ast.Ident)
	if !ok {
		return ""
	}

	// We only want exported methods on exported types where the type name ends in "Service".
	if !id.IsExported() || !fn.Name.IsExported() || !strings.HasSuffix(id.Name, "Service") {
		return ""
	}

	// Skip generated Iterator methods.
	if strings.HasSuffix(fn.Name.Name, "Iter") {
		return ""
	}

	serviceMethod := id.Name + "." + fn.Name.Name
	if skipServiceMethod[serviceMethod] {
		return ""
	}

	return serviceMethod
}

// See: https://github.com/google/go-github/issues/3894
var skipServiceMethod = map[string]bool{
	"BillingService.GetOrganizationPackagesBilling": true,
	"BillingService.GetOrganizationStorageBilling":  true,
	"BillingService.GetPackagesBilling":             true,
	"BillingService.GetStorageBilling":              true,
}
