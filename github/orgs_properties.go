// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CustomProperty represents an organization custom property object.
type CustomProperty struct {
	PropertyName *string `json:"property_name,omitempty"`
	// Possible values for ValueType are: string, single_select
	ValueType     string   `json:"value_type"`
	Required      *bool    `json:"required,omitempty"`
	DefaultValue  *string  `json:"default_value,omitempty"`
	Description   *string  `json:"description,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
}

// RepoCustomPropertyValue represents a repository custom property value.
type RepoCustomPropertyValue struct {
	RepositoryID       int64                  `json:"repository_id"`
	RepositoryName     string                 `json:"repository_name"`
	RepositoryFullName string                 `json:"repository_full_name"`
	Properties         []*CustomPropertyValue `json:"properties"`
}

// CustomPropertyValue represents a custom property value.
type CustomPropertyValue struct {
	PropertyName string `json:"property_name"`
	Value        string `json:"value"`
}

// CreateOrUpdateCustomProperty creates a new or updates an existing custom property that is defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#create-or-update-a-custom-property-for-an-organization
func (s *OrganizationsService) CreateOrUpdateCustomProperty(ctx context.Context, org, name string, property *CustomProperty) (*CustomProperty, *Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/schema/%v", org, name)

	req, err := s.client.NewRequest("PUT", u, property)
	if err != nil {
		return nil, nil, err
	}

	var customProperty *CustomProperty
	resp, err := s.client.Do(ctx, req, &customProperty)
	if err != nil {
		return nil, resp, err
	}

	return customProperty, resp, nil
}

// GetCustomProperty gets a custom property that is defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#get-a-custom-property-for-an-organization
func (s *OrganizationsService) GetCustomProperty(ctx context.Context, org, name string) (*CustomProperty, *Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/schema/%v", org, name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var customProperty *CustomProperty
	resp, err := s.client.Do(ctx, req, &customProperty)
	if err != nil {
		return nil, resp, err
	}

	return customProperty, resp, nil
}

// RemoveCustomProperty removes a custom property that is defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#remove-a-custom-property-for-an-organization
func (s *OrganizationsService) RemoveCustomProperty(ctx context.Context, org, name string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/schema/%v", org, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetAllCustomProperties gets all custom properties that is defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#get-all-custom-properties-for-an-organization
func (s *OrganizationsService) GetAllCustomProperties(ctx context.Context, org string) ([]*CustomProperty, *Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/schema", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var customProperties []*CustomProperty
	resp, err := s.client.Do(ctx, req, &customProperties)
	if err != nil {
		return nil, resp, err
	}

	return customProperties, resp, nil
}

// CreateOrUpdateCustomProperties creates new or updates existing custom properties that are defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#create-or-update-custom-properties-for-an-organization
func (s *OrganizationsService) CreateOrUpdateCustomProperties(ctx context.Context, org string, properties []*CustomProperty) ([]*CustomProperty, *Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/schema", org)

	params := struct {
		Properties []*CustomProperty `json:"properties"`
	}{
		Properties: properties,
	}

	req, err := s.client.NewRequest("PATCH", u, params)
	if err != nil {
		return nil, nil, err
	}

	var customProperties []*CustomProperty
	resp, err := s.client.Do(ctx, req, &customProperties)
	if err != nil {
		return nil, resp, err
	}

	return customProperties, resp, nil
}

// CreateOrUpdateCustomPropertyValuesForRepos creates new or updates existing custom property values across multiple repositories for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#create-or-update-custom-property-values-for-organization-repositories
func (s *OrganizationsService) CreateOrUpdateCustomPropertyValuesForRepos(ctx context.Context, org string, repoNames []string, properties []*CustomProperty) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/values", org)

	params := struct {
		RepositoryNames []string          `json:"repository_names"`
		Properties      []*CustomProperty `json:"properties"`
	}{
		RepositoryNames: repoNames,
		Properties:      properties,
	}

	req, err := s.client.NewRequest("PATCH", u, params)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListCustomPropertyValues creates new or updates existing custom properties that are defined for the specified organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/properties#list-custom-property-values-for-organization-repositories
func (s *OrganizationsService) ListCustomPropertyValues(ctx context.Context, org string, opts *ListOptions) ([]*RepoCustomPropertyValue, *Response, error) {
	u := fmt.Sprintf("orgs/%v/properties/values", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repoCustomPropertyValues []*RepoCustomPropertyValue
	resp, err := s.client.Do(ctx, req, &repoCustomPropertyValues)
	if err != nil {
		return nil, resp, err
	}

	return repoCustomPropertyValues, resp, nil
}
