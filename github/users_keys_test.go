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

func TestUsersService_ListKeys_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	keys, _, err := client.Users.ListKeys(ctx, "", opt)
	if err != nil {
		t.Errorf("Users.ListKeys returned error: %v", err)
	}

	want := []*Key{{ID: Int64(1)}}
	if !cmp.Equal(keys, want) {
		t.Errorf("Users.ListKeys returned %+v, want %+v", keys, want)
	}

	const methodName = "ListKeys"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.ListKeys(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListKeys(ctx, "", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_ListKeys_specifiedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	keys, _, err := client.Users.ListKeys(ctx, "u", nil)
	if err != nil {
		t.Errorf("Users.ListKeys returned error: %v", err)
	}

	want := []*Key{{ID: Int64(1)}}
	if !cmp.Equal(keys, want) {
		t.Errorf("Users.ListKeys returned %+v, want %+v", keys, want)
	}
}

func TestUsersService_ListKeys_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Users.ListKeys(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestUsersService_GetKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Users.GetKey(ctx, 1)
	if err != nil {
		t.Errorf("Users.GetKey returned error: %v", err)
	}

	want := &Key{ID: Int64(1)}
	if !cmp.Equal(key, want) {
		t.Errorf("Users.GetKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.GetKey(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetKey(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_CreateKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Key{Key: String("k"), Title: String("t")}

	mux.HandleFunc("/user/keys", func(w http.ResponseWriter, r *http.Request) {
		v := new(Key)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Users.CreateKey(ctx, input)
	if err != nil {
		t.Errorf("Users.GetKey returned error: %v", err)
	}

	want := &Key{ID: Int64(1)}
	if !cmp.Equal(key, want) {
		t.Errorf("Users.GetKey returned %+v, want %+v", key, want)
	}

	const methodName = "CreateKey"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.CreateKey(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_DeleteKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Users.DeleteKey(ctx, 1)
	if err != nil {
		t.Errorf("Users.DeleteKey returned error: %v", err)
	}

	const methodName = "DeleteKey"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.DeleteKey(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DeleteKey(ctx, 1)
	})
}

func TestKey_Marshal(t *testing.T) {
	testJSONMarshal(t, &Key{}, "{}")

	u := &Key{
		ID:        Int64(1),
		Key:       String("abc"),
		URL:       String("url"),
		Title:     String("title"),
		ReadOnly:  Bool(true),
		Verified:  Bool(true),
		CreatedAt: &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"key": "abc",
		"url": "url",
		"title": "title",
		"read_only": true,
		"verified": true,
		"created_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
