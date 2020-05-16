// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// PublicKey represents the public key that should be used to encrypt secrets.
type PublicKey struct {
	KeyID *string `json:"key_id"`
	Key   *string `json:"key"`
}

// GetRepositoryPublicKey gets a public key that should be used for secret encryption.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-a-repository-public-key
func (s *ActionsService) GetRepositoryPublicKey(ctx context.Context, owner, repo string) (*PublicKey, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/public-key", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pubKey := new(PublicKey)
	resp, err := s.client.Do(ctx, req, pubKey)
	if err != nil {
		return nil, resp, err
	}

	return pubKey, resp, nil
}

// GetOrganizationPublicKey gets a public key that should be used for secret encryption.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-an-organization-public-key
func (s *ActionsService) GetOrganizationPublicKey(ctx context.Context, org string) (*PublicKey, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/public-key", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pubKey := new(PublicKey)
	resp, err := s.client.Do(ctx, req, pubKey)
	if err != nil {
		return nil, resp, err
	}

	return pubKey, resp, nil
}

// Secret represents a repository action secret.
type Secret struct {
	Name                    string    `json:"name"`
	CreatedAt               Timestamp `json:"created_at"`
	UpdatedAt               Timestamp `json:"updated_at"`
	Visibility              string    `json:"visibility,omitempty"`
	SelectedRepositoriesURL string    `json:"selected_repositories_url,omitempty"`
}

// Secrets represents one item from the ListSecrets response.
type Secrets struct {
	TotalCount int       `json:"total_count"`
	Secrets    []*Secret `json:"secrets"`
}

// ListRepositorySecrets lists all secrets available in a repository
// without revealing their encrypted values.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-repository-secrets
func (s *ActionsService) ListRepositorySecrets(ctx context.Context, owner, repo string, opts *ListOptions) (*Secrets, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secrets := new(Secrets)
	resp, err := s.client.Do(ctx, req, &secrets)
	if err != nil {
		return nil, resp, err
	}

	return secrets, resp, nil
}

// GetRepositorySecret gets a single repository secret without revealing its encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-a-repository-secret
func (s *ActionsService) GetRepositorySecret(ctx context.Context, owner, repo, name string) (*Secret, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/%v", owner, repo, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secret := new(Secret)
	resp, err := s.client.Do(ctx, req, secret)
	if err != nil {
		return nil, resp, err
	}

	return secret, resp, nil
}

// SelectedRepositoryIDs are the repository IDs that have access to the secret.
type SelectedRepositoryIDs struct {
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitempty"`
}

// EncryptedSecret represents a secret that is encrypted using a public key.
//
// The value of EncryptedValue must be your secret, encrypted with
// LibSodium (see documentation here: https://libsodium.gitbook.io/doc/bindings_for_other_languages)
// using the public key retrieved using the GetPublicKey method.
type EncryptedSecret struct {
	Name           string `json:"-"`
	KeyID          string `json:"key_id"`
	EncryptedValue string `json:"encrypted_value"`
	Visibility     string `json:"visibility,omitempty"`
	SelectedRepositoryIDs
}

// CreateOrUpdateRepositorySecret creates or updates a repository secret with an encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#create-or-update-a-repository-secret
func (s *ActionsService) CreateOrUpdateRepositorySecret(ctx context.Context, owner, repo string, eSecret *EncryptedSecret) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/%v", owner, repo, eSecret.Name)

	req, err := s.client.NewRequest("PUT", u, eSecret)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteRepositorySecret deletes a secret in a repository using the secret name.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#delete-a-repository-secret
func (s *ActionsService) DeleteRepositorySecret(ctx context.Context, owner, repo, name string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/%v", owner, repo, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListOrganizationSecrets lists all secrets available in an organization
// without revealing their encrypted values.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-organization-secrets
func (s *ActionsService) ListOrganizationSecrets(ctx context.Context, org string, opts *ListOptions) (*Secrets, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secrets := new(Secrets)
	resp, err := s.client.Do(ctx, req, &secrets)
	if err != nil {
		return nil, resp, err
	}

	return secrets, resp, nil
}

// GetOrganizationSecret gets a single organization secret without revealing its encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-an-organization-secret
func (s *ActionsService) GetOrganizationSecret(ctx context.Context, org, name string) (*Secret, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", org, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secret := new(Secret)
	resp, err := s.client.Do(ctx, req, secret)
	if err != nil {
		return nil, resp, err
	}

	return secret, resp, nil
}

// CreateOrUpdateOrganizationSecret creates or updates an organization secret with an encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#create-or-update-an-organization-secret
func (s *ActionsService) CreateOrUpdateOrganizationSecret(ctx context.Context, org string, eSecret *EncryptedSecret) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", org, eSecret.Name)

	req, err := s.client.NewRequest("PUT", u, eSecret)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// SelectedRepositoriesList represents the list of repositories selected for an organization secret.
type SelectedRepositoriesList struct {
	TotalCount   *int          `json:"total_count,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
}

// ListSelectedRepositoriesForOrganizationSecret lists all repositories that have access to a secret.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-selected-repositories-for-an-organization-secret
func (s *ActionsService) ListSelectedRepositoriesForOrganizationSecret(ctx context.Context, org, name string) (*SelectedRepositoriesList, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories", org, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SelectedRepositoriesList)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetSelectedRepositoriesForOrganizationSecret sets the repositories that have access to a secret.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#set-selected-repositories-for-an-organization-secret
func (s *ActionsService) SetSelectedRepositoriesForOrganizationSecret(ctx context.Context, org, name string, ids SelectedRepositoryIDs) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories", org, name)
	req, err := s.client.NewRequest("PUT", u, ids)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddSelectedRepositoryToOrganizationSecret adds a repository to an organization secret.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#add-selected-repository-to-an-organization-secret
func (s *ActionsService) AddSelectedRepositoryToOrganizationSecret(ctx context.Context, org, name string, repo *Repository) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories/%v", org, name, *repo.ID)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RemoveSelectedRepositoryFromOrganizationSecret removes a repository from an organization secret.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#remove-selected-repository-from-an-organization-secret
func (s *ActionsService) RemoveSelectedRepositoryFromOrganizationSecret(ctx context.Context, org, name string, repo *Repository) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories/%v", org, name, *repo.ID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteOrganizationSecret deletes a secret in an organization using the secret name.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#delete-an-organization-secret
func (s *ActionsService) DeleteOrganizationSecret(ctx context.Context, org, name string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", org, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
