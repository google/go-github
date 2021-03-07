// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetAuditLogOptions sets up optional parameters to query audit-log endpoint.
type GetAuditLogOptions struct {
	Phrase  *string `json:"phrase,omitempty"`   // A search phrase. (Optional.)
	Include *string `json:"include,omitempty"`  // Event type includes. Can be one of "web", "git", "all". Default: "web". (Optional.)
	After   *string `json:"after,omitempty"`    // A cursor, as given in the Link header. If specified, the query only searches for events after this cursor. (Optional.)
	Before  *string `json:"before,omitempty"`   // A cursor, as given in the Link header. If specified, the query only searches for events before this cursor. (Optional.)
	Order   *string `json:"order,omitempty"`    // The order of audit log events. Can be one of "asc" or "desc". Default: "desc". (Optional.)
	PerPage *int64  `json:"per_page,omitempty"` // Results per page (max 100).
}

// AuditEntry describes the fields that may be represented by various audit-log "action" entries.
type AuditEntry struct {
	Timestamp            *int64     `json:"@timestamp,omitempty"`
	DocumentID           *string    `json:"_document_id,omitempty"`
	Action               *string    `json:"action,omitempty"`
	Active               *string    `json:"active,omitempty"`
	ActiveWas            *string    `json:"active_was,omitempty"`
	Actor                *string    `json:"actor,omitempty"`
	BlockedUser          *string    `json:"blocked_user,omitempty"`
	CancelledAt          *Timestamp `json:"cancelled_at,omitempty"`
	CompletedAt          *Timestamp `json:"completed_at,omitempty"`
	Conclusion           *string    `json:"conclusion,omitempty"`
	Config               *string    `json:"config,omitempty"`
	ConfigWas            *string    `json:"config_was,omitempty"`
	ContentType          *string    `json:"content_type,omitempty"`
	CreatedAt            *int64     `json:"created_at,omitempty"`
	DeployKeyFingerprint *string    `json:"deploy_key_fingerprint,omitempty"`
	Emoji                *string    `json:"emoji,omitempty"`
	EnvironmentName      *string    `json:"environment_name,omitempty"`
	Event                *string    `json:"event,omitempty"`
	Events               *string    `json:"events,omitempty"`
	EventsWere           *string    `json:"events_were,omitempty"`
	Explanation          *string    `json:"explanation,omitempty"`
	Fingerprint          *string    `json:"fingerprint,omitempty"`
	HeadBranch           *string    `json:"head_branch,omitempty"`
	HeadSHA              *string    `json:"head_sha,omitempty"`
	HookID               *string    `json:"hook_id,omitempty"`
	IsHostedRunner       *bool      `json:"is_hosted_runner,omitempty"`
	JobName              *string    `json:"job_name,omitempty"`
	LimitedAvailability  *string    `json:"limited_availability,omitempty"`
	Message              *string    `json:"message,omitempty"`
	Name                 *string    `json:"name,omitempty"`
	OldUser              *string    `json:"old_user,omitempty"`
	OpenSSHPublicKey     *string    `json:"openssh_public_key,omitempty"`
	Org                  *string    `json:"org,omitempty"`
	PreviousVisibility   *string    `json:"previous_visibility,omitempty"`
	ReadOnly             *string    `json:"read_only,omitempty"`
	Repo                 *string    `json:"repo,omitempty"`
	RunnerGroupID        *string    `json:"runner_group_id,omitempty"`
	RunnerGroupName      *string    `json:"runner_group_name,omitempty"`
	RunnerID             *string    `json:"runner_id,omitempty"`
	RunnerLabels         *[]string  `json:"runner_labels,omitempty"`
	RunnerName           *string    `json:"runner_name,omitempty"`
	SecretsPassed        *[]string  `json:"secrets_passed,omitempty"`
	SourceVersion        *string    `json:"source_version,omitempty"`
	StartedAt            *Timestamp `json:"started_at,omitempty"`
	TargetLogin          *string    `json:"target_login,omitempty"`
	TargetVersion        *string    `json:"target_version,omitempty"`
	Team                 *string    `json:"team,omitempty"`
	TriggerID            *int64     `json:"trigger_id,omitempty"`
	User                 *string    `json:"user,omitempty"`
	Visibility           *string    `json:"visibility,omitempty"`
	WorkflowID           *int64     `json:"workflow_id,omitempty"`
	WorkflowRunID        *int64     `json:"workflow_run_id,omitempty"`
	Business             *string    `json:"business,omitemtpy"`
	Repository           *string    `json:"repository,omitempty"`
	RepositoryPublic     *bool      `json:"repository_public,omitempty"`
	TransportProtocol    *int64     `json:"transport_protocol,omitempty"`
	TransportName        *string    `json:"transport_protocol_name,omitempty"`
}

// GetAuditLog gets the audit-log entries for an organization
//
// GitHub API docs: https://docs.github.com/en/rest/reference/orgs#get-the-audit-log-for-an-organization
func (s *OrganizationsService) GetAuditLog(ctx context.Context, org string, opts *GetAuditLogOptions) ([]*AuditEntry, *Response, error) {
	u := fmt.Sprintf("orgs/%v/audit-log", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)

	req.Header.Set("Accept", mediaTypeV3)

	var auditEntry []*AuditEntry
	resp, err := s.client.Do(ctx, req, &auditEntry)
	if err != nil {
		return nil, resp, err
	}

	return auditEntry, resp, nil
}
