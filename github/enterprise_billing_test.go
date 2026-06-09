// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseService_GetUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2023",
			"month": "8",
		})
		fmt.Fprint(w, `{
			"usageItems": [
				{
					"date": "2023-08-01",
					"product": "Actions",
					"sku": "Actions Linux",
					"quantity": 100,
					"unitType": "minutes",
					"pricePerUnit": 0.008,
					"grossAmount": 0.8,
					"discountAmount": 0,
					"netAmount": 0.8,
					"organizationName": "GitHub",
					"repositoryName": "github/example"
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageReportOptions{
		Year:  2023,
		Month: 8,
	}
	report, resp, err := client.Enterprise.GetUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageReport returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageReport returned nil response")
	}

	want := &EnterpriseUsageReport{
		UsageItems: []*EnterpriseUsageItem{
			{
				Date:             "2023-08-01",
				Product:          "Actions",
				SKU:              "Actions Linux",
				Quantity:         100.0,
				UnitType:         "minutes",
				PricePerUnit:     0.008,
				GrossAmount:      0.8,
				DiscountAmount:   0.0,
				NetAmount:        0.8,
				OrganizationName: "GitHub",
				RepositoryName:   Ptr("github/example"),
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("GetUsageReport returned %+v, want %+v", report, want)
	}
}

func TestEnterpriseService_GetUsageReport_WithCostCenter(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":           "2023",
			"cost_center_id": "cc-123",
		})
		fmt.Fprint(w, `{
			"usageItems": [
				{
					"date": "2023-08-01",
					"product": "Actions",
					"sku": "Actions Linux",
					"quantity": 250,
					"unitType": "minutes",
					"pricePerUnit": 0.008,
					"grossAmount": 2.0,
					"discountAmount": 0.2,
					"netAmount": 1.8,
					"organizationName": "GitHub"
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageReportOptions{
		Year:         2023,
		CostCenterID: "cc-123",
	}
	report, resp, err := client.Enterprise.GetUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageReport returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageReport returned nil response")
	}

	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
	if report.UsageItems[0].NetAmount != 1.8 {
		t.Errorf("Expected NetAmount 1.8, got %v", report.UsageItems[0].NetAmount)
	}
}

func TestEnterpriseService_GetUsageReport_InvalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.GetUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestEnterpriseService_GetUsageReport_NoOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"usageItems": []}`)
	})

	ctx := t.Context()
	report, resp, err := client.Enterprise.GetUsageReport(ctx, "test-enterprise", nil)
	if err != nil {
		t.Errorf("GetUsageReport returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageReport returned nil response")
	}

	if report == nil || len(report.UsageItems) != 0 {
		t.Error("Expected empty usage items")
	}
}

func TestEnterpriseService_GetUsageReport_APIError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message": "Forbidden"}`)
	})

	ctx := t.Context()
	_, resp, err := client.Enterprise.GetUsageReport(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetUsageReport should return error on 403 response")
	}
	if resp == nil {
		t.Error("GetUsageReport should return response even on error")
	}
}

func TestEnterpriseService_GetUsageReport_CanceledContext(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, resp, err := client.Enterprise.GetUsageReport(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetUsageReport should return error with canceled context")
	}
	_ = resp
}

