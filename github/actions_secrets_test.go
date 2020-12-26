// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestActionsService_GetRepoPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Actions.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: String("1234"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Actions.GetRepoPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetRepoPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepoPublicKey(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepoPublicKey(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListRepoSecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	secrets, _, err := client.Actions.ListRepoSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepoSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(secrets, want) {
		t.Errorf("Actions.ListRepoSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListRepoSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepoSecrets(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepoSecrets(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepoSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := context.Background()
	secret, _, err := client.Actions.GetRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.GetRepoSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Actions.GetRepoSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepoSecret(ctx, "o", "r", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateOrUpdateRepoSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &EncryptedSecret{
		Name:           "NAME",
		EncryptedValue: "QIv=",
		KeyID:          "1234",
	}
	ctx := context.Background()
	_, err := client.Actions.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	})
}

func TestActionsService_DeleteRepoSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Actions.DeleteRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteRepoSecret returned error: %v", err)
	}

	const methodName = "DeleteRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteRepoSecret(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteRepoSecret(ctx, "o", "r", "NAME")
	})
}

func TestActionsService_GetOrgPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Actions.GetOrgPublicKey(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetOrgPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: String("012345678"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Actions.GetOrgPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetOrgPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrgPublicKey(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrgPublicKey(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListOrgSecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	secrets, _, err := client.Actions.ListOrgSecrets(ctx, "o", opts)
	if err != nil {
		t.Errorf("Actions.ListOrgSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 3,
		Secrets: []*Secret{
			{Name: "GIST_ID", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "private"},
			{Name: "DEPLOY_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "all"},
			{Name: "GH_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "selected", SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"},
		},
	}
	if !reflect.DeepEqual(secrets, want) {
		t.Errorf("Actions.ListOrgSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListOrgSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListOrgSecrets(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListOrgSecrets(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}`)
	})

	ctx := context.Background()
	secret, _, err := client.Actions.GetOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.GetOrgSecret returned error: %v", err)
	}

	want := &Secret{
		Name:                    "NAME",
		CreatedAt:               Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt:               Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
		Visibility:              "selected",
		SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories",
	}
	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Actions.GetOrgSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrgSecret(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateOrUpdateOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv=","visibility":"selected","selected_repository_ids":[1296269,1269280]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &EncryptedSecret{
		Name:                  "NAME",
		EncryptedValue:        "QIv=",
		KeyID:                 "1234",
		Visibility:            "selected",
		SelectedRepositoryIDs: SelectedRepoIDs{1296269, 1269280},
	}
	ctx := context.Background()
	_, err := client.Actions.CreateOrUpdateOrgSecret(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateOrgSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateOrgSecret(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrUpdateOrgSecret(ctx, "o", input)
	})
}

func TestActionsService_ListSelectedReposForOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	ctx := context.Background()
	repos, _, err := client.Actions.ListSelectedReposForOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.ListSelectedReposForOrgSecret returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: Int(1),
		Repositories: []*Repository{
			{ID: Int64(1)},
		},
	}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Actions.ListSelectedReposForOrgSecret returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListSelectedReposForOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListSelectedReposForOrgSecret(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetSelectedReposForOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
	})

	ctx := context.Background()
	_, err := client.Actions.SetSelectedReposForOrgSecret(ctx, "o", "NAME", SelectedRepoIDs{64780797})
	if err != nil {
		t.Errorf("Actions.SetSelectedReposForOrgSecret returned error: %v", err)
	}

	const methodName = "SetSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetSelectedReposForOrgSecret(ctx, "\n", "\n", SelectedRepoIDs{64780797})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetSelectedReposForOrgSecret(ctx, "o", "NAME", SelectedRepoIDs{64780797})
	})
}

func TestActionsService_AddSelectedRepoToOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	repo := &Repository{ID: Int64(1234)}
	ctx := context.Background()
	_, err := client.Actions.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.AddSelectedRepoToOrgSecret returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_RemoveSelectedRepoFromOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	repo := &Repository{ID: Int64(1234)}
	ctx := context.Background()
	_, err := client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.RemoveSelectedRepoFromOrgSecret returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_DeleteOrgSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Actions.DeleteOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteOrgSecret returned error: %v", err)
	}

	const methodName = "DeleteOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteOrgSecret(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteOrgSecret(ctx, "o", "NAME")
	})
}
