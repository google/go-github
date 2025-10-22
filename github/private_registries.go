// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// PrivateRegistriesService handles communication with the private registries
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/private-registries
type PrivateRegistriesService service

// PrivateRegistryType represents the type of private registry.
type PrivateRegistryType string

const (
	PrivateRegistryTypeMavenRepository    PrivateRegistryType = "maven_repository"
	PrivateRegistryTypeNugetFeed          PrivateRegistryType = "nuget_feed"
	PrivateRegistryTypeGoProxyServer      PrivateRegistryType = "goproxy_server"
	PrivateRegistryTypeNpmRegistry        PrivateRegistryType = "npm_registry"
	PrivateRegistryTypeRubygemsServer     PrivateRegistryType = "rubygems_server"
	PrivateRegistryTypeCargoRegistry      PrivateRegistryType = "cargo_registry"
	PrivateRegistryTypeComposerRepository PrivateRegistryType = "composer_repository"
	PrivateRegistryTypeDockerRegistry     PrivateRegistryType = "docker_registry"
	PrivateRegistryTypeGitSource          PrivateRegistryType = "git_source"
	PrivateRegistryTypeHelmRegistry       PrivateRegistryType = "helm_registry"
	PrivateRegistryTypeHexOrganization    PrivateRegistryType = "hex_organization"
	PrivateRegistryTypeHexRepository      PrivateRegistryType = "hex_repository"
	PrivateRegistryTypePubRepository      PrivateRegistryType = "pub_repository"
	PrivateRegistryTypePythonIndex        PrivateRegistryType = "python_index"
	PrivateRegistryTypeTerraformRegistry  PrivateRegistryType = "terraform_registry"
)

// PrivateRegistryVisibility represents the visibility of a private registry.
type PrivateRegistryVisibility string

const (
	PrivateRegistryVisibilityPrivate  PrivateRegistryVisibility = "private"
	PrivateRegistryVisibilityAll      PrivateRegistryVisibility = "all"
	PrivateRegistryVisibilitySelected PrivateRegistryVisibility = "selected"
)

// PrivateRegistry represents a private registry configuration.
type PrivateRegistry struct {
	Name         *string    `json:"name,omitempty"`
	RegistryType *string    `json:"registry_type,omitempty"`
	Username     *string    `json:"username,omitempty"`
	CreatedAt    *Timestamp `json:"created_at,omitempty"`
	UpdatedAt    *Timestamp `json:"updated_at,omitempty"`
	Visibility   *string    `json:"visibility,omitempty"`
}

// PrivateRegistries represents a list of private registries.
type PrivateRegistries struct {
	TotalCount     *int               `json:"total_count,omitempty"`
	Configurations []*PrivateRegistry `json:"configurations,omitempty"`
}

// CreateOrganizationPrivateRegistry represents the payload to create a private registry.
type CreateOrganizationPrivateRegistry struct {
	RegistryType string `json:"registry_type"`
	URL          string `json:"url"`

	// The username to use when authenticating with the private registry.
	// This field should be omitted if the private registry does not require a username for authentication.
	Username *string `json:"username,omitempty"`

	// The value for your secret, encrypted with [LibSodium](https://libsodium.gitbook.io/doc/bindings_for_other_languages)
	// using the public key retrieved from the PrivateRegistriesService.GetOrganizationPrivateRegistriesPublicKey.
	EncryptedValue string                    `json:"encrypted_value"`
	KeyID          string                    `json:"key_id"`
	Visibility     PrivateRegistryVisibility `json:"visibility"`

	// An array of repository IDs that can access the organization private registry.
	// You can only provide a list of repository IDs when CreateOrganizationPrivateRegistry.Visibility is set to PrivateRegistryVisibilitySelected.
	// This field should be omitted if visibility is set to PrivateRegistryVisibilityAll or PrivateRegistryVisibilityPrivate.
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitempty"`
}

