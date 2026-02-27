// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEditChange_Marshal_TitleChange(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Title: &EditTitle{
			From: Ptr("TitleFrom"),
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
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Title: nil,
		Body: &EditBody{
			From: Ptr("BodyFrom"),
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
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	base := EditBase{
		Ref: &EditRef{
			From: Ptr("BaseRefFrom"),
		},
		SHA: &EditSHA{
			From: Ptr("BaseSHAFrom"),
		},
	}

	u := &EditChange{
		Title: nil,
		Body:  nil,
		Base:  &base,
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
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Repo: &EditRepo{
			Name: &RepoName{
				From: Ptr("old-repo-name"),
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
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Owner: &EditOwner{
			OwnerInfo: &OwnerInfo{
				User: &User{
					Login:     Ptr("l"),
					ID:        Ptr(int64(1)),
					NodeID:    Ptr("n"),
					URL:       Ptr("u"),
					ReposURL:  Ptr("r"),
					EventsURL: Ptr("e"),
					AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &EditChange{}, "{}")

	u := &EditChange{
		Owner: &EditOwner{
			OwnerInfo: &OwnerInfo{
				Org: &User{
					Login:     Ptr("l"),
					ID:        Ptr(int64(1)),
					NodeID:    Ptr("n"),
					URL:       Ptr("u"),
					ReposURL:  Ptr("r"),
					EventsURL: Ptr("e"),
					AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: &ProjectName{From: Ptr("NameFrom")},
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
	t.Parallel()
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: nil,
		Body: &ProjectBody{From: Ptr("BodyFrom")},
	}

	want := `{
		"body": {
			"from": "BodyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardChange_Marshal_NoteChange(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectCardChange{}, "{}")

	u := &ProjectCardChange{
		Note: &ProjectCardNote{From: Ptr("NoteFrom")},
	}

	want := `{
		"note": {
			"from": "NoteFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnChange_Marshal_NameChange(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectColumnChange{}, "{}")

	u := &ProjectColumnChange{
		Name: &ProjectColumnName{From: Ptr("NameFrom")},
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestBranchProtectionConfigurationEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &BranchProtectionConfigurationEvent{}, "{}")
	u := &BranchProtectionConfigurationEvent{
		Action: Ptr("enabled"),
		Repo: &Repository{
			ID:     Ptr(int64(12345)),
			NodeID: Ptr("MDEwOlJlcG9zaXRvcnkxMjM0NQ=="),
			Name:   Ptr("example-repo"),
		},
		Org: &Organization{
			Login: Ptr("example-org"),
			ID:    Ptr(int64(67890)),
		},
		Sender: &User{
			Login: Ptr("example-user"),
			ID:    Ptr(int64(1111)),
		},
		Installation: &Installation{
			ID: Ptr(int64(2222)),
		},
		Enterprise: &Enterprise{
			ID:   Ptr(3333),
			Slug: Ptr("example-enterprise"),
			Name: Ptr("Example Enterprise"),
		},
	}

	want := `{
		"action": "enabled",
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
		},
		"enterprise": {
			"id": 3333,
			"slug": "example-enterprise",
			"name": "Example Enterprise"
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestTeamAddEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TeamAddEvent{}, "{}")

	u := &TeamAddEvent{
		Team: &Team{
			ID:              Ptr(int64(1)),
			NodeID:          Ptr("n"),
			Name:            Ptr("n"),
			Description:     Ptr("d"),
			URL:             Ptr("u"),
			Slug:            Ptr("s"),
			Permission:      Ptr("p"),
			Privacy:         Ptr("p"),
			MembersCount:    Ptr(1),
			ReposCount:      Ptr(1),
			MembersURL:      Ptr("m"),
			RepositoriesURL: Ptr("r"),
			Organization: &Organization{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				AvatarURL: Ptr("a"),
				HTMLURL:   Ptr("h"),
				Name:      Ptr("n"),
				Company:   Ptr("c"),
				Blog:      Ptr("b"),
				Location:  Ptr("l"),
				Email:     Ptr("e"),
			},
			Parent: &Team{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("n"),
				Name:         Ptr("n"),
				Description:  Ptr("d"),
				URL:          Ptr("u"),
				Slug:         Ptr("s"),
				Permission:   Ptr("p"),
				Privacy:      Ptr("p"),
				MembersCount: Ptr(1),
				ReposCount:   Ptr(1),
			},
			LDAPDN: Ptr("l"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &StarEvent{}, "{}")

	u := &StarEvent{
		Action:    Ptr("a"),
		StarredAt: &Timestamp{referenceTime},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &TeamEvent{}, "{}")

	u := &TeamEvent{
		Action: Ptr("a"),
		Team: &Team{
			ID:              Ptr(int64(1)),
			NodeID:          Ptr("n"),
			Name:            Ptr("n"),
			Description:     Ptr("d"),
			URL:             Ptr("u"),
			Slug:            Ptr("s"),
			Permission:      Ptr("p"),
			Privacy:         Ptr("p"),
			MembersCount:    Ptr(1),
			ReposCount:      Ptr(1),
			MembersURL:      Ptr("m"),
			RepositoriesURL: Ptr("r"),
			Organization: &Organization{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				AvatarURL: Ptr("a"),
				HTMLURL:   Ptr("h"),
				Name:      Ptr("n"),
				Company:   Ptr("c"),
				Blog:      Ptr("b"),
				Location:  Ptr("l"),
				Email:     Ptr("e"),
			},
			Parent: &Team{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("n"),
				Name:         Ptr("n"),
				Description:  Ptr("d"),
				URL:          Ptr("u"),
				Slug:         Ptr("s"),
				Permission:   Ptr("p"),
				Privacy:      Ptr("p"),
				MembersCount: Ptr(1),
				ReposCount:   Ptr(1),
			},
			LDAPDN: Ptr("l"),
		},
		Changes: &TeamChange{
			Description: &TeamDescription{
				From: Ptr("from"),
			},
			Name: &TeamName{
				From: Ptr("from"),
			},
			Privacy: &TeamPrivacy{
				From: Ptr("from"),
			},
			Repository: &TeamRepository{
				Permissions: &TeamPermissions{
					From: &TeamPermissionsFrom{
						Admin: Ptr(true),
						Pull:  Ptr(true),
						Push:  Ptr(true),
					},
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &InstallationRepositoriesEvent{}, "{}")

	u := &InstallationRepositoriesEvent{
		Action: Ptr("a"),
		RepositoriesAdded: []*Repository{
			{
				ID:   Ptr(int64(1)),
				URL:  Ptr("s"),
				Name: Ptr("n"),
			},
		},
		RepositoriesRemoved: []*Repository{
			{
				ID:   Ptr(int64(1)),
				URL:  Ptr("s"),
				Name: Ptr("n"),
			},
		},
		RepositorySelection: Ptr("rs"),
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &InstallationTargetEvent{}, "{}")

	u := &InstallationTargetEvent{
		Account: &User{
			Login:     Ptr("u"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("l"),
		},
		Action: Ptr("a"),
		Changes: &InstallationChanges{
			Login: &InstallationLoginChange{
				From: Ptr("p"),
			},
			Slug: &InstallationSlugChange{
				From: Ptr("j"),
			},
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		TargetType: Ptr("running"),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &EditTitle{}, "{}")

	u := &EditTitle{
		From: Ptr("EditTitleFrom"),
	}

	want := `{
		"from": "EditTitleFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBody_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EditBody{}, "{}")

	u := &EditBody{
		From: Ptr("EditBodyFrom"),
	}

	want := `{
		"from": "EditBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBase_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EditBase{}, "{}")

	u := &EditBase{
		Ref: &EditRef{
			From: Ptr("EditRefFrom"),
		},
		SHA: &EditSHA{
			From: Ptr("EditSHAFrom"),
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
	t.Parallel()
	testJSONMarshal(t, &EditRef{}, "{}")

	u := &EditRef{
		From: Ptr("EditRefFrom"),
	}

	want := `{
		"from": "EditRefFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditSHA_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EditSHA{}, "{}")

	u := &EditSHA{
		From: Ptr("EditSHAFrom"),
	}

	want := `{
		"from": "EditSHAFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectName_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectName{}, "{}")

	u := &ProjectName{
		From: Ptr("ProjectNameFrom"),
	}

	want := `{
		"from": "ProjectNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectBody_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectBody{}, "{}")

	u := &ProjectBody{
		From: Ptr("ProjectBodyFrom"),
	}

	want := `{
		"from": "ProjectBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardNote_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectCardNote{}, "{}")

	u := &ProjectCardNote{
		From: Ptr("ProjectCardNoteFrom"),
	}

	want := `{
		"from": "ProjectCardNoteFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnName_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectColumnName{}, "{}")

	u := &ProjectColumnName{
		From: Ptr("ProjectColumnNameFrom"),
	}

	want := `{
		"from": "ProjectColumnNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamDescription_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TeamDescription{}, "{}")

	u := &TeamDescription{
		From: Ptr("TeamDescriptionFrom"),
	}

	want := `{
		"from": "TeamDescriptionFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamName_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TeamName{}, "{}")

	u := &TeamName{
		From: Ptr("TeamNameFrom"),
	}

	want := `{
		"from": "TeamNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPrivacy_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TeamPrivacy{}, "{}")

	u := &TeamPrivacy{
		From: Ptr("TeamPrivacyFrom"),
	}

	want := `{
		"from": "TeamPrivacyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamRepository_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TeamRepository{}, "{}")

	u := &TeamRepository{
		Permissions: &TeamPermissions{
			From: &TeamPermissionsFrom{
				Admin: Ptr(true),
				Pull:  Ptr(true),
				Push:  Ptr(true),
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
	t.Parallel()
	testJSONMarshal(t, &TeamPermissions{}, "{}")

	u := &TeamPermissions{
		From: &TeamPermissionsFrom{
			Admin: Ptr(true),
			Pull:  Ptr(true),
			Push:  Ptr(true),
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
	t.Parallel()
	testJSONMarshal(t, &TeamPermissionsFrom{}, "{}")

	u := &TeamPermissionsFrom{
		Admin: Ptr(true),
		Pull:  Ptr(true),
		Push:  Ptr(true),
	}

	want := `{
		"admin": true,
		"pull": true,
		"push": true
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoryVulnerabilityAlert_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryVulnerabilityAlert{}, "{}")

	u := &RepositoryVulnerabilityAlert{
		ID:                  Ptr(int64(1)),
		AffectedRange:       Ptr("ar"),
		AffectedPackageName: Ptr("apn"),
		ExternalReference:   Ptr("er"),
		ExternalIdentifier:  Ptr("ei"),
		FixedIn:             Ptr("fi"),
		Dismisser: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		DismissReason: Ptr("dr"),
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
	t.Parallel()
	testJSONMarshal(t, &Page{}, "{}")

	u := &Page{
		PageName: Ptr("p"),
		Title:    Ptr("t"),
		Summary:  Ptr("s"),
		Action:   Ptr("a"),
		SHA:      Ptr("s"),
		HTMLURL:  Ptr("h"),
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
	t.Parallel()
	testJSONMarshal(t, &TeamChange{}, "{}")

	u := &TeamChange{
		Description: &TeamDescription{
			From: Ptr("DescriptionFrom"),
		},
		Name: &TeamName{
			From: Ptr("NameFrom"),
		},
		Privacy: &TeamPrivacy{
			From: Ptr("PrivacyFrom"),
		},
		Repository: &TeamRepository{
			Permissions: &TeamPermissions{
				From: &TeamPermissionsFrom{
					Admin: Ptr(false),
					Pull:  Ptr(false),
					Push:  Ptr(false),
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
	t.Parallel()
	testJSONMarshal(t, &IssueCommentEvent{}, "{}")

	u := &IssueCommentEvent{
		Action:  Ptr("a"),
		Issue:   &Issue{ID: Ptr(int64(1))},
		Comment: &IssueComment{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &IssuesEvent{}, "{}")

	u := &IssuesEvent{
		Action: Ptr("a"),
		Issue:  &Issue{ID: Ptr(int64(1))},
		Assignee: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Label: &Label{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &LabelEvent{}, "{}")

	u := &LabelEvent{
		Action: Ptr("a"),
		Label:  &Label{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &MilestoneEvent{}, "{}")

	u := &MilestoneEvent{
		Action:    Ptr("a"),
		Milestone: &Milestone{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PublicEvent{}, "{}")

	u := &PublicEvent{
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PullRequestReviewEvent{}, "{}")

	u := &PullRequestReviewEvent{
		Action:      Ptr("a"),
		Review:      &PullRequestReview{ID: Ptr(int64(1))},
		PullRequest: &PullRequest{ID: Ptr(int64(1))},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PushEvent{}, "{}")

	u := &PushEvent{
		PushID: Ptr(int64(1)),
		Head:   Ptr("h"),
		Ref:    Ptr("ref"),
		Size:   Ptr(1),
		Commits: []*HeadCommit{
			{ID: Ptr("id")},
		},
		Before:       Ptr("b"),
		DistinctSize: Ptr(1),
		After:        Ptr("a"),
		Created:      Ptr(true),
		Deleted:      Ptr(true),
		Forced:       Ptr(true),
		BaseRef:      Ptr("a"),
		Compare:      Ptr("a"),
		Repo:         &PushEventRepository{ID: Ptr(int64(1))},
		HeadCommit:   &HeadCommit{ID: Ptr("id")},
		Pusher: &CommitAuthor{
			Login: Ptr("l"),
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("n"),
			Email: Ptr("e"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &StatusEvent{}, "{}")

	u := &StatusEvent{
		SHA:         Ptr("sha"),
		State:       Ptr("s"),
		Description: Ptr("d"),
		TargetURL:   Ptr("turl"),
		Branches: []*Branch{
			{
				Name:      Ptr("n"),
				Commit:    &RepositoryCommit{NodeID: Ptr("nid")},
				Protected: Ptr(false),
			},
		},
		ID:        Ptr(int64(1)),
		Name:      Ptr("n"),
		Context:   Ptr("c"),
		Commit:    &RepositoryCommit{NodeID: Ptr("nid")},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &MarketplacePurchaseEvent{}, "{}")

	u := &MarketplacePurchaseEvent{
		Action:        Ptr("a"),
		EffectiveDate: &Timestamp{referenceTime},
		MarketplacePurchase: &MarketplacePurchase{
			BillingCycle:    Ptr("bc"),
			NextBillingDate: &Timestamp{referenceTime},
			UnitCount:       Ptr(1),
			Plan: &MarketplacePlan{
				URL:                 Ptr("u"),
				AccountsURL:         Ptr("au"),
				ID:                  Ptr(int64(1)),
				Number:              Ptr(1),
				Name:                Ptr("n"),
				Description:         Ptr("d"),
				MonthlyPriceInCents: Ptr(1),
				YearlyPriceInCents:  Ptr(1),
				PriceModel:          Ptr("pm"),
				UnitName:            Ptr("un"),
				Bullets:             &[]string{"b"},
				State:               Ptr("s"),
				HasFreeTrial:        Ptr(false),
			},
			OnFreeTrial:     Ptr(false),
			FreeTrialEndsOn: &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
		},
		PreviousMarketplacePurchase: &MarketplacePurchase{
			BillingCycle:    Ptr("bc"),
			NextBillingDate: &Timestamp{referenceTime},
			UnitCount:       Ptr(1),
			Plan: &MarketplacePlan{
				URL:                 Ptr("u"),
				AccountsURL:         Ptr("au"),
				ID:                  Ptr(int64(1)),
				Number:              Ptr(1),
				Name:                Ptr("n"),
				Description:         Ptr("d"),
				MonthlyPriceInCents: Ptr(1),
				YearlyPriceInCents:  Ptr(1),
				PriceModel:          Ptr("pm"),
				UnitName:            Ptr("un"),
				Bullets:             &[]string{"b"},
				State:               Ptr("s"),
				HasFreeTrial:        Ptr(false),
			},
			OnFreeTrial:     Ptr(false),
			FreeTrialEndsOn: &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &OrganizationEvent{}, "{}")

	u := &OrganizationEvent{
		Action:     Ptr("a"),
		Invitation: &Invitation{ID: Ptr(int64(1))},
		Membership: &Membership{
			URL:             Ptr("url"),
			State:           Ptr("s"),
			Role:            Ptr("r"),
			OrganizationURL: Ptr("ou"),
			Organization: &Organization{
				BillingEmail:                         Ptr("be"),
				Blog:                                 Ptr("b"),
				Company:                              Ptr("c"),
				Email:                                Ptr("e"),
				TwitterUsername:                      Ptr("tu"),
				Location:                             Ptr("loc"),
				Name:                                 Ptr("n"),
				Description:                          Ptr("d"),
				IsVerified:                           Ptr(true),
				HasOrganizationProjects:              Ptr(true),
				HasRepositoryProjects:                Ptr(true),
				DefaultRepoPermission:                Ptr("drp"),
				MembersCanCreateRepos:                Ptr(true),
				MembersCanCreateInternalRepos:        Ptr(true),
				MembersCanCreatePrivateRepos:         Ptr(true),
				MembersCanCreatePublicRepos:          Ptr(false),
				MembersAllowedRepositoryCreationType: Ptr("marct"),
				MembersCanCreatePages:                Ptr(true),
				MembersCanCreatePublicPages:          Ptr(false),
				MembersCanCreatePrivatePages:         Ptr(true),
			},
			User: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PageBuildEvent{}, "{}")

	u := &PageBuildEvent{
		Build: &PagesBuild{URL: Ptr("url")},
		ID:    Ptr(int64(1)),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &CommitCommentEvent{}, "{}")

	u := &CommitCommentEvent{
		Comment: &RepositoryComment{
			HTMLURL:  Ptr("hurl"),
			URL:      Ptr("url"),
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			CommitID: Ptr("cid"),
			User: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Reactions: &Reactions{
				TotalCount: Ptr(1),
				PlusOne:    Ptr(1),
				MinusOne:   Ptr(1),
				Laugh:      Ptr(1),
				Confused:   Ptr(1),
				Heart:      Ptr(1),
				Hooray:     Ptr(1),
				Rocket:     Ptr(1),
				Eyes:       Ptr(1),
				URL:        Ptr("url"),
			},
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			Body:      Ptr("b"),
			Path:      Ptr("path"),
			Position:  Ptr(1),
		},
		Action: Ptr("a"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DeploymentEvent{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentEvent{
		Deployment: &Deployment{
			URL:         Ptr("url"),
			ID:          Ptr(int64(1)),
			SHA:         Ptr("sha"),
			Ref:         Ptr("ref"),
			Task:        Ptr("t"),
			Payload:     jsonMsg,
			Environment: Ptr("e"),
			Description: Ptr("d"),
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   Ptr("surl"),
			RepositoryURL: Ptr("rurl"),
			NodeID:        Ptr("nid"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Workflow: &Workflow{
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("nid"),
			Name:      Ptr("n"),
			Path:      Ptr("p"),
			State:     Ptr("s"),
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			URL:       Ptr("u"),
			HTMLURL:   Ptr("h"),
			BadgeURL:  Ptr("b"),
		},
		WorkflowRun: &WorkflowRun{
			ID:         Ptr(int64(1)),
			Name:       Ptr("n"),
			NodeID:     Ptr("nid"),
			HeadBranch: Ptr("hb"),
			HeadSHA:    Ptr("hs"),
			RunNumber:  Ptr(1),
			RunAttempt: Ptr(1),
			Event:      Ptr("e"),
			Status:     Ptr("s"),
			Conclusion: Ptr("c"),
			WorkflowID: Ptr(int64(1)),
			URL:        Ptr("u"),
			HTMLURL:    Ptr("h"),
			PullRequests: []*PullRequest{
				{
					URL:    Ptr("u"),
					ID:     Ptr(int64(1)),
					Number: Ptr(1),
					Head: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("s"),
							Name: Ptr("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("u"),
							Name: Ptr("n"),
						},
					},
				},
			},
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            Ptr("j"),
			LogsURL:            Ptr("l"),
			CheckSuiteURL:      Ptr("c"),
			ArtifactsURL:       Ptr("a"),
			CancelURL:          Ptr("c"),
			RerunURL:           Ptr("r"),
			PreviousAttemptURL: Ptr("p"),
			HeadCommit: &HeadCommit{
				Message: Ptr("m"),
				Author: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
				URL:       Ptr("u"),
				Distinct:  Ptr(false),
				SHA:       Ptr("s"),
				ID:        Ptr("i"),
				TreeID:    Ptr("tid"),
				Timestamp: &Timestamp{referenceTime},
				Committer: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
			},
			WorkflowURL: Ptr("w"),
			Repository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
			HeadRepository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
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
			"client_id": "cid",
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestDeploymentProtectionRuleEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DeploymentProtectionRuleEvent{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentProtectionRuleEvent{
		Action:                Ptr("a"),
		Environment:           Ptr("e"),
		DeploymentCallbackURL: Ptr("b"),
		Deployment: &Deployment{
			URL:         Ptr("url"),
			ID:          Ptr(int64(1)),
			SHA:         Ptr("sha"),
			Ref:         Ptr("ref"),
			Task:        Ptr("t"),
			Payload:     jsonMsg,
			Environment: Ptr("e"),
			Description: Ptr("d"),
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   Ptr("surl"),
			RepositoryURL: Ptr("rurl"),
			NodeID:        Ptr("nid"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		PullRequests: []*PullRequest{
			{
				URL:    Ptr("u"),
				ID:     Ptr(int64(1)),
				Number: Ptr(1),
				Head: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("s"),
						Name: Ptr("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("u"),
						Name: Ptr("n"),
					},
				},
			},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestDeploymentReviewEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DeploymentReviewEvent{}, "{}")

	u := &DeploymentReviewEvent{
		Action:      Ptr("a"),
		Environment: Ptr("e"),
		Requester: &User{
			AvatarURL:         Ptr("a"),
			Email:             Ptr("e"),
			EventsURL:         Ptr("e"),
			FollowersURL:      Ptr("f"),
			FollowingURL:      Ptr("f"),
			GistsURL:          Ptr("g"),
			GravatarID:        Ptr("g"),
			HTMLURL:           Ptr("h"),
			ID:                Ptr(int64(1)),
			Login:             Ptr("l"),
			Name:              Ptr("n"),
			NodeID:            Ptr("n"),
			OrganizationsURL:  Ptr("o"),
			ReceivedEventsURL: Ptr("r"),
			ReposURL:          Ptr("r"),
			SiteAdmin:         Ptr(false),
			StarredURL:        Ptr("s"),
			SubscriptionsURL:  Ptr("s"),
			Type:              Ptr("User"),
			URL:               Ptr("u"),
		},
		Reviewers: []*RequiredReviewer{
			{
				Type: Ptr("User"),
				Reviewer: &User{
					AvatarURL:         Ptr("a"),
					Email:             Ptr("e"),
					EventsURL:         Ptr("e"),
					FollowersURL:      Ptr("f"),
					FollowingURL:      Ptr("f"),
					GistsURL:          Ptr("g"),
					GravatarID:        Ptr("g"),
					HTMLURL:           Ptr("h"),
					ID:                Ptr(int64(1)),
					Login:             Ptr("l"),
					Name:              Ptr("n"),
					NodeID:            Ptr("n"),
					OrganizationsURL:  Ptr("o"),
					ReceivedEventsURL: Ptr("r"),
					ReposURL:          Ptr("r"),
					SiteAdmin:         Ptr(false),
					StarredURL:        Ptr("s"),
					SubscriptionsURL:  Ptr("s"),
					Type:              Ptr("User"),
					URL:               Ptr("u"),
				},
			},
			{
				Type: Ptr("Team"),
				Reviewer: &Team{
					ID:   Ptr(int64(1)),
					Name: Ptr("n"),
					Slug: Ptr("s"),
				},
			},
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Since: Ptr("s"),
		WorkflowJobRun: &WorkflowJobRun{
			ID:          Ptr(int64(1)),
			Conclusion:  Ptr("c"),
			Environment: Ptr("e"),
			HTMLURL:     Ptr("h"),
			Name:        Ptr("n"),
			Status:      Ptr("s"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		WorkflowRun: &WorkflowRun{
			ID:         Ptr(int64(1)),
			Name:       Ptr("n"),
			NodeID:     Ptr("nid"),
			HeadBranch: Ptr("hb"),
			HeadSHA:    Ptr("hs"),
			RunNumber:  Ptr(1),
			RunAttempt: Ptr(1),
			Event:      Ptr("e"),
			Status:     Ptr("s"),
			Conclusion: Ptr("c"),
			WorkflowID: Ptr(int64(1)),
			URL:        Ptr("u"),
			HTMLURL:    Ptr("h"),
			PullRequests: []*PullRequest{
				{
					URL:    Ptr("u"),
					ID:     Ptr(int64(1)),
					Number: Ptr(1),
					Head: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("s"),
							Name: Ptr("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("u"),
							Name: Ptr("n"),
						},
					},
				},
			},
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            Ptr("j"),
			LogsURL:            Ptr("l"),
			CheckSuiteURL:      Ptr("c"),
			ArtifactsURL:       Ptr("a"),
			CancelURL:          Ptr("c"),
			RerunURL:           Ptr("r"),
			PreviousAttemptURL: Ptr("p"),
			HeadCommit: &HeadCommit{
				Message: Ptr("m"),
				Author: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
				URL:       Ptr("u"),
				Distinct:  Ptr(false),
				SHA:       Ptr("s"),
				ID:        Ptr("i"),
				TreeID:    Ptr("tid"),
				Timestamp: &Timestamp{referenceTime},
				Committer: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
			},
			WorkflowURL: Ptr("w"),
			Repository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
			HeadRepository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DeploymentStatusEvent{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &DeploymentStatusEvent{
		Deployment: &Deployment{
			URL:         Ptr("url"),
			ID:          Ptr(int64(1)),
			SHA:         Ptr("sha"),
			Ref:         Ptr("ref"),
			Task:        Ptr("t"),
			Payload:     jsonMsg,
			Environment: Ptr("e"),
			Description: Ptr("d"),
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			CreatedAt:     &Timestamp{referenceTime},
			UpdatedAt:     &Timestamp{referenceTime},
			StatusesURL:   Ptr("surl"),
			RepositoryURL: Ptr("rurl"),
			NodeID:        Ptr("nid"),
		},
		DeploymentStatus: &DeploymentStatus{
			ID:    Ptr(int64(1)),
			State: Ptr("s"),
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Description:    Ptr("s"),
			Environment:    Ptr("s"),
			NodeID:         Ptr("s"),
			CreatedAt:      &Timestamp{referenceTime},
			UpdatedAt:      &Timestamp{referenceTime},
			TargetURL:      Ptr("s"),
			DeploymentURL:  Ptr("s"),
			RepositoryURL:  Ptr("s"),
			EnvironmentURL: Ptr("s"),
			LogURL:         Ptr("s"),
			URL:            Ptr("s"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestDiscussionCommentEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DiscussionCommentEvent{}, "{}")

	u := &DiscussionCommentEvent{
		Comment: &CommentDiscussion{
			AuthorAssociation: Ptr("aa"),
			Body:              Ptr("bo"),
			ChildCommentCount: Ptr(1),
			CreatedAt:         &Timestamp{referenceTime},
			DiscussionID:      Ptr(int64(1)),
			HTMLURL:           Ptr("hurl"),
			ID:                Ptr(int64(1)),
			NodeID:            Ptr("nid"),
			ParentID:          Ptr(int64(1)),
			Reactions: &Reactions{
				TotalCount: Ptr(1),
				PlusOne:    Ptr(1),
				MinusOne:   Ptr(1),
				Laugh:      Ptr(1),
				Confused:   Ptr(1),
				Heart:      Ptr(1),
				Hooray:     Ptr(1),
				Rocket:     Ptr(1),
				Eyes:       Ptr(1),
				URL:        Ptr("url"),
			},
			RepositoryURL: Ptr("rurl"),
			UpdatedAt:     &Timestamp{referenceTime},
			User: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
		},
		Discussion: &Discussion{
			RepositoryURL: Ptr("rurl"),
			DiscussionCategory: &DiscussionCategory{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("nid"),
				RepositoryID: Ptr(int64(1)),
				Emoji:        Ptr("emoji"),
				Name:         Ptr("name"),
				Description:  Ptr("description"),
				CreatedAt:    &Timestamp{referenceTime},
				UpdatedAt:    &Timestamp{referenceTime},
				Slug:         Ptr("slug"),
				IsAnswerable: Ptr(false),
			},
			HTMLURL: Ptr("hurl"),
			ID:      Ptr(int64(1)),
			NodeID:  Ptr("nurl"),
			Number:  Ptr(1),
			Title:   Ptr("title"),
			User: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			State:             Ptr("st"),
			Locked:            Ptr(false),
			Comments:          Ptr(1),
			CreatedAt:         &Timestamp{referenceTime},
			UpdatedAt:         &Timestamp{referenceTime},
			AuthorAssociation: Ptr("aa"),
			Body:              Ptr("bo"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DiscussionEvent{}, "{}")

	u := &DiscussionEvent{
		Discussion: &Discussion{
			RepositoryURL: Ptr("rurl"),
			DiscussionCategory: &DiscussionCategory{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("nid"),
				RepositoryID: Ptr(int64(1)),
				Emoji:        Ptr("emoji"),
				Name:         Ptr("name"),
				Description:  Ptr("description"),
				CreatedAt:    &Timestamp{referenceTime},
				UpdatedAt:    &Timestamp{referenceTime},
				Slug:         Ptr("slug"),
				IsAnswerable: Ptr(false),
			},
			HTMLURL: Ptr("hurl"),
			ID:      Ptr(int64(1)),
			NodeID:  Ptr("nurl"),
			Number:  Ptr(1),
			Title:   Ptr("title"),
			User: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			State:             Ptr("st"),
			Locked:            Ptr(false),
			Comments:          Ptr(1),
			CreatedAt:         &Timestamp{referenceTime},
			UpdatedAt:         &Timestamp{referenceTime},
			AuthorAssociation: Ptr("aa"),
			Body:              Ptr("bo"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PackageEvent{}, "{}")

	u := &PackageEvent{
		Action: Ptr("a"),
		Package: &Package{
			ID:          Ptr(int64(1)),
			Name:        Ptr("n"),
			PackageType: Ptr("pt"),
			HTMLURL:     Ptr("hurl"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			PackageVersion: &PackageVersion{ID: Ptr(int64(1))},
			Registry:       &PackageRegistry{Name: Ptr("n")},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &PersonalAccessTokenRequestEvent{}, "{}")

	event := &PersonalAccessTokenRequestEvent{
		Action: Ptr("a"),
		PersonalAccessTokenRequest: &PersonalAccessTokenRequest{
			ID:    Ptr(int64(1)),
			Owner: &User{Login: Ptr("l")},
			PermissionsAdded: &PersonalAccessTokenPermissions{
				Org:  map[string]string{"organization_events": "read"},
				Repo: map[string]string{"security_events": "write"},
			},
			CreatedAt:           &Timestamp{referenceTime},
			TokenExpired:        Ptr(false),
			TokenExpiresAt:      &Timestamp{referenceTime},
			TokenLastUsedAt:     &Timestamp{referenceTime},
			RepositoryCount:     Ptr(int64(1)),
			RepositorySelection: Ptr("rs"),
			Repositories: []*Repository{
				{
					Name: Ptr("n"),
				},
			},
		},
		Org: &Organization{Name: Ptr("n")},
		Sender: &User{
			Login: Ptr("l"),
		},
		Installation: &Installation{
			ID: Ptr(int64(1)),
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
	t.Parallel()
	testJSONMarshal(t, &PingEvent{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"
	hookConfig := new(HookConfig)

	u := &PingEvent{
		Zen:    Ptr("z"),
		HookID: Ptr(int64(1)),
		Hook: &Hook{
			CreatedAt:    &Timestamp{referenceTime},
			UpdatedAt:    &Timestamp{referenceTime},
			URL:          Ptr("url"),
			ID:           Ptr(int64(1)),
			Type:         Ptr("t"),
			Name:         Ptr("n"),
			TestURL:      Ptr("tu"),
			PingURL:      Ptr("pu"),
			LastResponse: l,
			Config:       hookConfig,
			Events:       []string{"a"},
			Active:       Ptr(true),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			},
			"events": [
				"a"
			],
			"active": true
		},
		"installation": {
			"id": 1,
			"node_id": "nid",
			"client_id": "cid",
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

func TestRegistryPackageEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RegistryPackageEvent{}, "{}")

	u := &RegistryPackageEvent{
		Action: Ptr("a"),
		RegistryPackage: &Package{
			ID:          Ptr(int64(1)),
			Name:        Ptr("n"),
			PackageType: Ptr("pt"),
			HTMLURL:     Ptr("hurl"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			PackageVersion: &PackageVersion{ID: Ptr(int64(1))},
			Registry:       &PackageRegistry{Name: Ptr("n")},
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
	}

	want := `{
		"action": "a",
		"registry_package": {
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

func TestRepositoryDispatchEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryDispatchEvent{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &RepositoryDispatchEvent{
		Action:        Ptr("a"),
		Branch:        Ptr("b"),
		ClientPayload: jsonMsg,
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestRepositoryImportEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryImportEvent{}, "{}")

	u := &RepositoryImportEvent{
		Status: Ptr("success"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &RepositoryEvent{}, "{}")

	u := &RepositoryEvent{
		Action: Ptr("a"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &ReleaseEvent{}, "{}")

	u := &ReleaseEvent{
		Action: Ptr("a"),
		Release: &RepositoryRelease{
			Name:                   Ptr("n"),
			DiscussionCategoryName: Ptr("dcn"),
			ID:                     Ptr(int64(2)),
			CreatedAt:              &Timestamp{referenceTime},
			PublishedAt:            &Timestamp{referenceTime},
			URL:                    Ptr("url"),
			HTMLURL:                Ptr("htmlurl"),
			AssetsURL:              Ptr("assetsurl"),
			Assets:                 []*ReleaseAsset{{ID: Ptr(int64(1))}},
			UploadURL:              Ptr("uploadurl"),
			ZipballURL:             Ptr("zipballurl"),
			TarballURL:             Ptr("tarballurl"),
			Author:                 &User{Name: Ptr("octocat")},
			NodeID:                 Ptr("nid"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

func TestRepositoryRulesetEvent_Unmarshal(t *testing.T) {
	t.Parallel()

	enterprise := &Enterprise{
		ID:     Ptr(1),
		NodeID: Ptr("n"),
		Slug:   Ptr("e"),
		Name:   Ptr("e"),
	}

	installation := &Installation{
		ID:      Ptr(int64(1)),
		NodeID:  Ptr("n"),
		AppID:   Ptr(int64(1)),
		AppSlug: Ptr("a"),
	}

	organization := &Organization{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("n"),
		Name:   Ptr("o"),
	}

	repository := &Repository{
		ID:       Ptr(int64(1)),
		NodeID:   Ptr("n"),
		Name:     Ptr("r"),
		FullName: Ptr("o/r"),
	}

	sender := &User{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("n"),
		Login:  Ptr("l"),
	}

	tests := []struct {
		name  string
		json  string
		event *RepositoryRulesetEvent
	}{
		{"empty", `{}`, &RepositoryRulesetEvent{}},
		{
			"created",
			fmt.Sprintf(
				`{"action":"created","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~ALL"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type":"required_linear_history"},{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase","merge"]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"high_or_higher","alerts_threshold":"errors"}]}}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("created"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~ALL"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:              &EmptyRuleParameters{},
						Update:                &UpdateRuleParameters{},
						Deletion:              &EmptyRuleParameters{},
						RequiredLinearHistory: &EmptyRuleParameters{},
						PullRequest: &PullRequestRuleParameters{
							AllowedMergeMethods: []PullRequestMergeMethod{
								PullRequestMergeMethodSquash,
								PullRequestMergeMethodRebase,
								PullRequestMergeMethodMerge,
							},
							DismissStaleReviewsOnPush:      false,
							RequireCodeOwnerReview:         false,
							RequireLastPushApproval:        false,
							RequiredApprovingReviewCount:   2,
							RequiredReviewThreadResolution: false,
						},
						CodeScanning: &CodeScanningRuleParameters{
							CodeScanningTools: []*RuleCodeScanningTool{
								{
									AlertsThreshold:         CodeScanningAlertsThresholdErrors,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
									Tool:                    "CodeQL",
								},
							},
						},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
		{
			"edited",
			fmt.Sprintf(
				`{"action":"edited","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type": "required_signatures"},{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase"]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"medium_or_higher","alerts_threshold":"errors"}]}}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"changes":{"rules":{"added":[{"type": "required_signatures"}],"updated":[{"rule":{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase"]}},"changes":{"configuration":{"from":"{\\\"required_reviewers\\\":[],\\\"allowed_merge_methods\\\":[\\\"squash\\\",\\\"rebase\\\",\\\"merge\\\"],\\\"require_code_owner_review\\\":false,\\\"require_last_push_approval\\\":false,\\\"dismiss_stale_reviews_on_push\\\":false,\\\"required_approving_review_count\\\":2,\\\"authorized_dismissal_actors_only\\\":false,\\\"required_review_thread_resolution\\\":false,\\\"ignore_approvals_from_contributors\\\":false}"}}},{"rule":{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"medium_or_higher","alerts_threshold":"errors"}]}},"changes":{"configuration":{"from":"{\\\"code_scanning_tools\\\":[{\\\"tool\\\":\\\"CodeQL\\\",\\\"alerts_threshold\\\":\\\"errors\\\",\\\"security_alerts_threshold\\\":\\\"high_or_higher\\\"}]}"}}}],"deleted":[{"type":"required_linear_history"}]},"conditions":{"updated":[{"condition":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"changes":{"include":{"from":["~ALL"]}}}],"deleted":[]}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("edited"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:           &EmptyRuleParameters{},
						Update:             &UpdateRuleParameters{},
						Deletion:           &EmptyRuleParameters{},
						RequiredSignatures: &EmptyRuleParameters{},
						PullRequest: &PullRequestRuleParameters{
							AllowedMergeMethods: []PullRequestMergeMethod{
								PullRequestMergeMethodSquash,
								PullRequestMergeMethodRebase,
							},
							DismissStaleReviewsOnPush:      false,
							RequireCodeOwnerReview:         false,
							RequireLastPushApproval:        false,
							RequiredApprovingReviewCount:   2,
							RequiredReviewThreadResolution: false,
						},
						CodeScanning: &CodeScanningRuleParameters{
							CodeScanningTools: []*RuleCodeScanningTool{
								{
									AlertsThreshold:         CodeScanningAlertsThresholdErrors,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdMediumOrHigher,
									Tool:                    "CodeQL",
								},
							},
						},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Changes: &RepositoryRulesetChanges{
					Conditions: &RepositoryRulesetChangedConditions{
						Updated: []*RepositoryRulesetUpdatedConditions{
							{
								Condition: &RepositoryRulesetConditions{
									RefName: &RepositoryRulesetRefConditionParameters{
										Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
										Exclude: []string{},
									},
								},
								Changes: &RepositoryRulesetUpdatedCondition{
									Include: &RepositoryRulesetChangeSources{
										From: []string{"~ALL"},
									},
								},
							},
						},
						Deleted: []*RepositoryRulesetConditions{},
					},
					Rules: &RepositoryRulesetChangedRules{
						Added: []*RepositoryRule{{Type: RulesetRuleTypeRequiredSignatures}},
						Updated: []*RepositoryRulesetUpdatedRules{
							{
								Rule: &RepositoryRule{
									Type: RulesetRuleTypePullRequest,
									Parameters: &PullRequestRuleParameters{
										AllowedMergeMethods: []PullRequestMergeMethod{
											PullRequestMergeMethodSquash,
											PullRequestMergeMethodRebase,
										},
										DismissStaleReviewsOnPush:      false,
										RequireCodeOwnerReview:         false,
										RequireLastPushApproval:        false,
										RequiredApprovingReviewCount:   2,
										RequiredReviewThreadResolution: false,
									},
								},
								Changes: &RepositoryRulesetChangedRule{
									Configuration: &RepositoryRulesetChangeSource{
										From: Ptr(
											`{\"required_reviewers\":[],\"allowed_merge_methods\":[\"squash\",\"rebase\",\"merge\"],\"require_code_owner_review\":false,\"require_last_push_approval\":false,\"dismiss_stale_reviews_on_push\":false,\"required_approving_review_count\":2,\"authorized_dismissal_actors_only\":false,\"required_review_thread_resolution\":false,\"ignore_approvals_from_contributors\":false}`,
										),
									},
								},
							},
							{
								Rule: &RepositoryRule{
									Type: RulesetRuleTypeCodeScanning,
									Parameters: &CodeScanningRuleParameters{
										CodeScanningTools: []*RuleCodeScanningTool{
											{
												AlertsThreshold:         CodeScanningAlertsThresholdErrors,
												SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdMediumOrHigher,
												Tool:                    "CodeQL",
											},
										},
									},
								},
								Changes: &RepositoryRulesetChangedRule{
									Configuration: &RepositoryRulesetChangeSource{
										From: Ptr(
											`{\"code_scanning_tools\":[{\"tool\":\"CodeQL\",\"alerts_threshold\":\"errors\",\"security_alerts_threshold\":\"high_or_higher\"}]}`,
										),
									},
								},
							},
						},
						Deleted: []*RepositoryRule{{Type: RulesetRuleTypeRequiredLinearHistory}},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
		{
			"deleted",
			fmt.Sprintf(
				`{"action":"deleted","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type":"required_linear_history"}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("deleted"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:              &EmptyRuleParameters{},
						Update:                &UpdateRuleParameters{},
						Deletion:              &EmptyRuleParameters{},
						RequiredLinearHistory: &EmptyRuleParameters{},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := &RepositoryRulesetEvent{}
			err := json.Unmarshal([]byte(test.json), got)
			if err != nil {
				t.Errorf("Unable to unmarshal JSON %v: %v", test.json, err)
			}

			if diff := cmp.Diff(test.event, got); diff != "" {
				t.Errorf("json.Unmarshal returned:\n%#v\nwant:\n%#v\ndiff:\n%v", got, test.event, diff)
			}
		})
	}
}

func TestContentReferenceEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ContentReferenceEvent{}, "{}")

	u := &ContentReferenceEvent{
		Action: Ptr("a"),
		ContentReference: &ContentReference{
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("nid"),
			Reference: Ptr("ref"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &MemberEvent{}, "{}")

	u := &MemberEvent{
		Action: Ptr("a"),
		Member: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Changes: &MemberChanges{
			Permission: &MemberChangesPermission{
				From: Ptr("f"),
				To:   Ptr("t"),
			},
			RoleName: &MemberChangesRoleName{
				From: Ptr("f"),
				To:   Ptr("t"),
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &MembershipEvent{}, "{}")

	u := &MembershipEvent{
		Action: Ptr("a"),
		Scope:  Ptr("s"),
		Member: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Team: &Team{
			ID:              Ptr(int64(1)),
			NodeID:          Ptr("n"),
			Name:            Ptr("n"),
			Description:     Ptr("d"),
			URL:             Ptr("u"),
			Slug:            Ptr("s"),
			Permission:      Ptr("p"),
			Privacy:         Ptr("p"),
			MembersCount:    Ptr(1),
			ReposCount:      Ptr(1),
			MembersURL:      Ptr("m"),
			RepositoriesURL: Ptr("r"),
			Organization: &Organization{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				AvatarURL: Ptr("a"),
				HTMLURL:   Ptr("h"),
				Name:      Ptr("n"),
				Company:   Ptr("c"),
				Blog:      Ptr("b"),
				Location:  Ptr("l"),
				Email:     Ptr("e"),
			},
			Parent: &Team{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("n"),
				Name:         Ptr("n"),
				Description:  Ptr("d"),
				URL:          Ptr("u"),
				Slug:         Ptr("s"),
				Permission:   Ptr("p"),
				Privacy:      Ptr("p"),
				MembersCount: Ptr(1),
				ReposCount:   Ptr(1),
			},
			LDAPDN: Ptr("l"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &MergeGroupEvent{}, "{}")

	u := &MergeGroupEvent{
		Action: Ptr("a"),
		Reason: Ptr("r"),
		MergeGroup: &MergeGroup{
			HeadSHA:    Ptr("hs"),
			HeadRef:    Ptr("hr"),
			BaseSHA:    Ptr("bs"),
			BaseRef:    Ptr("br"),
			HeadCommit: &Commit{NodeID: Ptr("nid")},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
	}

	want := `{
		"action": "a",
		"reason": "r",
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &OrgBlockEvent{}, "{}")

	u := &OrgBlockEvent{
		Action: Ptr("a"),
		BlockedUser: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &GollumEvent{}, "{}")

	u := &GollumEvent{
		Pages: []*Page{
			{
				PageName: Ptr("pn"),
				Title:    Ptr("t"),
				Summary:  Ptr("s"),
				Action:   Ptr("a"),
				SHA:      Ptr("sha"),
				HTMLURL:  Ptr("hu"),
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &WorkflowRunEvent{}, "{}")

	u := &WorkflowRunEvent{
		Action: Ptr("a"),
		Workflow: &Workflow{
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("nid"),
			Name:      Ptr("n"),
			Path:      Ptr("p"),
			State:     Ptr("s"),
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			URL:       Ptr("u"),
			HTMLURL:   Ptr("h"),
			BadgeURL:  Ptr("b"),
		},
		WorkflowRun: &WorkflowRun{
			ID:         Ptr(int64(1)),
			Name:       Ptr("n"),
			NodeID:     Ptr("nid"),
			HeadBranch: Ptr("hb"),
			HeadSHA:    Ptr("hs"),
			RunNumber:  Ptr(1),
			RunAttempt: Ptr(1),
			Event:      Ptr("e"),
			Status:     Ptr("s"),
			Conclusion: Ptr("c"),
			WorkflowID: Ptr(int64(1)),
			URL:        Ptr("u"),
			HTMLURL:    Ptr("h"),
			PullRequests: []*PullRequest{
				{
					URL:    Ptr("u"),
					ID:     Ptr(int64(1)),
					Number: Ptr(1),
					Head: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("s"),
							Name: Ptr("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("u"),
							Name: Ptr("n"),
						},
					},
				},
			},
			CreatedAt:          &Timestamp{referenceTime},
			UpdatedAt:          &Timestamp{referenceTime},
			RunStartedAt:       &Timestamp{referenceTime},
			JobsURL:            Ptr("j"),
			LogsURL:            Ptr("l"),
			CheckSuiteURL:      Ptr("c"),
			ArtifactsURL:       Ptr("a"),
			CancelURL:          Ptr("c"),
			RerunURL:           Ptr("r"),
			PreviousAttemptURL: Ptr("p"),
			HeadCommit: &HeadCommit{
				Message: Ptr("m"),
				Author: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
				URL:       Ptr("u"),
				Distinct:  Ptr(false),
				SHA:       Ptr("s"),
				ID:        Ptr("i"),
				TreeID:    Ptr("tid"),
				Timestamp: &Timestamp{referenceTime},
				Committer: &CommitAuthor{
					Name:  Ptr("n"),
					Email: Ptr("e"),
					Login: Ptr("l"),
				},
			},
			WorkflowURL: Ptr("w"),
			Repository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
			HeadRepository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &WorkflowDispatchEvent{}, "{}")

	i := make(map[string]any)
	i["key"] = "value"

	jsonMsg, _ := json.Marshal(i)
	u := &WorkflowDispatchEvent{
		Inputs:   jsonMsg,
		Ref:      Ptr("r"),
		Workflow: Ptr("w"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestWatchEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WatchEvent{}, "{}")

	u := &WatchEvent{
		Action: Ptr("a"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &UserEvent{}, "{}")

	u := &UserEvent{
		User: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		// The action performed. Possible values are: "created" or "deleted".
		Action: Ptr("a"),
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &CheckRunEvent{}, "{}")

	r := &CheckRunEvent{
		CheckRun: &CheckRun{
			ID:          Ptr(int64(1)),
			NodeID:      Ptr("n"),
			HeadSHA:     Ptr("h"),
			ExternalID:  Ptr("1"),
			URL:         Ptr("u"),
			HTMLURL:     Ptr("u"),
			DetailsURL:  Ptr("u"),
			Status:      Ptr("s"),
			Conclusion:  Ptr("c"),
			StartedAt:   &Timestamp{referenceTime},
			CompletedAt: &Timestamp{referenceTime},
			Output: &CheckRunOutput{
				Annotations: []*CheckRunAnnotation{
					{
						AnnotationLevel: Ptr("a"),
						EndLine:         Ptr(1),
						Message:         Ptr("m"),
						Path:            Ptr("p"),
						RawDetails:      Ptr("r"),
						StartLine:       Ptr(1),
						Title:           Ptr("t"),
					},
				},
				AnnotationsCount: Ptr(1),
				AnnotationsURL:   Ptr("a"),
				Images: []*CheckRunImage{
					{
						Alt:      Ptr("a"),
						ImageURL: Ptr("i"),
						Caption:  Ptr("c"),
					},
				},
				Title:   Ptr("t"),
				Summary: Ptr("s"),
				Text:    Ptr("t"),
			},
			Name: Ptr("n"),
			CheckSuite: &CheckSuite{
				ID: Ptr(int64(1)),
			},
			App: &App{
				ID:     Ptr(int64(1)),
				NodeID: Ptr("n"),
				Owner: &User{
					Login:     Ptr("l"),
					ID:        Ptr(int64(1)),
					NodeID:    Ptr("n"),
					URL:       Ptr("u"),
					ReposURL:  Ptr("r"),
					EventsURL: Ptr("e"),
					AvatarURL: Ptr("a"),
				},
				Name:        Ptr("n"),
				Description: Ptr("d"),
				HTMLURL:     Ptr("h"),
				ExternalURL: Ptr("u"),
				CreatedAt:   &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
			},
			PullRequests: []*PullRequest{
				{
					URL:    Ptr("u"),
					ID:     Ptr(int64(1)),
					Number: Ptr(1),
					Head: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("s"),
							Name: Ptr("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("u"),
							Name: Ptr("n"),
						},
					},
				},
			},
		},
		Action: Ptr("a"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &CheckSuiteEvent{}, "{}")

	r := &CheckSuiteEvent{
		CheckSuite: &CheckSuite{
			ID:         Ptr(int64(1)),
			NodeID:     Ptr("n"),
			HeadBranch: Ptr("h"),
			HeadSHA:    Ptr("h"),
			URL:        Ptr("u"),
			BeforeSHA:  Ptr("b"),
			AfterSHA:   Ptr("a"),
			Status:     Ptr("s"),
			Conclusion: Ptr("c"),
			App: &App{
				ID:     Ptr(int64(1)),
				NodeID: Ptr("n"),
				Owner: &User{
					Login:     Ptr("l"),
					ID:        Ptr(int64(1)),
					NodeID:    Ptr("n"),
					URL:       Ptr("u"),
					ReposURL:  Ptr("r"),
					EventsURL: Ptr("e"),
					AvatarURL: Ptr("a"),
				},
				Name:        Ptr("n"),
				Description: Ptr("d"),
				HTMLURL:     Ptr("h"),
				ExternalURL: Ptr("u"),
				CreatedAt:   &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
			},
			Repository: &Repository{
				ID: Ptr(int64(1)),
			},
			PullRequests: []*PullRequest{
				{
					URL:    Ptr("u"),
					ID:     Ptr(int64(1)),
					Number: Ptr(1),
					Head: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("s"),
							Name: Ptr("n"),
						},
					},
					Base: &PullRequestBranch{
						Ref: Ptr("r"),
						SHA: Ptr("s"),
						Repo: &Repository{
							ID:   Ptr(int64(1)),
							URL:  Ptr("u"),
							Name: Ptr("n"),
						},
					},
				},
			},
			HeadCommit: &Commit{
				SHA: Ptr("s"),
			},
		},
		Action: Ptr("a"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DeployKeyEvent{}, "{}")

	u := &DeployKeyEvent{
		Action: Ptr("a"),
		Key: &Key{
			ID:        Ptr(int64(1)),
			Key:       Ptr("k"),
			URL:       Ptr("k"),
			Title:     Ptr("k"),
			ReadOnly:  Ptr(false),
			Verified:  Ptr(false),
			CreatedAt: &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			AvatarURL: Ptr("a"),
			URL:       Ptr("u"),
			EventsURL: Ptr("e"),
			ReposURL:  Ptr("r"),
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
	t.Parallel()
	testJSONMarshal(t, &MetaEvent{}, "{}")

	v := make(map[string]any)
	v["a"] = "b"
	hookConfig := &HookConfig{
		ContentType: Ptr("json"),
	}

	u := &MetaEvent{
		Action: Ptr("a"),
		HookID: Ptr(int64(1)),
		Hook: &Hook{
			CreatedAt:    &Timestamp{referenceTime},
			UpdatedAt:    &Timestamp{referenceTime},
			URL:          Ptr("u"),
			ID:           Ptr(int64(1)),
			Type:         Ptr("t"),
			Name:         Ptr("n"),
			TestURL:      Ptr("tu"),
			PingURL:      Ptr("pu"),
			LastResponse: v,
			Config:       hookConfig,
			Events:       []string{"a"},
			Active:       Ptr(true),
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
	t.Parallel()
	testJSONMarshal(t, &RequestedAction{}, `{"identifier": ""}`)

	r := &RequestedAction{
		Identifier: "i",
	}

	want := `{
		"identifier": "i"
	}`

	testJSONMarshal(t, r, want)
}

func TestCreateEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CreateEvent{}, "{}")

	r := &CreateEvent{
		Ref:          Ptr("r"),
		RefType:      Ptr("rt"),
		MasterBranch: Ptr("mb"),
		Description:  Ptr("d"),
		PusherType:   Ptr("pt"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

func TestCustomPropertyEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CustomPropertyEvent{}, "{}")

	r := &CustomPropertyEvent{
		Action: Ptr("created"),
		Definition: &CustomProperty{
			PropertyName:     Ptr("name"),
			ValueType:        PropertyValueTypeSingleSelect,
			SourceType:       Ptr("enterprise"),
			Required:         Ptr(true),
			DefaultValue:     "production",
			Description:      Ptr("Prod or dev environment"),
			AllowedValues:    []string{"production", "development"},
			ValuesEditableBy: Ptr("org_actors"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
	}

	want := `{
		"action": "created",
		"definition": {
			"property_name": "name",
          	"source_type": "enterprise",
          	"value_type": "single_select",
          	"required": true,
          	"default_value": "production",
          	"description": "Prod or dev environment",
          	"allowed_values": [
            	"production",
            	"development"
          	],
          	"values_editable_by": "org_actors"
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
			"client_id": "cid",
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

func TestCustomPropertyValuesEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CustomPropertyValuesEvent{}, "{}")

	r := &CustomPropertyValuesEvent{
		Action: Ptr("updated"),
		NewPropertyValues: []*CustomPropertyValue{
			{
				PropertyName: "environment",
				Value:        "production",
			},
		},
		OldPropertyValues: []*CustomPropertyValue{
			{
				PropertyName: "environment",
				Value:        "staging",
			},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
	}

	want := `{
		"action": "updated",
		"new_property_values": [{
			"property_name": "environment",
			"value": "production"
        }],
		"old_property_values": [{
			"property_name": "environment",
			"value": "staging"
        }],
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DeleteEvent{}, "{}")

	r := &DeleteEvent{
		Ref:        Ptr("r"),
		RefType:    Ptr("rt"),
		PusherType: Ptr("pt"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &DependabotAlertEvent{}, "{}")

	e := &DependabotAlertEvent{
		Action: Ptr("a"),
		Alert: &DependabotAlert{
			Number: Ptr(1),
			State:  Ptr("s"),
			Dependency: &Dependency{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("e"),
					Name:      Ptr("n"),
				},
				ManifestPath: Ptr("mp"),
				Scope:        Ptr("s"),
			},
			SecurityAdvisory: &DependabotSecurityAdvisory{
				GHSAID:      Ptr("ghsaid"),
				CVEID:       Ptr("cveid"),
				Summary:     Ptr("s"),
				Description: Ptr("d"),
				Vulnerabilities: []*AdvisoryVulnerability{
					{
						Package: &VulnerabilityPackage{
							Ecosystem: Ptr("e"),
							Name:      Ptr("n"),
						},
						Severity: Ptr("s"),
					},
				},
				Severity: Ptr("s"),
				CVSS: &AdvisoryCVSS{
					Score:        Ptr(1.0),
					VectorString: Ptr("vs"),
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: Ptr("cweid"),
						Name:  Ptr("n"),
					},
				},
				Identifiers: []*AdvisoryIdentifier{
					{
						Value: Ptr("v"),
						Type:  Ptr("t"),
					},
				},
				References: []*AdvisoryReference{
					{
						URL: Ptr("u"),
					},
				},
				PublishedAt: &Timestamp{referenceTime},
				UpdatedAt:   &Timestamp{referenceTime},
				WithdrawnAt: &Timestamp{referenceTime},
			},
			SecurityVulnerability: &AdvisoryVulnerability{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("e"),
					Name:      Ptr("n"),
				},
				Severity:               Ptr("s"),
				VulnerableVersionRange: Ptr("vvr"),
				FirstPatchedVersion: &FirstPatchedVersion{
					Identifier: Ptr("i"),
				},
			},
			URL:         Ptr("u"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			DismissedAt: &Timestamp{referenceTime},
			DismissedBy: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			DismissedReason:  Ptr("dr"),
			DismissedComment: Ptr("dc"),
			FixedAt:          &Timestamp{referenceTime},
			AutoDismissedAt:  &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &ForkEvent{}, "{}")

	u := &ForkEvent{
		Forkee: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &GitHubAppAuthorizationEvent{}, "{}")

	u := &GitHubAppAuthorizationEvent{
		Action: Ptr("a"),
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
	t.Parallel()
	testJSONMarshal(t, &InstallationEvent{}, "{}")

	u := &InstallationEvent{
		Action: Ptr("a"),
		Repositories: []*Repository{
			{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &HeadCommit{}, "{}")

	u := &HeadCommit{
		Message: Ptr("m"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("n"),
			Email: Ptr("e"),
			Login: Ptr("u"),
		},
		URL:       Ptr("u"),
		Distinct:  Ptr(true),
		SHA:       Ptr("s"),
		ID:        Ptr("id"),
		TreeID:    Ptr("tid"),
		Timestamp: &Timestamp{referenceTime},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("n"),
			Email: Ptr("e"),
			Login: Ptr("u"),
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
	t.Parallel()
	testJSONMarshal(t, &PushEventRepository{}, "{}")

	u := &PushEventRepository{
		ID:       Ptr(int64(1)),
		NodeID:   Ptr("nid"),
		Name:     Ptr("n"),
		FullName: Ptr("fn"),
		Owner: &User{
			Login:       Ptr("l"),
			ID:          Ptr(int64(1)),
			AvatarURL:   Ptr("a"),
			GravatarID:  Ptr("g"),
			Name:        Ptr("n"),
			Company:     Ptr("c"),
			Blog:        Ptr("b"),
			Location:    Ptr("l"),
			Email:       Ptr("e"),
			Hireable:    Ptr(true),
			PublicRepos: Ptr(1),
			Followers:   Ptr(1),
			Following:   Ptr(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         Ptr("u"),
		},
		Private:         Ptr(true),
		Description:     Ptr("d"),
		Fork:            Ptr(true),
		CreatedAt:       &Timestamp{referenceTime},
		PushedAt:        &Timestamp{referenceTime},
		UpdatedAt:       &Timestamp{referenceTime},
		Homepage:        Ptr("h"),
		PullsURL:        Ptr("p"),
		Size:            Ptr(1),
		StargazersCount: Ptr(1),
		WatchersCount:   Ptr(1),
		Language:        Ptr("l"),
		HasIssues:       Ptr(true),
		HasDownloads:    Ptr(true),
		HasWiki:         Ptr(true),
		HasPages:        Ptr(true),
		ForksCount:      Ptr(1),
		Archived:        Ptr(true),
		Disabled:        Ptr(true),
		OpenIssuesCount: Ptr(1),
		DefaultBranch:   Ptr("d"),
		MasterBranch:    Ptr("m"),
		Organization:    Ptr("o"),
		URL:             Ptr("u"),
		ArchiveURL:      Ptr("a"),
		HTMLURL:         Ptr("h"),
		StatusesURL:     Ptr("s"),
		GitURL:          Ptr("g"),
		SSHURL:          Ptr("s"),
		CloneURL:        Ptr("c"),
		SVNURL:          Ptr("s"),
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
	t.Parallel()
	testJSONMarshal(t, &PushEventRepoOwner{}, "{}")

	u := &PushEventRepoOwner{
		Name:  Ptr("n"),
		Email: Ptr("e"),
	}

	want := `{
		"name": "n",
		"email": "e"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectV2Event_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2Event{}, "{}")

	u := &ProjectV2Event{
		Action: Ptr("a"),
		ProjectsV2: &ProjectV2{
			ID:     Ptr(int64(1)),
			NodeID: Ptr("nid"),
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Title:            Ptr("t"),
			Description:      Ptr("d"),
			Public:           Ptr(true),
			ClosedAt:         &Timestamp{referenceTime},
			CreatedAt:        &Timestamp{referenceTime},
			UpdatedAt:        &Timestamp{referenceTime},
			DeletedAt:        &Timestamp{referenceTime},
			Number:           Ptr(1),
			ShortDescription: Ptr("sd"),
			DeletedBy: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &ProjectV2ItemEvent{}, "{}")

	u := &ProjectV2ItemEvent{
		Action: Ptr("a"),
		Changes: &ProjectV2ItemChange{
			ArchivedAt: &ArchivedAt{
				From: &Timestamp{referenceTime},
				To:   &Timestamp{referenceTime},
			},
		},
		ProjectV2Item: &ProjectV2Item{
			ID:            Ptr(int64(1)),
			NodeID:        Ptr("nid"),
			ProjectNodeID: Ptr("pnid"),
			ContentNodeID: Ptr("cnid"),
			ContentType:   Ptr(ProjectV2ItemContentType("ct")),
			Creator: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			CreatedAt:  &Timestamp{referenceTime},
			UpdatedAt:  &Timestamp{referenceTime},
			ArchivedAt: &Timestamp{referenceTime},
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PullRequestEvent{}, "{}")

	u := &PullRequestEvent{
		Action: Ptr("a"),
		Assignee: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Number:      Ptr(1),
		PullRequest: &PullRequest{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		RequestedReviewer: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		RequestedTeam: &Team{ID: Ptr(int64(1))},
		Label:         &Label{ID: Ptr(int64(1))},
		Reason:        Ptr("CI_FAILURE"),
		Before:        Ptr("before"),
		After:         Ptr("after"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		PerformedViaGithubApp: &App{
			ID:          Ptr(int64(1)),
			NodeID:      Ptr("n"),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			Description: Ptr("d"),
			ExternalURL: Ptr("e"),
			HTMLURL:     Ptr("h"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
		"reason": "CI_FAILURE",
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
			"client_id": "cid",
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestPullRequestReviewCommentEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestReviewCommentEvent{}, "{}")

	u := &PullRequestReviewCommentEvent{
		Action:      Ptr("a"),
		PullRequest: &PullRequest{ID: Ptr(int64(1))},
		Comment:     &PullRequestComment{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PullRequestReviewThreadEvent{}, "{}")

	u := &PullRequestReviewThreadEvent{
		Action:      Ptr("a"),
		PullRequest: &PullRequest{ID: Ptr(int64(1))},
		Thread: &PullRequestThread{
			Comments: []*PullRequestComment{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &PullRequestTargetEvent{}, "{}")

	u := &PullRequestTargetEvent{
		Action: Ptr("a"),
		Assignee: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Number:      Ptr(1),
		PullRequest: &PullRequest{ID: Ptr(int64(1))},
		Changes: &EditChange{
			Title: &EditTitle{
				From: Ptr("TitleFrom"),
			},
			Body: &EditBody{
				From: Ptr("BodyFrom"),
			},
			Base: &EditBase{
				Ref: &EditRef{
					From: Ptr("BaseRefFrom"),
				},
				SHA: &EditSHA{
					From: Ptr("BaseSHAFrom"),
				},
			},
		},
		RequestedReviewer: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		RequestedTeam: &Team{ID: Ptr(int64(1))},
		Label:         &Label{ID: Ptr(int64(1))},
		Before:        Ptr("before"),
		After:         Ptr("after"),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		PerformedViaGithubApp: &App{
			ID:          Ptr(int64(1)),
			NodeID:      Ptr("n"),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			Description: Ptr("d"),
			ExternalURL: Ptr("e"),
			HTMLURL:     Ptr("h"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &RepositoryVulnerabilityAlertEvent{}, "{}")

	u := &RepositoryVulnerabilityAlertEvent{
		Action: Ptr("a"),
		Alert: &RepositoryVulnerabilityAlert{
			ID:                  Ptr(int64(1)),
			AffectedRange:       Ptr("ar"),
			AffectedPackageName: Ptr("apn"),
			ExternalReference:   Ptr("er"),
			ExternalIdentifier:  Ptr("ei"),
			FixedIn:             Ptr("fi"),
			Dismisser: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			DismissReason: Ptr("dr"),
			DismissedAt:   &Timestamp{referenceTime},
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
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
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlertEvent{}, "{}")

	u := &SecretScanningAlertEvent{
		Action: Ptr("a"),
		Alert: &SecretScanningAlert{
			Number:     Ptr(1),
			SecretType: Ptr("t"),
			Resolution: Ptr("r"),
			ResolvedBy: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			ResolvedAt: &Timestamp{referenceTime},
		},
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
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
			"client_id": "cid",
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

func TestSecretScanningAlertLocationEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlertLocationEvent{}, "{}")
	u := &SecretScanningAlertLocationEvent{
		Action: Ptr("created"),
		Alert: &SecretScanningAlert{
			Number:     Ptr(10),
			CreatedAt:  &Timestamp{referenceTime},
			UpdatedAt:  &Timestamp{referenceTime},
			URL:        Ptr("a"),
			HTMLURL:    Ptr("a"),
			SecretType: Ptr("mailchimp_api_key"),
		},
		Location: &SecretScanningAlertLocation{
			Type: Ptr("blob"),
			Details: &SecretScanningAlertLocationDetails{
				Path:        Ptr("path/to/file"),
				Startline:   Ptr(10),
				EndLine:     Ptr(20),
				StartColumn: Ptr(1),
				EndColumn:   Ptr(2),
				BlobSHA:     Ptr("d6e4c75c141dbacecc279b721b8bsomeSHA"),
				BlobURL:     Ptr("a"),
				CommitSHA:   Ptr("d6e4c75c141dbacecc279b721b8bsomeSHA"),
				CommitURL:   Ptr("a"),
			},
		},
		Repo: &Repository{
			ID:     Ptr(int64(12345)),
			NodeID: Ptr("MDEwOlJlcG9zaXRvcnkxMjM0NQ=="),
			Name:   Ptr("example-repo"),
		},
		Organization: &Organization{
			Login: Ptr("example-org"),
			ID:    Ptr(int64(67890)),
		},
		Sender: &User{
			Login: Ptr("example-user"),
			ID:    Ptr(int64(1111)),
		},
		Installation: &Installation{
			ID: Ptr(int64(2222)),
		},
	}

	want := `{
		"action": "created",
		"alert": {
			"number": 10,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"url": "a",
			"html_url": "a",
			"secret_type": "mailchimp_api_key"
		},
		"location": {

			"type": "blob",
			"details": {
				"path": "path/to/file",
				"start_line": 10,
				"end_line": 20,
				"start_column": 1,
				"end_column": 2,
				"blob_sha": "d6e4c75c141dbacecc279b721b8bsomeSHA",
				"blob_url": "a",
				"commit_sha": "d6e4c75c141dbacecc279b721b8bsomeSHA",
				"commit_url": "a"
			}
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

func TestSecurityAdvisoryEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecurityAdvisoryEvent{}, "{}")
	u := &SecurityAdvisoryEvent{
		Action: Ptr("published"),
		SecurityAdvisory: &SecurityAdvisory{
			CVSS: &AdvisoryCVSS{
				Score:        Ptr(1.0),
				VectorString: Ptr("vs"),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: Ptr("cweid"),
					Name:  Ptr("n"),
				},
			},
			GHSAID:      Ptr("GHSA-rf4j-j272-some"),
			Summary:     Ptr("Siuuuuuuuuu"),
			Description: Ptr("desc"),
			Severity:    Ptr("moderate"),
			Identifiers: []*AdvisoryIdentifier{
				{
					Value: Ptr("GHSA-rf4j-j272-some"),
					Type:  Ptr("GHSA"),
				},
			},
			References: []*AdvisoryReference{
				{
					URL: Ptr("https://some-url"),
				},
			},
			PublishedAt: &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
			WithdrawnAt: nil,
			Vulnerabilities: []*AdvisoryVulnerability{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: Ptr("ucl"),
						Name:      Ptr("penaldo"),
					},
					Severity:               Ptr("moderate"),
					VulnerableVersionRange: Ptr(">= 2.0.0, < 2.0.2"),
					FirstPatchedVersion: &FirstPatchedVersion{
						Identifier: Ptr("2.0.2"),
					},
				},
			},
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
			"client_id": "cid",
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
	t.Parallel()
	testJSONMarshal(t, &SecurityAndAnalysisEvent{}, "{}")

	u := &SecurityAndAnalysisEvent{
		Changes: &SecurityAndAnalysisChange{
			From: &SecurityAndAnalysisChangeFrom{
				SecurityAndAnalysis: &SecurityAndAnalysis{
					AdvancedSecurity: &AdvancedSecurity{
						Status: Ptr("enabled"),
					},
					SecretScanning: &SecretScanning{
						Status: Ptr("enabled"),
					},
					SecretScanningPushProtection: &SecretScanningPushProtection{
						Status: Ptr("enabled"),
					},
					DependabotSecurityUpdates: &DependabotSecurityUpdates{
						Status: Ptr("enabled"),
					},
				},
			},
		},
		Enterprise: &Enterprise{
			ID:          Ptr(1),
			Slug:        Ptr("s"),
			Name:        Ptr("n"),
			NodeID:      Ptr("nid"),
			AvatarURL:   Ptr("au"),
			Description: Ptr("d"),
			WebsiteURL:  Ptr("wu"),
			HTMLURL:     Ptr("hu"),
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		Installation: &Installation{
			ID:       Ptr(int64(1)),
			NodeID:   Ptr("nid"),
			ClientID: Ptr("cid"),
			AppID:    Ptr(int64(1)),
			AppSlug:  Ptr("as"),
			TargetID: Ptr(int64(1)),
			Account: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			AccessTokensURL:     Ptr("atu"),
			RepositoriesURL:     Ptr("ru"),
			HTMLURL:             Ptr("hu"),
			TargetType:          Ptr("tt"),
			SingleFileName:      Ptr("sfn"),
			RepositorySelection: Ptr("rs"),
			Events:              []string{"e"},
			SingleFilePaths:     []string{"s"},
			Permissions: &InstallationPermissions{
				Actions:                       Ptr("a"),
				Administration:                Ptr("ad"),
				Checks:                        Ptr("c"),
				Contents:                      Ptr("co"),
				ContentReferences:             Ptr("cr"),
				Deployments:                   Ptr("d"),
				Environments:                  Ptr("e"),
				Issues:                        Ptr("i"),
				Metadata:                      Ptr("md"),
				Members:                       Ptr("m"),
				OrganizationAdministration:    Ptr("oa"),
				OrganizationHooks:             Ptr("oh"),
				OrganizationPlan:              Ptr("op"),
				OrganizationPreReceiveHooks:   Ptr("opr"),
				OrganizationProjects:          Ptr("op"),
				OrganizationSecrets:           Ptr("os"),
				OrganizationSelfHostedRunners: Ptr("osh"),
				OrganizationUserBlocking:      Ptr("oub"),
				Packages:                      Ptr("pkg"),
				Pages:                         Ptr("pg"),
				PullRequests:                  Ptr("pr"),
				RepositoryHooks:               Ptr("rh"),
				RepositoryProjects:            Ptr("rp"),
				RepositoryPreReceiveHooks:     Ptr("rprh"),
				Secrets:                       Ptr("s"),
				SecretScanningAlerts:          Ptr("ssa"),
				SecurityEvents:                Ptr("se"),
				SingleFile:                    Ptr("sf"),
				Statuses:                      Ptr("s"),
				TeamDiscussions:               Ptr("td"),
				VulnerabilityAlerts:           Ptr("va"),
				Workflows:                     Ptr("w"),
			},
			CreatedAt:              &Timestamp{referenceTime},
			UpdatedAt:              &Timestamp{referenceTime},
			HasMultipleSingleFiles: Ptr(false),
			SuspendedBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			SuspendedAt: &Timestamp{referenceTime},
		},
		Organization: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Sender: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
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
			"client_id": "cid",
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

func TestCodeScanningAlertEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CodeScanningAlertEvent{}, "{}")

	u := &CodeScanningAlertEvent{
		Action: Ptr("reopened"),
		Alert: &Alert{
			Number: Ptr(10),
			Rule: &Rule{
				ID:              Ptr("Style/FrozenStringLiteralComment"),
				Severity:        Ptr("note"),
				Description:     Ptr("desc"),
				FullDescription: Ptr("full desc"),
				Tags:            []string{"style"},
				Help:            Ptr("help"),
			},
			Tool: &Tool{
				Name:    Ptr("Rubocop"),
				Version: nil,
			},
			CreatedAt: &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			FixedAt:   nil,
			State:     Ptr("open"),
			URL:       Ptr("a"),
			HTMLURL:   Ptr("a"),
			Instances: []*MostRecentInstance{
				{
					Ref:         Ptr("refs/heads/main"),
					AnalysisKey: Ptr(".github/workflows/workflow.yml:upload"),
					Environment: Ptr("{}"),
					State:       Ptr("open"),
				},
			},
			DismissedBy:     nil,
			DismissedAt:     nil,
			DismissedReason: nil,
		},
		Ref:       Ptr("refs/heads/main"),
		CommitOID: Ptr("d6e4c75c141dbacecc279b721b8bsomeSHA"),
		Repo: &Repository{
			ID:     Ptr(int64(1234234535)),
			NodeID: Ptr("MDEwOlJlcG9zaXRvcnkxODY4NT=="),
			Owner: &User{
				Login:             Ptr("Codertocat"),
				ID:                Ptr(int64(21031067)),
				NodeID:            Ptr("MDQ6VXNlcjIxMDMxMDY3"),
				AvatarURL:         Ptr("a"),
				GravatarID:        Ptr(""),
				URL:               Ptr("a"),
				HTMLURL:           Ptr("a"),
				Type:              Ptr("User"),
				SiteAdmin:         Ptr(false),
				FollowersURL:      Ptr("a"),
				FollowingURL:      Ptr("a"),
				EventsURL:         Ptr("a"),
				GistsURL:          Ptr("a"),
				OrganizationsURL:  Ptr("a"),
				ReceivedEventsURL: Ptr("a"),
				ReposURL:          Ptr("a"),
				StarredURL:        Ptr("a"),
				SubscriptionsURL:  Ptr("a"),
			},
			HTMLURL:          Ptr("a"),
			Name:             Ptr("Hello-World"),
			FullName:         Ptr("Codertocat/Hello-World"),
			Description:      nil,
			Fork:             Ptr(false),
			Homepage:         nil,
			DefaultBranch:    Ptr("main"),
			CreatedAt:        &Timestamp{referenceTime},
			PushedAt:         &Timestamp{referenceTime},
			UpdatedAt:        &Timestamp{referenceTime},
			CloneURL:         Ptr("a"),
			GitURL:           Ptr("a"),
			MirrorURL:        nil,
			SSHURL:           Ptr("a"),
			SVNURL:           Ptr("a"),
			Language:         nil,
			ForksCount:       Ptr(0),
			OpenIssuesCount:  Ptr(2),
			OpenIssues:       Ptr(2),
			StargazersCount:  Ptr(0),
			WatchersCount:    Ptr(0),
			Watchers:         Ptr(0),
			Size:             Ptr(0),
			Archived:         Ptr(false),
			Disabled:         Ptr(false),
			License:          nil,
			Private:          Ptr(false),
			HasIssues:        Ptr(true),
			HasWiki:          Ptr(true),
			HasPages:         Ptr(true),
			HasProjects:      Ptr(true),
			HasDownloads:     Ptr(true),
			URL:              Ptr("a"),
			ArchiveURL:       Ptr("a"),
			AssigneesURL:     Ptr("a"),
			BlobsURL:         Ptr("a"),
			BranchesURL:      Ptr("a"),
			CollaboratorsURL: Ptr("a"),
			CommentsURL:      Ptr("a"),
			CommitsURL:       Ptr("a"),
			CompareURL:       Ptr("a"),
			ContentsURL:      Ptr("a"),
			ContributorsURL:  Ptr("a"),
			DeploymentsURL:   Ptr("a"),
			DownloadsURL:     Ptr("a"),
			EventsURL:        Ptr("a"),
			ForksURL:         Ptr("a"),
			GitCommitsURL:    Ptr("a"),
			GitRefsURL:       Ptr("a"),
			GitTagsURL:       Ptr("a"),
			HooksURL:         Ptr("a"),
			IssueCommentURL:  Ptr("a"),
			IssueEventsURL:   Ptr("a"),
			IssuesURL:        Ptr("a"),
			KeysURL:          Ptr("a"),
			LabelsURL:        Ptr("a"),
			LanguagesURL:     Ptr("a"),
			MergesURL:        Ptr("a"),
			MilestonesURL:    Ptr("a"),
			NotificationsURL: Ptr("a"),
			PullsURL:         Ptr("a"),
			ReleasesURL:      Ptr("a"),
			StargazersURL:    Ptr("a"),
			StatusesURL:      Ptr("a"),
			SubscribersURL:   Ptr("a"),
			SubscriptionURL:  Ptr("a"),
			TagsURL:          Ptr("a"),
			TreesURL:         Ptr("a"),
			TeamsURL:         Ptr("a"),
		},
		Org: &Organization{
			Login:            Ptr("Octocoders"),
			ID:               Ptr(int64(6)),
			NodeID:           Ptr("MDEyOk9yZ2FuaXphdGlvbjY="),
			AvatarURL:        Ptr("a"),
			Description:      Ptr(""),
			URL:              Ptr("a"),
			EventsURL:        Ptr("a"),
			HooksURL:         Ptr("a"),
			IssuesURL:        Ptr("a"),
			MembersURL:       Ptr("a"),
			PublicMembersURL: Ptr("a"),
			ReposURL:         Ptr("a"),
		},
		Sender: &User{
			Login:             Ptr("github"),
			ID:                Ptr(int64(9919)),
			NodeID:            Ptr("MDEyOk9yZ2FuaXphdGlvbjk5MTk="),
			AvatarURL:         Ptr("a"),
			HTMLURL:           Ptr("a"),
			GravatarID:        Ptr(""),
			Type:              Ptr("Organization"),
			SiteAdmin:         Ptr(false),
			URL:               Ptr("a"),
			EventsURL:         Ptr("a"),
			FollowingURL:      Ptr("a"),
			FollowersURL:      Ptr("a"),
			GistsURL:          Ptr("a"),
			OrganizationsURL:  Ptr("a"),
			ReceivedEventsURL: Ptr("a"),
			ReposURL:          Ptr("a"),
			StarredURL:        Ptr("a"),
			SubscriptionsURL:  Ptr("a"),
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
			"name": "Rubocop"
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
		  "size": 0,
		  "stargazers_count": 0,
		  "watchers_count": 0,
		  "has_issues": true,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": true,
		  "forks_count": 0,
		  "archived": false,
		  "disabled": false,
		  "open_issues_count": 2,
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
	t.Parallel()
	testJSONMarshal(t, &SponsorshipEvent{}, "{}")

	u := &SponsorshipEvent{
		Action:        Ptr("created"),
		EffectiveDate: Ptr("2023-01-01T00:00:00Z"),
		Changes: &SponsorshipChanges{
			Tier: &SponsorshipTier{
				From: Ptr("basic"),
			},
			PrivacyLevel: Ptr("public"),
		},
		Repository: &Repository{
			ID:     Ptr(int64(12345)),
			NodeID: Ptr("MDEwOlJlcG9zaXRvcnkxMjM0NQ=="),
			Name:   Ptr("example-repo"),
		},
		Organization: &Organization{
			Login: Ptr("example-org"),
			ID:    Ptr(int64(67890)),
		},
		Sender: &User{
			Login: Ptr("example-user"),
			ID:    Ptr(int64(1111)),
		},
		Installation: &Installation{
			ID: Ptr(int64(2222)),
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
	t.Parallel()
	testJSONMarshal(t, &SponsorshipChanges{}, "{}")

	u := &SponsorshipChanges{
		Tier: &SponsorshipTier{
			From: Ptr("premium"),
		},
		PrivacyLevel: Ptr("private"),
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
	t.Parallel()
	testJSONMarshal(t, &SponsorshipTier{}, "{}")

	u := &SponsorshipTier{
		From: Ptr("gold"),
	}

	want := `{
		"from": "gold"
	}`

	testJSONMarshal(t, u, want)
}
