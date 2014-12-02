// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// Promote a user.
//
// GitHub API docs: https://developer.github.com/v3/users/administration/#promote-an-ordinary-user-to-a-site-administrator
func (s *UsersService) Promote(user string) (*Response, error) {
	u := fmt.Sprintf("/users/%v/site_admin", user)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Demote a user.
//
// GitHub API docs: https://developer.github.com/v3/users/administration/#demote-a-site-administrator-to-an-ordinary-user
func (s *UsersService) Unsuspend(user string) (*Response, error) {
	u := fmt.Sprintf("users/%v/site_admin", user)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Suspend a user.
//
// GitHub API docs: https://developer.github.com/v3/users/administration/#suspend-a-user
func (s *UsersService) Suspend(user string) (*Response, error) {
	u := fmt.Sprintf("users/%v/suspended", user)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Unsuspend a user.
//
// GitHub API docs: https://developer.github.com/v3/users/administration/#unsuspend-a-user
func (s *UsersService) Unsuspend(user string) (*Response, error) {
	u := fmt.Sprintf("users/%v/suspended", user)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
