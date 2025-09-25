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

// ListProjectsOptions specifies optional parameters to list projects for user / organization.
//
// Note: Pagination is powered by before/after cursor-style pagination. After the initial call,
// inspect the returned *Response. Use resp.After as the opts.After value to request
// the next page, and resp.Before as the opts.Before value to request the previous
// page. Set either Before or After for a request; if both are
// supplied GitHub API will return an error. PerPage controls the number of items
// per page (max 100 per GitHub API docs).
type ListProjectsOptions struct {
	// A cursor, as given in the Link header. If specified, the query only searches for events before this cursor.
	Before string `url:"before,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for events after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`

	// Q is an optional query string to limit results to projects of the specified type.
	Query string `url:"q,omitempty"`
}

// ProjectV2FieldOption represents an option for a project field of type single_select or multi_select.
// It defines the available choices that can be selected for dropdown-style fields.
//
// GitHub API docs: https://docs.github.com/rest/projects/fields
type ProjectV2FieldOption struct {
	ID          *int64 `json:"id,omitempty"`          // The unique identifier for this option.
	Name        string `json:"name,omitempty"`        // The display name of the option.
	Color       string `json:"color,omitempty"`       // The color associated with this option (e.g., "blue", "red").
	Description string `json:"description,omitempty"` // An optional description for this option.
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

// NewProjectV2Field represents a new field to be added to a GitHub Projects V2.
type NewProjectV2Field struct {
	ID    *int64 `json:"id,omitempty"`    // The unique identifier for this field.
	Value any    `json:"value,omitempty"` // The value of the field.
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

// ListProjectItemsOptions specifies optional parameters when listing project items.
// Note: Pagination uses before/after cursor-style pagination similar to ListProjectsOptions.
// "Fields" can be used to restrict which field values are returned (by their numeric IDs).
type ListProjectItemsOptions struct {
	// Embed ListProjectsOptions to reuse pagination and query parameters.
	ListProjectsOptions
	// Fields restricts which field values are returned by numeric field IDs.
	Fields []int64 `url:"fields,omitempty,comma"`
}

// GetProjectItemOptions specifies optional parameters when getting a project item.
type GetProjectItemOptions struct {
	// Fields restricts which field values are returned by numeric field IDs.
	Fields []int64 `url:"fields,omitempty,comma"`
}

// AddProjectItemOptions represents the payload to add an item (issue or pull request)
// to a project. The Type must be either "Issue" or "PullRequest" (as per API docs) and
// ID is the numerical ID of that issue or pull request.
type AddProjectItemOptions struct {
	Type string `json:"type,omitempty"`
	ID   int64  `json:"id,omitempty"`
}

// UpdateProjectItemOptions represents fields that can be modified for a project item.
// Currently the REST API allows archiving/unarchiving an item (archived boolean).
// This struct can be expanded in the future as the API grows.
type UpdateProjectItemOptions struct {
	// Archived indicates whether the item should be archived (true) or unarchived (false).
	Archived *bool `json:"archived,omitempty"`
	// Fields allows updating field values for the item. Each entry supplies a field ID and a value.
	Fields []*NewProjectV2Field `json:"fields,omitempty"`
}

// ListProjectItemsForOrg lists items for an organization owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#list-items-for-an-organization-owned-project
//
//meta:operation GET /orgs/{org}/projectsV2/{project_number}/items
func (s *ProjectsService) ListProjectItemsForOrg(ctx context.Context, org string, projectNumber int64, opts *ListProjectItemsOptions) ([]*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/items", org, projectNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var items []*ProjectV2Item
	resp, err := s.client.Do(ctx, req, &items)
	if err != nil {
		return nil, resp, err
	}
	return items, resp, nil
}

// AddProjectItemForOrg adds an issue or pull request item to an organization owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#add-item-to-organization-owned-project
//
//meta:operation POST /orgs/{org}/projectsV2/{project_number}/items
func (s *ProjectsService) AddProjectItemForOrg(ctx context.Context, org string, projectNumber int64, opts *AddProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/items", org, projectNumber)
	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}

	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// GetProjectItemForOrg gets a single item from an organization owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#get-an-item-for-an-organization-owned-project
//
//meta:operation GET /orgs/{org}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) GetProjectItemForOrg(ctx context.Context, org string, projectNumber, itemID int64, opts *GetProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/items/%v", org, projectNumber, itemID)
	req, err := s.client.NewRequest("GET", u, opts)
	if err != nil {
		return nil, nil, err
	}
	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// UpdateProjectItemForOrg updates an item in an organization owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#update-project-item-for-organization
//
//meta:operation PATCH /orgs/{org}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) UpdateProjectItemForOrg(ctx context.Context, org string, projectNumber, itemID int64, opts *UpdateProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/items/%v", org, projectNumber, itemID)
	req, err := s.client.NewRequest("PATCH", u, opts)
	if err != nil {
		return nil, nil, err
	}
	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// DeleteProjectItemForOrg deletes an item from an organization owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#delete-project-item-for-organization
//
//meta:operation DELETE /orgs/{org}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) DeleteProjectItemForOrg(ctx context.Context, org string, projectNumber, itemID int64) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/projectsV2/%v/items/%v", org, projectNumber, itemID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ListProjectItemsForUser lists items for a user owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#list-items-for-a-user-owned-project
//
//meta:operation GET /users/{username}/projectsV2/{project_number}/items
func (s *ProjectsService) ListProjectItemsForUser(ctx context.Context, username string, projectNumber int64, opts *ListProjectItemsOptions) ([]*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v/items", username, projectNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var items []*ProjectV2Item
	resp, err := s.client.Do(ctx, req, &items)
	if err != nil {
		return nil, resp, err
	}
	return items, resp, nil
}

// AddProjectItemForUser adds an issue or pull request item to a user owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#add-item-to-user-owned-project
//
//meta:operation POST /users/{username}/projectsV2/{project_number}/items
func (s *ProjectsService) AddProjectItemForUser(ctx context.Context, username string, projectNumber int64, opts *AddProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v/items", username, projectNumber)
	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}
	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// GetProjectItemForUser gets a single item from a user owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#get-an-item-for-a-user-owned-project
//
//meta:operation GET /users/{username}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) GetProjectItemForUser(ctx context.Context, username string, projectNumber, itemID int64, opts *GetProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v/items/%v", username, projectNumber, itemID)
	req, err := s.client.NewRequest("GET", u, opts)
	if err != nil {
		return nil, nil, err
	}
	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// UpdateProjectItemForUser updates an item in a user owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#update-project-item-for-user
//
//meta:operation PATCH /users/{username}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) UpdateProjectItemForUser(ctx context.Context, username string, projectNumber, itemID int64, opts *UpdateProjectItemOptions) (*ProjectV2Item, *Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v/items/%v", username, projectNumber, itemID)
	req, err := s.client.NewRequest("PATCH", u, opts)
	if err != nil {
		return nil, nil, err
	}
	item := new(ProjectV2Item)
	resp, err := s.client.Do(ctx, req, item)
	if err != nil {
		return nil, resp, err
	}
	return item, resp, nil
}

// DeleteProjectItemForUser deletes an item from a user owned project.
//
// GitHub API docs: https://docs.github.com/rest/projects/items#delete-project-item-for-user
//
//meta:operation DELETE /users/{username}/projectsV2/{project_number}/items/{item_id}
func (s *ProjectsService) DeleteProjectItemForUser(ctx context.Context, username string, projectNumber, itemID int64) (*Response, error) {
	u := fmt.Sprintf("users/%v/projectsV2/%v/items/%v", username, projectNumber, itemID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
