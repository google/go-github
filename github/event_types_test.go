// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"
)

func TestEditChange_Marshal_TitleChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	TitleFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("TitleFrom"),
	}

	u := &EditChange{
		Title: &TitleFrom,
		Body:  nil,
		Base:  nil,
	}

	want := `{
		"title": {
			"from": "TitleFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_BodyChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	BodyFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BodyFrom"),
	}

	u := &EditChange{
		Title: nil,
		Body:  &BodyFrom,
		Base:  nil,
	}

	want := `{
		"body": {
			"from": "BodyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_BaseChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	RefFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BaseRefFrom"),
	}

	SHAFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BaseSHAFrom"),
	}

	Base := struct {
		Ref *struct {
			From *string `json:"from,omitempty"`
		} `json:"ref,omitempty"`
		SHA *struct {
			From *string `json:"from,omitempty"`
		} `json:"sha,omitempty"`
	}{
		Ref: &RefFrom,
		SHA: &SHAFrom,
	}

	u := &EditChange{
		Title: nil,
		Body:  nil,
		Base:  &Base,
	}

	want := `{
		"base": {
			"ref": {
				"from": "BaseRefFrom"
			},
			"sha": {
				"from": "BaseSHAFrom"
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestWorkflowRunEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &WorkflowRunEvent{}, "{}")

	u := &WorkflowRunEvent{
		Action: String("a"),
		Workflow: &Workflow{
			ID:        Int64(1),
			NodeID:    String("nid"),
			Name:      String("n"),
			Path:      String("p"),
			State:     String("s"),
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			URL:       String("u"),
			HTMLURL:   String("h"),
			BadgeURL:  String("b"),
		},
		WorkflowRun: &WorkflowRun{
			ID:         Int64(1),
			Name:       String("n"),
			NodeID:     String("nid"),
			HeadBranch: String("hb"),
			HeadSHA:    String("hs"),
			RunNumber:  Int(1),
			Event:      String("e"),
			Status:     String("s"),
			Conclusion: String("c"),
			WorkflowID: Int64(1),
			URL:        String("u"),
			HTMLURL:    String("h"),
			PullRequests: []*PullRequest{
				{
					URL:    String("u"),
					ID:     Int64(1),
					Number: Int(1),
					Head: &PullRequestBranch{
						Ref: String("r"),
						SHA: String("s"),
						Repo: &Repository{
							ID:   Int64(1),
							URL:  String("s"),
							Name: String("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: String("r"),
						SHA: String("s"),
						Repo: &Repository{
							ID:   Int64(1),
							URL:  String("u"),
							Name: String("n"),
						},
					},
				},
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			JobsURL:       String("j"),
			LogsURL:       String("l"),
			CheckSuiteURL: String("c"),
			ArtifactsURL:  String("a"),
			CancelURL:     String("c"),
			RerunURL:      String("r"),
			HeadCommit: &HeadCommit{
				Message: String("m"),
				Author: &CommitAuthor{
					Name:  String("n"),
					Email: String("e"),
					Login: String("l"),
				},
				URL:       String("u"),
				Distinct:  Bool(false),
				SHA:       String("s"),
				ID:        String("i"),
				TreeID:    String("tid"),
				Timestamp: &Timestamp{referenceTime},
				Committer: &CommitAuthor{
					Name:  String("n"),
					Email: String("e"),
					Login: String("l"),
				},
			},
			WorkflowURL: String("w"),
			Repository: &Repository{
				ID:   Int64(1),
				URL:  String("u"),
				Name: String("n"),
			},
			HeadRepository: &Repository{
				ID:   Int64(1),
				URL:  String("u"),
				Name: String("n"),
			},
		},
		Org: &Organization{
			BillingEmail:                         String("be"),
			Blog:                                 String("b"),
			Company:                              String("c"),
			Email:                                String("e"),
			TwitterUsername:                      String("tu"),
			Location:                             String("loc"),
			Name:                                 String("n"),
			Description:                          String("d"),
			IsVerified:                           Bool(true),
			HasOrganizationProjects:              Bool(true),
			HasRepositoryProjects:                Bool(true),
			DefaultRepoPermission:                String("drp"),
			MembersCanCreateRepos:                Bool(true),
			MembersCanCreateInternalRepos:        Bool(true),
			MembersCanCreatePrivateRepos:         Bool(true),
			MembersCanCreatePublicRepos:          Bool(false),
			MembersAllowedRepositoryCreationType: String("marct"),
			MembersCanCreatePages:                Bool(true),
			MembersCanCreatePublicPages:          Bool(false),
			MembersCanCreatePrivatePages:         Bool(true),
		},
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Sender: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
	}

	want := `{
		"action": "a",
		"workflow": {
			"id": 1,
			"node_id": "nid",
			"name": "n",
			"path": "p",
			"state": "s",
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"url": "u",
			"html_url": "h",
			"badge_url": "b"
		},
		"workflow_run": {
			"id": 1,
			"name": "n",
			"node_id": "nid",
			"head_branch": "hb",
			"head_sha": "hs",
			"run_number": 1,
			"event": "e",
			"status": "s",
			"conclusion": "c",
			"workflow_id": 1,
			"url": "u",
			"html_url": "h",
			"pull_requests": [
				{
					"id": 1,
					"number": 1,
					"url": "u",
					"head": {
						"ref": "r",
						"sha": "s",
						"repo": {
							"id": 1,
							"name": "n",
							"url": "s"
						}
					},
					"base": {
						"ref": "r",
						"sha": "s",
						"repo": {
							"id": 1,
							"name": "n",
							"url": "u"
						}
					}
				}
			],
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"jobs_url": "j",
			"logs_url": "l",
			"check_suite_url": "c",
			"artifacts_url": "a",
			"cancel_url": "c",
			"rerun_url": "r",
			"head_commit": {
				"message": "m",
				"author": {
					"name": "n",
					"email": "e",
					"username": "l"
				},
				"url": "u",
				"distinct": false,
				"sha": "s",
				"id": "i",
				"tree_id": "tid",
				"timestamp": ` + referenceTimeStr + `,
				"committer": {
					"name": "n",
					"email": "e",
					"username": "l"
				}
			},
			"workflow_url": "w",
			"repository": {
				"id": 1,
				"name": "n",
				"url": "u"
			},
			"head_repository": {
				"id": 1,
				"name": "n",
				"url": "u"
			}
		},
		"organization": {
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "loc",
			"email": "e",
			"twitter_username": "tu",
			"description": "d",
			"billing_email": "be",
			"is_verified": true,
			"has_organization_projects": true,
			"has_repository_projects": true,
			"default_repository_permission": "drp",
			"members_can_create_repositories": true,
			"members_can_create_public_repositories": false,
			"members_can_create_private_repositories": true,
			"members_can_create_internal_repositories": true,
			"members_allowed_repository_creation_type": "marct",
			"members_can_create_pages": true,
			"members_can_create_public_pages": false,
			"members_can_create_private_pages": true
		},
		"repository": {
			"id": 1,
			"name": "n",
			"url": "s"
		},
		"sender": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		}
	}`

	testJSONMarshal(t, u, want)
}
