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

const actionsPolicyResponseJSON = `{
  "id": 1,
  "name": "Require approved actors",
  "target": "actions",
  "source_type": "Organization",
  "source": "o",
  "enforcement": "active",
  "conditions": {
    "repository_name": {"include": ["~ALL"], "exclude": []},
    "workflow_path": {"include": [".github/workflows/*.yml"], "exclude": []}
  },
  "rules": [
    {"type": "restrict_actions_actors", "parameters": {"allowed_actors": [{"id": 1234, "type": "Team"}]}},
    {"type": "restrict_action_events", "parameters": {"allowed_events": ["push", "pull_request"]}}
  ],
  "node_id": "RUL_lA"
}`

func testActionsPolicy() *ActionsPolicy {
	return &ActionsPolicy{
		ID:          1,
		Name:        "Require approved actors",
		Target:      ActionsPolicyTargetActions,
		SourceType:  RulesetSourceTypeOrganization,
		Source:      "o",
		Enforcement: RulesetEnforcementActive,
		Conditions: &ActionsPolicyConditions{
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include: []string{"~ALL"},
				Exclude: []string{},
			},
			WorkflowPath: &ActionsPolicyWorkflowPathConditionParameters{
				Include: []string{".github/workflows/*.yml"},
				Exclude: []string{},
			},
		},
		Rules: []*ActionsPolicyRule{
			{
				Type: ActionsPolicyRuleTypeRestrictActionsActors,
				Parameters: &ActionsPolicyRuleParameters{
					AllowedActors: []*ActionsPolicyActor{{ID: 1234, Type: ActionsPolicyActorTypeTeam}},
				},
			},
			{
				Type: ActionsPolicyRuleTypeRestrictActionEvents,
				Parameters: &ActionsPolicyRuleParameters{
					AllowedEvents: []string{"push", "pull_request"},
				},
			},
		},
		NodeID: new("RUL_lA"),
	}
}

