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

func TestCopilotService_ListEnterpriseCustomAgents(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/custom-agents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "10"})
		fmt.Fprint(w, `{
			"custom_agents": [
				{
					"name": "performance-optimizer",
					"file_path": "agents/performance-optimizer.agent.md",
					"url": "https://github.com/octo-org/.github-private/blob/main/agents/performance-optimizer.agent.md"
				},
				{
					"name": "security-reviewer",
					"file_path": "agents/security-reviewer.agent.md",
					"url": "https://github.com/octo-org/.github-private/blob/main/agents/security-reviewer.agent.md"
				}
			]
		}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 10}
	ctx := t.Context()
	got, _, err := client.Copilot.ListEnterpriseCustomAgents(ctx, "e", opts)
	if err != nil {
		t.Errorf("Copilot.ListEnterpriseCustomAgents returned error: %v", err)
	}

	want := &EnterpriseCustomAgents{
		CustomAgents: []*CopilotCustomAgent{
			{
				Name:     new("performance-optimizer"),
				FilePath: new("agents/performance-optimizer.agent.md"),
				URL:      new("https://github.com/octo-org/.github-private/blob/main/agents/performance-optimizer.agent.md"),
			},
			{
				Name:     new("security-reviewer"),
				FilePath: new("agents/security-reviewer.agent.md"),
				URL:      new("https://github.com/octo-org/.github-private/blob/main/agents/security-reviewer.agent.md"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.ListEnterpriseCustomAgents returned %+v, want %+v", got, want)
	}

	const methodName = "ListEnterpriseCustomAgents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.ListEnterpriseCustomAgents(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.ListEnterpriseCustomAgents(ctx, "e", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCopilotService_ListEnterpriseCustomAgents_NullAgents(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/custom-agents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"custom_agents": null}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.ListEnterpriseCustomAgents(ctx, "e", nil)
	if err != nil {
		t.Errorf("Copilot.ListEnterpriseCustomAgents returned error: %v", err)
	}

	want := &EnterpriseCustomAgents{CustomAgents: nil}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.ListEnterpriseCustomAgents returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_GetEnterpriseCustomAgentsSource(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/custom-agents/source", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"organization": {
				"id": 1,
				"login": "octo-org"
			},
			"repository": {
				"id": 2,
				"name": ".github-private",
				"full_name": "octo-org/.github-private"
			}
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseCustomAgentsSource(ctx, "e")
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseCustomAgentsSource returned error: %v", err)
	}

	want := &CopilotCustomAgentsSource{
		Organization: &CopilotCustomAgentsSourceOrganization{
			ID:    new(int64(1)),
			Login: new("octo-org"),
		},
		Repository: &CopilotCustomAgentsSourceRepository{
			ID:       new(int64(2)),
			Name:     new(".github-private"),
			FullName: new("octo-org/.github-private"),
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseCustomAgentsSource returned %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseCustomAgentsSource"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.GetEnterpriseCustomAgentsSource(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.GetEnterpriseCustomAgentsSource(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCopilotService_GetEnterpriseCustomAgentsSource_Null(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/custom-agents/source", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"organization": null, "repository": null}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.GetEnterpriseCustomAgentsSource(ctx, "e")
	if err != nil {
		t.Errorf("Copilot.GetEnterpriseCustomAgentsSource returned error: %v", err)
	}

	want := &CopilotCustomAgentsSource{}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.GetEnterpriseCustomAgentsSource returned %+v, want %+v", got, want)
	}
}

func TestCopilotService_SetEnterpriseCustomAgentsSource(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SetCopilotCustomAgentsSourceRequest{
		OrganizationID: 123,
		CreateRuleset:  new(false),
	}

	mux.HandleFunc("/enterprises/e/copilot/custom-agents/source", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"organization": {
				"id": 123,
				"login": "octo-org",
				"avatar_url": "https://github.com/images/error/octocat_happy.gif"
			},
			"repository": {
				"id": 2,
				"name": ".github-private",
				"full_name": "octo-org/.github-private"
			},
			"ruleset": {
				"id": 42,
				"name": "Protect custom agents",
				"enforcement": "active"
			}
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Copilot.SetEnterpriseCustomAgentsSource(ctx, "e", input)
	if err != nil {
		t.Errorf("Copilot.SetEnterpriseCustomAgentsSource returned error: %v", err)
	}

	want := &CopilotCustomAgentsSource{
		Organization: &CopilotCustomAgentsSourceOrganization{
			ID:        new(int64(123)),
			Login:     new("octo-org"),
			AvatarURL: new("https://github.com/images/error/octocat_happy.gif"),
		},
		Repository: &CopilotCustomAgentsSourceRepository{
			ID:       new(int64(2)),
			Name:     new(".github-private"),
			FullName: new("octo-org/.github-private"),
		},
		Ruleset: &CopilotCustomAgentsSourceRuleset{
			ID:          new(int64(42)),
			Name:        new("Protect custom agents"),
			Enforcement: new("active"),
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Copilot.SetEnterpriseCustomAgentsSource returned %+v, want %+v", got, want)
	}

	const methodName = "SetEnterpriseCustomAgentsSource"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Copilot.SetEnterpriseCustomAgentsSource(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Copilot.SetEnterpriseCustomAgentsSource(ctx, "e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCopilotService_DeleteEnterpriseCustomAgentsSource(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/copilot/custom-agents/source", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Copilot.DeleteEnterpriseCustomAgentsSource(ctx, "e")
	if err != nil {
		t.Errorf("Copilot.DeleteEnterpriseCustomAgentsSource returned error: %v", err)
	}

	const methodName = "DeleteEnterpriseCustomAgentsSource"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Copilot.DeleteEnterpriseCustomAgentsSource(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Copilot.DeleteEnterpriseCustomAgentsSource(ctx, "e")
	})
}
