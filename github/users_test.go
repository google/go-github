// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUser_Marshal(t *testing.T) {
	testJSONMarshal(t, &User{}, "{}")

	u := &User{
		Login:           String("l"),
		ID:              Int64(1),
		URL:             String("u"),
		AvatarURL:       String("a"),
		GravatarID:      String("g"),
		Name:            String("n"),
		Company:         String("c"),
		Blog:            String("b"),
		Location:        String("l"),
		Email:           String("e"),
		Hireable:        Bool(true),
		Bio:             String("b"),
		TwitterUsername: String("t"),
		PublicRepos:     Int(1),
		Followers:       Int(1),
		Following:       Int(1),
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
}

func TestUsersService_Get_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: Int64(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "u")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: Int64(1)}
	if !cmp.Equal(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Users.Get(ctx, "%")
	testURLParseError(t, err)
}

func TestUsersService_GetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("Users.GetByID returned error: %v", err)
	}

	want := &User{ID: Int64(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &User{Name: String("n")}

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		v := new(User)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Edit(ctx, input)
	if err != nil {
		t.Errorf("Users.Edit returned error: %v", err)
	}

	want := &User{ID: Int64(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/hovercard", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"subject_type": "repository", "subject_id": "20180408"})
		fmt.Fprint(w, `{"contexts": [{"message":"Owns this repository", "octicon": "repo"}]}`)
	})

	opt := &HovercardOptions{SubjectType: "repository", SubjectID: "20180408"}
	ctx := context.Background()
	hovercard, _, err := client.Users.GetHovercard(ctx, "u", opt)
	if err != nil {
		t.Errorf("Users.GetHovercard returned error: %v", err)
	}

	want := &Hovercard{Contexts: []*UserContext{{Message: String("Owns this repository"), Octicon: String("repo")}}}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1", "page": "2"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &UserListOptions{1, ListOptions{Page: 2}}
	ctx := context.Background()
	users, _, err := client.Users.ListAll(ctx, opt)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := []*User{{ID: Int64(2)}}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/repository_invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	ctx := context.Background()
	got, _, err := client.Users.ListInvitations(ctx, nil)
	if err != nil {
		t.Errorf("Users.ListInvitations returned error: %v", err)
	}

	want := []*RepositoryInvitation{{ID: Int64(1)}, {ID: Int64(2)}}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/repository_invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	ctx := context.Background()
	_, _, err := client.Users.ListInvitations(ctx, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Users.ListInvitations returned error: %v", err)
	}
}

func TestUsersService_AcceptInvitation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/repository_invitations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/repository_invitations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
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
	testJSONMarshal(t, &UserContext{}, "{}")

	u := &UserContext{
		Message: String("message"),
		Octicon: String("message"),
	}

	want := `{
		"message" : "message",
		"octicon" : "message"
	}`

	testJSONMarshal(t, u, want)
}

func TestHovercard_Marshal(t *testing.T) {
	testJSONMarshal(t, &Hovercard{}, "{}")

	h := &Hovercard{
		Contexts: []*UserContext{
			{
				Message: String("someMessage"),
				Octicon: String("someOcticon"),
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

func TestUserListOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &UserListOptions{}, "{}")

	u := &UserListOptions{
		Since: int64(1900),
		ListOptions: ListOptions{
			Page:    int(1),
			PerPage: int(10),
		},
	}

	want := `{
		"since" : 1900,
		"page": 1,
		"perPage": 10
	}`

	testJSONMarshal(t, u, want)
}
