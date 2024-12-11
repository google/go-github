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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListByAuthenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ListByAuthenticatedUser returned %+v, want %+v", got, want)
	}

	const methodName = "ListByAuthenticatedUser"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListByUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"sort":      "created",
			"direction": "asc",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListByUserOptions{
		Sort:        "created",
		Direction:   "asc",
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	repos, _, err := client.Repositories.ListByUser(ctx, "u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []*Repository{{ID: Ptr(int64(1))}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Repositories.ListByUser returned %+v, want %+v", repos, want)
	}

	const methodName = "ListByUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListByUser(ctx, "\n", &RepositoryListByUserOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListByUser(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListByUser_type(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"type": "owner",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListByUserOptions{
		Type: "owner",
	}
	ctx := context.Background()
	repos, _, err := client.Repositories.ListByUser(ctx, "u", opt)
	if err != nil {
		t.Errorf("Repositories.ListByUser returned error: %v", err)
	}

	want := []*Repository{{ID: Ptr(int64(1))}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Repositories.ListByUser returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_ListByUser_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListByUser(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListByOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	want := []*Repository{{ID: Ptr(int64(1))}}
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
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListByOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListAll(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	want := []*Repository{{ID: Ptr(int64(1))}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Repository{
		Name:     Ptr("n"),
		Archived: Ptr(true), // not passed along.
	}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(createRepoRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		want := &createRepoRequest{Name: Ptr("n")}
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

	want := &Repository{ID: Ptr(int64(1))}
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
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Repository{
		Name:     Ptr("n"),
		Archived: Ptr(true), // not passed along.
	}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(createRepoRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		want := &createRepoRequest{Name: Ptr("n")}
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

	want := &Repository{ID: Ptr(int64(1))}
	if !cmp.Equal(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_CreateFromTemplate(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	templateRepoReq := &TemplateRepoRequest{
		Name: Ptr("n"),
	}

	mux.HandleFunc("/repos/to/tr/generate", func(w http.ResponseWriter, r *http.Request) {
		v := new(TemplateRepoRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeRepositoryTemplatePreview)
		want := &TemplateRepoRequest{Name: Ptr("n")}
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

	want := &Repository{ID: Ptr(int64(1)), Name: Ptr("n")}
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
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview, mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"},"security_and_analysis":{"advanced_security":{"status":"enabled"},"secret_scanning":{"status":"enabled"},"secret_scanning_push_protection":{"status":"enabled"},"dependabot_security_updates":{"status": "enabled"}, "secret_scanning_validity_checks":{"status":"enabled"}}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.Get(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.Get returned error: %v", err)
	}

	want := &Repository{ID: Ptr(int64(1)), Name: Ptr("n"), Description: Ptr("d"), Owner: &User{Login: Ptr("l")}, License: &License{Key: Ptr("mit")}, SecurityAndAnalysis: &SecurityAndAnalysis{AdvancedSecurity: &AdvancedSecurity{Status: Ptr("enabled")}, SecretScanning: &SecretScanning{Ptr("enabled")}, SecretScanningPushProtection: &SecretScanningPushProtection{Ptr("enabled")}, DependabotSecurityUpdates: &DependabotSecurityUpdates{Ptr("enabled")}, SecretScanningValidityChecks: &SecretScanningValidityChecks{Ptr("enabled")}}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `{
            "code_of_conduct": {
  						"key": "key",
  						"name": "name",
  						"url": "url",
  						"body": "body"
            }}`,
		)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCodeOfConduct(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCodeOfConduct returned error: %v", err)
	}

	want := &CodeOfConduct{
		Key:  Ptr("key"),
		Name: Ptr("name"),
		URL:  Ptr("url"),
		Body: Ptr("body"),
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"}}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("Repositories.GetByID returned error: %v", err)
	}

	want := &Repository{ID: Ptr(int64(1)), Name: Ptr("n"), Description: Ptr("d"), Owner: &User{Login: Ptr("l")}, License: &License{Key: Ptr("mit")}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	i := true
	input := &Repository{HasIssues: &i}

	wantAcceptHeaders := []string{mediaTypeRepositoryTemplatePreview, mediaTypeRepositoryVisibilityPreview}
	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

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

	want := &Repository{ID: Ptr(int64(1))}
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.Get(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_Edit_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.Edit(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetVulnerabilityAlerts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.EnableAutomatedSecurityFixes(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.EnableAutomatedSecurityFixes returned error: %v", err)
	}
}

func TestRepositoriesService_GetAutomatedSecurityFixes(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled": true, "paused": false}`)
	})

	ctx := context.Background()
	fixes, _, err := client.Repositories.GetAutomatedSecurityFixes(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetAutomatedSecurityFixes returned errpr: #{err}")
	}

	want := &AutomatedSecurityFixes{
		Enabled: Ptr(true),
		Paused:  Ptr(false),
	}
	if !cmp.Equal(fixes, want) {
		t.Errorf("Repositories.GetAutomatedSecurityFixes returned #{fixes}, want #{want}")
	}

	const methodName = "GetAutomatedSecurityFixes"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetAutomatedSecurityFixes(ctx, "\n", "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAutomatedSecurityFixes(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DisableAutomatedSecurityFixes(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/automated-security-fixes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.DisableAutomatedSecurityFixes(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.DisableAutomatedSecurityFixes returned error: %v", err)
	}
}

func TestRepositoriesService_ListContributors(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	want := []*Contributor{{Contributions: Ptr(42)}}
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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

	want := []*Team{{ID: Ptr(int64(1))}}
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
	t.Parallel()
	client, mux, _ := setup(t)

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
			Name: Ptr("n"),
			Commit: &Commit{
				SHA: Ptr("s"),
				URL: Ptr("u"),
			},
			ZipballURL: Ptr("z"),
			TarballURL: Ptr("t"),
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
	t.Parallel()
	client, mux, _ := setup(t)

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

	want := []*Branch{{Name: Ptr("master"), Commit: &RepositoryCommit{SHA: Ptr("a57781"), URL: Ptr("https://api.github.com/repos/o/r/commits/a57781")}}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%"},
	}

	for _, test := range tests {
		mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"name":"n", "commit":{"sha":"s","commit":{"message":"m"}}, "protected":true, "protection":{"required_status_checks":{"contexts":["c"]}}}`)
		})

		ctx := context.Background()
		branch, _, err := client.Repositories.GetBranch(ctx, "o", "r", test.branch, 0)
		if err != nil {
			t.Errorf("Repositories.GetBranch returned error: %v", err)
		}

		want := &Branch{
			Name: Ptr("n"),
			Commit: &RepositoryCommit{
				SHA: Ptr("s"),
				Commit: &Commit{
					Message: Ptr("m"),
				},
			},
			Protected: Ptr(true),
			Protection: &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Contexts: &[]string{"c"},
				},
			},
		}

		if !cmp.Equal(branch, want) {
			t.Errorf("Repositories.GetBranch returned %+v, want %+v", branch, want)
		}

		const methodName = "GetBranch"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.Repositories.GetBranch(ctx, "\n", "\n", "\n", 0)
			return err
		})
	}
}

func TestRepositoriesService_GetBranch_BadJSONResponse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fmt.Fprint(w, `{"name":"n", "commit":{"sha":...truncated`)
			})

			ctx := context.Background()
			if _, _, err := client.Repositories.GetBranch(ctx, "o", "r", test.branch, 0); err == nil {
				t.Error("Repositories.GetBranch returned no error; wanted JSON error")
			}
		})
	}
}

func TestRepositoriesService_GetBranch_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	mux.HandleFunc("/repos/o/r/branches/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/repos/o/r/branches/br")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})
	mux.HandleFunc("/repos/o/r/branches/br", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"n", "commit":{"sha":"s","commit":{"message":"m"}}, "protected":true, "protection":{"required_status_checks":{"contexts":["c"]}}}`)
	})
	ctx := context.Background()
	branch, resp, err := client.Repositories.GetBranch(ctx, "o", "r", "b", 1)
	if err != nil {
		t.Errorf("Repositories.GetBranch returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.GetBranch returned status: %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := &Branch{
		Name: Ptr("n"),
		Commit: &RepositoryCommit{
			SHA: Ptr("s"),
			Commit: &Commit{
				Message: Ptr("m"),
			},
		},
		Protected: Ptr(true),
		Protection: &Protection{
			RequiredStatusChecks: &RequiredStatusChecks{
				Contexts: &[]string{"c"},
			},
		},
	}
	if !cmp.Equal(branch, want) {
		t.Errorf("Repositories.GetBranch returned %+v, want %+v", branch, want)
	}
}

func TestRepositoriesService_GetBranch_notFound(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat-branch-50%"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Error(w, "branch not found", http.StatusNotFound)
			})
			ctx := context.Background()
			_, resp, err := client.Repositories.GetBranch(ctx, "o", "r", test.branch, 1)
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
				_, _, err = client.Repositories.GetBranch(ctx, "\n", "\n", "\n", 1)
				return err
			})
		})
	}
}

