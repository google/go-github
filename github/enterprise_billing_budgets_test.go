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

func TestEnterpriseService_ListBudgets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"budgets": [
				{
					"id": "1",
					"budget_name": "Budget 1"
				}
			]
		}`)
	})

	ctx := t.Context()
	budgets, _, err := client.Enterprise.ListBudgets(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.ListBudgets returned error: %v", err)
	}

	want := &BudgetList{
		Budgets: []*Budget{
			{
				ID:         Ptr("1"),
				BudgetName: Ptr("Budget 1"),
			},
		},
	}
	if !cmp.Equal(budgets, want) {
		t.Errorf("Enterprise.ListBudgets returned %+v, want %+v", budgets, want)
	}
	const methodName = "ListBudgets"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListBudgets(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetBudget(t *testing.T) {
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
	budget, _, err := client.Enterprise.GetBudget(ctx, "e", "1")
	if err != nil {
		t.Errorf("Enterprise.GetBudget returned error: %v", err)
	}

	want := &Budget{
		ID:         Ptr("1"),
		BudgetName: Ptr("Budget 1"),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.GetBudget returned %+v, want %+v", budget, want)
	}
	const methodName = "GetBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetBudget(ctx, "e", "1")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateBudget(t *testing.T) {
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
	budget, _, err := client.Enterprise.CreateBudget(ctx, "e", input)
	if err != nil {
		t.Errorf("Enterprise.CreateBudget returned error: %v", err)
	}

	want := &Budget{
		ID:          Ptr("1"),
		BudgetName:  Ptr("New Budget"),
		LimitAmount: Ptr(500.0),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.CreateBudget returned %+v, want %+v", budget, want)
	}
	const methodName = "CreateBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateBudget(ctx, "e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Budget{
		BudgetName: Ptr("Updated Budget"),
	}

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"budget_name":"Updated Budget"}`+"\n")
		fmt.Fprint(w, `{
			"budget": {
				"id": "1",
				"budget_name": "Updated Budget"
			}
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Enterprise.UpdateBudget(ctx, "e", "1", input)
	if err != nil {
		t.Errorf("Enterprise.UpdateBudget returned error: %v", err)
	}

	want := &BudgetResponse{
		Budget: &Budget{
			ID:         Ptr("1"),
			BudgetName: Ptr("Updated Budget"),
		},
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.UpdateBudget returned %+v, want %+v", budget, want)
	}
	const methodName = "UpdateBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateBudget(ctx, "e", "1", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Enterprise.DeleteBudget(ctx, "e", "1")
	if err != nil {
		t.Errorf("Enterprise.DeleteBudget returned error: %v", err)
	}
	const methodName = "DeleteBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteBudget(ctx, "e", "1")
	})
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
