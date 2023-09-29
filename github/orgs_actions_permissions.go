// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// GetActionsPermissions gets the GitHub Actions permissions policy for repositories and allowed actions in an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#get-github-actions-permissions-for-an-organization
// Deprecated: please use `client.Actions.GetActionsPermissions` instead.
func (s *OrganizationsService) GetActionsPermissions(ctx context.Context, org string) (*ActionsPermissions, *Response, error) {
	s2 := (*ActionsService)(s)
	return s2.GetActionsPermissions(ctx, org)
}

// EditActionsPermissions sets the permissions policy for repositories and allowed actions in an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#set-github-actions-permissions-for-an-organization
// Deprecated: please use `client.Actions.EditActionsPermissions` instead.
func (s *OrganizationsService) EditActionsPermissions(ctx context.Context, org string, actionsPermissions ActionsPermissions) (*ActionsPermissions, *Response, error) {
	s2 := (*ActionsService)(s)
	return s2.EditActionsPermissions(ctx, org, actionsPermissions)
}
