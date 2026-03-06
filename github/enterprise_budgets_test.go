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
					"id": "2066deda-923f-43f9-88d2-62395a28c0cdd",
					"budget_type": "ProductPricing",
					"budget_product_skus": ["actions"],
					"budget_scope": "enterprise",
					"budget_amount": 1000,
					"prevent_further_usage": true,
					"budget_alerting": {
						"will_alert": true,
						"alert_recipients": ["enterprise-admin"]
					}
				}
			]
		}`)
	})

	ctx := t.Context()
	budgets, _, err := client.Enterprise.ListBudgets(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.ListBudgets returned error: %v", err)
	}

	want := &EnterpriseBudgets{
		Budgets: []*Budget{
			{
				ID:                  Ptr("2066deda-923f-43f9-88d2-62395a28c0cdd"),
				BudgetType:          Ptr("ProductPricing"),
				BudgetProductSkus:   []string{"actions"},
				BudgetScope:         Ptr("enterprise"),
				BudgetAmount:        Ptr(1000),
				PreventFurtherUsage: Ptr(true),
				BudgetAlerting: &BudgetAlerting{
					WillAlert:       Ptr(true),
					AlertRecipients: []string{"enterprise-admin"},
				},
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

func TestEnterpriseService_ListBudgets_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.ListBudgets(ctx, "%")
	testURLParseError(t, err)
}

func TestEnterpriseService_CreateBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"budget_type":"ProductPricing","budget_product_sku":"actions","budget_scope":"enterprise","budget_amount":200,"prevent_further_usage":true}`+"\n")
		fmt.Fprint(w, `{
			"message": "Budget successfully created."
		}`)
	})

	ctx := t.Context()
	req := &Budget{
		BudgetAmount:        Ptr(200),
		PreventFurtherUsage: Ptr(true),
		BudgetScope:         Ptr("enterprise"),
		BudgetType:          Ptr("ProductPricing"),
		BudgetProductSKU:    Ptr("actions"),
	}

	resp, _, err := client.Enterprise.CreateBudget(ctx, "e", req)
	if err != nil {
		t.Errorf("Enterprise.CreateBudget returned error: %v", err)
	}

	want := &BudgetActionResponse{
		Message: Ptr("Budget successfully created."),
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Enterprise.CreateBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "CreateBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateBudget(ctx, "e", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateBudget_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.CreateBudget(ctx, "%", &Budget{})
	testURLParseError(t, err)
}

func TestEnterpriseService_GetBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/2066deda-923f-43f9-88d2-62395a28c0cdd", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": "2066deda-923f-43f9-88d2-62395a28c0cdd",
			"budget_type": "ProductPricing",
			"budget_product_sku": "actions_linux",
			"budget_scope": "repository",
			"budget_amount": 0,
			"prevent_further_usage": true
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Enterprise.GetBudget(ctx, "e", "2066deda-923f-43f9-88d2-62395a28c0cdd")
	if err != nil {
		t.Errorf("Enterprise.GetBudget returned error: %v", err)
	}

	want := &Budget{
		ID:                  Ptr("2066deda-923f-43f9-88d2-62395a28c0cdd"),
		BudgetType:          Ptr("ProductPricing"),
		BudgetProductSKU:    Ptr("actions_linux"),
		BudgetScope:         Ptr("repository"),
		BudgetAmount:        Ptr(0),
		PreventFurtherUsage: Ptr(true),
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Enterprise.GetBudget returned %+v, want %+v", budget, want)
	}

	const methodName = "GetBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetBudget(ctx, "e", "2066deda-923f-43f9-88d2-62395a28c0cdd")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetBudget_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.GetBudget(ctx, "%", "b-123")
	testURLParseError(t, err)
}

