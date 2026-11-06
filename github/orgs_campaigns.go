// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// Campaign represents a security campaign for an organization.
type Campaign struct {
	Number       int64               `json:"number"`
	CreatedAt    Timestamp           `json:"created_at"`
	UpdatedAt    Timestamp           `json:"updated_at"`
	Name         *string             `json:"name,omitempty"`
	Description  string              `json:"description"`
	Managers     []*User             `json:"managers"`
	TeamManagers []*Team             `json:"team_managers,omitzero"`
	PublishedAt  *Timestamp          `json:"published_at,omitempty"`
	EndsAt       Timestamp           `json:"ends_at"`
	ClosedAt     *Timestamp          `json:"closed_at,omitempty"`
	State        string              `json:"state"`
	ContactLink  string              `json:"contact_link"`
	AlertStats   *CampaignAlertStats `json:"alert_stats,omitempty"`
}

// CampaignAlertStats represents the alert counts of a campaign.
type CampaignAlertStats struct {
	OpenCount       int `json:"open_count"`
	ClosedCount     int `json:"closed_count"`
	InProgressCount int `json:"in_progress_count"`
}

// CampaignCodeScanningAlert represents the code scanning alerts of a repository to include in a campaign.
type CampaignCodeScanningAlert struct {
	RepositoryID int64   `json:"repository_id"`
	AlertNumbers []int64 `json:"alert_numbers"`
}

// CreateCampaignRequest represents the body of a request to create a campaign.
type CreateCampaignRequest struct {
	Name               string                       `json:"name"`
	Description        string                       `json:"description"`
	Managers           []string                     `json:"managers,omitzero"`
	TeamManagers       []string                     `json:"team_managers,omitzero"`
	EndsAt             Timestamp                    `json:"ends_at"`
	ContactLink        *string                      `json:"contact_link,omitempty"`
	CodeScanningAlerts []*CampaignCodeScanningAlert `json:"code_scanning_alerts,omitzero"`
	GenerateIssues     *bool                        `json:"generate_issues,omitempty"`
}

// UpdateCampaignRequest represents the body of a request to update a campaign.
type UpdateCampaignRequest struct {
	Name         *string    `json:"name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Managers     []string   `json:"managers,omitzero"`
	TeamManagers []string   `json:"team_managers,omitzero"`
	EndsAt       *Timestamp `json:"ends_at,omitempty"`
	ContactLink  *string    `json:"contact_link,omitempty"`
	State        *string    `json:"state,omitempty"`
}

// ListCampaignsOptions specifies optional parameters to the
// OrganizationsService.ListCampaigns method.
type ListCampaignsOptions struct {
	// State filters campaigns by the state they are in. Possible values are: open, closed.
	State string `url:"state,omitempty"`

	// Sort determines how the campaigns are sorted. Possible values are:
	// created, updated, ends_at, published. Default is created.
	Sort string `url:"sort,omitempty"`

	// Direction in which to sort the campaigns. Possible values are: asc, desc.
	Direction string `url:"direction,omitempty"`

	ListOptions
}

// ListCampaigns lists the campaigns of an organization.
//
// GitHub API docs: https://docs.github.com/rest/campaigns/campaigns?apiVersion=2022-11-28#list-campaigns-for-an-organization
//
//meta:operation GET /orgs/{org}/campaigns
func (s *OrganizationsService) ListCampaigns(ctx context.Context, org string, opts *ListCampaignsOptions) ([]*Campaign, *Response, error) {
	u := fmt.Sprintf("orgs/%v/campaigns", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var campaigns []*Campaign
	resp, err := s.client.Do(req, &campaigns)
	if err != nil {
		return nil, resp, err
	}
	return campaigns, resp, nil
}

// CreateCampaign creates a campaign for an organization.
//
// GitHub API docs: https://docs.github.com/rest/campaigns/campaigns?apiVersion=2022-11-28#create-a-campaign-for-an-organization
//
//meta:operation POST /orgs/{org}/campaigns
func (s *OrganizationsService) CreateCampaign(ctx context.Context, org string, body CreateCampaignRequest) (*Campaign, *Response, error) {
	u := fmt.Sprintf("orgs/%v/campaigns", org)

	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var campaign *Campaign
	resp, err := s.client.Do(req, &campaign)
	if err != nil {
		return nil, resp, err
	}
	return campaign, resp, nil
}

// GetCampaign gets a campaign of an organization.
//
// GitHub API docs: https://docs.github.com/rest/campaigns/campaigns?apiVersion=2022-11-28#get-a-campaign-for-an-organization
//
//meta:operation GET /orgs/{org}/campaigns/{campaign_number}
func (s *OrganizationsService) GetCampaign(ctx context.Context, org string, campaignNumber int64) (*Campaign, *Response, error) {
	u := fmt.Sprintf("orgs/%v/campaigns/%v", org, campaignNumber)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var campaign *Campaign
	resp, err := s.client.Do(req, &campaign)
	if err != nil {
		return nil, resp, err
	}
	return campaign, resp, nil
}

// UpdateCampaign updates a campaign of an organization.
//
// GitHub API docs: https://docs.github.com/rest/campaigns/campaigns?apiVersion=2022-11-28#update-a-campaign
//
//meta:operation PATCH /orgs/{org}/campaigns/{campaign_number}
func (s *OrganizationsService) UpdateCampaign(ctx context.Context, org string, campaignNumber int64, body UpdateCampaignRequest) (*Campaign, *Response, error) {
	u := fmt.Sprintf("orgs/%v/campaigns/%v", org, campaignNumber)

	req, err := s.client.NewRequest(ctx, "PATCH", u, body)
	if err != nil {
		return nil, nil, err
	}

	var campaign *Campaign
	resp, err := s.client.Do(req, &campaign)
	if err != nil {
		return nil, resp, err
	}
	return campaign, resp, nil
}

// DeleteCampaign deletes a campaign of an organization.
//
// GitHub API docs: https://docs.github.com/rest/campaigns/campaigns?apiVersion=2022-11-28#delete-a-campaign-for-an-organization
//
//meta:operation DELETE /orgs/{org}/campaigns/{campaign_number}
func (s *OrganizationsService) DeleteCampaign(ctx context.Context, org string, campaignNumber int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/campaigns/%v", org, campaignNumber)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
