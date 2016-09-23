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
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"number":1}, {"number":2}]`)
	})

	opt := &ListOptions{Page: 2}
	projects, _, err := client.Repositories.ListProjects("o", "r", opt)
	if err != nil {
		t.Errorf("RepositoriesService.ListProjects returned error: %v", err)
	}

	want := []*RepositoryProject{{Number: Int(1)}, {Number: Int(2)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("RepositoriesService.ListProjects returned %+v, want %+v", projects, want)
	}
}

func TestRepositoriesService_GetProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"number":1}`)
	})

	project, _, err := client.Repositories.GetProject("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.GetProject returned %v", err)
	}

	want := &RepositoryProject{Number: Int(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("RepositoriesService.GetProject returned %+v, want %+v", project)
	}
}

func TestRepositoriesService_CreateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryProject{
		Name: String("n"),
		Body: String("b"),
	}

	mux.HandleFunc("/repos/o/r/projects", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryProject)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name":"n","body":"b"}`)
	})

	project, _, err := client.Repositories.CreateProject("o", "r", input)
	if err != nil {
		t.Errorf("RepositoriesService.CreateProject returned error: %v", err)
	}

	want := &RepositoryProject{Name: String("n"), Body: String("b")}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("RepositoriesService.CreateProject returned %+v, want %+v", project, want)
	}
}

func TestRepositoriesService_EditProject(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryProject{
		Name: String("n"),
		Body: String("b"),
	}

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryProject)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name":"n","body":"b"}`)
	})

	project, _, err := client.Repositories.EditProject("o", "r", 1, input)
	if err != nil {
		t.Errorf("RepositoriesService.EditProject returned error: %v", err)
	}

	want := &RepositoryProject{Name: String("n"), Body: String("b")}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("RepositoriesService.EditProject returned %+v, want %+v", project, want)
	}
}

func TestRepositoriesService_DeleteProject(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Repositories.DeleteProject("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.DeleteProject returned error: %v", err)
	}
}

func TestRepositoriesService_ListProjectColumns(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	columns, _, err := client.Repositories.ListProjectColumns("o", "r", 1, opt)
	if err != nil {
		t.Errorf("RepositoriesService.ListProjectColumns returned error: %v", err)
	}

	want := []*RepositoryColumn{{ID: Int(1)}, {ID: Int(2)}}
	if !reflect.DeepEqual(columns, want) {
		t.Errorf("RepositoriesService.ListProjectColumns returned %+v, want %+v", columns, want)
	}
}

func TestRepositoriesService_GetColumn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"id":1}`)
	})

	column, _, err := client.Repositories.GetColumn("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.GetColumn returned %v", err)
	}

	want := &RepositoryColumn{ID: Int(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("RepositoriesService.GetColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_CreateProjectColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryColumn{Name: String("n")}

	mux.HandleFunc("/repos/o/r/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryColumn)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name":"n"}`)
	})

	column, _, err := client.Repositories.CreateProjectColumn("o", "r", 1, input)
	if err != nil {
		t.Errorf("RepositoriesService.CreateProjectColumn returned error: %v", err)
	}

	want := &RepositoryColumn{Name: String("n")}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("RepositoriesService.CreateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_EditColumn(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryColumn{Name: String("n")}

	mux.HandleFunc("/repos/o/r/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryColumn)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name":"n"}`)
	})

	column, _, err := client.Repositories.EditColumn("o", "r", 1, input)
	if err != nil {
		t.Errorf("RepositoriesService.EditColumn returned error: %v", err)
	}

	want := &RepositoryColumn{Name: String("n")}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("RepositoriesService.EditColumn returned %+v, want %+v", column, want)
	}
}

func TestRepositoriesService_MoveColumnToFirstPlace(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("first")}

	mux.HandleFunc("/repos/o/r/projects/columns/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveColumnToFirstPlace("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.MoveColumnToFirstPlace returned error: %v", err)
	}
}

func TestRepositoriesService_MoveColumnToLastPlace(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("last")}

	mux.HandleFunc("/repos/o/r/projects/columns/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveColumnToLastPlace("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.MoveColumnToFirstPlace returned error: %v", err)
	}
}

func TestRepositoriesService_MoveColumnAfter(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("after:2")}

	mux.HandleFunc("/repos/o/r/projects/columns/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveColumnAfter("o", "r", 1, 2)
	if err != nil {
		t.Errorf("RepositoriesService.MoveColumnToFirstPlace returned error: %v", err)
	}
}

func TestRepositoriesService_ListColumnCards(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	cards, _, err := client.Repositories.ListColumnCards("o", "r", 1, opt)
	if err != nil {
		t.Errorf("RepositoriesService.ListColumnCards returned error: %v", err)
	}

	want := []*RepositoryCard{{ID: Int(1)}, {ID: Int(2)}}
	if !reflect.DeepEqual(cards, want) {
		t.Errorf("RepositoriesService.ListColumnCards returned %+v, want %+v", cards, want)
	}
}

func TestRepositoriesService_GetCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"id":1}}`)
	})

	card, _, err := client.Repositories.GetCard("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.GetCard returned error: %v")
	}

	want := &RepositoryCard{ID: Int(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("RepositoriesService.GetCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_CreateColumnCard(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryCard{Note: String("n")}

	mux.HandleFunc("/repos/o/r/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryCard)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"note":"n"}`)
	})

	card, _, err := client.Repositories.CreateColumnCard("o", "r", 1, input)
	if err != nil {
		t.Errorf("RepositoriesService.CreateColumnCard returned error: %v", err)
	}

	want := &RepositoryCard{Note: String("n")}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("RepositoriesService.CreateColumnCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_EditCard(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryCard{Note: String("n")}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryCard)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"note":"n"}`)
	})

	card, _, err := client.Repositories.EditCard("o", "r", 1, input)
	if err != nil {
		t.Errorf("RepositoriesService.EditCard returned error: %v", err)
	}

	want := &RepositoryCard{Note: String("n")}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("RepositoriesService.EditCard returned %+v, want %+v", card, want)
	}
}

func TestRepositoriesService_DeleteCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Repositories.DeleteCard("o", "r", 1)
	if err != nil {
		t.Error("RepositoriesService.DeleteCard returned error: %v", err)
	}
}

func TestRepositoriesService_MoveCardToTop(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("top")}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveCardToTop("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.MoveCardToTop returned error: %v", err)
	}
}

func TestRepositoriesService_MoveCardToBottom(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("bottom")}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveCardToBottom("o", "r", 1)
	if err != nil {
		t.Errorf("RepositoriesService.MoveCardToBottom returned error: %v", err)
	}
}

func TestRepositoriesService_MoveCardAfter(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		Position *string `json:"position,omitempty"`
	}

	input := &RequestBody{Position: String("after:2")}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveCardAfter("o", "r", 1, 2)
	if err != nil {
		t.Errorf("RepositoriesService.MoveCardAfter returned error: %v", err)
	}
}

func TestRepositoriesService_MoveCardToColumn(t *testing.T) {
	setup()
	defer teardown()

	type RequestBody struct {
		ColumnID *int `json="column_id,omitempty"`
	}

	input := &RequestBody{ColumnID: Int(2)}

	mux.HandleFunc("/repos/o/r/projects/columns/cards/1/moves", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestBody)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	_, err := client.Repositories.MoveCardToColumn("o", "r", 1, 2)
	if err != nil {
		t.Errorf("RepositoriesService.MoveCardToColumn returned error: %v", err)
	}
}
