// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryContent_GetContent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		encoding, content *string // input encoding and content
		want              string  // desired output
		wantErr           bool    // whether an error is expected
	}{
		{
			encoding: Ptr(""),
			content:  Ptr("hello"),
			want:     "hello",
			wantErr:  false,
		},
		{
			encoding: nil,
			content:  Ptr("hello"),
			want:     "hello",
			wantErr:  false,
		},
		{
			encoding: nil,
			content:  nil,
			want:     "",
			wantErr:  false,
		},
		{
			encoding: Ptr("base64"),
			content:  Ptr("aGVsbG8="),
			want:     "hello",
			wantErr:  false,
		},
		{
			encoding: Ptr("bad"),
			content:  Ptr("aGVsbG8="),
			want:     "",
			wantErr:  true,
		},
		{
			encoding: Ptr("none"),
			content:  nil,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		r := RepositoryContent{Encoding: tt.encoding, Content: tt.content}
		got, err := r.GetContent()
		if err != nil && !tt.wantErr {
			t.Errorf("RepositoryContent(%v, %v) returned unexpected error: %v",
				stringOrNil(tt.encoding), stringOrNil(tt.content), err)
		}
		if err == nil && tt.wantErr {
			t.Errorf("RepositoryContent(%v, %v) did not return unexpected error",
				stringOrNil(tt.encoding), stringOrNil(tt.content))
		}
		if want := tt.want; got != want {
			t.Errorf("RepositoryContent.GetContent returned %+v, want %+v", got, want)
		}
	}
}

// stringOrNil converts a potentially null string pointer to string.
// For non-nil input pointer, the returned string is enclosed in double-quotes.
func stringOrNil(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", *s)
}

func TestRepositoriesService_GetReadme(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/readme", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "encoding": "base64",
		  "size": 5362,
		  "name": "README.md",
		  "path": "README.md"
		}`)
	})
	ctx := t.Context()
	readme, _, err := client.Repositories.GetReadme(ctx, "o", "r", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetReadme returned error: %v", err)
	}
	want := &RepositoryContent{Type: Ptr("file"), Name: Ptr("README.md"), Size: Ptr(5362), Encoding: Ptr("base64"), Path: Ptr("README.md")}
	if !cmp.Equal(readme, want) {
		t.Errorf("Repositories.GetReadme returned %+v, want %+v", readme, want)
	}

	const methodName = "GetReadme"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetReadme(ctx, "\n", "\n", &RepositoryContentGetOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetReadme(ctx, "o", "r", &RepositoryContentGetOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DownloadContents_SuccessForFile(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f",
          "content": "foo",
		  "download_url": "`+serverURL+baseURLPath+`/download/f"
		}`)
	})

	ctx := t.Context()
	r, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContents returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("Repositories.DownloadContents returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo"; got != want {
		t.Errorf("Repositories.DownloadContents returned %v, want %v", got, want)
	}

	const methodName = "DownloadContents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.DownloadContents(ctx, "\n", "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DownloadContents_SuccessForDirectory(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f"
		}`)
	})
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		  "download_url": "`+serverURL+baseURLPath+`/download/f"
		}]`)
	})
	mux.HandleFunc("/download/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "foo")
	})

	ctx := t.Context()
	r, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContents returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("Repositories.DownloadContents returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo"; got != want {
		t.Errorf("Repositories.DownloadContents returned %v, want %v", got, want)
	}

	const methodName = "DownloadContents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.DownloadContents(ctx, "\n", "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DownloadContents_FailedResponse(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"type": "file",
			"name": "f"
		  }`)
	})
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"type": "file",
			"name": "f",
			"download_url": "`+serverURL+baseURLPath+`/download/f"
		  }]`)
	})
	mux.HandleFunc("/download/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "foo error")
	})

	ctx := t.Context()
	r, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContents returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusInternalServerError; got != want {
		t.Errorf("Repositories.DownloadContents returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo error"; got != want {
		t.Errorf("Repositories.DownloadContents returned %v, want %v", got, want)
	}
}

func TestRepositoriesService_DownloadContents_NoDownloadURL(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f",
		  "content": ""
		}`)
	})
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		  "content": ""
		}]`)
	})

	ctx := t.Context()
	reader, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Error("Repositories.DownloadContents did not return expected error")
	}

	if resp == nil {
		t.Error("Repositories.DownloadContents did not return expected response")
	}

	if reader != nil {
		t.Error("Repositories.DownloadContents did not return expected reader")
	}
}

