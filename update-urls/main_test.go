// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pmezard/go-difflib/difflib"
)

type pipelineSetup struct {
	// Fields filled in by the unit test:
	baseURL              string
	endpointsFromWebsite endpointsByFragmentID
	filename             string
	serviceName          string
	originalGoSource     string
	wantGoSource         string
	wantNumEndpoints     int

	// Fields filled in by setup:
	docCache     *fakeDocCache
	fileRewriter *fakeFileRewriter
	iter         *fakeAstFileIterator
	services     servicesMap
	wantFailure  bool
}

func (ps *pipelineSetup) setup(t *testing.T, stripURLs, destroyReceiverPointers bool) *pipelineSetup {
	t.Helper()

	if stripURLs {
		// For every GitHub API doc URL, remove it from the original source,
		// and alternate between stripping the previous blank comment line and not.
		for removeBlank := false; true; removeBlank = !removeBlank {
			var changes bool
			if removeBlank {
				ps.originalGoSource, changes = removeNextURLAndOptionalBlank(ps.originalGoSource)
			} else {
				ps.originalGoSource, changes = removeNextURLLineOnly(ps.originalGoSource)
			}
			if !changes {
				break
			}
		}
		// log.Printf("Modified Go Source:\n%v", ps.originalGoSource)
	}

	if destroyReceiverPointers {
		from := fmt.Sprintf(" *%v) ", ps.serviceName)
		to := fmt.Sprintf(" %v) ", ps.serviceName)
		ps.originalGoSource = strings.ReplaceAll(ps.originalGoSource, from, to)
		ps.wantFailure = true // receiver pointers must be fixed before running.
	}

	ps.docCache = &fakeDocCache{
		t:         t,
		baseURL:   ps.baseURL,
		endpoints: ps.endpointsFromWebsite,
	}
	fset := token.NewFileSet()
	ps.fileRewriter = &fakeFileRewriter{fset: fset, in: ps.originalGoSource}
	ps.services = servicesMap{ps.serviceName: &Service{serviceName: ps.serviceName}}
	astFile, err := parser.ParseFile(fset, ps.filename, ps.originalGoSource, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	ps.iter = &fakeAstFileIterator{
		fset: fset,
		orig: &filenameAstFilePair{
			filename: ps.filename,
			astFile:  astFile,
		},
	}

	return ps
}

func (ps *pipelineSetup) validate(t *testing.T) {
	t.Helper()

	// Call pipeline
	endpoints, err := findAllServiceEndpoints(ps.iter, ps.services)
	if ps.wantFailure {
		if err != nil {
			// test successful - receivers must be pointers first
			return
		}
		t.Fatalf("Expected non-pointer receivers to fail parsing, but no error was raised")
	}
	if err != nil {
		t.Fatalf("Fail detected but not expected: %v", err)
	}

	// log.Printf("endpoints=%#v (%v)", endpoints, len(endpoints))
	if len(endpoints) != ps.wantNumEndpoints {
		t.Errorf("got %v endpoints, want %v", len(endpoints), ps.wantNumEndpoints)
	}
	usedHelpers, endpointsByFilename := resolveHelpersAndCacheDocs(endpoints, ps.docCache)
	// log.Printf("endpointsByFilename=%#v (%v)", endpointsByFilename, len(endpointsByFilename[ps.filename]))
	if len(endpointsByFilename[ps.filename]) != ps.wantNumEndpoints {
		t.Errorf("got %v endpointsByFilename, want %v", len(endpointsByFilename[ps.filename]), ps.wantNumEndpoints)
	}
	validateRewriteURLs(usedHelpers, endpointsByFilename, ps.docCache, ps.fileRewriter)

	if ps.fileRewriter.out == "" {
		t.Fatalf("No modifications were made to the file")
	}

	if ps.fileRewriter.out != ps.wantGoSource {
		diff := difflib.ContextDiff{
			A:        difflib.SplitLines(ps.fileRewriter.out),
			B:        difflib.SplitLines(ps.wantGoSource),
			FromFile: "got",
			ToFile:   "want",
			Context:  1,
			Eol:      "\n",
		}
		result, _ := difflib.GetContextDiffString(diff)
		t.Errorf(strings.Replace(result, "\t", " ", -1))
	}
}

var (
	urlWithBlankCommentRE = regexp.MustCompile(`(//\n)?// GitHub API docs: [^\n]+\n`)
	urlLineOnlyRE         = regexp.MustCompile(`// GitHub API docs: [^\n]+\n`)
)

func removeNextURLAndOptionalBlank(s string) (string, bool) {
	parts := urlWithBlankCommentRE.Split(s, 2)
	if len(parts) == 1 {
		return parts[0], false
	}
	return parts[0] + parts[1], true
}

func TestRemoveNextURLAndOptionalBlank(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    string
		changes bool
	}{
		{name: "empty string"},
		{name: "no URLs", s: "// line 1\n//\n// line 3", want: "// line 1\n//\n// line 3"},
		{
			name:    "URL without prior blank comment",
			s:       "// line 1\n// GitHub API docs: yeah\nfunc MyFunc() {\n",
			want:    "// line 1\nfunc MyFunc() {\n",
			changes: true,
		},
		{
			name:    "URL with prior blank comment",
			s:       "// line 1\n//\n// GitHub API docs: yeah\nfunc MyFunc() {\n",
			want:    "// line 1\nfunc MyFunc() {\n",
			changes: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			got, changes := removeNextURLAndOptionalBlank(tt.s)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
			if changes != tt.changes {
				t.Errorf("got changes = %v, want %v", changes, tt.changes)
			}
		})
	}
}

