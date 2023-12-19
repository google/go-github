// Copyright 2023 The go-github AUTHORS. All rights reserved.
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
	"time"

	"github.com/google/go-cmp/cmp"
)

// Test invalid JSON responses, vlaid responses are covered in the other tests
func TestCopilotSeatDetails_UnmarshalJSON(t *testing.T) {
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
			err := json.Unmarshal([]byte(tc.data), seatDetails)
			if err == nil && tc.wantErr {
				t.Errorf("CopilotSeatDetails.UnmarshalJSON returned nil instead of an error")
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
		ID:   Int64(1),
		Type: String("User"),
	}

	if got, ok := seatDetails.GetUser(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetTeam returned %+v, want %+v", got, want)
	} else if !ok {
		t.Errorf("CopilotSeatDetails.GetUser returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "Organization",
					"id": 1
				}
			}`

	bad := &Organization{
		ID:   Int64(1),
		Type: String("Organization"),
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
		ID: Int64(1),
	}

	if got, ok := seatDetails.GetTeam(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetTeam returned %+v, want %+v", got, want)
	} else if !ok {
		t.Errorf("CopilotSeatDetails.GetTeam returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "User",
					"id": 1
				}
			}`

	bad := &User{
		ID:   Int64(1),
		Type: String("User"),
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
		ID:   Int64(1),
		Type: String("Organization"),
	}

	if got, ok := seatDetails.GetOrganization(); ok && !cmp.Equal(got, want) {
		t.Errorf("CopilotSeatDetails.GetOrganization returned %+v, want %+v", got, want)
	} else if !ok {
		t.Errorf("CopilotSeatDetails.GetOrganization returned false, expected true")
	}

	data = `{
				"assignee": {
					"type": "Team",
					"id": 1
				}
			}`

	bad := &Team{
		ID: Int64(1),
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/seats", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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

	ctx := context.Background()
	got, _, err := client.Copilot.ListCopilotSeats(ctx, "o", nil)
	if err != nil {
		t.Errorf("Copilot.ListCopilotSeats returned error: %v", err)
	}

	want := &ListCopilotSeatsResponse{
		TotalSeats: 4,
		Seats: []*CopilotSeatDetails{
			{
				Assignee: &User{
					Login:             String("octocat"),
					ID:                Int64(1),
					NodeID:            String("MDQ6VXNlcjE="),
					AvatarURL:         String("https://github.com/images/error/octocat_happy.gif"),
					GravatarID:        String(""),
					URL:               String("https://api.github.com/users/octocat"),
					HTMLURL:           String("https://github.com/octocat"),
					FollowersURL:      String("https://api.github.com/users/octocat/followers"),
					FollowingURL:      String("https://api.github.com/users/octocat/following{/other_user}"),
					GistsURL:          String("https://api.github.com/users/octocat/gists{/gist_id}"),
					StarredURL:        String("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
					SubscriptionsURL:  String("https://api.github.com/users/octocat/subscriptions"),
					OrganizationsURL:  String("https://api.github.com/users/octocat/orgs"),
					ReposURL:          String("https://api.github.com/users/octocat/repos"),
					EventsURL:         String("https://api.github.com/users/octocat/events{/privacy}"),
					ReceivedEventsURL: String("https://api.github.com/users/octocat/received_events"),
					Type:              String("User"),
					SiteAdmin:         Bool(false),
				},
				AssigningTeam: &Team{
					ID:              Int64(1),
					NodeID:          String("MDQ6VGVhbTE="),
					URL:             String("https://api.github.com/teams/1"),
					HTMLURL:         String("https://github.com/orgs/github/teams/justice-league"),
					Name:            String("Justice League"),
					Slug:            String("justice-league"),
					Description:     String("A great team."),
					Privacy:         String("closed"),
					Permission:      String("admin"),
					MembersURL:      String("https://api.github.com/teams/1/members{/member}"),
					RepositoriesURL: String("https://api.github.com/teams/1/repos"),
					Parent:          nil,
				},
				CreatedAt:               &createdAt1,
				UpdatedAt:               &updatedAt1,
				PendingCancellationDate: nil,
				LastActivityAt:          &lastActivityAt1,
				LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &User{
					Login:             String("octokitten"),
					ID:                Int64(1),
					NodeID:            String("MDQ76VNlcjE="),
					AvatarURL:         String("https://github.com/images/error/octokitten_happy.gif"),
					GravatarID:        String(""),
					URL:               String("https://api.github.com/users/octokitten"),
					HTMLURL:           String("https://github.com/octokitten"),
					FollowersURL:      String("https://api.github.com/users/octokitten/followers"),
					FollowingURL:      String("https://api.github.com/users/octokitten/following{/other_user}"),
					GistsURL:          String("https://api.github.com/users/octokitten/gists{/gist_id}"),
					StarredURL:        String("https://api.github.com/users/octokitten/starred{/owner}{/repo}"),
					SubscriptionsURL:  String("https://api.github.com/users/octokitten/subscriptions"),
					OrganizationsURL:  String("https://api.github.com/users/octokitten/orgs"),
					ReposURL:          String("https://api.github.com/users/octokitten/repos"),
					EventsURL:         String("https://api.github.com/users/octokitten/events{/privacy}"),
					ReceivedEventsURL: String("https://api.github.com/users/octokitten/received_events"),
					Type:              String("User"),
					SiteAdmin:         Bool(false),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: String("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &Team{
					ID:   Int64(1),
					Name: String("octokittens"),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: String("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
			},
			{
				Assignee: &Organization{
					ID:   Int64(1),
					Name: String("octocats"),
					Type: String("Organization"),
				},
				AssigningTeam:           nil,
				CreatedAt:               &createdAt2,
				UpdatedAt:               &updatedAt2,
				PendingCancellationDate: String("2021-11-01"),
				LastActivityAt:          &lastActivityAt2,
				LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.ListCopilotSeats returned %+v, want %+v", got, want)
	}

	const methodName = "ListCopilotSeats"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.ListCopilotSeats(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.ListCopilotSeats(ctx, "", nil)
		if got != nil {
			t.Errorf("Copilot.ListCopilotSeats returned %+v, want nil", got)
		}
		return resp, err
	})
}

func TestCopilotService_AddCopilotTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"selected_teams":["team1","team2"]}`+"\n")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"selected_teams":["team1","team2"]}`+"\n")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"selected_usernames":["user1","user2"]}`+"\n")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"selected_usernames":["user1","user2"]}`+"\n")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	got, _, err := client.Copilot.GetSeatDetails(ctx, "o", "u")
	if err != nil {
		t.Errorf("Copilot.GetSeatDetails returned error: %v", err)
	}

	want := &CopilotSeatDetails{
		Assignee: &User{
			Login:             String("octocat"),
			ID:                Int64(1),
			NodeID:            String("MDQ6VXNlcjE="),
			AvatarURL:         String("https://github.com/images/error/octocat_happy.gif"),
			GravatarID:        String(""),
			URL:               String("https://api.github.com/users/octocat"),
			HTMLURL:           String("https://github.com/octocat"),
			FollowersURL:      String("https://api.github.com/users/octocat/followers"),
			FollowingURL:      String("https://api.github.com/users/octocat/following{/other_user}"),
			GistsURL:          String("https://api.github.com/users/octocat/gists{/gist_id}"),
			StarredURL:        String("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
			SubscriptionsURL:  String("https://api.github.com/users/octocat/subscriptions"),
			OrganizationsURL:  String("https://api.github.com/users/octocat/orgs"),
			ReposURL:          String("https://api.github.com/users/octocat/repos"),
			EventsURL:         String("https://api.github.com/users/octocat/events{/privacy}"),
			ReceivedEventsURL: String("https://api.github.com/users/octocat/received_events"),
			Type:              String("User"),
			SiteAdmin:         Bool(false),
		},
		AssigningTeam: &Team{
			ID:              Int64(1),
			NodeID:          String("MDQ6VGVhbTE="),
			URL:             String("https://api.github.com/teams/1"),
			HTMLURL:         String("https://github.com/orgs/github/teams/justice-league"),
			Name:            String("Justice League"),
			Slug:            String("justice-league"),
			Description:     String("A great team."),
			Privacy:         String("closed"),
			Permission:      String("admin"),
			MembersURL:      String("https://api.github.com/teams/1/members{/member}"),
			RepositoriesURL: String("https://api.github.com/teams/1/repos"),
			Parent:          nil,
		},
		CreatedAt:               &createdAt,
		UpdatedAt:               &updatedAt,
		PendingCancellationDate: nil,
		LastActivityAt:          &lastActivityAt,
		LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
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