func TestActionsService_ListOrganizationPolicies(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"page": "2", "per_page": "50", "has_parents": "false"})
		fmt.Fprintf(w, `{"total_count":1,"policies":[%v]}`, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.ListOrganizationPolicies(ctx, "o", &ActionsPolicyListOptions{
		HasParents: new(false),
		ListOptions: ListOptions{
			Page:    2,
			PerPage: 50,
		},
	})
	if err != nil {
		t.Errorf("Actions.ListOrganizationPolicies returned error: %v", err)
	}
	want := &ActionsPolicyList{TotalCount: 1, Policies: []*ActionsPolicy{testActionsPolicy()}}
	if !cmp.Equal(got, want) {
		t.Errorf("Actions.ListOrganizationPolicies returned %+v, want %+v", got, want)
	}

	const methodName = "ListOrganizationPolicies"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListOrganizationPolicies(ctx, "\n", &ActionsPolicyListOptions{})
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListOrganizationPolicies(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateOrganizationPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := CreateActionsPolicyRequest{
		Name:        "Require approved actors",
		Enforcement: RulesetEnforcementActive,
		Rules: []*ActionsPolicyRule{{
			Type: ActionsPolicyRuleTypeRestrictActionsActors,
			Parameters: &ActionsPolicyRuleParameters{
				AllowedActors: []*ActionsPolicyActor{{ID: 1234, Type: ActionsPolicyActorTypeTeam}},
			},
		}},
	}

	mux.HandleFunc("/orgs/o/actions/policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.CreateOrganizationPolicy(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.CreateOrganizationPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.CreateOrganizationPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "CreateOrganizationPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.CreateOrganizationPolicy(ctx, "\n", input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.CreateOrganizationPolicy(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetOrganizationPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.GetOrganizationPolicy(ctx, "o", 1)
	if err != nil {
		t.Errorf("Actions.GetOrganizationPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.GetOrganizationPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrganizationPolicy(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrganizationPolicy(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_UpdateOrganizationPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := UpdateActionsPolicyRequest{Name: new("Updated policy"), Enforcement: new(RulesetEnforcementActive)}
	mux.HandleFunc("/orgs/o/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.UpdateOrganizationPolicy(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Actions.UpdateOrganizationPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.UpdateOrganizationPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "UpdateOrganizationPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.UpdateOrganizationPolicy(ctx, "\n", 1, input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.UpdateOrganizationPolicy(ctx, "o", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_DeleteOrganizationPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	if _, err := client.Actions.DeleteOrganizationPolicy(ctx, "o", 1); err != nil {
		t.Errorf("Actions.DeleteOrganizationPolicy returned error: %v", err)
	}

	const methodName = "DeleteOrganizationPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteOrganizationPolicy(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteOrganizationPolicy(ctx, "o", 1)
	})
}

func TestActionsService_ListRepositoryPolicies(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"has_parents": "true"})
		fmt.Fprintf(w, `{"total_count":1,"policies":[%v]}`, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.ListRepositoryPolicies(ctx, "o", "r", &ActionsPolicyListOptions{HasParents: new(true)})
	if err != nil {
		t.Errorf("Actions.ListRepositoryPolicies returned error: %v", err)
	}
	want := &ActionsPolicyList{TotalCount: 1, Policies: []*ActionsPolicy{testActionsPolicy()}}
	if !cmp.Equal(got, want) {
		t.Errorf("Actions.ListRepositoryPolicies returned %+v, want %+v", got, want)
	}

	const methodName = "ListRepositoryPolicies"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepositoryPolicies(ctx, "\n", "\n", nil)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepositoryPolicies(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateRepositoryPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := CreateActionsPolicyRequest{
		Name:        "Require approved events",
		Enforcement: RulesetEnforcementActive,
		Conditions: &ActionsPolicyConditions{
			WorkflowPath: &ActionsPolicyWorkflowPathConditionParameters{
				Include: []string{".github/workflows/*.yml"},
				Exclude: []string{},
			},
		},
		Rules: []*ActionsPolicyRule{{
			Type: ActionsPolicyRuleTypeRestrictActionEvents,
			Parameters: &ActionsPolicyRuleParameters{
				AllowedEvents: []string{"push"},
			},
		}},
	}

	mux.HandleFunc("/repos/o/r/actions/policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.CreateRepositoryPolicy(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateRepositoryPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.CreateRepositoryPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "CreateRepositoryPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.CreateRepositoryPolicy(ctx, "\n", "\n", input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.CreateRepositoryPolicy(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepositoryPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.GetRepositoryPolicy(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Actions.GetRepositoryPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.GetRepositoryPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "GetRepositoryPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepositoryPolicy(ctx, "\n", "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepositoryPolicy(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_UpdateRepositoryPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := UpdateActionsPolicyRequest{Name: new("Updated repository policy"), Enforcement: new(RulesetEnforcementEvaluate)}
	mux.HandleFunc("/repos/o/r/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		fmt.Fprint(w, actionsPolicyResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.Actions.UpdateRepositoryPolicy(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Actions.UpdateRepositoryPolicy returned error: %v", err)
	}
	if want := testActionsPolicy(); !cmp.Equal(got, want) {
		t.Errorf("Actions.UpdateRepositoryPolicy returned %+v, want %+v", got, want)
	}

	const methodName = "UpdateRepositoryPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.UpdateRepositoryPolicy(ctx, "\n", "\n", 1, input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.UpdateRepositoryPolicy(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_DeleteRepositoryPolicy(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	if _, err := client.Actions.DeleteRepositoryPolicy(ctx, "o", "r", 1); err != nil {
		t.Errorf("Actions.DeleteRepositoryPolicy returned error: %v", err)
	}

	const methodName = "DeleteRepositoryPolicy"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteRepositoryPolicy(ctx, "\n", "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteRepositoryPolicy(ctx, "o", "r", 1)
	})
}
