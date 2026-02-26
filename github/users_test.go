// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUser_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &User{}, "{}")

	u := &User{
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
	}
	want := `{
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
	}`
	testJSONMarshal(t, u, want)

	u2 := &User{
		Login:                   Ptr("testLogin"),
		ID:                      Ptr(int64(1)),
		NodeID:                  Ptr("testNode123"),
		AvatarURL:               Ptr("https://www.example.com"),
		HTMLURL:                 Ptr("https://www.example.com"),
		GravatarID:              Ptr("testGravatar123"),
		Name:                    Ptr("myName"),
		Company:                 Ptr("testCompany"),
		Blog:                    Ptr("test Blog"),
		Location:                Ptr("test location"),
		Email:                   Ptr("test@example.com"),
		Hireable:                Ptr(true),
		Bio:                     Ptr("my good bio"),
		TwitterUsername:         Ptr("https://www.example.com/test"),
		PublicRepos:             Ptr(1),
		PublicGists:             Ptr(2),
		Followers:               Ptr(100),
		Following:               Ptr(29),
		CreatedAt:               &Timestamp{referenceTime},
		UpdatedAt:               &Timestamp{referenceTime},
		SuspendedAt:             &Timestamp{referenceTime},
		Type:                    Ptr("test type"),
		SiteAdmin:               Ptr(false),
		TotalPrivateRepos:       Ptr(int64(2)),
		OwnedPrivateRepos:       Ptr(int64(1)),
		PrivateGists:            Ptr(1),
		DiskUsage:               Ptr(1),
		Collaborators:           Ptr(1),
		TwoFactorAuthentication: Ptr(false),
		Plan: &Plan{
			Name:          Ptr("silver"),
			Space:         Ptr(1024),
			Collaborators: Ptr(10),
			PrivateRepos:  Ptr(int64(4)),
			FilledSeats:   Ptr(24),
			Seats:         Ptr(1),
		},
		LdapDn: Ptr("test ldap"),
	}

	want2 := `{
		"login": "testLogin",
		"id": 1,
		"node_id":"testNode123",
		"avatar_url": "https://www.example.com",
		"html_url":"https://www.example.com",
		"gravatar_id": "testGravatar123",
		"name": "myName",
		"company": "testCompany",
		"blog": "test Blog",
		"location": "test location",
		"email": "test@example.com",
		"hireable": true,
		"bio": "my good bio",
		"twitter_username": "https://www.example.com/test",
		"public_repos": 1,
		"public_gists":2,
		"followers": 100,
		"following": 29,
		"created_at": ` + referenceTimeStr + `,
		"suspended_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"type": "test type",
		"site_admin": false,
		"total_private_repos": 2,
		"owned_private_repos": 1,
		"private_gists": 1,
		"disk_usage": 1,
		"collaborators": 1,
		"two_factor_authentication": false,
		"plan": {
			"name": "silver",
			"space": 1024,
			"collaborators": 10,
			"private_repos": 4,
			"filled_seats": 24,
			"seats": 1
		},
		"ldap_dn": "test ldap"
	}`
	testJSONMarshal(t, u2, want2)
}

func TestUsersService_Get_authenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := t.Context()
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: Ptr(int64(1))}
	if !cmp.Equal(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.Get(ctx, "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_Get_specifiedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := t.Context()
	user, _, err := client.Users.Get(ctx, "u")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: Ptr(int64(1))}
	if !cmp.Equal(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Users.Get(ctx, "%")
	testURLParseError(t, err)
}

func TestUsersService_GetByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := t.Context()
	user, _, err := client.Users.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("Users.GetByID returned error: %v", err)
	}

	want := &User{ID: Ptr(int64(1))}
	if !cmp.Equal(user, want) {
		t.Errorf("Users.GetByID returned %+v, want %+v", user, want)
	}

	const methodName = "GetByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.GetByID(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetByID(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_Edit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &User{Name: Ptr("n")}

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		v := new(User)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := t.Context()
	user, _, err := client.Users.Edit(ctx, input)
	if err != nil {
		t.Errorf("Users.Edit returned error: %v", err)
	}

	want := &User{ID: Ptr(int64(1))}
	if !cmp.Equal(user, want) {
		t.Errorf("Users.Edit returned %+v, want %+v", user, want)
	}

	const methodName = "Edit"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.Edit(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_GetHovercard(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/hovercard", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"subject_type": "repository", "subject_id": "20180408"})
		fmt.Fprint(w, `{"contexts": [{"message":"Owns this repository", "octicon": "repo"}]}`)
	})

	opt := &HovercardOptions{SubjectType: "repository", SubjectID: "20180408"}
	ctx := t.Context()
	hovercard, _, err := client.Users.GetHovercard(ctx, "u", opt)
	if err != nil {
		t.Errorf("Users.GetHovercard returned error: %v", err)
	}

	want := &Hovercard{Contexts: []*UserContext{{Message: Ptr("Owns this repository"), Octicon: Ptr("repo")}}}
	if !cmp.Equal(hovercard, want) {
		t.Errorf("Users.GetHovercard returned %+v, want %+v", hovercard, want)
	}

	const methodName = "GetHovercard"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.GetHovercard(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetHovercard(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_ListAll(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1", "per_page": "30"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &UserListOptions{Since: 1, PerPage: 30}
	ctx := t.Context()
	users, _, err := client.Users.ListAll(ctx, opt)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := []*User{{ID: Ptr(int64(2))}}
	if !cmp.Equal(users, want) {
		t.Errorf("Users.ListAll returned %+v, want %+v", users, want)
	}

	const methodName = "ListAll"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListAll(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_ListInvitations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/repository_invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	ctx := t.Context()
	got, _, err := client.Users.ListInvitations(ctx, nil)
	if err != nil {
		t.Errorf("Users.ListInvitations returned error: %v", err)
	}

	want := []*RepositoryInvitation{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(got, want) {
		t.Errorf("Users.ListInvitations = %+v, want %+v", got, want)
	}

	const methodName = "ListInvitations"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListInvitations(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_ListInvitations_withOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/repository_invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	ctx := t.Context()
	_, _, err := client.Users.ListInvitations(ctx, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Users.ListInvitations returned error: %v", err)
	}
}

func TestUsersService_AcceptInvitation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/repository_invitations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	if _, err := client.Users.AcceptInvitation(ctx, 1); err != nil {
		t.Errorf("Users.AcceptInvitation returned error: %v", err)
	}

	const methodName = "AcceptInvitation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.AcceptInvitation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.AcceptInvitation(ctx, 1)
	})
}

func TestUsersService_DeclineInvitation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/repository_invitations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	if _, err := client.Users.DeclineInvitation(ctx, 1); err != nil {
		t.Errorf("Users.DeclineInvitation returned error: %v", err)
	}

	const methodName = "DeclineInvitation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.DeclineInvitation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DeclineInvitation(ctx, 1)
	})
}

func TestUserContext_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UserContext{}, "{}")

	u := &UserContext{
		Message: Ptr("message"),
		Octicon: Ptr("message"),
	}

	want := `{
		"message" : "message",
		"octicon" : "message"
	}`

	testJSONMarshal(t, u, want)
}

func TestHovercard_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Hovercard{}, "{}")

	h := &Hovercard{
		Contexts: []*UserContext{
			{
				Message: Ptr("someMessage"),
				Octicon: Ptr("someOcticon"),
			},
		},
	}

	want := `{
		"contexts" : [
			{
				"message" : "someMessage",
				"octicon" : "someOcticon"
			}
		]
	}`

	testJSONMarshal(t, h, want)
}
