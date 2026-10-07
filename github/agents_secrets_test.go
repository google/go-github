// Copyright 2026 The go-github AUTHORS. All rights reserved.
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

func TestAgentsService_GetRepoPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234","id":1,"url":"https://api.github.com/repos/o/r/agents/secrets/public-key","title":"title","created_at":`+referenceTimeStr+`}`)
	})

	ctx := t.Context()
	key, _, err := client.Agents.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Agents.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{
		KeyID:     new("1234"),
		Key:       new("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"),
		ID:        new(int64(1)),
		URL:       new("https://api.github.com/repos/o/r/agents/secrets/public-key"),
		Title:     new("title"),
		CreatedAt: &referenceTimestamp,
	}
	if !cmp.Equal(key, want) {
		t.Errorf("Agents.GetRepoPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetRepoPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetRepoPublicKey(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetRepoPublicKey(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_GetRepoPublicKeyNumeric(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"key_id":1234,"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Agents.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Agents.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: new("1234"), Key: new("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Agents.GetRepoPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetRepoPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetRepoPublicKey(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetRepoPublicKey(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_ListRepoSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`},{"name":"B","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Agents.ListRepoSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Agents.ListRepoSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp},
			{Name: "B", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Agents.ListRepoSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListRepoSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListRepoSecrets(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListRepoSecrets(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_ListRepoOrgSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/organization-secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`},{"name":"B","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Agents.ListRepoOrgSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Agents.ListRepoOrgSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp},
			{Name: "B", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Agents.ListRepoOrgSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListRepoOrgSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListRepoOrgSecrets(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListRepoOrgSecrets(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_GetRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"name":"NAME","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}`)
	})

	ctx := t.Context()
	secret, _, err := client.Agents.GetRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Agents.GetRepoSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: referenceTimestamp,
		UpdatedAt: referenceTimestamp,
	}
	if !cmp.Equal(secret, want) {
		t.Errorf("Agents.GetRepoSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetRepoSecret(ctx, "o", "r", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_CreateOrUpdateRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretRequest{
		EncryptedValue: "QIv=",
		KeyID:          "1234",
	}

	mux.HandleFunc("/repos/o/r/agents/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
	})

	ctx := t.Context()
	_, err := client.Agents.CreateOrUpdateRepoSecret(ctx, "o", "r", "NAME", input)
	if err != nil {
		t.Errorf("Agents.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.CreateOrUpdateRepoSecret(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.CreateOrUpdateRepoSecret(ctx, "o", "r", "NAME", input)
	})
}

func TestAgentsService_DeleteRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/secrets/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
	})

	ctx := t.Context()
	_, err := client.Agents.DeleteRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Agents.DeleteRepoSecret returned error: %v", err)
	}

	const methodName = "DeleteRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.DeleteRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.DeleteRepoSecret(ctx, "o", "r", "NAME")
	})
}

func TestAgentsService_GetOrgPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"key_id":"012345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Agents.GetOrgPublicKey(ctx, "o")
	if err != nil {
		t.Errorf("Agents.GetOrgPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: new("012345678"), Key: new("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Agents.GetOrgPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetOrgPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetOrgPublicKey(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetOrgPublicKey(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_ListOrgSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"all"},{"name":"GH_TOKEN","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/agents/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Agents.ListOrgSecrets(ctx, "o", opts)
	if err != nil {
		t.Errorf("Agents.ListOrgSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 3,
		Secrets: []*Secret{
			{Name: "GIST_ID", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp, Visibility: "private"},
			{Name: "DEPLOY_TOKEN", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp, Visibility: "all"},
			{Name: "GH_TOKEN", CreatedAt: referenceTimestamp, UpdatedAt: referenceTimestamp, Visibility: "selected", SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/agents/secrets/SUPER_SECRET/repositories"},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Agents.ListOrgSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListOrgSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListOrgSecrets(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListOrgSecrets(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_GetOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"name":"NAME","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/agents/secrets/SUPER_SECRET/repositories"}`)
	})

	ctx := t.Context()
	secret, _, err := client.Agents.GetOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Agents.GetOrgSecret returned error: %v", err)
	}

	want := &Secret{
		Name:                    "NAME",
		CreatedAt:               referenceTimestamp,
		UpdatedAt:               referenceTimestamp,
		Visibility:              "selected",
		SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/agents/secrets/SUPER_SECRET/repositories",
	}
	if !cmp.Equal(secret, want) {
		t.Errorf("Agents.GetOrgSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetOrgSecret(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_CreateOrUpdateOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretOrgRequest{
		EncryptedValue:        "QIv=",
		KeyID:                 "1234",
		Visibility:            "selected",
		SelectedRepositoryIDs: []int64{1296269, 1269280},
	}

	mux.HandleFunc("/orgs/o/agents/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
	})

	ctx := t.Context()
	_, err := client.Agents.CreateOrUpdateOrgSecret(ctx, "o", "NAME", input)
	if err != nil {
		t.Errorf("Agents.CreateOrUpdateOrgSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.CreateOrUpdateOrgSecret(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.CreateOrUpdateOrgSecret(ctx, "o", "NAME", input)
	})
}

func TestAgentsService_ListSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	repos, _, err := client.Agents.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
	if err != nil {
		t.Errorf("Agents.ListSelectedReposForOrgSecret returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: new(1),
		Repositories: []*Repository{
			{ID: new(int64(1))},
		},
	}
	if !cmp.Equal(repos, want) {
		t.Errorf("Agents.ListSelectedReposForOrgSecret returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListSelectedReposForOrgSecret(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_SetSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := []int64{64780797}

	mux.HandleFunc("/orgs/o/agents/secrets/NAME/repositories", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, struct {
			SelectedIDs []int64 `json:"selected_repository_ids"`
		}{
			SelectedIDs: input,
		})
	})

	ctx := t.Context()
	_, err := client.Agents.SetSelectedReposForOrgSecret(ctx, "o", "NAME", input)
	if err != nil {
		t.Errorf("Agents.SetSelectedReposForOrgSecret returned error: %v", err)
	}

	const methodName = "SetSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.SetSelectedReposForOrgSecret(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.SetSelectedReposForOrgSecret(ctx, "o", "NAME", input)
	})
}

func TestAgentsService_AddSelectedRepoToOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
	})

	repoID := int64(1234)
	ctx := t.Context()
	_, err := client.Agents.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repoID)
	if err != nil {
		t.Errorf("Agents.AddSelectedRepoToOrgSecret returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", 0)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.AddSelectedRepoToOrgSecret(ctx, "\n", "\n", repoID)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repoID)
	})
}

func TestAgentsService_RemoveSelectedRepoFromOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
	})

	repoID := int64(1234)
	ctx := t.Context()
	_, err := client.Agents.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repoID)
	if err != nil {
		t.Errorf("Agents.RemoveSelectedRepoFromOrgSecret returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", 0)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.RemoveSelectedRepoFromOrgSecret(ctx, "\n", "\n", repoID)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repoID)
	})
}

func TestAgentsService_DeleteOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/secrets/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
	})

	ctx := t.Context()
	_, err := client.Agents.DeleteOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Agents.DeleteOrgSecret returned error: %v", err)
	}

	const methodName = "DeleteOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.DeleteOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.DeleteOrgSecret(ctx, "o", "NAME")
	})
}