func removeNextURLLineOnly(s string) (string, bool) {
	parts := urlLineOnlyRE.Split(s, 2)
	if len(parts) == 1 {
		return parts[0], false
	}
	return parts[0] + parts[1], true
}

func TestRemoveNextURLLineOnly(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    string
		changes bool
	}{
		{name: "empty string"},
		{name: "no URLs", s: "// line 1\n//\n// line 3", want: "// line 1\n//\n// line 3"},
		{
			name:    "URL without prior blank comment",
			s:       "// line 1\n// GitHub API docs: yeah\nfunc MyFunc() {\n",
			want:    "// line 1\nfunc MyFunc() {\n",
			changes: true,
		},
		{
			name:    "URL with prior blank comment",
			s:       "// line 1\n//\n// GitHub API docs: yeah\nfunc MyFunc() {\n",
			want:    "// line 1\n//\nfunc MyFunc() {\n",
			changes: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			got, changes := removeNextURLLineOnly(tt.s)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
			if changes != tt.changes {
				t.Errorf("got changes = %v, want %v", changes, tt.changes)
			}
		})
	}
}

type endpointsByFragmentID map[string][]*Endpoint

// fakeDocCache implements documentCacheReader and documentCacheWriter.
type fakeDocCache struct {
	t         *testing.T
	baseURL   string
	endpoints endpointsByFragmentID
}

func (f *fakeDocCache) URLByMethodAndPath(methodAndPath string) (string, bool) {
	f.t.Helper()
	for fragmentID, endpoints := range f.endpoints {
		for _, endpoint := range endpoints {
			for _, urlFormat := range endpoint.urlFormats {
				key := fmt.Sprintf("%v %v", endpoint.httpMethod, urlFormat)
				if key == methodAndPath {
					url := fmt.Sprintf("%v#%v", f.baseURL, fragmentID)
					// log.Printf("URLByMethodAndPath(%q) = (%q, true)", methodAndPath, url)
					return url, true
				}
			}
		}
	}
	f.t.Fatalf("fakeDocCache.URLByMethodAndPath: unable to find method %v", methodAndPath)
	return "", false
}

func (f *fakeDocCache) CacheDocFromInternet(url, filename string) {} // no-op

// fakeFileRewriter implements FileRewriter.
type fakeFileRewriter struct {
	fset *token.FileSet
	in   string
	out  string
}

func (f *fakeFileRewriter) Position(pos token.Pos) token.Position {
	return f.fset.Position(pos)
}

func (f *fakeFileRewriter) ReadFile(filename string) ([]byte, error) {
	return []byte(f.in), nil
}

func (f *fakeFileRewriter) WriteFile(filename string, buf []byte, mode os.FileMode) error {
	f.out = string(buf)
	return nil
}

// fakeAstFileIterator implements astFileIterator.
type fakeAstFileIterator struct {
	orig, next *filenameAstFilePair
	fset       *token.FileSet
}

