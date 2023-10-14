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
	"github.com/google/go-github/v56/github"
)

func TestUpdateURLs(t *testing.T) {
	res := runTest(t, "testdata/update-urls", "update-urls")
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestValidate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		res := runTest(t, "testdata/validate_invalid", "validate")
		res.assertErr("found 4 issues in")
		res.assertOutput("", `
Method AService.MissingFromMetadata does not exist in metadata.yaml. Please add it.
Method AService.Get has operation which is does not use the canonical name. You may be able to automatically fix this by running 'script/metadata.sh canonize': GET /a/{a_id_noncanonical}.
Name in override_operations does not exist in operations or openapi_operations: GET /a/{a_id_noncanonical2}
Name in override_operations does not exist in operations or openapi_operations: GET /fake/{a_id}
`)
		res.checkGolden()
	})

	t.Run("valid", func(t *testing.T) {
		res := runTest(t, "testdata/validate_valid", "validate")
		res.assertOutput("", "")
		res.assertNoErr()
		res.checkGolden()
	})
}

func TestUpdateMetadata(t *testing.T) {
	testServer := newTestServer(t, "main", map[string]interface{}{
		"api.github.com/api.github.com.json": openapi3.T{
			Paths: openapi3.Paths{
				"/a/{a_id}": &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				},
			},
		},
		"ghec/ghec.json": openapi3.T{
			Paths: openapi3.Paths{
				"/a/b/{a_id}": &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				},
			},
		},
		"ghes-3.9/ghes-3.9.json": openapi3.T{
			Paths: openapi3.Paths{
				"/a/b/{a_id}": &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				},
			},
		},
		"ghes-3.10/ghes-3.10.json": openapi3.T{
			Paths: openapi3.Paths{
				"/a/b/{a_id}": &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				},
			},
		},
		"ghes-2.22/ghes-2.22.json": openapi3.T{
			Paths: openapi3.Paths{
				"/a/b/{a_id}": &openapi3.PathItem{
					Get: &openapi3.Operation{
						ExternalDocs: &openapi3.ExternalDocs{
							URL: "https://docs.github.com/rest/reference/a",
						},
					},
				},
			},
		},
	})

	res := runTest(t, "testdata/update-metadata", "update-metadata", "--github-url", testServer.URL)
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestCanonize(t *testing.T) {
	res := runTest(t, "testdata/canonize", "canonize")
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestFormat(t *testing.T) {
	res := runTest(t, "testdata/format", "format")
	res.assertOutput("", "")
	res.assertNoErr()
	res.checkGolden()
}

func updateGoldenDir(t *testing.T, origDir, resultDir, goldenDir string) {
	t.Helper()
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
	if os.Getenv("UPDATE_GOLDEN") != "" {
		updateGoldenDir(t, origDir, resultDir, goldenDir)
		return
	}
	checked := map[string]bool{}
	_, err := os.Stat(goldenDir)
	if err == nil {
		assertNilError(t, filepath.Walk(goldenDir, func(wantPath string, info fs.FileInfo, err error) error {
			relPath := mustRel(t, goldenDir, wantPath)
			if err != nil || info.IsDir() {
				return err
			}
			assertEqualFiles(t, wantPath, filepath.Join(resultDir, relPath))
			checked[relPath] = true
			return nil
		}))
	}
	assertNilError(t, filepath.Walk(origDir, func(wantPath string, info fs.FileInfo, err error) error {
		relPath := mustRel(t, origDir, wantPath)
		if err != nil || info.IsDir() || checked[relPath] {
			return err
		}
		assertEqualFiles(t, wantPath, filepath.Join(resultDir, relPath))
		checked[relPath] = true
		return nil
	}))
	assertNilError(t, filepath.Walk(resultDir, func(resultPath string, info fs.FileInfo, err error) error {
		relPath := mustRel(t, resultDir, resultPath)
		if err != nil || info.IsDir() || checked[relPath] {
			return err
		}
		return fmt.Errorf("found unexpected file:\n%s", relPath)
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
	if !strings.Contains(r.err.Error(), want) {
		r.t.Errorf("expected error to contain %q, got %q", want, r.err.Error())
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

func newTestServer(t *testing.T, ref string, files map[string]interface{}) *httptest.Server {
	t.Helper()
	jsonHandler := func(wantQuery url.Values, val interface{}) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			gotQuery := r.URL.Query()
			queryDiff := cmp.Diff(wantQuery, gotQuery)
			if queryDiff != "" {
				t.Errorf("query mismatch for %s (-want +got):\n%s", r.URL.Path, queryDiff)
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
		jsonHandler(emptyQuery, &github.RepositoryCommit{SHA: github.String("s")}),
	)
	var descriptionsContent []*github.RepositoryContent
	for name, content := range files {
		descriptionsContent = append(descriptionsContent, &github.RepositoryContent{
			Name: github.String(path.Base(path.Dir(name))),
		})
		mux.HandleFunc(
			path.Join(repoPath, "contents/descriptions", path.Dir(name)),
			jsonHandler(refQuery, []*github.RepositoryContent{
				{
					Name:        github.String(path.Base(name)),
					DownloadURL: github.String(server.URL + "/dl/" + name),
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

func assertEqualFiles(t *testing.T, want, got string) {
	t.Helper()
	wantBytes, err := os.ReadFile(want)
	if !assertNilError(t, err) {
		return
	}
	wantBytes = bytes.ReplaceAll(wantBytes, []byte("\r\n"), []byte("\n"))
	gotBytes, err := os.ReadFile(got)
	if !assertNilError(t, err) {
		return
	}
	gotBytes = bytes.ReplaceAll(gotBytes, []byte("\r\n"), []byte("\n"))
	if bytes.Equal(wantBytes, gotBytes) {
		return
	}
	diff := cmp.Diff(string(wantBytes), string(gotBytes))
	t.Errorf("files %q and %q differ: %s", want, got, diff)
}

func assertNilError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Error(err)
		return false
	}
	return true
}