func TestRepositoriesService_RenameBranch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/rename"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/rename"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			renameBranchReq := "nn"

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(renameBranchRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "POST")
				want := &renameBranchRequest{NewName: renameBranchReq}
				if !cmp.Equal(v, want) {
					t.Errorf("Request body = %+v, want %+v", v, want)
				}

				fmt.Fprint(w, `{"protected":true,"name":"nn"}`)
			})

			ctx := context.Background()
			got, _, err := client.Repositories.RenameBranch(ctx, "o", "r", test.branch, renameBranchReq)
			if err != nil {
				t.Errorf("Repositories.RenameBranch returned error: %v", err)
			}

			want := &Branch{Name: Ptr("nn"), Protected: Ptr(true)}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.RenameBranch returned %+v, want %+v", got, want)
			}

			const methodName = "RenameBranch"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.RenameBranch(ctx, "\n", "\n", "\n", renameBranchReq)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.RenameBranch(ctx, "o", "r", test.branch, renameBranchReq)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_GetBranchProtection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch               string
		urlPath              string
		enforceAdminsURLPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection", enforceAdminsURLPath: "/repos/o/r/branches/b/protection/enforce_admins"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection", enforceAdminsURLPath: "/repos/o/r/branches/feat/branch-50%/protection/enforce_admins"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
						"required_status_checks":{
							"strict":true,
							"contexts":["continuous-integration"],
							"checks": [
								{
									"context": "continuous-integration",
									"app_id": null
								}
							]
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
								}],
								"apps":[{
									"id":5,
									"slug":"a"
								}]
							},
							"dismiss_stale_reviews":true,
							"require_code_owner_reviews":true,
							"require_last_push_approval":false,
							"required_approving_review_count":1
							},
							"enforce_admins":{
								"url":"%s",
								"enabled":true
							},
							"restrictions":{
								"users":[{"id":1,"login":"u"}],
								"teams":[{"id":2,"slug":"t"}],
								"apps":[{"id":3,"slug":"a"}]
							},
							"required_conversation_resolution": {
								"enabled": true
							},
							"block_creations": {
								"enabled": false
							},
							"lock_branch": {
								"enabled": false
							},
							"allow_fork_syncing": {
								"enabled": false
							}
						}`, test.enforceAdminsURLPath)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.GetBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{"continuous-integration"},
					Checks: &[]*RequiredStatusCheck{
						{
							Context: "continuous-integration",
						},
					},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("u"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("t"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("a"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews:      true,
					RequiredApprovingReviewCount: 1,
					RequireLastPushApproval:      false,
				},
				EnforceAdmins: &AdminEnforcement{
					URL:     Ptr(test.enforceAdminsURLPath),
					Enabled: true,
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
				RequiredConversationResolution: &RequiredConversationResolution{
					Enabled: true,
				},
				BlockCreations: &BlockCreations{
					Enabled: Ptr(false),
				},
				LockBranch: &LockBranch{
					Enabled: Ptr(false),
				},
				AllowForkSyncing: &AllowForkSyncing{
					Enabled: Ptr(false),
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
				got, resp, err := client.Repositories.GetBranchProtection(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_GetBranchProtection_noDismissalRestrictions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	tests := []struct {
		branch               string
		urlPath              string
		enforceAdminsURLPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection", enforceAdminsURLPath: "/repos/o/r/branches/b/protection/enforce_admins"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection", enforceAdminsURLPath: "/repos/o/r/branches/feat/branch-50%/protection/enforce_admins"},
	}

	for _, test := range tests {
		mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			// TODO: remove custom Accept header when this API fully launches
			testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
			fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":["continuous-integration"],
						"checks": [
							{
								"context": "continuous-integration",
								"app_id": null
							}
						]
					},
					"required_pull_request_reviews":{
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"required_approving_review_count":1
						},
						"enforce_admins":{
							"url":"%s",
							"enabled":true
						},
						"restrictions":{
							"users":[{"id":1,"login":"u"}],
							"teams":[{"id":2,"slug":"t"}]
						}
					}`, test.enforceAdminsURLPath)
		})

		ctx := context.Background()
		protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", test.branch)
		if err != nil {
			t.Errorf("Repositories.GetBranchProtection returned error: %v", err)
		}

		want := &Protection{
			RequiredStatusChecks: &RequiredStatusChecks{
				Strict:   true,
				Contexts: &[]string{"continuous-integration"},
				Checks: &[]*RequiredStatusCheck{
					{
						Context: "continuous-integration",
					},
				},
			},
			RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
				DismissStaleReviews:          true,
				DismissalRestrictions:        nil,
				RequireCodeOwnerReviews:      true,
				RequiredApprovingReviewCount: 1,
			},
			EnforceAdmins: &AdminEnforcement{
				URL:     Ptr(test.enforceAdminsURLPath),
				Enabled: true,
			},
			Restrictions: &BranchRestrictions{
				Users: []*User{
					{Login: Ptr("u"), ID: Ptr(int64(1))},
				},
				Teams: []*Team{
					{Slug: Ptr("t"), ID: Ptr(int64(2))},
				},
			},
		}
		if !cmp.Equal(protection, want) {
			t.Errorf("Repositories.GetBranchProtection returned %+v, want %+v", protection, want)
		}
	}
}

