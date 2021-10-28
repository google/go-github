// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"
	"testing"
)

func TestSCIMService_ListSCIMProvisionedIdentities(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	opts := &ListSCIMProvisionedIdentitiesOptions{}
	_, err := client.SCIM.ListSCIMProvisionedIdentities(ctx, "o", opts)
	if err != nil {
		t.Errorf("SCIM.ListSCIMProvisionedIdentities returned error: %v", err)
	}

	const methodName = "ListSCIMProvisionedIdentities"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SCIM.ListSCIMProvisionedIdentities(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.ListSCIMProvisionedIdentities(ctx, "o", opts)
	})
}

func TestSCIMService_ProvisionAndInviteSCIMUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
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
	_, err := client.SCIM.ProvisionAndInviteSCIMUser(ctx, "o", opts)
	if err != nil {
		t.Errorf("SCIM.ListSCIMProvisionedIdentities returned error: %v", err)
	}

	const methodName = "ProvisionAndInviteSCIMUser"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SCIM.ProvisionAndInviteSCIMUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.ProvisionAndInviteSCIMUser(ctx, "o", opts)
	})
}

func TestSCIMService_GetSCIMProvisioningInfoForUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	_, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
	if err != nil {
		t.Errorf("SCIM.GetSCIMProvisioningInfoForUser returned error: %v", err)
	}

	const methodName = "GetSCIMProvisioningInfoForUser"
	testBadOptions(t, methodName, func() error {
		_, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "\n", "123")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
	})
}

func TestSCIMService_UpdateProvisionedOrgMembership(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
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
	testJSONMarshal(t, &SCIMUserAttributes{}, `{
		"userName":"","name":{"givenName":"","familyName":""},"emails":null
	}`)

	u := &SCIMUserAttributes{
		UserName: "userName1",
		Name: SCIMUserName{
			GivenName:  "Name1",
			FamilyName: "Fname",
			Formatted:  String("formatted name"),
		},
		DisplayName: String("Name"),
		Emails: []*SCIMUserEmail{
			{
				Value:   "value",
				Primary: Bool(false),
				Type:    String("type"),
			},
		},
		Schemas:    []string{"schema1"},
		ExternalID: String("id"),
		Groups:     []string{"group1"},
		Active:     Bool(true),
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
