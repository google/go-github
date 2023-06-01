// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetAllOrganizationRepositoryRulesets gets all the repository rulesets for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/rules#get-all-organization-repository-rulesets
func (s *OrganizationsService) GetAllOrganizationRepositoryRulesets(ctx context.Context, org string) ([]*Ruleset, *Response, error) {
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

// CreateOrganizationRepositoryRuleset creates a repository ruleset for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/rules#create-an-organization-repository-ruleset
func (s *OrganizationsService) CreateOrganizationRepositoryRuleset(ctx context.Context, org string, rs *Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets", org)

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

// GetOrganizationRepositoryRuleset gets a repository ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/rules#get-an-organization-repository-ruleset
func (s *OrganizationsService) GetOrganizationRepositoryRuleset(ctx context.Context, org string, rulesetID int64) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rulesets *Ruleset
	resp, err := s.client.Do(ctx, req, &rulesets)
	if err != nil {
		return nil, resp, err
	}

	return rulesets, resp, nil
}

// UpdateOrganizationRepositoryRuleset updates a repository ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/rules#update-a-repository-ruleset
func (s *OrganizationsService) UpdateOrganizationRepositoryRuleset(ctx context.Context, org string, rulesetID int64, rs *Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	req, err := s.client.NewRequest("POST", u, rs)
	if err != nil {
		return nil, nil, err
	}

	var rulesets *Ruleset
	resp, err := s.client.Do(ctx, req, &rulesets)
	if err != nil {
		return nil, resp, err
	}

	return rulesets, resp, nil
}

// DeleteOrganizationRepositoryRuleset deletes a repository ruleset from the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/rules#delete-a-repository-ruleset
func (s *OrganizationsService) DeleteOrganizationRepositoryRuleset(ctx context.Context, org string, rulesetID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/rulesets/%v", org, rulesetID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
