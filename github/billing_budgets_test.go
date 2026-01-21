// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillingService_ListOrganizationBudgets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"id": "1",
				"budget_name": "Budget 1",
				"limit_amount": 100.5,
				"budget_alerting": {
					"will_alert": true,
					"alert_recipients": ["user1"]
				}
			}
		]`)
	})

	ctx := t.Context()
	budgets, _, err := client.Billing.ListOrganizationBudgets(ctx, "o")
	if err != nil {
		t.Errorf("Billing.ListOrganizationBudgets returned error: %v", err)
	}

	want := []*Budget{
		{
			ID:          Ptr("1"),
			BudgetName:  Ptr("Budget 1"),
			LimitAmount: Ptr(100.5),
			BudgetAlerting: &BudgetAlerting{
				WillAlert:       Ptr(true),
				AlertRecipients: []string{"user1"},
			},
		},
	}
	if !cmp.Equal(budgets, want) {
		t.Errorf("Billing.ListOrganizationBudgets returned %+v, want %+v", budgets, want)
	}
}

func TestBillingService_GetOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": "1",
			"budget_name": "Budget 1"
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Billing.GetOrganizationBudget(ctx, "o", "1")
	if err != nil {
		t.Errorf("Billing.GetOrganizationBudget returned error: %v", err)
	}

	want := &Budget{
		ID:         Ptr("1"),
		BudgetName: Ptr("Budget 1"),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Billing.GetOrganizationBudget returned %+v, want %+v", budget, want)
	}
}

func TestBillingService_UpdateOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Budget{
		BudgetName: Ptr("Updated Budget"),
	}

	mux.HandleFunc("/organizations/o/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"budget_name":"Updated Budget"}`+"\n")
		fmt.Fprint(w, `{
			"id": "1",
			"budget_name": "Updated Budget"
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Billing.UpdateOrganizationBudget(ctx, "o", "1", input)
	if err != nil {
		t.Errorf("Billing.UpdateOrganizationBudget returned error: %v", err)
	}

	want := &Budget{
		ID:         Ptr("1"),
		BudgetName: Ptr("Updated Budget"),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Billing.UpdateOrganizationBudget returned %+v, want %+v", budget, want)
	}
}

func TestBillingService_DeleteOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Billing.DeleteOrganizationBudget(ctx, "o", "1")
	if err != nil {
		t.Errorf("Billing.DeleteOrganizationBudget returned error: %v", err)
	}
}
