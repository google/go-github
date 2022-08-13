// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_ListTagProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/tags/protection", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `[{"id":1, "pattern":"tag1"},{"id":2, "pattern":"tag2"}]`)
	})

	ctx := context.Background()
	tagProtections, _, err := client.Repositories.ListTagProtection(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.ListTagProtection returned error: %v", err)
	}

	want := []*TagProtection{{ID: Int64(1), Pattern: String("tag1")}, {ID: Int64(2), Pattern: String("tag2")}}
	if !cmp.Equal(tagProtections, want) {
		t.Errorf("Repositories.ListTagProtection returned %+v, want %+v", tagProtections, want)
	}

	const methodName = "ListTagProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListTagProtection(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListTagProtection(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListTagProtection_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListTagProtection(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateTagProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	pattern := "tag*"

	mux.HandleFunc("/repos/o/r/tags/protection", func(w http.ResponseWriter, r *http.Request) {
		v := new(tagProtectionRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &tagProtectionRequest{Pattern: "tag*"}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1,"pattern":"tag*"}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.CreateTagProtection(ctx, "o", "r", pattern)
	if err != nil {
		t.Errorf("Repositories.CreateTagProtection returned error: %v", err)
	}

	want := &TagProtection{ID: Int64(1), Pattern: String("tag*")}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.CreateTagProtection returned %+v, want %+v", got, want)
	}

	const methodName = "CreateTagProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateTagProtection(ctx, "\n", "\n", pattern)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateTagProtection(ctx, "o", "r", pattern)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteTagProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/tags/protection/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteTagProtection(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteTagProtection returned error: %v", err)
	}

	const methodName = "DeleteTagProtection"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteTagProtection(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteTagProtection(ctx, "o", "r", 1)
	})
}
