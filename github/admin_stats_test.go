package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAdminService_GetAdminStats(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprise/stats/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
{
  "repos": {
    "total_repos": 212,
    "root_repos": 194,
    "fork_repos": 18,
    "org_repos": 51,
    "total_pushes": 3082,
    "total_wikis": 15
  },
  "hooks": {
    "total_hooks": 27,
    "active_hooks": 23,
    "inactive_hooks": 4
  },
  "pages": {
    "total_pages": 36
  },
  "orgs": {
    "total_orgs": 33,
    "disabled_orgs": 0,
    "total_teams": 60,
    "total_team_members": 314
  },
  "users": {
    "total_users": 254,
    "admin_users": 45,
    "suspended_users": 21
  },
  "pulls": {
    "total_pulls": 86,
    "merged_pulls": 60,
    "mergeable_pulls": 21,
    "unmergeable_pulls": 3
  },
  "issues": {
    "total_issues": 179,
    "open_issues": 83,
    "closed_issues": 96
  },
  "milestones": {
    "total_milestones": 7,
    "open_milestones": 6,
    "closed_milestones": 1
  },
  "gists": {
    "total_gists": 178,
    "private_gists": 151,
    "public_gists": 25
  },
  "comments": {
    "total_commit_comments": 6,
    "total_gist_comments": 28,
    "total_issue_comments": 366,
    "total_pull_request_comments": 30
  }
}
`)
	})

	ctx := context.Background()
	stats, _, err := client.Admin.GetAdminStats(ctx)
	if err != nil {
		t.Errorf("AdminService.GetAdminStats returned error: %v", err)
	}

	if want := testAdminStats; !reflect.DeepEqual(stats, want) {
		t.Errorf("AdminService.GetAdminStats returned %+v, want %+v", stats, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Admin.GetAdminStats(ctx)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' GetAdminStats = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' GetAdminStats resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' GetAdminStats err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Admin.GetAdminStats(ctx)
	if got != nil {
		t.Errorf("rate.Reset.Time > now GetAdminStats = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now GetAdminStats resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now GetAdminStats err = nil, want error")
	}
}

func TestAdminService_Stringify(t *testing.T) {
	want := "github.AdminStats{Issues:github.IssueStats{TotalIssues:179, OpenIssues:83, ClosedIssues:96}, Hooks:github.HookStats{TotalHooks:27, ActiveHooks:23, InactiveHooks:4}, Milestones:github.MilestoneStats{TotalMilestones:7, OpenMilestones:6, ClosedMilestones:1}, Orgs:github.OrgStats{TotalOrgs:33, DisabledOrgs:0, TotalTeams:60, TotalTeamMembers:314}, Comments:github.CommentStats{TotalCommitComments:6, TotalGistComments:28, TotalIssueComments:366, TotalPullRequestComments:30}, Pages:github.PageStats{TotalPages:36}, Users:github.UserStats{TotalUsers:254, AdminUsers:45, SuspendedUsers:21}, Gists:github.GistStats{TotalGists:178, PrivateGists:151, PublicGists:25}, Pulls:github.PullStats{TotalPulls:86, MergedPulls:60, MergablePulls:21, UnmergablePulls:3}, Repos:github.RepoStats{TotalRepos:212, RootRepos:194, ForkRepos:18, OrgRepos:51, TotalPushes:3082, TotalWikis:15}}"
	if got := testAdminStats.String(); got != want {
		t.Errorf("testAdminStats.String = %q, want %q", got, want)
	}

	want = "github.IssueStats{TotalIssues:179, OpenIssues:83, ClosedIssues:96}"
	if got := testAdminStats.Issues.String(); got != want {
		t.Errorf("testAdminStats.Issues.String = %q, want %q", got, want)
	}

	want = "github.HookStats{TotalHooks:27, ActiveHooks:23, InactiveHooks:4}"
	if got := testAdminStats.Hooks.String(); got != want {
		t.Errorf("testAdminStats.Hooks.String = %q, want %q", got, want)
	}

	want = "github.MilestoneStats{TotalMilestones:7, OpenMilestones:6, ClosedMilestones:1}"
	if got := testAdminStats.Milestones.String(); got != want {
		t.Errorf("testAdminStats.Milestones.String = %q, want %q", got, want)
	}

	want = "github.OrgStats{TotalOrgs:33, DisabledOrgs:0, TotalTeams:60, TotalTeamMembers:314}"
	if got := testAdminStats.Orgs.String(); got != want {
		t.Errorf("testAdminStats.Orgs.String = %q, want %q", got, want)
	}

	want = "github.CommentStats{TotalCommitComments:6, TotalGistComments:28, TotalIssueComments:366, TotalPullRequestComments:30}"
	if got := testAdminStats.Comments.String(); got != want {
		t.Errorf("testAdminStats.Comments.String = %q, want %q", got, want)
	}

	want = "github.PageStats{TotalPages:36}"
	if got := testAdminStats.Pages.String(); got != want {
		t.Errorf("testAdminStats.Pages.String = %q, want %q", got, want)
	}

	want = "github.UserStats{TotalUsers:254, AdminUsers:45, SuspendedUsers:21}"
	if got := testAdminStats.Users.String(); got != want {
		t.Errorf("testAdminStats.Users.String = %q, want %q", got, want)
	}

	want = "github.GistStats{TotalGists:178, PrivateGists:151, PublicGists:25}"
	if got := testAdminStats.Gists.String(); got != want {
		t.Errorf("testAdminStats.Gists.String = %q, want %q", got, want)
	}

	want = "github.PullStats{TotalPulls:86, MergedPulls:60, MergablePulls:21, UnmergablePulls:3}"
	if got := testAdminStats.Pulls.String(); got != want {
		t.Errorf("testAdminStats.Pulls.String = %q, want %q", got, want)
	}

	want = "github.RepoStats{TotalRepos:212, RootRepos:194, ForkRepos:18, OrgRepos:51, TotalPushes:3082, TotalWikis:15}"
	if got := testAdminStats.Repos.String(); got != want {
		t.Errorf("testAdminStats.Repos.String = %q, want %q", got, want)
	}
}

var testAdminStats = &AdminStats{
	Repos: &RepoStats{
		TotalRepos:  Int(212),
		RootRepos:   Int(194),
		ForkRepos:   Int(18),
		OrgRepos:    Int(51),
		TotalPushes: Int(3082),
		TotalWikis:  Int(15),
	},
	Hooks: &HookStats{
		TotalHooks:    Int(27),
		ActiveHooks:   Int(23),
		InactiveHooks: Int(4),
	},
	Pages: &PageStats{
		TotalPages: Int(36),
	},
	Orgs: &OrgStats{
		TotalOrgs:        Int(33),
		DisabledOrgs:     Int(0),
		TotalTeams:       Int(60),
		TotalTeamMembers: Int(314),
	},
	Users: &UserStats{
		TotalUsers:     Int(254),
		AdminUsers:     Int(45),
		SuspendedUsers: Int(21),
	},
	Pulls: &PullStats{
		TotalPulls:      Int(86),
		MergedPulls:     Int(60),
		MergablePulls:   Int(21),
		UnmergablePulls: Int(3),
	},
	Issues: &IssueStats{
		TotalIssues:  Int(179),
		OpenIssues:   Int(83),
		ClosedIssues: Int(96),
	},
	Milestones: &MilestoneStats{
		TotalMilestones:  Int(7),
		OpenMilestones:   Int(6),
		ClosedMilestones: Int(1),
	},
	Gists: &GistStats{
		TotalGists:   Int(178),
		PrivateGists: Int(151),
		PublicGists:  Int(25),
	},
	Comments: &CommentStats{
		TotalCommitComments:      Int(6),
		TotalGistComments:        Int(28),
		TotalIssueComments:       Int(366),
		TotalPullRequestComments: Int(30),
	},
}
