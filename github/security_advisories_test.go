// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSecurityAdvisoriesService_RequestCVE(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_ok/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_accepted/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	fork, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned error: %v", err)
	}

	want := &Repository{
		ID:     Int64(1),
		NodeID: String("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             String("owner"),
			ID:                Int64(2),
			NodeID:            String("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         String("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           String("https://github.com/xxxxx"),
			GravatarID:        String(""),
			Type:              String("User"),
			SiteAdmin:         Bool(false),
			URL:               String("https://api.github.com/users/owner"),
			EventsURL:         String("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      String("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      String("https://api.github.com/users/owner/followers"),
			GistsURL:          String("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  String("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: String("https://api.github.com/users/owner/received_events"),
			ReposURL:          String("https://api.github.com/users/owner/repos"),
			StarredURL:        String("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  String("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             String("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         String("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    String("master"),
		CreatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 41, 0, time.UTC)},
		PushedAt:         &Timestamp{time.Date(2023, time.December, 3, 11, 27, 8, 0, time.UTC)},
		UpdatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 42, 0, time.UTC)},
		HTMLURL:          String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           String("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           String("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             Bool(false),
		ForksCount:       Int(0),
		NetworkCount:     Int(0),
		OpenIssuesCount:  Int(0),
		OpenIssues:       Int(0),
		StargazersCount:  Int(0),
		SubscribersCount: Int(0),
		WatchersCount:    Int(0),
		Watchers:         Int(0),
		Size:             Int(0),
		Permissions: map[string]bool{
			"admin":    true,
			"maintain": true,
			"pull":     true,
			"push":     true,
			"triage":   true,
		},
		AllowForking:             Bool(true),
		WebCommitSignoffRequired: Bool(false),
		Archived:                 Bool(false),
		Disabled:                 Bool(false),
		Private:                  Bool(true),
		HasIssues:                Bool(false),
		HasWiki:                  Bool(false),
		HasPages:                 Bool(false),
		HasProjects:              Bool(false),
		HasDownloads:             Bool(false),
		HasDiscussions:           Bool(false),
		IsTemplate:               Bool(false),
		URL:                      String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               String("private"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	fork, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "o", "r", "ghsa_id")
	if _, ok := err.(*AcceptedError); !ok {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned error: %v (want AcceptedError)", err)
	}

	want := &Repository{
		ID:     Int64(1),
		NodeID: String("R_kgDPP3c6pQ"),
		Owner: &User{
			Login:             String("owner"),
			ID:                Int64(2),
			NodeID:            String("MDQ6VXFGcjYyMjcyMTQw"),
			AvatarURL:         String("https://avatars.githubusercontent.com/u/111111?v=4"),
			HTMLURL:           String("https://github.com/xxxxx"),
			GravatarID:        String(""),
			Type:              String("User"),
			SiteAdmin:         Bool(false),
			URL:               String("https://api.github.com/users/owner"),
			EventsURL:         String("https://api.github.com/users/owner/events{/privacy}"),
			FollowingURL:      String("https://api.github.com/users/owner/following{/other_user}"),
			FollowersURL:      String("https://api.github.com/users/owner/followers"),
			GistsURL:          String("https://api.github.com/users/owner/gists{/gist_id}"),
			OrganizationsURL:  String("https://api.github.com/users/owner/orgs"),
			ReceivedEventsURL: String("https://api.github.com/users/owner/received_events"),
			ReposURL:          String("https://api.github.com/users/owner/repos"),
			StarredURL:        String("https://api.github.com/users/owner/starred{/owner}{/repo}"),
			SubscriptionsURL:  String("https://api.github.com/users/owner/subscriptions"),
		},
		Name:             String("repo-ghsa-xxxx-xxxx-xxxx"),
		FullName:         String("owner/repo-ghsa-xxxx-xxxx-xxxx"),
		DefaultBranch:    String("master"),
		CreatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 41, 0, time.UTC)},
		PushedAt:         &Timestamp{time.Date(2023, time.December, 3, 11, 27, 8, 0, time.UTC)},
		UpdatedAt:        &Timestamp{time.Date(2023, time.December, 8, 17, 22, 42, 0, time.UTC)},
		HTMLURL:          String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		CloneURL:         String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		GitURL:           String("git://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SSHURL:           String("git@github.com:owner/repo-ghsa-xxxx-xxxx-xxxx.git"),
		SVNURL:           String("https://github.com/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		Fork:             Bool(false),
		ForksCount:       Int(0),
		NetworkCount:     Int(0),
		OpenIssuesCount:  Int(0),
		OpenIssues:       Int(0),
		StargazersCount:  Int(0),
		SubscribersCount: Int(0),
		WatchersCount:    Int(0),
		Watchers:         Int(0),
		Size:             Int(0),
		Permissions: map[string]bool{
			"admin":    true,
			"maintain": true,
			"pull":     true,
			"push":     true,
			"triage":   true,
		},
		AllowForking:             Bool(true),
		WebCommitSignoffRequired: Bool(false),
		Archived:                 Bool(false),
		Disabled:                 Bool(false),
		Private:                  Bool(true),
		HasIssues:                Bool(false),
		HasWiki:                  Bool(false),
		HasPages:                 Bool(false),
		HasProjects:              Bool(false),
		HasDownloads:             Bool(false),
		HasDiscussions:           Bool(false),
		IsTemplate:               Bool(false),
		URL:                      String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx"),
		ArchiveURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/{archive_format}{/ref}"),
		AssigneesURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/assignees{/user}"),
		BlobsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/blobs{/sha}"),
		BranchesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/branches{/branch}"),
		CollaboratorsURL:         String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/collaborators{/collaborator}"),
		CommentsURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/comments{/number}"),
		CommitsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/commits{/sha}"),
		CompareURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/compare/{base}...{head}"),
		ContentsURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contents/{+path}"),
		ContributorsURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/contributors"),
		DeploymentsURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/deployments"),
		DownloadsURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/downloads"),
		EventsURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/events"),
		ForksURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/forks"),
		GitCommitsURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/commits{/sha}"),
		GitRefsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/refs{/sha}"),
		GitTagsURL:               String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/git/tags{/sha}"),
		HooksURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/hooks"),
		IssueCommentURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/comments{/number}"),
		IssueEventsURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues/events{/number}"),
		IssuesURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/issues{/number}"),
		KeysURL:                  String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/keys{/key_id}"),
		LabelsURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/labels{/name}"),
		LanguagesURL:             String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/languages"),
		MergesURL:                String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/merges"),
		MilestonesURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/milestones{/number}"),
		NotificationsURL:         String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/notifications{?since,all,participating}"),
		PullsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/pulls{/number}"),
		ReleasesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/releases{/id}"),
		StargazersURL:            String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/stargazers"),
		StatusesURL:              String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/statuses/{sha}"),
		SubscribersURL:           String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscribers"),
		SubscriptionURL:          String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/subscription"),
		TagsURL:                  String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/tags"),
		TeamsURL:                 String("https://api.github.com/repos/owner/repo-ghsa-xxxx-xxxx-xxxx/teams"),
		Visibility:               String("private"),
	}
	if !cmp.Equal(fork, want) {
		t.Errorf("SecurityAdvisoriesService.CreateTemporaryPrivateFork returned %+v, want %+v", fork, want)
	}
}

func TestSecurityAdvisoriesService_CreateTemporaryPrivateFork_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.SecurityAdvisories.CreateTemporaryPrivateFork(ctx, "%", "r", "ghsa_id")
	testURLParseError(t, err)
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_BadRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_UnmarshalError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Errorf("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal number into Go struct field SecurityAdvisory.ghsa_id of type string") {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: String("GHSA-abcd-1234-efgh"),
			CVEID:  String("CVE-2050-00000"),
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisories returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_UnmarshalError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal number into Go struct field SecurityAdvisory.ghsa_id of type string") {
		t.Errorf("ListRepositorySecurityAdvisories returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisories returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: String("GHSA-abcd-1234-efgh"),
			CVEID:  String("CVE-2050-00000"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	opts := &ListGlobalSecurityAdvisoriesOptions{CVEID: String("CVE-xoxo-1234")}

	advisories, _, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, opts)
	if err != nil {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*GlobalSecurityAdvisory{
		{
			ID: Int64(1),
			SecurityAdvisory: SecurityAdvisory{
				GHSAID:      String("GHSA-xoxo-1234-xoxo"),
				CVEID:       String("CVE-xoxo-1234"),
				URL:         String("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
				HTMLURL:     String("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
				Severity:    String("high"),
				Summary:     String("Heartbleed security advisory"),
				Description: String("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
				Identifiers: []*AdvisoryIdentifier{
					{
						Type:  String("GHSA"),
						Value: String("GHSA-xoxo-1234-xoxo"),
					},
					{
						Type:  String("CVE"),
						Value: String("CVE-xoxo-1234"),
					},
				},
				PublishedAt: &date,
				UpdatedAt:   &date,
				WithdrawnAt: nil,
				CVSS: &AdvisoryCVSS{
					VectorString: String("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
					Score:        Float64(7.6),
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: String("CWE-400"),
						Name:  String("Uncontrolled Resource Consumption"),
					},
				},
			},
			References: []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
			Vulnerabilities: []*GlobalSecurityVulnerability{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: String("npm"),
						Name:      String("a-package"),
					},
					FirstPatchedVersion:    String("1.0.3"),
					VulnerableVersionRange: String("<=1.0.2"),
					VulnerableFunctions:    []string{"a_function"},
				},
			},
			RepositoryAdvisoryURL: String("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
			Type:                  String("reviewed"),
			SourceCodeLocation:    String("https://github.com/project/a-package"),
			GithubReviewedAt:      &date,
			NVDPublishedAt:        &date,
			Credits: []*Credit{
				{
					User: &User{
						Login:             String("user"),
						ID:                Int64(1),
						NodeID:            String("12="),
						AvatarURL:         String("a"),
						GravatarID:        String(""),
						URL:               String("a"),
						HTMLURL:           String("b"),
						FollowersURL:      String("b"),
						FollowingURL:      String("c"),
						GistsURL:          String("d"),
						StarredURL:        String("e"),
						SubscriptionsURL:  String("f"),
						OrganizationsURL:  String("g"),
						ReposURL:          String("h"),
						EventsURL:         String("i"),
						ReceivedEventsURL: String("j"),
						Type:              String("User"),
						SiteAdmin:         Bool(false),
					},
					Type: String("analyst"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	advisory, _, err := client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "GHSA-xoxo-1234-xoxo")
	if err != nil {
		t.Errorf("SecurityAdvisories.GetGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := &GlobalSecurityAdvisory{
		ID: Int64(1),
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      String("GHSA-xoxo-1234-xoxo"),
			CVEID:       String("CVE-xoxo-1234"),
			URL:         String("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     String("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    String("high"),
			Summary:     String("Heartbleed security advisory"),
			Description: String("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  String("GHSA"),
					Value: String("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  String("CVE"),
					Value: String("CVE-xoxo-1234"),
				},
			},
			PublishedAt: &date,
			UpdatedAt:   &date,
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: String("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        Float64(7.6),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: String("CWE-400"),
					Name:  String("Uncontrolled Resource Consumption"),
				},
			},
		},
		RepositoryAdvisoryURL: String("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
		Type:                  String("reviewed"),
		SourceCodeLocation:    String("https://github.com/project/a-package"),
		References:            []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
		GithubReviewedAt:      &date,
		NVDPublishedAt:        &date,

		Vulnerabilities: []*GlobalSecurityVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: String("npm"),
					Name:      String("a-package"),
				},
				FirstPatchedVersion:    String("1.0.3"),
				VulnerableVersionRange: String("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		Credits: []*Credit{
			{
				User: &User{
					Login:             String("user"),
					ID:                Int64(1),
					NodeID:            String("12="),
					AvatarURL:         String("a"),
					GravatarID:        String(""),
					URL:               String("a"),
					HTMLURL:           String("b"),
					FollowersURL:      String("b"),
					FollowingURL:      String("c"),
					GistsURL:          String("d"),
					StarredURL:        String("e"),
					SubscriptionsURL:  String("f"),
					OrganizationsURL:  String("g"),
					ReposURL:          String("h"),
					EventsURL:         String("i"),
					ReceivedEventsURL: String("j"),
					Type:              String("User"),
					SiteAdmin:         Bool(false),
				},
				Type: String("analyst"),
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
	testJSONMarshal(t, &SecurityAdvisorySubmission{}, `{}`)

	u := &SecurityAdvisorySubmission{
		Accepted: Bool(true),
	}

	w := `{
		"accepted": true
	}`

	testJSONMarshal(t, u, w)
}

func TestRepoAdvisoryCredit_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepoAdvisoryCredit{}, `{}`)

	u := &RepoAdvisoryCredit{
		Login: String("l"),
		Type:  String("t"),
	}

	w := `{
		"login": "l",
		"type": "t"
	}`

	testJSONMarshal(t, u, w)
}

func TestRepoAdvisoryCreditDetailed_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepoAdvisoryCreditDetailed{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &RepoAdvisoryCreditDetailed{
		Type:  String("t"),
		State: String("s"),
		User: &User{
			Name:                    String("u"),
			Company:                 String("c"),
			Blog:                    String("b"),
			Location:                String("l"),
			Email:                   String("e"),
			Hireable:                Bool(false),
			Bio:                     String("bio"),
			TwitterUsername:         String("tu"),
			PublicRepos:             Int(1),
			PublicGists:             Int(1),
			Followers:               Int(2),
			Following:               Int(2),
			CreatedAt:               testDate,
			UpdatedAt:               testDate,
			SuspendedAt:             testDate,
			Type:                    String("type"),
			SiteAdmin:               Bool(false),
			TotalPrivateRepos:       Int64(10),
			OwnedPrivateRepos:       Int64(10),
			PrivateGists:            Int(10),
			DiskUsage:               Int(10),
			Collaborators:           Int(10),
			TwoFactorAuthentication: Bool(true),
			Plan: &Plan{
				Name:          String("p"),
				Space:         Int(2),
				Collaborators: Int(2),
				PrivateRepos:  Int64(2),
				Seats:         Int(2),
				FilledSeats:   Int(1),
			},
			LdapDn:            String("l"),
			URL:               String("url"),
			EventsURL:         String("e"),
			FollowingURL:      String("f"),
			FollowersURL:      String("f"),
			GistsURL:          String("g"),
			OrganizationsURL:  String("o"),
			ReceivedEventsURL: String("r"),
			ReposURL:          String("rep"),
			StarredURL:        String("star"),
			SubscriptionsURL:  String("sub"),
			TextMatches: []*TextMatch{
				{
					ObjectURL:  String("u"),
					ObjectType: String("t"),
					Property:   String("p"),
					Fragment:   String("f"),
					Matches: []*Match{
						{
							Text:    String("t"),
							Indices: []int{1, 2},
						},
					},
				},
			},
			Permissions: map[string]bool{"p1": true},
			RoleName:    String("r"),
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
			"p1": true
		  },
		  "role_name": "r"
		}
	  }`

	testJSONMarshal(t, u, w)
}

func TestCredit_Marshal(t *testing.T) {
	testJSONMarshal(t, &Credit{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &Credit{
		Type: String("t"),
		User: &User{
			Name:                    String("u"),
			Company:                 String("c"),
			Blog:                    String("b"),
			Location:                String("l"),
			Email:                   String("e"),
			Hireable:                Bool(false),
			Bio:                     String("bio"),
			TwitterUsername:         String("tu"),
			PublicRepos:             Int(1),
			PublicGists:             Int(1),
			Followers:               Int(2),
			Following:               Int(2),
			CreatedAt:               testDate,
			UpdatedAt:               testDate,
			SuspendedAt:             testDate,
			Type:                    String("type"),
			SiteAdmin:               Bool(false),
			TotalPrivateRepos:       Int64(10),
			OwnedPrivateRepos:       Int64(10),
			PrivateGists:            Int(10),
			DiskUsage:               Int(10),
			Collaborators:           Int(10),
			TwoFactorAuthentication: Bool(true),
			Plan: &Plan{
				Name:          String("p"),
				Space:         Int(2),
				Collaborators: Int(2),
				PrivateRepos:  Int64(2),
				Seats:         Int(2),
				FilledSeats:   Int(1),
			},
			LdapDn:            String("l"),
			URL:               String("url"),
			EventsURL:         String("e"),
			FollowingURL:      String("f"),
			FollowersURL:      String("f"),
			GistsURL:          String("g"),
			OrganizationsURL:  String("o"),
			ReceivedEventsURL: String("r"),
			ReposURL:          String("rep"),
			StarredURL:        String("star"),
			SubscriptionsURL:  String("sub"),
			TextMatches: []*TextMatch{
				{
					ObjectURL:  String("u"),
					ObjectType: String("t"),
					Property:   String("p"),
					Fragment:   String("f"),
					Matches: []*Match{
						{
							Text:    String("t"),
							Indices: []int{1, 2},
						},
					},
				},
			},
			Permissions: map[string]bool{"p1": true},
			RoleName:    String("r"),
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
			"p1": true
			},
			"role_name": "r"
		}
	}`

	testJSONMarshal(t, u, w)
}

func TestGlobalSecurityAdvisory_Marshal(t *testing.T) {
	testJSONMarshal(t, &GlobalSecurityAdvisory{}, `{}`)

	testDate := &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}
	u := &GlobalSecurityAdvisory{
		ID:                    Int64(1),
		RepositoryAdvisoryURL: String("r"),
		Type:                  String("t"),
		SourceCodeLocation:    String("s"),
		References:            []string{"r"},
		Vulnerabilities: []*GlobalSecurityVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: String("npm"),
					Name:      String("a-package"),
				},
				FirstPatchedVersion:    String("1.0.3"),
				VulnerableVersionRange: String("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		GithubReviewedAt: testDate,
		NVDPublishedAt:   testDate,
		Credits: []*Credit{
			{
				Type: String("t"),
				User: &User{
					Name:                    String("u"),
					Company:                 String("c"),
					Blog:                    String("b"),
					Location:                String("l"),
					Email:                   String("e"),
					Hireable:                Bool(false),
					Bio:                     String("bio"),
					TwitterUsername:         String("tu"),
					PublicRepos:             Int(1),
					PublicGists:             Int(1),
					Followers:               Int(2),
					Following:               Int(2),
					CreatedAt:               testDate,
					UpdatedAt:               testDate,
					SuspendedAt:             testDate,
					Type:                    String("type"),
					SiteAdmin:               Bool(false),
					TotalPrivateRepos:       Int64(10),
					OwnedPrivateRepos:       Int64(10),
					PrivateGists:            Int(10),
					DiskUsage:               Int(10),
					Collaborators:           Int(10),
					TwoFactorAuthentication: Bool(true),
					Plan: &Plan{
						Name:          String("p"),
						Space:         Int(2),
						Collaborators: Int(2),
						PrivateRepos:  Int64(2),
						Seats:         Int(2),
						FilledSeats:   Int(1),
					},
					LdapDn:            String("l"),
					URL:               String("url"),
					EventsURL:         String("e"),
					FollowingURL:      String("f"),
					FollowersURL:      String("f"),
					GistsURL:          String("g"),
					OrganizationsURL:  String("o"),
					ReceivedEventsURL: String("r"),
					ReposURL:          String("rep"),
					StarredURL:        String("star"),
					SubscriptionsURL:  String("sub"),
					TextMatches: []*TextMatch{
						{
							ObjectURL:  String("u"),
							ObjectType: String("t"),
							Property:   String("p"),
							Fragment:   String("f"),
							Matches: []*Match{
								{
									Text:    String("t"),
									Indices: []int{1, 2},
								},
							},
						},
					},
					Permissions: map[string]bool{"p1": true},
					RoleName:    String("r"),
				},
			},
		},
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      String("GHSA-xoxo-1234-xoxo"),
			CVEID:       String("CVE-xoxo-1234"),
			URL:         String("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     String("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    String("high"),
			Summary:     String("Heartbleed security advisory"),
			Description: String("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  String("GHSA"),
					Value: String("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  String("CVE"),
					Value: String("CVE-xoxo-1234"),
				},
			},
			PublishedAt: testDate,
			UpdatedAt:   testDate,
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: String("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        Float64(7.6),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: String("CWE-400"),
					Name:  String("Uncontrolled Resource Consumption"),
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
				"p1": true
			  },
			  "role_name": "r"
			},
			"type": "t"
		  }
		]
	}`

	testJSONMarshal(t, u, w)
}
