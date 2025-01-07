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
	t.Parallel()
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
		"With Metadata": {
			data: `{
                    "type": "creation",
					"ruleset_source_type": "Repository",
					"ruleset_source": "google",
					"ruleset_id": 1984
           		   }`,
			want: &RepositoryRule{
				RulesetSource:     "google",
				RulesetSourceType: "Repository",
				RulesetID:         1984,
				Type:              "creation",
			},
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
		"Valid merge_queue": {
			data: `{"type":"merge_queue"}`,
			want: &RepositoryRule{
				Type:       "merge_queue",
				Parameters: nil,
			},
		},
		"Valid merge_queue with params": {
			data: `{
				"type":"merge_queue",
				"parameters":{
					"check_response_timeout_minutes": 35,
					"grouping_strategy": "HEADGREEN",
					"max_entries_to_build": 8,
					"max_entries_to_merge": 4,
					"merge_method": "SQUASH",
					"min_entries_to_merge": 2,
					"min_entries_to_merge_wait_minutes": 13
				}
			}`,
			want: NewMergeQueueRule(&MergeQueueRuleParameters{
				CheckResponseTimeoutMinutes:  35,
				GroupingStrategy:             "HEADGREEN",
				MaxEntriesToBuild:            8,
				MaxEntriesToMerge:            4,
				MergeMethod:                  "SQUASH",
				MinEntriesToMerge:            2,
				MinEntriesToMergeWaitMinutes: 13,
			}),
		},
		"Invalid merge_queue with params": {
			data: `{
				"type":"merge_queue",
				"parameters":{
					"check_response_timeout_minutes": "35",
					"grouping_strategy": "HEADGREEN",
					"max_entries_to_build": "8",
					"max_entries_to_merge": "4",
					"merge_method": "SQUASH",
					"min_entries_to_merge": "2",
					"min_entries_to_merge_wait_minutes": "13"
				}
			}`,
			want: &RepositoryRule{
				Type:       "merge_queue",
				Parameters: nil,
			},
			wantErr: true,
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
		"Valid file_path_restriction params": {
			data: `{"type":"file_path_restriction","parameters":{"restricted_file_paths":["/a/file"]}}`,
			want: NewFilePathRestrictionRule(&RuleFileParameters{
				RestrictedFilePaths: &[]string{"/a/file"},
			}),
		},
		"Invalid file_path_restriction params": {
			data: `{"type":"file_path_restriction","parameters":{"restricted_file_paths":true}}`,
			want: &RepositoryRule{
				Type:       "file_path_restriction",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid pull_request params": {
			data: `{
				"type":"pull_request",
				"parameters":{
					"allowed_merge_methods": ["rebase","squash"],
					"dismiss_stale_reviews_on_push": true,
					"require_code_owner_review": true,
					"require_last_push_approval": true,
					"required_approving_review_count": 1,
					"required_review_thread_resolution":true
				}
			}`,
			want: NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
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
			data: `{"type":"required_status_checks","parameters":{"required_status_checks":[{"context":"test","integration_id":1}],"strict_required_status_checks_policy":true,"do_not_enforce_on_create":true}}`,
			want: NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
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
		"Valid Required workflows params": {
			data: `{"type":"workflows","parameters":{"workflows":[{"path": ".github/workflows/test.yml", "repository_id": 1}]}}`,
			want: NewRequiredWorkflowsRule(&RequiredWorkflowsRuleParameters{
				RequiredWorkflows: []*RuleRequiredWorkflow{
					{
						Path:         ".github/workflows/test.yml",
						RepositoryID: Ptr(int64(1)),
					},
				},
			}),
		},
		"Invalid Required workflows params": {
			data: `{"type":"workflows","parameters":{"workflows":[{"path": ".github/workflows/test.yml", "repository_id": "test"}]}}`,
			want: &RepositoryRule{
				Type:       "workflows",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid Required code_scanning params": {
			data: `{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool": "CodeQL", "security_alerts_threshold": "high_or_higher", "alerts_threshold": "errors"}]}}`,
			want: NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
		"Invalid Required code_scanning params": {
			data: `{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool": 1}]}}`,
			want: &RepositoryRule{
				Type:       "code_scanning",
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
		"Valid max_file_path_length params": {
			data: `{"type":"max_file_path_length","parameters":{"max_file_path_length": 255}}`,
			want: NewMaxFilePathLengthRule(&RuleMaxFilePathLengthParameters{
				MaxFilePathLength: 255,
			}),
		},
		"Invalid max_file_path_length params": {
			data: `{"type":"max_file_path_length","parameters":{"max_file_path_length": "255"}}`,
			want: &RepositoryRule{
				Type:       "max_file_path_length",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid file_extension_restriction params": {
			data: `{"type":"file_extension_restriction","parameters":{"restricted_file_extensions":[".exe"]}}`,
			want: NewFileExtensionRestrictionRule(&RuleFileExtensionRestrictionParameters{
				RestrictedFileExtensions: []string{".exe"},
			}),
		},
		"Invalid file_extension_restriction params": {
			data: `{"type":"file_extension_restriction","parameters":{"restricted_file_extensions":true}}`,
			want: &RepositoryRule{
				Type:       "file_extension_restriction",
				Parameters: nil,
			},
			wantErr: true,
		},
		"Valid max_file_size params": {
			data: `{"type":"max_file_size","parameters":{"max_file_size": 1024}}`,
			want: NewMaxFileSizeRule(&RuleMaxFileSizeParameters{
				MaxFileSize: 1024,
			}),
		},
		"Invalid max_file_size params": {
			data: `{"type":"max_file_size","parameters":{"max_file_size": "1024"}}`,
			want: &RepositoryRule{
				Type:       "max_file_size",
				Parameters: nil,
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		rule := &RepositoryRule{}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rules/branches/branch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "ruleset_id": 42069,
			  "ruleset_source_type": "Repository",
			  "ruleset_source": "google",
			  "type": "creation"
			},
			{
			  "ruleset_id": 42069,
			  "ruleset_source_type": "Organization",
			  "ruleset_source": "google",
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
	creationRule.RulesetID = 42069
	creationRule.RulesetSource = "google"
	creationRule.RulesetSourceType = "Repository"
	updateRule := NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
		UpdateAllowsFetchAndMerge: true,
	})
	updateRule.RulesetID = 42069
	updateRule.RulesetSource = "google"
	updateRule.RulesetSourceType = "Organization"

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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "id": 42,
			  "name": "ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "enabled",
			  "created_at": `+referenceTimeStr+`,
			  "updated_at": `+referenceTimeStr+`
			},
			{
			  "id": 314,
			  "name": "Another ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "enabled",
			  "created_at": `+referenceTimeStr+`,
			  "updated_at": `+referenceTimeStr+`
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
			ID:          Ptr(int64(42)),
			Name:        "ruleset",
			SourceType:  Ptr("Repository"),
			Source:      "o/repo",
			Enforcement: "enabled",
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		{
			ID:          Ptr(int64(314)),
			Name:        "Another ruleset",
			SourceType:  Ptr("Repository"),
			Source:      "o/repo",
			Enforcement: "enabled",
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	ruleSet, _, err := client.Repositories.CreateRuleset(ctx, "o", "repo", Ruleset{
		Name:        "ruleset",
		Enforcement: "enabled",
	})
	if err != nil {
		t.Errorf("Repositories.CreateRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr("Repository"),
		Source:      "o/repo",
		Enforcement: "enabled",
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.CreateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "CreateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRuleset(ctx, "o", "repo", Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateRulesetWithPushRules(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "enabled",
			"target": "push",
			"rules": [
				{
					"type": "file_path_restriction",
					"parameters": {
						"restricted_file_paths": ["/a/file"]
					}
				},
				{
					"type": "max_file_path_length",
					"parameters": {
						"max_file_path_length": 255
					}
				},
				{
					"type": "file_extension_restriction",
					"parameters": {
						"restricted_file_extensions": [".exe"]
					}
				},
				{
					"type": "max_file_size",
					"parameters": {
						"max_file_size": 1024
					}
				}
			]
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.CreateRuleset(ctx, "o", "repo", Ruleset{
		Name:        "ruleset",
		Enforcement: "enabled",
	})
	if err != nil {
		t.Errorf("Repositories.CreateRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr("Repository"),
		Source:      "o/repo",
		Target:      Ptr("push"),
		Enforcement: "enabled",
		Rules: []*RepositoryRule{
			NewFilePathRestrictionRule(&RuleFileParameters{
				RestrictedFilePaths: &[]string{"/a/file"},
			}),
			NewMaxFilePathLengthRule(&RuleMaxFilePathLengthParameters{
				MaxFilePathLength: 255,
			}),
			NewFileExtensionRestrictionRule(&RuleFileExtensionRestrictionParameters{
				RestrictedFileExtensions: []string{".exe"},
			}),
			NewMaxFileSizeRule(&RuleMaxFileSizeParameters{
				MaxFileSize: 1024,
			}),
		},
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.CreateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "CreateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRuleset(ctx, "o", "repo", Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "enabled",
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.GetRuleset(ctx, "o", "repo", 42, true)
	if err != nil {
		t.Errorf("Repositories.GetRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "enabled",
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	ruleSet, _, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, Ruleset{
		Name:        "ruleset",
		Enforcement: "enabled",
	})
	if err != nil {
		t.Errorf("Repositories.UpdateRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr("Repository"),
		Source:      "o/repo",
		Enforcement: "enabled",
	}

	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.UpdateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "UpdateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateRulesetClearBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "enabled"
			"conditions": {
				"ref_name": {
					"include": [
						"refs/heads/main",
						"refs/heads/master"
					],
					"exclude": [
						"refs/heads/dev*"
					]
				}
			},
			"rules": [
			  {
					"type": "creation"
			  }
			]
		}`)
	})

	ctx := context.Background()

	_, err := client.Repositories.UpdateRulesetClearBypassActor(ctx, "o", "repo", 42)
	if err != nil {
		t.Errorf("Repositories.UpdateRulesetClearBypassActor returned error: %v \n", err)
	}

	const methodName = "UpdateRulesetClearBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.UpdateRulesetClearBypassActor(ctx, "o", "repo", 42)
	})
}

func TestRepositoriesService_UpdateRulesetNoBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	rs := Ruleset{
		Name:        "ruleset",
		Source:      "o/repo",
		Enforcement: "enabled",
	}

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

	ruleSet, _, err := client.Repositories.UpdateRulesetNoBypassActor(ctx, "o", "repo", 42, rs)
	if err != nil {
		t.Errorf("Repositories.UpdateRulesetNoBypassActor returned error: %v \n", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr("Repository"),
		Source:      "o/repo",
		Enforcement: "enabled",
	}

	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.UpdateRulesetNoBypassActor returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "UpdateRulesetNoBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRulesetNoBypassActor(ctx, "o", "repo", 42, Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
