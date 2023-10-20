// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "context"

// RateLimitService provides access to rate limit functions in the GitHub API.
type RateLimitService service

// Rate represents the rate limit for the current client.
type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
}

func (r Rate) String() string {
	return Stringify(r)
}

// RateLimits represents the rate limits for the current client.
type RateLimits struct {
	// The rate limit for non-search API requests. Unauthenticated
	// requests are limited to 60 per hour. Authenticated requests are
	// limited to 5,000 per hour.
	//
	// GitHub API docs: https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting
	Core *Rate `json:"core"`

	// The rate limit for search API requests. Unauthenticated requests
	// are limited to 10 requests per minutes. Authenticated requests are
	// limited to 30 per minute.
	//
	// GitHub API docs: https://docs.github.com/en/rest/search#rate-limit
	Search *Rate `json:"search"`

	// GitHub API docs: https://docs.github.com/en/graphql/overview/resource-limitations#rate-limit
	GraphQL *Rate `json:"graphql"`

	// GitHub API dos: https://docs.github.com/en/rest/rate-limit
	IntegrationManifest *Rate `json:"integration_manifest"`

	SourceImport              *Rate `json:"source_import"`
	CodeScanningUpload        *Rate `json:"code_scanning_upload"`
	ActionsRunnerRegistration *Rate `json:"actions_runner_registration"`
	SCIM                      *Rate `json:"scim"`
}

func (r RateLimits) String() string {
	return Stringify(r)
}

// RateLimits returns the rate limits for the current client.
func (s *RateLimitService) Get(ctx context.Context) (*RateLimits, *Response, error) {
	req, err := s.client.NewRequest("GET", "rate_limit", nil)
	if err != nil {
		return nil, nil, err
	}

	response := new(struct {
		Resources *RateLimits `json:"resources"`
	})

	// This resource is not subject to rate limits.
	ctx = context.WithValue(ctx, bypassRateLimitCheck, true)
	resp, err := s.client.Do(ctx, req, response)
	if err != nil {
		return nil, resp, err
	}

	if response.Resources != nil {
		s.client.rateMu.Lock()
		if response.Resources.Core != nil {
			s.client.rateLimits[coreCategory] = *response.Resources.Core
		}
		if response.Resources.Search != nil {
			s.client.rateLimits[searchCategory] = *response.Resources.Search
		}
		if response.Resources.GraphQL != nil {
			s.client.rateLimits[graphqlCategory] = *response.Resources.GraphQL
		}
		if response.Resources.IntegrationManifest != nil {
			s.client.rateLimits[integrationManifestCategory] = *response.Resources.IntegrationManifest
		}
		if response.Resources.SourceImport != nil {
			s.client.rateLimits[sourceImportCategory] = *response.Resources.SourceImport
		}
		if response.Resources.CodeScanningUpload != nil {
			s.client.rateLimits[codeScanningUploadCategory] = *response.Resources.CodeScanningUpload
		}
		if response.Resources.ActionsRunnerRegistration != nil {
			s.client.rateLimits[actionsRunnerRegistrationCategory] = *response.Resources.ActionsRunnerRegistration
		}
		if response.Resources.SCIM != nil {
			s.client.rateLimits[scimCategory] = *response.Resources.SCIM
		}
		s.client.rateMu.Unlock()
	}

	return response.Resources, resp, nil
}