// UpdateOrganizationPrivateRegistry represents the payload to update a private registry.
type UpdateOrganizationPrivateRegistry struct {
	RegistryType *string `json:"registry_type,omitempty"`
	URL          *string `json:"url,omitempty"`

	// The username to use when authenticating with the private registry.
	// This field should be omitted if the private registry does not require a username for authentication.
	Username *string `json:"username,omitempty"`

	// The value for your secret, encrypted with [LibSodium](https://libsodium.gitbook.io/doc/bindings_for_other_languages)
	// using the public key retrieved from the PrivateRegistriesService.GetOrganizationPrivateRegistriesPublicKey.
	EncryptedValue *string                    `json:"encrypted_value,omitempty"`
	KeyID          *string                    `json:"key_id,omitempty"`
	Visibility     *PrivateRegistryVisibility `json:"visibility,omitempty"`

	// An array of repository IDs that can access the organization private registry.
	// You can only provide a list of repository IDs when CreateOrganizationPrivateRegistry.Visibility is set to PrivateRegistryVisibilitySelected.
	// This field should be omitted if visibility is set to PrivateRegistryVisibilityAll or PrivateRegistryVisibilityPrivate.
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitempty"`
}

// ListOrganizationPrivateRegistries lists private registries for an organization.
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#list-private-registries-for-an-organization
//
//meta:operation GET /orgs/{org}/private-registries
func (s *PrivateRegistriesService) ListOrganizationPrivateRegistries(ctx context.Context, org string, opts *ListOptions) (*PrivateRegistries, *Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var privateRegistries PrivateRegistries
	resp, err := s.client.Do(ctx, req, &privateRegistries)
	if err != nil {
		return nil, resp, err
	}
	return &privateRegistries, resp, nil
}

// CreateOrganizationPrivateRegistries creates a private registry for an organization.
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#create-a-private-registry-for-an-organization
//
//meta:operation POST /orgs/{org}/private-registries
func (s *PrivateRegistriesService) CreateOrganizationPrivateRegistries(ctx context.Context, org string, privateRegistry CreateOrganizationPrivateRegistry) (*PrivateRegistry, *Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries", org)

	req, err := s.client.NewRequest("POST", u, privateRegistry)
	if err != nil {
		return nil, nil, err
	}

	var result PrivateRegistry
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// GetOrganizationPrivateRegistriesPublicKey retrieves the public key for encrypting secrets for an organization's private registries.
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#get-private-registries-public-key-for-an-organization
//
//meta:operation GET /orgs/{org}/private-registries/public-key
func (s *PrivateRegistriesService) GetOrganizationPrivateRegistriesPublicKey(ctx context.Context, org string) (*PublicKey, *Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries/public-key", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var publicKey PublicKey
	resp, err := s.client.Do(ctx, req, &publicKey)
	if err != nil {
		return nil, resp, err
	}
	return &publicKey, resp, nil
}

// GetOrganizationPrivateRegistry gets a specific private registry for an organization.
// name parameter is the name of the private registry to delete, PrivateRegistry.Name
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#get-a-private-registry-for-an-organization
//
//meta:operation GET /orgs/{org}/private-registries/{secret_name}
func (s *PrivateRegistriesService) GetOrganizationPrivateRegistry(ctx context.Context, org, name string) (*PrivateRegistry, *Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries/%v", org, name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var privateRegistry PrivateRegistry
	resp, err := s.client.Do(ctx, req, &privateRegistry)
	if err != nil {
		return nil, resp, err
	}

	return &privateRegistry, resp, nil
}

// UpdateOrganizationPrivateRegistry updates a specific private registry for an organization.
// name parameter is the name of the private registry to delete, PrivateRegistry.Name
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#update-a-private-registry-for-an-organization
//
//meta:operation PATCH /orgs/{org}/private-registries/{secret_name}
func (s *PrivateRegistriesService) UpdateOrganizationPrivateRegistry(ctx context.Context, org, name string, privateRegistry UpdateOrganizationPrivateRegistry) (*PrivateRegistry, *Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries/%v", org, name)

	req, err := s.client.NewRequest("PATCH", u, privateRegistry)
	if err != nil {
		return nil, nil, err
	}

	var updatedRegistry PrivateRegistry
	resp, err := s.client.Do(ctx, req, &updatedRegistry)
	if err != nil {
		return nil, resp, err
	}

	return &updatedRegistry, resp, nil
}

// DeleteOrganizationPrivateRegistry deletes a specific private registry for an organization.
// name parameter is the name of the private registry to delete, PrivateRegistry.Name
//
// GitHub API docs: https://docs.github.com/rest/private-registries/organization-configurations#delete-a-private-registry-for-an-organization
//
//meta:operation DELETE /orgs/{org}/private-registries/{secret_name}
func (s *PrivateRegistriesService) DeleteOrganizationPrivateRegistry(ctx context.Context, org, name string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/private-registries/%v", org, name)

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
