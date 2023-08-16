// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/google/go-cmp/cmp"
)

func TestRepositoryContent_GetContent(t *testing.T) {
	tests := []struct {
		encoding, content *string // input encoding and content
		want              string  // desired output
		wantErr           bool    // whether an error is expected
	}{
		{
			encoding: String(""),
			content:  String("hello"),
			want:     "hello",
			wantErr:  false,
		},
		{
			encoding: nil,
			content:  String("hello"),
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
			encoding: String("base64"),
			content:  String("aGVsbG8="),
			want:     "hello",
			wantErr:  false,
		},
		{
			encoding: String("bad"),
			content:  String("aGVsbG8="),
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		r := RepositoryContent{Encoding: tt.encoding, Content: tt.content}
		got, err := r.GetContent()
		if err != nil && !tt.wantErr {
			t.Errorf("RepositoryContent(%s, %s) returned unexpected error: %v",
				stringOrNil(tt.encoding), stringOrNil(tt.content), err)
		}
		if err == nil && tt.wantErr {
			t.Errorf("RepositoryContent(%s, %s) did not return unexpected error",
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
	client, mux, _, teardown := setup()
	defer teardown()
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
	ctx := context.Background()
	readme, _, err := client.Repositories.GetReadme(ctx, "o", "r", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetReadme returned error: %v", err)
	}
	want := &RepositoryContent{Type: String("file"), Name: String("README.md"), Size: Int(5362), Encoding: String("base64"), Path: String("README.md")}
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

func TestRepositoriesService_DownloadContents_Success(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()
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

	ctx := context.Background()
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
	client, mux, serverURL, teardown := setup()
	defer teardown()
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

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		}]`)
	})

	ctx := context.Background()
	_, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Errorf("Repositories.DownloadContents did not return expected error")
	}

	if resp == nil {
		t.Errorf("Repositories.DownloadContents did not return expected response")
	}
}

func TestRepositoriesService_DownloadContents_NoFile(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, resp, err := client.Repositories.DownloadContents(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Errorf("Repositories.DownloadContents did not return expected error")
	}

	if resp == nil {
		t.Errorf("Repositories.DownloadContents did not return expected response")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_Success(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()
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

	ctx := context.Background()
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
		t.Errorf("Returned RepositoryContent is null")
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

func TestRepositoriesService_DownloadContentsWithMeta_FailedResponse(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()
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

	ctx := context.Background()
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
		t.Errorf("Returned RepositoryContent is null")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_NoDownloadURL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		  "type": "file",
		  "name": "f",
		}]`)
	})

	ctx := context.Background()
	_, _, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Errorf("Repositories.DownloadContentsWithMeta did not return expected error")
	}

	if resp == nil {
		t.Errorf("Repositories.DownloadContentsWithMeta did not return expected response")
	}
}

func TestRepositoriesService_DownloadContentsWithMeta_NoFile(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, resp, err := client.Repositories.DownloadContentsWithMeta(ctx, "o", "r", "d/f", nil)
	if err == nil {
		t.Errorf("Repositories.DownloadContentsWithMeta did not return expected error")
	}

	if resp == nil {
		t.Errorf("Repositories.DownloadContentsWithMeta did not return expected response")
	}
}

func TestRepositoriesService_GetContents_File(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
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
	ctx := context.Background()
	fileContents, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "p", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := &RepositoryContent{Type: String("file"), Name: String("LICENSE"), Size: Int(20678), Encoding: String("base64"), Path: String("LICENSE")}
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
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/p#?%/中.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := context.Background()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "p#?%/中.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_DirectoryWithSpaces(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/some directory/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := context.Background()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some directory/file.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_PathWithParent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/some/../directory/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := context.Background()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some/../directory/file.go", &RepositoryContentGetOptions{})
	if err == nil {
		t.Fatal("Repositories.GetContents expected error but got none")
	}
}

func TestRepositoriesService_GetContents_DirectoryWithPlusChars(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/some directory+name/file.go", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})
	ctx := context.Background()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", "some directory+name/file.go", &RepositoryContentGetOptions{})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoriesService_GetContents_Directory(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
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
	ctx := context.Background()
	_, directoryContents, _, err := client.Repositories.GetContents(ctx, "o", "r", "p", &RepositoryContentGetOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := []*RepositoryContent{{Type: String("dir"), Name: String("lib"), Path: String("lib")},
		{Type: String("file"), Name: String("LICENSE"), Size: Int(20678), Path: String("LICENSE")}}
	if !cmp.Equal(directoryContents, want) {
		t.Errorf("Repositories.GetContents_Directory returned %+v, want %+v", directoryContents, want)
	}
}

func TestRepositoriesService_CreateFile(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
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
		Committer: &CommitAuthor{Name: String("n"), Email: String("e")},
	}
	ctx := context.Background()
	createResponse, _, err := client.Repositories.CreateFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.CreateFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: &RepositoryContent{Name: String("p")},
		Commit: Commit{
			Message: String("m"),
			SHA:     String("f5f369044773ff9c6383c087466d12adb6fa0828"),
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
	client, mux, _, teardown := setup()
	defer teardown()
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
		Committer: &CommitAuthor{Name: String("n"), Email: String("e")},
	}
	ctx := context.Background()
	updateResponse, _, err := client.Repositories.UpdateFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.UpdateFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: &RepositoryContent{Name: String("p")},
		Commit: Commit{
			Message: String("m"),
			SHA:     String("f5f369044773ff9c6383c087466d12adb6fa0828"),
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
	client, mux, _, teardown := setup()
	defer teardown()
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
		Committer: &CommitAuthor{Name: String("n"), Email: String("e")},
	}
	ctx := context.Background()
	deleteResponse, _, err := client.Repositories.DeleteFile(ctx, "o", "r", "p", repositoryContentsOptions)
	if err != nil {
		t.Errorf("Repositories.DeleteFile returned error: %v", err)
	}
	want := &RepositoryContentResponse{
		Content: nil,
		Commit: Commit{
			Message: String("m"),
			SHA:     String("f5f369044773ff9c6383c087466d12adb6fa0828"),
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
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/tarball/yo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})
	ctx := context.Background()
	url, resp, err := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{Ref: "yo"}, true)
	if err != nil {
		t.Errorf("Repositories.GetArchiveLink returned error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Repositories.GetArchiveLink returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}
	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Repositories.GetArchiveLink returned %+v, want %+v", url.String(), want)
	}

	const methodName = "GetArchiveLink"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetArchiveLink(ctx, "\n", "\n", Tarball, &RepositoryContentGetOptions{}, true)
		return err
	})

	// Add custom round tripper
	client.client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("failed to get archive link")
	})
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, true)
		return err
	})
}

