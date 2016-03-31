// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
)

// OAuthService handles communication with the OAuth related
// methods of the GitHub API.
// You can only access this API via Basic Authentication using your username and password, not tokens.
//
// GitHub API docs: http://developer.github.com/v3/oauth_authorizations/
type OAuthService struct {
	client *Client
}

// Authorization represents a GitHub authorized infomaration
type Authorization struct {
	ID             *int              `json:"id,omitempty"`
	URL            *string           `json:"url,omitempty"`
	Scopes         *[]string         `json:"scopes,omitempty"`
	Token          *string           `json:"token,omitempty"`
	TokenLastEight *string           `json:"token_last_eight,omitempty"`
	hashedToken    *string           `json:"hashed_token,omitempty"`
	App            *AuthorizationApp `json:"app,omitempty"`
	Note           *string           `json:"note,omitempty"`
	NoteURL        *string           `json:"note_url,omitempty"`
	UpdateAt       *Timestamp        `json:"updated_at,omitempty"`
	CreatedAt      *Timestamp        `json:"created_at,omitempty"`
	Fingerprint    *string           `json:"fingerprint,omitempty"`
}

// AuthorizationApp represents application name that github allows app to use api
type AuthorizationApp struct {
	URL      *string `json:"url,omitempty"`
	Name     *string `json:"name,omitempty"`
	ClientID *string `json:"client_id,omitempty"`
}

// AuthorizationRequest represents a request to Get/GetOrCreate an auth.
// It is separate from Authorization above because otherwise Labels
// and Assignee fail to serialize to the correct JSON.
type AuthorizationRequest struct {
	Scopes       *[]string `json:"scopes,omitempty"`
	Note         *string   `json:"note,omitempty"`
	NoteURL      *string   `json:"note_url,omitempty"`
	ClientID     *string   `json:"client_id,omitempty"`
	ClientSecret *string   `json:"client_secret,omitempty"`
	Fingerprint  *string   `json:"fingerprint,omitempty"`
}

// AuthorizationRequest represents a request to Update an auth.
// It is separate from Authorization above because otherwise Labels
// and Assignee fail to serialize to the correct JSON.
type AuthorizationUpdateRequest struct {
	Scopes       *[]string `json:"scopes,omitempty"`
	AddScopes    *[]string `json:"add_scopes,omitempty"`
	RemoveScopes *[]string `json:"remove_scopes,omitempty"`
	Note         *string   `json:"note,omitempty"`
	NoteURL      *string   `json:"note_url,omitempty"`
	Fingerprint  *string   `json:"fingerprint,omitempty"`
}

// ListAuthorizations lists the authorizations info the specified OAuth application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#list-your-authorizations
func (o *OAuthService) ListAuthorizations(opt *ListOptions) ([]Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	auths := new([]Authorization)
	resp, err := o.client.Do(req, auths)
	if err != nil {
		return nil, resp, err
	}
	return *auths, resp, err
}

// Get a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
func (o *OAuthService) Get(id int) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Create a new authorization for the specified OAuth application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
func (o *OAuthService) Create(auth *AuthorizationRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations")

	req, err := o.client.NewRequest("POST", u, auth)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Create a new authorization for the specified OAuth application, only if an authorization for
// that application doesn’t already exist for the user.
// client_id is the 20 character OAuth app client key for which to create the token.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app
func (o *OAuthService) GetOrCreate(client_id string, auth *AuthorizationRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/clients/%s", client_id)

	req, err := o.client.NewRequest("PUT", u, auth)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Create a new authorization for the specified OAuth application, only if an authorization for
// that application and fingerprint do not already exist for the user.
// client_id is the 20 character OAuth app client key for which to create the token.
// fingerprint is a unique string to distinguish an authorization from others created for the same client ID and user.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app-and-fingerprint
func (o *OAuthService) GetOrCreateFingerprint(client_id string, fingerprint string, auth *AuthorizationRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/clients/%s/%s", client_id, fingerprint)

	req, err := o.client.NewRequest("PUT", u, auth)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Update a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
func (o *OAuthService) Update(id int, auth *AuthorizationUpdateRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := o.client.NewRequest("PATCH", u, auth)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Delete a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#delete-an-authorization
func (o *OAuthService) Delete(id int) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := o.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Check a single authorization.
// OAuth applications can use a special API method for checking OAuth token validity without running afoul of
// normal rate limits for failed login attempts. Authentication works differently with this particular endpoint.
// You must use Basic Authentication when accessing it, where the username is the OAuth application client_id and
// the password is its client_secret.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#check-an-authorization
func (o *OAuthService) Check(client_id string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", client_id, token)

	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Reset a single authorization.
// OAuth applications can use this API method to reset a valid OAuth token without end user involvement.
// Applications must save the “token” property in the response, because changes take effect immediately.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#reset-an-authorization
func (o *OAuthService) Reset(client_id string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", client_id, token)

	req, err := o.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Revoke all authrizations for an application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#revoke-all-authorizations-for-an-application
func (o *OAuthService) RevokeAll(client_id string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens", client_id)

	req, err := o.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Revoke an authrization for an application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#revoke-an-authorization-for-an-application
func (o *OAuthService) Revoke(client_id string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", client_id, token)

	req, err := o.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}
	a := new(Authorization)
	resp, err := o.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}
