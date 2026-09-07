// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_EnablePagesLegacy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Pages{
		BuildType: new("legacy"),
		Source: &PagesSource{
			Branch: new("master"),
			Path:   new("/"),
		},
		CNAME: new("www.example.com"), // not passed along.
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeEnablePagesAPIPreview)
		want := &createPagesRequest{BuildType: new("legacy"), Source: &PagesSource{Branch: new("master"), Path: new("/")}}
		testJSONBody(t, r, want)

		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false,"html_url":"h","build_type": "legacy","source": {"branch":"master", "path":"/"}}`)
	})

	ctx := t.Context()
	page, _, err := client.Repositories.EnablePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.EnablePages returned error: %v", err)
	}

	want := &Pages{URL: new("u"), Status: new("s"), CNAME: new("c"), Custom404: new(false), HTMLURL: new("h"), BuildType: new("legacy"), Source: &PagesSource{Branch: new("master"), Path: new("/")}}

	if !cmp.Equal(page, want) {
		t.Errorf("Repositories.EnablePages returned %v, want %v", page, want)
	}

	const methodName = "EnablePages"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EnablePages(ctx, "o", "r", nil)
		return err
	})
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

func TestRepositoriesService_EnablePagesWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Pages{
		BuildType: new("workflow"),
		CNAME:     new("www.example.com"), // not passed along.
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeEnablePagesAPIPreview)
		want := &createPagesRequest{BuildType: new("workflow")}
		testJSONBody(t, r, want)
		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false,"html_url":"h","build_type": "workflow"}`)
	})

	ctx := t.Context()
	page, _, err := client.Repositories.EnablePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.EnablePages returned error: %v", err)
	}

	want := &Pages{URL: new("u"), Status: new("s"), CNAME: new("c"), Custom404: new(false), HTMLURL: new("h"), BuildType: new("workflow")}

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

func TestRepositoriesService_UpdatePagesLegacy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PagesUpdate{
		CNAME:     new("www.example.com"),
		BuildType: new("legacy"),
		Source:    &PagesSource{Branch: new("gh-pages")},
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"cname":"www.example.com","build_type":"legacy","source":{"branch":"gh-pages"}}`)
	})

	ctx := t.Context()
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

func TestRepositoriesService_UpdatePagesWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PagesUpdate{
		CNAME:     new("www.example.com"),
		BuildType: new("workflow"),
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"cname":"www.example.com","build_type":"workflow"}`)
	})

	ctx := t.Context()
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

func TestRepositoriesService_UpdatePagesGHES(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PagesUpdateWithoutCNAME{
		BuildType: new("workflow"),
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"build_type":"workflow"}`)
	})

	ctx := t.Context()
	_, err := client.Repositories.UpdatePagesGHES(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.UpdatePagesGHES returned error: %v", err)
	}

	const methodName = "UpdatePagesGHES"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.UpdatePagesGHES(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.UpdatePagesGHES(ctx, "o", "r", input)
	})
}

func TestRepositoriesService_UpdatePages_NullCNAME(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PagesUpdate{
		Source: &PagesSource{Branch: new("gh-pages")},
	}

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"cname":null,"source":{"branch":"gh-pages"}}`)
	})

	ctx := t.Context()
	_, err := client.Repositories.UpdatePages(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.UpdatePages returned error: %v", err)
	}
}

func TestRepositoriesService_DisablePages(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeEnablePagesAPIPreview)
	})

	ctx := t.Context()
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false,"html_url":"h","public":true, "https_certificate": {"state":"approved","description": "Certificate is approved","domains": ["developer.github.com"],"expires_at": "2021-05-22"},"https_enforced": true}`)
	})

	ctx := t.Context()
	page, _, err := client.Repositories.GetPagesInfo(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetPagesInfo returned error: %v", err)
	}

	want := &Pages{URL: new("u"), Status: new("s"), CNAME: new("c"), Custom404: new(false), HTMLURL: new("h"), Public: new(true), HTTPSCertificate: &PagesHTTPSCertificate{State: new("approved"), Description: new("Certificate is approved"), Domains: []string{"developer.github.com"}, ExpiresAt: new("2021-05-22")}, HTTPSEnforced: new(true)}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"url":"u","status":"s","commit":"c"}]`)
	})

	ctx := t.Context()
	pages, _, err := client.Repositories.ListPagesBuilds(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("Repositories.ListPagesBuilds returned error: %v", err)
	}

	want := []*PagesBuild{{URL: new("u"), Status: new("s"), Commit: new("c")}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := t.Context()
	_, _, err := client.Repositories.ListPagesBuilds(ctx, "o", "r", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Repositories.ListPagesBuilds returned error: %v", err)
	}
}

func TestRepositoriesService_GetLatestPagesBuild(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/builds/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","commit":"c"}`)
	})

	ctx := t.Context()
	build, _, err := client.Repositories.GetLatestPagesBuild(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetLatestPagesBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: new("u"), Status: new("s"), Commit: new("c")}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/builds/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","commit":"c"}`)
	})

	ctx := t.Context()
	build, _, err := client.Repositories.GetPageBuild(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetPageBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: new("u"), Status: new("s"), Commit: new("c")}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"url":"u","status":"s"}`)
	})

	ctx := t.Context()
	build, _, err := client.Repositories.RequestPageBuild(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.RequestPageBuild returned error: %v", err)
	}

	want := &PagesBuild{URL: new("u"), Status: new("s")}
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

func TestRepositoriesService_GetPageHealthCheck(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pages/health", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"domain":{"host":"example.com","uri":"http://example.com/","nameservers":"default","dns_resolves":true},"alt_domain":{"host":"www.example.com","uri":"http://www.example.com/","nameservers":"default","dns_resolves":true}}`)
	})

	ctx := t.Context()
	healthCheckResponse, _, err := client.Repositories.GetPageHealthCheck(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetPageHealthCheck returned error: %v", err)
	}

	want := &PagesHealthCheckResponse{
		Domain: &PagesDomain{
			Host:        new("example.com"),
			URI:         new("http://example.com/"),
			Nameservers: new("default"),
			DNSResolves: new(true),
		},
		AltDomain: &PagesDomain{
			Host:        new("www.example.com"),
			URI:         new("http://www.example.com/"),
			Nameservers: new("default"),
			DNSResolves: new(true),
		},
	}
	if !cmp.Equal(healthCheckResponse, want) {
		t.Errorf("Repositories.GetPageHealthCheck returned %+v, want %+v", healthCheckResponse, want)
	}

	const methodName = "GetPageHealthCheck"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetPageHealthCheck(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetPageHealthCheck(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
