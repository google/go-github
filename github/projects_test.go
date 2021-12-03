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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProject_Marshal(t *testing.T) {
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
		Name:    String("Project Name"),
		Body:    String("Project body."),
		State:   String("open"),
		Private: Bool(false),

		OrganizationPermission: String("read"),
	}

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.UpdateProject(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Projects.UpdateProject returned %+v, want %+v", project, want)
	}

	const methodName = "UpdateProject"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.UpdateProject(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateProject(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.GetProject(ctx, 1)
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}

	const methodName = "GetProject"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.GetProject(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProject(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	ctx := context.Background()
	_, err := client.Projects.DeleteProject(ctx, 1)
	if err != nil {
		t.Errorf("Projects.DeleteProject returned error: %v", err)
	}

	const methodName = "DeleteProject"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.DeleteProject(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteProject(ctx, 1)
	})
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
	ctx := context.Background()
	columns, _, err := client.Projects.ListProjectColumns(ctx, 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectColumns returned error: %v", err)
	}

	want := []*ProjectColumn{{ID: Int64(1)}}
	if !cmp.Equal(columns, want) {
		t.Errorf("Projects.ListProjectColumns returned %+v, want %+v", columns, want)
	}

	const methodName = "ListProjectColumns"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectColumns(ctx, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectColumns(ctx, 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	column, _, err := client.Projects.GetProjectColumn(ctx, 1)
	if err != nil {
		t.Errorf("Projects.GetProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !cmp.Equal(column, want) {
		t.Errorf("Projects.GetProjectColumn returned %+v, want %+v", column, want)
	}

	const methodName = "GetProjectColumn"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.GetProjectColumn(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectColumn(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	column, _, err := client.Projects.CreateProjectColumn(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !cmp.Equal(column, want) {
		t.Errorf("Projects.CreateProjectColumn returned %+v, want %+v", column, want)
	}

	const methodName = "CreateProjectColumn"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.CreateProjectColumn(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.CreateProjectColumn(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	column, _, err := client.Projects.UpdateProjectColumn(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectColumn returned error: %v", err)
	}

	want := &ProjectColumn{ID: Int64(1)}
	if !cmp.Equal(column, want) {
		t.Errorf("Projects.UpdateProjectColumn returned %+v, want %+v", column, want)
	}

	const methodName = "UpdateProjectColumn"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.UpdateProjectColumn(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateProjectColumn(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteProjectColumn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	ctx := context.Background()
	_, err := client.Projects.DeleteProjectColumn(ctx, 1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectColumn returned error: %v", err)
	}

	const methodName = "DeleteProjectColumn"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.DeleteProjectColumn(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteProjectColumn(ctx, 1)
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()
	_, err := client.Projects.MoveProjectColumn(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectColumn returned error: %v", err)
	}

	const methodName = "MoveProjectColumn"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.MoveProjectColumn(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.MoveProjectColumn(ctx, 1, input)
	})
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
	ctx := context.Background()
	cards, _, err := client.Projects.ListProjectCards(ctx, 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCards returned error: %v", err)
	}

	want := []*ProjectCard{{ID: Int64(1)}}
	if !cmp.Equal(cards, want) {
		t.Errorf("Projects.ListProjectCards returned %+v, want %+v", cards, want)
	}

	const methodName = "ListProjectCards"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectCards(ctx, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectCards(ctx, 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	card, _, err := client.Projects.GetProjectCard(ctx, 1)
	if err != nil {
		t.Errorf("Projects.GetProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1)}
	if !cmp.Equal(card, want) {
		t.Errorf("Projects.GetProjectCard returned %+v, want %+v", card, want)
	}

	const methodName = "GetProjectCard"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.GetProjectCard(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectCard(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	card, _, err := client.Projects.CreateProjectCard(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.CreateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1)}
	if !cmp.Equal(card, want) {
		t.Errorf("Projects.CreateProjectCard returned %+v, want %+v", card, want)
	}

	const methodName = "CreateProjectCard"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.CreateProjectCard(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.CreateProjectCard(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1, "archived":false}`)
	})

	ctx := context.Background()
	card, _, err := client.Projects.UpdateProjectCard(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.UpdateProjectCard returned error: %v", err)
	}

	want := &ProjectCard{ID: Int64(1), Archived: Bool(false)}
	if !cmp.Equal(card, want) {
		t.Errorf("Projects.UpdateProjectCard returned %+v, want %+v", card, want)
	}

	const methodName = "UpdateProjectCard"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.UpdateProjectCard(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateProjectCard(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteProjectCard(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/columns/cards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
	})

	ctx := context.Background()
	_, err := client.Projects.DeleteProjectCard(ctx, 1)
	if err != nil {
		t.Errorf("Projects.DeleteProjectCard returned error: %v", err)
	}

	const methodName = "DeleteProjectCard"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.DeleteProjectCard(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteProjectCard(ctx, 1)
	})
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()
	_, err := client.Projects.MoveProjectCard(ctx, 1, input)
	if err != nil {
		t.Errorf("Projects.MoveProjectCard returned error: %v", err)
	}

	const methodName = "MoveProjectCard"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.MoveProjectCard(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.MoveProjectCard(ctx, 1, input)
	})
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
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Projects.AddProjectCollaborator(ctx, 1, "u", opt)
	if err != nil {
		t.Errorf("Projects.AddProjectCollaborator returned error: %v", err)
	}

	const methodName = "AddProjectCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.AddProjectCollaborator(ctx, -1, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.AddProjectCollaborator(ctx, 1, "u", opt)
	})
}

func TestProjectsService_AddCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Projects.AddProjectCollaborator(ctx, 1, "%", nil)
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

	ctx := context.Background()
	_, err := client.Projects.RemoveProjectCollaborator(ctx, 1, "u")
	if err != nil {
		t.Errorf("Projects.RemoveProjectCollaborator returned error: %v", err)
	}

	const methodName = "RemoveProjectCollaborator"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Projects.RemoveProjectCollaborator(ctx, -1, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.RemoveProjectCollaborator(ctx, 1, "u")
	})
}

func TestProjectsService_RemoveCollaborator_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Projects.RemoveProjectCollaborator(ctx, 1, "%")
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
	ctx := context.Background()
	users, _, err := client.Projects.ListProjectCollaborators(ctx, 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(users, want) {
		t.Errorf("Projects.ListProjectCollaborators returned %+v, want %+v", users, want)
	}

	const methodName = "ListProjectCollaborators"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectCollaborators(ctx, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectCollaborators(ctx, 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	users, _, err := client.Projects.ListProjectCollaborators(ctx, 1, opt)
	if err != nil {
		t.Errorf("Projects.ListProjectCollaborators returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(users, want) {
		t.Errorf("Projects.ListProjectCollaborators returned %+v, want %+v", users, want)
	}
}

func TestProjectsService_ReviewProjectCollaboratorPermission(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/projects/1/collaborators/u/permission", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprintf(w, `{"permission":"admin","user":{"login":"u"}}`)
	})

	ctx := context.Background()
	ppl, _, err := client.Projects.ReviewProjectCollaboratorPermission(ctx, 1, "u")
	if err != nil {
		t.Errorf("Projects.ReviewProjectCollaboratorPermission returned error: %v", err)
	}

	want := &ProjectPermissionLevel{
		Permission: String("admin"),
		User: &User{
			Login: String("u"),
		},
	}

	if !cmp.Equal(ppl, want) {
		t.Errorf("Projects.ReviewProjectCollaboratorPermission returned %+v, want %+v", ppl, want)
	}

	const methodName = "ReviewProjectCollaboratorPermission"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ReviewProjectCollaboratorPermission(ctx, -1, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ReviewProjectCollaboratorPermission(ctx, 1, "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectOptions{}, "{}")

	u := &ProjectOptions{
		Name:                   String("name"),
		Body:                   String("body"),
		State:                  String("state"),
		OrganizationPermission: String("op"),
		Private:                Bool(false),
	}

	want := `{
		"name": "name",
		"body": "body",
		"state": "state",
		"organization_permission": "op",
		"private": false
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumn_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumn{}, "{}")

	u := &ProjectColumn{
		ID:         Int64(1),
		Name:       String("name"),
		URL:        String("url"),
		ProjectURL: String("purl"),
		CardsURL:   String("curl"),
		CreatedAt:  &Timestamp{referenceTime},
		UpdatedAt:  &Timestamp{referenceTime},
		NodeID:     String("onidp"),
	}

	want := `{
		"id": 1,
		"name": "name",
		"url": "url",
		"project_url": "purl",
		"cards_url": "curl",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"node_id": "onidp"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnOptions{}, "{}")

	u := &ProjectColumnOptions{
		Name: "name",
	}

	want := `{
		"name": "name"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnMoveOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnMoveOptions{}, "{}")

	u := &ProjectColumnMoveOptions{
		Position: "pos",
	}

	want := `{
		"position": "pos"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCard_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCard{}, "{}")

	u := &ProjectCard{
		URL:        String("url"),
		ColumnURL:  String("curl"),
		ContentURL: String("conurl"),
		ID:         Int64(1),
		Note:       String("note"),
		Creator: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		CreatedAt:          &Timestamp{referenceTime},
		UpdatedAt:          &Timestamp{referenceTime},
		NodeID:             String("nid"),
		Archived:           Bool(true),
		ColumnID:           Int64(1),
		ProjectID:          Int64(1),
		ProjectURL:         String("purl"),
		ColumnName:         String("cn"),
		PreviousColumnName: String("pcn"),
	}

	want := `{
		"url": "url",
		"column_url": "curl",
		"content_url": "conurl",
		"id": 1,
		"note": "note",
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"node_id": "nid",
		"archived": true,
		"column_id": 1,
		"project_id": 1,
		"project_url": "purl",
		"column_name": "cn",
		"previous_column_name": "pcn"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCardOptions{}, "{}")

	u := &ProjectCardOptions{
		Note:        "note",
		ContentID:   1,
		ContentType: "ct",
		Archived:    Bool(false),
	}

	want := `{
		"note": "note",
		"content_id": 1,
		"content_type": "ct",
		"archived": false
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardMoveOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCardMoveOptions{}, "{}")

	u := &ProjectCardMoveOptions{
		Position: "pos",
		ColumnID: 1,
	}

	want := `{
		"position": "pos",
		"column_id": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCollaboratorOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCollaboratorOptions{}, "{}")

	u := &ProjectCollaboratorOptions{
		Permission: String("per"),
	}

	want := `{
		"permission": "per"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectPermissionLevel_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectPermissionLevel{}, "{}")

	u := &ProjectPermissionLevel{
		Permission: String("per"),
		User: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
	}

	want := `{
		"permission": "per",
		"user": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		}
	}`

	testJSONMarshal(t, u, want)
}
