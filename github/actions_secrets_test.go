// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPublicKey_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		data          []byte
		wantPublicKey PublicKey
		wantErr       bool
	}{
		"Empty": {
			data:          []byte("{}"),
			wantPublicKey: PublicKey{},
			wantErr:       false,
		},
		"Invalid JSON": {
			data:          []byte("{"),
			wantPublicKey: PublicKey{},
			wantErr:       true,
		},
		"Numeric KeyID": {
			data:          []byte(`{"key_id":1234,"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`),
			wantPublicKey: PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")},
			wantErr:       false,
		},
		"String KeyID": {
			data:          []byte(`{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`),
			wantPublicKey: PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")},
			wantErr:       false,
		},
		"Invalid KeyID": {
			data:          []byte(`{"key_id":["1234"],"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`),
			wantPublicKey: PublicKey{KeyID: nil, Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")},
			wantErr:       true,
		},
		"Invalid Key": {
			data:          []byte(`{"key":123}`),
			wantPublicKey: PublicKey{KeyID: nil, Key: nil},
			wantErr:       true,
		},
		"Nil": {
			data:          nil,
			wantPublicKey: PublicKey{KeyID: nil, Key: nil},
			wantErr:       true,
		},
		"Empty String": {
			data:          []byte(""),
			wantPublicKey: PublicKey{KeyID: nil, Key: nil},
			wantErr:       true,
		},
		"Missing Key": {
			data:          []byte(`{"key_id":"1234"}`),
			wantPublicKey: PublicKey{KeyID: Ptr("1234")},
			wantErr:       false,
		},
		"Missing KeyID": {
			data:          []byte(`{"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`),
			wantPublicKey: PublicKey{Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")},
			wantErr:       false,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			pk := PublicKey{}
			err := json.Unmarshal(tt.data, &pk)
			if err == nil && tt.wantErr {
				t.Error("PublicKey.UnmarshalJSON returned nil instead of an error")
			}
			if err != nil && !tt.wantErr {
				t.Errorf("PublicKey.UnmarshalJSON returned an unexpected error: %+v", err)
			}
			if !cmp.Equal(tt.wantPublicKey, pk) {
				t.Errorf("PublicKey.UnmarshalJSON expected public key %+v, got %+v", tt.wantPublicKey, pk)
			}
		})
	}
}

func TestActionsService_GetRepoPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Actions.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
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

