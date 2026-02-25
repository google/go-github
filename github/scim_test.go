// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSCIMService_ListSCIMProvisionedIdentities(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"startIndex": "1", "count": "10", "filter": `userName="Octocat"`})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"schemas": [
			  "urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 1,
			"itemsPerPage": 1,
			"startIndex": 1,
			"Resources": [
			  {
				"schemas": [
				  "urn:ietf:params:scim:schemas:core:2.0:User"
				],
				"id": "5fc0c238-1112-11e8-8e45-920c87bdbd75",
				"externalId": "00u1dhhb1fkIGP7RL1d8",
				"userName": "octocat@github.com",
				"displayName": "Mona Octocat",
				"name": {
				  "givenName": "Mona",
				  "familyName": "Octocat",
				  "formatted": "Mona Octocat"
				},
				"emails": [
				  {
					"value": "octocat@github.com",
					"primary": true
				  }
				],
				"active": true,
				"meta": {
				  "resourceType": "User",
				  "created": "2018-02-13T15:05:24.000-00:00",
				  "lastModified": "2018-02-13T15:05:24.000-00:00",
				  "location": "https://api.github.com/scim/v2/organizations/octo-org/Users/5fc0c238-1112-11e8-8e45-920c87bdbd75"
				}
			  }
			]
		  }`))
	})

	ctx := t.Context()
	opts := &ListSCIMProvisionedIdentitiesOptions{
		StartIndex: Ptr(1),
		Count:      Ptr(10),
		Filter:     Ptr(`userName="Octocat"`),
	}
	identities, _, err := client.SCIM.ListSCIMProvisionedIdentities(ctx, "o", opts)
	if err != nil {
		t.Errorf("SCIM.ListSCIMProvisionedIdentities returned error: %v", err)
	}

	date := Timestamp{time.Date(2018, time.February, 13, 15, 5, 24, 0, time.UTC)}
	want := SCIMProvisionedIdentities{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(1),
		StartIndex:   Ptr(1),
		Resources: []*SCIMUserAttributes{
			{
				ID: Ptr("5fc0c238-1112-11e8-8e45-920c87bdbd75"),
				Meta: &SCIMMeta{
					ResourceType: Ptr("User"),
					Created:      &date,
					LastModified: &date,
					Location:     Ptr("https://api.github.com/scim/v2/organizations/octo-org/Users/5fc0c238-1112-11e8-8e45-920c87bdbd75"),
				},
				UserName: "octocat@github.com",
				Name: SCIMUserName{
					GivenName:  "Mona",
					FamilyName: "Octocat",
					Formatted:  Ptr("Mona Octocat"),
				},
				DisplayName: Ptr("Mona Octocat"),
				Emails: []*SCIMUserEmail{
					{
						Value:   "octocat@github.com",
						Primary: Ptr(true),
					},
				},
				Schemas:    []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
				ExternalID: Ptr("00u1dhhb1fkIGP7RL1d8"),
				Groups:     nil,
				Active:     Ptr(true),
			},
		},
	}

	if !cmp.Equal(identities, &want) {
		diff := cmp.Diff(identities, want)
		t.Errorf("SCIM.ListSCIMProvisionedIdentities returned %+v, want %+v: diff %+v", identities, want, diff)
	}

	const methodName = "ListSCIMProvisionedIdentities"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SCIM.ListSCIMProvisionedIdentities(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, r, err := client.SCIM.ListSCIMProvisionedIdentities(ctx, "o", nil)
		return r, err
	})
}

func TestSCIMService_ProvisionAndInviteSCIMUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"1234567890","userName":"userName"}`)
	})

	ctx := t.Context()
	opts := &SCIMUserAttributes{
		UserName: "userName",
		Name: SCIMUserName{
			GivenName:  "givenName",
			FamilyName: "familyName",
		},
		Emails: []*SCIMUserEmail{
			{
				Value: "octocat@github.com",
			},
		},
	}
	user, _, err := client.SCIM.ProvisionAndInviteSCIMUser(ctx, "o", opts)
	if err != nil {
		t.Errorf("SCIM.ProvisionAndInviteSCIMUser returned error: %v", err)
	}

	want := &SCIMUserAttributes{
		ID:       Ptr("1234567890"),
		UserName: "userName",
	}
	if !cmp.Equal(user, want) {
		t.Errorf("SCIM.ProvisionAndInviteSCIMUser returned %+v, want %+v", user, want)
	}

	const methodName = "ProvisionAndInviteSCIMUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SCIM.ProvisionAndInviteSCIMUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SCIM.ProvisionAndInviteSCIMUser(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSCIMService_GetSCIMProvisioningInfoForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"schemas": [
			  "urn:ietf:params:scim:schemas:core:2.0:User"
			],
			"id": "edefdfedf-050c-11e7-8d32",
			"externalId": "a7d0f98382",
			"userName": "mona.octocat@okta.example.com",
			"displayName": "Mona Octocat",
			"name": {
			  "givenName": "Mona",
			  "familyName": "Octocat",
			  "formatted": "Mona Octocat"
			},
			"emails": [
			  {
				"value": "mona.octocat@okta.example.com",
				"primary": true
			  },
			  {
				"value": "mona@octocat.github.com"
			  }
			],
			"active": true,
			"meta": {
			  "resourceType": "User",
			  "created": "2017-03-09T16:11:13-00:00",
			  "lastModified": "2017-03-09T16:11:13-00:00",
			  "location": "https://api.github.com/scim/v2/organizations/octo-org/Users/edefdfedf-050c-11e7-8d32"
			}
		  }`))
	})

	ctx := t.Context()
	user, _, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
	if err != nil {
		t.Errorf("SCIM.GetSCIMProvisioningInfoForUser returned error: %v", err)
	}

	date := Timestamp{time.Date(2017, time.March, 9, 16, 11, 13, 0, time.UTC)}
	want := SCIMUserAttributes{
		ID: Ptr("edefdfedf-050c-11e7-8d32"),
		Meta: &SCIMMeta{
			ResourceType: Ptr("User"),
			Created:      &date,
			LastModified: &date,
			Location:     Ptr("https://api.github.com/scim/v2/organizations/octo-org/Users/edefdfedf-050c-11e7-8d32"),
		},
		UserName: "mona.octocat@okta.example.com",
		Name: SCIMUserName{
			GivenName:  "Mona",
			FamilyName: "Octocat",
			Formatted:  Ptr("Mona Octocat"),
		},
		DisplayName: Ptr("Mona Octocat"),
		Emails: []*SCIMUserEmail{
			{
				Value:   "mona.octocat@okta.example.com",
				Primary: Ptr(true),
			},
			{
				Value: "mona@octocat.github.com",
			},
		},
		Schemas:    []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		ExternalID: Ptr("a7d0f98382"),
		Groups:     nil,
		Active:     Ptr(true),
	}

	if !cmp.Equal(user, &want) {
		diff := cmp.Diff(user, want)
		t.Errorf("SCIM.ListSCIMProvisionedIdentities returned %+v, want %+v: diff %+v", user, want, diff)
	}

	const methodName = "GetSCIMProvisioningInfoForUser"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "\n", "123")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, r, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
		return r, err
	})
}

func TestSCIMService_UpdateProvisionedOrgMembership(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusOK)
	})

	ctx := t.Context()
	opts := &SCIMUserAttributes{
		UserName: "userName",
		Name: SCIMUserName{
			GivenName:  "givenName",
			FamilyName: "familyName",
		},
		Emails: []*SCIMUserEmail{
			{
				Value: "octocat@github.com",
			},
		},
	}
	_, err := client.SCIM.UpdateProvisionedOrgMembership(ctx, "o", "123", opts)
	if err != nil {
		t.Errorf("SCIM.UpdateProvisionedOrgMembership returned error: %v", err)
	}

	const methodName = "UpdateProvisionedOrgMembership"
	testBadOptions(t, methodName, func() error {
		_, err := client.SCIM.UpdateProvisionedOrgMembership(ctx, "\n", "123", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.UpdateProvisionedOrgMembership(ctx, "o", "123", opts)
	})
}

func TestSCIMService_UpdateAttributeForSCIMUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	opts := &UpdateAttributeForSCIMUserOptions{}
	_, err := client.SCIM.UpdateAttributeForSCIMUser(ctx, "o", "123", opts)
	if err != nil {
		t.Errorf("SCIM.UpdateAttributeForSCIMUser returned error: %v", err)
	}

	const methodName = "UpdateAttributeForSCIMUser"
	testBadOptions(t, methodName, func() error {
		_, err := client.SCIM.UpdateAttributeForSCIMUser(ctx, "\n", "123", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.UpdateAttributeForSCIMUser(ctx, "o", "123", opts)
	})
}

func TestSCIMService_DeleteSCIMUserFromOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.SCIM.DeleteSCIMUserFromOrg(ctx, "o", "123")
	if err != nil {
		t.Errorf("SCIM.DeleteSCIMUserFromOrg returned error: %v", err)
	}

	const methodName = "DeleteSCIMUserFromOrg"
	testBadOptions(t, methodName, func() error {
		_, err := client.SCIM.DeleteSCIMUserFromOrg(ctx, "\n", "")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.DeleteSCIMUserFromOrg(ctx, "o", "123")
	})
}

func TestSCIMUserAttributes_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMUserAttributes{}, `{
		"userName":"","name":{"givenName":"","familyName":""},"emails":null
	}`)

	u := &SCIMUserAttributes{
		UserName: "userName1",
		Name: SCIMUserName{
			GivenName:  "Name1",
			FamilyName: "Fname",
			Formatted:  Ptr("formatted name"),
		},
		DisplayName: Ptr("Name"),
		Emails: []*SCIMUserEmail{
			{
				Value:   "value",
				Primary: Ptr(false),
				Type:    Ptr("type"),
			},
		},
		Schemas:    []string{"schema1"},
		ExternalID: Ptr("id"),
		Groups:     []string{"group1"},
		Active:     Ptr(true),
	}

	want := `{
		"userName": "userName1",
		"name": {
			"givenName": "Name1",
			"familyName": "Fname",
			"formatted": "formatted name"
		},
		"displayName": "Name",
		"emails": [{
			"value": "value",
			"primary": false,
			"type": "type"
		}],
		"schemas": ["schema1"],
		"externalId": "id",
		"groups": ["group1"],
		"active": true
	}`

	testJSONMarshal(t, u, want)
}

func TestUpdateAttributeForSCIMUserOperations_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UpdateAttributeForSCIMUserOperations{}, `{}`)

	u := &UpdateAttributeForSCIMUserOperations{
		Op:   "TestOp",
		Path: Ptr("path"),
	}

	want := `{
		"op": "TestOp",
		"path": "path"
	}`

	testJSONMarshal(t, u, want)
}

func TestUpdateAttributeForSCIMUserOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UpdateAttributeForSCIMUserOptions{}, `{}`)

	u := &UpdateAttributeForSCIMUserOptions{
		Schemas: []string{"test", "schema"},
		Operations: UpdateAttributeForSCIMUserOperations{
			Op:   "TestOp",
			Path: Ptr("path"),
		},
	}

	want := `{
		"schemas": ["test", "schema"],
		"operations": {
			"op": "TestOp",
			"path": "path"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestSCIMUserName_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMUserName{}, `{
		"givenName":"","familyName":""
	}`)

	u := &SCIMUserName{
		GivenName:  "Name1",
		FamilyName: "Fname",
		Formatted:  Ptr("formatted name"),
	}

	want := `{
			"givenName": "Name1",
			"familyName": "Fname",
			"formatted": "formatted name"
	}`
	testJSONMarshal(t, u, want)
}