func TestRepositoriesService_GetBranchProtection_branchNotProtected(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")

				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{
					"message": %q,
					"documentation_url": "https://docs.github.com/rest/repos#get-branch-protection"
					}`, githubBranchNotProtected)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.GetBranchProtection(ctx, "o", "r", test.branch)

			if protection != nil {
				t.Errorf("Repositories.GetBranchProtection returned non-nil protection data")
			}

			if err != ErrBranchNotProtected {
				t.Errorf("Repositories.GetBranchProtection returned an invalid error: %v", err)
			}
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_Contexts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{"continuous-integration"},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					DismissStaleReviews: true,
					DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
						Users: &[]string{"uu"},
						Teams: &[]string{"tt"},
						Apps:  &[]string{"aa"},
					},
					BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
						Users: []string{"uuu"},
						Teams: []string{"ttt"},
						Apps:  []string{"aaa"},
					},
				},
				Restrictions: &BranchRestrictionsRequest{
					Users: []string{"u"},
					Teams: []string{"t"},
					Apps:  []string{"a"},
				},
				BlockCreations:   Ptr(true),
				LockBranch:       Ptr(true),
				AllowForkSyncing: Ptr(true),
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":["continuous-integration"],
						"checks": [
							{
								"context": "continuous-integration",
								"app_id": null
							}
						]
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
							}],
							"apps":[{
								"id":5,
								"slug":"aa"
							}]
						},
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"bypass_pull_request_allowances": {
							"users":[{"id":10,"login":"uuu"}],
							"teams":[{"id":20,"slug":"ttt"}],
							"apps":[{"id":30,"slug":"aaa"}]
						}
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					},
					"block_creations": {
						"enabled": true
					},
					"lock_branch": {
						"enabled": true
					},
					"allow_fork_syncing": {
						"enabled": true
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{"continuous-integration"},
					Checks: &[]*RequiredStatusCheck{
						{
							Context: "continuous-integration",
						},
					},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("uu"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("tt"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("aa"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews: true,
					BypassPullRequestAllowances: &BypassPullRequestAllowances{
						Users: []*User{
							{Login: Ptr("uuu"), ID: Ptr(int64(10))},
						},
						Teams: []*Team{
							{Slug: Ptr("ttt"), ID: Ptr(int64(20))},
						},
						Apps: []*App{
							{Slug: Ptr("aaa"), ID: Ptr(int64(30))},
						},
					},
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
				BlockCreations: &BlockCreations{
					Enabled: Ptr(true),
				},
				LockBranch: &LockBranch{
					Enabled: Ptr(true),
				},
				AllowForkSyncing: &AllowForkSyncing{
					Enabled: Ptr(true),
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
				got, resp, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_EmptyContexts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					DismissStaleReviews: true,
					DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
						Users: &[]string{"uu"},
						Teams: &[]string{"tt"},
						Apps:  &[]string{"aa"},
					},
					BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
						Users: []string{"uuu"},
						Teams: []string{"ttt"},
						Apps:  []string{"aaa"},
					},
				},
				Restrictions: &BranchRestrictionsRequest{
					Users: []string{"u"},
					Teams: []string{"t"},
					Apps:  []string{"a"},
				},
				BlockCreations:   Ptr(true),
				LockBranch:       Ptr(true),
				AllowForkSyncing: Ptr(true),
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":[],
						"checks": null
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
							}],
							"apps":[{
								"id":5,
								"slug":"aa"
							}]
						},
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"bypass_pull_request_allowances": {
							"users":[{"id":10,"login":"uuu"}],
							"teams":[{"id":20,"slug":"ttt"}],
							"apps":[{"id":30,"slug":"aaa"}]
						}
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					},
					"block_creations": {
						"enabled": true
					},
					"lock_branch": {
						"enabled": true
					},
					"allow_fork_syncing": {
						"enabled": true
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("uu"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("tt"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("aa"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews: true,
					BypassPullRequestAllowances: &BypassPullRequestAllowances{
						Users: []*User{
							{Login: Ptr("uuu"), ID: Ptr(int64(10))},
						},
						Teams: []*Team{
							{Slug: Ptr("ttt"), ID: Ptr(int64(20))},
						},
						Apps: []*App{
							{Slug: Ptr("aaa"), ID: Ptr(int64(30))},
						},
					},
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
				BlockCreations: &BlockCreations{
					Enabled: Ptr(true),
				},
				LockBranch: &LockBranch{
					Enabled: Ptr(true),
				},
				AllowForkSyncing: &AllowForkSyncing{
					Enabled: Ptr(true),
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
				got, resp, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_Checks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict: true,
					Checks: &[]*RequiredStatusCheck{
						{
							Context: "continuous-integration",
						},
					},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					DismissStaleReviews: true,
					DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
						Users: &[]string{"uu"},
						Teams: &[]string{"tt"},
						Apps:  &[]string{"aa"},
					},
					BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
						Users: []string{"uuu"},
						Teams: []string{"ttt"},
						Apps:  []string{"aaa"},
					},
				},
				Restrictions: &BranchRestrictionsRequest{
					Users: []string{"u"},
					Teams: []string{"t"},
					Apps:  []string{"a"},
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":["continuous-integration"],
						"checks": [
							{
								"context": "continuous-integration",
								"app_id": null
							}
						]
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
							}],
							"apps":[{
								"id":5,
								"slug":"aa"
							}]
						},
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"bypass_pull_request_allowances": {
							"users":[{"id":10,"login":"uuu"}],
							"teams":[{"id":20,"slug":"ttt"}],
							"apps":[{"id":30,"slug":"aaa"}]
						}
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{"continuous-integration"},
					Checks: &[]*RequiredStatusCheck{
						{
							Context: "continuous-integration",
						},
					},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("uu"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("tt"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("aa"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews: true,
					BypassPullRequestAllowances: &BypassPullRequestAllowances{
						Users: []*User{
							{Login: Ptr("uuu"), ID: Ptr(int64(10))},
						},
						Teams: []*Team{
							{Slug: Ptr("ttt"), ID: Ptr(int64(20))},
						},
						Apps: []*App{
							{Slug: Ptr("aaa"), ID: Ptr(int64(30))},
						},
					},
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
			}
			if !cmp.Equal(protection, want) {
				t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
			}
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_EmptyChecks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict: true,
					Checks: &[]*RequiredStatusCheck{},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					DismissStaleReviews: true,
					DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
						Users: &[]string{"uu"},
						Teams: &[]string{"tt"},
						Apps:  &[]string{"aa"},
					},
					BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
						Users: []string{"uuu"},
						Teams: []string{"ttt"},
						Apps:  []string{"aaa"},
					},
				},
				Restrictions: &BranchRestrictionsRequest{
					Users: []string{"u"},
					Teams: []string{"t"},
					Apps:  []string{"a"},
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":null,
						"checks": []
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
							}],
							"apps":[{
								"id":5,
								"slug":"aa"
							}]
						},
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"bypass_pull_request_allowances": {
							"users":[{"id":10,"login":"uuu"}],
							"teams":[{"id":20,"slug":"ttt"}],
							"apps":[{"id":30,"slug":"aaa"}]
						}
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict: true,
					Checks: &[]*RequiredStatusCheck{},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("uu"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("tt"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("aa"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews: true,
					BypassPullRequestAllowances: &BypassPullRequestAllowances{
						Users: []*User{
							{Login: Ptr("uuu"), ID: Ptr(int64(10))},
						},
						Teams: []*Team{
							{Slug: Ptr("ttt"), ID: Ptr(int64(20))},
						},
						Apps: []*App{
							{Slug: Ptr("aaa"), ID: Ptr(int64(30))},
						},
					},
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
			}
			if !cmp.Equal(protection, want) {
				t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
			}
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_StrictNoChecks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict: true,
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					DismissStaleReviews: true,
					DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
						Users: &[]string{"uu"},
						Teams: &[]string{"tt"},
						Apps:  &[]string{"aa"},
					},
					BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
						Users: []string{"uuu"},
						Teams: []string{"ttt"},
						Apps:  []string{"aaa"},
					},
				},
				Restrictions: &BranchRestrictionsRequest{
					Users: []string{"u"},
					Teams: []string{"t"},
					Apps:  []string{"a"},
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"required_status_checks":{
						"strict":true,
						"contexts":[]
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
							}],
							"apps":[{
								"id":5,
								"slug":"aa"
							}]
						},
						"dismiss_stale_reviews":true,
						"require_code_owner_reviews":true,
						"require_last_push_approval":false,
						"bypass_pull_request_allowances": {
							"users":[{"id":10,"login":"uuu"}],
							"teams":[{"id":20,"slug":"ttt"}],
							"apps":[{"id":30,"slug":"aaa"}]
						}
					},
					"restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredStatusChecks: &RequiredStatusChecks{
					Strict:   true,
					Contexts: &[]string{},
				},
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					DismissStaleReviews: true,
					DismissalRestrictions: &DismissalRestrictions{
						Users: []*User{
							{Login: Ptr("uu"), ID: Ptr(int64(3))},
						},
						Teams: []*Team{
							{Slug: Ptr("tt"), ID: Ptr(int64(4))},
						},
						Apps: []*App{
							{Slug: Ptr("aa"), ID: Ptr(int64(5))},
						},
					},
					RequireCodeOwnerReviews: true,
					BypassPullRequestAllowances: &BypassPullRequestAllowances{
						Users: []*User{
							{Login: Ptr("uuu"), ID: Ptr(int64(10))},
						},
						Teams: []*Team{
							{Slug: Ptr("ttt"), ID: Ptr(int64(20))},
						},
						Apps: []*App{
							{Slug: Ptr("aaa"), ID: Ptr(int64(30))},
						},
					},
				},
				Restrictions: &BranchRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
					},
				},
			}
			if !cmp.Equal(protection, want) {
				t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
			}
		})
	}
}

func TestRepositoriesService_UpdateBranchProtection_RequireLastPushApproval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &ProtectionRequest{
				RequiredPullRequestReviews: &PullRequestReviewsEnforcementRequest{
					RequireLastPushApproval: Ptr(true),
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(ProtectionRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PUT")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}

				fmt.Fprintf(w, `{
					"required_pull_request_reviews":{
						"require_last_push_approval":true
					}
				}`)
			})

			ctx := context.Background()
			protection, _, err := client.Repositories.UpdateBranchProtection(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateBranchProtection returned error: %v", err)
			}

			want := &Protection{
				RequiredPullRequestReviews: &PullRequestReviewsEnforcement{
					RequireLastPushApproval: true,
				},
			}
			if !cmp.Equal(protection, want) {
				t.Errorf("Repositories.UpdateBranchProtection returned %+v, want %+v", protection, want)
			}
		})
	}
}

func TestRepositoriesService_RemoveBranchProtection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := context.Background()
			_, err := client.Repositories.RemoveBranchProtection(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.RemoveBranchProtection returned error: %v", err)
			}

			const methodName = "RemoveBranchProtection"
			testBadOptions(t, methodName, func() (err error) {
				_, err = client.Repositories.RemoveBranchProtection(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				return client.Repositories.RemoveBranchProtection(ctx, "o", "r", test.branch)
			})
		})
	}
}

func TestRepositoriesService_ListLanguages_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListLanguages(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestRepositoriesService_License(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
		Name: Ptr("LICENSE"),
		Path: Ptr("LICENSE"),
		License: &License{
			Name:     Ptr("MIT License"),
			Key:      Ptr("mit"),
			SPDXID:   Ptr("MIT"),
			URL:      Ptr("https://api.github.com/licenses/mit"),
			Featured: Ptr(true),
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
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fmt.Fprint(w, `{
					"strict": true,
					"contexts": ["x","y","z"],
					"checks": [
						{
							"context": "x",
							"app_id": null
						},
						{
							"context": "y",
							"app_id": null
						},
						{
							"context": "z",
							"app_id": null
						}
					]
				}`)
			})

			ctx := context.Background()
			checks, _, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.GetRequiredStatusChecks returned error: %v", err)
			}

			want := &RequiredStatusChecks{
				Strict:   true,
				Contexts: &[]string{"x", "y", "z"},
				Checks: &[]*RequiredStatusCheck{
					{
						Context: "x",
					},
					{
						Context: "y",
					},
					{
						Context: "z",
					},
				},
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
				got, resp, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_GetRequiredStatusChecks_branchNotProtected(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")

				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{
			"message": %q,
			"documentation_url": "https://docs.github.com/rest/repos#get-branch-protection"
			}`, githubBranchNotProtected)
			})

			ctx := context.Background()
			checks, _, err := client.Repositories.GetRequiredStatusChecks(ctx, "o", "r", test.branch)

			if checks != nil {
				t.Errorf("Repositories.GetRequiredStatusChecks returned non-nil status-checks data")
			}

			if err != ErrBranchNotProtected {
				t.Errorf("Repositories.GetRequiredStatusChecks returned an invalid error: %v", err)
			}
		})
	}
}