func TestRepositoriesService_DownloadContents_NoFile(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f",
		  "content": ""
		}`)
	})

	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})

	ctx := t.Context()
	reader, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Error("Repositories.DownloadContents did not return expected error")
	}

	if resp == nil {
		t.Error("Repositories.DownloadContents did not return expected response")
	}

	if reader != nil {
		t.Error("Repositories.DownloadContents did not return expected reader")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_SuccessForFile(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f",
		  "download_url": "`+serverURL+baseURLPath+`/download/f",
          "content": "foo"
		}`)
	})

	ctx := t.Context()
	r, c, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContentsWithMeta returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo"; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned %v, want %v", got, want)
	}

	if c != nil && c.Name != nil {
		if got, want := *c.Name, "f"; got != want {
			t.Errorf("Repositories.DownloadContentsWithMeta returned content name %v, want %v", got, want)
		}
	} else {
		t.Error("Returned RepositoryContent is null")
	}

	const methodName = "DownloadContentsWithMeta"
	testBadOptions(t, methodName, func() (err error) {
		_, _, _, err = client.Repositories.DownloadContentsWithMeta(ctx, "\n", "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, cot, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		if cot != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, cot)
		}
		return resp, err
	})
}

func TestRepositoriesService_DownloadContentsWithMeta_SuccessForDirectory(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		  "download_url": "`+serverURL+baseURLPath+`/download/f"
		}]`)
	})
	mux.HandleFunc("/download/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "foo")
	})

	ctx := t.Context()
	r, c, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContentsWithMeta returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo"; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned %v, want %v", got, want)
	}

	if c != nil && c.Name != nil {
		if got, want := *c.Name, "f"; got != want {
			t.Errorf("Repositories.DownloadContentsWithMeta returned content name %v, want %v", got, want)
		}
	} else {
		t.Error("Returned RepositoryContent is null")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_FailedResponse(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	downloadURL := fmt.Sprintf("%v%v/download/f", serverURL, baseURLPath)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"type": "file",
			"name": "f",
			"download_url": "`+downloadURL+`"
		  }`)
	})
	mux.HandleFunc("/download/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "foo error")
	})
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"type": "file",
			"name": "f",
			"download_url": "`+downloadURL+`"
		  }]`)
	})

	ctx := t.Context()
	r, c, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err != nil {
		t.Errorf("Repositories.DownloadContentsWithMeta returned error: %v", err)
	}

	if got, want := resp.Response.StatusCode, http.StatusInternalServerError; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned status code %v, want %v", got, want)
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	r.Close()

	if got, want := string(bytes), "foo error"; got != want {
		t.Errorf("Repositories.DownloadContentsWithMeta returned %v, want %v", got, want)
	}

	if c != nil && c.Name != nil {
		if got, want := *c.Name, "f"; got != want {
			t.Errorf("Repositories.DownloadContentsWithMeta returned content name %v, want %v", got, want)
		}
	} else {
		t.Error("Returned RepositoryContent is null")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_NoDownloadURL(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d/f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "name": "f",
		}`)
	})
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		  "content": ""
		}]`)
	})

	ctx := t.Context()
	reader, contents, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected error")
	}

	if reader != nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected reader")
	}

	if resp == nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected response")
	}

	if contents == nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected content")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_NoFile(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})

	ctx := t.Context()
	_, _, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected error")
	}

	if resp == nil {
		t.Error("Repositories.DownloadContentsWithMeta did not return expected response")
	}
}

func TestRepositoriesService_GetContents_File(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "encoding": "base64",
		  "size": 20678,
		  "name": "LICENSE",
		  "path": "LICENSE"
		}`)
	})
	ctx := t.Context()
	fileContents, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "p", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := &RepositoryContent{Type: Ptr("file"), Name: Ptr("LICENSE"), Size: Ptr(20678), Encoding: Ptr("base64"), Path: Ptr("LICENSE")}
	if !cmp.Equal(fileContents, want) {
		t.Errorf("Repositories.GetContents returned %+v, want %+v", fileContents, want)
	}

	const methodName = "GetContents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, _, err = client.Repositories.GetContents(ctx, "\n", "\n", "\n", &RepositoryContentGetOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, _, resp, err := client.Repositories.GetContents(ctx, "o", "r", "p", &RepositoryContentGetOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetContents_FilenameNeedsEscape(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p#?%/中.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := t.Context()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "p#?%/中.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_DirectoryWithSpaces(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/some%20directory/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := t.Context()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some directory/file.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_PathWithParent(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/some/../directory/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := t.Context()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some/../directory/file.go", &RepositoryContentGetOptions{})
	if err == nil {
		t.Fatal("Repositories.GetContents expected error but got none")
	}
}

