// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ProjectsService handles communication with the project V2
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects
type ProjectsService service

func (p ProjectV2) String() string { return Stringify(p) }

// ListProjectsOptions specifies optional parameters to list projects for user / organization.
type ListProjectsOptions struct {
	// A cursor, as given in the Link header. If specified, the query only searches for events before this cursor.
	Before string `url:"before,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for events after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`

	// Q is an optional query string to filter/search projects (when supported).
	Query string `url:"q,omitempty"`
}

// ListOrganizationProjects lists Projects V2 for an organization.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#list-organization-projects
//
//meta:operation GET /orgs/{org}/projectsV2
func (s *ProjectsService) ListOrganizationProjects(ctx context.Context, org string, opts *ListProjectsOptions) ([]*ProjectV2, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	var projects []*ProjectV2
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// GetByOrg gets a Projects V2 project for an organization by ID.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#get-project-for-organization
//
//meta:operation GET /orgs/{org}/projectsV2/{project_id}
func (s *ProjectsService) GetByOrg(ctx context.Context, org string, projectID int64) (*ProjectV2, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v", org, projectID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := new(ProjectV2)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// ListByUser lists Projects V2 for a user.
//
// GitHub API docs: https://docs.github.com/en/rest/projects/projects#list-projects-for-user
//
//meta:operation GET /users/{username}/projectsV2
func (s *ProjectsService) ListByUser(ctx context.Context, username string, opts *ListProjectsOptions) ([]*ProjectV2, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2", username)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	var projects []*ProjectV2
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// GetUserProject gets a Projects V2 project for a user by ID.
//
// GitHub API docs: https://docs.github.com/en/rest/projects/projects#get-project-for-user
//
//meta:operation GET /users/{username}/projectsV2/{project_id}
func (s *ProjectsService) GetUserProject(ctx context.Context, username string, projectID int64) (*ProjectV2, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v", username, projectID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := new(ProjectV2)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}
