// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v84/github"
)

func TestUpdateGo(t *testing.T) {
	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		res := runTest(t, "testdata/update-go/valid", "update-go")
		res.assertOutput("", "")
		res.assertNoErr()
		res.checkGolden()
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		res := runTest(t, "testdata/update-go/invalid", "update-go")
		res.assertOutput("", "")
		res.assertErr(`
no operations defined for AService.NoOperation
no operations defined for AService.NoComment
ambiguous operation "GET /ambiguous/{}" could match any of: [GET /ambiguous/{id} GET /ambiguous/{name}]
could not find operation "GET /missing/{id}" in openapi_operations.yaml
duplicate operation: GET /a/{a_id}
`)
		res.checkGolden()
	})
}

func TestUnused(t *testing.T) {
	t.Parallel()
	res := runTest(t, "testdata/unused", "unused")
	res.assertOutput(`
Found 3 unused operations

GET /a/{a_id}
doc:     https://docs.github.com/rest/a/a#overridden-get-a

POST /a/{a_id}
doc:     https://docs.github.com/rest/a/a#update-a

GET /undocumented/{undocumented_id}
`, "")
}

//nolint:paralleltest // cannot use t.Parallel() when helper calls t.Setenv
func TestUpdateOpenAPI(t *testing.T) {
	testServer := newTestServer(t, "main", map[string]any{
		"api.github.com/api.github.com.json": openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath("/a/{a_id}", &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				})),
		},
		"ghec/ghec.json": openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath("/a/b/{a_id}", &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				})),
		},
		"ghes-3.9/ghes-3.9.json": openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath("/a/b/{a_id}", &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				})),
		},
		"ghes-3.10/ghes-3.10.json": openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath("/a/b/{a_id}", &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				})),
		},
		"ghes-2.22/ghes-2.22.json": openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath("/a/b/{a_id}", &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				})),
		},
	})

	res := runTest(t, "testdata/update-openapi", "update-openapi", "--github-url", testServer.URL)
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestFormat(t *testing.T) {
	t.Parallel()
	res := runTest(t, "testdata/format", "format")
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func updateGoldenDir(t *testing.T, origDir, resultDir, goldenDir string) {
	t.Helper()
	if os.Getenv("UPDATE_GOLDEN") == "" {
		return
	}
	assertNilError(t, os.RemoveAll(goldenDir))
	assertNilError(t, filepath.WalkDir(resultDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		relName := mustRel(t, resultDir, path)
		origName := filepath.Join(origDir, relName)
		_, err = os.Stat(origName)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(filepath.Dir(filepath.Join(goldenDir, relName)), d.Type())
				if err != nil {
					return err
				}
				return copyFile(path, filepath.Join(goldenDir, relName))
			}
			return err
		}
		resContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		origContent, err := os.ReadFile(origName)
		if err != nil {
			return err
		}
		if bytes.Equal(resContent, origContent) {
			return nil
		}
		return copyFile(path, filepath.Join(goldenDir, relName))
	}))
}

func checkGoldenDir(t *testing.T, origDir, resultDir, goldenDir string) {
	t.Helper()
	golden := true
	t.Cleanup(func() {
		t.Helper()
		if !golden {
			t.Log("To regenerate golden files run `UPDATE_GOLDEN=1 script/test.sh`")
		}
	})
	updateGoldenDir(t, origDir, resultDir, goldenDir)
	checked := map[string]bool{}
	_, err := os.Stat(goldenDir)
	if err == nil {
		assertNilError(t, filepath.Walk(goldenDir, func(wantPath string, info fs.FileInfo, err error) error {
			relPath := mustRel(t, goldenDir, wantPath)
			if err != nil || info.IsDir() {
				return err
			}
			if !assertEqualFiles(t, wantPath, filepath.Join(resultDir, relPath)) {
				golden = false
			}
			checked[relPath] = true
			return nil
		}))
	}
	assertNilError(t, filepath.Walk(origDir, func(wantPath string, info fs.FileInfo, err error) error {
		relPath := mustRel(t, origDir, wantPath)
		if err != nil || info.IsDir() || checked[relPath] {
			return err
		}
		if !assertEqualFiles(t, wantPath, filepath.Join(resultDir, relPath)) {
			golden = false
		}
		checked[relPath] = true
		return nil
	}))
	assertNilError(t, filepath.Walk(resultDir, func(resultPath string, info fs.FileInfo, err error) error {
		relPath := mustRel(t, resultDir, resultPath)
		if err != nil || info.IsDir() || checked[relPath] {
			return err
		}
		golden = false
		return fmt.Errorf("found unexpected file:\n%v", relPath)
	}))
}

func mustRel(t *testing.T, base, target string) string {
	t.Helper()
	rel, err := filepath.Rel(base, target)
	assertNilError(t, err)
	return rel
}