func TestRepositoriesService_GetContents_DirectoryWithPlusChars(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/some%20directory%2Bname/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := t.Context()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some directory+name/file.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_Directory(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "dir",
		  "name": "lib",
		  "path": "lib"
		},
		{
		  "type": "file",
		  "size": 20678,
		  "name": "LICENSE",
		  "path": "LICENSE"
		}]`)
	})
	ctx := t.Context()
	_, directoryContents, _, err := client.Repositories.GetContents(ctx, "o", "r", "p", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := []*RepositoryContent{
		{Type: Ptr("dir"), Name: Ptr("lib"), Path: Ptr("lib")},
		{Type: Ptr("file"), Name: Ptr("LICENSE"), Size: Ptr(20678), Path: Ptr("LICENSE")},
	}
	if !cmp.Equal(directoryContents, want) {
		t.Errorf("Repositories.GetContents_Directory returned %+v, want %+v", directoryContents, want)
	}
}

func TestRepositoriesService_CreateFile(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"content":{
				"name":"p"
			},
			"commit":{
				"message":"m",
				"sha":"f5f369044773ff9c6383c087466d12adb6fa0828"
			}
		}`)
	})
	message := "m"
	content := []byte("c")
	repositoryContentsOptions := &RepositoryContentFileOptions{
		Message:   &message,
		Content:   content,
		Committer: &CommitAuthor{Name: Ptr("n"), Email: Ptr("e")},
	}
	ctx := t.Context()
	createResponse, _, err := client.Repositories.CreateFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.CreateFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: &RepositoryContent{Name: Ptr("p")},
		Commit: Commit{
			Message: Ptr("m"),
			SHA:     Ptr("f5f369044773ff9c6383c087466d12adb6fa0828"),
		},
	}
	if !cmp.Equal(createResponse, want) {
		t.Errorf("Repositories.CreateFile returned %+v, want %+v", createResponse, want)
	}

	const methodName = "CreateFile"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateFile(ctx, "\n", "\n", "\n", repositoryContentsOptions)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateFile(ctx, "o", "r", "p", repositoryContentsOptions)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateFile(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"content":{
				"name":"p"
			},
			"commit":{
				"message":"m",
				"sha":"f5f369044773ff9c6383c087466d12adb6fa0828"
			}
		}`)
	})
	message := "m"
	content := []byte("c")
	sha := "f5f369044773ff9c6383c087466d12adb6fa0828"
	repositoryContentsOptions := &RepositoryContentFileOptions{
		Message:   &message,
		Content:   content,
		SHA:       &sha,
		Committer: &CommitAuthor{Name: Ptr("n"), Email: Ptr("e")},
	}
	ctx := t.Context()
	updateResponse, _, err := client.Repositories.UpdateFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.UpdateFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: &RepositoryContent{Name: Ptr("p")},
		Commit: Commit{
			Message: Ptr("m"),
			SHA:     Ptr("f5f369044773ff9c6383c087466d12adb6fa0828"),
		},
	}
	if !cmp.Equal(updateResponse, want) {
		t.Errorf("Repositories.UpdateFile returned %+v, want %+v", updateResponse, want)
	}

	const methodName = "UpdateFile"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdateFile(ctx, "\n", "\n", "\n", repositoryContentsOptions)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateFile(ctx, "o", "r", "p", repositoryContentsOptions)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteFile(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"content": null,
			"commit":{
				"message":"m",
				"sha":"f5f369044773ff9c6383c087466d12adb6fa0828"
			}
		}`)
	})
	message := "m"
	sha := "f5f369044773ff9c6383c087466d12adb6fa0828"
	repositoryContentsOptions := &RepositoryContentFileOptions{
		Message:   &message,
		SHA:       &sha,
		Committer: &CommitAuthor{Name: Ptr("n"), Email: Ptr("e")},
	}
	ctx := t.Context()
	deleteResponse, _, err := client.Repositories.DeleteFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.DeleteFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: nil,
		Commit: Commit{
			Message: Ptr("m"),
			SHA:     Ptr("f5f369044773ff9c6383c087466d12adb6fa0828"),
		},
	}
	if !cmp.Equal(deleteResponse, want) {
		t.Errorf("Repositories.DeleteFile returned %+v, want %+v", deleteResponse, want)
	}

	const methodName = "DeleteFile"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.DeleteFile(ctx, "\n", "\n", "\n", repositoryContentsOptions)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.DeleteFile(ctx, "o", "r", "p", repositoryContentsOptions)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetArchiveLink(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/tarball/yo", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})
			ctx := t.Context()
			url, resp, err := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{Ref: "yo"}, 1)
			if err != nil {
				t.Errorf("Repositories.GetArchiveLink returned error: %v", err)
			}
			if resp.StatusCode != http.StatusFound {
				t.Errorf("Repositories.GetArchiveLink returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}
			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Repositories.GetArchiveLink returned %+v, want %+v", url, want)
			}

			const methodName = "GetArchiveLink"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.GetArchiveLink(ctx, "\n", "\n", Tarball, &RepositoryContentGetOptions{}, 1)
				return err
			})

			// Add custom round tripper
			client.client.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("failed to get archive link")
			})
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, 1)
				return err
			})
		})
	}
}

