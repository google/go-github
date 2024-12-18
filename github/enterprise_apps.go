// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// InstallableOrganization represents a GitHub organization that can be used to install enteprise apps.
type InstallableOrganization struct {
	ID                        *int64  `json:"id,omitempty"`
	Login                     *string `json:"login,omitempty"`
	AccessibleRepositoriesUrl *string `json:"accessible_repositories_url,omitempty"`
}

// ListRunnerApplicationDownloads lists self-hosted runner application binaries that can be downloaded and run.
//
// GitHub API docs: TBD
//
//meta:operation GET /enterprises/{enterprise}/apps/installable_organization
func (s *EnterpriseService) ListOrganizations(ctx context.Context, enterprise string) ([]*InstallableOrganization, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/installable_organization", enterprise)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var orgs []*InstallableOrganization
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, resp, err
	}

	return orgs, resp, nil
}

// ListOrgInstallations list all apps installed on the specified organization.
//
// GitHub API docs: TBD
//
//meta:operation GET /enterprises/{enterprise}/apps/organizations/{org1}/installations
func (s *EnterpriseService) ListOrgInstallations(ctx context.Context, enterprise string, org string) ([]*Installation, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations", enterprise, org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var installations []*Installation
	resp, err := s.client.Do(ctx, req, installations)
	if err != nil {
		return nil, resp, err
	}

	return installations, resp, nil
}

type InstallAppRequest struct {
	// ClientID is the client ID of the GitHub App that should be installed.
	ClientID string
	// RepositorySelection is the type of repository selection requested.
	// Possible values are: "none", "all", "subset".
	RepositorySelection *string `json:"repository_selection"`
}

// InstallApp installs a GitHub App on the specified organization.
//
// GitHub API docs: TBD
//
//meta:operation POST /enterprises/{enterprise}/apps/organizations/{org1}/installation
func (s *EnterpriseService) InstallApp(ctx context.Context, enterprise string, org string, request *InstallAppRequest) (*Installation, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installation", enterprise, org)

	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}

	installation := new(Installation)
	resp, err := s.client.Do(ctx, req, installation)
	if err != nil {
		return nil, resp, err
	}

	return installation, resp, nil
}

// UninstallApp uninstalls the GitHub App for the Enterprise-owned organization
//
// GitHub API docs: TBD
//
//meta:operation DELETE /enterprises/{enterprise}/apps/organizations/{org}/installations/{installation_id}
func (s *EnterpriseService) UninstallApp(ctx context.Context, enterprise string, org string, installationID string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations/%v", enterprise, org, installationID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
