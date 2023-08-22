// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// OIDCSubjectClaimCustomTemplate represents an OIDC subject claim customization template.
type OIDCSubjectClaimCustomTemplate struct {
	UseDefault       *bool    `json:"use_default,omitempty"`
	IncludeClaimKeys []string `json:"include_claim_keys,omitempty"`
}

// GetOrgOIDCSubjectClaimCustomTemplate gets the subject claim customization template for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#get-the-customization-template-for-an-oidc-subject-claim-for-an-organization
func (s *ActionsService) GetOrgOIDCSubjectClaimCustomTemplate(ctx context.Context, org string) (*OIDCSubjectClaimCustomTemplate, *Response, error) {
	u, err := newURLString("orgs/%v/actions/oidc/customization/sub", org)
	if err != nil {
		return nil, nil, err
	}
	return s.getOIDCSubjectClaimCustomTemplate(ctx, u)
}

// GetRepoOIDCSubjectClaimCustomTemplate gets the subject claim customization template for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#get-the-customization-template-for-an-oidc-subject-claim-for-a-repository
func (s *ActionsService) GetRepoOIDCSubjectClaimCustomTemplate(ctx context.Context, owner, repo string) (*OIDCSubjectClaimCustomTemplate, *Response, error) {
	u, err := newURLString("repos/%v/%v/actions/oidc/customization/sub", owner, repo)
	if err != nil {
		return nil, nil, err
	}
	return s.getOIDCSubjectClaimCustomTemplate(ctx, u)
}

func (s *ActionsService) getOIDCSubjectClaimCustomTemplate(ctx context.Context, url string) (*OIDCSubjectClaimCustomTemplate, *Response, error) {
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	tmpl := new(OIDCSubjectClaimCustomTemplate)
	resp, err := s.client.Do(ctx, req, tmpl)
	if err != nil {
		return nil, resp, err
	}

	return tmpl, resp, nil
}

// SetOrgOIDCSubjectClaimCustomTemplate sets the subject claim customization for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#set-the-customization-template-for-an-oidc-subject-claim-for-an-organization
func (s *ActionsService) SetOrgOIDCSubjectClaimCustomTemplate(ctx context.Context, org string, template *OIDCSubjectClaimCustomTemplate) (*Response, error) {
	u, err := newURLString("orgs/%v/actions/oidc/customization/sub", org)
	if err != nil {
		return nil, err
	}
	return s.setOIDCSubjectClaimCustomTemplate(ctx, u, template)
}

// SetRepoOIDCSubjectClaimCustomTemplate sets the subject claim customization for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/oidc#set-the-customization-template-for-an-oidc-subject-claim-for-a-repository
func (s *ActionsService) SetRepoOIDCSubjectClaimCustomTemplate(ctx context.Context, owner, repo string, template *OIDCSubjectClaimCustomTemplate) (*Response, error) {
	u, err := newURLString("repos/%v/%v/actions/oidc/customization/sub", owner, repo)
	if err != nil {
		return nil, err
	}
	return s.setOIDCSubjectClaimCustomTemplate(ctx, u, template)
}

func (s *ActionsService) setOIDCSubjectClaimCustomTemplate(ctx context.Context, url string, template *OIDCSubjectClaimCustomTemplate) (*Response, error) {
	req, err := s.client.NewRequest("PUT", url, template)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
