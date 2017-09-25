package github

import (
	"context"
	"fmt"
)

type AdminStats struct {
	Issues     *IssueStats     `json:"issues,omitempty"`
	Hooks      *HookStats      `json:"hooks,omitempty"`
	Milestones *MilestoneStats `json:"milestones,omitempty"`
	Orgs       *OrgStats       `json:"orgs,omitempty"`
	Comments   *CommentStats   `json:"comments,omitempty"`
	Pages      *PageStats      `json:"pages,omitempty"`
	Users      *UserStats      `json:"users,omitempty"`
	Gists      *GistStats      `json:"gists,omitempty"`
	Pulls      *PullStats      `json:"pulls,omitempty"`
	Repos      *RepoStats      `json:"repos,omitempty"`
}

type IssueStats struct {
	TotalIssues  *int `json:"total_issues,omitempty"`
	OpenIssues   *int `json:"open_issues,omitempty"`
	ClosedIssues *int `json:"closed_issues,omitempty"`
}

type HookStats struct {
	TotalHooks    *int `json:"total_hooks,omitempty"`
	ActiveHooks   *int `json:"active_hooks,omitempty"`
	InactiveHooks *int `json:"inactive_hooks,omitempty"`
}

type MilestoneStats struct {
	TotalMilestones  *int `json:"total_milestones,omitempty"`
	OpenMilestones   *int `json:"open_milestones,omitempty"`
	ClosedMilestones *int `json:"closed_milestones,omitempty"`
}

type OrgStats struct {
	TotalOrgs        *int `json:"total_orgs,omitempty"`
	DisabledOrgs     *int `json:"disabled_orgs,omitempty"`
	TotalTeams       *int `json:"total_teams,omitempty"`
	TotalTeamMembers *int `json:"total_team_members,omitempty"`
}

type CommentStats struct {
	TotalCommitComments      *int `json:"total_commit_comments,omitempty"`
	TotalGistComments        *int `json:"total_gist_comments,omitempty"`
	TotalIssueComments       *int `json:"total_issue_comments,omitempty"`
	TotalPullRequestComments *int `json:"total_pull_request_comments,omitempty"`
}

type PageStats struct {
	TotalPages *int `json:"total_pages,omitempty"`
}

type UserStats struct {
	TotalUsers     *int `json:"total_users,omitempty"`
	AdminUsers     *int `json:"admin_users,omitempty"`
	SuspendedUsers *int `json:"suspended_users,omitempty"`
}

type GistStats struct {
	TotalGists   *int `json:"total_gists,omitempty"`
	PrivateGists *int `json:"private_gists,omitempty"`
	PublicGists  *int `json:"public_gists,omitempty"`
}

type PullStats struct {
	TotalPulls      *int `json:"total_pulls,omitempty"`
	MergedPulls     *int `json:"merged_pulls,omitempty"`
	MergablePulls   *int `json:"mergeable_pulls,omitempty"`
	UnmergablePulls *int `json:"unmergeable_pulls,omitempty"`
}

type RepoStats struct {
	TotalRepos  *int `json:"total_repos,omitempty"`
	RootRepos   *int `json:"root_repos,omitempty"`
	ForkRepos   *int `json:"fork_repos,omitempty"`
	OrgRepos    *int `json:"org_repos,omitempty"`
	TotalPushes *int `json:"total_pushes,omitempty"`
	TotalWikis  *int `json:"total_wikis,omitempty"`
}

func (s *AdminService) GetAdminStats(ctx context.Context) (*AdminStats, *Response, error) {
	u := fmt.Sprintf("enterprise/stats/all")
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	m := new(AdminStats)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}
