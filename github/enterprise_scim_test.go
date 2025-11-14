// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSCIMEnterpriseGroups_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMEnterpriseGroups{}, "{}")

	u := &SCIMEnterpriseGroups{
		Schemas:      []string{SCIMSchemasURINamespacesListResponse},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(1),
		StartIndex:   Ptr(1),
		Resources: []*SCIMEnterpriseGroupAttributes{{
			DisplayName: Ptr("gn1"),
			Members: []*SCIMEnterpriseDisplayReference{{
				Value:   "idm1",
				Ref:     "https://api.github.com/scim/v2/enterprises/ee/Users/idm1",
				Display: Ptr("m1"),
			}},
			Schemas:    []string{SCIMSchemasURINamespacesGroups},
			ExternalID: Ptr("eidgn1"),
			ID:         Ptr("idgn1"),
			Meta: &SCIMEnterpriseMeta{
				ResourceType: "Group",
				Created:      &Timestamp{referenceTime},
				LastModified: &Timestamp{referenceTime},
				Location:     Ptr("https://api.github.com/scim/v2/enterprises/ee/Groups/idgn1"),
			},
		}},
	}

	want := `{
		"schemas": ["` + SCIMSchemasURINamespacesListResponse + `"],
        "totalResults": 1,
        "itemsPerPage": 1,
        "startIndex": 1,
        "Resources": [{
            "schemas": ["` + SCIMSchemasURINamespacesGroups + `"],
            "id": "idgn1",
            "externalId": "eidgn1",
            "displayName": "gn1",
            "meta": {
                "resourceType": "Group",
                "created": ` + referenceTimeStr + `,
                "lastModified": ` + referenceTimeStr + `,
                "location": "https://api.github.com/scim/v2/enterprises/ee/Groups/idgn1"
            },
            "members": [{
                "value": "idm1",
                "$ref": "https://api.github.com/scim/v2/enterprises/ee/Users/idm1",
                "display": "m1"
            }]
        }]
    }`

	testJSONMarshal(t, u, want)
}

func TestSCIMEnterpriseUsers_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMEnterpriseUsers{}, "{}")

	u := &SCIMEnterpriseUsers{
		Schemas:      []string{SCIMSchemasURINamespacesListResponse},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(1),
		StartIndex:   Ptr(1),
		Resources: []*SCIMEnterpriseUserAttributes{{
			Active: true,
			Emails: []*SCIMEnterpriseUserEmail{{
				Primary: true,
				Type:    "work",
				Value:   "un1@email.com",
			}},
			Roles: []*SCIMEnterpriseUserRole{{
				Display: Ptr("rd1"),
				Primary: Ptr(true),
				Type:    Ptr("rt1"),
				Value:   "rv1",
			}},
			Schemas:  []string{SCIMSchemasURINamespacesUser},
			UserName: "un1",
			Groups: []*SCIMEnterpriseDisplayReference{{
				Value:   "idgn1",
				Ref:     "https://api.github.com/scim/v2/enterprises/ee/Groups/idgn1",
				Display: Ptr("gn1"),
			}},
			ID:          Ptr("idun1"),
			ExternalID:  "eidun1",
			DisplayName: "dun1",
			Meta: &SCIMEnterpriseMeta{
				ResourceType: "User",
				Created:      &Timestamp{referenceTime},
				LastModified: &Timestamp{referenceTime},
				Location:     Ptr("https://api.github.com/scim/v2/enterprises/ee/User/idun1"),
			},
			Name: &SCIMEnterpriseUserName{
				GivenName:  "gnn1",
				FamilyName: "fnn1",
				Formatted:  Ptr("f1"),
				MiddleName: Ptr("mn1"),
			},
		}},
	}

	want := `{
        "schemas": ["` + SCIMSchemasURINamespacesListResponse + `"],
        "TotalResults": 1,
        "itemsPerPage": 1,
        "StartIndex": 1,
        "Resources": [{
            "active": true,
            "emails": [{
                "primary": true,
                "type": "work",
                "value": "un1@email.com"
            }],
            "roles": [{
                "display": "rd1",
                "primary": true,
                "type": "rt1",
                "value": "rv1"
            }],
            "schemas": ["` + SCIMSchemasURINamespacesUser + `"],
            "userName": "un1",
            "groups": [{
                "value": "idgn1",
                "$ref": "https://api.github.com/scim/v2/enterprises/ee/Groups/idgn1",
                "display": "gn1"
            }],
            "id": "idun1",
            "externalId": "eidun1",
            "name": {
                "givenName": "gnn1",
                "familyName": "fnn1",
                "formatted": "f1",
                "middleName": "mn1"
            },
            "displayName": "dun1",
            "meta": {
                "resourceType": "User",
                "created": ` + referenceTimeStr + `,
                "lastModified": ` + referenceTimeStr + `,
                "location": "https://api.github.com/scim/v2/enterprises/ee/User/idun1"
            }
        }]
    }`

	testJSONMarshal(t, u, want)
}

