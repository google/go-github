// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetRepoDependabotPublicKey gets a public key that should be used for Dependabot secret encryption.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#get-a-repository-public-key
func (s *ActionsService) GetRepoDependabotPublicKey(ctx context.Context, owner, repo string) (*PublicKey, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependabot/secrets/public-key", owner, repo)
	return s.getPublicKey(ctx, url)
}

// GetOrgDependabotPublicKey gets a public key that should be used for Dependabot secret encryption.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#get-an-organization-public-key
func (s *ActionsService) GetOrgDependabotPublicKey(ctx context.Context, org string) (*PublicKey, *Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/public-key", org)
	return s.getPublicKey(ctx, url)
}

// ListRepoDependabotSecrets lists all Dependabot secrets available in a repository
// without revealing their encrypted values.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#list-repository-secrets
func (s *ActionsService) ListRepoDependabotSecrets(ctx context.Context, owner, repo string, opts *ListOptions) (*Secrets, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependabot/secrets", owner, repo)
	return s.listSecrets(ctx, url, opts)
}

// ListOrgDependabotSecrets lists all Dependabot secrets available in an organization
// without revealing their encrypted values.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#list-organization-secrets
func (s *ActionsService) ListOrgDependabotSecrets(ctx context.Context, org string, opts *ListOptions) (*Secrets, *Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets", org)
	return s.listSecrets(ctx, url, opts)
}

// GetRepoDependabotSecret gets a single repository Dependabot secret without revealing its encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#get-a-repository-secret
func (s *ActionsService) GetRepoDependabotSecret(ctx context.Context, owner, repo, name string) (*Secret, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependabot/secrets/%v", owner, repo, name)
	return s.getSecret(ctx, url)
}

// GetOrgDependabotSecret gets a single organization Dependabot secret without revealing its encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#get-an-organization-secret
func (s *ActionsService) GetOrgDependabotSecret(ctx context.Context, org, name string) (*Secret, *Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v", org, name)
	return s.getSecret(ctx, url)
}

// CreateOrUpdateRepoDependabotSecret creates or updates a repository Dependabot secret with an encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#create-or-update-a-repository-secret
func (s *ActionsService) CreateOrUpdateRepoDependabotSecret(ctx context.Context, owner, repo string, eSecret *EncryptedSecret) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependabot/secrets/%v", owner, repo, eSecret.Name)
	return s.putSecret(ctx, url, eSecret)
}

// CreateOrUpdateOrgDependabotSecret creates or updates an organization Dependabot secret with an encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#create-or-update-an-organization-secret
func (s *ActionsService) CreateOrUpdateOrgDependabotSecret(ctx context.Context, org string, eSecret *EncryptedSecret) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v", org, eSecret.Name)
	return s.putSecret(ctx, url, eSecret)
}

// DeleteRepoDependabotSecret deletes a Dependabot secret in a repository using the secret name.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#delete-a-repository-secret
func (s *ActionsService) DeleteRepoDependabotSecret(ctx context.Context, owner, repo, name string) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/dependabot/secrets/%v", owner, repo, name)
	return s.deleteSecret(ctx, url)
}

// DeleteOrgDependabotSecret deletes a Dependabot secret in an organization using the secret name.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#delete-an-organization-secret
func (s *ActionsService) DeleteOrgDependabotSecret(ctx context.Context, org, name string) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v", org, name)
	return s.deleteSecret(ctx, url)
}

// ListSelectedReposForOrgDependabotSecret lists all repositories that have access to a Dependabot secret.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#list-selected-repositories-for-an-organization-secret
func (s *ActionsService) ListSelectedReposForOrgDependabotSecret(ctx context.Context, org, name string, opts *ListOptions) (*SelectedReposList, *Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v/repositories", org, name)
	return s.listSelectedReposForSecret(ctx, url, opts)
}

// SetSelectedReposForOrgDependabotSecret sets the repositories that have access to a Dependabot secret.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#set-selected-repositories-for-an-organization-secret
func (s *ActionsService) SetSelectedReposForOrgDependabotSecret(ctx context.Context, org, name string, ids SelectedRepoIDs) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v/repositories", org, name)
	return s.setSelectedReposForSecret(ctx, url, ids)
}

// AddSelectedRepoToOrgDependabotSecret adds a repository to an organization Dependabot secret.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#add-selected-repository-to-an-organization-secret
func (s *ActionsService) AddSelectedRepoToOrgDependabotSecret(ctx context.Context, org, name string, repo *Repository) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v/repositories/%v", org, name, *repo.ID)
	return s.addSelectedRepoToSecret(ctx, url)
}

// RemoveSelectedRepoFromOrgDependabotSecret removes a repository from an organization Dependabot secret.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/dependabot#remove-selected-repository-from-an-organization-secret
func (s *ActionsService) RemoveSelectedRepoFromOrgDependabotSecret(ctx context.Context, org, name string, repo *Repository) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/dependabot/secrets/%v/repositories/%v", org, name, *repo.ID)
	return s.removeSelectedRepoFromSecret(ctx, url)
}
