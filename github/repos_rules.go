// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// BypassActor represents the bypass actors from a repository ruleset.
type BypassActor struct {
	ActorID int64 `json:"actor_id,omitempty"`
	// Possible values for ActorType are: Team, Integration
	ActorType string `json:"actor_type,omitempty"`
}

// RulesetLink represents a single link object from GitHub ruleset request _links.
type RulesetLink struct {
	HRef *string `json:"href,omitempty"`
}

// RulesetLinks represents the "_links" object in a Ruleset.
type RulesetLinks struct {
	Self *RulesetLink `json:"self,omitempty"`
}

// RulesetRefConditionParameters represents the conditions object for ref_names.
type RulesetRefConditionParameters struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// RulesetRepositoryConditionParameters represents the conditions object for repository_names.
type RulesetRepositoryConditionParameters struct {
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	Protected *bool    `json:"protected,omitempty"`
}

// RulesetCondition represents the conditions object in a ruleset.
type RulesetConditions struct {
	RefName        *RulesetRefConditionParameters        `json:"ref_name,omitempty"`
	RepositoryName *RulesetRepositoryConditionParameters `json:"repository_name,omitempty"`
}

// RulePatternParameters represents the rule pattern parameter.
type RulePatternParameters struct {
	Name *string `json:"name,omitempty"`
	// If Negate is true, the rule will fail if the pattern matches.
	Negate *bool `json:"negate,omitempty"`
	// Possible values for Operator are: starts_with, ends_with, contains, regex
	Operator string `json:"operator"`
	Pattern  string `json:"pattern"`
}

// UpdateAllowsFetchAndMergeRuleParameters represents the update rule parameters.
type UpdateAllowsFetchAndMergeRuleParameters struct {
	UpdateAllowsFetchAndMerge bool `json:"update_allows_fetch_and_merge"`
}

// RequiredDeploymentEnvironmentsRuleParameters represents the required_deployments rule parameters.
type RequiredDeploymentEnvironmentsRuleParameters struct {
	RequiredDeploymentEnvironments []string `json:"required_deployment_environments"`
}

// PullRequestRuleParameters represents the pull_request rule parameters.
type PullRequestRuleParameters struct {
	DismissStaleReviewsOnPush      bool `json:"dismiss_stale_reviews_on_push"`
	RequireCodeOwnerReview         bool `json:"require_code_owner_review"`
	RequireLastPushApproval        bool `json:"require_last_push_approval"`
	RequiredApprovingReviewCount   int  `json:"required_approving_review_count"`
	RequiredReviewThreadResolution bool `json:"required_review_thread_resolution"`
}

// RuleRequiredStatusChecks represents the RequiredStatusChecks for the RequiredStatusChecksRuleParameters object.
type RuleRequiredStatusChecks struct {
	Context       string `json:"context"`
	IntegrationID *int64 `json:"integration_id,omitempty"`
}

// RequiredStatusChecksRuleParameters represents the required_status_checks rule parameters.
type RequiredStatusChecksRuleParameters struct {
	RequiredStatusChecks             []RuleRequiredStatusChecks `json:"required_status_checks"`
	StrictRequiredStatusChecksPolicy bool                       `json:"strict_required_status_checks_policy"`
}

