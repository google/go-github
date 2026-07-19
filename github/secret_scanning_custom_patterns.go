// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// SecretScanningCustomPattern represents a custom pattern for secret scanning,
// as returned by the GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/secret-scanning/secret-scanning
type SecretScanningCustomPattern struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Slug    string `json:"slug"`

	// State is the publish state of the pattern. Possible values are:
	// "published" or "unpublished".
	State *string `json:"state,omitempty"`

	PushProtectionEnabled bool     `json:"push_protection_enabled"`
	StartDelimiter        *string  `json:"start_delimiter,omitempty"`
	EndDelimiter          *string  `json:"end_delimiter,omitempty"`
	MustMatch             []string `json:"must_match,omitempty"`
	MustNotMatch          []string `json:"must_not_match,omitempty"`

	// CustomPatternVersion is used to confirm you're updating the current
	// version of the pattern and to avoid unintentionally overriding
	// someone else's update.
	CustomPatternVersion *string    `json:"custom_pattern_version,omitempty"`
	CreatedAt            *Timestamp `json:"created_at,omitempty"`
	UpdatedAt            *Timestamp `json:"updated_at,omitempty"`
}

// SecretScanningCustomPatternRequest represents a custom pattern to be
// created for secret scanning.
type SecretScanningCustomPatternRequest struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`

	StartDelimiter *string  `json:"start_delimiter,omitempty"`
	EndDelimiter   *string  `json:"end_delimiter,omitempty"`
	MustMatch      []string `json:"must_match,omitempty"`
	MustNotMatch   []string `json:"must_not_match,omitempty"`
}

// SecretScanningCustomPatternsCreateRequest represents the bulk request body
// used to create one or more custom patterns.
type SecretScanningCustomPatternsCreateRequest struct {
	Patterns []*SecretScanningCustomPatternRequest `json:"patterns"`
}

// SecretScanningCustomPatternsCreateResponse represents the bulk response
// returned after creating one or more custom patterns.
type SecretScanningCustomPatternsCreateResponse struct {
	CreatedPatterns []*SecretScanningCustomPattern `json:"created_patterns"`
}

// SecretScanningCustomPatternToDelete identifies a single custom pattern to
// remove in a bulk delete operation.
type SecretScanningCustomPatternToDelete struct {
	PatternID int64 `json:"pattern_id"`

	// CustomPatternVersion is used to confirm you're deleting the current
	// version of the pattern.
	CustomPatternVersion *string `json:"custom_pattern_version,omitempty"`
}

// SecretScanningCustomPatternsDeleteRequest represents the bulk request body
// used to delete one or more custom patterns.
type SecretScanningCustomPatternsDeleteRequest struct {
	Patterns []*SecretScanningCustomPatternToDelete `json:"patterns"`

	// PostDeleteAction controls what happens to alerts associated with the
	// deleted patterns. Possible values are: "delete_alerts" (default) or
	// "resolve_alerts".
	PostDeleteAction *string `json:"post_delete_action,omitempty"`
}

// SecretScanningCustomPatternUpdateRequest represents the fields that can be
// updated on an existing custom pattern. At least one of Pattern,
// StartDelimiter, EndDelimiter, MustMatch, or MustNotMatch must be set.
type SecretScanningCustomPatternUpdateRequest struct {
	Pattern        *string  `json:"pattern,omitempty"`
	StartDelimiter *string  `json:"start_delimiter,omitempty"`
	EndDelimiter   *string  `json:"end_delimiter,omitempty"`
	MustMatch      []string `json:"must_match,omitempty"`
	MustNotMatch   []string `json:"must_not_match,omitempty"`

	// CustomPatternVersion is required and is used to confirm you're
	// updating the current version of the pattern.
	CustomPatternVersion *string `json:"custom_pattern_version"`
}

// ListCustomPatternsForRepo lists the secret scanning custom patterns
// defined at the repository level.
//
// GitHub API docs: https://docs.github.com/rest/secret-scanning/custom-patterns?apiVersion=2022-11-28#list-repository-custom-patterns
//
//meta:operation GET /repos/{owner}/{repo}/secret-scanning/custom-patterns
func (s *SecretScanningService) ListCustomPatternsForRepo(ctx context.Context, owner, repo string) ([]*SecretScanningCustomPattern, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/secret-scanning/custom-patterns", owner, repo)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var patterns []*SecretScanningCustomPattern
	resp, err := s.client.Do(req, &patterns)
	if err != nil {
		return nil, resp, err
	}

	return patterns, resp, nil
}

// CreateCustomPatternsForRepo creates one or more secret scanning custom
// patterns at the repository level.
//
// GitHub API docs: https://docs.github.com/rest/secret-scanning/custom-patterns?apiVersion=2022-11-28#bulk-create-repository-custom-patterns
//
//meta:operation POST /repos/{owner}/{repo}/secret-scanning/custom-patterns
func (s *SecretScanningService) CreateCustomPatternsForRepo(ctx context.Context, owner, repo string, body SecretScanningCustomPatternsCreateRequest) (*SecretScanningCustomPatternsCreateResponse, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/secret-scanning/custom-patterns", owner, repo)

	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var result *SecretScanningCustomPatternsCreateResponse
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateCustomPatternForRepo updates a single secret scanning custom pattern
// at the repository level.
//
// GitHub API docs: https://docs.github.com/rest/secret-scanning/custom-patterns?apiVersion=2022-11-28#update-a-repository-custom-pattern
//
//meta:operation PATCH /repos/{owner}/{repo}/secret-scanning/custom-patterns/{pattern_id}
func (s *SecretScanningService) UpdateCustomPatternForRepo(ctx context.Context, owner, repo string, patternID int64, body SecretScanningCustomPatternUpdateRequest) (*SecretScanningCustomPattern, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/secret-scanning/custom-patterns/%v", owner, repo, patternID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, body)
	if err != nil {
		return nil, nil, err
	}

	var pattern *SecretScanningCustomPattern
	resp, err := s.client.Do(req, &pattern)
	if err != nil {
		return nil, resp, err
	}

	return pattern, resp, nil
}

// DeleteCustomPatternsForRepo deletes one or more secret scanning custom
// patterns at the repository level.
//
// GitHub API docs: https://docs.github.com/rest/secret-scanning/custom-patterns?apiVersion=2022-11-28#bulk-delete-repository-custom-patterns
//
//meta:operation DELETE /repos/{owner}/{repo}/secret-scanning/custom-patterns
func (s *SecretScanningService) DeleteCustomPatternsForRepo(ctx context.Context, owner, repo string, patterns *SecretScanningCustomPatternsDeleteRequest) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/secret-scanning/custom-patterns", owner, repo)

	req, err := s.client.NewRequest(ctx, "DELETE", u, patterns)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
