// Copyright 2021 The go-github AUTHORS. All rights reserved.
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

func TestBillingService_GetPackagesBillingOrg(t *testing.T) {
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

	ctx := context.Background()
	hook, _, err := client.Billing.GetPackagesBillingOrg(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetPackagesBillingOrg returned error: %v", err)
	}

	want := &PackageBilling{
		TotalGigabytesBandwidthUsed:     50,
		TotalPaidGigabytesBandwidthUsed: 40,
		IncludedGigabytesBandwidth:      10,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetPackagesBillingOrg returned %+v, want %+v", hook, want)
	}

	const methodName = "GetPackagesBillingOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPackagesBillingOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPackagesBillingOrg(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPackagesBillingOrg_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetPackagesBillingOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetStorageBillingOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/settings/billing/shared-storage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"days_left_in_billing_cycle": 20,
				"estimated_paid_storage_for_month": 15.25,
				"estimated_storage_for_month": 40
			}`)
	})

	ctx := context.Background()
	hook, _, err := client.Billing.GetStorageBillingOrg(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetStorageBillingOrg returned error: %v", err)
	}

	want := &StorageBilling{
		DaysLeftInBillingCycle:       20,
		EstimatedPaidStorageForMonth: 15.25,
		EstimatedStorageForMonth:     40,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetStorageBillingOrg returned %+v, want %+v", hook, want)
	}

	const methodName = "GetStorageBillingOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetStorageBillingOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetStorageBillingOrg(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetStorageBillingOrg_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetStorageBillingOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetPackagesBillingUser(t *testing.T) {
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

	ctx := context.Background()
	hook, _, err := client.Billing.GetPackagesBillingUser(ctx, "u")
	if err != nil {
		t.Errorf("Billing.GetPackagesBillingUser returned error: %v", err)
	}

	want := &PackageBilling{
		TotalGigabytesBandwidthUsed:     50,
		TotalPaidGigabytesBandwidthUsed: 40,
		IncludedGigabytesBandwidth:      10,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetPackagesBillingUser returned %+v, want %+v", hook, want)
	}

	const methodName = "GetPackagesBillingUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPackagesBillingUser(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPackagesBillingUser(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPackagesBillingUser_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetPackagesBillingUser(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetStorageBillingUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/settings/billing/shared-storage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"days_left_in_billing_cycle": 20,
				"estimated_paid_storage_for_month": 15.25,
				"estimated_storage_for_month": 40
			}`)
	})

	ctx := context.Background()
	hook, _, err := client.Billing.GetStorageBillingUser(ctx, "u")
	if err != nil {
		t.Errorf("Billing.GetStorageBillingUser returned error: %v", err)
	}

	want := &StorageBilling{
		DaysLeftInBillingCycle:       20,
		EstimatedPaidStorageForMonth: 15.25,
		EstimatedStorageForMonth:     40,
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetStorageBillingUser returned %+v, want %+v", hook, want)
	}

	const methodName = "GetStorageBillingUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetStorageBillingUser(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetStorageBillingUser(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetStorageBillingUser_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetStorageBillingUser(ctx, "%")
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
	testJSONMarshal(t, &PackageBilling{}, "{}")

	u := &PackageBilling{
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

func TestBillingService_GetAdvancedSecurityActiveCommittersOrg(t *testing.T) {
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
          "last_pushed_date": "2021-10-25"
        }
      ]
    }
  ]
}`)
	})

	ctx := context.Background()
	opts := &ListOptions{Page: 2, PerPage: 50}
	hook, _, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetAdvancedSecurityActiveCommittersOrg	 returned error: %v", err)
	}

	want := &ActiveCommitters{
		TotalAdvancedSecurityCommitters:     2,
		TotalCount:                          2,
		MaximumAdvancedSecurityCommitters:   3,
		PurchasedAdvancedSecurityCommitters: 4,
		Repositories: []*RepositoryActiveCommitters{
			{
				Name:                       Ptr("octocat-org/Hello-World"),
				AdvancedSecurityCommitters: Ptr(2),
				AdvancedSecurityCommittersBreakdown: []*AdvancedSecurityCommittersBreakdown{
					{
						UserLogin:      Ptr("octokitten"),
						LastPushedDate: Ptr("2021-10-25"),
					},
				},
			},
		},
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetAdvancedSecurityActiveCommittersOrg returned %+v, want %+v", hook, want)
	}

	const methodName = "GetAdvancedSecurityActiveCommittersOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetAdvancedSecurityActiveCommittersOrg_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetUsageReportOrg(t *testing.T) {
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

	ctx := context.Background()
	opts := &UsageReportOptions{
		Year:  Ptr(2023),
		Month: Ptr(8),
	}
	report, _, err := client.Billing.GetUsageReportOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetUsageReportOrg returned error: %v", err)
	}

	want := &UsageReport{
		UsageItems: []*UsageItem{
			{
				Date:             Ptr("2023-08-01"),
				Product:          Ptr("Actions"),
				SKU:              Ptr("Actions Linux"),
				Quantity:         Ptr(100.0),
				UnitType:         Ptr("minutes"),
				PricePerUnit:     Ptr(0.008),
				GrossAmount:      Ptr(0.8),
				DiscountAmount:   Ptr(0.0),
				NetAmount:        Ptr(0.8),
				OrganizationName: Ptr("GitHub"),
				RepositoryName:   Ptr("github/example"),
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetUsageReportOrg returned %+v, want %+v", report, want)
	}

	const methodName = "GetUsageReportOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetUsageReportOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetUsageReportOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetUsageReportOrg_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetUsageReportOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetUsageReportUser(t *testing.T) {
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

	ctx := context.Background()
	opts := &UsageReportOptions{
		Day: Ptr(15),
	}
	report, _, err := client.Billing.GetUsageReportUser(ctx, "u", opts)
	if err != nil {
		t.Errorf("Billing.GetUsageReportUser returned error: %v", err)
	}

	want := &UsageReport{
		UsageItems: []*UsageItem{
			{
				Date:           Ptr("2023-08-15"),
				Product:        Ptr("Codespaces"),
				SKU:            Ptr("Codespaces Linux"),
				Quantity:       Ptr(50.0),
				UnitType:       Ptr("hours"),
				PricePerUnit:   Ptr(0.18),
				GrossAmount:    Ptr(9.0),
				DiscountAmount: Ptr(1.0),
				NetAmount:      Ptr(8.0),
				RepositoryName: Ptr("user/example"),
			},
		},
	}
	if !cmp.Equal(report, want) {
		t.Errorf("Billing.GetUsageReportUser returned %+v, want %+v", report, want)
	}

	const methodName = "GetUsageReportUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetUsageReportUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetUsageReportUser(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetUsageReportUser_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetUsageReportUser(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetPremiumRequestUsageReportOrg(t *testing.T) {
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
	ctx := context.Background()
	opts := &PremiumRequestUsageReportOptions{
		UsageReportOptions: UsageReportOptions{
			Year:  Ptr(2025),
			Month: Ptr(10),
		},
		User: Ptr("testuser"),
	}
	report, _, err := client.Billing.GetPremiumRequestUsageReportOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.GetPremiumRequestUsageReportOrg returned error: %v", err)
	}
	want := &PremiumRequestUsageReport{
		TimePeriod: TimePeriod{
			Year:  2025,
			Month: Ptr(10),
		},
		Organization: "GitHub",
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
		t.Errorf("Billing.GetPremiumRequestUsageReportOrg returned %+v, want %+v", report, want)
	}

	const methodName = "GetPremiumRequestUsageReportOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPremiumRequestUsageReportOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPremiumRequestUsageReportOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPremiumRequestUsageReportOrg_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetPremiumRequestUsageReportOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetPremiumRequestUsageReportUser(t *testing.T) {
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
			"organization": "UserOrg",
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
	ctx := context.Background()
	opts := &PremiumRequestUsageReportOptions{
		UsageReportOptions: UsageReportOptions{
			Year: Ptr(2025),
			Day:  Ptr(15),
		},
	}
	report, _, err := client.Billing.GetPremiumRequestUsageReportUser(ctx, "u", opts)
	if err != nil {
		t.Errorf("Billing.GetPremiumRequestUsageReportUser returned error: %v", err)
	}
	want := &PremiumRequestUsageReport{
		TimePeriod: TimePeriod{
			Year: 2025,
			Day:  Ptr(15),
		},
		Organization: "UserOrg",
		Product:      Ptr("Copilot"),
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
		t.Errorf("Billing.GetPremiumRequestUsageReportUser returned %+v, want %+v", report, want)
	}

	const methodName = "GetPremiumRequestUsageReportUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetPremiumRequestUsageReportUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetPremiumRequestUsageReportUser(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetPremiumRequestUsageReportUser_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Billing.GetPremiumRequestUsageReportUser(ctx, "%", nil)
	testURLParseError(t, err)
}
