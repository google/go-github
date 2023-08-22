// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// GetAuditLog gets the audit-log entries for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-admin/audit-log#get-the-audit-log-for-an-enterprise
func (s *EnterpriseService) GetAuditLog(ctx context.Context, enterprise string, opts *GetAuditLogOptions) ([]*AuditEntry, *Response, error) {
	u, err := newURLString("enterprises/%v/audit-log", enterprise)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var auditEntries []*AuditEntry
	resp, err := s.client.Do(ctx, req, &auditEntries)
	if err != nil {
		return nil, resp, err
	}

	return auditEntries, resp, nil
}
