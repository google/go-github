// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoriesService_List_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview}
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", got, want)
	}

	// Test addOptions failure
	_, _, err = client.Repositories.List(ctx, "\n", &RepositoryListOptions{})
	if err == nil {
		t.Error("bad options List err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err := client.Repositories.List(ctx, "", nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now List = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now List resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now List err = nil, want error")
	}
}

func TestRepositoriesService_List_specifiedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview}
	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{
			"visibility":  "public",
			"affiliation": "owner,collaborator",
			"sort":        "created",
			"direction":   "asc",
			"page":        "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
		ListOptions: ListOptions{Page: 2},
	}
	repos, _, err := client.Repositories.List(context.Background(), "u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_specifiedUser_type(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview}
	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{
			"type": "owner",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListOptions{
		Type: "owner",
	}
	repos, _, err := client.Repositories.List(context.Background(), "u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Repositories.List(context.Background(), "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListByOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview}
	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{
			"type": "forks",
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	opt := &RepositoryListByOrgOptions{
		Type:        "forks",
		ListOptions: ListOptions{Page: 2},
	}
	got, _, err := client.Repositories.ListByOrg(ctx, "o", opt)
	if err != nil {
		t.Errorf("Repositories.ListByOrg returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ListByOrg returned %+v, want %+v", got, want)
	}

	// Test addOptions failure
	_, _, err = client.Repositories.ListByOrg(ctx, "\n", opt)
	if err == nil {
		t.Error("bad options ListByOrg err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err := client.Repositories.ListByOrg(ctx, "o", opt)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListByOrg = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListByOrg resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListByOrg err = nil, want error")
	}
}

func TestRepositoriesService_ListByOrg_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Repositories.ListByOrg(context.Background(), "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "1",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	opt := &RepositoryListAllOptions{1}
	got, _, err := client.Repositories.ListAll(ctx, opt)
	if err != nil {
		t.Errorf("Repositories.ListAll returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ListAll returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.ListAll(ctx, &RepositoryListAllOptions{1})
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListAll = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListAll resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListAll err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.ListAll(ctx, &RepositoryListAllOptions{1})
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListAll = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListAll resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListAll err = nil, want error")
	}
}

func TestRepositoriesService_Create_user(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Repository{
		Name:     String("n"),
		Archived: Bool(true), // not passed along.
	}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(createRepoRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		want := &createRepoRequest{Name: String("n")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Create(ctx, "", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.Create(ctx, "", input)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' Create = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' Create resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' Create err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.Create(ctx, "", input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now Create = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now Create resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now Create err = nil, want error")
	}
}

func TestRepositoriesService_Create_org(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Repository{
		Name:     String("n"),
		Archived: Bool(true), // not passed along.
	}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(createRepoRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		want := &createRepoRequest{Name: String("n")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Repositories.Create(context.Background(), "o", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_CreateFromTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	templateRepoReq := &TemplateRepoRequest{
		Name: String("n"),
	}

	mux.HandleFunc("/repos/to/tr/generate", func(w http.ResponseWriter, r *http.Request) {
		v := new(TemplateRepoRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeRepositoryTemplatePreview)
		want := &TemplateRepoRequest{Name: String("n")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1,"name":"n"}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.CreateFromTemplate(ctx, "to", "tr", templateRepoReq)
	if err != nil {
		t.Errorf("Repositories.CreateFromTemplate returned error: %v", err)
	}

	want := &Repository{ID: Int64(1), Name: String("n")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.CreateFromTemplate returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.CreateFromTemplate(ctx, "to", "tr", templateRepoReq)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' CreateFromTemplate = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' CreateFromTemplate resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' CreateFromTemplate err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.CreateFromTemplate(ctx, "to", "tr", templateRepoReq)
	if got != nil {
		t.Errorf("rate.Reset.Time > now CreateFromTemplate = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now CreateFromTemplate resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now CreateFromTemplate err = nil, want error")
	}
}

func TestRepositoriesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview, mediaTypeRepositoryTemplatePreview}
	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Get(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.Get returned error: %v", err)
	}

	want := &Repository{ID: Int64(1), Name: String("n"), Description: String("d"), Owner: &User{Login: String("l")}, License: &License{Key: String("mit")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.Get returned %+v, want %+v", got, want)
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err := client.Repositories.Get(ctx, "o", "r")
	if got != nil {
		t.Errorf("rate.Reset.Time > now Get = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now Get resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now Get err = nil, want error")
	}
}

func TestRepositoriesService_GetCodeOfConduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/community/code_of_conduct", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `{
						"key": "key",
						"name": "name",
						"url": "url",
						"body": "body"}`,
		)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCodeOfConduct(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCodeOfConduct returned error: %v", err)
	}

	want := &CodeOfConduct{
		Key:  String("key"),
		Name: String("name"),
		URL:  String("url"),
		Body: String("body"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.GetCodeOfConduct returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.GetCodeOfConduct(ctx, "o", "r")
	if got != nil {
		t.Errorf("client.BaseURL.Path='' GetCodeOfConduct = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' GetCodeOfConduct resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' GetCodeOfConduct err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.GetCodeOfConduct(ctx, "o", "r")
	if got != nil {
		t.Errorf("rate.Reset.Time > now GetCodeOfConduct = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now GetCodeOfConduct resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now GetCodeOfConduct err = nil, want error")
	}
}

func TestRepositoriesService_GetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("Repositories.GetByID returned error: %v", err)
	}

	want := &Repository{ID: Int64(1), Name: String("n"), Description: String("d"), Owner: &User{Login: String("l")}, License: &License{Key: String("mit")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.GetByID returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.GetByID(ctx, 1)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' GetByID = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' GetByID resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' GetByID err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.GetByID(ctx, 1)
	if got != nil {
		t.Errorf("rate.Reset.Time > now GetByID = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now GetByID resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now GetByID err = nil, want error")
	}
}

func TestRepositoriesService_Edit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	i := true
	input := &Repository{HasIssues: &i}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Edit(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.Edit returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.Edit returned %+v, want %+v", got, want)
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err := client.Repositories.Edit(ctx, "o", "r", input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now Edit = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now Edit resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now Edit err = nil, want error")
	}
}

func TestRepositoriesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.Delete(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.Delete returned error: %v", err)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	resp, err := client.Repositories.Delete(ctx, "o", "r")
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' Delete resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' Delete err = nil, want error")
	}
}

func TestRepositoriesService_Get_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Repositories.Get(context.Background(), "%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_Edit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Repositories.Edit(context.Background(), "%", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetVulnerabilityAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/vulnerability-alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeRequiredVulnerabilityAlertsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	vulnerabilityAlertsEnabled, _, err := client.Repositories.GetVulnerabilityAlerts(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetVulnerabilityAlerts returned error: %v", err)
	}

	if want := true; vulnerabilityAlertsEnabled != want {
		t.Errorf("Repositories.GetVulnerabilityAlerts returned %+v, want %+v", vulnerabilityAlertsEnabled, want)
	}
}

func TestRepositoriesService_EnableVulnerabilityAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/vulnerability-alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeRequiredVulnerabilityAlertsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	if _, err := client.Repositories.EnableVulnerabilityAlerts(context.Background(), "o", "r"); err != nil {
		t.Errorf("Repositories.EnableVulnerabilityAlerts returned error: %v", err)
	}
}

func TestRepositoriesService_DisableVulnerabilityAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/vulnerability-alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeRequiredVulnerabilityAlertsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	if _, err := client.Repositories.DisableVulnerabilityAlerts(context.Background(), "o", "r"); err != nil {
		t.Errorf("Repositories.DisableVulnerabilityAlerts returned error: %v", err)
	}
}

func TestRepositoriesService_EnableAutomatedSecurityFixes(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeRequiredAutomatedSecurityFixesPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	if _, err := client.Repositories.EnableAutomatedSecurityFixes(context.Background(), "o", "r"); err != nil {
		t.Errorf("Repositories.EnableAutomatedSecurityFixes returned error: %v", err)
	}
}

func TestRepositoriesService_DisableAutomatedSecurityFixes(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeRequiredAutomatedSecurityFixesPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	if _, err := client.Repositories.DisableAutomatedSecurityFixes(context.Background(), "o", "r"); err != nil {
		t.Errorf("Repositories.DisableAutomatedSecurityFixes returned error: %v", err)
	}
}

func TestRepositoriesService_ListContributors(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contributors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"anon": "true",
			"page": "2",
		})
		fmt.Fprint(w, `[{"contributions":42}]`)
	})

	opts := &ListContributorsOptions{Anon: "true", ListOptions: ListOptions{Page: 2}}
	contributors, _, err := client.Repositories.ListContributors(context.Background(), "o", "r", opts)
	if err != nil {
		t.Errorf("Repositories.ListContributors returned error: %v", err)
	}

	want := []*Contributor{{Contributions: Int(42)}}
	if !reflect.DeepEqual(contributors, want) {
		t.Errorf("Repositories.ListContributors returned %+v, want %+v", contributors, want)
	}
}

func TestRepositoriesService_ListLanguages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/languages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"go":1}`)
	})

	languages, _, err := client.Repositories.ListLanguages(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Repositories.ListLanguages returned error: %v", err)
	}

	want := map[string]int{"go": 1}
	if !reflect.DeepEqual(languages, want) {
		t.Errorf("Repositories.ListLanguages returned %+v, want %+v", languages, want)
	}
}

func TestRepositoriesService_ListTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	teams, _, err := client.Repositories.ListTeams(context.Background(), "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Repositories.ListTeams returned %+v, want %+v", teams, want)
	}
}

func TestRepositoriesService_ListTags(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name":"n", "commit" : {"sha" : "s", "url" : "u"}, "zipball_url": "z", "tarball_url": "t"}]`)
	})

	opt := &ListOptions{Page: 2}
	tags, _, err := client.Repositories.ListTags(context.Background(), "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListTags returned error: %v", err)
	}

	want := []*RepositoryTag{
		{
			Name: String("n"),
			Commit: &Commit{
				SHA: String("s"),
				URL: String("u"),
			},
			ZipballURL: String("z"),
			TarballURL: String("t"),
		},
	}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Repositories.ListTags returned %+v, want %+v", tags, want)
	}
}

