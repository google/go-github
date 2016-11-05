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

func TestProjectsService_UpdateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectOptions{Name: "Project Name", Body: "Project body."}

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Projects.UpdateProject(1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProject returned error: %v", err)
	}

	want := &Project{ID: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.UpdateProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Projects.GetProject(1)
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &Project{ID: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_DeleteProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProject(1)
	if err != nil {
		t.Errorf("Projects.DeleteProject returned error: %v", err)
	}
}

func TestProjectsService_ListProjectColumns(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	columns, _, err := client.Projects.ListProjectColumns(1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectColumns returned error: %v", err)
	}

	want := []*ProjectColumn{{ID: Int(1)}}
	if !reflect.DeepEqual(columns, want) {
		t.Errorf("Projects.ListProjectColumns returned %+v, want %+v", columns, want)
	}
}

func TestProjectsService_GetProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Projects.GetProjectColumn(1)
	if err != nil {
		t.Errorf("Projects.GetProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.GetProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_CreateProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnOptions{Name: "Column Name"}

	mux.HandleFunc("/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Projects.CreateProjectColumn(1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.CreateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_UpdateProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnOptions{Name: "Column Name"}

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Projects.UpdateProjectColumn(1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.UpdateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_DeleteProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProjectColumn(1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectColumn returned error: %v", err)
	}
}

func TestProjectsService_MoveProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectColumnMoveOptions{Position: "after:12345"}

	mux.HandleFunc("/projects/columns/1/moves", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectColumnMoveOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Projects.MoveProjectColumn(1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectColumn returned error: %v", err)
	}
}

func TestProjectsService_ListProjectCards(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	cards, _, err := client.Projects.ListProjectCards(1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCards returned error: %v", err)
	}

	want := []*ProjectCard{{ID: Int(1)}}
	if !reflect.DeepEqual(cards, want) {
		t.Errorf("Projects.ListProjectCards returned %+v, want %+v", cards, want)
	}
}

func TestProjectsService_GetProjectCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Projects.GetProjectCard(1)
	if err != nil {
		t.Errorf("Projects.GetProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.GetProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_CreateProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardOptions{
		ContentID:   12345,
		ContentType: "Issue",
	}

	mux.HandleFunc("/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Projects.CreateProjectCard(1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.CreateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_UpdateProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardOptions{
		ContentID:   12345,
		ContentType: "Issue",
	}

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Projects.UpdateProjectCard(1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.UpdateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_DeleteProjectCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProjectCard(1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectCard returned error: %v", err)
	}
}

func TestProjectsService_MoveProjectCard(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectCardMoveOptions{Position: "after:12345"}

	mux.HandleFunc("/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCardMoveOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Projects.MoveProjectCard(1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectCard returned error: %v", err)
	}
}
