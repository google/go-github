// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDependabotService_GetRepoPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Dependabot.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Dependabot.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Dependabot.GetRepoPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetRepoPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetRepoPublicKey(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetRepoPublicKey(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_GetRepoPublicKeyNumeric(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":1234,"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Dependabot.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Dependabot.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Dependabot.GetRepoPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetRepoPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetRepoPublicKey(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetRepoPublicKey(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_ListRepoSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	secrets, _, err := client.Dependabot.ListRepoSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Dependabot.ListRepoSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Dependabot.ListRepoSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListRepoSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListRepoSecrets(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListRepoSecrets(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_GetRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := context.Background()
	secret, _, err := client.Dependabot.GetRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Dependabot.GetRepoSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !cmp.Equal(secret, want) {
		t.Errorf("Dependabot.GetRepoSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetRepoSecret(ctx, "o", "r", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_CreateOrUpdateRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &DependabotEncryptedSecret{
		Name:           "NAME",
		EncryptedValue: "QIv=",
		KeyID:          "1234",
	}
	ctx := context.Background()
	_, err := client.Dependabot.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Dependabot.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.CreateOrUpdateRepoSecret(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	})
}

func TestDependabotService_DeleteRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Dependabot.DeleteRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Dependabot.DeleteRepoSecret returned error: %v", err)
	}

	const methodName = "DeleteRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.DeleteRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.DeleteRepoSecret(ctx, "o", "r", "NAME")
	})
}

func TestDependabotService_GetOrgPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Dependabot.GetOrgPublicKey(ctx, "o")
	if err != nil {
		t.Errorf("Dependabot.GetOrgPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("012345678"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Dependabot.GetOrgPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetOrgPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetOrgPublicKey(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetOrgPublicKey(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_ListOrgSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/dependabot/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	secrets, _, err := client.Dependabot.ListOrgSecrets(ctx, "o", opts)
	if err != nil {
		t.Errorf("Dependabot.ListOrgSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 3,
		Secrets: []*Secret{
			{Name: "GIST_ID", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "private"},
			{Name: "DEPLOY_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "all"},
			{Name: "GH_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "selected", SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/dependabot/secrets/SUPER_SECRET/repositories"},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Dependabot.ListOrgSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListOrgSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListOrgSecrets(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListOrgSecrets(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_GetOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/dependabot/secrets/SUPER_SECRET/repositories"}`)
	})

	ctx := context.Background()
	secret, _, err := client.Dependabot.GetOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Dependabot.GetOrgSecret returned error: %v", err)
	}

	want := &Secret{
		Name:                    "NAME",
		CreatedAt:               Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt:               Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
		Visibility:              "selected",
		SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/dependabot/secrets/SUPER_SECRET/repositories",
	}
	if !cmp.Equal(secret, want) {
		t.Errorf("Dependabot.GetOrgSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetOrgSecret(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_CreateOrUpdateOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv=","visibility":"selected","selected_repository_ids":["1296269","1269280"]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &DependabotEncryptedSecret{
		Name:                  "NAME",
		EncryptedValue:        "QIv=",
		KeyID:                 "1234",
		Visibility:            "selected",
		SelectedRepositoryIDs: DependabotSecretsSelectedRepoIDs{1296269, 1269280},
	}
	ctx := context.Background()
	_, err := client.Dependabot.CreateOrUpdateOrgSecret(ctx, "o", input)
	if err != nil {
		t.Errorf("Dependabot.CreateOrUpdateOrgSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.CreateOrUpdateOrgSecret(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.CreateOrUpdateOrgSecret(ctx, "o", input)
	})
}

func TestDependabotService_ListSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	repos, _, err := client.Dependabot.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
	if err != nil {
		t.Errorf("Dependabot.ListSelectedReposForOrgSecret returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: Ptr(1),
		Repositories: []*Repository{
			{ID: Ptr(int64(1))},
		},
	}
	if !cmp.Equal(repos, want) {
		t.Errorf("Dependabot.ListSelectedReposForOrgSecret returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListSelectedReposForOrgSecret(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_SetSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
	})

	ctx := context.Background()
	_, err := client.Dependabot.SetSelectedReposForOrgSecret(ctx, "o", "NAME", DependabotSecretsSelectedRepoIDs{64780797})
	if err != nil {
		t.Errorf("Dependabot.SetSelectedReposForOrgSecret returned error: %v", err)
	}

	const methodName = "SetSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.SetSelectedReposForOrgSecret(ctx, "\n", "\n", DependabotSecretsSelectedRepoIDs{64780797})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.SetSelectedReposForOrgSecret(ctx, "o", "NAME", DependabotSecretsSelectedRepoIDs{64780797})
	})
}

func TestDependabotService_AddSelectedRepoToOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := context.Background()
	_, err := client.Dependabot.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Dependabot.AddSelectedRepoToOrgSecret returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.AddSelectedRepoToOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestDependabotService_RemoveSelectedRepoFromOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := context.Background()
	_, err := client.Dependabot.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Dependabot.RemoveSelectedRepoFromOrgSecret returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.RemoveSelectedRepoFromOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestDependabotService_DeleteOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Dependabot.DeleteOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Dependabot.DeleteOrgSecret returned error: %v", err)
	}

	const methodName = "DeleteOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Dependabot.DeleteOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Dependabot.DeleteOrgSecret(ctx, "o", "NAME")
	})
}