func TestRepositoriesService_UpdateRequiredStatusChecks_Contexts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &RequiredStatusChecksRequest{
				Strict:   Ptr(true),
				Contexts: []string{"continuous-integration"},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(RequiredStatusChecksRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PATCH")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}
				testHeader(t, r, "Accept", mediaTypeV3)
				fmt.Fprintf(w, `{
					"strict":true,
					"contexts":["continuous-integration"],
					"checks": [
						{
							"context": "continuous-integration",
							"app_id": null
						}
					]
				}`)
			})

			ctx := context.Background()
			statusChecks, _, err := client.Repositories.UpdateRequiredStatusChecks(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateRequiredStatusChecks returned error: %v", err)
			}

			want := &RequiredStatusChecks{
				Strict:   true,
				Contexts: &[]string{"continuous-integration"},
				Checks: &[]*RequiredStatusCheck{
					{
						Context: "continuous-integration",
					},
				},
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
				got, resp, err := client.Repositories.UpdateRequiredStatusChecks(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_UpdateRequiredStatusChecks_Checks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			appID := int64(123)
			noAppID := int64(-1)
			input := &RequiredStatusChecksRequest{
				Strict: Ptr(true),
				Checks: []*RequiredStatusCheck{
					{
						Context: "continuous-integration",
					},
					{
						Context: "continuous-integration2",
						AppID:   &appID,
					},
					{
						Context: "continuous-integration3",
						AppID:   &noAppID,
					},
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(RequiredStatusChecksRequest)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PATCH")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}
				testHeader(t, r, "Accept", mediaTypeV3)
				fmt.Fprintf(w, `{
					"strict":true,
					"contexts":["continuous-integration"],
					"checks": [
						{
							"context": "continuous-integration",
							"app_id": null
						},
						{
							"context": "continuous-integration2",
							"app_id": 123
						},
						{
							"context": "continuous-integration3",
							"app_id": null
						}
					]
				}`)
			})

			ctx := context.Background()
			statusChecks, _, err := client.Repositories.UpdateRequiredStatusChecks(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdateRequiredStatusChecks returned error: %v", err)
			}

			want := &RequiredStatusChecks{
				Strict:   true,
				Contexts: &[]string{"continuous-integration"},
				Checks: &[]*RequiredStatusCheck{
					{
						Context: "continuous-integration",
					},
					{
						Context: "continuous-integration2",
						AppID:   &appID,
					},
					{
						Context: "continuous-integration3",
					},
				},
			}
			if !cmp.Equal(statusChecks, want) {
				t.Errorf("Repositories.UpdateRequiredStatusChecks returned %+v, want %+v", statusChecks, want)
			}
		})
	}
}

