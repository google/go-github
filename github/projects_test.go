// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestProject_marshall(t *testing.T) {
	testJSONMarshal(t, &Project{}, "{}")

	u := &Project{
		ID:         Int64(1),
		URL:        String("u"),
		HTMLURL:    String("h"),
		ColumnsURL: String("c"),
		OwnerURL:   String("o"),
		Name:       String("n"),
		Body:       String("b"),
		Number:     Int(1),
		State:      String("s"),
		CreatedAt:  &Timestamp{referenceTime},
		UpdatedAt:  &Timestamp{referenceTime},
		NodeID:     String("n"),
		Creator: &User{
			Login:       String("l"),
			ID:          Int64(1),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		},
	}
	want := `{
		"id": 1,
		"url": "u",
		"html_url": "h",
		"columns_url": "c",
		"owner_url": "o",
		"name": "n",
		"body": "b",
		"number": 1,
		"state": "s",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"node_id": "n",
		"creator": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"url": "u"
		}
	}`
	testJSONMarshal(t, u, want)
}
func TestProjectsService_UpdateProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ProjectOptions{
		Name:   String("Project Name"),
		Body:   String("Project body."),
		State:  String("open"),
		Public: Bool(true),

		OrganizationPermission: String("read"),
	}

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

	project, _, err := client.Projects.UpdateProject(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.UpdateProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Projects.GetProject(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_DeleteProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProject(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.DeleteProject returned error: %v", err)
	}
}

func TestProjectsService_ListProjectColumns(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/projects/1/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	columns, _, err := client.Projects.ListProjectColumns(context.Background(), 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectColumns returned error: %v", err)
	}

	want := []*ProjectColumn{{ID: Int64(1)}}
	if !reflect.DeepEqual(columns, want) {
		t.Errorf("Projects.ListProjectColumns returned %+v, want %+v", columns, want)
	}
}

func TestProjectsService_GetProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	column, _, err := client.Projects.GetProjectColumn(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.GetProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.GetProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_CreateProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
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

	column, _, err := client.Projects.CreateProjectColumn(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.CreateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_UpdateProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
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

	column, _, err := client.Projects.UpdateProjectColumn(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !reflect.DeepEqual(column, want) {
		t.Errorf("Projects.UpdateProjectColumn returned %+v, want %+v", column, want)
	}
}

func TestProjectsService_DeleteProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProjectColumn(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectColumn returned error: %v", err)
	}
}

func TestProjectsService_MoveProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
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

	_, err := client.Projects.MoveProjectColumn(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectColumn returned error: %v", err)
	}
}

func TestProjectsService_ListProjectCards(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1/cards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{
			"archived_state": "all",
			"page":           "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ProjectCardListOptions{
		ArchivedState: String("all"),
		ListOptions:   ListOptions{Page: 2}}
	cards, _, err := client.Projects.ListProjectCards(context.Background(), 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCards returned error: %v", err)
	}

	want := []*ProjectCard{{ID: Int64(1)}}
	if !reflect.DeepEqual(cards, want) {
		t.Errorf("Projects.ListProjectCards returned %+v, want %+v", cards, want)
	}
}

func TestProjectsService_GetProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	card, _, err := client.Projects.GetProjectCard(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.GetProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.GetProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_CreateProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
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

	card, _, err := client.Projects.CreateProjectCard(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.CreateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_UpdateProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
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

		fmt.Fprint(w, `{"id":1, "archived":false}`)
	})

	card, _, err := client.Projects.UpdateProjectCard(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1), Archived: Bool(false)}
	if !reflect.DeepEqual(card, want) {
		t.Errorf("Projects.UpdateProjectCard returned %+v, want %+v", card, want)
	}
}

func TestProjectsService_DeleteProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	_, err := client.Projects.DeleteProjectCard(context.Background(), 1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectCard returned error: %v", err)
	}
}

func TestProjectsService_MoveProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
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

	_, err := client.Projects.MoveProjectCard(context.Background(), 1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectCard returned error: %v", err)
	}
}

func TestProjectsService_AddProjectCollaborator(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &ProjectCollaboratorOptions{
		Permission: String("admin"),
	}

	mux.HandleFunc("/projects/1/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectCollaboratorOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Projects.AddProjectCollaborator(context.Background(), 1, "u", opt)
	if err != nil {
		t.Errorf("Projects.AddProjectCollaborator returned error: %v", err)
	}
}

func TestProjectsService_AddCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Projects.AddProjectCollaborator(context.Background(), 1, "%", nil)
	testURLParseError(t, err)
}

func TestProjectsService_RemoveCollaborator(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1/collaborators/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Projects.RemoveProjectCollaborator(context.Background(), 1, "u")
	if err != nil {
		t.Errorf("Projects.RemoveProjectCollaborator returned error: %v", err)
	}
}

func TestProjectsService_RemoveCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Projects.RemoveProjectCollaborator(context.Background(), 1, "%")
	testURLParseError(t, err)
}

func TestProjectsService_ListCollaborators(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1/collaborators", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListCollaboratorOptions{
		ListOptions: ListOptions{Page: 2},
	}
	users, _, err := client.Projects.ListProjectCollaborators(context.Background(), 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Projects.ListProjectCollaborators returned %+v, want %+v", users, want)
	}
}

func TestProjectsService_ListCollaborators_withAffiliation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1/collaborators", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"affiliation": "all", "page": "2"})
		fmt.Fprintf(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListCollaboratorOptions{
		ListOptions: ListOptions{Page: 2},
		Affiliation: String("all"),
	}
	users, _, err := client.Projects.ListProjectCollaborators(context.Background(), 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Projects.ListProjectCollaborators returned %+v, want %+v", users, want)
	}
}

func TestProjectsService_GetPermissionLevel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1/collaborators/u/permission", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprintf(w, `{"permission":"admin","user":{"login":"u"}}`)
	})

	ppl, _, err := client.Projects.ReviewProjectCollaboratorPermission(context.Background(), 1, "u")
	if err != nil {
		t.Errorf("Projects.ReviewProjectCollaboratorPermission returned error: %v", err)
	}

	want := &ProjectPermissionLevel{
		Permission: String("admin"),
		User: &User{
			Login: String("u"),
		},
	}

	if !reflect.DeepEqual(ppl, want) {
		t.Errorf("Projects.ReviewProjectCollaboratorPermission returned %+v, want %+v", ppl, want)
	}
}