func TestRepositoriesService_ListBranches(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name":"master", "commit" : {"sha" : "a57781", "url" : "https://api.github.com/repos/o/r/commits/a57781"}}]`)
	})

	opt := &BranchListOptions{
		Protected:   nil,
		ListOptions: ListOptions{Page: 2},
	}
	branches, _, err := client.Repositories.ListBranches(context.Background(), "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListBranches returned error: %v", err)
	}

	want := []*Branch{{Name: String("master"), Commit: &RepositoryCommit{SHA: String("a57781"), URL: String("https://api.github.com/repos/o/r/commits/a57781")}}}
	if !reflect.DeepEqual(branches, want) {
		t.Errorf("Repositories.ListBranches returned %+v, want %+v", branches, want)
	}
}

func TestRepositoriesService_GetBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprint(w, `{"name":"n", "commit":{"sha":"s","commit":{"message":"m"}}, "protected":true}`)
	})

	branch, _, err := client.Repositories.GetBranch(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetBranch returned error: %v", err)
	}

	want := &Branch{
		Name: String("n"),
		Commit: &RepositoryCommit{
			SHA: String("s"),
			Commit: &Commit{
				Message: String("m"),
			},
		},
		Protected: Bool(true),
	}

	if !reflect.DeepEqual(branch, want) {
		t.Errorf("Repositories.GetBranch returned %+v, want %+v", branch, want)
	}
}

