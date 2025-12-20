// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CodespaceMachine  represent a description of the machine powering a codespace.
type CodespaceMachine struct {
	Name            *string `json:"name,omitempty"`
	DisplayName     *string `json:"display_name,omitempty"`
	OperatingSystem *string `json:"operating_system,omitempty"`
	StorageInBytes  *int64  `json:"storage_in_bytes,omitempty"`
	MemoryInBytes   *int64  `json:"memory_in_bytes,omitempty"`
	CPUs            *int    `json:"cpus,omitempty"`
	// PrebuildAvailability represents whether a prebuild is currently available when creating a codespace for this machine and repository.
	// Value will be "null" if prebuilds are not supported or prebuild availability could not be determined.
	// Value will be "none" if no prebuild is available.
	// Latest values "ready" and "in_progress" indicate the prebuild availability status.
	PrebuildAvailability *string `json:"prebuild_availability,omitempty"`
}

// CodespaceMachines represent a list of machines.
type CodespaceMachines struct {
	TotalCount int64               `json:"total_count"`
	Machines   []*CodespaceMachine `json:"machines"`
}

// ListMachinesOptions represent options for ListMachineTypesForRepository.
type ListMachinesOptions struct {
	// Ref represent the branch or commit to check for prebuild availability and devcontainer restrictions.
	Ref *string `json:"ref,omitempty"`
	// Location represent the location to check for available machines. Assigned by IP if not provided.
	Location *string `json:"location,omitempty"`
	// ClientIP represent the IP for location auto-detection when proxying a request
	ClientIP *string `json:"client_ip,omitempty"`
}

// ListMachineTypesForRepository lists the machine types available for a given repository based on its configuration.
//
// GitHub API docs: https://docs.github.com/rest/codespaces/machines#list-available-machine-types-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/codespaces/machines
func (s *CodespacesService) ListMachineTypesForRepository(ctx context.Context, owner, repo string, opts *ListMachinesOptions) (*CodespaceMachines, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/codespaces/machines", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var machines *CodespaceMachines
	resp, err := s.client.Do(ctx, req, &machines)
	if err != nil {
		return nil, resp, err
	}

	return machines, resp, nil
}

// ListMachineTypesForCodespace lists the machine types a codespace can transition to use.
//
// GitHub API docs: https://docs.github.com/rest/codespaces/machines#list-machine-types-for-a-codespace
//
//meta:operation GET /user/codespaces/{codespace_name}/machines
func (s *CodespacesService) ListMachineTypesForCodespace(ctx context.Context, codespaceName string) (*CodespaceMachines, *Response, error) {
	u := fmt.Sprintf("user/codespaces/%v/machines", codespaceName)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var machines *CodespaceMachines
	resp, err := s.client.Do(ctx, req, &machines)
	if err != nil {
		return nil, resp, err
	}

	return machines, resp, nil
}
