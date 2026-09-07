// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

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
			"created_at": `+refTimeStr(1136178000)+`,
			"pushed_at": `+refTimeStr(1136178001)+`,
			"updated_at": `+refTimeStr(1136178002)+`,
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
		ID:     new(int64(1)),
		NodeID: new("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             new("owner"),
			ID:                new(int64(2)),
			NodeID:            new("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         new("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           new("https://github.com/xxxxx"),
			GravatarID:        new(""),
			Type:              new("User"),
			SiteAdmin:         new(false),
			URL:               new("https://api.github.com/users/owner"),
			EventsURL:         new("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      new("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      new("https://api.github.com/users/owner/followers"),
			GistsURL:          new("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  new("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: new("https://api.github.com/users/owner/received_events"),
			ReposURL:          new("https://api.github.com/users/owner/repos"),
			StarredURL:        new("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  new("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             new("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         new("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    new("master"),
		CreatedAt:        refTimestamp(1136178000),
		PushedAt:         refTimestamp(1136178001),
		UpdatedAt:        refTimestamp(1136178002),
		HTMLURL:          new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           new("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           new("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             new(false),
		ForksCount:       new(0),
		NetworkCount:     new(0),
		OpenIssuesCount:  new(0),
		OpenIssues:       new(0),
		StargazersCount:  new(0),
		SubscribersCount: new(0),
		WatchersCount:    new(0),
		Watchers:         new(0),
		Size:             new(0),
		Permissions: &RepositoryPermissions{
			Admin:    new(true),
			Maintain: new(true),
			Pull:     new(true),
			Push:     new(true),
			Triage:   new(true),
		},
		AllowForking:             new(true),
		WebCommitSignoffRequired: new(false),
		Archived:                 new(false),
		Disabled:                 new(false),
		Private:                  new(true),
		HasIssues:                new(false),
		HasWiki:                  new(false),
		HasPages:                 new(false),
		HasProjects:              new(false),
		HasDownloads:             new(false),
		HasDiscussions:           new(false),
		IsTemplate:               new(false),
		URL:                      new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               new("private"),
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
			"created_at": `+refTimeStr(1136178000)+`,
			"pushed_at": `+refTimeStr(1136178001)+`,
			"updated_at": `+refTimeStr(1136178002)+`,
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
		ID:     new(int64(1)),
		NodeID: new("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             new("owner"),
			ID:                new(int64(2)),
			NodeID:            new("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         new("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           new("https://github.com/xxxxx"),
			GravatarID:        new(""),
			Type:              new("User"),
			SiteAdmin:         new(false),
			URL:               new("https://api.github.com/users/owner"),
			EventsURL:         new("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      new("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      new("https://api.github.com/users/owner/followers"),
			GistsURL:          new("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  new("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: new("https://api.github.com/users/owner/received_events"),
			ReposURL:          new("https://api.github.com/users/owner/repos"),
			StarredURL:        new("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  new("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             new("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         new("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    new("master"),
		CreatedAt:        refTimestamp(1136178000),
		PushedAt:         refTimestamp(1136178001),
		UpdatedAt:        refTimestamp(1136178002),
		HTMLURL:          new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           new("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           new("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           new("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             new(false),
		ForksCount:       new(0),
		NetworkCount:     new(0),
		OpenIssuesCount:  new(0),
		OpenIssues:       new(0),
		StargazersCount:  new(0),
		SubscribersCount: new(0),
		WatchersCount:    new(0),
		Watchers:         new(0),
		Size:             new(0),
		Permissions: &RepositoryPermissions{
			Admin:    new(true),
			Maintain: new(true),
			Pull:     new(true),
			Push:     new(true),
			Triage:   new(true),
		},
		AllowForking:             new(true),
		WebCommitSignoffRequired: new(false),
		Archived:                 new(false),
		Disabled:                 new(false),
		Private:                  new(true),
		HasIssues:                new(false),
		HasWiki:                  new(false),
		HasPages:                 new(false),
		HasProjects:              new(false),
		HasDownloads:             new(false),
		HasDiscussions:           new(false),
		IsTemplate:               new(false),
		URL:                      new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 new("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               new("private"),
	}
	if !cmp.Equal(fork, want) {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned %+v, want %+v", fork, want)
	}
}

func TestSecurityAdvisoriesService_CreateTemporaryPrivateFork_deferred_badBody(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{invalid json`)
	})

	ctx := t.Context()
	fork, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
	if err == nil {
		t.Fatal("SecurityAdvisories.CreateTemporaryPrivateFork returned nil error")
	}
	if fork != nil {
		t.Errorf("SecurityAdvisories.CreateTemporaryPrivateFork returned non-nil fork: %+v", fork)
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
		testFormValues(t, r, values{"state": "draft"})

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
			GHSAID: new("GHSA-abcd-1234-efgh"),
			CVEID:  new("CVE-2050-00000"),
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
		testFormValues(t, r, values{"state": "draft"})

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
			GHSAID: new("GHSA-abcd-1234-efgh"),
			CVEID:  new("CVE-2050-00000"),
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
				"published_at": `+refTimeStr(1136178000)+`,
				"updated_at": `+refTimeStr(1136178001)+`,
				"github_reviewed_at": `+refTimeStr(1136178002)+`,
				"nvd_published_at": `+refTimeStr(1136178003)+`,
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
				"cvss_severities": {
					"cvss_v3": {"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H", "score": 7.6},
					"cvss_v4": {"vector_string": "CVSS:4.0/AV:N/AC:L/AT:N/PR:H/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N", "score": 8.7}
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
	opts := &ListGlobalSecurityAdvisoriesOptions{CVEID: new("CVE-xoxo-1234")}

	advisories, _, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, opts)
	if err != nil {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories returned error: %v", err)
	}

	want := []*GlobalSecurityAdvisory{
		{
			ID: new(int64(1)),
			SecurityAdvisory: SecurityAdvisory{
				GHSAID:      new("GHSA-xoxo-1234-xoxo"),
				CVEID:       new("CVE-xoxo-1234"),
				URL:         new("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
				HTMLURL:     new("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
				Severity:    new("high"),
				Summary:     new("Heartbleed security advisory"),
				Description: new("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
				Identifiers: []*AdvisoryIdentifier{
					{
						Type:  new("GHSA"),
						Value: new("GHSA-xoxo-1234-xoxo"),
					},
					{
						Type:  new("CVE"),
						Value: new("CVE-xoxo-1234"),
					},
				},
				PublishedAt: refTimestamp(1136178000),
				UpdatedAt:   refTimestamp(1136178001),
				WithdrawnAt: nil,
				CVSS: &AdvisoryCVSS{
					VectorString: new("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
					Score:        new(7.6),
				},
				CVSSSeverities: &AdvisoryCVSSSeverities{
					CVSSV3: &AdvisoryCVSS{
						VectorString: new("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
						Score:        new(7.6),
					},
					CVSSV4: &AdvisoryCVSS{
						VectorString: new("CVSS:4.0/AV:N/AC:L/AT:N/PR:H/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N"),
						Score:        new(8.7),
					},
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: new("CWE-400"),
						Name:  new("Uncontrolled Resource Consumption"),
					},
				},
			},
			References: []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
			Vulnerabilities: []*GlobalSecurityVulnerability{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: new("npm"),
						Name:      new("a-package"),
					},
					FirstPatchedVersion:    new("1.0.3"),
					VulnerableVersionRange: new("<=1.0.2"),
					VulnerableFunctions:    []string{"a_function"},
				},
			},
			RepositoryAdvisoryURL: new("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
			Type:                  new("reviewed"),
			SourceCodeLocation:    new("https://github.com/project/a-package"),
			GithubReviewedAt:      refTimestamp(1136178002),
			NVDPublishedAt:        refTimestamp(1136178003),
			Credits: []*Credit{
				{
					User: &User{
						Login:             new("user"),
						ID:                new(int64(1)),
						NodeID:            new("12="),
						AvatarURL:         new("a"),
						GravatarID:        new(""),
						URL:               new("a"),
						HTMLURL:           new("b"),
						FollowersURL:      new("b"),
						FollowingURL:      new("c"),
						GistsURL:          new("d"),
						StarredURL:        new("e"),
						SubscriptionsURL:  new("f"),
						OrganizationsURL:  new("g"),
						ReposURL:          new("h"),
						EventsURL:         new("i"),
						ReceivedEventsURL: new("j"),
						Type:              new("User"),
						SiteAdmin:         new(false),
					},
					Type: new("analyst"),
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
			"published_at": `+refTimeStr(1136178000)+`,
			"updated_at": `+refTimeStr(1136178001)+`,
			"github_reviewed_at": `+refTimeStr(1136178002)+`,
			"nvd_published_at": `+refTimeStr(1136178003)+`,
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
			"cvss_severities": {
				"cvss_v3": {"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H", "score": 7.6},
				"cvss_v4": {"vector_string": "CVSS:4.0/AV:N/AC:L/AT:N/PR:H/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N", "score": 8.7}
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

	want := &GlobalSecurityAdvisory{
		ID: new(int64(1)),
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      new("GHSA-xoxo-1234-xoxo"),
			CVEID:       new("CVE-xoxo-1234"),
			URL:         new("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     new("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    new("high"),
			Summary:     new("Heartbleed security advisory"),
			Description: new("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  new("GHSA"),
					Value: new("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  new("CVE"),
					Value: new("CVE-xoxo-1234"),
				},
			},
			PublishedAt: refTimestamp(1136178000),
			UpdatedAt:   refTimestamp(1136178001),
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: new("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        new(7.6),
			},
			CVSSSeverities: &AdvisoryCVSSSeverities{
				CVSSV3: &AdvisoryCVSS{
					VectorString: new("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
					Score:        new(7.6),
				},
				CVSSV4: &AdvisoryCVSS{
					VectorString: new("CVSS:4.0/AV:N/AC:L/AT:N/PR:H/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N"),
					Score:        new(8.7),
				},
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: new("CWE-400"),
					Name:  new("Uncontrolled Resource Consumption"),
				},
			},
		},
		RepositoryAdvisoryURL: new("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
		Type:                  new("reviewed"),
		SourceCodeLocation:    new("https://github.com/project/a-package"),
		References:            []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
		GithubReviewedAt:      refTimestamp(1136178002),
		NVDPublishedAt:        refTimestamp(1136178003),

		Vulnerabilities: []*GlobalSecurityVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: new("npm"),
					Name:      new("a-package"),
				},
				FirstPatchedVersion:    new("1.0.3"),
				VulnerableVersionRange: new("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		Credits: []*Credit{
			{
				User: &User{
					Login:             new("user"),
					ID:                new(int64(1)),
					NodeID:            new("12="),
					AvatarURL:         new("a"),
					GravatarID:        new(""),
					URL:               new("a"),
					HTMLURL:           new("b"),
					FollowersURL:      new("b"),
					FollowingURL:      new("c"),
					GistsURL:          new("d"),
					StarredURL:        new("e"),
					SubscriptionsURL:  new("f"),
					OrganizationsURL:  new("g"),
					ReposURL:          new("h"),
					EventsURL:         new("i"),
					ReceivedEventsURL: new("j"),
					Type:              new("User"),
					SiteAdmin:         new(false),
				},
				Type: new("analyst"),
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