func TestRepositoriesService_GetBranchProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{
				"required_status_checks":{
					"strict":true,
					"contexts":["continuous-integration"]
				},
				"required_pull_request_reviews":{
					"dismissal_restrictions":{
						"users":[{
							"id":3,
							"login":"u"
						}],
						"teams":[{
							"id":4,
							"slug":"t"
						}]
					},
					"dismiss_stale_reviews":true,
					"require_code_owner_reviews":true,
					"required_approving_review_count":1
					},
					"enforce_admins":{
						"url":"/repos/o/r/branches/b/protection/enforce_admins",
						"enabled":true
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}]
					}
				}`)
	})

	protection, _, err := client.Repositories.GetBranchProtection(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetBranchProtection returned error: %v", err)
	}

	want := &Protection{
		RequiredStatusChecks: &RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
			DismissStaleReviews: true,
			DismissalRestrictions: &DismissalRestrictions{
				Users: []*User{
					{Login: String("u"), ID: Int64(3)},
				},
				Teams: []*Team{
					{Slug: String("t"), ID: Int64(4)},
				},
			},
			RequireCodeOwnerReviews:      true,
			RequiredApprovingReviewCount: 1,
		},
		EnforceAdmins: &AdminEnforcement{
			URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
			Enabled: true,
		},
		Restrictions: &BranchRestrictions{
			Users: []*User{
				{Login: String("u"), ID: Int64(1)},
			},
			Teams: []*Team{
				{Slug: String("t"), ID: Int64(2)},
			},
		},
	}
	if !reflect.DeepEqual(protection, want) {
		t.Errorf("Repositories.GetBranchProtection returned %+v, want %+v", protection, want)
	}
}

func TestRepositoriesService_GetBranchProtection_noDismissalRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{
				"required_status_checks":{
					"strict":true,
					"contexts":["continuous-integration"]
				},
				"required_pull_request_reviews":{
					"dismiss_stale_reviews":true,
					"require_code_owner_reviews":true,
					"required_approving_review_count":1
					},
					"enforce_admins":{
						"url":"/repos/o/r/branches/b/protection/enforce_admins",
						"enabled":true
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}]
					}
				}`)
	})

	protection, _, err := client.Repositories.GetBranchProtection(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetBranchProtection returned error: %v", err)
	}

	want := &Protection{
		RequiredStatusChecks: &RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
			DismissStaleReviews:          true,
			DismissalRestrictions:        nil,
			RequireCodeOwnerReviews:      true,
			RequiredApprovingReviewCount: 1,
		},
		EnforceAdmins: &AdminEnforcement{
			URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
			Enabled: true,
		},
		Restrictions: &BranchRestrictions{
			Users: []*User{
				{Login: String("u"), ID: Int64(1)},
			},
			Teams: []*Team{
				{Slug: String("t"), ID: Int64(2)},
			},
		},
	}
	if !reflect.DeepEqual(protection, want) {
		t.Errorf("Repositories.GetBranchProtection returned %+v, want %+v", protection, want)
	}
}

