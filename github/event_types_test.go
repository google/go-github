// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
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

func TestWorkflowDispatchEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &WorkflowDispatchEvent{}, "{}")

	i := make(map[string]interface{})
	i["key"] = "value"

	jsonMsg, _ := json.Marshal(i)
	u := &WorkflowDispatchEvent{
		Inputs:   jsonMsg,
		Ref:      String("r"),
		Workflow: String("w"),
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
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
		"inputs": {
			"key": "value"
		},
		"ref": "r",
		"workflow": "w",
		"repository": {
			"id": 1,
			"name": "n",
			"url": "s"
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

func TestWatchEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &WatchEvent{}, "{}")

	u := &WatchEvent{
		Action: String("a"),
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
		Installation: &Installation{
			ID:       Int64(1),
			NodeID:   String("nid"),
			AppID:    Int64(1),
			AppSlug:  String("as"),
			TargetID: Int64(1),
			Account: &User{
				Login:           String("l"),
				ID:              Int64(1),
				URL:             String("u"),
				AvatarURL:       String("a"),
				GravatarID:      String("g"),
				Name:            String("n"),
				Company:         String("c"),
				Blog:            String("b"),
				Location:        String("l"),
				Email:           String("e"),
				Hireable:        Bool(true),
				Bio:             String("b"),
				TwitterUsername: String("t"),
				PublicRepos:     Int(1),
				Followers:       Int(1),
				Following:       Int(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     String("atu"),
			RepositoriesURL:     String("ru"),
			HTMLURL:             String("hu"),
			TargetType:          String("tt"),
			SingleFileName:      String("sfn"),
			RepositorySelection: String("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       String("a"),
				Administration:                String("ad"),
				Checks:                        String("c"),
				Contents:                      String("co"),
				ContentReferences:             String("cr"),
				Deployments:                   String("d"),
				Environments:                  String("e"),
				Issues:                        String("i"),
				Metadata:                      String("md"),
				Members:                       String("m"),
				OrganizationAdministration:    String("oa"),
				OrganizationHooks:             String("oh"),
				OrganizationPlan:              String("op"),
				OrganizationPreReceiveHooks:   String("opr"),
				OrganizationProjects:          String("op"),
				OrganizationSecrets:           String("os"),
				OrganizationSelfHostedRunners: String("osh"),
				OrganizationUserBlocking:      String("oub"),
				Packages:                      String("pkg"),
				Pages:                         String("pg"),
				PullRequests:                  String("pr"),
				RepositoryHooks:               String("rh"),
				RepositoryProjects:            String("rp"),
				RepositoryPreReceiveHooks:     String("rprh"),
				Secrets:                       String("s"),
				SecretScanningAlerts:          String("ssa"),
				SecurityEvents:                String("se"),
				SingleFile:                    String("sf"),
				Statuses:                      String("s"),
				TeamDiscussions:               String("td"),
				VulnerabilityAlerts:           String("va"),
				Workflows:                     String("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Bool(false),
			SuspendedBy: &User{
				Login:           String("l"),
				ID:              Int64(1),
				URL:             String("u"),
				AvatarURL:       String("a"),
				GravatarID:      String("g"),
				Name:            String("n"),
				Company:         String("c"),
				Blog:            String("b"),
				Location:        String("l"),
				Email:           String("e"),
				Hireable:        Bool(true),
				Bio:             String("b"),
				TwitterUsername: String("t"),
				PublicRepos:     Int(1),
				Followers:       Int(1),
				Following:       Int(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
	}

	want := `{
		"action": "a",
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
		},
		"installation": {
			"id": 1,
			"node_id": "nid",
			"app_id": 1,
			"app_slug": "as",
			"target_id": 1,
			"account": {
				"login": "l",
				"id": 1,
				"avatar_url": "a",
				"gravatar_id": "g",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e",
				"hireable": true,
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			},
			"access_tokens_url": "atu",
			"repositories_url": "ru",
			"html_url": "hu",
			"target_type": "tt",
			"single_file_name": "sfn",
			"repository_selection": "rs",
			"events": [
				"e"
			],
			"single_file_paths": [
				"s"
			],
			"permissions": {
				"actions": "a",
				"administration": "ad",
				"checks": "c",
				"contents": "co",
				"content_references": "cr",
				"deployments": "d",
				"environments": "e",
				"issues": "i",
				"metadata": "md",
				"members": "m",
				"organization_administration": "oa",
				"organization_hooks": "oh",
				"organization_plan": "op",
				"organization_pre_receive_hooks": "opr",
				"organization_projects": "op",
				"organization_secrets": "os",
				"organization_self_hosted_runners": "osh",
				"organization_user_blocking": "oub",
				"packages": "pkg",
				"pages": "pg",
				"pull_requests": "pr",
				"repository_hooks": "rh",
				"repository_projects": "rp",
				"repository_pre_receive_hooks": "rprh",
				"secrets": "s",
				"secret_scanning_alerts": "ssa",
				"security_events": "se",
				"single_file": "sf",
				"statuses": "s",
				"team_discussions": "td",
				"vulnerability_alerts": "va",
				"workflows": "w"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"has_multiple_single_files": false,
			"suspended_by": {
				"login": "l",
				"id": 1,
				"avatar_url": "a",
				"gravatar_id": "g",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e",
				"hireable": true,
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			},
			"suspended_at": ` + referenceTimeStr + `
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestUserEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &UserEvent{}, "{}")

	u := &UserEvent{
		User: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		// The action performed. Possible values are: "created" or "deleted".
		Action: String("a"),
		Enterprise: &Enterprise{
			ID:          Int(1),
			Slug:        String("s"),
			Name:        String("n"),
			NodeID:      String("nid"),
			AvatarURL:   String("au"),
			Description: String("d"),
			WebsiteURL:  String("wu"),
			HTMLURL:     String("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
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
		"user": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"action": "a",
		"enterprise": {
			"id": 1,
			"slug": "s",
			"name": "n",
			"node_id": "nid",
			"avatar_url": "au",
			"description": "d",
			"website_url": "wu",
			"html_url": "hu",
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `
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
