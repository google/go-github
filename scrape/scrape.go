// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scrape provides a client for interacting with GitHub using screen
// scraping.
//
// This is intended to be used as a supplement to the standard go-github
// library to access data that is not currently exposed by either the official
// REST or GraphQL APIs.
//
// Because of the nature of screen scraping, this package should be treated as
// HIGHLY EXPERIMENTAL, and potentially unstable. We make no promises relating
// to compatibility or stability of the exported API. Even though this package
// is distributed as part of the go-github library, it is explicitly exempt
// from any stability promises that my be implied by the library version
// number.
package scrape

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xlzd/gotp"
	"golang.org/x/net/publicsuffix"
)

var defaultBaseURL = "https://github.com/"

// Client is a GitHub scraping client.
type Client struct {
	*http.Client

	// base URL for github.com pages.  Exposed primarily for testing.  Also
	// used for saving and restoring cookies on the Client.
	baseURL *url.URL
}

// NewClient constructs a new Client.  If transport is nil, a default transport is used.
func NewClient(transport http.RoundTripper) *Client {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatalf("error creating cookiejar: %v", err)
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		Client: &http.Client{
			Transport: transport,
			Jar:       jar,
		},
		baseURL: baseURL,
	}
}

// SaveCookies returns an encoded form of the github.com cookies set on this
// client.  If Authenticate() has been called, this should include the
// github.com session cookie.  These cookies can be loaded onto a new client by
// calling LoadCookies.
//
// GitHub session cookies are bearer tokens that are not tied to any particular
// client, so should be treated with the same sensitivity as the account
// credentials.
func (c *Client) SaveCookies() ([]byte, error) {
	cookies := c.Client.Jar.Cookies(c.baseURL)

	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(cookies)
	return b.Bytes(), err
}

// LoadCookies loads the provided cookies for github.com.
func (c *Client) LoadCookies(v []byte) error {
	var cookies []*http.Cookie
	r := bytes.NewReader(v)
	err := gob.NewDecoder(r).Decode(&cookies)
	if err != nil {
		return err
	}

	c.Client.Jar.SetCookies(c.baseURL, cookies)
	return nil
}

// get fetches a urlStr (a GitHub URL relative to the client's baseURL), and
// returns the parsed response document.
func (c *Client) get(urlStr string, a ...interface{}) (*goquery.Document, error) {
	u, err := c.baseURL.Parse(fmt.Sprintf(urlStr, a...))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %q: %v", urlStr, err)
	}
	resp, err := c.Client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching url %q: %v", u, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("received %v response fetching URL %q", resp.StatusCode, u)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return doc, nil
}

// Authenticate client to GitHub with the provided username, password, and if
// two-factor auth is enabled for the account, otpseed.
//
// otpseed is the OTP Secret provided from GitHub as part of two-factor
// application enrollment.  When registering the application, click the "enter
// this text code" link on the QR Code page to see the raw OTP Secret.
func (c *Client) Authenticate(username, password, otpseed string) error {
	setPassword := func(values url.Values) {
		values.Set("login", username)
		values.Set("password", password)
	}
	resp, err := fetchAndSubmitForm(c.Client, "https://github.com/login", setPassword)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v response submitting login form", resp.StatusCode)
	}

	if otpseed == "" {
		return nil
	}

	setOTP := func(values url.Values) {
		otp := gotp.NewDefaultTOTP(strings.ToUpper(otpseed)).Now()
		values.Set("otp", otp)
	}
	resp, err = fetchAndSubmitForm(c.Client, "https://github.com/sessions/two-factor", setOTP)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v response submitting otp form", resp.StatusCode)
	}

	return nil
}