func TestRepositoriesService_UpdateBranchProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ProtectionRequest{
		RequiredStatusChecks: &RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
			DismissStaleReviews: true,
			DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
				Users: &[]string{"uu"},
				Teams: &[]string{"tt"},
			},
		},
		Restrictions: &BranchRestrictionsRequest{
			Users: []string{"u"},
			Teams: []string{"t"},
			Apps:  []string{"a"},
		},
	}

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{
			"required_status_checks":{
				"strict":true,
				"contexts":["continuous-integration"]
			},
			"required_pull_request_reviews":{
				"dismissal_restrictions":{
					"users":[{
						"id":3,
						"login":"uu"
					}],
					"teams":[{
						"id":4,
						"slug":"tt"
					}]
				},
				"dismiss_stale_reviews":true,
				"require_code_owner_reviews":true
			},
			"restrictions":{
				"users":[{"id":1,"login":"u"}],
				"teams":[{"id":2,"slug":"t"}],
				"apps":[{"id":3,"slug":"a"}]
			}
		}`)
	})

	protection, _, err := client.Repositories.UpdateBranchProtection(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
	}

	want := &Protection{
		RequiredStatusChecks: &RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
			DismissStaleReviews: true,
			DismissalRestrictions: &DismissalRestrictions{
				Users: []*User{
					{Login: String("uu"), ID: Int64(3)},
				},
				Teams: []*Team{
					{Slug: String("tt"), ID: Int64(4)},
				},
			},
			RequireCodeOwnerReviews: true,
		},
		Restrictions: &BranchRestrictions{
			Users: []*User{
				{Login: String("u"), ID: Int64(1)},
			},
			Teams: []*Team{
				{Slug: String("t"), ID: Int64(2)},
			},
			Apps: []*App{
				{Slug: String("a"), ID: Int64(3)},
			},
		},
	}
	if !reflect.DeepEqual(protection, want) {
		t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
	}
}

func TestRepositoriesService_RemoveBranchProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Repositories.RemoveBranchProtection(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemoveBranchProtection returned error: %v", err)
	}
}

func TestRepositoriesService_ListLanguages_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Repositories.ListLanguages(context.Background(), "%", "%")
	testURLParseError(t, err)
}

func TestRepositoriesService_License(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/license", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name": "LICENSE", "path": "LICENSE", "license":{"key":"mit","name":"MIT License","spdx_id":"MIT","url":"https://api.github.com/licenses/mit","featured":true}}`)
	})

	got, _, err := client.Repositories.License(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Repositories.License returned error: %v", err)
	}

	want := &RepositoryLicense{
		Name: String("LICENSE"),
		Path: String("LICENSE"),
		License: &License{
			Name:     String("MIT License"),
			Key:      String("mit"),
			SPDXID:   String("MIT"),
			URL:      String("https://api.github.com/licenses/mit"),
			Featured: Bool(true),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.License returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_GetRequiredStatusChecks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprint(w, `{"strict": true,"contexts": ["x","y","z"]}`)
	})

	checks, _, err := client.Repositories.GetRequiredStatusChecks(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetRequiredStatusChecks returned error: %v", err)
	}

	want := &RequiredStatusChecks{
		Strict:   true,
		Contexts: []string{"x", "y", "z"},
	}
	if !reflect.DeepEqual(checks, want) {
		t.Errorf("Repositories.GetRequiredStatusChecks returned %+v, want %+v", checks, want)
	}
}