func TestEnterpriseService_GetUsageSummary(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":    "2025",
			"product": "Actions",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"usageItems": [
				{
					"product": "Actions",
					"sku": "actions_linux",
					"unitType": "minutes",
					"pricePerUnit": 0.008,
					"grossQuantity": 1000,
					"grossAmount": 8,
					"discountQuantity": 0,
					"discountAmount": 0,
					"netQuantity": 1000,
					"netAmount": 8
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageSummaryOptions{
		Year:    2025,
		Product: "Actions",
	}
	report, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageSummary returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageSummary returned nil response")
	}

	want := &EnterpriseAggregatedUsageReport{
		TimePeriod: EnterpriseUsageTimePeriod{
			Year: 2025,
		},
		Enterprise: "GitHub",
		UsageItems: []*EnterpriseAggregatedUsageItem{
			{
				Product:          "Actions",
				SKU:              "actions_linux",
				UnitType:         "minutes",
				PricePerUnit:     0.008,
				GrossQuantity:    1000.0,
				GrossAmount:      8.0,
				DiscountQuantity: 0.0,
				DiscountAmount:   0.0,
				NetQuantity:      1000.0,
				NetAmount:        8.0,
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("GetUsageSummary returned %+v, want %+v", report, want)
	}
}

func TestEnterpriseService_GetUsageSummary_MultipleItems(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year": "2025",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"usageItems": [
				{
					"product": "Actions",
					"sku": "actions_linux",
					"unitType": "minutes",
					"pricePerUnit": 0.008,
					"grossQuantity": 1000,
					"grossAmount": 8,
					"discountQuantity": 100,
					"discountAmount": 0.8,
					"netQuantity": 900,
					"netAmount": 7.2
				},
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"unitType": "credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 500,
					"grossAmount": 5,
					"discountQuantity": 0,
					"discountAmount": 0,
					"netQuantity": 500,
					"netAmount": 5
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageSummaryOptions{
		Year: 2025,
	}
	report, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageSummary returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageSummary returned nil response")
	}

	if len(report.UsageItems) != 2 {
		t.Errorf("Expected 2 usage items, got %v", len(report.UsageItems))
	}
	if report.UsageItems[0].Product != "Actions" {
		t.Errorf("Expected first product to be Actions, got %v", report.UsageItems[0].Product)
	}
	if report.UsageItems[1].Product != "Copilot" {
		t.Errorf("Expected second product to be Copilot, got %v", report.UsageItems[1].Product)
	}
}

func TestEnterpriseService_GetUsageSummary_WithRepository(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"repository": "acme/api",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"repository": "acme/api",
			"usageItems": [
				{
					"product": "Actions",
					"sku": "actions_linux",
					"unitType": "minutes",
					"pricePerUnit": 0.008,
					"grossQuantity": 500,
					"grossAmount": 4.0,
					"discountQuantity": 0,
					"discountAmount": 0,
					"netQuantity": 500,
					"netAmount": 4.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageSummaryOptions{
		Repository: "acme/api",
	}
	report, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageSummary returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageSummary returned nil response")
	}

	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
}

func TestEnterpriseService_GetUsageSummary_WithAllFilters(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":           "2025",
			"month":          "6",
			"day":            "15",
			"organization":   "acme",
			"repository":     "acme/api",
			"product":        "Actions",
			"sku":            "actions_linux",
			"cost_center_id": "cc-001",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 6,
				"day": 15
			},
			"enterprise": "GitHub",
			"organization": "acme",
			"repository": "acme/api",
			"product": "Actions",
			"usageItems": []
		}`)
	})

	ctx := t.Context()
	opts := &EnterpriseUsageSummaryOptions{
		Year:         2025,
		Month:        6,
		Day:          15,
		Organization: "acme",
		Repository:   "acme/api",
		Product:      "Actions",
		SKU:          "actions_linux",
		CostCenterID: "cc-001",
	}
	report, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetUsageSummary returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetUsageSummary returned nil response")
	}

	if report == nil {
		t.Error("Expected non-nil report")
	}
}

func TestEnterpriseService_GetUsageSummary_InvalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.GetUsageSummary(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestEnterpriseService_GetUsageSummary_APIError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/usage/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message": "Not Found"}`)
	})

	ctx := t.Context()
	_, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetUsageSummary should return error on 404 response")
	}
	if resp == nil {
		t.Error("GetUsageSummary should return response even on error")
	}
}

func TestEnterpriseService_GetUsageSummary_CanceledContext(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, resp, err := client.Enterprise.GetUsageSummary(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetUsageSummary should return error with canceled context")
	}
	_ = resp
}

