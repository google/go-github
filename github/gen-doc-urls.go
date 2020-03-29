// -*- compile-command: "go run gen-doc-urls.go -v"; -*-
// DO NOT COMMIT WITH COMPILE-COMMAND.

// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-doc-urls generates GitHub URL docs for each service endpoint.
//
// It is meant to be used by go-github contributors in conjunction with the
// go generate tool before sending a PR to GitHub.
// Please see the CONTRIBUTING.md file for more information.
package main

import (
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
	skipPrefix = "gen-"

	enterpriseURL = "developer.github.com/enterprise"
	stdURL        = "developer.github.com"

	enterpriseRefFmt = "// GitHub Enterprise API docs: %v"
	stdRefFmt        = "// GitHub API docs: %v"
)

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")

	// methodBlacklist holds methods that do not have GitHub v3 API URLs.
	methodBlacklist = map[string]bool{
		"MarketplaceService.marketplaceURI":    true,
		"RepositoriesService.DownloadContents": true,
		"SearchService.search":                 true,
		"UsersService.GetByID":                 true,
	}

	helperOverrides = map[string]overrideFunc{
		"s.search": func(arg string) (httpMethod, url string) {
			return "GET", fmt.Sprintf("search/%v", arg)
		},
	}

	// methodOverrides contains overrides for troublesome endpoints.
	methodOverrides = map[string]string{
		"OrganizationsService.EditOrgMembership: method orgs/%v/memberships/%v":   "PUT",
		"OrganizationsService.EditOrgMembership: method user/memberships/orgs/%v": "PATCH",
	}

	paramRE = regexp.MustCompile(`:[a-z_]+`)
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

	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Step 1 - get a map of all services.
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

	// Step 2 - find all the API service endpoints.
	endpoints := endpointsMap{}
	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			if filename == "github.go" {
				continue
			}

			if filename != "activity_events.go" { // DEBUGGING ONLY!!!
				continue
			}

			logf("Step 2 - Processing %v ...", filename)
			if err := processAST(filename, f, services, endpoints); err != nil {
				log.Fatal(err)
			}
		}
	}

	apiDocs := map[string]map[string][]*Endpoint{} // cached by URL, then mapped by web fragment identifier.
	urlByMethodAndPath := map[string]string{}
	cacheDocs := func(s string) {
		url := getURL(s)
		if _, ok := apiDocs[url]; ok {
			return
		}

		// TODO: Enterprise URLs are currently causing problems - for example:
		// GET https://developer.github.com/enterprise/v3/enterprise-admin/users/
		// returns StatusCode=404
		if strings.Contains(url, "enterprise") {
			logf("Skipping troublesome Enterprise URL: %v", url)
			return
		}

		logf("GET %q ...", url)
		resp, err := http.Get(url)
		check("Unable to get URL: %v: %v", url, err)
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("url %v - StatusCode=%v", url, resp.StatusCode)
		}

		b, err := ioutil.ReadAll(resp.Body)
		check("Unable to read body of URL: %v, %v", url, err)
		check("Unable to close body of URL: %v, %v", url, resp.Body.Close())
		apiDocs[url] = parseWebPageEndpoints(string(b))

		// Now reverse-map the methods+paths to URLs.
		for fragID, v := range apiDocs[url] {
			for _, endpoint := range v {
				for _, path := range endpoint.urlFormats {
					methodAndPath := fmt.Sprintf("%v %v", endpoint.httpMethod, path)
					urlByMethodAndPath[methodAndPath] = fmt.Sprintf("%v#%v", url, fragID)
					logf("urlByMethodAndPath[%q] = %q", methodAndPath, urlByMethodAndPath[methodAndPath])
				}
			}
		}
	}

	// Step 3 - resolve all missing httpMethods from helperMethods.
	// Additionally, use existing URLs as hints to pre-cache all apiDocs.
	usedHelpers := map[string]bool{}
	endpointsByFilename := map[string][]*Endpoint{}
	for k, v := range endpoints {
		if _, ok := endpointsByFilename[v.filename]; !ok {
			endpointsByFilename[v.filename] = []*Endpoint{}
		}
		endpointsByFilename[v.filename] = append(endpointsByFilename[v.filename], v)

		for _, cmt := range v.enterpriseRefLines {
			cacheDocs(cmt.Text)
		}
		for _, cmt := range v.stdRefLines {
			cacheDocs(cmt.Text)
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

	// Step 4 - validate and rewrite all URLs, skipping used helper methods.
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
				methodAndPath := fmt.Sprintf("%v %v", endpoint.httpMethod, path)
				url, ok := urlByMethodAndPath[methodAndPath]
				if !ok {
					log.Printf("WARNING: Unable to find documentation for %v - %q: (%v)", filename, fullName, methodAndPath)
					continue
				}
				logf("found %q for: %q (%v)", url, fullName, methodAndPath)

				// Make sure URL is up-to-date.
				switch {
				case len(endpoint.enterpriseRefLines) > 1:
					log.Printf("WARNING: multiple Enterprise GitHub URLs found - skipping: %#v", endpoint.enterpriseRefLines)
				case len(endpoint.enterpriseRefLines) > 0:
					line := fmt.Sprintf(enterpriseRefFmt, url)
					cmt := endpoint.enterpriseRefLines[0]
					if cmt.Text != line {
						log.Printf("At token.pos=%v:\nFOUND %q\nWANT: %q", cmt.Pos(), cmt.Text, line)
						fileEdits = append(fileEdits, &FileEdit{
							pos:      fset.Position(cmt.Pos()),
							fromText: cmt.Text,
							toText:   line,
						})
					}
				case len(endpoint.stdRefLines) > 1:
					log.Printf("WARNING: multiple GitHub URLs found - skipping: %#v", endpoint.stdRefLines)
				case len(endpoint.stdRefLines) > 0:
					line := fmt.Sprintf(stdRefFmt, url)
					cmt := endpoint.stdRefLines[0]
					if cmt.Text != line {
						log.Printf("At token.pos=%v:\nFOUND %q\nWANT: %q", cmt.Pos(), cmt.Text, line)
						fileEdits = append(fileEdits, &FileEdit{
							pos:      fset.Position(cmt.Pos()),
							fromText: cmt.Text,
							toText:   line,
						})
					}
				default: // Missing documentation - add it.
					log.Printf("TODO: Add missing documentation: %q for: %q (%v)", url, fullName, methodAndPath)
				}
			}
		}

		if len(fileEdits) > 0 {
			logf("Performing %v edits on file %v", len(fileEdits), filename)
			// Sort edits from last to first in the file.
			sort.Slice(fileEdits, func(a, b int) bool { return fileEdits[b].pos.Offset < fileEdits[a].pos.Offset })

			b, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatalf("ReadFile: %v", err)
			}

			var lastOffset int
			for _, edit := range fileEdits {
				if edit.pos.Offset == lastOffset {
					logf("TODO: At offset %v, inserting second URL of %v bytes", edit.pos.Offset, len(edit.fromText), len(edit.toText))
					continue
				}
				lastOffset = edit.pos.Offset
				logf("At offset %v, replacing %v bytes with %v bytes", edit.pos.Offset, len(edit.fromText), len(edit.toText))
				before := b[0:edit.pos.Offset]
				after := b[edit.pos.Offset+len(edit.fromText):]
				b = []byte(fmt.Sprintf("%s%v%s", before, edit.toText, after))
			}

			if err := ioutil.WriteFile(filename, b, 0644); err != nil {
				log.Fatalf("WriteFile: %v", err)
			}
		}
	}

	logf("Done.")
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
		return s[i:]
	}
	return s[i:j]
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
}

