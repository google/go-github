// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// update-urls updates GitHub URL docs for each service endpoint.
//
// It is meant to be used periodically by go-github repo maintainers
// to update stale GitHub Developer v3 API documentation URLs.
//
// Usage (from go-github directory):
//   go run ./update-urls/main.go
//   go generate ./...
//   go test ./...
//   go vet ./...
//
// When confronted with "PLEASE CHECK MANUALLY AND FIX", the problematic
// URL needs to be debugged. To debug a specific file, run like this:
//   go run ./update-urls/main.go -v -d enterprise_actions_runners.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	codeLegacySplitString = `<code>`
	codeSplitString       = `<code class="hljs language-javascript"><span class="hljs-keyword">await</span> octokit.request(<span class="hljs-string">&apos;`
	fragmentIDString      = `<h3 id="`
	skipPrefix            = "gen-"

	// enterpriseURL = "docs.github.com"
	stdURL = "docs.github.com"

	enterpriseRefFmt = "// GitHub Enterprise API docs: %v"
	stdRefFmt        = "// GitHub API docs: %v"
)

var (
	verbose   = flag.Bool("v", false, "Print verbose log messages")
	debugFile = flag.String("d", "", "Debug named file only")

	// skipMethods holds methods which are skipped because they do not have GitHub v3
	// API URLs or are otherwise problematic in parsing, discovering, and/or fixing.
	skipMethods = map[string]bool{
		"ActionsService.DownloadArtifact":              true,
		"AdminService.CreateOrg":                       true,
		"AdminService.CreateUser":                      true,
		"AdminService.CreateUserImpersonation":         true,
		"AdminService.DeleteUserImpersonation":         true,
		"AdminService.GetAdminStats":                   true,
		"AdminService.RenameOrg":                       true,
		"AdminService.RenameOrgByName":                 true,
		"AdminService.UpdateTeamLDAPMapping":           true,
		"AdminService.UpdateUserLDAPMapping":           true,
		"AppsService.FindRepositoryInstallationByID":   true,
		"AuthorizationsService.CreateImpersonation":    true,
		"AuthorizationsService.DeleteImpersonation":    true,
		"IssueImportService.CheckStatus":               true,
		"IssueImportService.CheckStatusSince":          true,
		"IssueImportService.Create":                    true,
		"MarketplaceService.marketplaceURI":            true,
		"OrganizationsService.GetByID":                 true,
		"RepositoriesService.DeletePreReceiveHook":     true,
		"RepositoriesService.DownloadContents":         true,
		"RepositoriesService.DownloadContentsWithMeta": true,
		"RepositoriesService.GetArchiveLink":           true,
		"RepositoriesService.GetByID":                  true,
		"RepositoriesService.GetPreReceiveHook":        true,
		"RepositoriesService.ListPreReceiveHooks":      true,
		"RepositoriesService.UpdatePreReceiveHook":     true,
		"SearchService.search":                         true,
		"TeamsService.ListTeamMembersByID":             true,
		"UsersService.DemoteSiteAdmin":                 true,
		"UsersService.GetByID":                         true,
		"UsersService.PromoteSiteAdmin":                true,
		"UsersService.Suspend":                         true,
		"UsersService.Unsuspend":                       true,
	}

	helperOverrides = map[string]overrideFunc{
		"s.search": func(arg string) (httpMethod, url string) {
			return "GET", fmt.Sprintf("search/%v", arg)
		},
	}

	// methodOverrides contains overrides for troublesome endpoints.
	methodOverrides = map[string]string{
		"OrganizationsService.EditOrgMembership: method orgs/%v/memberships/%v": "PUT",
		"OrganizationsService.EditOrgMembership: PUT user/memberships/orgs/%v":  "PATCH",
	}

	paramLegacyRE = regexp.MustCompile(`:[a-z_]+`)
	paramRE       = regexp.MustCompile(`{[a-z_]+}`)
)

type overrideFunc func(arg string) (httpMethod, url string)

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

type servicesMap map[string]*Service
type endpointsMap map[string]*Endpoint

func main() {
	flag.Parse()
	fset := token.NewFileSet()

	sourceFilter := func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasPrefix(fi.Name(), skipPrefix)
	}

	if err := os.Chdir("./github"); err != nil {
		if err := os.Chdir("../github"); err != nil {
			log.Fatalf("Please run this from the go-github directory.")
		}
	}

	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Step 1 - get a map of all services.
	services := findAllServices(pkgs)

	// Step 2 - find all the API service endpoints.
	iter := &realAstFileIterator{fset: fset, pkgs: pkgs}
	endpoints, err := findAllServiceEndpoints(iter, services)
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	// Step 3 - resolve all missing httpMethods from helperMethods.
	// Additionally, use existing URLs as hints to pre-cache all apiDocs.
	docCache := &documentCache{}
	usedHelpers, endpointsByFilename := resolveHelpersAndCacheDocs(endpoints, docCache)

	// Step 4 - validate and rewrite all URLs, skipping used helper methods.
	frw := &liveFileRewriter{fset: fset}
	validateRewriteURLs(usedHelpers, endpointsByFilename, docCache, frw)

	logf("Done.")
}

type usedHelpersMap map[string]bool
type endpointsByFilenameMap map[string][]*Endpoint

// FileRewriter read/writes files and converts AST token positions.
type FileRewriter interface {
	Position(token.Pos) token.Position
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, buf []byte, mode os.FileMode) error
}

// liveFileRewriter implements FileRewriter.
type liveFileRewriter struct {
	fset *token.FileSet
}

