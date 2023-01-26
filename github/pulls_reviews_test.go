// Copyright 2016 The go-github AUTHORS. All rights reserved.
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

func TestPullRequestsService_ListReviews(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	reviews, _, err := client.PullRequests.ListReviews(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("PullRequests.ListReviews returned error: %v", err)
	}

	want := []*PullRequestReview{
		{ID: Int64(1)},
		{ID: Int64(2)},
	}
	if !cmp.Equal(reviews, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", reviews, want)
	}

	const methodName = "ListReviews"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListReviews(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListReviews(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListReviews_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviews(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_GetReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.GetReview(ctx, "o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.GetReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(review, want) {
		t.Errorf("PullRequests.GetReview returned %+v, want %+v", review, want)
	}

	const methodName = "GetReview"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.GetReview(ctx, "\n", "\n", -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.GetReview(ctx, "o", "r", 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_GetReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.GetReview(ctx, "%", "r", 1, 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_DeletePendingReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.DeletePendingReview(ctx, "o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.DeletePendingReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(review, want) {
		t.Errorf("PullRequests.DeletePendingReview returned %+v, want %+v", review, want)
	}

	const methodName = "DeletePendingReview"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.DeletePendingReview(ctx, "\n", "\n", -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.DeletePendingReview(ctx, "o", "r", 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_DeletePendingReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.DeletePendingReview(ctx, "%", "r", 1, 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_ListReviewComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	comments, _, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, nil)
	if err != nil {
		t.Errorf("PullRequests.ListReviewComments returned error: %v", err)
	}

	want := []*PullRequestComment{
		{ID: Int64(1)},
		{ID: Int64(2)},
	}
	if !cmp.Equal(comments, want) {
		t.Errorf("PullRequests.ListReviewComments returned %+v, want %+v", comments, want)
	}

	const methodName = "ListReviewComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListReviewComments(ctx, "\n", "\n", -1, -1, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListReviewComments_withOptions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("PullRequests.ListReviewComments returned error: %v", err)
	}

	const methodName = "ListReviewComments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListReviewComments(ctx, "\n", "\n", -1, -1, &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestReviewRequest_isComfortFadePreview(t *testing.T) {
	path := "path/to/file.go"
	body := "this is a comment body"
	left, right := "LEFT", "RIGHT"
	pos1, pos2, pos3 := 1, 2, 3
	line1, line2, line3 := 11, 22, 33

	tests := []struct {
		name     string
		review   *PullRequestReviewRequest
		wantErr  error
		wantBool bool
	}{{
		name:     "empty review",
		review:   &PullRequestReviewRequest{},
		wantBool: false,
	}, {
		name:     "nil comment",
		review:   &PullRequestReviewRequest{Comments: []*DraftReviewComment{nil}},
		wantBool: false,
	}, {
		name: "old-style review",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos2,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos3,
			}},
		},
		wantBool: false,
	}, {
		name: "new-style review",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}, {
				Path: &path,
				Body: &body,
				Side: &left,
				Line: &line2,
			}, {
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line3,
			}},
		},
		wantBool: true,
	}, {
		name: "blended comment",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1, // can't have both styles.
				Side:     &right,
				Line:     &line1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}, {
		name: "position then line",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}, {
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}, {
		name: "line then position",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotBool, gotErr := tc.review.isComfortFadePreview()
			if tc.wantErr != nil {
				if gotErr != tc.wantErr {
					t.Errorf("isComfortFadePreview() = %v, wanted %v", gotErr, tc.wantErr)
				}
			} else {
				if gotBool != tc.wantBool {
					t.Errorf("isComfortFadePreview() = %v, wanted %v", gotBool, tc.wantBool)
				}
			}
		})
	}
}

