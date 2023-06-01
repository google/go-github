// Copyright 2017 The go-github AUTHORS. All rights reserved.
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
