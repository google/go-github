// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// IssueArtifact represents an artifact attached to an issue.
// Number identifies the artifact within its issue; ID is its global identifier.
type IssueArtifact struct {
	ID             int64      `json:"id"`
	Number         int        `json:"number"`
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	Body           string     `json:"body"`
	BodyHTML       *string    `json:"body_html,omitempty"`
	Description    *string    `json:"description,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	Creator        *User      `json:"creator,omitempty"`
	UpdatedByActor *User      `json:"updated_by_actor,omitempty"`
	CreatedAt      *Timestamp `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp `json:"updated_at,omitempty"`
}

func (a IssueArtifact) String() string {
	return Stringify(a)
}

// IssueArtifactVersion represents a saved version of an issue artifact.
type IssueArtifactVersion struct {
	Version   int        `json:"version"`
	Name      string     `json:"name"`
	Body      string     `json:"body"`
	BodyHTML  *string    `json:"body_html,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	Actor     *User      `json:"actor,omitempty"`
}

func (v IssueArtifactVersion) String() string {
	return Stringify(v)
}

// IssueArtifactWithVersions contains an artifact and its edit history.
type IssueArtifactWithVersions struct {
	Artifact *IssueArtifact          `json:"artifact"`
	Versions []*IssueArtifactVersion `json:"versions"`
}

func (a IssueArtifactWithVersions) String() string {
	return Stringify(a)
}

// IssueArtifactListOptions specifies optional parameters to
// IssuesService.ListArtifacts.
type IssueArtifactListOptions struct {
	Type string `url:"type,omitempty"`
	ListOptions
}

// CreateIssueArtifactRequest specifies an artifact to create.
type CreateIssueArtifactRequest struct {
	// Type can be "generic", "plan", "link", or "embedded_html".
	Type *string `json:"type,omitempty"`
	Name string  `json:"name"`
	Body string  `json:"body"`
}

// UpdateIssueArtifactRequest specifies changes to an artifact.
// Nil fields preserve existing values; non-nil pointers send even empty strings.
// An artifact's type cannot be changed.
type UpdateIssueArtifactRequest struct {
	Name *string `json:"name,omitempty"`
	Body *string `json:"body,omitempty"`
}

// ListArtifacts lists artifacts attached to an issue.
//
// The response uses standard Link headers for pagination, which are exposed
// through Response and supported by ListArtifactsIter. Offset pagination has
// no snapshot guarantee.
//
// Note: ListArtifacts uses the undocumented GitHub API endpoint "GET /repos/{owner}/{repo}/issues/{issue_number}/artifacts".
//
//meta:operation GET /repos/{owner}/{repo}/issues/{issue_number}/artifacts
func (s *IssuesService) ListArtifacts(ctx context.Context, owner, repo string, issueNumber int, opts *IssueArtifactListOptions) ([]*IssueArtifact, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/artifacts", owner, repo, issueNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var artifacts []*IssueArtifact
	resp, err := s.client.Do(req, &artifacts)
	if err != nil {
		return nil, resp, err
	}

	return artifacts, resp, nil
}

// GetArtifact gets an artifact and its edit history using its issue-scoped number.
//
// Note: GetArtifact uses the undocumented GitHub API endpoint "GET /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}".
//
//meta:operation GET /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}
func (s *IssuesService) GetArtifact(ctx context.Context, owner, repo string, issueNumber, artifactNumber int) (*IssueArtifactWithVersions, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/artifacts/%v", owner, repo, issueNumber, artifactNumber)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var artifact *IssueArtifactWithVersions
	resp, err := s.client.Do(req, &artifact)
	if err != nil {
		return nil, resp, err
	}

	return artifact, resp, nil
}

// CreateArtifact creates an artifact attached to an issue.
//
// Note: CreateArtifact uses the undocumented GitHub API endpoint "POST /repos/{owner}/{repo}/issues/{issue_number}/artifacts".
//
//meta:operation POST /repos/{owner}/{repo}/issues/{issue_number}/artifacts
func (s *IssuesService) CreateArtifact(ctx context.Context, owner, repo string, issueNumber int, body CreateIssueArtifactRequest) (*IssueArtifact, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/artifacts", owner, repo, issueNumber)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var artifact *IssueArtifact
	resp, err := s.client.Do(req, &artifact)
	if err != nil {
		return nil, resp, err
	}

	return artifact, resp, nil
}

// UpdateArtifact updates an artifact using its issue-scoped number.
//
// Note: UpdateArtifact uses the undocumented GitHub API endpoint "PATCH /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}".
//
//meta:operation PATCH /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}
func (s *IssuesService) UpdateArtifact(ctx context.Context, owner, repo string, issueNumber, artifactNumber int, body UpdateIssueArtifactRequest) (*IssueArtifact, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/artifacts/%v", owner, repo, issueNumber, artifactNumber)
	req, err := s.client.NewRequest(ctx, "PATCH", u, body)
	if err != nil {
		return nil, nil, err
	}

	var artifact *IssueArtifact
	resp, err := s.client.Do(req, &artifact)
	if err != nil {
		return nil, resp, err
	}

	return artifact, resp, nil
}

// DeleteArtifact deletes an artifact using its issue-scoped number.
//
// Note: DeleteArtifact uses the undocumented GitHub API endpoint "DELETE /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}".
//
//meta:operation DELETE /repos/{owner}/{repo}/issues/{issue_number}/artifacts/{artifact_number}
func (s *IssuesService) DeleteArtifact(ctx context.Context, owner, repo string, issueNumber, artifactNumber int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/artifacts/%v", owner, repo, issueNumber, artifactNumber)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
