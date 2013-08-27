// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Commit represents a commit in a repo.
// Note that it's wrapping a Commit, so author/committer information is in two places,
// but contain different details about them: in RepositoryCommit "github details", in Commit - "git details".
type RepositoryCommit struct {
	SHA       *string  `json:"sha,omitempty"`
	Commit    *Commit  `json:"commit,omitempty"`
	Author    *User    `json:"author,omitempty"`
	Committer *User    `json:"committer,omitempty"`
	Parents   []Commit `json:"parents,omitempty"`

	// Details about how many changes were made in this commit. Only filled in during GetCommit!
	Stats *CommitStats `json:"stats,omitempty"`
	// Details about which files, and how this commit touched. Only filled in during GetCommit!
	Files []CommitFile `json:"files,omitempty"`
}

// CommitsListOptions specifies the optional parameters to the
// RepositoriesService.RepositoriesList method.
type CommitsListOptions struct {
	// Sha or branch to start listing commits from.
	SHA *string
	// Only commits containing this file path will be returned.
	Path *string
	// GitHub login, name, or email by which to filter by commit author
	Author *string
	// Only commits after this date will be returned
	Since *time.Time
	// Only commits before this date will be returned
	Until *time.Time
}

type CommitStats struct {
	Additions *int `json:"additions,omitempty"`
	Deletions *int `json:"deletions,omitempty"`
	Total     *int `json:"total,omitempty"`
}

type CommitFile struct {
	SHA       *string `json:"sha,omitempty"`
	Filename  *string `json:"filename,omitempty"`
	Additions *int    `json:"additions,omitempty"`
	Deletions *int    `json:"deletions,omitempty"`
	Changes   *int    `json:"changes,omitempty"`
	Status    *string `json:"status,omitempty"`
	Patch     *string `json:"patch,omitempty"`
}

// CommitsComparison is the result of comparing two commits.
// See CompareCommits for details.
type CommitsComparison struct {
	BaseCommit *BaseComparedCommit `json:"base_commit,omitempty"`

	// Head can be 'behind' or 'ahead'
	Status       *string `json:"status,omitempty"`
	AheadBy      *int    `json:"ahead_by,omitempty"`
	BehindBy     *int    `json:"behind_by,omitempty"`
	TotalCommits *int    `json:"total_commits,omitempty"`

	Commits []RepositoryCommit `json:"commits,omitempty"`

	Files []CommitFile `json:"files,omitempty"`
}

// Like Commit but wraps around it, provifing additional author/committer User information.
type BaseComparedCommit struct {
	SHA       *string  `json:"sha,omitempty"`
	Commit    *Commit  `json:"commit,omitempty"`
	Author    *User    `json:"author,omitempty"`
	Committer *User    `json:"committer,omitempty"`
	Message   *string  `json:"message,omitempty"`
	Tree      *Tree    `json:"tree,omitempty"`
	Parents   []Commit `json:"parents,omitempty"`
}

// ListCommits lists the commits of a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/commits/#list
func (s *RepositoriesService) ListCommits(owner, repo string, opts *CommitsListOptions) ([]RepositoryCommit, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits", owner, repo)

	if opts != nil {
		params := url.Values{}
		params.Add("sha", *opts.SHA)
		params.Add("path", *opts.Path)
		params.Add("author", *opts.Author)
		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format(time.RFC3339))
		}
		if !opts.Until.IsZero() {
			params.Add("until", opts.Until.Format(time.RFC3339))
		}

		u = appendUrlParams(u, params)
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	commits := new([]RepositoryCommit)
	resp, err := s.client.Do(req, commits)
	return *commits, resp, err
}

// GetCommit fetches the specified commit, including all details about it.
// todo: support media formats - https://github.com/google/go-github/issues/6
//
// GitHub API docs: http://developer.github.com/v3/repos/commits/#get-a-single-commit
// See also: http://developer.github.com//v3/git/commits/#get-a-single-commit provides the same functionality
func (s *RepositoriesService) GetCommit(owner, repo, sha string) (*RepositoryCommit, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits/%v", owner, repo, sha)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	commit := new(RepositoryCommit)
	resp, err := s.client.Do(req, commit)
	return commit, resp, err
}

// Compare a range of commits with each other.
// todo: support media formats - https://github.com/google/go-github/issues/6
//
// GitHub API docs: http://developer.github.com/v3/repos/commits/index.html#compare-two-commits
func (s *RepositoriesService) CompareCommits(owner, repo string, base, head string) (*CommitsComparison, *Response, error) { // todo I'm sure I missspelled this
	u := fmt.Sprintf("repos/%v/%v/compare/%v...%v", owner, repo, base, head)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	comp := new(CommitsComparison)
	resp, err := s.client.Do(req, comp)
	return comp, resp, err
}

func appendUrlParams(baseUrl string, params url.Values) string {
	if strings.Contains(baseUrl, "?") {
		return baseUrl + "&" + params.Encode()
	}
	return baseUrl + "?" + params.Encode()
}
