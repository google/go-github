// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "context"

// SecretScanningPatternConfigs represents a collection of GitHub secret scanning patterns
// and their settings related to push protection.
type SecretScanningPatternConfigs struct {
	//TODO:
}

// SecretScanningPatternConfigsUpdate represents a secret scanning pattern configurations update.
type SecretScanningPatternConfigsUpdate struct {
	//TODO:
}

// SecretScanningPatternConfigsUpdateOptions specifies optional parameters to
// the SecretScanningService.UpdatePatternConfigsForEnterprise method and
// the SecretScanningService.UpdatePatternConfigsForOrg method.
type SecretScanningPatternConfigsUpdateOptions struct {
	//TODO:
}

// ListPatternConfigsForEnterprise lists the secret scanning pattern configurations for an enterprise.
//
// Github API docs: https://docs.github.com/enterprise-cloud@latest/rest/secret-scanning/push-protection?apiVersion=2022-11-28#list-enterprise-pattern-configurations
//
//meta:operation GET /enterprises/{enterprise}/secret-scanning/pattern-configurations
func (s *SecretScanningService) ListPatternConfigsForEnterprise(ctx context.Context, enterprise string) (*SecretScanningPatternConfigs, *Response, error) {
	//TODO:
}

// UpdatePatternConfigsForEnterprise updates the secret scanning pattern configurations for an enterprise.
//
// Github API docs: https://docs.github.com/enterprise-cloud@latest/rest/secret-scanning/push-protection?apiVersion=2022-11-28#update-enterprise-pattern-configurations
//
//meta:operation PATCH /enterprises/{enterprise}/secret-scanning/pattern-configurations
func (s *SecretScanningService) UpdatePatternConfigsForEnterprise(ctx context.Context, enterprise string, opts *SecretScanningPatternConfigsUpdateOptions) (*SecretScanningPatternConfigsUpdate, *Response, error) {
	//TODO:
}

// ListPatternConfigsForOrg lists the secret scanning pattern configurations for an organization.
//
// Github API docs: https://docs.github.com/enterprise-cloud@latest/rest/secret-scanning/push-protection?apiVersion=2022-11-28#list-organization-pattern-configurations
//
//meta:operation GET /orgs/{org}/secret-scanning/pattern-configurations
func (s *SecretScanningService) ListPatternConfigsForOrg(ctx context.Context, org string) (*SecretScanningPatternConfigs, *Response, error) {
	//TODO:
}

// UpdatePatternConfigsForOrg updates the secret scanning pattern configurations for an organization.
//
// Github API docs: https://docs.github.com/enterprise-cloud@latest/rest/secret-scanning/push-protection?apiVersion=2022-11-28#update-organization-pattern-configurations
//
//meta:operation PATCH /orgs/{org}/secret-scanning/pattern-configurations
func (s *SecretScanningService) UpdatePatternConfigsForOrg(ctx context.Context, org string, opts *SecretScanningPatternConfigsUpdateOptions) (*SecretScanningPatternConfigsUpdate, *Response, error) {
	//TODO:
}
