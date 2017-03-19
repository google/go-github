// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListProjectsOptions specifies the optional parameters to the
// OrganizationsService.ListProjects method.
type ListProjectsOptions struct {
	// Indicates the state of the projects to return. Can be either open, closed, or all. Default: open
	State string `url:"state,omitempty"`

	ListOptions
}

// ListProjects lists the projects for an organization.
//
// GitHub API docs: https://developer.github.com/v3/projects/#list-organization-projects
func (s *OrganizationsService) ListProjects(ctx context.Context, org string, opt *ListProjectsOptions) ([]*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/projects", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	projects := []*Project{}
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}

	return projects, resp, nil
}

// CreateProject creates a GitHub Project for the specified organization.
//
// GitHub API docs: https://developer.github.com/v3/projects/#create-an-organization-project
func (s *OrganizationsService) CreateProject(ctx context.Context, org string, opt *ProjectOptions) (*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/projects", org)
	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := &Project{}
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}
