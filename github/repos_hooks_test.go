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

func TestRepositoriesService_CreateHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Hook{CreatedAt: &referenceTime}

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(createHookRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &createHookRequest{Name: "web"}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Repositories.CreateHook(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !cmp.Equal(hook, want) {
		t.Errorf("Repositories.CreateHook returned %+v, want %+v", hook, want)
	}

	const methodName = "CreateHook"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateHook(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateHook(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListHooks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}

	ctx := context.Background()
	hooks, _, err := client.Repositories.ListHooks(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListHooks returned error: %v", err)
	}

	want := []*Hook{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(hooks, want) {
		t.Errorf("Repositories.ListHooks returned %+v, want %+v", hooks, want)
	}

	const methodName = "ListHooks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListHooks(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListHooks(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListHooks_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListHooks(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Repositories.GetHook(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !cmp.Equal(hook, want) {
		t.Errorf("Repositories.GetHook returned %+v, want %+v", hook, want)
	}

	const methodName = "GetHook"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetHook(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetHook(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetHook_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.GetHook(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoriesService_EditHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Hook{}

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Repositories.EditHook(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !cmp.Equal(hook, want) {
		t.Errorf("Repositories.EditHook returned %+v, want %+v", hook, want)
	}

	const methodName = "EditHook"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditHook(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditHook(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EditHook_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.EditHook(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_DeleteHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteHook(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteHook returned error: %v", err)
	}

	const methodName = "DeleteHook"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteHook(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteHook(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_DeleteHook_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Repositories.DeleteHook(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoriesService_PingHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/pings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Repositories.PingHook(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.PingHook returned error: %v", err)
	}

	const methodName = "PingHook"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.PingHook(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.PingHook(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_TestHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Repositories.TestHook(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.TestHook returned error: %v", err)
	}

	const methodName = "TestHook"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.TestHook(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.TestHook(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_TestHook_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Repositories.TestHook(ctx, "%", "%", 1)
	testURLParseError(t, err)
}