func TestActionsService_GetRepoPublicKeyNumeric(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":1234,"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Actions.GetRepoPublicKey(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepoPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Actions.ListRepoSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepoSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(secrets, want) {
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

func TestActionsService_ListRepoOrgSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/organization-secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Actions.ListRepoOrgSecrets(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepoOrgSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Actions.ListRepoOrgSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListRepoOrgSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepoOrgSecrets(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepoOrgSecrets(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	secret, _, err := client.Actions.GetRepoSecret(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.GetRepoSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
	}
	if !cmp.Equal(secret, want) {
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	_, err := client.Actions.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateRepoSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, "o", "r", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrUpdateRepoSecret(ctx, "o", "r", input)
	})
}

func TestActionsService_DeleteRepoSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Actions.GetOrgPublicKey(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetOrgPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("012345678"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
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
	if !cmp.Equal(secrets, want) {
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}`)
	})

	ctx := t.Context()
	secret, _, err := client.Actions.GetOrgSecret(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.GetOrgSecret returned error: %v", err)
	}

	want := &Secret{
		Name:                    "NAME",
		CreatedAt:               Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt:               Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
		Visibility:              "selected",
		SelectedRepositoriesURL: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories",
	}
	if !cmp.Equal(secret, want) {
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	_, err := client.Actions.CreateOrUpdateOrgSecret(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateOrgSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateOrgSecret(ctx, "o", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateOrgSecret(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrUpdateOrgSecret(ctx, "o", input)
	})
}

func TestActionsService_ListSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	repos, _, err := client.Actions.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
	if err != nil {
		t.Errorf("Actions.ListSelectedReposForOrgSecret returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: Ptr(1),
		Repositories: []*Repository{
			{ID: Ptr(int64(1))},
		},
	}
	if !cmp.Equal(repos, want) {
		t.Errorf("Actions.ListSelectedReposForOrgSecret returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListSelectedReposForOrgSecret(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetSelectedReposForOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
	})

	ctx := t.Context()
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := t.Context()
	_, err := client.Actions.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.AddSelectedRepoToOrgSecret returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_RemoveSelectedRepoFromOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := t.Context()
	_, err := client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.RemoveSelectedRepoFromOrgSecret returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_DeleteOrgSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
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

func TestActionsService_GetEnvPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Actions.GetEnvPublicKey(ctx, 1, "e")
	if err != nil {
		t.Errorf("Actions.GetEnvPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Actions.GetEnvPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetEnvPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetEnvPublicKey(ctx, 0.0, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetEnvPublicKey(ctx, 1, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetEnvPublicKeyNumeric(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":1234,"key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := t.Context()
	key, _, err := client.Actions.GetEnvPublicKey(ctx, 1, "e")
	if err != nil {
		t.Errorf("Actions.GetEnvPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !cmp.Equal(key, want) {
		t.Errorf("Actions.GetEnvPublicKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetEnvPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetEnvPublicKey(ctx, 0.0, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetEnvPublicKey(ctx, 1, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListEnvSecrets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	secrets, _, err := client.Actions.ListEnvSecrets(ctx, 1, "e", opts)
	if err != nil {
		t.Errorf("Actions.ListEnvSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(secrets, want) {
		t.Errorf("Actions.ListEnvSecrets returned %+v, want %+v", secrets, want)
	}

	const methodName = "ListEnvSecrets"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListEnvSecrets(ctx, 0.0, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListEnvSecrets(ctx, 1, "e", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetEnvSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets/secret", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"secret","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	secret, _, err := client.Actions.GetEnvSecret(ctx, 1, "e", "secret")
	if err != nil {
		t.Errorf("Actions.GetEnvSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "secret",
		CreatedAt: Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
	}
	if !cmp.Equal(secret, want) {
		t.Errorf("Actions.GetEnvSecret returned %+v, want %+v", secret, want)
	}

	const methodName = "GetEnvSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetEnvSecret(ctx, 0.0, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetEnvSecret(ctx, 1, "e", "secret")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateOrUpdateEnvSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets/secret", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &EncryptedSecret{
		Name:           "secret",
		EncryptedValue: "QIv=",
		KeyID:          "1234",
	}
	ctx := t.Context()
	_, err := client.Actions.CreateOrUpdateEnvSecret(ctx, 1, "e", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateEnvSecret returned error: %v", err)
	}

	const methodName = "CreateOrUpdateEnvSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateEnvSecret(ctx, 1, "e", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrUpdateEnvSecret(ctx, 0.0, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrUpdateEnvSecret(ctx, 1, "e", input)
	})
}

func TestActionsService_DeleteEnvSecret(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/environments/e/secrets/secret", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Actions.DeleteEnvSecret(ctx, 1, "e", "secret")
	if err != nil {
		t.Errorf("Actions.DeleteEnvSecret returned error: %v", err)
	}

	const methodName = "DeleteEnvSecret"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteEnvSecret(ctx, 0.0, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteEnvSecret(ctx, 1, "r", "secret")
	})
}

func TestPublicKey_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PublicKey{}, `{
		"key": null,
		"key_id": null
	}`)

	u := &PublicKey{
		KeyID: Ptr("kid"),
		Key:   Ptr("k"),
	}

	want := `{
		"key_id": "kid",
		"key": "k"
	}`

	testJSONMarshal(t, u, want)
}

func TestSecret_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Secret{}, `{
		"name": "",
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`)

	u := &Secret{
		Name:                    "n",
		CreatedAt:               Timestamp{referenceTime},
		UpdatedAt:               Timestamp{referenceTime},
		Visibility:              "v",
		SelectedRepositoriesURL: "s",
	}

	want := `{
		"name": "n",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"visibility": "v",
		"selected_repositories_url": "s"
	}`

	testJSONMarshal(t, u, want)
}

func TestSecrets_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Secrets{}, `{
		"total_count": 0,
		"secrets": null
	}`)

	u := &Secrets{
		TotalCount: 1,
		Secrets: []*Secret{
			{
				Name:                    "n",
				CreatedAt:               Timestamp{referenceTime},
				UpdatedAt:               Timestamp{referenceTime},
				Visibility:              "v",
				SelectedRepositoriesURL: "s",
			},
		},
	}

	want := `{
		"total_count": 1,
		"secrets": [
			{
				"name": "n",
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"visibility": "v",
				"selected_repositories_url": "s"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestEncryptedSecret_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EncryptedSecret{}, `{
		"key_id": "",
		"encrypted_value": ""
	}`)

	u := &EncryptedSecret{
		Name:                  "n",
		KeyID:                 "kid",
		EncryptedValue:        "e",
		Visibility:            "v",
		SelectedRepositoryIDs: []int64{1},
	}

	want := `{
		"key_id": "kid",
		"encrypted_value": "e",
		"visibility": "v",
		"selected_repository_ids": [1]
	}`

	testJSONMarshal(t, u, want, cmpIgnoreFieldOption("Name"))
}

func TestSelectedReposList_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SelectedReposList{}, "{}")

	u := &SelectedReposList{
		TotalCount: Ptr(1),
		Repositories: []*Repository{
			{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
		},
	}

	want := `{
		"total_count": 1,
		"repositories": [
			{
				"id": 1,
				"url": "u",
				"name": "n"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
