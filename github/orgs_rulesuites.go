// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// RuleSuite represents a GitHub rule suite object returned by the rule-suites API.
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

// RuleEvaluation represents the result of evaluating a single rule within a suite.
type RuleEvaluation struct {
	RuleSource  *RuleEvaluationSource `json:"rule_source,omitempty"`
	Enforcement *string               `json:"enforcement,omitempty"`
	Result      *string               `json:"result,omitempty"`
	RuleType    *string               `json:"rule_type,omitempty"`
	Details     *string               `json:"details,omitempty"`
}

// RuleEvaluationSource identifies where a rule came from (ruleset, protected_branch, etc.).
type RuleEvaluationSource struct {
	Type *string `json:"type,omitempty"`
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ListOrganizationRuleSuites lists rule suites for the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rule-suites?apiVersion=2022-11-28#list-organization-rule-suites
//
//meta:operation GET /orgs/{org}/rulesets/rule-suites
func (s *OrganizationsService) ListOrganizationRuleSuites(ctx context.Context, org string, opts *ListOptions) ([]*RuleSuite, *Response, error) {
	endpoint := fmt.Sprintf("orgs/%v/rulesets/rule-suites", org)

	u, err := addOptions(endpoint, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var suites []*RuleSuite
	resp, err := s.client.Do(req, &suites)
	if err != nil {
		return nil, resp, err
	}

	return suites, resp, nil
}

// GetOrganizationRuleSuite gets a single rule suite for the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rule-suites?apiVersion=2022-11-28#get-an-organization-rule-suite
//
//meta:operation GET /orgs/{org}/rulesets/rule-suites/{rule_suite_id}
func (s *OrganizationsService) GetOrganizationRuleSuite(ctx context.Context, org string, ruleSuiteID int64) (*RuleSuite, *Response, error) {
	endpoint := fmt.Sprintf("orgs/%v/rulesets/rule-suites/%v", org, ruleSuiteID)

	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var suite *RuleSuite
	resp, err := s.client.Do(req, &suite)
	if err != nil {
		return nil, resp, err
	}

	return suite, resp, nil
}
