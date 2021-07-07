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

func TestCheckRunEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckRunEvent{}, "{}")

	r := &CheckRunEvent{
		CheckRun: &CheckRun{
			ID:          Int64(1),
			NodeID:      String("n"),
			HeadSHA:     String("h"),
			ExternalID:  String("1"),
			URL:         String("u"),
			HTMLURL:     String("u"),
			DetailsURL:  String("u"),
			Status:      String("s"),
			Conclusion:  String("c"),
			StartedAt:   &Timestamp{referenceTime},
			CompletedAt: &Timestamp{referenceTime},
			Output: &CheckRunOutput{
				Annotations: []*CheckRunAnnotation{
					{
						AnnotationLevel: String("a"),
						EndLine:         Int(1),
						Message:         String("m"),
						Path:            String("p"),
						RawDetails:      String("r"),
						StartLine:       Int(1),
						Title:           String("t"),
					},
				},
				AnnotationsCount: Int(1),
				AnnotationsURL:   String("a"),
				Images: []*CheckRunImage{
					{
						Alt:      String("a"),
						ImageURL: String("i"),
						Caption:  String("c"),
					},
				},
				Title:   String("t"),
				Summary: String("s"),
				Text:    String("t"),
			},
			Name: String("n"),
			CheckSuite: &CheckSuite{
				ID: Int64(1),
			},
			App: &App{
				ID:     Int64(1),
				NodeID: String("n"),
				Owner: &User{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					URL:       String("u"),
					ReposURL:  String("r"),
					EventsURL: String("e"),
					AvatarURL: String("a"),
				},
				Name:        String("n"),
				Description: String("d"),
				HTMLURL:     String("h"),
				ExternalURL: String("u"),
				CreatedAt:   &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
			},
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
		},
		Action: String("a"),
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
		RequestedAction: &RequestedAction{
			Identifier: "i",
		},
	}

	want := `{
		"check_run": {
			"id": 1,
			"node_id": "n",
			"head_sha": "h",
			"external_id": "1",
			"url": "u",
			"html_url": "u",
			"details_url": "u",
			"status": "s",
			"conclusion": "c",
			"started_at": ` + referenceTimeStr + `,
			"completed_at": ` + referenceTimeStr + `,
			"output": {
				"title": "t",
				"summary": "s",
				"text": "t",
				"annotations_count": 1,
				"annotations_url": "a",
				"annotations": [
					{
						"path": "p",
						"start_line": 1,
						"end_line": 1,
						"annotation_level": "a",
						"message": "m",
						"title": "t",
						"raw_details": "r"
					}
				],
				"images": [
					{
						"alt": "a",
						"image_url": "i",
						"caption": "c"
					}
				]
			},
			"name": "n",
			"check_suite": {
				"id": 1
			},
			"app": {
				"id": 1,
				"node_id": "n",
				"owner": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"url": "u",
					"events_url": "e",
					"repos_url": "r"
				},
				"name": "n",
				"description": "d",
				"external_url": "u",
				"html_url": "h",
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `
			},
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
			]
		},
		"action": "a",
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
		},
		"requested_action": {
			"identifier": "i"
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestCheckSuiteEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckSuiteEvent{}, "{}")

	r := &CheckSuiteEvent{
		CheckSuite: &CheckSuite{
			ID:         Int64(1),
			NodeID:     String("n"),
			HeadBranch: String("h"),
			HeadSHA:    String("h"),
			URL:        String("u"),
			BeforeSHA:  String("b"),
			AfterSHA:   String("a"),
			Status:     String("s"),
			Conclusion: String("c"),
			App: &App{
				ID:     Int64(1),
				NodeID: String("n"),
				Owner: &User{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					URL:       String("u"),
					ReposURL:  String("r"),
					EventsURL: String("e"),
					AvatarURL: String("a"),
				},
				Name:        String("n"),
				Description: String("d"),
				HTMLURL:     String("h"),
				ExternalURL: String("u"),
				CreatedAt:   &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
			},
			Repository: &Repository{
				ID: Int64(1),
			},
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
			HeadCommit: &Commit{
				SHA: String("s"),
			},
		},
		Action: String("a"),
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
		"check_suite": {
			"id": 1,
			"node_id": "n",
			"head_branch": "h",
			"head_sha": "h",
			"url": "u",
			"before": "b",
			"after": "a",
			"status": "s",
			"conclusion": "c",
			"app": {
				"id": 1,
				"node_id": "n",
				"owner": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"url": "u",
					"events_url": "e",
					"repos_url": "r"
				},
				"name": "n",
				"description": "d",
				"external_url": "u",
				"html_url": "h",
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `
			},
			"repository": {
				"id": 1
			},
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
		"head_commit": {
			"sha": "s"
		}
		},
		"action": "a",
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

	testJSONMarshal(t, r, want)
}