func TestEnterpriseService_UpdateBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/b-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"budget_amount":10,"prevent_further_usage":false}`+"\n")
		fmt.Fprint(w, `{
			"message": "Budget successfully updated.",
			"budget": {
				"id": "b-123",
				"budget_amount": 10,
				"prevent_further_usage": false
			}
		}`)
	})

	ctx := t.Context()
	req := &Budget{
		BudgetAmount:        Ptr(10),
		PreventFurtherUsage: Ptr(false),
	}

	resp, _, err := client.Enterprise.UpdateBudget(ctx, "e", "b-123", req)
	if err != nil {
		t.Errorf("Enterprise.UpdateBudget returned error: %v", err)
	}

	want := &BudgetActionResponse{
		Message: Ptr("Budget successfully updated."),
		Budget: &Budget{
			ID:                  Ptr("b-123"),
			BudgetAmount:        Ptr(10),
			PreventFurtherUsage: Ptr(false),
		},
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Enterprise.UpdateBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "UpdateBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateBudget(ctx, "e", "b-123", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateBudget_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.UpdateBudget(ctx, "%", "b-123", &Budget{})
	testURLParseError(t, err)
}

func TestEnterpriseService_DeleteBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/settings/billing/budgets/b-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message": "Budget successfully deleted.",
			"budget_id": "b-123"
		}`)
	})

	ctx := t.Context()
	resp, _, err := client.Enterprise.DeleteBudget(ctx, "e", "b-123")
	if err != nil {
		t.Errorf("Enterprise.DeleteBudget returned error: %v", err)
	}

	want := &BudgetActionResponse{
		Message:  Ptr("Budget successfully deleted."),
		BudgetID: Ptr("b-123"),
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Enterprise.DeleteBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "DeleteBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.DeleteBudget(ctx, "e", "b-123")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteBudget_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.DeleteBudget(ctx, "%", "b-123")
	testURLParseError(t, err)
}

func TestBudget_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Budget{}, "{}")

	u := &Budget{
		ID:                  Ptr("b-123"),
		BudgetType:          Ptr("ProductPricing"),
		BudgetProductSkus:   []string{"actions"},
		BudgetProductSKU:    Ptr("actions"),
		BudgetScope:         Ptr("enterprise"),
		BudgetEntityName:    Ptr("org-name"),
		BudgetAmount:        Ptr(100),
		PreventFurtherUsage: Ptr(true),
		BudgetAlerting: &BudgetAlerting{
			WillAlert:       Ptr(true),
			AlertRecipients: []string{"mona"},
		},
	}

	want := `{
		"id": "b-123",
		"budget_type": "ProductPricing",
		"budget_product_skus": ["actions"],
		"budget_product_sku": "actions",
		"budget_scope": "enterprise",
		"budget_entity_name": "org-name",
		"budget_amount": 100,
		"prevent_further_usage": true,
		"budget_alerting": {
			"will_alert": true,
			"alert_recipients": ["mona"]
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestBudgetAlerting_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &BudgetAlerting{}, "{}")

	u := &BudgetAlerting{
		WillAlert:       Ptr(true),
		AlertRecipients: []string{"admin1", "admin2"},
	}

	want := `{
		"will_alert": true,
		"alert_recipients": ["admin1", "admin2"]
	}`

	testJSONMarshal(t, u, want)
}

func TestEnterpriseBudgets_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EnterpriseBudgets{}, "{}")

	u := &EnterpriseBudgets{
		Budgets: []*Budget{
			{
				ID: Ptr("1"),
			},
		},
	}

	want := `{
		"budgets": [
			{
				"id": "1"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestBudgetActionResponse_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &BudgetActionResponse{}, "{}")

	u := &BudgetActionResponse{
		Message:  Ptr("Success"),
		BudgetID: Ptr("123"),
		Budget: &Budget{
			ID: Ptr("123"),
		},
	}

	want := `{
		"message": "Success",
		"budget_id": "123",
		"budget": {
			"id": "123"
		}
	}`

	testJSONMarshal(t, u, want)
}
