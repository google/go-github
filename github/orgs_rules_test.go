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

func TestOrganizationsService_GetAllOrganizationRepositoryRulesets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			}
		}]`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetAllOrganizationRepositoryRulesets(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetAllOrganizationRepositoryRulesets returned error: %v", err)
	}

	want := []*Ruleset{{
		ID:          26110,
		Name:        "test ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassMode:  String("none"),
		NodeID:      String("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: String("https://api.github.com/orgs/o/rulesets/26110")},
		},
	}}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetAllOrganizationRepositoryRulesets returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetAllOrganizationRepositoryRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAllOrganizationRepositoryRulesets(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrganizationRepositoryRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "ref_name": {
				"include": [
				  "refs/heads/main",
				  "refs/heads/master"
				],
				"exclude": [
				  "refs/heads/dev*"
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
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateOrganizationRepositoryRuleset(ctx, "o", &Ruleset{
		ID:          21,
		Name:        "ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: &[]BypassActor{
			{
				ActorID:   234,
				ActorType: "Team",
			},
		},
		Conditions: &RulesetCondition{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Bool(true),
			},
		},
		Rules: &[]RulesetRule{
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
						IntegrationID: Int64(1),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     String("avoid test commits"),
				Negate:   Bool(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     String("avoid commit emails"),
				Negate:   Bool(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     String("avoid branch names"),
				Negate:   Bool(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     String("avoid tag names"),
				Negate:   Bool(true),
				Operator: "contains",
				Pattern:  "github",
			}),
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          21,
		Name:        "ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: &[]BypassActor{
			{
				ActorID:   234,
				ActorType: "Team",
			},
		},
		Conditions: &RulesetCondition{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Bool(true),
			},
		},
		Rules: &[]RulesetRule{
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
						IntegrationID: Int64(1),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     String("avoid test commits"),
				Negate:   Bool(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     String("avoid commit emails"),
				Negate:   Bool(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     String("avoid branch names"),
				Negate:   Bool(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     String("avoid tag names"),
				Negate:   Bool(true),
				Operator: "contains",
				Pattern:  "github",
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateOrganizationRepositoryRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateOrganizationRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrganizationRepositoryRuleset(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetOrganizationRepositoryRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
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
	rulesets, _, err := client.Organizations.GetOrganizationRepositoryRuleset(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          26110,
		Name:        "test ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassMode:  String("none"),
		NodeID:      String("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: String("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetCondition{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Bool(true),
			},
		},
		Rules: &[]RulesetRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetOrganizationRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrganizationRepositoryRuleset(ctx, "o", 26110)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateOrganizationRepositoryRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
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
	rulesets, _, err := client.Organizations.UpdateOrganizationRepositoryRuleset(ctx, "o", 26110, &Ruleset{
		Name:        "test ruleset",
		Target:      String("branch"),
		Enforcement: "active",
		BypassMode:  String("none"),
		Conditions: &RulesetCondition{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Bool(true),
			},
		},
		Rules: &[]RulesetRule{
			NewCreationRule(),
		},
	})

	if err != nil {
		t.Errorf("Organizations.UpdateOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          26110,
		Name:        "test ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassMode:  String("none"),
		NodeID:      String("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: String("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetCondition{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Bool(true),
			},
		},
		Rules: &[]RulesetRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.UpdateOrganizationRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateOrganizationRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateOrganizationRepositoryRuleset(ctx, "o", 26110, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteOrganizationRepositoryRuleset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.DeleteOrganizationRepositoryRuleset(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.DeleteOrganizationRepositoryRuleset returned error: %v", err)
	}

	const methodName = "DeleteOrganizationRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.DeleteOrganizationRepositoryRuleset(ctx, "0", 26110)
	})
}
