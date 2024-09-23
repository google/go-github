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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListRoles(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "roles": [
			{
				"id": 1,
				"name": "Auditor",
				"permissions": ["read_audit_logs"],
				"organization": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"html_url": "h",
					"name": "n",
					"company": "c",
					"blog": "b",
					"location": "l",
					"email": "e"
				},
				"created_at": "2024-07-21T19:33:08Z",
				"updated_at": "2024-07-21T19:33:08Z",
				"source": "Organization",
				"base_role": "admin"
			}
			]
		}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListRoles(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.ListRoles returned error: %v", err)
	}

	want := &OrganizationCustomRoles{
		TotalCount: Int(1),
		CustomRepoRoles: []*CustomOrgRoles{
			{
				ID:          Int64(1),
				Name:        String("Auditor"),
				Permissions: []string{"read_audit_logs"},
				Org: &Organization{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					AvatarURL: String("a"),
					HTMLURL:   String("h"),
					Name:      String("n"),
					Company:   String("c"),
					Blog:      String("b"),
					Location:  String("l"),
					Email:     String("e"),
				},
				CreatedAt: &Timestamp{time.Date(2024, time.July, 21, 19, 33, 8, 0, time.UTC)},
				UpdatedAt: &Timestamp{time.Date(2024, time.July, 21, 19, 33, 8, 0, time.UTC)},
				Source:    String("Organization"),
				BaseRole:  String("admin"),
			},
		},
	}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListRoles returned %+v, want %+v", apps, want)
	}

	const methodName = "ListRoles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListRoles(ctx, "\no")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListRoles(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	// Test built-in org role
	mux.HandleFunc("/orgs/o/organization-roles/8132", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 8132,
			"name": "all_repo_read",
			"description": "Grants read access to all repositories in the organization.",
			"permissions": [],
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`,
			"source": "Predefined",
			"base_role": "read"
		}`)
	})

	ctx := context.Background()

	gotBuiltInRole, _, err := client.Organizations.GetOrgRole(ctx, "o", 8132)
	if err != nil {
		t.Errorf("Organizations.GetOrgRole returned error: %v", err)
	}

	wantBuiltInRole := &CustomOrgRoles{
		ID:          Int64(8132),
		Name:        String("all_repo_read"),
		Description: String("Grants read access to all repositories in the organization."),
		Permissions: []string{},
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
		Source:      String("Predefined"),
		BaseRole:    String("read"),
	}

	if !cmp.Equal(gotBuiltInRole, wantBuiltInRole) {
		t.Errorf("Organizations.GetOrgRole returned %+v, want %+v", gotBuiltInRole, wantBuiltInRole)
	}

	// Test custom org role
	mux.HandleFunc("/orgs/o/organization-roles/123456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 123456,
			"name": "test-role",
			"description": "test-role",
			"permissions": [
				"read_organization_custom_org_role",
				"read_organization_custom_repo_role",
				"write_organization_custom_org_role"
			],
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`,
			"source": "Organization",
			"base_role": null
			}`)
	})

	gotCustomRole, _, err := client.Organizations.GetOrgRole(ctx, "o", 123456)
	if err != nil {
		t.Errorf("Organizations.GetOrgRole returned error: %v", err)
	}

	wantCustomRole := &CustomOrgRoles{
		ID:          Int64(123456),
		Name:        String("test-role"),
		Description: String("test-role"),
		Permissions: []string{
			"read_organization_custom_org_role",
			"read_organization_custom_repo_role",
			"write_organization_custom_org_role",
		},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		Source:    String("Organization"),
		BaseRole:  nil,
	}

	if !cmp.Equal(gotCustomRole, wantCustomRole) {
		t.Errorf("Organizations.GetOrgRole returned %+v, want %+v", gotCustomRole, wantCustomRole)
	}

	const methodName = "GetOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetOrgRole(ctx, "\no", -8132)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrgRole(ctx, "o", 8132)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateCustomOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":8030,"name":"Reader","description":"A role for reading custom org roles","permissions":["read_organization_custom_org_role"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateOrgRoleOptions{
		Name:        String("Reader"),
		Description: String("A role for reading custom org roles"),
		Permissions: []string{"read_organization_custom_org_role"},
	}
	gotRoles, _, err := client.Organizations.CreateCustomOrgRole(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.CreateCustomOrgRole returned error: %v", err)
	}

	want := &CustomOrgRoles{ID: Int64(8030), Name: String("Reader"), Permissions: []string{"read_organization_custom_org_role"}, Description: String("A role for reading custom org roles")}

	if !cmp.Equal(gotRoles, want) {
		t.Errorf("Organizations.CreateCustomOrgRole returned %+v, want %+v", gotRoles, want)
	}

	const methodName = "CreateCustomOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.CreateCustomOrgRole(ctx, "\no", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateCustomOrgRole(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateCustomOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":8030,"name":"Updated Name","description":"Updated Description","permissions":["read_organization_custom_org_role"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateOrgRoleOptions{
		Name:        String("Updated Name"),
		Description: String("Updated Description"),
	}
	gotRoles, _, err := client.Organizations.UpdateCustomOrgRole(ctx, "o", 8030, opts)
	if err != nil {
		t.Errorf("Organizations.UpdateCustomOrgRole returned error: %v", err)
	}

	want := &CustomOrgRoles{ID: Int64(8030), Name: String("Updated Name"), Permissions: []string{"read_organization_custom_org_role"}, Description: String("Updated Description")}

	if !cmp.Equal(gotRoles, want) {
		t.Errorf("Organizations.UpdateCustomOrgRole returned %+v, want %+v", gotRoles, want)
	}

	const methodName = "UpdateCustomOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.UpdateCustomOrgRole(ctx, "\no", 8030, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateCustomOrgRole(ctx, "o", 8030, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteCustomOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()

	resp, err := client.Organizations.DeleteCustomOrgRole(ctx, "o", 8030)
	if err != nil {
		t.Errorf("Organizations.DeleteCustomOrgRole returned error: %v", err)
	}

	if !cmp.Equal(resp.StatusCode, 204) {
		t.Errorf("Organizations.DeleteCustomOrgRole returned  status code %+v, want %+v", resp.StatusCode, "204")
	}

	const methodName = "DeleteCustomOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DeleteCustomOrgRole(ctx, "\no", 8030)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.DeleteCustomOrgRole(ctx, "o", 8030)
	})
}