func TestSCIMMeta_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMMeta{}, `{}`)

	u := &SCIMMeta{
		ResourceType: Ptr("test"),
		Location:     Ptr("test"),
	}

	want := `{
		"resourceType": "test",
		"location": "test"
	}`

	testJSONMarshal(t, u, want)
}

func TestSCIMUserRole_Marshal(t *testing.T) {
	t.Parallel()

	testJSONMarshal(t, &SCIMUserRole{
		Value:   "enterprise_owner",
		Primary: Ptr(true),
	}, `{
		"value": "enterprise_owner",
		"primary": true
	}`)

	r := &SCIMUserRole{
		Value: "billing_manager",
	}
	want := `{"value": "billing_manager"}`
	testJSONMarshal(t, r, want)
}

func TestSCIMProvisionedIdentities_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMProvisionedIdentities{}, `{}`)

	u := &SCIMProvisionedIdentities{
		Schemas:      []string{"test", "schema"},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(2),
		StartIndex:   Ptr(1),
		Resources: []*SCIMUserAttributes{
			{
				UserName: "SCIM",
				Name: SCIMUserName{
					GivenName:  "scim",
					FamilyName: "test",
					Formatted:  Ptr("SCIM"),
				},
				DisplayName: Ptr("Test SCIM"),
				Emails: []*SCIMUserEmail{
					{
						Value:   "test",
						Primary: Ptr(true),
						Type:    Ptr("test"),
					},
				},
				Schemas:    []string{"schema1"},
				ExternalID: Ptr("id"),
				Groups:     []string{"group1"},
				Active:     Ptr(true),
			},
		},
	}

	want := `{
		"schemas": ["test", "schema"],
		"totalResults": 1,
		"itemsPerPage": 2,
		"startIndex": 1,
		"Resources": [{
			"userName": "SCIM",
			"name": {
				"givenName": "scim",
				"familyName": "test",
				"formatted": "SCIM"
			},
			"displayName": "Test SCIM",
			"emails": [{
				"value": "test",
				"primary": true,
				"type": "test"
			}],
			"schemas": ["schema1"],
			"externalId": "id",
			"groups": ["group1"],
			"active": true
		}]
	}`

	testJSONMarshal(t, u, want)
}
