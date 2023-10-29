// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CustomProperty represents the organization custom property object.
type CustomProperty struct {
	PropertyName *string `json:"property_name,omitempty"`
	// Possible values for ValueType are: string, single_select
	ValueType     string   `json:"value_type"`
	Required      *bool    `json:"required,omitempty"`
	DefaultValue  *string  `json:"default_value,omitempty"`
	Description   *string  `json:"description,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
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
