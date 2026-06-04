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

func TestBillingService_GetAICreditUsage(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	opts := &AICreditUsageOptions{
		Year:  2026,
		Month: 6,
	}

	mux.HandleFunc("/organizations/o/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2026",
			"month": "6",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2026,
				"month": 6
			},
			"organization": "o",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "Auto: Claude Haiku 4.5",
					"unitType": "ai-credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 3956.1799545,
					"grossAmount": 39.561799545,
					"discountQuantity": 3956.1799545,
					"discountAmount": 39.561799545,
					"netQuantity": 0.0,
					"netAmount": 0.0
				},
				{
					"product": "Copilot",
					"sku": "Copilot AI Credits",
					"model": "Auto: Claude Sonnet 4.6",
					"unitType": "ai-credits",
					"pricePerUnit": 0.01,
					"grossQuantity": 9498.6969165,
					"grossAmount": 94.986969165,
					"discountQuantity": 9564.710256,
					"discountAmount": 95.64710256,
					"netQuantity": 0.0,
					"netAmount": 0.0
				}
			]
		}`)
	})

	ctx := t.Context()
	usage, _, err := client.Billing.GetAICreditUsage(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetAICreditUsage returned error: %v", err)
	}

	month := 6
	want := &AICreditUsage{
		TimePeriod: AICreditTimePeriod{
			Year:  2026,
			Month: &month,
		},
		Organization: "o",
		UsageItems: []*AICreditUsageItem{
			{
				Product:          "Copilot",
				SKU:              "Copilot AI Credits",
				Model:            "Auto: Claude Haiku 4.5",
				UnitType:         "ai-credits",
				PricePerUnit:     0.01,
				GrossQuantity:    3956.1799545,
				GrossAmount:      39.561799545,
				DiscountQuantity: 3956.1799545,
				DiscountAmount:   39.561799545,
				NetQuantity:      0.0,
				NetAmount:        0.0,
			},
			{
				Product:          "Copilot",
				SKU:              "Copilot AI Credits",
				Model:            "Auto: Claude Sonnet 4.6",
				UnitType:         "ai-credits",
				PricePerUnit:     0.01,
				GrossQuantity:    9498.6969165,
				GrossAmount:      94.986969165,
				DiscountQuantity: 9564.710256,
				DiscountAmount:   95.64710256,
				NetQuantity:      0.0,
				NetAmount:        0.0,
			},
		},
	}
	if !cmp.Equal(usage, want) {
		t.Errorf("Billing.GetAICreditUsage returned %+v, want %+v", usage, want)
	}

	const methodName = "GetAICreditUsage"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetAICreditUsage(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetAICreditUsage(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetAICreditUsage_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetAICreditUsage(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetAICreditUsage_WithOptionalParams(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	opts := &AICreditUsageOptions{
		Year:  2026,
		Model: "Claude Sonnet 4.6",
	}

	mux.HandleFunc("/organizations/o/settings/billing/ai_credit/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year":  "2026",
			"model": "Claude Sonnet 4.6",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2026,
				"month": 6
			},
			"organization": "o",
			"usageItems": []
		}`)
	})

	ctx := t.Context()
	usage, _, err := client.Billing.GetAICreditUsage(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetAICreditUsage with optional params returned error: %v", err)
	}

	if usage.Organization != "o" {
		t.Errorf("Expected organization 'o', got %v", usage.Organization)
	}
}

func TestAICreditTimePeriod_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AICreditTimePeriod{}, `{"year":0}`)

	month := 6
	u := &AICreditTimePeriod{
		Year:  2026,
		Month: &month,
	}

	want := `{
		"year": 2026,
		"month": 6
	}`

	testJSONMarshal(t, u, want)
}

func TestAICreditTimePeriod_Marshal_WithDay(t *testing.T) {
	t.Parallel()
	month := 6
	day := 15
	u := &AICreditTimePeriod{
		Year:  2026,
		Month: &month,
		Day:   &day,
	}

	want := `{
		"year": 2026,
		"month": 6,
		"day": 15
	}`

	testJSONMarshal(t, u, want)
}

func TestAICreditUsageItem_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AICreditUsageItem{}, `{"product":"","sku":"","model":"","unitType":"","pricePerUnit":0,"grossQuantity":0,"grossAmount":0,"discountQuantity":0,"discountAmount":0,"netQuantity":0,"netAmount":0}`)

	u := &AICreditUsageItem{
		Product:          "Copilot",
		SKU:              "Copilot AI Credits",
		Model:            "Auto: Claude Haiku 4.5",
		UnitType:         "ai-credits",
		PricePerUnit:     0.01,
		GrossQuantity:    3956.1799545,
		GrossAmount:      39.561799545,
		DiscountQuantity: 3956.1799545,
		DiscountAmount:   39.561799545,
		NetQuantity:      0.0,
		NetAmount:        0.0,
	}

	want := `{
		"product": "Copilot",
		"sku": "Copilot AI Credits",
		"model": "Auto: Claude Haiku 4.5",
		"unitType": "ai-credits",
		"pricePerUnit": 0.01,
		"grossQuantity": 3956.1799545,
		"grossAmount": 39.561799545,
		"discountQuantity": 3956.1799545,
		"discountAmount": 39.561799545,
		"netQuantity": 0,
		"netAmount": 0
	}`

	testJSONMarshal(t, u, want)
}

func TestAICreditUsage_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AICreditUsage{}, `{"timePeriod":{"year":0},"organization":""}`)

	month := 6
	u := &AICreditUsage{
		TimePeriod: AICreditTimePeriod{
			Year:  2026,
			Month: &month,
		},
		Organization: "test-org",
		UsageItems: []*AICreditUsageItem{
			{
				Product:          "Copilot",
				SKU:              "Copilot AI Credits",
				Model:            "Auto: Claude Haiku 4.5",
				UnitType:         "ai-credits",
				PricePerUnit:     0.01,
				GrossQuantity:    3956.1799545,
				GrossAmount:      39.561799545,
				DiscountQuantity: 3956.1799545,
				DiscountAmount:   39.561799545,
				NetQuantity:      0.0,
				NetAmount:        0.0,
			},
		},
	}

	want := `{
		"timePeriod": {
			"year": 2026,
			"month": 6
		},
		"organization": "test-org",
		"usageItems": [
			{
				"product": "Copilot",
				"sku": "Copilot AI Credits",
				"model": "Auto: Claude Haiku 4.5",
				"unitType": "ai-credits",
				"pricePerUnit": 0.01,
				"grossQuantity": 3956.1799545,
				"grossAmount": 39.561799545,
				"discountQuantity": 3956.1799545,
				"discountAmount": 39.561799545,
				"netQuantity": 0,
				"netAmount": 0
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
