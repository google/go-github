// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

// ListProjectsPaginationOptions specifies optional parameters to list projects for user / organization.
//
// Note: Pagination is powered by before/after cursor-style pagination. After the initial call,
// inspect the returned *Response. Use resp.After as the opts.After value to request
// the next page, and resp.Before as the opts.Before value to request the previous
// page. Set either Before or After for a request; if both are
// supplied GitHub API will return an error. PerPage controls the number of items
// per page (max 100 per GitHub API docs).
type ListProjectsPaginationOptions struct {
	// A cursor, as given in the Link header. If specified, the query only searches for events before this cursor.
	Before string `url:"before,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for events after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// ListProjectsOptions specifies optional parameters to list projects for user / organization.
type ListProjectsOptions struct {
	ListProjectsPaginationOptions

	// Q is an optional query string to limit results to projects of the specified type.
	Query string `url:"q,omitempty"`
}

// ProjectV2FieldOption represents an option for a project field of type single_select or multi_select.
// It defines the available choices that can be selected for dropdown-style fields.
//
// GitHub API docs: https://docs.github.com/rest/projects/fields
type ProjectV2FieldOption struct {
	// The unique identifier for this option.
	ID string `json:"id,omitempty"`
	// The display name of the option.
	Name string `json:"name,omitempty"`
	// The color associated with this option (e.g., "blue", "red").
	Color string `json:"color,omitempty"`
	// An optional description for this option.
	Description string `json:"description,omitempty"`
}

// ProjectV2Field represents a field in a GitHub Projects V2 project.
// Fields define the structure and data types for project items.
//
// GitHub API docs: https://docs.github.com/rest/projects/fields
type ProjectV2Field struct {
	ID        *int64                  `json:"id,omitempty"`         // The unique identifier for this field.
	NodeID    string                  `json:"node_id,omitempty"`    // The GraphQL node ID for this field.
	Name      string                  `json:"name,omitempty"`       // The display name of the field.
	DataType  string                  `json:"dataType,omitempty"`   // The data type of the field (e.g., "text", "number", "date", "single_select", "multi_select").
	URL       string                  `json:"url,omitempty"`        // The API URL for this field.
	Options   []*ProjectV2FieldOption `json:"options,omitempty"`    // Available options for single_select and multi_select fields.
	CreatedAt *Timestamp              `json:"created_at,omitempty"` // The time when this field was created.
	UpdatedAt *Timestamp              `json:"updated_at,omitempty"` // The time when this field was last updated.
}

// ListProjectsForOrg lists Projects V2 for an organization.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#list-projects-for-organization
//
//meta:operation GET /orgs/{org}/projectsV2
func (s *ProjectsService) ListProjectsForOrg(ctx context.Context, org string, opts *ListProjectsOptions) ([]*ProjectV2, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var projects []*ProjectV2
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// GetProjectForOrg gets a Projects V2 project for an organization by ID.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#get-project-for-organization
//
//meta:operation GET /orgs/{org}/projectsV2/{project_number}
func (s *ProjectsService) GetProjectForOrg(ctx context.Context, org string, projectNumber int64) (*ProjectV2, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v", org, projectNumber)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(ProjectV2)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// ListProjectsForUser lists Projects V2 for a user.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#list-projects-for-user
//
//meta:operation GET /users/{username}/projectsV2
func (s *ProjectsService) ListProjectsForUser(ctx context.Context, username string, opts *ListProjectsOptions) ([]*ProjectV2, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2", username)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var projects []*ProjectV2
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// GetProjectForUser gets a Projects V2 project for a user by ID.
//
// GitHub API docs: https://docs.github.com/rest/projects/projects#get-project-for-user
//
//meta:operation GET /users/{username}/projectsV2/{project_number}
func (s *ProjectsService) GetProjectForUser(ctx context.Context, username string, projectNumber int64) (*ProjectV2, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v", username, projectNumber)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(ProjectV2)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// ListProjectFieldsForOrg lists Projects V2 for an organization.
//
// GitHub API docs: https://docs.github.com/rest/projects/fields#list-project-fields-for-organization
//
//meta:operation GET /orgs/{org}/projectsV2/{project_number}/fields
func (s *ProjectsService) ListProjectFieldsForOrg(ctx context.Context, org string, projectNumber int64, opts *ListProjectsOptions) ([]*ProjectV2Field, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/fields", org, projectNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var fields []*ProjectV2Field
	resp, err := s.client.Do(ctx, req, &fields)
	if err != nil {
		return nil, resp, err
	}
	return fields, resp, nil
}
