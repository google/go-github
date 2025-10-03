// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// BillingService provides access to the billing related functions
// in the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/billing
type BillingService service

// MinutesUsedBreakdown counts the actions minutes used by machine type (e.g. UBUNTU, WINDOWS, MACOS).
type MinutesUsedBreakdown = map[string]int

// PackageBilling represents a GitHub Package billing.
type PackageBilling struct {
	TotalGigabytesBandwidthUsed     int     `json:"total_gigabytes_bandwidth_used"`
	TotalPaidGigabytesBandwidthUsed int     `json:"total_paid_gigabytes_bandwidth_used"`
	IncludedGigabytesBandwidth      float64 `json:"included_gigabytes_bandwidth"`
}

// StorageBilling represents a GitHub Storage billing.
type StorageBilling struct {
	DaysLeftInBillingCycle       int     `json:"days_left_in_billing_cycle"`
	EstimatedPaidStorageForMonth float64 `json:"estimated_paid_storage_for_month"`
	EstimatedStorageForMonth     float64 `json:"estimated_storage_for_month"`
}

// ActiveCommitters represents the total active committers across all repositories in an Organization.
type ActiveCommitters struct {
	TotalAdvancedSecurityCommitters     int                           `json:"total_advanced_security_committers"`
	TotalCount                          int                           `json:"total_count"`
	MaximumAdvancedSecurityCommitters   int                           `json:"maximum_advanced_security_committers"`
	PurchasedAdvancedSecurityCommitters int                           `json:"purchased_advanced_security_committers"`
	Repositories                        []*RepositoryActiveCommitters `json:"repositories,omitempty"`
}

// RepositoryActiveCommitters represents active committers on each repository.
type RepositoryActiveCommitters struct {
	Name                                *string                                `json:"name,omitempty"`
	AdvancedSecurityCommitters          *int                                   `json:"advanced_security_committers,omitempty"`
	AdvancedSecurityCommittersBreakdown []*AdvancedSecurityCommittersBreakdown `json:"advanced_security_committers_breakdown,omitempty"`
}

// AdvancedSecurityCommittersBreakdown represents the user activity breakdown for ActiveCommitters.
type AdvancedSecurityCommittersBreakdown struct {
	UserLogin      *string `json:"user_login,omitempty"`
	LastPushedDate *string `json:"last_pushed_date,omitempty"`
}

// UsageReportOptions specifies optional parameters for the enhanced billing platform usage report.
type UsageReportOptions struct {
	// If specified, only return results for a single year. The value of year is an integer with four digits representing a year. For example, 2025.
	// Default value is the current year.
	Year *int `url:"year,omitempty"`

	// If specified, only return results for a single month. The value of month is an integer between 1 and 12.
	// If no year is specified the default year is used.
	Month *int `url:"month,omitempty"`

	// If specified, only return results for a single day. The value of day is an integer between 1 and 31.
	// If no year or month is specified, the default year and month are used.
	Day *int `url:"day,omitempty"`

	// If specified, only return results for a single hour. The value of hour is an integer between 0 and 23.
	// If no year, month, or day is specified, the default year, month, and day are used.
	Hour *int `url:"hour,omitempty"`
}

// PremiumRequestUsageReportOptions specifies optional parameters
// for the enhanced billing platform premiun request usage report.
type PremiumRequestUsageReportOptions struct {
	UsageReportOptions

	// The user name to query usage for. The name is not case sensitive.
	User *string `url:"user,omitempty"`

	// The model name to query usage for. The name is not case sensitive.
	Model *string `url:"model,omitempty"`

	// The product name to query usage for. The name is not case sensitive.
	Product *string `url:"product,omitempty"`
}

