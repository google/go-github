// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type SecurityAdvisoriesService service

// SecurityAdvisorySubmission represents the Security Advisory Submission.
type SecurityAdvisorySubmission struct {
	// Accepted represents whether a private vulnerability report was accepted by the repository's administrators.
	Accepted *bool `json:"accepted,omitempty"`
}

// RepoAdvisoryCredit represents the credit object for a repository Security Advisory.
type RepoAdvisoryCredit struct {
	Login *string `json:"login,omitempty"`
	Type  *string `json:"type,omitempty"`
}

// RepoAdvisoryCreditDetailed represents a credit given to a user for a repository Security Advisory.
type RepoAdvisoryCreditDetailed struct {
	User  *User   `json:"user,omitempty"`
	Type  *string `json:"type,omitempty"`
	State *string `json:"state,omitempty"`
}

// ListRepositorySecurityAdvisoriesOptions specifies the optional parameters to list the repository security advisories.
type ListRepositorySecurityAdvisoriesOptions struct {
	ListCursorOptions

	// Direction in which to sort advisories. Possible values are: asc, desc.
	// Default is "asc".
	Direction string `url:"direction,omitempty"`

	// Sort specifies how to sort advisories. Possible values are: created, updated,
	// and published. Default value is "created".
	Sort string `url:"sort,omitempty"`

	// State filters advisories based on their state. Possible values are: triage, draft, published, closed.
	State string `url:"state,omitempty"`
}

// RequestCVE requests a Common Vulnerabilities and Exposures (CVE) for a repository security advisory.
// The ghsaID is the GitHub Security Advisory identifier of the advisory.
//
// GitHub API docs: https://docs.github.com/en/rest/security-advisories/repository-advisories#request-a-cve-for-a-repository-security-advisory
func (s *SecurityAdvisoriesService) RequestCVE(ctx context.Context, owner, repo, ghsaID string) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/security-advisories/%v/cve", owner, repo, ghsaID)

	req, err := s.client.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if _, ok := err.(*AcceptedError); ok {
			return resp, nil
		}

		return resp, err
	}

	return resp, nil
}

// ListRepositorySecurityAdvisoriesForOrg lists the repository security advisories for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/security-advisories/repository-advisories?apiVersion=2022-11-28#list-repository-security-advisories-for-an-organization
func (s *SecurityAdvisoriesService) ListRepositorySecurityAdvisoriesForOrg(ctx context.Context, org string, opt *ListRepositorySecurityAdvisoriesOptions) ([]*SecurityAdvisory, *Response, error) {
	url := fmt.Sprintf("orgs/%v/security-advisories", org)
	url, err := addOptions(url, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var advisories []*SecurityAdvisory
	resp, err := s.client.Do(ctx, req, &advisories)
	if err != nil {
		return nil, resp, err
	}

	return advisories, resp, nil
}

// ListRepositorySecurityAdvisories lists the security advisories in a repository.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/security-advisories/repository-advisories?apiVersion=2022-11-28#list-repository-security-advisories
func (s *SecurityAdvisoriesService) ListRepositorySecurityAdvisories(ctx context.Context, owner, repo string, opt *ListRepositorySecurityAdvisoriesOptions) ([]*SecurityAdvisory, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/security-advisories", owner, repo)
	url, err := addOptions(url, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var advisories []*SecurityAdvisory
	resp, err := s.client.Do(ctx, req, &advisories)
	if err != nil {
		return nil, resp, err
	}

	return advisories, resp, nil
}