func TestRepositoriesService_UpdateRequiredStatusChecks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RequiredStatusChecksRequest{
		Strict:   Bool(true),
		Contexts: []string{"continuous-integration"},
	}

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequiredStatusChecksRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		testHeader(t, r, "Accept", mediaTypeV3)
		fmt.Fprintf(w, `{"strict":true,"contexts":["continuous-integration"]}`)
	})

	statusChecks, _, err := client.Repositories.UpdateRequiredStatusChecks(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.UpdateRequiredStatusChecks returned error: %v", err)
	}

	want := &RequiredStatusChecks{
		Strict:   true,
		Contexts: []string{"continuous-integration"},
	}
	if !reflect.DeepEqual(statusChecks, want) {
		t.Errorf("Repositories.UpdateRequiredStatusChecks returned %+v, want %+v", statusChecks, want)
	}
}

func TestRepositoriesService_ListRequiredStatusChecksContexts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks/contexts", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprint(w, `["x", "y", "z"]`)
	})

	contexts, _, err := client.Repositories.ListRequiredStatusChecksContexts(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned error: %v", err)
	}

	want := []string{"x", "y", "z"}
	if !reflect.DeepEqual(contexts, want) {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned %+v, want %+v", contexts, want)
	}
}

func TestRepositoriesService_GetPullRequestReviewEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_pull_request_reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{
			"dismissal_restrictions":{
				"users":[{"id":1,"login":"u"}],
				"teams":[{"id":2,"slug":"t"}]
			},
			"dismiss_stale_reviews":true,
			"require_code_owner_reviews":true,
			"required_approving_review_count":1
		}`)
	})

	enforcement, _, err := client.Repositories.GetPullRequestReviewEnforcement(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetPullRequestReviewEnforcement returned error: %v", err)
	}

	want := &PullRequestReviewsEnforcement{
		DismissStaleReviews: true,
		DismissalRestrictions: &DismissalRestrictions{
			Users: []*User{
				{Login: String("u"), ID: Int64(1)},
			},
			Teams: []*Team{
				{Slug: String("t"), ID: Int64(2)},
			},
		},
		RequireCodeOwnerReviews:      true,
		RequiredApprovingReviewCount: 1,
	}

	if !reflect.DeepEqual(enforcement, want) {
		t.Errorf("Repositories.GetPullRequestReviewEnforcement returned %+v, want %+v", enforcement, want)
	}
}

func TestRepositoriesService_UpdatePullRequestReviewEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewsEnforcementUpdate{
		DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
			Users: &[]string{"u"},
			Teams: &[]string{"t"},
		},
	}

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_pull_request_reviews", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewsEnforcementUpdate)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{
			"dismissal_restrictions":{
				"users":[{"id":1,"login":"u"}],
				"teams":[{"id":2,"slug":"t"}]
			},
			"dismiss_stale_reviews":true,
			"require_code_owner_reviews":true,
			"required_approving_review_count":3
		}`)
	})

	enforcement, _, err := client.Repositories.UpdatePullRequestReviewEnforcement(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.UpdatePullRequestReviewEnforcement returned error: %v", err)
	}

	want := &PullRequestReviewsEnforcement{
		DismissStaleReviews: true,
		DismissalRestrictions: &DismissalRestrictions{
			Users: []*User{
				{Login: String("u"), ID: Int64(1)},
			},
			Teams: []*Team{
				{Slug: String("t"), ID: Int64(2)},
			},
		},
		RequireCodeOwnerReviews:      true,
		RequiredApprovingReviewCount: 3,
	}
	if !reflect.DeepEqual(enforcement, want) {
		t.Errorf("Repositories.UpdatePullRequestReviewEnforcement returned %+v, want %+v", enforcement, want)
	}
}

