// Copyright 2017 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListOutsideCollaborators(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/outside_collaborators", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter": "2fa_disabled",
			"page":   "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOutsideCollaboratorsOptions{
		Filter:      "2fa_disabled",
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	members, _, err := client.Organizations.ListOutsideCollaborators(ctx, "o", opt)
	if err != nil {
		t.Errorf("Organizations.ListOutsideCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Ptr(int64(1))}}
	if !cmp.Equal(members, want) {
		t.Errorf("Organizations.ListOutsideCollaborators returned %+v, want %+v", members, want)
	}

	const methodName = "ListOutsideCollaborators"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListOutsideCollaborators(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListOutsideCollaborators(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_ListOutsideCollaborators_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.ListOutsideCollaborators(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_RemoveOutsideCollaborator(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	}
	mux.HandleFunc("/orgs/o/outside_collaborators/u", handler)

	ctx := context.Background()
	_, err := client.Organizations.RemoveOutsideCollaborator(ctx, "o", "u")
	if err != nil {
		t.Errorf("Organizations.RemoveOutsideCollaborator returned error: %v", err)
	}

	const methodName = "RemoveOutsideCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.RemoveOutsideCollaborator(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.RemoveOutsideCollaborator(ctx, "o", "u")
	})
}

func TestOrganizationsService_RemoveOutsideCollaborator_NonMember(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	}
	mux.HandleFunc("/orgs/o/outside_collaborators/u", handler)

	ctx := context.Background()
	_, err := client.Organizations.RemoveOutsideCollaborator(ctx, "o", "u")
	if err, ok := err.(*ErrorResponse); !ok {
		t.Errorf("Organizations.RemoveOutsideCollaborator did not return an error")
	} else if err.Response.StatusCode != http.StatusNotFound {
		t.Errorf("Organizations.RemoveOutsideCollaborator did not return 404 status code")
	}
}

func TestOrganizationsService_RemoveOutsideCollaborator_Member(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	mux.HandleFunc("/orgs/o/outside_collaborators/u", handler)

	ctx := context.Background()
	_, err := client.Organizations.RemoveOutsideCollaborator(ctx, "o", "u")
	if err, ok := err.(*ErrorResponse); !ok {
		t.Errorf("Organizations.RemoveOutsideCollaborator did not return an error")
	} else if err.Response.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Organizations.RemoveOutsideCollaborator did not return 422 status code")
	}
}

func TestOrganizationsService_ConvertMemberToOutsideCollaborator(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	}
	mux.HandleFunc("/orgs/o/outside_collaborators/u", handler)

	ctx := context.Background()
	_, err := client.Organizations.ConvertMemberToOutsideCollaborator(ctx, "o", "u")
	if err != nil {
		t.Errorf("Organizations.ConvertMemberToOutsideCollaborator returned error: %v", err)
	}

	const methodName = "ConvertMemberToOutsideCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.ConvertMemberToOutsideCollaborator(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.ConvertMemberToOutsideCollaborator(ctx, "o", "u")
	})
}

func TestOrganizationsService_ConvertMemberToOutsideCollaborator_NonMemberOrLastOwner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusForbidden)
	}
	mux.HandleFunc("/orgs/o/outside_collaborators/u", handler)

	ctx := context.Background()
	_, err := client.Organizations.ConvertMemberToOutsideCollaborator(ctx, "o", "u")
	if err, ok := err.(*ErrorResponse); !ok {
		t.Errorf("Organizations.ConvertMemberToOutsideCollaborator did not return an error")
	} else if err.Response.StatusCode != http.StatusForbidden {
		t.Errorf("Organizations.ConvertMemberToOutsideCollaborator did not return 403 status code")
	}
}