func (lfr *liveFileRewriter) Position(pos token.Pos) token.Position { return lfr.fset.Position(pos) }
func (lfr *liveFileRewriter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
func (lfr *liveFileRewriter) WriteFile(filename string, buf []byte, mode os.FileMode) error {
	return ioutil.WriteFile(filename, buf, mode)
}

func validateRewriteURLs(usedHelpers usedHelpersMap, endpointsByFilename endpointsByFilenameMap, docCache documentCacheReader, fileRewriter FileRewriter) {
	for filename, slc := range endpointsByFilename {
		logf("Step 4 - Processing %v methods in %v ...", len(slc), filename)

		var fileEdits []*FileEdit
		for _, endpoint := range slc {
			fullName := fmt.Sprintf("%v.%v", endpoint.serviceName, endpoint.endpointName)
			if usedHelpers[fullName] {
				logf("Step 4 - skipping used helper method %q", fullName)
				continue
			}

			// First, find the correct GitHub v3 API URL by httpMethod and urlFormat.
			for _, path := range endpoint.urlFormats {
				path = strings.ReplaceAll(path, "%d", "%v")
				path = strings.ReplaceAll(path, "%s", "%v")

				// Check the overrides.
				endpoint.checkHTTPMethodOverride(path)

				methodAndPath := fmt.Sprintf("%v %v", endpoint.httpMethod, path)
				url, ok := docCache.URLByMethodAndPath(methodAndPath)
				if !ok {
					if i := len(endpoint.endpointComments); i > 0 {
						pos := fileRewriter.Position(endpoint.endpointComments[i-1].Pos())
						fmt.Printf("%v:%v:%v: WARNING: unable to find online docs for %q: (%v)\nPLEASE CHECK MANUALLY AND FIX.\n", pos.Filename, pos.Line, pos.Column, fullName, methodAndPath)
					} else {
						fmt.Printf("%v: WARNING: unable to find online docs for %q: (%v)\nPLEASE CHECK MANUALLY AND FIX.\n", filename, fullName, methodAndPath)
					}
					continue
				}
				logf("found %q for: %q (%v)", url, fullName, methodAndPath)

				// Make sure URL is up-to-date.
				switch {
				case len(endpoint.enterpriseRefLines) > 1:
					log.Printf("WARNING: multiple Enterprise GitHub URLs found - skipping:")
					for i, refLine := range endpoint.enterpriseRefLines {
						log.Printf("line %v: %#v", i, refLine)
					}
				case len(endpoint.enterpriseRefLines) > 0:
					line := fmt.Sprintf(enterpriseRefFmt, url)
					cmt := endpoint.enterpriseRefLines[0]
					if cmt.Text != line {
						pos := fileRewriter.Position(cmt.Pos())
						logf("At byte offset %v:\nFOUND %q\nWANT: %q", pos.Offset, cmt.Text, line)
						fileEdits = append(fileEdits, &FileEdit{
							pos:      pos,
							fromText: cmt.Text,
							toText:   line,
						})
					}
				case len(endpoint.stdRefLines) > 1:
					var foundMatch bool
					line := fmt.Sprintf(stdRefFmt, url)
					for i, stdRefLine := range endpoint.stdRefLines {
						if stdRefLine.Text == line {
							foundMatch = true
							logf("found match with %v, not editing and removing from list", line)
							// Remove matching line
							endpoint.stdRefLines = append(endpoint.stdRefLines[:i], endpoint.stdRefLines[i+1:]...)
							break
						}
					}
					if !foundMatch { // Edit last stdRefLine, then remove it.
						cmt := endpoint.stdRefLines[len(endpoint.stdRefLines)-1]
						pos := fileRewriter.Position(cmt.Pos())
						logf("stdRefLines=%v: At byte offset %v:\nFOUND %q\nWANT: %q", len(endpoint.stdRefLines), pos.Offset, cmt.Text, line)
						fileEdits = append(fileEdits, &FileEdit{
							pos:      pos,
							fromText: cmt.Text,
							toText:   line,
						})
						endpoint.stdRefLines = endpoint.stdRefLines[:len(endpoint.stdRefLines)-1]
					}
				case len(endpoint.stdRefLines) > 0:
					line := fmt.Sprintf(stdRefFmt, url)
					cmt := endpoint.stdRefLines[0]
					if cmt.Text != line {
						pos := fileRewriter.Position(cmt.Pos())
						logf("stdRefLines=1: At byte offset %v:\nFOUND %q\nWANT: %q", pos.Offset, cmt.Text, line)
						fileEdits = append(fileEdits, &FileEdit{
							pos:      pos,
							fromText: cmt.Text,
							toText:   line,
						})
					}
					endpoint.stdRefLines = nil
				case len(endpoint.endpointComments) > 0:
					lastCmt := endpoint.endpointComments[len(endpoint.endpointComments)-1]
					// logf("lastCmt.Text=%q (len=%v)", lastCmt.Text, len(lastCmt.Text))
					pos := fileRewriter.Position(lastCmt.Pos())
					pos.Offset += len(lastCmt.Text)
					line := "\n" + fmt.Sprintf(stdRefFmt, url)
					if lastCmt.Text != "//" {
						line = "\n//" + line // Add blank comment line before URL.
					}
					// logf("line=%q (len=%v)", line, len(line))
					// logf("At byte offset %v: adding missing documentation:\n%q", pos.Offset, line)
					fileEdits = append(fileEdits, &FileEdit{
						pos:      pos,
						fromText: "",
						toText:   line,
					})
				default: // Missing documentation - add it.
					log.Printf("WARNING: file %v has no godoc comment string for method %v", fullName, methodAndPath)
				}
			}
		}

		if len(fileEdits) > 0 {
			b, err := fileRewriter.ReadFile(filename)
			if err != nil {
				log.Fatalf("ReadFile: %v", err)
			}

			log.Printf("Performing %v edits on file %v", len(fileEdits), filename)
			b = performBufferEdits(b, fileEdits)

			if err := fileRewriter.WriteFile(filename, b, 0644); err != nil {
				log.Fatalf("WriteFile: %v", err)
			}
		}
	}
}

func performBufferEdits(b []byte, fileEdits []*FileEdit) []byte {
	fileEdits = sortAndMergeFileEdits(fileEdits)

	for _, edit := range fileEdits {
		prelude := b[0:edit.pos.Offset]
		postlude := b[edit.pos.Offset+len(edit.fromText):]
		logf("At byte offset %v, replacing %v bytes with %v bytes\nBEFORE: %v\nAFTER : %v", edit.pos.Offset, len(edit.fromText), len(edit.toText), edit.fromText, edit.toText)
		b = []byte(fmt.Sprintf("%s%v%s", prelude, edit.toText, postlude))
	}

	return b
}

func sortAndMergeFileEdits(fileEdits []*FileEdit) []*FileEdit {
	// Sort edits from last to first in the file.
	// If the offsets are identical, sort the comment "toText" strings, ascending.
	var foundDups bool
	sort.Slice(fileEdits, func(a, b int) bool {
		if fileEdits[a].pos.Offset == fileEdits[b].pos.Offset {
			foundDups = true
			return fileEdits[a].toText < fileEdits[b].toText
		}
		return fileEdits[a].pos.Offset > fileEdits[b].pos.Offset
	})

	if !foundDups {
		return fileEdits
	}

	// Merge the duplicate edits.
	var mergedEdits []*FileEdit
	var dupOffsets []*FileEdit

	mergeFunc := func() {
		if len(dupOffsets) > 1 {
			isInsert := dupOffsets[0].fromText == ""
			var hasBlankCommentLine bool

			// Merge dups
			var lines []string
			for _, dup := range dupOffsets {
				if isInsert && strings.HasPrefix(dup.toText, "\n//\n//") {
					lines = append(lines, strings.TrimPrefix(dup.toText, "\n//"))
					hasBlankCommentLine = true
				} else {
					lines = append(lines, dup.toText)
				}
			}
			sort.Strings(lines)

			var joinStr string
			// if insert, no extra newlines
			if !isInsert { // if replacement - add newlines
				joinStr = "\n"
			}
			toText := strings.Join(lines, joinStr)
			if hasBlankCommentLine { // Add back in
				toText = "\n//" + toText
			}
			mergedEdits = append(mergedEdits, &FileEdit{
				pos:      dupOffsets[0].pos,
				fromText: dupOffsets[0].fromText,
				toText:   toText,
			})
		} else if len(dupOffsets) > 0 {
			// Move non-dup to final output
			mergedEdits = append(mergedEdits, dupOffsets[0])
		}
		dupOffsets = nil
	}

	lastOffset := -1
	for _, fileEdit := range fileEdits {
		if fileEdit.pos.Offset != lastOffset {
			mergeFunc()
		}
		dupOffsets = append(dupOffsets, fileEdit)
		lastOffset = fileEdit.pos.Offset
	}
	mergeFunc()
	return mergedEdits
}

// astFileIterator iterates over all files in an ast.Package.
type astFileIterator interface {
	// Finds the position of a token.
	Position(token.Pos) token.Position
	// Reset resets the iterator.
	Reset()
	// Next returns the next filenameAstFilePair pair or nil if done.
	Next() *filenameAstFilePair
}

type filenameAstFilePair struct {
	filename string
	astFile  *ast.File
}

// realAstFileIterator implements astFileIterator.
type realAstFileIterator struct {
	fset   *token.FileSet
	pkgs   map[string]*ast.Package
	ch     chan *filenameAstFilePair
	closed bool
}

func (rafi *realAstFileIterator) Position(pos token.Pos) token.Position {
	return rafi.fset.Position(pos)
}

func (rafi *realAstFileIterator) Reset() {
	if !rafi.closed && rafi.ch != nil {
		logf("Closing old channel on Reset")
		close(rafi.ch)
	}
	rafi.ch = make(chan *filenameAstFilePair, 10)
	rafi.closed = false

	go func() {
		var count int
		for _, pkg := range rafi.pkgs {
			for filename, f := range pkg.Files {
				// logf("Sending file #%v: %v to channel", count, filename)
				rafi.ch <- &filenameAstFilePair{filename: filename, astFile: f}
				count++
			}
		}
		rafi.closed = true
		close(rafi.ch)
		logf("Closed channel after sending %v files", count)
		if count == 0 {
			log.Fatalf("Processed no files. Did you run this from the go-github directory?")
		}
	}()
}

func (rafi *realAstFileIterator) Next() *filenameAstFilePair {
	for pair := range rafi.ch {
		// logf("Next: returning file %v", pair.filename)
		return pair
	}
	return nil
}

func findAllServices(pkgs map[string]*ast.Package) servicesMap {
	services := servicesMap{}
	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			if filename != "github.go" {
				continue
			}

			logf("Step 1 - Processing %v ...", filename)
			if err := findClientServices(filename, f, services); err != nil {
				log.Fatal(err)
			}
		}
	}
	return services
}

