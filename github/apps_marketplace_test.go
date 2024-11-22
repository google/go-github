// Copyright 2017 The go-github AUTHORS. All rights reserved.
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

func TestMarketplaceService_ListPlans(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/plans", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = false
	ctx := context.Background()
	plans, _, err := client.Marketplace.ListPlans(ctx, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlans returned error: %v", err)
	}

	want := []*MarketplacePlan{{ID: Ptr(int64(1))}}
	if !cmp.Equal(plans, want) {
		t.Errorf("Marketplace.ListPlans returned %+v, want %+v", plans, want)
	}

	const methodName = "ListPlans"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Marketplace.ListPlans(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarketplaceService_Stubbed_ListPlans(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/stubbed/plans", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = true
	ctx := context.Background()
	plans, _, err := client.Marketplace.ListPlans(ctx, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlans (Stubbed) returned error: %v", err)
	}

	want := []*MarketplacePlan{{ID: Ptr(int64(1))}}
	if !cmp.Equal(plans, want) {
		t.Errorf("Marketplace.ListPlans (Stubbed) returned %+v, want %+v", plans, want)
	}
}

func TestMarketplaceService_ListPlanAccountsForPlan(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/plans/1/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = false
	ctx := context.Background()
	accounts, _, err := client.Marketplace.ListPlanAccountsForPlan(ctx, 1, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlanAccountsForPlan returned error: %v", err)
	}

	want := []*MarketplacePlanAccount{{ID: Ptr(int64(1))}}
	if !cmp.Equal(accounts, want) {
		t.Errorf("Marketplace.ListPlanAccountsForPlan returned %+v, want %+v", accounts, want)
	}

	const methodName = "ListPlanAccountsForPlan"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Marketplace.ListPlanAccountsForPlan(ctx, 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarketplaceService_Stubbed_ListPlanAccountsForPlan(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/stubbed/plans/1/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = true
	ctx := context.Background()
	accounts, _, err := client.Marketplace.ListPlanAccountsForPlan(ctx, 1, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlanAccountsForPlan (Stubbed) returned error: %v", err)
	}

	want := []*MarketplacePlanAccount{{ID: Ptr(int64(1))}}
	if !cmp.Equal(accounts, want) {
		t.Errorf("Marketplace.ListPlanAccountsForPlan (Stubbed) returned %+v, want %+v", accounts, want)
	}
}

func TestMarketplaceService_GetPlanAccountForAccount(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "marketplace_pending_change": {"id": 77}}`)
	})

	client.Marketplace.Stubbed = false
	ctx := context.Background()
	account, _, err := client.Marketplace.GetPlanAccountForAccount(ctx, 1)
	if err != nil {
		t.Errorf("Marketplace.GetPlanAccountForAccount returned error: %v", err)
	}

	want := &MarketplacePlanAccount{ID: Ptr(int64(1)), MarketplacePendingChange: &MarketplacePendingChange{ID: Ptr(int64(77))}}
	if !cmp.Equal(account, want) {
		t.Errorf("Marketplace.GetPlanAccountForAccount returned %+v, want %+v", account, want)
	}

	const methodName = "GetPlanAccountForAccount"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Marketplace.GetPlanAccountForAccount(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarketplaceService_Stubbed_GetPlanAccountForAccount(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/marketplace_listing/stubbed/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	client.Marketplace.Stubbed = true
	ctx := context.Background()
	account, _, err := client.Marketplace.GetPlanAccountForAccount(ctx, 1)
	if err != nil {
		t.Errorf("Marketplace.GetPlanAccountForAccount (Stubbed) returned error: %v", err)
	}

	want := &MarketplacePlanAccount{ID: Ptr(int64(1))}
	if !cmp.Equal(account, want) {
		t.Errorf("Marketplace.GetPlanAccountForAccount (Stubbed) returned %+v, want %+v", account, want)
	}
}

func TestMarketplaceService_ListMarketplacePurchasesForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/marketplace_purchases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"billing_cycle":"monthly"}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = false
	ctx := context.Background()
	purchases, _, err := client.Marketplace.ListMarketplacePurchasesForUser(ctx, opt)
	if err != nil {
		t.Errorf("Marketplace.ListMarketplacePurchasesForUser returned error: %v", err)
	}

	want := []*MarketplacePurchase{{BillingCycle: Ptr("monthly")}}
	if !cmp.Equal(purchases, want) {
		t.Errorf("Marketplace.ListMarketplacePurchasesForUser returned %+v, want %+v", purchases, want)
	}

	const methodName = "ListMarketplacePurchasesForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Marketplace.ListMarketplacePurchasesForUser(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarketplaceService_Stubbed_ListMarketplacePurchasesForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/marketplace_purchases/stubbed", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"billing_cycle":"monthly"}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	client.Marketplace.Stubbed = true
	ctx := context.Background()
	purchases, _, err := client.Marketplace.ListMarketplacePurchasesForUser(ctx, opt)
	if err != nil {
		t.Errorf("Marketplace.ListMarketplacePurchasesForUser returned error: %v", err)
	}

	want := []*MarketplacePurchase{{BillingCycle: Ptr("monthly")}}
	if !cmp.Equal(purchases, want) {
		t.Errorf("Marketplace.ListMarketplacePurchasesForUser returned %+v, want %+v", purchases, want)
	}
}

func TestMarketplacePlan_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MarketplacePlan{}, "{}")

	u := &MarketplacePlan{
		URL:                 Ptr("u"),
		AccountsURL:         Ptr("au"),
		ID:                  Ptr(int64(1)),
		Number:              Ptr(1),
		Name:                Ptr("n"),
		Description:         Ptr("d"),
		MonthlyPriceInCents: Ptr(1),
		YearlyPriceInCents:  Ptr(1),
		PriceModel:          Ptr("pm"),
		UnitName:            Ptr("un"),
		Bullets:             &[]string{"b"},
		State:               Ptr("s"),
		HasFreeTrial:        Ptr(false),
	}

	want := `{
		"url": "u",
		"accounts_url": "au",
		"id": 1,
		"number": 1,
		"name": "n",
		"description": "d",
		"monthly_price_in_cents": 1,
		"yearly_price_in_cents": 1,
		"price_model": "pm",
		"unit_name": "un",
		"bullets": ["b"],
		"state": "s",
		"has_free_trial": false
	}`

	testJSONMarshal(t, u, want)
}

func TestMarketplacePurchase_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MarketplacePurchase{}, "{}")

	u := &MarketplacePurchase{
		BillingCycle:    Ptr("bc"),
		NextBillingDate: &Timestamp{referenceTime},
		UnitCount:       Ptr(1),
		Plan: &MarketplacePlan{
			URL:                 Ptr("u"),
			AccountsURL:         Ptr("au"),
			ID:                  Ptr(int64(1)),
			Number:              Ptr(1),
			Name:                Ptr("n"),
			Description:         Ptr("d"),
			MonthlyPriceInCents: Ptr(1),
			YearlyPriceInCents:  Ptr(1),
			PriceModel:          Ptr("pm"),
			UnitName:            Ptr("un"),
			Bullets:             &[]string{"b"},
			State:               Ptr("s"),
			HasFreeTrial:        Ptr(false),
		},
		OnFreeTrial:     Ptr(false),
		FreeTrialEndsOn: &Timestamp{referenceTime},
		UpdatedAt:       &Timestamp{referenceTime},
	}

	want := `{
		"billing_cycle": "bc",
		"next_billing_date": ` + referenceTimeStr + `,
		"unit_count": 1,
		"plan": {
			"url": "u",
			"accounts_url": "au",
			"id": 1,
			"number": 1,
			"name": "n",
			"description": "d",
			"monthly_price_in_cents": 1,
			"yearly_price_in_cents": 1,
			"price_model": "pm",
			"unit_name": "un",
			"bullets": ["b"],
			"state": "s",
			"has_free_trial": false
			},
		"on_free_trial": false,
		"free_trial_ends_on": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestMarketplacePendingChange_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MarketplacePendingChange{}, "{}")

	u := &MarketplacePendingChange{
		EffectiveDate: &Timestamp{referenceTime},
		UnitCount:     Ptr(1),
		ID:            Ptr(int64(1)),
		Plan: &MarketplacePlan{
			URL:                 Ptr("u"),
			AccountsURL:         Ptr("au"),
			ID:                  Ptr(int64(1)),
			Number:              Ptr(1),
			Name:                Ptr("n"),
			Description:         Ptr("d"),
			MonthlyPriceInCents: Ptr(1),
			YearlyPriceInCents:  Ptr(1),
			PriceModel:          Ptr("pm"),
			UnitName:            Ptr("un"),
			Bullets:             &[]string{"b"},
			State:               Ptr("s"),
			HasFreeTrial:        Ptr(false),
		},
	}

	want := `{
		"effective_date": ` + referenceTimeStr + `,
		"unit_count": 1,
		"id": 1,
		"plan": {
			"url": "u",
			"accounts_url": "au",
			"id": 1,
			"number": 1,
			"name": "n",
			"description": "d",
			"monthly_price_in_cents": 1,
			"yearly_price_in_cents": 1,
			"price_model": "pm",
			"unit_name": "un",
			"bullets": ["b"],
			"state": "s",
			"has_free_trial": false
			}
	}`

	testJSONMarshal(t, u, want)
}

func TestMarketplacePlanAccount_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &MarketplacePlanAccount{}, "{}")

	u := &MarketplacePlanAccount{
		URL:                      Ptr("u"),
		Type:                     Ptr("t"),
		ID:                       Ptr(int64(1)),
		Login:                    Ptr("l"),
		OrganizationBillingEmail: Ptr("obe"),
		MarketplacePurchase: &MarketplacePurchase{
			BillingCycle:    Ptr("bc"),
			NextBillingDate: &Timestamp{referenceTime},
			UnitCount:       Ptr(1),
			Plan: &MarketplacePlan{
				URL:                 Ptr("u"),
				AccountsURL:         Ptr("au"),
				ID:                  Ptr(int64(1)),
				Number:              Ptr(1),
				Name:                Ptr("n"),
				Description:         Ptr("d"),
				MonthlyPriceInCents: Ptr(1),
				YearlyPriceInCents:  Ptr(1),
				PriceModel:          Ptr("pm"),
				UnitName:            Ptr("un"),
				Bullets:             &[]string{"b"},
				State:               Ptr("s"),
				HasFreeTrial:        Ptr(false),
			},
			OnFreeTrial:     Ptr(false),
			FreeTrialEndsOn: &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
		},
		MarketplacePendingChange: &MarketplacePendingChange{
			EffectiveDate: &Timestamp{referenceTime},
			UnitCount:     Ptr(1),
			ID:            Ptr(int64(1)),
			Plan: &MarketplacePlan{
				URL:                 Ptr("u"),
				AccountsURL:         Ptr("au"),
				ID:                  Ptr(int64(1)),
				Number:              Ptr(1),
				Name:                Ptr("n"),
				Description:         Ptr("d"),
				MonthlyPriceInCents: Ptr(1),
				YearlyPriceInCents:  Ptr(1),
				PriceModel:          Ptr("pm"),
				UnitName:            Ptr("un"),
				Bullets:             &[]string{"b"},
				State:               Ptr("s"),
				HasFreeTrial:        Ptr(false),
			},
		},
	}

	want := `{
		"url": "u",
		"type": "t",
		"id": 1,
		"login": "l",
		"organization_billing_email": "obe",
		"marketplace_purchase": {
			"billing_cycle": "bc",
			"next_billing_date": ` + referenceTimeStr + `,
			"unit_count": 1,
			"plan": {
				"url": "u",
				"accounts_url": "au",
				"id": 1,
				"number": 1,
				"name": "n",
				"description": "d",
				"monthly_price_in_cents": 1,
				"yearly_price_in_cents": 1,
				"price_model": "pm",
				"unit_name": "un",
				"bullets": ["b"],
				"state": "s",
				"has_free_trial": false
				},
			"on_free_trial": false,
			"free_trial_ends_on": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `
		},
		"marketplace_pending_change": {
			"effective_date": ` + referenceTimeStr + `,
			"unit_count": 1,
			"id": 1,
			"plan": {
				"url": "u",
				"accounts_url": "au",
				"id": 1,
				"number": 1,
				"name": "n",
				"description": "d",
				"monthly_price_in_cents": 1,
				"yearly_price_in_cents": 1,
				"price_model": "pm",
				"unit_name": "un",
				"bullets": ["b"],
				"state": "s",
				"has_free_trial": false
			}
		}
	}`

	testJSONMarshal(t, u, want)
}
