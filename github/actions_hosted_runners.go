// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// HostedRunnerPublicIP represents the details of a public IP for github-hosted runner.
type HostedRunnerPublicIP struct {
	Enabled bool   `json:"enabled"`
	Prefix  string `json:"prefix"`
	Length  int    `json:"length"`
}

// HostedRunnerMachineSpec represents the details of a particular machine specification for github-hosted runner.
type HostedRunnerMachineSpec struct {
	ID        string `json:"id"`
	CPUCores  int    `json:"cpu_cores"`
	MemoryGB  int    `json:"memory_gb"`
	StorageGB int    `json:"storage_gb"`
}

// HostedRunner represents a single github-hosted runner with additional details.
type HostedRunner struct {
	ID                 *int64                   `json:"id,omitempty"`
	Name               *string                  `json:"name,omitempty"`
	RunnerGroupID      *int64                   `json:"runner_group_id,omitempty"`
	Platform           *string                  `json:"platform,omitempty"`
	Image              *HostedRunnerImageDetail `json:"image,omitempty"`
	MachineSizeDetails *HostedRunnerMachineSpec `json:"machine_size_details,omitempty"`
	Status             *string                  `json:"status,omitempty"`
	MaximumRunners     *int64                   `json:"maximum_runners,omitempty"`
	PublicIPEnabled    *bool                    `json:"public_ip_enabled,omitempty"`
	PublicIPs          []*HostedRunnerPublicIP  `json:"public_ips,omitempty"`
	LastActiveOn       *string                  `json:"last_active_on,omitempty"`
}

// HostedRunnerImageDetail represents the image details of a github-hosted runners.
type HostedRunnerImageDetail struct {
	ID   *string `json:"id,omitempty"`
	Size *int    `json:"size,omitempty"`
}

// HostedRunners represents a collection of github-hosted runners for an organization.
type HostedRunners struct {
	TotalCount int             `json:"total_count"`
	Runners    []*HostedRunner `json:"runners"`
}

// ListHostedRunners lists all the github-hosted runners for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#list-github-hosted-runners-for-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners
func (s *ActionsService) ListHostedRunners(ctx context.Context, org string, opts *ListOptions) (*HostedRunners, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	runners := &HostedRunners{}
	resp, err := s.client.Do(ctx, req, &runners)
	if err != nil {
		return nil, resp, err
	}

	return runners, resp, nil
}

// HostedRunnerImage represents the image of github-hosted runners
// To list all available images, use GET /actions/hosted-runners/images/github-owned or GET /actions/hosted-runners/images/partner
type HostedRunnerImage struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Version string `json:"version"`
}

// CreateHostedRunnerRequest specifies body parameters to Hosted Runner configuration.
type CreateHostedRunnerRequest struct {
	Name           string            `json:"name"`
	Image          HostedRunnerImage `json:"image"`
	RunnerGroupID  int64             `json:"runner_group_id"`
	Size           string            `json:"size"`
	MaximumRunners int64             `json:"maximum_runners,omitempty"`
	EnableStaticIP bool              `json:"enable_static_ip,omitempty"`
}

// CreateHostedRunner creates a github-hosted runner for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#create-a-github-hosted-runner-for-an-organization
//
//meta:operation POST /orgs/{org}/actions/hosted-runners
func (s *ActionsService) CreateHostedRunner(ctx context.Context, org string, request *CreateHostedRunnerRequest) (*HostedRunner, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners", org)
	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}

	hostedRunner := new(HostedRunner)
	resp, err := s.client.Do(ctx, req, hostedRunner)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunner, resp, nil
}

// HostedRunnerImageSpecs represents the details of a github-hosted runner image.
type HostedRunnerImageSpecs struct {
	ID          string `json:"id"`
	Platform    string `json:"platform"`
	SizeGB      int    `json:"size_gb"`
	DisplayName string `json:"display_name"`
	Source      string `json:"source"`
}

// HostedRunnerImages represents the response containing the total count and details of runner images.
type HostedRunnerImages struct {
	TotalCount int                       `json:"total_count"`
	Images     []*HostedRunnerImageSpecs `json:"images"`
}

// GetHostedRunnerGithubOwnedImages gets the list of GitHub-owned images available for github-hosted runners for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-github-owned-images-for-github-hosted-runners-in-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/images/github-owned
func (s *ActionsService) GetHostedRunnerGithubOwnedImages(ctx context.Context, org string) (*HostedRunnerImages, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/images/github-owned", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	hostedRunnerImages := new(HostedRunnerImages)
	resp, err := s.client.Do(ctx, req, hostedRunnerImages)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunnerImages, resp, nil
}

// GetHostedRunnerPartnerImages gets the list of partner images available for github-hosted runners for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-partner-images-for-github-hosted-runners-in-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/images/partner
func (s *ActionsService) GetHostedRunnerPartnerImages(ctx context.Context, org string) (*HostedRunnerImages, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/images/partner", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	hostedRunnerImages := new(HostedRunnerImages)
	resp, err := s.client.Do(ctx, req, hostedRunnerImages)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunnerImages, resp, nil
}

