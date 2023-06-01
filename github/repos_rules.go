// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
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
type RulesetCondition struct {
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

// RulesetRule represents a GitHub Rule within a Ruleset.
type RulesetRule struct {
	Type       string      `json:"type"`
	Parameters interface{} `json:"parameters,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This helps us handle the fact that RulesetRule parameter field can be of numerous types.
func (rsr *RulesetRule) UnmarshalJSON(data []byte) error {
	type rule RulesetRule
	var rulesetRule rule
	if err := json.Unmarshal(data, &rulesetRule); err != nil {
		return err
	}

	rsr.Type = rulesetRule.Type

	switch rulesetRule.Type {
	case "creation", "deletion", "required_linear_history", "required_signatures", "non_fast_forward":
		rsr.Parameters = nil
	case "update":
		rulesetRule.Parameters = &UpdateAllowsFetchAndMergeRuleParameters{}
		if err := json.Unmarshal(data, &rulesetRule); err != nil {
			return err
		}
		rsr.Parameters = rulesetRule.Parameters
	case "required_deployments":
		rulesetRule.Parameters = &RequiredDeploymentEnvironmentsRuleParameters{}
		if err := json.Unmarshal(data, &rulesetRule); err != nil {
			return err
		}
		rsr.Parameters = rulesetRule.Parameters
	case "commit_message_pattern", "commit_author_email_pattern", "committer_email_pattern", "branch_name_pattern", "tag_name_pattern":
		rulesetRule.Parameters = &RulePatternParameters{}
		if err := json.Unmarshal(data, &rulesetRule); err != nil {
			return err
		}
		rsr.Parameters = rulesetRule.Parameters
	case "pull_request":
		rulesetRule.Parameters = &PullRequestRuleParameters{}
		if err := json.Unmarshal(data, &rulesetRule); err != nil {
			return err
		}
		rsr.Parameters = rulesetRule.Parameters
	case "required_status_checks":
		rulesetRule.Parameters = &RequiredStatusChecksRuleParameters{}
		if err := json.Unmarshal(data, &rulesetRule); err != nil {
			return err
		}
		rsr.Parameters = rulesetRule.Parameters
	default:
		rsr.Type = ""
		rsr.Parameters = nil
		return fmt.Errorf("rulesetRule.Type %T is not yet implemented, unable to unmarshal", rulesetRule.Type)
	}

	return nil
}

func NewCreationRule() (rule RulesetRule) {
	return RulesetRule{
		Type: "creation",
	}
}

func NewUpdateRule(params *UpdateAllowsFetchAndMergeRuleParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "update",
		Parameters: params,
	}
}

func NewDeletionRule() (rule RulesetRule) {
	return RulesetRule{
		Type: "deletion",
	}
}

func NewRequiredLinearHistoryRule() (rule RulesetRule) {
	return RulesetRule{
		Type: "required_linear_history",
	}
}

func NewRequiredDeploymentsRule(params *RequiredDeploymentEnvironmentsRuleParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "required_deployments",
		Parameters: params,
	}
}

func NewRequiredSignaturesRule() (rule RulesetRule) {
	return RulesetRule{
		Type: "required_signatures",
	}
}

func NewPullRequestRule(params *PullRequestRuleParameters) (
	rule RulesetRule) {
	return RulesetRule{
		Type:       "pull_request",
		Parameters: params,
	}
}

func NewRequiredStatusChecksRule(params *RequiredStatusChecksRuleParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "required_status_checks",
		Parameters: params,
	}
}

func NewNonFastForwardRule() (rule RulesetRule) {
	return RulesetRule{
		Type: "non_fast_forward",
	}
}

func NewCommitMessagePatternRule(pattern *RulePatternParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "commit_message_pattern",
		Parameters: pattern,
	}
}

func NewCommitAuthorEmailPatternRule(pattern *RulePatternParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "commit_author_email_pattern",
		Parameters: pattern,
	}
}

func NewCommitterEmailPatternRule(pattern *RulePatternParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "committer_email_pattern",
		Parameters: pattern,
	}
}

func NewBranchNamePatternRule(pattern *RulePatternParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "branch_name_pattern",
		Parameters: pattern,
	}
}

func NewTagNamePatternRule(pattern *RulePatternParameters) (rule RulesetRule) {
	return RulesetRule{
		Type:       "tag_name_pattern",
		Parameters: pattern,
	}
}

// Ruleset represents a GitHub rules request.
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
	BypassMode   *string           `json:"bypass_mode,omitempty"`
	BypassActors *[]BypassActor    `json:"bypass_actors,omitempty"`
	NodeID       *string           `json:"node_id,omitempty"`
	Links        *RulesetLinks     `json:"_links,omitempty"`
	Conditions   *RulesetCondition `json:"conditions,omitempty"`
	Rules        *[]RulesetRule    `json:"rules,omitempty"`
}
