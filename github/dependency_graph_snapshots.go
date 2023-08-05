// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type DependencyRelationship string

const (
	DIRECT   DependencyRelationship = "direct"
	INDIRECT                        = "indirect"
)

type DependencyScope string

const (
	RUNTIME     DependencyScope = "runtime"
	DEVELOPMENT                 = "development"
)

type SnapshotCreationResult string

const (
	SUCCESS  SnapshotCreationResult = "SUCCESS"
	ACCEPTED                        = "ACCEPTED"
	INVALID                         = "INVALID"
)

type Resolved struct {
	PackageUrl   *string                `json:"package_url,omitempty"`
	Relationship DependencyRelationship `json:"relationship,omitempty"`
	Scope        DependencyScope        `json:"scope,omitempty"`
	Dependencies []string               `json:"dependencies,omitempty"`
}

type Job struct {
	Correlator *string `json:"correlator,omitempty"`
	ID         *string `json:"id,omitempty"`
	HtmlUrl    *string `json:"html_url,omitempty"`
}

type Detector struct {
	Name    *string `json:"name,omitempty"`
	Version *string `json:"version,omitempty"`
	URL     *string `json:"url,omitempty"`
}

type File struct {
	SourceLocation *string `json:"source_location,omitempty"`
}

type Manifest struct {
	Name     *string              `json:"name,omitempty"`
	File     *File                `json:"file,omitempty"`
	Resolved map[string]*Resolved `json:"resolved,omitempty"`
}

type Snapshot struct {
	Version   int                  `json:"version"`
	Sha       *string              `json:"sha,omitempty"`
	Ref       *string              `json:"ref,omitempty"`
	Job       *Job                 `json:"job,omitempty"`
	Detector  *Detector            `json:"detector,omitempty"`
	Scanned   *Timestamp           `json:"scanned,omitempty"`
	Manifests map[string]*Manifest `json:"manifests,omitempty"`
}

type SnapshotCreationData struct {
	ID        int                    `json:"id"`
	CreatedAt *Timestamp             `json:"created_at"`
	Message   *string                `json:"message"`
	Result    SnapshotCreationResult `json:"result"`
}

// Create a new snapshot of a repository's dependencies.
//
// GitHub API docs: https://docs.github.com/en/rest/dependency-graph/dependency-submission#create-a-snapshot-of-dependencies-for-a-repository
func (s *DependencyGraphService) CreateSnapshot(ctx context.Context, owner, repo string, snapshot *Snapshot) (*SnapshotCreationData, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependency-graph/snapshots", owner, repo)

	req, err := s.client.NewRequest("POST", url, snapshot)
	if err != nil {
		return nil, nil, err
	}

	var spanshotCreationData *SnapshotCreationData
	resp, err := s.client.Do(ctx, req, &spanshotCreationData)
	if err != nil {
		return nil, resp, err
	}

	return spanshotCreationData, resp, nil
}