// UsageItem represents a single usage item in the enhanced billing platform report.
type UsageItem struct {
	Date           *string  `json:"date"`
	Product        *string  `json:"product"`
	SKU            *string  `json:"sku"`
	Quantity       *float64 `json:"quantity"`
	UnitType       *string  `json:"unitType"`
	PricePerUnit   *float64 `json:"pricePerUnit"`
	GrossAmount    *float64 `json:"grossAmount"`
	DiscountAmount *float64 `json:"discountAmount"`
	NetAmount      *float64 `json:"netAmount"`
	RepositoryName *string  `json:"repositoryName,omitempty"`
	// Organization name is only used for organization-level reports.
	OrganizationName *string `json:"organizationName,omitempty"`
}

// UsageReport represents the enhanced billing platform usage report response.
type UsageReport struct {
	UsageItems []*UsageItem `json:"usageItems,omitempty"`
}

// PremiumRequestUsageItem represents a single usage line item in premium request usage reports.
type PremiumRequestUsageItem struct {
	Product          string  `json:"product"`
	SKU              string  `json:"sku"`
	Model            string  `json:"model"`
	UnitType         string  `json:"unitType"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	GrossQuantity    int     `json:"grossQuantity"`
	GrossAmount      float64 `json:"grossAmount"`
	DiscountQuantity int     `json:"discountQuantity"`
	DiscountAmount   float64 `json:"discountAmount"`
	NetQuantity      int     `json:"netQuantity"`
	NetAmount        float64 `json:"netAmount"`
}

// PremiumRequestUsageReport represents the premium request usage report response.
type PremiumRequestUsageReport struct {
	TimePeriod struct {
		Year  int  `json:"year"`
		Month *int `json:"month,omitempty"`
		Day   *int `json:"day,omitempty"`
	} `json:"timePeriod"`
	Organization string                     `json:"organization"`
	User         *string                    `json:"user,omitempty"`
	Product      *string                    `json:"product,omitempty"`
	Model        *string                    `json:"model,omitempty"`
	UsageItems   []*PremiumRequestUsageItem `json:"usageItems"`
}

// GetPackagesBillingOrg returns the free and paid storage used for GitHub Packages in gigabytes for an Org.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-github-packages-billing-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/billing/packages
func (s *BillingService) GetPackagesBillingOrg(ctx context.Context, org string) (*PackageBilling, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/billing/packages", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	packagesOrgBilling := new(PackageBilling)
	resp, err := s.client.Do(ctx, req, packagesOrgBilling)
	if err != nil {
		return nil, resp, err
	}

	return packagesOrgBilling, resp, nil
}

// GetStorageBillingOrg returns the estimated paid and estimated total storage used for GitHub Actions
// and GitHub Packages in gigabytes for an Org.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-shared-storage-billing-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/billing/shared-storage
func (s *BillingService) GetStorageBillingOrg(ctx context.Context, org string) (*StorageBilling, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/billing/shared-storage", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	storageOrgBilling := new(StorageBilling)
	resp, err := s.client.Do(ctx, req, storageOrgBilling)
	if err != nil {
		return nil, resp, err
	}

	return storageOrgBilling, resp, nil
}

// GetAdvancedSecurityActiveCommittersOrg returns the GitHub Advanced Security active committers for an organization per repository.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/billing/billing#get-github-advanced-security-active-committers-for-an-organization
//
//meta:operation GET /orgs/{org}/settings/billing/advanced-security
func (s *BillingService) GetAdvancedSecurityActiveCommittersOrg(ctx context.Context, org string, opts *ListOptions) (*ActiveCommitters, *Response, error) {
	u := fmt.Sprintf("orgs/%v/settings/billing/advanced-security", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	activeOrgCommitters := new(ActiveCommitters)
	resp, err := s.client.Do(ctx, req, activeOrgCommitters)
	if err != nil {
		return nil, resp, err
	}

	return activeOrgCommitters, resp, nil
}

// GetPackagesBillingUser returns the free and paid storage used for GitHub Packages in gigabytes for a user.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-github-packages-billing-for-a-user
//
//meta:operation GET /users/{username}/settings/billing/packages
func (s *BillingService) GetPackagesBillingUser(ctx context.Context, user string) (*PackageBilling, *Response, error) {
	u := fmt.Sprintf("users/%v/settings/billing/packages", user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	packagesUserBilling := new(PackageBilling)
	resp, err := s.client.Do(ctx, req, packagesUserBilling)
	if err != nil {
		return nil, resp, err
	}

	return packagesUserBilling, resp, nil
}

// GetStorageBillingUser returns the estimated paid and estimated total storage used for GitHub Actions
// and GitHub Packages in gigabytes for a user.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-shared-storage-billing-for-a-user
//
//meta:operation GET /users/{username}/settings/billing/shared-storage
func (s *BillingService) GetStorageBillingUser(ctx context.Context, user string) (*StorageBilling, *Response, error) {
	u := fmt.Sprintf("users/%v/settings/billing/shared-storage", user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	storageUserBilling := new(StorageBilling)
	resp, err := s.client.Do(ctx, req, storageUserBilling)
	if err != nil {
		return nil, resp, err
	}

	return storageUserBilling, resp, nil
}

// GetUsageReportOrg returns a report of the total usage for an organization using the enhanced billing platform.
//
// Note: This endpoint is only available to organizations with access to the enhanced billing platform.
//
// GitHub API docs: https://docs.github.com/rest/billing/enhanced-billing#get-billing-usage-report-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/usage
func (s *BillingService) GetUsageReportOrg(ctx context.Context, org string, opts *UsageReportOptions) (*UsageReport, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/usage", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	usageReport := new(UsageReport)
	resp, err := s.client.Do(ctx, req, usageReport)
	if err != nil {
		return nil, resp, err
	}

	return usageReport, resp, nil
}

// GetUsageReportUser returns a report of the total usage for a user using the enhanced billing platform.
//
// Note: This endpoint is only available to users with access to the enhanced billing platform.
//
// GitHub API docs: https://docs.github.com/rest/billing/enhanced-billing#get-billing-usage-report-for-a-user
//
//meta:operation GET /users/{username}/settings/billing/usage
func (s *BillingService) GetUsageReportUser(ctx context.Context, user string, opts *UsageReportOptions) (*UsageReport, *Response, error) {
	u := fmt.Sprintf("users/%v/settings/billing/usage", user)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	usageReport := new(UsageReport)
	resp, err := s.client.Do(ctx, req, usageReport)
	if err != nil {
		return nil, resp, err
	}

	return usageReport, resp, nil
}

// GetPremiumRequestUsageReportOrg returns a report of the premium request
// usage for an organization using the enhanced billing platform.
//
// Note: This endpoint is only available to organizations with access to the enhanced billing platform.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-premium-request-usage-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/premium_request/usage
func (s *BillingService) GetPremiumRequestUsageReportOrg(ctx context.Context, org string, opts *PremiumRequestUsageReportOptions) (*PremiumRequestUsageReport, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/premium_request/usage", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	premiumRequestUsageReport := new(PremiumRequestUsageReport)
	resp, err := s.client.Do(ctx, req, premiumRequestUsageReport)
	if err != nil {
		return nil, resp, err
	}

	return premiumRequestUsageReport, resp, nil
}

// GetPremiumRequestUsageReportUser returns a report of the premium request
// usage for a user using the enhanced billing platform.
//
// Note: This endpoint is only available to users with access to the enhanced billing platform.
//
// GitHub API docs: https://docs.github.com/rest/billing/billing#get-premium-request-usage-for-a-user
//
//meta:operation GET /users/{username}/settings/billing/premium_request/usage
func (s *BillingService) GetPremiumRequestUsageReportUser(ctx context.Context, user string, opts *PremiumRequestUsageReportOptions) (*PremiumRequestUsageReport, *Response, error) {
	u := fmt.Sprintf("users/%v/settings/billing/premium_request/usage", user)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	premiumRequestUsageReport := new(PremiumRequestUsageReport)
	resp, err := s.client.Do(ctx, req, premiumRequestUsageReport)
	if err != nil {
		return nil, resp, err
	}

	return premiumRequestUsageReport, resp, nil
}
