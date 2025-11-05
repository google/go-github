// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CostCenter represents an enterprise cost center.
type CostCenter struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	Resources         []*CostCenterResource `json:"resources"`
	State             *string               `json:"state,omitempty"`
	AzureSubscription *string               `json:"azure_subscription,omitempty"`
}

// CostCenterResource represents a resource assigned to a cost center.
type CostCenterResource struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// CostCenters represents a list of cost centers.
type CostCenters struct {
	CostCenters []*CostCenter `json:"costCenters,omitempty"`
}

// CostCenterListOptions specifies optional parameters to the EnterpriseService.ListCostCenters method.
type CostCenterListOptions struct {
	State *string `url:"state,omitempty"`
}

// CostCenterRequest represents a request to create or update a cost center.
type CostCenterRequest struct {
	Name *string `json:"name,omitempty"`
}

// CostCenterResourceRequest represents a request to add or remove resources from a cost center.
type CostCenterResourceRequest struct {
	Users         []string `json:"users,omitempty"`
	Organizations []string `json:"organizations,omitempty"`
	Repositories  []string `json:"repositories,omitempty"`
}

// CostCenterAddResourceResponse represents a response from adding resources to a cost center.
type CostCenterAddResourceResponse struct {
	Message             *string               `json:"message,omitempty"`
	ReassignedResources []*ReassignedResource `json:"reassigned_resources,omitempty"`
}

// ReassignedResource represents a resource that was reassigned from another cost center.
type ReassignedResource struct {
	ResourceType       *string `json:"resource_type,omitempty"`
	Name               *string `json:"name,omitempty"`
	PreviousCostCenter *string `json:"previous_cost_center,omitempty"`
}

// CostCenterRemoveResourceResponse represents a response from removing resources from a cost center.
type CostCenterRemoveResourceResponse struct {
	Message *string `json:"message,omitempty"`
}

// CostCenterDeleteResponse represents a response from deleting a cost center.
type CostCenterDeleteResponse struct {
	Message         *string `json:"message,omitempty"`
	ID              *string `json:"id,omitempty"`
	Name            *string `json:"name,omitempty"`
	CostCenterState *string `json:"costCenterState,omitempty"`
}

// ListCostCenters lists all cost centers for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#get-all-cost-centers-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/cost-centers
func (s *EnterpriseService) ListCostCenters(ctx context.Context, enterprise string, opts *CostCenterListOptions) (*CostCenters, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers", enterprise)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	costCenters := &CostCenters{}
	resp, err := s.client.Do(ctx, req, costCenters)
	if err != nil {
		return nil, resp, err
	}

	return costCenters, resp, nil
}

// CreateCostCenter creates a new cost center for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#create-a-new-cost-center
//
//meta:operation POST /enterprises/{enterprise}/settings/billing/cost-centers
func (s *EnterpriseService) CreateCostCenter(ctx context.Context, enterprise string, costCenter CostCenterRequest) (*CostCenter, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers", enterprise)

	req, err := s.client.NewRequest("POST", u, costCenter)
	if err != nil {
		return nil, nil, err
	}

	result := &CostCenter{}
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetCostCenter gets a cost center by ID for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#get-a-cost-center-by-id
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/cost-centers/{cost_center_id}
func (s *EnterpriseService) GetCostCenter(ctx context.Context, enterprise, costCenterID string) (*CostCenter, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers/%v", enterprise, costCenterID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	costCenter := &CostCenter{}
	resp, err := s.client.Do(ctx, req, costCenter)
	if err != nil {
		return nil, resp, err
	}

	return costCenter, resp, nil
}

// UpdateCostCenter updates the name of a cost center.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#update-a-cost-center-name
//
//meta:operation PATCH /enterprises/{enterprise}/settings/billing/cost-centers/{cost_center_id}
func (s *EnterpriseService) UpdateCostCenter(ctx context.Context, enterprise, costCenterID string, costCenter CostCenterRequest) (*CostCenter, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers/%v", enterprise, costCenterID)

	req, err := s.client.NewRequest("PATCH", u, costCenter)
	if err != nil {
		return nil, nil, err
	}

	result := &CostCenter{}
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteCostCenter deletes a cost center.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#delete-a-cost-center
//
//meta:operation DELETE /enterprises/{enterprise}/settings/billing/cost-centers/{cost_center_id}
func (s *EnterpriseService) DeleteCostCenter(ctx context.Context, enterprise, costCenterID string) (*CostCenterDeleteResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers/%v", enterprise, costCenterID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := &CostCenterDeleteResponse{}
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AddResourcesToCostCenter adds resources to a cost center.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#add-resources-to-a-cost-center
//
//meta:operation POST /enterprises/{enterprise}/settings/billing/cost-centers/{cost_center_id}/resource
func (s *EnterpriseService) AddResourcesToCostCenter(ctx context.Context, enterprise, costCenterID string, resources CostCenterResourceRequest) (*CostCenterAddResourceResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers/%v/resource", enterprise, costCenterID)

	req, err := s.client.NewRequest("POST", u, resources)
	if err != nil {
		return nil, nil, err
	}

	result := &CostCenterAddResourceResponse{}
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RemoveResourcesFromCostCenter removes resources from a cost center.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/billing#remove-resources-from-a-cost-center
//
//meta:operation DELETE /enterprises/{enterprise}/settings/billing/cost-centers/{cost_center_id}/resource
func (s *EnterpriseService) RemoveResourcesFromCostCenter(ctx context.Context, enterprise, costCenterID string, resources CostCenterResourceRequest) (*CostCenterRemoveResourceResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/cost-centers/%v/resource", enterprise, costCenterID)

	req, err := s.client.NewRequest("DELETE", u, resources)
	if err != nil {
		return nil, nil, err
	}

	result := &CostCenterRemoveResourceResponse{}
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
