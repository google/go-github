package github

import (
	"context"
	"fmt"
)

// PublicKey represents the public key that should be used to encrypt secrets.
type OrganizationPublicKey struct {
	KeyID *string `json:"key_id"`
	Key   *string `json:"key"`
}

// GetPublicKey gets a public key that should be used for secret encryption.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-an-organization-public-key
func (s *OrganizationsService) GetPublicKey(ctx context.Context, owner string) (*OrganizationPublicKey, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/public-key", owner)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pubKey := new(OrganizationPublicKey)
	resp, err := s.client.Do(ctx, req, pubKey)
	if err != nil {
		return nil, resp, err
	}

	return pubKey, resp, nil
}

type OrganizationSecret struct {
	Name                    string    `json:"name"`
	CreatedAt               Timestamp `json:"created_at"`
	UpdatedAt               Timestamp `json:"updated_at"`
	Visibility              string    `json:"visibility"`
	SelectedRepositoriesUrl string    `json:"selected_repositories_url"`
}

type OrganizationSecrets struct {
	TotalCount int                   `json:"total_count"`
	Secrets    []*OrganizationSecret `json:"secrets"`
}

// ListSecrets lists all secrets available in an Organization
// without revealing their encrypted values.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-organization-secrets
func (s *OrganizationsService) ListSecrets(ctx context.Context, owner string, opts *ListOptions) (*OrganizationSecrets, *Response, error) {
	u := fmt.Sprintf("orgs/%s/actions/secrets", owner)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secrets := new(OrganizationSecrets)
	resp, err := s.client.Do(ctx, req, &secrets)
	if err != nil {
		return nil, resp, err
	}

	return secrets, resp, nil
}

// GetSecret gets a single secret without revealing its encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-an-organization-secret
func (s *OrganizationsService) GetSecret(ctx context.Context, owner, name string) (*OrganizationSecret, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", owner, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secret := new(OrganizationSecret)
	resp, err := s.client.Do(ctx, req, secret)
	if err != nil {
		return nil, resp, err
	}

	return secret, resp, nil
}

// OrganizationEncryptedSecret represents an Organization secret that is encrypted using a public key.
//
// The value of EncryptedValue must be your secret, encrypted with
// LibSodium (see documentation here: https://libsodium.gitbook.io/doc/bindings_for_other_languages)
// using the public key retrieved using the GetPublicKey method.
type OrganizationEncryptedSecret struct {
	Name                  string   `json:"-"`
	KeyID                 string   `json:"key_id"`
	EncryptedValue        string   `json:"encrypted_value"`
	Visibility            string   `json:"visibility"`
	SelectedRepositoryIDs []string `json:"selected_repository_ids,omitempty"`
}

// CreateOrUpdateSecret creates or updates a secret with an encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#create-or-update-an-organization-secret
func (s *OrganizationsService) CreateOrUpdateSecret(ctx context.Context, owner string, eSecret *OrganizationEncryptedSecret) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", owner, eSecret.Name)

	req, err := s.client.NewRequest("PUT", u, eSecret)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteSecret deletes a secret in a repository using the secret name.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#delete-an-organization-secret
func (s *OrganizationsService) DeleteSecret(ctx context.Context, owner, name string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v", owner, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// OrganizationSecretSelectedRepositories represents all repositories that have been selected to  have access to an OrganizationSecret
type OrganizationSecretSelectedRepositories struct {
	TotalCount   *int64        `json:"total_count,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
}

// Lists all repositories that have been selected when the visibility for repository access to a secret is set to selected.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-selected-repositories-for-an-organization-secret
func (s *OrganizationsService) ListSecretSelectedRepositories(ctx context.Context, owner, name string) (*OrganizationSecretSelectedRepositories, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories", owner, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secretSelectedRepositories := new(OrganizationSecretSelectedRepositories)
	resp, err := s.client.Do(ctx, req, secretSelectedRepositories)
	if err != nil {
		return nil, resp, err
	}

	return secretSelectedRepositories, resp, nil
}

// Replaces all repositories for an organization secret when the visibility for repository access is set to selected.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#set-selected-repositories-for-an-organization-secret
func (s *OrganizationsService) SetSecretSelectedRepositories(ctx context.Context, owner, name string, repositoryIDs []int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories", owner, name)

	req, err := s.client.NewRequest("PUT", u, repositoryIDs)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Adds a repository to an organization secret when the visibility for repository access is set to selected.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#add-selected-repository-to-an-organization-secret
func (s *OrganizationsService) AddSelectedRepositoryToSecret(ctx context.Context, owner, name string, repositoryID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories/%v", owner, name, repositoryID)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Removes a repository from an organization secret when the visibility for repository access is set to selected
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#remove-selected-repository-from-an-organization-secret
func (s *OrganizationsService) RemoveSelectedRepositoryFromSecret(ctx context.Context, owner, name string, repositoryID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/%v/repositories/%v", owner, name, repositoryID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