func findAllServiceEndpoints(iter astFileIterator, services servicesMap) (endpointsMap, error) {
	endpoints := endpointsMap{}
	iter.Reset()
	var errs []string // Collect all the errors and return in a big batch.
	for next := iter.Next(); next != nil; next = iter.Next() {
		filename, f := next.filename, next.astFile
		if filename == "github.go" {
			continue
		}

		if *debugFile != "" && !strings.Contains(filename, *debugFile) {
			continue
		}

		logf("Step 2 - Processing %v ...", filename)
		if err := processAST(filename, f, services, endpoints, iter); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return endpoints, nil
}

func resolveHelpersAndCacheDocs(endpoints endpointsMap, docCache documentCacheWriter) (usedHelpers usedHelpersMap, endpointsByFilename endpointsByFilenameMap) {
	usedHelpers = usedHelpersMap{}
	endpointsByFilename = endpointsByFilenameMap{}
	for k, v := range endpoints {
		if _, ok := endpointsByFilename[v.filename]; !ok {
			endpointsByFilename[v.filename] = []*Endpoint{}
		}
		endpointsByFilename[v.filename] = append(endpointsByFilename[v.filename], v)

		for _, cmt := range v.enterpriseRefLines {
			docCache.CacheDocFromInternet(cmt.Text, v.filename)
		}
		for _, cmt := range v.stdRefLines {
			docCache.CacheDocFromInternet(cmt.Text, v.filename)
		}

		if v.httpMethod == "" && v.helperMethod != "" {
			fullName := fmt.Sprintf("%v.%v", v.serviceName, v.helperMethod)
			hm, ok := endpoints[fullName]
			if !ok {
				log.Fatalf("Unable to find helper method %q for %q", fullName, k)
			}
			if hm.httpMethod == "" {
				log.Fatalf("Helper method %q for %q has empty httpMethod: %#v", fullName, k, hm)
			}
			v.httpMethod = hm.httpMethod
			usedHelpers[fullName] = true
		}
	}

	return usedHelpers, endpointsByFilename
}

type documentCacheReader interface {
	URLByMethodAndPath(string) (string, bool)
}

type documentCacheWriter interface {
	CacheDocFromInternet(urlWithFragmentID, filename string)
}

// documentCache implements documentCacheReader and documentCachWriter.
type documentCache struct {
	apiDocs            map[string]map[string][]*Endpoint // cached by URL, then mapped by web fragment identifier.
	urlByMethodAndPath map[string]string
}

func (dc *documentCache) URLByMethodAndPath(methodAndPath string) (string, bool) {
	url, ok := dc.urlByMethodAndPath[methodAndPath]
	return url, ok
}

func (dc *documentCache) CacheDocFromInternet(urlWithID, filename string) {
	if dc.apiDocs == nil {
		dc.apiDocs = map[string]map[string][]*Endpoint{} // cached by URL, then mapped by web fragment identifier.
		dc.urlByMethodAndPath = map[string]string{}
	}

	url := getURL(urlWithID)
	if _, ok := dc.apiDocs[url]; ok {
		return // already cached
	}

	logf("GET %q ...", url)
	resp, err := http.Get(url)
	check("Unable to get URL: %v: %v", url, err)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("filename: %v - url %v - StatusCode=%v", filename, url, resp.StatusCode)
	}

	finalURL := resp.Request.URL.String()
	url = getURL(finalURL)
	logf("The final URL is: %v; url=%v\n", finalURL, url)

	b, err := ioutil.ReadAll(resp.Body)
	check("Unable to read body of URL: %v, %v", url, err)
	check("Unable to close body of URL: %v, %v", url, resp.Body.Close())
	dc.apiDocs[url] = parseWebPageEndpoints(string(b))
	logf("Found %v web page fragment identifiers.", len(dc.apiDocs[url]))
	if len(dc.apiDocs[url]) == 0 {
		logf("webage text: %s", b)
	}

	// Now reverse-map the methods+paths to URLs.
	for fragID, v := range dc.apiDocs[url] {
		logf("For fragID=%q, found %v endpoints.", fragID, len(v))
		for _, endpoint := range v {
			logf("For fragID=%q, endpoint=%q, found %v paths.", fragID, endpoint, len(endpoint.urlFormats))
			for _, path := range endpoint.urlFormats {
				methodAndPath := fmt.Sprintf("%v %v", endpoint.httpMethod, path)
				dc.urlByMethodAndPath[methodAndPath] = fmt.Sprintf("%v#%v", url, fragID)
				logf("urlByMethodAndPath[%q] = %q", methodAndPath, dc.urlByMethodAndPath[methodAndPath])
			}
		}
	}
}

