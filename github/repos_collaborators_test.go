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

func TestRepositoriesService_ListCollaborators(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListCollaboratorsOptions{
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	users, _, err := client.Repositories.ListCollaborators(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(users, want) {
		t.Errorf("Repositories.ListCollaborators returned %+v, want %+v", users, want)
	}

	const methodName = "ListCollaborators"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCollaborators(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCollaborators(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListCollaborators_withAffiliation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"affiliation": "all", "page": "2"})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListCollaboratorsOptions{
		ListOptions: ListOptions{Page: 2},
		Affiliation: "all",
	}
	ctx := context.Background()
	users, _, err := client.Repositories.ListCollaborators(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(users, want) {
		t.Errorf("Repositories.ListCollaborators returned %+v, want %+v", users, want)
	}

	const methodName = "ListCollaborators"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCollaborators(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCollaborators(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListCollaborators_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListCollaborators(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_IsCollaborator_True(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	isCollab, _, err := client.Repositories.IsCollaborator(ctx, "o", "r", "u")
	if err != nil {
		t.Errorf("Repositories.IsCollaborator returned error: %v", err)
	}

	if !isCollab {
		t.Errorf("Repositories.IsCollaborator returned false, want true")
	}

	const methodName = "IsCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.IsCollaborator(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.IsCollaborator(ctx, "o", "r", "u")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_IsCollaborator_False(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	isCollab, _, err := client.Repositories.IsCollaborator(ctx, "o", "r", "u")
	if err != nil {
		t.Errorf("Repositories.IsCollaborator returned error: %v", err)
	}

	if isCollab {
		t.Errorf("Repositories.IsCollaborator returned true, want false")
	}

	const methodName = "IsCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.IsCollaborator(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.IsCollaborator(ctx, "o", "r", "u")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_IsCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.IsCollaborator(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestRepositoryService_GetPermissionLevel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators/u/permission", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"permission":"admin","user":{"login":"u"}}`)
	})

	ctx := context.Background()
	rpl, _, err := client.Repositories.GetPermissionLevel(ctx, "o", "r", "u")
	if err != nil {
		t.Errorf("Repositories.GetPermissionLevel returned error: %v", err)
	}

	want := &RepositoryPermissionLevel{
		Permission: String("admin"),
		User: &User{
			Login: String("u"),
		},
	}

	if !cmp.Equal(rpl, want) {
		t.Errorf("Repositories.GetPermissionLevel returned %+v, want %+v", rpl, want)
	}

	const methodName = "GetPermissionLevel"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetPermissionLevel(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetPermissionLevel(ctx, "o", "r", "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_AddCollaborator(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &RepositoryAddCollaboratorOptions{Permission: "admin"}
	mux.HandleFunc("/repos/o/r/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryAddCollaboratorOptions)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"permissions": "write","url": "https://api.github.com/user/repository_invitations/1296269","html_url": "https://github.com/octocat/Hello-World/invitations","id":1,"permissions":"write","repository":{"url":"s","name":"r","id":1},"invitee":{"login":"u"},"inviter":{"login":"o"}}`))
	})
	ctx := context.Background()
	collaboratorInvitation, _, err := client.Repositories.AddCollaborator(ctx, "o", "r", "u", opt)
	if err != nil {
		t.Errorf("Repositories.AddCollaborator returned error: %v", err)
	}
	want := &CollaboratorInvitation{
		ID: Int64(1),
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("r"),
		},
		Invitee: &User{
			Login: String("u"),
		},
		Inviter: &User{
			Login: String("o"),
		},
		Permissions: String("write"),
		URL:         String("https://api.github.com/user/repository_invitations/1296269"),
		HTMLURL:     String("https://github.com/octocat/Hello-World/invitations"),
	}

	if !cmp.Equal(collaboratorInvitation, want) {
		t.Errorf("AddCollaborator returned %+v, want %+v", collaboratorInvitation, want)
	}

	const methodName = "AddCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.AddCollaborator(ctx, "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.AddCollaborator(ctx, "o", "r", "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_AddCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.AddCollaborator(ctx, "%", "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_RemoveCollaborator(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.RemoveCollaborator(ctx, "o", "r", "u")
	if err != nil {
		t.Errorf("Repositories.RemoveCollaborator returned error: %v", err)
	}

	const methodName = "RemoveCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.RemoveCollaborator(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.RemoveCollaborator(ctx, "o", "r", "u")
	})
}

func TestRepositoriesService_RemoveCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Repositories.RemoveCollaborator(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestRepositoryAddCollaboratorOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryAddCollaboratorOptions{}, "{}")

	r := &RepositoryAddCollaboratorOptions{
		Permission: "permission",
	}

	want := `{
		"permission": "permission"
	}`

	testJSONMarshal(t, r, want)
}

func TestRepositoryPermissionLevel_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryPermissionLevel{}, "{}")

	r := &RepositoryPermissionLevel{
		Permission: String("permission"),
		User: &User{
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
		},
	}

	want := `{
		"permission": "permission",
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

	testJSONMarshal(t, r, want)
}

func TestCollaboratorInvitation_Marshal(t *testing.T) {
	testJSONMarshal(t, &CollaboratorInvitation{}, "{}")

	r := &CollaboratorInvitation{
		ID: Int64(1),
		Repo: &Repository{

			ID:   Int64(1),
			URL:  String("url"),
			Name: String("n"),
		},
		Invitee: &User{
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
		},
		Inviter: &User{
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
		},
		Permissions: String("per"),
		CreatedAt:   &Timestamp{referenceTime},
		URL:         String("url"),
		HTMLURL:     String("hurl"),
	}

	want := `{
		"id": 1,
		"repository": {
			"id": 1,
			"name": "n",
			"url": "url"
		},
		"invitee": {
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
		},
		"inviter": {
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
		},
		"permissions": "per",
		"created_at": ` + referenceTimeStr + `,
		"url": "url",
		"html_url": "hurl"
	}`

	testJSONMarshal(t, r, want)
}
