// Copyright 2019 The go-github AUTHORS. All rights reserved.
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
	"testing"
	"time"
)

func TestAdminUsers_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		v := new(createUserRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &createUserRequest{Login: String("github"), Email: String("email@domain.com")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"login":"github","id":1}`)
	})

	org, _, err := client.Admin.CreateUser(context.Background(), "github", "email@domain.com")
	if err != nil {
		t.Errorf("Admin.CreateUser returned error: %v", err)
	}

	want := &User{ID: Int64(1), Login: String("github")}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Admin.CreateUser returned %+v, want %+v", org, want)
	}
}

func TestAdminUsers_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/github", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Admin.DeleteUser(context.Background(), "github")
	if err != nil {
		t.Errorf("Admin.DeleteUser returned error: %v", err)
	}
}

func TestUserImpersonation_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/github/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"scopes":["repo"]}`+"\n")
		fmt.Fprint(w, `{"id": 1234,
		"url": "https://git.company.com/api/v3/authorizations/1234",
		"app": {
		  "name": "GitHub Site Administrator",
		  "url": "https://developer.github.com/v3/enterprise/users/",
		  "client_id": "1234"
		},
		"token": "1234",
		"hashed_token": "1234",
		"token_last_eight": "1234",
		"note": null,
		"note_url": null,
		"created_at": "2018-01-01T00:00:00Z",
		"updated_at": "2018-01-01T00:00:00Z",
		"scopes": [
		  "repo"
		],
		"fingerprint": null}`)
	})

	opt := &ImpersonateUserOptions{Scopes: []string{"repo"}}
	auth, _, err := client.Admin.CreateUserImpersonation(context.Background(), "github", opt)
	if err != nil {
		t.Errorf("Admin.CreateUserImpersonation returned error: %v", err)
	}

	date := Timestamp{Time: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}
	want := &UserAuthorization{
		ID:  Int64(1234),
		URL: String("https://git.company.com/api/v3/authorizations/1234"),
		App: &OAuthAPP{
			Name:     String("GitHub Site Administrator"),
			URL:      String("https://developer.github.com/v3/enterprise/users/"),
			ClientID: String("1234"),
		},
		Token:          String("1234"),
		HashedToken:    String("1234"),
		TokenLastEight: String("1234"),
		Note:           nil,
		NoteURL:        nil,
		CreatedAt:      &date,
		UpdatedAt:      &date,
		Scopes:         []string{"repo"},
		Fingerprint:    nil,
	}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("Admin.CreateUserImpersonation returned %+v, want %+v", auth, want)
	}
}

func TestUserImpersonation_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/github/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Admin.DeleteUserImpersonation(context.Background(), "github")
	if err != nil {
		t.Errorf("Admin.DeleteUserImpersonation returned error: %v", err)
	}
}
