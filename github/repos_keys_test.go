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

func TestRepositoriesService_ListKeys(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	keys, _, err := client.Repositories.ListKeys(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListKeys returned error: %v", err)
	}

	want := []*Key{{ID: Ptr(int64(1))}}
	if !cmp.Equal(keys, want) {
		t.Errorf("Repositories.ListKeys returned %+v, want %+v", keys, want)
	}

	const methodName = "ListKeys"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListKeys(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListKeys(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListKeys_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.ListKeys(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Repositories.GetKey(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetKey returned error: %v", err)
	}

	want := &Key{ID: Ptr(int64(1))}
	if !cmp.Equal(key, want) {
		t.Errorf("Repositories.GetKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetKey(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetKey(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetKey_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.GetKey(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Key{Key: Ptr("k"), Title: Ptr("t")}

	mux.HandleFunc("/repos/o/r/keys", func(w http.ResponseWriter, r *http.Request) {
		v := new(Key)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Repositories.CreateKey(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.GetKey returned error: %v", err)
	}

	want := &Key{ID: Ptr(int64(1))}
	if !cmp.Equal(key, want) {
		t.Errorf("Repositories.GetKey returned %+v, want %+v", key, want)
	}

	const methodName = "CreateKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateKey(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateKey(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateKey_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.CreateKey(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_DeleteKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteKey(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteKey returned error: %v", err)
	}

	const methodName = "DeleteKey"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteKey(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteKey(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_DeleteKey_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Repositories.DeleteKey(ctx, "%", "%", 1)
	testURLParseError(t, err)
}