// FileEdit represents an edit that needs to be performed on a file.
type FileEdit struct {
	pos      token.Position
	fromText string
	toText   string
}

func getURL(s string) string {
	i := strings.Index(s, "http")
	if i < 0 {
		return ""
	}
	j := strings.Index(s, "#")
	if j < i {
		s = s[i:]
	} else {
		s = s[i:j]
	}
	if !strings.HasSuffix(s, "/") { // Prevent unnecessary redirects if possible.
		s += "/"
	}
	return s
}

// Service represents a go-github service.
type Service struct {
	serviceName string
}

// Endpoint represents an API endpoint in this repo.
type Endpoint struct {
	endpointName string
	filename     string
	serviceName  string
	urlFormats   []string
	httpMethod   string
	helperMethod string // If populated, httpMethod lives in helperMethod.

	enterpriseRefLines []*ast.Comment
	stdRefLines        []*ast.Comment
	endpointComments   []*ast.Comment
}

// String helps with debugging by providing an easy-to-read summary of the endpoint.
func (e *Endpoint) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  filename: %v\n", e.filename))
	b.WriteString(fmt.Sprintf("  serviceName: %v\n", e.serviceName))
	b.WriteString(fmt.Sprintf("  endpointName: %v\n", e.endpointName))
	b.WriteString(fmt.Sprintf("  httpMethod: %v\n", e.httpMethod))
	if e.helperMethod != "" {
		b.WriteString(fmt.Sprintf("  helperMethod: %v\n", e.helperMethod))
	}
	for i := 0; i < len(e.urlFormats); i++ {
		b.WriteString(fmt.Sprintf("  urlFormats[%v]: %v\n", i, e.urlFormats[i]))
	}
	for i := 0; i < len(e.enterpriseRefLines); i++ {
		b.WriteString(fmt.Sprintf("  enterpriseRefLines[%v]: comment: %v\n", i, e.enterpriseRefLines[i].Text))
	}
	for i := 0; i < len(e.stdRefLines); i++ {
		b.WriteString(fmt.Sprintf("  stdRefLines[%v]: comment: %v\n", i, e.stdRefLines[i].Text))
	}
	return b.String()
}

