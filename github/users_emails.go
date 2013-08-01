// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// ListEmails lists all authenticated user email addresses
//
// GitHub API docs: http://developer.github.com/v3/users/emails/#list-email-addresses-for-a-user
func (s *UsersService) ListEmails() ([]UserEmail, error) {
	u := "user/emails"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	emails := new([]UserEmail)
	_, err = s.client.Do(req, emails)
	return *emails, err
}

// AddEmails adds email addresses of authenticated user
//
// GitHub API docs: http://developer.github.com/v3/users/emails/#add-email-addresses
func (s *UsersService) AddEmails(emails []UserEmail) ([]UserEmail, error) {
	u := "user/emails"
	req, err := s.client.NewRequest("POST", u, emails)
	if err != nil {
		return nil, err
	}

	e := new([]UserEmail)
	_, err = s.client.Do(req, e)
	return *e, err
}

// DeleteEmails deletes email addresses from authenticated user
//
// GitHub API docs: http://developer.github.com/v3/users/emails/#delete-email-addresses
func (s *UsersService) DeleteEmails(emails []UserEmail) error {
	u := "user/emails"
	req, err := s.client.NewRequest("DELETE", u, emails)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
