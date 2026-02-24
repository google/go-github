// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSecurityAdvisoriesService_RequestCVE(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_ok/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_accepted/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := t.Context()
	_, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_ok")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	_, err = client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_accepted")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	const methodName = "RequestCVE"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecurityAdvisories.RequestCVE(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id")
		if err == nil {
			t.Errorf("testNewRequestAndDoFailure %v should have return err", methodName)
		}
		return resp, err
	})
}

func TestSecurityAdvisoriesService_CreateTemporaryPrivateFork(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 1,
			"node_id": "R_kgDPP3c6pQ",
			"owner": {
				"login": "owner",
				"id": 2,
				"node_id": "MDQ6VXFGcjYyMjcyMTQw",
				"avatar_url": "https://avatars.githubusercontent.com/u/111111?v=4",
				"html_url": "https://github.com/xxxxx",
				"gravatar_id": "",
				"type": "User",
				"site_admin": false,
				"url": "https://api.github.com/users/owner",
				"events_url": "https://api.github.com/users/owner/events{/privacy}",
				"following_url": "https://api.github.com/users/owner/following{/other_user}",
				"followers_url": "https://api.github.com/users/owner/followers",
				"gists_url": "https://api.github.com/users/owner/gists{/gist_id}",
				"organizations_url": "https://api.github.com/users/owner/orgs",
				"received_events_url": "https://api.github.com/users/owner/received_events",
				"repos_url": "https://api.github.com/users/owner/repos",
				"starred_url": "https://api.github.com/users/owner/starred{/owner}{/repo}",
				"subscriptions_url": "https://api.github.com/users/owner/subscriptions"
			},
			"name": "repo-ghsa-xxxx-xxxx-xxxx",
			"full_name": "owner/repo-ghsa-xxxx-xxxx-xxxx",
			"default_branch": "master",
			"created_at": "2023-12-08T17:22:41Z",
			"pushed_at": "2023-12-03T11:27:08Z",
			"updated_at": "2023-12-08T17:22:42Z",
			"html_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"clone_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"git_url": "git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"ssh_url": "git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"svn_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"fork": false,
			"forks_count": 0,
			"network_count": 0,
			"open_issues_count": 0,
			"open_issues": 0,
			"stargazers_count": 0,
			"subscribers_count": 0,
			"watchers_count": 0,
			"watchers": 0,
			"size": 0,
			"permissions": {
				"admin": true,
				"maintain": true,
				"pull": true,
				"push": true,
				"triage": true
			},
			"allow_forking": true,
			"web_commit_signoff_required": false,
			"archived": false,
			"disabled": false,
			"private": true,
			"has_issues": false,
			"has_wiki": false,
			"has_pages": false,
			"has_projects": false,
			"has_downloads": false,
			"has_discussions": false,
			"is_template": false,
			"url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"archive_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}",
			"assignees_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}",
			"blobs_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}",
			"branches_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}",
			"collaborators_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}",
			"comments_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}",
			"commits_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}",
			"compare_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}",
			"contents_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}",
			"contributors_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors",
			"deployments_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments",
			"downloads_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads",
			"events_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events",
			"forks_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks",
			"git_commits_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}",
			"git_refs_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}",
			"git_tags_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}",
			"hooks_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks",
			"issue_comment_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}",
			"issue_events_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}",
			"issues_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}",
			"keys_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}",
			"labels_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}",
			"languages_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages",
			"merges_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges",
			"milestones_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}",
			"notifications_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}",
			"pulls_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}",
			"releases_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}",
			"stargazers_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers",
			"statuses_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}",
			"subscribers_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers",
			"subscription_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription",
			"tags_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags",
			"teams_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams",
			"visibility": "private"
		}`)
	})

	ctx := t.Context()
	fork, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned error: %v", err)
	}

	want := &Repository{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             Ptr("owner"),
			ID:                Ptr(int64(2)),
			NodeID:            Ptr("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         Ptr("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           Ptr("https://github.com/xxxxx"),
			GravatarID:        Ptr(""),
			Type:              Ptr("User"),
			SiteAdmin:         Ptr(false),
			URL:               Ptr("https://api.github.com/users/owner"),
			EventsURL:         Ptr("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      Ptr("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      Ptr("https://api.github.com/users/owner/followers"),
			GistsURL:          Ptr("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  Ptr("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: Ptr("https://api.github.com/users/owner/received_events"),
			ReposURL:          Ptr("https://api.github.com/users/owner/repos"),
			StarredURL:        Ptr("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  Ptr("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             Ptr("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         Ptr("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    Ptr("master"),
		CreatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 41, 0, time.UTC)},
		PushedAt:         &Timestamp{time.Date(2023, time.December, 3, 11, 27, 8, 0, time.UTC)},
		UpdatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 42, 0, time.UTC)},
		HTMLURL:          Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           Ptr("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           Ptr("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             Ptr(false),
		ForksCount:       Ptr(0),
		NetworkCount:     Ptr(0),
		OpenIssuesCount:  Ptr(0),
		OpenIssues:       Ptr(0),
		StargazersCount:  Ptr(0),
		SubscribersCount: Ptr(0),
		WatchersCount:    Ptr(0),
		Watchers:         Ptr(0),
		Size:             Ptr(0),
		Permissions: &RepositoryPermissions{
			Admin:    Ptr(true),
			Maintain: Ptr(true),
			Pull:     Ptr(true),
			Push:     Ptr(true),
			Triage:   Ptr(true),
		},
		AllowForking:             Ptr(true),
		WebCommitSignoffRequired: Ptr(false),
		Archived:                 Ptr(false),
		Disabled:                 Ptr(false),
		Private:                  Ptr(true),
		HasIssues:                Ptr(false),
		HasWiki:                  Ptr(false),
		HasPages:                 Ptr(false),
		HasProjects:              Ptr(false),
		HasDownloads:             Ptr(false),
		HasDiscussions:           Ptr(false),
		IsTemplate:               Ptr(false),
		URL:                      Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               Ptr("private"),
	}
	if !cmp.Equal(fork, want) {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned %+v, want %+v", fork, want)
	}

	const methodName = "CreateTemporaryPrivateFork"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSecurityAdvisoriesService_CreateTemporaryPrivateFork_deferred(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{
			"id": 1,
			"node_id": "R_kgDPP3c6pQ",
			"owner": {
				"login": "owner",
				"id": 2,
				"node_id": "MDQ6VXFGcjYyMjcyMTQw",
				"avatar_url": "https://avatars.githubusercontent.com/u/111111?v=4",
				"html_url": "https://github.com/xxxxx",
				"gravatar_id": "",
				"type": "User",
				"site_admin": false,
				"url": "https://api.github.com/users/owner",
				"events_url": "https://api.github.com/users/owner/events{/privacy}",
				"following_url": "https://api.github.com/users/owner/following{/other_user}",
				"followers_url": "https://api.github.com/users/owner/followers",
				"gists_url": "https://api.github.com/users/owner/gists{/gist_id}",
				"organizations_url": "https://api.github.com/users/owner/orgs",
				"received_events_url": "https://api.github.com/users/owner/received_events",
				"repos_url": "https://api.github.com/users/owner/repos",
				"starred_url": "https://api.github.com/users/owner/starred{/owner}{/repo}",
				"subscriptions_url": "https://api.github.com/users/owner/subscriptions"
			},
			"name": "repo-ghsa-xxxx-xxxx-xxxx",
			"full_name": "owner/repo-ghsa-xxxx-xxxx-xxxx",
			"default_branch": "master",
			"created_at": "2023-12-08T17:22:41Z",
			"pushed_at": "2023-12-03T11:27:08Z",
			"updated_at": "2023-12-08T17:22:42Z",
			"html_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"clone_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"git_url": "git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"ssh_url": "git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git",
			"svn_url": "https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"fork": false,
			"forks_count": 0,
			"network_count": 0,
			"open_issues_count": 0,
			"open_issues": 0,
			"stargazers_count": 0,
			"subscribers_count": 0,
			"watchers_count": 0,
			"watchers": 0,
			"size": 0,
			"permissions": {
				"admin": true,
				"maintain": true,
				"pull": true,
				"push": true,
				"triage": true
			},
			"allow_forking": true,
			"web_commit_signoff_required": false,
			"archived": false,
			"disabled": false,
			"private": true,
			"has_issues": false,
			"has_wiki": false,
			"has_pages": false,
			"has_projects": false,
			"has_downloads": false,
			"has_discussions": false,
			"is_template": false,
			"url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx",
			"archive_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}",
			"assignees_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}",
			"blobs_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}",
			"branches_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}",
			"collaborators_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}",
			"comments_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}",
			"commits_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}",
			"compare_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}",
			"contents_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}",
			"contributors_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors",
			"deployments_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments",
			"downloads_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads",
			"events_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events",
			"forks_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks",
			"git_commits_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}",
			"git_refs_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}",
			"git_tags_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}",
			"hooks_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks",
			"issue_comment_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}",
			"issue_events_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}",
			"issues_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}",
			"keys_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}",
			"labels_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}",
			"languages_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages",
			"merges_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges",
			"milestones_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}",
			"notifications_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}",
			"pulls_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}",
			"releases_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}",
			"stargazers_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers",
			"statuses_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}",
			"subscribers_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers",
			"subscription_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription",
			"tags_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags",
			"teams_url": "https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams",
			"visibility": "private"
		}`)
	})

	ctx := t.Context()
	fork, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
	if !errors.As(err, new(*AcceptedError)) {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned error: %v (want AcceptedError)", err)
	}

	want := &Repository{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             Ptr("owner"),
			ID:                Ptr(int64(2)),
			NodeID:            Ptr("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         Ptr("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           Ptr("https://github.com/xxxxx"),
			GravatarID:        Ptr(""),
			Type:              Ptr("User"),
			SiteAdmin:         Ptr(false),
			URL:               Ptr("https://api.github.com/users/owner"),
			EventsURL:         Ptr("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      Ptr("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      Ptr("https://api.github.com/users/owner/followers"),
			GistsURL:          Ptr("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  Ptr("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: Ptr("https://api.github.com/users/owner/received_events"),
			ReposURL:          Ptr("https://api.github.com/users/owner/repos"),
			StarredURL:        Ptr("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  Ptr("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             Ptr("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         Ptr("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    Ptr("master"),
		CreatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 41, 0, time.UTC)},
		PushedAt:         &Timestamp{time.Date(2023, time.December, 3, 11, 27, 8, 0, time.UTC)},
		UpdatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 42, 0, time.UTC)},
		HTMLURL:          Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           Ptr("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           Ptr("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           Ptr("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             Ptr(false),
		ForksCount:       Ptr(0),
		NetworkCount:     Ptr(0),
		OpenIssuesCount:  Ptr(0),
		OpenIssues:       Ptr(0),
		StargazersCount:  Ptr(0),
		SubscribersCount: Ptr(0),
		WatchersCount:    Ptr(0),
		Watchers:         Ptr(0),
		Size:             Ptr(0),
		Permissions: &RepositoryPermissions{
			Admin:    Ptr(true),
			Maintain: Ptr(true),
			Pull:     Ptr(true),
			Push:     Ptr(true),
			Triage:   Ptr(true),
		},
		AllowForking:             Ptr(true),
		WebCommitSignoffRequired: Ptr(false),
		Archived:                 Ptr(false),
		Disabled:                 Ptr(false),
		Private:                  Ptr(true),
		HasIssues:                Ptr(false),
		HasWiki:                  Ptr(false),
		HasPages:                 Ptr(false),
		HasProjects:              Ptr(false),
		HasDownloads:             Ptr(false),
		HasDiscussions:           Ptr(false),
		IsTemplate:               Ptr(false),
		URL:                      Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 Ptr("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               Ptr("private"),
	}
	if !cmp.Equal(fork, want) {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned %+v, want %+v", fork, want)
	}
}

func TestSecurityAdvisoriesService_CreateTemporaryPrivateFork_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "%", "r", "ghsa_id")
	testURLParseError(t, err)
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_BadRequest(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Error("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_NotFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Error("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_UnmarshalError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Error("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal") || !strings.Contains(err.Error(), "ghsa_id") {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[
			{
				"ghsa_id": "GHSA-abcd-1234-efgh",
				"cve_id": "CVE-2050-00000"
			}
		]`))
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %v, want %v", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: Ptr("GHSA-abcd-1234-efgh"),
			CVEID:  Ptr("CVE-2050-00000"),
		},
	}
	if !cmp.Equal(advisories, want) {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned %+v, want %+v", advisories, want)
	}

	methodName := "ListRepositorySecurityAdvisoriesForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "\n", &ListRepositorySecurityAdvisoriesOptions{
			Sort: "\n",
		})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_BadRequest(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Error("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_NotFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisories returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Error("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_UnmarshalError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Error("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal") || !strings.Contains(err.Error(), "ghsa_id") {
		t.Errorf("ListRepositorySecurityAdvisories returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %v, want %v", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[
			{
				"ghsa_id": "GHSA-abcd-1234-efgh",
				"cve_id": "CVE-2050-00000"
			}
		]`))
	})

	ctx := t.Context()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisories returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %v, want %v", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: Ptr("GHSA-abcd-1234-efgh"),
			CVEID:  Ptr("CVE-2050-00000"),
		},
	}
	if !cmp.Equal(advisories, want) {
		t.Errorf("ListRepositorySecurityAdvisories returned %+v, want %+v", advisories, want)
	}

	methodName := "ListRepositorySecurityAdvisories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "\n", "\n", &ListRepositorySecurityAdvisoriesOptions{
			Sort: "\n",
		})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestListGlobalSecurityAdvisories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"cve_id": "CVE-xoxo-1234"})

		fmt.Fprint(w, `[{
				"id": 1,
				"ghsa_id": "GHSA-xoxo-1234-xoxo",
				"cve_id": "CVE-xoxo-1234",
				"url": "https://api.github.com/advisories/GHSA-xoxo-1234-xoxo",
				"html_url": "https://github.com/advisories/GHSA-xoxo-1234-xoxo",
				"repository_advisory_url": "https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo",
				"summary": "Heartbleed security advisory",
				"description": "This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information.",
				"type": "reviewed",
				"severity": "high",
				"source_code_location": "https://github.com/project/a-package",
				"identifiers": [
					{
						"type": "GHSA",
						"value": "GHSA-xoxo-1234-xoxo"
					},
					{
						"type": "CVE",
						"value": "CVE-xoxo-1234"
					}
				],
				"references": ["https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"],
				"published_at": "1996-06-20T00:00:00Z",
				"updated_at": "1996-06-20T00:00:00Z",
				"github_reviewed_at": "1996-06-20T00:00:00Z",
				"nvd_published_at": "1996-06-20T00:00:00Z",
				"withdrawn_at": null,
				"vulnerabilities": [
					{
						"package": {
							"ecosystem": "npm",
							"name": "a-package"
						},
						"first_patched_version": "1.0.3",
						"vulnerable_version_range": "<=1.0.2",
						"vulnerable_functions": ["a_function"]
					}
				],
				"cvss": {
					"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H",
					"score": 7.6
				},
				"cwes": [
					{
						"cwe_id": "CWE-400",
						"name": "Uncontrolled Resource Consumption"
					}
				],
				"credits": [
					{
						"user": {
							"login": "user",
							"id": 1,
							"node_id": "12=",
							"avatar_url": "a",
							"gravatar_id": "",
							"url": "a",
							"html_url": "b",
							"followers_url": "b",
							"following_url": "c",
							"gists_url": "d",
							"starred_url": "e",
							"subscriptions_url": "f",
							"organizations_url": "g",
							"repos_url": "h",
							"events_url": "i",
							"received_events_url": "j",
							"type": "User",
							"site_admin": false
						},
						"type": "analyst"
					}
				]
			}
		]`)
	})

	ctx := t.Context()
	opts := &ListGlobalSecurityAdvisoriesOptions{CVEID: Ptr("CVE-xoxo-1234")}

	advisories, _, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, opts)
	if err != nil {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 0, 0, 0, 0, time.UTC)}
	want := []*GlobalSecurityAdvisory{
		{
			ID: Ptr(int64(1)),
			SecurityAdvisory: SecurityAdvisory{
				GHSAID:      Ptr("GHSA-xoxo-1234-xoxo"),
				CVEID:       Ptr("CVE-xoxo-1234"),
				URL:         Ptr("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
				HTMLURL:     Ptr("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
				Severity:    Ptr("high"),
				Summary:     Ptr("Heartbleed security advisory"),
				Description: Ptr("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
				Identifiers: []*AdvisoryIdentifier{
					{
						Type:  Ptr("GHSA"),
						Value: Ptr("GHSA-xoxo-1234-xoxo"),
					},
					{
						Type:  Ptr("CVE"),
						Value: Ptr("CVE-xoxo-1234"),
					},
				},
				PublishedAt: &date,
				UpdatedAt:   &date,
				WithdrawnAt: nil,
				CVSS: &AdvisoryCVSS{
					VectorString: Ptr("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
					Score:        Ptr(7.6),
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: Ptr("CWE-400"),
						Name:  Ptr("Uncontrolled Resource Consumption"),
					},
				},
			},
			References: []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
			Vulnerabilities: []*GlobalSecurityVulnerability{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: Ptr("npm"),
						Name:      Ptr("a-package"),
					},
					FirstPatchedVersion:    Ptr("1.0.3"),
					VulnerableVersionRange: Ptr("<=1.0.2"),
					VulnerableFunctions:    []string{"a_function"},
				},
			},
			RepositoryAdvisoryURL: Ptr("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
			Type:                  Ptr("reviewed"),
			SourceCodeLocation:    Ptr("https://github.com/project/a-package"),
			GithubReviewedAt:      &date,
			NVDPublishedAt:        &date,
			Credits: []*Credit{
				{
					User: &User{
						Login:             Ptr("user"),
						ID:                Ptr(int64(1)),
						NodeID:            Ptr("12="),
						AvatarURL:         Ptr("a"),
						GravatarID:        Ptr(""),
						URL:               Ptr("a"),
						HTMLURL:           Ptr("b"),
						FollowersURL:      Ptr("b"),
						FollowingURL:      Ptr("c"),
						GistsURL:          Ptr("d"),
						StarredURL:        Ptr("e"),
						SubscriptionsURL:  Ptr("f"),
						OrganizationsURL:  Ptr("g"),
						ReposURL:          Ptr("h"),
						EventsURL:         Ptr("i"),
						ReceivedEventsURL: Ptr("j"),
						Type:              Ptr("User"),
						SiteAdmin:         Ptr(false),
					},
					Type: Ptr("analyst"),
				},
			},
		},
	}

	if !cmp.Equal(advisories, want) {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories %+v, want %+v", advisories, want)
	}

	const methodName = "ListGlobalSecurityAdvisories"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, nil)
		return resp, err
	})
}

