// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_List_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview, mediaTypeRepositoryVisibilityPreview}
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", got, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.List(ctx, "\n", &RepositoryListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.List(ctx, "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_List_specifiedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview, mediaTypeRepositoryVisibilityPreview}
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
	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, "u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_specifiedUser_type(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview, mediaTypeRepositoryVisibilityPreview}
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
	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, "u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.List(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListByOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeTopicsPreview, mediaTypeRepositoryVisibilityPreview}
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ListByOrg returned %+v, want %+v", got, want)
	}

	const methodName = "ListByOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListByOrg(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListByOrg(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListByOrg_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListByOrg(ctx, "%", nil)
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ListAll returned %+v, want %+v", got, want)
	}

	const methodName = "ListAll"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListAll(ctx, &RepositoryListAllOptions{1})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, want) {
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", got, want)
	}

	const methodName = "Create"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.Create(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.Create(ctx, "", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	repo, _, err := client.Repositories.Create(ctx, "o", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !cmp.Equal(repo, want) {
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
		if !cmp.Equal(v, want) {
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.CreateFromTemplate returned %+v, want %+v", got, want)
	}

	const methodName = "CreateFromTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateFromTemplate(ctx, "\n", "\n", templateRepoReq)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateFromTemplate(ctx, "to", "tr", templateRepoReq)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview, mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"},"security_and_analysis":{"advanced_security":{"status":"enabled"},"secret_scanning":{"status":"enabled"}}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Get(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.Get returned error: %v", err)
	}

	want := &Repository{ID: Int64(1), Name: String("n"), Description: String("d"), Owner: &User{Login: String("l")}, License: &License{Key: String("mit")}, SecurityAndAnalysis: &SecurityAndAnalysis{AdvancedSecurity: &AdvancedSecurity{Status: String("enabled")}, SecretScanning: &SecretScanning{String("enabled")}}}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.Get returned %+v, want %+v", got, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.Get(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.Get(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.GetCodeOfConduct returned %+v, want %+v", got, want)
	}

	const methodName = "GetCodeOfConduct"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCodeOfConduct(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCodeOfConduct(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.GetByID returned %+v, want %+v", got, want)
	}

	const methodName = "GetByID"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetByID(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
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
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.Edit returned %+v, want %+v", got, want)
	}

	const methodName = "Edit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.Edit(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.Edit(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	const methodName = "Delete"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.Delete(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.Delete(ctx, "o", "r")
	})
}

func TestRepositoriesService_Get_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.Get(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_Edit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.Edit(ctx, "%", "r", nil)
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

	ctx := context.Background()
	vulnerabilityAlertsEnabled, _, err := client.Repositories.GetVulnerabilityAlerts(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetVulnerabilityAlerts returned error: %v", err)
	}

	if want := true; vulnerabilityAlertsEnabled != want {
		t.Errorf("Repositories.GetVulnerabilityAlerts returned %+v, want %+v", vulnerabilityAlertsEnabled, want)
	}

	const methodName = "GetVulnerabilityAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetVulnerabilityAlerts(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetVulnerabilityAlerts(ctx, "o", "r")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EnableVulnerabilityAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/vulnerability-alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeRequiredVulnerabilityAlertsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.EnableVulnerabilityAlerts(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.EnableVulnerabilityAlerts returned error: %v", err)
	}

	const methodName = "EnableVulnerabilityAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.EnableVulnerabilityAlerts(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.EnableVulnerabilityAlerts(ctx, "o", "r")
	})
}

func TestRepositoriesService_DisableVulnerabilityAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/vulnerability-alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeRequiredVulnerabilityAlertsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.DisableVulnerabilityAlerts(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.DisableVulnerabilityAlerts returned error: %v", err)
	}

	const methodName = "DisableVulnerabilityAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DisableVulnerabilityAlerts(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DisableVulnerabilityAlerts(ctx, "o", "r")
	})
}

