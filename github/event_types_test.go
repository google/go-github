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

	u := &EditChange{
		Title: &EditTitle{
			From: String("TitleFrom"),
		},
		Body: nil,
		Base: nil,
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

	u := &EditChange{
		Title: nil,
		Body: &EditBody{
			From: String("BodyFrom"),
		},
		Base: nil,
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

	Base := EditBase{
		Ref: &EditRef{
			From: String("BaseRefFrom"),
		},
		SHA: &EditSHA{
			From: String("BaseSHAFrom"),
		},
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

func TestEditChange_Marshal_Repo(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Repo: &EditRepo{
			Name: &RepoName{
				From: String("old-repo-name"),
			},
		},
		Topics: &EditTopics{
			From: []string{"topic1", "topic2"},
		},
	}

	want := `{
		"repository": {
			"name": {
				"from": "old-repo-name"
			}
		},
		"topics": {
			"from": [
				"topic1",
				"topic2"
			]
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_TransferFromUser(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Owner: &EditOwner{
			OwnerInfo: &OwnerInfo{
				User: &User{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					URL:       String("u"),
					ReposURL:  String("r"),
					EventsURL: String("e"),
					AvatarURL: String("a"),
				},
			},
		},
	}

	want := `{
		"owner": {
			"from": {
				"user": {
					"login": "l",
          			"id": 1,
         		 	"node_id": "n",
          			"avatar_url": "a",
          			"url": "u",
          			"repos_url": "r",
          			"events_url": "e"
				}
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_TransferFromOrg(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Owner: &EditOwner{
			OwnerInfo: &OwnerInfo{
				Org: &User{
					Login:     String("l"),
					ID:        Int64(1),
					NodeID:    String("n"),
					URL:       String("u"),
					ReposURL:  String("r"),
					EventsURL: String("e"),
					AvatarURL: String("a"),
				},
			},
		},
	}

	want := `{
		"owner": {
			"from": {
				"organization": {
					"login": "l",
          			"id": 1,
         		 	"node_id": "n",
          			"avatar_url": "a",
          			"url": "u",
          			"repos_url": "r",
          			"events_url": "e"
				}
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: &ProjectName{From: String("NameFrom")},
		Body: nil,
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectChange_Marshal_BodyChange(t *testing.T) {
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: nil,
		Body: &ProjectBody{From: String("BodyFrom")},
	}

	want := `{
		"body": {
			"from": "BodyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardChange_Marshal_NoteChange(t *testing.T) {
	testJSONMarshal(t, &ProjectCardChange{}, "{}")

	u := &ProjectCardChange{
		Note: &ProjectCardNote{From: String("NoteFrom")},
	}

	want := `{
		"note": {
			"from": "NoteFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnChange{}, "{}")

	u := &ProjectColumnChange{
		Name: &ProjectColumnName{From: String("NameFrom")},
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamAddEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamAddEvent{}, "{}")

	u := &TeamAddEvent{
		Team: &Team{
			ID:              Int64(1),
			NodeID:          String("n"),
			Name:            String("n"),
			Description:     String("d"),
			URL:             String("u"),
			Slug:            String("s"),
			Permission:      String("p"),
			Privacy:         String("p"),
			MembersCount:    Int(1),
			ReposCount:      Int(1),
			MembersURL:      String("m"),
			RepositoriesURL: String("r"),
			Organization: &Organization{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				AvatarURL: String("a"),
				HTMLURL:   String("h"),
				Name:      String("n"),
				Company:   String("c"),
				Blog:      String("b"),
				Location:  String("l"),
				Email:     String("e"),
			},
			Parent: &Team{
				ID:           Int64(1),
				NodeID:       String("n"),
				Name:         String("n"),
				Description:  String("d"),
				URL:          String("u"),
				Slug:         String("s"),
				Permission:   String("p"),
				Privacy:      String("p"),
				MembersCount: Int(1),
				ReposCount:   Int(1),
			},
			LDAPDN: String("l"),
		},
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
		"team": {
			"id": 1,
			"node_id": "n",
			"name": "n",
			"description": "d",
			"url": "u",
			"slug": "s",
			"permission": "p",
			"privacy": "p",
			"members_count": 1,
			"repos_count": 1,
			"organization": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"html_url": "h",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e"
			},
			"members_url": "m",
			"repositories_url": "r",
			"parent": {
				"id": 1,
				"node_id": "n",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "p",
				"members_count": 1,
				"repos_count": 1
			},
			"ldap_dn": "l"
		},
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

	testJSONMarshal(t, u, want)
}

func TestStarEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &StarEvent{}, "{}")

	u := &StarEvent{
		Action:    String("a"),
		StarredAt: &Timestamp{referenceTime},
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
		"starred_at": ` + referenceTimeStr + `,
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

func TestTeamEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamEvent{}, "{}")

	u := &TeamEvent{
		Action: String("a"),
		Team: &Team{
			ID:              Int64(1),
			NodeID:          String("n"),
			Name:            String("n"),
			Description:     String("d"),
			URL:             String("u"),
			Slug:            String("s"),
			Permission:      String("p"),
			Privacy:         String("p"),
			MembersCount:    Int(1),
			ReposCount:      Int(1),
			MembersURL:      String("m"),
			RepositoriesURL: String("r"),
			Organization: &Organization{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				AvatarURL: String("a"),
				HTMLURL:   String("h"),
				Name:      String("n"),
				Company:   String("c"),
				Blog:      String("b"),
				Location:  String("l"),
				Email:     String("e"),
			},
			Parent: &Team{
				ID:           Int64(1),
				NodeID:       String("n"),
				Name:         String("n"),
				Description:  String("d"),
				URL:          String("u"),
				Slug:         String("s"),
				Permission:   String("p"),
				Privacy:      String("p"),
				MembersCount: Int(1),
				ReposCount:   Int(1),
			},
			LDAPDN: String("l"),
		},
		Changes: &TeamChange{
			Description: &TeamDescription{
				From: String("from"),
			},
			Name: &TeamName{
				From: String("from"),
			},
			Privacy: &TeamPrivacy{
				From: String("from"),
			},
			Repository: &TeamRepository{
				Permissions: &TeamPermissions{
					From: &TeamPermissionsFrom{
						Admin: Bool(true),
						Pull:  Bool(true),
						Push:  Bool(true),
					},
				},
			},
		},
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
		"action": "a",
		"team": {
			"id": 1,
			"node_id": "n",
			"name": "n",
			"description": "d",
			"url": "u",
			"slug": "s",
			"permission": "p",
			"privacy": "p",
			"members_count": 1,
			"repos_count": 1,
			"organization": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"html_url": "h",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e"
			},
			"members_url": "m",
			"repositories_url": "r",
			"parent": {
				"id": 1,
				"node_id": "n",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "p",
				"members_count": 1,
				"repos_count": 1
			},
			"ldap_dn": "l"
		},
		"changes": {
			"description": {
				"from": "from"
			},
			"name": {
				"from": "from"
			},
			"privacy": {
				"from": "from"
			},
			"repository": {
				"permissions": {
					"from": {
						"admin": true,
						"pull": true,
						"push": true
					}
				}
			}
		},
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

	testJSONMarshal(t, u, want)
}

func TestInstallationRepositoriesEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &InstallationRepositoriesEvent{}, "{}")

	u := &InstallationRepositoriesEvent{
		Action: String("a"),
		RepositoriesAdded: []*Repository{
			{
				ID:   Int64(1),
				URL:  String("s"),
				Name: String("n"),
			},
		},
		RepositoriesRemoved: []*Repository{
			{
				ID:   Int64(1),
				URL:  String("s"),
				Name: String("n"),
			},
		},
		RepositorySelection: String("rs"),
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
		"repositories_added": [
			{
				"id": 1,
				"name": "n",
				"url": "s"
			}
		],
		"repositories_removed": [
			{
				"id": 1,
				"name": "n",
				"url": "s"
			}
		],
		"repository_selection": "rs",
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

func TestInstallationTargetEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &InstallationTargetEvent{}, "{}")

	u := &InstallationTargetEvent{
		Account: &User{
			Login:     String("u"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("l"),
		},
		Action: String("a"),
		Changes: &InstallationChanges{
			Login: &InstallationLoginChange{
				From: String("p"),
			},
			Slug: &InstallationSlugChange{
				From: String("j"),
			},
		},
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
		Organization: &Organization{
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
		Repository: &Repository{
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
		TargetType: String("running"),
	}

	want := `{
		"account": {
			"login": "u",
			"id": 1,
			"node_id": "n",
			"avatar_url": "l",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"action": "a",
		"changes": {
			"login": {
				"from": "p"
			},
			"slug": {
				"from": "j"
			}
		},
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
			"url": "s",
			"name": "n"
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
		"target_type": "running"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditTitle_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditTitle{}, "{}")

	u := &EditTitle{
		From: String("EditTitleFrom"),
	}

	want := `{
		"from": "EditTitleFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBody_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditBody{}, "{}")

	u := &EditBody{
		From: String("EditBodyFrom"),
	}

	want := `{
		"from": "EditBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBase_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditBase{}, "{}")

	u := &EditBase{
		Ref: &EditRef{
			From: String("EditRefFrom"),
		},
		SHA: &EditSHA{
			From: String("EditSHAFrom"),
		},
	}

	want := `{
		"ref": {
			"from": "EditRefFrom"
		},
		"sha": {
			"from": "EditSHAFrom"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestEditRef_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditRef{}, "{}")

	u := &EditRef{
		From: String("EditRefFrom"),
	}

	want := `{
		"from": "EditRefFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditSHA_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditSHA{}, "{}")

	u := &EditSHA{
		From: String("EditSHAFrom"),
	}

	want := `{
		"from": "EditSHAFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectName_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectName{}, "{}")

	u := &ProjectName{
		From: String("ProjectNameFrom"),
	}

	want := `{
		"from": "ProjectNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectBody_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectBody{}, "{}")

	u := &ProjectBody{
		From: String("ProjectBodyFrom"),
	}

	want := `{
		"from": "ProjectBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardNote_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCardNote{}, "{}")

	u := &ProjectCardNote{
		From: String("ProjectCardNoteFrom"),
	}

	want := `{
		"from": "ProjectCardNoteFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnName_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnName{}, "{}")

	u := &ProjectColumnName{
		From: String("ProjectColumnNameFrom"),
	}

	want := `{
		"from": "ProjectColumnNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamDescription_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamDescription{}, "{}")

	u := &TeamDescription{
		From: String("TeamDescriptionFrom"),
	}

	want := `{
		"from": "TeamDescriptionFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamName_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamName{}, "{}")

	u := &TeamName{
		From: String("TeamNameFrom"),
	}

	want := `{
		"from": "TeamNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPrivacy_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPrivacy{}, "{}")

	u := &TeamPrivacy{
		From: String("TeamPrivacyFrom"),
	}

	want := `{
		"from": "TeamPrivacyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamRepository_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamRepository{}, "{}")

	u := &TeamRepository{
		Permissions: &TeamPermissions{
			From: &TeamPermissionsFrom{
				Admin: Bool(true),
				Pull:  Bool(true),
				Push:  Bool(true),
			},
		},
	}

	want := `{
		"permissions": {
			"from": {
				"admin": true,
				"pull": true,
				"push": true
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPermissions_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPermissions{}, "{}")

	u := &TeamPermissions{
		From: &TeamPermissionsFrom{
			Admin: Bool(true),
			Pull:  Bool(true),
			Push:  Bool(true),
		},
	}

	want := `{
		"from": {
			"admin": true,
			"pull": true,
			"push": true
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPermissionsFrom_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPermissionsFrom{}, "{}")

	u := &TeamPermissionsFrom{
		Admin: Bool(true),
		Pull:  Bool(true),
		Push:  Bool(true),
	}

	want := `{
		"admin": true,
		"pull": true,
		"push": true
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoryVulnerabilityAlert_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryVulnerabilityAlert{}, "{}")

	u := &RepositoryVulnerabilityAlert{
		ID:                  Int64(1),
		AffectedRange:       String("ar"),
		AffectedPackageName: String("apn"),
		ExternalReference:   String("er"),
		ExternalIdentifier:  String("ei"),
		FixedIn:             String("fi"),
		Dismisser: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		DismissReason: String("dr"),
		DismissedAt:   &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"affected_range": "ar",
		"affected_package_name": "apn",
		"external_reference": "er",
		"external_identifier": "ei",
		"fixed_in": "fi",
		"dismisser": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"dismiss_reason": "dr",
		"dismissed_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestPage_Marshal(t *testing.T) {
	testJSONMarshal(t, &Page{}, "{}")

	u := &Page{
		PageName: String("p"),
		Title:    String("t"),
		Summary:  String("s"),
		Action:   String("a"),
		SHA:      String("s"),
		HTMLURL:  String("h"),
	}

	want := `{
		"page_name": "p",
		"title": "t",
		"summary": "s",
		"action": "a",
		"sha": "s",
		"html_url": "h"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamChange_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	u := &TeamChange{
		Description: &TeamDescription{
			From: String("DescriptionFrom"),
		},
		Name: &TeamName{
			From: String("NameFrom"),
		},
		Privacy: &TeamPrivacy{
			From: String("PrivacyFrom"),
		},
		Repository: &TeamRepository{
			Permissions: &TeamPermissions{
				From: &TeamPermissionsFrom{
					Admin: Bool(false),
					Pull:  Bool(false),
					Push:  Bool(false),
				},
			},
		},
	}

	want := `{
		"description": {
			"from": "DescriptionFrom"
		},
		"name": {
			"from": "NameFrom"
		},
		"privacy": {
			"from": "PrivacyFrom"
		},
		"repository": {
			"permissions": {
				"from": {
					"admin": false,
					"pull": false,
					"push": false
				}
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueCommentEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssueCommentEvent{}, "{}")

	u := &IssueCommentEvent{
		Action:  String("a"),
		Issue:   &Issue{ID: Int64(1)},
		Comment: &IssueComment{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
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
		Organization: &Organization{
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
	}

	want := `{
		"action": "a",
		"issue": {
			"id": 1
		},
		"comment": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
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
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestIssuesEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssuesEvent{}, "{}")

	u := &IssuesEvent{
		Action: String("a"),
		Issue:  &Issue{ID: Int64(1)},
		Assignee: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Label: &Label{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
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
		"issue": {
			"id": 1
		},
		"assignee": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"label": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
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

func TestLabelEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &LabelEvent{}, "{}")

	u := &LabelEvent{
		Action: String("a"),
		Label:  &Label{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
		},
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
		"label": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
		},
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

func TestMilestoneEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MilestoneEvent{}, "{}")

	u := &MilestoneEvent{
		Action:    String("a"),
		Milestone: &Milestone{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
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
		"milestone": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
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

func TestPublicEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PublicEvent{}, "{}")

	u := &PublicEvent{
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

func TestPullRequestReviewEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReviewEvent{}, "{}")

	u := &PullRequestReviewEvent{
		Action:      String("a"),
		Review:      &PullRequestReview{ID: Int64(1)},
		PullRequest: &PullRequest{ID: Int64(1)},
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
		Organization: &Organization{
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
	}

	want := `{
		"action": "a",
		"review": {
			"id": 1
		},
		"pull_request": {
			"id": 1
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
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPushEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PushEvent{}, "{}")

	u := &PushEvent{
		PushID: Int64(1),
		Head:   String("h"),
		Ref:    String("ref"),
		Size:   Int(1),
		Commits: []*HeadCommit{
			{ID: String("id")},
		},
		Before:       String("b"),
		DistinctSize: Int(1),
		After:        String("a"),
		Created:      Bool(true),
		Deleted:      Bool(true),
		Forced:       Bool(true),
		BaseRef:      String("a"),
		Compare:      String("a"),
		Repo:         &PushEventRepository{ID: Int64(1)},
		HeadCommit:   &HeadCommit{ID: String("id")},
		Pusher: &CommitAuthor{
			Login: String("l"),
			Date:  &Timestamp{referenceTime},
			Name:  String("n"),
			Email: String("e"),
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
		Organization: &Organization{
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
	}

	want := `{
		"push_id": 1,
		"head": "h",
		"ref": "ref",
		"size": 1,
		"commits": [
			{
				"id": "id"
			}
		],
		"before": "b",
		"distinct_size": 1,
		"after": "a",
		"created": true,
		"deleted": true,
		"forced": true,
		"base_ref": "a",
		"compare": "a",
		"repository": {
			"id": 1
		},
		"head_commit": {
			"id": "id"
		},
		"pusher": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "l"
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
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestStatusEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &StatusEvent{}, "{}")

	u := &StatusEvent{
		SHA:         String("sha"),
		State:       String("s"),
		Description: String("d"),
		TargetURL:   String("turl"),
		Branches: []*Branch{
			{
				Name:      String("n"),
				Commit:    &RepositoryCommit{NodeID: String("nid")},
				Protected: Bool(false),
			},
		},
		ID:        Int64(1),
		Name:      String("n"),
		Context:   String("c"),
		Commit:    &RepositoryCommit{NodeID: String("nid")},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
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
		"sha": "sha",
		"state": "s",
		"description": "d",
		"target_url": "turl",
		"branches": [
			{
				"name": "n",
				"commit": {
					"node_id": "nid"
				},
				"protected": false
			}
		],
		"id": 1,
		"name": "n",
		"context": "c",
		"commit": {
			"node_id": "nid"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
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

func TestMarketplacePurchaseEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MarketplacePurchaseEvent{}, "{}")

	u := &MarketplacePurchaseEvent{
		Action:        String("a"),
		EffectiveDate: &Timestamp{referenceTime},
		MarketplacePurchase: &MarketplacePurchase{
			BillingCycle:    String("bc"),
			NextBillingDate: &Timestamp{referenceTime},
			UnitCount:       Int(1),
			Plan: &MarketplacePlan{
				URL:                 String("u"),
				AccountsURL:         String("au"),
				ID:                  Int64(1),
				Number:              Int(1),
				Name:                String("n"),
				Description:         String("d"),
				MonthlyPriceInCents: Int(1),
				YearlyPriceInCents:  Int(1),
				PriceModel:          String("pm"),
				UnitName:            String("un"),
				Bullets:             &[]string{"b"},
				State:               String("s"),
				HasFreeTrial:        Bool(false),
			},
			OnFreeTrial:     Bool(false),
			FreeTrialEndsOn: &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
		},
		PreviousMarketplacePurchase: &MarketplacePurchase{
			BillingCycle:    String("bc"),
			NextBillingDate: &Timestamp{referenceTime},
			UnitCount:       Int(1),
			Plan: &MarketplacePlan{
				URL:                 String("u"),
				AccountsURL:         String("au"),
				ID:                  Int64(1),
				Number:              Int(1),
				Name:                String("n"),
				Description:         String("d"),
				MonthlyPriceInCents: Int(1),
				YearlyPriceInCents:  Int(1),
				PriceModel:          String("pm"),
				UnitName:            String("un"),
				Bullets:             &[]string{"b"},
				State:               String("s"),
				HasFreeTrial:        Bool(false),
			},
			OnFreeTrial:     Bool(false),
			FreeTrialEndsOn: &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
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
		"effective_date": ` + referenceTimeStr + `,
		"marketplace_purchase": {
			"billing_cycle": "bc",
			"next_billing_date": ` + referenceTimeStr + `,
			"unit_count": 1,
			"plan": {
				"url": "u",
				"accounts_url": "au",
				"id": 1,
				"number": 1,
				"name": "n",
				"description": "d",
				"monthly_price_in_cents": 1,
				"yearly_price_in_cents": 1,
				"price_model": "pm",
				"unit_name": "un",
				"bullets": [
					"b"
				],
				"state": "s",
				"has_free_trial": false
			},
			"on_free_trial": false,
			"free_trial_ends_on": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `
		},
		"previous_marketplace_purchase": {
			"billing_cycle": "bc",
			"next_billing_date": ` + referenceTimeStr + `,
			"unit_count": 1,
			"plan": {
				"url": "u",
				"accounts_url": "au",
				"id": 1,
				"number": 1,
				"name": "n",
				"description": "d",
				"monthly_price_in_cents": 1,
				"yearly_price_in_cents": 1,
				"price_model": "pm",
				"unit_name": "un",
				"bullets": [
					"b"
				],
				"state": "s",
				"has_free_trial": false
			},
			"on_free_trial": false,
			"free_trial_ends_on": ` + referenceTimeStr + `,
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

func TestOrganizationEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &OrganizationEvent{}, "{}")

	u := &OrganizationEvent{
		Action:     String("a"),
		Invitation: &Invitation{ID: Int64(1)},
		Membership: &Membership{
			URL:             String("url"),
			State:           String("s"),
			Role:            String("r"),
			OrganizationURL: String("ou"),
			Organization: &Organization{
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
			User: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
		},
		Organization: &Organization{
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
		"action": "a",
		"invitation": {
			"id": 1
		},
		"membership": {
			"url": "url",
			"state": "s",
			"role": "r",
			"organization_url": "ou",
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
			"user": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
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

func TestPageBuildEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PageBuildEvent{}, "{}")

	u := &PageBuildEvent{
		Build: &PagesBuild{URL: String("url")},
		ID:    Int64(1),
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
		"build": {
			"url": "url"
		},
		"id": 1,
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

func TestCommitCommentEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &CommitCommentEvent{}, "{}")

	u := &CommitCommentEvent{
		Comment: &RepositoryComment{
			HTMLURL:  String("hurl"),
			URL:      String("url"),
			ID:       Int64(1),
			NodeID:   String("nid"),
			CommitID: String("cid"),
			User: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Reactions: &Reactions{
				TotalCount: Int(1),
				PlusOne:    Int(1),
				MinusOne:   Int(1),
				Laugh:      Int(1),
				Confused:   Int(1),
				Heart:      Int(1),
				Hooray:     Int(1),
				Rocket:     Int(1),
				Eyes:       Int(1),
				URL:        String("url"),
			},
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			Body:      String("b"),
			Path:      String("path"),
			Position:  Int(1),
		},
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
		"comment": {
			"html_url": "hurl",
			"url": "url",
			"id": 1,
			"node_id": "nid",
			"commit_id": "cid",
			"user": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"reactions": {
				"total_count": 1,
				"+1": 1,
				"-1": 1,
				"laugh": 1,
				"confused": 1,
				"heart": 1,
				"hooray": 1,
				"rocket": 1,
				"eyes": 1,
				"url": "url"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"body": "b",
			"path": "path",
			"position": 1
		},
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

func TestDeploymentEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeploymentEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentEvent{
		Deployment: &Deployment{
			URL:         String("url"),
			ID:          Int64(1),
			SHA:         String("sha"),
			Ref:         String("ref"),
			Task:        String("t"),
			Payload:     jsonMsg,
			Environment: String("e"),
			Description: String("d"),
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   String("surl"),
			RepositoryURL: String("rurl"),
			NodeID:        String("nid"),
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
			RunAttempt: Int(1),
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
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            String("j"),
			LogsURL:            String("l"),
			CheckSuiteURL:      String("c"),
			ArtifactsURL:       String("a"),
			CancelURL:          String("c"),
			RerunURL:           String("r"),
			PreviousAttemptURL: String("p"),
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
	}

	want := `{
		"deployment": {
			"url": "url",
			"id": 1,
			"sha": "sha",
			"ref": "ref",
			"task": "t",
			"payload": {
				"key": "value"
			},
			"environment": "e",
			"description": "d",
			"creator": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"statuses_url": "surl",
			"repository_url": "rurl",
			"node_id": "nid"
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
			"run_attempt": 1,
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
			"run_started_at": ` + referenceTimeStr + `,
			"jobs_url": "j",
			"logs_url": "l",
			"check_suite_url": "c",
			"artifacts_url": "a",
			"cancel_url": "c",
			"rerun_url": "r",
			"previous_attempt_url": "p",
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
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestDeploymentProtectionRuleEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeploymentProtectionRuleEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentProtectionRuleEvent{
		Action:                String("a"),
		Environment:           String("e"),
		DeploymentCallbackURL: String("b"),
		Deployment: &Deployment{
			URL:         String("url"),
			ID:          Int64(1),
			SHA:         String("sha"),
			Ref:         String("ref"),
			Task:        String("t"),
			Payload:     jsonMsg,
			Environment: String("e"),
			Description: String("d"),
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   String("surl"),
			RepositoryURL: String("rurl"),
			NodeID:        String("nid"),
		},
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Organization: &Organization{
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
		"environment": "e",
		"deployment_callback_url": "b",
		"deployment": {
			"url": "url",
			"id": 1,
			"sha": "sha",
			"ref": "ref",
			"task": "t",
			"payload": {
				"key": "value"
			},
			"environment": "e",
			"description": "d",
			"creator": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"statuses_url": "surl",
			"repository_url": "rurl",
			"node_id": "nid"
		},
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

func TestDeploymentReviewEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeploymentReviewEvent{}, "{}")

	u := &DeploymentReviewEvent{
		Action:      String("a"),
		Environment: String("e"),
		Requester: &User{
			AvatarURL:         String("a"),
			Email:             String("e"),
			EventsURL:         String("e"),
			FollowersURL:      String("f"),
			FollowingURL:      String("f"),
			GistsURL:          String("g"),
			GravatarID:        String("g"),
			HTMLURL:           String("h"),
			ID:                Int64(1),
			Login:             String("l"),
			Name:              String("n"),
			NodeID:            String("n"),
			OrganizationsURL:  String("o"),
			ReceivedEventsURL: String("r"),
			ReposURL:          String("r"),
			SiteAdmin:         Bool(false),
			StarredURL:        String("s"),
			SubscriptionsURL:  String("s"),
			Type:              String("User"),
			URL:               String("u"),
		},
		Reviewers: []*RequiredReviewer{
			{
				Type: String("User"),
				Reviewer: &User{
					AvatarURL:         String("a"),
					Email:             String("e"),
					EventsURL:         String("e"),
					FollowersURL:      String("f"),
					FollowingURL:      String("f"),
					GistsURL:          String("g"),
					GravatarID:        String("g"),
					HTMLURL:           String("h"),
					ID:                Int64(1),
					Login:             String("l"),
					Name:              String("n"),
					NodeID:            String("n"),
					OrganizationsURL:  String("o"),
					ReceivedEventsURL: String("r"),
					ReposURL:          String("r"),
					SiteAdmin:         Bool(false),
					StarredURL:        String("s"),
					SubscriptionsURL:  String("s"),
					Type:              String("User"),
					URL:               String("u"),
				},
			},
			{
				Type: String("Team"),
				Reviewer: &Team{
					ID:   Int64(1),
					Name: String("n"),
					Slug: String("s"),
				},
			},
		},
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
		Organization: &Organization{
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
		Since: String("s"),
		WorkflowJobRun: &WorkflowJobRun{
			ID:          Int64(1),
			Conclusion:  String("c"),
			Environment: String("e"),
			HTMLURL:     String("h"),
			Name:        String("n"),
			Status:      String("s"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		WorkflowRun: &WorkflowRun{
			ID:         Int64(1),
			Name:       String("n"),
			NodeID:     String("nid"),
			HeadBranch: String("hb"),
			HeadSHA:    String("hs"),
			RunNumber:  Int(1),
			RunAttempt: Int(1),
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
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            String("j"),
			LogsURL:            String("l"),
			CheckSuiteURL:      String("c"),
			ArtifactsURL:       String("a"),
			CancelURL:          String("c"),
			RerunURL:           String("r"),
			PreviousAttemptURL: String("p"),
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
	}

	want := `{
		"action": "a",
		"requester": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"html_url": "h",
			"followers_url": "f",
			"following_url": "f",
			"gists_url": "g",
			"starred_url": "s",
			"subscriptions_url": "s",
			"organizations_url": "o",
			"repos_url": "r",
			"events_url": "e",
			"received_events_url": "r",
			"type": "User",
			"site_admin": false,
			"name": "n",
			"email": "e",
			"gravatar_id": "g"
		},
		"reviewers": [
			{
				"type": "User",
				"reviewer": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"url": "u",
					"html_url": "h",
					"followers_url": "f",
					"following_url": "f",
					"gists_url": "g",
					"starred_url": "s",
					"subscriptions_url": "s",
					"organizations_url": "o",
					"repos_url": "r",
					"events_url": "e",
					"received_events_url": "r",
					"type": "User",
					"site_admin": false,
					"name": "n",
					"email": "e",
					"gravatar_id": "g"
				}
			},
			{
				"type": "Team",
				"reviewer": {
					"id": 1,
					"name": "n",
					"slug": "s"
				}
			}
		],
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
		"environment": "e",
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
		"since": "s",
		"workflow_job_run": {
			"conclusion": "c",
			"created_at": "2006-01-02T15:04:05Z",
			"environment": "e",
			"html_url": "h",
			"id": 1,
			"name": "n",
			"status": "s",
			"updated_at": "2006-01-02T15:04:05Z"
		},
		"workflow_run": {
			"id": 1,
			"name": "n",
			"node_id": "nid",
			"head_branch": "hb",
			"head_sha": "hs",
			"run_number": 1,
			"run_attempt": 1,
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
			"run_started_at": ` + referenceTimeStr + `,
			"jobs_url": "j",
			"logs_url": "l",
			"check_suite_url": "c",
			"artifacts_url": "a",
			"cancel_url": "c",
			"rerun_url": "r",
			"previous_attempt_url": "p",
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
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestDeploymentStatusEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeploymentStatusEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentStatusEvent{
		Deployment: &Deployment{
			URL:         String("url"),
			ID:          Int64(1),
			SHA:         String("sha"),
			Ref:         String("ref"),
			Task:        String("t"),
			Payload:     jsonMsg,
			Environment: String("e"),
			Description: String("d"),
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   String("surl"),
			RepositoryURL: String("rurl"),
			NodeID:        String("nid"),
		},
		DeploymentStatus: &DeploymentStatus{
			ID:    Int64(1),
			State: String("s"),
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Description:    String("s"),
			Environment:    String("s"),
			NodeID:         String("s"),
			CreatedAt:      &Timestamp{referenceTime},
			UpdatedAt:      &Timestamp{referenceTime},
			TargetURL:      String("s"),
			DeploymentURL:  String("s"),
			RepositoryURL:  String("s"),
			EnvironmentURL: String("s"),
			LogURL:         String("s"),
			URL:            String("s"),
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
		"deployment": {
			"url": "url",
			"id": 1,
			"sha": "sha",
			"ref": "ref",
			"task": "t",
			"payload": {
				"key": "value"
			},
			"environment": "e",
			"description": "d",
			"creator": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"statuses_url": "surl",
			"repository_url": "rurl",
			"node_id": "nid"
		},
		"deployment_status": {
			"id": 1,
			"state": "s",
			"creator": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"description": "s",
			"environment": "s",
			"node_id": "s",
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"target_url": "s",
			"deployment_url": "s",
			"repository_url": "s",
			"environment_url": "s",
			"log_url": "s",
			"url": "s"
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

func TestDiscussionCommentEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DiscussionCommentEvent{}, "{}")

	u := &DiscussionCommentEvent{
		Comment: &CommentDiscussion{
			AuthorAssociation: String("aa"),
			Body:              String("bo"),
			ChildCommentCount: Int(1),
			CreatedAt:         &Timestamp{referenceTime},
			DiscussionID:      Int64(1),
			HTMLURL:           String("hurl"),
			ID:                Int64(1),
			NodeID:            String("nid"),
			ParentID:          Int64(1),
			Reactions: &Reactions{
				TotalCount: Int(1),
				PlusOne:    Int(1),
				MinusOne:   Int(1),
				Laugh:      Int(1),
				Confused:   Int(1),
				Heart:      Int(1),
				Hooray:     Int(1),
				Rocket:     Int(1),
				Eyes:       Int(1),
				URL:        String("url"),
			},
			RepositoryURL: String("rurl"),
			UpdatedAt:     &Timestamp{referenceTime},
			User: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
		},
		Discussion: &Discussion{
			RepositoryURL: String("rurl"),
			DiscussionCategory: &DiscussionCategory{
				ID:           Int64(1),
				NodeID:       String("nid"),
				RepositoryID: Int64(1),
				Emoji:        String("emoji"),
				Name:         String("name"),
				Description:  String("description"),
				CreatedAt:    &Timestamp{referenceTime},
				UpdatedAt:    &Timestamp{referenceTime},
				Slug:         String("slug"),
				IsAnswerable: Bool(false),
			},
			HTMLURL: String("hurl"),
			ID:      Int64(1),
			NodeID:  String("nurl"),
			Number:  Int(1),
			Title:   String("title"),
			User: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			State:             String("st"),
			Locked:            Bool(false),
			Comments:          Int(1),
			CreatedAt:         &Timestamp{referenceTime},
			UpdatedAt:         &Timestamp{referenceTime},
			AuthorAssociation: String("aa"),
			Body:              String("bo"),
		},
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
		"comment": {
			"author_association": "aa",
			"body": "bo",
			"child_comment_count": 1,
			"created_at": ` + referenceTimeStr + `,
			"discussion_id": 1,
			"html_url": "hurl",
			"id": 1,
			"node_id": "nid",
			"parent_id": 1,
			"reactions": {
				"total_count": 1,
				"+1": 1,
				"-1": 1,
				"laugh": 1,
				"confused": 1,
				"heart": 1,
				"hooray": 1,
				"rocket": 1,
				"eyes": 1,
				"url": "url"
			},
			"repository_url": "rurl",
			"updated_at": ` + referenceTimeStr + `,
			"user": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			}
		},
		"discussion": {
			"repository_url": "rurl",
			"category": {
				"id": 1,
				"node_id": "nid",
				"repository_id": 1,
				"emoji": "emoji",
				"name": "name",
				"description": "description",
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"slug": "slug",
				"is_answerable": false
			},
			"html_url": "hurl",
			"id": 1,
			"node_id": "nurl",
			"number": 1,
			"title": "title",
			"user": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"state": "st",
			"locked": false,
			"comments": 1,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"author_association": "aa",
			"body": "bo"
		},
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

	testJSONMarshal(t, u, want)
}

func TestDiscussionEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DiscussionEvent{}, "{}")

	u := &DiscussionEvent{
		Discussion: &Discussion{
			RepositoryURL: String("rurl"),
			DiscussionCategory: &DiscussionCategory{
				ID:           Int64(1),
				NodeID:       String("nid"),
				RepositoryID: Int64(1),
				Emoji:        String("emoji"),
				Name:         String("name"),
				Description:  String("description"),
				CreatedAt:    &Timestamp{referenceTime},
				UpdatedAt:    &Timestamp{referenceTime},
				Slug:         String("slug"),
				IsAnswerable: Bool(false),
			},
			HTMLURL: String("hurl"),
			ID:      Int64(1),
			NodeID:  String("nurl"),
			Number:  Int(1),
			Title:   String("title"),
			User: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			State:             String("st"),
			Locked:            Bool(false),
			Comments:          Int(1),
			CreatedAt:         &Timestamp{referenceTime},
			UpdatedAt:         &Timestamp{referenceTime},
			AuthorAssociation: String("aa"),
			Body:              String("bo"),
		},
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
		"discussion": {
			"repository_url": "rurl",
			"category": {
				"id": 1,
				"node_id": "nid",
				"repository_id": 1,
				"emoji": "emoji",
				"name": "name",
				"description": "description",
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"slug": "slug",
				"is_answerable": false
			},
			"html_url": "hurl",
			"id": 1,
			"node_id": "nurl",
			"number": 1,
			"title": "title",
			"user": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"state": "st",
			"locked": false,
			"comments": 1,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"author_association": "aa",
			"body": "bo"
		},
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

	testJSONMarshal(t, u, want)
}

func TestPackageEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PackageEvent{}, "{}")

	u := &PackageEvent{
		Action: String("a"),
		Package: &Package{
			ID:          Int64(1),
			Name:        String("n"),
			PackageType: String("pt"),
			HTMLURL:     String("hurl"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			Owner: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			PackageVersion: &PackageVersion{ID: Int64(1)},
			Registry:       &PackageRegistry{Name: String("n")},
		},
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
		"action": "a",
		"package": {
			"id": 1,
			"name": "n",
			"package_type": "pt",
			"html_url": "hurl",
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"owner": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"package_version": {
				"id": 1
			},
			"registry": {
				"name": "n"
			}
		},
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

func TestPersonalAccessTokenRequestEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PersonalAccessTokenRequestEvent{}, "{}")

	event := &PersonalAccessTokenRequestEvent{
		Action: String("a"),
		PersonalAccessTokenRequest: &PersonalAccessTokenRequest{
			ID:    Int64(1),
			Owner: &User{Login: String("l")},
			PermissionsAdded: &PersonalAccessTokenPermissions{
				Org:  map[string]string{"organization_events": "read"},
				Repo: map[string]string{"security_events": "write"},
			},
			CreatedAt:           &Timestamp{referenceTime},
			TokenExpired:        Bool(false),
			TokenExpiresAt:      &Timestamp{referenceTime},
			TokenLastUsedAt:     &Timestamp{referenceTime},
			RepositoryCount:     Int64(1),
			RepositorySelection: String("rs"),
			Repositories: []*Repository{
				{
					Name: String("n"),
				},
			},
		},
		Org: &Organization{Name: String("n")},
		Sender: &User{
			Login: String("l"),
		},
		Installation: &Installation{
			ID: Int64(1),
		},
	}

	want := `{
		"action": "a",
		"personal_access_token_request": {
			"id": 1,
			"owner": {
				"login": "l"
			},
			"permissions_added": {
				"organization": {
					"organization_events": "read"
				},
				"repository": {
					"security_events": "write"
				}
			},
			"created_at": ` + referenceTimeStr + `,
			"token_expired": false,
			"token_expires_at": ` + referenceTimeStr + `,
			"token_last_used_at": ` + referenceTimeStr + `,
			"repository_count": 1,
			"repository_selection": "rs",
			"repositories": [
				{
					"name": "n"
				}
			]
		},
		"organization": {
			"name": "n"
		},
		"sender": {
			"login": "l"
		},
		"installation": {
			"id": 1
		}
	}`

	testJSONMarshal(t, event, want)
}

func TestPingEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PingEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"
	hookConfig := new(HookConfig)

	u := &PingEvent{
		Zen:    String("z"),
		HookID: Int64(1),
		Hook: &Hook{
			CreatedAt:    &Timestamp{referenceTime},
			UpdatedAt:    &Timestamp{referenceTime},
			URL:          String("url"),
			ID:           Int64(1),
			Type:         String("t"),
			Name:         String("n"),
			TestURL:      String("tu"),
			PingURL:      String("pu"),
			LastResponse: l,
			Config:       hookConfig,
			Events:       []string{"a"},
			Active:       Bool(true),
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
		"zen": "z",
		"hook_id": 1,
		"hook": {
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"url": "url",
			"id": 1,
			"type": "t",
			"name": "n",
			"test_url": "tu",
			"ping_url": "pu",
			"last_response": {
				"key": "value"
			},
			"config": {
				"key": "value"
			},
			"events": [
				"a"
			],
			"active": true
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

func TestRepositoryDispatchEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryDispatchEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &RepositoryDispatchEvent{
		Action:        String("a"),
		Branch:        String("b"),
		ClientPayload: jsonMsg,
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
		"action": "a",
		"branch": "b",
		"client_payload": {
			"key": "value"
		},
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

	testJSONMarshal(t, u, want)
}

func TestRepositoryImportEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryImportEvent{}, "{}")

	u := &RepositoryImportEvent{
		Status: String("success"),
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
		"status": "success",
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

func TestRepositoryEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryEvent{}, "{}")

	u := &RepositoryEvent{
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

	testJSONMarshal(t, u, want)
}

func TestReleaseEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ReleaseEvent{}, "{}")

	u := &ReleaseEvent{
		Action: String("a"),
		Release: &RepositoryRelease{
			Name:                   String("n"),
			DiscussionCategoryName: String("dcn"),
			ID:                     Int64(2),
			CreatedAt:              &Timestamp{referenceTime},
			PublishedAt:            &Timestamp{referenceTime},
			URL:                    String("url"),
			HTMLURL:                String("htmlurl"),
			AssetsURL:              String("assetsurl"),
			Assets:                 []*ReleaseAsset{{ID: Int64(1)}},
			UploadURL:              String("uploadurl"),
			ZipballURL:             String("zipballurl"),
			TarballURL:             String("tarballurl"),
			Author:                 &User{Name: String("octocat")},
			NodeID:                 String("nid"),
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
		"release": {
			"name": "n",
			"discussion_category_name": "dcn",
			"id": 2,
			"created_at": ` + referenceTimeStr + `,
			"published_at": ` + referenceTimeStr + `,
			"url": "url",
			"html_url": "htmlurl",
			"assets_url": "assetsurl",
			"assets": [
				{
					"id": 1
				}
			],
			"upload_url": "uploadurl",
			"zipball_url": "zipballurl",
			"tarball_url": "tarballurl",
			"author": {
				"name": "octocat"
			},
			"node_id": "nid"
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

func TestRepositoryRulesetEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryRulesetEvent{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &RepositoryRulesetEvent{
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
		Organization: &Organization{
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
		Repository: &Repository{
			ID:   Int64(1),
			URL:  String("u"),
			Name: String("n"),
		},
		RepositoryRuleset: &RepositoryRuleset{
			ID:          1,
			Name:        "n",
			Target:      String("branch"),
			SourceType:  String("Repository"),
			Source:      "s",
			Enforcement: "disabled",
			BypassActors: []*BypassActor{
				{
					ActorID:    Int64(234),
					ActorType:  String("Team"),
					BypassMode: String("Always"),
				},
			},
			CurrentUserCanBypass: String("always"),
			NodeID:               String("n"),
			Links: &RepositoryRulesetLink{
				Self: &RulesetLink{
					HRef: String("href"),
				},
				HTML: &RulesetLink{
					HRef: String("href"),
				},
			},
			Conditions: json.RawMessage(jsonMsg),
			Rules: []*RepositoryRulesetRule{
				{
					Creation: &RepositoryRulesetRuleType{
						Type: "creation",
					},
					Update: &RepositoryRulesetUpdateRule{
						Type: "update",
						Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
							UpdateAllowsFetchAndMerge: true,
						},
					},
					Deletion: &RepositoryRulesetRuleType{
						Type: "deletion",
					},
					RequireLinearHistory: &RepositoryRulesetRuleType{
						Type: "require_linear_history",
					},
					MergeQueue: &RepositoryRulesetMergeQueueRule{
						Type: "merge_queue",
						Parameters: &MergeQueueRuleParameters{
							CheckResponseTimeoutMinutes:  35,
							GroupingStrategy:             "HEADGREEN",
							MaxEntriesToBuild:            8,
							MaxEntriesToMerge:            4,
							MergeMethod:                  "SQUASH",
							MinEntriesToMerge:            2,
							MinEntriesToMergeWaitMinutes: 13,
						},
					},
					RequireDeployments: &RepositoryRulesetRequireDeploymentsRule{
						Type: "required_deployments",
						Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
							RequiredDeploymentEnvironments: []string{"test"},
						},
					},
					RequiredSignatures: &RepositoryRulesetRuleType{
						Type: "required_signatures",
					},
					PullRequest: &RepositoryRulesetPullRequestRule{
						Type: "pull_request",
						Parameters: &PullRequestRuleParameters{
							RequireCodeOwnerReview:         true,
							RequireLastPushApproval:        true,
							RequiredApprovingReviewCount:   1,
							RequiredReviewThreadResolution: true,
							DismissStaleReviewsOnPush:      true,
						},
					},
					RequiredStatusChecks: &RepositoryRulesetRequiredStatusChecksRule{
						Type: "required_status_checks",
						Parameters: &RequiredStatusChecksRuleParameters{
							RequiredStatusChecks: []RuleRequiredStatusChecks{
								{
									Context:       "test",
									IntegrationID: Int64(1),
								},
							},
							StrictRequiredStatusChecksPolicy: true,
						},
					},
					NonFastForward: &RepositoryRulesetRuleType{
						Type: "non_fast_forward",
					},
					CommitMessagePattern: &RepositoryRulesetPatternRule{
						Type: "commit_message_pattern",
						Parameters: &RulePatternParameters{
							Name:     String("avoid test commits"),
							Negate:   Bool(true),
							Operator: "starts_with",
							Pattern:  "[test]",
						},
					},
					CommitAuthorEmailPattern: &RepositoryRulesetPatternRule{
						Type: "commit_author_email_pattern",
						Parameters: &RulePatternParameters{
							Operator: "contains",
							Pattern:  "github",
						},
					},
					CommitterEmailPattern: &RepositoryRulesetPatternRule{
						Type: "committer_email_pattern",
						Parameters: &RulePatternParameters{
							Name:     String("avoid commit emails"),
							Negate:   Bool(true),
							Operator: "ends_with",
							Pattern:  "abc",
						},
					},
					BranchNamePattern: &RepositoryRulesetPatternRule{
						Type: "branch_name_pattern",
						Parameters: &RulePatternParameters{
							Name:     String("avoid branch names"),
							Negate:   Bool(true),
							Operator: "regex",
							Pattern:  "github$",
						},
					},
					TagNamePattern: &RepositoryRulesetPatternRule{
						Type: "tag_name_pattern",
						Parameters: &RulePatternParameters{
							Name:     String("avoid tag names"),
							Negate:   Bool(true),
							Operator: "contains",
							Pattern:  "github",
						},
					},
					FilePathRestriction: &RepositoryRulesetFilePathRestrictionRule{
						Type: "file_path_restriction",
						Parameters: &RuleFileParameters{
							RestrictedFilePaths: &[]string{"/a/file"},
						},
					},
					MaxFilePathLength: &RepositoryRulesetMaxFilePathLengthRule{
						Type: "max_file_path_length",
						Parameters: &RuleMaxFilePathLengthParameters{
							MaxFilePathLength: 255,
						},
					},
					FileExtensionRestriction: &RepositoryRulesetFileExtensionRestrictionRule{
						Type: "file_extension_restriction",
						Parameters: &RuleFileExtensionRestrictionParameters{
							RestrictedFileExtensions: []string{".exe"},
						},
					},
					MaxFileSize: &RepositoryRulesetMaxFileSizeRule{
						Type: "max_file_size",
						Parameters: &RuleMaxFileSizeParameters{
							MaxFileSize: 1024,
						},
					},
					Workflows: &RepositoryRulesetWorkflowsRule{
						Type: "workflows",
						Parameters: &RequiredWorkflowsRuleParameters{
							RequiredWorkflows: []*RuleRequiredWorkflow{
								{
									Path:         ".github/workflows/test.yml",
									RepositoryID: Int64(1),
								},
							},
						},
					},
					CodeScanning: &RepositoryRulesetCodeScanningRule{
						Type: "code_scanning",
						Parameters: &RuleCodeScanningParameters{
							CodeScanningTools: []CodeScanningTool{{
								AlertsThreshold:         "alert",
								SecurityAlertsThreshold: "security",
								Tool:                    "tool",
							}},
						},
					},
				},
			},
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
		},
		Changes: &RepositoryRulesetEditedChanges{
			Name: &RepositoryRulesetEditedSource{
				From: String("f"),
			},
			Enforcement: &RepositoryRulesetEditedSource{
				From: String("e"),
			},
			Conditions: &RepositoryRulesetEditedConditions{
				Added: []*RepositoryRulesetRefCondition{
					{
						RefName: &RulesetRefConditionParameters{
							Include: []string{"refs/heads/main", "refs/heads/master"},
							Exclude: []string{"refs/heads/dev*"},
						},
					},
				},
				Deleted: []*RepositoryRulesetRefCondition{
					{
						RefName: &RulesetRefConditionParameters{
							Include: []string{"refs/heads/main", "refs/heads/master"},
							Exclude: []string{"refs/heads/dev*"},
						},
					},
				},
				Updated: []*RepositoryRulesetEditedUpdatedConditions{
					{
						Condition: &RepositoryRulesetRefCondition{
							RefName: &RulesetRefConditionParameters{
								Include: []string{"refs/heads/main", "refs/heads/master"},
								Exclude: []string{"refs/heads/dev*"},
							},
						},
						Changes: &RepositoryRulesetUpdatedConditionsEdited{
							ConditionType: &RepositoryRulesetEditedSource{
								From: String("c"),
							},
							Target: &RepositoryRulesetEditedSource{
								From: String("t"),
							},
							Include: &RepositoryRulesetEditedSources{
								From: []string{"from"},
							},
							Exclude: &RepositoryRulesetEditedSources{
								From: []string{"to"},
							},
						},
					},
				},
			},
			Rules: &RepositoryRulesetEditedRules{
				Added: []*RepositoryRulesetRule{
					//Creating just one object with all the possible rules for testing
					{
						Creation: &RepositoryRulesetRuleType{
							Type: "creation",
						},
						Update: &RepositoryRulesetUpdateRule{
							Type: "update",
							Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
								UpdateAllowsFetchAndMerge: true,
							},
						},
						Deletion: &RepositoryRulesetRuleType{
							Type: "deletion",
						},
						RequireLinearHistory: &RepositoryRulesetRuleType{
							Type: "require_linear_history",
						},
						MergeQueue: &RepositoryRulesetMergeQueueRule{
							Type: "merge_queue",
							Parameters: &MergeQueueRuleParameters{
								CheckResponseTimeoutMinutes:  35,
								GroupingStrategy:             "HEADGREEN",
								MaxEntriesToBuild:            8,
								MaxEntriesToMerge:            4,
								MergeMethod:                  "SQUASH",
								MinEntriesToMerge:            2,
								MinEntriesToMergeWaitMinutes: 13,
							},
						},
						RequireDeployments: &RepositoryRulesetRequireDeploymentsRule{
							Type: "required_deployments",
							Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
								RequiredDeploymentEnvironments: []string{"test"},
							},
						},
						RequiredSignatures: &RepositoryRulesetRuleType{
							Type: "required_signatures",
						},
						PullRequest: &RepositoryRulesetPullRequestRule{
							Type: "pull_request",
							Parameters: &PullRequestRuleParameters{
								RequireCodeOwnerReview:         true,
								RequireLastPushApproval:        true,
								RequiredApprovingReviewCount:   1,
								RequiredReviewThreadResolution: true,
								DismissStaleReviewsOnPush:      true,
							},
						},
						RequiredStatusChecks: &RepositoryRulesetRequiredStatusChecksRule{
							Type: "required_status_checks",
							Parameters: &RequiredStatusChecksRuleParameters{
								RequiredStatusChecks: []RuleRequiredStatusChecks{
									{
										Context:       "test",
										IntegrationID: Int64(1),
									},
								},
								StrictRequiredStatusChecksPolicy: true,
							},
						},
						NonFastForward: &RepositoryRulesetRuleType{
							Type: "non_fast_forward",
						},
						CommitMessagePattern: &RepositoryRulesetPatternRule{
							Type: "commit_message_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid test commits"),
								Negate:   Bool(true),
								Operator: "starts_with",
								Pattern:  "[test]",
							},
						},
						CommitAuthorEmailPattern: &RepositoryRulesetPatternRule{
							Type: "commit_author_email_pattern",
							Parameters: &RulePatternParameters{
								Operator: "contains",
								Pattern:  "github",
							},
						},
						CommitterEmailPattern: &RepositoryRulesetPatternRule{
							Type: "committer_email_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid commit emails"),
								Negate:   Bool(true),
								Operator: "ends_with",
								Pattern:  "abc",
							},
						},
						BranchNamePattern: &RepositoryRulesetPatternRule{
							Type: "branch_name_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid branch names"),
								Negate:   Bool(true),
								Operator: "regex",
								Pattern:  "github$",
							},
						},
						TagNamePattern: &RepositoryRulesetPatternRule{
							Type: "tag_name_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid tag names"),
								Negate:   Bool(true),
								Operator: "contains",
								Pattern:  "github",
							},
						},
						FilePathRestriction: &RepositoryRulesetFilePathRestrictionRule{
							Type: "file_path_restriction",
							Parameters: &RuleFileParameters{
								RestrictedFilePaths: &[]string{"/a/file"},
							},
						},
						MaxFilePathLength: &RepositoryRulesetMaxFilePathLengthRule{
							Type: "max_file_path_length",
							Parameters: &RuleMaxFilePathLengthParameters{
								MaxFilePathLength: 255,
							},
						},
						FileExtensionRestriction: &RepositoryRulesetFileExtensionRestrictionRule{
							Type: "file_extension_restriction",
							Parameters: &RuleFileExtensionRestrictionParameters{
								RestrictedFileExtensions: []string{".exe"},
							},
						},
						MaxFileSize: &RepositoryRulesetMaxFileSizeRule{
							Type: "max_file_size",
							Parameters: &RuleMaxFileSizeParameters{
								MaxFileSize: 1024,
							},
						},
						Workflows: &RepositoryRulesetWorkflowsRule{
							Type: "workflows",
							Parameters: &RequiredWorkflowsRuleParameters{
								RequiredWorkflows: []*RuleRequiredWorkflow{
									{
										Path:         ".github/workflows/test.yml",
										RepositoryID: Int64(1),
									},
								},
							},
						},
						CodeScanning: &RepositoryRulesetCodeScanningRule{
							Type: "code_scanning",
							Parameters: &RuleCodeScanningParameters{
								CodeScanningTools: []CodeScanningTool{{
									AlertsThreshold:         "alert",
									SecurityAlertsThreshold: "security",
									Tool:                    "tool",
								}},
							},
						},
					},
				},
				Deleted: []*RepositoryRulesetRule{
					//Creating just one object with all the possible rules for testing
					{
						Creation: &RepositoryRulesetRuleType{
							Type: "creation",
						},
						Update: &RepositoryRulesetUpdateRule{
							Type: "update",
							Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
								UpdateAllowsFetchAndMerge: true,
							},
						},
						Deletion: &RepositoryRulesetRuleType{
							Type: "deletion",
						},
						RequireLinearHistory: &RepositoryRulesetRuleType{
							Type: "require_linear_history",
						},
						MergeQueue: &RepositoryRulesetMergeQueueRule{
							Type: "merge_queue",
							Parameters: &MergeQueueRuleParameters{
								CheckResponseTimeoutMinutes:  35,
								GroupingStrategy:             "HEADGREEN",
								MaxEntriesToBuild:            8,
								MaxEntriesToMerge:            4,
								MergeMethod:                  "SQUASH",
								MinEntriesToMerge:            2,
								MinEntriesToMergeWaitMinutes: 13,
							},
						},
						RequireDeployments: &RepositoryRulesetRequireDeploymentsRule{
							Type: "required_deployments",
							Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
								RequiredDeploymentEnvironments: []string{"test"},
							},
						},
						RequiredSignatures: &RepositoryRulesetRuleType{
							Type: "required_signatures",
						},
						PullRequest: &RepositoryRulesetPullRequestRule{
							Type: "pull_request",
							Parameters: &PullRequestRuleParameters{
								RequireCodeOwnerReview:         true,
								RequireLastPushApproval:        true,
								RequiredApprovingReviewCount:   1,
								RequiredReviewThreadResolution: true,
								DismissStaleReviewsOnPush:      true,
							},
						},
						RequiredStatusChecks: &RepositoryRulesetRequiredStatusChecksRule{
							Type: "required_status_checks",
							Parameters: &RequiredStatusChecksRuleParameters{
								RequiredStatusChecks: []RuleRequiredStatusChecks{
									{
										Context:       "test",
										IntegrationID: Int64(1),
									},
								},
								StrictRequiredStatusChecksPolicy: true,
							},
						},
						NonFastForward: &RepositoryRulesetRuleType{
							Type: "non_fast_forward",
						},
						CommitMessagePattern: &RepositoryRulesetPatternRule{
							Type: "commit_message_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid test commits"),
								Negate:   Bool(true),
								Operator: "starts_with",
								Pattern:  "[test]",
							},
						},
						CommitAuthorEmailPattern: &RepositoryRulesetPatternRule{
							Type: "commit_author_email_pattern",
							Parameters: &RulePatternParameters{
								Operator: "contains",
								Pattern:  "github",
							},
						},
						CommitterEmailPattern: &RepositoryRulesetPatternRule{
							Type: "committer_email_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid commit emails"),
								Negate:   Bool(true),
								Operator: "ends_with",
								Pattern:  "abc",
							},
						},
						BranchNamePattern: &RepositoryRulesetPatternRule{
							Type: "branch_name_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid branch names"),
								Negate:   Bool(true),
								Operator: "regex",
								Pattern:  "github$",
							},
						},
						TagNamePattern: &RepositoryRulesetPatternRule{
							Type: "tag_name_pattern",
							Parameters: &RulePatternParameters{
								Name:     String("avoid tag names"),
								Negate:   Bool(true),
								Operator: "contains",
								Pattern:  "github",
							},
						},
						FilePathRestriction: &RepositoryRulesetFilePathRestrictionRule{
							Type: "file_path_restriction",
							Parameters: &RuleFileParameters{
								RestrictedFilePaths: &[]string{"/a/file"},
							},
						},
						MaxFilePathLength: &RepositoryRulesetMaxFilePathLengthRule{
							Type: "max_file_path_length",
							Parameters: &RuleMaxFilePathLengthParameters{
								MaxFilePathLength: 255,
							},
						},
						FileExtensionRestriction: &RepositoryRulesetFileExtensionRestrictionRule{
							Type: "file_extension_restriction",
							Parameters: &RuleFileExtensionRestrictionParameters{
								RestrictedFileExtensions: []string{".exe"},
							},
						},
						MaxFileSize: &RepositoryRulesetMaxFileSizeRule{
							Type: "max_file_size",
							Parameters: &RuleMaxFileSizeParameters{
								MaxFileSize: 1024,
							},
						},
						Workflows: &RepositoryRulesetWorkflowsRule{
							Type: "workflows",
							Parameters: &RequiredWorkflowsRuleParameters{
								RequiredWorkflows: []*RuleRequiredWorkflow{
									{
										Path:         ".github/workflows/test.yml",
										RepositoryID: Int64(1),
									},
								},
							},
						},
						CodeScanning: &RepositoryRulesetCodeScanningRule{
							Type: "code_scanning",
							Parameters: &RuleCodeScanningParameters{
								CodeScanningTools: []CodeScanningTool{{
									AlertsThreshold:         "alert",
									SecurityAlertsThreshold: "security",
									Tool:                    "tool",
								}},
							},
						},
					},
				},
				Updated: []*RepositoryRulesetUpdatedRules{
					{
						Rule: &RepositoryRulesetRule{
							Creation: &RepositoryRulesetRuleType{
								Type: "creation",
							},
							Update: &RepositoryRulesetUpdateRule{
								Type: "update",
								Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
									UpdateAllowsFetchAndMerge: true,
								},
							},
							Deletion: &RepositoryRulesetRuleType{
								Type: "deletion",
							},
							RequireLinearHistory: &RepositoryRulesetRuleType{
								Type: "require_linear_history",
							},
							MergeQueue: &RepositoryRulesetMergeQueueRule{
								Type: "merge_queue",
								Parameters: &MergeQueueRuleParameters{
									CheckResponseTimeoutMinutes:  35,
									GroupingStrategy:             "HEADGREEN",
									MaxEntriesToBuild:            8,
									MaxEntriesToMerge:            4,
									MergeMethod:                  "SQUASH",
									MinEntriesToMerge:            2,
									MinEntriesToMergeWaitMinutes: 13,
								},
							},
							RequireDeployments: &RepositoryRulesetRequireDeploymentsRule{
								Type: "required_deployments",
								Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
									RequiredDeploymentEnvironments: []string{"test"},
								},
							},
							RequiredSignatures: &RepositoryRulesetRuleType{
								Type: "required_signatures",
							},
							PullRequest: &RepositoryRulesetPullRequestRule{
								Type: "pull_request",
								Parameters: &PullRequestRuleParameters{
									RequireCodeOwnerReview:         true,
									RequireLastPushApproval:        true,
									RequiredApprovingReviewCount:   1,
									RequiredReviewThreadResolution: true,
									DismissStaleReviewsOnPush:      true,
								},
							},
							RequiredStatusChecks: &RepositoryRulesetRequiredStatusChecksRule{
								Type: "required_status_checks",
								Parameters: &RequiredStatusChecksRuleParameters{
									RequiredStatusChecks: []RuleRequiredStatusChecks{
										{
											Context:       "test",
											IntegrationID: Int64(1),
										},
									},
									StrictRequiredStatusChecksPolicy: true,
								},
							},
							NonFastForward: &RepositoryRulesetRuleType{
								Type: "non_fast_forward",
							},
							CommitMessagePattern: &RepositoryRulesetPatternRule{
								Type: "commit_message_pattern",
								Parameters: &RulePatternParameters{
									Name:     String("avoid test commits"),
									Negate:   Bool(true),
									Operator: "starts_with",
									Pattern:  "[test]",
								},
							},
							CommitAuthorEmailPattern: &RepositoryRulesetPatternRule{
								Type: "commit_author_email_pattern",
								Parameters: &RulePatternParameters{
									Operator: "contains",
									Pattern:  "github",
								},
							},
							CommitterEmailPattern: &RepositoryRulesetPatternRule{
								Type: "committer_email_pattern",
								Parameters: &RulePatternParameters{
									Name:     String("avoid commit emails"),
									Negate:   Bool(true),
									Operator: "ends_with",
									Pattern:  "abc",
								},
							},
							BranchNamePattern: &RepositoryRulesetPatternRule{
								Type: "branch_name_pattern",
								Parameters: &RulePatternParameters{
									Name:     String("avoid branch names"),
									Negate:   Bool(true),
									Operator: "regex",
									Pattern:  "github$",
								},
							},
							TagNamePattern: &RepositoryRulesetPatternRule{
								Type: "tag_name_pattern",
								Parameters: &RulePatternParameters{
									Name:     String("avoid tag names"),
									Negate:   Bool(true),
									Operator: "contains",
									Pattern:  "github",
								},
							},
							FilePathRestriction: &RepositoryRulesetFilePathRestrictionRule{
								Type: "file_path_restriction",
								Parameters: &RuleFileParameters{
									RestrictedFilePaths: &[]string{"/a/file"},
								},
							},
							MaxFilePathLength: &RepositoryRulesetMaxFilePathLengthRule{
								Type: "max_file_path_length",
								Parameters: &RuleMaxFilePathLengthParameters{
									MaxFilePathLength: 255,
								},
							},
							FileExtensionRestriction: &RepositoryRulesetFileExtensionRestrictionRule{
								Type: "file_extension_restriction",
								Parameters: &RuleFileExtensionRestrictionParameters{
									RestrictedFileExtensions: []string{".exe"},
								},
							},
							MaxFileSize: &RepositoryRulesetMaxFileSizeRule{
								Type: "max_file_size",
								Parameters: &RuleMaxFileSizeParameters{
									MaxFileSize: 1024,
								},
							},
							Workflows: &RepositoryRulesetWorkflowsRule{
								Type: "workflows",
								Parameters: &RequiredWorkflowsRuleParameters{
									RequiredWorkflows: []*RuleRequiredWorkflow{
										{
											Path:         ".github/workflows/test.yml",
											RepositoryID: Int64(1),
										},
									},
								},
							},
							CodeScanning: &RepositoryRulesetCodeScanningRule{
								Type: "code_scanning",
								Parameters: &RuleCodeScanningParameters{
									CodeScanningTools: []CodeScanningTool{{
										AlertsThreshold:         "alert",
										SecurityAlertsThreshold: "security",
										Tool:                    "tool",
									}},
								},
							},
						},
						Changes: &RepositoryRulesetEditedRuleChanges{
							Configuration: &RepositoryRulesetEditedSources{
								From: []string{"from"},
							},
							RuleType: &RepositoryRulesetEditedSources{
								From: []string{"from"},
							},
							Pattern: &RepositoryRulesetEditedSources{
								From: []string{"from"},
							},
						},
					},
				},
			},
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
          "enterprise": {
            "id": 1,
            "slug": "s",
            "name": "n",
            "node_id": "nid",
            "avatar_url": "au",
            "description": "d",
            "website_url": "wu",
            "html_url": "hu",
            "created_at":  ` + referenceTimeStr + `,
            "updated_at":  ` + referenceTimeStr + `
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
              "created_at":  ` + referenceTimeStr + `,
              "suspended_at":  ` + referenceTimeStr + `,
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
            "created_at":  ` + referenceTimeStr + `,
            "updated_at":  ` + referenceTimeStr + `,
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
              "created_at":  ` + referenceTimeStr + `,
              "suspended_at":  ` + referenceTimeStr + `,
              "url": "u"
            },
            "suspended_at":  ` + referenceTimeStr + `
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
            "url": "u"
          },
          "repository_ruleset": {
            "id": 1,
            "name": "n",
            "target": "branch",
            "source_type": "Repository",
            "source": "s",
            "enforcement": "disabled",
            "bypass_actors": [
              {
                "actor_id": 234,
                "actor_type": "Team",
                "bypass_mode": "Always"
              }
            ],
            "current_user_can_bypass": "always",
            "node_id": "n",
            "_links": {
              "self": {
                "href": "href"
              },
              "html": {
                "href": "href"
              }
            },
            "conditions": {
              "key": "value"
            },
            "rules": [
              {
                "creation": {
                  "type": "creation"
                },
                "update": {
                  "type": "update",
                  "parameters": {
                    "update_allows_fetch_and_merge": true
                  }
                },
                "deletion": {
                  "type": "deletion"
                },
                "required_linear_history": {
                  "type": "require_linear_history"
                },
                "merge_queue": {
                  "type": "merge_queue",
                  "parameters": {
                    "check_response_timeout_minutes": 35,
                    "grouping_strategy": "HEADGREEN",
                    "max_entries_to_build": 8,
                    "max_entries_to_merge": 4,
                    "merge_method": "SQUASH",
                    "min_entries_to_merge": 2,
                    "min_entries_to_merge_wait_minutes": 13
                  }
                },
                "required_deployments": {
                  "type": "required_deployments",
                  "parameters": {
                    "required_deployment_environments": [
                      "test"
                    ]
                  }
                },
                "required_signatures": {
                  "type": "required_signatures"
                },
                "pull_request": {
                  "type": "pull_request",
                  "parameters": {
                    "dismiss_stale_reviews_on_push": true,
                    "require_code_owner_review": true,
                    "require_last_push_approval": true,
                    "required_approving_review_count": 1,
                    "required_review_thread_resolution": true
                  }
                },
                "required_status_checks": {
                  "type": "required_status_checks",
                  "parameters": {
                    "required_status_checks": [
                      {
                        "context": "test",
                        "integration_id": 1
                      }
                    ],
                    "strict_required_status_checks_policy": true
                  }
                },
                "non_fast_forward": {
                  "type": "non_fast_forward"
                },
                "commit_message_pattern": {
                  "type": "commit_message_pattern",
                  "parameters": {
                    "name": "avoid test commits",
                    "negate": true,
                    "operator": "starts_with",
                    "pattern": "[test]"
                  }
                },
                "commit_author_email_pattern": {
                  "type": "commit_author_email_pattern",
                  "parameters": {
                    "operator": "contains",
                    "pattern": "github"
                  }
                },
                "committer_email_pattern": {
                  "type": "committer_email_pattern",
                  "parameters": {
                    "name": "avoid commit emails",
                    "negate": true,
                    "operator": "ends_with",
                    "pattern": "abc"
                  }
                },
                "branch_name_pattern": {
                  "type": "branch_name_pattern",
                  "parameters": {
                    "name": "avoid branch names",
                    "negate": true,
                    "operator": "regex",
                    "pattern": "github$"
                  }
                },
                "tag_name_pattern": {
                  "type": "tag_name_pattern",
                  "parameters": {
                    "name": "avoid tag names",
                    "negate": true,
                    "operator": "contains",
                    "pattern": "github"
                  }
                },
                "file_path_restriction": {
                  "type": "file_path_restriction",
                  "parameters": {
                    "restricted_file_paths": [
                      "/a/file"
                    ]
                  }
                },
                "max_file_path_length": {
                  "type": "max_file_path_length",
                  "parameters": {
                    "max_file_path_length": 255
                  }
                },
                "file_extension_restriction": {
                  "type": "file_extension_restriction",
                  "parameters": {
                    "restricted_file_extensions": [
                      ".exe"
                    ]
                  }
                },
                "max_file_size": {
                  "type": "max_file_size",
                  "parameters": {
                    "max_file_size": 1024
                  }
                },
                "workflows": {
                  "type": "workflows",
                  "parameters": {
                    "workflows": [
                      {
                        "path": ".github/workflows/test.yml",
                        "repository_id": 1
                      }
                    ]
                  }
                },
                "code_scanning": {
                  "type": "code_scanning",
                  "parameters": {
                    "code_scanning_tools": [
                      {
                        "alerts_threshold": "alert",
                        "security_alerts_threshold": "security",
                        "tool": "tool"
                      }
                    ]
                  }
                }
              }
            ],
            "created_at":  ` + referenceTimeStr + `,
            "updated_at":  ` + referenceTimeStr + `
          },
          "changes": {
            "name": {
              "from": "f"
            },
            "enforcement": {
              "from": "e"
            },
            "conditions": {
              "added": [
                {
                  "ref_name": {
                    "include": [
                      "refs/heads/main",
                      "refs/heads/master"
                    ],
                    "exclude": [
                      "refs/heads/dev*"
                    ]
                  }
                }
              ],
              "deleted": [
                {
                  "ref_name": {
                    "include": [
                      "refs/heads/main",
                      "refs/heads/master"
                    ],
                    "exclude": [
                      "refs/heads/dev*"
                    ]
                  }
                }
              ],
              "updated": [
                {
                  "condition": {
                    "ref_name": {
                      "include": [
                        "refs/heads/main",
                        "refs/heads/master"
                      ],
                      "exclude": [
                        "refs/heads/dev*"
                      ]
                    }
                  },
                  "changes": {
                    "condition_type": {
                      "from": "c"
                    },
                    "target": {
                      "from": "t"
                    },
                    "include": {
                      "from": [
                        "from"
                      ]
                    },
                    "exclude": {
                      "from": [
                        "to"
                      ]
                    }
                  }
                }
              ]
            },
            "rules": {
              "added": [
                {
                  "creation": {
                    "type": "creation"
                  },
                  "update": {
                    "type": "update",
                    "parameters": {
                      "update_allows_fetch_and_merge": true
                    }
                  },
                  "deletion": {
                    "type": "deletion"
                  },
                  "required_linear_history": {
                    "type": "require_linear_history"
                  },
                  "merge_queue": {
                    "type": "merge_queue",
                    "parameters": {
                      "check_response_timeout_minutes": 35,
                      "grouping_strategy": "HEADGREEN",
                      "max_entries_to_build": 8,
                      "max_entries_to_merge": 4,
                      "merge_method": "SQUASH",
                      "min_entries_to_merge": 2,
                      "min_entries_to_merge_wait_minutes": 13
                    }
                  },
                  "required_deployments": {
                    "type": "required_deployments",
                    "parameters": {
                      "required_deployment_environments": [
                        "test"
                      ]
                    }
                  },
                  "required_signatures": {
                    "type": "required_signatures"
                  },
                  "pull_request": {
                    "type": "pull_request",
                    "parameters": {
                      "dismiss_stale_reviews_on_push": true,
                      "require_code_owner_review": true,
                      "require_last_push_approval": true,
                      "required_approving_review_count": 1,
                      "required_review_thread_resolution": true
                    }
                  },
                  "required_status_checks": {
                    "type": "required_status_checks",
                    "parameters": {
                      "required_status_checks": [
                        {
                          "context": "test",
                          "integration_id": 1
                        }
                      ],
                      "strict_required_status_checks_policy": true
                    }
                  },
                  "non_fast_forward": {
                    "type": "non_fast_forward"
                  },
                  "commit_message_pattern": {
                    "type": "commit_message_pattern",
                    "parameters": {
                      "name": "avoid test commits",
                      "negate": true,
                      "operator": "starts_with",
                      "pattern": "[test]"
                    }
                  },
                  "commit_author_email_pattern": {
                    "type": "commit_author_email_pattern",
                    "parameters": {
                      "operator": "contains",
                      "pattern": "github"
                    }
                  },
                  "committer_email_pattern": {
                    "type": "committer_email_pattern",
                    "parameters": {
                      "name": "avoid commit emails",
                      "negate": true,
                      "operator": "ends_with",
                      "pattern": "abc"
                    }
                  },
                  "branch_name_pattern": {
                    "type": "branch_name_pattern",
                    "parameters": {
                      "name": "avoid branch names",
                      "negate": true,
                      "operator": "regex",
                      "pattern": "github$"
                    }
                  },
                  "tag_name_pattern": {
                    "type": "tag_name_pattern",
                    "parameters": {
                      "name": "avoid tag names",
                      "negate": true,
                      "operator": "contains",
                      "pattern": "github"
                    }
                  },
                  "file_path_restriction": {
                    "type": "file_path_restriction",
                    "parameters": {
                      "restricted_file_paths": [
                        "/a/file"
                      ]
                    }
                  },
                  "max_file_path_length": {
                    "type": "max_file_path_length",
                    "parameters": {
                      "max_file_path_length": 255
                    }
                  },
                  "file_extension_restriction": {
                    "type": "file_extension_restriction",
                    "parameters": {
                      "restricted_file_extensions": [
                        ".exe"
                      ]
                    }
                  },
                  "max_file_size": {
                    "type": "max_file_size",
                    "parameters": {
                      "max_file_size": 1024
                    }
                  },
                  "workflows": {
                    "type": "workflows",
                    "parameters": {
                      "workflows": [
                        {
                          "path": ".github/workflows/test.yml",
                          "repository_id": 1
                        }
                      ]
                    }
                  },
                  "code_scanning": {
                    "type": "code_scanning",
                    "parameters": {
                      "code_scanning_tools": [
                        {
                          "alerts_threshold": "alert",
                          "security_alerts_threshold": "security",
                          "tool": "tool"
                        }
                      ]
                    }
                  }
                }
              ],
              "deleted": [
                {
                  "creation": {
                    "type": "creation"
                  },
                  "update": {
                    "type": "update",
                    "parameters": {
                      "update_allows_fetch_and_merge": true
                    }
                  },
                  "deletion": {
                    "type": "deletion"
                  },
                  "required_linear_history": {
                    "type": "require_linear_history"
                  },
                  "merge_queue": {
                    "type": "merge_queue",
                    "parameters": {
                      "check_response_timeout_minutes": 35,
                      "grouping_strategy": "HEADGREEN",
                      "max_entries_to_build": 8,
                      "max_entries_to_merge": 4,
                      "merge_method": "SQUASH",
                      "min_entries_to_merge": 2,
                      "min_entries_to_merge_wait_minutes": 13
                    }
                  },
                  "required_deployments": {
                    "type": "required_deployments",
                    "parameters": {
                      "required_deployment_environments": [
                        "test"
                      ]
                    }
                  },
                  "required_signatures": {
                    "type": "required_signatures"
                  },
                  "pull_request": {
                    "type": "pull_request",
                    "parameters": {
                      "dismiss_stale_reviews_on_push": true,
                      "require_code_owner_review": true,
                      "require_last_push_approval": true,
                      "required_approving_review_count": 1,
                      "required_review_thread_resolution": true
                    }
                  },
                  "required_status_checks": {
                    "type": "required_status_checks",
                    "parameters": {
                      "required_status_checks": [
                        {
                          "context": "test",
                          "integration_id": 1
                        }
                      ],
                      "strict_required_status_checks_policy": true
                    }
                  },
                  "non_fast_forward": {
                    "type": "non_fast_forward"
                  },
                  "commit_message_pattern": {
                    "type": "commit_message_pattern",
                    "parameters": {
                      "name": "avoid test commits",
                      "negate": true,
                      "operator": "starts_with",
                      "pattern": "[test]"
                    }
                  },
                  "commit_author_email_pattern": {
                    "type": "commit_author_email_pattern",
                    "parameters": {
                      "operator": "contains",
                      "pattern": "github"
                    }
                  },
                  "committer_email_pattern": {
                    "type": "committer_email_pattern",
                    "parameters": {
                      "name": "avoid commit emails",
                      "negate": true,
                      "operator": "ends_with",
                      "pattern": "abc"
                    }
                  },
                  "branch_name_pattern": {
                    "type": "branch_name_pattern",
                    "parameters": {
                      "name": "avoid branch names",
                      "negate": true,
                      "operator": "regex",
                      "pattern": "github$"
                    }
                  },
                  "tag_name_pattern": {
                    "type": "tag_name_pattern",
                    "parameters": {
                      "name": "avoid tag names",
                      "negate": true,
                      "operator": "contains",
                      "pattern": "github"
                    }
                  },
                  "file_path_restriction": {
                    "type": "file_path_restriction",
                    "parameters": {
                      "restricted_file_paths": [
                        "/a/file"
                      ]
                    }
                  },
                  "max_file_path_length": {
                    "type": "max_file_path_length",
                    "parameters": {
                      "max_file_path_length": 255
                    }
                  },
                  "file_extension_restriction": {
                    "type": "file_extension_restriction",
                    "parameters": {
                      "restricted_file_extensions": [
                        ".exe"
                      ]
                    }
                  },
                  "max_file_size": {
                    "type": "max_file_size",
                    "parameters": {
                      "max_file_size": 1024
                    }
                  },
                  "workflows": {
                    "type": "workflows",
                    "parameters": {
                      "workflows": [
                        {
                          "path": ".github/workflows/test.yml",
                          "repository_id": 1
                        }
                      ]
                    }
                  },
                  "code_scanning": {
                    "type": "code_scanning",
                    "parameters": {
                      "code_scanning_tools": [
                        {
                          "alerts_threshold": "alert",
                          "security_alerts_threshold": "security",
                          "tool": "tool"
                        }
                      ]
                    }
                  }
                }
              ],
              "updated": [
                {
                  "rule": {
                    "creation": {
                      "type": "creation"
                    },
                    "update": {
                      "type": "update",
                      "parameters": {
                        "update_allows_fetch_and_merge": true
                      }
                    },
                    "deletion": {
                      "type": "deletion"
                    },
                    "required_linear_history": {
                      "type": "require_linear_history"
                    },
                    "merge_queue": {
                      "type": "merge_queue",
                      "parameters": {
                        "check_response_timeout_minutes": 35,
                        "grouping_strategy": "HEADGREEN",
                        "max_entries_to_build": 8,
                        "max_entries_to_merge": 4,
                        "merge_method": "SQUASH",
                        "min_entries_to_merge": 2,
                        "min_entries_to_merge_wait_minutes": 13
                      }
                    },
                    "required_deployments": {
                      "type": "required_deployments",
                      "parameters": {
                        "required_deployment_environments": [
                          "test"
                        ]
                      }
                    },
                    "required_signatures": {
                      "type": "required_signatures"
                    },
                    "pull_request": {
                      "type": "pull_request",
                      "parameters": {
                        "dismiss_stale_reviews_on_push": true,
                        "require_code_owner_review": true,
                        "require_last_push_approval": true,
                        "required_approving_review_count": 1,
                        "required_review_thread_resolution": true
                      }
                    },
                    "required_status_checks": {
                      "type": "required_status_checks",
                      "parameters": {
                        "required_status_checks": [
                          {
                            "context": "test",
                            "integration_id": 1
                          }
                        ],
                        "strict_required_status_checks_policy": true
                      }
                    },
                    "non_fast_forward": {
                      "type": "non_fast_forward"
                    },
                    "commit_message_pattern": {
                      "type": "commit_message_pattern",
                      "parameters": {
                        "name": "avoid test commits",
                        "negate": true,
                        "operator": "starts_with",
                        "pattern": "[test]"
                      }
                    },
                    "commit_author_email_pattern": {
                      "type": "commit_author_email_pattern",
                      "parameters": {
                        "operator": "contains",
                        "pattern": "github"
                      }
                    },
                    "committer_email_pattern": {
                      "type": "committer_email_pattern",
                      "parameters": {
                        "name": "avoid commit emails",
                        "negate": true,
                        "operator": "ends_with",
                        "pattern": "abc"
                      }
                    },
                    "branch_name_pattern": {
                      "type": "branch_name_pattern",
                      "parameters": {
                        "name": "avoid branch names",
                        "negate": true,
                        "operator": "regex",
                        "pattern": "github$"
                      }
                    },
                    "tag_name_pattern": {
                      "type": "tag_name_pattern",
                      "parameters": {
                        "name": "avoid tag names",
                        "negate": true,
                        "operator": "contains",
                        "pattern": "github"
                      }
                    },
                    "file_path_restriction": {
                      "type": "file_path_restriction",
                      "parameters": {
                        "restricted_file_paths": [
                          "/a/file"
                        ]
                      }
                    },
                    "max_file_path_length": {
                      "type": "max_file_path_length",
                      "parameters": {
                        "max_file_path_length": 255
                      }
                    },
                    "file_extension_restriction": {
                      "type": "file_extension_restriction",
                      "parameters": {
                        "restricted_file_extensions": [
                          ".exe"
                        ]
                      }
                    },
                    "max_file_size": {
                      "type": "max_file_size",
                      "parameters": {
                        "max_file_size": 1024
                      }
                    },
                    "workflows": {
                      "type": "workflows",
                      "parameters": {
                        "workflows": [
                          {
                            "path": ".github/workflows/test.yml",
                            "repository_id": 1
                          }
                        ]
                      }
                    },
                    "code_scanning": {
                      "type": "code_scanning",
                      "parameters": {
                        "code_scanning_tools": [
                          {
                            "alerts_threshold": "alert",
                            "security_alerts_threshold": "security",
                            "tool": "tool"
                          }
                        ]
                      }
                    }
                  },
                  "changes": {
                    "configuration": {
                      "from": [
                        "from"
                      ]
                    },
                    "rule_type": {
                      "from": [
                        "from"
                      ]
                    },
                    "pattern": {
                      "from": [
                        "from"
                      ]
                    }
                  }
                }
              ]
            }
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

func TestContentReferenceEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ContentReferenceEvent{}, "{}")

	u := &ContentReferenceEvent{
		Action: String("a"),
		ContentReference: &ContentReference{
			ID:        Int64(1),
			NodeID:    String("nid"),
			Reference: String("ref"),
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
		"content_reference": {
			"id": 1,
			"node_id": "nid",
			"reference": "ref"
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

func TestMemberEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MemberEvent{}, "{}")

	u := &MemberEvent{
		Action: String("a"),
		Member: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Changes: &MemberChanges{
			Permission: &MemberChangesPermission{
				From: String("f"),
				To:   String("t"),
			},
			RoleName: &MemberChangesRoleName{
				From: String("f"),
				To:   String("t"),
			},
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
		"member": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"changes": {
			"permission": {
				"from": "f",
				"to": "t"
			},
			"role_name": {
				"from": "f",
				"to": "t"
			}
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

func TestMembershipEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MembershipEvent{}, "{}")

	u := &MembershipEvent{
		Action: String("a"),
		Scope:  String("s"),
		Member: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Team: &Team{
			ID:              Int64(1),
			NodeID:          String("n"),
			Name:            String("n"),
			Description:     String("d"),
			URL:             String("u"),
			Slug:            String("s"),
			Permission:      String("p"),
			Privacy:         String("p"),
			MembersCount:    Int(1),
			ReposCount:      Int(1),
			MembersURL:      String("m"),
			RepositoriesURL: String("r"),
			Organization: &Organization{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				AvatarURL: String("a"),
				HTMLURL:   String("h"),
				Name:      String("n"),
				Company:   String("c"),
				Blog:      String("b"),
				Location:  String("l"),
				Email:     String("e"),
			},
			Parent: &Team{
				ID:           Int64(1),
				NodeID:       String("n"),
				Name:         String("n"),
				Description:  String("d"),
				URL:          String("u"),
				Slug:         String("s"),
				Permission:   String("p"),
				Privacy:      String("p"),
				MembersCount: Int(1),
				ReposCount:   Int(1),
			},
			LDAPDN: String("l"),
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
		"action": "a",
		"scope": "s",
		"member": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"team": {
			"id": 1,
			"node_id": "n",
			"name": "n",
			"description": "d",
			"url": "u",
			"slug": "s",
			"permission": "p",
			"privacy": "p",
			"members_count": 1,
			"repos_count": 1,
			"organization": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"html_url": "h",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e"
			},
			"members_url": "m",
			"repositories_url": "r",
			"parent": {
				"id": 1,
				"node_id": "n",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "p",
				"members_count": 1,
				"repos_count": 1
			},
			"ldap_dn": "l"
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

	testJSONMarshal(t, u, want)
}

func TestMergeGroupEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MergeGroupEvent{}, "{}")

	u := &MergeGroupEvent{
		Action: String("a"),
		MergeGroup: &MergeGroup{
			HeadSHA:    String("hs"),
			HeadRef:    String("hr"),
			BaseSHA:    String("bs"),
			BaseRef:    String("br"),
			HeadCommit: &Commit{NodeID: String("nid")},
		},
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
		"action": "a",
		"merge_group": {
			"head_sha": "hs",
			"head_ref": "hr",
			"base_sha": "bs",
			"base_ref": "br",
			"head_commit": {
				"node_id": "nid"
			}
		},
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

	testJSONMarshal(t, u, want)
}

func TestOrgBlockEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &OrgBlockEvent{}, "{}")

	u := &OrgBlockEvent{
		Action: String("a"),
		BlockedUser: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Organization: &Organization{
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
		"action": "a",
		"blocked_user": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
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

	testJSONMarshal(t, u, want)
}

func TestGollumEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &GollumEvent{}, "{}")

	u := &GollumEvent{
		Pages: []*Page{
			{
				PageName: String("pn"),
				Title:    String("t"),
				Summary:  String("s"),
				Action:   String("a"),
				SHA:      String("sha"),
				HTMLURL:  String("hu"),
			},
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
		"pages": [
			{
				"page_name": "pn",
				"title": "t",
				"summary": "s",
				"action": "a",
				"sha": "sha",
				"html_url": "hu"
			}
		],
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
			RunAttempt: Int(1),
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
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            String("j"),
			LogsURL:            String("l"),
			CheckSuiteURL:      String("c"),
			ArtifactsURL:       String("a"),
			CancelURL:          String("c"),
			RerunURL:           String("r"),
			PreviousAttemptURL: String("p"),
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
			"run_attempt": 1,
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
			"run_started_at": ` + referenceTimeStr + `,
			"jobs_url": "j",
			"logs_url": "l",
			"check_suite_url": "c",
			"artifacts_url": "a",
			"cancel_url": "c",
			"rerun_url": "r",
			"previous_attempt_url": "p",
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

func TestDeployKeyEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeployKeyEvent{}, "{}")

	u := &DeployKeyEvent{
		Action: String("a"),
		Key: &Key{
			ID:        Int64(1),
			Key:       String("k"),
			URL:       String("k"),
			Title:     String("k"),
			ReadOnly:  Bool(false),
			Verified:  Bool(false),
			CreatedAt: &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Organization: &Organization{
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
			AvatarURL: String("a"),
			URL:       String("u"),
			EventsURL: String("e"),
			ReposURL:  String("r"),
		},
	}

	want := `{
		"action": "a",
		"key": {
			"id": 1,
			"key": "k",
			"url": "k",
			"title": "k",
			"read_only": false,
			"verified": false,
			"created_at": ` + referenceTimeStr + `
		},
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

func TestMetaEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &MetaEvent{}, "{}")

	v := make(map[string]interface{})
	v["a"] = "b"
	hookConfig := &HookConfig{
		ContentType: String("json"),
	}

	u := &MetaEvent{
		Action: String("a"),
		HookID: Int64(1),
		Hook: &Hook{
			CreatedAt:    &Timestamp{referenceTime},
			UpdatedAt:    &Timestamp{referenceTime},
			URL:          String("u"),
			ID:           Int64(1),
			Type:         String("t"),
			Name:         String("n"),
			TestURL:      String("tu"),
			PingURL:      String("pu"),
			LastResponse: v,
			Config:       hookConfig,
			Events:       []string{"a"},
			Active:       Bool(true),
		},
	}

	want := `{
		"action": "a",
		"hook_id": 1,
		"hook": {
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"url": "u",
			"id": 1,
			"type": "t",
			"name": "n",
			"test_url": "tu",
			"ping_url": "pu",
			"last_response": {
				"a": "b"
			},
			"config": {
				"content_type": "json"
			},
			"events": [
				"a"
			],
			"active": true
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestRequestedAction_Marshal(t *testing.T) {
	testJSONMarshal(t, &RequestedAction{}, "{}")

	r := &RequestedAction{
		Identifier: "i",
	}

	want := `{
		"identifier": "i"
	}`

	testJSONMarshal(t, r, want)
}

func TestCreateEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &CreateEvent{}, "{}")

	r := &CreateEvent{
		Ref:          String("r"),
		RefType:      String("rt"),
		MasterBranch: String("mb"),
		Description:  String("d"),
		PusherType:   String("pt"),
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
		"ref": "r",
		"ref_type": "rt",
		"master_branch": "mb",
		"description": "d",
		"pusher_type": "pt",
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

	testJSONMarshal(t, r, want)
}

func TestDeleteEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DeleteEvent{}, "{}")

	r := &DeleteEvent{
		Ref:        String("r"),
		RefType:    String("rt"),
		PusherType: String("pt"),
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
		"ref": "r",
		"ref_type": "rt",
		"pusher_type": "pt",
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

	testJSONMarshal(t, r, want)
}

func TestDependabotAlertEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependabotAlertEvent{}, "{}")

	e := &DependabotAlertEvent{
		Action: String("a"),
		Alert: &DependabotAlert{
			Number: Int(1),
			State:  String("s"),
			Dependency: &Dependency{
				Package: &VulnerabilityPackage{
					Ecosystem: String("e"),
					Name:      String("n"),
				},
				ManifestPath: String("mp"),
				Scope:        String("s"),
			},
			SecurityAdvisory: &DependabotSecurityAdvisory{
				GHSAID:      String("ghsaid"),
				CVEID:       String("cveid"),
				Summary:     String("s"),
				Description: String("d"),
				Vulnerabilities: []*AdvisoryVulnerability{
					{
						Package: &VulnerabilityPackage{
							Ecosystem: String("e"),
							Name:      String("n"),
						},
						Severity: String("s"),
					},
				},
				Severity: String("s"),
				CVSS: &AdvisoryCVSS{
					Score:        Float64(1.0),
					VectorString: String("vs"),
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: String("cweid"),
						Name:  String("n"),
					},
				},
				Identifiers: []*AdvisoryIdentifier{
					{
						Value: String("v"),
						Type:  String("t"),
					},
				},
				References: []*AdvisoryReference{
					{
						URL: String("u"),
					},
				},
				PublishedAt: &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
				WithdrawnAt: &Timestamp{referenceTime},
			},
			SecurityVulnerability: &AdvisoryVulnerability{
				Package: &VulnerabilityPackage{
					Ecosystem: String("e"),
					Name:      String("n"),
				},
				Severity:               String("s"),
				VulnerableVersionRange: String("vvr"),
				FirstPatchedVersion: &FirstPatchedVersion{
					Identifier: String("i"),
				},
			},
			URL:         String("u"),
			HTMLURL:     String("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			DismissedAt: &Timestamp{referenceTime},
			DismissedBy: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			DismissedReason:  String("dr"),
			DismissedComment: String("dc"),
			FixedAt:          &Timestamp{referenceTime},
			AutoDismissedAt:  &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Organization: &Organization{
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
		"alert": {
			"number": 1,
			"state": "s",
			"dependency": {
				"package": {
					"ecosystem": "e",
					"name": "n"
				},
				"manifest_path": "mp",
				"scope": "s"
			},
			"security_advisory": {
				"ghsa_id": "ghsaid",
				"cve_id": "cveid",
				"summary": "s",
				"description": "d",
				"vulnerabilities": [
					{
						"package": {
							"ecosystem": "e",
							"name": "n"
						},
						"severity": "s"
					}
				],
				"severity": "s",
				"cvss": {
					"score": 1.0,
					"vector_string": "vs"
				},
				"cwes": [
					{
						"cwe_id": "cweid",
						"name": "n"
					}
				],
				"identifiers": [
					{
						"value": "v",
						"type": "t"
					}
				],
				"references": [
					{
						"url": "u"
					}
				],
				"published_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"withdrawn_at": ` + referenceTimeStr + `
			},
			"security_vulnerability": {
				"package": {
					"ecosystem": "e",
					"name": "n"
				},
				"severity": "s",
				"vulnerable_version_range": "vvr",
				"first_patched_version": {
					"identifier": "i"
				}
			},
			"url": "u",
			"html_url": "hu",
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"dismissed_at": ` + referenceTimeStr + `,
			"dismissed_by": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"dismissed_reason": "dr",
			"dismissed_comment": "dc",
			"fixed_at": ` + referenceTimeStr + `,
			"auto_dismissed_at": ` + referenceTimeStr + `
		},
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

	testJSONMarshal(t, e, want)
}

func TestForkEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ForkEvent{}, "{}")

	u := &ForkEvent{
		Forkee: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
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
		"forkee": {
			"id": 1,
			"name": "n",
			"url": "s"
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

func TestGitHubAppAuthorizationEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &GitHubAppAuthorizationEvent{}, "{}")

	u := &GitHubAppAuthorizationEvent{
		Action: String("a"),
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

func TestInstallationEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &InstallationEvent{}, "{}")

	u := &InstallationEvent{
		Action: String("a"),
		Repositories: []*Repository{
			{
				ID:   Int64(1),
				URL:  String("u"),
				Name: String("n"),
			},
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
		"repositories": [
			{
				"id":1,
				"name":"n",
				"url":"u"
			}
		],
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

func TestHeadCommit_Marshal(t *testing.T) {
	testJSONMarshal(t, &HeadCommit{}, "{}")

	u := &HeadCommit{
		Message: String("m"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		URL:       String("u"),
		Distinct:  Bool(true),
		SHA:       String("s"),
		ID:        String("id"),
		TreeID:    String("tid"),
		Timestamp: &Timestamp{referenceTime},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		Added:    []string{"a"},
		Removed:  []string{"r"},
		Modified: []string{"m"},
	}

	want := `{
		"message": "m",
		"author": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"url": "u",
		"distinct": true,
		"sha": "s",
		"id": "id",
		"tree_id": "tid",
		"timestamp": ` + referenceTimeStr + `,
		"committer": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"added": [
			"a"
		],
		"removed":  [
			"r"
		],
		"modified":  [
			"m"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestPushEventRepository_Marshal(t *testing.T) {
	testJSONMarshal(t, &PushEventRepository{}, "{}")

	u := &PushEventRepository{
		ID:       Int64(1),
		NodeID:   String("nid"),
		Name:     String("n"),
		FullName: String("fn"),
		Owner: &User{
			Login:       String("l"),
			ID:          Int64(1),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		},
		Private:         Bool(true),
		Description:     String("d"),
		Fork:            Bool(true),
		CreatedAt:       &Timestamp{referenceTime},
		PushedAt:        &Timestamp{referenceTime},
		UpdatedAt:       &Timestamp{referenceTime},
		Homepage:        String("h"),
		PullsURL:        String("p"),
		Size:            Int(1),
		StargazersCount: Int(1),
		WatchersCount:   Int(1),
		Language:        String("l"),
		HasIssues:       Bool(true),
		HasDownloads:    Bool(true),
		HasWiki:         Bool(true),
		HasPages:        Bool(true),
		ForksCount:      Int(1),
		Archived:        Bool(true),
		Disabled:        Bool(true),
		OpenIssuesCount: Int(1),
		DefaultBranch:   String("d"),
		MasterBranch:    String("m"),
		Organization:    String("o"),
		URL:             String("u"),
		ArchiveURL:      String("a"),
		HTMLURL:         String("h"),
		StatusesURL:     String("s"),
		GitURL:          String("g"),
		SSHURL:          String("s"),
		CloneURL:        String("c"),
		SVNURL:          String("s"),
		Topics:          []string{"octocat", "api"},
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"name": "n",
		"full_name": "fn",
		"owner": {
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
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"private": true,
		"description": "d",
		"fork": true,
		"created_at": ` + referenceTimeStr + `,
		"pushed_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"homepage": "h",
		"pulls_url": "p",
		"size": 1,
		"stargazers_count": 1,
		"watchers_count": 1,
		"language": "l",
		"has_issues": true,
		"has_downloads": true,
		"has_wiki": true,
		"has_pages": true,
		"forks_count": 1,
		"archived": true,
		"disabled": true,
		"open_issues_count": 1,
		"default_branch": "d",
		"master_branch": "m",
		"organization": "o",
		"url": "u",
		"archive_url": "a",
		"html_url": "h",
		"statuses_url": "s",
		"git_url": "g",
		"ssh_url": "s",
		"clone_url": "c",
		"svn_url": "s",
		"topics": ["octocat","api"]
    }`

	testJSONMarshal(t, u, want)
}

func TestPushEventRepoOwner_Marshal(t *testing.T) {
	testJSONMarshal(t, &PushEventRepoOwner{}, "{}")

	u := &PushEventRepoOwner{
		Name:  String("n"),
		Email: String("e"),
	}

	want := `{
		"name": "n",
		"email": "e"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectEvent{}, "{}")

	u := &ProjectEvent{
		Project: &Project{ID: Int64(1)},
		Action:  String("a"),
		Changes: &ProjectChange{
			Name: &ProjectName{From: String("NameFrom")},
			Body: &ProjectBody{From: String("BodyFrom")},
		},
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
		"action": "a",
		"changes": {
			"name": {
				"from": "NameFrom"
			},
			"body": {
				"from": "BodyFrom"
			}
		},
		"project": {
			"id": 1
		},
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

	testJSONMarshal(t, u, want)
}

func TestProjectCardEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCardEvent{}, "{}")

	u := &ProjectCardEvent{
		Action: String("a"),
		Changes: &ProjectCardChange{
			Note: &ProjectCardNote{From: String("NoteFrom")},
		},
		AfterID:     Int64(1),
		ProjectCard: &ProjectCard{ID: Int64(1)},
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
		"action": "a",
		"changes": {
			"note": {
				"from": "NoteFrom"
			}
		},
		"after_id": 1,
		"project_card": {
			"id": 1
		},
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

	testJSONMarshal(t, u, want)
}

func TestProjectColumnEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnEvent{}, "{}")

	u := &ProjectColumnEvent{
		Action: String("a"),
		Changes: &ProjectColumnChange{
			Name: &ProjectColumnName{From: String("NameFrom")},
		},
		AfterID:       Int64(1),
		ProjectColumn: &ProjectColumn{ID: Int64(1)},
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
		"action": "a",
		"changes": {
			"name": {
				"from": "NameFrom"
			}
		},
		"after_id": 1,
		"project_column": {
			"id": 1
		},
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

	testJSONMarshal(t, u, want)
}

func TestProjectV2Event_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectV2Event{}, "{}")

	u := &ProjectV2Event{
		Action: String("a"),
		ProjectsV2: &ProjectsV2{
			ID:     Int64(1),
			NodeID: String("nid"),
			Owner: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Title:            String("t"),
			Description:      String("d"),
			Public:           Bool(true),
			ClosedAt:         &Timestamp{referenceTime},
			CreatedAt:        &Timestamp{referenceTime},
			UpdatedAt:        &Timestamp{referenceTime},
			DeletedAt:        &Timestamp{referenceTime},
			Number:           Int(1),
			ShortDescription: String("sd"),
			DeletedBy: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
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
		},
	}

	want := `{
		"action": "a",
		"projects_v2": {
			"id": 1,
			"node_id": "nid",
			"owner": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"creator": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"title": "t",
			"description": "d",
			"public": true,
			"closed_at": ` + referenceTimeStr + `,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"deleted_at": ` + referenceTimeStr + `,
			"number": 1,
			"short_description": "sd",
			"deleted_by": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
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
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectV2ItemEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectV2ItemEvent{}, "{}")

	u := &ProjectV2ItemEvent{
		Action: String("a"),
		Changes: &ProjectV2ItemChange{
			ArchivedAt: &ArchivedAt{
				From: &Timestamp{referenceTime},
				To:   &Timestamp{referenceTime},
			},
		},
		ProjectV2Item: &ProjectV2Item{
			ID:            Int64(1),
			NodeID:        String("nid"),
			ProjectNodeID: String("pnid"),
			ContentNodeID: String("cnid"),
			ContentType:   String("ct"),
			Creator: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			CreatedAt:  &Timestamp{referenceTime},
			UpdatedAt:  &Timestamp{referenceTime},
			ArchivedAt: &Timestamp{referenceTime},
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
		},
	}

	want := `{
		"action":  "a",
		"changes": {
			"archived_at": {
				"from": ` + referenceTimeStr + `,
				"to": ` + referenceTimeStr + `
			}
		},
		"projects_v2_item": {
			"id": 1,
			"node_id": "nid",
			"project_node_id": "pnid",
			"content_node_id": "cnid",
			"content_type": "ct",
			"creator":  {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"archived_at": ` + referenceTimeStr + `
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
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestEvent{}, "{}")

	u := &PullRequestEvent{
		Action: String("a"),
		Assignee: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Number:      Int(1),
		PullRequest: &PullRequest{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
		},
		RequestedReviewer: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		RequestedTeam: &Team{ID: Int64(1)},
		Label:         &Label{ID: Int64(1)},
		Before:        String("before"),
		After:         String("after"),
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		PerformedViaGithubApp: &App{
			ID:          Int64(1),
			NodeID:      String("n"),
			Slug:        String("s"),
			Name:        String("n"),
			Description: String("d"),
			ExternalURL: String("e"),
			HTMLURL:     String("h"),
		},
		Organization: &Organization{
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
		"action": "a",
		"assignee": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"number": 1,
		"pull_request": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
		},
		"requested_reviewer": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"requested_team": {
			"id": 1
		},
		"label": {
			"id": 1
		},
		"before": "before",
		"after": "after",
		"repository": {
			"id": 1,
			"name": "n",
			"url": "s"
		},
		"performed_via_github_app": {
			"id": 1,
			"node_id": "n",
			"slug": "s",
			"name": "n",
			"description": "d",
			"external_url": "e",
			"html_url": "h"
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

	testJSONMarshal(t, u, want)
}

func TestPullRequestReviewCommentEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReviewCommentEvent{}, "{}")

	u := &PullRequestReviewCommentEvent{
		Action:      String("a"),
		PullRequest: &PullRequest{ID: Int64(1)},
		Comment:     &PullRequestComment{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
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
		"pull_request": {
			"id": 1
		},
		"comment": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
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

func TestPullRequestReviewThreadEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestReviewThreadEvent{}, "{}")

	u := &PullRequestReviewThreadEvent{
		Action:      String("a"),
		PullRequest: &PullRequest{ID: Int64(1)},
		Thread: &PullRequestThread{
			Comments: []*PullRequestComment{{ID: Int64(1)}, {ID: Int64(2)}},
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
		"pull_request": {
			"id": 1
		},
		"thread": {
			"comments": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			]
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

func TestPullRequestTargetEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &PullRequestTargetEvent{}, "{}")

	u := &PullRequestTargetEvent{
		Action: String("a"),
		Assignee: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Number:      Int(1),
		PullRequest: &PullRequest{ID: Int64(1)},
		Changes: &EditChange{
			Title: &EditTitle{
				From: String("TitleFrom"),
			},
			Body: &EditBody{
				From: String("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: String("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: String("BaseSHAFrom"),
				},
			},
		},
		RequestedReviewer: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		RequestedTeam: &Team{ID: Int64(1)},
		Label:         &Label{ID: Int64(1)},
		Before:        String("before"),
		After:         String("after"),
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		PerformedViaGithubApp: &App{
			ID:          Int64(1),
			NodeID:      String("n"),
			Slug:        String("s"),
			Name:        String("n"),
			Description: String("d"),
			ExternalURL: String("e"),
			HTMLURL:     String("h"),
		},
		Organization: &Organization{
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
		"action": "a",
		"assignee": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"number": 1,
		"pull_request": {
			"id": 1
		},
		"changes": {
			"title": {
				"from": "TitleFrom"
			},
			"body": {
				"from": "BodyFrom"
			},
			"base": {
				"ref": {
					"from": "BaseRefFrom"
				},
				"sha": {
					"from": "BaseSHAFrom"
				}
			}
		},
		"requested_reviewer": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"requested_team": {
			"id": 1
		},
		"label": {
			"id": 1
		},
		"before": "before",
		"after": "after",
		"repository": {
			"id": 1,
			"name": "n",
			"url": "s"
		},
		"performed_via_github_app": {
			"id": 1,
			"node_id": "n",
			"slug": "s",
			"name": "n",
			"description": "d",
			"external_url": "e",
			"html_url": "h"
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

	testJSONMarshal(t, u, want)
}

func TestRepositoryVulnerabilityAlertEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryVulnerabilityAlertEvent{}, "{}")

	u := &RepositoryVulnerabilityAlertEvent{
		Action: String("a"),
		Alert: &RepositoryVulnerabilityAlert{
			ID:                  Int64(1),
			AffectedRange:       String("ar"),
			AffectedPackageName: String("apn"),
			ExternalReference:   String("er"),
			ExternalIdentifier:  String("ei"),
			FixedIn:             String("fi"),
			Dismisser: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			DismissReason: String("dr"),
			DismissedAt:   &Timestamp{referenceTime},
		},
		Repository: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
	}

	want := `{
		"action": "a",
		"alert": {
			"id": 1,
			"affected_range": "ar",
			"affected_package_name": "apn",
			"external_reference": "er",
			"external_identifier": "ei",
			"fixed_in": "fi",
			"dismisser": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"dismiss_reason": "dr",
			"dismissed_at": ` + referenceTimeStr + `
		},
		"repository": {
			"id": 1,
			"name": "n",
			"url": "s"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestSecretScanningAlertEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &SecretScanningAlertEvent{}, "{}")

	u := &SecretScanningAlertEvent{
		Action: String("a"),
		Alert: &SecretScanningAlert{
			Number:     Int(1),
			SecretType: String("t"),
			Resolution: String("r"),
			ResolvedBy: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			ResolvedAt: &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Organization: &Organization{
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
		"alert": {
			"number": 1,
			"secret_type": "t",
			"resolution": "r",
			"resolved_by": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"resolved_at": ` + referenceTimeStr + `
		},
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

func TestSecurityAdvisoryEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &SecurityAdvisoryEvent{}, "{}")
	u := &SecurityAdvisoryEvent{
		Action: String("published"),
		SecurityAdvisory: &SecurityAdvisory{
			CVSS: &AdvisoryCVSS{
				Score:        Float64(1.0),
				VectorString: String("vs"),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: String("cweid"),
					Name:  String("n"),
				},
			},
			GHSAID:      String("GHSA-rf4j-j272-some"),
			Summary:     String("Siuuuuuuuuu"),
			Description: String("desc"),
			Severity:    String("moderate"),
			Identifiers: []*AdvisoryIdentifier{
				{
					Value: String("GHSA-rf4j-j272-some"),
					Type:  String("GHSA"),
				},
			},
			References: []*AdvisoryReference{
				{
					URL: String("https://some-url"),
				},
			},
			PublishedAt: &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			WithdrawnAt: nil,
			Vulnerabilities: []*AdvisoryVulnerability{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: String("ucl"),
						Name:      String("penaldo"),
					},
					Severity:               String("moderate"),
					VulnerableVersionRange: String(">= 2.0.0, < 2.0.2"),
					FirstPatchedVersion: &FirstPatchedVersion{
						Identifier: String("2.0.2"),
					},
				},
			},
		},
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
		Organization: &Organization{
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
		Repository: &Repository{
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
		"action": "published",
		"security_advisory": {
		  "ghsa_id": "GHSA-rf4j-j272-some",
		  "summary": "Siuuuuuuuuu",
		  "cvss": {
			"score": 1.0,
			"vector_string": "vs"
		  },
		  "cwes": [
			{
				"cwe_id": "cweid",
				"name": "n"
			}
		  ],
		  "description": "desc",
		  "severity": "moderate",
		  "identifiers": [
			{
			  "value": "GHSA-rf4j-j272-some",
			  "type": "GHSA"
			}
		  ],
		  "references": [
			{
			  "url": "https://some-url"
			}
		  ],
		  "published_at": ` + referenceTimeStr + `,
		  "updated_at": ` + referenceTimeStr + `,
		  "withdrawn_at": null,
		  "vulnerabilities": [
			{
			  "package": {
				"ecosystem": "ucl",
				"name": "penaldo"
			  },
			  "severity": "moderate",
			  "vulnerable_version_range": ">= 2.0.0, < 2.0.2",
			  "first_patched_version": {
				"identifier": "2.0.2"
			  }
			}
		  ]
		},
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
			"url": "s",
			"name": "n"
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

func TestSecurityAndAnalysisEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &SecurityAndAnalysisEvent{}, "{}")

	u := &SecurityAndAnalysisEvent{
		Changes: &SecurityAndAnalysisChange{
			From: &SecurityAndAnalysisChangeFrom{
				SecurityAndAnalysis: &SecurityAndAnalysis{
					AdvancedSecurity: &AdvancedSecurity{
						Status: String("enabled"),
					},
					SecretScanning: &SecretScanning{
						Status: String("enabled"),
					},
					SecretScanningPushProtection: &SecretScanningPushProtection{
						Status: String("enabled"),
					},
					DependabotSecurityUpdates: &DependabotSecurityUpdates{
						Status: String("enabled"),
					},
				},
			},
		},
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
		Organization: &Organization{
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
		Repository: &Repository{
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
		"changes": {
			"from": {
				"security_and_analysis": {
					"advanced_security": {
						"status": "enabled"
					},
					"secret_scanning": {
						"status": "enabled"
					},
					"secret_scanning_push_protection": {
						"status": "enabled"
					},
					"dependabot_security_updates": {
						"status": "enabled"
					}
				}
			}
		},
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
			"url": "s",
			"name": "n"
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
		"target_type": "running"
	}`

	testJSONMarshal(t, u, want)
}

func TestCodeScanningAlertEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &CodeScanningAlertEvent{}, "{}")

	u := &CodeScanningAlertEvent{
		Action: String("reopened"),
		Alert: &Alert{
			Number: Int(10),
			Rule: &Rule{
				ID:              String("Style/FrozenStringLiteralComment"),
				Severity:        String("note"),
				Description:     String("desc"),
				FullDescription: String("full desc"),
				Tags:            []string{"style"},
				Help:            String("help"),
			},
			Tool: &Tool{
				Name:    String("Rubocop"),
				Version: nil,
			},
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			FixedAt:   nil,
			State:     String("open"),
			URL:       String("a"),
			HTMLURL:   String("a"),
			Instances: []*MostRecentInstance{
				{
					Ref:         String("refs/heads/main"),
					AnalysisKey: String(".github/workflows/workflow.yml:upload"),
					Environment: String("{}"),
					State:       String("open"),
				},
			},
			DismissedBy:     nil,
			DismissedAt:     nil,
			DismissedReason: nil,
		},
		Ref:       String("refs/heads/main"),
		CommitOID: String("d6e4c75c141dbacecc279b721b8bsomeSHA"),
		Repo: &Repository{
			ID:     Int64(1234234535),
			NodeID: String("MDEwOlJlcG9zaXRvcnkxODY4NT=="),
			Owner: &User{
				Login:             String("Codertocat"),
				ID:                Int64(21031067),
				NodeID:            String("MDQ6VXNlcjIxMDMxMDY3"),
				AvatarURL:         String("a"),
				GravatarID:        String(""),
				URL:               String("a"),
				HTMLURL:           String("a"),
				Type:              String("User"),
				SiteAdmin:         Bool(false),
				FollowersURL:      String("a"),
				FollowingURL:      String("a"),
				EventsURL:         String("a"),
				GistsURL:          String("a"),
				OrganizationsURL:  String("a"),
				ReceivedEventsURL: String("a"),
				ReposURL:          String("a"),
				StarredURL:        String("a"),
				SubscriptionsURL:  String("a"),
			},
			HTMLURL:          String("a"),
			Name:             String("Hello-World"),
			FullName:         String("Codertocat/Hello-World"),
			Description:      nil,
			Fork:             Bool(false),
			Homepage:         nil,
			DefaultBranch:    String("main"),
			CreatedAt:        &Timestamp{referenceTime},
			PushedAt:         &Timestamp{referenceTime},
			UpdatedAt:        &Timestamp{referenceTime},
			CloneURL:         String("a"),
			GitURL:           String("a"),
			MirrorURL:        nil,
			SSHURL:           String("a"),
			SVNURL:           String("a"),
			Language:         nil,
			ForksCount:       Int(0),
			OpenIssuesCount:  Int(2),
			OpenIssues:       Int(2),
			StargazersCount:  Int(0),
			WatchersCount:    Int(0),
			Watchers:         Int(0),
			Size:             Int(0),
			Archived:         Bool(false),
			Disabled:         Bool(false),
			License:          nil,
			Private:          Bool(false),
			HasIssues:        Bool(true),
			HasWiki:          Bool(true),
			HasPages:         Bool(true),
			HasProjects:      Bool(true),
			HasDownloads:     Bool(true),
			URL:              String("a"),
			ArchiveURL:       String("a"),
			AssigneesURL:     String("a"),
			BlobsURL:         String("a"),
			BranchesURL:      String("a"),
			CollaboratorsURL: String("a"),
			CommentsURL:      String("a"),
			CommitsURL:       String("a"),
			CompareURL:       String("a"),
			ContentsURL:      String("a"),
			ContributorsURL:  String("a"),
			DeploymentsURL:   String("a"),
			DownloadsURL:     String("a"),
			EventsURL:        String("a"),
			ForksURL:         String("a"),
			GitCommitsURL:    String("a"),
			GitRefsURL:       String("a"),
			GitTagsURL:       String("a"),
			HooksURL:         String("a"),
			IssueCommentURL:  String("a"),
			IssueEventsURL:   String("a"),
			IssuesURL:        String("a"),
			KeysURL:          String("a"),
			LabelsURL:        String("a"),
			LanguagesURL:     String("a"),
			MergesURL:        String("a"),
			MilestonesURL:    String("a"),
			NotificationsURL: String("a"),
			PullsURL:         String("a"),
			ReleasesURL:      String("a"),
			StargazersURL:    String("a"),
			StatusesURL:      String("a"),
			SubscribersURL:   String("a"),
			SubscriptionURL:  String("a"),
			TagsURL:          String("a"),
			TreesURL:         String("a"),
			TeamsURL:         String("a"),
		},
		Org: &Organization{
			Login:            String("Octocoders"),
			ID:               Int64(6),
			NodeID:           String("MDEyOk9yZ2FuaXphdGlvbjY="),
			AvatarURL:        String("a"),
			Description:      String(""),
			URL:              String("a"),
			EventsURL:        String("a"),
			HooksURL:         String("a"),
			IssuesURL:        String("a"),
			MembersURL:       String("a"),
			PublicMembersURL: String("a"),
			ReposURL:         String("a"),
		},
		Sender: &User{
			Login:             String("github"),
			ID:                Int64(9919),
			NodeID:            String("MDEyOk9yZ2FuaXphdGlvbjk5MTk="),
			AvatarURL:         String("a"),
			HTMLURL:           String("a"),
			GravatarID:        String(""),
			Type:              String("Organization"),
			SiteAdmin:         Bool(false),
			URL:               String("a"),
			EventsURL:         String("a"),
			FollowingURL:      String("a"),
			FollowersURL:      String("a"),
			GistsURL:          String("a"),
			OrganizationsURL:  String("a"),
			ReceivedEventsURL: String("a"),
			ReposURL:          String("a"),
			StarredURL:        String("a"),
			SubscriptionsURL:  String("a"),
		},
	}

	want := `{
		"action": "reopened",
		"alert": {
		  "number": 10,
		  "created_at": ` + referenceTimeStr + `,
		  "updated_at": ` + referenceTimeStr + `,
		  "url": "a",
		  "html_url": "a",
		  "instances": [
			{
			  "ref": "refs/heads/main",
			  "analysis_key": ".github/workflows/workflow.yml:upload",
			  "environment": "{}",
			  "state": "open"
			}
		  ],
		  "state": "open",
		  "fixed_at": null,
		  "dismissed_by": null,
		  "dismissed_at": null,
		  "dismissed_reason": null,
		  "rule": {
			"id": "Style/FrozenStringLiteralComment",
			"severity": "note",
			"description": "desc",
			"full_description": "full desc",
			"tags": [
			  "style"
			],
			"help": "help"
		  },
		  "tool": {
			"name": "Rubocop",
			"version": null
		  }
		},
		"ref": "refs/heads/main",
		"commit_oid": "d6e4c75c141dbacecc279b721b8bsomeSHA",
		"repository": {
		  "id": 1234234535,
		  "node_id": "MDEwOlJlcG9zaXRvcnkxODY4NT==",
		  "name": "Hello-World",
		  "full_name": "Codertocat/Hello-World",
		  "private": false,
		  "owner": {
			"login": "Codertocat",
			"id": 21031067,
			"node_id": "MDQ6VXNlcjIxMDMxMDY3",
			"avatar_url": "a",
			"gravatar_id": "",
			"url": "a",
			"html_url": "a",
			"followers_url": "a",
			"following_url": "a",
			"gists_url": "a",
			"starred_url": "a",
			"subscriptions_url": "a",
			"organizations_url": "a",
			"repos_url": "a",
			"events_url": "a",
			"received_events_url": "a",
			"type": "User",
			"site_admin": false
		  },
		  "html_url": "a",
		  "description": null,
		  "fork": false,
		  "url": "a",
		  "forks_url": "a",
		  "keys_url": "a",
		  "collaborators_url": "a",
		  "teams_url": "a",
		  "hooks_url": "a",
		  "issue_events_url": "a",
		  "events_url": "a",
		  "assignees_url": "a",
		  "branches_url": "a",
		  "tags_url": "a",
		  "blobs_url": "a",
		  "git_tags_url": "a",
		  "git_refs_url": "a",
		  "trees_url": "a",
		  "statuses_url": "a",
		  "languages_url": "a",
		  "stargazers_url": "a",
		  "contributors_url": "a",
		  "subscribers_url": "a",
		  "subscription_url": "a",
		  "commits_url": "a",
		  "git_commits_url": "a",
		  "comments_url": "a",
		  "issue_comment_url": "a",
		  "contents_url": "a",
		  "compare_url": "a",
		  "merges_url": "a",
		  "archive_url": "a",
		  "downloads_url": "a",
		  "issues_url": "a",
		  "pulls_url": "a",
		  "milestones_url": "a",
		  "notifications_url": "a",
		  "labels_url": "a",
		  "releases_url": "a",
		  "deployments_url": "a",
		  "created_at": ` + referenceTimeStr + `,
		  "updated_at": ` + referenceTimeStr + `,
		  "pushed_at": ` + referenceTimeStr + `,
		  "git_url": "a",
		  "ssh_url": "a",
		  "clone_url": "a",
		  "svn_url": "a",
		  "homepage": null,
		  "size": 0,
		  "stargazers_count": 0,
		  "watchers_count": 0,
		  "language": null,
		  "has_issues": true,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": true,
		  "forks_count": 0,
		  "mirror_url": null,
		  "archived": false,
		  "disabled": false,
		  "open_issues_count": 2,
		  "license": null,
		  "forks": 0,
		  "open_issues": 2,
		  "watchers": 0,
		  "default_branch": "main"
		},
		"organization": {
		  "login": "Octocoders",
		  "id": 6,
		  "node_id": "MDEyOk9yZ2FuaXphdGlvbjY=",
		  "url": "a",
		  "repos_url": "a",
		  "events_url": "a",
		  "hooks_url": "a",
		  "issues_url": "a",
		  "members_url": "a",
		  "public_members_url": "a",
		  "avatar_url": "a",
		  "description": ""
		},
		"sender": {
		  "login": "github",
		  "id": 9919,
		  "node_id": "MDEyOk9yZ2FuaXphdGlvbjk5MTk=",
		  "avatar_url": "a",
		  "gravatar_id": "",
		  "url": "a",
		  "html_url": "a",
		  "followers_url": "a",
		  "following_url": "a",
		  "gists_url": "a",
		  "starred_url": "a",
		  "subscriptions_url": "a",
		  "organizations_url": "a",
		  "repos_url": "a",
		  "events_url": "a",
		  "received_events_url": "a",
		  "type": "Organization",
		  "site_admin": false
		}
	  }`

	testJSONMarshal(t, u, want)
}

func TestSponsorshipEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &SponsorshipEvent{}, "{}")

	u := &SponsorshipEvent{
		Action:        String("created"),
		EffectiveDate: String("2023-01-01T00:00:00Z"),
		Changes: &SponsorshipChanges{
			Tier: &SponsorshipTier{
				From: String("basic"),
			},
			PrivacyLevel: String("public"),
		},
		Repository: &Repository{
			ID:     Int64(12345),
			NodeID: String("MDEwOlJlcG9zaXRvcnkxMjM0NQ=="),
			Name:   String("example-repo"),
		},
		Organization: &Organization{
			Login: String("example-org"),
			ID:    Int64(67890),
		},
		Sender: &User{
			Login: String("example-user"),
			ID:    Int64(1111),
		},
		Installation: &Installation{
			ID: Int64(2222),
		},
	}

	want := `{
		"action": "created",
		"effective_date": "2023-01-01T00:00:00Z",
		"changes": {
			"tier": {
				"from": "basic"
			},
			"privacy_level": "public"
		},
		"repository": {
			"id": 12345,
			"node_id": "MDEwOlJlcG9zaXRvcnkxMjM0NQ==",
			"name": "example-repo"
		},
		"organization": {
			"login": "example-org",
			"id": 67890
		},
		"sender": {
			"login": "example-user",
			"id": 1111
		},
		"installation": {
			"id": 2222
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestSponsorshipChanges_Marshal(t *testing.T) {
	testJSONMarshal(t, &SponsorshipChanges{}, "{}")

	u := &SponsorshipChanges{
		Tier: &SponsorshipTier{
			From: String("premium"),
		},
		PrivacyLevel: String("private"),
	}

	want := `{
		"tier": {
			"from": "premium"
		},
		"privacy_level": "private"
	}`

	testJSONMarshal(t, u, want)
}

func TestSponsorshipTier_Marshal(t *testing.T) {
	testJSONMarshal(t, &SponsorshipTier{}, "{}")

	u := &SponsorshipTier{
		From: String("gold"),
	}

	want := `{
		"from": "gold"
	}`

	testJSONMarshal(t, u, want)
}
