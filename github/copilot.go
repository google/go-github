// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// CopilotService provides access to the Copilot-related functions
// in the GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/copilot/
type CopilotService service

// OrganizationDetails represents the details of an organization's Copilot for Business supbscription.
type OrganizationCopilotDetails struct {
	SeatBreakdown         *SeatBreakdown `json:"seat_breakdown"`
	PublicCodeSuggestions string         `json:"public_code_suggestions"`
	CopilotChat           string         `json:"copilot_chat"`
	SeatManagementSetting string         `json:"seat_management_setting"`
}

// SeatBreakdown represents the breakdown of Copilot for Business seats for the organization.
type SeatBreakdown struct {
	Total               int64 `json:"total"`
	AddedThisCycle      int64 `json:"added_this_cycle"`
	PendingCancellation int64 `json:"pending_cancellation"`
	PendingInvitation   int64 `json:"pending_invitation"`
	ActiveThisCycle     int64 `json:"active_this_cycle"`
	InactiveThisCycle   int64 `json:"inactive_this_cycle"`
}

type CopilotSeats struct {
	TotalSeats int64                `json:"total_seats"`
	Seats      []CopilotSeatDetails `json:"seats"`
}

// CopilotSeatDetail represents the details of a Copilot for Business seat.
// Assignee can either be a User, Team, or Organization.
type CopilotSeatDetails struct {
	Assignee                interface{} `json:"assignee"`
	AssigningTeam           *Team       `json:"assigning_team,omitempty"`
	PendingCancellationDate *string     `json:"pending_cancellation_date,omitempty"`
	LastActivityAt          *string     `json:"last_activity_at,omitempty"`
	LastActivityEditor      *string     `json:"last_activity_editor,omitempty"`
	CreatedAt               string      `json:"created_at"`
	UpdatedAt               string      `json:"updated_at,omitempty"`
}

// SelectedTeams represents the teams selected for the Copilot for Business subscription.
type SelectedTeams struct {
	SelectedTeams []string `json:"selected_teams"`
}

// SelectedUsers represents the users selected for the Copilot for Business subscription.
type SelectedUsers struct {
	SelectedUsers []string `json:"selected_users"`
}

// SeatAssignments represents the number of seats assigned.
type SeatAssignments struct {
	SeatsCreated int64 `json:"seats_created"`
}

// SeatCancellations represents the number of seats cancelled.
type SeatCancellations struct {
	SeatsCancelled int64 `json:"seats_cancelled"`
}

func (cp *CopilotSeatDetails) UnmarshalJSON(data []byte) error {
	type Alias CopilotSeatDetails
	var seatDetail Alias

	if err := json.Unmarshal(data, &seatDetail); err != nil {
		return err
	}

	cp.AssigningTeam = seatDetail.AssigningTeam
	cp.PendingCancellationDate = seatDetail.PendingCancellationDate
	cp.LastActivityAt = seatDetail.LastActivityAt
	cp.LastActivityEditor = seatDetail.LastActivityEditor
	cp.CreatedAt = seatDetail.CreatedAt
	cp.UpdatedAt = seatDetail.UpdatedAt

	switch v := seatDetail.Assignee.(type) {
	case map[string]interface{}:
		jsonData, err := json.Marshal(seatDetail.Assignee)
		if err != nil {
			return err
		}
		if v["type"].(string) == "User" {
			user := &User{}
			if err := json.Unmarshal(jsonData, user); err != nil {
				return err
			}
			cp.Assignee = user
		} else if v["type"].(string) == "Team" {
			team := &Team{}
			if err := json.Unmarshal(jsonData, team); err != nil {
				return err
			}
			cp.Assignee = team
		} else if v["type"].(string) == "Organization" {
			organization := &Organization{}
			if err := json.Unmarshal(jsonData, organization); err != nil {
				return err
			}
			cp.Assignee = organization
		} else {
			return fmt.Errorf("unsupported assignee type %s", v["type"].(string))
		}
	default:
		return fmt.Errorf("unsupported assignee type %T", v)
	}

	return nil
}

