// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListForks(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTopicsPreview)
		testFormValues(t, r, values{
			"sort": "newest",
			"page": "3",
		})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: ListOptions{Page: 3},
	}
	ctx := context.Background()
	repos, _, err := client.Repositories.ListForks(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListForks returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Repositories.ListForks returned %+v, want %+v", repos, want)
	}

	const methodName = "ListForks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListForks(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListForks(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListForks_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListForks(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateFork(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"organization":"o","name":"n","default_branch_only":true}`+"\n")
		fmt.Fprint(w, `{"id":1}`)
	})

	opt := &RepositoryCreateForkOptions{Organization: "o", Name: "n", DefaultBranchOnly: true}
	ctx := context.Background()
	repo, _, err := client.Repositories.CreateFork(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.CreateFork returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !cmp.Equal(repo, want) {
		t.Errorf("Repositories.CreateFork returned %+v, want %+v", repo, want)
	}

	const methodName = "CreateFork"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateFork(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateFork(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateFork_deferred(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"organization":"o","name":"n","default_branch_only":true}`+"\n")
		// This response indicates the fork will happen asynchronously.
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"id":1}`)
	})

	opt := &RepositoryCreateForkOptions{Organization: "o", Name: "n", DefaultBranchOnly: true}
	ctx := context.Background()
	repo, _, err := client.Repositories.CreateFork(ctx, "o", "r", opt)
	if _, ok := err.(*AcceptedError); !ok {
		t.Errorf("Repositories.CreateFork returned error: %v (want AcceptedError)", err)
	}

	want := &Repository{ID: Int64(1)}
	if !cmp.Equal(repo, want) {
		t.Errorf("Repositories.CreateFork returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_CreateFork_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.CreateFork(ctx, "%", "r", nil)
	testURLParseError(t, err)
}
