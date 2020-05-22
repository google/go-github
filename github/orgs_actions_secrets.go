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
