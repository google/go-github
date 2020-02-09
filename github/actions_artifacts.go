// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// Artifact reprents a GitHub artifact.  Artifacts allow sharing
// data between jobs in a workflow and provide storage for data
// once a workflow is complete.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/
type Artifact struct {
	ID                 *int64  `json:"id,omitempty"`
	NodeID             *string `json:"node_id,omitempty"`
	Name               *string `json:"name,omitempty"`
	SizeInBytes        *int64  `json:"size_in_bytes,omitempty"`
	ArchiveDownloadURL *string `json:"archive_download_url,omitempty"`
	Expired            *bool   `json:"expired,omitempty"`
	CreatedAt          *string `json:"created_at,omitempty"`
	ExpiresAt          *string `json:"expires_at,omitempty"`
}

// ArtifactList represents a list of GitHub artifacts.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/
type ArtifactList struct {
	TotalCount *int64      `json:"total_count,omitempty"`
	Artifacts  []*Artifact `json:"artifacts,omitempty"`
}

// ListWorkflowRunArtifacts lists all artifacts that belong to a workflow run.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/#list-workflow-run-artifacts
func (s *ActionsService) ListWorkflowRunArtifacts(ctx context.Context, owner, repo string, runID int64, opt *ListOptions) (*ArtifactList, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/runs/%v/artifacts", owner, repo, runID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	artifactList := new(ArtifactList)
	resp, err := s.client.Do(ctx, req, artifactList)
	if err != nil {
		return nil, resp, err
	}

	return artifactList, resp, nil
}

// GetArtifact gets a specific artifact for a workflow run.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/#get-an-artifact
func (s *ActionsService) GetArtifact(ctx context.Context, owner, repo string, artifactID int64) (*Artifact, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/artifacts/%v", owner, repo, artifactID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	artifact := new(Artifact)
	resp, err := s.client.Do(ctx, req, artifact)
	if err != nil {
		return nil, resp, err
	}

	return artifact, resp, nil
}

// DownloadArtifact gets a redirect URL to download an archive for a repository.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/#download-an-artifact
func (s *ActionsService) DownloadArtifact(ctx context.Context, owner, repo string, artifactID int64, archiveType string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/artifacts/%v/%v", owner, repo, artifactID, archiveType)

	if archiveType != "zip" {
		err := fmt.Errorf("Unsupported archive type %v used, can only use zip", archiveType)
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	location := String("")
	resp, err := s.client.Do(ctx, req, location)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteArtifact deletes a workflow run artifact.
//
// GitHub API docs: https://developer.github.com/v3/actions/artifacts/#delete-an-artifact
func (s *ActionsService) DeleteArtifact(ctx context.Context, owner, repo string, artifactID int64) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/artifacts/%v", owner, repo, artifactID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
