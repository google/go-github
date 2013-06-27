// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github provides a client for using the GitHub API.

Access different parts of the GitHub API using the various services on a GitHub
Client:

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, err := client.Organizations.List("willnorris", nil)

Set optional parameters for an API method by passing an Options object.

	// list recently updated repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Sort: "updated"}
	repos, err := client.Repositories.ListByOrg("github", opt)

Make authenticated API calls by constructing a GitHub client using an OAuth
capable http.Client:

	import "code.google.com/p/goauth2/oauth"

	// simple OAuth transport if you already have an access token;
	// see goauth2 library for full usage
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "..."},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	repos, err := client.Repositories.List(nil)

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

The full GitHub API is documented at http://developer.github.com/v3/.
*/
package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.github.com/"
	userAgent      = "go-github/" + libraryVersion

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
)

// A Client manages communication with the GitHub API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to the public GitHub API, but can be
	// set to a domain endpoint to use with GitHub Enterprise.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the GitHub API.
	UserAgent string

	// Rate specifies the current rate limit for the client as determined by the
	// most recent API call.  If the client is used in a multi-user application,
	// this rate may not always be up-to-date.  Call RateLimit() to check the
	// current rate.
	Rate Rate

	// Services used for talking to different parts of the API

	Issues        *IssuesService
	Organizations *OrganizationsService
	PullRequests  *PullRequestsService
	Repositories  *RepositoriesService
	Users         *UsersService
	Gists         *GistsService
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int
}

// NewClient returns a new GitHub API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Issues = &IssuesService{client: c}
	c.Organizations = &OrganizationsService{client: c}
	c.PullRequests = &PullRequestsService{client: c}
	c.Repositories = &RepositoriesService{client: c}
	c.Users = &UsersService{client: c}
	c.Gists = &GistsService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// update rate limit
	if limit := resp.Header.Get(headerRateLimit); limit != "" {
		c.Rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := resp.Header.Get(headerRateRemaining); remaining != "" {
		c.Rate.Remaining, _ = strconv.Atoi(remaining)
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

/*
An ErrorResponse reports one or more errors caused by an API request.

GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:message` // error message
	Errors   []Error        `json:errors`  // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Errors)
}

/*
An Error reports more details on an individual error in an ErrorResponse.
These are the possible validation error codes:

    missing:
        resource does not exist
    missing_field:
        a required field on a resource has not been set
    invalid:
        the formatting of a field is invalid
    already_exists:
        another resource has the same valid as this field

GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type Error struct {
	Resource string `json:resource` // resource on which the error occurred
	Field    string `json:field`    // field on which the error occurred
	Code     string `json:code`     // validation error code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// parseBoolResponse determines the boolean result from a GitHub API response.
// Several GitHub API methods return boolean responses indicated by the HTTP
// status code in the response (true indicated by a 204, false indicated by a
// 404).  This helper function will determine that result and hide the 404
// error if present.  Any other error will be returned through as-is.
func parseBoolResponse(err error) (bool, error) {
	if err == nil {
		return true, nil
	}

	if err, ok := err.(*ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
		// Simply false.  In this one case, we do not pass the error through.
		return false, nil
	}

	// some other real error occurred
	return false, err
}

// API response wrapper to a rate limit request.
type rateResponse struct {
	Rate *Rate `json:rate`
}

// Rate represents the rate limit for the current client.  Unauthenticated
// requests are limited to 60 per hour.  Authenticated requests are limited to
// 5,000 per hour.
type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:limit`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:remaining`
}

// RateLimit returns the rate limit for the current client.
func (c *Client) RateLimit() (*Rate, error) {
	req, err := c.NewRequest("GET", "rate_limit", nil)
	if err != nil {
		return nil, err
	}

	response := new(rateResponse)
	_, err = c.Do(req, response)
	return response.Rate, err
}

/*
UnauthenticatedRateLimitedTransport allows you to make unauthenticated calls
that need to use a higher rate limit associated with your OAuth application.

	t := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     "your app's client ID",
		ClientSecret: "your app's client secret",
	}
	client := github.NewClient(t.Client())

This will append the querystring params client_id=xxx&client_secret=yyy to all
requests.

See http://developer.github.com/v3/#unauthenticated-rate-limited-requests for
more information.
*/
type UnauthenticatedRateLimitedTransport struct {
	// ClientID is the GitHub OAuth client ID of the current application, which
	// can be found by selecting its entry in the list at
	// https://github.com/settings/applications.
	ClientID string

	// ClientSecret is the GitHub OAuth client secret of the current
	// application.
	ClientSecret string

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *UnauthenticatedRateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ClientID == "" {
		return nil, errors.New("ClientID is empty")
	}
	if t.ClientSecret == "" {
		return nil, errors.New("ClientSecret is empty")
	}

	// To set extra querystring params, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	req = cloneRequest(req)
	q := req.URL.Query()
	q.Set("client_id", t.ClientID)
	q.Set("client_secret", t.ClientSecret)
	req.URL.RawQuery = q.Encode()

	// Make the HTTP request.
	return t.transport().RoundTrip(req)
}

// Client returns an *http.Client that makes requests which are subject to the
// rate limit of your OAuth application.
func (t *UnauthenticatedRateLimitedTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *UnauthenticatedRateLimitedTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request. The clone is a
// shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