func TestRepositoriesService_DisableDismissalRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_pull_request_reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		// TODO: remove custom Accept header when this API fully launches
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		testBody(t, r, `{"dismissal_restrictions":{}}`+"\n")
		fmt.Fprintf(w, `{"dismiss_stale_reviews":true,"require_code_owner_reviews":true,"required_approving_review_count":1}`)
	})

	enforcement, _, err := client.Repositories.DisableDismissalRestrictions(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.DisableDismissalRestrictions returned error: %v", err)
	}

	want := &PullRequestReviewsEnforcement{
		DismissStaleReviews:          true,
		DismissalRestrictions:        nil,
		RequireCodeOwnerReviews:      true,
		RequiredApprovingReviewCount: 1,
	}
	if !reflect.DeepEqual(enforcement, want) {
		t.Errorf("Repositories.DisableDismissalRestrictions returned %+v, want %+v", enforcement, want)
	}
}

func TestRepositoriesService_RemovePullRequestReviewEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_pull_request_reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Repositories.RemovePullRequestReviewEnforcement(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemovePullRequestReviewEnforcement returned error: %v", err)
	}
}

func TestRepositoriesService_GetAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
	})

	enforcement, _, err := client.Repositories.GetAdminEnforcement(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetAdminEnforcement returned error: %v", err)
	}

	want := &AdminEnforcement{
		URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
		Enabled: true,
	}

	if !reflect.DeepEqual(enforcement, want) {
		t.Errorf("Repositories.GetAdminEnforcement returned %+v, want %+v", enforcement, want)
	}
}

