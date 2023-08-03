// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryRule_UnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		data    string
		want    *RepositoryRule
		wantErr bool
	}{
		"Invalid JSON": {
			data: `{`,
			want: &RepositoryRule{
				Type:       "",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid creation": {
			data: `{"type":"creation"}`,
			want: NewCreationRule(),
		},
		"Valid deletion": {
			data: `{"type":"deletion"}`,
			want: &RepositoryRule{
				Type:       "deletion",
				Parameters: nil,
			},
		},
		"Valid required_linear_history": {
			data: `{"type":"required_linear_history"}`,
			want: &RepositoryRule{
				Type:       "required_linear_history",
				Parameters: nil,
			},
		},
		"Valid required_signatures": {
			data: `{"type":"required_signatures"}`,
			want: &RepositoryRule{
				Type:       "required_signatures",
				Parameters: nil,
			},
		},
		"Valid non_fast_forward": {
			data: `{"type":"non_fast_forward"}`,
			want: &RepositoryRule{
				Type:       "non_fast_forward",
				Parameters: nil,
			},
		},
		"Valid update params": {
			data: `{"type":"update","parameters":{"update_allows_fetch_and_merge":true}}`,
			want: NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{UpdateAllowsFetchAndMerge: true}),
		},
		"Invalid update params": {
			data: `{"type":"update","parameters":{"update_allows_fetch_and_merge":"true"}}`,
			want: &RepositoryRule{
				Type:       "update",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid required_deployments params": {
			data: `{"type":"required_deployments","parameters":{"required_deployment_environments":["test"]}}`,
			want: NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
		},
		"Invalid required_deployments params": {
			data: `{"type":"required_deployments","parameters":{"required_deployment_environments":true}}`,
			want: &RepositoryRule{
				Type:       "required_deployments",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid commit_message_pattern params": {
			data: `{"type":"commit_message_pattern","parameters":{"operator":"starts_with","pattern":"github"}}`,
			want: NewCommitMessagePatternRule(&RulePatternParameters{
				Operator: "starts_with",
				Pattern:  "github",
			}),
		},
		"Invalid commit_message_pattern params": {
			data: `{"type":"commit_message_pattern","parameters":{"operator":"starts_with","pattern":1}}`,
			want: &RepositoryRule{
				Type:       "commit_message_pattern",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid commit_author_email_pattern params": {
			data: `{"type":"commit_author_email_pattern","parameters":{"operator":"starts_with","pattern":"github"}}`,
			want: NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "starts_with",
				Pattern:  "github",
			}),
		},
		"Invalid commit_author_email_pattern params": {
			data: `{"type":"commit_author_email_pattern","parameters":{"operator":"starts_with","pattern":1}}`,
			want: &RepositoryRule{
				Type:       "commit_author_email_pattern",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid committer_email_pattern params": {
			data: `{"type":"committer_email_pattern","parameters":{"operator":"starts_with","pattern":"github"}}`,
			want: NewCommitterEmailPatternRule(&RulePatternParameters{
				Operator: "starts_with",
				Pattern:  "github",
			}),
		},
		"Invalid committer_email_pattern params": {
			data: `{"type":"committer_email_pattern","parameters":{"operator":"starts_with","pattern":1}}`,
			want: &RepositoryRule{
				Type:       "committer_email_pattern",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid branch_name_pattern params": {
			data: `{"type":"branch_name_pattern","parameters":{"operator":"starts_with","pattern":"github"}}`,
			want: NewBranchNamePatternRule(&RulePatternParameters{
				Operator: "starts_with",
				Pattern:  "github",
			}),
		},
		"Invalid branch_name_pattern params": {
			data: `{"type":"branch_name_pattern","parameters":{"operator":"starts_with","pattern":1}}`,
			want: &RepositoryRule{
				Type:       "branch_name_pattern",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid tag_name_pattern params": {
			data: `{"type":"tag_name_pattern","parameters":{"operator":"starts_with","pattern":"github"}}`,
			want: NewTagNamePatternRule(&RulePatternParameters{
				Operator: "starts_with",
				Pattern:  "github",
			}),
		},
		"Invalid tag_name_pattern params": {
			data: `{"type":"tag_name_pattern","parameters":{"operator":"starts_with","pattern":1}}`,
			want: &RepositoryRule{
				Type:       "tag_name_pattern",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid pull_request params": {
			data: `{
				"type":"pull_request",
				"parameters":{
					"dismiss_stale_reviews_on_push": true,
					"require_code_owner_review": true,
					"require_last_push_approval": true,
					"required_approving_review_count": 1,
					"required_review_thread_resolution":true
				}
			}`,
			want: NewPullRequestRule(&PullRequestRuleParameters{
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
		},
		"Invalid pull_request params": {
			data: `{"type":"pull_request","parameters": {"dismiss_stale_reviews_on_push":"true"}}`,
			want: &RepositoryRule{
				Type:       "pull_request",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid required_status_checks params": {
			data: `{"type":"required_status_checks","parameters":{"required_status_checks":[{"context":"test","integration_id":1}],"strict_required_status_checks_policy":true}}`,
			want: NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Int64(1),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
		},
		"Invalid required_status_checks params": {
			data: `{"type":"required_status_checks",
			"parameters": {
				"required_status_checks": [
				  {
					"context": 1
				  }
				]
			  }}`,
			want: &RepositoryRule{
				Type:       "required_status_checks",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Invalid type": {
			data: `{"type":"unknown"}`,
			want: &RepositoryRule{
				Type:       "",
				Parameters: nil,
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		rule := &RepositoryRule{}

		t.Run(name, func(t *testing.T) {
			err := rule.UnmarshalJSON([]byte(tc.data))
			if err == nil && tc.wantErr {
				t.Errorf("RepositoryRule.UnmarshalJSON returned nil instead of an error")
			}
			if err != nil && !tc.wantErr {
				t.Errorf("RepositoryRule.UnmarshalJSON returned an unexpected error: %+v", err)
			}
			if !cmp.Equal(tc.want, rule) {
				t.Errorf("RepositoryRule.UnmarshalJSON expected rule %+v, got %+v", tc.want, rule)
			}
		})
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
		creationRule,
		updateRule,
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

func TestRepositoriesService_GetRulesForBranchEmptyUpdateRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/repo/rules/branches/branch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "type": "update"
			}
		]`)
	})

	ctx := context.Background()
	rules, _, err := client.Repositories.GetRulesForBranch(ctx, "o", "repo", "branch")
	if err != nil {
		t.Errorf("Repositories.GetRulesForBranch returned error: %v", err)
	}

	updateRule := NewUpdateRule(nil)

	want := []*RepositoryRule{
		updateRule,
	}
	if !cmp.Equal(rules, want) {
		t.Errorf("Repositories.GetRulesForBranch returned %+v, want %+v", Stringify(rules), Stringify(want))
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
			ID:          Int64(42),
			Name:        "ruleset",
			SourceType:  String("Repository"),
			Source:      "o/repo",
			Enforcement: "enabled",
		},
		{
			ID:          Int64(314),
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
		ID:          Int64(42),
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
		ID:          Int64(42),
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
		ID:          Int64(42),
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
