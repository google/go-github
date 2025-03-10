// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// NetworkConfigurations represents a hosted compute network configuration.
type NetworkConfigurations struct {
	TotalCount            int                     `json:"total_count,omitempty"`
	NetworkConfigurations []*NetworkConfiguration `json:"network_configurations,omitempty"`
}

// NetworkConfiguration represents a hosted compute network configurations.
type NetworkConfiguration struct {
	ID                 *string         `json:"id,omitempty"`
	Name               *string         `json:"name,omitempty"`
	ComputeService     *ComputeService `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string        `json:"network_settings_ids,omitempty"`
	CreatedOn          *Timestamp      `json:"created_on"`
}

// NetworkSettingsResource represents a hosted compute network settings resource.
type NetworkSettingsResource struct {
	ID                     *string `json:"id,omitempty"`
	NetworkConfigurationID *string `json:"network_configuration_id,omitempty"`
	Name                   *string `json:"name,omitempty"`
	SubnetID               *string `json:"subnet_id,omitempty"`
	Region                 *string `json:"region,omitempty"`
}

// NetworkConfigurationRequest represents a request to create or update a network configuration for an organization.
type NetworkConfigurationRequest struct {
	Name               *string         `json:"name,omitempty"`
	ComputeService     *ComputeService `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string        `json:"network_settings_ids,omitempty"`
}

// ListNetworkConfigurations lists all hosted compute network configurations configured in an organization.
//
// GitHub API-docs: https://docs.github.com/rest/orgs/network-configurations#list-hosted-compute-network-configurations-for-an-organization
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#list-hosted-compute-network-configurations-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/network-configurations
func (o *OrganizationsService) ListNetworkConfigurations(ctx context.Context, org string, opts *ListOptions) (*NetworkConfigurations, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-configurations", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	configurations := &NetworkConfigurations{}
	resp, err := o.client.Do(ctx, req, configurations)
	if err != nil {
		return nil, resp, err
	}
	return configurations, resp, nil
}

// CreateNetworkConfiguration creates a hosted compute network configuration for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#create-a-hosted-compute-network-configuration-for-an-organization
//
//meta:operation POST /orgs/{org}/settings/network-configurations
func (o *OrganizationsService) CreateNetworkConfiguration(ctx context.Context, org string, createReq NetworkConfigurationRequest) (*NetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-configurations", org)
	req, err := o.client.NewRequest("POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	configuration := &NetworkConfiguration{}
	resp, err := o.client.Do(ctx, req, configuration)
	if err != nil {
		return nil, resp, err
	}
	return configuration, resp, nil
}

// GetNetworkConfiguration gets a hosted compute network configuration configured in an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#get-a-hosted-compute-network-configuration-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/network-configurations/{network_configuration_id}
func (o *OrganizationsService) GetNetworkConfiguration(ctx context.Context, org, networkID string) (*NetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-configurations/%v", org, networkID)
	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	configuration := &NetworkConfiguration{}
	resp, err := o.client.Do(ctx, req, configuration)
	if err != nil {
		return nil, resp, err
	}
	return configuration, resp, nil
}

// UpdateNetworkConfiguration updates a hosted compute network configuration for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#update-a-hosted-compute-network-configuration-for-an-organization
//
//meta:operation PATCH /orgs/{org}/settings/network-configurations/{network_configuration_id}
func (o *OrganizationsService) UpdateNetworkConfiguration(ctx context.Context, org, networkID string, updateReq NetworkConfigurationRequest) (*NetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-configurations/%v", org, networkID)
	req, err := o.client.NewRequest("PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	configuration := &NetworkConfiguration{}
	resp, err := o.client.Do(ctx, req, configuration)
	if err != nil {
		return nil, resp, err
	}
	return configuration, resp, nil
}

// DeleteNetworkConfigurations deletes a hosted compute network configuration from an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#delete-a-hosted-compute-network-configuration-from-an-organization
//
//meta:operation DELETE /orgs/{org}/settings/network-configurations/{network_configuration_id}
func (o *OrganizationsService) DeleteNetworkConfigurations(ctx context.Context, org, networkID string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-configurations/%v", org, networkID)
	req, err := o.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	configuration := &NetworkConfiguration{}
	resp, err := o.client.Do(ctx, req, configuration)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetNetworkConfigurationResource gets a hosted compute network settings resource configured for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/network-configurations#get-a-hosted-compute-network-settings-resource-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/network-settings/{network_settings_id}
func (o *OrganizationsService) GetNetworkConfigurationResource(ctx context.Context, org, networkID string) (*NetworkSettingsResource, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/network-settings/%v", org, networkID)
	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resource := &NetworkSettingsResource{}
	resp, err := o.client.Do(ctx, req, resource)
	if err != nil {
		return nil, resp, err
	}
	return resource, resp, nil
}
