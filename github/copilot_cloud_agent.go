// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// CopilotCloudAgentConfiguration represents the Copilot cloud agent configuration for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/copilot/copilot-cloud-agent-management?apiVersion=2026-03-10#get-copilot-cloud-agent-configuration-for-a-repository
type CopilotCloudAgentConfiguration struct {
	McpConfiguration                      *json.RawMessage `json:"mcp_configuration,omitempty"`
	EnabledTools                          *EnabledTools    `json:"enabled_tools,omitempty"`
	RequireActionsWorkflowApproval        *bool            `json:"require_actions_workflow_approval,omitempty"`
	IsFirewallEnabled                     *bool            `json:"is_firewall_enabled,omitempty"`
	IsFirewallRecommendedAllowlistEnabled *bool            `json:"is_firewall_recommended_allowlist_enabled,omitempty"`
	CustomAllowlist                       []string         `json:"custom_allowlist,omitempty"`
}

// EnabledTools represents the enabled review tools for Copilot cloud agent.
type EnabledTools struct {
	Codeql                        *bool `json:"codeql,omitempty"`
	CopilotCodeReview             *bool `json:"copilot_code_review,omitempty"`
	SecretScanning                *bool `json:"secret_scanning,omitempty"`
	DependencyVulnerabilityChecks *bool `json:"dependency_vulnerability_checks,omitempty"`
}

// GetCopilotCloudAgentConfiguration gets the Copilot cloud agent configuration for a repository.
//
// GitHub API docs: https://docs.github.com/rest/copilot/copilot-cloud-agent-management?apiVersion=2026-03-10#get-copilot-cloud-agent-configuration-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/copilot/cloud-agent/configuration
func (s *CopilotService) GetCopilotCloudAgentConfiguration(ctx context.Context, owner, repo string) (*CopilotCloudAgentConfiguration, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/copilot/cloud-agent/configuration", owner, repo)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var config CopilotCloudAgentConfiguration
	resp, err := s.client.Do(req, &config)
	if err != nil {
		return nil, resp, err
	}

	return &config, resp, nil
}
