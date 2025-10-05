// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_GetConsumedLicenses(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/consumed-licenses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "10"})
		fmt.Fprint(w, `{
			"total_seats_consumed": 20,
			"total_seats_purchased": 25,
			"users": [{
				"github_com_login": "user1",
				"github_com_name": "User One",
				"enterprise_server_user_ids": ["123", "456"],
				"github_com_user": true,
				"enterprise_server_user": false,
				"visual_studio_subscription_user": false,
				"license_type": "Enterprise",
				"github_com_profile": "https://github.com/user1",
				"github_com_member_roles": ["member"],
				"github_com_enterprise_roles": ["member"],
				"github_com_verified_domain_emails": ["user1@example.com"],
				"github_com_saml_name_id": "saml123",
				"github_com_orgs_with_pending_invites": [],
				"github_com_two_factor_auth": ["enabled"],
				"enterprise_server_emails": ["user1@enterprise.local"],
				"visual_studio_license_status": "active",
				"visual_studio_subscription_email": "user1@visualstudio.com",
				"total_user_accounts": 1,
				"github_com_enterprise_role": "owner"
			}]
		}`)
	})

	opt := &ListOptions{Page: 2, PerPage: 10}
	ctx := t.Context()
	licenses, _, err := client.Enterprise.GetConsumedLicenses(ctx, "e", opt)
	if err != nil {
		t.Errorf("Enterprise.GetConsumedLicenses returned error: %v", err)
	}

	want := &EnterpriseConsumedLicenses{
		TotalSeatsConsumed:  20,
		TotalSeatsPurchased: 25,
		Users: []*EnterpriseLicensedUsers{
			{
				GithubComLogin:                  "user1",
				GithubComName:                   "User One",
				EnterpriseServerUserIDs:         []string{"123", "456"},
				GithubComUser:                   true,
				EnterpriseServerUser:            false,
				VisualStudioSubscriptionUser:    false,
				LicenseType:                     "Enterprise",
				GithubComProfile:                "https://github.com/user1",
				GithubComMemberRoles:            []string{"member"},
				GithubComEnterpriseRoles:        []string{"member"},
				GithubComVerifiedDomainEmails:   []string{"user1@example.com"},
				GithubComSamlNameID:             "saml123",
				GithubComOrgsWithPendingInvites: []string{},
				GithubComTwoFactorAuth:          []string{"enabled"},
				EnterpriseServerEmails:          []string{"user1@enterprise.local"},
				VisualStudioLicenseStatus:       "active",
				VisualStudioSubscriptionEmail:   "user1@visualstudio.com",
				TotalUserAccounts:               1,
				GithubComEnterpriseRole:         "owner",
			},
		},
	}

	if !cmp.Equal(licenses, want) {
		t.Errorf("Enterprise.GetConsumedLicenses returned %+v, want %+v", licenses, want)
	}

	const methodName = "GetConsumedLicenses"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetConsumedLicenses(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetConsumedLicenses(ctx, "e", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetLicenseSyncStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/license-sync-status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"title": "Enterprise License Sync Status",
			"description": "Status of license synchronization",
			"properties": {
				"server_instances": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"server_id": "ghes-1",
							"hostname": "github.enterprise.local",
							"last_sync": {
								"type": "object",
								"properties": {
									"date": "2025-10-30T10:30:00Z",
									"status": "success",
									"error": ""
								}
							}
						}
					}
				}
			}
		}`)
	})

	ctx := t.Context()
	syncStatus, _, err := client.Enterprise.GetLicenseSyncStatus(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.GetLicenseSyncStatus returned error: %v", err)
	}

	want := &EnterpriseLicenseSyncStatus{
		Title:       "Enterprise License Sync Status",
		Description: "Status of license synchronization",
		Properties: &ServerInstanceProperties{
			ServerInstances: &ServerInstances{
				Type: "array",
				Items: &ServiceInstanceItems{
					Type: "object",
					Properties: &ServerItemProperties{
						ServerID: "ghes-1",
						Hostname: "github.enterprise.local",
						LastSync: &LastLicenseSync{
							Type: "object",
							Properties: &LastLicenseSyncProperties{
								Date:   &Timestamp{time.Date(2025, 10, 30, 10, 30, 0, 0, time.UTC)},
								Status: "success",
								Error:  "",
							},
						},
					},
				},
			},
		},
	}

	fmt.Printf("%v\n", cmp.Diff(want, syncStatus))

	if !cmp.Equal(syncStatus, want) {
		t.Errorf("Enterprise.GetLicenseSyncStatus returned %+v, want %+v", syncStatus, want)
	}

	const methodName = "GetLicenseSyncStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetLicenseSyncStatus(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetLicenseSyncStatus(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
