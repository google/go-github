// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CreateIssueTypeRequest represents the parameters to create an issue type.
type CreateIssueTypeRequest struct {
	Name        string  `json:"name"`                  // Name of the issue type. (Required.)
	IsEnabled   bool    `json:"is_enabled"`            // Whether or not the issue type is enabled at the organization level. (Required.)
	Description *string `json:"description,omitempty"` // Description of the issue type. (Optional.)
	Color       *string `json:"color,omitempty"`       // Color for the issue type. Can be one of "gray", "blue", "green", "yellow", "orange", "red", "pink", "purple". (Optional.)
}

// UpdateIssueTypeRequest represents the parameters to update an issue type.
type UpdateIssueTypeRequest struct {
	Name        string  `json:"name"`                  // Name of the issue type. (Required.)
	IsEnabled   bool    `json:"is_enabled"`            // Whether or not the issue type is enabled at the organization level. (Required.)
	Description *string `json:"description,omitempty"` // Description of the issue type. (Optional.)
	Color       *string `json:"color,omitempty"`       // Color for the issue type. Can be one of "gray", "blue", "green", "yellow", "orange", "red", "pink", "purple". (Optional.)
}

// ListIssueTypes lists all issue types for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/issue-types?apiVersion=2022-11-28#list-issue-types-for-an-organization
//
//meta:operation GET /orgs/{org}/issue-types
func (s *OrganizationsService) ListIssueTypes(ctx context.Context, org string) ([]*IssueType, *Response, error) {
	u := fmt.Sprintf("orgs/%v/issue-types", org)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}

// CreateIssueType creates a new issue type for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/issue-types?apiVersion=2022-11-28#create-issue-type-for-an-organization
//
//meta:operation POST /orgs/{org}/issue-types
func (s *OrganizationsService) CreateIssueType(ctx context.Context, org string, body CreateIssueTypeRequest) (*IssueType, *Response, error) {
	u := fmt.Sprintf("orgs/%v/issue-types", org)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var issueType *IssueType
	resp, err := s.client.Do(req, &issueType)
	if err != nil {
		return nil, resp, err
	}

	return issueType, resp, nil
}

// UpdateIssueType updates an issue type for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/issue-types?apiVersion=2022-11-28#update-issue-type-for-an-organization
//
//meta:operation PUT /orgs/{org}/issue-types/{issue_type_id}
func (s *OrganizationsService) UpdateIssueType(ctx context.Context, org string, issueTypeID int64, body UpdateIssueTypeRequest) (*IssueType, *Response, error) {
	u := fmt.Sprintf("orgs/%v/issue-types/%v", org, issueTypeID)
	req, err := s.client.NewRequest(ctx, "PUT", u, body)
	if err != nil {
		return nil, nil, err
	}

	var issueType *IssueType
	resp, err := s.client.Do(req, &issueType)
	if err != nil {
		return nil, resp, err
	}

	return issueType, resp, nil
}

// DeleteIssueType deletes an issue type for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/issue-types?apiVersion=2022-11-28#delete-issue-type-for-an-organization
//
//meta:operation DELETE /orgs/{org}/issue-types/{issue_type_id}
func (s *OrganizationsService) DeleteIssueType(ctx context.Context, org string, issueTypeID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/issue-types/%v", org, issueTypeID)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