func TestPullRequestsService_ListReviewComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviewComments(ctx, "%", "r", 1, 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewRequest{
		CommitID: String("commit_id"),
		Body:     String("b"),
		Event:    String("APPROVE"),
	}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.CreateReview(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.CreateReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(review, want) {
		t.Errorf("PullRequests.CreateReview returned %+v, want %+v", review, want)
	}

	const methodName = "CreateReview"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.CreateReview(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.CreateReview(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_CreateReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.CreateReview(ctx, "%", "r", 1, &PullRequestReviewRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateReview_badReview(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	path := "path/to/file.go"
	body := "this is a comment body"
	right := "RIGHT"
	pos1 := 1
	line1 := 11
	badReview := &PullRequestReviewRequest{
		Comments: []*DraftReviewComment{{
			Path: &path,
			Body: &body,
			Side: &right,
			Line: &line1,
		}, {
			Path:     &path,
			Body:     &body,
			Position: &pos1,
		}}}

	_, _, err := client.PullRequests.CreateReview(ctx, "o", "r", 1, badReview)
	if err == nil {
		t.Errorf("CreateReview badReview err = nil, want err")
	}
}

func TestPullRequestsService_CreateReview_addHeader(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	path := "path/to/file.go"
	body := "this is a comment body"
	left, right := "LEFT", "RIGHT"
	line1, line2, line3 := 11, 22, 33
	input := &PullRequestReviewRequest{
		Comments: []*DraftReviewComment{{
			Path: &path,
			Body: &body,
			Side: &right,
			Line: &line1,
		}, {
			Path: &path,
			Body: &body,
			Side: &left,
			Line: &line2,
		}, {
			Path: &path,
			Body: &body,
			Side: &right,
			Line: &line3,
		}},
	}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()

	_, _, err := client.PullRequests.CreateReview(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("CreateReview addHeader err = %v, want nil", err)
	}
}

func TestPullRequestsService_UpdateReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.PullRequests.UpdateReview(ctx, "o", "r", 1, 1, "updated_body")
	if err != nil {
		t.Errorf("PullRequests.UpdateReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("PullRequests.UpdateReview = %+v, want %+v", got, want)
	}

	const methodName = "UpdateReview"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.UpdateReview(ctx, "\n", "\n", -1, -1, "updated_body")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.UpdateReview(ctx, "o", "r", 1, 1, "updated_body")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_SubmitReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewRequest{
		Body:  String("b"),
		Event: String("APPROVE"),
	}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/events", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.SubmitReview(ctx, "o", "r", 1, 1, input)
	if err != nil {
		t.Errorf("PullRequests.SubmitReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(review, want) {
		t.Errorf("PullRequests.SubmitReview returned %+v, want %+v", review, want)
	}

	const methodName = "SubmitReview"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.SubmitReview(ctx, "\n", "\n", -1, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.SubmitReview(ctx, "o", "r", 1, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_SubmitReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.SubmitReview(ctx, "%", "r", 1, 1, &PullRequestReviewRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_DismissReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewDismissalRequest{Message: String("m")}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/dismissals", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewDismissalRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.DismissReview(ctx, "o", "r", 1, 1, input)
	if err != nil {
		t.Errorf("PullRequests.DismissReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !cmp.Equal(review, want) {
		t.Errorf("PullRequests.DismissReview returned %+v, want %+v", review, want)
	}

	const methodName = "ListReviews"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.DismissReview(ctx, "\n", "\n", -1, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.DismissReview(ctx, "o", "r", 1, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_DismissReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.DismissReview(ctx, "%", "r", 1, 1, &PullRequestReviewDismissalRequest{})
	testURLParseError(t, err)
}

func TestPullRequestReviewDismissalRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReviewDismissalRequest{}, "{}")

	u := &PullRequestReviewDismissalRequest{
		Message: String("msg"),
	}

	want := `{
		"message": "msg"
	}`

	testJSONMarshal(t, u, want)
}

func TestDraftReviewComment_Marshal(t *testing.T) {
	testJSONMarshal(t, &DraftReviewComment{}, "{}")

	u := &DraftReviewComment{
		Path:      String("path"),
		Position:  Int(1),
		Body:      String("body"),
		StartSide: String("ss"),
		Side:      String("side"),
		StartLine: Int(1),
		Line:      Int(1),
	}

	want := `{
		"path": "path",
		"position": 1,
		"body": "body",
		"start_side": "ss",
		"side": "side",
		"start_line": 1,
		"line": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestReviewRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReviewRequest{}, "{}")

	u := &PullRequestReviewRequest{
		NodeID:   String("nodeid"),
		CommitID: String("cid"),
		Body:     String("body"),
		Event:    String("event"),
		Comments: []*DraftReviewComment{
			{
				Path:      String("path"),
				Position:  Int(1),
				Body:      String("body"),
				StartSide: String("ss"),
				Side:      String("side"),
				StartLine: Int(1),
				Line:      Int(1),
			},
		},
	}

	want := `{
		"node_id": "nodeid",
		"commit_id": "cid",
		"body": "body",
		"event": "event",
		"comments": [
			{
				"path": "path",
				"position": 1,
				"body": "body",
				"start_side": "ss",
				"side": "side",
				"start_line": 1,
				"line": 1
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestReview_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReview{}, "{}")

	u := &PullRequestReview{
		ID:     Int64(1),
		NodeID: String("nid"),
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
		Body:              String("body"),
		SubmittedAt:       &Timestamp{referenceTime},
		CommitID:          String("cid"),
		HTMLURL:           String("hurl"),
		PullRequestURL:    String("prurl"),
		State:             String("state"),
		AuthorAssociation: String("aa"),
	}

	want := `{
		"id": 1,
		"node_id": "nid",
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
		"body": "body",
		"submitted_at": ` + referenceTimeStr + `,
		"commit_id": "cid",
		"html_url": "hurl",
		"pull_request_url": "prurl",
		"state": "state",
		"author_association": "aa"
	}`

	testJSONMarshal(t, u, want)
}