func TestRepositoriesService_GetArchiveLink_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/tarball", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusMovedPermanently)
			})
			ctx := t.Context()
			_, resp, _ := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, 0)
			if resp.StatusCode != http.StatusMovedPermanently {
				t.Errorf("Repositories.GetArchiveLink returned status: %v, want %v", resp.StatusCode, http.StatusMovedPermanently)
			}
		})
	}
}

func TestRepositoriesService_GetArchiveLink_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, serverURL := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			// Mock a redirect link, which leads to an archive link
			mux.HandleFunc("/repos/o/r/tarball", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
				http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
			})
			mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})
			ctx := t.Context()
			url, resp, err := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, 1)
			if err != nil {
				t.Errorf("Repositories.GetArchiveLink returned error: %v", err)
			}
			if resp.StatusCode != http.StatusFound {
				t.Errorf("Repositories.GetArchiveLink returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}
			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Repositories.GetArchiveLink returned %+v, want %+v", url, want)
			}
		})
	}
}

func TestRepositoriesService_GetContents_NoTrailingSlashInDirectoryApiPath(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/contents/.github", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"ref": "mybranch"})
		fmt.Fprint(w, `{}`)
	})
	ctx := t.Context()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", ".github/", &RepositoryContentGetOptions{
		Ref: "mybranch",
	})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoryContent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryContent{}, "{}")

	r := &RepositoryContent{
		Type:            Ptr("type"),
		Target:          Ptr("target"),
		Encoding:        Ptr("encoding"),
		Size:            Ptr(1),
		Name:            Ptr("name"),
		Path:            Ptr("path"),
		Content:         Ptr("content"),
		SHA:             Ptr("sha"),
		URL:             Ptr("url"),
		GitURL:          Ptr("gurl"),
		HTMLURL:         Ptr("hurl"),
		DownloadURL:     Ptr("durl"),
		SubmoduleGitURL: Ptr("smgurl"),
	}

	want := `{
		"type": "type",
		"target": "target",
		"encoding": "encoding",
		"size": 1,
		"name": "name",
		"path": "path",
		"content": "content",
		"sha": "sha",
		"url": "url",
		"git_url": "gurl",
		"html_url": "hurl",
		"download_url": "durl",
		"submodule_git_url": "smgurl"
	}`

	testJSONMarshal(t, r, want)
}

