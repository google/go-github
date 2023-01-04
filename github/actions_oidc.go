// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// OIDCSubjectClaimCustomizationTemplate represents an OIDC subject claim customization template.
type OIDCSubjectClaimCustomizationTemplate struct {
	UseDefault       *bool    `json:"use_default,omitempty"`
	IncludeClaimKeys []string `json:"include_claim_keys"`
}

// GetOrgOIDCSubjectClaimCustomizationTemplate gets the subject claim customization template for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#get-the-customization-template-for-an-oidc-subject-claim-for-an-organization
func (s *ActionsService) GetOrgOIDCSubjectClaimCustomizationTemplate(ctx context.Context, org string) (*OIDCSubjectClaimCustomizationTemplate, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/oidc/customization/sub", org)
	return s.getOIDCSubjectClaimCustomizationTemplate(ctx, u)
}

// GetRepoOIDCSubjectClaimCustomizationTemplate gets the subject claim customization template for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#get-the-customization-template-for-an-oidc-subject-claim-for-a-repository
func (s *ActionsService) GetRepoOIDCSubjectClaimCustomizationTemplate(ctx context.Context, owner, repo string) (*OIDCSubjectClaimCustomizationTemplate, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/oidc/customization/sub", owner, repo)
	return s.getOIDCSubjectClaimCustomizationTemplate(ctx, u)
}

func (s *ActionsService) getOIDCSubjectClaimCustomizationTemplate(ctx context.Context, url string) (*OIDCSubjectClaimCustomizationTemplate, *Response, error) {
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	tmpl := new(OIDCSubjectClaimCustomizationTemplate)
	resp, err := s.client.Do(ctx, req, tmpl)
	if err != nil {
		return nil, resp, err
	}

	return tmpl, resp, nil
}

// SetOrgOIDCSubjectClaimCustomizationTemplate sets the subject claim customization for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#set-the-customization-template-for-an-oidc-subject-claim-for-an-organization
func (s *ActionsService) SetOrgOIDCSubjectClaimCustomizationTemplate(ctx context.Context, org string, template *OIDCSubjectClaimCustomizationTemplate) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/oidc/customization/sub", org)
	return s.setOIDCSubjectClaimCustomizationTemplate(ctx, u, template)
}

// SetRepoOIDCSubjectClaimCustomizationTemplate sets the subject claim customization for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#set-the-customization-template-for-an-oidc-subject-claim-for-a-repository
func (s *ActionsService) SetRepoOIDCSubjectClaimCustomizationTemplate(ctx context.Context, owner, repo string, template *OIDCSubjectClaimCustomizationTemplate) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/oidc/customization/sub", owner, repo)
	return s.setOIDCSubjectClaimCustomizationTemplate(ctx, u, template)
}

func (s *ActionsService) setOIDCSubjectClaimCustomizationTemplate(ctx context.Context, url string, template *OIDCSubjectClaimCustomizationTemplate) (*Response, error) {
	req, err := s.client.NewRequest("PUT", url, template)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
