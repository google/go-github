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

func TestOrganization_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Organization{}, "{}")

	o := &Organization{
		BillingEmail:                         Ptr("support@github.com"),
		Blog:                                 Ptr("https://github.com/blog"),
		Company:                              Ptr("GitHub"),
		Email:                                Ptr("support@github.com"),
		TwitterUsername:                      Ptr("github"),
		Location:                             Ptr("San Francisco"),
		Name:                                 Ptr("github"),
		Description:                          Ptr("GitHub, the company."),
		IsVerified:                           Ptr(true),
		HasOrganizationProjects:              Ptr(true),
		HasRepositoryProjects:                Ptr(true),
		DefaultRepoPermission:                Ptr("read"),
		MembersCanCreateRepos:                Ptr(true),
		MembersCanCreateInternalRepos:        Ptr(true),
		MembersCanCreatePrivateRepos:         Ptr(true),
		MembersCanCreatePublicRepos:          Ptr(false),
		MembersAllowedRepositoryCreationType: Ptr("all"),
		MembersCanCreatePages:                Ptr(true),
		MembersCanCreatePublicPages:          Ptr(false),
		MembersCanCreatePrivatePages:         Ptr(true),
	}
	want := `
		{
			"billing_email": "support@github.com",
			"blog": "https://github.com/blog",
			"company": "GitHub",
			"email": "support@github.com",
			"twitter_username": "github",
			"location": "San Francisco",
			"name": "github",
			"description": "GitHub, the company.",
			"is_verified": true,
			"has_organization_projects": true,
			"has_repository_projects": true,
			"default_repository_permission": "read",
			"members_can_create_repositories": true,
			"members_can_create_public_repositories": false,
			"members_can_create_private_repositories": true,
			"members_can_create_internal_repositories": true,
			"members_allowed_repository_creation_type": "all",
			"members_can_create_pages": true,
			"members_can_create_public_pages": false,
			"members_can_create_private_pages": true
		}
	`
	testJSONMarshal(t, o, want)
}

func TestOrganizationsService_ListAll(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	since := int64(1342004)
	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1342004"})
		fmt.Fprint(w, `[{"id":4314092}]`)
	})

	opt := &OrganizationsListOptions{Since: since}
	ctx := context.Background()
	orgs, _, err := client.Organizations.ListAll(ctx, opt)
	if err != nil {
		t.Errorf("Organizations.ListAll returned error: %v", err)
	}

	want := []*Organization{{ID: Ptr(int64(4314092))}}
	if !cmp.Equal(orgs, want) {
		t.Errorf("Organizations.ListAll returned %+v, want %+v", orgs, want)
	}

	const methodName = "ListAll"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListAll(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_List_authenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	orgs, _, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []*Organization{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.List(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.List(ctx, "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_List_specifiedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	orgs, _, err := client.Organizations.List(ctx, "u", opt)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []*Organization{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.List(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.List(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_List_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.List(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_Get(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMemberAllowedRepoCreationTypePreview)
		fmt.Fprint(w, `{"id":1, "login":"l", "url":"u", "avatar_url": "a", "location":"l"}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.Get(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	want := &Organization{ID: Ptr(int64(1)), Login: Ptr("l"), URL: Ptr("u"), AvatarURL: Ptr("a"), Location: Ptr("l")}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.Get returned %+v, want %+v", org, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.Get(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_Get_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.Get(ctx, "%")
	testURLParseError(t, err)
}

func TestOrganizationsService_GetByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "login":"l", "url":"u", "avatar_url": "a", "location":"l"}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("Organizations.GetByID returned error: %v", err)
	}

	want := &Organization{ID: Ptr(int64(1)), Login: Ptr("l"), URL: Ptr("u"), AvatarURL: Ptr("a"), Location: Ptr("l")}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.GetByID returned %+v, want %+v", org, want)
	}

	const methodName = "GetByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetByID(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetByID(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_Edit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Organization{Login: Ptr("l")}

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		v := new(Organization)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testHeader(t, r, "Accept", mediaTypeMemberAllowedRepoCreationTypePreview)
		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.Edit(ctx, "o", input)
	if err != nil {
		t.Errorf("Organizations.Edit returned error: %v", err)
	}

	want := &Organization{ID: Ptr(int64(1))}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.Edit returned %+v, want %+v", org, want)
	}

	const methodName = "Edit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.Edit(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.Edit(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_Edit_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.Edit(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_Delete(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.Delete(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.Delete returned error: %v", err)
	}

	const methodName = "Delete"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.Delete(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.Delete(ctx, "o")
	})
}

func TestOrganizationsService_ListInstallations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "installations": [{ "id": 1, "app_id": 5}]}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListInstallations(ctx, "o", nil)
	if err != nil {
		t.Errorf("Organizations.ListInstallations returned error: %v", err)
	}

	want := &OrganizationInstallations{TotalCount: Ptr(1), Installations: []*Installation{{ID: Ptr(int64(1)), AppID: Ptr(int64(5))}}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListInstallations returned %+v, want %+v", apps, want)
	}

	const methodName = "ListInstallations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListInstallations(ctx, "\no", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListInstallations(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_ListInstallations_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.ListInstallations(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_ListInstallations_withListOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `{"total_count": 2, "installations": [{ "id": 2, "app_id": 10}]}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListInstallations(ctx, "o", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Organizations.ListInstallations returned error: %v", err)
	}

	want := &OrganizationInstallations{TotalCount: Ptr(2), Installations: []*Installation{{ID: Ptr(int64(2)), AppID: Ptr(int64(10))}}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListInstallations returned %+v, want %+v", apps, want)
	}

	// Test ListOptions failure
	_, _, err = client.Organizations.ListInstallations(ctx, "%", &ListOptions{})
	if err == nil {
		t.Error("Organizations.ListInstallations returned error: nil")
	}

	const methodName = "ListInstallations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListInstallations(ctx, "\n", &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListInstallations(ctx, "o", &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationInstallations_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &OrganizationInstallations{}, "{}")

	o := &OrganizationInstallations{
		TotalCount:    Ptr(1),
		Installations: []*Installation{{ID: Ptr(int64(1))}},
	}
	want := `{
		"total_count": 1,
		"installations": [
			{
				"id": 1
			}
		]
	}`

	testJSONMarshal(t, o, want)
}

func TestPlan_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Plan{}, "{}")

	o := &Plan{
		Name:          Ptr("name"),
		Space:         Ptr(1),
		Collaborators: Ptr(1),
		PrivateRepos:  Ptr(int64(1)),
		FilledSeats:   Ptr(1),
		Seats:         Ptr(1),
	}
	want := `{
		"name": "name",
		"space": 1,
		"collaborators": 1,
		"private_repos": 1,
		"filled_seats": 1,
		"seats": 1
	}`

	testJSONMarshal(t, o, want)
}
