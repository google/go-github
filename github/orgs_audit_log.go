// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetAuditLogOptions sets up optional parameters to query audit-log endpoint.
type GetAuditLogOptions struct {
	Phrase  *string `url:"phrase,omitempty"`  // A search phrase. (Optional.)
	Include *string `url:"include,omitempty"` // Event type includes. Can be one of "web", "git", "all". Default: "web". (Optional.)
	Order   *string `url:"order,omitempty"`   // The order of audit log events. Can be one of "asc" or "desc". Default: "desc". (Optional.)

	ListCursorOptions
}

// ActorLocation contains information about reported location for an actor.
type ActorLocation struct {
	CountryCode *string `json:"country_code,omitempty"`
}

// AuditEntry describes the fields that may be represented by various audit-log "action" entries.
// There are many other fields that may be present depending on the action. You can access those
// in AdditionalFields.
// For a list of actions see - https://docs.github.com/github/setting-up-and-managing-organizations-and-teams/reviewing-the-audit-log-for-your-organization#audit-log-actions
type AuditEntry struct {
	Action                   *string        `json:"action,omitempty"` // The name of the action that was performed, for example `user.login` or `repo.create`.
	Actor                    *string        `json:"actor,omitempty"`  // The actor who performed the action.
	ActorID                  *int64         `json:"actor_id,omitempty"`
	ActorLocation            *ActorLocation `json:"actor_location,omitempty"`
	Business                 *string        `json:"business,omitempty"`
	BusinessID               *int64         `json:"business_id,omitempty"`
	CreatedAt                *Timestamp     `json:"created_at,omitempty"`
	DocumentID               *string        `json:"_document_id,omitempty"`
	ExternalIdentityNameID   *string        `json:"external_identity_nameid,omitempty"`
	ExternalIdentityUsername *string        `json:"external_identity_username,omitempty"`
	HashedToken              *string        `json:"hashed_token,omitempty"`
	Org                      *string        `json:"org,omitempty"`
	OrgID                    *int64         `json:"org_id,omitempty"`
	Timestamp                *Timestamp     `json:"@timestamp,omitempty"` // The time the audit log event occurred, given as a [Unix timestamp](https://en.wikipedia.org/wiki/Unix_time).
	TokenID                  *int64         `json:"token_id,omitempty"`
	TokenScopes              *string        `json:"token_scopes,omitempty"`
	User                     *string        `json:"user,omitempty"` // The user that was affected by the action performed (if available).
	UserID                   *int64         `json:"user_id,omitempty"`

	// Some events types have a data field that contains additional information about the event.
	Data map[string]any `json:"data,omitempty"`

	// All fields that are not explicitly defined in the struct are captured here.
	AdditionalFields map[string]any `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// GitHub's audit-log API occasionally returns "org" as a JSON array of strings
// and "org_id" as a JSON array of integers instead of the documented scalar
// types.  This implementation normalises both fields to their scalar forms
// (joining multiple org names with a comma, and using the first org_id) so
// callers always receive a consistent type regardless of the API response shape.
func (a *AuditEntry) UnmarshalJSON(data []byte) error {
	// rawEntry shadows Org and OrgID so we can inspect their raw JSON tokens
	// before deciding how to decode them.
	type entryAlias AuditEntry
	var raw struct {
		entryAlias
		Org   json.RawMessage `json:"org,omitempty"`
		OrgID json.RawMessage `json:"org_id,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	v := raw.entryAlias

	// Normalise "org": accept both "string" and ["string", ...].
	if len(raw.Org) > 0 && string(raw.Org) != "null" {
		org, err := unmarshalStringOrStringArray(raw.Org)
		if err != nil {
			return fmt.Errorf("AuditEntry.Org: %w", err)
		}
		v.Org = org
	}

	// Normalise "org_id": accept both integer and [integer, ...].
	if len(raw.OrgID) > 0 && string(raw.OrgID) != "null" {
		orgID, err := unmarshalInt64OrInt64Array(raw.OrgID)
		if err != nil {
			return fmt.Errorf("AuditEntry.OrgID: %w", err)
		}
		v.OrgID = orgID
	}

	rawDefinedFields, err := json.Marshal(v)
	if err != nil {
		return err
	}
	definedFields := map[string]any{}
	if err := json.Unmarshal(rawDefinedFields, &definedFields); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v.AdditionalFields); err != nil {
		return err
	}

	for key, val := range v.AdditionalFields {
		if _, ok := definedFields[key]; ok || val == nil {
			delete(v.AdditionalFields, key)
		}
	}

	*a = AuditEntry(v)
	if len(v.AdditionalFields) == 0 {
		a.AdditionalFields = nil
	}
	return nil
}

// unmarshalStringOrStringArray decodes a JSON value that is either a plain
// string or an array of strings.  Arrays are joined with ", ".
func unmarshalStringOrStringArray(raw json.RawMessage) (*string, error) {
	// Try scalar string first (the common case).
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return &s, nil
	}
	// Fall back to array of strings.
	var arr []string
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, err
	}
	joined := strings.Join(arr, ", ")
	return &joined, nil
}

// unmarshalInt64OrInt64Array decodes a JSON value that is either a plain
// integer or an array of integers.  Arrays use the first element.
func unmarshalInt64OrInt64Array(raw json.RawMessage) (*int64, error) {
	// Try scalar integer first (the common case).
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n, nil
	}
	// Fall back to array of integers; use the first element.
	var arr []int64
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, nil
	}
	return &arr[0], nil
}

// MarshalJSON implements the json.Marshaler interface.
func (a AuditEntry) MarshalJSON() ([]byte, error) {
	type entryAlias AuditEntry
	v := entryAlias(a)
	defBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(a.AdditionalFields) == 0 {
		return defBytes, err
	}
	resMap := map[string]any{}
	if err := json.Unmarshal(defBytes, &resMap); err != nil {
		return nil, err
	}
	for key, val := range a.AdditionalFields {
		if val == nil {
			continue
		}
		if _, ok := resMap[key]; ok {
			return nil, fmt.Errorf("unexpected field in AdditionalFields: %v", key)
		}
		resMap[key] = val
	}
	return json.Marshal(resMap)
}

// GetAuditLog gets the audit-log entries for an organization.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/orgs/orgs?apiVersion=2022-11-28#get-the-audit-log-for-an-organization
//
//meta:operation GET /orgs/{org}/audit-log
func (s *OrganizationsService) GetAuditLog(ctx context.Context, org string, opts *GetAuditLogOptions) ([]*AuditEntry, *Response, error) {
	u := fmt.Sprintf("orgs/%v/audit-log", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var auditEntries []*AuditEntry
	resp, err := s.client.Do(req, &auditEntries)
	if err != nil {
		return nil, resp, err
	}

	return auditEntries, resp, nil
}
