// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type SecurityAdvisoriesService service

// SecurityAdvisorySubmission represents the Security Advisory Submission.
type SecurityAdvisorySubmission struct {
	// Accepted represents whether a private vulnerability report was accepted by the repository's administrators.
	Accepted *bool `json:"accepted,omitempty"`
}

// Vulnerability represents the vulnerability object for a Security Advisory.
type Vulnerability struct {
	Package                *VulnerabilityPackage `json:"package,omitempty"`
	VulnerableVersionRange *string               `json:"vulnerable_version_range,omitempty"`
	PatchedVersions        *string               `json:"patched_versions,omitempty"`
	VulnerableFunctions    []string              `json:"vulnerable_functions,omitempty"`
}

// RepoAdvisoryCredit represents the credit object for a repository Security Advisory.
type RepoAdvisoryCredit struct {
	Login *string `json:"login,omitempty"`
	Type  *string `json:"type,omitempty"`
}

// RepoAdvisoryCreditDetailed represents a credit given to a user for a repository Security Advisory.
type RepoAdvisoryCreditDetailed struct {
	User  *User   `json:"user,omitempty"`
	Type  *string `json:"type,omitempty"`
	State *string `json:"state,omitempty"`
}

// Permissions represent a team's permissions.
type Permissions struct {
	TeamPermissionsFrom
	Triage   *bool `json:"triage,omitempty"`
	Maintain *bool `json:"maintain,omitempty"`
}

// TeamParent represents a team's parent team.
type TeamParent struct {
	ID                  *int64  `json:"id,omitempty"`
	NodeID              *string `json:"node_id,omitempty"`
	URL                 *string `json:"url,omitempty"`
	MembersURL          *string `json:"members_url,omitempty"`
	Name                *string `json:"name,omitempty"`
	Description         *string `json:"description,omitempty"`
	Permission          *string `json:"permission,omitempty"`
	Privacy             *string `json:"privacy,omitempty"`
	NotificationSetting *string `json:"notification_setting,omitempty"`
	HTMLURL             *string `json:"html_url,omitempty"`
	RepositoriesURL     *string `json:"repositories_url,omitempty"`
	Slug                *string `json:"slug,omitempty"`
	LDAPDN              *string `json:"ldap_dn,omitempty"`
}

// PrivateFork represents a temporary private fork of the advisory's repository for collaborating on a fix.
type PrivateFork struct {
	ID               *int64  `json:"id,omitempty"`
	NodeID           *string `json:"node_id,omitempty"`
	Name             *string `json:"name,omitempty"`
	FullName         *string `json:"full_name,omitempty"`
	Owner            *User   `json:"owner,omitempty"`
	Private          *bool   `json:"private,omitempty"`
	HTMLURL          *string `json:"html_url,omitempty"`
	Description      *string `json:"description,omitempty"`
	Fork             *bool   `json:"fork,omitempty"`
	URL              *string `json:"url,omitempty"`
	ArchiveURL       *string `json:"archive_url,omitempty"`
	AssigneesURL     *string `json:"assignees_url,omitempty"`
	BlobsURL         *string `json:"blobs_url,omitempty"`
	BranchesURL      *string `json:"branches_url,omitempty"`
	CollaboratorsURL *string `json:"collaborators_url,omitempty"`
	CommentsURL      *string `json:"comments_url,omitempty"`
	CommitsURL       *string `json:"commits_url,omitempty"`
	CompareURL       *string `json:"compare_url,omitempty"`
	ContentsURL      *string `json:"contents_url,omitempty"`
	ContributorsURL  *string `json:"contributors_url,omitempty"`
	DeploymentsURL   *string `json:"deployments_url,omitempty"`
	DownloadsURL     *string `json:"downloads_url,omitempty"`
	EventsURL        *string `json:"events_url,omitempty"`
	ForksURL         *string `json:"forks_url,omitempty"`
	GitCommitsURL    *string `json:"git_commits_url,omitempty"`
	GitRefsURL       *string `json:"git_refs_url,omitempty"`
	GitTagsURL       *string `json:"git_tags_url,omitempty"`
	IssueCommentURL  *string `json:"issue_comment_url,omitempty"`
	IssueEventsURL   *string `json:"issue_events_url,omitempty"`
	IssuesURL        *string `json:"issues_url,omitempty"`
	KeysURL          *string `json:"keys_url,omitempty"`
	LabelsURL        *string `json:"labels_url,omitempty"`
	LanguagesURL     *string `json:"languages_url,omitempty"`
	MergesURL        *string `json:"merges_url,omitempty"`
	MilestonesURL    *string `json:"milestones_url,omitempty"`
	NotificationsURL *string `json:"notifications_url,omitempty"`
	PullsURL         *string `json:"pulls_url,omitempty"`
	ReleasesURL      *string `json:"releases_url,omitempty"`
	StargazersURL    *string `json:"stargazers_url,omitempty"`
	StatusesURL      *string `json:"statuses_url,omitempty"`
	SubscribersURL   *string `json:"subscribers_url,omitempty"`
	SubscriptionURL  *string `json:"subscription_url,omitempty"`
	TagsURL          *string `json:"tags_url,omitempty"`
	TeamsURL         *string `json:"teams_url,omitempty"`
	TreesURL         *string `json:"trees_url,omitempty"`
	HooksURL         *string `json:"hooks_url,omitempty"`
}

