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
