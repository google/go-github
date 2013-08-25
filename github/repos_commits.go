// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

// Commit represents a commit in a repo.
// Note that it's wrapping a Commit, so author/committer information is in two places,
// but contain different details about them: in RepositoryCommit "github details", in Commit - "git details".
type RepositoryCommit struct {
	SHA       *string                 `json:"sha,omitempty"`
	Commit    *Commit                 `json:"commit,omitempty"`
	Author    *RepositoryCommitAuthor `json:"author,omitempty"`
	Committer *RepositoryCommitAuthor `json:"committer,omitempty"`
	Parents   []Commit                `json:"parents,omitempty"`

	// Details about how many changes were made in this commit. Only filled in during GetCommit!
	Stats *CommitStats `json:"stats,omitempty"`
	// Details about which files, and how this commit touched. Only filled in during GetCommit!
	Files []CommitFile `json:"files,omitempty"`
}

// RepositoryCommitAuthor represents the author or committer of a commit.
// It should correspond to a GitHub user, unlike CommitAuthor.
type RepositoryCommitAuthor struct {
	Login      *string `json:"login,omitempty"`
	Id         *int    `json:"id,omitempty"`
	AvatarUrl  *string `json:"avatar_url,omitempty"`
	GravatarId *string `json:"gravatar_id,omitempty"`
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

// just like a Commit
type BaseComparedCommit struct {
	SHA       *string                 `json:"sha,omitempty"`
	Commit    *Commit                 `json:"commit,omitempty"`
	Author    *RepositoryCommitAuthor `json:"author,omitempty"`
	Committer *RepositoryCommitAuthor `json:"committer,omitempty"`
	Message   *string                 `json:"message,omitempty"`
	Tree      *Tree                   `json:"tree,omitempty"`
	Parents   []Commit                `json:"parents,omitempty"`
}

// ListCommits lists the commits of a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/commits/#list
func (s *RepositoriesService) ListCommits(owner, repo string, opts *CommitsListOptions) ([]RepositoryCommit, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits", owner, repo)
	if opts != nil {
		u = withOptionalParameter(u, "sha", opts.SHA)
		u = withOptionalParameter(u, "path", opts.Path)
		u = withOptionalParameter(u, "author", opts.Author)
		u = withOptionalTimeParameter(u, "since", *opts.Since)
		u = withOptionalTimeParameter(u, "until", *opts.Until) // todo think if this is nice or not hm hm -- ktoso
	}

	// todo test if we really call the API I think we are..... ;-)
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
