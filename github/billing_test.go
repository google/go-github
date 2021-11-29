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

func TestBillingService_GetActionsBillingOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/settings/billing/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"total_minutes_used": 305,
				"total_paid_minutes_used": 0,
				"included_minutes": 3000,
				"minutes_used_breakdown": {
					"UBUNTU": 205,
					"MACOS": 10,
					"WINDOWS": 90
				}
			}`)
	})

	ctx := context.Background()
	hook, _, err := client.Billing.GetActionsBillingOrg(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetActionsBillingOrg returned error: %v", err)
	}

	want := &ActionBilling{
		TotalMinutesUsed:     305,
		TotalPaidMinutesUsed: 0,
		IncludedMinutes:      3000,
		MinutesUsedBreakdown: MinutesUsedBreakdown{
			Ubuntu:  205,
			MacOS:   10,
			Windows: 90,
		},
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetActionsBillingOrg returned %+v, want %+v", hook, want)
	}

	const methodName = "GetActionsBillingOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetActionsBillingOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetActionsBillingOrg(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetActionsBillingOrg_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetActionsBillingOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetPackagesBillingOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetPackagesBillingOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetStorageBillingOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetStorageBillingOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetActionsBillingUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/settings/billing/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"total_minutes_used": 10,
				"total_paid_minutes_used": 0,
				"included_minutes": 3000,
				"minutes_used_breakdown": {
					"UBUNTU": 205,
					"MACOS": 10,
					"WINDOWS": 90
				}
			}`)
	})

	ctx := context.Background()
	hook, _, err := client.Billing.GetActionsBillingUser(ctx, "u")
	if err != nil {
		t.Errorf("Billing.GetActionsBillingUser returned error: %v", err)
	}

	want := &ActionBilling{
		TotalMinutesUsed:     10,
		TotalPaidMinutesUsed: 0,
		IncludedMinutes:      3000,
		MinutesUsedBreakdown: MinutesUsedBreakdown{
			Ubuntu:  205,
			MacOS:   10,
			Windows: 90,
		},
	}
	if !cmp.Equal(hook, want) {
		t.Errorf("Billing.GetActionsBillingUser returned %+v, want %+v", hook, want)
	}

	const methodName = "GetActionsBillingUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetActionsBillingOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetActionsBillingUser(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetActionsBillingUser_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetActionsBillingUser(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetPackagesBillingUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetPackagesBillingUser(ctx, "%")
	testURLParseError(t, err)
}

func TestBillingService_GetStorageBillingUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetStorageBillingUser(ctx, "%")
	testURLParseError(t, err)
}

func TestMinutesUsedBreakdown_Marshal(t *testing.T) {
	testJSONMarshal(t, &MinutesUsedBreakdown{}, "{}")

	u := &MinutesUsedBreakdown{
		Ubuntu:  1,
		MacOS:   1,
		Windows: 1,
	}

	want := `{
		"UBUNTU": 1,
		"MACOS": 1,
		"WINDOWS": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestActionBilling_Marshal(t *testing.T) {
	testJSONMarshal(t, &MinutesUsedBreakdown{}, "{}")

	u := &ActionBilling{
		TotalMinutesUsed:     1,
		TotalPaidMinutesUsed: 1,
		IncludedMinutes:      1,
		MinutesUsedBreakdown: MinutesUsedBreakdown{
			Ubuntu:  1,
			MacOS:   1,
			Windows: 1,
		},
	}

	want := `{
		"total_minutes_used": 1,
		"total_paid_minutes_used": 1,
		"included_minutes": 1,
		"minutes_used_breakdown": {
			"UBUNTU": 1,
			"MACOS": 1,
			"WINDOWS": 1
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPackageBilling_Marshal(t *testing.T) {
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/settings/billing/advanced-security", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "total_advanced_security_committers": 2,
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
	hook, _, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "o")
	if err != nil {
		t.Errorf("Billing.GetAdvancedSecurityActiveCommittersOrg	 returned error: %v", err)
	}

	want := &ActiveCommitters{
		TotalAdvancedSecurityCommitters: 2,
		Repositories: []*RepositoryActiveCommitters{
			{
				Name:                       String("octocat-org/Hello-World"),
				AdvancedSecurityCommitters: Int(2),
				AdvancedSecurityCommittersBreakdown: []*AdvancedSecurityCommittersBreakdown{
					{
						UserLogin:      String("octokitten"),
						LastPushedDate: String("2021-10-25"),
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
		_, _, err = client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestBillingService_GetAdvancedSecurityActiveCommittersOrg_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Billing.GetAdvancedSecurityActiveCommittersOrg(ctx, "%")
	testURLParseError(t, err)
}
