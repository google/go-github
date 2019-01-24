// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// createUserRequest is a subset of User and is used internally
// by Create to pass only the known fields for the endpoint
//
type createUserRequest struct {
	Login *string `json:"login"`
	Email *string `json:"email"`
}

// CreateUser creates a new user in Github Enterprise.
//
// Github Enterprise API docs: https://developer.github.com/enterprise/2.16/v3/enterprise-admin/users/#create-a-new-user
func (s *AdminService) CreateUser(ctx context.Context, login, email string) (*User, *Response, error) {
	u := "admin/users"

	userReq := &createUserRequest{
		Login: &login,
		Email: &email,
	}

	req, err := s.client.NewRequest("POST", u, userReq)
	if err != nil {
		return nil, nil, err
	}

	var user User
	resp, err := s.client.Do(ctx, req, &user)
	if err != nil {
		return nil, resp, err
	}

	return &user, resp, nil
}

// DeleteUser deletes a user in Github Enterprise.
//
// https://developer.github.com/enterprise/2.16/v3/enterprise-admin/users/#delete-a-user
func (s *AdminService) DeleteUser(ctx context.Context, username string) (*Response, error) {
	u := "admin/users/" + username

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
