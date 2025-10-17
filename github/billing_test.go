// Copyright 2021 The go-github AUTHORS. All rights reserved.
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

func TestBillingService_GetOrganizationPackagesBilling(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/settings/billing/packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"total_gigabytes_bandwidth_used": 50,
				"total_paid_gigabytes_bandwidth_used": 40,
				"included_gigabytes_bandwidth": 10
			}`)
	})

	ctx := t.Context()
	hook, _, err := client.Billing.GetOrganizationPackagesBilling(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetOrganizationPackagesBilling returned error: %v", err)
	}

	want := &PackagesBilling{
		TotalGigabytesBandwidthUsed:     50,
		TotalPaidGigabytesBandwidthUsed: 40,
		IncludedGigabytesBandwidth:      10,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetOrganizationPackagesBilling returned %+v, want %+v", hook, want)
	}

	const methodName = "GetOrganizationPackagesBilling"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationPackagesBilling(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationPackagesBilling(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetOrganizationPackagesBilling_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationPackagesBilling(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetOrganizationStorageBilling(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/settings/billing/shared-storage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"days_left_in_billing_cycle": 20,
				"estimated_paid_storage_for_month": 15,
				"estimated_storage_for_month": 40
			}`)
	})

	ctx := t.Context()
	hook, _, err := client.Billing.GetOrganizationStorageBilling(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetOrganizationStorageBilling returned error: %v", err)
	}

	want := &StorageBilling{
		DaysLeftInBillingCycle:       20,
		EstimatedPaidStorageForMonth: 15,
		EstimatedStorageForMonth:     40,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetOrganizationStorageBilling returned %+v, want %+v", hook, want)
	}

	const methodName = "GetOrganizationStorageBilling"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationStorageBilling(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationStorageBilling(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetOrganizationStorageBilling_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationStorageBilling(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetPackagesBilling(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/settings/billing/packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"total_gigabytes_bandwidth_used": 50,
				"total_paid_gigabytes_bandwidth_used": 40,
				"included_gigabytes_bandwidth": 10
			}`)
	})

	ctx := t.Context()
	hook, _, err := client.Billing.GetPackagesBilling(ctx, "u")
	if err != nil {
		t.Errorf("Billing.GetPackagesBilling returned error: %v", err)
	}

	want := &PackagesBilling{
		TotalGigabytesBandwidthUsed:     50,
		TotalPaidGigabytesBandwidthUsed: 40,
		IncludedGigabytesBandwidth:      10,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetPackagesBilling returned %+v, want %+v", hook, want)
	}

	const methodName = "GetPackagesBilling"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPackagesBilling(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPackagesBilling(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPackagesBilling_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetPackagesBilling(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetStorageBilling(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/settings/billing/shared-storage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"days_left_in_billing_cycle": 20,
				"estimated_paid_storage_for_month": 15,
				"estimated_storage_for_month": 40
			}`)
	})

	ctx := t.Context()
	hook, _, err := client.Billing.GetStorageBilling(ctx, "u")
	if err != nil {
		t.Errorf("Billing.GetStorageBilling returned error: %v", err)
	}

	want := &StorageBilling{
		DaysLeftInBillingCycle:       20,
		EstimatedPaidStorageForMonth: 15,
		EstimatedStorageForMonth:     40,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetStorageBilling returned %+v, want %+v", hook, want)
	}

	const methodName = "GetStorageBilling"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetStorageBilling(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetStorageBilling(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetStorageBilling_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetStorageBilling(ctx, "%")
	testURLParseError(t, err)
}

func TestMinutesUsedBreakdown_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MinutesUsedBreakdown{}, "{}")

	u := &MinutesUsedBreakdown{
		"UBUNTU":  1,
		"MACOS":   1,
		"WINDOWS": 1,
	}

	want := `{
		"UBUNTU": 1,
		"MACOS": 1,
		"WINDOWS": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestPackageBilling_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackagesBilling{}, "{}")

	u := &PackagesBilling{
		TotalGigabytesBandwidthUsed:     1,
		TotalPaidGigabytesBandwidthUsed: 1,
		IncludedGigabytesBandwidth:      1,
	}

	want := `{
		"total_gigabytes_bandwidth_used": 1,
		"total_paid_gigabytes_bandwidth_used": 1,
		"included_gigabytes_bandwidth": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestStorageBilling_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &StorageBilling{}, "{}")

	u := &StorageBilling{
		DaysLeftInBillingCycle:       1,
		EstimatedPaidStorageForMonth: 1,
		EstimatedStorageForMonth:     1,
	}

	want := `{
		"days_left_in_billing_cycle": 1,
		"estimated_paid_storage_for_month": 1,
		"estimated_storage_for_month": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestBillingService_GetOrganizationAdvancedSecurityActiveCommitters(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/settings/billing/advanced-security", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "total_advanced_security_committers": 2,
  "total_count": 2,
  "maximum_advanced_security_committers": 3,
  "purchased_advanced_security_committers": 4,
  "repositories": [
    {
      "name": "octocat-org/Hello-World",
      "advanced_security_committers": 2,
      "advanced_security_committers_breakdown": [
        {
          "user_login": "octokitten",
          "last_pushed_date": "2021-10-25",
          "last_pushed_email": "octokitten@example.com"
        }
      ]
    }
  ]
}`)
	})

	ctx := t.Context()
	opts := &ActiveCommittersListOptions{
		nil,
		ListOptions{Page: 2, PerPage: 50},
	}
	hook, _, err := client.Billing.GetOrganizationAdvancedSecurityActiveCommitters(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetOrganizationAdvancedSecurityActiveCommitters returned error: %v", err)
	}

	want := &ActiveCommitters{
		TotalAdvancedSecurityCommitters:     Ptr(2),
		TotalCount:                          Ptr(2),
		MaximumAdvancedSecurityCommitters:   Ptr(3),
		PurchasedAdvancedSecurityCommitters: Ptr(4),
		Repositories: []*RepositoryActiveCommitters{
			{
				Name:                       "octocat-org/Hello-World",
				AdvancedSecurityCommitters: 2,
				AdvancedSecurityCommittersBreakdown: []*AdvancedSecurityCommittersBreakdown{
					{
						UserLogin:       "octokitten",
						LastPushedDate:  "2021-10-25",
						LastPushedEmail: "octokitten@example.com",
					},
				},
			},
		},
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetOrganizationAdvancedSecurityActiveCommitters returned %+v, want %+v", hook, want)
	}

	const methodName = "GetOrganizationAdvancedSecurityActiveCommitters"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationAdvancedSecurityActiveCommitters(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationAdvancedSecurityActiveCommitters(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetOrganizationAdvancedSecurityActiveCommitters_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationAdvancedSecurityActiveCommitters(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetOrganizationUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
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
	opts := &UsageReportOptions{
		Year:  Ptr(2023),
		Month: Ptr(8),
	}
	report, _, err := client.Billing.GetOrganizationUsageReport(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetOrganizationUsageReport returned error: %v", err)
	}

	want := &UsageReport{
		UsageItems: []*UsageItem{
			{
				Date:             "2023-08-01",
				Product:          "Actions",
				SKU:              "Actions Linux",
				Quantity:         100,
				UnitType:         "minutes",
				PricePerUnit:     0.008,
				GrossAmount:      0.8,
				DiscountAmount:   0.0,
				NetAmount:        0.8,
				OrganizationName: Ptr("GitHub"),
				RepositoryName:   Ptr("github/example"),
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetOrganizationUsageReport returned %+v, want %+v", report, want)
	}

	const methodName = "GetOrganizationUsageReport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationUsageReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationUsageReport(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetOrganizationUsageReport_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/settings/billing/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"day": "15",
		})
		fmt.Fprint(w, `{
			"usageItems": [
				{
					"date": "2023-08-15",
					"product": "Codespaces",
					"sku": "Codespaces Linux",
					"quantity": 50,
					"unitType": "hours",
					"pricePerUnit": 0.18,
					"grossAmount": 9.0,
					"discountAmount": 1.0,
					"netAmount": 8.0,
					"repositoryName": "user/example"
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &UsageReportOptions{
		Day: Ptr(15),
	}
	report, _, err := client.Billing.GetUsageReport(ctx, "u", opts)
	if err != nil {
		t.Errorf("Billing.GetUsageReport returned error: %v", err)
	}

	want := &UsageReport{
		UsageItems: []*UsageItem{
			{
				Date:           "2023-08-15",
				Product:        "Codespaces",
				SKU:            "Codespaces Linux",
				Quantity:       50,
				UnitType:       "hours",
				PricePerUnit:   0.18,
				GrossAmount:    9.0,
				DiscountAmount: 1.0,
				NetAmount:      8.0,
				RepositoryName: Ptr("user/example"),
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetUsageReport returned %+v, want %+v", report, want)
	}

	const methodName = "GetUsageReport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetUsageReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetUsageReport(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetUsageReport_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetOrganizationPremiumRequestUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/organizations/o/settings/billing/premium_request/usage", func(w http.ResponseWriter, r *http.Request) {
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
	opts := &PremiumRequestUsageReportOptions{
		Year:  Ptr(2025),
		Month: Ptr(10),
		User:  Ptr("testuser"),
	}
	report, _, err := client.Billing.GetOrganizationPremiumRequestUsageReport(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetOrganizationPremiumRequestUsageReport returned error: %v", err)
	}
	want := &PremiumRequestUsageReport{
		TimePeriod: PremiumRequestUsageTimePeriod{
			Year:  2025,
			Month: Ptr(10),
		},
		Organization: Ptr("GitHub"),
		User:         Ptr("testuser"),
		Product:      Ptr("Copilot"),
		Model:        Ptr("GPT-5"),
		UsageItems: []*PremiumRequestUsageItem{
			{
				Product:          "Copilot",
				SKU:              "Copilot Premium Request",
				Model:            "GPT-5",
				UnitType:         "requests",
				PricePerUnit:     0.04,
				GrossQuantity:    100,
				GrossAmount:      4.0,
				DiscountQuantity: 0,
				DiscountAmount:   0.0,
				NetQuantity:      100,
				NetAmount:        4.0,
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetOrganizationPremiumRequestUsageReport returned %+v, want %+v", report, want)
	}

	const methodName = "GetOrganizationPremiumRequestUsageReport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationPremiumRequestUsageReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationPremiumRequestUsageReport(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetOrganizationPremiumRequestUsageReport_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationPremiumRequestUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetPremiumRequestUsageReport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/settings/billing/premium_request/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"year": "2025",
			"day":  "15",
		})
		fmt.Fprint(w, `{
			"timePeriod": {
				"year": 2025,
				"day": 15
			},
			"user": "User",
			"product": "Copilot",
			"usageItems": [
				{
					"product": "Copilot",
					"sku": "Copilot Premium Request",
					"model": "GPT-4",
					"unitType": "requests",
					"pricePerUnit": 0.02,
					"grossQuantity": 50,
					"grossAmount": 1.0,
					"discountQuantity": 5,
					"discountAmount": 0.1,
					"netQuantity": 45,
					"netAmount": 0.9
				}
			]
		}`)
	})
	ctx := t.Context()
	opts := &PremiumRequestUsageReportOptions{
		Year: Ptr(2025),
		Day:  Ptr(15),
	}
	report, _, err := client.Billing.GetPremiumRequestUsageReport(ctx, "u", opts)
	if err != nil {
		t.Errorf("Billing.GetPremiumRequestUsageReport returned error: %v", err)
	}
	want := &PremiumRequestUsageReport{
		TimePeriod: PremiumRequestUsageTimePeriod{
			Year: 2025,
			Day:  Ptr(15),
		},
		User:    Ptr("User"),
		Product: Ptr("Copilot"),
		UsageItems: []*PremiumRequestUsageItem{
			{
				Product:          "Copilot",
				SKU:              "Copilot Premium Request",
				Model:            "GPT-4",
				UnitType:         "requests",
				PricePerUnit:     0.02,
				GrossQuantity:    50,
				GrossAmount:      1.0,
				DiscountQuantity: 5,
				DiscountAmount:   0.1,
				NetQuantity:      45,
				NetAmount:        0.9,
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetPremiumRequestUsageReport returned %+v, want %+v", report, want)
	}

	const methodName = "GetPremiumRequestUsageReport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPremiumRequestUsageReport(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPremiumRequestUsageReport(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPremiumRequestUsageReport_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetPremiumRequestUsageReport(ctx, "%", nil)
	testURLParseError(t, err)
}