func TestEnterpriseService_GetPremiumRequestUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/premium_request/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2025",
			"month": "10",
			"user":  "testuser",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 10
			},
			"enterprise": "GitHub",
			"organization": "GitHub",
			"user": "testuser",
			"product": "Copilot",
			"model": "GPT-5",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot Premium Request",
					"model": "GPT-5",
					"unitType": "requests",
					"pricePerUnit": 0.04,
					"grossQuantity": 100,
					"grossAmount": 4.0,
					"discountQuantity": 0,
					"discountAmount": 0.0,
					"netQuantity": 100,
					"netAmount": 4.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Year:  2025,
		Month: 10,
		User:  "testuser",
	}
	report, resp, err := client.Enterprise.GetPremiumRequestUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetPremiumRequestUsageReport returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetPremiumRequestUsageReport returned nil response")
	}

	if report.Enterprise != "GitHub" {
		t.Errorf("Expected enterprise GitHub, got %v", report.Enterprise)
	}
	if report.User == nil || *report.User != "testuser" {
		t.Errorf("Expected user testuser, got %v", report.User)
	}
	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
}

func TestEnterpriseService_GetPremiumRequestUsageReport_WithAllOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/premium_request/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":           "2025",
			"month":          "5",
			"day":            "10",
			"organization":   "acme-org",
			"user":           "alice",
			"model":          "GPT-4",
			"product":        "Copilot",
			"cost_center_id": "cc-sales",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 5,
				"day": 10
			},
			"enterprise": "GitHub",
			"organization": "acme-org",
			"user": "alice",
			"product": "Copilot",
			"model": "GPT-4",
			"usageItems": []
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Year:         2025,
		Month:        5,
		Day:          10,
		Organization: "acme-org",
		User:         "alice",
		Model:        "GPT-4",
		Product:      "Copilot",
		CostCenterID: "cc-sales",
	}
	report, resp, err := client.Enterprise.GetPremiumRequestUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetPremiumRequestUsageReport returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetPremiumRequestUsageReport returned nil response")
	}

	if report == nil {
		t.Error("Expected non-nil report")
	}
}

