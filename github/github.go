// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

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
		Config: &oauth.Config{},
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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.github.com/"
	userAgent      = "go-github/" + libraryVersion
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

	// Services used for talking to different parts of the API

	Organizations *OrganizationsService
	Repositories  *RepositoriesService
	Users         *UsersService
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
	c.Organizations = &OrganizationsService{client: c}
	c.Repositories = &RepositoriesService{client: c}
	c.Users = &UsersService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urls,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urls string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urls)
	if err != nil {
		return nil, err
	}

	url_ := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url_.String(), buf)
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
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
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
