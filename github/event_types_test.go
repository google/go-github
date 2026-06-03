// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryRulesetEvent_Unmarshal(t *testing.T) {
	t.Parallel()

	enterprise := &Enterprise{
		ID:     Ptr(1),
		NodeID: Ptr("n"),
		Slug:   Ptr("e"),
		Name:   Ptr("e"),
	}

	installation := &Installation{
		ID:      Ptr(int64(1)),
		NodeID:  Ptr("n"),
		AppID:   Ptr(int64(1)),
		AppSlug: Ptr("a"),
	}

	organization := &Organization{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("n"),
		Name:   Ptr("o"),
	}

	repository := &Repository{
		ID:       Ptr(int64(1)),
		NodeID:   Ptr("n"),
		Name:     Ptr("r"),
		FullName: Ptr("o/r"),
	}

	sender := &User{
		ID:     Ptr(int64(1)),
		NodeID: Ptr("n"),
		Login:  Ptr("l"),
	}

	tests := []struct {
		name  string
		json  string
		event *RepositoryRulesetEvent
	}{
		{"empty", `{}`, &RepositoryRulesetEvent{}},
		{
			"created",
			fmt.Sprintf(
				`{"action":"created","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~ALL"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type":"required_linear_history"},{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase","merge"]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"high_or_higher","alerts_threshold":"errors"}]}}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("created"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~ALL"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:              &EmptyRuleParameters{},
						Update:                &UpdateRuleParameters{},
						Deletion:              &EmptyRuleParameters{},
						RequiredLinearHistory: &EmptyRuleParameters{},
						PullRequest: &PullRequestRuleParameters{
							AllowedMergeMethods: []PullRequestMergeMethod{
								PullRequestMergeMethodSquash,
								PullRequestMergeMethodRebase,
								PullRequestMergeMethodMerge,
							},
							DismissStaleReviewsOnPush:      false,
							RequireCodeOwnerReview:         false,
							RequireLastPushApproval:        false,
							RequiredApprovingReviewCount:   2,
							RequiredReviewThreadResolution: false,
						},
						CodeScanning: &CodeScanningRuleParameters{
							CodeScanningTools: []*RuleCodeScanningTool{
								{
									AlertsThreshold:         CodeScanningAlertsThresholdErrors,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
									Tool:                    "CodeQL",
								},
							},
						},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
		{
			"edited",
			fmt.Sprintf(
				`{"action":"edited","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type": "required_signatures"},{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase"]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"medium_or_higher","alerts_threshold":"errors"}]}}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"changes":{"rules":{"added":[{"type": "required_signatures"}],"updated":[{"rule":{"type":"pull_request","parameters":{"required_approving_review_count":2,"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_last_push_approval":false,"required_review_thread_resolution":false,"allowed_merge_methods":["squash","rebase"]}},"changes":{"configuration":{"from":"{\\\"required_reviewers\\\":[],\\\"allowed_merge_methods\\\":[\\\"squash\\\",\\\"rebase\\\",\\\"merge\\\"],\\\"require_code_owner_review\\\":false,\\\"require_last_push_approval\\\":false,\\\"dismiss_stale_reviews_on_push\\\":false,\\\"required_approving_review_count\\\":2,\\\"authorized_dismissal_actors_only\\\":false,\\\"required_review_thread_resolution\\\":false,\\\"ignore_approvals_from_contributors\\\":false}"}}},{"rule":{"type":"code_scanning","parameters":{"code_scanning_tools":[{"tool":"CodeQL","security_alerts_threshold":"medium_or_higher","alerts_threshold":"errors"}]}},"changes":{"configuration":{"from":"{\\\"code_scanning_tools\\\":[{\\\"tool\\\":\\\"CodeQL\\\",\\\"alerts_threshold\\\":\\\"errors\\\",\\\"security_alerts_threshold\\\":\\\"high_or_higher\\\"}]}"}}}],"deleted":[{"type":"required_linear_history"}]},"conditions":{"updated":[{"condition":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"changes":{"include":{"from":["~ALL"]}}}],"deleted":[]}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("edited"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:           &EmptyRuleParameters{},
						Update:             &UpdateRuleParameters{},
						Deletion:           &EmptyRuleParameters{},
						RequiredSignatures: &EmptyRuleParameters{},
						PullRequest: &PullRequestRuleParameters{
							AllowedMergeMethods: []PullRequestMergeMethod{
								PullRequestMergeMethodSquash,
								PullRequestMergeMethodRebase,
							},
							DismissStaleReviewsOnPush:      false,
							RequireCodeOwnerReview:         false,
							RequireLastPushApproval:        false,
							RequiredApprovingReviewCount:   2,
							RequiredReviewThreadResolution: false,
						},
						CodeScanning: &CodeScanningRuleParameters{
							CodeScanningTools: []*RuleCodeScanningTool{
								{
									AlertsThreshold:         CodeScanningAlertsThresholdErrors,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdMediumOrHigher,
									Tool:                    "CodeQL",
								},
							},
						},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Changes: &RepositoryRulesetChanges{
					Conditions: &RepositoryRulesetChangedConditions{
						Updated: []*RepositoryRulesetUpdatedConditions{
							{
								Condition: &RepositoryRulesetConditions{
									RefName: &RepositoryRulesetRefConditionParameters{
										Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
										Exclude: []string{},
									},
								},
								Changes: &RepositoryRulesetUpdatedCondition{
									Include: &RepositoryRulesetChangeSources{
										From: []string{"~ALL"},
									},
								},
							},
						},
						Deleted: []*RepositoryRulesetConditions{},
					},
					Rules: &RepositoryRulesetChangedRules{
						Added: []*RepositoryRule{{Type: RulesetRuleTypeRequiredSignatures}},
						Updated: []*RepositoryRulesetUpdatedRules{
							{
								Rule: &RepositoryRule{
									Type: RulesetRuleTypePullRequest,
									Parameters: &PullRequestRuleParameters{
										AllowedMergeMethods: []PullRequestMergeMethod{
											PullRequestMergeMethodSquash,
											PullRequestMergeMethodRebase,
										},
										DismissStaleReviewsOnPush:      false,
										RequireCodeOwnerReview:         false,
										RequireLastPushApproval:        false,
										RequiredApprovingReviewCount:   2,
										RequiredReviewThreadResolution: false,
									},
								},
								Changes: &RepositoryRulesetChangedRule{
									Configuration: &RepositoryRulesetChangeSource{
										From: Ptr(
											`{\"required_reviewers\":[],\"allowed_merge_methods\":[\"squash\",\"rebase\",\"merge\"],\"require_code_owner_review\":false,\"require_last_push_approval\":false,\"dismiss_stale_reviews_on_push\":false,\"required_approving_review_count\":2,\"authorized_dismissal_actors_only\":false,\"required_review_thread_resolution\":false,\"ignore_approvals_from_contributors\":false}`,
										),
									},
								},
							},
							{
								Rule: &RepositoryRule{
									Type: RulesetRuleTypeCodeScanning,
									Parameters: &CodeScanningRuleParameters{
										CodeScanningTools: []*RuleCodeScanningTool{
											{
												AlertsThreshold:         CodeScanningAlertsThresholdErrors,
												SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdMediumOrHigher,
												Tool:                    "CodeQL",
											},
										},
									},
								},
								Changes: &RepositoryRulesetChangedRule{
									Configuration: &RepositoryRulesetChangeSource{
										From: Ptr(
											`{\"code_scanning_tools\":[{\"tool\":\"CodeQL\",\"alerts_threshold\":\"errors\",\"security_alerts_threshold\":\"high_or_higher\"}]}`,
										),
									},
								},
							},
						},
						Deleted: []*RepositoryRule{{Type: RulesetRuleTypeRequiredLinearHistory}},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
		{
			"deleted",
			fmt.Sprintf(
				`{"action":"deleted","repository_ruleset":{"id":1,"name":"r","target":"branch","source_type":"Repository","source":"o/r","enforcement":"active","conditions":{"ref_name":{"exclude":[],"include":["~DEFAULT_BRANCH","refs/heads/dev-*"]}},"rules":[{"type":"deletion"},{"type":"creation"},{"type":"update"},{"type":"required_linear_history"}],"node_id":"n","created_at":%[1]s,"updated_at":%[1]s,"_links":{"self":{"href":"a"},"html":{"href":"a"}}},"repository":{"id":1,"node_id":"n","name":"r","full_name":"o/r"},"organization":{"id":1,"node_id":"n","name":"o"},"enterprise":{"id":1,"node_id":"n","slug":"e","name":"e"},"installation":{"id":1,"node_id":"n","app_id":1,"app_slug":"a"},"sender":{"id":1,"node_id":"n","login":"l"}}`,
				referenceTimeStr,
			),
			&RepositoryRulesetEvent{
				Action: Ptr("deleted"),
				RepositoryRuleset: &RepositoryRuleset{
					ID:          Ptr(int64(1)),
					Name:        "r",
					Target:      Ptr(RulesetTargetBranch),
					SourceType:  Ptr(RulesetSourceTypeRepository),
					Source:      "o/r",
					Enforcement: RulesetEnforcementActive,
					Conditions: &RepositoryRulesetConditions{
						RefName: &RepositoryRulesetRefConditionParameters{
							Include: []string{"~DEFAULT_BRANCH", "refs/heads/dev-*"},
							Exclude: []string{},
						},
					},
					Rules: &RepositoryRulesetRules{
						Creation:              &EmptyRuleParameters{},
						Update:                &UpdateRuleParameters{},
						Deletion:              &EmptyRuleParameters{},
						RequiredLinearHistory: &EmptyRuleParameters{},
					},
					NodeID:    Ptr("n"),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
					Links: &RepositoryRulesetLinks{
						Self: &RepositoryRulesetLink{HRef: Ptr("a")},
						HTML: &RepositoryRulesetLink{HRef: Ptr("a")},
					},
				},
				Repository:   repository,
				Organization: organization,
				Enterprise:   enterprise,
				Installation: installation,
				Sender:       sender,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := &RepositoryRulesetEvent{}
			err := json.Unmarshal([]byte(test.json), got)
			if err != nil {
				t.Errorf("Unable to unmarshal JSON %v: %v", test.json, err)
			}

			if diff := cmp.Diff(test.event, got); diff != "" {
				t.Errorf("json.Unmarshal returned:\n%#v\nwant:\n%#v\ndiff:\n%v", got, test.event, diff)
			}
		})
	}
}
