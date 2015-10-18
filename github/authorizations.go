// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// AuthorizationsService handles communication with the OAuth related
// methods of the GitHub API.
//
// This service requires HTTP Basic Authentication; it cannot be accessed using
// an OAuth token.
//
// GitHub API docs: http://developer.github.com/v3/oauth_authorizations/
type AuthorizationsService struct {
	client *Client
}

// Authorization represents a GitHub authorized infomaration.
type Authorization struct {
	ID             *int              `json:"id,omitempty"`
	URL            *string           `json:"url,omitempty"`
	Scopes         []string          `json:"scopes,omitempty"`
	Token          *string           `json:"token,omitempty"`
	TokenLastEight *string           `json:"token_last_eight,omitempty"`
	HashedToken    *string           `json:"hashed_token,omitempty"`
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

// AuthorizationRequest represents a request to create an authorization.
type AuthorizationRequest struct {
	Scopes       []string `json:"scopes,omitempty"`
	Note         *string  `json:"note,omitempty"`
	NoteURL      *string  `json:"note_url,omitempty"`
	ClientID     *string  `json:"client_id,omitempty"`
	ClientSecret *string  `json:"client_secret,omitempty"`
	Fingerprint  *string  `json:"fingerprint,omitempty"`
}

// AuthorizationUpdateRequest represents a request to update an authorization.
type AuthorizationUpdateRequest struct {
	Scopes       []string `json:"scopes,omitempty"`
	AddScopes    []string `json:"add_scopes,omitempty"`
	RemoveScopes []string `json:"remove_scopes,omitempty"`
	Note         *string  `json:"note,omitempty"`
	NoteURL      *string  `json:"note_url,omitempty"`
	Fingerprint  *string  `json:"fingerprint,omitempty"`
}

// List the authorizations info the specified OAuth application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#list-your-authorizations
func (s *AuthorizationsService) List(opt *ListOptions) ([]Authorization, *Response, error) {
	u := "authorizations"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	auths := new([]Authorization)
	resp, err := s.client.Do(req, auths)
	if err != nil {
		return nil, resp, err
	}
	return *auths, resp, err
}

// Get a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
func (s *AuthorizationsService) Get(id int) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Create a new authorization for the specified OAuth application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
func (s *AuthorizationsService) Create(auth *AuthorizationRequest) (*Authorization, *Response, error) {
	u := "authorizations"

	req, err := s.client.NewRequest("POST", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// GetOrCreateForApp creates a new authorization for the specified OAuth
// application, only if an authorization for that application doesn’t already
// exist for the user.
//
// clientID is the OAuth Client ID with which to create the token.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app
func (s *AuthorizationsService) GetOrCreateForApp(clientID string, auth *AuthorizationRequest) (*Authorization, *Response, error) {
	var u string
	if auth.Fingerprint == nil || *auth.Fingerprint == "" {
		u = fmt.Sprintf("authorizations/clients/%v", clientID)
	} else {
		u = fmt.Sprintf("authorizations/clients/%v/%v", clientID, *auth.Fingerprint)
	}

	req, err := s.client.NewRequest("PUT", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Edit a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
func (s *AuthorizationsService) Edit(id int, auth *AuthorizationUpdateRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("PATCH", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Delete a single authorization.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#delete-an-authorization
func (s *AuthorizationsService) Delete(id int) (*Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Check a single authorization.
//
// OAuth applications can use this API method for checking OAuth token validity
// without running afoul of normal rate limits for failed login attempts.
// Authentication works differently with this particular endpoint.  You must
// use Basic Authentication when accessing it, where the username is the OAuth
// application clientID and the password is its client_secret.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#check-an-authorization
func (s *AuthorizationsService) Check(clientID string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", clientID, token)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Reset a single authorization.
//
// OAuth applications can use this API method to reset a valid OAuth token
// without end user involvement.  Applications must save the "token" property
// in the response, because changes take effect immediately.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#reset-an-authorization
func (s *AuthorizationsService) Reset(clientID string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", clientID, token)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Revoke an authorization for an application.
//
// GitHub API Docs: https://developer.github.com/v3/oauth_authorizations/#revoke-an-authorization-for-an-application
func (s *AuthorizationsService) Revoke(clientID string, token string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/tokens/%s", clientID, token)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