func processAST(filename string, f *ast.File, services servicesMap, endpoints endpointsMap) error {
	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl: // Doc, Recv, Name, Type, Body
			if decl.Recv == nil || len(decl.Recv.List) != 1 || decl.Name == nil || decl.Body == nil {
				continue
			}
			recv := decl.Recv.List[0]
			se, ok := recv.Type.(*ast.StarExpr) // Star, X
			if !ok || se.X == nil || len(recv.Names) != 1 {
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
			if methodBlacklist[fullName] {
				logf("skipping %v", fullName)
				continue
			}

			receiverName := recv.Names[0].Name

			logf("ast.FuncDecl: %#v", *decl)           // Doc, Recv, Name, Type, Body
			logf("ast.FuncDecl.Name: %#v", *decl.Name) // NamePos, Name, Obj(nil)
			// logf("ast.FuncDecl.Recv: %#v", *decl.Recv)  // Opening, List, Closing
			logf("ast.FuncDecl.Recv.List[0]: %#v", *recv) // Doc, Names, Type, Tag, Comment
			// for i, name := range decl.Recv.List[0].Names {
			// 	logf("recv.name[%v] = %v", i, name.Name)
			// }
			logf("recvType = %#v", recvType)
			var enterpriseRefLines []*ast.Comment
			var stdRefLines []*ast.Comment
			if decl.Doc != nil {
				for i, comment := range decl.Doc.List {
					logf("doc.comment[%v] = %#v", i, *comment)
					if strings.Contains(comment.Text, enterpriseURL) {
						enterpriseRefLines = append(enterpriseRefLines, comment)
					} else if strings.Contains(comment.Text, stdURL) {
						stdRefLines = append(stdRefLines, comment)
					}
				}
				logf("%v comment lines, %v enterprise URLs, %v standard URLs", len(decl.Doc.List), len(enterpriseRefLines), len(stdRefLines))
			}

			bd := &bodyData{receiverName: receiverName}
			if err := bd.parseBody(decl.Body); err != nil { // Lbrace, List, Rbrace
				return fmt.Errorf("parseBody: %v", err)
			}

			if len(bd.urlFormats) == 1 {
				lookupOverride := fmt.Sprintf("%v.%v: %v %v", serviceName, endpointName, bd.httpMethod, bd.urlFormats[0])
				if v, ok := methodOverrides[lookupOverride]; ok {
					logf("overriding method for %v to %q", lookupOverride, v)
					bd.httpMethod = v
				}
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
			}
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

	// for _, comment := range f.Comments {
	// 	log.Printf("Found %v comments, starting with: %#v", len(comment.List), comment.List[0])
	// }

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
			assignments = append(assignments, lhsrhs{lhs: lhs[i], rhs: strings.Trim(expr.Value, `"`)})
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
		default:
			log.Fatalf("processCallExpr: unhandled X receiver type: %T", x)
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

	parts := strings.Split(buf, "<code>")
	var lastFragmentID string
	for _, part := range parts {
		for _, method := range httpMethods {
			if strings.HasPrefix(part, method) {
				eol := strings.Index(part, "\n")
				if eol < 0 {
					eol = len(part)
				}
				if v := strings.Index(part, "<"); v > len(method) && v < eol {
					eol = v
				}
				path := strings.TrimSpace(part[len(method):eol])
				path = paramRE.ReplaceAllString(path, "%v")
				// strip leading garbage
				if i := strings.Index(path, "/"); i >= 0 {
					path = path[i+1:]
				}
				path = strings.TrimSuffix(path, ".")
				logf("Found %v %v", method, path)
				result[lastFragmentID] = append(result[lastFragmentID], &Endpoint{
					urlFormats: []string{path},
					httpMethod: method,
				})
			}
		}

		if i := strings.LastIndex(part, "<a id="); i >= 0 {
			b := part[i+7:]
			i = strings.Index(b, `"`)
			if i >= 0 {
				lastFragmentID = b[:i]
				logf("Found lastFragmentID: %v", lastFragmentID)
			}
		}
	}

	return result
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
