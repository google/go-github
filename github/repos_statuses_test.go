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

func TestRepositoriesService_ListStatuses(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/commits/r/statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	statuses, _, err := client.Repositories.ListStatuses(ctx, "o", "r", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListStatuses returned error: %v", err)
	}

	want := []*RepoStatus{{ID: Ptr(int64(1))}}
	if !cmp.Equal(statuses, want) {
		t.Errorf("Repositories.ListStatuses returned %+v, want %+v", statuses, want)
	}

	const methodName = "ListStatuses"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListStatuses(ctx, "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListStatuses(ctx, "o", "r", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListStatuses_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListStatuses(ctx, "%", "r", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &RepoStatus{State: Ptr("s"), TargetURL: Ptr("t"), Description: Ptr("d")}

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepoStatus)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	status, _, err := client.Repositories.CreateStatus(ctx, "o", "r", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateStatus returned error: %v", err)
	}

	want := &RepoStatus{ID: Ptr(int64(1))}
	if !cmp.Equal(status, want) {
		t.Errorf("Repositories.CreateStatus returned %+v, want %+v", status, want)
	}

	const methodName = "CreateStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateStatus(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateStatus(ctx, "o", "r", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateStatus_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.CreateStatus(ctx, "%", "r", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetCombinedStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/commits/r/status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `{"state":"success", "statuses":[{"id":1}]}`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	status, _, err := client.Repositories.GetCombinedStatus(ctx, "o", "r", "r", opt)
	if err != nil {
		t.Errorf("Repositories.GetCombinedStatus returned error: %v", err)
	}

	want := &CombinedStatus{State: Ptr("success"), Statuses: []*RepoStatus{{ID: Ptr(int64(1))}}}
	if !cmp.Equal(status, want) {
		t.Errorf("Repositories.GetCombinedStatus returned %+v, want %+v", status, want)
	}

	const methodName = "GetCombinedStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCombinedStatus(ctx, "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCombinedStatus(ctx, "o", "r", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepoStatus_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepoStatus{}, "{}")

	u := &RepoStatus{
		ID:          Ptr(int64(1)),
		NodeID:      Ptr("nid"),
		URL:         Ptr("url"),
		State:       Ptr("state"),
		TargetURL:   Ptr("turl"),
		Description: Ptr("desc"),
		Context:     Ptr("ctx"),
		AvatarURL:   Ptr("aurl"),
		Creator:     &User{ID: Ptr(int64(1))},
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"url": "url",
		"state": "state",
		"target_url": "turl",
		"description": "desc",
		"context": "ctx",
		"avatar_url": "aurl",
		"creator": {
			"id": 1
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestCombinedStatus_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CombinedStatus{}, "{}")

	u := &CombinedStatus{
		State:      Ptr("state"),
		Name:       Ptr("name"),
		SHA:        Ptr("sha"),
		TotalCount: Ptr(1),
		Statuses: []*RepoStatus{
			{
				ID:          Ptr(int64(1)),
				NodeID:      Ptr("nid"),
				URL:         Ptr("url"),
				State:       Ptr("state"),
				TargetURL:   Ptr("turl"),
				Description: Ptr("desc"),
				Context:     Ptr("ctx"),
				AvatarURL:   Ptr("aurl"),
				Creator:     &User{ID: Ptr(int64(1))},
				CreatedAt:   &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
			},
		},
		CommitURL:     Ptr("curl"),
		RepositoryURL: Ptr("rurl"),
	}

	want := `{
		"state": "state",
		"name": "name",
		"sha": "sha",
		"total_count": 1,
		"statuses": [
			{
				"id": 1,
				"node_id": "nid",
				"url": "url",
				"state": "state",
				"target_url": "turl",
				"description": "desc",
				"context": "ctx",
				"avatar_url": "aurl",
				"creator": {
					"id": 1
				},
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `
			}
		],
		"commit_url": "curl",
		"repository_url": "rurl"
	}`

	testJSONMarshal(t, u, want)
}
