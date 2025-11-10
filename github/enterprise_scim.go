// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// The URIs that are used to indicate the namespaces of the SCIM schemas (only core schemas are supported).
const SCIMSchemasURINamespacesGroups string = "urn:ietf:params:scim:schemas:core:2.0:Group"

// SCIMEnterpriseGroupAttributes represents supported SCIM Enterprise group attributes.
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/scim#supported-scim-group-attributes
type SCIMEnterpriseGroupAttributes struct {
	DisplayName *string                           `json:"displayName,omitempty"` // Human-readable name for a group.
	Members     []*SCIMEnterpriseDisplayReference `json:"members,omitempty"`     // (Optional.)
	ExternalID  *string                           `json:"externalId,omitempty"`  // (Optional.)
	// Only populated as a result of calling ListSCIMProvisionedIdentitiesOptions:
	Schemas []string  `json:"schemas"` // (Optional.)
	ID      *string   `json:"id,omitempty"`
	Meta    *SCIMMeta `json:"meta,omitempty"`
}

// SCIMEnterpriseDisplayReference represents a JSON SCIM (System for Cross-domain Identity Management) resource.
type SCIMEnterpriseDisplayReference struct {
	Value   string  `json:"value"`                 // (Required.)
	Ref     string  `json:"$+ref"`                 // (Required.)
	Display *string `json:"displayName,omitempty"` // (Optional.)
}

// SCIMEnterpriseMeta represents metadata about the SCIM resource.
type SCIMEnterpriseMeta struct {
	ResourceType *string    `json:"resourceType,omitempty"`
	Created      *Timestamp `json:"created,omitempty"`
	LastModified *Timestamp `json:"lastModified,omitempty"`
	Location     *string    `json:"location,omitempty"`
}

// ListProvisionedSCIMGroupsForEnterprise lists provisioned SCIM groups in an enterprise.
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/scim#list-provisioned-scim-groups-for-an-enterprise
//
//meta:operation GET /scim/v2/enterprises/{enterprise}/Groups
// func (s *SCIMService) ListProvisionedSCIMGroupsForEnterprise(ctx context.Context, enterprise string, opts *ListOptions) ([]*SCIMEnterpriseGroupAttributes, *Response, error) {
//}
