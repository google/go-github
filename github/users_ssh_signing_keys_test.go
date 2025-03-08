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

func TestUsersService_ListSSHSigningKeys_authenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/ssh_signing_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	keys, _, err := client.Users.ListSSHSigningKeys(ctx, "", opt)
	if err != nil {
		t.Errorf("Users.ListSSHSigningKeys returned error: %v", err)
	}

	want := []*SSHSigningKey{{ID: Ptr(int64(1))}}
	if !cmp.Equal(keys, want) {
		t.Errorf("Users.ListSSHSigningKeys returned %+v, want %+v", keys, want)
	}

	const methodName = "ListSSHSigningKeys"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.ListSSHSigningKeys(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListSSHSigningKeys(ctx, "", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_ListSSHSigningKeys_specifiedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/ssh_signing_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	keys, _, err := client.Users.ListSSHSigningKeys(ctx, "u", nil)
	if err != nil {
		t.Errorf("Users.ListSSHSigningKeys returned error: %v", err)
	}

	want := []*SSHSigningKey{{ID: Ptr(int64(1))}}
	if !cmp.Equal(keys, want) {
		t.Errorf("Users.ListSSHSigningKeys returned %+v, want %+v", keys, want)
	}
}

func TestUsersService_ListSSHSigningKeys_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Users.ListSSHSigningKeys(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestUsersService_GetSSHSigningKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/ssh_signing_keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Users.GetSSHSigningKey(ctx, 1)
	if err != nil {
		t.Errorf("Users.GetSSHSigningKey returned error: %v", err)
	}

	want := &SSHSigningKey{ID: Ptr(int64(1))}
	if !cmp.Equal(key, want) {
		t.Errorf("Users.GetSSHSigningKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetSSHSigningKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.GetSSHSigningKey(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetSSHSigningKey(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_CreateSSHSigningKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Key{Key: Ptr("k"), Title: Ptr("t")}

	mux.HandleFunc("/user/ssh_signing_keys", func(w http.ResponseWriter, r *http.Request) {
		v := new(Key)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	key, _, err := client.Users.CreateSSHSigningKey(ctx, input)
	if err != nil {
		t.Errorf("Users.CreateSSHSigningKey returned error: %v", err)
	}

	want := &SSHSigningKey{ID: Ptr(int64(1))}
	if !cmp.Equal(key, want) {
		t.Errorf("Users.CreateSSHSigningKey returned %+v, want %+v", key, want)
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

func TestUsersService_DeleteSSHSigningKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/ssh_signing_keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Users.DeleteSSHSigningKey(ctx, 1)
	if err != nil {
		t.Errorf("Users.DeleteSSHSigningKey returned error: %v", err)
	}

	const methodName = "DeleteSSHSigningKey"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.DeleteSSHSigningKey(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DeleteSSHSigningKey(ctx, 1)
	})
}

func TestSSHSigningKey_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SSHSigningKey{}, "{}")

	u := &Key{
		ID:        Ptr(int64(1)),
		Key:       Ptr("abc"),
		Title:     Ptr("title"),
		CreatedAt: &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"key": "abc",
		"title": "title",
		"created_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
