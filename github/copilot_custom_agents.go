// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CopilotCustomAgent represents a custom agent defined in an enterprise source repository.
type CopilotCustomAgent struct {
	Name     *string `json:"name,omitempty"`
	FilePath *string `json:"file_path,omitempty"`
	URL      *string `json:"url,omitempty"`
}

// EnterpriseCustomAgents represents the list of custom agents for an enterprise.
type EnterpriseCustomAgents struct {
	CustomAgents []*CopilotCustomAgent `json:"custom_agents"`
}

// CopilotCustomAgentsSourceOrganization represents the organization configured as
// the custom agents source for an enterprise.
type CopilotCustomAgentsSourceOrganization struct {
	ID        *int64  `json:"id,omitempty"`
	Login     *string `json:"login,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// CopilotCustomAgentsSourceRepository represents the repository that stores custom agent definitions.
type CopilotCustomAgentsSourceRepository struct {
	ID       *int64  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

// CopilotCustomAgentsSourceRuleset represents the ruleset created to protect agent definition files.
type CopilotCustomAgentsSourceRuleset struct {
	ID          *int64  `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Enforcement *string `json:"enforcement,omitempty"`
}

// CopilotCustomAgentsSource represents the source organization and repository for
// enterprise custom agents (and optionally a protecting ruleset).
type CopilotCustomAgentsSource struct {
	Organization *CopilotCustomAgentsSourceOrganization `json:"organization"`
	Repository   *CopilotCustomAgentsSourceRepository   `json:"repository"`
	Ruleset      *CopilotCustomAgentsSourceRuleset      `json:"ruleset,omitempty"`
}

// SetCopilotCustomAgentsSourceRequest is the request body for setting the custom agents source.
type SetCopilotCustomAgentsSourceRequest struct {
	// OrganizationID is the ID of the organization to use as the custom agents source.
	OrganizationID int64 `json:"organization_id"`
	// CreateRuleset controls whether to create a ruleset to protect agent definition files.
	// Defaults to true when omitted.
	CreateRuleset *bool `json:"create_ruleset,omitempty"`
}

// ListEnterpriseCustomAgents gets the list of custom agents defined for an enterprise.
//
// If no source repository has been configured, CustomAgents may be nil.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/copilot/copilot-custom-agents?apiVersion=2022-11-28#get-custom-agents-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/copilot/custom-agents
func (s *CopilotService) ListEnterpriseCustomAgents(ctx context.Context, enterprise string, opts *ListOptions) (*EnterpriseCustomAgents, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/copilot/custom-agents", enterprise)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var agents *EnterpriseCustomAgents
	resp, err := s.client.Do(req, &agents)
	if err != nil {
		return nil, resp, err
	}

	return agents, resp, nil
}

// GetEnterpriseCustomAgentsSource gets the organization and repository configured as
// the source for custom agent definitions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/copilot/copilot-custom-agents?apiVersion=2022-11-28#get-the-source-organization-for-custom-agents-in-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/copilot/custom-agents/source
func (s *CopilotService) GetEnterpriseCustomAgentsSource(ctx context.Context, enterprise string) (*CopilotCustomAgentsSource, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/copilot/custom-agents/source", enterprise)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var source *CopilotCustomAgentsSource
	resp, err := s.client.Do(req, &source)
	if err != nil {
		return nil, resp, err
	}

	return source, resp, nil
}

// SetEnterpriseCustomAgentsSource sets an organization as the source for custom agent
// definitions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/copilot/copilot-custom-agents?apiVersion=2022-11-28#set-the-source-organization-for-custom-agents-in-an-enterprise
//
//meta:operation PUT /enterprises/{enterprise}/copilot/custom-agents/source
func (s *CopilotService) SetEnterpriseCustomAgentsSource(ctx context.Context, enterprise string, request SetCopilotCustomAgentsSourceRequest) (*CopilotCustomAgentsSource, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/copilot/custom-agents/source", enterprise)

	req, err := s.client.NewRequest(ctx, "PUT", u, request)
	if err != nil {
		return nil, nil, err
	}

	var source *CopilotCustomAgentsSource
	resp, err := s.client.Do(req, &source)
	if err != nil {
		return nil, resp, err
	}

	return source, resp, nil
}

// DeleteEnterpriseCustomAgentsSource removes the custom agents source configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/copilot/copilot-custom-agents?apiVersion=2022-11-28#delete-the-custom-agents-source-for-an-enterprise
//
//meta:operation DELETE /enterprises/{enterprise}/copilot/custom-agents/source
func (s *CopilotService) DeleteEnterpriseCustomAgentsSource(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/copilot/custom-agents/source", enterprise)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
