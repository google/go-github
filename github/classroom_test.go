// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
			"invite_link": "https://example.com/a/Lx7jiUgx",
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
			"deadline": `+referenceTimeStr+`,
			"starter_code_repository": {
				"id": 1296269,
				"full_name": "octocat/Hello-World",
				"html_url": "https://example.com/octocat/Hello-World",
				"node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
				"private": false,
				"default_branch": "main"
			},
			"classroom": {
				"id": 1296269,
				"name": "Programming Elixir",
				"archived": false,
				"url": "https://example.com/classrooms/programming"
			}
		}`)
	})

	ctx := t.Context()
	assignment, _, err := client.Classroom.GetAssignment(ctx, 12)
	if err != nil {
		t.Errorf("Classroom.GetAssignment returned error: %v", err)
	}

	want := &ClassroomAssignment{
		ID:                          new(int64(12)),
		PublicRepo:                  new(false),
		Title:                       new("Intro to Binaries"),
		Type:                        new("individual"),
		InviteLink:                  new("https://example.com/a/Lx7jiUgx"),
		InvitationsEnabled:          new(true),
		Slug:                        new("intro-to-binaries"),
		StudentsAreRepoAdmins:       new(false),
		FeedbackPullRequestsEnabled: new(true),
		MaxTeams:                    new(0),
		MaxMembers:                  new(0),
		Editor:                      new("codespaces"),
		Accepted:                    new(100),
		Submitted:                   new(40),
		Passing:                     new(10),
		Language:                    new("ruby"),
		Deadline:                    &referenceTimestamp,
		StarterCodeRepository: &Repository{
			ID:            new(int64(1296269)),
			FullName:      new("octocat/Hello-World"),
			HTMLURL:       new("https://example.com/octocat/Hello-World"),
			NodeID:        new("MDEwOlJlcG9zaXRvcnkxMjk2MjY5"),
			Private:       new(false),
			DefaultBranch: new("main"),
		},
		Classroom: &Classroom{
			ID:       new(int64(1296269)),
			Name:     new("Programming Elixir"),
			Archived: new(false),
			URL:      new("https://example.com/classrooms/programming"),
		},
	}

	if !cmp.Equal(assignment, want) {
		t.Errorf("Classroom.GetAssignment returned %+v, want %+v", assignment, want)
	}

	const methodName = "GetAssignment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Classroom.GetAssignment(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.GetAssignment(ctx, 12)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestClassroomService_GetClassroom(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/classrooms/1296269", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 1296269,
			"name": "Programming Elixir",
			"archived": false,
			"organization": {
				"id": 1,
				"login": "programming-elixir",
				"node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
				"html_url": "https://example.com/programming-elixir",
				"name": "Learn how to build fault tolerant applications",
				"avatar_url": "https://example.com/avatars/u/9919?v=4"
			},
			"url": "https://example.com/classrooms/programming"
		}`)
	})

	ctx := t.Context()
	classroom, _, err := client.Classroom.GetClassroom(ctx, 1296269)
	if err != nil {
		t.Errorf("Classroom.GetClassroom returned error: %v", err)
	}

	want := &Classroom{
		ID:       new(int64(1296269)),
		Name:     new("Programming Elixir"),
		Archived: new(false),
		Organization: &Organization{
			ID:        new(int64(1)),
			Login:     new("programming-elixir"),
			NodeID:    new("MDEyOk9yZ2FuaXphdGlvbjE="),
			HTMLURL:   new("https://example.com/programming-elixir"),
			Name:      new("Learn how to build fault tolerant applications"),
			AvatarURL: new("https://example.com/avatars/u/9919?v=4"),
		},
		URL: new("https://example.com/classrooms/programming"),
	}

	if !cmp.Equal(classroom, want) {
		t.Errorf("Classroom.GetClassroom returned %+v, want %+v", classroom, want)
	}

	const methodName = "GetClassroom"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Classroom.GetClassroom(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.GetClassroom(ctx, 1296269)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestClassroomService_ListClassrooms(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/classrooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "2"})
		fmt.Fprint(w, `[
			{
				"id": 1296269,
				"name": "Programming Elixir",
				"archived": false,
				"url": "https://example.com/classrooms/programming"
			},
			{
				"id": 1296270,
				"name": "Advanced Programming",
				"archived": true,
				"url": "https://example.com/classrooms/2-advanced-programming"
			}
		]`)
	})

	opt := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	classrooms, _, err := client.Classroom.ListClassrooms(ctx, opt)
	if err != nil {
		t.Errorf("Classroom.ListClassrooms returned error: %v", err)
	}

	want := []*Classroom{
		{
			ID:       new(int64(1296269)),
			Name:     new("Programming Elixir"),
			Archived: new(false),
			URL:      new("https://example.com/classrooms/programming"),
		},
		{
			ID:       new(int64(1296270)),
			Name:     new("Advanced Programming"),
			Archived: new(true),
			URL:      new("https://example.com/classrooms/2-advanced-programming"),
		},
	}

	if !cmp.Equal(classrooms, want) {
		t.Errorf("Classroom.ListClassrooms returned %+v, want %+v", classrooms, want)
	}

	const methodName = "ListClassrooms"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.ListClassrooms(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestClassroomService_ListClassroomAssignments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/classrooms/1296269/assignments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "2"})
		fmt.Fprint(w, `[
			{
				"id": 12,
				"public_repo": false,
				"title": "Intro to Binaries",
				"type": "individual",
				"invite_link": "https://example.com/a/Lx7jiUgx",
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
				"deadline": `+referenceTimeStr+`,
				"classroom": {
					"id": 1296269,
					"name": "Programming Elixir",
					"archived": false,
					"url": "https://example.com/classrooms/programming"
				}
			},
			{
				"id": 13,
				"public_repo": true,
				"title": "Advanced Programming",
				"type": "group",
				"invite_link": "https://example.com/a/AdvancedProg",
				"invitations_enabled": true,
				"slug": "advanced-programming",
				"students_are_repo_admins": true,
				"feedback_pull_requests_enabled": false,
				"max_teams": 5,
				"max_members": 3,
				"editor": "vscode",
				"accepted": 50,
				"submitted": 25,
				"passing": 20,
				"language": "python",
				"deadline": `+referenceTimeStr+`,
				"classroom": {
					"id": 1296269,
					"name": "Programming Elixir",
					"archived": false,
					"url": "https://example.com/classrooms/programming"
				}
			}
		]`)
	})

	opt := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	assignments, _, err := client.Classroom.ListClassroomAssignments(ctx, 1296269, opt)
	if err != nil {
		t.Errorf("Classroom.ListClassroomAssignments returned error: %v", err)
	}

	want := []*ClassroomAssignment{
		{
			ID:                          new(int64(12)),
			PublicRepo:                  new(false),
			Title:                       new("Intro to Binaries"),
			Type:                        new("individual"),
			InviteLink:                  new("https://example.com/a/Lx7jiUgx"),
			InvitationsEnabled:          new(true),
			Slug:                        new("intro-to-binaries"),
			StudentsAreRepoAdmins:       new(false),
			FeedbackPullRequestsEnabled: new(true),
			MaxTeams:                    new(0),
			MaxMembers:                  new(0),
			Editor:                      new("codespaces"),
			Accepted:                    new(100),
			Submitted:                   new(40),
			Passing:                     new(10),
			Language:                    new("ruby"),
			Deadline:                    &referenceTimestamp,
			Classroom: &Classroom{
				ID:       new(int64(1296269)),
				Name:     new("Programming Elixir"),
				Archived: new(false),
				URL:      new("https://example.com/classrooms/programming"),
			},
		},
		{
			ID:                          new(int64(13)),
			PublicRepo:                  new(true),
			Title:                       new("Advanced Programming"),
			Type:                        new("group"),
			InviteLink:                  new("https://example.com/a/AdvancedProg"),
			InvitationsEnabled:          new(true),
			Slug:                        new("advanced-programming"),
			StudentsAreRepoAdmins:       new(true),
			FeedbackPullRequestsEnabled: new(false),
			MaxTeams:                    new(5),
			MaxMembers:                  new(3),
			Editor:                      new("vscode"),
			Accepted:                    new(50),
			Submitted:                   new(25),
			Passing:                     new(20),
			Language:                    new("python"),
			Deadline:                    &referenceTimestamp,
			Classroom: &Classroom{
				ID:       new(int64(1296269)),
				Name:     new("Programming Elixir"),
				Archived: new(false),
				URL:      new("https://example.com/classrooms/programming"),
			},
		},
	}

	if !cmp.Equal(assignments, want) {
		t.Errorf("Classroom.ListClassroomAssignments returned %+v, want %+v", assignments, want)
	}

	const methodName = "ListClassroomAssignments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Classroom.ListClassroomAssignments(ctx, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.ListClassroomAssignments(ctx, 1296269, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestClassroomService_ListAcceptedAssignments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/assignments/12/accepted_assignments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "2"})
		fmt.Fprint(w, `[
			{
				"id": 42,
				"submitted": true,
				"passing": true,
				"commit_count": 5,
				"grade": "10/10",
				"students": [
					{
						"id": 1,
						"login": "octocat",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"html_url": "https://github.com/octocat"
					}
				],
				"repository": {
					"id": 1296269,
					"full_name": "octocat/Hello-World",
					"html_url": "https://github.com/octocat/Hello-World",
					"node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
					"private": false,
					"default_branch": "main"
				},
				"assignment": {
					"id": 12,
					"public_repo": false,
					"title": "Intro to Binaries",
					"type": "individual",
					"invite_link": "https://example.com/a/Lx7jiUgx",
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
					"deadline": `+referenceTimeStr+`,
					"classroom": {
						"id": 1296269,
						"name": "Programming Elixir",
						"archived": false,
						"url": "https://example.com/classrooms/programming"
					}
				}
			},
			{
				"id": 43,
				"submitted": false,
				"passing": false,
				"commit_count": 2,
				"grade": "5/10",
				"students": [
					{
						"id": 2,
						"login": "monalisa",
						"avatar_url": "https://github.com/images/error/monalisa_happy.gif",
						"html_url": "https://github.com/monalisa"
					}
				],
				"repository": {
					"id": 1296270,
					"full_name": "monalisa/Hello-World",
					"html_url": "https://github.com/monalisa/Hello-World",
					"node_id": "MDEwOlJlcG9zaXRvcnkxMjk2Mjcw",
					"private": true,
					"default_branch": "main"
				},
				"assignment": {
					"id": 12,
					"public_repo": false,
					"title": "Intro to Binaries",
					"type": "individual",
					"invite_link": "https://example.com/a/Lx7jiUgx",
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
					"deadline": `+referenceTimeStr+`,
					"classroom": {
						"id": 1296269,
						"name": "Programming Elixir",
						"archived": false,
						"url": "https://example.com/classrooms/programming"
					}
				}
			}
		]`)
	})

	ctx := t.Context()
	opt := &ListOptions{Page: 2, PerPage: 2}
	acceptedAssignments, _, err := client.Classroom.ListAcceptedAssignments(ctx, 12, opt)
	if err != nil {
		t.Errorf("Classroom.ListAcceptedAssignments returned error: %v", err)
	}

	want := []*AcceptedAssignment{
		{
			ID:          new(int64(42)),
			Submitted:   new(true),
			Passing:     new(true),
			CommitCount: new(5),
			Grade:       new("10/10"),
			Students: []*ClassroomUser{
				{
					ID:        new(int64(1)),
					Login:     new("octocat"),
					AvatarURL: new("https://github.com/images/error/octocat_happy.gif"),
					HTMLURL:   new("https://github.com/octocat"),
				},
			},
			Repository: &Repository{
				ID:            new(int64(1296269)),
				FullName:      new("octocat/Hello-World"),
				HTMLURL:       new("https://github.com/octocat/Hello-World"),
				NodeID:        new("MDEwOlJlcG9zaXRvcnkxMjk2MjY5"),
				Private:       new(false),
				DefaultBranch: new("main"),
			},
			Assignment: &ClassroomAssignment{
				ID:                          new(int64(12)),
				PublicRepo:                  new(false),
				Title:                       new("Intro to Binaries"),
				Type:                        new("individual"),
				InviteLink:                  new("https://example.com/a/Lx7jiUgx"),
				InvitationsEnabled:          new(true),
				Slug:                        new("intro-to-binaries"),
				StudentsAreRepoAdmins:       new(false),
				FeedbackPullRequestsEnabled: new(true),
				MaxTeams:                    new(0),
				MaxMembers:                  new(0),
				Editor:                      new("codespaces"),
				Accepted:                    new(100),
				Submitted:                   new(40),
				Passing:                     new(10),
				Language:                    new("ruby"),
				Deadline:                    &referenceTimestamp,
				Classroom: &Classroom{
					ID:       new(int64(1296269)),
					Name:     new("Programming Elixir"),
					Archived: new(false),
					URL:      new("https://example.com/classrooms/programming"),
				},
			},
		},
		{
			ID:          new(int64(43)),
			Submitted:   new(false),
			Passing:     new(false),
			CommitCount: new(2),
			Grade:       new("5/10"),
			Students: []*ClassroomUser{
				{
					ID:        new(int64(2)),
					Login:     new("monalisa"),
					AvatarURL: new("https://github.com/images/error/monalisa_happy.gif"),
					HTMLURL:   new("https://github.com/monalisa"),
				},
			},
			Repository: &Repository{
				ID:            new(int64(1296270)),
				FullName:      new("monalisa/Hello-World"),
				HTMLURL:       new("https://github.com/monalisa/Hello-World"),
				NodeID:        new("MDEwOlJlcG9zaXRvcnkxMjk2Mjcw"),
				Private:       new(true),
				DefaultBranch: new("main"),
			},
			Assignment: &ClassroomAssignment{
				ID:                          new(int64(12)),
				PublicRepo:                  new(false),
				Title:                       new("Intro to Binaries"),
				Type:                        new("individual"),
				InviteLink:                  new("https://example.com/a/Lx7jiUgx"),
				InvitationsEnabled:          new(true),
				Slug:                        new("intro-to-binaries"),
				StudentsAreRepoAdmins:       new(false),
				FeedbackPullRequestsEnabled: new(true),
				MaxTeams:                    new(0),
				MaxMembers:                  new(0),
				Editor:                      new("codespaces"),
				Accepted:                    new(100),
				Submitted:                   new(40),
				Passing:                     new(10),
				Language:                    new("ruby"),
				Deadline:                    &referenceTimestamp,
				Classroom: &Classroom{
					ID:       new(int64(1296269)),
					Name:     new("Programming Elixir"),
					Archived: new(false),
					URL:      new("https://example.com/classrooms/programming"),
				},
			},
		},
	}

	if !cmp.Equal(acceptedAssignments, want) {
		t.Errorf("Classroom.ListAcceptedAssignments returned %+v, want %+v", acceptedAssignments, want)
	}

	const methodName = "ListAcceptedAssignments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Classroom.ListAcceptedAssignments(ctx, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.ListAcceptedAssignments(ctx, 12, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestClassroomService_GetAssignmentGrades(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/assignments/12/grades", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"assignment_name": "Intro to Binaries",
				"assignment_url": "https://classroom.github.com/assignments/12",
				"starter_code_url": "https://github.com/octocat/Hello-World",
				"github_username": "octocat",
				"roster_identifier": "student123",
				"student_repository_name": "octocat/intro-to-binaries",
				"student_repository_url": "https://github.com/octocat/intro-to-binaries",
				"submission_timestamp": `+referenceTimeStr+`,
				"points_awarded": 10,
				"points_available": 10,
				"group_name": "Team Alpha"
			},
			{
				"assignment_name": "Intro to Binaries",
				"assignment_url": "https://classroom.github.com/assignments/12",
				"starter_code_url": "https://github.com/octocat/Hello-World",
				"github_username": "monalisa",
				"roster_identifier": "student456",
				"student_repository_name": "monalisa/intro-to-binaries",
				"student_repository_url": "https://github.com/monalisa/intro-to-binaries",
				"submission_timestamp": `+referenceTimeStr+`,
				"points_awarded": 8,
				"points_available": 10,
				"group_name": "Team Beta"
			}
		]`)
	})

	ctx := t.Context()
	grades, _, err := client.Classroom.GetAssignmentGrades(ctx, 12)
	if err != nil {
		t.Errorf("Classroom.GetAssignmentGrades returned error: %v", err)
	}

	want := []*AssignmentGrade{
		{
			AssignmentName:        new("Intro to Binaries"),
			AssignmentURL:         new("https://classroom.github.com/assignments/12"),
			StarterCodeURL:        new("https://github.com/octocat/Hello-World"),
			GithubUsername:        new("octocat"),
			RosterIdentifier:      new("student123"),
			StudentRepositoryName: new("octocat/intro-to-binaries"),
			StudentRepositoryURL:  new("https://github.com/octocat/intro-to-binaries"),
			SubmissionTimestamp:   &referenceTimestamp,
			PointsAwarded:         new(10),
			PointsAvailable:       new(10),
			GroupName:             new("Team Alpha"),
		},
		{
			AssignmentName:        new("Intro to Binaries"),
			AssignmentURL:         new("https://classroom.github.com/assignments/12"),
			StarterCodeURL:        new("https://github.com/octocat/Hello-World"),
			GithubUsername:        new("monalisa"),
			RosterIdentifier:      new("student456"),
			StudentRepositoryName: new("monalisa/intro-to-binaries"),
			StudentRepositoryURL:  new("https://github.com/monalisa/intro-to-binaries"),
			SubmissionTimestamp:   &referenceTimestamp,
			PointsAwarded:         new(8),
			PointsAvailable:       new(10),
			GroupName:             new("Team Beta"),
		},
	}

	if !cmp.Equal(grades, want) {
		t.Errorf("Classroom.GetAssignmentGrades returned %+v, want %+v", grades, want)
	}

	const methodName = "GetAssignmentGrades"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Classroom.GetAssignmentGrades(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Classroom.GetAssignmentGrades(ctx, 12)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
