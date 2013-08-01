// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Hook represents a GitHub (web and service) hook for a repository.
type Hook struct {
  CreatedAt *time.Time             `json:"created_at,omitempty"`
  UpdatedAt *time.Time             `json:"updated_at,omitempty"`
  Name      string                 `json:"name,omitempty"`
  Events    []string               `json:"events,omitempty"`
  Active    bool                   `json:"active,omitempty"`
  Config    map[string]interface{} `json:"config,omitempty"`
  ID        int                    `json:"id,omitempty"`
}

// CreateHook creates a Hook for the specified repository.
// Name and Config are required fields.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#create-a-hook
func (s *RepositoriesService) CreateHook(owner, repo string, hook *Hook) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks", owner, repo)
	req, err := s.client.NewRequest("POST", u, hook)
	if err != nil {
		return nil, err
	}
	h := new(Hook)
	_, err = s.client.Do(req, h)
	return h, err
}

// ListHooks lists all Hooks for the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#list
func (s *RepositoriesService) ListHooks(owner, repo string, opt *ListOptions) ([]Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks", owner, repo)

	if opt != nil {
		params := url.Values{
			"page": []string{strconv.Itoa(opt.Page)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	hooks := new([]Hook)
	_, err = s.client.Do(req, hooks)
	return *hooks, err
}

// GetHook returns a single specified Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#get-single-hook
func (s *RepositoriesService) GetHook(owner, repo string, id int) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	hook := new(Hook)
	_, err = s.client.Do(req, hook)
	return hook, err
}

// EditHook updates a specified Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#edit-a-hook
func (s *RepositoriesService) EditHook(owner, repo string, id int, hook *Hook) (*Hook, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, hook)
	if err != nil {
		return nil, err
	}
	h := new(Hook)
	_, err = s.client.Do(req, h)
	return h, err
}

// DeleteHook deletes a specified Hook.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#delete-a-hook
func (s *RepositoriesService) DeleteHook(owner, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// TestHook triggers a test Hook by github.
//
// GitHub API docs: http://developer.github.com/v3/repos/hooks/#test-a-push-hook
func (s *RepositoriesService) TestHook(owner, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d/tests", owner, repo, id)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
