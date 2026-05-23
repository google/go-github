// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCopilotService_GetCopilotCloudAgentConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		responseBody string
		want         *CopilotCloudAgentConfiguration
		wantErr      bool
	}{
		{
			name: "with null mcp_configuration",
			responseBody: `{
				"mcp_configuration": null,
				"enabled_tools": {
					"codeql": true,
					"copilot_code_review": true,
					"secret_scanning": true,
					"dependency_vulnerability_checks": true
				},
				"require_actions_workflow_approval": true,
				"is_firewall_enabled": true,
				"is_firewall_recommended_allowlist_enabled": true,
				"custom_allowlist": []
			}`,
			want: &CopilotCloudAgentConfiguration{
				McpConfiguration: nil,
				EnabledTools: EnabledTools{
					Codeql:                        true,
					CopilotCodeReview:             true,
					SecretScanning:                true,
					DependencyVulnerabilityChecks: true,
				},
				RequireActionsWorkflowApproval:        true,
				IsFirewallEnabled:                     true,
				IsFirewallRecommendedAllowlistEnabled: true,
				CustomAllowlist:                       []string{},
			},
			wantErr: false,
		},
		{
			name: "with custom allowlist",
			responseBody: `{
				"mcp_configuration": null,
				"enabled_tools": {
					"codeql": false,
					"copilot_code_review": true,
					"secret_scanning": false,
					"dependency_vulnerability_checks": true
				},
				"require_actions_workflow_approval": false,
				"is_firewall_enabled": true,
				"is_firewall_recommended_allowlist_enabled": true,
				"custom_allowlist": ["192.168.0.0/16", "10.0.0.0/8"]
			}`,
			want: &CopilotCloudAgentConfiguration{
				McpConfiguration: nil,
				EnabledTools: EnabledTools{
					Codeql:                        false,
					CopilotCodeReview:             true,
					SecretScanning:                false,
					DependencyVulnerabilityChecks: true,
				},
				RequireActionsWorkflowApproval:        false,
				IsFirewallEnabled:                     true,
				IsFirewallRecommendedAllowlistEnabled: true,
				CustomAllowlist:                       []string{"192.168.0.0/16", "10.0.0.0/8"},
			},
			wantErr: false,
		},
		{
			name: "all tools disabled",
			responseBody: `{
				"mcp_configuration": null,
				"enabled_tools": {
					"codeql": false,
					"copilot_code_review": false,
					"secret_scanning": false,
					"dependency_vulnerability_checks": false
				},
				"require_actions_workflow_approval": false,
				"is_firewall_enabled": false,
				"is_firewall_recommended_allowlist_enabled": false,
				"custom_allowlist": []
			}`,
			want: &CopilotCloudAgentConfiguration{
				McpConfiguration: nil,
				EnabledTools: EnabledTools{
					Codeql:                        false,
					CopilotCodeReview:             false,
					SecretScanning:                false,
					DependencyVulnerabilityChecks: false,
				},
				RequireActionsWorkflowApproval:        false,
				IsFirewallEnabled:                     false,
				IsFirewallRecommendedAllowlistEnabled: false,
				CustomAllowlist:                       []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc("/repos/o/r/copilot/cloud-agent/configuration", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fmt.Fprint(w, tt.responseBody)
			})

			ctx := t.Context()
			config, _, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "r")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCopilotCloudAgentConfiguration returned error: %v, wantErr: %v", err, tt.wantErr)
			}

			if !cmp.Equal(config, tt.want) {
				t.Errorf("GetCopilotCloudAgentConfiguration returned %+v, want %+v", config, tt.want)
			}
		})
	}
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_BadOptions(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	const methodName = "GetCopilotCloudAgentConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "\n", "\n")
		return err
	})
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_NewRequestFailure(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	const methodName = "GetCopilotCloudAgentConfiguration"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		ctx := t.Context()
		got, resp, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_InvalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_InvalidRepo(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "%")
	testURLParseError(t, err)
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_NotFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/copilot/cloud-agent/configuration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := t.Context()
	config, resp, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "r")
	if err == nil {
		t.Error("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("GetCopilotCloudAgentConfiguration return status %v, want %v", got, want)
	}
	if config != nil {
		t.Errorf("GetCopilotCloudAgentConfiguration return %+v, want nil", config)
	}
}

func TestCopilotService_GetCopilotCloudAgentConfiguration_Forbidden(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/copilot/cloud-agent/configuration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusForbidden)
	})

	ctx := t.Context()
	config, resp, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "r")
	if err == nil {
		t.Error("Expected HTTP 403 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusForbidden; got != want {
		t.Errorf("GetCopilotCloudAgentConfiguration return status %v, want %v", got, want)
	}
	if config != nil {
		t.Errorf("GetCopilotCloudAgentConfiguration return %+v, want nil", config)
	}
}