func TestRepositoriesService_AddAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
	})

	enforcement, _, err := client.Repositories.AddAdminEnforcement(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.AddAdminEnforcement returned error: %v", err)
	}

	want := &AdminEnforcement{
		URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
		Enabled: true,
	}
	if !reflect.DeepEqual(enforcement, want) {
		t.Errorf("Repositories.AddAdminEnforcement returned %+v, want %+v", enforcement, want)
	}
}

func TestRepositoriesService_RemoveAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Repositories.RemoveAdminEnforcement(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemoveAdminEnforcement returned error: %v", err)
	}
}

func TestRepositoriesService_GetSignaturesProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":false}`)
	})

	signature, _, err := client.Repositories.GetSignaturesProtectedBranch(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetSignaturesProtectedBranch returned error: %v", err)
	}

	want := &SignaturesProtectedBranch{
		URL:     String("/repos/o/r/branches/b/protection/required_signatures"),
		Enabled: Bool(false),
	}

	if !reflect.DeepEqual(signature, want) {
		t.Errorf("Repositories.GetSignaturesProtectedBranch returned %+v, want %+v", signature, want)
	}
}

func TestRepositoriesService_RequireSignaturesOnProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":true}`)
	})

	signature, _, err := client.Repositories.RequireSignaturesOnProtectedBranch(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RequireSignaturesOnProtectedBranch returned error: %v", err)
	}

	want := &SignaturesProtectedBranch{
		URL:     String("/repos/o/r/branches/b/protection/required_signatures"),
		Enabled: Bool(true),
	}

	if !reflect.DeepEqual(signature, want) {
		t.Errorf("Repositories.RequireSignaturesOnProtectedBranch returned %+v, want %+v", signature, want)
	}
}

func TestRepositoriesService_OptionalSignaturesOnProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Repositories.OptionalSignaturesOnProtectedBranch(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.OptionalSignaturesOnProtectedBranch returned error: %v", err)
	}
}

func TestPullRequestReviewsEnforcementRequest_MarshalJSON_nilDismissalRestirctions(t *testing.T) {
	req := PullRequestReviewsEnforcementRequest{}

	got, err := json.Marshal(req)
	if err != nil {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned error: %v", err)
	}

	want := `{"dismiss_stale_reviews":false,"require_code_owner_reviews":false,"required_approving_review_count":0}`
	if want != string(got) {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned %+v, want %+v", string(got), want)
	}

	req = PullRequestReviewsEnforcementRequest{
		DismissalRestrictionsRequest: &DismissalRestrictionsRequest{},
	}

	got, err = json.Marshal(req)
	if err != nil {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned error: %v", err)
	}

	want = `{"dismissal_restrictions":{},"dismiss_stale_reviews":false,"require_code_owner_reviews":false,"required_approving_review_count":0}`
	if want != string(got) {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned %+v, want %+v", string(got), want)
	}

	req = PullRequestReviewsEnforcementRequest{
		DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
			Users: &[]string{},
			Teams: &[]string{},
		},
	}

	got, err = json.Marshal(req)
	if err != nil {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned error: %v", err)
	}

	want = `{"dismissal_restrictions":{"users":[],"teams":[]},"dismiss_stale_reviews":false,"require_code_owner_reviews":false,"required_approving_review_count":0}`
	if want != string(got) {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned %+v, want %+v", string(got), want)
	}
}

