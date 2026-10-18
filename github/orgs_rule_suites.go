// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// RuleSuite represents a rule suite: the evaluation of the rules that apply to a push.
type RuleSuite struct {
	ID               *int64            `json:"id,omitempty"`
	ActorID          *int64            `json:"actor_id,omitempty"`
	ActorName        *string           `json:"actor_name,omitempty"`
	BeforeSHA        *string           `json:"before_sha,omitempty"`
	AfterSHA         *string           `json:"after_sha,omitempty"`
	Ref              *string           `json:"ref,omitempty"`
	RepositoryID     *int64            `json:"repository_id,omitempty"`
	RepositoryName   *string           `json:"repository_name,omitempty"`
	PushedAt         *Timestamp        `json:"pushed_at,omitempty"`
	Result           *string           `json:"result,omitempty"`
	EvaluationResult *string           `json:"evaluation_result,omitempty"`
	RuleEvaluations  []*RuleEvaluation `json:"rule_evaluations,omitempty"`
}

// RuleEvaluation represents the result of evaluating a single rule within a rule suite.
type RuleEvaluation struct {
	RuleSource  *RuleEvaluationSource   `json:"rule_source,omitempty"`
	Enforcement *string                 `json:"enforcement,omitempty"`
	Result      *string                 `json:"result,omitempty"`
	RuleType    *string                 `json:"rule_type,omitempty"`
	Details     *string                 `json:"details,omitempty"`
	Metadata    *RuleEvaluationMetadata `json:"metadata,omitempty"`
}

// RuleEvaluationMetadata represents additional metadata about a rule evaluation.
// PullRequest is set for pull_request rules, and Checks is set for required_status_checks rules.
type RuleEvaluationMetadata struct {
	PullRequest *RuleEvaluationPullRequest   `json:"pull_request,omitempty"`
	Checks      []*RuleEvaluationStatusCheck `json:"checks,omitempty"`
}

// RuleEvaluationPullRequest represents the pull request associated with a rule evaluation.
type RuleEvaluationPullRequest struct {
	ID              *int64                         `json:"id,omitempty"`
	Number          *int                           `json:"number,omitempty"`
	User            *User                          `json:"user,omitempty"`
	Reviews         []*PullRequestReview           `json:"reviews,omitempty"`
	ProofOfPresence *RuleEvaluationProofOfPresence `json:"proof_of_presence,omitempty"`
}

// RuleEvaluationProofOfPresence represents the proof of presence evaluation for a pull request.
type RuleEvaluationProofOfPresence struct {
	State         *string `json:"state,omitempty"`
	Code          *string `json:"code,omitempty"`
	BaseSHA       *string `json:"base_sha,omitempty"`
	HeadSHA       *string `json:"head_sha,omitempty"`
	RequiredLevel *string `json:"required_level,omitempty"`
	VerifiedLevel *string `json:"verified_level,omitempty"`
	VerifiedAt    *string `json:"verified_at,omitempty"`
	Verifier      *User   `json:"verifier,omitempty"`
}

// RuleEvaluationStatusCheck represents a status check associated with a rule evaluation.
type RuleEvaluationStatusCheck struct {
	ID      *int64  `json:"id,omitempty"`
	Context *string `json:"context,omitempty"`
	State   *string `json:"state,omitempty"`
	Type    *string `json:"type,omitempty"`
	App     *App    `json:"app,omitempty"`
}

// RuleEvaluationSource represents the source of an evaluated rule.
type RuleEvaluationSource struct {
	Type *string `json:"type,omitempty"`
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ListRuleSuitesOptions specifies the optional parameters to the
// OrganizationsService.ListRuleSuites method.
type ListRuleSuitesOptions struct {
	// Ref filters rule suites by the name of the ref, which cannot contain wildcard characters.
	// Prefix it with "refs/heads/" to limit to branches or "refs/tags/" to limit to tags.
	Ref string `url:"ref,omitempty"`

	// RepositoryName filters rule suites by the name of the repository.
	RepositoryName string `url:"repository_name,omitempty"`

	// TimePeriod filters rule suites by the time period in which they occurred.
	// Can be one of: "hour", "day", "week", "month". Default: "day".
	TimePeriod string `url:"time_period,omitempty"`

	// ActorName filters rule suites by the handle of the user who triggered them.
	ActorName string `url:"actor_name,omitempty"`

	// RuleSuiteResult filters rule suites by their result.
	// Can be one of: "pass", "fail", "bypass", "all". Default: "all".
	RuleSuiteResult string `url:"rule_suite_result,omitempty"`

	// EvaluateStatus filters rule suites by the evaluate status of the rulesets that produced them.
	// Can be one of: "all", "active", "evaluate". Default: "all".
	EvaluateStatus string `url:"evaluate_status,omitempty"`

	ListOptions
}

// ListRuleSuites lists the rule suites for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rule-suites?apiVersion=2022-11-28#list-organization-rule-suites
//
//meta:operation GET /orgs/{org}/rulesets/rule-suites
func (s *OrganizationsService) ListRuleSuites(ctx context.Context, org string, opts *ListRuleSuitesOptions) ([]*RuleSuite, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/rule-suites", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ruleSuites []*RuleSuite
	resp, err := s.client.Do(req, &ruleSuites)
	if err != nil {
		return nil, resp, err
	}

	return ruleSuites, resp, nil
}

// GetRuleSuiteOptions specifies the optional parameters to the
// OrganizationsService.GetRuleSuite method.
type GetRuleSuiteOptions struct {
	// IncludeMetadata, when true, includes additional metadata for pull request
	// and required status check rule evaluations. Default: false.
	IncludeMetadata bool `url:"include_metadata,omitempty"`
}

// GetRuleSuite gets a rule suite for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rule-suites?apiVersion=2022-11-28#get-an-organization-rule-suite
//
//meta:operation GET /orgs/{org}/rulesets/rule-suites/{rule_suite_id}
func (s *OrganizationsService) GetRuleSuite(ctx context.Context, org string, ruleSuiteID int64, opts *GetRuleSuiteOptions) (*RuleSuite, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/rule-suites/%v", org, ruleSuiteID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ruleSuite *RuleSuite
	resp, err := s.client.Do(req, &ruleSuite)
	if err != nil {
		return nil, resp, err
	}

	return ruleSuite, resp, nil
}