func TestRepositoriesService_RemoveRequiredStatusChecks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				testHeader(t, r, "Accept", mediaTypeV3)
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := context.Background()
			_, err := client.Repositories.RemoveRequiredStatusChecks(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.RemoveRequiredStatusChecks returned error: %v", err)
			}

			const methodName = "RemoveRequiredStatusChecks"
			testBadOptions(t, methodName, func() (err error) {
				_, err = client.Repositories.RemoveRequiredStatusChecks(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				return client.Repositories.RemoveRequiredStatusChecks(ctx, "o", "r", test.branch)
			})
		})
	}
}

func TestRepositoriesService_ListRequiredStatusChecksContexts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks/contexts"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks/contexts"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fmt.Fprint(w, `["x", "y", "z"]`)
			})

			ctx := context.Background()
			contexts, _, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", test.branch)
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
				got, resp, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ListRequiredStatusChecksContexts_branchNotProtected(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_status_checks/contexts"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_status_checks/contexts"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")

				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{
			"message": %q,
			"documentation_url": "https://docs.github.com/rest/repos#get-branch-protection"
			}`, githubBranchNotProtected)
			})

			ctx := context.Background()
			contexts, _, err := client.Repositories.ListRequiredStatusChecksContexts(ctx, "o", "r", test.branch)

			if contexts != nil {
				t.Errorf("Repositories.ListRequiredStatusChecksContexts returned non-nil contexts data")
			}

			if err != ErrBranchNotProtected {
				t.Errorf("Repositories.ListRequiredStatusChecksContexts returned an invalid error: %v", err)
			}
		})
	}
}

func TestRepositoriesService_GetPullRequestReviewEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_pull_request_reviews"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_pull_request_reviews"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
			"dismissal_restrictions":{
				"users":[{"id":1,"login":"u"}],
				"teams":[{"id":2,"slug":"t"}],
				"apps":[{"id":3,"slug":"a"}]
			},
			"dismiss_stale_reviews":true,
			"require_code_owner_reviews":true,
			"required_approving_review_count":1
		}`)
			})

			ctx := context.Background()
			enforcement, _, err := client.Repositories.GetPullRequestReviewEnforcement(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.GetPullRequestReviewEnforcement returned error: %v", err)
			}

			want := &PullRequestReviewsEnforcement{
				DismissStaleReviews: true,
				DismissalRestrictions: &DismissalRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
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
				got, resp, err := client.Repositories.GetPullRequestReviewEnforcement(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_UpdatePullRequestReviewEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_pull_request_reviews"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_pull_request_reviews"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			input := &PullRequestReviewsEnforcementUpdate{
				DismissalRestrictionsRequest: &DismissalRestrictionsRequest{
					Users: &[]string{"u"},
					Teams: &[]string{"t"},
					Apps:  &[]string{"a"},
				},
			}

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				v := new(PullRequestReviewsEnforcementUpdate)
				assertNilError(t, json.NewDecoder(r.Body).Decode(v))

				testMethod(t, r, "PATCH")
				if !cmp.Equal(v, input) {
					t.Errorf("Request body = %+v, want %+v", v, input)
				}
				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				fmt.Fprintf(w, `{
					"dismissal_restrictions":{
						"users":[{"id":1,"login":"u"}],
						"teams":[{"id":2,"slug":"t"}],
						"apps":[{"id":3,"slug":"a"}]
					},
					"dismiss_stale_reviews":true,
					"require_code_owner_reviews":true,
					"required_approving_review_count":3
				}`)
			})

			ctx := context.Background()
			enforcement, _, err := client.Repositories.UpdatePullRequestReviewEnforcement(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.UpdatePullRequestReviewEnforcement returned error: %v", err)
			}

			want := &PullRequestReviewsEnforcement{
				DismissStaleReviews: true,
				DismissalRestrictions: &DismissalRestrictions{
					Users: []*User{
						{Login: Ptr("u"), ID: Ptr(int64(1))},
					},
					Teams: []*Team{
						{Slug: Ptr("t"), ID: Ptr(int64(2))},
					},
					Apps: []*App{
						{Slug: Ptr("a"), ID: Ptr(int64(3))},
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
				got, resp, err := client.Repositories.UpdatePullRequestReviewEnforcement(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_DisableDismissalRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_pull_request_reviews"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_pull_request_reviews"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PATCH")
				// TODO: remove custom Accept header when this API fully launches
				testHeader(t, r, "Accept", mediaTypeRequiredApprovingReviewsPreview)
				testBody(t, r, `{"dismissal_restrictions":{}}`+"\n")
				fmt.Fprintf(w, `{"dismiss_stale_reviews":true,"require_code_owner_reviews":true,"required_approving_review_count":1}`)
			})

			ctx := context.Background()
			enforcement, _, err := client.Repositories.DisableDismissalRestrictions(ctx, "o", "r", test.branch)
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
				got, resp, err := client.Repositories.DisableDismissalRestrictions(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RemovePullRequestReviewEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_pull_request_reviews"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_pull_request_reviews"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := context.Background()
			_, err := client.Repositories.RemovePullRequestReviewEnforcement(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.RemovePullRequestReviewEnforcement returned error: %v", err)
			}

			const methodName = "RemovePullRequestReviewEnforcement"
			testBadOptions(t, methodName, func() (err error) {
				_, err = client.Repositories.RemovePullRequestReviewEnforcement(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				return client.Repositories.RemovePullRequestReviewEnforcement(ctx, "o", "r", test.branch)
			})
		})
	}
}

