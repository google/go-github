// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ComputeService represents a hosted compute service the network configuration supports.
type ComputeService string

const (
	ComputeServiceNone    ComputeService = "none"
	ComputeServiceActions ComputeService = "actions"
)

// EnterpriseNetworkConfiguration represents a hosted compute network configuration.
type EnterpriseNetworkConfiguration struct {
	ID                 *string         `json:"id,omitempty"`
	Name               *string         `json:"name,omitempty"`
	ComputeService     *ComputeService `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string        `json:"network_settings_ids,omitempty"`
	CreatedOn          *Timestamp      `json:"created_on,omitempty"`
}

// EnterpriseNetworkConfigurations represents a hosted compute network configurations.
type EnterpriseNetworkConfigurations struct {
	TotalCount            *int64                            `json:"total_count,omitempty"`
	NetworkConfigurations []*EnterpriseNetworkConfiguration `json:"network_configurations,omitempty"`
}

// EnterpriseNetworkSettingsResource represents a hosted compute network settings resource.
type EnterpriseNetworkSettingsResource struct {
	ID                     *string `json:"id,omitempty"`
	Name                   *string `json:"name,omitempty"`
	NetworkConfigurationID *string `json:"network_configuration_id,omitempty"`
	SubnetID               *string `json:"subnet_id,omitempty"`
	Region                 *string `json:"region,omitempty"`
}

// EnterpriseNetworkConfigurationRequest represents a request to create or update a network configuration for an enterprise.
type EnterpriseNetworkConfigurationRequest struct {
	Name               *string         `json:"name,omitempty"`
	ComputeService     *ComputeService `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string        `json:"network_settings_ids,omitempty"`
}

// ListEnterpriseNetworkConfigurations lists all hosted compute network configurations configured in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#list-hosted-compute-network-configurations-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/network-configurations
func (s *EnterpriseService) ListEnterpriseNetworkConfigurations(ctx context.Context, enterprise string, opts *ListOptions) (*EnterpriseNetworkConfigurations, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-configurations", enterprise)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	networks := &EnterpriseNetworkConfigurations{}
	resp, err := s.client.Do(ctx, req, networks)
	if err != nil {
		return nil, resp, err
	}
	return networks, resp, nil
}

// CreateEnterpriseNetworkConfiguration creates a hosted compute network configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#create-a-hosted-compute-network-configuration-for-an-enterprise
//
//meta:operation POST /enterprises/{enterprise}/network-configurations
func (s *EnterpriseService) CreateEnterpriseNetworkConfiguration(ctx context.Context, enterprise string, createReq EnterpriseNetworkConfigurationRequest) (*EnterpriseNetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-configurations", enterprise)
	req, err := s.client.NewRequest("POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	network := &EnterpriseNetworkConfiguration{}
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network, resp, nil
}

// GetEnterpriseNetworkConfiguration gets a hosted compute network configuration configured in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#get-a-hosted-compute-network-configuration-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/network-configurations/{network_configuration_id}
func (s *EnterpriseService) GetEnterpriseNetworkConfiguration(ctx context.Context, enterprise, networkID string) (*EnterpriseNetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-configurations/%v", enterprise, networkID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	network := &EnterpriseNetworkConfiguration{}
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}
	return network, resp, nil
}

// UpdateEnterpriseNetworkConfiguration updates a hosted compute network configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#update-a-hosted-compute-network-configuration-for-an-enterprise
//
//meta:operation PATCH /enterprises/{enterprise}/network-configurations/{network_configuration_id}
func (s *EnterpriseService) UpdateEnterpriseNetworkConfiguration(ctx context.Context, enterprise, networkID string, updateReq EnterpriseNetworkConfigurationRequest) (*EnterpriseNetworkConfiguration, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-configurations/%v", enterprise, networkID)
	req, err := s.client.NewRequest("PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	network := &EnterpriseNetworkConfiguration{}
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}
	return network, resp, nil
}

// DeleteEnterpriseNetworkConfiguration deletes a hosted compute network configuration from an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#delete-a-hosted-compute-network-configuration-from-an-enterprise
//
//meta:operation DELETE /enterprises/{enterprise}/network-configurations/{network_configuration_id}
func (s *EnterpriseService) DeleteEnterpriseNetworkConfiguration(ctx context.Context, enterprise, networkID string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-configurations/%v", enterprise, networkID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// GetEnterpriseNetworkSettingsResource gets a hosted compute network settings resource configured for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/network-configurations#get-a-hosted-compute-network-settings-resource-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/network-settings/{network_settings_id}
func (s *EnterpriseService) GetEnterpriseNetworkSettingsResource(ctx context.Context, enterprise, networkID string) (*EnterpriseNetworkSettingsResource, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/network-settings/%v", enterprise, networkID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resource := &EnterpriseNetworkSettingsResource{}
	resp, err := s.client.Do(ctx, req, resource)
	if err != nil {
		return nil, resp, err
	}
	return resource, resp, err
}