func copyDir(t *testing.T, dst, src string) error {
	fmt.Println("dst", dst)
	dst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	return filepath.Walk(src, func(srcPath string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		dstPath := filepath.Join(dst, mustRel(t, src, srcPath))
		err = copyFile(srcPath, dstPath)
		return err
	})
}

func copyFile(src, dst string) (errOut error) {
	srcDirStat, err := os.Stat(filepath.Dir(src))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(dst), srcDirStat.Mode())
	if err != nil {
		return err
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		e := dstFile.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		e := srcFile.Close()
		if errOut == nil {
			errOut = e
		}
	}()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

type testRun struct {
	t       *testing.T
	workDir string
	srcDir  string
	stdOut  bytes.Buffer
	stdErr  bytes.Buffer
	err     error
}

func (r testRun) checkGolden() {
	r.t.Helper()
	checkGoldenDir(r.t, r.srcDir, r.workDir, filepath.Join("testdata", "golden", r.t.Name()))
}

func (r testRun) assertOutput(stdout, stderr string) {
	r.t.Helper()
	assertEqualStrings(r.t, strings.TrimSpace(stdout), strings.TrimSpace(r.stdOut.String()))
	assertEqualStrings(r.t, strings.TrimSpace(stderr), strings.TrimSpace(r.stdErr.String()))
}

func (r testRun) assertNoErr() {
	r.t.Helper()
	assertNilError(r.t, r.err)
}

func (r testRun) assertErr(want string) {
	r.t.Helper()
	if r.err == nil {
		r.t.Error("expected error")
		return
	}
	if strings.TrimSpace(r.err.Error()) != strings.TrimSpace(want) {
		r.t.Errorf("unexpected error:\nwant:\n%v\ngot:\n%v", want, r.err.Error())
	}
}

func runTest(t *testing.T, srcDir string, args ...string) testRun {
	t.Helper()
	srcDir = filepath.FromSlash(srcDir)
	res := testRun{
		t:       t,
		workDir: t.TempDir(),
		srcDir:  srcDir,
	}
	err := copyDir(t, res.workDir, srcDir)
	if err != nil {
		t.Error(err)
		return res
	}
	res.err = run(
		append(args, "-C", res.workDir),
		[]kong.Option{kong.Writers(&res.stdOut, &res.stdErr)},
	)
	return res
}

func newTestServer(t *testing.T, ref string, files map[string]any) *httptest.Server {
	t.Helper()
	jsonHandler := func(wantQuery url.Values, val any) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			gotQuery := r.URL.Query()
			queryDiff := cmp.Diff(wantQuery, gotQuery)
			if queryDiff != "" {
				t.Errorf("query mismatch for %v (-want +got):\n%v", r.URL.Path, queryDiff)
			}
			w.WriteHeader(200)
			err := json.NewEncoder(w).Encode(val)
			if err != nil {
				panic(err)
			}
		}
	}
	repoPath := "/api/v3/repos/github/rest-api-description"
	emptyQuery := url.Values{}
	refQuery := url.Values{"ref": []string{ref}}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(
		path.Join(repoPath, "commits", ref),
		jsonHandler(emptyQuery, &github.RepositoryCommit{SHA: github.Ptr("s")}),
	)
	var descriptionsContent []*github.RepositoryContent
	for name, content := range files {
		descriptionsContent = append(descriptionsContent, &github.RepositoryContent{
			Name: github.Ptr(path.Base(path.Dir(name))),
		})
		mux.HandleFunc(
			path.Join(repoPath, "contents/descriptions", path.Dir(name)),
			jsonHandler(refQuery, []*github.RepositoryContent{
				{
					Name:        github.Ptr(path.Base(name)),
					DownloadURL: github.Ptr(server.URL + "/dl/" + name),
				},
			}),
		)
		mux.HandleFunc(
			path.Join("/dl", name),
			jsonHandler(emptyQuery, content),
		)
	}
	mux.HandleFunc(
		path.Join(repoPath, "contents/descriptions"),
		jsonHandler(refQuery, descriptionsContent),
	)
	t.Cleanup(server.Close)
	t.Setenv("GITHUB_TOKEN", "fake token")
	return server
}

func assertEqualStrings(t *testing.T, want, got string) {
	t.Helper()
	diff := cmp.Diff(want, got)
	if diff != "" {
		t.Error(diff)
	}
}

func assertEqualFiles(t *testing.T, want, got string) bool {
	t.Helper()
	wantBytes, err := os.ReadFile(want)
	if !assertNilError(t, err) {
		return false
	}
	wantBytes = bytes.ReplaceAll(wantBytes, []byte("\r\n"), []byte("\n"))
	gotBytes, err := os.ReadFile(got)
	if !assertNilError(t, err) {
		return false
	}
	gotBytes = bytes.ReplaceAll(gotBytes, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(wantBytes, gotBytes) {
		diff := cmp.Diff(string(wantBytes), string(gotBytes))
		t.Errorf("files %q and %q differ: %v", want, got, diff)
		return false
	}
	return true
}

func assertNilError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Error(err)
		return false
	}
	return true
}
