// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// apps.go contains functions for accessing data about applications installed
// on a GitHub organization.

package scrape

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// OrgPaymentInformation returns payment information for the specified org.
func (c *Client) OrgPaymentInformation(org string) (PaymentInformation, error) {
	var info PaymentInformation

	doc, err := c.get("/organizations/%s/settings/billing/payment_information", org)
	if err != nil {
		return info, err
	}

	doc.Find("main h4.mb-1").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(strings.ToLower(s.Text()))
		value := strings.Join(strings.Fields(strings.TrimSpace(s.NextFiltered("p").Text())), " ")

		switch name {
		case "payment method":
			info.PaymentMethod = value
		case "last payment":
			info.LastPayment = value
		case "coupon":
			info.Coupon = value
		case "extra information":
			info.ExtraInformation = value
		}
	})

	return info, nil
}

// PaymentInformation for an organization on a paid plan.
type PaymentInformation struct {
	PaymentMethod    string
	LastPayment      string
	Coupon           string
	ExtraInformation string
}
