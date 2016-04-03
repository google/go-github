// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
)

// Scope models a GitHub authorization scope.
//
// GitHub API docs:https://developer.github.com/v3/oauth/#scopes
type Scope string

// This is the set of scopes for GitHub API V3
const (
	ScopeNone           Scope = "(no scope)" // REVISIT: is this actually returned, or just a documentation artifact?
	ScopeUser           Scope = "user"
	ScopeUserEmail      Scope = "user:email"
	ScopeUserFollow     Scope = "user:follow"
	ScopePublicRepo     Scope = "public_repo"
	ScopeRepo           Scope = "repo"
	ScopeRepoDeployment Scope = "repo_deployment"
	ScopeRepoStatus     Scope = "repo:status"
	ScopeDeleteRepo     Scope = "delete_repo"
	ScopeNotifications  Scope = "notifications"
	ScopeGist           Scope = "gist"
	ScopeReadRepoHook   Scope = "read:repo_hook"
	ScopeWriteRepoHook  Scope = "write:repo_hook"
	ScopeAdminRepoHook  Scope = "admin:repo_hook"
	ScopeAdminOrgHook   Scope = "admin:org_hook"
	ScopeReadOrg        Scope = "read:org"
	ScopeWriteOrg       Scope = "write:org"
	ScopeAdminOrg       Scope = "admin:org"
	ScopeReadPublicKey  Scope = "read:public_key"
	ScopeWritePublicKey Scope = "write:public_key"
	ScopeAdminPublicKey Scope = "admin:public_key"
)

// AuthorizationsService handles communication with the authorization related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/
type AuthorizationsService struct {
	client *Client
}

// Authorization models an individual GitHub authorization.
// Note that the User field is only present for the CheckAppAuthorization
// and ResetAppAuthorization operations.
type Authorization struct {
	ID             *int    `json:"id,omitempty"`
	URL            *string `json:"url,omitempty"`
	App            *App    `json:"app,omitempty"`
	Token          *string `json:"token,omitempty"`
	HashedToken    *string `json:"hashed_token,omitempty"`
	TokenLastEight *string `json:"token_last_eight,omitempty"`
	Note           *string `json:"note,omitempty"`
	NoteURL        *string `json:"note_url,omitempty"`
	CreatedAt      *string `json:"created_at,omitempty"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
	Scopes         []Scope `json:"scopes,omitempty"`
	Fingerprint    *string `json:"fingerprint,omitempty"`
	User           *User   `json:"user,omitempty"`
}

func (a Authorization) String() string {
	return Stringify(a)
}

// App models an individual GitHub app (in the context of authorization).
type App struct {
	Name     *string `json:"name,omitempty"`
	URL      *string `json:"url,omitempty"`
	ClientID *string `json:"client_id,omitempty"`
}

func (a App) String() string {
	return Stringify(a)
}

// AuthorizationRequest is used to create a new GitHub authorization. Note that not all fields
// are used for any one operation.
type AuthorizationRequest struct {
	ClientID     *string `json:"client_id,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
	Scopes       []Scope `json:"scopes,omitempty"`
	Note         *string `json:"note,omitempty"`
	NoteURL      *string `json:"note_url,omitempty"`
	Fingerprint  *string `json:"fingerprint,omitempty"`
}

func (a AuthorizationRequest) String() string {
	return Stringify(a)
}

// AuthorizationUpdate is used to update an existing GitHub authorization.
// Note that for any one update, you must only provide one of the "scopes" fields.
// That is, you may provide only one of "Scopes", or "AddScopes", or "RemoveScopes".
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
type AuthorizationUpdate struct {
	// ID is not serialized as the update operations take the ID in the path
	ID           *int    `json:"-"`
	Scopes       []Scope `json:"scopes,omitempty"`
	AddScopes    []Scope `json:"add_scopes,omitempty"`
	RemoveScopes []Scope `json:"remove_scopes,omitempty"`
	Note         *string `json:"note,omitempty"`
	NoteURL      *string `json:"note_url,omitempty"`
	Fingerprint  *string `json:"fingerprint,omitempty"`
}

func (a AuthorizationUpdate) String() string {
	return Stringify(a)
}

// List lists your GitHub authorizations.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#list-your-authorizations
func (s *AuthorizationsService) List() ([]Authorization, *Response, error) {

	// GET /authorizations
	req, err := s.client.NewRequest("GET", "authorizations", nil)
	if err != nil {
		return nil, nil, err
	}

	auths := new([]Authorization)
	resp, err := s.client.Do(req, auths)

	return *auths, resp, err
}

