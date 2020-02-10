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

func TestActionsService_GetPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	key, _, err := client.Actions.GetPublicKey(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Actions.GetPublicKey returned error: %v", err)
	}

	want := &PublicKey{KeyID: String("1234"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Actions.GetPublicKey returned %+v, want %+v", key, want)
	}
}

func TestActionsService_ListSecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	secrets, _, err := client.Actions.ListSecrets(context.Background(), "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListSecrets returned error: %v", err)
	}

	want := &Secrets{
		TotalCount: 4,
		Secrets: []*Secret{
			{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(secrets, want) {
		t.Errorf("Actions.ListSecrets returned %+v, want %+v", secrets, want)
	}
}

func TestActionsService_GetSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	secret, _, err := client.Actions.GetSecret(context.Background(), "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.GetSecret returned error: %v", err)
	}

	want := &Secret{
		Name:      "NAME",
		CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Actions.GetSecret returned %+v, want %+v", secret, want)
	}
}

func TestActionsService_CreateOrUpdateSecret(t *testing.T) {
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
	_, err := client.Actions.CreateOrUpdateSecret(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateOrUpdateSecret returned error: %v", err)
	}
}

func TestActionsService_DeleteSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Actions.DeleteSecret(context.Background(), "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteSecret returned error: %v", err)
	}
}
