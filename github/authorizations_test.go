// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAuthorizationsService_Check(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Authorizations.Check(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Check returned error: %v", err)
	}

	want := &Authorization{ID: Ptr(int64(1))}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.Check returned auth %+v, want %+v", got, want)
	}

	const methodName = "Check"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.Check(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.Check(ctx, "id", "a")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_Reset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		fmt.Fprint(w, `{"ID":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Authorizations.Reset(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Reset returned error: %v", err)
	}

	want := &Authorization{ID: Ptr(int64(1))}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.Reset returned auth %+v, want %+v", got, want)
	}

	const methodName = "Reset"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.Reset(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.Reset(ctx, "id", "a")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_Revoke(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Authorizations.Revoke(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Revoke returned error: %v", err)
	}

	const methodName = "Revoke"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.Revoke(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.Revoke(ctx, "id", "a")
	})
}

func TestDeleteGrant(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/applications/id/grant", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
	})

	ctx := context.Background()
	_, err := client.Authorizations.DeleteGrant(ctx, "id", "a")
	if err != nil {
		t.Errorf("OAuthAuthorizations.DeleteGrant returned error: %v", err)
	}

	const methodName = "DeleteGrant"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.DeleteGrant(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.DeleteGrant(ctx, "id", "a")
	})
}

func TestAuthorizationsService_CreateImpersonation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1}`)
	})

	req := &AuthorizationRequest{Scopes: []Scope{ScopePublicRepo}}
	ctx := context.Background()
	got, _, err := client.Authorizations.CreateImpersonation(ctx, "u", req)
	if err != nil {
		t.Errorf("Authorizations.CreateImpersonation returned error: %+v", err)
	}

	want := &Authorization{ID: Ptr(int64(1))}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.CreateImpersonation returned %+v, want %+v", *got.ID, *want.ID)
	}

	const methodName = "CreateImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.CreateImpersonation(ctx, "\n", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.CreateImpersonation(ctx, "u", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_DeleteImpersonation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Authorizations.DeleteImpersonation(ctx, "u")
	if err != nil {
		t.Errorf("Authorizations.DeleteImpersonation returned error: %+v", err)
	}

	const methodName = "DeleteImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.DeleteImpersonation(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.DeleteImpersonation(ctx, "u")
	})
}

func TestAuthorizationUpdateRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AuthorizationUpdateRequest{}, "{}")

	u := &AuthorizationUpdateRequest{
		Scopes:       []string{"s"},
		AddScopes:    []string{"a"},
		RemoveScopes: []string{"r"},
		Note:         Ptr("n"),
		NoteURL:      Ptr("nu"),
		Fingerprint:  Ptr("f"),
	}

	want := `{
		"scopes": ["s"],
		"add_scopes": ["a"],
		"remove_scopes": ["r"],
		"note": "n",
		"note_url": "nu",
		"fingerprint": "f"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuthorizationRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AuthorizationRequest{}, "{}")

	u := &AuthorizationRequest{
		Scopes:       []Scope{"s"},
		ClientID:     Ptr("cid"),
		ClientSecret: Ptr("cs"),
		Note:         Ptr("n"),
		NoteURL:      Ptr("nu"),
		Fingerprint:  Ptr("f"),
	}

	want := `{
		"scopes": ["s"],
		"client_id": "cid",
		"client_secret": "cs",
		"note": "n",
		"note_url": "nu",
		"fingerprint": "f"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuthorizationApp_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AuthorizationApp{}, "{}")

	u := &AuthorizationApp{
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

func TestGrant_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Grant{}, "{}")

	u := &Grant{
		ID:  Ptr(int64(1)),
		URL: Ptr("u"),
		App: &AuthorizationApp{
			URL:      Ptr("u"),
			Name:     Ptr("n"),
			ClientID: Ptr("cid"),
		},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		Scopes:    []string{"s"},
	}

	want := `{
		"id": 1,
		"url": "u",
		"app": {
			"url": "u",
			"name": "n",
			"client_id": "cid"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"scopes": ["s"]
	}`

	testJSONMarshal(t, u, want)
}

func TestAuthorization_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Authorization{}, "{}")

	u := &Authorization{
		ID:             Ptr(int64(1)),
		URL:            Ptr("u"),
		Scopes:         []Scope{"s"},
		Token:          Ptr("t"),
		TokenLastEight: Ptr("tle"),
		HashedToken:    Ptr("ht"),
		App: &AuthorizationApp{
			URL:      Ptr("u"),
			Name:     Ptr("n"),
			ClientID: Ptr("cid"),
		},
		Note:        Ptr("n"),
		NoteURL:     Ptr("nu"),
		UpdatedAt:   &Timestamp{referenceTime},
		CreatedAt:   &Timestamp{referenceTime},
		Fingerprint: Ptr("f"),
		User: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
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
		"fingerprint": "f",
		"user": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		}
	}`

	testJSONMarshal(t, u, want)
}