func TestRepositoriesService_ListAllTopics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		fmt.Fprint(w, `{"names":["go", "go-github", "github"]}`)
	})

	got, _, err := client.Repositories.ListAllTopics(context.Background(), "o", "r")
	if err != nil {
		t.Fatalf("Repositories.ListAllTopics returned error: %v", err)
	}

	want := []string{"go", "go-github", "github"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ListAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ListAllTopics_emptyTopics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		fmt.Fprint(w, `{"names":[]}`)
	})

	got, _, err := client.Repositories.ListAllTopics(context.Background(), "o", "r")
	if err != nil {
		t.Fatalf("Repositories.ListAllTopics returned error: %v", err)
	}

	want := []string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ListAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ReplaceAllTopics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		fmt.Fprint(w, `{"names":["go", "go-github", "github"]}`)
	})

	got, _, err := client.Repositories.ReplaceAllTopics(context.Background(), "o", "r", []string{"go", "go-github", "github"})
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{"go", "go-github", "github"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ReplaceAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ReplaceAllTopics_nilSlice(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		testBody(t, r, `{"names":[]}`+"\n")
		fmt.Fprint(w, `{"names":[]}`)
	})

	got, _, err := client.Repositories.ReplaceAllTopics(context.Background(), "o", "r", nil)
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ReplaceAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ReplaceAllTopics_emptySlice(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		testBody(t, r, `{"names":[]}`+"\n")
		fmt.Fprint(w, `{"names":[]}`)
	})

	got, _, err := client.Repositories.ReplaceAllTopics(context.Background(), "o", "r", []string{})
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ReplaceAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ListApps(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	_, _, err := client.Repositories.ListApps(context.Background(), "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.ListApps returned error: %v", err)
	}
}

func TestRepositoriesService_ReplaceAppRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
	})
	input := []string{"octocat"}
	got, _, err := client.Repositories.ReplaceAppRestrictions(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.ReplaceAppRestrictions returned error: %v", err)
	}
	want := []*App{
		{Name: String("octocat")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ReplaceAppRestrictions returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_AddAppRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
	})
	input := []string{"octocat"}
	got, _, err := client.Repositories.AddAppRestrictions(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.AddAppRestrictions returned error: %v", err)
	}
	want := []*App{
		{Name: String("octocat")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.AddAppRestrictions returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_RemoveAppRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `[]`)
	})
	input := []string{"octocat"}
	got, _, err := client.Repositories.RemoveAppRestrictions(context.Background(), "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.RemoveAppRestrictions returned error: %v", err)
	}
	want := []*App{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.RemoveAppRestrictions returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_Transfer(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := TransferRequest{NewOwner: "a", TeamID: []int64{123}}

	mux.HandleFunc("/repos/o/r/transfer", func(w http.ResponseWriter, r *http.Request) {
		var v TransferRequest
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"owner":{"login":"a"}}`)
	})

	got, _, err := client.Repositories.Transfer(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.Transfer returned error: %v", err)
	}

	want := &Repository{Owner: &User{Login: String("a")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.Transfer returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_Dispatch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var input DispatchRequestOptions

	mux.HandleFunc("/repos/o/r/dispatches", func(w http.ResponseWriter, r *http.Request) {
		var v DispatchRequestOptions
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"owner":{"login":"a"}}`)
	})

	ctx := context.Background()

	testCases := []interface{}{
		nil,
		struct {
			Foo string
		}{
			Foo: "test",
		},
		struct {
			Bar int
		}{
			Bar: 42,
		},
		struct {
			Foo string
			Bar int
			Baz bool
		}{
			Foo: "test",
			Bar: 42,
			Baz: false,
		},
	}
	for _, tc := range testCases {

		if tc == nil {
			input = DispatchRequestOptions{EventType: "go"}
		} else {
			bytes, _ := json.Marshal(tc)
			payload := json.RawMessage(bytes)
			input = DispatchRequestOptions{EventType: "go", ClientPayload: &payload}
		}

		got, _, err := client.Repositories.Dispatch(ctx, "o", "r", input)
		if err != nil {
			t.Errorf("Repositories.Dispatch returned error: %v", err)
		}

		want := &Repository{Owner: &User{Login: String("a")}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Repositories.Dispatch returned %+v, want %+v", got, want)
		}
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Repositories.Dispatch(ctx, "o", "r", input)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' Dispatch = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' Dispatch resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' Dispatch err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Repositories.Dispatch(ctx, "o", "r", input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now Dispatch = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now Dispatch resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now Dispatch err = nil, want error")
	}
}
