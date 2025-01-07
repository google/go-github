// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetAllOrganizationRulesets gets all the rulesets for the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#get-all-organization-repository-rulesets
//
//meta:operation GET /orgs/{org}/rulesets
func (s *OrganizationsService) GetAllOrganizationRulesets(ctx context.Context, org string) ([]*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rulesets []*Ruleset
	resp, err := s.client.Do(ctx, req, &rulesets)
	if err != nil {
		return nil, resp, err
	}

	return rulesets, resp, nil
}

// CreateOrganizationRuleset creates a ruleset for the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#create-an-organization-repository-ruleset
//
//meta:operation POST /orgs/{org}/rulesets
func (s *OrganizationsService) CreateOrganizationRuleset(ctx context.Context, org string, ruleset Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets", org)

	req, err := s.client.NewRequest("POST", u, ruleset)
	if err != nil {
		return nil, nil, err
	}

	var rs *Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, nil
}

// GetOrganizationRuleset gets a ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#get-an-organization-repository-ruleset
//
//meta:operation GET /orgs/{org}/rulesets/{ruleset_id}
func (s *OrganizationsService) GetOrganizationRuleset(ctx context.Context, org string, rulesetID int64) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

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

// UpdateOrganizationRuleset updates a ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#update-an-organization-repository-ruleset
//
//meta:operation PUT /orgs/{org}/rulesets/{ruleset_id}
func (s *OrganizationsService) UpdateOrganizationRuleset(ctx context.Context, org string, rulesetID int64, ruleset Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	req, err := s.client.NewRequest("PUT", u, ruleset)
	if err != nil {
		return nil, nil, err
	}

	var rs *Ruleset
	resp, err := s.client.Do(ctx, req, &ruleset)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, nil
}

// UpdateOrganizationRulesetClearBypassActor clears the ruleset bypass actors for a ruleset for the specified repository.
//
// This function is necessary as the UpdateOrganizationRuleset function does not marshal ByPassActor if passed as an empty array.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#update-an-organization-repository-ruleset
//
//meta:operation PUT /orgs/{org}/rulesets/{ruleset_id}
func (s *OrganizationsService) UpdateOrganizationRulesetClearBypassActor(ctx context.Context, org string, rulesetID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	rsClearBypassActor := rulesetClearBypassActors{}

	req, err := s.client.NewRequest("PUT", u, rsClearBypassActor)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteOrganizationRuleset deletes a ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/rules#delete-an-organization-repository-ruleset
//
//meta:operation DELETE /orgs/{org}/rulesets/{ruleset_id}
func (s *OrganizationsService) DeleteOrganizationRuleset(ctx context.Context, org string, rulesetID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