func TestGetGlobalSecurityAdvisories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/advisories/GHSA-xoxo-1234-xoxo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"id": 1,
			"ghsa_id": "GHSA-xoxo-1234-xoxo",
			"cve_id": "CVE-xoxo-1234",
			"url": "https://api.github.com/advisories/GHSA-xoxo-1234-xoxo",
			"html_url": "https://github.com/advisories/GHSA-xoxo-1234-xoxo",
			"repository_advisory_url": "https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo",
			"summary": "Heartbleed security advisory",
			"description": "This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information.",
			"type": "reviewed",
			"severity": "high",
			"source_code_location": "https://github.com/project/a-package",
			"identifiers": [
				{
					"type": "GHSA",
					"value": "GHSA-xoxo-1234-xoxo"
				},
				{
					"type": "CVE",
					"value": "CVE-xoxo-1234"
				}
			],
			"references": ["https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"],
			"published_at": "1996-06-20T00:00:00Z",
			"updated_at": "1996-06-20T00:00:00Z",
			"github_reviewed_at": "1996-06-20T00:00:00Z",
			"nvd_published_at": "1996-06-20T00:00:00Z",
			"withdrawn_at": null,
			"vulnerabilities": [
				{
					"package": {
						"ecosystem": "npm",
						"name": "a-package"
					},
					"first_patched_version": "1.0.3",
					"vulnerable_version_range": "<=1.0.2",
					"vulnerable_functions": ["a_function"]
				}
			],
			"cvss": {
				"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H",
				"score": 7.6
			},
			"cwes": [
				{
					"cwe_id": "CWE-400",
					"name": "Uncontrolled Resource Consumption"
				}
			],
			"credits": [
				{
					"user": {
						"login": "user",
						"id": 1,
						"node_id": "12=",
						"avatar_url": "a",
						"gravatar_id": "",
						"url": "a",
						"html_url": "b",
						"followers_url": "b",
						"following_url": "c",
						"gists_url": "d",
						"starred_url": "e",
						"subscriptions_url": "f",
						"organizations_url": "g",
						"repos_url": "h",
						"events_url": "i",
						"received_events_url": "j",
						"type": "User",
						"site_admin": false
					},
					"type": "analyst"
				}
			]
		}`)
	})

	ctx := t.Context()
	advisory, _, err := client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "GHSA-xoxo-1234-xoxo")
	if err != nil {
		t.Errorf("SecurityAdvisories.GetGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 0, 0, 0, 0, time.UTC)}
	want := &GlobalSecurityAdvisory{
		ID: Ptr(int64(1)),
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      Ptr("GHSA-xoxo-1234-xoxo"),
			CVEID:       Ptr("CVE-xoxo-1234"),
			URL:         Ptr("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     Ptr("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    Ptr("high"),
			Summary:     Ptr("Heartbleed security advisory"),
			Description: Ptr("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  Ptr("GHSA"),
					Value: Ptr("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  Ptr("CVE"),
					Value: Ptr("CVE-xoxo-1234"),
				},
			},
			PublishedAt: &date,
			UpdatedAt:   &date,
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: Ptr("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        Ptr(7.6),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: Ptr("CWE-400"),
					Name:  Ptr("Uncontrolled Resource Consumption"),
				},
			},
		},
		RepositoryAdvisoryURL: Ptr("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
		Type:                  Ptr("reviewed"),
		SourceCodeLocation:    Ptr("https://github.com/project/a-package"),
		References:            []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
		GithubReviewedAt:      &date,
		NVDPublishedAt:        &date,

		Vulnerabilities: []*GlobalSecurityVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("npm"),
					Name:      Ptr("a-package"),
				},
				FirstPatchedVersion:    Ptr("1.0.3"),
				VulnerableVersionRange: Ptr("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		Credits: []*Credit{
			{
				User: &User{
					Login:             Ptr("user"),
					ID:                Ptr(int64(1)),
					NodeID:            Ptr("12="),
					AvatarURL:         Ptr("a"),
					GravatarID:        Ptr(""),
					URL:               Ptr("a"),
					HTMLURL:           Ptr("b"),
					FollowersURL:      Ptr("b"),
					FollowingURL:      Ptr("c"),
					GistsURL:          Ptr("d"),
					StarredURL:        Ptr("e"),
					SubscriptionsURL:  Ptr("f"),
					OrganizationsURL:  Ptr("g"),
					ReposURL:          Ptr("h"),
					EventsURL:         Ptr("i"),
					ReceivedEventsURL: Ptr("j"),
					Type:              Ptr("User"),
					SiteAdmin:         Ptr(false),
				},
				Type: Ptr("analyst"),
			},
		},
	}

	if !cmp.Equal(advisory, want) {
		t.Errorf("SecurityAdvisories.GetGlobalSecurityAdvisories %+v, want %+v", advisory, want)
	}

	const methodName = "GetGlobalSecurityAdvisories"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "CVE-\n-1234")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSecurityAdvisorySubmission_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecurityAdvisorySubmission{}, `{}`)

	u := &SecurityAdvisorySubmission{
		Accepted: Ptr(true),
	}

	w := `{
		"accepted": true
	}`

	testJSONMarshal(t, u, w)
}

func TestRepoAdvisoryCredit_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepoAdvisoryCredit{}, `{}`)

	u := &RepoAdvisoryCredit{
		Login: Ptr("l"),
		Type:  Ptr("t"),
	}

	w := `{
		"login": "l",
		"type": "t"
	}`

	testJSONMarshal(t, u, w)
}

func TestRepoAdvisoryCreditDetailed_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepoAdvisoryCreditDetailed{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &RepoAdvisoryCreditDetailed{
		Type:  Ptr("t"),
		State: Ptr("s"),
		User: &User{
			Name:                    Ptr("u"),
			Company:                 Ptr("c"),
			Blog:                    Ptr("b"),
			Location:                Ptr("l"),
			Email:                   Ptr("e"),
			Hireable:                Ptr(false),
			Bio:                     Ptr("bio"),
			TwitterUsername:         Ptr("tu"),
			PublicRepos:             Ptr(1),
			PublicGists:             Ptr(1),
			Followers:               Ptr(2),
			Following:               Ptr(2),
			CreatedAt:               testDate,
			UpdatedAt:               testDate,
			SuspendedAt:             testDate,
			Type:                    Ptr("type"),
			SiteAdmin:               Ptr(false),
			TotalPrivateRepos:       Ptr(int64(10)),
			OwnedPrivateRepos:       Ptr(int64(10)),
			PrivateGists:            Ptr(10),
			DiskUsage:               Ptr(10),
			Collaborators:           Ptr(10),
			TwoFactorAuthentication: Ptr(true),
			Plan: &Plan{
				Name:          Ptr("p"),
				Space:         Ptr(2),
				Collaborators: Ptr(2),
				PrivateRepos:  Ptr(int64(2)),
				Seats:         Ptr(2),
				FilledSeats:   Ptr(1),
			},
			LdapDn:            Ptr("l"),
			URL:               Ptr("url"),
			EventsURL:         Ptr("e"),
			FollowingURL:      Ptr("f"),
			FollowersURL:      Ptr("f"),
			GistsURL:          Ptr("g"),
			OrganizationsURL:  Ptr("o"),
			ReceivedEventsURL: Ptr("r"),
			ReposURL:          Ptr("rep"),
			StarredURL:        Ptr("star"),
			SubscriptionsURL:  Ptr("sub"),
			TextMatches: []*TextMatch{
				{
					ObjectURL:  Ptr("u"),
					ObjectType: Ptr("t"),
					Property:   Ptr("p"),
					Fragment:   Ptr("f"),
					Matches: []*Match{
						{
							Text:    Ptr("t"),
							Indices: []int{1, 2},
						},
					},
				},
			},
			Permissions: &RepositoryPermissions{Pull: Ptr(true)},
			RoleName:    Ptr("r"),
		},
	}

	w := `{
		"type": "t",
		"state": "s",
		"user": {
		  "name": "u",
		  "company": "c",
		  "blog": "b",
		  "location": "l",
		  "email": "e",
		  "hireable": false,
		  "bio": "bio",
		  "twitter_username": "tu",
		  "public_repos": 1,
		  "public_gists": 1,
		  "followers": 2,
		  "following": 2,
		  "created_at": "2019-08-10T14:59:22Z",
		  "updated_at": "2019-08-10T14:59:22Z",
		  "suspended_at": "2019-08-10T14:59:22Z",
		  "type": "type",
		  "site_admin": false,
		  "total_private_repos": 10,
		  "owned_private_repos": 10,
		  "private_gists": 10,
		  "disk_usage": 10,
		  "collaborators": 10,
		  "two_factor_authentication": true,
		  "plan": {
			"name": "p",
			"space": 2,
			"collaborators": 2,
			"private_repos": 2,
			"seats": 2,
			"filled_seats": 1
		  },
		  "ldap_dn": "l",
		  "url": "url",
		  "events_url": "e",
		  "following_url": "f",
		  "followers_url": "f",
		  "gists_url": "g",
		  "organizations_url": "o",
		  "received_events_url": "r",
		  "repos_url": "rep",
		  "starred_url": "star",
		  "subscriptions_url": "sub",
		  "text_matches": [
			{
			  "object_url": "u",
			  "object_type": "t",
			  "property": "p",
			  "fragment": "f",
			  "matches": [
				{
				  "text": "t",
				  "indices": [1, 2]
				}
			  ]
			}
		  ],
		  "permissions": {
			"pull": true
		  },
		  "role_name": "r"
		}
	  }`

	testJSONMarshal(t, u, w)
}

func TestCredit_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Credit{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &Credit{
		Type: Ptr("t"),
		User: &User{
			Name:                    Ptr("u"),
			Company:                 Ptr("c"),
			Blog:                    Ptr("b"),
			Location:                Ptr("l"),
			Email:                   Ptr("e"),
			Hireable:                Ptr(false),
			Bio:                     Ptr("bio"),
			TwitterUsername:         Ptr("tu"),
			PublicRepos:             Ptr(1),
			PublicGists:             Ptr(1),
			Followers:               Ptr(2),
			Following:               Ptr(2),
			CreatedAt:               testDate,
			UpdatedAt:               testDate,
			SuspendedAt:             testDate,
			Type:                    Ptr("type"),
			SiteAdmin:               Ptr(false),
			TotalPrivateRepos:       Ptr(int64(10)),
			OwnedPrivateRepos:       Ptr(int64(10)),
			PrivateGists:            Ptr(10),
			DiskUsage:               Ptr(10),
			Collaborators:           Ptr(10),
			TwoFactorAuthentication: Ptr(true),
			Plan: &Plan{
				Name:          Ptr("p"),
				Space:         Ptr(2),
				Collaborators: Ptr(2),
				PrivateRepos:  Ptr(int64(2)),
				Seats:         Ptr(2),
				FilledSeats:   Ptr(1),
			},
			LdapDn:            Ptr("l"),
			URL:               Ptr("url"),
			EventsURL:         Ptr("e"),
			FollowingURL:      Ptr("f"),
			FollowersURL:      Ptr("f"),
			GistsURL:          Ptr("g"),
			OrganizationsURL:  Ptr("o"),
			ReceivedEventsURL: Ptr("r"),
			ReposURL:          Ptr("rep"),
			StarredURL:        Ptr("star"),
			SubscriptionsURL:  Ptr("sub"),
			TextMatches: []*TextMatch{
				{
					ObjectURL:  Ptr("u"),
					ObjectType: Ptr("t"),
					Property:   Ptr("p"),
					Fragment:   Ptr("f"),
					Matches: []*Match{
						{
							Text:    Ptr("t"),
							Indices: []int{1, 2},
						},
					},
				},
			},
			Permissions: &RepositoryPermissions{Pull: Ptr(true)},
			RoleName:    Ptr("r"),
		},
	}

	w := `{
		"type": "t",
		"user": {
			"name": "u",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": false,
			"bio": "bio",
			"twitter_username": "tu",
			"public_repos": 1,
			"public_gists": 1,
			"followers": 2,
			"following": 2,
			"created_at": "2019-08-10T14:59:22Z",
			"updated_at": "2019-08-10T14:59:22Z",
			"suspended_at": "2019-08-10T14:59:22Z",
			"type": "type",
			"site_admin": false,
			"total_private_repos": 10,
			"owned_private_repos": 10,
			"private_gists": 10,
			"disk_usage": 10,
			"collaborators": 10,
			"two_factor_authentication": true,
			"plan": {
			"name": "p",
			"space": 2,
			"collaborators": 2,
			"private_repos": 2,
			"seats": 2,
			"filled_seats": 1
			},
			"ldap_dn": "l",
			"url": "url",
			"events_url": "e",
			"following_url": "f",
			"followers_url": "f",
			"gists_url": "g",
			"organizations_url": "o",
			"received_events_url": "r",
			"repos_url": "rep",
			"starred_url": "star",
			"subscriptions_url": "sub",
			"text_matches": [
			{
				"object_url": "u",
				"object_type": "t",
				"property": "p",
				"fragment": "f",
				"matches": [
				{
					"text": "t",
					"indices": [1, 2]
				}
				]
			}
			],
			"permissions": {
			"pull": true
			},
			"role_name": "r"
		}
	}`

	testJSONMarshal(t, u, w)
}

func TestGlobalSecurityAdvisory_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GlobalSecurityAdvisory{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &GlobalSecurityAdvisory{
		ID:                    Ptr(int64(1)),
		RepositoryAdvisoryURL: Ptr("r"),
		Type:                  Ptr("t"),
		SourceCodeLocation:    Ptr("s"),
		References:            []string{"r"},
		Vulnerabilities: []*GlobalSecurityVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("npm"),
					Name:      Ptr("a-package"),
				},
				FirstPatchedVersion:    Ptr("1.0.3"),
				VulnerableVersionRange: Ptr("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		GithubReviewedAt: testDate,
		NVDPublishedAt:   testDate,
		Credits: []*Credit{
			{
				Type: Ptr("t"),
				User: &User{
					Name:                    Ptr("u"),
					Company:                 Ptr("c"),
					Blog:                    Ptr("b"),
					Location:                Ptr("l"),
					Email:                   Ptr("e"),
					Hireable:                Ptr(false),
					Bio:                     Ptr("bio"),
					TwitterUsername:         Ptr("tu"),
					PublicRepos:             Ptr(1),
					PublicGists:             Ptr(1),
					Followers:               Ptr(2),
					Following:               Ptr(2),
					CreatedAt:               testDate,
					UpdatedAt:               testDate,
					SuspendedAt:             testDate,
					Type:                    Ptr("type"),
					SiteAdmin:               Ptr(false),
					TotalPrivateRepos:       Ptr(int64(10)),
					OwnedPrivateRepos:       Ptr(int64(10)),
					PrivateGists:            Ptr(10),
					DiskUsage:               Ptr(10),
					Collaborators:           Ptr(10),
					TwoFactorAuthentication: Ptr(true),
					Plan: &Plan{
						Name:          Ptr("p"),
						Space:         Ptr(2),
						Collaborators: Ptr(2),
						PrivateRepos:  Ptr(int64(2)),
						Seats:         Ptr(2),
						FilledSeats:   Ptr(1),
					},
					LdapDn:            Ptr("l"),
					URL:               Ptr("url"),
					EventsURL:         Ptr("e"),
					FollowingURL:      Ptr("f"),
					FollowersURL:      Ptr("f"),
					GistsURL:          Ptr("g"),
					OrganizationsURL:  Ptr("o"),
					ReceivedEventsURL: Ptr("r"),
					ReposURL:          Ptr("rep"),
					StarredURL:        Ptr("star"),
					SubscriptionsURL:  Ptr("sub"),
					TextMatches: []*TextMatch{
						{
							ObjectURL:  Ptr("u"),
							ObjectType: Ptr("t"),
							Property:   Ptr("p"),
							Fragment:   Ptr("f"),
							Matches: []*Match{
								{
									Text:    Ptr("t"),
									Indices: []int{1, 2},
								},
							},
						},
					},
					Permissions: &RepositoryPermissions{Pull: Ptr(true)},
					RoleName:    Ptr("r"),
				},
			},
		},
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      Ptr("GHSA-xoxo-1234-xoxo"),
			CVEID:       Ptr("CVE-xoxo-1234"),
			URL:         Ptr("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     Ptr("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    Ptr("high"),
			Summary:     Ptr("Heartbleed security advisory"),
			Description: Ptr("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  Ptr("GHSA"),
					Value: Ptr("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  Ptr("CVE"),
					Value: Ptr("CVE-xoxo-1234"),
				},
			},
			PublishedAt: testDate,
			UpdatedAt:   testDate,
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: Ptr("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        Ptr(7.6),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: Ptr("CWE-400"),
					Name:  Ptr("Uncontrolled Resource Consumption"),
				},
			},
		},
	}

	w := `{
		"cvss": {
		  "score": 7.6,
		  "vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"
		},
		"cwes": [
		  {
			"cwe_id": "CWE-400",
			"name": "Uncontrolled Resource Consumption"
		  }
		],
		"ghsa_id": "GHSA-xoxo-1234-xoxo",
		"summary": "Heartbleed security advisory",
		"description": "This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information.",
		"severity": "high",
		"identifiers": [
		  {
			"value": "GHSA-xoxo-1234-xoxo",
			"type": "GHSA"
		  },
		  {
			"value": "CVE-xoxo-1234",
			"type": "CVE"
		  }
		],
		"published_at": "2019-08-10T14:59:22Z",
		"updated_at": "2019-08-10T14:59:22Z",
		"cve_id": "CVE-xoxo-1234",
		"url": "https://api.github.com/advisories/GHSA-xoxo-1234-xoxo",
		"html_url": "https://github.com/advisories/GHSA-xoxo-1234-xoxo",
		"id": 1,
		"repository_advisory_url": "r",
		"type": "t",
		"source_code_location": "s",
		"references": [
		  "r"
		],
		"vulnerabilities": [
		  {
			"package": {
			  "ecosystem": "npm",
			  "name": "a-package"
			},
			"first_patched_version": "1.0.3",
			"vulnerable_version_range": "\u003c=1.0.2",
			"vulnerable_functions": [
			  "a_function"
			]
		  }
		],
		"github_reviewed_at": "2019-08-10T14:59:22Z",
		"nvd_published_at": "2019-08-10T14:59:22Z",
		"credits": [
		  {
			"user": {
			  "name": "u",
			  "company": "c",
			  "blog": "b",
			  "location": "l",
			  "email": "e",
			  "hireable": false,
			  "bio": "bio",
			  "twitter_username": "tu",
			  "public_repos": 1,
			  "public_gists": 1,
			  "followers": 2,
			  "following": 2,
			  "created_at": "2019-08-10T14:59:22Z",
			  "updated_at": "2019-08-10T14:59:22Z",
			  "suspended_at": "2019-08-10T14:59:22Z",
			  "type": "type",
			  "site_admin": false,
			  "total_private_repos": 10,
			  "owned_private_repos": 10,
			  "private_gists": 10,
			  "disk_usage": 10,
			  "collaborators": 10,
			  "two_factor_authentication": true,
			  "plan": {
				"name": "p",
				"space": 2,
				"collaborators": 2,
				"private_repos": 2,
				"filled_seats": 1,
				"seats": 2
			  },
			  "ldap_dn": "l",
			  "url": "url",
			  "events_url": "e",
			  "following_url": "f",
			  "followers_url": "f",
			  "gists_url": "g",
			  "organizations_url": "o",
			  "received_events_url": "r",
			  "repos_url": "rep",
			  "starred_url": "star",
			  "subscriptions_url": "sub",
			  "text_matches": [
				{
				  "object_url": "u",
				  "object_type": "t",
				  "property": "p",
				  "fragment": "f",
				  "matches": [
					{
					  "text": "t",
					  "indices": [
						1,
						2
					  ]
					}
				  ]
				}
			  ],
			  "permissions": {
				"pull": true
			  },
			  "role_name": "r"
			},
			"type": "t"
		  }
		]
	}`

	testJSONMarshal(t, u, w)
}
