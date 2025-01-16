// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// NodeQueryOptions specifies the optional parameters to the EnterpriseService
// Node management APIS.
type NodeQueryOptions struct {
	// UUID filters issues based on the node UUID.
	UUID string `url:"uuid,omitempty"`

	// ClusterRoles filters The cluster roles from the cluster configuration file.
	ClusterRoles string `url:"cluster_roles,omitempty"`
}

// ClusterStatus represents a response from the GetClusterStatus and GetReplicationStatus method.
type ClusterStatus struct {
	Status *string               `json:"status"`
	Nodes  []*ClusterStatusNodes `json:"nodes"`
}

// ClusterStatusNodes represents the status of a cluster node.
type ClusterStatusNodes struct {
	Hostname *string                       `json:"hostname"`
	Status   *string                       `json:"status"`
	Services []*ClusterStatusNodesServices `json:"services,omitempty"`
}

// ClusterStatusNodesServices represents the status of a service running on a cluster node.
type ClusterStatusNodesServices struct {
	Status  *string `json:"status"`
	Name    *string `json:"name"`
	Details *string `json:"details"`
}

// SystemRequirements represents a response from the GetCheckSystemRequirements method.
type SystemRequirements struct {
	Status *string                    `json:"status"`
	Nodes  []*SystemRequirementsNodes `json:"nodes"`
}

// SystemRequirementsNodes represents the status of a system node.
type SystemRequirementsNodes struct {
	Hostname    *string                               `json:"hostname"`
	Status      *string                               `json:"status"`
	RolesStatus []*SystemRequirementsNodesRolesStatus `json:"roles_status"`
}

// SystemRequirementsNodesRolesStatus represents the status of a role on a system node.
type SystemRequirementsNodesRolesStatus struct {
	Status *string `json:"status"`
	Role   *string `json:"role"`
}

// NodeReleaseVersions represents a response from the GetReplicationStatus method.
type NodeReleaseVersions struct {
	Hostname *string          `json:"hostname"`
	Version  *ReleaseVersions `json:"version"`
}

// ReleaseVersions holds the release version information of the node.
type ReleaseVersions struct {
	Version   *string `json:"version"`
	Platform  *string `json:"platform"`
	BuildID   *string `json:"build_id"`
	BuildDate *string `json:"build_date"`
}

// CheckSystemRequirements checks if GHES system nodes meet the system requirements.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#get-the-system-requirement-check-results-for-configured-cluster-nodes
//
//meta:operation GET /manage/v1/checks/system-requirements
func (s *EnterpriseService) CheckSystemRequirements(ctx context.Context) (*SystemRequirements, *Response, error) {
	u := "manage/v1/checks/system-requirements"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	systemRequirements := new(SystemRequirements)
	resp, err := s.client.Do(ctx, req, systemRequirements)
	if err != nil {
		return nil, resp, err
	}

	return systemRequirements, resp, nil
}

// ClusterStatus gets the status of all services running on each cluster node.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#get-the-status-of-services-running-on-all-cluster-nodes
//
//meta:operation GET /manage/v1/cluster/status
func (s *EnterpriseService) ClusterStatus(ctx context.Context) (*ClusterStatus, *Response, error) {
	u := "manage/v1/cluster/status"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	clusterStatus := new(ClusterStatus)
	resp, err := s.client.Do(ctx, req, clusterStatus)
	if err != nil {
		return nil, resp, err
	}

	return clusterStatus, resp, nil
}

// ReplicationStatus gets the status of all services running on each replica node.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#get-the-status-of-services-running-on-all-replica-nodes
//
//meta:operation GET /manage/v1/replication/status
func (s *EnterpriseService) ReplicationStatus(ctx context.Context, opts *NodeQueryOptions) (*ClusterStatus, *Response, error) {
	u, err := addOptions("manage/v1/replication/status", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	replicationStatus := new(ClusterStatus)
	resp, err := s.client.Do(ctx, req, replicationStatus)
	if err != nil {
		return nil, resp, err
	}

	return replicationStatus, resp, nil
}

// Versions gets the versions information deployed to each node.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#get-all-ghes-release-versions-for-all-nodes
//
//meta:operation GET /manage/v1/version
func (s *EnterpriseService) Versions(ctx context.Context, opts *NodeQueryOptions) ([]*NodeReleaseVersions, *Response, error) {
	u, err := addOptions("manage/v1/version", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var releaseVersions []*NodeReleaseVersions
	resp, err := s.client.Do(ctx, req, &releaseVersions)
	if err != nil {
		return nil, resp, err
	}

	return releaseVersions, resp, nil
}