// Get Copilot for Business seat information and settings for an organization
//
// GitHub API docs: https://docs.github.com/en/rest/copilot/copilot-for-business#get-copilot-for-business-seat-information-and-settings-for-an-organization
func (s *CopilotService) GetCopilotBilling(ctx context.Context, org string) (*OrganizationCopilotDetails, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var copilotDetails *OrganizationCopilotDetails
	resp, err := s.client.Do(ctx, req, &copilotDetails)
	if err != nil {
		return nil, resp, err
	}

	return copilotDetails, resp, nil
}

// Get Copilot for Business seat assignments for an organization
//
// GitHub API docs: https://docs.github.com/en/rest/copilot/copilot-for-business#list-all-copilot-for-business-seat-assignments-for-an-organization
func (s *CopilotService) ListCopilotSeats(ctx context.Context, org string) (*CopilotSeats, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing/seats", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var copilotSeats *CopilotSeats
	resp, err := s.client.Do(ctx, req, &copilotSeats)
	if err != nil {
		return nil, resp, err
	}

	return copilotSeats, resp, nil
}

// Add teams to the Copilot for Business subscription for an organization
//
// https://docs.github.com/en/rest/copilot/copilot-for-business#add-teams-to-the-copilot-for-business-subscription-for-an-organization
func (s *CopilotService) AddCopilotTeams(ctx context.Context, org string, teams SelectedTeams) (*SeatAssignments, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing/selected_teams", org)

	req, err := s.client.NewRequest("PUT", u, teams)
	if err != nil {
		return nil, nil, err
	}

	var seatAssignments *SeatAssignments
	resp, err := s.client.Do(ctx, req, &seatAssignments)
	if err != nil {
		return nil, resp, err
	}

	return seatAssignments, resp, nil
}

// Remove teams from the Copilot for Business subscription for an organization
//
// https://docs.github.com/en/rest/copilot/copilot-for-business#remove-teams-from-the-copilot-for-business-subscription-for-an-organization

func (s *CopilotService) RemoveCopilotTeams(ctx context.Context, org string, teams SelectedTeams) (*SeatCancellations, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing/selected_teams", org)

	req, err := s.client.NewRequest("DELETE", u, teams)
	if err != nil {
		return nil, nil, err
	}

	var SeatCancellations *SeatCancellations
	resp, err := s.client.Do(ctx, req, &SeatCancellations)
	if err != nil {
		return nil, resp, err
	}

	return SeatCancellations, resp, nil
}

// Add users to the Copilot for Business subscription for an organization
//
// https://docs.github.com/en/rest/copilot/copilot-for-business#add-users-to-the-copilot-for-business-subscription-for-an-organization
func (s *CopilotService) AddCopilotUsers(ctx context.Context, org string, users SelectedUsers) (*SeatAssignments, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing/selected_users", org)

	req, err := s.client.NewRequest("PUT", u, users)
	if err != nil {
		return nil, nil, err
	}

	var seatAssignments *SeatAssignments
	resp, err := s.client.Do(ctx, req, &seatAssignments)
	if err != nil {
		return nil, resp, err
	}

	return seatAssignments, resp, nil
}

// Remove users from the Copilot for Business subscription for an organization
//
// https://docs.github.com/en/rest/copilot/copilot-for-business#remove-users-from-the-copilot-for-business-subscription-for-an-organization
func (s *CopilotService) RemoveCopilotUsers(ctx context.Context, org string, users SelectedUsers) (*SeatCancellations, *Response, error) {
	u := fmt.Sprintf("orgs/%v/copilot/billing/selected_users", org)

	req, err := s.client.NewRequest("DELETE", u, users)
	if err != nil {
		return nil, nil, err
	}

	var SeatCancellations *SeatCancellations
	resp, err := s.client.Do(ctx, req, &SeatCancellations)
	if err != nil {
		return nil, resp, err
	}

	return SeatCancellations, resp, nil
}

// Get Copilot for Business seat assignment details for a user
//
// https://docs.github.com/en/rest/copilot/copilot-for-business#get-copilot-for-business-seat-assignment-details-for-a-user
func (s *CopilotService) GetSeatDetails(ctx context.Context, org string, user string) (*CopilotSeatDetails, *Response, error) {
	u := fmt.Sprintf("orgs/%v/members/%v/copilot", org, user)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var seatDetails *CopilotSeatDetails
	resp, err := s.client.Do(ctx, req, &seatDetails)
	if err != nil {
		return nil, resp, err
	}

	return seatDetails, resp, nil
}
