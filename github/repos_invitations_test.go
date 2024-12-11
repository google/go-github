// Copyright 2016 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_ListInvitations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	got, _, err := client.Repositories.ListInvitations(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListInvitations returned error: %v", err)
	}

	want := []*RepositoryInvitation{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.ListInvitations = %+v, want %+v", got, want)
	}

	const methodName = "ListInvitations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListInvitations(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListInvitations(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteInvitation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/invitations/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteInvitation(ctx, "o", "r", 2)
	if err != nil {
		t.Errorf("Repositories.DeleteInvitation returned error: %v", err)
	}

	const methodName = "DeleteInvitation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteInvitation(ctx, "\n", "\n", 2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteInvitation(ctx, "o", "r", 2)
	})
}

func TestRepositoriesService_UpdateInvitation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/invitations/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.UpdateInvitation(ctx, "o", "r", 2, "write")
	if err != nil {
		t.Errorf("Repositories.UpdateInvitation returned error: %v", err)
	}

	want := &RepositoryInvitation{ID: Ptr(int64(1))}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.UpdateInvitation = %+v, want %+v", got, want)
	}

	const methodName = "UpdateInvitation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.UpdateInvitation(ctx, "\n", "\n", 2, "write")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateInvitation(ctx, "o", "r", 2, "write")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoryInvitation_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryInvitation{}, "{}")

	r := &RepositoryInvitation{
		ID: Ptr(int64(1)),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			Name: Ptr("n"),
			URL:  Ptr("u"),
		},
		Invitee: &User{
			ID:   Ptr(int64(1)),
			Name: Ptr("n"),
			URL:  Ptr("u"),
		},
		Inviter: &User{
			ID:   Ptr(int64(1)),
			Name: Ptr("n"),
			URL:  Ptr("u"),
		},
		Permissions: Ptr("p"),
		CreatedAt:   &Timestamp{referenceTime},
		URL:         Ptr("u"),
		HTMLURL:     Ptr("h"),
	}

	want := `{
		"id":1,
		"repository":{
			"id":1,
			"name":"n",
			"url":"u"
		},
		"invitee":{
			"id":1,
			"name":"n",
			"url":"u"
		},
		"inviter":{
			"id":1,
			"name":"n",
			"url":"u"
		},
		"permissions":"p",
		"created_at":` + referenceTimeStr + `,
		"url":"u",
		"html_url":"h"
	}`

	testJSONMarshal(t, r, want)
}
