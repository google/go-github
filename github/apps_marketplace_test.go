// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMarketplaceService_ListPlans(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/apps/marketplace_listing/plans", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMarketplacePreview)
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	plans, _, err := client.Marketplace.ListPlans(context.Background(), false, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlans returned error: %v", err)
	}

	want := []*MarketplacePlan{{ID: Int(1)}}
	if !reflect.DeepEqual(plans, want) {
		t.Errorf("Marketplace.ListPlans returned %+v, want %+v", plans, want)
	}
}

func TestMarketplaceService_ListPlansStubbed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/apps/marketplace_listing/stubbed/plans", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMarketplacePreview)
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	plans, _, err := client.Marketplace.ListPlans(context.Background(), true, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlans returned error: %v", err)
	}

	want := []*MarketplacePlan{{ID: Int(1)}}
	if !reflect.DeepEqual(plans, want) {
		t.Errorf("Marketplace.ListPlans returned %+v, want %+v", plans, want)
	}
}

func TestMarketplaceService_ListPlanAccounts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/apps/marketplace_listing/plans/1/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMarketplacePreview)
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	accounts, _, err := client.Marketplace.ListPlanAccounts(context.Background(), 1, false, opt)
	if err != nil {
		t.Errorf("Marketplace.ListPlanAccounts returned error: %v", err)
	}

	want := []*MarketplacePlanAccount{{ID: Int(1)}}
	if !reflect.DeepEqual(accounts, want) {
		t.Errorf("Marketplace.ListPlanAccounts returned %+v, want %+v", accounts, want)
	}
}
