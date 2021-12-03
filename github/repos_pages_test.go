// Copyright 2014 The go-github AUTHORS. All rights reserved.
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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_EnablePages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Pages{
		Source: &PagesSource{
			Branch: String("master"),
			Path:   String("/"),
		},
		CNAME: String("www.my-domain.com"), // not passed along.
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		v := new(createPagesRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeEnablePagesAPIPreview)
		want := &createPagesRequest{Source: &PagesSource{Branch: String("master"), Path: String("/")}}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false,"html_url":"h", "source": {"branch":"master", "path":"/"}}`)
	})

	ctx := context.Background()
	page, _, err := client.Repositories.EnablePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.EnablePages returned error: %v", err)
	}

	want := &Pages{URL: String("u"), Status: String("s"), CNAME: String("c"), Custom404: Bool(false), HTMLURL: String("h"), Source: &PagesSource{Branch: String("master"), Path: String("/")}}

	if !cmp.Equal(page, want) {
		t.Errorf("Repositories.EnablePages returned %v, want %v", page, want)
	}

	const methodName = "EnablePages"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EnablePages(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EnablePages(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdatePages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PagesUpdate{
		CNAME:  String("www.my-domain.com"),
		Source: String("gh-pages"),
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		v := new(PagesUpdate)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		want := &PagesUpdate{CNAME: String("www.my-domain.com"), Source: String("gh-pages")}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"cname":"www.my-domain.com","source":"gh-pages"}`)
	})

	ctx := context.Background()
	_, err := client.Repositories.UpdatePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.UpdatePages returned error: %v", err)
	}

	const methodName = "UpdatePages"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.UpdatePages(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.UpdatePages(ctx, "o", "r", input)
	})
}

func TestRepositoriesService_UpdatePages_NullCNAME(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PagesUpdate{
		Source: String("gh-pages"),
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unable to read body: %v", err)
		}

		want := []byte(`{"cname":null,"source":"gh-pages"}` + "\n")
		if !bytes.Equal(got, want) {
			t.Errorf("Request body = %+v, want %+v", got, want)
		}

		fmt.Fprint(w, `{"cname":null,"source":"gh-pages"}`)
	})

	ctx := context.Background()
	_, err := client.Repositories.UpdatePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.UpdatePages returned error: %v", err)
	}
}

func TestRepositoriesService_DisablePages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeEnablePagesAPIPreview)
	})

	ctx := context.Background()
	_, err := client.Repositories.DisablePages(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.DisablePages returned error: %v", err)
	}

	const methodName = "DisablePages"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DisablePages(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DisablePages(ctx, "o", "r")
	})
}

func TestRepositoriesService_GetPagesInfo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false,"html_url":"h","public":true, "https_certificate": {"state":"approved","description": "Certificate is approved","domains": ["developer.github.com"],"expires_at": "2021-05-22"},"https_enforced": true}`)
	})

	ctx := context.Background()
	page, _, err := client.Repositories.GetPagesInfo(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetPagesInfo returned error: %v", err)
	}

	want := &Pages{URL: String("u"), Status: String("s"), CNAME: String("c"), Custom404: Bool(false), HTMLURL: String("h"), Public: Bool(true), HTTPSCertificate: &PagesHTTPSCertificate{State: String("approved"), Description: String("Certificate is approved"), Domains: []string{"developer.github.com"}, ExpiresAt: String("2021-05-22")}, HTTPSEnforced: Bool(true)}
	if !cmp.Equal(page, want) {
		t.Errorf("Repositories.GetPagesInfo returned %+v, want %+v", page, want)
	}

	const methodName = "GetPagesInfo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetPagesInfo(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetPagesInfo(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListPagesBuilds(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"url":"u","status":"s","commit":"c"}]`)
	})

	ctx := context.Background()
	pages, _, err := client.Repositories.ListPagesBuilds(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("Repositories.ListPagesBuilds returned error: %v", err)
	}

	want := []*PagesBuild{{URL: String("u"), Status: String("s"), Commit: String("c")}}
	if !cmp.Equal(pages, want) {
		t.Errorf("Repositories.ListPagesBuilds returned %+v, want %+v", pages, want)
	}

	const methodName = "ListPagesBuilds"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListPagesBuilds(ctx, "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListPagesBuilds(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListPagesBuilds_withOptions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, err := client.Repositories.ListPagesBuilds(ctx, "o", "r", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Repositories.ListPagesBuilds returned error: %v", err)
	}
}