func TestRepositoriesService_GetArchiveLink_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/tarball", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusMovedPermanently)
	})
	ctx := context.Background()
	_, resp, _ := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, false)
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Repositories.GetArchiveLink returned status: %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}
}

func TestRepositoriesService_GetArchiveLink_StatusMovedPermanently_followRedirects(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()
	// Mock a redirect link, which leads to an archive link
	mux.HandleFunc("/repos/o/r/tarball", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})
	ctx := context.Background()
	url, resp, err := client.Repositories.GetArchiveLink(ctx, "o", "r", Tarball, &RepositoryContentGetOptions{}, true)
	if err != nil {
		t.Errorf("Repositories.GetArchiveLink returned error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Repositories.GetArchiveLink returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}
	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Repositories.GetArchiveLink returned %+v, want %+v", url.String(), want)
	}
}

func TestRepositoriesService_GetContents_NoTrailingSlashInDirectoryApiPath(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/.github", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		query := r.URL.Query()
		if query.Get("ref") != "mybranch" {
			t.Errorf("Repositories.GetContents returned %+v, want %+v", query.Get("ref"), "mybranch")
		}
		fmt.Fprint(w, `{}`)
	})
	ctx := context.Background()
	_, _, _, err := client.Repositories.GetContents(ctx, "o", "r", ".github/", &RepositoryContentGetOptions{
		Ref: "mybranch",
	})
	if err != nil {
		t.Fatalf("Repositories.GetContents returned error: %v", err)
	}
}

func TestRepositoryContent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryContent{}, "{}")

	r := &RepositoryContent{
		Type:            String("type"),
		Target:          String("target"),
		Encoding:        String("encoding"),
		Size:            Int(1),
		Name:            String("name"),
		Path:            String("path"),
		Content:         String("content"),
		SHA:             String("sha"),
		URL:             String("url"),
		GitURL:          String("gurl"),
		HTMLURL:         String("hurl"),
		DownloadURL:     String("durl"),
		SubmoduleGitURL: String("smgurl"),
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
	testJSONMarshal(t, &RepositoryContentResponse{}, "{}")

	r := &RepositoryContentResponse{
		Content: &RepositoryContent{
			Type:            String("type"),
			Target:          String("target"),
			Encoding:        String("encoding"),
			Size:            Int(1),
			Name:            String("name"),
			Path:            String("path"),
			Content:         String("content"),
			SHA:             String("sha"),
			URL:             String("url"),
			GitURL:          String("gurl"),
			HTMLURL:         String("hurl"),
			DownloadURL:     String("durl"),
			SubmoduleGitURL: String("smgurl"),
		},
		Commit: Commit{
			SHA: String("s"),
			Author: &CommitAuthor{
				Date:  &Timestamp{referenceTime},
				Name:  String("n"),
				Email: String("e"),
				Login: String("u"),
			},
			Committer: &CommitAuthor{
				Date:  &Timestamp{referenceTime},
				Name:  String("n"),
				Email: String("e"),
				Login: String("u"),
			},
			Message: String("m"),
			Tree: &Tree{
				SHA: String("s"),
				Entries: []*TreeEntry{{
					SHA:     String("s"),
					Path:    String("p"),
					Mode:    String("m"),
					Type:    String("t"),
					Size:    Int(1),
					Content: String("c"),
					URL:     String("u"),
				}},
				Truncated: Bool(false),
			},
			Parents: nil,
			Stats: &CommitStats{
				Additions: Int(1),
				Deletions: Int(1),
				Total:     Int(1),
			},
			HTMLURL: String("h"),
			URL:     String("u"),
			Verification: &SignatureVerification{
				Verified:  Bool(false),
				Reason:    String("r"),
				Signature: String("s"),
				Payload:   String("p"),
			},
			NodeID:       String("n"),
			CommentCount: Int(1),
			SigningKey:   &openpgp.Entity{},
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
			"stats": {
				"additions": 1,
				"deletions": 1,
				"total": 1
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
	testJSONMarshal(t, &RepositoryContentFileOptions{}, "{}")

	r := &RepositoryContentFileOptions{
		Message: String("type"),
		Content: []byte{1},
		SHA:     String("type"),
		Branch:  String("type"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
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