// RepoSecurityAdvisory represents a repository security advisory.
type RepoSecurityAdvisory struct {
	GHSAID             *string                       `json:"ghsa_id,omitempty"`
	CVEID              *string                       `json:"cve_id,omitempty"`
	URL                *string                       `json:"url,omitempty"`
	HTMLURL            *string                       `json:"html_url,omitempty"`
	Summary            *string                       `json:"summary,omitempty"`
	Description        *string                       `json:"description,omitempty"`
	Severity           *string                       `json:"severity,omitempty"`
	Author             *User                         `json:"author,omitempty"`
	Publisher          *User                         `json:"publisher,omitempty"`
	Identifiers        []*AdvisoryIdentifier         `json:"identifiers,omitempty"`
	State              *string                       `json:"state,omitempty"`
	CreatedAt          *Timestamp                    `json:"created_at,omitempty"`
	UpdatedAt          *Timestamp                    `json:"updated_at,omitempty"`
	PublishedAt        *Timestamp                    `json:"published_at,omitempty"`
	ClosedAt           *Timestamp                    `json:"closed_at,omitempty"`
	WithdrawnAt        *Timestamp                    `json:"withdrawn_at,omitempty"`
	Submission         *SecurityAdvisorySubmission   `json:"submission,omitempty"`
	Vulnerabilities    []*Vulnerability              `json:"vulnerabilities,omitempty"`
	CVSs               *AdvisoryCVSs                 `json:"cvss,omitempty"`
	CWEs               []*AdvisoryCWEs               `json:"cwes,omitempty"`
	CWEIDs             []string                      `json:"cwe_ids,omitempty"`
	Credits            []*RepoAdvisoryCredit         `json:"credits,omitempty"`
	CreditsDetailed    []*RepoAdvisoryCreditDetailed `json:"credits_detailed,omitempty"`
	CollaboratingUsers []*User                       `json:"collaborating_users,omitempty"`
	CollaboratingTeams []*Team                       `json:"collaborating_teams,omitempty"`
	PrivateFork        *PrivateFork                  `json:"private_fork,omitempty"`
}

// ListRepositorySecurityAdvisoriesOptions specifies the optional parameters to lists the repository security advisories.
type ListRepositorySecurityAdvisoriesOptions struct {
	// Direction in which to sort advisories. Possible values are: asc, desc.
	// Default is "asc".
	Direction string `url:"direction,omitempty"`

	// Sort specifies how to sort advisories. Possible values are: created, updated,
	// and published. Default value is "created".
	Sort string `url:"sort,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for events before this cursor.
	Before string `url:"before,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for events after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of advisories to include per page.
	PerPage int `url:"per_page,omitempty"`

	// State filters advisories based on their state. Possible values are: triage, draft, published, closed.
	State string `url:"state,omitempty"`
}

// RequestCVE requests a Common Vulnerabilities and Exposures (CVE) for a repository security advisory.
// The ghsaID is the GitHub Security Advisory identifier of the advisory.
//
// GitHub API docs: https://docs.github.com/en/rest/security-advisories/repository-advisories#request-a-cve-for-a-repository-security-advisory
func (s *SecurityAdvisoriesService) RequestCVE(ctx context.Context, owner, repo, ghsaID string) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/security-advisories/%v/cve", owner, repo, ghsaID)

	req, err := s.client.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if _, ok := err.(*AcceptedError); ok {
			return resp, nil
		}

		return resp, err
	}

	return resp, nil
}

// ListRepositorySecurityAdvisoriesForOrg lists the repository security advisories for an organization.
//
// Github API docs: https://docs.github.com/en/rest/security-advisories/repository-advisories?apiVersion=2022-11-28#list-repository-security-advisories-for-an-organization
func (s *SecurityAdvisoriesService) ListRepositorySecurityAdvisoriesForOrg(ctx context.Context, org string, opt *ListRepositorySecurityAdvisoriesOptions) ([]*RepoSecurityAdvisory, *Response, error) {
	url := fmt.Sprintf("orgs/%v/security-advisories", org)
	url, err := addOptions(url, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var advisories []*RepoSecurityAdvisory
	resp, err := s.client.Do(ctx, req, &advisories)
	if err != nil {
		return nil, resp, err
	}

	return advisories, resp, nil
}

// ListRepositorySecurityAdvisories lists the security advisories in a repository.
//
// Github API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/security-advisories/repository-advisories?apiVersion=2022-11-28#list-repository-security-advisories
func (s *SecurityAdvisoriesService) ListRepositorySecurityAdvisories(ctx context.Context, owner string, repo string, opt *ListRepositorySecurityAdvisoriesOptions) ([]*RepoSecurityAdvisory, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/security-advisories", owner, repo)
	url, err := addOptions(url, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var advisories []*RepoSecurityAdvisory
	resp, err := s.client.Do(ctx, req, &advisories)
	if err != nil {
		return nil, resp, err
	}

	return advisories, resp, nil
}
