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

// "Team Discussion Comments" endpoint, when using a teamID.
func tdcEndpointByID(orgID, teamID, discussionNumber, commentNumber string) string {
	out := fmt.Sprintf("/organizations/%v/team/%v/discussions/%v/comments", orgID, teamID, discussionNumber)
	if commentNumber != "" {
		return fmt.Sprintf("%v/%v", out, commentNumber)
	}
	return out
}

// "Team Discussion Comments" endpoint, when using a team slug.
func tdcEndpointBySlug(org, slug, dicsuccionsNumber, commentNumber string) string {
	out := fmt.Sprintf("/orgs/%v/teams/%v/discussions/%v/comments", org, slug, dicsuccionsNumber)
	if commentNumber != "" {
		return fmt.Sprintf("%v/%v", out, commentNumber)
	}
	return out
}

func TestTeamsService_ListComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"direction": "desc",
		})
		fmt.Fprintf(w,
			`[
				{
					"author": {
						"login": "author",
						"id": 0,
						"avatar_url": "https://avatars1.githubusercontent.com/u/0?v=4",
						"gravatar_id": "",
						"url": "https://api.github.com/users/author",
						"html_url": "https://github.com/author",
						"followers_url": "https://api.github.com/users/author/followers",
						"following_url": "https://api.github.com/users/author/following{/other_user}",
						"gists_url": "https://api.github.com/users/author/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/author/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/author/subscriptions",
						"organizations_url": "https://api.github.com/users/author/orgs",
						"repos_url": "https://api.github.com/users/author/repos",
						"events_url": "https://api.github.com/users/author/events{/privacy}",
						"received_events_url": "https://api.github.com/users/author/received_events",
						"type": "User",
						"site_admin": false
					},
					"body": "comment",
					"body_html": "<p>comment</p>",
					"body_version": "version",
					"created_at": "2018-01-01T00:00:00Z",
					"last_edited_at": null,
					"discussion_url": "https://api.github.com/teams/2/discussions/3",
					"html_url": "https://github.com/orgs/1/teams/2/discussions/3/comments/4",
					"node_id": "node",
					"number": 4,
					"updated_at": "2018-01-01T00:00:00Z",
					"url": "https://api.github.com/teams/2/discussions/3/comments/4"
				}
			]`)
	}

	want := []*DiscussionComment{
		{
			Author: &User{
				Login:             String("author"),
				ID:                Int64(0),
				AvatarURL:         String("https://avatars1.githubusercontent.com/u/0?v=4"),
				GravatarID:        String(""),
				URL:               String("https://api.github.com/users/author"),
				HTMLURL:           String("https://github.com/author"),
				FollowersURL:      String("https://api.github.com/users/author/followers"),
				FollowingURL:      String("https://api.github.com/users/author/following{/other_user}"),
				GistsURL:          String("https://api.github.com/users/author/gists{/gist_id}"),
				StarredURL:        String("https://api.github.com/users/author/starred{/owner}{/repo}"),
				SubscriptionsURL:  String("https://api.github.com/users/author/subscriptions"),
				OrganizationsURL:  String("https://api.github.com/users/author/orgs"),
				ReposURL:          String("https://api.github.com/users/author/repos"),
				EventsURL:         String("https://api.github.com/users/author/events{/privacy}"),
				ReceivedEventsURL: String("https://api.github.com/users/author/received_events"),
				Type:              String("User"),
				SiteAdmin:         Bool(false),
			},
			Body:          String("comment"),
			BodyHTML:      String("<p>comment</p>"),
			BodyVersion:   String("version"),
			CreatedAt:     &Timestamp{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)},
			LastEditedAt:  nil,
			DiscussionURL: String("https://api.github.com/teams/2/discussions/3"),
			HTMLURL:       String("https://github.com/orgs/1/teams/2/discussions/3/comments/4"),
			NodeID:        String("node"),
			Number:        Int(4),
			UpdatedAt:     &Timestamp{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)},
			URL:           String("https://api.github.com/teams/2/discussions/3/comments/4"),
		},
	}

	e := tdcEndpointByID("1", "2", "3", "")
	mux.HandleFunc(e, handleFunc)

	commentsByID, _, err := client.Teams.ListCommentsByID(context.Background(), 1, 2, 3,
		&DiscussionCommentListOptions{Direction: "desc"})
	if err != nil {
		t.Errorf("Teams.ListCommentsByID returned error: %v", err)
	}

	if !reflect.DeepEqual(commentsByID, want) {
		t.Errorf("Teams.ListCommentsByID returned %+v, want %+v", commentsByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "")
	mux.HandleFunc(e, handleFunc)

	commentsBySlug, _, err := client.Teams.ListCommentsBySlug(context.Background(), "a", "b", 3,
		&DiscussionCommentListOptions{Direction: "desc"})
	if err != nil {
		t.Errorf("Teams.ListCommentsBySlug returned error: %v", err)
	}

	if !reflect.DeepEqual(commentsBySlug, want) {
		t.Errorf("Teams.ListCommentsBySlug returned %+v, want %+v", commentsBySlug, want)
	}
}

