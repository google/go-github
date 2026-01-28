// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ArtifactDeploymentRecord represents a GitHub artifact deployment record.
type ArtifactDeploymentRecord struct {
	ID                  *int64            `json:"id,omitempty"`
	Digest              *string           `json:"digest,omitempty"`
	Name                *string           `json:"name,omitempty"`
	Version             *string           `json:"version,omitempty"`
	Status              *string           `json:"status,omitempty"`
	LogicalEnvironment  *string           `json:"logical_environment,omitempty"`
	PhysicalEnvironment *string           `json:"physical_environment,omitempty"`
	Cluster             *string           `json:"cluster,omitempty"`
	DeploymentName      *string           `json:"deployment_name,omitempty"`
	Tags                map[string]string `json:"tags,omitempty"`
	RuntimeRisks        []string          `json:"runtime_risks,omitempty"`
	GithubRepository    *string           `json:"github_repository,omitempty"`
	AttestationID       *int64            `json:"attestation_id,omitempty"`
	CreatedAt           *Timestamp        `json:"created_at,omitempty"`
	UpdatedAt           *Timestamp        `json:"updated_at,omitempty"`
}

func (r ArtifactDeploymentRecord) String() string { return Stringify(r) }

// ArtifactDeploymentResponse represents the response for deployment records.
type ArtifactDeploymentResponse struct {
	TotalCount        *int                        `json:"total_count,omitempty"`
	DeploymentRecords []*ArtifactDeploymentRecord `json:"deployment_records,omitempty"`
}

func (r ArtifactDeploymentResponse) String() string { return Stringify(r) }

// ClusterDeploymentRecordsRequest represents the request body for setting cluster deployment records.
type ClusterDeploymentRecordsRequest struct {
	LogicalEnvironment  *string                     `json:"logical_environment,omitempty"`
	PhysicalEnvironment *string                     `json:"physical_environment,omitempty"`
	Deployments         []*ArtifactDeploymentRecord `json:"deployments,omitempty"`
}

func (r ClusterDeploymentRecordsRequest) String() string { return Stringify(r) }

// ArtifactStorageRecord represents a GitHub artifact storage record.
type ArtifactStorageRecord struct {
	ID               *int64     `json:"id,omitempty"`
	Name             *string    `json:"name,omitempty"`
	Digest           *string    `json:"digest,omitempty"`
	Version          *string    `json:"version,omitempty"`
	ArtifactURL      *string    `json:"artifact_url,omitempty"`
	Path             *string    `json:"path,omitempty"`
	RegistryURL      *string    `json:"registry_url,omitempty"`
	Repository       *string    `json:"repository,omitempty"`
	Status           *string    `json:"status,omitempty"`
	GithubRepository *string    `json:"github_repository,omitempty"`
	CreatedAt        *Timestamp `json:"created_at,omitempty"`
	UpdatedAt        *Timestamp `json:"updated_at,omitempty"`
}

func (r ArtifactStorageRecord) String() string { return Stringify(r) }

// ArtifactStorageResponse represents the response for storage records.
type ArtifactStorageResponse struct {
	TotalCount     *int                     `json:"total_count,omitempty"`
	StorageRecords []*ArtifactStorageRecord `json:"storage_records,omitempty"`
}

func (r ArtifactStorageResponse) String() string { return Stringify(r) }

// CreateArtifactDeploymentRecord creates an artifact deployment record for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/artifact-metadata#create-an-artifact-deployment-record
//
//meta:operation POST /orgs/{org}/artifacts/metadata/deployment-record
func (s *OrganizationsService) CreateArtifactDeploymentRecord(ctx context.Context, org string, record *ArtifactDeploymentRecord) (*ArtifactDeploymentResponse, *Response, error) {
	u := fmt.Sprintf("orgs/%v/artifacts/metadata/deployment-record", org)
	req, err := s.client.NewRequest("POST", u, record)
	if err != nil {
		return nil, nil, err
	}
	v := new(ArtifactDeploymentResponse)
	resp, err := s.client.Do(ctx, req, v)
	return v, resp, err
}

// SetClusterDeploymentRecords sets deployment records for a given cluster.
//
// GitHub API docs: https://docs.github.com/rest/orgs/artifact-metadata#set-cluster-deployment-records
//
//meta:operation POST /orgs/{org}/artifacts/metadata/deployment-record/cluster/{cluster}
func (s *OrganizationsService) SetClusterDeploymentRecords(ctx context.Context, org, cluster string, request *ClusterDeploymentRecordsRequest) (*ArtifactDeploymentResponse, *Response, error) {
	u := fmt.Sprintf("orgs/%v/artifacts/metadata/deployment-record/cluster/%v", org, cluster)
	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}
	v := new(ArtifactDeploymentResponse)
	resp, err := s.client.Do(ctx, req, v)
	return v, resp, err
}

// CreateArtifactStorageRecord creates metadata storage records for artifacts.
//
// GitHub API docs: https://docs.github.com/rest/orgs/artifact-metadata#create-artifact-metadata-storage-record
//
//meta:operation POST /orgs/{org}/artifacts/metadata/storage-record
func (s *OrganizationsService) CreateArtifactStorageRecord(ctx context.Context, org string, record *ArtifactStorageRecord) (*ArtifactStorageResponse, *Response, error) {
	u := fmt.Sprintf("orgs/%v/artifacts/metadata/storage-record", org)
	req, err := s.client.NewRequest("POST", u, record)
	if err != nil {
		return nil, nil, err
	}
	v := new(ArtifactStorageResponse)
	resp, err := s.client.Do(ctx, req, v)
	return v, resp, err
}

// ListArtifactDeploymentRecords lists deployment records for an artifact metadata.
//
// GitHub API docs: https://docs.github.com/rest/orgs/artifact-metadata#list-artifact-deployment-records
//
//meta:operation GET /orgs/{org}/artifacts/{subject_digest}/metadata/deployment-records
func (s *OrganizationsService) ListArtifactDeploymentRecords(ctx context.Context, org, subjectDigest string) (*ArtifactDeploymentResponse, *Response, error) {
	u := fmt.Sprintf("orgs/%v/artifacts/%v/metadata/deployment-records", org, subjectDigest)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(ArtifactDeploymentResponse)
	resp, err := s.client.Do(ctx, req, v)
	return v, resp, err
}

// ListArtifactStorageRecords lists artifact storage records with a given subject digest.
//
// GitHub API docs: https://docs.github.com/rest/orgs/artifact-metadata#list-artifact-storage-records
//
//meta:operation GET /orgs/{org}/artifacts/{subject_digest}/metadata/storage-records
func (s *OrganizationsService) ListArtifactStorageRecords(ctx context.Context, org, subjectDigest string) (*ArtifactStorageResponse, *Response, error) {
	u := fmt.Sprintf("orgs/%v/artifacts/%v/metadata/storage-records", org, subjectDigest)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(ArtifactStorageResponse)
	resp, err := s.client.Do(ctx, req, v)
	return v, resp, err
}