func TestRepositoriesService_EnableAutomatedSecurityFixes(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeRequiredAutomatedSecurityFixesPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.EnableAutomatedSecurityFixes(ctx, "o", "r"); err != nil {
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

	ctx := context.Background()
	if _, err := client.Repositories.DisableAutomatedSecurityFixes(ctx, "o", "r"); err != nil {
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
	ctx := context.Background()
	contributors, _, err := client.Repositories.ListContributors(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Repositories.ListContributors returned error: %v", err)
	}

	want := []*Contributor{{Contributions: Int(42)}}
	if !cmp.Equal(contributors, want) {
		t.Errorf("Repositories.ListContributors returned %+v, want %+v", contributors, want)
	}

	const methodName = "ListContributors"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListContributors(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListContributors(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListLanguages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/languages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"go":1}`)
	})

	ctx := context.Background()
	languages, _, err := client.Repositories.ListLanguages(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.ListLanguages returned error: %v", err)
	}

	want := map[string]int{"go": 1}
	if !cmp.Equal(languages, want) {
		t.Errorf("Repositories.ListLanguages returned %+v, want %+v", languages, want)
	}

	const methodName = "ListLanguages"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListLanguages(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListLanguages(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	teams, _, err := client.Repositories.ListTeams(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !cmp.Equal(teams, want) {
		t.Errorf("Repositories.ListTeams returned %+v, want %+v", teams, want)
	}

	const methodName = "ListTeams"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListTeams(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListTeams(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	tags, _, err := client.Repositories.ListTags(ctx, "o", "r", opt)
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
	if !cmp.Equal(tags, want) {
		t.Errorf("Repositories.ListTags returned %+v, want %+v", tags, want)
	}

	const methodName = "ListTags"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListTags(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListTags(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListBranches(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name":"master", "commit" : {"sha" : "a57781", "url" : "https://api.github.com/repos/o/r/commits/a57781"}}]`)
	})

	opt := &BranchListOptions{
		Protected:   nil,
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	branches, _, err := client.Repositories.ListBranches(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListBranches returned error: %v", err)
	}

	want := []*Branch{{Name: String("master"), Commit: &RepositoryCommit{SHA: String("a57781"), URL: String("https://api.github.com/repos/o/r/commits/a57781")}}}
	if !cmp.Equal(branches, want) {
		t.Errorf("Repositories.ListBranches returned %+v, want %+v", branches, want)
	}

	const methodName = "ListBranches"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListBranches(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListBranches(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"n", "commit":{"sha":"s","commit":{"message":"m"}}, "protected":true}`)
	})

	ctx := context.Background()
	branch, _, err := client.Repositories.GetBranch(ctx, "o", "r", "b", false)
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

	if !cmp.Equal(branch, want) {
		t.Errorf("Repositories.GetBranch returned %+v, want %+v", branch, want)
	}

	const methodName = "GetBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetBranch(ctx, "\n", "\n", "\n", false)
		return err
	})
}

func TestRepositoriesService_GetBranch_StatusMovedPermanently_followRedirects(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/repos/o/r/branches/br")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})
	mux.HandleFunc("/repos/o/r/branches/br", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"n", "commit":{"sha":"s","commit":{"message":"m"}}, "protected":true}`)
	})
	ctx := context.Background()
	branch, resp, err := client.Repositories.GetBranch(ctx, "o", "r", "b", true)
	if err != nil {
		t.Errorf("Repositories.GetBranch returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.GetBranch returned status: %d, want %d", resp.StatusCode, http.StatusOK)
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
	if !cmp.Equal(branch, want) {
		t.Errorf("Repositories.GetBranch returned %+v, want %+v", branch, want)
	}
}

func TestRepositoriesService_GetBranch_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "branch not found", http.StatusNotFound)
	})
	ctx := context.Background()
	_, resp, err := client.Repositories.GetBranch(ctx, "o", "r", "b", true)
	if err == nil {
		t.Error("Repositories.GetBranch returned error: nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Repositories.GetBranch returned status: %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	// Add custom round tripper
	client.client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("failed to get branch")
	})

	const methodName = "GetBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetBranch(ctx, "o", "r", "b", true)
		return err
	})
}

