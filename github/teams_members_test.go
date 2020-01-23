// Copyright 2018 The go-github AUTHORS. All rights reserved.
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
	"time"
)

func TestTeamsService__ListTeamMembersByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, _, err := client.Teams.ListTeamMembersByID(context.Background(), 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersByID returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamMembersByID returned %+v, want %+v", members, want)
	}
}

func TestTeamsService__ListTeamMembersByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, _, err := client.Teams.ListTeamMembersByName(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersByName returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamMembers returned %+v, want %+v", members, want)
	}
}

func TestTeamsService__IsTeamMember_true(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	member, _, err := client.Teams.IsTeamMember(context.Background(), 1, "u")
	if err != nil {
		t.Errorf("Teams.IsTeamMember returned error: %v", err)
	}
	if want := true; member != want {
		t.Errorf("Teams.IsTeamMember returned %+v, want %+v", member, want)
	}
}

// ensure that a 404 response is interpreted as "false" and not an error
func TestTeamsService__IsTeamMember_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	member, _, err := client.Teams.IsTeamMember(context.Background(), 1, "u")
	if err != nil {
		t.Errorf("Teams.IsTeamMember returned error: %+v", err)
	}
	if want := false; member != want {
		t.Errorf("Teams.IsTeamMember returned %+v, want %+v", member, want)
	}
}

// ensure that a 400 response is interpreted as an actual error, and not simply
// as "false" like the above case of a 404
func TestTeamsService__IsTeamMember_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	member, _, err := client.Teams.IsTeamMember(context.Background(), 1, "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; member != want {
		t.Errorf("Teams.IsTeamMember returned %+v, want %+v", member, want)
	}
}

func TestTeamsService__IsTeamMember_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.IsTeamMember(context.Background(), 1, "%")
	testURLParseError(t, err)
}

func TestTeamsService__GetTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	membership, _, err := client.Teams.GetTeamMembershipByID(context.Background(), 1, 1, "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.GetTeamMembershipByID returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__GetTeamMembershipByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	membership, _, err := client.Teams.GetTeamMembershipByName(context.Background(), "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipByName returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.GetTeamMembershipByName returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__AddTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/organizations/1/team/1/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	membership, _, err := client.Teams.AddTeamMembershipByID(context.Background(), 1, 1, "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.AddTeamMembershipByID returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__AddTeamMembershipByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	membership, _, err := client.Teams.AddTeamMembershipByName(context.Background(), "o", "s", "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipByName returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.AddTeamMembershipByName returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__RemoveTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamMembershipByID(context.Background(), 1, 1, "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned error: %v", err)
	}
}

func TestTeamsService__RemoveTeamMembershipByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamMembershipByName(context.Background(), "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipByName returned error: %v", err)
	}
}

func TestTeamsService_ListPendingTeamInvitationsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1"})
		fmt.Fprint(w, `[
				{
    					"id": 1,
    					"login": "monalisa",
    					"email": "octocat@github.com",
    					"role": "direct_member",
    					"created_at": "2017-01-21T00:00:00Z",
    					"inviter": {
      						"login": "other_user",
      						"id": 1,
      						"avatar_url": "https://github.com/images/error/other_user_happy.gif",
      						"gravatar_id": "",
      						"url": "https://api.github.com/users/other_user",
      						"html_url": "https://github.com/other_user",
      						"followers_url": "https://api.github.com/users/other_user/followers",
      						"following_url": "https://api.github.com/users/other_user/following/other_user",
      						"gists_url": "https://api.github.com/users/other_user/gists/gist_id",
      						"starred_url": "https://api.github.com/users/other_user/starred/owner/repo",
      						"subscriptions_url": "https://api.github.com/users/other_user/subscriptions",
      						"organizations_url": "https://api.github.com/users/other_user/orgs",
      						"repos_url": "https://api.github.com/users/other_user/repos",
      						"events_url": "https://api.github.com/users/other_user/events/privacy",
      						"received_events_url": "https://api.github.com/users/other_user/received_events/privacy",
      						"type": "User",
      						"site_admin": false
    					}
  				}
			]`)
	})

	opt := &ListOptions{Page: 1}
	invitations, _, err := client.Teams.ListPendingTeamInvitationsByID(context.Background(), 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned error: %v", err)
	}

	createdAt := time.Date(2017, time.January, 21, 0, 0, 0, 0, time.UTC)
	want := []*Invitation{
		{
			ID:        Int64(1),
			Login:     String("monalisa"),
			Email:     String("octocat@github.com"),
			Role:      String("direct_member"),
			CreatedAt: &createdAt,
			Inviter: &User{
				Login:             String("other_user"),
				ID:                Int64(1),
				AvatarURL:         String("https://github.com/images/error/other_user_happy.gif"),
				GravatarID:        String(""),
				URL:               String("https://api.github.com/users/other_user"),
				HTMLURL:           String("https://github.com/other_user"),
				FollowersURL:      String("https://api.github.com/users/other_user/followers"),
				FollowingURL:      String("https://api.github.com/users/other_user/following/other_user"),
				GistsURL:          String("https://api.github.com/users/other_user/gists/gist_id"),
				StarredURL:        String("https://api.github.com/users/other_user/starred/owner/repo"),
				SubscriptionsURL:  String("https://api.github.com/users/other_user/subscriptions"),
				OrganizationsURL:  String("https://api.github.com/users/other_user/orgs"),
				ReposURL:          String("https://api.github.com/users/other_user/repos"),
				EventsURL:         String("https://api.github.com/users/other_user/events/privacy"),
				ReceivedEventsURL: String("https://api.github.com/users/other_user/received_events/privacy"),
				Type:              String("User"),
				SiteAdmin:         Bool(false),
			},
		}}

	if !reflect.DeepEqual(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned %+v, want %+v", invitations, want)
	}
}

