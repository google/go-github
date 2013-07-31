// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

// HooksService handles communication with the Hooks related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/
type HooksService struct {
	client *Client
}

// Hook represents a GitHub Hook for a repository.
type Hook struct {
	CreatedAt *time.Time        `json:"created_at,omitempty"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty"`
	Name      string            `json:"name"`
	Events    []string          `json:"events,omitempty"`
	Active    bool              `json:"active"`
	Config    map[string]string `json:"config,omitempty"`
	Id        int               `json:"id,omitempty"`
}

// Create a Hook for the specified repository.
// Name and Config are required fields.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#create-a-hook
func (s *HooksService) Create(owner string, repo string, hook *Hook) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks", owner, repo)
	req, err := s.client.NewRequest("POST", u, hook)
	if err != nil {
		return nil, err
	}
	h := new(Hook)
	_, err = s.client.Do(req, h)
	return h, err
}

// List all Hooks for the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#list
func (s *HooksService) List(owner string, repo string) ([]Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	hooks := new([]Hook)
	_, err = s.client.Do(req, hooks)
	return *hooks, err
}

// Get a single Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#get-single-hook
func (s *HooksService) Get(owner string, repo string, id int) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	hook := new(Hook)
	_, err = s.client.Do(req, hook)
	return hook, err
}

// Edit a Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#edit-a-hook
func (s *HooksService) Edit(owner string, repo string, id int, hook *Hook) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, hook)
	if err != nil {
		return nil, err
	}
	h := new(Hook)
	_, err = s.client.Do(req, h)
	return h, err
}

// Delete a Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#delete-a-hook
func (s *HooksService) Delete(owner string, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Test a Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#test-a-push-hook
func (s *HooksService) Test(owner string, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d/tests", owner, repo, id)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
