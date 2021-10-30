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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	ctx := context.Background()
	commentsByID, _, err := client.Teams.ListCommentsByID(ctx, 1, 2, 3,
		&DiscussionCommentListOptions{Direction: "desc"})
	if err != nil {
		t.Errorf("Teams.ListCommentsByID returned error: %v", err)
	}

	if !cmp.Equal(commentsByID, want) {
		t.Errorf("Teams.ListCommentsByID returned %+v, want %+v", commentsByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "")
	mux.HandleFunc(e, handleFunc)

	commentsBySlug, _, err := client.Teams.ListCommentsBySlug(ctx, "a", "b", 3,
		&DiscussionCommentListOptions{Direction: "desc"})
	if err != nil {
		t.Errorf("Teams.ListCommentsBySlug returned error: %v", err)
	}

	if !cmp.Equal(commentsBySlug, want) {
		t.Errorf("Teams.ListCommentsBySlug returned %+v, want %+v", commentsBySlug, want)
	}

	methodName := "ListCommentsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListCommentsByID(ctx, -1, -2, -3,
			&DiscussionCommentListOptions{Direction: "desc"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListCommentsByID(ctx, 1, 2, 3,
			&DiscussionCommentListOptions{Direction: "desc"})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	methodName = "ListCommentsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListCommentsBySlug(ctx, "a\na", "b\nb", -3,
			&DiscussionCommentListOptions{Direction: "desc"})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListCommentsBySlug(ctx, "a", "b", 3,
			&DiscussionCommentListOptions{Direction: "desc"})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	commentByID, _, err := client.Teams.GetCommentByID(ctx, 1, 2, 3, 4)
	if err != nil {
		t.Errorf("Teams.GetCommentByID returned error: %v", err)
	}

	if !cmp.Equal(commentByID, want) {
		t.Errorf("Teams.GetCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.GetCommentBySlug(ctx, "a", "b", 3, 4)
	if err != nil {
		t.Errorf("Teams.GetCommentBySlug returned error: %v", err)
	}

	if !cmp.Equal(commentBySlug, want) {
		t.Errorf("Teams.GetCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}

	methodName := "GetCommentByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetCommentByID(ctx, -1, -2, -3, -4)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetCommentByID(ctx, 1, 2, 3, 4)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	methodName = "ListCommentsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetCommentBySlug(ctx, "a\na", "b\nb", -3, -4)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetCommentBySlug(ctx, "a", "b", 3, 4)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := DiscussionComment{Body: String("c")}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		v := new(DiscussionComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":4}`)
	}
	want := &DiscussionComment{Number: Int(4)}

	e := tdcEndpointByID("1", "2", "3", "")
	mux.HandleFunc(e, handlerFunc)

	ctx := context.Background()
	commentByID, _, err := client.Teams.CreateCommentByID(ctx, 1, 2, 3, input)
	if err != nil {
		t.Errorf("Teams.CreateCommentByID returned error: %v", err)
	}

	if !cmp.Equal(commentByID, want) {
		t.Errorf("Teams.CreateCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.CreateCommentBySlug(ctx, "a", "b", 3, input)
	if err != nil {
		t.Errorf("Teams.CreateCommentBySlug returned error: %v", err)
	}

	if !cmp.Equal(commentBySlug, want) {
		t.Errorf("Teams.CreateCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}

	methodName := "CreateCommentByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.CreateCommentByID(ctx, -1, -2, -3, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.CreateCommentByID(ctx, 1, 2, 3, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	methodName = "CreateCommentBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.CreateCommentBySlug(ctx, "a\na", "b\nb", -3, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.CreateCommentBySlug(ctx, "a", "b", 3, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_EditComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := DiscussionComment{Body: String("e")}
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		v := new(DiscussionComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":4}`)
	}
	want := &DiscussionComment{Number: Int(4)}

	e := tdcEndpointByID("1", "2", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	ctx := context.Background()
	commentByID, _, err := client.Teams.EditCommentByID(ctx, 1, 2, 3, 4, input)
	if err != nil {
		t.Errorf("Teams.EditCommentByID returned error: %v", err)
	}

	if !cmp.Equal(commentByID, want) {
		t.Errorf("Teams.EditCommentByID returned %+v, want %+v", commentByID, want)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	commentBySlug, _, err := client.Teams.EditCommentBySlug(ctx, "a", "b", 3, 4, input)
	if err != nil {
		t.Errorf("Teams.EditCommentBySlug returned error: %v", err)
	}

	if !cmp.Equal(commentBySlug, want) {
		t.Errorf("Teams.EditCommentBySlug returned %+v, want %+v", commentBySlug, want)
	}

	methodName := "EditCommentByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.EditCommentByID(ctx, -1, -2, -3, -4, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.EditCommentByID(ctx, 1, 2, 3, 4, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	methodName = "EditCommentBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.EditCommentBySlug(ctx, "a\na", "b\nb", -3, -4, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.EditCommentBySlug(ctx, "a", "b", 3, 4, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	}

	e := tdcEndpointByID("1", "2", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	ctx := context.Background()
	_, err := client.Teams.DeleteCommentByID(ctx, 1, 2, 3, 4)
	if err != nil {
		t.Errorf("Teams.DeleteCommentByID returned error: %v", err)
	}

	e = tdcEndpointBySlug("a", "b", "3", "4")
	mux.HandleFunc(e, handlerFunc)

	_, err = client.Teams.DeleteCommentBySlug(ctx, "a", "b", 3, 4)
	if err != nil {
		t.Errorf("Teams.DeleteCommentBySlug returned error: %v", err)
	}

	methodName := "DeleteCommentByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.DeleteCommentByID(ctx, -1, -2, -3, -4)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.Teams.DeleteCommentByID(ctx, 1, 2, 3, 4)
		return resp, err
	})

	methodName = "DeleteCommentBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.DeleteCommentBySlug(ctx, "a\na", "b\nb", -3, -4)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.Teams.DeleteCommentBySlug(ctx, "a", "b", 3, 4)
		return resp, err
	})
}

func TestDiscussionComment_Marshal(t *testing.T) {
	testJSONMarshal(t, &DiscussionComment{}, "{}")

	u := &DiscussionComment{
		Author:        &User{},
		Body:          String("body"),
		BodyHTML:      String("body html"),
		BodyVersion:   String("body version"),
		CreatedAt:     &Timestamp{referenceTime},
		LastEditedAt:  &Timestamp{referenceTime},
		DiscussionURL: String("url"),
		HTMLURL:       String("html url"),
		NodeID:        String("node"),
		Number:        Int(1),
		UpdatedAt:     &Timestamp{referenceTime},
		URL:           String("url"),
		Reactions: &Reactions{
			TotalCount: Int(10),
			PlusOne:    Int(1),
			MinusOne:   Int(1),
			Laugh:      Int(1),
			Confused:   Int(1),
			Heart:      Int(2),
			Hooray:     Int(5),
			Rocket:     Int(3),
			Eyes:       Int(9),
			URL:        String("url"),
		},
	}

	want := `{
		"author":{},
		"body":"body",
		"body_html":"body html",
		"body_version":"body version",
		"created_at":` + referenceTimeStr + `,
		"last_edited_at":` + referenceTimeStr + `,
		"discussion_url":"url",
		"html_url":"html url",
		"node_id":"node",
		"number":1,
		"updated_at":` + referenceTimeStr + `,
		"url":"url",
		"reactions":{
			"total_count": 10,
			"+1": 1,
			"-1": 1,
			"laugh": 1,
			"confused": 1,
			"heart": 2,
			"hooray": 5,
			"rocket": 3,
			"eyes": 9,
			"url":"url"
		}
	}`

	testJSONMarshal(t, u, want)
}
