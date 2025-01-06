// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_CreateEnterpriseRuleset_OrgNameRepoName(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_actors": [
				{
					"actor_id": 234,
					"actor_type": "Team"
				}
			],
			"conditions": {
				"organization_name": {
					"include": [
						"important_organization",
						"another_important_organization"
					],
					"exclude": [
						"unimportant_organization"
					]
				},
			  "repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
					],
					"exclude": [
						"unimportant_repository"
					],
					"protected": true
				},
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
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", &Ruleset{
		Name:        "ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateEnterpriseRuleset_OrgNameRepoProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_actors": [
				{
					"actor_id": 234,
					"actor_type": "Team"
				}
			],
			"conditions": {
				"organization_name": {
					"include": [
						"important_organization",
						"another_important_organization"
					],
					"exclude": [
						"unimportant_organization"
					]
				},
			  "repository_property": {
					"include": [
						{
							"name": "testIncludeProp",
							"source": "custom",
							"property_values": [
								"true"
							]
						}
					],
					"exclude": [
						{
							"name": "testExcludeProp",
							"property_values": [
								"false"
							]
						}
					]
				},
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
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", &Ruleset{
		Name:        "ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateEnterpriseRuleset_OrgIdRepoName(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_actors": [
				{
					"actor_id": 234,
					"actor_type": "Team"
				}
			],
			"conditions": {
				"organization_id": {
					"organization_ids": [1001, 1002]
				},
			  "repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
					],
					"exclude": [
						"unimportant_repository"
					],
					"protected": true
				},
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
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", &Ruleset{
		Name:        "ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationID: &RulesetOrganizationIDsConditionParameters{
				OrganizationIDs: []int64{1001, 1002},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationID: &RulesetOrganizationIDsConditionParameters{
				OrganizationIDs: []int64{1001, 1002},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateEnterpriseRuleset_OrgIdRepoProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_actors": [
				{
					"actor_id": 234,
					"actor_type": "Team"
				}
			],
			"conditions": {
				"organization_id": {
					"organization_ids": [1001, 1002]
				},
			  "repository_property": {
					"include": [
						{
							"name": "testIncludeProp",
							"source": "custom",
							"property_values": [
								"true"
							]
						}
					],
					"exclude": [
						{
							"name": "testExcludeProp",
							"property_values": [
								"false"
							]
						}
					]
				},
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
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", &Ruleset{
		Name:        "ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationID: &RulesetOrganizationIDsConditionParameters{
				OrganizationIDs: []int64{1001, 1002},
			},
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			OrganizationID: &RulesetOrganizationIDsConditionParameters{
				OrganizationIDs: []int64{1001, 1002},
			},
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
				DismissStaleReviewsOnPush:      true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Enterprise.CreateEnterpriseRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseRuleset(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetEnterpriseRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
					"href": "https://api.github.com/enterprises/e/rulesets/26110"
				}
			},
			"conditions": {
				"organization_name": {
					"include": [
						"important_organization",
						"another_important_organization"
					],
					"exclude": [
						"unimportant_organization"
					]
				},
				"repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
					],
					"exclude": [
						"unimportant_repository"
				  ],
				  "protected": true
				},
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
	rulesets, _, err := client.Enterprise.GetEnterpriseRuleset(ctx, "e", 26110)
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/enterprises/e/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Enterprise.GetEnterpriseRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseRuleset(ctx, "e", 26110)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateEnterpriseRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Enterprise",
			"source": "e",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/enterprises/e/rulesets/26110"
			  }
			},
			"conditions": {
				"organization_name": {
					"include": [
						"important_organization",
						"another_important_organization"
					],
					"exclude": [
						"unimportant_organization"
					]
				},
				"repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
				  ],
					"exclude": [
						"unimportant_repository"
					],
				  "protected": true
				},
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
	rulesets, _, err := client.Enterprise.UpdateEnterpriseRuleset(ctx, "e", 26110, &Ruleset{
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	})
	if err != nil {
		t.Errorf("Enterprise.UpdateEnterpriseRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Enterprise"),
		Source:      "e",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/enterprises/e/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			OrganizationName: &RulesetOrganizationNamesConditionParameters{
				Include: []string{"important_organization", "another_important_organization"},
				Exclude: []string{"unimportant_organization"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Enterprise.UpdateEnterpriseRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateEnterpriseRuleset(ctx, "e", 26110, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteEnterpriseRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.DeleteEnterpriseRuleset(ctx, "e", 26110)
	if err != nil {
		t.Errorf("Enterprise.DeleteEnterpriseRuleset returned error: %v", err)
	}

	const methodName = "DeleteEnterpriseRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteEnterpriseRuleset(ctx, "e", 26110)
	})
}
