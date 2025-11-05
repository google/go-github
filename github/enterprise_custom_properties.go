// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// Property represents a custom property definition in an enterprise.
type Property struct {
	// Name of the custom property.
	Name string `json:"name"`
	// The URL that can be used to fetch, update, or delete info about this property via the API.
	URL *string `json:"url,omitempty"`
	// The source type of the property.
	SourceType *string `json:"source_type,omitempty"`
	// The type of the value for the property.
	ValueType string `json:"value_type"`
	// Whether the property is required.
	Required *bool `json:"required,omitempty"`
	// The default value for the property.
	DefaultValue *string `json:"default_value,omitempty"`
	// Short description of the property.
	Description *string `json:"description,omitempty"`
	// An ordered list of the allowed values of the property.
	// The property can have up to 200 allowed values.
	AllowedValues []string `json:"allowed_values,omitempty"`
	// Who can edit the values of the property.
	ValuesEditableBy *string `json:"values_editable_by,omitempty"`
}

// EnterpriseCustomPropertySchema represents the schema response for GetEnterpriseCustomPropertiesSchema.
type EnterpriseCustomPropertySchema struct {
	// An ordered list of the custom property defined in the enterprise.
	Properties []*Property `json:"properties,omitempty"`
}

// PropertyValue represents the custom property values for an organization.
type PropertyValue struct {
	// The name of the property
	PropertyName *string `json:"property_name,omitempty"`
	// The value assigned to the property
	Value *string `json:"value,omitempty"`
}

// CustomPropertiesValues represents the custom properties values for an organization within an enterprise.
type CustomPropertiesValues struct {
	// The Organization ID that the custom property values will be applied to.
	OrganizationID *int64 `json:"organization_id,omitempty"`
	// The names of organizations that the custom property values will be applied to.
	OrganizationLogin *string `json:"organization_login,omitempty"`
	// List of custom property names and associated values to apply to the organizations.
	Properties []*PropertyValue `json:"properties,omitempty"`
}

// CustomPropertiesValuesUpdate represents the request to update custom property values for organizations within an enterprise.
type CustomPropertiesValuesUpdate struct {
	// The names of organizations that the custom property values will be applied to.
	// OrganizationLogin specifies the organization name when updating multiple organizations.
	OrganizationLogin []string `json:"organization_login"`
	// List of custom property names and associated values to apply to the organizations.
	Properties []*PropertyValue `json:"properties"`
}

// GetEnterpriseCustomPropertySchema gives all organization custom property definitions that are defined on an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#get-organization-custom-properties-schema-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/org-properties/schema
func (s *EnterpriseService) GetEnterpriseCustomPropertySchema(ctx context.Context, enterprise string) (*EnterpriseCustomPropertySchema, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/schema", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var schema *EnterpriseCustomPropertySchema
	resp, err := s.client.Do(ctx, req, &schema)
	if err != nil {
		return nil, resp, err
	}

	return schema, resp, nil
}

// UpdateEnterpriseCustomPropertySchema creates new or updates existing organization custom properties defined on an enterprise in a batch.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#create-or-update-organization-custom-property-definitions-on-an-enterprise
//
//meta:operation PATCH /enterprises/{enterprise}/org-properties/schema
func (s *EnterpriseService) UpdateEnterpriseCustomPropertySchema(ctx context.Context, enterprise string, schema *EnterpriseCustomPropertySchema) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/schema", enterprise)
	req, err := s.client.NewRequest("PATCH", u, schema)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetEnterpriseCustomProperty retrieves a specific organization custom property definition from an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#get-an-organization-custom-property-definition-from-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/org-properties/schema/{custom_property_name}
func (s *EnterpriseService) GetEnterpriseCustomProperty(ctx context.Context, enterprise, customPropertyName string) (*Property, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/schema/%v", enterprise, customPropertyName)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var property *Property
	resp, err := s.client.Do(ctx, req, &property)
	if err != nil {
		return nil, resp, err
	}

	return property, resp, nil
}

// UpdateEnterpriseCustomProperty creates a new or updates an existing organization custom property definition that is defined on an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#create-or-update-an-organization-custom-property-definition-on-an-enterprise
//
//meta:operation PUT /enterprises/{enterprise}/org-properties/schema/{custom_property_name}
func (s *EnterpriseService) UpdateEnterpriseCustomProperty(ctx context.Context, enterprise, customPropertyName string, property *Property) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/schema/%v", enterprise, customPropertyName)
	req, err := s.client.NewRequest("PUT", u, property)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteEnterpriseCustomProperty removes an organization custom property definition that is defined on an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#remove-an-organization-custom-property-definition-from-an-enterprise
//
//meta:operation DELETE /enterprises/{enterprise}/org-properties/schema/{custom_property_name}
func (s *EnterpriseService) DeleteEnterpriseCustomProperty(ctx context.Context, enterprise, customPropertyName string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/schema/%v", enterprise, customPropertyName)
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

// GetEnterpriseCustomPropertyValues lists enterprise organizations with all of their custom property values.
// Returns a list of organizations and their custom property values defined in the enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#list-custom-property-values-for-organizations-in-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/org-properties/values
func (s *EnterpriseService) GetEnterpriseCustomPropertyValues(ctx context.Context, enterprise string, opts *ListOptions) ([]*CustomPropertiesValues, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/values", enterprise)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var values []*CustomPropertiesValues
	resp, err := s.client.Do(ctx, req, &values)
	if err != nil {
		return nil, resp, err
	}

	return values, resp, nil
}

// UpdateEnterpriseCustomPropertyValues create or update custom property values for organizations in an enterprise.
// To remove a custom property value from an organization, set the property value to null.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/custom-properties-for-orgs#create-or-update-custom-property-values-for-organizations-in-an-enterprise
//
//meta:operation PATCH /enterprises/{enterprise}/org-properties/values
func (s *EnterpriseService) UpdateEnterpriseCustomPropertyValues(ctx context.Context, enterprise string, values CustomPropertiesValuesUpdate) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/org-properties/values", enterprise)
	req, err := s.client.NewRequest("PATCH", u, values)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