func (e *Endpoint) checkHTTPMethodOverride(path string) {
	lookupOverride := fmt.Sprintf("%v.%v: %v %v", e.serviceName, e.endpointName, e.httpMethod, path)
	logf("Looking up override for %q", lookupOverride)
	if v, ok := methodOverrides[lookupOverride]; ok {
		logf("overriding method for %v to %q", lookupOverride, v)
		e.httpMethod = v
		return
	}
}

func processAST(filename string, f *ast.File, services servicesMap, endpoints endpointsMap, iter astFileIterator) error {
	var errs []string

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl: // Doc, Recv, Name, Type, Body
			if decl.Recv == nil || len(decl.Recv.List) != 1 || decl.Name == nil || decl.Body == nil {
				continue
			}

			recv := decl.Recv.List[0]
			se, ok := recv.Type.(*ast.StarExpr) // Star, X
			if !ok || se.X == nil || len(recv.Names) != 1 {
				if decl.Name.Name != "String" && decl.Name.Name != "Equal" && decl.Name.Name != "IsPullRequest" {
					pos := iter.Position(recv.Pos())
					if id, ok := recv.Type.(*ast.Ident); ok {
						pos = iter.Position(id.Pos())
					}
					errs = append(errs, fmt.Sprintf("%v:%v:%v: method %v does not use a pointer receiver and needs fixing!", pos.Filename, pos.Line, pos.Column, decl.Name))
				}
				continue
			}
			recvType, ok := se.X.(*ast.Ident) // NamePos, Name, Obj
			if !ok {
				return fmt.Errorf("unhandled se.X = %T", se.X)
			}
			serviceName := recvType.Name
			if _, ok := services[serviceName]; !ok {
				continue
			}
			endpointName := decl.Name.Name
			fullName := fmt.Sprintf("%v.%v", serviceName, endpointName)
			if skipMethods[fullName] {
				logf("skipping %v", fullName)
				continue
			}

			receiverName := recv.Names[0].Name

			logf("\n\nast.FuncDecl: %#v", *decl)       // Doc, Recv, Name, Type, Body
			logf("ast.FuncDecl.Name: %#v", *decl.Name) // NamePos, Name, Obj(nil)
			// logf("ast.FuncDecl.Recv: %#v", *decl.Recv)  // Opening, List, Closing
			logf("ast.FuncDecl.Recv.List[0]: %#v", *recv) // Doc, Names, Type, Tag, Comment
			// for i, name := range decl.Recv.List[0].Names {
			// 	logf("recv.name[%v] = %v", i, name.Name)
			// }
			logf("recvType = %#v", recvType)
			var enterpriseRefLines []*ast.Comment
			var stdRefLines []*ast.Comment
			var endpointComments []*ast.Comment
			if decl.Doc != nil {
				endpointComments = decl.Doc.List
				for i, comment := range decl.Doc.List {
					logf("doc.comment[%v] = %#v", i, *comment)
					// if strings.Contains(comment.Text, enterpriseURL) {
					// 	enterpriseRefLines = append(enterpriseRefLines, comment)
					// } else
					if strings.Contains(comment.Text, stdURL) {
						stdRefLines = append(stdRefLines, comment)
					}
				}
				logf("%v comment lines, %v enterprise URLs, %v standard URLs", len(decl.Doc.List), len(enterpriseRefLines), len(stdRefLines))
			}

			bd := &bodyData{receiverName: receiverName}
			if err := bd.parseBody(decl.Body); err != nil { // Lbrace, List, Rbrace
				return fmt.Errorf("parseBody: %v", err)
			}

			ep := &Endpoint{
				endpointName:       endpointName,
				filename:           filename,
				serviceName:        serviceName,
				urlFormats:         bd.urlFormats,
				httpMethod:         bd.httpMethod,
				helperMethod:       bd.helperMethod,
				enterpriseRefLines: enterpriseRefLines,
				stdRefLines:        stdRefLines,
				endpointComments:   endpointComments,
			}
			// ep.checkHTTPMethodOverride("")
			endpoints[fullName] = ep
			logf("endpoints[%q] = %#v", fullName, endpoints[fullName])
			if ep.httpMethod == "" && (ep.helperMethod == "" || len(ep.urlFormats) == 0) {
				return fmt.Errorf("could not find body info: %#v", *ep)
			}
		case *ast.GenDecl:
		default:
			return fmt.Errorf("unhandled decl type: %T", decl)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// bodyData contains information found in a BlockStmt.
type bodyData struct {
	receiverName string // receiver name of method to help identify helper methods.
	httpMethod   string
	urlVarName   string
	urlFormats   []string
	assignments  []lhsrhs
	helperMethod string // If populated, httpMethod lives in helperMethod.
}

func (b *bodyData) parseBody(body *ast.BlockStmt) error {
	logf("body=%#v", *body)

	// Find the variable used for the format string, its one-or-more values,
	// and the httpMethod used for the NewRequest.
	for _, stmt := range body.List {
		switch stmt := stmt.(type) {
		case *ast.AssignStmt:
			hm, uvn, hlp, asgn := processAssignStmt(b.receiverName, stmt)
			if b.httpMethod != "" && hm != "" && b.httpMethod != hm {
				return fmt.Errorf("found two httpMethod values: %q and %q", b.httpMethod, hm)
			}
			if hm != "" {
				b.httpMethod = hm
				// logf("parseBody: httpMethod=%v", b.httpMethod)
			}
			if hlp != "" {
				b.helperMethod = hlp
			}
			b.assignments = append(b.assignments, asgn...)
			// logf("assignments=%#v", b.assignments)
			if b.urlVarName == "" && uvn != "" {
				b.urlVarName = uvn
				// logf("parseBody: urlVarName=%v", b.urlVarName)
				// By the time the urlVarName is found, all assignments should
				// have already taken place so that we can find the correct
				// ones and determine the urlFormats.
				for _, lr := range b.assignments {
					if lr.lhs == b.urlVarName {
						b.urlFormats = append(b.urlFormats, lr.rhs)
						logf("found urlFormat: %v", lr.rhs)
					}
				}
			}
		case *ast.DeclStmt:
			logf("*ast.DeclStmt: %#v", *stmt)
		case *ast.DeferStmt:
			logf("*ast.DeferStmt: %#v", *stmt)
		case *ast.ExprStmt:
			logf("*ast.ExprStmt: %#v", *stmt)
		case *ast.IfStmt:
			if err := b.parseIf(stmt); err != nil {
				return err
			}
		case *ast.RangeStmt:
			logf("*ast.RangeStmt: %#v", *stmt)
		case *ast.ReturnStmt: // Return Results
			logf("*ast.ReturnStmt: %#v", *stmt)
			if len(stmt.Results) > 0 {
				ce, ok := stmt.Results[0].(*ast.CallExpr)
				if ok {
					recv, funcName, args := processCallExpr(ce)
					logf("return CallExpr: recv=%q, funcName=%q, args=%#v", recv, funcName, args)
					// If the httpMethod has not been found at this point, but
					// this method is calling a helper function, then see if
					// any of its arguments match a previous assignment, then
					// record the urlFormat and remember the helper method.
					if b.httpMethod == "" && len(args) > 1 && recv == b.receiverName {
						if args[0] != "ctx" {
							return fmt.Errorf("expected helper function to get ctx as first arg: %#v, %#v", args, *b)
						}
						if len(b.assignments) == 0 && len(b.urlFormats) == 0 {
							b.urlFormats = append(b.urlFormats, strings.Trim(args[1], `"`))
							b.helperMethod = funcName
							logf("found urlFormat: %v and helper method: %v", b.urlFormats[0], b.helperMethod)
						} else {
							for _, lr := range b.assignments {
								if lr.lhs == args[1] { // Multiple matches are possible. Loop over all assignments.
									b.urlVarName = args[1]
									b.urlFormats = append(b.urlFormats, lr.rhs)
									b.helperMethod = funcName
									logf("found urlFormat: %v and helper method: %v", lr.rhs, b.helperMethod)
								}
							}
						}
					}
				}
			}
		case *ast.SwitchStmt:
			logf("*ast.SwitchStmt: %#v", *stmt)
		default:
			return fmt.Errorf("unhandled stmt type: %T", stmt)
		}
	}
	logf("parseBody: assignments=%#v", b.assignments)

	return nil
}

func (b *bodyData) parseIf(stmt *ast.IfStmt) error {
	logf("*ast.IfStmt: %#v", *stmt)
	if err := b.parseBody(stmt.Body); err != nil {
		return err
	}
	logf("if body: b=%#v", *b)
	if stmt.Else != nil {
		switch els := stmt.Else.(type) {
		case *ast.BlockStmt:
			if err := b.parseBody(els); err != nil {
				return err
			}
			logf("if else: b=%#v", *b)
		case *ast.IfStmt:
			if err := b.parseIf(els); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled else stmt type %T", els)
		}
	}

	return nil
}

// lhsrhs represents an assignment with a variable name on the left
// and a string on the right - used to find the URL format string.
type lhsrhs struct {
	lhs string
	rhs string
}

func processAssignStmt(receiverName string, stmt *ast.AssignStmt) (httpMethod, urlVarName, helperMethod string, assignments []lhsrhs) {
	logf("*ast.AssignStmt: %#v", *stmt) // Lhs, TokPos, Tok, Rhs
	var lhs []string
	for _, expr := range stmt.Lhs {
		switch expr := expr.(type) {
		case *ast.Ident: // NamePos, Name, Obj
			logf("processAssignStmt: *ast.Ident: %#v", expr)
			lhs = append(lhs, expr.Name)
		case *ast.SelectorExpr: // X, Sel
			logf("processAssignStmt: *ast.SelectorExpr: %#v", expr)
		default:
			log.Fatalf("unhandled AssignStmt Lhs type: %T", expr)
		}
	}

	for i, expr := range stmt.Rhs {
		switch expr := expr.(type) {
		case *ast.BasicLit: // ValuePos, Kind, Value
			v := strings.Trim(expr.Value, `"`)
			if !strings.HasPrefix(v, "?") { // Hack to remove "?recursive=1"
				assignments = append(assignments, lhsrhs{lhs: lhs[i], rhs: v})
			}
		case *ast.BinaryExpr:
			logf("processAssignStmt: *ast.BinaryExpr: %#v", *expr)
		case *ast.CallExpr: // Fun, Lparen, Args, Ellipsis, Rparen
			recv, funcName, args := processCallExpr(expr)
			logf("processAssignStmt: CallExpr: recv=%q, funcName=%q, args=%#v", recv, funcName, args)
			switch funcName {
			case "addOptions":
				if v := strings.Trim(args[0], `"`); v != args[0] {
					assignments = append(assignments, lhsrhs{lhs: lhs[i], rhs: v})
					urlVarName = lhs[i]
				} else {
					urlVarName = args[0]
				}
			case "Sprintf":
				assignments = append(assignments, lhsrhs{lhs: lhs[i], rhs: strings.Trim(args[0], `"`)})
			case "NewRequest":
				httpMethod = strings.Trim(args[0], `"`)
				urlVarName = args[1]
			case "NewUploadRequest":
				httpMethod = "POST"
				urlVarName = args[0]
			}
			if recv == receiverName && len(args) > 1 && args[0] == "ctx" { // This might be a helper method.
				fullName := fmt.Sprintf("%v.%v", recv, funcName)
				logf("checking for override: fullName=%v", fullName)
				if fn, ok := helperOverrides[fullName]; ok {
					logf("found helperOverride for %v", fullName)
					hm, url := fn(strings.Trim(args[1], `"`))
					httpMethod = hm
					urlVarName = "u" // arbitrary
					assignments = []lhsrhs{{lhs: urlVarName, rhs: url}}
				} else {
					urlVarName = args[1] // For this to work correctly, the URL must be the second arg to the helper method!
					helperMethod = funcName
					logf("found possible helper method: funcName=%v, urlVarName=%v", funcName, urlVarName)
				}
			}
		case *ast.CompositeLit: // Type, Lbrace, Elts, Rbrace, Incomplete
			logf("processAssignStmt: *ast.CompositeLit: %#v", *expr)
		case *ast.FuncLit:
			logf("processAssignStmt: *ast.FuncLit: %#v", *expr)
		case *ast.SelectorExpr:
			logf("processAssignStmt: *ast.SelectorExpr: %#v", *expr)
		case *ast.UnaryExpr: // OpPos, Op, X
			logf("processAssignStmt: *ast.UnaryExpr: %#v", *expr)
		case *ast.TypeAssertExpr: // X, Lparen, Type, Rparen
			logf("processAssignStmt: *ast.TypeAssertExpr: %#v", *expr)
		case *ast.Ident: // NamePos, Name, Obj
			logf("processAssignStmt: *ast.Ident: %#v", *expr)
		default:
			log.Fatalf("unhandled AssignStmt Rhs type: %T", expr)
		}
	}
	logf("urlVarName=%v, assignments=%#v", urlVarName, assignments)

	return httpMethod, urlVarName, helperMethod, assignments
}

func processCallExpr(expr *ast.CallExpr) (recv, funcName string, args []string) {
	logf("*ast.CallExpr: %#v", *expr)

	for _, arg := range expr.Args {
		switch arg := arg.(type) {
		case *ast.ArrayType:
			logf("processCallExpr: *ast.ArrayType: %#v", arg)
		case *ast.BasicLit: // ValuePos, Kind, Value
			args = append(args, arg.Value) // Do not trim quotes here so as to identify it later as a string literal.
		case *ast.CallExpr: // Fun, Lparen, Args, Ellipsis, Rparen
			logf("processCallExpr: *ast.CallExpr: %#v", arg)
			r, fn, as := processCallExpr(arg)
			if r == "fmt" && fn == "Sprintf" && len(as) > 0 { // Special case - return format string.
				args = append(args, as[0])
			}
		case *ast.CompositeLit:
			logf("processCallExpr: *ast.CompositeLit: %#v", arg) // Type, Lbrace, Elts, Rbrace, Incomplete
		case *ast.Ident: // NamePos, Name, Obj
			args = append(args, arg.Name)
		case *ast.MapType:
			logf("processCallExpr: *ast.MapType: %#v", arg)
		case *ast.SelectorExpr: // X, Sel
			logf("processCallExpr: *ast.SelectorExpr: %#v", arg)
			x, ok := arg.X.(*ast.Ident)
			if ok { // special case
				switch name := fmt.Sprintf("%v.%v", x.Name, arg.Sel.Name); name {
				case "http.MethodGet":
					args = append(args, http.MethodGet)
				case "http.MethodHead":
					args = append(args, http.MethodHead)
				case "http.MethodPost":
					args = append(args, http.MethodPost)
				case "http.MethodPut":
					args = append(args, http.MethodPut)
				case "http.MethodPatch":
					args = append(args, http.MethodPatch)
				case "http.MethodDelete":
					args = append(args, http.MethodDelete)
				case "http.MethodConnect":
					args = append(args, http.MethodConnect)
				case "http.MethodOptions":
					args = append(args, http.MethodOptions)
				case "http.MethodTrace":
					args = append(args, http.MethodTrace)
				default:
					args = append(args, name)
				}
			}
		case *ast.StarExpr:
			logf("processCallExpr: *ast.StarExpr: %#v", arg)
		case *ast.StructType:
			logf("processCallExpr: *ast.StructType: %#v", arg)
		case *ast.UnaryExpr: // OpPos, Op, X
			switch x := arg.X.(type) {
			case *ast.Ident:
				args = append(args, x.Name)
			case *ast.CompositeLit: // Type, Lbrace, Elts, Rbrace, Incomplete
				logf("processCallExpr: *ast.CompositeLit: %#v", x)
			default:
				log.Fatalf("processCallExpr: unhandled UnaryExpr.X arg type: %T", arg.X)
			}
		default:
			log.Fatalf("processCallExpr: unhandled arg type: %T", arg)
		}
	}

	switch fun := expr.Fun.(type) {
	case *ast.Ident: // NamePos, Name, Obj
		funcName = fun.Name
	case *ast.SelectorExpr: // X, Sel
		funcName = fun.Sel.Name
		switch x := fun.X.(type) {
		case *ast.Ident: // NamePos, Name, Obj
			logf("processCallExpr: X recv *ast.Ident=%#v", x)
			recv = x.Name
		case *ast.ParenExpr:
			logf("processCallExpr: X recv *ast.ParenExpr: %#v", x)
		case *ast.SelectorExpr: // X, Sel
			logf("processCallExpr: X recv *ast.SelectorExpr: %#v", x.Sel)
			recv = x.Sel.Name
		case *ast.CallExpr: // Fun, LParen, Args, Ellipsis, RParen
			logf("processCallExpr: X recv *ast.CallExpr: %#v", x)
		default:
			log.Fatalf("processCallExpr: unhandled X receiver type: %T, funcName=%q", x, funcName)
		}
	default:
		log.Fatalf("processCallExpr: unhandled Fun: %T", expr.Fun)
	}

	return recv, funcName, args
}

// findClientServices finds all go-github services from the Client struct.
func findClientServices(filename string, f *ast.File, services servicesMap) error {
	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok != token.TYPE || len(decl.Specs) != 1 {
				continue
			}
			ts, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok || decl.Doc == nil || ts.Name == nil || ts.Type == nil || ts.Name.Name != "Client" {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || st.Fields == nil || len(st.Fields.List) == 0 {
				continue
			}

			for _, field := range st.Fields.List {
				se, ok := field.Type.(*ast.StarExpr)
				if !ok || se.X == nil || len(field.Names) != 1 {
					continue
				}
				id, ok := se.X.(*ast.Ident)
				if !ok {
					continue
				}
				name := id.Name
				if !strings.HasSuffix(name, "Service") {
					continue
				}

				services[name] = &Service{serviceName: name}
			}

			return nil // Found all services in Client struct.
		}
	}

	return fmt.Errorf("unable to find Client struct in github.go")
}

