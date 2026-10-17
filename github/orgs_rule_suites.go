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
	RuleSource  *RuleEvaluationSource `json:"rule_source,omitempty"`
	Enforcement *string               `json:"enforcement,omitempty"`
	Result      *string               `json:"result,omitempty"`
	RuleType    *string               `json:"rule_type,omitempty"`
	Details     *string               `json:"details,omitempty"`
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
	Ref             string `url:"ref,omitempty"`
	RepositoryName  string `url:"repository_name,omitempty"`
	TimePeriod      string `url:"time_period,omitempty"`
	ActorName       string `url:"actor_name,omitempty"`
	RuleSuiteResult string `url:"rule_suite_result,omitempty"`
	EvaluateStatus  string `url:"evaluate_status,omitempty"`

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

// GetRuleSuite gets a rule suite for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rule-suites?apiVersion=2022-11-28#get-an-organization-rule-suite
//
//meta:operation GET /orgs/{org}/rulesets/rule-suites/{rule_suite_id}
func (s *OrganizationsService) GetRuleSuite(ctx context.Context, org string, ruleSuiteID int64) (*RuleSuite, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/rule-suites/%v", org, ruleSuiteID)

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