func TestRepositoriesService_RenameBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	renameBranchReq := "nn"

	mux.HandleFunc("/repos/o/r/branches/b/rename", func(w http.ResponseWriter, r *http.Request) {
		v := new(renameBranchRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &renameBranchRequest{NewName: "nn"}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"protected":true,"name":"nn"}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.RenameBranch(ctx, "o", "r", "b", renameBranchReq)
	if err != nil {
		t.Errorf("Repositories.RenameBranch returned error: %v", err)
	}

	want := &Branch{Name: String("nn"), Protected: Bool(true)}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.RenameBranch returned %+v, want %+v", got, want)
	}

	const methodName = "RenameBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.RenameBranch(ctx, "\n", "\n", "\n", renameBranchReq)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.RenameBranch(ctx, "o", "r", "b", renameBranchReq)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
					},
					"required_conversation_resolution": {
						"enabled": true
					}
				}`)
	})

	ctx := context.Background()
	protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", "b")
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
		RequiredConversationResolution: &RequiredConversationResolution{
			Enabled: true,
		},
	}
	if !cmp.Equal(protection, want) {
		t.Errorf("Repositories.GetBranchProtection returned %+v, want %+v", protection, want)
	}

	const methodName = "GetBranchProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetBranchProtection(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetBranchProtection(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", "b")
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
	if !cmp.Equal(protection, want) {
		t.Errorf("Repositories.GetBranchProtection returned %+v, want %+v", protection, want)
	}
}

func TestRepositoriesService_GetBranchProtection_branchNotProtected(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{
			"message": %q,
			"documentation_url": "https://docs.github.com/rest/reference/repos#get-branch-protection"
			}`, githubBranchNotProtected)
	})

	ctx := context.Background()
	protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", "b")

	if protection != nil {
		t.Errorf("Repositories.GetBranchProtection returned non-nil protection data")
	}

	if err != ErrBranchNotProtected {
		t.Errorf("Repositories.GetBranchProtection returned an invalid error: %v", err)
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
		if !cmp.Equal(v, input) {
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

	ctx := context.Background()
	protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", "b", input)
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
	if !cmp.Equal(protection, want) {
		t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
	}

	const methodName = "UpdateBranchProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdateBranchProtection(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RemoveBranchProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.RemoveBranchProtection(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemoveBranchProtection returned error: %v", err)
	}

	const methodName = "RemoveBranchProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.RemoveBranchProtection(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.RemoveBranchProtection(ctx, "o", "r", "b")
	})
}

func TestRepositoriesService_ListLanguages_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListLanguages(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestRepositoriesService_License(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/license", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name": "LICENSE", "path": "LICENSE", "license":{"key":"mit","name":"MIT License","spdx_id":"MIT","url":"https://api.github.com/licenses/mit","featured":true}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.License(ctx, "o", "r")
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

	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.License returned %+v, want %+v", got, want)
	}

	const methodName = "License"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.License(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.License(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRequiredStatusChecks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"strict": true,"contexts": ["x","y","z"]}`)
	})

	ctx := context.Background()
	checks, _, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetRequiredStatusChecks returned error: %v", err)
	}

	want := &RequiredStatusChecks{
		Strict:   true,
		Contexts: []string{"x", "y", "z"},
	}
	if !cmp.Equal(checks, want) {
		t.Errorf("Repositories.GetRequiredStatusChecks returned %+v, want %+v", checks, want)
	}

	const methodName = "GetRequiredStatusChecks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetRequiredStatusChecks(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRequiredStatusChecks_branchNotProtected(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{
			"message": %q,
			"documentation_url": "https://docs.github.com/rest/reference/repos#get-branch-protection"
			}`, githubBranchNotProtected)
	})

	ctx := context.Background()
	checks, _, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", "b")

	if checks != nil {
		t.Errorf("Repositories.GetRequiredStatusChecks returned non-nil status-checks data")
	}

	if err != ErrBranchNotProtected {
		t.Errorf("Repositories.GetRequiredStatusChecks returned an invalid error: %v", err)
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		testHeader(t, r, "Accept", mediaTypeV3)
		fmt.Fprintf(w, `{"strict":true,"contexts":["continuous-integration"]}`)
	})

	ctx := context.Background()
	statusChecks, _, err := client.Repositories.UpdateRequiredStatusChecks(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.UpdateRequiredStatusChecks returned error: %v", err)
	}

	want := &RequiredStatusChecks{
		Strict:   true,
		Contexts: []string{"continuous-integration"},
	}
	if !cmp.Equal(statusChecks, want) {
		t.Errorf("Repositories.UpdateRequiredStatusChecks returned %+v, want %+v", statusChecks, want)
	}

	const methodName = "UpdateRequiredStatusChecks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdateRequiredStatusChecks(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRequiredStatusChecks(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RemoveRequiredStatusChecks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeV3)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.RemoveRequiredStatusChecks(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemoveRequiredStatusChecks returned error: %v", err)
	}

	const methodName = "RemoveRequiredStatusChecks"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.RemoveRequiredStatusChecks(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.RemoveRequiredStatusChecks(ctx, "o", "r", "b")
	})
}