func TestCopilotCloudAgentConfiguration_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    *CopilotCloudAgentConfiguration
		want string
	}{
		{
			name: "empty configuration",
			u:    &CopilotCloudAgentConfiguration{},
			want: `{"mcp_configuration":null,"enabled_tools":{"codeql":false,"copilot_code_review":false,"secret_scanning":false,"dependency_vulnerability_checks":false},"require_actions_workflow_approval":false,"is_firewall_enabled":false,"is_firewall_recommended_allowlist_enabled":false,"custom_allowlist":null}`,
		},
		{
			name: "with all settings configured",
			u: &CopilotCloudAgentConfiguration{
				McpConfiguration: nil,
				EnabledTools: EnabledTools{
					Codeql:                        true,
					CopilotCodeReview:             false,
					SecretScanning:                true,
					DependencyVulnerabilityChecks: false,
				},
				RequireActionsWorkflowApproval:        true,
				IsFirewallEnabled:                     false,
				IsFirewallRecommendedAllowlistEnabled: true,
				CustomAllowlist:                       []string{"192.168.0.0/16"},
			},
			want: `{"mcp_configuration":null,"enabled_tools":{"codeql":true,"copilot_code_review":false,"secret_scanning":true,"dependency_vulnerability_checks":false},"require_actions_workflow_approval":true,"is_firewall_enabled":false,"is_firewall_recommended_allowlist_enabled":true,"custom_allowlist":["192.168.0.0/16"]}`,
		},
		{
			name: "with mcp configuration",
			u: &CopilotCloudAgentConfiguration{
				McpConfiguration: func() *json.RawMessage {
					raw := json.RawMessage(`{"type":"resource","uri":"stdio://server"}`)
					return &raw
				}(),
				EnabledTools: EnabledTools{
					Codeql:                        true,
					CopilotCodeReview:             true,
					SecretScanning:                true,
					DependencyVulnerabilityChecks: true,
				},
				RequireActionsWorkflowApproval:        true,
				IsFirewallEnabled:                     true,
				IsFirewallRecommendedAllowlistEnabled: false,
				CustomAllowlist:                       []string{},
			},
			want: `{"mcp_configuration":{"type":"resource","uri":"stdio://server"},"enabled_tools":{"codeql":true,"copilot_code_review":true,"secret_scanning":true,"dependency_vulnerability_checks":true},"require_actions_workflow_approval":true,"is_firewall_enabled":true,"is_firewall_recommended_allowlist_enabled":false,"custom_allowlist":[]}`,
		},
		{
			name: "with multiple allowlist entries",
			u: &CopilotCloudAgentConfiguration{
				McpConfiguration: nil,
				EnabledTools: EnabledTools{
					Codeql:                        false,
					CopilotCodeReview:             false,
					SecretScanning:                false,
					DependencyVulnerabilityChecks: false,
				},
				RequireActionsWorkflowApproval:        false,
				IsFirewallEnabled:                     true,
				IsFirewallRecommendedAllowlistEnabled: true,
				CustomAllowlist:                       []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12"},
			},
			want: `{"mcp_configuration":null,"enabled_tools":{"codeql":false,"copilot_code_review":false,"secret_scanning":false,"dependency_vulnerability_checks":false},"require_actions_workflow_approval":false,"is_firewall_enabled":true,"is_firewall_recommended_allowlist_enabled":true,"custom_allowlist":["192.168.0.0/16","10.0.0.0/8","172.16.0.0/12"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testJSONMarshal(t, tt.u, tt.want)
		})
	}
}

func TestEnabledTools_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    *EnabledTools
		want string
	}{
		{
			name: "all enabled",
			u: &EnabledTools{
				Codeql:                        true,
				CopilotCodeReview:             true,
				SecretScanning:                true,
				DependencyVulnerabilityChecks: true,
			},
			want: `{"codeql":true,"copilot_code_review":true,"secret_scanning":true,"dependency_vulnerability_checks":true}`,
		},
		{
			name: "all disabled",
			u: &EnabledTools{
				Codeql:                        false,
				CopilotCodeReview:             false,
				SecretScanning:                false,
				DependencyVulnerabilityChecks: false,
			},
			want: `{"codeql":false,"copilot_code_review":false,"secret_scanning":false,"dependency_vulnerability_checks":false}`,
		},
		{
			name: "mixed settings",
			u: &EnabledTools{
				Codeql:                        true,
				CopilotCodeReview:             false,
				SecretScanning:                true,
				DependencyVulnerabilityChecks: false,
			},
			want: `{"codeql":true,"copilot_code_review":false,"secret_scanning":true,"dependency_vulnerability_checks":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testJSONMarshal(t, tt.u, tt.want)
		})
	}
}
