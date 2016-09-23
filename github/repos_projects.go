// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// Projects represents a Project's configuration of a GitHub Repository
type RepositoryProject struct {
	ID        *int       `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Body      *string    `json:"body,omitempty"`
	URL       *string    `json:"url,omitempty"`
	OwnerURL  *string    `json:"url,omitempty"`
	Number    *int       `json:"number,omitempty"`
	Creator   *User      `json:"creator,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

func (r RepositoryProject) String() string {
	return Stringify(r)
}

// ListProjects lists all the projects for the repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-projects
func (s *RepositoriesService) ListProjects(owner, repo string, opt *ListOptions) ([]*RepositoryProject, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects", owner, repo)
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

	projects := new([]*RepositoryProject)
	resp, err := s.client.Do(req, projects)
	if err != nil {
		return nil, resp, err
	}

	return *projects, resp, nil
}

// GetProject gets a single project from a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-a-project
func (s *RepositoriesService) GetProject(owner, repo string, number int) (*RepositoryProject, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	p := new(RepositoryProject)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// CreateProject creates a new project on the specified repository.
//
// GitHub Api Docs: https://developer.github.com/v3/repos/projects/#create-a-project
func (s *RepositoriesService) CreateProject(owner, repo string, project *RepositoryProject) (*RepositoryProject, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects", owner, repo)
	req, err := s.client.NewRequest("POST", u, project)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	p := new(RepositoryProject)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// EditProject edits a project of a repository.
//
// GitHub Api Docs: https://developer.github.com/v3/repos/projects/#update-a-project
func (s *RepositoriesService) EditProject(owner, repo string, number int, project *RepositoryProject) (*RepositoryProject, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, project)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	p := new(RepositoryProject)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// DeleteProject deletes a project of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#delete-a-project
func (s *RepositoriesService) DeleteProject(owner, repo string, number int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// Projects represents a Column's configuration of a Project of a GitHub Repository
type RepositoryColumn struct {
	ID         *int       `json:"id,omitempty"`
	Name       *string    `json:"name,omitempty"`
	ProjectUrl *string    `json:"project_url,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp `json:"updated_at,omitempty"`
}

func (r RepositoryColumn) String() string {
	return Stringify(r)
}

// ListProjectColumns lists all the columns for a specified project of a repository.
//
// GitHub Api Docs: https://developer.github.com/v3/repos/projects/#list-columns
func (s *RepositoriesService) ListProjectColumns(owner, repo string, number int, opt *ListOptions) ([]*RepositoryColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v/columns", owner, repo, number)
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

	columns := new([]*RepositoryColumn)
	resp, err := s.client.Do(req, columns)
	if err != nil {
		return nil, resp, err
	}

	return *columns, resp, nil
}

// GetColumn gets a single column of a repository.
//
// GitHub Api Docs: https://developer.github.com/v3/repos/projects/#get-a-column
func (s *RepositoriesService) GetColumn(owner, repo string, id int) (*RepositoryColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(RepositoryColumn)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// CreateProjectColumn creates a new colum on the specified project of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-column
func (s *RepositoriesService) CreateProjectColumn(owner, repo string, number int, column *RepositoryColumn) (*RepositoryColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v/columns", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, column)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(RepositoryColumn)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// EditColumn edits a column of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#update-a-column
func (s *RepositoriesService) EditColumn(owner, repo string, id int, column *RepositoryColumn) (*RepositoryColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, column)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(RepositoryColumn)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// DeleteColumn deletes a column of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#delete-a-column
func (s *RepositoriesService) DeleteColumn(owner, repo string, id int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// MoveColumn moves a column of a repository.
//
// GitHun API docs: https://developer.github.com/v3/repos/projects/#move-a-column
func (s *RepositoriesService) MoveColumn(owner, repo string, id int, position string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/moves", owner, repo, id)
	data := struct {
		Position *string `json="position,omitempty"`
	}{
		&position,
	}

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// MoveColumnToFirstPlace moves a column of a project on a repository to the first place.
//
// GitHun API docs: https://developer.github.com/v3/repos/projects/#move-a-column
func (s *RepositoriesService) MoveColumnToFirstPlace(owner, repo string, id int) (*Response, error) {
	return s.MoveColumn(owner, repo, id, "first")
}

// MoveColumnToLastPlace moves a column of a project on a repository to the last place.
//
// GitHun API docs: https://developer.github.com/v3/repos/projects/#move-a-column
func (s *RepositoriesService) MoveColumnToLastPlace(owner, repo string, id int) (*Response, error) {
	return s.MoveColumn(owner, repo, id, "last")
}

// MoveColumnAfter moves a column of a project on a repository after a selected column.
//
// GitHun API docs: https://developer.github.com/v3/repos/projects/#move-a-column
func (s *RepositoriesService) MoveColumnAfter(owner, repo string, id, after int) (*Response, error) {
	position := fmt.Sprintf("after:%v", after)

	return s.MoveColumn(owner, repo, id, position)
}

// Projects represents a Card's configuration of a Project of a GitHub Repository
type RepositoryCard struct {
	ColumnUrl  *string    `json:"column_url,omitempty"`
	ContentUrl *string    `json:"column_url,omitempty"`
	ID         *int       `json:"id,omitempty"`
	Note       *string    `json:"note,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp `json:"updated_at,omitempty"`
}

func (r RepositoryCard) String() string {
	return Stringify(r)
}

// ListProjectCards lists all the cards on a column of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-projects-cards
func (s *RepositoriesService) ListColumnCards(owner, repo string, column int, opt *ListOptions) ([]*RepositoryCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/cards", owner, repo, column)
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

	cards := new([]*RepositoryCard)
	resp, err := s.client.Do(req, cards)
	if err != nil {
		return nil, resp, err
	}

	return *cards, resp, nil
}

// GetCard gets a single card of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-a-project-card
func (s *RepositoriesService) GetCard(owner, repo string, id int) (*RepositoryCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, id)
	u, err := addOptions(u, nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(*RepositoryCard)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return *c, resp, nil
}

// CreateColumnCard creates a card on the specified column of a project of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-project-card
func (s *RepositoriesService) CreateColumnCard(owner, repo string, column int, card *RepositoryCard) (*RepositoryCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/cards", owner, repo, column)
	req, err := s.client.NewRequest("POST", u, card)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(RepositoryCard)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// EditCard edits a card of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#update-a-project-card
func (s *RepositoriesService) EditCard(owner, repo string, id int, card *RepositoryCard) (*RepositoryCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, card)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	c := new(RepositoryCard)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// DeleteCard deletes a card of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-project-card
func (s *RepositoriesService) DeleteCard(owner, repo string, id int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// MoveCard moves a card of a repository to the given position of it's the column.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveCard(owner, repo string, id int, position string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v/moves", owner, repo, id)
	data := struct {
		Position *string `json="position,omitempty"`
	}{
		&position,
	}

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// MoveCardToTop moves a card of a repository to the top of it's column.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveCardToTop(owner, repo string, id int) (*Response, error) {
	return s.MoveCard(owner, repo, id, "top")
}

// MoveCardToBottom moves a card of a repository to the bottom of it's column.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveCardToBottom(owner, repo string, id int) (*Response, error) {
	return s.MoveCard(owner, repo, id, "bottom")
}

// MoveCardToBottom moves a card of a repository after a selected card.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveCardAfter(owner, repo string, id, card int) (*Response, error) {
	position := fmt.Sprintf("after:%v", card)

	return s.MoveCard(owner, repo, id, position)
}

// MoveCardToBottom moves a card of a repository to the specified column.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveCardToColumn(owner, repo string, id, column int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v/moves", owner, repo, id)
	data := struct {
		ColumnID *int `json="column_id,omitempty"`
	}{
		&column,
	}

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}
