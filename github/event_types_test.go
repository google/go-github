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
