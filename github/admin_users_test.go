// Copyright 2019 The go-github AUTHORS. All rights reserved.
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

func TestAdminUsers_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreateUserRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		want := &CreateUserRequest{Login: "github", Email: Ptr("email@example.com"), Suspended: Ptr(false)}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"login":"github","id":1}`)
	})

	ctx := t.Context()
	org, _, err := client.Admin.CreateUser(ctx, CreateUserRequest{
		Login:     "github",
		Email:     Ptr("email@example.com"),
		Suspended: Ptr(false),
	})
	if err != nil {
		t.Errorf("Admin.CreateUser returned error: %v", err)
	}

	want := &User{ID: Ptr(int64(1)), Login: Ptr("github")}
	if !cmp.Equal(org, want) {
		t.Errorf("Admin.CreateUser returned %+v, want %+v", org, want)
	}

	const methodName = "CreateUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Admin.CreateUser(ctx, CreateUserRequest{
			Login:     "github",
			Email:     Ptr("email@example.com"),
			Suspended: Ptr(false),
		})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAdminUsers_Delete(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users/github", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Admin.DeleteUser(ctx, "github")
	if err != nil {
		t.Errorf("Admin.DeleteUser returned error: %v", err)
	}

	const methodName = "DeleteUser"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Admin.DeleteUser(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Admin.DeleteUser(ctx, "github")
	})
}

func TestUserImpersonation_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users/github/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"scopes":["repo"]}`+"\n")
		fmt.Fprint(w, `{"id": 1234,
		"url": "https://example.com/authorizations",
		"app": {
		  "name": "GitHub Site Administrator",
		  "url": "https://docs.github.com/en/rest/enterprise/users/",
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
	ctx := t.Context()
	auth, _, err := client.Admin.CreateUserImpersonation(ctx, "github", opt)
	if err != nil {
		t.Errorf("Admin.CreateUserImpersonation returned error: %v", err)
	}

	date := Timestamp{Time: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}
	want := &UserAuthorization{
		ID:  Ptr(int64(1234)),
		URL: Ptr("https://example.com/authorizations"),
		App: &OAuthAPP{
			Name:     Ptr("GitHub Site Administrator"),
			URL:      Ptr("https://docs.github.com/en/rest/enterprise/users/"),
			ClientID: Ptr("1234"),
		},
		Token:          Ptr("1234"),
		HashedToken:    Ptr("1234"),
		TokenLastEight: Ptr("1234"),
		Note:           nil,
		NoteURL:        nil,
		CreatedAt:      &date,
		UpdatedAt:      &date,
		Scopes:         []string{"repo"},
		Fingerprint:    nil,
	}
	if !cmp.Equal(auth, want) {
		t.Errorf("Admin.CreateUserImpersonation returned %+v, want %+v", auth, want)
	}

	const methodName = "CreateUserImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Admin.CreateUserImpersonation(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Admin.CreateUserImpersonation(ctx, "github", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUserImpersonation_Delete(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users/github/authorizations", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Admin.DeleteUserImpersonation(ctx, "github")
	if err != nil {
		t.Errorf("Admin.DeleteUserImpersonation returned error: %v", err)
	}

	const methodName = "DeleteUserImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Admin.DeleteUserImpersonation(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Admin.DeleteUserImpersonation(ctx, "github")
	})
}

func TestCreateUserRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CreateUserRequest{}, `{"login": ""}`)

	u := &CreateUserRequest{
		Login: "l",
		Email: Ptr("e"),
	}

	want := `{
		"login": "l",
		"email": "e"
	}`

	testJSONMarshal(t, u, want)
}

func TestImpersonateUserOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ImpersonateUserOptions{}, "{}")

	u := &ImpersonateUserOptions{
		Scopes: []string{
			"s",
		},
	}

	want := `{
		"scopes": ["s"]
	}`

	testJSONMarshal(t, u, want)
}

func TestOAuthAPP_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &OAuthAPP{}, "{}")

	u := &OAuthAPP{
		URL:      Ptr("u"),
		Name:     Ptr("n"),
		ClientID: Ptr("cid"),
	}

	want := `{
		"url": "u",
		"name": "n",
		"client_id": "cid"
	}`

	testJSONMarshal(t, u, want)
}

func TestUserAuthorization_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UserAuthorization{}, "{}")

	u := &UserAuthorization{
		ID:  Ptr(int64(1)),
		URL: Ptr("u"),
		Scopes: []string{
			"s",
		},
		Token:          Ptr("t"),
		TokenLastEight: Ptr("tle"),
		HashedToken:    Ptr("ht"),
		App: &OAuthAPP{
			URL:      Ptr("u"),
			Name:     Ptr("n"),
			ClientID: Ptr("cid"),
		},
		Note:        Ptr("n"),
		NoteURL:     Ptr("nu"),
		UpdatedAt:   &Timestamp{referenceTime},
		CreatedAt:   &Timestamp{referenceTime},
		Fingerprint: Ptr("f"),
	}

	want := `{
		"id": 1,
		"url": "u",
		"scopes": ["s"],
		"token": "t",
		"token_last_eight": "tle",
		"hashed_token": "ht",
		"app": {
			"url": "u",
			"name": "n",
			"client_id": "cid"
		},
		"note": "n",
		"note_url": "nu",
		"updated_at": ` + referenceTimeStr + `,
		"created_at": ` + referenceTimeStr + `,
		"fingerprint": "f"
	}`

	testJSONMarshal(t, u, want)
}