// Get retrieves a single GitHub authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
func (s *AuthorizationsService) Get(id int) (*Authorization, *Response, error) {

	if id == 0 {
		return nil, nil, errors.New("You must provide an id parameter.")
	}

	// GET /authorizations/:id
	u := fmt.Sprintf("authorizations/%v", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	auth := new(Authorization)
	resp, err := s.client.Do(req, auth)

	return auth, resp, err
}

// Create creates a new GitHub authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
func (s *AuthorizationsService) Create(auth *AuthorizationRequest) (*Authorization, *Response, error) {

	// POST /authorizations
	req, err := s.client.NewRequest("POST", "authorizations", auth)
	if err != nil {
		return nil, nil, err
	}

	authResponse := new(Authorization)
	resp, err := s.client.Do(req, authResponse)

	return authResponse, resp, err
}

// GetOrCreateForApp will create a new GitHub authorization for the specified OAuth application,
// (or optionally app/fingerprint combination) only if an authorization for that app/fingerprint doesn't already exist for that user. If a new token is
// created, the HTTP status code will be "201 Created", and the returned Authorization.Token field
// will be populated. If an existing token is returned, the status code will be "200 OK" and the
// Authorization.Token field will be empty.
//
// GitHub API docs:
// - https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app
// - https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app-and-fingerprint
func (s *AuthorizationsService) GetOrCreateForApp(auth *AuthorizationRequest) (*Authorization, *Response, error) {

	if *auth.ClientID == "" || *auth.ClientSecret == "" {
		return nil, nil, errors.New("You must provide the ClientID and ClientSecret parameters (as part of AuthorizationRequest)")
	}

	// PUT /authorizations/clients/:client_id
	u := fmt.Sprintf("authorizations/clients/%v", *auth.ClientID)

	// We want to set the ClientID to nil as this operation does not expect it in the body
	// But we also don't want to mess up the struct we received, so we make a copy
	authCopy := *auth
	authCopy.ClientID = nil

	req, err := s.client.NewRequest("PUT", u, authCopy)
	if err != nil {
		return nil, nil, err
	}

	authResponse := new(Authorization)
	resp, err := s.client.Do(req, authResponse)

	return authResponse, resp, err
}

// Edit updates an existing GitHub authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
func (s *AuthorizationsService) Edit(auth *AuthorizationUpdate) (*Authorization, *Response, error) {

	// PATCH /authorizations/:id
	u := fmt.Sprintf("authorizations/%v", *auth.ID)
	req, err := s.client.NewRequest("PATCH", u, auth)
	if err != nil {
		return nil, nil, err
	}

	authResponse := new(Authorization)
	resp, err := s.client.Do(req, authResponse)

	return authResponse, resp, err
}

// Delete deletes a GitHub authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#delete-an-authorization
func (s *AuthorizationsService) Delete(id int) (*Response, error) {

	if id == 0 {
		return nil, errors.New("You must provide an id parameter.")
	}
	// DELETE /authorizations/:id
	u := fmt.Sprintf("authorizations/%v", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// CheckAppAuthorization checks if an OAuth token is valid for a specific app.
// Note that this operation requires the use of BasicAuth, but where the username
// is the OAuth application clientID, and the password is its clientSecret. Invalid
// tokens will return a 404 Not Found.
//
// Note that for this operation, the returned Authorization.User field will be populated.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#check-an-authorization
func (s *AuthorizationsService) CheckAppAuthorization(clientID string, token string) (*Authorization, *Response, error) {

	// GET /applications/:client_id/tokens/:access_token
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	auth := new(Authorization)
	resp, err := s.client.Do(req, auth)

	return auth, resp, err
}

// ResetAppAuthorization is used to reset a valid OAuth token without end user involvement.
// Note that this operation requires the use of BasicAuth, but where the username
// is the OAuth application clientID, and the password is its clientSecret. Invalid
// tokens will return a 404 Not Found.
//
// Note that for this operation, the returned Authorization.User field will be populated.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#reset-an-authorization
func (s *AuthorizationsService) ResetAppAuthorization(clientID string, token string) (*Authorization, *Response, error) {

	// POST /applications/:client_id/tokens/:access_token
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)
	req, err := s.client.NewRequest("POST", u, nil)

	if err != nil {
		return nil, nil, err
	}

	auth := new(Authorization)
	resp, err := s.client.Do(req, auth)

	return auth, resp, err
}

// RevokeAppAuthorization is used to revoke a single token for an OAuth application.
// Note that this operation requires the use of BasicAuth, but where the username
// is the OAuth application clientID, and the password is its clientSecret. Invalid
// tokens will return a 404 Not Found.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#revoke-an-authorization-for-an-application
func (s *AuthorizationsService) RevokeAppAuthorization(clientID string, token string) (*Response, error) {

	// DELETE /applications/:client_id/tokens/:access_token
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)
	req, err := s.client.NewRequest("DELETE", u, nil)

	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