func TestListProvisionedSCIMGroupsEnterpriseOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ListProvisionedSCIMGroupsEnterpriseOptions{}, "{}")

	u := &ListProvisionedSCIMGroupsEnterpriseOptions{
		Filter:             "f",
		ExcludedAttributes: "ea",
		StartIndex:         5,
		Count:              9,
	}

	want := `{
        "filter": "f",
        "excludedAttributes": "ea",
        "startIndex": 5,
        "count": 9
	}`

	testJSONMarshal(t, u, want)
}

func TestListProvisionedSCIMUsersEnterpriseOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ListProvisionedSCIMUsersEnterpriseOptions{}, "{}")

	u := &ListProvisionedSCIMUsersEnterpriseOptions{
		Filter:     "f",
		StartIndex: 3,
		Count:      7,
	}

	want := `{
        "filter": "f",
        "startIndex": 3,
        "count": 7
	}`

	testJSONMarshal(t, u, want)
}

func TestSCIMEnterpriseGroupAttributes_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMEnterpriseGroupAttributes{}, "{}")

	u := &SCIMEnterpriseGroupAttributes{
		DisplayName: Ptr("dn"),
		Members: []*SCIMEnterpriseDisplayReference{{
			Value:   "v",
			Ref:     "r",
			Display: Ptr("d"),
		}},
		ExternalID: Ptr("eid"),
		ID:         Ptr("id"),
		Schemas:    []string{"s1"},
		Meta: &SCIMEnterpriseMeta{
			ResourceType: "rt",
			Created:      &Timestamp{referenceTime},
			LastModified: &Timestamp{referenceTime},
			Location:     Ptr("l"),
		},
	}

	want := `{
        "schemas": ["s1"],
        "externalId": "eid",
		"displayName": "dn",
        "members" : [{
            "value": "v",
            "$ref": "r",
            "display": "d"
        }],
        "id": "id",
		"meta": {
			"resourceType": "rt",
			"created": ` + referenceTimeStr + `,
			"lastModified": ` + referenceTimeStr + `,
			"location": "l"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestEnterpriseService_ListProvisionedSCIMEnterpriseGroups(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/scim/v2/enterprises/ee/Groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"startIndex":         "1",
			"excludedAttributes": "members,meta",
			"count":              "3",
			"filter":             `externalId eq "914a"`,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"schemas": ["` + SCIMSchemasURINamespacesListResponse + `"],
            "totalResults": 1,
            "itemsPerPage": 1,
            "startIndex": 1,
            "Resources": [{
                "schemas": ["` + SCIMSchemasURINamespacesGroups + `"],
                "id": "914a",
                "externalId": "de88",
                "displayName": "gn1",
                "meta": {
                    "resourceType": "Group",
                    "created": ` + referenceTimeStr + `,
                    "lastModified": ` + referenceTimeStr + `,
                    "location": "https://api.github.com/scim/v2/enterprises/ee/Groups/914a"
                },
                "members": [{
                    "value": "e7f9",
                    "$ref": "https://api.github.com/scim/v2/enterprises/ee/Users/e7f9",
                    "display": "d1"
                }]
            }]
        }`))
	})

	ctx := t.Context()
	opts := &ListProvisionedSCIMGroupsEnterpriseOptions{
		StartIndex:         1,
		ExcludedAttributes: "members,meta",
		Count:              3,
		Filter:             `externalId eq "914a"`,
	}
	groups, _, err := client.Enterprise.ListProvisionedSCIMGroups(ctx, "ee", opts)
	if err != nil {
		t.Errorf("Enterprise.ListProvisionedSCIMGroups returned error: %v", err)
	}

	want := SCIMEnterpriseGroups{
		Schemas:      []string{SCIMSchemasURINamespacesListResponse},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(1),
		StartIndex:   Ptr(1),
		Resources: []*SCIMEnterpriseGroupAttributes{{
			ID: Ptr("914a"),
			Meta: &SCIMEnterpriseMeta{
				ResourceType: "Group",
				Created:      &Timestamp{referenceTime},
				LastModified: &Timestamp{referenceTime},
				Location:     Ptr("https://api.github.com/scim/v2/enterprises/ee/Groups/914a"),
			},
			DisplayName: Ptr("gn1"),
			Schemas:     []string{SCIMSchemasURINamespacesGroups},
			ExternalID:  Ptr("de88"),
			Members: []*SCIMEnterpriseDisplayReference{{
				Value:   "e7f9",
				Ref:     "https://api.github.com/scim/v2/enterprises/ee/Users/e7f9",
				Display: Ptr("d1"),
			}},
		}},
	}

	if diff := cmp.Diff(want, *groups); diff != "" {
		t.Errorf("Enterprise.ListProvisionedSCIMGroups diff mismatch (-want +got):\n%v", diff)
	}

	const methodName = "ListProvisionedSCIMGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListProvisionedSCIMGroups(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, r, err := client.Enterprise.ListProvisionedSCIMGroups(ctx, "o", opts)
		return r, err
	})
}
