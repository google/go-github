// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CodeScanningService handles communication with the code scanning related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/code-scanning/
type CodeScanningService service

type Alert struct {
	RuleID          *string    `json:"rule_id,omitempty"`
	RuleSeverity    *string    `json:"rule_severity,omitempty"`
	RuleDescription *string    `json:"rule_description,omitempty"`
	Tool            *string    `json:"tool,omitempty"`
	CreatedAt       *Timestamp `json:"created_at,omitempty"`
	Open            *bool      `json:"open,omitempty"`
	ClosedBy        *string    `json:"closed_by,omitempty"`
	ClosedAt        *Timestamp `json:"closed_at,omitempty"`
	URL             *string    `json:"url,omitempty"`
	HTMLURL         *string    `json:"html_url,omitempty"`
}

// AlertListOptions specifies optional parameters to the CodeScanningService.ListAlerts
// method.
type AlertListOptions struct {
	// State of the code scanning alerts to list. Set to closed to list only closed code scanning alerts. Default: open
	State string `url:"state,omitempty"`

	// Return code scanning alerts for a specific branch reference. The ref must be formatted as heads/<branch name>.
	Ref string `url:"ref,omitempty"`
}

// ListAlertsForRepo lists code scanning alerts for a repository.
//
// Lists all open code scanning alerts for the default branch (usually master) and protected branches in a repository.
// You must use an access token with the security_events scope to use this endpoint. GitHub Apps must have the security_events
// read permission to use this endpoint.
//
// GitHub API docs: https://developer.github.com/v3/code-scanning/#list-code-scanning-alerts-for-a-repository
func (s *CodeScanningService) ListAlertsForRepo(ctx context.Context, owner, repo string, opts *AlertListOptions) ([]*Alert, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/code-scanning/alerts", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var alerts []*Alert
	resp, err := s.client.Do(ctx, req, &alerts)
	if err != nil {
		return nil, resp, err
	}

	return alerts, resp, nil
}

// GetAlert gets a single code scanning alert for a repository.
//
// You must use an access token with the security_events scope to use this endpoint.
// GitHub Apps must have the security_events read permission to use this endpoint.
//
// The security alert_id is the number at the end of the security alert's URL.
//
// GitHub API docs: https://developer.github.com/v3/code-scanning/#get-a-code-scanning-alert
func (s *CodeScanningService) GetAlert(ctx context.Context, owner, repo string, id int64) (*Alert, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/code-scanning/alerts/%v", owner, repo, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Alert)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}