func TestRepositoriesService_ListRequiredStatusChecksContexts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks/contexts", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "GET")
		fmt.Fprint(w, `["x", "y", "z"]`)
	})

	ctx := context.Background()
	contexts, _, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned error: %v", err)
	}

	want := []string{"x", "y", "z"}
	if !cmp.Equal(contexts, want) {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned %+v, want %+v", contexts, want)
	}

	const methodName = "ListRequiredStatusChecksContexts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListRequiredStatusChecksContexts(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListRequiredStatusChecksContexts_branchNotProtected(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_status_checks/contexts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{
			"message": %q,
			"documentation_url": "https://docs.github.com/rest/reference/repos#get-branch-protection"
			}`, githubBranchNotProtected)
	})

	ctx := context.Background()
	contexts, _, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", "b")

	if contexts != nil {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned non-nil contexts data")
	}

	if err != ErrBranchNotProtected {
		t.Errorf("Repositories.ListRequiredStatusChecksContexts returned an invalid error: %v", err)
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

	ctx := context.Background()
	enforcement, _, err := client.Repositories.GetPullRequestReviewEnforcement(ctx, "o", "r", "b")
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

	if !cmp.Equal(enforcement, want) {
		t.Errorf("Repositories.GetPullRequestReviewEnforcement returned %+v, want %+v", enforcement, want)
	}

	const methodName = "GetPullRequestReviewEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetPullRequestReviewEnforcement(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetPullRequestReviewEnforcement(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
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

	ctx := context.Background()
	enforcement, _, err := client.Repositories.UpdatePullRequestReviewEnforcement(ctx, "o", "r", "b", input)
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
	if !cmp.Equal(enforcement, want) {
		t.Errorf("Repositories.UpdatePullRequestReviewEnforcement returned %+v, want %+v", enforcement, want)
	}

	const methodName = "UpdatePullRequestReviewEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdatePullRequestReviewEnforcement(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdatePullRequestReviewEnforcement(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	enforcement, _, err := client.Repositories.DisableDismissalRestrictions(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.DisableDismissalRestrictions returned error: %v", err)
	}

	want := &PullRequestReviewsEnforcement{
		DismissStaleReviews:          true,
		DismissalRestrictions:        nil,
		RequireCodeOwnerReviews:      true,
		RequiredApprovingReviewCount: 1,
	}
	if !cmp.Equal(enforcement, want) {
		t.Errorf("Repositories.DisableDismissalRestrictions returned %+v, want %+v", enforcement, want)
	}

	const methodName = "DisableDismissalRestrictions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.DisableDismissalRestrictions(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.DisableDismissalRestrictions(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RemovePullRequestReviewEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_pull_request_reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.RemovePullRequestReviewEnforcement(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemovePullRequestReviewEnforcement returned error: %v", err)
	}

	const methodName = "RemovePullRequestReviewEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.RemovePullRequestReviewEnforcement(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.RemovePullRequestReviewEnforcement(ctx, "o", "r", "b")
	})
}

func TestRepositoriesService_GetAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
	})

	ctx := context.Background()
	enforcement, _, err := client.Repositories.GetAdminEnforcement(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetAdminEnforcement returned error: %v", err)
	}

	want := &AdminEnforcement{
		URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
		Enabled: true,
	}

	if !cmp.Equal(enforcement, want) {
		t.Errorf("Repositories.GetAdminEnforcement returned %+v, want %+v", enforcement, want)
	}

	const methodName = "GetAdminEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetAdminEnforcement(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAdminEnforcement(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_AddAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
	})

	ctx := context.Background()
	enforcement, _, err := client.Repositories.AddAdminEnforcement(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.AddAdminEnforcement returned error: %v", err)
	}

	want := &AdminEnforcement{
		URL:     String("/repos/o/r/branches/b/protection/enforce_admins"),
		Enabled: true,
	}
	if !cmp.Equal(enforcement, want) {
		t.Errorf("Repositories.AddAdminEnforcement returned %+v, want %+v", enforcement, want)
	}

	const methodName = "AddAdminEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.AddAdminEnforcement(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.AddAdminEnforcement(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RemoveAdminEnforcement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/enforce_admins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.RemoveAdminEnforcement(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RemoveAdminEnforcement returned error: %v", err)
	}

	const methodName = "RemoveAdminEnforcement"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.RemoveAdminEnforcement(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.RemoveAdminEnforcement(ctx, "o", "r", "b")
	})
}

func TestRepositoriesService_GetSignaturesProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":false}`)
	})

	ctx := context.Background()
	signature, _, err := client.Repositories.GetSignaturesProtectedBranch(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.GetSignaturesProtectedBranch returned error: %v", err)
	}

	want := &SignaturesProtectedBranch{
		URL:     String("/repos/o/r/branches/b/protection/required_signatures"),
		Enabled: Bool(false),
	}

	if !cmp.Equal(signature, want) {
		t.Errorf("Repositories.GetSignaturesProtectedBranch returned %+v, want %+v", signature, want)
	}

	const methodName = "GetSignaturesProtectedBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetSignaturesProtectedBranch(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetSignaturesProtectedBranch(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RequireSignaturesOnProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":true}`)
	})

	ctx := context.Background()
	signature, _, err := client.Repositories.RequireSignaturesOnProtectedBranch(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.RequireSignaturesOnProtectedBranch returned error: %v", err)
	}

	want := &SignaturesProtectedBranch{
		URL:     String("/repos/o/r/branches/b/protection/required_signatures"),
		Enabled: Bool(true),
	}

	if !cmp.Equal(signature, want) {
		t.Errorf("Repositories.RequireSignaturesOnProtectedBranch returned %+v, want %+v", signature, want)
	}

	const methodName = "RequireSignaturesOnProtectedBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.RequireSignaturesOnProtectedBranch(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.RequireSignaturesOnProtectedBranch(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_OptionalSignaturesOnProtectedBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/required_signatures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeSignaturePreview)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.OptionalSignaturesOnProtectedBranch returned error: %v", err)
	}

	const methodName = "OptionalSignaturesOnProtectedBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "o", "r", "b")
	})
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

	ctx := context.Background()
	got, _, err := client.Repositories.ListAllTopics(ctx, "o", "r")
	if err != nil {
		t.Fatalf("Repositories.ListAllTopics returned error: %v", err)
	}

	want := []string{"go", "go-github", "github"}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ListAllTopics returned %+v, want %+v", got, want)
	}

	const methodName = "ListAllTopics"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListAllTopics(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListAllTopics(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListAllTopics_emptyTopics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		fmt.Fprint(w, `{"names":[]}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.ListAllTopics(ctx, "o", "r")
	if err != nil {
		t.Fatalf("Repositories.ListAllTopics returned error: %v", err)
	}

	want := []string{}
	if !cmp.Equal(got, want) {
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

	ctx := context.Background()
	got, _, err := client.Repositories.ReplaceAllTopics(ctx, "o", "r", []string{"go", "go-github", "github"})
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{"go", "go-github", "github"}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ReplaceAllTopics returned %+v, want %+v", got, want)
	}

	const methodName = "ReplaceAllTopics"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ReplaceAllTopics(ctx, "\n", "\n", []string{"\n", "\n", "\n"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ReplaceAllTopics(ctx, "o", "r", []string{"go", "go-github", "github"})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	got, _, err := client.Repositories.ReplaceAllTopics(ctx, "o", "r", nil)
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{}
	if !cmp.Equal(got, want) {
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

	ctx := context.Background()
	got, _, err := client.Repositories.ReplaceAllTopics(ctx, "o", "r", []string{})
	if err != nil {
		t.Fatalf("Repositories.ReplaceAllTopics returned error: %v", err)
	}

	want := []string{}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ReplaceAllTopics returned %+v, want %+v", got, want)
	}
}

func TestRepositoriesService_ListApps(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	ctx := context.Background()
	_, _, err := client.Repositories.ListApps(ctx, "o", "r", "b")
	if err != nil {
		t.Errorf("Repositories.ListApps returned error: %v", err)
	}

	const methodName = "ListApps"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListApps(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListApps(ctx, "o", "r", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	got, _, err := client.Repositories.ReplaceAppRestrictions(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.ReplaceAppRestrictions returned error: %v", err)
	}
	want := []*App{
		{Name: String("octocat")},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ReplaceAppRestrictions returned %+v, want %+v", got, want)
	}

	const methodName = "ReplaceAppRestrictions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ReplaceAppRestrictions(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ReplaceAppRestrictions(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	got, _, err := client.Repositories.AddAppRestrictions(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.AddAppRestrictions returned error: %v", err)
	}
	want := []*App{
		{Name: String("octocat")},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.AddAppRestrictions returned %+v, want %+v", got, want)
	}

	const methodName = "AddAppRestrictions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.AddAppRestrictions(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.AddAppRestrictions(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_RemoveAppRestrictions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/branches/b/protection/restrictions/apps", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `[]`)
	})
	input := []string{"octocat"}
	ctx := context.Background()
	got, _, err := client.Repositories.RemoveAppRestrictions(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Repositories.RemoveAppRestrictions returned error: %v", err)
	}
	want := []*App{}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.RemoveAppRestrictions returned %+v, want %+v", got, want)
	}

	const methodName = "RemoveAppRestrictions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.RemoveAppRestrictions(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.RemoveAppRestrictions(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_Transfer(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := TransferRequest{NewOwner: "a", TeamID: []int64{123}}

	mux.HandleFunc("/repos/o/r/transfer", func(w http.ResponseWriter, r *http.Request) {
		var v TransferRequest
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"owner":{"login":"a"}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Transfer(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.Transfer returned error: %v", err)
	}

	want := &Repository{Owner: &User{Login: String("a")}}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.Transfer returned %+v, want %+v", got, want)
	}

	const methodName = "Transfer"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.Transfer(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.Transfer(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_Dispatch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var input DispatchRequestOptions

	mux.HandleFunc("/repos/o/r/dispatches", func(w http.ResponseWriter, r *http.Request) {
		var v DispatchRequestOptions
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
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
		if !cmp.Equal(got, want) {
			t.Errorf("Repositories.Dispatch returned %+v, want %+v", got, want)
		}
	}

	const methodName = "Dispatch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.Dispatch(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.Dispatch(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAdvancedSecurity_Marshal(t *testing.T) {
	testJSONMarshal(t, &AdvancedSecurity{}, "{}")

	u := &AdvancedSecurity{
		Status: String("status"),
	}

	want := `{
		"status": "status"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuthorizedActorsOnly_Marshal(t *testing.T) {
	testJSONMarshal(t, &AuthorizedActorsOnly{}, "{}")

	u := &AuthorizedActorsOnly{
		From: Bool(true),
	}

	want := `{
		"from" : true
	}`

	testJSONMarshal(t, u, want)
}