// HostedRunnerPublicIPLimits represents the static public IP limits for github-hosted runners.
type HostedRunnerPublicIPLimits struct {
	PublicIPs *PublicIPUsage `json:"public_ips"`
}

// PublicIPUsage provides details of static public IP limits for github-hosted runners.
type PublicIPUsage struct {
	Maximum      int `json:"maximum"`
	CurrentUsage int `json:"current_usage"`
}

// GetHostedRunnerLimits gets the github-hosted runners Static public IP Limits for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-limits-on-github-hosted-runners-for-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/limits
func (s *ActionsService) GetHostedRunnerLimits(ctx context.Context, org string) (*HostedRunnerPublicIPLimits, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/limits", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	publicIPLimits := new(HostedRunnerPublicIPLimits)
	resp, err := s.client.Do(ctx, req, publicIPLimits)
	if err != nil {
		return nil, resp, err
	}

	return publicIPLimits, resp, nil
}

// HostedRunnerMachineSpecs represents the response containing the total count and details of machine specs for github-hosted runners.
type HostedRunnerMachineSpecs struct {
	TotalCount   int                        `json:"total_count"`
	MachineSpecs []*HostedRunnerMachineSpec `json:"machine_specs"`
}

// GetHostedRunnerMachineSpecs gets the list of machine specs available for github-hosted runners for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-github-hosted-runners-machine-specs-for-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/machine-sizes
func (s *ActionsService) GetHostedRunnerMachineSpecs(ctx context.Context, org string) (*HostedRunnerMachineSpecs, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/machine-sizes", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	machineSpecs := new(HostedRunnerMachineSpecs)
	resp, err := s.client.Do(ctx, req, machineSpecs)
	if err != nil {
		return nil, resp, err
	}

	return machineSpecs, resp, nil
}

// HostedRunnerPlatforms represents the response containing the total count and platforms for github-hosted runners.
type HostedRunnerPlatforms struct {
	TotalCount int      `json:"total_count"`
	Platforms  []string `json:"platforms"`
}

// GetHostedRunnerPlatforms gets list of platforms available for github-hosted runners for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-platforms-for-github-hosted-runners-in-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/platforms
func (s *ActionsService) GetHostedRunnerPlatforms(ctx context.Context, org string) (*HostedRunnerPlatforms, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/platforms", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	platforms := new(HostedRunnerPlatforms)
	resp, err := s.client.Do(ctx, req, platforms)
	if err != nil {
		return nil, resp, err
	}

	return platforms, resp, nil
}

// GetHostedRunner gets a github-hosted runner in an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#get-a-github-hosted-runner-for-an-organization
//
//meta:operation GET /orgs/{org}/actions/hosted-runners/{hosted_runner_id}
func (s *ActionsService) GetHostedRunner(ctx context.Context, org string, runnerID int64) (*HostedRunner, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/%v", org, runnerID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	hostedRunner := new(HostedRunner)
	resp, err := s.client.Do(ctx, req, hostedRunner)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunner, resp, nil
}

// UpdateHostedRunnerRequest specifies the parameters for updating a runner specifications.
type UpdateHostedRunnerRequest struct {
	Name           string `json:"name"`
	RunnerGroupID  int64  `json:"runner_group_id"`
	MaximumRunners int64  `json:"maximum_runners"`
	EnableStaticIP bool   `json:"enable_static_ip"`
	ImageVersion   string `json:"image_version"`
}

// UpdateHostedRunner updates a github-hosted runner for an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#update-a-github-hosted-runner-for-an-organization
//
//meta:operation PATCH /orgs/{org}/actions/hosted-runners/{hosted_runner_id}
func (s *ActionsService) UpdateHostedRunner(ctx context.Context, org string, runnerID int64, updateReq UpdateHostedRunnerRequest) (*HostedRunner, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/%v", org, runnerID)
	req, err := s.client.NewRequest("PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	hostedRunner := new(HostedRunner)
	resp, err := s.client.Do(ctx, req, hostedRunner)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunner, resp, nil
}

// DeleteHostedRunner deletes github-hosted runner from an organization.
//
// GitHub API docs: https://docs.github.com/rest/actions/hosted-runners#delete-a-github-hosted-runner-for-an-organization
//
//meta:operation DELETE /orgs/{org}/actions/hosted-runners/{hosted_runner_id}
func (s *ActionsService) DeleteHostedRunner(ctx context.Context, org string, runnerID int64) (*HostedRunner, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/hosted-runners/%v", org, runnerID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	hostedRunner := new(HostedRunner)
	resp, err := s.client.Do(ctx, req, hostedRunner)
	if err != nil {
		return nil, resp, err
	}

	return hostedRunner, resp, nil
}