func TestRepositoriesService_GetLatestPagesBuild(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","commit":"c"}`)
	})

	ctx := context.Background()
	build, _, err := client.Repositories.GetLatestPagesBuild(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetLatestPagesBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: String("u"), Status: String("s"), Commit: String("c")}
	if !cmp.Equal(build, want) {
		t.Errorf("Repositories.GetLatestPagesBuild returned %+v, want %+v", build, want)
	}

	const methodName = "GetLatestPagesBuild"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetLatestPagesBuild(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetLatestPagesBuild(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetPageBuild(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","commit":"c"}`)
	})

	ctx := context.Background()
	build, _, err := client.Repositories.GetPageBuild(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetPageBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: String("u"), Status: String("s"), Commit: String("c")}
	if !cmp.Equal(build, want) {
		t.Errorf("Repositories.GetPageBuild returned %+v, want %+v", build, want)
	}

	const methodName = "GetPageBuild"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetPageBuild(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetPageBuild(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RequestPageBuild(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"url":"u","status":"s"}`)
	})

	ctx := context.Background()
	build, _, err := client.Repositories.RequestPageBuild(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.RequestPageBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: String("u"), Status: String("s")}
	if !cmp.Equal(build, want) {
		t.Errorf("Repositories.RequestPageBuild returned %+v, want %+v", build, want)
	}

	const methodName = "RequestPageBuild"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.RequestPageBuild(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.RequestPageBuild(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPagesSource_Marshal(t *testing.T) {
	testJSONMarshal(t, &PagesSource{}, "{}")

	u := &PagesSource{
		Branch: String("branch"),
		Path:   String("path"),
	}

	want := `{
		"branch": "branch",
		"path": "path"
	}`

	testJSONMarshal(t, u, want)
}

func TestPagesError_Marshal(t *testing.T) {
	testJSONMarshal(t, &PagesError{}, "{}")

	u := &PagesError{
		Message: String("message"),
	}

	want := `{
		"message": "message"
	}`

	testJSONMarshal(t, u, want)
}

func TestPagesUpdate_Marshal(t *testing.T) {
	testJSONMarshal(t, &PagesUpdate{}, "{}")

	u := &PagesUpdate{
		CNAME:  String("cname"),
		Source: String("src"),
	}

	want := `{
		"cname": "cname",
		"source": "src"
	}`

	testJSONMarshal(t, u, want)
}

func TestPages_Marshal(t *testing.T) {
	testJSONMarshal(t, &Pages{}, "{}")

	u := &Pages{
		URL:       String("url"),
		Status:    String("status"),
		CNAME:     String("cname"),
		Custom404: Bool(false),
		HTMLURL:   String("hurl"),
		Source: &PagesSource{
			Branch: String("branch"),
			Path:   String("path"),
		},
	}

	want := `{
		"url": "url",
		"status": "status",
		"cname": "cname",
		"custom_404": false,
		"html_url": "hurl",
		"source": {
			"branch": "branch",
			"path": "path"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPagesBuild_Marshal(t *testing.T) {
	testJSONMarshal(t, &PagesBuild{}, "{}")

	u := &PagesBuild{
		URL:    String("url"),
		Status: String("status"),
		Error: &PagesError{
			Message: String("message"),
		},
		Pusher:    &User{ID: Int64(1)},
		Commit:    String("commit"),
		Duration:  Int(1),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
	}

	want := `{
		"url": "url",
		"status": "status",
		"error": {
			"message": "message"
		},
		"pusher": {
			"id": 1
		},
		"commit": "commit",
		"duration": 1,
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestCreatePagesRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &createPagesRequest{}, "{}")

	u := &createPagesRequest{
		Source: &PagesSource{
			Branch: String("branch"),
			Path:   String("path"),
		},
	}

	want := `{
		"source": {
			"branch": "branch",
			"path": "path"
		}
	}`

	testJSONMarshal(t, u, want)
}
