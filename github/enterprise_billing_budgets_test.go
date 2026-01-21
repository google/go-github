// Copyright 2026 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_ListEnterpriseBudgets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"id": "1",
				"budget_name": "Budget 1"
			}
		]`)
	})

	ctx := t.Context()
	budgets, _, err := client.Enterprise.ListEnterpriseBudgets(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.ListEnterpriseBudgets returned error: %v", err)
	}

	want := []*Budget{
		{
			ID:         Ptr("1"),
			BudgetName: Ptr("Budget 1"),
		},
	}
	if !cmp.Equal(budgets, want) {
		t.Errorf("Enterprise.ListEnterpriseBudgets returned %+v, want %+v", budgets, want)
	}
}

func TestEnterpriseService_GetEnterpriseBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": "1",
			"budget_name": "Budget 1"
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Enterprise.GetEnterpriseBudget(ctx, "e", "1")
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseBudget returned error: %v", err)
	}

	want := &Budget{
		ID:         Ptr("1"),
		BudgetName: Ptr("Budget 1"),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.GetEnterpriseBudget returned %+v, want %+v", budget, want)
	}
}

func TestEnterpriseService_CreateEnterpriseBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Budget{
		BudgetName:  Ptr("New Budget"),
		LimitAmount: Ptr(500.0),
	}

	mux.HandleFunc("/enterprises/e/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"budget_name":"New Budget","limit_amount":500}`+"\n")
		fmt.Fprint(w, `{
			"id": "1",
			"budget_name": "New Budget",
			"limit_amount": 500
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Enterprise.CreateEnterpriseBudget(ctx, "e", input)
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseBudget returned error: %v", err)
	}

	want := &Budget{
		ID:          Ptr("1"),
		BudgetName:  Ptr("New Budget"),
		LimitAmount: Ptr(500.0),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.CreateEnterpriseBudget returned %+v, want %+v", budget, want)
	}
}

func TestEnterpriseService_UpdateEnterpriseBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Budget{
		BudgetName: Ptr("Updated Budget"),
	}

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"budget_name":"Updated Budget"}`+"\n")
		fmt.Fprint(w, `{
			"id": "1",
			"budget_name": "Updated Budget"
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Enterprise.UpdateEnterpriseBudget(ctx, "e", "1", input)
	if err != nil {
		t.Errorf("Enterprise.UpdateEnterpriseBudget returned error: %v", err)
	}

	want := &Budget{
		ID:         Ptr("1"),
		BudgetName: Ptr("Updated Budget"),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.UpdateEnterpriseBudget returned %+v, want %+v", budget, want)
	}
}

func TestEnterpriseService_DeleteEnterpriseBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Enterprise.DeleteEnterpriseBudget(ctx, "e", "1")
	if err != nil {
		t.Errorf("Enterprise.DeleteEnterpriseBudget returned error: %v", err)
	}
}

func TestBudget_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Budget{}, "{}")

	b := &Budget{
		ID:          Ptr("1"),
		BudgetName:  Ptr("Budget"),
		LimitAmount: Ptr(100.0),
		BudgetAlerting: &BudgetAlerting{
			WillAlert:       Ptr(true),
			AlertRecipients: []string{"u1"},
		},
	}

	want := `{
		"id": "1",
		"budget_name": "Budget",
		"limit_amount": 100,
		"budget_alerting": {
			"will_alert": true,
			"alert_recipients": ["u1"]
		}
	}`

	testJSONMarshal(t, b, want)
}