func TestRepositoryContentResponse_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryContentResponse{}, `{"commit": {}}`)

	r := &RepositoryContentResponse{
		Content: &RepositoryContent{
			Type:            Ptr("type"),
			Target:          Ptr("target"),
			Encoding:        Ptr("encoding"),
			Size:            Ptr(1),
			Name:            Ptr("name"),
			Path:            Ptr("path"),
			Content:         Ptr("content"),
			SHA:             Ptr("sha"),
			URL:             Ptr("url"),
			GitURL:          Ptr("gurl"),
			HTMLURL:         Ptr("hurl"),
			DownloadURL:     Ptr("durl"),
			SubmoduleGitURL: Ptr("smgurl"),
		},
		Commit: Commit{
			SHA: Ptr("s"),
			Author: &CommitAuthor{
				Date:  &Timestamp{referenceTime},
				Name:  Ptr("n"),
				Email: Ptr("e"),
				Login: Ptr("u"),
			},
			Committer: &CommitAuthor{
				Date:  &Timestamp{referenceTime},
				Name:  Ptr("n"),
				Email: Ptr("e"),
				Login: Ptr("u"),
			},
			Message: Ptr("m"),
			Tree: &Tree{
				SHA: Ptr("s"),
				Entries: []*TreeEntry{{
					SHA:     Ptr("s"),
					Path:    Ptr("p"),
					Mode:    Ptr("m"),
					Type:    Ptr("t"),
					Size:    Ptr(1),
					Content: Ptr("c"),
					URL:     Ptr("u"),
				}},
				Truncated: Ptr(false),
			},
			Parents: nil,
			HTMLURL: Ptr("h"),
			URL:     Ptr("u"),
			Verification: &SignatureVerification{
				Verified:  Ptr(false),
				Reason:    Ptr("r"),
				Signature: Ptr("s"),
				Payload:   Ptr("p"),
			},
			NodeID:       Ptr("n"),
			CommentCount: Ptr(1),
		},
	}

	want := `{
		"content": {
			"type": "type",
			"target": "target",
			"encoding": "encoding",
			"size": 1,
			"name": "name",
			"path": "path",
			"content": "content",
			"sha": "sha",
			"url": "url",
			"git_url": "gurl",
			"html_url": "hurl",
			"download_url": "durl",
			"submodule_git_url": "smgurl"
		},
		"commit": {
			"sha": "s",
			"author": {
				"date": ` + referenceTimeStr + `,
				"name": "n",
				"email": "e",
				"username": "u"
			},
			"committer": {
				"date": ` + referenceTimeStr + `,
				"name": "n",
				"email": "e",
				"username": "u"
			},
			"message": "m",
			"tree": {
				"sha": "s",
				"tree": [
					{
						"sha": "s",
						"path": "p",
						"mode": "m",
						"type": "t",
						"size": 1,
						"content": "c",
						"url": "u"
					}
				],
				"truncated": false
			},
			"html_url": "h",
			"url": "u",
			"verification": {
				"verified": false,
				"reason": "r",
				"signature": "s",
				"payload": "p"
			},
			"node_id": "n",
			"comment_count": 1
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestRepositoryContentFileOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryContentFileOptions{}, `{"content": null}`)

	r := &RepositoryContentFileOptions{
		Message: Ptr("type"),
		Content: []byte{1},
		SHA:     Ptr("type"),
		Branch:  Ptr("type"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("name"),
			Email: Ptr("email"),
			Login: Ptr("login"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("name"),
			Email: Ptr("email"),
			Login: Ptr("login"),
		},
	}

	want := `{
		"message": "type",
		"content": "AQ==",
		"sha": "type",
		"branch": "type",
		"author": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		},
		"committer": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		}
	}`

	testJSONMarshal(t, r, want)
}
