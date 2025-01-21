// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRulesetRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		rules *RepositoryRulesetRules
		json  string
	}{
		{"empty", &RepositoryRulesetRules{}, `[]`},
		{
			"single_rule_with_empty_params",
			&RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
			`[{"type":"creation"}]`,
		},
		{
			"single_rule_with_required_params",
			&RepositoryRulesetRules{
				RequiredDeployments: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test"},
				},
			},
			`[{"type":"required_deployments","parameters":{"required_deployment_environments":["test"]}}]`,
		},
		{
			"all_rules_with_required_params",
			&RepositoryRulesetRules{
				Creation:              &EmptyRuleParameters{},
				Update:                &UpdateRuleParameters{},
				Deletion:              &EmptyRuleParameters{},
				RequiredLinearHistory: &EmptyRuleParameters{},
				MergeQueue: &MergeQueueRuleParameters{
					CheckResponseTimeoutMinutes:  5,
					GroupingStrategy:             MergeGroupingStrategyAllGreen,
					MaxEntriesToBuild:            10,
					MaxEntriesToMerge:            20,
					MergeMethod:                  MergeMethodSquash,
					MinEntriesToMerge:            1,
					MinEntriesToMergeWaitMinutes: 15,
				},
				RequiredDeployments: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test1", "test2"},
				},
				RequiredSignatures: &EmptyRuleParameters{},
				PullRequest: &PullRequestRuleParameters{
					AllowedMergeMethods: []MergeMethod{
						MergeMethodSquash,
						MergeMethodRebase,
					},
					DismissStaleReviewsOnPush:      true,
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   2,
					RequiredReviewThreadResolution: true,
				},
				RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
					RequiredStatusChecks: []*RuleStatusCheck{
						{Context: "test1"},
						{Context: "test2"},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
				NonFastForward: &EmptyRuleParameters{},
				CommitMessagePattern: &PatternRuleParameters{
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				CommitAuthorEmailPattern: &PatternRuleParameters{
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				CommitterEmailPattern: &PatternRuleParameters{
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				BranchNamePattern: &PatternRuleParameters{
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				TagNamePattern: &PatternRuleParameters{
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				FilePathRestriction: &FilePathRestrictionRuleParameters{
					RestrictedFilePaths: []string{"test1", "test2"},
				},
				MaxFilePathLength: &MaxFilePathLengthRuleParameters{MaxFilePathLength: 512},
				FileExtensionRestriction: &FileExtensionRestrictionRuleParameters{
					RestrictedFileExtensions: []string{".exe", ".pkg"},
				},
				MaxFileSize: &MaxFileSizeRuleParameters{MaxFileSize: 1024},
				Workflows: &WorkflowsRuleParameters{
					Workflows: []*RuleWorkflow{
						{Path: ".github/workflows/test1.yaml"},
						{Path: ".github/workflows/test2.yaml"},
					},
				},
				CodeScanning: &CodeScanningRuleParameters{
					CodeScanningTools: []*RuleCodeScanningTool{
						{
							AlertsThreshold:         CodeScanningAlertsThresholdAll,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdAll,
							Tool:                    "test",
						},
						{
							AlertsThreshold:         CodeScanningAlertsThresholdNone,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdNone,
							Tool:                    "test",
						},
					},
				},
			},
			`[{"type":"creation"},{"type":"update"},{"type":"deletion"},{"type":"required_linear_history"},{"type":"merge_queue","parameters":{"check_response_timeout_minutes":5,"grouping_strategy":"ALLGREEN","max_entries_to_build":10,"max_entries_to_merge":20,"merge_method":"squash","min_entries_to_merge":1,"min_entries_to_merge_wait_minutes":15}},{"type":"required_deployments","parameters":{"required_deployment_environments":["test1","test2"]}},{"type":"required_signatures"},{"type":"pull_request","parameters":{"allowed_merge_methods":["squash","rebase"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":2,"required_review_thread_resolution":true}},{"type":"required_status_checks","parameters":{"required_status_checks":[{"context":"test1"},{"context":"test2"}],"strict_required_status_checks_policy":true}},{"type":"non_fast_forward"},{"type":"commit_message_pattern","parameters":{"operator":"starts_with","pattern":"test"}},{"type":"commit_author_email_pattern","parameters":{"operator":"starts_with","pattern":"test"}},{"type":"committer_email_pattern","parameters":{"operator":"starts_with","pattern":"test"}},{"type":"branch_name_pattern","parameters":{"operator":"starts_with","pattern":"test"}},{"type":"tag_name_pattern","parameters":{"operator":"starts_with","pattern":"test"}},{"type":"file_path_restriction","parameters":{"restricted_file_paths":["test1","test2"]}},{"type":"max_file_path_length","parameters":{"max_file_path_length":512}},{"type":"file_extension_restriction","parameters":{"restricted_file_extensions":[".exe",".pkg"]}},{"type":"max_file_size","parameters":{"max_file_size":1024}},{"type":"workflows","parameters":{"workflows":[{"path":".github/workflows/test1.yaml"},{"path":".github/workflows/test2.yaml"}]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"alerts_threshold":"all","security_alerts_threshold":"all","tool":"test"},{"alerts_threshold":"none","security_alerts_threshold":"none","tool":"test"}]}}]`,
		},
		{
			"all_rules_with_all_params",
			&RepositoryRulesetRules{
				Creation:              &EmptyRuleParameters{},
				Update:                &UpdateRuleParameters{UpdateAllowsFetchAndMerge: true},
				Deletion:              &EmptyRuleParameters{},
				RequiredLinearHistory: &EmptyRuleParameters{},
				MergeQueue: &MergeQueueRuleParameters{
					CheckResponseTimeoutMinutes:  5,
					GroupingStrategy:             MergeGroupingStrategyAllGreen,
					MaxEntriesToBuild:            10,
					MaxEntriesToMerge:            20,
					MergeMethod:                  MergeMethodSquash,
					MinEntriesToMerge:            1,
					MinEntriesToMergeWaitMinutes: 15,
				},
				RequiredDeployments: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test1", "test2"},
				},
				RequiredSignatures: &EmptyRuleParameters{},
				PullRequest: &PullRequestRuleParameters{
					AllowedMergeMethods: []MergeMethod{
						MergeMethodSquash,
						MergeMethodRebase,
					},
					DismissStaleReviewsOnPush:      true,
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   2,
					RequiredReviewThreadResolution: true,
				},
				RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
					DoNotEnforceOnCreate: Ptr(true),
					RequiredStatusChecks: []*RuleStatusCheck{
						{Context: "test1", IntegrationID: Ptr(int64(1))},
						{Context: "test2", IntegrationID: Ptr(int64(2))},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
				NonFastForward: &EmptyRuleParameters{},
				CommitMessagePattern: &PatternRuleParameters{
					Name:     Ptr("cmp"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				CommitAuthorEmailPattern: &PatternRuleParameters{
					Name:     Ptr("caep"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				CommitterEmailPattern: &PatternRuleParameters{
					Name:     Ptr("cep"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				BranchNamePattern: &PatternRuleParameters{
					Name:     Ptr("bp"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				TagNamePattern: &PatternRuleParameters{
					Name:     Ptr("tp"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
				FilePathRestriction: &FilePathRestrictionRuleParameters{
					RestrictedFilePaths: []string{"test1", "test2"},
				},
				MaxFilePathLength: &MaxFilePathLengthRuleParameters{MaxFilePathLength: 512},
				FileExtensionRestriction: &FileExtensionRestrictionRuleParameters{
					RestrictedFileExtensions: []string{".exe", ".pkg"},
				},
				MaxFileSize: &MaxFileSizeRuleParameters{MaxFileSize: 1024},
				Workflows: &WorkflowsRuleParameters{
					DoNotEnforceOnCreate: Ptr(true),
					Workflows: []*RuleWorkflow{
						{
							Path:         ".github/workflows/test1.yaml",
							Ref:          Ptr("main"),
							RepositoryID: Ptr(int64(1)),
							SHA:          Ptr("aaaa"),
						},
						{
							Path:         ".github/workflows/test2.yaml",
							Ref:          Ptr("main"),
							RepositoryID: Ptr(int64(2)),
							SHA:          Ptr("bbbb"),
						},
					},
				},
				CodeScanning: &CodeScanningRuleParameters{
					CodeScanningTools: []*RuleCodeScanningTool{
						{
							AlertsThreshold:         CodeScanningAlertsThresholdAll,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdAll,
							Tool:                    "test",
						},
						{
							AlertsThreshold:         CodeScanningAlertsThresholdNone,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdNone,
							Tool:                    "test",
						},
					},
				},
			},
			`[{"type":"creation"},{"type":"update","parameters":{"update_allows_fetch_and_merge":true}},{"type":"deletion"},{"type":"required_linear_history"},{"type":"merge_queue","parameters":{"check_response_timeout_minutes":5,"grouping_strategy":"ALLGREEN","max_entries_to_build":10,"max_entries_to_merge":20,"merge_method":"squash","min_entries_to_merge":1,"min_entries_to_merge_wait_minutes":15}},{"type":"required_deployments","parameters":{"required_deployment_environments":["test1","test2"]}},{"type":"required_signatures"},{"type":"pull_request","parameters":{"allowed_merge_methods":["squash","rebase"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":2,"required_review_thread_resolution":true}},{"type":"required_status_checks","parameters":{"do_not_enforce_on_create":true,"required_status_checks":[{"context":"test1","integration_id":1},{"context":"test2","integration_id":2}],"strict_required_status_checks_policy":true}},{"type":"non_fast_forward"},{"type":"commit_message_pattern","parameters":{"name":"cmp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"commit_author_email_pattern","parameters":{"name":"caep","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"committer_email_pattern","parameters":{"name":"cep","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"branch_name_pattern","parameters":{"name":"bp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"tag_name_pattern","parameters":{"name":"tp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"file_path_restriction","parameters":{"restricted_file_paths":["test1","test2"]}},{"type":"max_file_path_length","parameters":{"max_file_path_length":512}},{"type":"file_extension_restriction","parameters":{"restricted_file_extensions":[".exe",".pkg"]}},{"type":"max_file_size","parameters":{"max_file_size":1024}},{"type":"workflows","parameters":{"do_not_enforce_on_create":true,"workflows":[{"path":".github/workflows/test1.yaml","ref":"main","repository_id":1,"sha":"aaaa"},{"path":".github/workflows/test2.yaml","ref":"main","repository_id":2,"sha":"bbbb"}]}},{"type":"code_scanning","parameters":{"code_scanning_tools":[{"alerts_threshold":"all","security_alerts_threshold":"all","tool":"test"},{"alerts_threshold":"none","security_alerts_threshold":"none","tool":"test"}]}}]`,
		},
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got, err := json.Marshal(test.rules)
				if err != nil {
					t.Errorf("Unable to marshal JSON for %#v", test.rules)
				}

				if diff := cmp.Diff(test.json, string(got)); diff != "" {
					t.Errorf(
						"json.Marshal returned:\n%s\nwant:\n%s\ndiff:\n%v",
						got,
						test.json,
						diff,
					)
				}
			})
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := &RepositoryRulesetRules{}
				err := json.Unmarshal([]byte(test.json), got)
				if err != nil {
					t.Errorf("Unable to unmarshal JSON %v: %v", test.json, err)
				}

				if diff := cmp.Diff(test.rules, got); diff != "" {
					t.Errorf(
						"json.Unmarshal returned:\n%#v\nwant:\n%#v\ndiff:\n%v",
						got,
						test.rules,
						diff,
					)
				}
			})
		}
	})
}

func TestBranchRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		rules *BranchRules
		json  string
	}{
		{"empty", &BranchRules{}, `[]`},
		{
			"single_rule_type_single_rule_empty_params",
			&BranchRules{
				Creation: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
			},
			`[{"type":"creation","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1}]`,
		},
		{
			"single_rule_type_single_rule_with_params",
			&BranchRules{
				Update: []*UpdateBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: UpdateRuleParameters{UpdateAllowsFetchAndMerge: true},
					},
				},
			},
			`[{"type":"update","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"update_allows_fetch_and_merge":true}}]`,
		},
		{
			"all_rule_types_with_all_parameters",
			&BranchRules{
				Creation: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
				Update: []*UpdateBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: UpdateRuleParameters{UpdateAllowsFetchAndMerge: true},
					},
				},
				Deletion: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
				RequiredLinearHistory: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
				MergeQueue: []*MergeQueueBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: MergeQueueRuleParameters{
							CheckResponseTimeoutMinutes:  5,
							GroupingStrategy:             MergeGroupingStrategyAllGreen,
							MaxEntriesToBuild:            10,
							MaxEntriesToMerge:            20,
							MergeMethod:                  MergeMethodSquash,
							MinEntriesToMerge:            1,
							MinEntriesToMergeWaitMinutes: 15,
						},
					},
				},
				RequiredDeployments: []*RequiredDeploymentsBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: RequiredDeploymentsRuleParameters{
							RequiredDeploymentEnvironments: []string{"test1", "test2"},
						},
					},
				},
				RequiredSignatures: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
				PullRequest: []*PullRequestBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PullRequestRuleParameters{
							AllowedMergeMethods: []MergeMethod{
								MergeMethodSquash,
								MergeMethodRebase,
							},
							DismissStaleReviewsOnPush:      true,
							RequireCodeOwnerReview:         true,
							RequireLastPushApproval:        true,
							RequiredApprovingReviewCount:   2,
							RequiredReviewThreadResolution: true,
						},
					},
				},
				RequiredStatusChecks: []*RequiredStatusChecksBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: RequiredStatusChecksRuleParameters{
							DoNotEnforceOnCreate: Ptr(true),
							RequiredStatusChecks: []*RuleStatusCheck{
								{Context: "test1", IntegrationID: Ptr(int64(1))},
								{Context: "test2", IntegrationID: Ptr(int64(2))},
							},
							StrictRequiredStatusChecksPolicy: true,
						},
					},
				},
				NonFastForward: []*BranchRuleMetadata{
					{
						RulesetSourceType: RulesetSourceTypeRepository,
						RulesetSource:     "test/test",
						RulesetID:         1,
					},
				},
				CommitMessagePattern: []*PatternBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PatternRuleParameters{
							Name:     Ptr("cmp"),
							Negate:   Ptr(false),
							Operator: PatternRuleOperatorStartsWith,
							Pattern:  "test",
						},
					},
				},
				CommitAuthorEmailPattern: []*PatternBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PatternRuleParameters{
							Name:     Ptr("caep"),
							Negate:   Ptr(false),
							Operator: PatternRuleOperatorStartsWith,
							Pattern:  "test",
						},
					},
				},
				CommitterEmailPattern: []*PatternBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PatternRuleParameters{
							Name:     Ptr("cep"),
							Negate:   Ptr(false),
							Operator: PatternRuleOperatorStartsWith,
							Pattern:  "test",
						},
					},
				},
				BranchNamePattern: []*PatternBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PatternRuleParameters{
							Name:     Ptr("bp"),
							Negate:   Ptr(false),
							Operator: PatternRuleOperatorStartsWith,
							Pattern:  "test",
						},
					},
				},
				TagNamePattern: []*PatternBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: PatternRuleParameters{
							Name:     Ptr("tp"),
							Negate:   Ptr(false),
							Operator: PatternRuleOperatorStartsWith,
							Pattern:  "test",
						},
					},
				},
				FilePathRestriction: []*FilePathRestrictionBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: FilePathRestrictionRuleParameters{
							RestrictedFilePaths: []string{"test1", "test2"},
						},
					},
				},
				MaxFilePathLength: []*MaxFilePathLengthBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: MaxFilePathLengthRuleParameters{MaxFilePathLength: 512},
					},
				},
				FileExtensionRestriction: []*FileExtensionRestrictionBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: FileExtensionRestrictionRuleParameters{
							RestrictedFileExtensions: []string{".exe", ".pkg"},
						},
					},
				},
				MaxFileSize: []*MaxFileSizeBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: MaxFileSizeRuleParameters{MaxFileSize: 1024},
					},
				},
				Workflows: []*WorkflowsBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: WorkflowsRuleParameters{
							DoNotEnforceOnCreate: Ptr(true),
							Workflows: []*RuleWorkflow{
								{
									Path:         ".github/workflows/test1.yaml",
									Ref:          Ptr("main"),
									RepositoryID: Ptr(int64(1)),
									SHA:          Ptr("aaaa"),
								},
								{
									Path:         ".github/workflows/test2.yaml",
									Ref:          Ptr("main"),
									RepositoryID: Ptr(int64(2)),
									SHA:          Ptr("bbbb"),
								},
							},
						},
					},
				},
				CodeScanning: []*CodeScanningBranchRule{
					{
						BranchRuleMetadata: BranchRuleMetadata{
							RulesetSourceType: RulesetSourceTypeRepository,
							RulesetSource:     "test/test",
							RulesetID:         1,
						},
						Parameters: CodeScanningRuleParameters{
							CodeScanningTools: []*RuleCodeScanningTool{
								{
									AlertsThreshold:         CodeScanningAlertsThresholdAll,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdAll,
									Tool:                    "test",
								},
								{
									AlertsThreshold:         CodeScanningAlertsThresholdNone,
									SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdNone,
									Tool:                    "test",
								},
							},
						},
					},
				},
			},
			`[{"type":"creation","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1},{"type":"update","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"update_allows_fetch_and_merge":true}},{"type":"deletion","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1},{"type":"required_linear_history","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1},{"type":"merge_queue","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"check_response_timeout_minutes":5,"grouping_strategy":"ALLGREEN","max_entries_to_build":10,"max_entries_to_merge":20,"merge_method":"squash","min_entries_to_merge":1,"min_entries_to_merge_wait_minutes":15}},{"type":"required_deployments","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"required_deployment_environments":["test1","test2"]}},{"type":"required_signatures","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1},{"type":"pull_request","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"allowed_merge_methods":["squash","rebase"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":2,"required_review_thread_resolution":true}},{"type":"required_status_checks","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"do_not_enforce_on_create":true,"required_status_checks":[{"context":"test1","integration_id":1},{"context":"test2","integration_id":2}],"strict_required_status_checks_policy":true}},{"type":"non_fast_forward","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1},{"type":"commit_message_pattern","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"name":"cmp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"commit_author_email_pattern","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"name":"caep","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"committer_email_pattern","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"name":"cep","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"branch_name_pattern","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"name":"bp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"tag_name_pattern","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"name":"tp","negate":false,"operator":"starts_with","pattern":"test"}},{"type":"file_path_restriction","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"restricted_file_paths":["test1","test2"]}},{"type":"max_file_path_length","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"max_file_path_length":512}},{"type":"file_extension_restriction","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"restricted_file_extensions":[".exe",".pkg"]}},{"type":"max_file_size","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"max_file_size":1024}},{"type":"workflows","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"do_not_enforce_on_create":true,"workflows":[{"path":".github/workflows/test1.yaml","ref":"main","repository_id":1,"sha":"aaaa"},{"path":".github/workflows/test2.yaml","ref":"main","repository_id":2,"sha":"bbbb"}]}},{"type":"code_scanning","ruleset_source_type":"Repository","ruleset_source":"test/test","ruleset_id":1,"parameters":{"code_scanning_tools":[{"alerts_threshold":"all","security_alerts_threshold":"all","tool":"test"},{"alerts_threshold":"none","security_alerts_threshold":"none","tool":"test"}]}}]`,
		},
	}

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := &BranchRules{}
				err := json.Unmarshal([]byte(test.json), got)
				if err != nil {
					t.Errorf("Unable to unmarshal JSON %v: %v", test.json, err)
				}

				if diff := cmp.Diff(test.rules, got); diff != "" {
					t.Errorf(
						"json.Unmarshal returned:\n%#v\nwant:\n%#v\ndiff:\n%v",
						got,
						test.rules,
						diff,
					)
				}
			})
		}
	})
}