func check(fmtStr string, args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		log.Fatalf(fmtStr, args...)
	}
}

// parseWebPageEndpoints returns endpoint information, mapped by
// web page fragment identifier.
func parseWebPageEndpoints(buf string) map[string][]*Endpoint {
	result := map[string][]*Endpoint{}

	// The GitHub v3 API web pages do not appear to be auto-generated
	// and therefore, the XML decoder is too strict to reliably parse them.
	// Here is a tiny example where the XML decoder completely fails
	// due to mal-formed HTML:
	//
	//   <optgroup label="Overview">
	//     <option value="/v3/">API Overview</a></h3>
	//     <option value="/v3/media/">Media Types</option>
	//   ...
	//   </optgroup>

	parts := splitHTML(buf)
	var lastFragmentID string
	for _, part := range parts {
		for _, method := range httpMethods {
			if strings.HasPrefix(part, method) {
				endpoint := parseEndpoint(part, method)
				if lastFragmentID == "" {
					log.Fatalf("parseWebPageEndpoints: empty lastFragmentID")
				}
				result[lastFragmentID] = append(result[lastFragmentID], endpoint)
			}
		}

		if i := strings.LastIndex(part, fragmentIDString); i >= 0 {
			b := part[i+len(fragmentIDString):]
			i = strings.Index(b, `"`)
			if i >= 0 {
				lastFragmentID = b[:i]
				logf("Found lastFragmentID: %v", lastFragmentID)
			}
		}
	}

	return result
}

