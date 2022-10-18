// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListCustomRepoRoles(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom_roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "custom_roles": [{ "id": 1, "name": "Developer", "base_role": "write", "permissions": ["delete_alerts_code_scanning"]}]}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListCustomRepoRoles(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.ListCustomRepoRoles returned error: %v", err)
	}

	want := &OrganizationCustomRepoRoles{TotalCount: Int(1), CustomRepoRoles: []*CustomRepoRoles{{ID: Int64(1), Name: String("Developer"), BaseRole: String("write"), Permissions: []string{"delete_alerts_code_scanning"}}}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListCustomRepoRoles returned %+v, want %+v", apps, want)
	}

	const methodName = "ListCustomRepoRoles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCustomRepoRoles(ctx, "\no")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListCustomRepoRoles(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateCustomRepoRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom_roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":8030,"name":"Labeler","description":"A role for issue and PR labelers","base_role":"read","permissions":["add_label"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateCustomRoleOptions{
		Name:        String("Labeler"),
		Description: String("A role for issue and PR labelers"),
		BaseRole:    String("read"),
		Permissions: []string{"add_label"},
	}
	apps, _, err := client.Organizations.CreateCustomRepoRole(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.CreateCustomRepoRole returned error: %v", err)
	}

	want := &CustomRepoRoles{ID: Int64(8030), Name: String("Labeler"), BaseRole: String("read"), Permissions: []string{"add_label"}, Description: String("A role for issue and PR labelers")}

	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.CreateCustomRepoRole returned %+v, want %+v", apps, want)
	}

	const methodName = "CreateCustomRepoRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.CreateCustomRepoRole(ctx, "\no", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateCustomRepoRole(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateCustomRepoRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom_roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":8030,"name":"Updated Name","description":"Updated Description","base_role":"read","permissions":["add_label"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateCustomRoleOptions{
		Name:        String("Updated Name"),
		Description: String("Updated Description"),
	}
	apps, _, err := client.Organizations.UpdateCustomRepoRole(ctx, "o", "8030", opts)
	if err != nil {
		t.Errorf("Organizations.UpdateCustomRepoRole returned error: %v", err)
	}

	want := &CustomRepoRoles{ID: Int64(8030), Name: String("Updated Name"), BaseRole: String("read"), Permissions: []string{"add_label"}, Description: String("Updated Description")}

	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.UpdateCustomRepoRole returned %+v, want %+v", apps, want)
	}

	const methodName = "UpdateCustomRepoRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.UpdateCustomRepoRole(ctx, "\no", "8030", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateCustomRepoRole(ctx, "o", "8030", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteCustomRepoRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom_roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()

	resp, err := client.Organizations.DeleteCustomRepoRole(ctx, "o", "8030")
	if err != nil {
		t.Errorf("Organizations.DeleteCustomRepoRole returned error: %v", err)
	}

	if !cmp.Equal(resp.StatusCode, 204) {
		t.Errorf("Organizations.DeleteCustomRepoRole returned  status code %+v, want %+v", resp.StatusCode, "204")
	}

	const methodName = "DeleteCustomRepoRole"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DeleteCustomRepoRole(ctx, "\no", "8030")
		return err
	})
}
