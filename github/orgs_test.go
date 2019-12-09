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
	"reflect"
	"testing"
)

func TestOrganization_marshal(t *testing.T) {
	testJSONMarshal(t, &Organization{}, "{}")

	o := &Organization{
		BillingEmail:                         String("support@github.com"),
		Blog:                                 String("https://github.com/blog"),
		Company:                              String("GitHub"),
		Email:                                String("support@github.com"),
		Location:                             String("San Francisco"),
		Name:                                 String("github"),
		Description:                          String("GitHub, the company."),
		DefaultRepoPermission:                String("read"),
		MembersCanCreateRepos:                Bool(true),
		MembersCanCreateInternalRepos:        Bool(true),
		MembersCanCreatePrivateRepos:         Bool(true),
		MembersCanCreatePublicRepos:          Bool(false),
		MembersAllowedRepositoryCreationType: String("all"),
	}
	want := `
		{
			"billing_email": "support@github.com",
			"blog": "https://github.com/blog",
			"company": "GitHub",
			"email": "support@github.com",
			"location": "San Francisco",
			"name": "github",
			"description": "GitHub, the company.",
			"default_repository_permission": "read",
			"members_can_create_repositories": true,
			"members_can_create_public_repositories": false,
			"members_can_create_private_repositories": true,
			"members_can_create_internal_repositories": true,
			"members_allowed_repository_creation_type": "all"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestOrganizationsService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	since := int64(1342004)
	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1342004"})
		fmt.Fprint(w, `[{"id":4314092}]`)
	})

	opt := &OrganizationsListOptions{Since: since}
	orgs, _, err := client.Organizations.ListAll(context.Background(), opt)
	if err != nil {
		t.Errorf("Organizations.ListAll returned error: %v", err)
	}

	want := []*Organization{{ID: Int64(4314092)}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.ListAll returned %+v, want %+v", orgs, want)
	}
}

func TestOrganizationsService_List_authenticatedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	orgs, _, err := client.Organizations.List(context.Background(), "", nil)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []*Organization{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}
}

func TestOrganizationsService_List_specifiedUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	orgs, _, err := client.Organizations.List(context.Background(), "u", opt)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []*Organization{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}
}

func TestOrganizationsService_List_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Organizations.List(context.Background(), "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMemberAllowedRepoCreationTypePreview)
		fmt.Fprint(w, `{"id":1, "login":"l", "url":"u", "avatar_url": "a", "location":"l"}`)
	})

	org, _, err := client.Organizations.Get(context.Background(), "o")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	want := &Organization{ID: Int64(1), Login: String("l"), URL: String("u"), AvatarURL: String("a"), Location: String("l")}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Organizations.Get returned %+v, want %+v", org, want)
	}
}

func TestOrganizationsService_Get_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Organizations.Get(context.Background(), "%")
	testURLParseError(t, err)
}

func TestOrganizationsService_GetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "login":"l", "url":"u", "avatar_url": "a", "location":"l"}`)
	})

	org, _, err := client.Organizations.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Organizations.GetByID returned error: %v", err)
	}

	want := &Organization{ID: Int64(1), Login: String("l"), URL: String("u"), AvatarURL: String("a"), Location: String("l")}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Organizations.GetByID returned %+v, want %+v", org, want)
	}
}

func TestOrganizationsService_Edit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Organization{Login: String("l")}

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		v := new(Organization)
		json.NewDecoder(r.Body).Decode(v)

		testHeader(t, r, "Accept", mediaTypeMemberAllowedRepoCreationTypePreview)
		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	org, _, err := client.Organizations.Edit(context.Background(), "o", input)
	if err != nil {
		t.Errorf("Organizations.Edit returned error: %v", err)
	}

	want := &Organization{ID: Int64(1)}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Organizations.Edit returned %+v, want %+v", org, want)
	}
}

func TestOrganizationsService_Edit_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Organizations.Edit(context.Background(), "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_ListInstallations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIntegrationPreview)
		fmt.Fprint(w, `{"total_count": 1, "installations": [{ "id": 1, "app_id": 5}]}`)
	})

	apps, _, err := client.Organizations.ListInstallations(context.Background(), "o", nil)
	if err != nil {
		t.Errorf("Organizations.ListInstallations returned error: %v", err)
	}

	want := &OrganizationInstallations{TotalCount: Int(1), Installations: []*Installation{{ID: Int64(1), AppID: Int64(5)}}}
	if !reflect.DeepEqual(apps, want) {
		t.Errorf("Organizations.ListInstallations returned %+v, want %+v", apps, want)
	}
}

func TestOrganizationsService_ListInstallations_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Organizations.ListInstallations(context.Background(), "%", nil)
	testURLParseError(t, err)

}

func TestOrganizationsService_ListInstallations_withListOptions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIntegrationPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `{"total_count": 2, "installations": [{ "id": 2, "app_id": 10}]}`)
	})

	apps, _, err := client.Organizations.ListInstallations(context.Background(), "o", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Organizations.ListInstallations returned error: %v", err)
	}

	want := &OrganizationInstallations{TotalCount: Int(2), Installations: []*Installation{{ID: Int64(2), AppID: Int64(10)}}}
	if !reflect.DeepEqual(apps, want) {
		t.Errorf("Organizations.ListInstallations returned %+v, want %+v", apps, want)
	}

	// Test ListOptions failure
	_, _, err = client.Organizations.ListInstallations(context.Background(), "%", &ListOptions{})
	if err == nil {
		t.Error("Organizations.ListInstallations returned error: nil")
	}
}
