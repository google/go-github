// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestClassroom_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Classroom{}, "{}")

	c := &Classroom{
		ID:       Ptr(int64(1296269)),
		Name:     Ptr("Programming Elixir"),
		Archived: Ptr(false),
		Organization: &Organization{
			ID:        Ptr(int64(1)),
			Login:     Ptr("programming-elixir"),
			NodeID:    Ptr("MDEyOk9yZ2FuaXphdGlvbjE="),
			HTMLURL:   Ptr("https://github.com/programming-elixir"),
			Name:      Ptr("Learn how to build fault tolerant applications"),
			AvatarURL: Ptr("https://avatars.githubusercontent.com/u/9919?v=4"),
		},
		URL: Ptr("https://classroom.github.com/classrooms/1-programming-elixir"),
	}

	want := `{
		"id": 1296269,
		"name": "Programming Elixir",
		"archived": false,
		"organization": {
			"id": 1,
			"login": "programming-elixir",
			"node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
			"html_url": "https://github.com/programming-elixir",
			"name": "Learn how to build fault tolerant applications",
			"avatar_url": "https://avatars.githubusercontent.com/u/9919?v=4"
		},
		"url": "https://classroom.github.com/classrooms/1-programming-elixir"
	}`

	testJSONMarshal(t, c, want)
}

func TestClassroomAssignment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ClassroomAssignment{}, "{}")

	a := &ClassroomAssignment{
		ID:                          Ptr(int64(12)),
		PublicRepo:                  Ptr(false),
		Title:                       Ptr("Intro to Binaries"),
		Type:                        Ptr("individual"),
		InviteLink:                  Ptr("https://classroom.github.com/a/Lx7jiUgx"),
		InvitationsEnabled:          Ptr(true),
		Slug:                        Ptr("intro-to-binaries"),
		StudentsAreRepoAdmins:       Ptr(false),
		FeedbackPullRequestsEnabled: Ptr(true),
		MaxTeams:                    Ptr(0),
		MaxMembers:                  Ptr(0),
		Editor:                      Ptr("codespaces"),
		Accepted:                    Ptr(100),
		Submitted:                   Ptr(40),
		Passing:                     Ptr(10),
		Language:                    Ptr("ruby"),
		Deadline:                    &Timestamp{referenceTime},
		StarterCodeRepository: &Repository{
			ID:       Ptr(int64(1296269)),
			FullName: Ptr("octocat/Hello-World"),
		},
		Classroom: &Classroom{
			ID:   Ptr(int64(1296269)),
			Name: Ptr("Programming Elixir"),
		},
	}

	want := `{
		"id": 12,
		"public_repo": false,
		"title": "Intro to Binaries",
		"type": "individual",
		"invite_link": "https://classroom.github.com/a/Lx7jiUgx",
		"invitations_enabled": true,
		"slug": "intro-to-binaries",
		"students_are_repo_admins": false,
		"feedback_pull_requests_enabled": true,
		"max_teams": 0,
		"max_members": 0,
		"editor": "codespaces",
		"accepted": 100,
		"submitted": 40,
		"passing": 10,
		"language": "ruby",
		"deadline": ` + referenceTimeStr + `,
		"starter_code_repository": {
			"id": 1296269,
			"full_name": "octocat/Hello-World"
		},
		"classroom": {
			"id": 1296269,
			"name": "Programming Elixir"
		}
	}`

	testJSONMarshal(t, a, want)
}

func TestClassroomService_GetAssignment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/assignments/12", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 12,
			"public_repo": false,
			"title": "Intro to Binaries",
			"type": "individual",
			"invite_link": "https://classroom.github.com/a/Lx7jiUgx",
			"invitations_enabled": true,
			"slug": "intro-to-binaries",
			"students_are_repo_admins": false,
			"feedback_pull_requests_enabled": true,
			"max_teams": 0,
			"max_members": 0,
			"editor": "codespaces",
			"accepted": 100,
			"submitted": 40,
			"passing": 10,
			"language": "ruby",
			"deadline": "2011-01-26T19:06:43Z",
			"starter_code_repository": {
				"id": 1296269,
				"full_name": "octocat/Hello-World",
				"html_url": "https://github.com/octocat/Hello-World",
				"node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
				"private": false,
				"default_branch": "main"
			},
			"classroom": {
				"id": 1296269,
				"name": "Programming Elixir",
				"archived": false,
				"url": "https://classroom.github.com/classrooms/1-programming-elixir"
			}
		}`)
	})

	ctx := context.Background()
	assignment, _, err := client.Classroom.GetAssignment(ctx, 12)
	if err != nil {
		t.Errorf("Classroom.GetAssignment returned error: %v", err)
	}

	want := &ClassroomAssignment{
		ID:                          Ptr(int64(12)),
		PublicRepo:                  Ptr(false),
		Title:                       Ptr("Intro to Binaries"),
		Type:                        Ptr("individual"),
		InviteLink:                  Ptr("https://classroom.github.com/a/Lx7jiUgx"),
		InvitationsEnabled:          Ptr(true),
		Slug:                        Ptr("intro-to-binaries"),
		StudentsAreRepoAdmins:       Ptr(false),
		FeedbackPullRequestsEnabled: Ptr(true),
		MaxTeams:                    Ptr(0),
		MaxMembers:                  Ptr(0),
		Editor:                      Ptr("codespaces"),
		Accepted:                    Ptr(100),
		Submitted:                   Ptr(40),
		Passing:                     Ptr(10),
		Language:                    Ptr("ruby"),
		Deadline:                    func() *Timestamp { t, _ := time.Parse(time.RFC3339, "2011-01-26T19:06:43Z"); return &Timestamp{t} }(),
		StarterCodeRepository: &Repository{
			ID:            Ptr(int64(1296269)),
			FullName:      Ptr("octocat/Hello-World"),
			HTMLURL:       Ptr("https://github.com/octocat/Hello-World"),
			NodeID:        Ptr("MDEwOlJlcG9zaXRvcnkxMjk2MjY5"),
			Private:       Ptr(false),
			DefaultBranch: Ptr("main"),
		},
		Classroom: &Classroom{
			ID:       Ptr(int64(1296269)),
			Name:     Ptr("Programming Elixir"),
			Archived: Ptr(false),
			URL:      Ptr("https://classroom.github.com/classrooms/1-programming-elixir"),
		},
	}

	if !cmp.Equal(assignment, want) {
		t.Errorf("Classroom.GetAssignment returned %+v, want %+v", assignment, want)
	}

	const methodName = "GetAssignment"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.GetAssignment(ctx, 12)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
