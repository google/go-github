// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type DependencyGraphService service

// SBOM represents a software bill of materials, which describes the
// packages/libraries that a repository depends on.
type SBOM struct {
	SBOM *SBOMInfo `json:"sbom,omitempty"`
}

// CreationInfo represents when the SBOM was created and who created it.
type CreationInfo struct {
	Created  *Timestamp `json:"created,omitempty"`
	Creators []string   `json:"creators,omitempty"`
}

// RepoDependencies represents the dependencies of a repo.
type RepoDependencies struct {
	SPDXID *string `json:"SPDXID,omitempty"`
	// Package name
	Name             *string        `json:"name,omitempty"`
	VersionInfo      *string        `json:"versionInfo,omitempty"`
	DownloadLocation *string        `json:"downloadLocation,omitempty"`
	FilesAnalyzed    *bool          `json:"filesAnalyzed,omitempty"`
	LicenseConcluded *string        `json:"licenseConcluded,omitempty"`
	LicenseDeclared  *string        `json:"licenseDeclared,omitempty"`
	ExternalRefs     []*ExternalRef `json:"externalRefs"`
}

// ExternalRef represents an external reference (e.g., PURL/SWID/CPE) for a package in the SBOM.
type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceType     string `json:"referenceType"`
	ReferenceLocator  string `json:"referenceLocator"`
}

// Relationship represents a relationship between two packages in the SBOM.
type Relationship struct {
	// Element ID
	SpdxElementId *string `json:"spdxElementId,omitempty"`

	// Related Element ID
	RelatedSpdxElement *string `json:"relatedSpdxElement,omitempty"`

	// Relationship type, e.g., "DEPENDS_ON", "CONTAINS", etc.
	RelationshipType *string `json:"relationshipType,omitempty"`
}

// SBOMInfo represents a software bill of materials (SBOM) using SPDX.
// SPDX is an open standard for SBOMs that
// identifies and catalogs components, licenses, copyrights, security
// references, and other metadata relating to software.
type SBOMInfo struct {
	SPDXID       *string       `json:"SPDXID,omitempty"`
	SPDXVersion  *string       `json:"spdxVersion,omitempty"`
	CreationInfo *CreationInfo `json:"creationInfo,omitempty"`

	// Repo name
	Name              *string  `json:"name,omitempty"`
	DataLicense       *string  `json:"dataLicense,omitempty"`
	DocumentDescribes []string `json:"documentDescribes,omitempty"`
	DocumentNamespace *string  `json:"documentNamespace,omitempty"`

	// List of packages dependencies
	Packages []*RepoDependencies `json:"packages,omitempty"`

	// List of relationships between packages
	Relationships []*Relationship `json:"relationships,omitempty"`
}

func (s SBOM) String() string {
	return Stringify(s)
}

// GetSBOM fetches the software bill of materials for a repository.
//
// GitHub API docs: https://docs.github.com/rest/dependency-graph/sboms#export-a-software-bill-of-materials-sbom-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/dependency-graph/sbom
func (s *DependencyGraphService) GetSBOM(ctx context.Context, owner, repo string) (*SBOM, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/dependency-graph/sbom", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var sbom *SBOM
	resp, err := s.client.Do(ctx, req, &sbom)
	if err != nil {
		return nil, resp, err
	}

	return sbom, resp, nil
}