func splitHTML(buf string) []string {
	var result []string
	for buf != "" {
		i := strings.Index(buf, codeLegacySplitString)
		j := strings.Index(buf, codeSplitString)
		switch {
		case i < 0 && j < 0:
			result = append(result, buf)
			buf = ""
		case j < 0, i >= 0 && j >= 0 && i < j:
			result = append(result, buf[:i])
			buf = buf[i+len(codeLegacySplitString):]
		case i < 0, i >= 0 && j >= 0 && j < i:
			result = append(result, buf[:j])
			buf = buf[j+len(codeSplitString):]
		default:
			log.Fatalf("splitHTML: i=%v, j=%v", i, j)
		}
	}
	return result
}

func parseEndpoint(s, method string) *Endpoint {
	eol := strings.Index(s, "\n")
	if eol < 0 {
		eol = len(s)
	}
	if v := strings.Index(s, "&apos;"); v > len(method) && v < eol {
		eol = v
	}
	if v := strings.Index(s, "<"); v > len(method) && v < eol {
		eol = v
	}
	// if v := strings.Index(s, "{"); v > len(method) && v < eol {
	// 	eol = v
	// }
	path := strings.TrimSpace(s[len(method):eol])
	path = strings.TrimPrefix(path, "{server}")
	path = paramLegacyRE.ReplaceAllString(path, "%v")
	path = paramRE.ReplaceAllString(path, "%v")
	// strip leading garbage
	if i := strings.Index(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	path = strings.TrimSuffix(path, ".")
	logf("Found endpoint: %v %v", method, path)
	return &Endpoint{
		urlFormats: []string{path},
		httpMethod: method,
	}
}

var httpMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}
