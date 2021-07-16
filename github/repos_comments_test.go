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

func TestRepositoriesService_ListComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	comments, _, err := client.Repositories.ListComments(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListComments returned error: %v", err)
	}

	want := []*RepositoryComment{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(comments, want) {
		t.Errorf("Repositories.ListComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListComments(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListComments(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListComments(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListCommitComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/commits/s/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	comments, _, err := client.Repositories.ListCommitComments(ctx, "o", "r", "s", opt)
	if err != nil {
		t.Errorf("Repositories.ListCommitComments returned error: %v", err)
	}

	want := []*RepositoryComment{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(comments, want) {
		t.Errorf("Repositories.ListCommitComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListCommitComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCommitComments(ctx, "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCommitComments(ctx, "o", "r", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListCommitComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListCommitComments(ctx, "%", "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RepositoryComment{Body: String("b")}

	mux.HandleFunc("/repos/o/r/commits/s/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Repositories.CreateComment(ctx, "o", "r", "s", input)
	if err != nil {
		t.Errorf("Repositories.CreateComment returned error: %v", err)
	}

	want := &RepositoryComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Repositories.CreateComment returned %+v, want %+v", comment, want)
	}

	const methodName = "CreateComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateComment(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateComment(ctx, "o", "r", "s", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.CreateComment(ctx, "%", "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Repositories.GetComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetComment returned error: %v", err)
	}

	want := &RepositoryComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Repositories.GetComment returned %+v, want %+v", comment, want)
	}

	const methodName = "GetComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetComment(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetComment(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.GetComment(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoriesService_UpdateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RepositoryComment{Body: String("b")}

	mux.HandleFunc("/repos/o/r/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Repositories.UpdateComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.UpdateComment returned error: %v", err)
	}

	want := &RepositoryComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Repositories.UpdateComment returned %+v, want %+v", comment, want)
	}

	const methodName = "UpdateComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdateComment(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateComment(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.UpdateComment(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteComment returned error: %v", err)
	}

	const methodName = "DeleteComment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteComment(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteComment(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_DeleteComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Repositories.DeleteComment(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoryComment_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryComment{}, "{}")

	r := &RepositoryComment{
		HTMLURL:  String("hurl"),
		URL:      String("url"),
		ID:       Int64(1),
		NodeID:   String("nid"),
		CommitID: String("cid"),
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
		Reactions: &Reactions{
			TotalCount: Int(1),
			PlusOne:    Int(1),
			MinusOne:   Int(1),
			Laugh:      Int(1),
			Confused:   Int(1),
			Heart:      Int(1),
			Hooray:     Int(1),
			Rocket:     Int(1),
			Eyes:       Int(1),
			URL:        String("u"),
		},
		CreatedAt: &referenceTime,
		UpdatedAt: &referenceTime,
		Body:      String("body"),
		Path:      String("path"),
		Position:  Int(1),
	}

	want := `{
		"html_url": "hurl",
		"url": "url",
		"id": 1,
		"node_id": "nid",
		"commit_id": "cid",
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
		},
		"reactions": {
			"total_count": 1,
			"+1": 1,
			"-1": 1,
			"laugh": 1,
			"confused": 1,
			"heart": 1,
			"hooray": 1,
			"rocket": 1,
			"eyes": 1,
			"url": "u"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"body": "body",
		"path": "path",
		"position": 1
	}`

	testJSONMarshal(t, r, want)
}
