// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// EnterpriseSecurityAnalysisSettings represents security analysis settings for an enterprise.
type EnterpriseSecurityAnalysisSettings struct {
	AdvancedSecurityEnabledForNewRepositories             bool   `json:"advanced_security_enabled_for_new_repositories,omitempty"`
	SecretScanningEnabledForNewRepositories               bool   `json:"secret_scanning_enabled_for_new_repositories,omitempty"`
	SecretScanningPushProtectionEnabledForNewRepositories bool   `json:"secret_scanning_push_protection_enabled_for_new_repositories,omitempty"`
	SecretScanningPushProtectionCustomLink                string `json:"secret_scanning_push_protection_custom_link,omitempty"`
}

// GetCodeSecurityAndAnalysis gets code security and analysis features for an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#get-code-security-and-analysis-features-for-an-enterprise
func (s *EnterpriseService) GetCodeSecurityAndAnalysis(ctx context.Context, enterprise string) (*EnterpriseSecurityAnalysisSettings, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/code_security_and_analysis", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	enterpriseSecurityAnalysisSettings := new(EnterpriseSecurityAnalysisSettings)
	resp, err := s.client.Do(ctx, req, enterpriseSecurityAnalysisSettings)
	if err != nil {
		return nil, resp, err
	}

	return enterpriseSecurityAnalysisSettings, resp, nil
}

// UpdateCodeSecurityAndAnalysis updates code security and analysis features for new repositories in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#update-code-security-and-analysis-features-for-an-enterprise
func (s *EnterpriseService) UpdateCodeSecurityAndAnalysis(ctx context.Context, enterprise string, enterpriseSecurityAnalysisSettings *EnterpriseSecurityAnalysisSettings) (*Response, error) {

	u := fmt.Sprintf("enterprises/%v/code_security_and_analysis", enterprise)
	req, err := s.client.NewRequest("PATCH", u, enterpriseSecurityAnalysisSettings)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// EnableAdvancedSecurity enables advanced security for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) EnableAdvancedSecurity(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/advanced_security/enable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DisableAdvancedSecurity disables advanced security for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) DisableAdvancedSecurity(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/advanced_security/disable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// EnableSecretScanning enables secret scanning for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) EnableSecretScanning(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/secret_scanning/enable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DisableSecretScanning disables secret scanning for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) DisableSecretScanning(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/secret_scanning/disable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// EnableSecretScanningPushProtection enables secret scanning push protection for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) EnableSecretScanningPushProtection(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/secret_scanning_push_protection/enable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DisableSecretScanningPushProtection disables secret scanning push protection for all repositories in an enterprise.
//
// GitHub API docs:https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/code-security-and-analysis?apiVersion=2022-11-28#enable-or-disable-a-security-feature
func (s *EnterpriseService) DisableSecretScanningPushProtection(ctx context.Context, enterprise string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/secret_scanning_push_protection/disable_all", enterprise)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
