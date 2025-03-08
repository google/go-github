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
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPullComments_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestComment{}, "{}")

	createdAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	updatedAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	reactions := &Reactions{
		TotalCount: Ptr(1),
		PlusOne:    Ptr(1),
		MinusOne:   Ptr(0),
		Laugh:      Ptr(0),
		Confused:   Ptr(0),
		Heart:      Ptr(0),
		Hooray:     Ptr(0),
		Rocket:     Ptr(0),
		Eyes:       Ptr(0),
		URL:        Ptr("u"),
	}

	u := &PullRequestComment{
		ID:                  Ptr(int64(10)),
		InReplyTo:           Ptr(int64(8)),
		Body:                Ptr("Test comment"),
		Path:                Ptr("file1.txt"),
		DiffHunk:            Ptr("@@ -16,33 +16,40 @@ fmt.Println()"),
		PullRequestReviewID: Ptr(int64(42)),
		Position:            Ptr(1),
		OriginalPosition:    Ptr(4),
		StartLine:           Ptr(2),
		Line:                Ptr(3),
		OriginalLine:        Ptr(2),
		OriginalStartLine:   Ptr(2),
		Side:                Ptr("RIGHT"),
		StartSide:           Ptr("LEFT"),
		CommitID:            Ptr("ab"),
		OriginalCommitID:    Ptr("9c"),
		User: &User{
			Login:       Ptr("ll"),
			ID:          Ptr(int64(123)),
			AvatarURL:   Ptr("a"),
			GravatarID:  Ptr("g"),
			Name:        Ptr("n"),
			Company:     Ptr("c"),
			Blog:        Ptr("b"),
			Location:    Ptr("l"),
			Email:       Ptr("e"),
			Hireable:    Ptr(true),
			PublicRepos: Ptr(1),
			Followers:   Ptr(1),
			Following:   Ptr(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         Ptr("u"),
		},
		Reactions:      reactions,
		CreatedAt:      &createdAt,
		UpdatedAt:      &updatedAt,
		URL:            Ptr("pullrequestcommentUrl"),
		HTMLURL:        Ptr("pullrequestcommentHTMLUrl"),
		PullRequestURL: Ptr("pullrequestcommentPullRequestURL"),
	}

	want := `{
		"id": 10,
		"in_reply_to_id": 8,
		"body": "Test comment",
		"path": "file1.txt",
		"diff_hunk": "@@ -16,33 +16,40 @@ fmt.Println()",
		"pull_request_review_id": 42,
		"position": 1,
		"original_position": 4,
		"start_line": 2,
		"line": 3,
		"original_line": 2,
		"original_start_line": 2,
		"side": "RIGHT",
		"start_side": "LEFT",
		"commit_id": "ab",
		"original_commit_id": "9c",
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
		"reactions": {
			"total_count": 1,
			"+1": 1,
			"-1": 0,
			"laugh": 0,
			"confused": 0,
			"heart": 0,
			"hooray": 0,
			"rocket": 0,
			"eyes": 0,
			"url": "u"
		},
		"created_at": "2002-02-10T15:30:00Z",
		"updated_at": "2002-02-10T15:30:00Z",
		"url": "pullrequestcommentUrl",
		"html_url": "pullrequestcommentHTMLUrl",
		"pull_request_url": "pullrequestcommentPullRequestURL"
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestsService_ListComments_allPulls(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{
			"sort":      "updated",
			"direction": "desc",
			"since":     "2002-02-10T15:30:00Z",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &PullRequestListCommentsOptions{
		Sort:        "updated",
		Direction:   "desc",
		Since:       time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	pulls, _, err := client.PullRequests.ListComments(ctx, "o", "r", 0, opt)
	if err != nil {
		t.Errorf("PullRequests.ListComments returned error: %v", err)
	}

	want := []*PullRequestComment{{ID: Ptr(int64(1))}}
	if !cmp.Equal(pulls, want) {
		t.Errorf("PullRequests.ListComments returned %+v, want %+v", pulls, want)
	}

	const methodName = "ListComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListComments(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListComments(ctx, "o", "r", 0, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListComments_specificPull(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1, "pull_request_review_id":42}]`)
	})

	ctx := context.Background()
	pulls, _, err := client.PullRequests.ListComments(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("PullRequests.ListComments returned error: %v", err)
	}

	want := []*PullRequestComment{{ID: Ptr(int64(1)), PullRequestReviewID: Ptr(int64(42))}}
	if !cmp.Equal(pulls, want) {
		t.Errorf("PullRequests.ListComments returned %+v, want %+v", pulls, want)
	}
}

func TestPullRequestsService_ListComments_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.PullRequests.ListComments(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_GetComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.PullRequests.GetComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.GetComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("PullRequests.GetComment returned %+v, want %+v", comment, want)
	}

	const methodName = "GetComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.GetComment(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.GetComment(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_GetComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.PullRequests.GetComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PullRequestComment{Body: Ptr("b")}

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestComment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		// TODO: remove custom Accept header assertion when the API fully launches.
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.PullRequests.CreateComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.CreateComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("PullRequests.CreateComment returned %+v, want %+v", comment, want)
	}

	const methodName = "CreateComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.CreateComment(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.CreateComment(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_CreateComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.PullRequests.CreateComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateCommentInReplyTo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PullRequestComment{Body: Ptr("b")}

	mux.HandleFunc("/repos/o/r/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestComment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.PullRequests.CreateCommentInReplyTo(ctx, "o", "r", 1, "b", 2)
	if err != nil {
		t.Errorf("PullRequests.CreateCommentInReplyTo returned error: %v", err)
	}

	want := &PullRequestComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("PullRequests.CreateCommentInReplyTo returned %+v, want %+v", comment, want)
	}

	const methodName = "CreateCommentInReplyTo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.CreateCommentInReplyTo(ctx, "\n", "\n", -1, "\n", -2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.CreateCommentInReplyTo(ctx, "o", "r", 1, "b", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_EditComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PullRequestComment{Body: Ptr("b")}

	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestComment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.PullRequests.EditComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.EditComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Ptr(int64(1))}
	if !cmp.Equal(comment, want) {
		t.Errorf("PullRequests.EditComment returned %+v, want %+v", comment, want)
	}

	const methodName = "EditComment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.EditComment(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.EditComment(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_EditComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.PullRequests.EditComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_DeleteComment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.PullRequests.DeleteComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.DeleteComment returned error: %v", err)
	}

	const methodName = "DeleteComment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.PullRequests.DeleteComment(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.PullRequests.DeleteComment(ctx, "o", "r", 1)
	})
}

func TestPullRequestsService_DeleteComment_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.PullRequests.DeleteComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}
