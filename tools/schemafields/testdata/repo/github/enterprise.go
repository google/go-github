// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// EnterpriseRunnerGroupsService is a stand-in for the real service.
type EnterpriseRunnerGroupsService struct{}

// CreateEnterpriseRunnerGroupRequest is documented only by the GHEC plan.
type CreateEnterpriseRunnerGroupRequest struct {
	// Name is required by the schema.
	Name string `json:"name"`

	// SelectedWorkflows is optional.
	SelectedWorkflows []string `json:"selected_workflows"`
}

//meta:operation POST /enterprises/{enterprise}/runner-groups
func (s *EnterpriseRunnerGroupsService) CreateEnterpriseRunnerGroup(body CreateEnterpriseRunnerGroupRequest) error {
	return nil
}