func TestTeamsService_ListPendingTeamInvitationsByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1"})
		fmt.Fprint(w, `[
				{
    					"id": 1,
    					"login": "monalisa",
    					"email": "octocat@github.com",
    					"role": "direct_member",
    					"created_at": "2017-01-21T00:00:00Z",
    					"inviter": {
      						"login": "other_user",
      						"id": 1,
      						"avatar_url": "https://github.com/images/error/other_user_happy.gif",
      						"gravatar_id": "",
      						"url": "https://api.github.com/users/other_user",
      						"html_url": "https://github.com/other_user",
      						"followers_url": "https://api.github.com/users/other_user/followers",
      						"following_url": "https://api.github.com/users/other_user/following/other_user",
      						"gists_url": "https://api.github.com/users/other_user/gists/gist_id",
      						"starred_url": "https://api.github.com/users/other_user/starred/owner/repo",
      						"subscriptions_url": "https://api.github.com/users/other_user/subscriptions",
      						"organizations_url": "https://api.github.com/users/other_user/orgs",
      						"repos_url": "https://api.github.com/users/other_user/repos",
      						"events_url": "https://api.github.com/users/other_user/events/privacy",
      						"received_events_url": "https://api.github.com/users/other_user/received_events/privacy",
      						"type": "User",
      						"site_admin": false
    					}
  				}
			]`)
	})

	opt := &ListOptions{Page: 1}
	invitations, _, err := client.Teams.ListPendingTeamInvitationsByName(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByName returned error: %v", err)
	}

	createdAt := time.Date(2017, time.January, 21, 0, 0, 0, 0, time.UTC)
	want := []*Invitation{
		{
			ID:        Int64(1),
			Login:     String("monalisa"),
			Email:     String("octocat@github.com"),
			Role:      String("direct_member"),
			CreatedAt: &createdAt,
			Inviter: &User{
				Login:             String("other_user"),
				ID:                Int64(1),
				AvatarURL:         String("https://github.com/images/error/other_user_happy.gif"),
				GravatarID:        String(""),
				URL:               String("https://api.github.com/users/other_user"),
				HTMLURL:           String("https://github.com/other_user"),
				FollowersURL:      String("https://api.github.com/users/other_user/followers"),
				FollowingURL:      String("https://api.github.com/users/other_user/following/other_user"),
				GistsURL:          String("https://api.github.com/users/other_user/gists/gist_id"),
				StarredURL:        String("https://api.github.com/users/other_user/starred/owner/repo"),
				SubscriptionsURL:  String("https://api.github.com/users/other_user/subscriptions"),
				OrganizationsURL:  String("https://api.github.com/users/other_user/orgs"),
				ReposURL:          String("https://api.github.com/users/other_user/repos"),
				EventsURL:         String("https://api.github.com/users/other_user/events/privacy"),
				ReceivedEventsURL: String("https://api.github.com/users/other_user/received_events/privacy"),
				Type:              String("User"),
				SiteAdmin:         Bool(false),
			},
		}}

	if !reflect.DeepEqual(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByName returned %+v, want %+v", invitations, want)
	}
}
