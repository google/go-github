// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// MaintenanceOperationStatus represents the message to be displayed when the instance gets a maintenance operation request.
type MaintenanceOperationStatus struct {
	Hostname *string `json:"hostname"`
	UUID     *string `json:"uuid"`
	Message  *string `json:"message"`
}

// MaintenanceStatus represents the status of maintenance mode for all nodes.
type MaintenanceStatus struct {
	Hostname               *string               `json:"hostname"`
	UUID                   *string               `json:"uuid"`
	Status                 *string               `json:"status"`
	ScheduledTime          *Timestamp            `json:"scheduled_time"`
	ConnectionServices     []*ConnectionServices `json:"connection_services"`
	CanUnsetMaintenance    *bool                 `json:"can_unset_maintenance"`
	IPExceptionList        []*string             `json:"ip_exception_list"`
	MaintenanceModeMessage *string               `json:"maintenance_mode_message"`
}

// ConnectionServices represents the connection services for the maintenance status.
type ConnectionServices struct {
	Name   *string `json:"name"`
	Number *int    `json:"number"`
}

// MaintenanceOptions represents the options for setting the maintenance mode for the instance.
// When can be a string, so we cant use a Timestamp type.
type MaintenanceOptions struct {
	Enabled                *bool     `json:"enabled,omitempty"`
	UUID                   *string   `json:"uuid,omitempty"`
	When                   *string   `json:"when,omitempty"`
	IPExceptionList        []*string `json:"ip_exception_list,omitempty"`
	MaintenanceModeMessage *string   `json:"maintenance_mode_message,omitempty"`
}

// GetMaintenanceStatus gets the status of maintenance mode for all nodes.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#get-the-status-of-maintenance-mode
//
//meta:operation GET /manage/v1/maintenance
func (s *EnterpriseService) GetMaintenanceStatus(ctx context.Context, opts *NodeQueryOptions) ([]*MaintenanceStatus, *Response, error) {
	u, err := addOptions("manage/v1/maintenance", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var status []*MaintenanceStatus
	resp, err := s.client.Do(ctx, req, &status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, nil
}

// CreateMaintenance sets the maintenance mode for the instance.
// With the enable parameter we can control to put instance into maintenance mode or not. With False we can disable the maintenance mode.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.15/rest/enterprise-admin/manage-ghes#set-the-status-of-maintenance-mode
//
//meta:operation POST /manage/v1/maintenance
func (s *EnterpriseService) CreateMaintenance(ctx context.Context, enable bool, opts *MaintenanceOptions) ([]*MaintenanceOperationStatus, *Response, error) {
	u := "manage/v1/maintenance"

	opts.Enabled = &enable

	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var i []*MaintenanceOperationStatus
	resp, err := s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}
