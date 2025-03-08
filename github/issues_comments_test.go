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

func TestIssuesService_ListComments_allIssues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{
			"sort":      "updated",
			"direction": "desc",
			"since":     "2002-02-10T15:30:00Z",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	since := time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)
	opt := &IssueListCommentsOptions{
		Sort:        Ptr("updated"),
		Direction:   Ptr("desc"),
		Since:       &since,
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, "o", "r", 0, opt)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []*IssueComment{{ID: Ptr(int64(1))}}
	if !cmp.Equal(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListComments(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListComments(ctx, "o", "r", 0, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListComments_specificIssue(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []*IssueComment{{ID: Ptr(int64(1))}}
	if !cmp.Equal(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListComments(ctx, "\n", "\n", -1, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListComments(ctx, "o", "r", 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListComments_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Issues.ListComments(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_GetComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.GetComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetComment returned error: %v", err)
	}

	want := &IssueComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("Issues.GetComment returned %+v, want %+v", comment, want)
	}

	const methodName = "GetComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.GetComment(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.GetComment(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_GetComment_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Issues.GetComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssuesService_CreateComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &IssueComment{Body: Ptr("b")}

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.CreateComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.CreateComment returned error: %v", err)
	}

	want := &IssueComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("Issues.CreateComment returned %+v, want %+v", comment, want)
	}

	const methodName = "CreateComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.CreateComment(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.CreateComment(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_CreateComment_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Issues.CreateComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_EditComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &IssueComment{Body: Ptr("b")}

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.EditComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.EditComment returned error: %v", err)
	}

	want := &IssueComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("Issues.EditComment returned %+v, want %+v", comment, want)
	}

	const methodName = "EditComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.EditComment(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.EditComment(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_EditComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Issues.EditComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_DeleteComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Issues.DeleteComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.DeleteComments returned error: %v", err)
	}

	const methodName = "DeleteComment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Issues.DeleteComment(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Issues.DeleteComment(ctx, "o", "r", 1)
	})
}

func TestIssuesService_DeleteComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Issues.DeleteComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssueComment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueComment{}, "{}")

	u := &IssueComment{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("nid"),
		Body:   Ptr("body"),
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
		Reactions:         &Reactions{TotalCount: Ptr(1)},
		CreatedAt:         &Timestamp{referenceTime},
		UpdatedAt:         &Timestamp{referenceTime},
		AuthorAssociation: Ptr("aa"),
		URL:               Ptr("url"),
		HTMLURL:           Ptr("hurl"),
		IssueURL:          Ptr("iurl"),
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"body": "body",
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
			"total_count": 1
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"author_association": "aa",
		"url": "url",
		"html_url": "hurl",
		"issue_url": "iurl"
	}`

	testJSONMarshal(t, u, want)
}