func TestOrganizationsService_ListCustomRepoRoles(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom-repository-roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "custom_roles": [
			{
				"id": 1,
				"name": "Developer",
				"base_role": "write",
				"permissions": ["delete_alerts_code_scanning"],
				"organization": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"html_url": "h",
					"name": "n",
					"company": "c",
					"blog": "b",
					"location": "l",
					"email": "e"
				},
				"created_at": "2024-07-21T19:33:08Z",
				"updated_at": "2024-07-21T19:33:08Z"
			}
		  ]
		}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListCustomRepoRoles(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.ListCustomRepoRoles returned error: %v", err)
	}

	want := &OrganizationCustomRepoRoles{
		TotalCount: Int(1),
		CustomRepoRoles: []*CustomRepoRoles{
			{
				ID:          Int64(1),
				Name:        String("Developer"),
				BaseRole:    String("write"),
				Permissions: []string{"delete_alerts_code_scanning"},
				Org: &Organization{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					AvatarURL: String("a"),
					HTMLURL:   String("h"),
					Name:      String("n"),
					Company:   String("c"),
					Blog:      String("b"),
					Location:  String("l"),
					Email:     String("e"),
				},
				CreatedAt: &Timestamp{time.Date(2024, time.July, 21, 19, 33, 8, 0, time.UTC)},
				UpdatedAt: &Timestamp{time.Date(2024, time.July, 21, 19, 33, 8, 0, time.UTC)},
			},
		},
	}
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

	mux.HandleFunc("/orgs/o/custom-repository-roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":8030,"name":"Labeler","description":"A role for issue and PR labelers","base_role":"read","permissions":["add_label"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateCustomRepoRoleOptions{
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

	mux.HandleFunc("/orgs/o/custom-repository-roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":8030,"name":"Updated Name","description":"Updated Description","base_role":"read","permissions":["add_label"]}`)
	})

	ctx := context.Background()

	opts := &CreateOrUpdateCustomRepoRoleOptions{
		Name:        String("Updated Name"),
		Description: String("Updated Description"),
	}
	apps, _, err := client.Organizations.UpdateCustomRepoRole(ctx, "o", 8030, opts)
	if err != nil {
		t.Errorf("Organizations.UpdateCustomRepoRole returned error: %v", err)
	}

	want := &CustomRepoRoles{ID: Int64(8030), Name: String("Updated Name"), BaseRole: String("read"), Permissions: []string{"add_label"}, Description: String("Updated Description")}

	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.UpdateCustomRepoRole returned %+v, want %+v", apps, want)
	}

	const methodName = "UpdateCustomRepoRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.UpdateCustomRepoRole(ctx, "\no", 8030, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateCustomRepoRole(ctx, "o", 8030, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteCustomRepoRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom-repository-roles/8030", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()

	resp, err := client.Organizations.DeleteCustomRepoRole(ctx, "o", 8030)
	if err != nil {
		t.Errorf("Organizations.DeleteCustomRepoRole returned error: %v", err)
	}

	if !cmp.Equal(resp.StatusCode, 204) {
		t.Errorf("Organizations.DeleteCustomRepoRole returned  status code %+v, want %+v", resp.StatusCode, "204")
	}

	const methodName = "DeleteCustomRepoRole"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DeleteCustomRepoRole(ctx, "\no", 8030)
		return err
	})
}

func TestOrganizationsService_ListTeamsAssignedToOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles/1729/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})
	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	apps, _, err := client.Organizations.ListTeamsAssignedToOrgRole(ctx, "o", 1729, opt)
	if err != nil {
		t.Errorf("Organizations.ListTeamsAssignedToOrgRole returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListTeamsAssignedToOrgRole returned %+v, want %+v", apps, want)
	}

	const methodName = "ListTeamsAssignedToOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListTeamsAssignedToOrgRole(ctx, "\no", 1729, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListTeamsAssignedToOrgRole(ctx, "o", 1729, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_ListUsersAssignedToOrgRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/organization-roles/1729/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})
	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	apps, _, err := client.Organizations.ListUsersAssignedToOrgRole(ctx, "o", 1729, opt)
	if err != nil {
		t.Errorf("Organizations.ListUsersAssignedToOrgRole returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListUsersAssignedToOrgRole returned %+v, want %+v", apps, want)
	}

	const methodName = "ListUsersAssignedToOrgRole"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListUsersAssignedToOrgRole(ctx, "\no", 1729, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListUsersAssignedToOrgRole(ctx, "o", 1729, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
