// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListReleases(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	releases, _, err := client.Repositories.ListReleases(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListReleases returned error: %v", err)
	}
	want := []*RepositoryRelease{{ID: Int64(1)}}
	if !cmp.Equal(releases, want) {
		t.Errorf("Repositories.ListReleases returned %+v, want %+v", releases, want)
	}

	const methodName = "ListReleases"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListReleases(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListReleases(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GenerateReleaseNotes(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/generate-notes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"tag_name":"v1.0.0"}`+"\n")
		fmt.Fprint(w, `{"name":"v1.0.0","body":"**Full Changelog**: https://github.com/o/r/compare/v0.9.0...v1.0.0"}`)
	})

	opt := &GenerateNotesOptions{
		TagName: "v1.0.0",
	}
	ctx := context.Background()
	releases, _, err := client.Repositories.GenerateReleaseNotes(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.GenerateReleaseNotes returned error: %v", err)
	}
	want := &RepositoryReleaseNotes{
		Name: "v1.0.0",
		Body: "**Full Changelog**: https://github.com/o/r/compare/v0.9.0...v1.0.0",
	}
	if !cmp.Equal(releases, want) {
		t.Errorf("Repositories.GenerateReleaseNotes returned %+v, want %+v", releases, want)
	}

	const methodName = "GenerateReleaseNotes"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GenerateReleaseNotes(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GenerateReleaseNotes(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRelease(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"author":{"login":"l"}}`)
	})

	ctx := context.Background()
	release, resp, err := client.Repositories.GetRelease(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetRelease returned error: %v\n%v", err, resp.Body)
	}

	want := &RepositoryRelease{ID: Int64(1), Author: &User{Login: String("l")}}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.GetRelease returned %+v, want %+v", release, want)
	}

	const methodName = "GetRelease"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetRelease(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRelease(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetLatestRelease(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3}`)
	})

	ctx := context.Background()
	release, resp, err := client.Repositories.GetLatestRelease(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetLatestRelease returned error: %v\n%v", err, resp.Body)
	}

	want := &RepositoryRelease{ID: Int64(3)}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.GetLatestRelease returned %+v, want %+v", release, want)
	}

	const methodName = "GetLatestRelease"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetLatestRelease(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetLatestRelease(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetReleaseByTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/tags/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":13}`)
	})

	ctx := context.Background()
	release, resp, err := client.Repositories.GetReleaseByTag(ctx, "o", "r", "foo")
	if err != nil {
		t.Errorf("Repositories.GetReleaseByTag returned error: %v\n%v", err, resp.Body)
	}

	want := &RepositoryRelease{ID: Int64(13)}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.GetReleaseByTag returned %+v, want %+v", release, want)
	}

	const methodName = "GetReleaseByTag"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetReleaseByTag(ctx, "\n", "\n", "foo")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetReleaseByTag(ctx, "o", "r", "foo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateRelease(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RepositoryRelease{
		Name:                   String("v1.0"),
		DiscussionCategoryName: String("General"),
		GenerateReleaseNotes:   Bool(true),
		// Fields to be removed:
		ID:          Int64(2),
		CreatedAt:   &Timestamp{referenceTime},
		PublishedAt: &Timestamp{referenceTime},
		URL:         String("http://url/"),
		HTMLURL:     String("http://htmlurl/"),
		AssetsURL:   String("http://assetsurl/"),
		Assets:      []*ReleaseAsset{{ID: Int64(5)}},
		UploadURL:   String("http://uploadurl/"),
		ZipballURL:  String("http://zipballurl/"),
		TarballURL:  String("http://tarballurl/"),
		Author:      &User{Name: String("octocat")},
		NodeID:      String("nodeid"),
	}

	mux.HandleFunc("/repos/o/r/releases", func(w http.ResponseWriter, r *http.Request) {
		v := new(repositoryReleaseRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &repositoryReleaseRequest{
			Name:                   String("v1.0"),
			DiscussionCategoryName: String("General"),
			GenerateReleaseNotes:   Bool(true),
		}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	release, _, err := client.Repositories.CreateRelease(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateRelease returned error: %v", err)
	}

	want := &RepositoryRelease{ID: Int64(1)}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.CreateRelease returned %+v, want %+v", release, want)
	}

	const methodName = "CreateRelease"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateRelease(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRelease(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EditRelease(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RepositoryRelease{
		Name:                   String("n"),
		DiscussionCategoryName: String("General"),
		GenerateReleaseNotes:   Bool(true),
		// Fields to be removed:
		ID:          Int64(2),
		CreatedAt:   &Timestamp{referenceTime},
		PublishedAt: &Timestamp{referenceTime},
		URL:         String("http://url/"),
		HTMLURL:     String("http://htmlurl/"),
		AssetsURL:   String("http://assetsurl/"),
		Assets:      []*ReleaseAsset{{ID: Int64(5)}},
		UploadURL:   String("http://uploadurl/"),
		ZipballURL:  String("http://zipballurl/"),
		TarballURL:  String("http://tarballurl/"),
		Author:      &User{Name: String("octocat")},
		NodeID:      String("nodeid"),
	}

	mux.HandleFunc("/repos/o/r/releases/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(repositoryReleaseRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		want := &repositoryReleaseRequest{
			Name:                   String("n"),
			DiscussionCategoryName: String("General"),
			GenerateReleaseNotes:   Bool(true),
		}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	release, _, err := client.Repositories.EditRelease(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditRelease returned error: %v", err)
	}
	want := &RepositoryRelease{ID: Int64(1)}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.EditRelease returned = %+v, want %+v", release, want)
	}

	const methodName = "EditRelease"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditRelease(ctx, "\n", "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditRelease(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteRelease(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteRelease(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteRelease returned error: %v", err)
	}

	const methodName = "DeleteRelease"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteRelease(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteRelease(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_ListReleaseAssets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/1/assets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Repositories.ListReleaseAssets returned error: %v", err)
	}
	want := []*ReleaseAsset{{ID: Int64(1)}}
	if !cmp.Equal(assets, want) {
		t.Errorf("Repositories.ListReleaseAssets returned %+v, want %+v", assets, want)
	}

	const methodName = "ListReleaseAssets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListReleaseAssets(ctx, "\n", "\n", 1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListReleaseAssets(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetReleaseAsset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	asset, _, err := client.Repositories.GetReleaseAsset(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetReleaseAsset returned error: %v", err)
	}
	want := &ReleaseAsset{ID: Int64(1)}
	if !cmp.Equal(asset, want) {
		t.Errorf("Repositories.GetReleaseAsset returned %+v, want %+v", asset, want)
	}

	const methodName = "GetReleaseAsset"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetReleaseAsset(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetReleaseAsset(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DownloadReleaseAsset_Stream(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=hello-world.txt")
		fmt.Fprint(w, "Hello World")
	})

	ctx := context.Background()
	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("Repositories.DownloadReleaseAsset returned error: %v", err)
	}
	want := []byte("Hello World")
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("Repositories.DownloadReleaseAsset returned bad reader: %v", err)
	}
	if !bytes.Equal(want, content) {
		t.Errorf("Repositories.DownloadReleaseAsset returned %+v, want %+v", content, want)
	}

	const methodName = "DownloadReleaseAsset"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.DownloadReleaseAsset(ctx, "\n", "\n", -1, nil)
		return err
	})
}

