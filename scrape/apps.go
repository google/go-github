// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// apps.go contains functions for accessing data about applications installed
// on a GitHub organization.

package scrape

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// AppRestrictionsEnabled returns whether the specified organization has
// restricted third-party application access.
func (c *Client) AppRestrictionsEnabled(org string) (bool, error) {
	doc, err := c.get("/organizations/%s/settings/oauth_application_policy", org)
	if err != nil {
		return false, err
	}

	s := doc.Find(".oauth-application-whitelist svg").First()
	if s.Length() == 0 {
		return false, errors.New("unable to find expected markup")
	}

	if s.HasClass("octicon-check") {
		return true, nil
	}
	if s.HasClass("octicon-alert") {
		return false, nil
	}

	return false, errors.New("unable to find expected markup")
}

// ListOAuthApps lists the reviewed OAuth Applications for the
// specified organization (whether approved or denied).
func (c *Client) ListOAuthApps(org string) ([]OAuthApp, error) {
	doc, err := c.get("/organizations/%s/settings/oauth_application_policy", org)
	if err != nil {
		return nil, err
	}

	var apps []OAuthApp
	doc.Find(".oauth-application-whitelist ul > li").Each(func(i int, s *goquery.Selection) {
		var app OAuthApp
		app.Name = s.Find(".request-info strong").First().Text()
		app.Description = s.Find(".request-info .application-description").Text()

		if editURL, ok := s.Find(".request-indicator a.edit-link").Attr("href"); ok {
			app.ID = intFromLastPathSegment(editURL)
		}

		if r := s.Find(".request-indicator .requestor"); r.Length() > 0 {
			app.State = OAuthAppRequested
			app.RequestedBy = r.Text()
			if editURL, ok := s.Find(".request-indicator a").Last().Attr("href"); ok {
				app.ID = intFromLastPathSegment(editURL)
			}
		} else if r := s.Find(".request-indicator .approved-request"); r.Length() > 0 {
			app.State = OAuthAppApproved
		} else if r := s.Find(".request-indicator .denied-request"); r.Length() > 0 {
			app.State = OAuthAppDenied
		}
		apps = append(apps, app)
	})

	return apps, nil
}

func intFromLastPathSegment(s string) int {
	seg := strings.Split(s, "/")
	if len(seg) > 0 {
		i, _ := strconv.Atoi(seg[len(seg)-1])
		return i
	}
	return 0
}

// OAuthAppReviewState indicates the current state of a requested OAuth Application.
type OAuthAppReviewState int

const (
	// OAuthAppRequested indicates access has been requested, but not reviewed
	OAuthAppRequested OAuthAppReviewState = iota + 1
	// OAuthAppApproved indicates access has been approved
	OAuthAppApproved
	// OAuthAppDenied indicates access has been denied
	OAuthAppDenied
)

// OAuthApp represents an OAuth application that has been reviewed for access to organization data.
type OAuthApp struct {
	ID          int
	Name        string
	Description string
	State       OAuthAppReviewState
	RequestedBy string
}