func (f *fakeAstFileIterator) Position(pos token.Pos) token.Position { return f.fset.Position(pos) }
func (f *fakeAstFileIterator) Reset()                                { f.next = f.orig }
func (f *fakeAstFileIterator) Next() *filenameAstFilePair {
	v := f.next
	f.next = nil
	return v
}

func TestSortAndMergeFileEdits(t *testing.T) {
	tests := []struct {
		name      string
		fileEdits []*FileEdit
		want      []*FileEdit
	}{
		{name: "no edits"},
		{
			name: "one edit",
			fileEdits: []*FileEdit{
				{toText: "one edit"},
			},
			want: []*FileEdit{
				{toText: "one edit"},
			},
		},
		{
			name: "two inserts at same offset - no extra blank comment",
			fileEdits: []*FileEdit{
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\n// one insert"},
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\n// second insert"},
			},
			want: []*FileEdit{
				{pos: token.Position{Offset: 2}, toText: "\n// one insert\n// second insert"},
			},
		},
		{
			name: "two inserts at same offset - strip extra blank comment",
			fileEdits: []*FileEdit{
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\n//\n// one insert"},
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\n//\n// second insert"},
			},
			want: []*FileEdit{
				{pos: token.Position{Offset: 2}, toText: "\n//\n// one insert\n// second insert"},
			},
		},
		{
			name: "two non-overlapping edits, low offset to high",
			fileEdits: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 0}, toText: "edit one"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit two"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit two"},
				{fromText: ".", pos: token.Position{Offset: 0}, toText: "edit one"},
			},
		},
		{
			name: "two non-overlapping edits, high offset to low",
			fileEdits: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit two"},
				{fromText: ".", pos: token.Position{Offset: 0}, toText: "edit one"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit two"},
				{fromText: ".", pos: token.Position{Offset: 0}, toText: "edit one"},
			},
		},
		{
			name: "two overlapping edits, text low to high",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 0"},
				{fromText: ".", toText: "edit 1"},
			},
			want: []*FileEdit{
				{fromText: ".", toText: "edit 0\nedit 1"},
			},
		},
		{
			name: "two overlapping edits, text high to low",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 1"},
				{fromText: ".", toText: "edit 0"},
			},
			want: []*FileEdit{
				{fromText: ".", toText: "edit 0\nedit 1"},
			},
		},
		{
			name: "dup, non-dup",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 1"},
				{fromText: ".", toText: "edit 0"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 2"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 2"},
				{fromText: ".", toText: "edit 0\nedit 1"},
			},
		},
		{
			name: "non-dup, dup",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 2"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 1"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 0"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 0\nedit 1"},
				{fromText: ".", toText: "edit 2"},
			},
		},
		{
			name: "dup, non-dup, dup",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 1"},
				{fromText: ".", toText: "edit 0"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 2"},
				{fromText: ".", pos: token.Position{Offset: 200}, toText: "edit 4"},
				{fromText: ".", pos: token.Position{Offset: 200}, toText: "edit 3"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 200}, toText: "edit 3\nedit 4"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 2"},
				{fromText: ".", toText: "edit 0\nedit 1"},
			},
		},
		{
			name: "non-dup, dup, non-dup",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 2"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 1"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 0"},
				{fromText: ".", pos: token.Position{Offset: 200}, toText: "edit 3"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 200}, toText: "edit 3"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 0\nedit 1"},
				{fromText: ".", toText: "edit 2"},
			},
		},
		{
			name: "triplet, non-dup",
			fileEdits: []*FileEdit{
				{fromText: ".", toText: "edit 1"},
				{fromText: ".", toText: "edit 0"},
				{fromText: ".", toText: "edit 2"},
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 3"},
			},
			want: []*FileEdit{
				{fromText: ".", pos: token.Position{Offset: 100}, toText: "edit 3"},
				{fromText: ".", toText: "edit 0\nedit 1\nedit 2"},
			},
		},
	}

	fileEditEqual := cmp.Comparer(func(a, b *FileEdit) bool {
		return a.fromText == b.fromText && a.pos == b.pos && a.toText == b.toText
	})

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			got := sortAndMergeFileEdits(tt.fileEdits)

			if len(got) != len(tt.want) {
				t.Errorf("len(got) = %v, len(want) = %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got); i++ {
				var wantFileEdit *FileEdit
				if i < len(tt.want) {
					wantFileEdit = tt.want[i]
				}
				if !cmp.Equal(got[i], wantFileEdit, fileEditEqual) {
					t.Errorf("got[%v] =\n%#v\nwant[%v]:\n%#v", i, got[i], i, wantFileEdit)
				}
			}
		})
	}
}

