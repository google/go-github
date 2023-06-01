// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
)

func TestRulesetRule_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in      *RulesetRule
		want    string
		wantErr bool
	}{
		{
			in: &RulesetRule{
				Type: "update",
				Parameters: &UpdateAllowsFetchAndMergeRuleParameters{
					UpdateAllowsFetchAndMerge: true,
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RulesetRule{
				Type: "required_deployments",
				Parameters: &RequiredDeploymentEnvironmentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test"},
				},
			},
			want:    `{`,
			wantErr: true,
		},
		{
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
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
			in: &RulesetRule{
				Type: "unknown",
			},
			want:    `{`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		err := json.Unmarshal([]byte(tc.want), tc.in)
		if err == nil && tc.wantErr {
			t.Errorf("RulesetRule.UnmarshalJSON returned nil instead of an error")
		}
		if err != nil && !tc.wantErr {
			t.Errorf("RulesetRule.UnmarshalJSON returned an unexpected error: %+v", err)
		}
	}

}