func TestRepositoriesService_GetAdminEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/enforce_admins"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/enforce_admins"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
			})

			ctx := context.Background()
			enforcement, _, err := client.Repositories.GetAdminEnforcement(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.GetAdminEnforcement returned error: %v", err)
			}

			want := &AdminEnforcement{
				URL:     Ptr("/repos/o/r/branches/b/protection/enforce_admins"),
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
				got, resp, err := client.Repositories.GetAdminEnforcement(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_AddAdminEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/enforce_admins"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/enforce_admins"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/enforce_admins","enabled":true}`)
			})

			ctx := context.Background()
			enforcement, _, err := client.Repositories.AddAdminEnforcement(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.AddAdminEnforcement returned error: %v", err)
			}

			want := &AdminEnforcement{
				URL:     Ptr("/repos/o/r/branches/b/protection/enforce_admins"),
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
				got, resp, err := client.Repositories.AddAdminEnforcement(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RemoveAdminEnforcement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/enforce_admins"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/enforce_admins"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := context.Background()
			_, err := client.Repositories.RemoveAdminEnforcement(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.RemoveAdminEnforcement returned error: %v", err)
			}

			const methodName = "RemoveAdminEnforcement"
			testBadOptions(t, methodName, func() (err error) {
				_, err = client.Repositories.RemoveAdminEnforcement(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				return client.Repositories.RemoveAdminEnforcement(ctx, "o", "r", test.branch)
			})
		})
	}
}

func TestRepositoriesService_GetSignaturesProtectedBranch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_signatures"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_signatures"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testHeader(t, r, "Accept", mediaTypeSignaturePreview)
				fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":false}`)
			})

			ctx := context.Background()
			signature, _, err := client.Repositories.GetSignaturesProtectedBranch(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.GetSignaturesProtectedBranch returned error: %v", err)
			}

			want := &SignaturesProtectedBranch{
				URL:     Ptr("/repos/o/r/branches/b/protection/required_signatures"),
				Enabled: Ptr(false),
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
				got, resp, err := client.Repositories.GetSignaturesProtectedBranch(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RequireSignaturesOnProtectedBranch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_signatures"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_signatures"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				testHeader(t, r, "Accept", mediaTypeSignaturePreview)
				fmt.Fprintf(w, `{"url":"/repos/o/r/branches/b/protection/required_signatures","enabled":true}`)
			})

			ctx := context.Background()
			signature, _, err := client.Repositories.RequireSignaturesOnProtectedBranch(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.RequireSignaturesOnProtectedBranch returned error: %v", err)
			}

			want := &SignaturesProtectedBranch{
				URL:     Ptr("/repos/o/r/branches/b/protection/required_signatures"),
				Enabled: Ptr(true),
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
				got, resp, err := client.Repositories.RequireSignaturesOnProtectedBranch(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_OptionalSignaturesOnProtectedBranch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/required_signatures"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/required_signatures"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				testHeader(t, r, "Accept", mediaTypeSignaturePreview)
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := context.Background()
			_, err := client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.OptionalSignaturesOnProtectedBranch returned error: %v", err)
			}

			const methodName = "OptionalSignaturesOnProtectedBranch"
			testBadOptions(t, methodName, func() (err error) {
				_, err = client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				return client.Repositories.OptionalSignaturesOnProtectedBranch(ctx, "o", "r", test.branch)
			})
		})
	}
}

func TestPullRequestReviewsEnforcementRequest_MarshalJSON_nilDismissalRestirctions(t *testing.T) {
	t.Parallel()
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
			Apps:  &[]string{},
		},
		RequireLastPushApproval: Ptr(true),
	}

	got, err = json.Marshal(req)
	if err != nil {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned error: %v", err)
	}

	want = `{"dismissal_restrictions":{"users":[],"teams":[],"apps":[]},"dismiss_stale_reviews":false,"require_code_owner_reviews":false,"required_approving_review_count":0,"require_last_push_approval":true}`
	if want != string(got) {
		t.Errorf("PullRequestReviewsEnforcementRequest.MarshalJSON returned %+v, want %+v", string(got), want)
	}
}