func TestPerformBufferEdits(t *testing.T) {
	tests := []struct {
		name      string
		fileEdits []*FileEdit
		s         string
		want      string
	}{
		{name: "no edits", s: "my\nshort\nfile\n", want: "my\nshort\nfile\n"},
		{
			name: "one edit",
			fileEdits: []*FileEdit{
				{pos: token.Position{Offset: 3}, fromText: "short", toText: "one edit"},
			},
			s:    "my\nshort\nfile\n",
			want: "my\none edit\nfile\n",
		},
		{
			name: "one insert",
			fileEdits: []*FileEdit{
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\none insert"},
			},
			s:    "my\nshort\nfile\n",
			want: "my\none insert\nshort\nfile\n",
		},
		{
			name: "two inserts at same offset",
			fileEdits: []*FileEdit{
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\none insert"},
				{pos: token.Position{Offset: 2}, fromText: "", toText: "\nsecond insert"},
			},
			s:    "my\nshort\nfile\n",
			want: "my\none insert\nsecond insert\nshort\nfile\n",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			b := performBufferEdits([]byte(tt.s), tt.fileEdits)
			got := string(b)

			if len(got) != len(tt.want) {
				t.Errorf("len(got) = %v, len(want) = %v", len(got), len(tt.want))
			}
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestGitURL(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{name: "empty string"},
		{name: "non-http", s: "howdy"},
		{
			name: "normal URL, no slash",
			s:    "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events",
			want: "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events/",
		},
		{
			name: "normal URL, with slash",
			s:    "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events/",
			want: "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events/",
		},
		{
			name: "normal URL, with fragment identifier",
			s:    "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events/#list-public-events",
			want: "https://docs.github.com/en/free-pro-team@latest/rest/reference/activity/events/",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			got := getURL(tt.s)
			if got != tt.want {
				t.Errorf("getURL = %v ; want %v", got, tt.want)
			}
		})
	}
}

var endpointEqual = cmp.Comparer(func(a, b *Endpoint) bool {
	if a.httpMethod != b.httpMethod {
		return false
	}
	return cmp.Equal(a.urlFormats, b.urlFormats)
})

func testWebPageHelper(t *testing.T, got, want map[string][]*Endpoint) {
	t.Helper()

	for k := range got {
		w, ok := want[k]
		if len(got[k]) != len(w) {
			t.Errorf("len(got[%q]) = %v, len(want[%q]) = %v", k, len(got[k]), k, len(w))
		}
		for i := 0; i < len(got[k]); i++ {
			var wantEndpoint *Endpoint
			if ok && i < len(w) {
				wantEndpoint = w[i]
			}
			if !cmp.Equal(got[k][i], wantEndpoint, endpointEqual) {
				t.Errorf("got[%q][%v] =\n%#v\nwant[%q][%v]:\n%#v", k, i, got[k][i], k, i, wantEndpoint)
			}
		}
	}
	for k := range want {
		if _, ok := got[k]; !ok {
			t.Errorf("got[%q] = nil\nwant[%q]:\n%#v", k, k, want[k])
		}
	}
}

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		method string
		want   *Endpoint
	}{
		{
			name: "orgs_projects: list-repository-projects",
			s: `GET /repos/{owner}/{repo}/projects&apos;</span>, {
`,
			method: "GET",
			want:   &Endpoint{urlFormats: []string{"repos/%v/%v/projects"}, httpMethod: "GET"},
		},
		{
			name: "orgs_projects: ListProjects",
			s: `GET /orgs/{org}/projects&apos;</span>, {
`,
			method: "GET",
			want:   &Endpoint{urlFormats: []string{"orgs/%v/projects"}, httpMethod: "GET"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v: %v", i, tt.name), func(t *testing.T) {
			got := parseEndpoint(tt.s, tt.method)

			if !cmp.Equal(got, tt.want, endpointEqual) {
				t.Errorf("parseEndpoint = %#v, want %#v", got, tt.want)
			}
		})
	}
}
