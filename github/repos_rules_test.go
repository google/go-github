// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryRule_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in      *RepositoryRule
		want    string
		wantErr bool
	}{
		{
			in: &RepositoryRule{
				Type: "update",
				Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
					UpdateAllowsFetchAndMerge: true,
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "required_deployments",
				Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test"},
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "commit_message_pattern",
				Parameters: &RulePatternParameters{
					Name:     String("avoid test commits"),
					Negate:   Bool(true),
					Operator: "starts_with",
					Pattern:  "[test]",
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "commit_author_email_pattern",
				Parameters: &RulePatternParameters{
					Operator: "contains",
					Pattern:  "github",
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "committer_email_pattern",
				Parameters: &RulePatternParameters{
					Name:     String("avoid commit emails"),
					Negate:   Bool(true),
					Operator: "ends_with",
					Pattern:  "abc",
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "branch_name_pattern",
				Parameters: &RulePatternParameters{
					Name:     String("avoid branch names"),
					Negate:   Bool(true),
					Operator: "regex",
					Pattern:  "github$",
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "tag_name_pattern",
				Parameters: &RulePatternParameters{
					Name:     String("avoid tag names"),
					Negate:   Bool(true),
					Operator: "contains",
					Pattern:  "github",
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "pull_request",
				Parameters: &PullRequestRuleParameters{
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   1,
					RequiredReviewThreadResolution: true,
					DismissStaleReviewsOnPush:      true,
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "required_status_checks",
				Parameters: &RequiredStatusChecksRuleParameters{
					RequiredStatusChecks: []RuleRequiredStatusChecks{
						{
							Context:       "test",
							IntegrationID: Int64(1),
						},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RepositoryRule{
				Type: "unknown",
			},
			want:    `{`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		err := json.Unmarshal([]byte(tc.want), tc.in)
		if err == nil && tc.wantErr {
			t.Errorf("RepositoryRule.UnmarshalJSON returned nil instead of an error")
		}
		if err != nil && !tc.wantErr {
			t.Errorf("RepositoryRule.UnmarshalJSON returned an unexpected error: %+v", err)
		}
	}
}

func TestRepositoriesService_GetRulesForBranch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rules/branches/branch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "type": "creation"
			},
			{
			  "type": "update",
			  "parameters": {
			    "update_allows_fetch_and_merge": true
			  }
			}
		]`)
	})

	ctx := context.Background()
	rules, _, err := client.Repositories.GetRulesForBranch(ctx, "o", "repo", "branch")
	if err != nil {
		t.Errorf("Repositories.GetRulesForBranch returned error: %v", err)
	}

	creationRule := NewCreationRule()
	updateRule := NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
		UpdateAllowsFetchAndMerge: true,
	})

	want := []*RepositoryRule{
		&creationRule,
		&updateRule,
	}
	if !cmp.Equal(rules, want) {
		t.Errorf("Repositories.GetRulesForBranch returned %+v, want %+v", rules, want)
	}

	const methodName = "GetRulesForBranch"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRulesForBranch(ctx, "o", "repo", "branch")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetAllRulesets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "id": 42,
			  "name": "ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "enabled"
			},
			{
			  "id": 314,
			  "name": "Another ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "enabled"
			}
		]`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.GetAllRulesets(ctx, "o", "repo", false)
	if err != nil {
		t.Errorf("Repositories.GetAllRulesets returned error: %v", err)
	}

	want := []*Ruleset{
		{
			ID:          42,
			Name:        "ruleset",
			SourceType:  String("Repository"),
			Source:      "o/repo",
			Enforcement: "enabled",
		},
		{
			ID:          314,
			Name:        "Another ruleset",
			SourceType:  String("Repository"),
			Source:      "o/repo",
			Enforcement: "enabled",
		},
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.GetAllRulesets returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "GetAllRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAllRulesets(ctx, "o", "repo", false)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "enabled"
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.CreateRuleset(ctx, "o", "repo", &Ruleset{
		Name:        "ruleset",
		Enforcement: "enabled",
	})
	if err != nil {
		t.Errorf("Repositories.CreateRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          42,
		Name:        "ruleset",
		SourceType:  String("Repository"),
		Source:      "o/repo",
		Enforcement: "enabled",
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.CreateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "CreateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRuleset(ctx, "o", "repo", &Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "enabled"
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.GetRuleset(ctx, "o", "repo", 42, true)
	if err != nil {
		t.Errorf("Repositories.GetRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          42,
		Name:        "ruleset",
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "enabled",
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.GetRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "GetRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRuleset(ctx, "o", "repo", 42, true)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "enabled"
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, &Ruleset{
		Name:        "ruleset",
		Enforcement: "enabled",
	})
	if err != nil {
		t.Errorf("Repositories.UpdateRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          42,
		Name:        "ruleset",
		SourceType:  String("Repository"),
		Source:      "o/repo",
		Enforcement: "enabled",
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.UpdateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "UpdateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteRuleset(ctx, "o", "repo", 42)
	if err != nil {
		t.Errorf("Repositories.DeleteRuleset returned error: %v", err)
	}

	const methodName = "DeleteRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteRuleset(ctx, "o", "repo", 42)
	})
}
