// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testPullRequestStackResponse() string {
	return `{
		"id":1,
		"number":42,
		"node_id":"S_kwDOABCDEF4AAAAA",
		"url":"https://api.github.com/repos/o/r/stacks/42",
		"base":{"ref":"main"},
		"open":true,
		"created_at":` + referenceTimeStr + `,
		"pull_requests":[{
			"number":101,
			"state":"open",
			"draft":false,
			"merged_at":null,
			"head":{"ref":"feature","sha":"abc123"}
		}]
	}`
}

func testPullRequestStack() *PullRequestStack {
	return &PullRequestStack{
		ID:        Ptr(int64(1)),
		Number:    Ptr(42),
		NodeID:    Ptr("S_kwDOABCDEF4AAAAA"),
		URL:       Ptr("https://api.github.com/repos/o/r/stacks/42"),
		Base:      &PullRequestStackBase{Ref: "main"},
		Open:      Ptr(true),
		CreatedAt: &referenceTimestamp,
		PullRequests: []*PullRequest{{
			Number: Ptr(101),
			State:  Ptr("open"),
			Draft:  Ptr(false),
			Head: &PullRequestBranch{
				Ref: Ptr("feature"),
				SHA: Ptr("abc123"),
			},
		}},
	}
}

func TestPullRequestsService_ListStacks(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stacks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"pull_request": "101",
			"page":         "2",
			"per_page":     "50",
		})
		fmt.Fprintf(w, "[%s]", testPullRequestStackResponse())
	})

	opts := &PullRequestListStacksOptions{
		PullRequest: 101,
		ListOptions: ListOptions{Page: 2, PerPage: 50},
	}
	ctx := t.Context()
	stacks, _, err := client.PullRequests.ListStacks(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("PullRequests.ListStacks returned error: %v", err)
	}

	want := []*PullRequestStack{testPullRequestStack()}
	if !cmp.Equal(stacks, want) {
		t.Errorf("PullRequests.ListStacks returned %+v, want %+v", stacks, want)
	}

	const methodName = "ListStacks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListStacks(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListStacks(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListStacks_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.PullRequests.ListStacks(t.Context(), "%", "%", nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateStack(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	input := CreatePullRequestStackRequest{PullRequests: []int{101, 102}}

	mux.HandleFunc("/repos/o/r/stacks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, testPullRequestStackResponse())
	})

	ctx := t.Context()
	stack, _, err := client.PullRequests.CreateStack(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("PullRequests.CreateStack returned error: %v", err)
	}
	if want := testPullRequestStack(); !cmp.Equal(stack, want) {
		t.Errorf("PullRequests.CreateStack returned %+v, want %+v", stack, want)
	}

	const methodName = "CreateStack"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.CreateStack(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.CreateStack(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_CreateStack_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.PullRequests.CreateStack(t.Context(), "%", "%", CreatePullRequestStackRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_GetStack(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stacks/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, testPullRequestStackResponse())
	})

	ctx := t.Context()
	stack, _, err := client.PullRequests.GetStack(ctx, "o", "r", 42)
	if err != nil {
		t.Errorf("PullRequests.GetStack returned error: %v", err)
	}
	if want := testPullRequestStack(); !cmp.Equal(stack, want) {
		t.Errorf("PullRequests.GetStack returned %+v, want %+v", stack, want)
	}

	const methodName = "GetStack"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.GetStack(ctx, "\n", "\n", 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.GetStack(ctx, "o", "r", 42)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_GetStack_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.PullRequests.GetStack(t.Context(), "%", "%", 42)
	testURLParseError(t, err)
}

func TestPullRequestsService_AddToStack(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	input := AddPullRequestsToStackRequest{PullRequests: []int{103, 104}}

	mux.HandleFunc("/repos/o/r/stacks/42/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		fmt.Fprint(w, testPullRequestStackResponse())
	})

	ctx := t.Context()
	stack, _, err := client.PullRequests.AddToStack(ctx, "o", "r", 42, input)
	if err != nil {
		t.Errorf("PullRequests.AddToStack returned error: %v", err)
	}
	if want := testPullRequestStack(); !cmp.Equal(stack, want) {
		t.Errorf("PullRequests.AddToStack returned %+v, want %+v", stack, want)
	}

	const methodName = "AddToStack"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.AddToStack(ctx, "\n", "\n", 42, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.AddToStack(ctx, "o", "r", 42, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_AddToStack_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.PullRequests.AddToStack(t.Context(), "%", "%", 42, AddPullRequestsToStackRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_Unstack(t *testing.T) {
	t.Parallel()

	t.Run("returns updated stack", func(t *testing.T) {
		t.Parallel()
		client, mux, _ := setup(t)
		mux.HandleFunc("/repos/o/r/stacks/42/unstack", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, testPullRequestStackResponse())
		})

		stack, _, err := client.PullRequests.Unstack(t.Context(), "o", "r", 42)
		if err != nil {
			t.Errorf("PullRequests.Unstack returned error: %v", err)
		}
		if want := testPullRequestStack(); !cmp.Equal(stack, want) {
			t.Errorf("PullRequests.Unstack returned %+v, want %+v", stack, want)
		}
	})

	t.Run("returns nil when stack is dissolved", func(t *testing.T) {
		t.Parallel()
		client, mux, _ := setup(t)
		mux.HandleFunc("/repos/o/r/stacks/42/unstack", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			w.WriteHeader(http.StatusNoContent)
		})

		stack, resp, err := client.PullRequests.Unstack(t.Context(), "o", "r", 42)
		if err != nil {
			t.Errorf("PullRequests.Unstack returned error: %v", err)
		}
		if stack != nil {
			t.Errorf("PullRequests.Unstack returned %+v, want nil", stack)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("PullRequests.Unstack returned status %v, want %v", resp.StatusCode, http.StatusNoContent)
		}
	})

	client, _, _ := setup(t)
	ctx := t.Context()
	const methodName = "Unstack"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.Unstack(ctx, "\n", "\n", 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.Unstack(ctx, "o", "r", 42)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_Unstack_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.PullRequests.Unstack(t.Context(), "%", "%", 42)
	testURLParseError(t, err)
}

func TestPullRequestStack_unmarshal(t *testing.T) {
	t.Parallel()
	testJSONUnmarshalOnly(t, testPullRequestStack(), testPullRequestStackResponse())
}

func TestCreatePullRequestStackRequest_marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CreatePullRequestStackRequest{PullRequests: []int{101, 102}}, `{"pull_requests":[101,102]}`)
}

func TestAddPullRequestsToStackRequest_marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AddPullRequestsToStackRequest{PullRequests: []int{103}}, `{"pull_requests":[103]}`)
}
