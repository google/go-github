// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// Budget represents a GitHub budget.
type Budget struct {
	ID                      *string         `json:"id,omitempty"`
	BudgetName              *string         `json:"budget_name,omitempty"`
	TargetSubAccount        *string         `json:"target_sub_account,omitempty"`
	TargetType              *string         `json:"target_type,omitempty"`
	TargetID                *int64          `json:"target_id,omitempty"`
	TargetName              *string         `json:"target_name,omitempty"`
	PricingModel            *string         `json:"pricing_model,omitempty"`
	PricingModelID          *string         `json:"pricing_model_id,omitempty"`
	PricingModelDisplayName *string         `json:"pricing_model_display_name,omitempty"`
	BudgetType              *string         `json:"budget_type,omitempty"`
	LimitAmount             *float64        `json:"limit_amount,omitempty"`
	CurrentAmount           *float64        `json:"current_amount,omitempty"`
	Currency                *string         `json:"currency,omitempty"`
	ExcludeCostCenterUsage  *bool           `json:"exclude_cost_center_usage,omitempty"`
	BudgetAlerting          *BudgetAlerting `json:"budget_alerting,omitempty"`
}

// BudgetAlerting represents the alerting configuration for a budget.
type BudgetAlerting struct {
	WillAlert       *bool    `json:"will_alert,omitempty"`
	AlertRecipients []string `json:"alert_recipients,omitempty"`
}

// BudgetList represents a list of budgets.
type BudgetList struct {
	Budgets     []*Budget `json:"budgets"`
	HasNextPage *bool     `json:"has_next_page,omitempty"`
}

// BudgetResponse represents the response when updating a budget.
type BudgetResponse struct {
	Budget  *Budget `json:"budget"`
	Message *string `json:"message,omitempty"`
}

// ListOrganizationBudgets lists all budgets for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets#get-all-budgets-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/budgets
func (s *BillingService) ListOrganizationBudgets(ctx context.Context, org string) (*BudgetList, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	budgets := new(BudgetList)
	resp, err := s.client.Do(ctx, req, budgets)
	if err != nil {
		return nil, resp, err
	}

	return budgets, resp, nil
}

// GetOrganizationBudget gets a specific budget for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets#get-a-budget-by-id-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) GetOrganizationBudget(ctx context.Context, org, budgetID string) (*Budget, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	budget := new(Budget)
	resp, err := s.client.Do(ctx, req, budget)
	if err != nil {
		return nil, resp, err
	}

	return budget, resp, nil
}

// UpdateOrganizationBudget updates a specific budget for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets#update-a-budget-for-an-organization
//
//meta:operation PATCH /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) UpdateOrganizationBudget(ctx context.Context, org, budgetID string, budget *Budget) (*BudgetResponse, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)
	req, err := s.client.NewRequest("PATCH", u, budget)
	if err != nil {
		return nil, nil, err
	}

	updatedBudget := new(BudgetResponse)
	resp, err := s.client.Do(ctx, req, updatedBudget)
	if err != nil {
		return nil, resp, err
	}

	return updatedBudget, resp, nil
}

// DeleteOrganizationBudget deletes a specific budget for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets#delete-a-budget-for-an-organization
//
//meta:operation DELETE /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) DeleteOrganizationBudget(ctx context.Context, org, budgetID string) (*Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
