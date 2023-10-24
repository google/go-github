package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

	want := &OrganizationCopilotDetails{
		SeatBreakdown: &SeatBreakdown{
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
}

func TestCopilotService_ListCopilotSeats(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/seats", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"total_seats": 2,
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
					}
				]
			}`)
	})

	ctx := context.Background()
	got, _, err := client.Copilot.ListCopilotSeats(ctx, "o")
	if err != nil {
		t.Errorf("Copilot.ListCopilotSeats returned error: %v", err)
	}

	want := &CopilotSeats{
		TotalSeats: 2,
		Seats: []CopilotSeatDetails{
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
				CreatedAt:               "2021-08-03T18:00:00-06:00",
				UpdatedAt:               "2021-09-23T15:00:00-06:00",
				PendingCancellationDate: nil,
				LastActivityAt:          String("2021-10-14T00:53:32-06:00"),
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
				CreatedAt:               "2021-09-23T18:00:00-06:00",
				UpdatedAt:               "2021-09-23T15:00:00-06:00",
				PendingCancellationDate: String("2021-11-01"),
				LastActivityAt:          String("2021-10-13T00:53:32-06:00"),
				LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.ListCopilotSeats returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_AddCopilotTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := context.Background()
	got, _, err := client.Copilot.AddCopilotTeams(ctx, "o", SelectedTeams{SelectedTeams: []string{"team1", "team2"}})
	if err != nil {
		t.Errorf("Copilot.AddCopilotTeams returned error: %v", err)
	}

	want := &SeatAssignments{SeatsCreated: 2}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.AddCopilotTeams returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_RemoveCopilotTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := context.Background()
	got, _, err := client.Copilot.RemoveCopilotTeams(ctx, "o", SelectedTeams{SelectedTeams: []string{"team1", "team2"}})
	if err != nil {
		t.Errorf("Copilot.RemoveCopilotTeams returned error: %v", err)
	}

	want := &SeatCancellations{SeatsCancelled: 2}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.RemoveCopilotTeams returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_AddCopilotUsers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"seats_created": 2}`)
	})

	ctx := context.Background()
	got, _, err := client.Copilot.AddCopilotUsers(ctx, "o", SelectedUsers{SelectedUsers: []string{"user1", "user2"}})
	if err != nil {
		t.Errorf("Copilot.AddCopilotUsers returned error: %v", err)
	}

	want := &SeatAssignments{SeatsCreated: 2}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.AddCopilotUsers returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_RemoveCopilotUsers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/copilot/billing/selected_users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"seats_cancelled": 2}`)
	})

	ctx := context.Background()
	got, _, err := client.Copilot.RemoveCopilotUsers(ctx, "o", SelectedUsers{SelectedUsers: []string{"user1", "user2"}})
	if err != nil {
		t.Errorf("Copilot.RemoveCopilotUsers returned error: %v", err)
	}

	want := &SeatCancellations{SeatsCancelled: 2}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.RemoveCopilotUsers returned %+v, want %+v", got, want)
	}
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
		CreatedAt:               "2021-08-03T18:00:00-06:00",
		UpdatedAt:               "2021-09-23T15:00:00-06:00",
		PendingCancellationDate: nil,
		LastActivityAt:          String("2021-10-14T00:53:32-06:00"),
		LastActivityEditor:      String("vscode/1.77.3/copilot/1.86.82"),
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetSeatDetails returned %+v, want %+v", got, want)
	}
}
