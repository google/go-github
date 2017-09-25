package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAdminService_GetAdminStats(t *testing.T) {
	setup()
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

	stats, _, err := client.Admin.GetAdminStats(context.Background())
	if err != nil {
		t.Errorf("AdminService.GetAdminStats returned error: %v", err)
	}

	want := &AdminStats{
		Repos: &RepoStats{
			Int(212),
			Int(194),
			Int(18),
			Int(51),
			Int(3082),
			Int(15),
		},
		Hooks: &HookStats{
			Int(27),
			Int(23),
			Int(4),
		},
		Pages: &PageStats{
			Int(36),
		},
		Orgs: &OrgStats{
			Int(33),
			Int(0),
			Int(60),
			Int(314),
		},
		Users: &UserStats{
			Int(254),
			Int(45),
			Int(21),
		},
		Pulls: &PullStats{
			Int(86),
			Int(60),
			Int(21),
			Int(3),
		},
		Issues: &IssueStats{
			Int(179),
			Int(83),
			Int(96),
		},
		Milestones: &MilestoneStats{
			Int(7),
			Int(6),
			Int(1),
		},
		Gists: &GistStats{
			Int(178),
			Int(151),
			Int(25),
		},
		Comments: &CommentStats{
			Int(6),
			Int(28),
			Int(366),
			Int(30),
		},
	}
	if !reflect.DeepEqual(stats, want) {
		t.Errorf("AdminService.GetAdminStats returned %+v, want %+v", stats, want)
	}
}
