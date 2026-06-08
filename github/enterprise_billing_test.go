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

func TestBillingService_GetEnterpriseUsageReport(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseUsageReport returned error: %v", err)
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
		t.Errorf("GetEnterpriseUsageReport returned %+v, want %+v", report, want)
	}
}

func TestBillingService_GetEnterpriseUsageReport_WithCostCenter(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseUsageReport returned error: %v", err)
	}

	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
	if report.UsageItems[0].NetAmount != 1.8 {
		t.Errorf("Expected NetAmount 1.8, got %v", report.UsageItems[0].NetAmount)
	}
}

func TestBillingService_GetEnterpriseUsageReport_InvalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetEnterpriseUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetEnterpriseUsageSummary(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseUsageSummary returned error: %v", err)
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
		t.Errorf("GetEnterpriseUsageSummary returned %+v, want %+v", report, want)
	}
}

func TestBillingService_GetEnterpriseUsageSummary_MultipleItems(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseUsageSummary returned error: %v", err)
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

func TestBillingService_GetEnterpriseUsageSummary_WithRepository(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseUsageSummary(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseUsageSummary returned error: %v", err)
	}

	if len(report.UsageItems) != 1 {
		t.Errorf("Expected 1 usage item, got %v", len(report.UsageItems))
	}
}

func TestBillingService_GetEnterprisePremiumRequestUsageReport(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterprisePremiumRequestUsageReport(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterprisePremiumRequestUsageReport returned error: %v", err)
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

func TestBillingService_GetEnterpriseAICreditUsage(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseAICreditUsage returned error: %v", err)
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

func TestBillingService_GetEnterpriseAICreditUsage_FloatQuantities(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Fatalf("GetEnterpriseAICreditUsage returned error: %v", err)
	}

	if report.UsageItems[0].GrossQuantity != 1500.5 {
		t.Errorf("Expected GrossQuantity 1500.5, got %v", report.UsageItems[0].GrossQuantity)
	}
	if report.UsageItems[0].NetAmount != 14.0 {
		t.Errorf("Expected NetAmount 14.0, got %v", report.UsageItems[0].NetAmount)
	}
}

func TestBillingService_GetEnterpriseAICreditUsage_WithCostCenter(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseAICreditUsage returned error: %v", err)
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

func TestBillingService_GetEnterpriseAICreditUsage_WithOrganization(t *testing.T) {
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
	report, _, err := client.Billing.GetEnterpriseAICreditUsage(ctx, "test-enterprise", opts)
	if err != nil {
		t.Errorf("GetEnterpriseAICreditUsage returned error: %v", err)
	}

	if report.Organization == nil || *report.Organization != "acme-corp" {
		t.Errorf("Expected organization acme-corp, got %v", report.Organization)
	}
}
