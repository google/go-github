// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// DependencyGraphSnapshotResolvedDependencyRelationship represents whether the dependency is requested directly by the manifest or is a dependency of another dependency.
//
// Can have the following values:
//   - "direct": indicates that the dependency is requested directly by the manifest.
//   - "indirect": indicates that the dependency is a dependency of another dependency.
type DependencyGraphSnapshotResolvedDependencyRelationship string

// DependencyGraphSnapshotResolvedDependencyScope represents whether the dependency is required for the primary build artifact or is only used for development.
//
// Can have the following values:
//   - "runtime": indicates that the dependency is required for the primary build artifact.
//   - "development": indicates that the dependency is only used for development.
type DependencyGraphSnapshotResolvedDependencyScope string

// DependencyGraphSnapshotCreationResult represents the snapshot creation result.
//
// Can have the following values:
//   - "SUCCESS": indicates that the snapshot was successfully created and the repository's dependencies were updated.
//   - "ACCEPTED": indicates that the snapshot was successfully created, but the repository's dependencies were not updated.
//   - "INVALID": indicates that the snapshot was malformed.
type DependencyGraphSnapshotCreationResult string

type DependencyGraphSnapshotResolvedDependency struct {
	PackageUrl   *string                                               `json:"package_url,omitempty"`
	Relationship DependencyGraphSnapshotResolvedDependencyRelationship `json:"relationship,omitempty"`
	Scope        DependencyGraphSnapshotResolvedDependencyScope        `json:"scope,omitempty"`
	Dependencies []string                                              `json:"dependencies,omitempty"`
}

type DependencyGraphSnapshotJob struct {
	Correlator *string `json:"correlator,omitempty"`
	ID         *string `json:"id,omitempty"`
	HtmlUrl    *string `json:"html_url,omitempty"`
}

type DependencyGraphSnapshotDetector struct {
	Name    *string `json:"name,omitempty"`
	Version *string `json:"version,omitempty"`
	URL     *string `json:"url,omitempty"`
}

type DependencyGraphSnapshotManifestFile struct {
	SourceLocation *string `json:"source_location,omitempty"`
}

type DependencyGraphSnapshotManifest struct {
	Name     *string                                               `json:"name,omitempty"`
	File     *DependencyGraphSnapshotManifestFile                  `json:"file,omitempty"`
	Resolved map[string]*DependencyGraphSnapshotResolvedDependency `json:"resolved,omitempty"`
}

type DependencyGraphSnapshot struct {
	Version   int                                         `json:"version"`
	Sha       *string                                     `json:"sha,omitempty"`
	Ref       *string                                     `json:"ref,omitempty"`
	Job       *DependencyGraphSnapshotJob                 `json:"job,omitempty"`
	Detector  *DependencyGraphSnapshotDetector            `json:"detector,omitempty"`
	Scanned   *Timestamp                                  `json:"scanned,omitempty"`
	Manifests map[string]*DependencyGraphSnapshotManifest `json:"manifests,omitempty"`
}

type DependencyGraphSnapshotCreationData struct {
	ID        int                                   `json:"id"`
	CreatedAt *Timestamp                            `json:"created_at"`
	Message   *string                               `json:"message"`
	Result    DependencyGraphSnapshotCreationResult `json:"result"`
}

// CreateSnapshot creates a new snapshot of a repository's dependencies.
//
// GitHub API docs: https://docs.github.com/en/rest/dependency-graph/dependency-submission#create-a-snapshot-of-dependencies-for-a-repository
func (s *DependencyGraphService) CreateSnapshot(ctx context.Context, owner, repo string, dependencyGraphSnapshot *DependencyGraphSnapshot) (*DependencyGraphSnapshotCreationData, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependency-graph/snapshots", owner, repo)

	req, err := s.client.NewRequest("POST", url, dependencyGraphSnapshot)
	if err != nil {
		return nil, nil, err
	}

	var snapshotCreationData *DependencyGraphSnapshotCreationData
	resp, err := s.client.Do(ctx, req, &snapshotCreationData)
	if err != nil {
		return nil, resp, err
	}

	return snapshotCreationData, resp, nil
}
