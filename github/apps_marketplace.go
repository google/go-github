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
	URL                 *string   `json:"url,omitempty"`
	AccountsURL         *string   `json:"accounts_url,omitempty"`
	ID                  *int      `json:"id,omitempty"`
	Name                *string   `json:"name,omitempty"`
	Description         *string   `json:"description,omitempty"`
	MonthlyPriceInCents *int      `json:"monthly_price_in_cents,omitempty"`
	YearlyPriceInCents  *int      `json:"yearly_price_in_cents,omitempty"`
	PriceModel          *string   `json:"price_model,omitempty"`
	UnitName            *string   `json:"unit_name,omitempty"`
	Bullets             *[]string `json:"bullets,omitempty"`
}

// MarketplacePurchase represents a GitHub Apps Marketplace Purchase.
type MarketplacePurchase struct {
	BillingCycle           *string                 `json:"billing_cycle,omitempty"`
	NextBillingDate        *string                 `json:"next_billing_date,omitempty"`
	UnitCount              *int                    `json:"unit_count,omitempty"`
	Plan                   *MarketplacePlan        `json:"plan,omitempty"`
	MarketplacePlanAccount *MarketplacePlanAccount `json:"account,omitempty"`
}

// MarketplacePlanAccount represents a GitHub Account (user or organization) on a specific plan.
type MarketplacePlanAccount struct {
	URL                      *string              `json:"url,omitempty"`
	AccountType              *string              `json:"type,omitempty"`
	ID                       *int                 `json:"id,omitempty"`
	Login                    *string              `json:"login,omitempty"`
	Email                    *string              `json:"email,omitempty"`
	OrganizationBillingEmail *string              `json:"organization_billing_email,omitempty"`
	MarketplacePurchase      *MarketplacePurchase `json:"marketplace_purchase,omitempty"`
}

// ListPlans lists all plans for your Marketplace listing.
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

// ListPlanAccountsForPlan lists all GitHub accounts (user or organization) on a specific plan.
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/#list-all-github-accounts-user-or-organization-on-a-specific-plan
func (s *MarketplaceService) ListPlanAccountsForPlan(ctx context.Context, planID int, stubbed bool, opt *ListOptions) ([]*MarketplacePlanAccount, *Response, error) {
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

// ListPlanAccountsForAccount lists all GitHub accounts (user or organization) associated with an account.
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/#check-if-a-github-account-is-associated-with-any-marketplace-listing
func (s *MarketplaceService) ListPlanAccountsForAccount(ctx context.Context, accountID int, stubbed bool, opt *ListOptions) ([]*MarketplacePlanAccount, *Response, error) {
	u, err := addOptions(fmt.Sprintf("%v/accounts/%v", marketplaceURI(stubbed), accountID), opt)
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

// ListMarketplacePurchasesForUser lists all GitHub marketplace purchases made by a user.
//
// GitHub API docs: https://developer.github.com/v3/apps/marketplace/#get-a-users-marketplace-purchases
func (s *MarketplaceService) ListMarketplacePurchasesForUser(ctx context.Context, stubbed bool, opt *ListOptions) ([]*MarketplacePurchase, *Response, error) {
	u, err := addOptions("apps/user/marketplace_purchases", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeMarketplacePreview)

	var i []*MarketplacePurchase
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