// RepositoryRule represents a GitHub Rule within a Ruleset.
type RepositoryRule struct {
	Type       string      `json:"type"`
	Parameters interface{} `json:"parameters,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This helps us handle the fact that RepositoryRule parameter field can be of numerous types.
func (r *RepositoryRule) UnmarshalJSON(data []byte) error {
	type rule RepositoryRule
	var RepositoryRule rule
	if err := json.Unmarshal(data, &RepositoryRule); err != nil {
		return err
	}

	r.Type = RepositoryRule.Type

	switch RepositoryRule.Type {
	case "creation", "deletion", "required_linear_history", "required_signatures", "non_fast_forward":
		r.Parameters = nil
	case "update":
		RepositoryRule.Parameters = &UpdateAllowsFetchAndMergeRuleParameters{}
		if err := json.Unmarshal(data, &RepositoryRule); err != nil {
			return err
		}
		r.Parameters = RepositoryRule.Parameters
	case "required_deployments":
		RepositoryRule.Parameters = &RequiredDeploymentEnvironmentsRuleParameters{}
		if err := json.Unmarshal(data, &RepositoryRule); err != nil {
			return err
		}
		r.Parameters = RepositoryRule.Parameters
	case "commit_message_pattern", "commit_author_email_pattern", "committer_email_pattern", "branch_name_pattern", "tag_name_pattern":
		RepositoryRule.Parameters = &RulePatternParameters{}
		if err := json.Unmarshal(data, &RepositoryRule); err != nil {
			return err
		}
		r.Parameters = RepositoryRule.Parameters
	case "pull_request":
		RepositoryRule.Parameters = &PullRequestRuleParameters{}
		if err := json.Unmarshal(data, &RepositoryRule); err != nil {
			return err
		}
		r.Parameters = RepositoryRule.Parameters
	case "required_status_checks":
		RepositoryRule.Parameters = &RequiredStatusChecksRuleParameters{}
		if err := json.Unmarshal(data, &RepositoryRule); err != nil {
			return err
		}
		r.Parameters = RepositoryRule.Parameters
	default:
		r.Type = ""
		r.Parameters = nil
		return fmt.Errorf("RepositoryRule.Type %T is not yet implemented, unable to unmarshal", RepositoryRule.Type)
	}

	return nil
}

// NewCreationRule creates a rule as part of a GitHub ruleset to only allow users with bypass permission to create matching refs.
func NewCreationRule() (rule RepositoryRule) {
	return RepositoryRule{
		Type: "creation",
	}
}

// NewUpdateRule creates a rule as part of a GitHub ruleset to only allow users with bypass permission to update matching refs.
func NewUpdateRule(params *UpdateAllowsFetchAndMergeRuleParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "update",
		Parameters: params,
	}
}

// NewDeletionRule creates a rule as part of a GitHub ruleset to only allow users with bypass permissions to delete matching refs.
func NewDeletionRule() (rule RepositoryRule) {
	return RepositoryRule{
		Type: "deletion",
	}
}

// NewRequiredLinearHistoryRule creates a rule as part of a GitHub ruleset to prevent merge commits from being pushed to matching branches.
func NewRequiredLinearHistoryRule() (rule RepositoryRule) {
	return RepositoryRule{
		Type: "required_linear_history",
	}
}

// NewRequiredDeploymentsRule creates a rule as part of a GitHub ruleset to require environments to be successfully deployed before they can be merged into the matching branches.
func NewRequiredDeploymentsRule(params *RequiredDeploymentEnvironmentsRuleParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "required_deployments",
		Parameters: params,
	}
}

// NewRequiredSignaturesRule creates a rule as part of a GitHub ruleset to require commits pushed to matching branches to have verified signatures.
func NewRequiredSignaturesRule() (rule RepositoryRule) {
	return RepositoryRule{
		Type: "required_signatures",
	}
}

// NewPullRequestRule creates a rule as part of a GitHub ruleset to require all commits be made to a non-target branch and submitted via a pull request before they can be merged.
func NewPullRequestRule(params *PullRequestRuleParameters) (
	rule RepositoryRule) {
	return RepositoryRule{
		Type:       "pull_request",
		Parameters: params,
	}
}

// NewRequiredStatusChecksRule creates a rule as part of a GitHub ruleset to require which status checks must pass before branches can be merged into a branch rule.
func NewRequiredStatusChecksRule(params *RequiredStatusChecksRuleParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "required_status_checks",
		Parameters: params,
	}
}

// NewNonFastForwardRule creates a rule as part of a GitHub ruleset to prevent users with push access from force pushing to matching branches.
func NewNonFastForwardRule() (rule RepositoryRule) {
	return RepositoryRule{
		Type: "non_fast_forward",
	}
}

// NewCommitMessagePatternRule creates a rule as part of a GitHub ruleset to restrict commit message patterns being pushed to matching branches.
func NewCommitMessagePatternRule(pattern *RulePatternParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "commit_message_pattern",
		Parameters: pattern,
	}
}

// NewCommitAuthorEmailPatternRule creates a rule as part of a GitHub ruleset to restrict commits with author email patterns being merged into matching branches.
func NewCommitAuthorEmailPatternRule(pattern *RulePatternParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "commit_author_email_pattern",
		Parameters: pattern,
	}
}

// NewCommitterEmailPatternRule creates a rule as part of a GitHub ruleset to restrict commits with committer email patterns being merged into matching branches.
func NewCommitterEmailPatternRule(pattern *RulePatternParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "committer_email_pattern",
		Parameters: pattern,
	}
}

func NewBranchNamePatternRule(pattern *RulePatternParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "branch_name_pattern",
		Parameters: pattern,
	}
}

func NewTagNamePatternRule(pattern *RulePatternParameters) (rule RepositoryRule) {
	return RepositoryRule{
		Type:       "tag_name_pattern",
		Parameters: pattern,
	}
}

// Ruleset represents a GitHub ruleset object.
type Ruleset struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// Possible values for Target are branch, tag
	Target *string `json:"target,omitempty"`
	// Possible values for SourceType are: Repository, Organization
	SourceType *string `json:"source_type,omitempty"`
	Source     string  `json:"source"`
	// Possible values for Enforcement are: disabled, active, evaluate
	Enforcement string `json:"enforcement"`
	// Possible values for BypassMode are: none, repository, organization
	BypassMode   *string            `json:"bypass_mode,omitempty"`
	BypassActors *[]BypassActor     `json:"bypass_actors,omitempty"`
	NodeID       *string            `json:"node_id,omitempty"`
	Links        *RulesetLinks      `json:"_links,omitempty"`
	Conditions   *RulesetConditions `json:"conditions,omitempty"`
	Rules        *[]RepositoryRule  `json:"rules,omitempty"`
}

// GetRulesForBranch gets all the repository rules that apply to the specified branch.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#get-rules-for-a-branch
func (s *RepositoriesService) GetRulesForBranch(ctx context.Context, owner, repo, branch string) ([]*RepositoryRule, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rules/branches/%v", owner, repo, branch)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rules []*RepositoryRule
	resp, err := s.client.Do(ctx, req, &rules)
	if err != nil {
		return nil, resp, err
	}

	return rules, resp, nil
}

// GetAllRepositoryRulesets gets all the repository rules that apply to the specified repository.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#get-all-repository-rulesets
func (s *RepositoriesService) GetAllRulesets(ctx context.Context, owner, repo string, includesParent bool) ([]*Ruleset, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rulesets?includes_parents=%v", owner, repo, includesParent)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ruleset []*Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return ruleset, resp, nil
}

// CreateRuleset creates a ruleset for the specified repository.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#create-a-repository-ruleset
func (s *RepositoriesService) CreateRuleset(ctx context.Context, owner, repo string, rs *Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rulesets", owner, repo)

	req, err := s.client.NewRequest("POST", u, rs)
	if err != nil {
		return nil, nil, err
	}

	var ruleset *Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return ruleset, resp, nil
}

// GetRuleset gets a ruleset for the specified repository.
// If includesParent is true, rulesets configured at organisation level that apply to the repository can be retrieved.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#get-a-repository-ruleset
func (s *RepositoriesService) GetRuleset(ctx context.Context, owner, repo string, rulesetID int64, includesParent bool) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rulesets/%v?includes_parents=%v", owner, repo, rulesetID, includesParent)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ruleset *Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return ruleset, resp, nil
}

// UpdateRuleset updates a ruleset for the specified repository.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#update-a-repository-ruleset
func (s *RepositoriesService) UpdateRuleset(ctx context.Context, owner, repo string, rulesetID int64, rs *Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rulesets/%v", owner, repo, rulesetID)

	req, err := s.client.NewRequest("PUT", u, rs)
	if err != nil {
		return nil, nil, err
	}

	var ruleset *Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return ruleset, resp, nil
}

// DeleteRuleset deletes a ruleset for the specified repository.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/rules#delete-a-repository-ruleset
func (s *RepositoriesService) DeleteRuleset(ctx context.Context, owner, repo string, rulesetID int64) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/rulesets/%v", owner, repo, rulesetID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
