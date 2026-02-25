// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Test invalid JSON responses, valid responses are covered in the other tests.
func TestCopilotSeatDetails_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		data    string
		want    *CopilotSeatDetails
		wantErr bool
	}{
		{
			name: "Invalid JSON",
			data: `{`,
			want: &CopilotSeatDetails{
				Assignee: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid top level type",
			data: `{
					"assignee": {
						"type": "User",
						"name": "octokittens",
						"id": 1
					},
					"assigning_team": "this should be an object"
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "No Type Field",
			data: `{
					"assignee": {
						"name": "octokittens",
						"id": 1
					}
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "Null Assignee",
			data: `{
					"assignee": null
				}`,
			want: &CopilotSeatDetails{
				Assignee: nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid Assignee Field Type",
			data: `{
					"assignee": "test"
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "Invalid Assignee Type",
			data: `{
					"assignee": {
						"name": "octokittens",
						"id": 1,
						"type": []
					}
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "Invalid User",
			data: `{
					"assignee": {
						"type": "User",
						"id": "bad"
					}
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "Invalid Team",
			data: `{
					"assignee": {
						"type": "Team",
						"id": "bad"
					}
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
		{
			name: "Invalid Organization",
			data: `{
					"assignee": {
						"type": "Organization",
						"id": "bad"
					}
				}`,
			want:    &CopilotSeatDetails{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		seatDetails := &CopilotSeatDetails{}

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tc.data), seatDetails)
			if err == nil && tc.wantErr {
				t.Error("CopilotSeatDetails.UnmarshalJSON returned nil instead of an error")
			}
			if err != nil && !tc.wantErr {
				t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
			}
			if !cmp.Equal(tc.want, seatDetails) {
				t.Errorf("CopilotSeatDetails.UnmarshalJSON expected %+v, got %+v", tc.want, seatDetails)
			}
		})
	}
}

func TestCopilotService_GetSeatDetailsUser(t *testing.T) {
	t.Parallel()
	data := `{
				"assignee": {
					"type": "User",
					"id": 1
				}
			}`

	seatDetails := &CopilotSeatDetails{}

	err := json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	want := &User{
		ID:   Ptr(int64(1)),
		Type: Ptr("User"),
	}

	if got, ok := seatDetails.GetUser(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetTeam returned %+v, want %+v", got, want)
	} else if !ok {
		t.Error("CopilotSeatDetails.GetUser returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "Organization",
					"id": 1
				}
			}`

	bad := &Organization{
		ID:   Ptr(int64(1)),
		Type: Ptr("Organization"),
	}

	err = json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	if got, ok := seatDetails.GetUser(); ok {
		t.Errorf("CopilotSeatDetails.GetUser returned true, expected false. Returned %v, expected %v", got, bad)
	}
}

func TestCopilotService_GetSeatDetailsTeam(t *testing.T) {
	t.Parallel()
	data := `{
				"assignee": {
					"type": "Team",
					"id": 1
				}
			}`

	seatDetails := &CopilotSeatDetails{}

	err := json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	want := &Team{
		ID:   Ptr(int64(1)),
		Type: Ptr("Team"),
	}

	if got, ok := seatDetails.GetTeam(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetTeam returned %+v, want %+v", got, want)
	} else if !ok {
		t.Error("CopilotSeatDetails.GetTeam returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "User",
					"id": 1
				}
			}`

	bad := &User{
		ID:   Ptr(int64(1)),
		Type: Ptr("User"),
	}

	err = json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	if got, ok := seatDetails.GetTeam(); ok {
		t.Errorf("CopilotSeatDetails.GetTeam returned true, expected false. Returned %v, expected %v", got, bad)
	}
}

func TestCopilotService_GetSeatDetailsOrganization(t *testing.T) {
	t.Parallel()
	data := `{
				"assignee": {
					"type": "Organization",
					"id": 1
				}
			}`

	seatDetails := &CopilotSeatDetails{}

	err := json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	want := &Organization{
		ID:   Ptr(int64(1)),
		Type: Ptr("Organization"),
	}

	if got, ok := seatDetails.GetOrganization(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetOrganization returned %+v, want %+v", got, want)
	} else if !ok {
		t.Error("CopilotSeatDetails.GetOrganization returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "Team",
					"id": 1
				}
			}`

	bad := &Team{
		ID: Ptr(int64(1)),
	}

	err = json.Unmarshal([]byte(data), seatDetails)
	if err != nil {
		t.Errorf("CopilotSeatDetails.UnmarshalJSON returned an unexpected error: %v", err)
	}

	if got, ok := seatDetails.GetOrganization(); ok {
		t.Errorf("CopilotSeatDetails.GetOrganization returned true, expected false. Returned %v, expected %v", got, bad)
	}
}

func TestCopilotService_GetCopilotBilling(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"seat_breakdown": {
				"total": 12,
				"added_this_cycle": 9,
				"pending_invitation": 0,
				"pending_cancellation": 0,
				"active_this_cycle": 12,
				"inactive_this_cycle": 11
			},
			"seat_management_setting": "assign_selected",
			"public_code_suggestions": "block"
			}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetCopilotBilling(ctx, "o")
	if err != nil {
		t.Errorf("Copilot.GetCopilotBilling returned error: %v", err)
	}

	want := &CopilotOrganizationDetails{
		SeatBreakdown: &CopilotSeatBreakdown{
			Total:               12,
			AddedThisCycle:      9,
			PendingInvitation:   0,
			PendingCancellation: 0,
			ActiveThisCycle:     12,
			InactiveThisCycle:   11,
		},
		PublicCodeSuggestions: "block",
		CopilotChat:           "",
		SeatManagementSetting: "assign_selected",
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetCopilotBilling returned %+v, want %+v", got, want)
	}

	const methodName = "GetCopilotBilling"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetCopilotBilling(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetCopilotBilling(ctx, "o")
		if got != nil {
			t.Errorf("Copilot.GetCopilotBilling returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_ListCopilotSeats(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing/seats", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"per_page": "100",
			"page":     "1",
		})
		fmt.Fprint(w, `{
				"total_seats": 4,
				"seats": [
					{
						"created_at": "2021-08-03T18:00:00-06:00",
						"updated_at": "2021-09-23T15:00:00-06:00",
						"pending_cancellation_date": null,
						"last_activity_at": "2021-10-14T00:53:32-06:00",
						"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
						"assignee": {
							"login": "octocat",
							"id": 1,
							"node_id": "MDQ6VXNlcjE=",
							"avatar_url": "https://github.com/images/error/octocat_happy.gif",
							"gravatar_id": "",
							"url": "https://api.github.com/users/octocat",
							"html_url": "https://github.com/octocat",
							"followers_url": "https://api.github.com/users/octocat/followers",
							"following_url": "https://api.github.com/users/octocat/following{/other_user}",
							"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
							"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
							"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
							"organizations_url": "https://api.github.com/users/octocat/orgs",
							"repos_url": "https://api.github.com/users/octocat/repos",
							"events_url": "https://api.github.com/users/octocat/events{/privacy}",
							"received_events_url": "https://api.github.com/users/octocat/received_events",
							"type": "User",
							"site_admin": false
						},
						"assigning_team": {
							"id": 1,
							"node_id": "MDQ6VGVhbTE=",
							"url": "https://api.github.com/teams/1",
							"html_url": "https://github.com/orgs/github/teams/justice-league",
							"name": "Justice League",
							"slug": "justice-league",
							"description": "A great team.",
							"privacy": "closed",
							"notification_setting": "notifications_enabled",
							"permission": "admin",
							"members_url": "https://api.github.com/teams/1/members{/member}",
							"repositories_url": "https://api.github.com/teams/1/repos",
							"parent": null
						}
					},
					{
						"created_at": "2021-09-23T18:00:00-06:00",
						"updated_at": "2021-09-23T15:00:00-06:00",
						"pending_cancellation_date": "2021-11-01",
						"last_activity_at": "2021-10-13T00:53:32-06:00",
						"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
						"assignee": {
							"login": "octokitten",
							"id": 1,
							"node_id": "MDQ76VNlcjE=",
							"avatar_url": "https://github.com/images/error/octokitten_happy.gif",
							"gravatar_id": "",
							"url": "https://api.github.com/users/octokitten",
							"html_url": "https://github.com/octokitten",
							"followers_url": "https://api.github.com/users/octokitten/followers",
							"following_url": "https://api.github.com/users/octokitten/following{/other_user}",
							"gists_url": "https://api.github.com/users/octokitten/gists{/gist_id}",
							"starred_url": "https://api.github.com/users/octokitten/starred{/owner}{/repo}",
							"subscriptions_url": "https://api.github.com/users/octokitten/subscriptions",
							"organizations_url": "https://api.github.com/users/octokitten/orgs",
							"repos_url": "https://api.github.com/users/octokitten/repos",
							"events_url": "https://api.github.com/users/octokitten/events{/privacy}",
							"received_events_url": "https://api.github.com/users/octokitten/received_events",
							"type": "User",
							"site_admin": false
						}
					},
					{
						"created_at": "2021-09-23T18:00:00-06:00",
						"updated_at": "2021-09-23T15:00:00-06:00",
						"pending_cancellation_date": "2021-11-01",
						"last_activity_at": "2021-10-13T00:53:32-06:00",
						"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
						"assignee": {
							"name": "octokittens",
							"id": 1,
							"type": "Team"
						}
					},
					{
						"created_at": "2021-09-23T18:00:00-06:00",
						"updated_at": "2021-09-23T15:00:00-06:00",
						"pending_cancellation_date": "2021-11-01",
						"last_activity_at": "2021-10-13T00:53:32-06:00",
						"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
						"assignee": {
							"name": "octocats",
							"id": 1,
							"type": "Organization"
						}
					}
				]
			}`)
	})

	tmp, err := time.Parse(time.RFC3339, "2021-08-03T18:00:00-06:00")
	if err != nil {
		panic(err)
	}
	createdAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T15:00:00-06:00")
	if err != nil {
		panic(err)
	}
	updatedAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-10-14T00:53:32-06:00")
	if err != nil {
		panic(err)
	}
	lastActivityAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T18:00:00-06:00")
	if err != nil {
		panic(err)
	}
	createdAt2 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T15:00:00-06:00")
	if err != nil {
		panic(err)
	}
	updatedAt2 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-10-13T00:53:32-06:00")
	if err != nil {
		panic(err)
	}
	lastActivityAt2 := Timestamp{tmp}

	ctx := t.Context()
	opts := &ListOptions{Page: 1, PerPage: 100}
	got, _, err := client.Copilot.ListCopilotSeats(ctx, "o", opts)
	if err != nil {
		t.Errorf("Copilot.ListCopilotSeats returned error: %v", err)
	}

	want := &ListCopilotSeatsResponse{
		TotalSeats: 4,
		Seats: []*CopilotSeatDetails{
			{
				Assignee: &User{
					Login:             Ptr("octocat"),
					ID:                Ptr(int64(1)),
					NodeID:            Ptr("MDQ6VXNlcjE="),
					AvatarURL:         Ptr("https://github.com/images/error/octocat_happy.gif"),
					GravatarID:        Ptr(""),
					URL:               Ptr("https://api.github.com/users/octocat"),
					HTMLURL:           Ptr("https://github.com/octocat"),
					FollowersURL:      Ptr("https://api.github.com/users/octocat/followers"),
					FollowingURL:      Ptr("https://api.github.com/users/octocat/following{/other_user}"),
					GistsURL:          Ptr("https://api.github.com/users/octocat/gists{/gist_id}"),
					StarredURL:        Ptr("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
					SubscriptionsURL:  Ptr("https://api.github.com/users/octocat/subscriptions"),
					OrganizationsURL:  Ptr("https://api.github.com/users/octocat/orgs"),
					ReposURL:          Ptr("https://api.github.com/users/octocat/repos"),
					EventsURL:         Ptr("https://api.github.com/users/octocat/events{/privacy}"),
					ReceivedEventsURL: Ptr("https://api.github.com/users/octocat/received_events"),
					Type:              Ptr("User"),
					SiteAdmin:         Ptr(false),
				},
				AssigningTeam: &Team{
					ID:                  Ptr(int64(1)),
					NodeID:              Ptr("MDQ6VGVhbTE="),
					URL:                 Ptr("https://api.github.com/teams/1"),
					HTMLURL:             Ptr("https://github.com/orgs/github/teams/justice-league"),
					Name:                Ptr("Justice League"),
					Slug:                Ptr("justice-league"),
					Description:         Ptr("A great team."),
					Privacy:             Ptr("closed"),
					Permission:          Ptr("admin"),
					NotificationSetting: Ptr("notifications_enabled"),
					MembersURL:          Ptr("https://api.github.com/teams/1/members{/member}"),
					RepositoriesURL:     Ptr("https://api.github.com/teams/1/repos"),
					Parent:              nil,
				},
				CreatedAt:               &createdAt1,
				UpdatedAt:               &updatedAt1,
				PendingCancellationDate: nil,
				LastActivityAt:          &lastActivityAt1,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &User{
					Login:             Ptr("octokitten"),
					ID:                Ptr(int64(1)),
					NodeID:            Ptr("MDQ76VNlcjE="),
					AvatarURL:         Ptr("https://github.com/images/error/octokitten_happy.gif"),
					GravatarID:        Ptr(""),
					URL:               Ptr("https://api.github.com/users/octokitten"),
					HTMLURL:           Ptr("https://github.com/octokitten"),
					FollowersURL:      Ptr("https://api.github.com/users/octokitten/followers"),
					FollowingURL:      Ptr("https://api.github.com/users/octokitten/following{/other_user}"),
					GistsURL:          Ptr("https://api.github.com/users/octokitten/gists{/gist_id}"),
					StarredURL:        Ptr("https://api.github.com/users/octokitten/starred{/owner}{/repo}"),
					SubscriptionsURL:  Ptr("https://api.github.com/users/octokitten/subscriptions"),
					OrganizationsURL:  Ptr("https://api.github.com/users/octokitten/orgs"),
					ReposURL:          Ptr("https://api.github.com/users/octokitten/repos"),
					EventsURL:         Ptr("https://api.github.com/users/octokitten/events{/privacy}"),
					ReceivedEventsURL: Ptr("https://api.github.com/users/octokitten/received_events"),
					Type:              Ptr("User"),
					SiteAdmin:         Ptr(false),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: Ptr("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &Team{
					ID:   Ptr(int64(1)),
					Name: Ptr("octokittens"),
					Type: Ptr("Team"),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: Ptr("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &Organization{
					ID:   Ptr(int64(1)),
					Name: Ptr("octocats"),
					Type: Ptr("Organization"),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: Ptr("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
			},
		},
	}

	assertNoDiff(t, want, got)

	const methodName = "ListCopilotSeats"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.ListCopilotSeats(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.ListCopilotSeats(ctx, "o", opts)
		if got != nil {
			t.Errorf("Copilot.ListCopilotSeats returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_ListCopilotEnterpriseSeats(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/billing/seats", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"per_page": "100",
			"page":     "1",
		})
		fmt.Fprint(w, `{
			"total_seats": 2,
			"seats": [
				{
					"created_at": "2021-08-03T18:00:00-06:00",
					"updated_at": "2021-09-23T15:00:00-06:00",
					"pending_cancellation_date": null,
					"last_activity_at": "2021-10-14T00:53:32-06:00",
					"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
					"plan_type": "business",
					"assignee": {
						"login": "octocat",
						"id": 1,
						"node_id": "MDQ6VXNlcjE=",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octocat",
						"html_url": "https://github.com/octocat",
						"followers_url": "https://api.github.com/users/octocat/followers",
						"following_url": "https://api.github.com/users/octocat/following{/other_user}",
						"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
						"organizations_url": "https://api.github.com/users/octocat/orgs",
						"repos_url": "https://api.github.com/users/octocat/repos",
						"events_url": "https://api.github.com/users/octocat/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octocat/received_events",
						"type": "User",
						"site_admin": false
					},
					"assigning_team": {
						"id": 1,
						"node_id": "MDQ6VGVhbTE=",
						"url": "https://api.github.com/teams/1",
						"html_url": "https://github.com/orgs/github/teams/justice-league",
						"name": "Justice League",
						"slug": "justice-league",
						"description": "A great team.",
						"privacy": "closed",
						"notification_setting": "notifications_enabled",
						"permission": "admin",
						"members_url": "https://api.github.com/teams/1/members{/member}",
						"repositories_url": "https://api.github.com/teams/1/repos",
						"parent": null
					}
				},
				{
					"created_at": "2021-09-23T18:00:00-06:00",
					"updated_at": "2021-09-23T15:00:00-06:00",
					"pending_cancellation_date": "2021-11-01",
					"last_activity_at": "2021-10-13T00:53:32-06:00",
					"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
					"assignee": {
						"login": "octokitten",
						"id": 1,
						"node_id": "MDQ76VNlcjE=",
						"avatar_url": "https://github.com/images/error/octokitten_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octokitten",
						"html_url": "https://github.com/octokitten",
						"followers_url": "https://api.github.com/users/octokitten/followers",
						"following_url": "https://api.github.com/users/octokitten/following{/other_user}",
						"gists_url": "https://api.github.com/users/octokitten/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octokitten/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octokitten/subscriptions",
						"organizations_url": "https://api.github.com/users/octokitten/orgs",
						"repos_url": "https://api.github.com/users/octokitten/repos",
						"events_url": "https://api.github.com/users/octokitten/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octokitten/received_events",
						"type": "User",
						"site_admin": false
					}
				}
			]
		}`)
	})

	tmp, err := time.Parse(time.RFC3339, "2021-08-03T18:00:00-06:00")
	if err != nil {
		panic(err)
	}
	createdAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T15:00:00-06:00")
	if err != nil {
		panic(err)
	}
	updatedAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-10-14T00:53:32-06:00")
	if err != nil {
		panic(err)
	}
	lastActivityAt1 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T18:00:00-06:00")
	if err != nil {
		panic(err)
	}
	createdAt2 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T15:00:00-06:00")
	if err != nil {
		panic(err)
	}
	updatedAt2 := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-10-13T00:53:32-06:00")
	if err != nil {
		panic(err)
	}
	lastActivityAt2 := Timestamp{tmp}

	ctx := t.Context()
	opts := &ListOptions{Page: 1, PerPage: 100}
	got, _, err := client.Copilot.ListCopilotEnterpriseSeats(ctx, "e", opts)
	if err != nil {
		t.Errorf("Copilot.ListCopilotEnterpriseSeats returned error: %v", err)
	}

	want := &ListCopilotSeatsResponse{
		TotalSeats: 2,
		Seats: []*CopilotSeatDetails{
			{
				Assignee: &User{
					Login:             Ptr("octocat"),
					ID:                Ptr(int64(1)),
					NodeID:            Ptr("MDQ6VXNlcjE="),
					AvatarURL:         Ptr("https://github.com/images/error/octocat_happy.gif"),
					GravatarID:        Ptr(""),
					URL:               Ptr("https://api.github.com/users/octocat"),
					HTMLURL:           Ptr("https://github.com/octocat"),
					FollowersURL:      Ptr("https://api.github.com/users/octocat/followers"),
					FollowingURL:      Ptr("https://api.github.com/users/octocat/following{/other_user}"),
					GistsURL:          Ptr("https://api.github.com/users/octocat/gists{/gist_id}"),
					StarredURL:        Ptr("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
					SubscriptionsURL:  Ptr("https://api.github.com/users/octocat/subscriptions"),
					OrganizationsURL:  Ptr("https://api.github.com/users/octocat/orgs"),
					ReposURL:          Ptr("https://api.github.com/users/octocat/repos"),
					EventsURL:         Ptr("https://api.github.com/users/octocat/events{/privacy}"),
					ReceivedEventsURL: Ptr("https://api.github.com/users/octocat/received_events"),
					Type:              Ptr("User"),
					SiteAdmin:         Ptr(false),
				},
				AssigningTeam: &Team{
					ID:                  Ptr(int64(1)),
					NodeID:              Ptr("MDQ6VGVhbTE="),
					URL:                 Ptr("https://api.github.com/teams/1"),
					HTMLURL:             Ptr("https://github.com/orgs/github/teams/justice-league"),
					Name:                Ptr("Justice League"),
					Slug:                Ptr("justice-league"),
					Description:         Ptr("A great team."),
					Privacy:             Ptr("closed"),
					NotificationSetting: Ptr("notifications_enabled"),
					Permission:          Ptr("admin"),
					MembersURL:          Ptr("https://api.github.com/teams/1/members{/member}"),
					RepositoriesURL:     Ptr("https://api.github.com/teams/1/repos"),
					Parent:              nil,
				},
				CreatedAt:               &createdAt1,
				UpdatedAt:               &updatedAt1,
				PendingCancellationDate: nil,
				LastActivityAt:          &lastActivityAt1,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
				PlanType:                Ptr("business"),
			},
			{
				Assignee: &User{
					Login:             Ptr("octokitten"),
					ID:                Ptr(int64(1)),
					NodeID:            Ptr("MDQ76VNlcjE="),
					AvatarURL:         Ptr("https://github.com/images/error/octokitten_happy.gif"),
					GravatarID:        Ptr(""),
					URL:               Ptr("https://api.github.com/users/octokitten"),
					HTMLURL:           Ptr("https://github.com/octokitten"),
					FollowersURL:      Ptr("https://api.github.com/users/octokitten/followers"),
					FollowingURL:      Ptr("https://api.github.com/users/octokitten/following{/other_user}"),
					GistsURL:          Ptr("https://api.github.com/users/octokitten/gists{/gist_id}"),
					StarredURL:        Ptr("https://api.github.com/users/octokitten/starred{/owner}{/repo}"),
					SubscriptionsURL:  Ptr("https://api.github.com/users/octokitten/subscriptions"),
					OrganizationsURL:  Ptr("https://api.github.com/users/octokitten/orgs"),
					ReposURL:          Ptr("https://api.github.com/users/octokitten/repos"),
					EventsURL:         Ptr("https://api.github.com/users/octokitten/events{/privacy}"),
					ReceivedEventsURL: Ptr("https://api.github.com/users/octokitten/received_events"),
					Type:              Ptr("User"),
					SiteAdmin:         Ptr(false),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: Ptr("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
				PlanType:                nil,
			},
		},
	}

	if !cmp.Equal(got, want) {
		log.Printf("got: %+v", got.Seats[1])
		log.Printf("want: %+v", want.Seats[1])
		t.Errorf("Copilot.ListCopilotEnterpriseSeats returned %+v, want %+v", got, want)
	}

	const methodName = "ListCopilotEnterpriseSeats"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.ListCopilotEnterpriseSeats(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.ListCopilotEnterpriseSeats(ctx, "e", opts)
		if got != nil {
			t.Errorf("Copilot.ListCopilotEnterpriseSeats returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_AddCopilotTeams(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"selected_teams":["team1","team2"]}`+"\n")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.AddCopilotTeams(ctx, "o", []string{"team1", "team2"})
	if err != nil {
		t.Errorf("Copilot.AddCopilotTeams returned error: %v", err)
	}

	want := &SeatAssignments{SeatsCreated: 2}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.AddCopilotTeams returned %+v, want %+v", got, want)
	}

	const methodName = "AddCopilotTeams"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.AddCopilotTeams(ctx, "\n", []string{"team1", "team2"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.AddCopilotTeams(ctx, "o", []string{"team1", "team2"})
		if got != nil {
			t.Errorf("Copilot.AddCopilotTeams returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_RemoveCopilotTeams(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"selected_teams":["team1","team2"]}`+"\n")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.RemoveCopilotTeams(ctx, "o", []string{"team1", "team2"})
	if err != nil {
		t.Errorf("Copilot.RemoveCopilotTeams returned error: %v", err)
	}

	want := &SeatCancellations{SeatsCancelled: 2}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.RemoveCopilotTeams returned %+v, want %+v", got, want)
	}

	const methodName = "RemoveCopilotTeams"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.RemoveCopilotTeams(ctx, "\n", []string{"team1", "team2"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.RemoveCopilotTeams(ctx, "o", []string{"team1", "team2"})
		if got != nil {
			t.Errorf("Copilot.RemoveCopilotTeams returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_AddCopilotUsers(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"selected_usernames":["user1","user2"]}`+"\n")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.AddCopilotUsers(ctx, "o", []string{"user1", "user2"})
	if err != nil {
		t.Errorf("Copilot.AddCopilotUsers returned error: %v", err)
	}

	want := &SeatAssignments{SeatsCreated: 2}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.AddCopilotUsers returned %+v, want %+v", got, want)
	}

	const methodName = "AddCopilotUsers"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.AddCopilotUsers(ctx, "\n", []string{"user1", "user2"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.AddCopilotUsers(ctx, "o", []string{"user1", "user2"})
		if got != nil {
			t.Errorf("Copilot.AddCopilotUsers returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_RemoveCopilotUsers(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"selected_usernames":["user1","user2"]}`+"\n")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.RemoveCopilotUsers(ctx, "o", []string{"user1", "user2"})
	if err != nil {
		t.Errorf("Copilot.RemoveCopilotUsers returned error: %v", err)
	}

	want := &SeatCancellations{SeatsCancelled: 2}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.RemoveCopilotUsers returned %+v, want %+v", got, want)
	}

	const methodName = "RemoveCopilotUsers"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.RemoveCopilotUsers(ctx, "\n", []string{"user1", "user2"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.RemoveCopilotUsers(ctx, "o", []string{"user1", "user2"})
		if got != nil {
			t.Errorf("Copilot.RemoveCopilotUsers returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetSeatDetails(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/members/u/copilot", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"created_at": "2021-08-03T18:00:00-06:00",
				"updated_at": "2021-09-23T15:00:00-06:00",
				"pending_cancellation_date": null,
				"last_activity_at": "2021-10-14T00:53:32-06:00",
				"last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
				"assignee": {
					"login": "octocat",
					"id": 1,
					"node_id": "MDQ6VXNlcjE=",
					"avatar_url": "https://github.com/images/error/octocat_happy.gif",
					"gravatar_id": "",
					"url": "https://api.github.com/users/octocat",
					"html_url": "https://github.com/octocat",
					"followers_url": "https://api.github.com/users/octocat/followers",
					"following_url": "https://api.github.com/users/octocat/following{/other_user}",
					"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
					"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
					"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
					"organizations_url": "https://api.github.com/users/octocat/orgs",
					"repos_url": "https://api.github.com/users/octocat/repos",
					"events_url": "https://api.github.com/users/octocat/events{/privacy}",
					"received_events_url": "https://api.github.com/users/octocat/received_events",
					"type": "User",
					"site_admin": false
				},
				"assigning_team": {
					"id": 1,
					"node_id": "MDQ6VGVhbTE=",
					"url": "https://api.github.com/teams/1",
					"html_url": "https://github.com/orgs/github/teams/justice-league",
					"name": "Justice League",
					"slug": "justice-league",
					"description": "A great team.",
					"privacy": "closed",
					"notification_setting": "notifications_enabled",
					"permission": "admin",
					"members_url": "https://api.github.com/teams/1/members{/member}",
					"repositories_url": "https://api.github.com/teams/1/repos",
					"parent": null
				}
			}`)
	})

	tmp, err := time.Parse(time.RFC3339, "2021-08-03T18:00:00-06:00")
	if err != nil {
		panic(err)
	}
	createdAt := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-09-23T15:00:00-06:00")
	if err != nil {
		panic(err)
	}
	updatedAt := Timestamp{tmp}

	tmp, err = time.Parse(time.RFC3339, "2021-10-14T00:53:32-06:00")
	if err != nil {
		panic(err)
	}
	lastActivityAt := Timestamp{tmp}

	ctx := t.Context()
	got, _, err := client.Copilot.GetSeatDetails(ctx, "o", "u")
	if err != nil {
		t.Errorf("Copilot.GetSeatDetails returned error: %v", err)
	}

	want := &CopilotSeatDetails{
		Assignee: &User{
			Login:             Ptr("octocat"),
			ID:                Ptr(int64(1)),
			NodeID:            Ptr("MDQ6VXNlcjE="),
			AvatarURL:         Ptr("https://github.com/images/error/octocat_happy.gif"),
			GravatarID:        Ptr(""),
			URL:               Ptr("https://api.github.com/users/octocat"),
			HTMLURL:           Ptr("https://github.com/octocat"),
			FollowersURL:      Ptr("https://api.github.com/users/octocat/followers"),
			FollowingURL:      Ptr("https://api.github.com/users/octocat/following{/other_user}"),
			GistsURL:          Ptr("https://api.github.com/users/octocat/gists{/gist_id}"),
			StarredURL:        Ptr("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
			SubscriptionsURL:  Ptr("https://api.github.com/users/octocat/subscriptions"),
			OrganizationsURL:  Ptr("https://api.github.com/users/octocat/orgs"),
			ReposURL:          Ptr("https://api.github.com/users/octocat/repos"),
			EventsURL:         Ptr("https://api.github.com/users/octocat/events{/privacy}"),
			ReceivedEventsURL: Ptr("https://api.github.com/users/octocat/received_events"),
			Type:              Ptr("User"),
			SiteAdmin:         Ptr(false),
		},
		AssigningTeam: &Team{
			ID:                  Ptr(int64(1)),
			NodeID:              Ptr("MDQ6VGVhbTE="),
			URL:                 Ptr("https://api.github.com/teams/1"),
			HTMLURL:             Ptr("https://github.com/orgs/github/teams/justice-league"),
			Name:                Ptr("Justice League"),
			Slug:                Ptr("justice-league"),
			Description:         Ptr("A great team."),
			Privacy:             Ptr("closed"),
			NotificationSetting: Ptr("notifications_enabled"),
			Permission:          Ptr("admin"),
			MembersURL:          Ptr("https://api.github.com/teams/1/members{/member}"),
			RepositoriesURL:     Ptr("https://api.github.com/teams/1/repos"),
			Parent:              nil,
		},
		CreatedAt:               &createdAt,
		UpdatedAt:               &updatedAt,
		PendingCancellationDate: nil,
		LastActivityAt:          &lastActivityAt,
		LastActivityEditor:      Ptr("vscode/1.77.3/copilot/1.86.82"),
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetSeatDetails returned %+v, want %+v", got, want)
	}

	const methodName = "GetSeatDetails"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetSeatDetails(ctx, "\n", "u")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetSeatDetails(ctx, "o", "u")
		if got != nil {
			t.Errorf("Copilot.GetSeatDetails returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseMetrics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"date": "2024-06-24",
				"total_active_users": 24,
				"total_engaged_users": 20,
				"copilot_ide_code_completions": {
				"total_engaged_users": 20,
				"languages": [
					{
					"name": "python",
					"total_engaged_users": 10
					},
					{
					"name": "ruby",
					"total_engaged_users": 10
					}
				],
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 13,
						"languages": [
							{
							"name": "python",
							"total_engaged_users": 6,
							"total_code_suggestions": 249,
							"total_code_acceptances": 123,
							"total_code_lines_suggested": 225,
							"total_code_lines_accepted": 135
							},
							{
							"name": "ruby",
							"total_engaged_users": 7,
							"total_code_suggestions": 496,
							"total_code_acceptances": 253,
							"total_code_lines_suggested": 520,
							"total_code_lines_accepted": 270
							}
						]
						}
					]
					},
					{
					"name": "neovim",
					"total_engaged_users": 7,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"languages": [
							{
							"name": "typescript",
							"total_engaged_users": 3,
							"total_code_suggestions": 112,
							"total_code_acceptances": 56,
							"total_code_lines_suggested": 143,
							"total_code_lines_accepted": 61
							},
							{
							"name": "go",
							"total_engaged_users": 4,
							"total_code_suggestions": 132,
							"total_code_acceptances": 67,
							"total_code_lines_suggested": 154,
							"total_code_lines_accepted": 72
							}
						]
						}
					]
					}
				]
				},
				"copilot_ide_chat": {
				"total_engaged_users": 13,
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 12,
						"total_chats": 45,
						"total_chat_insertion_events": 12,
						"total_chat_copy_events": 16
						},
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_engaged_users": 1,
						"total_chats": 10,
						"total_chat_insertion_events": 11,
						"total_chat_copy_events": 3
						}
					]
					}
				]
				},
				"copilot_dotcom_chat": {
				"total_engaged_users": 14,
				"models": [
					{
					"name": "default",
					"is_custom_model": false,
					"custom_model_training_date": null,
					"total_engaged_users": 14,
					"total_chats": 38
					}
				]
				},
				"copilot_dotcom_pull_requests": {
				"total_engaged_users": 12,
				"repositories": [
					{
					"name": "demo/repo1",
					"total_engaged_users": 8,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_pr_summaries_created": 6,
						"total_engaged_users": 8
						}
					]
					},
					{
					"name": "demo/repo2",
					"total_engaged_users": 4,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_pr_summaries_created": 10,
						"total_engaged_users": 4
						}
					]
					}
				]
				}
			}
		]`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseMetrics(ctx, "e", &CopilotMetricsListOptions{})
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseMetrics returned error: %v", err)
	}

	totalActiveUsers := 24
	totalEngagedUsers := 20
	want := []*CopilotMetrics{
		{
			Date:              "2024-06-24",
			TotalActiveUsers:  &totalActiveUsers,
			TotalEngagedUsers: &totalEngagedUsers,
			CopilotIDECodeCompletions: &CopilotIDECodeCompletions{
				TotalEngagedUsers: 20,
				Languages: []*CopilotIDECodeCompletionsLanguage{
					{
						Name:              "python",
						TotalEngagedUsers: 10,
					},
					{
						Name:              "ruby",
						TotalEngagedUsers: 10,
					},
				},
				Editors: []*CopilotIDECodeCompletionsEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalEngagedUsers:       13,
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "python",
										TotalEngagedUsers:       6,
										TotalCodeSuggestions:    249,
										TotalCodeAcceptances:    123,
										TotalCodeLinesSuggested: 225,
										TotalCodeLinesAccepted:  135,
									},
									{
										Name:                    "ruby",
										TotalEngagedUsers:       7,
										TotalCodeSuggestions:    496,
										TotalCodeAcceptances:    253,
										TotalCodeLinesSuggested: 520,
										TotalCodeLinesAccepted:  270,
									},
								},
							},
						},
					},
					{
						Name:              "neovim",
						TotalEngagedUsers: 7,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "typescript",
										TotalEngagedUsers:       3,
										TotalCodeSuggestions:    112,
										TotalCodeAcceptances:    56,
										TotalCodeLinesSuggested: 143,
										TotalCodeLinesAccepted:  61,
									},
									{
										Name:                    "go",
										TotalEngagedUsers:       4,
										TotalCodeSuggestions:    132,
										TotalCodeAcceptances:    67,
										TotalCodeLinesSuggested: 154,
										TotalCodeLinesAccepted:  72,
									},
								},
							},
						},
					},
				},
			},
			CopilotIDEChat: &CopilotIDEChat{
				TotalEngagedUsers: 13,
				Editors: []*CopilotIDEChatEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDEChatModel{
							{
								Name:                     "default",
								IsCustomModel:            false,
								CustomModelTrainingDate:  nil,
								TotalEngagedUsers:        12,
								TotalChats:               45,
								TotalChatInsertionEvents: 12,
								TotalChatCopyEvents:      16,
							},
							{
								Name:                     "a-custom-model",
								IsCustomModel:            true,
								CustomModelTrainingDate:  Ptr("2024-02-01"),
								TotalEngagedUsers:        1,
								TotalChats:               10,
								TotalChatInsertionEvents: 11,
								TotalChatCopyEvents:      3,
							},
						},
					},
				},
			},
			CopilotDotcomChat: &CopilotDotcomChat{
				TotalEngagedUsers: 14,
				Models: []*CopilotDotcomChatModel{
					{
						Name:                    "default",
						IsCustomModel:           false,
						CustomModelTrainingDate: nil,
						TotalEngagedUsers:       14,
						TotalChats:              38,
					},
				},
			},
			CopilotDotcomPullRequests: &CopilotDotcomPullRequests{
				TotalEngagedUsers: 12,
				Repositories: []*CopilotDotcomPullRequestsRepository{
					{
						Name:              "demo/repo1",
						TotalEngagedUsers: 8,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalPRSummariesCreated: 6,
								TotalEngagedUsers:       8,
							},
						},
					},
					{
						Name:              "demo/repo2",
						TotalEngagedUsers: 4,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								TotalPRSummariesCreated: 10,
								TotalEngagedUsers:       4,
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseMetrics returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseMetrics"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseMetrics(ctx, "\n", &CopilotMetricsListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseMetrics(ctx, "e", &CopilotMetricsListOptions{})
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseMetrics returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseTeamMetrics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/team/t/copilot/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"date": "2024-06-24",
				"total_active_users": 24,
				"total_engaged_users": 20,
				"copilot_ide_code_completions": {
				"total_engaged_users": 20,
				"languages": [
					{
					"name": "python",
					"total_engaged_users": 10
					},
					{
					"name": "ruby",
					"total_engaged_users": 10
					}
				],
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 13,
						"languages": [
							{
							"name": "python",
							"total_engaged_users": 6,
							"total_code_suggestions": 249,
							"total_code_acceptances": 123,
							"total_code_lines_suggested": 225,
							"total_code_lines_accepted": 135
							},
							{
							"name": "ruby",
							"total_engaged_users": 7,
							"total_code_suggestions": 496,
							"total_code_acceptances": 253,
							"total_code_lines_suggested": 520,
							"total_code_lines_accepted": 270
							}
						]
						}
					]
					},
					{
					"name": "neovim",
					"total_engaged_users": 7,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"languages": [
							{
							"name": "typescript",
							"total_engaged_users": 3,
							"total_code_suggestions": 112,
							"total_code_acceptances": 56,
							"total_code_lines_suggested": 143,
							"total_code_lines_accepted": 61
							},
							{
							"name": "go",
							"total_engaged_users": 4,
							"total_code_suggestions": 132,
							"total_code_acceptances": 67,
							"total_code_lines_suggested": 154,
							"total_code_lines_accepted": 72
							}
						]
						}
					]
					}
				]
				},
				"copilot_ide_chat": {
				"total_engaged_users": 13,
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 12,
						"total_chats": 45,
						"total_chat_insertion_events": 12,
						"total_chat_copy_events": 16
						},
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_engaged_users": 1,
						"total_chats": 10,
						"total_chat_insertion_events": 11,
						"total_chat_copy_events": 3
						}
					]
					}
				]
				},
				"copilot_dotcom_chat": {
				"total_engaged_users": 14,
				"models": [
					{
					"name": "default",
					"is_custom_model": false,
					"custom_model_training_date": null,
					"total_engaged_users": 14,
					"total_chats": 38
					}
				]
				},
				"copilot_dotcom_pull_requests": {
				"total_engaged_users": 12,
				"repositories": [
					{
					"name": "demo/repo1",
					"total_engaged_users": 8,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_pr_summaries_created": 6,
						"total_engaged_users": 8
						}
					]
					},
					{
					"name": "demo/repo2",
					"total_engaged_users": 4,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_pr_summaries_created": 10,
						"total_engaged_users": 4
						}
					]
					}
				]
				}
			}
		]`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseTeamMetrics(ctx, "e", "t", &CopilotMetricsListOptions{})
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseTeamMetrics returned error: %v", err)
	}

	totalActiveUsers := 24
	totalEngagedUsers := 20
	want := []*CopilotMetrics{
		{
			Date:              "2024-06-24",
			TotalActiveUsers:  &totalActiveUsers,
			TotalEngagedUsers: &totalEngagedUsers,
			CopilotIDECodeCompletions: &CopilotIDECodeCompletions{
				TotalEngagedUsers: 20,
				Languages: []*CopilotIDECodeCompletionsLanguage{
					{
						Name:              "python",
						TotalEngagedUsers: 10,
					},
					{
						Name:              "ruby",
						TotalEngagedUsers: 10,
					},
				},
				Editors: []*CopilotIDECodeCompletionsEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalEngagedUsers:       13,
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "python",
										TotalEngagedUsers:       6,
										TotalCodeSuggestions:    249,
										TotalCodeAcceptances:    123,
										TotalCodeLinesSuggested: 225,
										TotalCodeLinesAccepted:  135,
									},
									{
										Name:                    "ruby",
										TotalEngagedUsers:       7,
										TotalCodeSuggestions:    496,
										TotalCodeAcceptances:    253,
										TotalCodeLinesSuggested: 520,
										TotalCodeLinesAccepted:  270,
									},
								},
							},
						},
					},
					{
						Name:              "neovim",
						TotalEngagedUsers: 7,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "typescript",
										TotalEngagedUsers:       3,
										TotalCodeSuggestions:    112,
										TotalCodeAcceptances:    56,
										TotalCodeLinesSuggested: 143,
										TotalCodeLinesAccepted:  61,
									},
									{
										Name:                    "go",
										TotalEngagedUsers:       4,
										TotalCodeSuggestions:    132,
										TotalCodeAcceptances:    67,
										TotalCodeLinesSuggested: 154,
										TotalCodeLinesAccepted:  72,
									},
								},
							},
						},
					},
				},
			},
			CopilotIDEChat: &CopilotIDEChat{
				TotalEngagedUsers: 13,
				Editors: []*CopilotIDEChatEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDEChatModel{
							{
								Name:                     "default",
								IsCustomModel:            false,
								CustomModelTrainingDate:  nil,
								TotalEngagedUsers:        12,
								TotalChats:               45,
								TotalChatInsertionEvents: 12,
								TotalChatCopyEvents:      16,
							},
							{
								Name:                     "a-custom-model",
								IsCustomModel:            true,
								CustomModelTrainingDate:  Ptr("2024-02-01"),
								TotalEngagedUsers:        1,
								TotalChats:               10,
								TotalChatInsertionEvents: 11,
								TotalChatCopyEvents:      3,
							},
						},
					},
				},
			},
			CopilotDotcomChat: &CopilotDotcomChat{
				TotalEngagedUsers: 14,
				Models: []*CopilotDotcomChatModel{
					{
						Name:                    "default",
						IsCustomModel:           false,
						CustomModelTrainingDate: nil,
						TotalEngagedUsers:       14,
						TotalChats:              38,
					},
				},
			},
			CopilotDotcomPullRequests: &CopilotDotcomPullRequests{
				TotalEngagedUsers: 12,
				Repositories: []*CopilotDotcomPullRequestsRepository{
					{
						Name:              "demo/repo1",
						TotalEngagedUsers: 8,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalPRSummariesCreated: 6,
								TotalEngagedUsers:       8,
							},
						},
					},
					{
						Name:              "demo/repo2",
						TotalEngagedUsers: 4,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								TotalPRSummariesCreated: 10,
								TotalEngagedUsers:       4,
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseTeamMetrics returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseTeamMetrics"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseTeamMetrics(ctx, "\n", "t", &CopilotMetricsListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseTeamMetrics(ctx, "e", "t", &CopilotMetricsListOptions{})
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseTeamMetrics returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationMetrics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"date": "2024-06-24",
				"total_active_users": 24,
				"total_engaged_users": 20,
				"copilot_ide_code_completions": {
				"total_engaged_users": 20,
				"languages": [
					{
					"name": "python",
					"total_engaged_users": 10
					},
					{
					"name": "ruby",
					"total_engaged_users": 10
					}
				],
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 13,
						"languages": [
							{
							"name": "python",
							"total_engaged_users": 6,
							"total_code_suggestions": 249,
							"total_code_acceptances": 123,
							"total_code_lines_suggested": 225,
							"total_code_lines_accepted": 135
							},
							{
							"name": "ruby",
							"total_engaged_users": 7,
							"total_code_suggestions": 496,
							"total_code_acceptances": 253,
							"total_code_lines_suggested": 520,
							"total_code_lines_accepted": 270
							}
						]
						}
					]
					},
					{
					"name": "neovim",
					"total_engaged_users": 7,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"languages": [
							{
							"name": "typescript",
							"total_engaged_users": 3,
							"total_code_suggestions": 112,
							"total_code_acceptances": 56,
							"total_code_lines_suggested": 143,
							"total_code_lines_accepted": 61
							},
							{
							"name": "go",
							"total_engaged_users": 4,
							"total_code_suggestions": 132,
							"total_code_acceptances": 67,
							"total_code_lines_suggested": 154,
							"total_code_lines_accepted": 72
							}
						]
						}
					]
					}
				]
				},
				"copilot_ide_chat": {
				"total_engaged_users": 13,
				"editors": [
					{
					"name": "vscode",
					"total_engaged_users": 13,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_engaged_users": 12,
						"total_chats": 45,
						"total_chat_insertion_events": 12,
						"total_chat_copy_events": 16
						},
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_engaged_users": 1,
						"total_chats": 10,
						"total_chat_insertion_events": 11,
						"total_chat_copy_events": 3
						}
					]
					}
				]
				},
				"copilot_dotcom_chat": {
				"total_engaged_users": 14,
				"models": [
					{
					"name": "default",
					"is_custom_model": false,
					"custom_model_training_date": null,
					"total_engaged_users": 14,
					"total_chats": 38
					}
				]
				},
				"copilot_dotcom_pull_requests": {
				"total_engaged_users": 12,
				"repositories": [
					{
					"name": "demo/repo1",
					"total_engaged_users": 8,
					"models": [
						{
						"name": "default",
						"is_custom_model": false,
						"custom_model_training_date": null,
						"total_pr_summaries_created": 6,
						"total_engaged_users": 8
						}
					]
					},
					{
					"name": "demo/repo2",
					"total_engaged_users": 4,
					"models": [
						{
						"name": "a-custom-model",
						"is_custom_model": true,
						"custom_model_training_date": "2024-02-01",
						"total_pr_summaries_created": 10,
						"total_engaged_users": 4
						}
					]
					}
				]
				}
			}
		]`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetOrganizationMetrics(ctx, "o", &CopilotMetricsListOptions{})
	if err != nil {
		t.Errorf("Copilot.GetOrganizationMetrics returned error: %v", err)
	}

	totalActiveUsers := 24
	totalEngagedUsers := 20
	want := []*CopilotMetrics{
		{
			Date:              "2024-06-24",
			TotalActiveUsers:  &totalActiveUsers,
			TotalEngagedUsers: &totalEngagedUsers,
			CopilotIDECodeCompletions: &CopilotIDECodeCompletions{
				TotalEngagedUsers: 20,
				Languages: []*CopilotIDECodeCompletionsLanguage{
					{
						Name:              "python",
						TotalEngagedUsers: 10,
					},
					{
						Name:              "ruby",
						TotalEngagedUsers: 10,
					},
				},
				Editors: []*CopilotIDECodeCompletionsEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalEngagedUsers:       13,
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "python",
										TotalEngagedUsers:       6,
										TotalCodeSuggestions:    249,
										TotalCodeAcceptances:    123,
										TotalCodeLinesSuggested: 225,
										TotalCodeLinesAccepted:  135,
									},
									{
										Name:                    "ruby",
										TotalEngagedUsers:       7,
										TotalCodeSuggestions:    496,
										TotalCodeAcceptances:    253,
										TotalCodeLinesSuggested: 520,
										TotalCodeLinesAccepted:  270,
									},
								},
							},
						},
					},
					{
						Name:              "neovim",
						TotalEngagedUsers: 7,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "typescript",
										TotalEngagedUsers:       3,
										TotalCodeSuggestions:    112,
										TotalCodeAcceptances:    56,
										TotalCodeLinesSuggested: 143,
										TotalCodeLinesAccepted:  61,
									},
									{
										Name:                    "go",
										TotalEngagedUsers:       4,
										TotalCodeSuggestions:    132,
										TotalCodeAcceptances:    67,
										TotalCodeLinesSuggested: 154,
										TotalCodeLinesAccepted:  72,
									},
								},
							},
						},
					},
				},
			},
			CopilotIDEChat: &CopilotIDEChat{
				TotalEngagedUsers: 13,
				Editors: []*CopilotIDEChatEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDEChatModel{
							{
								Name:                     "default",
								IsCustomModel:            false,
								CustomModelTrainingDate:  nil,
								TotalEngagedUsers:        12,
								TotalChats:               45,
								TotalChatInsertionEvents: 12,
								TotalChatCopyEvents:      16,
							},
							{
								Name:                     "a-custom-model",
								IsCustomModel:            true,
								CustomModelTrainingDate:  Ptr("2024-02-01"),
								TotalEngagedUsers:        1,
								TotalChats:               10,
								TotalChatInsertionEvents: 11,
								TotalChatCopyEvents:      3,
							},
						},
					},
				},
			},
			CopilotDotcomChat: &CopilotDotcomChat{
				TotalEngagedUsers: 14,
				Models: []*CopilotDotcomChatModel{
					{
						Name:                    "default",
						IsCustomModel:           false,
						CustomModelTrainingDate: nil,
						TotalEngagedUsers:       14,
						TotalChats:              38,
					},
				},
			},
			CopilotDotcomPullRequests: &CopilotDotcomPullRequests{
				TotalEngagedUsers: 12,
				Repositories: []*CopilotDotcomPullRequestsRepository{
					{
						Name:              "demo/repo1",
						TotalEngagedUsers: 8,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalPRSummariesCreated: 6,
								TotalEngagedUsers:       8,
							},
						},
					},
					{
						Name:              "demo/repo2",
						TotalEngagedUsers: 4,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								TotalPRSummariesCreated: 10,
								TotalEngagedUsers:       4,
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationMetrics returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationMetrics"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationMetrics(ctx, "\n", &CopilotMetricsListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationMetrics(ctx, "o", &CopilotMetricsListOptions{})
		if got != nil {
			t.Errorf("Copilot.GetOrganizationMetrics returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationTeamMetrics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/team/t/copilot/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"date": "2024-06-24",
				"total_active_users": 24,
				"total_engaged_users": 20,
				"copilot_ide_code_completions": {
					"total_engaged_users": 20,
					"languages": [
						{
							"name": "python",
							"total_engaged_users": 10
						},
						{
							"name": "ruby",
							"total_engaged_users": 10
						}
					],
					"editors": [
						{
							"name": "vscode",
							"total_engaged_users": 13,
							"models": [
								{
									"name": "default",
									"is_custom_model": false,
									"custom_model_training_date": null,
									"total_engaged_users": 13,
									"languages": [
										{
											"name": "python",
											"total_engaged_users": 6,
											"total_code_suggestions": 249,
											"total_code_acceptances": 123,
											"total_code_lines_suggested": 225,
											"total_code_lines_accepted": 135
										},
										{
											"name": "ruby",
											"total_engaged_users": 7,
											"total_code_suggestions": 496,
											"total_code_acceptances": 253,
											"total_code_lines_suggested": 520,
											"total_code_lines_accepted": 270
										}
									]
								}
							]
						},
						{
							"name": "neovim",
							"total_engaged_users": 7,
							"models": [
								{
									"name": "a-custom-model",
									"is_custom_model": true,
									"custom_model_training_date": "2024-02-01",
									"languages": [
										{
											"name": "typescript",
											"total_engaged_users": 3,
											"total_code_suggestions": 112,
											"total_code_acceptances": 56,
											"total_code_lines_suggested": 143,
											"total_code_lines_accepted": 61
										},
										{
											"name": "go",
											"total_engaged_users": 4,
											"total_code_suggestions": 132,
											"total_code_acceptances": 67,
											"total_code_lines_suggested": 154,
											"total_code_lines_accepted": 72
										}
									]
								}
							]
						}
					]
				},
				"copilot_ide_chat": {
					"total_engaged_users": 13,
					"editors": [
						{
							"name": "vscode",
							"total_engaged_users": 13,
							"models": [
								{
									"name": "default",
									"is_custom_model": false,
									"custom_model_training_date": null,
									"total_engaged_users": 12,
									"total_chats": 45,
									"total_chat_insertion_events": 12,
									"total_chat_copy_events": 16
								},
								{
									"name": "a-custom-model",
									"is_custom_model": true,
									"custom_model_training_date": "2024-02-01",
									"total_engaged_users": 1,
									"total_chats": 10,
									"total_chat_insertion_events": 11,
									"total_chat_copy_events": 3
								}
							]
						}
					]
				},
				"copilot_dotcom_chat": {
					"total_engaged_users": 14,
					"models": [
						{
							"name": "default",
							"is_custom_model": false,
							"custom_model_training_date": null,
							"total_engaged_users": 14,
							"total_chats": 38
						}
					]
				},
				"copilot_dotcom_pull_requests": {
					"total_engaged_users": 12,
					"repositories": [
						{
							"name": "demo/repo1",
							"total_engaged_users": 8,
							"models": [
								{
									"name": "default",
									"is_custom_model": false,
									"custom_model_training_date": null,
									"total_pr_summaries_created": 6,
									"total_engaged_users": 8
								}
							]
						},
						{
							"name": "demo/repo2",
							"total_engaged_users": 4,
							"models": [
								{
									"name": "a-custom-model",
									"is_custom_model": true,
									"custom_model_training_date": "2024-02-01",
									"total_pr_summaries_created": 10,
									"total_engaged_users": 4
								}
							]
						}
					]
				}
			}
		]`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetOrganizationTeamMetrics(ctx, "o", "t", &CopilotMetricsListOptions{})
	if err != nil {
		t.Errorf("Copilot.GetOrganizationTeamMetrics returned error: %v", err)
	}

	totalActiveUsers := 24
	totalEngagedUsers := 20
	want := []*CopilotMetrics{
		{
			Date:              "2024-06-24",
			TotalActiveUsers:  &totalActiveUsers,
			TotalEngagedUsers: &totalEngagedUsers,
			CopilotIDECodeCompletions: &CopilotIDECodeCompletions{
				TotalEngagedUsers: 20,
				Languages: []*CopilotIDECodeCompletionsLanguage{
					{
						Name:              "python",
						TotalEngagedUsers: 10,
					},
					{
						Name:              "ruby",
						TotalEngagedUsers: 10,
					},
				},
				Editors: []*CopilotIDECodeCompletionsEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalEngagedUsers:       13,
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "python",
										TotalEngagedUsers:       6,
										TotalCodeSuggestions:    249,
										TotalCodeAcceptances:    123,
										TotalCodeLinesSuggested: 225,
										TotalCodeLinesAccepted:  135,
									},
									{
										Name:                    "ruby",
										TotalEngagedUsers:       7,
										TotalCodeSuggestions:    496,
										TotalCodeAcceptances:    253,
										TotalCodeLinesSuggested: 520,
										TotalCodeLinesAccepted:  270,
									},
								},
							},
						},
					},
					{
						Name:              "neovim",
						TotalEngagedUsers: 7,
						Models: []*CopilotIDECodeCompletionsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								Languages: []*CopilotIDECodeCompletionsModelLanguage{
									{
										Name:                    "typescript",
										TotalEngagedUsers:       3,
										TotalCodeSuggestions:    112,
										TotalCodeAcceptances:    56,
										TotalCodeLinesSuggested: 143,
										TotalCodeLinesAccepted:  61,
									},
									{
										Name:                    "go",
										TotalEngagedUsers:       4,
										TotalCodeSuggestions:    132,
										TotalCodeAcceptances:    67,
										TotalCodeLinesSuggested: 154,
										TotalCodeLinesAccepted:  72,
									},
								},
							},
						},
					},
				},
			},
			CopilotIDEChat: &CopilotIDEChat{
				TotalEngagedUsers: 13,
				Editors: []*CopilotIDEChatEditor{
					{
						Name:              "vscode",
						TotalEngagedUsers: 13,
						Models: []*CopilotIDEChatModel{
							{
								Name:                     "default",
								IsCustomModel:            false,
								CustomModelTrainingDate:  nil,
								TotalEngagedUsers:        12,
								TotalChats:               45,
								TotalChatInsertionEvents: 12,
								TotalChatCopyEvents:      16,
							},
							{
								Name:                     "a-custom-model",
								IsCustomModel:            true,
								CustomModelTrainingDate:  Ptr("2024-02-01"),
								TotalEngagedUsers:        1,
								TotalChats:               10,
								TotalChatInsertionEvents: 11,
								TotalChatCopyEvents:      3,
							},
						},
					},
				},
			},
			CopilotDotcomChat: &CopilotDotcomChat{
				TotalEngagedUsers: 14,
				Models: []*CopilotDotcomChatModel{
					{
						Name:                    "default",
						IsCustomModel:           false,
						CustomModelTrainingDate: nil,
						TotalEngagedUsers:       14,
						TotalChats:              38,
					},
				},
			},
			CopilotDotcomPullRequests: &CopilotDotcomPullRequests{
				TotalEngagedUsers: 12,
				Repositories: []*CopilotDotcomPullRequestsRepository{
					{
						Name:              "demo/repo1",
						TotalEngagedUsers: 8,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "default",
								IsCustomModel:           false,
								CustomModelTrainingDate: nil,
								TotalPRSummariesCreated: 6,
								TotalEngagedUsers:       8,
							},
						},
					},
					{
						Name:              "demo/repo2",
						TotalEngagedUsers: 4,
						Models: []*CopilotDotcomPullRequestsModel{
							{
								Name:                    "a-custom-model",
								IsCustomModel:           true,
								CustomModelTrainingDate: Ptr("2024-02-01"),
								TotalPRSummariesCreated: 10,
								TotalEngagedUsers:       4,
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationTeamMetrics returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationTeamMetrics"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationTeamMetrics(ctx, "\n", "\n", &CopilotMetricsListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationTeamMetrics(ctx, "o", "t", &CopilotMetricsListOptions{})
		if got != nil {
			t.Errorf("Copilot.GetOrganizationTeamMetrics returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseDailyMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/metrics/reports/enterprise-1-day", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"day": "2025-07-01"})
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_day": "2025-07-01"
		}`)
	})

	ctx := t.Context()
	opts := &CopilotMetricsReportOptions{Day: "2025-07-01"}
	got, _, err := client.Copilot.GetEnterpriseDailyMetricsReport(ctx, "e", opts)
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseDailyMetricsReport returned error: %v", err)
	}

	want := &CopilotDailyMetricsReport{
		DownloadLinks: []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportDay:     "2025-07-01",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseDailyMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseDailyMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseDailyMetricsReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseDailyMetricsReport(ctx, "e", opts)
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseDailyMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/metrics/reports/enterprise-28-day/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_start_day": "2025-07-01",
			"report_end_day": "2025-07-28"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseMetricsReport(ctx, "e")
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseMetricsReport returned error: %v", err)
	}

	want := &CopilotMetricsReport{
		DownloadLinks:  []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportStartDay: "2025-07-01",
		ReportEndDay:   "2025-07-28",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseMetricsReport(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseMetricsReport(ctx, "e")
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseUsersDailyMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/metrics/reports/users-1-day", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"day": "2025-07-01"})
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_day": "2025-07-01"
		}`)
	})

	ctx := t.Context()
	opts := &CopilotMetricsReportOptions{Day: "2025-07-01"}
	got, _, err := client.Copilot.GetEnterpriseUsersDailyMetricsReport(ctx, "e", opts)
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseUsersDailyMetricsReport returned error: %v", err)
	}

	want := &CopilotDailyMetricsReport{
		DownloadLinks: []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportDay:     "2025-07-01",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseUsersDailyMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseUsersDailyMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseUsersDailyMetricsReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseUsersDailyMetricsReport(ctx, "e", opts)
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseUsersDailyMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseUsersMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/metrics/reports/users-28-day/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_start_day": "2025-07-01",
			"report_end_day": "2025-07-28"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseUsersMetricsReport(ctx, "e")
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseUsersMetricsReport returned error: %v", err)
	}

	want := &CopilotMetricsReport{
		DownloadLinks:  []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportStartDay: "2025-07-01",
		ReportEndDay:   "2025-07-28",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseUsersMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseUsersMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseUsersMetricsReport(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseUsersMetricsReport(ctx, "e")
		if got != nil {
			t.Errorf("Copilot.GetEnterpriseUsersMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationDailyMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/metrics/reports/organization-1-day", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"day": "2025-07-01"})
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_day": "2025-07-01"
		}`)
	})

	ctx := t.Context()
	opts := &CopilotMetricsReportOptions{Day: "2025-07-01"}
	got, _, err := client.Copilot.GetOrganizationDailyMetricsReport(ctx, "o", opts)
	if err != nil {
		t.Errorf("Copilot.GetOrganizationDailyMetricsReport returned error: %v", err)
	}

	want := &CopilotDailyMetricsReport{
		DownloadLinks: []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportDay:     "2025-07-01",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationDailyMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationDailyMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationDailyMetricsReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationDailyMetricsReport(ctx, "o", opts)
		if got != nil {
			t.Errorf("Copilot.GetOrganizationDailyMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/metrics/reports/organization-28-day/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_start_day": "2025-07-01",
			"report_end_day": "2025-07-28"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetOrganizationMetricsReport(ctx, "o")
	if err != nil {
		t.Errorf("Copilot.GetOrganizationMetricsReport returned error: %v", err)
	}

	want := &CopilotMetricsReport{
		DownloadLinks:  []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportStartDay: "2025-07-01",
		ReportEndDay:   "2025-07-28",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationMetricsReport(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationMetricsReport(ctx, "o")
		if got != nil {
			t.Errorf("Copilot.GetOrganizationMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationUsersDailyMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/metrics/reports/users-1-day", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"day": "2025-07-01"})
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_day": "2025-07-01"
		}`)
	})

	ctx := t.Context()
	opts := &CopilotMetricsReportOptions{Day: "2025-07-01"}
	got, _, err := client.Copilot.GetOrganizationUsersDailyMetricsReport(ctx, "o", opts)
	if err != nil {
		t.Errorf("Copilot.GetOrganizationUsersDailyMetricsReport returned error: %v", err)
	}

	want := &CopilotDailyMetricsReport{
		DownloadLinks: []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportDay:     "2025-07-01",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationUsersDailyMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationUsersDailyMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationUsersDailyMetricsReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationUsersDailyMetricsReport(ctx, "o", opts)
		if got != nil {
			t.Errorf("Copilot.GetOrganizationUsersDailyMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_GetOrganizationUsersMetricsReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/copilot/metrics/reports/users-28-day/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"download_links": ["https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"],
			"report_start_day": "2025-07-01",
			"report_end_day": "2025-07-28"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetOrganizationUsersMetricsReport(ctx, "o")
	if err != nil {
		t.Errorf("Copilot.GetOrganizationUsersMetricsReport returned error: %v", err)
	}

	want := &CopilotMetricsReport{
		DownloadLinks:  []string{"https://example.com/copilot-usage-report-1.json", "https://example.com/copilot-usage-report-2.json"},
		ReportStartDay: "2025-07-01",
		ReportEndDay:   "2025-07-28",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetOrganizationUsersMetricsReport returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationUsersMetricsReport"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetOrganizationUsersMetricsReport(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetOrganizationUsersMetricsReport(ctx, "o")
		if got != nil {
			t.Errorf("Copilot.GetOrganizationUsersMetricsReport returned %+v, want nil", got)
		}
		return resp, err
	})
}
