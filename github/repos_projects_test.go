// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoriesService_ListProjects(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	projects, _, err := client.Repositories.ListProjects("o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListProjects returned error: %v", err)
	}

	want := []*Project{{ID: Int(1)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Repositories.ListProjects returned %+v, want %+v", projects, want)
	}
}

func TestRepositoriesService_GetProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Repositories.GetProject("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetProject returned error: %v", err)
	}

	want := &Project{ID: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Repositories.GetProject returned %+v, want %+v", project, want)
	}
}

func TestRepositoriesService_CreateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectOptions{Name: "Project Name", Body: "Project body."}

	mux.HandleFunc("/repos/o/r/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Repositories.CreateProject("o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateProject returned error: %v", err)
	}

	want := &Project{ID: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Repositories.CreateProject returned %+v, want %+v", project, want)
	}
}

func TestRepositoriesService_UpdateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectOptions{Name: "Project Name", Body: "Project body."}

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Repositories.UpdateProject("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.UpdateProject returned error: %v", err)
	}

	want := &Project{ID: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Repositories.UpdateProject returned %+v, want %+v", project, want)
	}
}

func TestRepositoriesService_DeleteProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Repositories.DeleteProject("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteProject returned error: %v", err)
	}
}

func TestRepositoriesService_ListProjectColumns(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	columns, _, err := client.Repositories.ListProjectColumns("o", "r", 1, opt)
	if err != nil {
		t.Errorf("Repositories.ListProjectColumns returned error: %v", err)
	}

	want := []*ProjectColumn{{ID: Int(1)}}
	if !reflect.DeepEqual(columns, want) {
		t.Errorf("Repositories.ListProjectColumns returned %+v, want %+v", columns, want)
	}
}

func TestRepositoriesService_GetProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Repositories.GetProjectColumn("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Repositories.GetProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_CreateProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnOptions{Name: "Column Name"}

	mux.HandleFunc("/repos/o/r/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Repositories.CreateProjectColumn("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.CreateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Repositories.CreateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_UpdateProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnOptions{Name: "Column Name"}

	mux.HandleFunc("/repos/o/r/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Repositories.UpdateProjectColumn("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.UpdateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Repositories.UpdateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_DeleteProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Repositories.DeleteProjectColumn("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteProjectColumn returned error: %v", err)
	}
}

func TestRepositoriesService_MoveProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnMoveOptions{Position: "after:12345"}

	mux.HandleFunc("/repos/o/r/projects/columns/1/moves", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnMoveOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveProjectColumn("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.MoveProjectColumn returned error: %v", err)
	}
}

func TestRepositoriesService_ListProjectCards(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	cards, _, err := client.Repositories.ListProjectCards("o", "r", 1, opt)
	if err != nil {
		t.Errorf("Repositories.ListProjectCards returned error: %v", err)
	}

	want := []*ProjectCard{{ID: Int(1)}}
	if !reflect.DeepEqual(cards, want) {
		t.Errorf("Repositories.ListProjectCards returned %+v, want %+v", cards, want)
	}
}

func TestRepositoriesService_GetProjectCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Repositories.GetProjectCard("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Repositories.GetProjectCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_CreateProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardOptions{
		ContentID:   12345,
		ContentType: "Issue",
	}

	mux.HandleFunc("/repos/o/r/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Repositories.CreateProjectCard("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.CreateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Repositories.CreateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_UpdateProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardOptions{
		ContentID:   12345,
		ContentType: "Issue",
	}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Repositories.UpdateProjectCard("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.UpdateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Repositories.UpdateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_DeleteProjectCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Repositories.DeleteProjectCard("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteProjectCard returned error: %v", err)
	}
}

func TestRepositoriesService_MoveProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardMoveOptions{Position: "after:12345"}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardMoveOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveProjectCard("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.MoveProjectCard returned error: %v", err)
	}
}