func TestRepositoriesService_DownloadReleaseAsset_Redirect(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		http.Redirect(w, r, "/yo", http.StatusFound)
	})

	ctx := context.Background()
	_, got, err := client.Repositories.DownloadReleaseAsset(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("Repositories.DownloadReleaseAsset returned error: %v", err)
	}
	want := "/yo"
	if !strings.HasSuffix(got, want) {
		t.Errorf("Repositories.DownloadReleaseAsset returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_DownloadReleaseAsset_FollowRedirect(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		// /yo, below will be served as baseURLPath/yo
		http.Redirect(w, r, baseURLPath+"/yo", http.StatusFound)
	})
	mux.HandleFunc("/yo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "*/*")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=hello-world.txt")
		fmt.Fprint(w, "Hello World")
	})

	ctx := context.Background()
	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, "o", "r", 1, http.DefaultClient)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("Repositories.DownloadReleaseAsset returned error: %v", err)
	}
	reader.Close()
	want := []byte("Hello World")
	if !bytes.Equal(want, content) {
		t.Errorf("Repositories.DownloadReleaseAsset returned %+v, want %+v", content, want)
	}
}

func TestRepositoriesService_DownloadReleaseAsset_APIError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"Not Found","documentation_url":"https://developer.github.com/v3"}`)
	})

	ctx := context.Background()
	resp, loc, err := client.Repositories.DownloadReleaseAsset(ctx, "o", "r", 1, nil)
	if err == nil {
		t.Error("Repositories.DownloadReleaseAsset did not return an error")
	}

	if resp != nil {
		resp.Close()
		t.Error("Repositories.DownloadReleaseAsset returned stream, want nil")
	}

	if loc != "" {
		t.Errorf(`Repositories.DownloadReleaseAsset returned "%s", want empty ""`, loc)
	}
}

func TestRepositoriesService_EditReleaseAsset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ReleaseAsset{Name: String("n")}

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(ReleaseAsset)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	asset, _, err := client.Repositories.EditReleaseAsset(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditReleaseAsset returned error: %v", err)
	}
	want := &ReleaseAsset{ID: Int64(1)}
	if !cmp.Equal(asset, want) {
		t.Errorf("Repositories.EditReleaseAsset returned = %+v, want %+v", asset, want)
	}

	const methodName = "EditReleaseAsset"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditReleaseAsset(ctx, "\n", "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditReleaseAsset(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteReleaseAsset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteReleaseAsset(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteReleaseAsset returned error: %v", err)
	}

	const methodName = "DeleteReleaseAsset"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteReleaseAsset(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteReleaseAsset(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_UploadReleaseAsset(t *testing.T) {
	var (
		defaultUploadOptions     = &UploadOptions{Name: "n"}
		defaultExpectedFormValue = values{"name": "n"}
		mediaTypeTextPlain       = "text/plain; charset=utf-8"
	)
	uploadTests := []struct {
		uploadOpts         *UploadOptions
		fileName           string
		expectedFormValues values
		expectedMediaType  string
	}{
		// No file extension and no explicit media type.
		{
			defaultUploadOptions,
			"upload",
			defaultExpectedFormValue,
			defaultMediaType,
		},
		// File extension and no explicit media type.
		{
			defaultUploadOptions,
			"upload.txt",
			defaultExpectedFormValue,
			mediaTypeTextPlain,
		},
		// No file extension and explicit media type.
		{
			&UploadOptions{Name: "n", MediaType: "image/png"},
			"upload",
			defaultExpectedFormValue,
			"image/png",
		},
		// File extension and explicit media type.
		{
			&UploadOptions{Name: "n", MediaType: "image/png"},
			"upload.png",
			defaultExpectedFormValue,
			"image/png",
		},
		// Label provided.
		{
			&UploadOptions{Name: "n", Label: "l"},
			"upload.txt",
			values{"name": "n", "label": "l"},
			mediaTypeTextPlain,
		},
		// No label provided.
		{
			defaultUploadOptions,
			"upload.txt",
			defaultExpectedFormValue,
			mediaTypeTextPlain,
		},
	}

	client, mux, _, teardown := setup()
	defer teardown()

	for key, test := range uploadTests {
		releaseEndpoint := fmt.Sprintf("/repos/o/r/releases/%d/assets", key)
		mux.HandleFunc(releaseEndpoint, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "Content-Type", test.expectedMediaType)
			testHeader(t, r, "Content-Length", "12")
			testFormValues(t, r, test.expectedFormValues)
			testBody(t, r, "Upload me !\n")

			fmt.Fprintf(w, `{"id":1}`)
		})

		file, dir, err := openTestFile(test.fileName, "Upload me !\n")
		if err != nil {
			t.Fatalf("Unable to create temp file: %v", err)
		}
		defer os.RemoveAll(dir)

		ctx := context.Background()
		asset, _, err := client.Repositories.UploadReleaseAsset(ctx, "o", "r", int64(key), test.uploadOpts, file)
		if err != nil {
			t.Errorf("Repositories.UploadReleaseAssert returned error: %v", err)
		}
		want := &ReleaseAsset{ID: Int64(1)}
		if !cmp.Equal(asset, want) {
			t.Errorf("Repositories.UploadReleaseAssert returned %+v, want %+v", asset, want)
		}

		const methodName = "UploadReleaseAsset"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.Repositories.UploadReleaseAsset(ctx, "\n", "\n", int64(key), test.uploadOpts, file)
			return err
		})
	}
}

func TestRepositoryReleaseRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &repositoryReleaseRequest{}, "{}")

	u := &repositoryReleaseRequest{
		TagName:                String("tn"),
		TargetCommitish:        String("tc"),
		Name:                   String("name"),
		Body:                   String("body"),
		Draft:                  Bool(false),
		Prerelease:             Bool(false),
		DiscussionCategoryName: String("dcn"),
	}

	want := `{
		"tag_name": "tn",
		"target_commitish": "tc",
		"name": "name",
		"body": "body",
		"draft": false,
		"prerelease": false,
		"discussion_category_name": "dcn"
	}`

	testJSONMarshal(t, u, want)
}

func TestReleaseAsset_Marshal(t *testing.T) {
	testJSONMarshal(t, &ReleaseAsset{}, "{}")

	u := &ReleaseAsset{
		ID:                 Int64(1),
		URL:                String("url"),
		Name:               String("name"),
		Label:              String("label"),
		State:              String("state"),
		ContentType:        String("ct"),
		Size:               Int(1),
		DownloadCount:      Int(1),
		CreatedAt:          &Timestamp{referenceTime},
		UpdatedAt:          &Timestamp{referenceTime},
		BrowserDownloadURL: String("bdu"),
		Uploader:           &User{ID: Int64(1)},
		NodeID:             String("nid"),
	}

	want := `{
		"id": 1,
		"url": "url",
		"name": "name",
		"label": "label",
		"state": "state",
		"content_type": "ct",
		"size": 1,
		"download_count": 1,
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"browser_download_url": "bdu",
		"uploader": {
			"id": 1
		},
		"node_id": "nid"
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoryRelease_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryRelease{}, "{}")

	u := &RepositoryRelease{
		TagName:                String("tn"),
		TargetCommitish:        String("tc"),
		Name:                   String("name"),
		Body:                   String("body"),
		Draft:                  Bool(false),
		Prerelease:             Bool(false),
		DiscussionCategoryName: String("dcn"),
		ID:                     Int64(1),
		CreatedAt:              &Timestamp{referenceTime},
		PublishedAt:            &Timestamp{referenceTime},
		URL:                    String("url"),
		HTMLURL:                String("hurl"),
		AssetsURL:              String("aurl"),
		Assets:                 []*ReleaseAsset{{ID: Int64(1)}},
		UploadURL:              String("uurl"),
		ZipballURL:             String("zurl"),
		TarballURL:             String("turl"),
		Author:                 &User{ID: Int64(1)},
		NodeID:                 String("nid"),
	}

	want := `{
		"tag_name": "tn",
		"target_commitish": "tc",
		"name": "name",
		"body": "body",
		"draft": false,
		"prerelease": false,
		"discussion_category_name": "dcn",
		"id": 1,
		"created_at": ` + referenceTimeStr + `,
		"published_at": ` + referenceTimeStr + `,
		"url": "url",
		"html_url": "hurl",
		"assets_url": "aurl",
		"assets": [
			{
				"id": 1
			}
		],
		"upload_url": "uurl",
		"zipball_url": "zurl",
		"tarball_url": "turl",
		"author": {
			"id": 1
		},
		"node_id": "nid"
	}`

	testJSONMarshal(t, u, want)
}

func TestGenerateNotesOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &GenerateNotesOptions{}, "{}")

	u := &GenerateNotesOptions{
		TagName:         "tag_name",
		PreviousTagName: String("previous_tag_name"),
		TargetCommitish: String("target_commitish"),
	}

	want := `{
		"tag_name":          "tag_name",
		"previous_tag_name": "previous_tag_name",
		"target_commitish":  "target_commitish"
	}`

	testJSONMarshal(t, u, want)
}