func TestRepositoriesService_ListAllTopics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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

func TestRepositoriesService_ListAppRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/apps"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/apps"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
			})

			ctx := context.Background()
			_, _, err := client.Repositories.ListAppRestrictions(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.ListAppRestrictions returned error: %v", err)
			}

			const methodName = "ListAppRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.ListAppRestrictions(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.ListAppRestrictions(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ReplaceAppRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/apps"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/apps"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PUT")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.ReplaceAppRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.ReplaceAppRestrictions returned error: %v", err)
			}
			want := []*App{
				{Name: Ptr("octocat")},
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
				got, resp, err := client.Repositories.ReplaceAppRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_AddAppRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/apps"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/apps"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.AddAppRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.AddAppRestrictions returned error: %v", err)
			}
			want := []*App{
				{Name: Ptr("octocat")},
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
				got, resp, err := client.Repositories.AddAppRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RemoveAppRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/apps"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/apps"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				fmt.Fprint(w, `[]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.RemoveAppRestrictions(ctx, "o", "r", test.branch, input)
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
				got, resp, err := client.Repositories.RemoveAppRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ListTeamRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/teams"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/teams"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
			})

			ctx := context.Background()
			_, _, err := client.Repositories.ListTeamRestrictions(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.ListTeamRestrictions returned error: %v", err)
			}

			const methodName = "ListTeamRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.ListTeamRestrictions(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.ListTeamRestrictions(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ReplaceTeamRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/teams"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/teams"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PUT")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.ReplaceTeamRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.ReplaceTeamRestrictions returned error: %v", err)
			}
			want := []*Team{
				{Name: Ptr("octocat")},
			}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.ReplaceTeamRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "ReplaceTeamRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.ReplaceTeamRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.ReplaceTeamRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_AddTeamRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/teams"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/teams"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.AddTeamRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.AddTeamRestrictions returned error: %v", err)
			}
			want := []*Team{
				{Name: Ptr("octocat")},
			}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.AddTeamRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "AddTeamRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.AddTeamRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.AddTeamRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RemoveTeamRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/teams"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/teams"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				fmt.Fprint(w, `[]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.RemoveTeamRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.RemoveTeamRestrictions returned error: %v", err)
			}
			want := []*Team{}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.RemoveTeamRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "RemoveTeamRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.RemoveTeamRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.RemoveTeamRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ListUserRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/users"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/users"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
			})

			ctx := context.Background()
			_, _, err := client.Repositories.ListUserRestrictions(ctx, "o", "r", test.branch)
			if err != nil {
				t.Errorf("Repositories.ListUserRestrictions returned error: %v", err)
			}

			const methodName = "ListUserRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.ListUserRestrictions(ctx, "\n", "\n", "\n")
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.ListUserRestrictions(ctx, "o", "r", test.branch)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_ReplaceUserRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/users"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/users"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PUT")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.ReplaceUserRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.ReplaceUserRestrictions returned error: %v", err)
			}
			want := []*User{
				{Name: Ptr("octocat")},
			}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.ReplaceUserRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "ReplaceUserRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.ReplaceUserRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.ReplaceUserRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_AddUserRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/users"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/users"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprint(w, `[{
				"name": "octocat"
			}]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.AddUserRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.AddUserRestrictions returned error: %v", err)
			}
			want := []*User{
				{Name: Ptr("octocat")},
			}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.AddUserRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "AddUserRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.AddUserRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.AddUserRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_RemoveUserRestrictions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch  string
		urlPath string
	}{
		{branch: "b", urlPath: "/repos/o/r/branches/b/protection/restrictions/users"},
		{branch: "feat/branch-50%", urlPath: "/repos/o/r/branches/feat/branch-50%/protection/restrictions/users"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.branch, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc(test.urlPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "DELETE")
				fmt.Fprint(w, `[]`)
			})
			input := []string{"octocat"}
			ctx := context.Background()
			got, _, err := client.Repositories.RemoveUserRestrictions(ctx, "o", "r", test.branch, input)
			if err != nil {
				t.Errorf("Repositories.RemoveUserRestrictions returned error: %v", err)
			}
			want := []*User{}
			if !cmp.Equal(got, want) {
				t.Errorf("Repositories.RemoveUserRestrictions returned %+v, want %+v", got, want)
			}

			const methodName = "RemoveUserRestrictions"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Repositories.RemoveUserRestrictions(ctx, "\n", "\n", "\n", input)
				return err
			})

			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Repositories.RemoveUserRestrictions(ctx, "o", "r", test.branch, input)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestRepositoriesService_Transfer(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := TransferRequest{NewOwner: "a", NewName: Ptr("b"), TeamID: []int64{123}}

	mux.HandleFunc("/repos/o/r/transfer", func(w http.ResponseWriter, r *http.Request) {
		var v TransferRequest
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

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

	want := &Repository{Owner: &User{Login: Ptr("a")}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	var input DispatchRequestOptions

	mux.HandleFunc("/repos/o/r/dispatches", func(w http.ResponseWriter, r *http.Request) {
		var v DispatchRequestOptions
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

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

		want := &Repository{Owner: &User{Login: Ptr("a")}}
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
	t.Parallel()
	testJSONMarshal(t, &AdvancedSecurity{}, "{}")

	u := &AdvancedSecurity{
		Status: Ptr("status"),
	}

	want := `{
		"status": "status"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuthorizedActorsOnly_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AuthorizedActorsOnly{}, "{}")

	u := &AuthorizedActorsOnly{
		From: Ptr(true),
	}

	want := `{
		"from" : true
	}`

	testJSONMarshal(t, u, want)
}

func TestDispatchRequestOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DispatchRequestOptions{}, "{}")

	cp := json.RawMessage(`{"testKey":"testValue"}`)
	u := &DispatchRequestOptions{
		EventType:     "test_event_type",
		ClientPayload: &cp,
	}

	want := `{
		"event_type": "test_event_type",
		"client_payload": {
		  "testKey": "testValue"
		}
	  }`

	testJSONMarshal(t, u, want)
}

func TestTransferRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TransferRequest{}, "{}")

	u := &TransferRequest{
		NewOwner: "testOwner",
		NewName:  Ptr("testName"),
		TeamID:   []int64{1, 2},
	}

	want := `{
		"new_owner": "testOwner",
		"new_name": "testName",
		"team_ids": [1,2]
	}`

	testJSONMarshal(t, u, want)
}

func TestSignaturesProtectedBranch_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SignaturesProtectedBranch{}, "{}")

	u := &SignaturesProtectedBranch{
		URL:     Ptr("https://www.testURL.in"),
		Enabled: Ptr(false),
	}

	want := `{
		"url": "https://www.testURL.in",
		"enabled": false
	}`

	testJSONMarshal(t, u, want)

	u2 := &SignaturesProtectedBranch{
		URL:     Ptr("testURL"),
		Enabled: Ptr(true),
	}

	want2 := `{
		"url": "testURL",
		"enabled": true
	}`

	testJSONMarshal(t, u2, want2)
}

func TestDismissalRestrictionsRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DismissalRestrictionsRequest{}, "{}")

	u := &DismissalRestrictionsRequest{
		Users: &[]string{"user1", "user2"},
		Teams: &[]string{"team1", "team2"},
		Apps:  &[]string{"app1", "app2"},
	}

	want := `{
		"users": ["user1","user2"],
		"teams": ["team1","team2"],
		"apps": ["app1","app2"]
	}`

	testJSONMarshal(t, u, want)
}

func TestAdminEnforcement_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AdminEnforcement{}, "{}")

	u := &AdminEnforcement{
		URL:     Ptr("https://www.test-url.in"),
		Enabled: false,
	}

	want := `{
		"url": "https://www.test-url.in",
		"enabled": false
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestReviewsEnforcementUpdate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestReviewsEnforcementUpdate{}, "{}")

	u := &PullRequestReviewsEnforcementUpdate{
		BypassPullRequestAllowancesRequest: &BypassPullRequestAllowancesRequest{
			Users: []string{"user1", "user2"},
			Teams: []string{"team1", "team2"},
			Apps:  []string{"app1", "app2"},
		},
		DismissStaleReviews:          Ptr(false),
		RequireCodeOwnerReviews:      Ptr(true),
		RequiredApprovingReviewCount: 2,
	}

	want := `{
		"bypass_pull_request_allowances": {
			"users": ["user1","user2"],
			"teams": ["team1","team2"],
			"apps": ["app1","app2"]
		},
		"dismiss_stale_reviews": false,
		"require_code_owner_reviews": true,
		"required_approving_review_count": 2
	}`

	testJSONMarshal(t, u, want)
}

func TestRequiredStatusCheck_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RequiredStatusCheck{}, "{}")

	u := &RequiredStatusCheck{
		Context: "ctx",
		AppID:   Ptr(int64(1)),
	}

	want := `{
		"context": "ctx",
		"app_id": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoryTag_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryTag{}, "{}")

	u := &RepositoryTag{
		Name: Ptr("v0.1"),
		Commit: &Commit{
			SHA: Ptr("sha"),
			URL: Ptr("url"),
		},
		ZipballURL: Ptr("zball"),
		TarballURL: Ptr("tball"),
	}

	want := `{
		"name": "v0.1",
		"commit": {
			"sha": "sha",
			"url": "url"
		},
		"zipball_url": "zball",
		"tarball_url": "tball"
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoriesService_EnablePrivateReporting(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/private-vulnerability-reporting", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.EnablePrivateReporting(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("Repositories.EnablePrivateReporting returned error: %v", err)
	}

	const methodName = "EnablePrivateReporting"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.EnablePrivateReporting(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.EnablePrivateReporting(ctx, "owner", "repo")
	})
}

func TestRepositoriesService_DisablePrivateReporting(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/private-vulnerability-reporting", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.DisablePrivateReporting(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("Repositories.DisablePrivateReporting returned error: %v", err)
	}

	const methodName = "DisablePrivateReporting"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DisablePrivateReporting(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DisablePrivateReporting(ctx, "owner", "repo")
	})
}

func TestRepositoriesService_IsPrivateReportingEnabled(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/private-vulnerability-reporting", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled": true}`)
	})

	ctx := context.Background()
	enabled, _, err := client.Repositories.IsPrivateReportingEnabled(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("Repositories.IsPrivateReportingEnabled returned error: %v", err)
	}
	if want := true; enabled != want {
		t.Errorf("Repositories.IsPrivateReportingEnabled returned %+v, want %+v", enabled, want)
	}

	const methodName = "IsPrivateReportingEnabled"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.IsPrivateReportingEnabled(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.IsPrivateReportingEnabled(ctx, "owner", "repo")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestRepository_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var testCases = map[string]struct {
		data           []byte
		wantRepository Repository
		wantErr        bool
	}{
		"Empty": {
			data:           []byte("{}"),
			wantRepository: Repository{},
			wantErr:        false,
		},
		"Invalid JSON": {
			data:           []byte("{"),
			wantRepository: Repository{},
			wantErr:        true,
		},
		"Partial project": {
			data:           []byte(`{"id":10270722,"name":"go-github","private":false,"owner":{"login":"google"},"created_at":"2013-05-24T16:42:58Z","license":{},"topics":["github"],"permissions":{"pull":true},"custom_properties":{},"organization":{"login":"google"}}`),
			wantRepository: Repository{ID: Ptr(int64(10270722)), Name: Ptr("go-github"), Private: Ptr(false), Owner: &User{Login: Ptr("google")}, CreatedAt: &Timestamp{time.Date(2013, 5, 24, 16, 42, 58, 0, time.UTC)}, License: &License{}, Topics: []string{"github"}, Permissions: map[string]bool{"pull": true}, CustomProperties: map[string]interface{}{}, Organization: &Organization{Login: Ptr("google")}},
			wantErr:        false,
		},
		"With custom properties": {
			data:           []byte(`{"custom_properties":{"boolean":"false","text":"a","single-select":"a","multi-select":["a","b","c"]}}`),
			wantRepository: Repository{CustomProperties: map[string]interface{}{"boolean": "false", "text": "a", "single-select": "a", "multi-select": []interface{}{"a", "b", "c"}}},
			wantErr:        false,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			pk := Repository{}
			err := json.Unmarshal(tt.data, &pk)
			if err == nil && tt.wantErr {
				t.Errorf("Repository.UnmarshalJSON returned nil instead of an error")
			}
			if err != nil && !tt.wantErr {
				t.Errorf("Repository.UnmarshalJSON returned an unexpected error: %+v", err)
			}
			if !cmp.Equal(tt.wantRepository, pk) {
				t.Errorf("Repository.UnmarshalJSON expected repository %+v, got %+v", tt.wantRepository, pk)
			}
		})
	}
}
