// Copyright 2021 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_ListAutolinks(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/autolinks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprintf(w, `[{"id":1, "key_prefix": "TICKET-", "url_template": "https://example.com/TICKET?query=<num>"}, {"id":2, "key_prefix": "STORY-", "url_template": "https://example.com/STORY?query=<num>"}]`)
	})

	opt := &ListOptions{
		Page: 2,
	}
	ctx := context.Background()
	autolinks, _, err := client.Repositories.ListAutolinks(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListAutolinks returned error: %v", err)
	}

	want := []*Autolink{
		{ID: Ptr(int64(1)), KeyPrefix: Ptr("TICKET-"), URLTemplate: Ptr("https://example.com/TICKET?query=<num>")},
		{ID: Ptr(int64(2)), KeyPrefix: Ptr("STORY-"), URLTemplate: Ptr("https://example.com/STORY?query=<num>")},
	}

	if !cmp.Equal(autolinks, want) {
		t.Errorf("Repositories.ListAutolinks returned %+v, want %+v", autolinks, want)
	}

	const methodName = "ListAutolinks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListAutolinks(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListAutolinks(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_AddAutolink(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	opt := &AutolinkOptions{
		KeyPrefix:      Ptr("TICKET-"),
		URLTemplate:    Ptr("https://example.com/TICKET?query=<num>"),
		IsAlphanumeric: Ptr(true),
	}
	mux.HandleFunc("/repos/o/r/autolinks", func(w http.ResponseWriter, r *http.Request) {
		v := new(AutolinkOptions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))
		testMethod(t, r, "POST")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}
		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`
			{
				"key_prefix": "TICKET-",
				"url_template": "https://example.com/TICKET?query=<num>",
				"is_alphanumeric": true
			}
		`))
	})
	ctx := context.Background()
	autolink, _, err := client.Repositories.AddAutolink(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.AddAutolink returned error: %v", err)
	}
	want := &Autolink{
		KeyPrefix:      Ptr("TICKET-"),
		URLTemplate:    Ptr("https://example.com/TICKET?query=<num>"),
		IsAlphanumeric: Ptr(true),
	}

	if !cmp.Equal(autolink, want) {
		t.Errorf("AddAutolink returned %+v, want %+v", autolink, want)
	}

	const methodName = "AddAutolink"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.AddAutolink(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.AddAutolink(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetAutolink(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/autolinks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"id":1, "key_prefix": "TICKET-", "url_template": "https://example.com/TICKET?query=<num>"}`)
	})

	ctx := context.Background()
	autolink, _, err := client.Repositories.GetAutolink(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetAutolink returned error: %v", err)
	}

	want := &Autolink{ID: Ptr(int64(1)), KeyPrefix: Ptr("TICKET-"), URLTemplate: Ptr("https://example.com/TICKET?query=<num>")}
	if !cmp.Equal(autolink, want) {
		t.Errorf("Repositories.GetAutolink returned %+v, want %+v", autolink, want)
	}

	const methodName = "GetAutolink"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAutolink(ctx, "o", "r", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteAutolink(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/autolinks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteAutolink(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteAutolink returned error: %v", err)
	}

	const methodName = "DeleteAutolink"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteAutolink(ctx, "o", "r", 2)
	})
}

func TestAutolinkOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AutolinkOptions{}, "{}")

	r := &AutolinkOptions{
		KeyPrefix:      Ptr("kp"),
		URLTemplate:    Ptr("URLT"),
		IsAlphanumeric: Ptr(true),
	}

	want := `{
		"key_prefix": "kp",
		"url_template": "URLT",
		"is_alphanumeric": true
	}`

	testJSONMarshal(t, r, want)
}

func TestAutolink_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Autolink{}, "{}")

	r := &Autolink{
		ID:             Ptr(int64(1)),
		KeyPrefix:      Ptr("kp"),
		URLTemplate:    Ptr("URLT"),
		IsAlphanumeric: Ptr(true),
	}

	want := `{
		"id": 1,
		"key_prefix": "kp",
		"url_template": "URLT",
		"is_alphanumeric": true
	}`

	testJSONMarshal(t, r, want)
}
