// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCopilotService_GetCopilotCloudAgentConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/copilot/cloud-agent/configuration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
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
		}`)
	})

	ctx := t.Context()
	config, _, err := client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "o", "r")
	if err != nil {
		t.Errorf("GetCopilotCloudAgentConfiguration returned error: %v", err)
	}

	want := &CopilotCloudAgentConfiguration{
		McpConfiguration: nil,
		EnabledTools: &EnabledTools{
			Codeql:                        Ptr(true),
			CopilotCodeReview:             Ptr(true),
			SecretScanning:                Ptr(true),
			DependencyVulnerabilityChecks: Ptr(true),
		},
		RequireActionsWorkflowApproval:        Ptr(true),
		IsFirewallEnabled:                     Ptr(true),
		IsFirewallRecommendedAllowlistEnabled: Ptr(true),
		CustomAllowlist:                       []string{},
	}

	if !cmp.Equal(config, want) {
		t.Errorf("GetCopilotCloudAgentConfiguration returned %+v, want %+v", config, want)
	}

	const methodName = "GetCopilotCloudAgentConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetCopilotCloudAgentConfiguration(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
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

func TestCopilotCloudAgentConfiguration_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CopilotCloudAgentConfiguration{}, "{}")

	u := &CopilotCloudAgentConfiguration{
		McpConfiguration: nil,
		EnabledTools: &EnabledTools{
			Codeql:                        Ptr(true),
			CopilotCodeReview:             Ptr(false),
			SecretScanning:                Ptr(true),
			DependencyVulnerabilityChecks: Ptr(false),
		},
		RequireActionsWorkflowApproval:        Ptr(true),
		IsFirewallEnabled:                     Ptr(false),
		IsFirewallRecommendedAllowlistEnabled: Ptr(true),
		CustomAllowlist:                       []string{"192.168.0.0/16"},
	}

	want := `{
		"enabled_tools": {
			"codeql": true,
			"copilot_code_review": false,
			"secret_scanning": true,
			"dependency_vulnerability_checks": false
		},
		"require_actions_workflow_approval": true,
		"is_firewall_enabled": false,
		"is_firewall_recommended_allowlist_enabled": true,
		"custom_allowlist": ["192.168.0.0/16"]
	}`

	testJSONMarshal(t, u, want)
}
