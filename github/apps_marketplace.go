// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// MarketplaceService handles communication with the marketplace related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/
type MarketplaceService service

// MarketplacePlan represents a GitHub Apps Marketplace Listing Plan.
type MarketplacePlan struct {
	URL                 *string   `json:"url"`
	AccountsURL         *string   `json:"accounts_url"`
	ID                  *int      `json:"id"`
	Name                *string   `json:"name"`
	Description         *string   `json:"description"`
	MonthlyPriceInCents *int      `json:"monthly_price_in_cents"`
	YearlyPriceInCents  *int      `json:"yearly_price_in_cents"`
	PriceModel          *string   `json:"price_model"`
	UnitName            *string   `json:"unit_name"`
	Bullets             *[]string `json:"bullets"`
}

// MarketplacePurchase represents a GitHub Apps Marketplace Purchase.
type MarketplacePurchase struct {
	BillingCycle    *string          `json:"billing_cycle"`
	NextBillingDate *string          `json:"next_billing_date"`
	UnitCount       *int             `json:"unit_count"`
	Plan            *MarketplacePlan `json:"plan"`
}

// MarketplacePlanAccount represents a GitHub Account (user or organization) on a specific plan
type MarketplacePlanAccount struct {
	URL                      *string              `json:"url"`
	AccountType              *string              `json:"type"`
	ID                       *int                 `json:"id"`
	Login                    *string              `json:"login"`
	Email                    *string              `json:"email"`
	OrganizationBillingEmail *string              `json:"organization_billing_email"`
	MarketplacePurchase      *MarketplacePurchase `json:"marketplace_purchase"`
}

// ListPlans lists all plans for your Marketplace listing
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/#list-all-plans-for-your-marketplace-listing
func (s *MarketplaceService) ListPlans(ctx context.Context, stubbed bool, opt *ListOptions) ([]*MarketplacePlan, *Response, error) {
	u, err := addOptions(fmt.Sprintf("%v/plans", marketplaceURI(stubbed)), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeMarketplacePreview)

	var i []*MarketplacePlan
	resp, err := s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}

// ListPlanAccounts lists all GitHub accounts (user or organization) on a specific plan
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/#list-all-github-accounts-user-or-organization-on-a-specific-plan
func (s *MarketplaceService) ListPlanAccounts(ctx context.Context, planID int, stubbed bool, opt *ListOptions) ([]*MarketplacePlanAccount, *Response, error) {
	u, err := addOptions(fmt.Sprintf("%v/plans/%v/accounts", marketplaceURI(stubbed), planID), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeMarketplacePreview)

	var i []*MarketplacePlanAccount
	resp, err := s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}

func marketplaceURI(stubbed bool) string {
	if stubbed {
		return "apps/marketplace_listing/stubbed"
	}
	return "apps/marketplace_listing"
}
