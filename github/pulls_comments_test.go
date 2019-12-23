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
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPullComments_marshall(t *testing.T) {
	testJSONMarshal(t, &PullRequestComment{}, "{}")

	createdAt := time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)
	reactions := &Reactions{
		TotalCount: Int(1),
		PlusOne:    Int(1),
		MinusOne:   Int(0),
		Laugh:      Int(0),
		Confused:   Int(0),
		Heart:      Int(0),
		Hooray:     Int(0),
		URL:        String("u"),
	}

	u := &PullRequestComment{
		ID:                  Int64(10),
		InReplyTo:           Int64(8),
		Body:                String("Test comment"),
		Path:                String("file1.txt"),
		DiffHunk:            String("@@ -16,33 +16,40 @@ fmt.Println()"),
		PullRequestReviewID: Int64(42),
		Position:            Int(1),
		OriginalPosition:    Int(4),
		StartLine:           Int(2),
		Line:                Int(3),
		OriginalLine:        Int(2),
		OriginalStartLine:   Int(2),
		Side:                String("RIGHT"),
		StartSide:           String("LEFT"),
		CommitID:            String("ab"),
		OriginalCommitID:    String("9c"),
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
		Reactions:      reactions,
		CreatedAt:      &createdAt,
		UpdatedAt:      &updatedAt,
		URL:            String("pullrequestcommentUrl"),
		HTMLURL:        String("pullrequestcommentHTMLUrl"),
		PullRequestURL: String("pullrequestcommentPullRequestURL"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
	pulls, _, err := client.PullRequests.ListComments(context.Background(), "o", "r", 0, opt)
	if err != nil {
		t.Errorf("PullRequests.ListComments returned error: %v", err)
	}

	want := []*PullRequestComment{{ID: Int64(1)}}
	if !reflect.DeepEqual(pulls, want) {
		t.Errorf("PullRequests.ListComments returned %+v, want %+v", pulls, want)
	}
}

func TestPullRequestsService_ListComments_specificPull(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1, "pull_request_review_id":42}]`)
	})

	pulls, _, err := client.PullRequests.ListComments(context.Background(), "o", "r", 1, nil)
	if err != nil {
		t.Errorf("PullRequests.ListComments returned error: %v", err)
	}

	want := []*PullRequestComment{{ID: Int64(1), PullRequestReviewID: Int64(42)}}
	if !reflect.DeepEqual(pulls, want) {
		t.Errorf("PullRequests.ListComments returned %+v, want %+v", pulls, want)
	}
}

func TestPullRequestsService_ListComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.PullRequests.ListComments(context.Background(), "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_GetComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	comment, _, err := client.PullRequests.GetComment(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.GetComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("PullRequests.GetComment returned %+v, want %+v", comment, want)
	}
}

func TestPullRequestsService_GetComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.PullRequests.GetComment(context.Background(), "%", "r", 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestComment{Body: String("b")}

	wantAcceptHeaders := []string{mediaTypeReactionsPreview, mediaTypeMultiLineCommentsPreview}
	mux.HandleFunc("/repos/o/r/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestComment)
		json.NewDecoder(r.Body).Decode(v)

		// TODO: remove custom Accept header assertion when the API fully launches.
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	comment, _, err := client.PullRequests.CreateComment(context.Background(), "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.CreateComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("PullRequests.CreateComment returned %+v, want %+v", comment, want)
	}
}

func TestPullRequestsService_CreateComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.PullRequests.CreateComment(context.Background(), "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_EditComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestComment{Body: String("b")}

	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	comment, _, err := client.PullRequests.EditComment(context.Background(), "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.EditComment returned error: %v", err)
	}

	want := &PullRequestComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("PullRequests.EditComment returned %+v, want %+v", comment, want)
	}
}

func TestPullRequestsService_EditComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.PullRequests.EditComment(context.Background(), "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.PullRequests.DeleteComment(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.DeleteComment returned error: %v", err)
	}
}

func TestPullRequestsService_DeleteComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.PullRequests.DeleteComment(context.Background(), "%", "r", 1)
	testURLParseError(t, err)
}