func TestEnterpriseService_GetPremiumRequestUsageReport_InvalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.GetPremiumRequestUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestEnterpriseService_GetPremiumRequestUsageReport_APIError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/premium_request/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"message": "Unauthorized"}`)
	})

	ctx := t.Context()
	_, resp, err := client.Enterprise.GetPremiumRequestUsageReport(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetPremiumRequestUsageReport should return error on 401 response")
	}
	if resp == nil {
		t.Error("GetPremiumRequestUsageReport should return response even on error")
	}
}

func TestEnterpriseService_GetPremiumRequestUsageReport_CanceledContext(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, resp, err := client.Enterprise.GetPremiumRequestUsageReport(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetPremiumRequestUsageReport should return error with canceled context")
	}
	_ = resp
}

func TestEnterpriseService_GetAICreditUsage(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2025",
			"month": "6",
			"user":  "testuser",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 6
			},
			"enterprise": "GitHub",
			"organization": "GitHub",
			"user": "testuser",
			"product": "Copilot",
			"model": "GPT-5",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "GPT-5",
					"unitType": "credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 100,
					"grossAmount": 1.0,
					"discountQuantity": 0,
					"discountAmount": 0.0,
					"netQuantity": 100,
					"netAmount": 1.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Year:  2025,
		Month: 6,
		User:  "testuser",
	}
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if report.Enterprise != "GitHub" {
		t.Errorf("Expected enterprise GitHub, got %v", report.Enterprise)
	}
	if report.User == nil || *report.User != "testuser" {
		t.Errorf("Expected user testuser, got %v", report.User)
	}
	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
	if report.UsageItems[0].UnitType != "credits" {
		t.Errorf("Expected unitType credits, got %v", report.UsageItems[0].UnitType)
	}
}

func TestEnterpriseService_GetAICreditUsage_FloatQuantities(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2025",
			"month": "3",
			"day":   "15",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 3,
				"day": 15
			},
			"enterprise": "testenterprise",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "GPT-5",
					"unitType": "credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 1500.5,
					"grossAmount": 15.005,
					"discountQuantity": 100.5,
					"discountAmount": 1.005,
					"netQuantity": 1400.0,
					"netAmount": 14.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Year:  2025,
		Month: 3,
		Day:   15,
	}
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Fatalf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if report.UsageItems[0].GrossQuantity != 1500.5 {
		t.Errorf("Expected GrossQuantity 1500.5, got %v", report.UsageItems[0].GrossQuantity)
	}
	if report.UsageItems[0].NetAmount != 14.0 {
		t.Errorf("Expected NetAmount 14.0, got %v", report.UsageItems[0].NetAmount)
	}
}

func TestEnterpriseService_GetAICreditUsage_WithCostCenter(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"cost_center_id": "cc-engineering",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"costCenter": {
				"id": "cc-engineering",
				"name": "Engineering Team"
			},
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "GPT-5",
					"unitType": "credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 5000,
					"grossAmount": 50.0,
					"discountQuantity": 500,
					"discountAmount": 5.0,
					"netQuantity": 4500,
					"netAmount": 45.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		CostCenterID: "cc-engineering",
	}
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if report.CostCenter == nil {
		t.Error("Expected CostCenter to be set")
	}
	if report.CostCenter.ID != "cc-engineering" {
		t.Errorf("Expected CostCenter ID cc-engineering, got %v", report.CostCenter.ID)
	}
	if report.CostCenter.Name != "Engineering Team" {
		t.Errorf("Expected CostCenter Name Engineering Team, got %v", report.CostCenter.Name)
	}
}

func TestEnterpriseService_GetAICreditUsage_WithOrganization(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"organization": "acme-corp",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"organization": "acme-corp",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "GPT-5",
					"unitType": "credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 2000,
					"grossAmount": 20.0,
					"discountQuantity": 200,
					"discountAmount": 2.0,
					"netQuantity": 1800,
					"netAmount": 18.0
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Organization: "acme-corp",
	}
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if report.Organization == nil || *report.Organization != "acme-corp" {
		t.Errorf("Expected organization acme-corp, got %v", report.Organization)
	}
}

func TestEnterpriseService_GetAICreditUsage_InvalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Enterprise.GetAICreditUsage(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestEnterpriseService_GetAICreditUsage_APIError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "Bad Request"}`)
	})

	ctx := t.Context()
	_, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetAICreditUsage should return error on 400 response")
	}
	if resp == nil {
		t.Error("GetAICreditUsage should return response even on error")
	}
}

func TestEnterpriseService_GetAICreditUsage_NoItems(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025
			},
			"enterprise": "GitHub",
			"usageItems": []
		}`)
	})

	ctx := t.Context()
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", nil)
	if err != nil {
		t.Errorf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if len(report.UsageItems) != 0 {
		t.Errorf("Expected 0 usage items, got %v", len(report.UsageItems))
	}
}

func TestEnterpriseService_GetAICreditUsage_CanceledContext(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", nil)
	if err == nil {
		t.Error("GetAICreditUsage should return error with canceled context")
	}
	_ = resp
}

func TestEnterpriseService_GetAICreditUsage_WithAllOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/test-enterprise/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":           "2025",
			"month":          "2",
			"day":            "28",
			"organization":   "tech-org",
			"user":           "bob",
			"model":          "GPT-4",
			"product":        "Copilot",
			"cost_center_id": "cc-dev",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"month": 2,
				"day": 28
			},
			"enterprise": "GitHub",
			"organization": "tech-org",
			"user": "bob",
			"product": "Copilot",
			"model": "GPT-4",
			"costCenter": {
				"id": "cc-dev",
				"name": "Development Team"
			},
			"usageItems": []
		}`)
	})

	ctx := t.Context()
	opts := &EnterprisePremiumRequestUsageReportOptions{
		Year:         2025,
		Month:        2,
		Day:          28,
		Organization: "tech-org",
		User:         "bob",
		Model:        "GPT-4",
		Product:      "Copilot",
		CostCenterID: "cc-dev",
	}
	report, resp, err := client.Enterprise.GetAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetAICreditUsage returned error: %v", err)
	}

	if resp == nil {
		t.Error("GetAICreditUsage returned nil response")
	}

	if report == nil {
		t.Error("Expected non-nil report")
	}
}
