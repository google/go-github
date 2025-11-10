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

func TestSCIMProvisionedGroups_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SCIMEnterpriseGroups{}, "{}")

	u := &SCIMEnterpriseGroups{
		Schemas:      []string{"s1"},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(2),
		StartIndex:   Ptr(3),
		Resources: []*SCIMEnterpriseGroupAttributes{
			{
				DisplayName: Ptr("dn"),
				Members: []*SCIMEnterpriseDisplayReference{
					{
						Value:   "v",
						Ref:     "r",
						Display: Ptr("d"),
					},
				},
				Schemas:    []string{"s2"},
				ExternalID: Ptr("eid"),
				ID:         Ptr("id"),
				Meta: &SCIMEnterpriseMeta{
					ResourceType: Ptr("rt"),
					Created:      &Timestamp{referenceTime},
					LastModified: &Timestamp{referenceTime},
					Location:     Ptr("l"),
				},
			},
		},
	}

	want := `{
		"schemas": ["s1"],
		"totalResults": 1,
		"itemsPerPage": 2,
		"startIndex": 3,
		"Resources": [{
			"displayName": "dn",
			"members": [{
				"value": "v",
				"$ref": "r",
				"display": "d"
			}]
			,"schemas": ["s2"],
			"externalId": "eid",
			"id": "id",
			"meta": {
				"resourceType": "rt",
				"created": ` + referenceTimeStr + `,
				"lastModified": ` + referenceTimeStr + `,
				"location": "l"
			}
		}]
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
			ResourceType: Ptr("rt"),
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
		}]
	}`

	testJSONMarshal(t, u, want)
}

func TestEnterpriseService_ListProvisionedSCIMGroupsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("scim/v2/enterprises/ee/Groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"startIndex":         "1",
			"excludedAttributes": "members,meta",
			"count":              "3",
			"filter":             `externalId eq "8aa1"`,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"schemas": [` + SCIMSchemasURINamespacesListResponse + `],
			"totalResults": 1,
			"Resources": [
			  {
				"schemas": [` + SCIMSchemasURINamespacesGroups + `],
				"externalId": "8aa1",
				"id": "24b2",
				"displayName": "dn",
				"members": [
				  {
					"value": "879d",
					"$+ref": "https://api.github.localhost/scim/v2/Users/879d",
					"display": "u1"
				  },
				  {
					"value": "0db5",
					"$+ref": "https://api.github.localhost/scim/v2/Users/0db5",
					"display": "u2"
				  }
				],
                "meta": {
				  "resourceType": "Group",
				  "created": ` + referenceTimeStr + `,
				  "lastModified": ` + referenceTimeStr + `,
				  "location": "https://api.github.localhost/scim/v2/Groups/24b2"
				}
			  }
			],
			"startIndex": 1,
            "itemsPerPage": 1
		  }`))
	})

	ctx := t.Context()
	opts := &ListProvisionedSCIMGroupsForEnterpriseOptions{
		StartIndex:         Ptr(1),
		ExcludedAttributes: Ptr("members,meta"),
		Count:              Ptr(3),
		Filter:             Ptr(`externalId eq "8aa1"`),
	}
	groups, _, err := client.Enterprise.ListProvisionedSCIMGroupsForEnterprise(ctx, "ee", opts)
	if err != nil {
		t.Errorf("Enterprise.ListSCIMProvisionedIdentities returned error: %v", err)
	}

	want := SCIMEnterpriseGroups{
		Schemas:      []string{SCIMSchemasURINamespacesListResponse},
		TotalResults: Ptr(1),
		ItemsPerPage: Ptr(1),
		StartIndex:   Ptr(1),
		Resources: []*SCIMEnterpriseGroupAttributes{
			{
				ID: Ptr("123e4567-e89b-12d3-a456-426614174000"),
				Meta: &SCIMEnterpriseMeta{
					ResourceType: Ptr("Group"),
					Created:      &Timestamp{referenceTime},
					LastModified: &Timestamp{referenceTime},
					Location:     Ptr("https://api.github.com/scim/v2/enterprises/octo/Groups/123e4567-e89b-12d3-a456-426614174000"),
				},

				DisplayName: Ptr("Mona Octocat"),
				Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
				ExternalID:  Ptr("00u1dhhb1fkIGP7RL1d8"),
				Members: []*SCIMEnterpriseDisplayReference{
					{
						Value:   "5fc0c238-1112-11e8-8e45-920c87bdbd75",
						Ref:     "https://api.github.com/scim/v2/enterprises/octo/Users/5fc0c238-1112-11e8-8e45-920c87bdbd75",
						Display: Ptr("Mona Octocat"),
					},
				},
			},
		},
	}

	if !cmp.Equal(groups, &want) {
		diff := cmp.Diff(groups, want)
		t.Errorf("SCIM.ListSCIMProvisionedGroupsForEnterprise returned %+v, want %+v: diff %+v", groups, want, diff)
	}

	const methodName = "ListSCIMProvisionedGroupsForEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListProvisionedSCIMGroupsForEnterprise(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, r, err := client.Enterprise.ListProvisionedSCIMGroupsForEnterprise(ctx, "o", opts)
		return r, err
	})
}