func TestRepositoryRule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		rule *RepositoryRule
		json string
	}{
		{
			"empty",
			&RepositoryRule{},
			`{}`,
		},
		{
			"creation",
			&RepositoryRule{Type: RulesetRuleTypeCreation, Parameters: nil},
			`{"type":"creation"}`,
		},
		{
			"update",
			&RepositoryRule{Type: RulesetRuleTypeUpdate, Parameters: &UpdateRuleParameters{}},
			`{"type":"update"}`,
		},
		{
			"update_params_empty",
			&RepositoryRule{Type: RulesetRuleTypeUpdate, Parameters: &UpdateRuleParameters{}},
			`{"type":"update","parameters":{}}`,
		},
		{
			"update_params_set",
			&RepositoryRule{
				Type:       RulesetRuleTypeUpdate,
				Parameters: &UpdateRuleParameters{UpdateAllowsFetchAndMerge: true},
			},
			`{"type":"update","parameters":{"update_allows_fetch_and_merge":true}}`,
		},
		{
			"deletion",
			&RepositoryRule{Type: RulesetRuleTypeDeletion, Parameters: nil},
			`{"type":"deletion"}`,
		},
		{
			"required_linear_history",
			&RepositoryRule{Type: RulesetRuleTypeRequiredLinearHistory, Parameters: nil},
			`{"type":"required_linear_history"}`,
		},
		{
			"merge_queue",
			&RepositoryRule{
				Type: RulesetRuleTypeMergeQueue,
				Parameters: &MergeQueueRuleParameters{
					CheckResponseTimeoutMinutes:  5,
					GroupingStrategy:             MergeGroupingStrategyAllGreen,
					MaxEntriesToBuild:            10,
					MaxEntriesToMerge:            20,
					MergeMethod:                  MergeMethodSquash,
					MinEntriesToMerge:            1,
					MinEntriesToMergeWaitMinutes: 15,
				},
			},
			`{"type":"merge_queue","parameters":{"check_response_timeout_minutes":5,"grouping_strategy":"ALLGREEN","max_entries_to_build":10,"max_entries_to_merge":20,"merge_method":"squash","min_entries_to_merge":1,"min_entries_to_merge_wait_minutes":15}}`,
		},
		{
			"required_deployments",
			&RepositoryRule{
				Type: RulesetRuleTypeRequiredDeployments,
				Parameters: &RequiredDeploymentsRuleParameters{
					RequiredDeploymentEnvironments: []string{"test1", "test2"},
				},
			},
			`{"type":"required_deployments","parameters":{"required_deployment_environments":["test1","test2"]}}`,
		},
		{
			"required_signatures",
			&RepositoryRule{Type: RulesetRuleTypeRequiredSignatures, Parameters: nil},
			`{"type":"required_signatures"}`,
		},
		{
			"pull_request",
			&RepositoryRule{
				Type: RulesetRuleTypePullRequest,
				Parameters: &PullRequestRuleParameters{
					AllowedMergeMethods: []MergeMethod{
						MergeMethodSquash,
						MergeMethodRebase,
					},
					DismissStaleReviewsOnPush:      true,
					RequireCodeOwnerReview:         true,
					RequireLastPushApproval:        true,
					RequiredApprovingReviewCount:   2,
					RequiredReviewThreadResolution: true,
				},
			},
			`{"type":"pull_request","parameters":{"allowed_merge_methods":["squash","rebase"],"dismiss_stale_reviews_on_push":true,"require_code_owner_review":true,"require_last_push_approval":true,"required_approving_review_count":2,"required_review_thread_resolution":true}}`,
		},
		{
			"required_status_checks",
			&RepositoryRule{
				Type: RulesetRuleTypeRequiredStatusChecks,
				Parameters: &RequiredStatusChecksRuleParameters{
					RequiredStatusChecks: []*RuleStatusCheck{
						{Context: "test1"},
						{Context: "test2"},
					},
					StrictRequiredStatusChecksPolicy: true,
				},
			},
			`{"type":"required_status_checks","parameters":{"required_status_checks":[{"context":"test1"},{"context":"test2"}],"strict_required_status_checks_policy":true}}`,
		},
		{
			"non_fast_forward",
			&RepositoryRule{Type: RulesetRuleTypeNonFastForward, Parameters: nil},
			`{"type":"non_fast_forward"}`,
		},
		{
			"commit_message_pattern",
			&RepositoryRule{
				Type: RulesetRuleTypeCommitMessagePattern,
				Parameters: &PatternRuleParameters{
					Name:     Ptr("test"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
			},
			`{"type":"commit_message_pattern","parameters":{"name":"test","negate":false,"operator":"starts_with","pattern":"test"}}`,
		},
		{
			"commit_author_email_pattern",
			&RepositoryRule{
				Type: RulesetRuleTypeCommitAuthorEmailPattern,
				Parameters: &PatternRuleParameters{
					Name:     Ptr("test"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
			},
			`{"type":"commit_author_email_pattern","parameters":{"name":"test","negate":false,"operator":"starts_with","pattern":"test"}}`,
		},
		{
			"committer_email_pattern",
			&RepositoryRule{
				Type: RulesetRuleTypeCommitterEmailPattern,
				Parameters: &PatternRuleParameters{
					Name:     Ptr("test"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
			},
			`{"type":"committer_email_pattern","parameters":{"name":"test","negate":false,"operator":"starts_with","pattern":"test"}}`,
		},
		{
			"branch_name_pattern",
			&RepositoryRule{
				Type: RulesetRuleTypeBranchNamePattern,
				Parameters: &PatternRuleParameters{
					Name:     Ptr("test"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
			},
			`{"type":"branch_name_pattern","parameters":{"name":"test","negate":false,"operator":"starts_with","pattern":"test"}}`,
		},
		{
			"tag_name_pattern",
			&RepositoryRule{
				Type: RulesetRuleTypeTagNamePattern,
				Parameters: &PatternRuleParameters{
					Name:     Ptr("test"),
					Negate:   Ptr(false),
					Operator: PatternRuleOperatorStartsWith,
					Pattern:  "test",
				},
			},
			`{"type":"tag_name_pattern","parameters":{"name":"test","negate":false,"operator":"starts_with","pattern":"test"}}`,
		},
		{
			"file_path_restriction",
			&RepositoryRule{
				Type: RulesetRuleTypeFilePathRestriction,
				Parameters: &FilePathRestrictionRuleParameters{
					RestrictedFilePaths: []string{"test1", "test2"},
				},
			},
			`{"type":"file_path_restriction","parameters":{"restricted_file_paths":["test1","test2"]}}`,
		},
		{
			"max_file_path_length",
			&RepositoryRule{
				Type:       RulesetRuleTypeMaxFilePathLength,
				Parameters: &MaxFilePathLengthRuleParameters{MaxFilePathLength: 512},
			},
			`{"type":"max_file_path_length","parameters":{"max_file_path_length":512}}`,
		},
		{
			"file_extension_restriction",
			&RepositoryRule{
				Type: RulesetRuleTypeFileExtensionRestriction,
				Parameters: &FileExtensionRestrictionRuleParameters{
					RestrictedFileExtensions: []string{".exe", ".pkg"},
				},
			},
			`{"type":"file_extension_restriction","parameters":{"restricted_file_extensions":[".exe",".pkg"]}}`,
		},
		{
			"max_file_size",
			&RepositoryRule{
				Type:       RulesetRuleTypeMaxFileSize,
				Parameters: &MaxFileSizeRuleParameters{MaxFileSize: 1024},
			},
			`{"type":"max_file_size","parameters":{"max_file_size":1024}}`,
		},
		{
			"workflows",
			&RepositoryRule{
				Type: RulesetRuleTypeWorkflows,
				Parameters: &WorkflowsRuleParameters{
					Workflows: []*RuleWorkflow{
						{Path: ".github/workflows/test1.yaml"},
						{Path: ".github/workflows/test2.yaml"},
					},
				},
			},
			`{"type":"workflows","parameters":{"workflows":[{"path":".github/workflows/test1.yaml"},{"path":".github/workflows/test2.yaml"}]}}`,
		},
		{
			"code_scanning",
			&RepositoryRule{
				Type: RulesetRuleTypeCodeScanning,
				Parameters: &CodeScanningRuleParameters{
					CodeScanningTools: []*RuleCodeScanningTool{
						{
							AlertsThreshold:         CodeScanningAlertsThresholdAll,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdAll,
							Tool:                    "test",
						},
						{
							AlertsThreshold:         CodeScanningAlertsThresholdNone,
							SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdNone,
							Tool:                    "test",
						},
					},
				},
			},
			`{"type":"code_scanning","parameters":{"code_scanning_tools":[{"alerts_threshold":"all","security_alerts_threshold":"all","tool":"test"},{"alerts_threshold":"none","security_alerts_threshold":"none","tool":"test"}]}}`,
		},
	}

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := &RepositoryRule{}
				err := json.Unmarshal([]byte(test.json), got)
				if err != nil {
					t.Errorf("Unable to unmarshal JSON %v: %v", test.json, err)
				}

				if diff := cmp.Diff(test.rule, got); diff != "" {
					t.Errorf(
						"json.Unmarshal returned:\n%#v\nwant:\n%#v\ndiff:\n%v",
						got,
						test.rule,
						diff,
					)
				}
			})
		}
	})
}
