package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAdminService_GetAdminStats(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	if want := testAdminStats; !cmp.Equal(stats, want) {
		t.Errorf("AdminService.GetAdminStats returned %+v, want %+v", stats, want)
	}

	const methodName = "GetAdminStats"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Admin.GetAdminStats(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAdminService_Stringify(t *testing.T) {
	t.Parallel()
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
		TotalRepos:  Ptr(212),
		RootRepos:   Ptr(194),
		ForkRepos:   Ptr(18),
		OrgRepos:    Ptr(51),
		TotalPushes: Ptr(3082),
		TotalWikis:  Ptr(15),
	},
	Hooks: &HookStats{
		TotalHooks:    Ptr(27),
		ActiveHooks:   Ptr(23),
		InactiveHooks: Ptr(4),
	},
	Pages: &PageStats{
		TotalPages: Ptr(36),
	},
	Orgs: &OrgStats{
		TotalOrgs:        Ptr(33),
		DisabledOrgs:     Ptr(0),
		TotalTeams:       Ptr(60),
		TotalTeamMembers: Ptr(314),
	},
	Users: &UserStats{
		TotalUsers:     Ptr(254),
		AdminUsers:     Ptr(45),
		SuspendedUsers: Ptr(21),
	},
	Pulls: &PullStats{
		TotalPulls:      Ptr(86),
		MergedPulls:     Ptr(60),
		MergablePulls:   Ptr(21),
		UnmergablePulls: Ptr(3),
	},
	Issues: &IssueStats{
		TotalIssues:  Ptr(179),
		OpenIssues:   Ptr(83),
		ClosedIssues: Ptr(96),
	},
	Milestones: &MilestoneStats{
		TotalMilestones:  Ptr(7),
		OpenMilestones:   Ptr(6),
		ClosedMilestones: Ptr(1),
	},
	Gists: &GistStats{
		TotalGists:   Ptr(178),
		PrivateGists: Ptr(151),
		PublicGists:  Ptr(25),
	},
	Comments: &CommentStats{
		TotalCommitComments:      Ptr(6),
		TotalGistComments:        Ptr(28),
		TotalIssueComments:       Ptr(366),
		TotalPullRequestComments: Ptr(30),
	},
}

func TestIssueStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueStats{}, "{}")

	u := &IssueStats{
		TotalIssues:  Ptr(1),
		OpenIssues:   Ptr(1),
		ClosedIssues: Ptr(1),
	}

	want := `{
		"total_issues": 1,
		"open_issues": 1,
		"closed_issues": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestHookStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &HookStats{}, "{}")

	u := &HookStats{
		TotalHooks:    Ptr(1),
		ActiveHooks:   Ptr(1),
		InactiveHooks: Ptr(1),
	}

	want := `{
		"total_hooks": 1,
		"active_hooks": 1,
		"inactive_hooks": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestMilestoneStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MilestoneStats{}, "{}")

	u := &MilestoneStats{
		TotalMilestones:  Ptr(1),
		OpenMilestones:   Ptr(1),
		ClosedMilestones: Ptr(1),
	}

	want := `{
		"total_milestones": 1,
		"open_milestones": 1,
		"closed_milestones": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestOrgStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &OrgStats{}, "{}")

	u := &OrgStats{
		TotalOrgs:        Ptr(1),
		DisabledOrgs:     Ptr(1),
		TotalTeams:       Ptr(1),
		TotalTeamMembers: Ptr(1),
	}

	want := `{
		"total_orgs": 1,
		"disabled_orgs": 1,
		"total_teams": 1,
		"total_team_members": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestCommentStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CommentStats{}, "{}")

	u := &CommentStats{
		TotalCommitComments:      Ptr(1),
		TotalGistComments:        Ptr(1),
		TotalIssueComments:       Ptr(1),
		TotalPullRequestComments: Ptr(1),
	}

	want := `{
		"total_commit_comments": 1,
		"total_gist_comments": 1,
		"total_issue_comments": 1,
		"total_pull_request_comments": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestPageStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PageStats{}, "{}")

	u := &PageStats{
		TotalPages: Ptr(1),
	}

	want := `{
		"total_pages": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestUserStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UserStats{}, "{}")

	u := &UserStats{
		TotalUsers:     Ptr(1),
		AdminUsers:     Ptr(1),
		SuspendedUsers: Ptr(1),
	}

	want := `{
		"total_users": 1,
		"admin_users": 1,
		"suspended_users": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestGistStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GistStats{}, "{}")

	u := &GistStats{
		TotalGists:   Ptr(1),
		PrivateGists: Ptr(1),
		PublicGists:  Ptr(1),
	}

	want := `{
		"total_gists": 1,
		"private_gists": 1,
		"public_gists": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestPullStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullStats{}, "{}")

	u := &PullStats{
		TotalPulls:      Ptr(1),
		MergedPulls:     Ptr(1),
		MergablePulls:   Ptr(1),
		UnmergablePulls: Ptr(1),
	}

	want := `{
		"total_pulls": 1,
		"merged_pulls": 1,
		"mergeable_pulls": 1,
		"unmergeable_pulls": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestRepoStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepoStats{}, "{}")

	u := &RepoStats{
		TotalRepos:  Ptr(1),
		RootRepos:   Ptr(1),
		ForkRepos:   Ptr(1),
		OrgRepos:    Ptr(1),
		TotalPushes: Ptr(1),
		TotalWikis:  Ptr(1),
	}

	want := `{
		"total_repos": 1,
		"root_repos": 1,
		"fork_repos": 1,
		"org_repos": 1,
		"total_pushes": 1,
		"total_wikis": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestAdminStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AdminStats{}, "{}")

	u := &AdminStats{
		Repos: &RepoStats{
			TotalRepos:  Ptr(212),
			RootRepos:   Ptr(194),
			ForkRepos:   Ptr(18),
			OrgRepos:    Ptr(51),
			TotalPushes: Ptr(3082),
			TotalWikis:  Ptr(15),
		},
		Hooks: &HookStats{
			TotalHooks:    Ptr(27),
			ActiveHooks:   Ptr(23),
			InactiveHooks: Ptr(4),
		},
		Pages: &PageStats{
			TotalPages: Ptr(36),
		},
		Orgs: &OrgStats{
			TotalOrgs:        Ptr(33),
			DisabledOrgs:     Ptr(0),
			TotalTeams:       Ptr(60),
			TotalTeamMembers: Ptr(314),
		},
		Users: &UserStats{
			TotalUsers:     Ptr(254),
			AdminUsers:     Ptr(45),
			SuspendedUsers: Ptr(21),
		},
		Pulls: &PullStats{
			TotalPulls:      Ptr(86),
			MergedPulls:     Ptr(60),
			MergablePulls:   Ptr(21),
			UnmergablePulls: Ptr(3),
		},
		Issues: &IssueStats{
			TotalIssues:  Ptr(179),
			OpenIssues:   Ptr(83),
			ClosedIssues: Ptr(96),
		},
		Milestones: &MilestoneStats{
			TotalMilestones:  Ptr(7),
			OpenMilestones:   Ptr(6),
			ClosedMilestones: Ptr(1),
		},
		Gists: &GistStats{
			TotalGists:   Ptr(178),
			PrivateGists: Ptr(151),
			PublicGists:  Ptr(25),
		},
		Comments: &CommentStats{
			TotalCommitComments:      Ptr(6),
			TotalGistComments:        Ptr(28),
			TotalIssueComments:       Ptr(366),
			TotalPullRequestComments: Ptr(30),
		},
	}

	want := `{
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
	}`

	testJSONMarshal(t, u, want)
}
