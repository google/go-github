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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGistComments_Marshal(t *testing.T) {
	testJSONMarshal(t, &GistComment{}, "{}")

	createdAt := time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)

	u := &GistComment{
		ID:   Int64(1),
		URL:  String("u"),
		Body: String("test gist comment"),
		User: &User{
			Login:       String("ll"),
			ID:          Int64(123),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		},
		CreatedAt: &createdAt,
	}

	want := `{
		"id": 1,
		"url": "u",
		"body": "test gist comment",
		"user": {
			"login": "ll",
			"id": 123,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"created_at": "2002-02-10T15:30:00Z"
	}`

	testJSONMarshal(t, u, want)
}
func TestGistsService_ListComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id": 1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	comments, _, err := client.Gists.ListComments(ctx, "1", opt)
	if err != nil {
		t.Errorf("Gists.Comments returned error: %v", err)
	}

	want := []*GistComment{{ID: Int64(1)}}
	if !cmp.Equal(comments, want) {
		t.Errorf("Gists.ListComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.ListComments(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListComments(ctx, "1", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListComments_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.ListComments(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_GetComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": 1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Gists.GetComment(ctx, "1", 2)
	if err != nil {
		t.Errorf("Gists.GetComment returned error: %v", err)
	}

	want := &GistComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Gists.GetComment returned %+v, want %+v", comment, want)
	}

	const methodName = "GetComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.GetComment(ctx, "\n", -2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.GetComment(ctx, "1", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_GetComment_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.GetComment(ctx, "%", 1)
	testURLParseError(t, err)
}

func TestGistsService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &GistComment{ID: Int64(1), Body: String("b")}

	mux.HandleFunc("/gists/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(GistComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Gists.CreateComment(ctx, "1", input)
	if err != nil {
		t.Errorf("Gists.CreateComment returned error: %v", err)
	}

	want := &GistComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Gists.CreateComment returned %+v, want %+v", comment, want)
	}

	const methodName = "CreateComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.CreateComment(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.CreateComment(ctx, "1", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_CreateComment_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.CreateComment(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_EditComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &GistComment{ID: Int64(1), Body: String("b")}

	mux.HandleFunc("/gists/1/comments/2", func(w http.ResponseWriter, r *http.Request) {
		v := new(GistComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Gists.EditComment(ctx, "1", 2, input)
	if err != nil {
		t.Errorf("Gists.EditComment returned error: %v", err)
	}

	want := &GistComment{ID: Int64(1)}
	if !cmp.Equal(comment, want) {
		t.Errorf("Gists.EditComment returned %+v, want %+v", comment, want)
	}

	const methodName = "EditComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.EditComment(ctx, "\n", -2, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.EditComment(ctx, "1", 2, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_EditComment_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.EditComment(ctx, "%", 1, nil)
	testURLParseError(t, err)
}

func TestGistsService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Gists.DeleteComment(ctx, "1", 2)
	if err != nil {
		t.Errorf("Gists.Delete returned error: %v", err)
	}

	const methodName = "DeleteComment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Gists.DeleteComment(ctx, "\n", -2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Gists.DeleteComment(ctx, "1", 2)
	})
}

func TestGistsService_DeleteComment_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Gists.DeleteComment(ctx, "%", 1)
	testURLParseError(t, err)
}
