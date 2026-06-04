// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// AICreditTimePeriod represents the billing time period for AI credit usage.
type AICreditTimePeriod struct {
	Year  int  `json:"year"`
	Month *int `json:"month,omitempty"`
	Day   *int `json:"day,omitempty"`
}

// AICreditUsageItem represents a single line item in the AI credit usage report.
type AICreditUsageItem struct {
	Product          string  `json:"product"`
	SKU              string  `json:"sku"`
	Model            string  `json:"model"`
	UnitType         string  `json:"unitType"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	GrossQuantity    float64 `json:"grossQuantity"`
	GrossAmount      float64 `json:"grossAmount"`
	DiscountQuantity float64 `json:"discountQuantity"`
	DiscountAmount   float64 `json:"discountAmount"`
	NetQuantity      float64 `json:"netQuantity"`
	NetAmount        float64 `json:"netAmount"`
}

// AICreditUsage represents the AI credit billing usage for an organization.
type AICreditUsage struct {
	TimePeriod   AICreditTimePeriod   `json:"timePeriod"`
	Organization string               `json:"organization"`
	UsageItems   []*AICreditUsageItem `json:"usageItems"`
}

// AICreditUsageOptions specifies the optional parameters to the
// BillingService.GetAICreditUsage endpoint.
type AICreditUsageOptions struct {
	Year    int    `url:"year,omitempty"`
	Month   int    `url:"month,omitempty"`
	Day     int    `url:"day,omitempty"`
	User    string `url:"user,omitempty"`
	Model   string `url:"model,omitempty"`
	Product string `url:"product,omitempty"`
}

// GetAICreditUsage returns the AI credit billing usage for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/usage?apiVersion=2022-11-28#get-billing-ai-credit-usage-report-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/ai_credit/usage
func (s *BillingService) GetAICreditUsage(ctx context.Context, org string, opts *AICreditUsageOptions) (*AICreditUsage, *Response, error) {
	u := fmt.Sprintf("/organizations/%v/settings/billing/ai_credit/usage", org)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result *AICreditUsage
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