func TestTeamsService_GetComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":4}`)
	}
	want := &DiscussionComment{Number: Int(4)}

	e := tdcEndpointByID("1", "2", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentByID, _, err := client.Teams.GetCommentByID(context.Background(), 1, 2, 3, 4)
	if err != nil {
		t.Errorf("Teams.GetCommentByID returned error: %v", err)
	}

	if !reflect.DeepEqual(commentByID, want) {
		t.Errorf("Teams.GetCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.GetCommentBySlug(context.Background(), "a", "b", 3, 4)
	if err != nil {
		t.Errorf("Teams.GetCommentBySlug returned error: %v", err)
	}

	if !reflect.DeepEqual(commentBySlug, want) {
		t.Errorf("Teams.GetCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}

}

func TestTeamsService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := DiscussionComment{Body: String("c")}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		v := new(DiscussionComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":4}`)
	}
	want := &DiscussionComment{Number: Int(4)}

	e := tdcEndpointByID("1", "2", "3", "")
	mux.HandleFunc(e, handlerFunc)

	commentByID, _, err := client.Teams.CreateCommentByID(context.Background(), 1, 2, 3, input)
	if err != nil {
		t.Errorf("Teams.CreateCommentByID returned error: %v", err)
	}

	if !reflect.DeepEqual(commentByID, want) {
		t.Errorf("Teams.CreateCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.CreateCommentBySlug(context.Background(), "a", "b", 3, input)
	if err != nil {
		t.Errorf("Teams.CreateCommentBySlug returned error: %v", err)
	}

	if !reflect.DeepEqual(commentBySlug, want) {
		t.Errorf("Teams.CreateCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}
}

func TestTeamsService_EditComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := DiscussionComment{Body: String("e")}
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		v := new(DiscussionComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":4}`)
	}
	want := &DiscussionComment{Number: Int(4)}

	e := tdcEndpointByID("1", "2", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentByID, _, err := client.Teams.EditCommentByID(context.Background(), 1, 2, 3, 4, input)
	if err != nil {
		t.Errorf("Teams.EditCommentByID returned error: %v", err)
	}

	if !reflect.DeepEqual(commentByID, want) {
		t.Errorf("Teams.EditCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.EditCommentBySlug(context.Background(), "a", "b", 3, 4, input)
	if err != nil {
		t.Errorf("Teams.EditCommentBySlug returned error: %v", err)
	}

	if !reflect.DeepEqual(commentBySlug, want) {
		t.Errorf("Teams.EditCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}
}

func TestTeamsService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	}

	e := tdcEndpointByID("1", "2", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	_, err := client.Teams.DeleteCommentByID(context.Background(), 1, 2, 3, 4)
	if err != nil {
		t.Errorf("Teams.DeleteCommentByID returned error: %v", err)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	_, err = client.Teams.DeleteCommentBySlug(context.Background(), "a", "b", 3, 4)
	if err != nil {
		t.Errorf("Teams.DeleteCommentBySlug returned error: %v", err)
	}
}
