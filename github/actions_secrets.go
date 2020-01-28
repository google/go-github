// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type PublicKey struct {
	ID  *string `json:"key_id,omitempty"`
	Key *string `json:"key,omitempty"`
}

// GetPublicKey gets a public key that should be used for secret encryption.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-your-public-key
func (s *ActionsService) GetPublicKey(ctx context.Context, owner string, repo string) (*PublicKey, *Response, error) {
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

// Secret represents a repository action secret.
type Secret struct {
	Name      *string    `json:"name,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

type secrets struct {
	Secrets []*Secret `json:"secrets,omitempty"`
}

// ListSecrets lists all secrets available in a repository
// without revealing their encrypted values.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#list-secrets-for-a-repository
func (s *ActionsService) ListSecrets(ctx context.Context, owner, repo string, opt *ListOptions) ([]*Secret, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/secrets", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	secrets := new(secrets)
	resp, err := s.client.Do(ctx, req, &secrets)
	if err != nil {
		return nil, resp, err
	}

	return secrets.Secrets, resp, nil
}

// GetSecret gets a single secret without revealing its encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#get-a-secret
func (s *ActionsService) GetSecret(ctx context.Context, owner, repo, name string) (*Secret, *Response, error) {
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

type EncryptedSecret struct {
	Name           *string `json:"-"`
	KeyId          *string `json:"key_id,omitempty"`
	EncryptedValue *string `json:"encrypted_value,omitempty"`
}

// CreateOrUpdateSecret creates or updates a secret with an encrypted value.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#create-or-update-a-secret-for-a-repository
func (s *ActionsService) CreateOrUpdateSecret(ctx context.Context, owner, repo string, eSecret *EncryptedSecret) (*EncryptedSecret, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/%v", owner, repo, *eSecret.Name)

	req, err := s.client.NewRequest("PUT", u, eSecret)
	if err != nil {
		return nil, nil, err
	}

	es := new(EncryptedSecret)
	resp, err := s.client.Do(ctx, req, es)
	if err != nil {
		return nil, resp, err
	}

	return es, resp, nil
}

// DeleteSecret deletes a secret in a repository using the secret name.
//
// GitHub API docs: https://developer.github.com/v3/actions/secrets/#delete-a-secret-from-a-repository
func (s *ActionsService) DeleteSecret(ctx context.Context, owner, repo, name string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/secrets/%v", owner, repo, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
