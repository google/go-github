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

func TestActionsService_GetRepositoryPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	key, _, err := client.Actions.GetRepositoryPublicKey(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepositoryPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: String("1234"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Actions.GetRepositoryPublicKey returned %+v, want %+v", key, want)
	}
}

func TestActionsService_ListRepositorySecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	secrets, _, err := client.Actions.ListRepositorySecrets(context.Background(), "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepositorySecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(secrets, want) {
		t.Errorf("Actions.ListRepositorySecrets returned %+v, want %+v", secrets, want)
	}
}

func TestActionsService_GetRepositorySecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	secret, _, err := client.Actions.GetRepositorySecret(context.Background(), "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.GetRepositorySecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Actions.GetRepositorySecret returned %+v, want %+v", secret, want)
	}
}

func TestActionsService_CreateOrUpdateRepositorySecret(t *testing.T) {
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
	_, err := client.Actions.CreateOrUpdateRepositorySecret(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateRepositorySecret returned error: %v", err)
	}
}

func TestActionsService_DeleteRepositorySecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Actions.DeleteRepositorySecret(context.Background(), "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteRepositorySecret returned error: %v", err)
	}
}

func TestActionsService_GetOrganizationPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	key, _, err := client.Actions.GetOrganizationPublicKey(context.Background(), "o")
	if err != nil {
		t.Errorf("Actions.GetOrganizationPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: String("012345678"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Actions.GetOrganizationPublicKey returned %+v, want %+v", key, want)
	}
}

func TestActionsService_ListOrganizationSecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	secrets, _, err := client.Actions.ListOrganizationSecrets(context.Background(), "o", opts)
	if err != nil {
		t.Errorf("Actions.ListOrganizationSecrets returned error: %v", err)
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
		t.Errorf("Actions.ListOrganizationSecrets returned %+v, want %+v", secrets, want)
	}
}

func TestActionsService_GetOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}`)
	})

	secret, _, err := client.Actions.GetOrganizationSecret(context.Background(), "o", "NAME")
	if err != nil {
		t.Errorf("Actions.GetRepositorySecret returned error: %v", err)
	}

	want := &Secret{
		Name:                    "NAME",
		CreatedAt:               Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt:               Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
		Visibility:              "selected",
		SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories",
	}
	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Actions.GetOrganizationSecret returned %+v, want %+v", secret, want)
	}
}

func TestActionsService_CreateOrUpdateOrganizationSecret(t *testing.T) {
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
		SelectedRepositoryIDs: SelectedRepositoryIDs{[]int64{1296269, 1269280}},
	}
	_, err := client.Actions.CreateOrUpdateOrganizationSecret(context.Background(), "o", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateOrganizationSecret returned error: %v", err)
	}
}

func TestActionsService_ListSelectedRepositoriesForOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	repos, _, err := client.Actions.ListSelectedRepositoriesForOrganizationSecret(context.Background(), "o", "NAME")
	if err != nil {
		t.Errorf("Actions.ListSelectedRepositoriesForOrganizationSecret returned error: %v", err)
	}

	want := &SelectedRepositoriesList{
		TotalCount: Int(1),
		Repositories: []*Repository{
			{ID: Int64(1)},
		},
	}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Actions.ListSelectedRepositoriesForOrganizationSecret returned %+v, want %+v", repos, want)
	}
}

func TestActionsService_SetSelectedRepositoriesForOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
	})

	_, err := client.Actions.SetSelectedRepositoriesForOrganizationSecret(context.Background(), "o", "NAME", SelectedRepositoryIDs{[]int64{64780797}})
	if err != nil {
		t.Errorf("Actions.SetSelectedRepositoriesForOrganizationSecret returned error: %v", err)
	}
}

func TestActionsService_AddSelectedRepositoryToOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	repo := &Repository{ID: Int64(1234)}
	_, err := client.Actions.AddSelectedRepositoryToOrganizationSecret(context.Background(), "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.AddSelectedRepositoryToOrganizationSecret returned error: %v", err)
	}
}

func TestActionsService_RemoveSelectedRepositoryFromOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	repo := &Repository{ID: Int64(1234)}
	_, err := client.Actions.RemoveSelectedRepositoryFromOrganizationSecret(context.Background(), "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.RemoveSelectedRepositoryFromOrganizationSecret returned error: %v", err)
	}
}

func TestActionsService_DeleteOrganizationSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Actions.DeleteOrganizationSecret(context.Background(), "o", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteOrganizationSecret returned error: %v", err)
	}
}
