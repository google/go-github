// Copyright 2024 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
)

func TestRepositoryRulesetRules_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       RepositoryRulesetRules
		expected    string
		expectError bool
	}{
		{
			name:     "Empty RepositoryRulesetRules returns empty array",
			input:    RepositoryRulesetRules{},
			expected: "[]",
		},
		{
			name: "RepositoryRulesetRules with Creation rule only",
			input: RepositoryRulesetRules{
				Creation: &EmptyRuleParameters{},
			},
			expected: `[{"type":"creation"}]`,
		},
		{
			name: "RepositoryRulesetRules with Deletion rule only",
			input: RepositoryRulesetRules{
				Deletion: &EmptyRuleParameters{},
			},
			expected: `[{"type":"deletion"}]`,
		},
		{
			name: "RepositoryRulesetRules with basic rules",
			input: RepositoryRulesetRules{
				Creation:              &EmptyRuleParameters{},
				Deletion:              &EmptyRuleParameters{},
				RequiredLinearHistory: &EmptyRuleParameters{},
				NonFastForward:        &EmptyRuleParameters{},
			},
			expected: `[{"type":"creation"},{"type":"deletion"},{"type":"required_linear_history"},{"type":"non_fast_forward"}]`,
		},
		{
			name: "RepositoryRulesetRules with Update rule",
			input: RepositoryRulesetRules{
				Update: &UpdateRuleParameters{
					UpdateAllowsFetchAndMerge: true,
				},
			},
			expected: `[{"type":"update","parameters":{"update_allows_fetch_and_merge":true}}]`,
		},
		{
			name: "RepositoryRulesetRules with RequiredDeployments",
			input: RepositoryRulesetRules{
				RequiredDeployments: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"production", "staging"},
				},
			},
			expected: `[{"type":"required_deployments","parameters":{"required_deployment_environments":["production","staging"]}}]`,
		},
		{
			name: "RepositoryRulesetRules with RequiredSignatures",
			input: RepositoryRulesetRules{
				RequiredSignatures: &EmptyRuleParameters{},
			},
			expected: `[{"type":"required_signatures"}]`,
		},
		{
			name: "RepositoryRulesetRules with PullRequest rule",
			input: RepositoryRulesetRules{
				PullRequest: &PullRequestRuleParameters{
					AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
					DismissStaleReviewsOnPush:      true,
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   2,
					RequiredReviewThreadResolution: true,
				},
			},
			expected: `[{"type":"pull_request","parameters":{"allowed_merge_methods":["rebase","squash"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":2,"required_review_thread_resolution":true}}]`,
		},
		{
			name: "RepositoryRulesetRules with RequiredStatusChecks",
			input: RepositoryRulesetRules{
				RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
					DoNotEnforceOnCreate: Ptr(true),
					RequiredStatusChecks: []*RuleStatusCheck{
						{
							Context:       "build",
							IntegrationID: Ptr(int64(1)),
						},
						{
							Context:       "lint",
							IntegrationID: Ptr(int64(2)),
						},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
			},
			expected: `[{"type":"required_status_checks","parameters":{"do_not_enforce_on_create":true,"required_status_checks":[{"context":"build","integration_id":1},{"context":"lint","integration_id":2}],"strict_required_status_checks_policy":true}}]`,
		},
		{
			name: "RepositoryRulesetRules with CommitMessagePattern",
			input: RepositoryRulesetRules{
				CommitMessagePattern: &PatternRuleParameters{
					Name:     Ptr("avoid test commits"),
					Negate:   Ptr(true),
					Operator: "starts_with",
					Pattern:  "[test]",
				},
			},
			expected: `[{"type":"commit_message_pattern","parameters":{"name":"avoid test commits","negate":true,"operator":"starts_with","pattern":"[test]"}}]`,
		},
		{
			name: "RepositoryRulesetRules with CommitAuthorEmailPattern",
			input: RepositoryRulesetRules{
				CommitAuthorEmailPattern: &PatternRuleParameters{
					Operator: "contains",
					Pattern:  "example.com",
				},
			},
			expected: `[{"type":"commit_author_email_pattern","parameters":{"operator":"contains","pattern":"example.com"}}]`,
		},
		{
			name: "RepositoryRulesetRules with CommitterEmailPattern",
			input: RepositoryRulesetRules{
				CommitterEmailPattern: &PatternRuleParameters{
					Name:     Ptr("require org email"),
					Operator: "ends_with",
					Pattern:  "@company.com",
				},
			},
			expected: `[{"type":"committer_email_pattern","parameters":{"name":"require org email","operator":"ends_with","pattern":"@company.com"}}]`,
		},
		{
			name: "RepositoryRulesetRules with BranchNamePattern",
			input: RepositoryRulesetRules{
				BranchNamePattern: &PatternRuleParameters{
					Name:     Ptr("enforce naming convention"),
					Negate:   Ptr(false),
					Operator: "regex",
					Pattern:  "^(main|develop|release/)",
				},
			},
			expected: `[{"type":"branch_name_pattern","parameters":{"name":"enforce naming convention","negate":false,"operator":"regex","pattern":"^(main|develop|release/)"}}]`,
		},
		{
			name: "RepositoryRulesetRules with TagNamePattern",
			input: RepositoryRulesetRules{
				TagNamePattern: &PatternRuleParameters{
					Operator: "contains",
					Pattern:  "v",
				},
			},
			expected: `[{"type":"tag_name_pattern","parameters":{"operator":"contains","pattern":"v"}}]`,
		},
		{
			name: "RepositoryRulesetRules with CodeScanning",
			input: RepositoryRulesetRules{
				CodeScanning: &CodeScanningRuleParameters{
					CodeScanningTools: []*RuleCodeScanningTool{
						{
							Tool:                    "CodeQL",
							AlertsThreshold:         CodeScanningAlertsThresholdErrors,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						},
					},
				},
			},
			expected: `[{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","alerts_threshold":"errors","security_alerts_threshold":"high_or_higher"}]}}]`,
		},
		{
			name: "RepositoryRulesetRules with all rules populated",
			input: RepositoryRulesetRules{
				Creation: &EmptyRuleParameters{},
				Update: &UpdateRuleParameters{
					UpdateAllowsFetchAndMerge: true,
				},
				Deletion:              &EmptyRuleParameters{},
				RequiredLinearHistory: &EmptyRuleParameters{},
				RequiredDeployments: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"production"},
				},
				RequiredSignatures: &EmptyRuleParameters{},
				PullRequest: &PullRequestRuleParameters{
					AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase},
					DismissStaleReviewsOnPush:      true,
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   1,
					RequiredReviewThreadResolution: true,
				},
				RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
					DoNotEnforceOnCreate: Ptr(true),
					RequiredStatusChecks: []*RuleStatusCheck{
						{
							Context:       "build",
							IntegrationID: Ptr(int64(1)),
						},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
				NonFastForward: &EmptyRuleParameters{},
				CommitMessagePattern: &PatternRuleParameters{
					Name:     Ptr("commit message"),
					Operator: "contains",
					Pattern:  "required",
				},
				CommitAuthorEmailPattern: &PatternRuleParameters{
					Operator: "contains",
					Pattern:  "example.com",
				},
				CommitterEmailPattern: &PatternRuleParameters{
					Operator: "ends_with",
					Pattern:  "@example.com",
				},
				BranchNamePattern: &PatternRuleParameters{
					Operator: "regex",
					Pattern:  "^main$",
				},
				TagNamePattern: &PatternRuleParameters{
					Operator: "contains",
					Pattern:  "v",
				},
				CodeScanning: &CodeScanningRuleParameters{
					CodeScanningTools: []*RuleCodeScanningTool{
						{
							Tool:                    "CodeQL",
							AlertsThreshold:         CodeScanningAlertsThresholdErrors,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						},
					},
				},
			},
			expected: `[{"type":"creation"},{"type":"update","parameters":{"update_allows_fetch_and_merge":true}},{"type":"deletion"},{"type":"required_linear_history"},{"type":"required_deployments","parameters":{"required_deployment_environments":["production"]}},{"type":"required_signatures"},{"type":"pull_request","parameters":{"allowed_merge_methods":["rebase"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":1,"required_review_thread_resolution":true}},{"type":"required_status_checks","parameters":{"do_not_enforce_on_create":true,"required_status_checks":[{"context":"build","integration_id":1}],"strict_required_status_checks_policy":true}},{"type":"non_fast_forward"},{"type":"commit_message_pattern","parameters":{"name":"commit message","operator":"contains","pattern":"required"}},{"type":"commit_author_email_pattern","parameters":{"operator":"contains","pattern":"example.com"}},{"type":"committer_email_pattern","parameters":{"operator":"ends_with","pattern":"@example.com"}},{"type":"branch_name_pattern","parameters":{"operator":"regex","pattern":"^main$"}},{"type":"tag_name_pattern","parameters":{"operator":"contains","pattern":"v"}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","alerts_threshold":"errors","security_alerts_threshold":"high_or_higher"}]}}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the input to JSON using a pointer to ensure MarshalJSON is called
			data, err := json.Marshal(&tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("MarshalJSON error = %v, expectError %v", err, tt.expectError)
				return
			}

			if err != nil {
				return
			}

			// Parse both expected and actual as JSON arrays for normalized comparison
			var expectedArray, actualArray interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedArray); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			if err := json.Unmarshal(data, &actualArray); err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			// Convert back to JSON strings for comparison
			expectedJSON, _ := json.Marshal(expectedArray)
			actualJSON, _ := json.Marshal(actualArray)

			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("MarshalJSON() = %s, want %s", string(actualJSON), string(expectedJSON))
			}
		})
	}
}
