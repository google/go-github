// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"
)

type DependencyGraphService service

// Sbom represents software bill of materials, which descibes the
// packages/libraries that a repository depends on.
type Sbom struct {
	Sbom *SbomInfo `json:"sbom,omitempty"`
}

// When was the SBOM created and who created it
type CreationInfo struct {
	Created  *time.Time `json:"created,omitempty"`
	Creators []*string  `json:"creators,omitempty"`
}

type RepoDependencies struct {
	Spdxid *string `json:"SPDXID,omitempty"`
	// Package name
	Name             *string `json:"name,omitempty"`
	VersionInfo      *string `json:"versionInfo,omitempty"`
	DownloadLocation *string `json:"downloadLocation,omitempty"`
	FilesAnalyzed    *bool   `json:"filesAnalyzed,omitempty"`
	LicenseConcluded *string `json:"licenseConcluded,omitempty"`
	LicenseDeclared  *string `json:"licenseDeclared,omitempty"`
}

// SPDX is an open standard for software bill of materials (SBOM) that
// identifies and catalogs components, licenses, copyrights, security
// references, and other metadata relating to software
type SbomInfo struct {
	Spdxid       *string       `json:"SPDXID,omitempty"`
	SpdxVersion  *string       `json:"spdxVersion,omitempty"`
	CreationInfo *CreationInfo `json:"creationInfo,omitempty"`

	// Repo name
	Name              *string   `json:"name,omitempty"`
	DataLicense       *string   `json:"dataLicense,omitempty"`
	DocumentDescribes []*string `json:"documentDescribes,omitempty"`
	DocumentNamespace *string   `json:"documentNamespace,omitempty"`

	// List of packages dependencies
	Packages []*RepoDependencies `json:"packages,omitempty"`
}

func (s Sbom) String() string {
	return Stringify(s)
}

// GetSbom fetches the Software bill of materials for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/dependency-graph/sboms
func (s *DependencyGraphService) GetSbom(ctx context.Context, owner string, repo string) (*Sbom, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/dependency-graph/sbom", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var sbom *Sbom
	resp, err := s.client.Do(ctx, req, &sbom)
	if err != nil {
		return nil, resp, err
	}

	return sbom, resp, nil
}
