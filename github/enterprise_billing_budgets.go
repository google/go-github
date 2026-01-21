// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListEnterpriseBudgets lists all budgets for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#get-all-budgets
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/budgets
func (s *EnterpriseService) ListEnterpriseBudgets(ctx context.Context, enterprise string) ([]*Budget, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets", enterprise)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var budgets []*Budget
	resp, err := s.client.Do(ctx, req, &budgets)
	if err != nil {
		return nil, resp, err
	}

	return budgets, resp, nil
}

// GetEnterpriseBudget gets a specific budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#get-a-budget-by-id
//
//meta:operation GET /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) GetEnterpriseBudget(ctx context.Context, enterprise, budgetID string) (*Budget, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)
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

// CreateEnterpriseBudget creates a specific budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#create-a-budget
//
//meta:operation POST /enterprises/{enterprise}/settings/billing/budgets
func (s *EnterpriseService) CreateEnterpriseBudget(ctx context.Context, enterprise string, budget *Budget) (*Budget, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets", enterprise)
	req, err := s.client.NewRequest("POST", u, budget)
	if err != nil {
		return nil, nil, err
	}

	createdBudget := new(Budget)
	resp, err := s.client.Do(ctx, req, createdBudget)
	if err != nil {
		return nil, resp, err
	}

	return createdBudget, resp, nil
}

// UpdateEnterpriseBudget updates a specific budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#update-a-budget
//
//meta:operation PATCH /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) UpdateEnterpriseBudget(ctx context.Context, enterprise, budgetID string, budget *Budget) (*Budget, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)
	req, err := s.client.NewRequest("PATCH", u, budget)
	if err != nil {
		return nil, nil, err
	}

	updatedBudget := new(Budget)
	resp, err := s.client.Do(ctx, req, updatedBudget)
	if err != nil {
		return nil, resp, err
	}

	return updatedBudget, resp, nil
}

// DeleteEnterpriseBudget deletes a specific budget for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/budgets#delete-a-budget
//
//meta:operation DELETE /enterprises/{enterprise}/settings/billing/budgets/{budget_id}
func (s *EnterpriseService) DeleteEnterpriseBudget(ctx context.Context, enterprise, budgetID string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/settings/billing/budgets/%v", enterprise, budgetID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
