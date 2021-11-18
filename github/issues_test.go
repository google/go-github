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

func TestIssuesService_List_all(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{
			"filter":    "all",
			"state":     "closed",
			"labels":    "a,b",
			"sort":      "updated",
			"direction": "asc",
			"since":     "2002-02-10T15:30:00Z",
			"page":      "1",
			"per_page":  "2",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &IssueListOptions{
		"all", "closed", []string{"a", "b"}, "updated", "asc",
		time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
		ListOptions{Page: 1, PerPage: 2},
	}
	ctx := context.Background()
	issues, _, err := client.Issues.List(ctx, true, opt)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []*Issue{{Number: Int(1)}}
	if !cmp.Equal(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.List(ctx, true, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_List_owned(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `[{"number":1}]`)
	})

	ctx := context.Background()
	issues, _, err := client.Issues.List(ctx, false, nil)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []*Issue{{Number: Int(1)}}
	if !cmp.Equal(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}
}

func TestIssuesService_ListByOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `[{"number":1}]`)
	})

	ctx := context.Background()
	issues, _, err := client.Issues.ListByOrg(ctx, "o", nil)
	if err != nil {
		t.Errorf("Issues.ListByOrg returned error: %v", err)
	}

	want := []*Issue{{Number: Int(1)}}
	if !cmp.Equal(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}

	const methodName = "ListByOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListByOrg(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListByOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListByOrg_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListByOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestIssuesService_ListByRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{
			"milestone": "*",
			"state":     "closed",
			"assignee":  "a",
			"creator":   "c",
			"mentioned": "m",
			"labels":    "a,b",
			"sort":      "updated",
			"direction": "asc",
			"since":     "2002-02-10T15:30:00Z",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &IssueListByRepoOptions{
		"*", "closed", "a", "c", "m", []string{"a", "b"}, "updated", "asc",
		time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
		ListOptions{0, 0},
	}
	ctx := context.Background()
	issues, _, err := client.Issues.ListByRepo(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListByOrg returned error: %v", err)
	}

	want := []*Issue{{Number: Int(1)}}
	if !cmp.Equal(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}

	const methodName = "ListByRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListByRepo(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListByRepo(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListByRepo_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListByRepo(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestIssuesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `{"number":1, "author_association": "MEMBER","labels": [{"url": "u", "name": "n", "color": "c"}]}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.Get returned error: %v", err)
	}

	want := &Issue{
		Number:            Int(1),
		AuthorAssociation: String("MEMBER"),
		Labels: []*Label{{
			URL:   String("u"),
			Name:  String("n"),
			Color: String("c"),
		}},
	}
	if !cmp.Equal(issue, want) {
		t.Errorf("Issues.Get returned %+v, want %+v", issue, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.Get(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.Get(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_Get_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.Get(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssuesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &IssueRequest{
		Title:    String("t"),
		Body:     String("b"),
		Assignee: String("a"),
		Labels:   &[]string{"l1", "l2"},
	}

	mux.HandleFunc("/repos/o/r/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.Create(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Issues.Create returned error: %v", err)
	}

	want := &Issue{Number: Int(1)}
	if !cmp.Equal(issue, want) {
		t.Errorf("Issues.Create returned %+v, want %+v", issue, want)
	}

	const methodName = "Create"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.Create(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.Create(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_Create_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.Create(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestIssuesService_Edit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &IssueRequest{Title: String("t")}

	mux.HandleFunc("/repos/o/r/issues/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.Edit(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.Edit returned error: %v", err)
	}

	want := &Issue{Number: Int(1)}
	if !cmp.Equal(issue, want) {
		t.Errorf("Issues.Edit returned %+v, want %+v", issue, want)
	}

	const methodName = "Edit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.Edit(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.Edit(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_RemoveMilestone(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/issues/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.RemoveMilestone(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.RemoveMilestone returned error: %v", err)
	}

	want := &Issue{Number: Int(1)}
	if !cmp.Equal(issue, want) {
		t.Errorf("Issues.RemoveMilestone returned %+v, want %+v", issue, want)
	}

	const methodName = "RemoveMilestone"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.RemoveMilestone(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.RemoveMilestone(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_Edit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.Edit(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_Lock(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/lock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Issues.Lock(ctx, "o", "r", 1, nil); err != nil {
		t.Errorf("Issues.Lock returned error: %v", err)
	}

	const methodName = "Lock"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Issues.Lock(ctx, "\n", "\n", -1, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Issues.Lock(ctx, "o", "r", 1, nil)
	})
}

func TestIssuesService_LockWithReason(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/lock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	opt := &LockIssueOptions{LockReason: "off-topic"}

	ctx := context.Background()
	if _, err := client.Issues.Lock(ctx, "o", "r", 1, opt); err != nil {
		t.Errorf("Issues.Lock returned error: %v", err)
	}
}

func TestIssuesService_Unlock(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/lock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Issues.Unlock(ctx, "o", "r", 1); err != nil {
		t.Errorf("Issues.Unlock returned error: %v", err)
	}

	const methodName = "Unlock"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Issues.Unlock(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Issues.Unlock(ctx, "o", "r", 1)
	})
}

func TestIsPullRequest(t *testing.T) {
	i := new(Issue)
	if i.IsPullRequest() == true {
		t.Errorf("expected i.IsPullRequest (%v) to return false, got true", i)
	}
	i.PullRequestLinks = &PullRequestLinks{URL: String("http://example.com")}
	if i.IsPullRequest() == false {
		t.Errorf("expected i.IsPullRequest (%v) to return true, got false", i)
	}
}

func TestLockIssueOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &LockIssueOptions{}, "{}")

	u := &LockIssueOptions{
		LockReason: "lr",
	}

	want := `{
		"lock_reason": "lr"
		}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestLinks_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestLinks{}, "{}")

	u := &PullRequestLinks{
		URL:      String("url"),
		HTMLURL:  String("hurl"),
		DiffURL:  String("durl"),
		PatchURL: String("purl"),
	}

	want := `{
		"url": "url",
		"html_url": "hurl",
		"diff_url": "durl",
		"patch_url": "purl"
		}`

	testJSONMarshal(t, u, want)
}

func TestIssueRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssueRequest{}, "{}")

	u := &IssueRequest{
		Title:     String("url"),
		Body:      String("url"),
		Labels:    &[]string{"l"},
		Assignee:  String("url"),
		State:     String("url"),
		Milestone: Int(1),
		Assignees: &[]string{"a"},
	}

	want := `{
		"title": "url",
		"body": "url",
		"labels": [
			"l"
		],
		"assignee": "url",
		"state": "url",
		"milestone": 1,
		"assignees": [
			"a"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestIssue_Marshal(t *testing.T) {
	testJSONMarshal(t, &Issue{}, "{}")

	u := &Issue{
		ID:                Int64(1),
		Number:            Int(1),
		State:             String("s"),
		Locked:            Bool(false),
		Title:             String("title"),
		Body:              String("body"),
		AuthorAssociation: String("aa"),
		User:              &User{ID: Int64(1)},
		Labels:            []*Label{{ID: Int64(1)}},
		Assignee:          &User{ID: Int64(1)},
		Comments:          Int(1),
		ClosedAt:          &referenceTime,
		CreatedAt:         &referenceTime,
		UpdatedAt:         &referenceTime,
		ClosedBy:          &User{ID: Int64(1)},
		URL:               String("url"),
		HTMLURL:           String("hurl"),
		CommentsURL:       String("curl"),
		EventsURL:         String("eurl"),
		LabelsURL:         String("lurl"),
		RepositoryURL:     String("rurl"),
		Milestone:         &Milestone{ID: Int64(1)},
		PullRequestLinks:  &PullRequestLinks{URL: String("url")},
		Repository:        &Repository{ID: Int64(1)},
		Reactions:         &Reactions{TotalCount: Int(1)},
		Assignees:         []*User{{ID: Int64(1)}},
		NodeID:            String("nid"),
		TextMatches:       []*TextMatch{{ObjectURL: String("ourl")}},
		ActiveLockReason:  String("alr"),
	}

	want := `{
		"id": 1,
		"number": 1,
		"state": "s",
		"locked": false,
		"title": "title",
		"body": "body",
		"author_association": "aa",
		"user": {
			"id": 1
		},
		"labels": [
			{
				"id": 1
			}
		],
		"assignee": {
			"id": 1
		},
		"comments": 1,
		"closed_at": ` + referenceTimeStr + `,
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"closed_by": {
			"id": 1
		},
		"url": "url",
		"html_url": "hurl",
		"comments_url": "curl",
		"events_url": "eurl",
		"labels_url": "lurl",
		"repository_url": "rurl",
		"milestone": {
			"id": 1
		},
		"pull_request": {
			"url": "url"
		},
		"repository": {
			"id": 1
		},
		"reactions": {
			"total_count": 1
		},
		"assignees": [
			{
				"id": 1
			}
		],
		"node_id": "nid",
		"text_matches": [
			{
				"object_url": "ourl"
			}
		],
		"active_lock_reason": "alr"
	}`

	testJSONMarshal(t, u, want)
}
