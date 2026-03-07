// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// EnterpriseBudgetAlerting represents alerting settings for a GitHub enterprise budget.
type EnterpriseBudgetAlerting struct {
	WillAlert       *bool    `json:"will_alert,omitempty"`
	AlertRecipients []string `json:"alert_recipients,omitempty"`
}

// EnterpriseBudget represents a GitHub enterprise budget.
type EnterpriseBudget struct {
	ID                  *string                   `json:"id,omitempty"`
	BudgetType          *string                   `json:"budget_type,omitempty"`
	BudgetProductSkus   []string                  `json:"budget_product_skus,omitempty"`
	BudgetProductSKU    *string                   `json:"budget_product_sku,omitempty"`
	BudgetScope         *string                   `json:"budget_scope,omitempty"`
	BudgetEntityName    *string                   `json:"budget_entity_name,omitempty"`
	BudgetAmount        *int                      `json:"budget_amount,omitempty"`
	PreventFurtherUsage *bool                     `json:"prevent_further_usage,omitempty"`
	BudgetAlerting      *EnterpriseBudgetAlerting `json:"budget_alerting,omitempty"`
}

func (b EnterpriseBudget) String() string {
	return Stringify(b)
}

// EnterpriseBudgets represents a collection of GitHub enterprise budgets.
type EnterpriseBudgets struct {
	Budgets []*EnterpriseBudget `json:"budgets,omitempty"`
}

// EnterpriseBudgetActionResponse represents the response when creating, updating, or deleting a budget.
type EnterpriseBudgetActionResponse struct {
	Message  *string           `json:"message,omitempty"`
	BudgetID *string           `json:"budget_id,omitempty"`
	Budget   *EnterpriseBudget `json:"budget,omitempty"`
}

// ListBudgets gets all budgets for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#get-all-budgets
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/budgets
func (s *EnterpriseService) ListBudgets(ctx context.Context, enterprise string) (*EnterpriseBudgets, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var budgets *EnterpriseBudgets
	resp, err := s.client.Do(ctx, req, &budgets)
	if err != nil {
		return nil, resp, err
	}

	return budgets, resp, nil
}

// CreateBudget creates a new budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#create-a-budget
//
//meta:operation POST /enterprises/{enterprise}/settings/billing/budgets
func (s *EnterpriseService) CreateBudget(ctx context.Context, enterprise string, budget EnterpriseBudget) (*EnterpriseBudgetActionResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets", enterprise)

	req, err := s.client.NewRequest("POST", u, budget)
	if err != nil {
		return nil, nil, err
	}

	var actionResponse *EnterpriseBudgetActionResponse
	resp, err := s.client.Do(ctx, req, &actionResponse)
	if err != nil {
		return nil, resp, err
	}

	return actionResponse, resp, nil
}

// GetBudget gets a budget by ID for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#get-a-budget-by-id
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) GetBudget(ctx context.Context, enterprise, budgetID string) (*EnterpriseBudget, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var budget *EnterpriseBudget
	resp, err := s.client.Do(ctx, req, &budget)
	if err != nil {
		return nil, resp, err
	}

	return budget, resp, nil
}

// UpdateBudget updates an existing budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#update-a-budget
//
//meta:operation PATCH /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) UpdateBudget(ctx context.Context, enterprise, budgetID string, budget *EnterpriseBudget) (*EnterpriseBudgetActionResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)

	req, err := s.client.NewRequest("PATCH", u, budget)
	if err != nil {
		return nil, nil, err
	}

	var actionResponse *EnterpriseBudgetActionResponse
	resp, err := s.client.Do(ctx, req, &actionResponse)
	if err != nil {
		return nil, resp, err
	}

	return actionResponse, resp, nil
}

// DeleteBudget deletes a budget by ID for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#delete-a-budget
//
//meta:operation DELETE /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) DeleteBudget(ctx context.Context, enterprise, budgetID string) (*EnterpriseBudgetActionResponse, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var actionResponse *EnterpriseBudgetActionResponse
	resp, err := s.client.Do(ctx, req, &actionResponse)
	if err != nil {
		return nil, resp, err
	}

	return actionResponse, resp, nil
}
