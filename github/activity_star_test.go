// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestActivityService_ListStargazers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/stargazers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeStarringPreview)
		testFormValues(t, r, values{
			"page": "2",
		})

		fmt.Fprint(w, `[{"starred_at":"2002-02-10T15:30:00Z","user":{"id":1}}]`)
	})

	ctx := context.Background()
	stargazers, _, err := client.Activity.ListStargazers(ctx, "o", "r", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Activity.ListStargazers returned error: %v", err)
	}

	want := []*Stargazer{{StarredAt: &Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}, User: &User{ID: Int64(1)}}}
	if !cmp.Equal(stargazers, want) {
		t.Errorf("Activity.ListStargazers returned %+v, want %+v", stargazers, want)
	}

	const methodName = "ListStargazers"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListStargazers(ctx, "\n", "\n", &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListStargazers(ctx, "o", "r", &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListStarred_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/starred", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join([]string{mediaTypeStarringPreview, mediaTypeTopicsPreview}, ", "))
		fmt.Fprint(w, `[{"starred_at":"2002-02-10T15:30:00Z","repo":{"id":1}}]`)
	})

	ctx := context.Background()
	repos, _, err := client.Activity.ListStarred(ctx, "", nil)
	if err != nil {
		t.Errorf("Activity.ListStarred returned error: %v", err)
	}

	want := []*StarredRepository{{StarredAt: &Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}, Repository: &Repository{ID: Int64(1)}}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Activity.ListStarred returned %+v, want %+v", repos, want)
	}

	const methodName = "ListStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListStarred(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListStarred(ctx, "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListStarred_specifiedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/starred", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join([]string{mediaTypeStarringPreview, mediaTypeTopicsPreview}, ", "))
		testFormValues(t, r, values{
			"sort":      "created",
			"direction": "asc",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"starred_at":"2002-02-10T15:30:00Z","repo":{"id":2}}]`)
	})

	opt := &ActivityListStarredOptions{"created", "asc", ListOptions{Page: 2}}
	ctx := context.Background()
	repos, _, err := client.Activity.ListStarred(ctx, "u", opt)
	if err != nil {
		t.Errorf("Activity.ListStarred returned error: %v", err)
	}

	want := []*StarredRepository{{StarredAt: &Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}, Repository: &Repository{ID: Int64(2)}}}
	if !cmp.Equal(repos, want) {
		t.Errorf("Activity.ListStarred returned %+v, want %+v", repos, want)
	}

	const methodName = "ListStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListStarred(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListStarred(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListStarred_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListStarred(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestActivityService_IsStarred_hasStar(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/starred/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	star, _, err := client.Activity.IsStarred(ctx, "o", "r")
	if err != nil {
		t.Errorf("Activity.IsStarred returned error: %v", err)
	}
	if want := true; star != want {
		t.Errorf("Activity.IsStarred returned %+v, want %+v", star, want)
	}

	const methodName = "IsStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.IsStarred(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.IsStarred(ctx, "o", "r")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_IsStarred_noStar(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/starred/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	star, _, err := client.Activity.IsStarred(ctx, "o", "r")
	if err != nil {
		t.Errorf("Activity.IsStarred returned error: %v", err)
	}
	if want := false; star != want {
		t.Errorf("Activity.IsStarred returned %+v, want %+v", star, want)
	}

	const methodName = "IsStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.IsStarred(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.IsStarred(ctx, "o", "r")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_IsStarred_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.IsStarred(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestActivityService_Star(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/starred/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	ctx := context.Background()
	_, err := client.Activity.Star(ctx, "o", "r")
	if err != nil {
		t.Errorf("Activity.Star returned error: %v", err)
	}

	const methodName = "Star"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Activity.Star(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Activity.Star(ctx, "o", "r")
	})
}

func TestActivityService_Star_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Activity.Star(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestActivityService_Unstar(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/starred/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Activity.Unstar(ctx, "o", "r")
	if err != nil {
		t.Errorf("Activity.Unstar returned error: %v", err)
	}

	const methodName = "Unstar"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Activity.Unstar(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Activity.Unstar(ctx, "o", "r")
	})
}

func TestActivityService_Unstar_invalidID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Activity.Unstar(ctx, "%", "%")
	testURLParseError(t, err)
}
