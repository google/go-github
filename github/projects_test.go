// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestProjectsService_ListOrganizationProjects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// Combined handler: supports initial test case and dual before/after validation scenario.
	mux.HandleFunc("/orgs/o/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" {
			fmt.Fprint(w, `[]`)
			return
		}
		// default expectation for main part of test
		testFormValues(t, r, values{"q": "alpha", "after": "2", "before": "1"})
		fmt.Fprint(w, `[{"id":1,"title":"T1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
	})

	opts := &ListProjectsOptions{Query: Ptr("alpha"), ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr("2"), Before: Ptr("1")}}
	ctx := t.Context()
	projects, _, err := client.Projects.ListOrganizationProjects(ctx, "o", opts)
	if err != nil {
		t.Fatalf("Projects.ListOrganizationProjects returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 1 || projects[0].GetTitle() != "T1" {
		t.Fatalf("Projects.ListOrganizationProjects returned %+v", projects)
	}

	const methodName = "ListOrganizationProjects"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListOrganizationProjects(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListOrganizationProjects(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass := context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListOrganizationProjects(ctxBypass, "o", &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("b"), After: Ptr("a")}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_GetOrganizationProject(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"title":"OrgProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	project, _, err := client.Projects.GetOrganizationProject(ctx, "o", 1)
	if err != nil {
		t.Fatalf("Projects.GetOrganizationProject returned error: %v", err)
	}
	if project.GetID() != 1 || project.GetTitle() != "OrgProj" {
		t.Fatalf("Projects.GetOrganizationProject returned %+v", project)
	}

	const methodName = "GetOrganizationProject"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetOrganizationProject(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListUserProjects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// Combined handler: supports initial test case and dual before/after scenario.
	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" {
			fmt.Fprint(w, `[]`)
			return
		}
		testFormValues(t, r, values{"q": "beta", "before": "1", "after": "2", "per_page": "2"})
		fmt.Fprint(w, `[{"id":2,"title":"UProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
	})

	opts := &ListProjectsOptions{Query: Ptr("beta"), ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("1"), After: Ptr("2"), PerPage: Ptr(2)}}
	ctx := t.Context()
	var ctxBypass context.Context
	projects, _, err := client.Projects.ListUserProjects(ctx, "u", opts)
	if err != nil {
		t.Fatalf("Projects.ListUserProjects returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 2 || projects[0].GetTitle() != "UProj" {
		t.Fatalf("Projects.ListUserProjects returned %+v", projects)
	}

	const methodName = "ListUserProjects"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListUserProjects(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListUserProjects(ctx, "u", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass = context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListUserProjects(ctxBypass, "u", &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("b"), After: Ptr("a")}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_GetUserProject(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3,"title":"UserProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	project, _, err := client.Projects.GetUserProject(ctx, "u", 3)
	if err != nil {
		t.Fatalf("Projects.GetUserProject returned error: %v", err)
	}
	if project.GetID() != 3 || project.GetTitle() != "UserProj" {
		t.Fatalf("Projects.GetUserProject returned %+v", project)
	}

	const methodName = "GetUserProject"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetUserProject(ctx, "u", 3)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListOrganizationProjectFields(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" { // bypass scenario
			fmt.Fprint(w, `[]`)
			return
		}
		testFormValues(t, r, values{"after": "2", "before": "1", "q": "text"})
		fmt.Fprint(w, `[
		{
			"id": 1,
			"node_id": "node_1",
			"name": "Status",
			"data_type": "single_select",
			"url": "https://api.github.com/projects/1/fields/field1",
			"options": [
				{"id": "1", "name": {"raw": "Todo", "html": "Todo"}, "color": "blue", "description": {"raw": "Tasks to be done", "html": "Tasks to be done"}},
				{"id": "2", "name": {"raw": "In Progress", "html": "In Progress"}, "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		},
		{
			"id": 2,
			"node_id": "node_2",
			"name": "Priority",
			"data_type": "text",
			"url": "https://api.github.com/projects/1/fields/field2",
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		}
		]`)
	})

	opts := &ListProjectsOptions{Query: Ptr("text"), ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr("2"), Before: Ptr("1")}}
	ctx := t.Context()
	fields, _, err := client.Projects.ListOrganizationProjectFields(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListOrganizationProjectFields returned error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("Projects.ListOrganizationProjectFields returned %d fields, want 2", len(fields))
	}
	if fields[0].ID == nil || *fields[0].ID != 1 || fields[1].ID == nil || *fields[1].ID != 2 {
		t.Fatalf("unexpected field IDs: %+v", fields)
	}

	const methodName = "ListOrganizationProjectFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListOrganizationProjectFields(ctx, "\n", 1, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListOrganizationProjectFields(ctx, "o", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListOrganizationProjectFields(ctxBypass, "o", 1, &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("b"), After: Ptr("a")}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_ListUserProjectFields(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2/1/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" { // bypass scenario
			fmt.Fprint(w, `[]`)
			return
		}
		testFormValues(t, r, values{"after": "2", "before": "1", "q": "text"})
		fmt.Fprint(w, `[
		{
			"id": 1,
			"node_id": "node_1",
			"name": "Status",
			"data_type": "single_select",
			"url": "https://api.github.com/projects/1/fields/field1",
			"options": [
				{"id": "1", "name": {"raw": "Todo", "html": "Todo"}, "color": "blue", "description": {"raw": "Tasks to be done", "html": "Tasks to be done"}},
				{"id": "2", "name": {"raw": "In Progress", "html": "In Progress"}, "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		},
		{
			"id": 2,
			"node_id": "node_2",
			"name": "Priority",
			"data_type": "text",
			"url": "https://api.github.com/projects/1/fields/field2",
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		}
		]`)
	})

	opts := &ListProjectsOptions{Query: Ptr("text"), ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr("2"), Before: Ptr("1")}}
	ctx := t.Context()
	fields, _, err := client.Projects.ListUserProjectFields(ctx, "u", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListUserProjectFields returned error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("Projects.ListUserProjectFields returned %d fields, want 2", len(fields))
	}
	if fields[0].ID == nil || *fields[0].ID != 1 || fields[1].ID == nil || *fields[1].ID != 2 {
		t.Fatalf("unexpected field IDs: %+v", fields)
	}

	const methodName = "ListUserProjectFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListUserProjectFields(ctx, "\n", 1, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListUserProjectFields(ctx, "u", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListUserProjectFields(ctxBypass, "u", 1, &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("b"), After: Ptr("a")}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_GetOrganizationProjectField(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1/fields/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		{
			"id": 1,
			"node_id": "node_1",
			"name": "Status",
			"data_type": "single_select",
			"url": "https://api.github.com/projects/1/fields/field1",
			"options": [
				{"id": "1", "name": {"raw": "Todo", "html": "Todo"}, "color": "blue", "description": {"raw": "Tasks to be done", "html": "Tasks to be done"}},
				{"id": "2", "name": {"raw": "In Progress", "html": "In Progress"}, "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		}`)
	})

	ctx := t.Context()
	field, _, err := client.Projects.GetOrganizationProjectField(ctx, "o", 1, 1)
	if err != nil {
		t.Fatalf("Projects.GetOrganizationProjectField returned error: %v", err)
	}
	if field == nil || field.ID == nil || *field.ID != 1 {
		t.Fatalf("unexpected field: %+v", field)
	}

	const methodName = "GetOrganizationProjectField"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetOrganizationProjectField(ctx, "o", 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetUserProjectField(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2/1/fields/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		{
			"id": 3,
			"node_id": "node_3",
			"name": "Status",
			"data_type": "single_select",
			"url": "https://api.github.com/projects/1/fields/field3",
			"options": [
				{"id": "1", "name": {"raw": "Done", "html": "Done"}, "color": "red", "description": {"raw": "Done task", "html": "Done task"}},
				{"id": "2", "name": {"raw": "In Progress", "html": "In Progress"}, "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		}`)
	})

	ctx := t.Context()
	field, _, err := client.Projects.GetUserProjectField(ctx, "u", 1, 3)
	if err != nil {
		t.Fatalf("Projects.GetUserProjectField returned error: %v", err)
	}
	if field == nil || field.ID == nil || *field.ID != 3 {
		t.Fatalf("unexpected field: %+v", field)
	}

	const methodName = "GetUserProjectField"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetUserProjectField(ctx, "u", 1, 3)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListUserProjects_pagination(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		after := q.Get("after")
		before := q.Get("before")
		if after == "" && before == "" {
			w.Header().Set("Link", "<http://example.org/users/u/projectsV2?after=ucursor2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":10,"title":"UP1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "ucursor2" {
			w.Header().Set("Link", "<http://example.org/users/u/projectsV2?before=ucursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":11,"title":"UP2","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})
	ctx := t.Context()
	first, resp, err := client.Projects.ListUserProjects(ctx, "u", nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].GetID() != 10 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "ucursor2" {
		t.Fatalf("expected resp.After=ucursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr(resp.After)}}
	second, resp2, err := client.Projects.ListUserProjects(ctx, "u", opts)
	if err != nil {
		t.Fatalf("second page error: %v", err)
	}
	if len(second) != 1 || second[0].GetID() != 11 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "ucursor2" {
		t.Fatalf("expected resp2.Before=ucursor2 got %q", resp2.Before)
	}
}

func TestProjectsService_ListUserProjects_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})
	ctx := t.Context()
	const methodName = "ListUserProjects"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListUserProjects(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	// bad options (bad username) should error
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListUserProjects(ctx, "\n", nil)
		return err
	})
}

func TestProjectsService_ListOrganizationProjectFields_pagination(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// First page returns a Link header with rel="next" containing an after cursor
	mux.HandleFunc("/orgs/o/projectsV2/1/fields", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		after := q.Get("after")
		before := q.Get("before")
		if after == "" && before == "" {
			// first request
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2/1/fields?after=cursor2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":1,"name":"Status","data_type":"single_select","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "cursor2" {
			// second request simulates a previous link
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2/1/fields?before=cursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":2,"name":"Priority","data_type":"text","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		// unexpected state
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})

	ctx := t.Context()
	first, resp, err := client.Projects.ListOrganizationProjectFields(ctx, "o", 1, nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].ID == nil || *first[0].ID != 1 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "cursor2" {
		t.Fatalf("expected resp.After=cursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr(resp.After)}}
	second, resp2, err := client.Projects.ListOrganizationProjectFields(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("second page error: %v", err)
	}
	if len(second) != 1 || second[0].ID == nil || *second[0].ID != 2 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "cursor2" {
		t.Fatalf("expected resp2.Before=cursor2 got %q", resp2.Before)
	}
}

func TestProjectsService_ListOrganizationProjects_pagination(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		after := q.Get("after")
		before := q.Get("before")
		if after == "" && before == "" {
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2?after=ocursor2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":20,"title":"OP1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "ocursor2" {
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2?before=ocursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":21,"title":"OP2","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})

	ctx := t.Context()
	first, resp, err := client.Projects.ListOrganizationProjects(ctx, "o", nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].GetID() != 20 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "ocursor2" {
		t.Fatalf("expected resp.After=ocursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr(resp.After)}}
	second, resp2, err := client.Projects.ListOrganizationProjects(ctx, "o", opts)
	if err != nil {
		t.Fatalf("second page error: %v", err)
	}
	if len(second) != 1 || second[0].GetID() != 21 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "ocursor2" {
		t.Fatalf("expected resp2.Before=ocursor2 got %q", resp2.Before)
	}
}

// Marshal test ensures V2 fields marshal correctly.
func TestProjectV2_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2{}, "{}")

	p := &ProjectV2{
		ID:          Ptr(int64(10)),
		Title:       Ptr("Title"),
		Description: Ptr("Desc"),
		Public:      Ptr(true),
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
	}

	want := `{
        "id": 10,
        "title": "Title",
        "description": "Desc",
        "public": true,
        "created_at": ` + referenceTimeStr + `,
        "updated_at": ` + referenceTimeStr + `
    }`

	testJSONMarshal(t, p, want)
}

// Marshal test ensures V2 field structures marshal correctly.
func TestProjectV2Field_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2Field{}, "{}")
	testJSONMarshal(t, &ProjectV2FieldOption{}, "{}")

	field := &ProjectV2Field{
		ID:         Ptr(int64(2)),
		NodeID:     Ptr("node_1"),
		Name:       Ptr("Status"),
		DataType:   Ptr("single_select"),
		ProjectURL: Ptr("https://api.github.com/projects/67890"),
		Options: []*ProjectV2FieldOption{
			{
				ID:          Ptr("1"),
				Name:        &ProjectV2TextContent{Raw: Ptr("Todo"), HTML: Ptr("Todo")},
				Color:       Ptr("blue"),
				Description: &ProjectV2TextContent{Raw: Ptr("Tasks to be done"), HTML: Ptr("Tasks to be done")},
			},
		},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
	}

	want := `{
        "id": 2,
        "node_id": "node_1",
        "name": "Status",
        "data_type": "single_select",
        "project_url": "https://api.github.com/projects/67890",
        "options": [
            {
                "id": "1",
                "color": "blue",
                "description": {
                    "raw": "Tasks to be done",
                    "html": "Tasks to be done"
                },
                "name": {
                    "raw": "Todo",
                    "html": "Todo"
                }
            }
        ],
        "created_at": ` + referenceTimeStr + `,
        "updated_at": ` + referenceTimeStr + `
    }`

	testJSONMarshal(t, field, want)
}

// Marshal test ensures ProjectV2FieldConfiguration marshals correctly.
func TestProjectV2FieldConfiguration_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2FieldConfiguration{}, "{}")
	testJSONMarshal(t, &ProjectV2FieldIteration{}, "{}")

	// Test a field with configuration (iteration field)
	fieldWithConfiguration := &ProjectV2Field{
		ID:         Ptr(int64(3)),
		NodeID:     Ptr("node_3"),
		Name:       Ptr("Sprint"),
		DataType:   Ptr("iteration"),
		ProjectURL: Ptr("https://api.github.com/projects/67890"),
		Configuration: &ProjectV2FieldConfiguration{
			Duration: Ptr(1209600), // 2 weeks in seconds
			StartDay: Ptr(1),       // Monday
			Iterations: []*ProjectV2FieldIteration{
				{
					ID:        Ptr("iter_1"),
					Title:     &ProjectV2TextContent{Raw: Ptr("Sprint 1"), HTML: Ptr("Sprint 1")},
					StartDate: Ptr("2025-01-06"),
					Duration:  Ptr(1209600),
				},
				{
					ID:        Ptr("iter_2"),
					Title:     &ProjectV2TextContent{Raw: Ptr("Sprint 2"), HTML: Ptr("Sprint 2")},
					StartDate: Ptr("2025-01-20"),
					Duration:  Ptr(1209600),
				},
			},
		},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
	}

	want := `{
        "id": 3,
        "node_id": "node_3",
        "name": "Sprint",
        "data_type": "iteration",
        "project_url": "https://api.github.com/projects/67890",
        "configuration": {
            "duration": 1209600,
            "start_day": 1,
            "iterations": [
                {
                    "id": "iter_1",
                    "title": {
                        "raw": "Sprint 1",
                        "html": "Sprint 1"
                    },
                    "start_date": "2025-01-06",
                    "duration": 1209600
                },
                {
                    "id": "iter_2",
                    "title": {
                        "raw": "Sprint 2",
                        "html": "Sprint 2"
                    },
                    "start_date": "2025-01-20",
                    "duration": 1209600
                }
            ]
        },
        "created_at": ` + referenceTimeStr + `,
        "updated_at": ` + referenceTimeStr + `
    }`

	testJSONMarshal(t, fieldWithConfiguration, want)

	// Test just the configuration struct by itself
	config := &ProjectV2FieldConfiguration{
		Duration: Ptr(604800), // 1 week in seconds
		StartDay: Ptr(0),      // Sunday
		Iterations: []*ProjectV2FieldIteration{
			{
				ID:        Ptr("config_iter_1"),
				Title:     &ProjectV2TextContent{Raw: Ptr("Week 1"), HTML: Ptr("Week 1")},
				StartDate: Ptr("2025-01-01"),
				Duration:  Ptr(604800),
			},
		},
	}

	configWant := `{
        "duration": 604800,
        "start_day": 0,
        "iterations": [
            {
                "id": "config_iter_1",
                "title": {
                    "raw": "Week 1",
                    "html": "Week 1"
                },
                "start_date": "2025-01-01",
                "duration": 604800
            }
        ]
    }`

	testJSONMarshal(t, config, configWant)

	// Test iteration struct by itself
	iteration := &ProjectV2FieldIteration{
		ID:        Ptr("single_iter"),
		Title:     &ProjectV2TextContent{Raw: Ptr("Test Iteration"), HTML: Ptr("Test Iteration")},
		StartDate: Ptr("2025-02-01"),
		Duration:  Ptr(1209600),
	}

	iterationWant := `{
        "id": "single_iter",
        "title": {
            "raw": "Test Iteration",
            "html": "Test Iteration"
        },
        "start_date": "2025-02-01",
        "duration": 1209600
    }`

	testJSONMarshal(t, iteration, iterationWant)
}

func TestProjectsService_ListOrganizationProjectItems(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" { // bypass scenario
			fmt.Fprint(w, `[]`)
			return
		}
		testFormValues(t, r, values{"after": "2", "before": "1", "per_page": "50", "fields": "10,11", "q": "status:open"})
		fmt.Fprint(w, `[{"id":17,"node_id":"PVTI_node"}]`)
	})

	opts := &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: Ptr("2"), Before: Ptr("1"), PerPage: Ptr(50)}, Query: Ptr("status:open")}, Fields: []int64{10, 11}}
	ctx := t.Context()
	items, _, err := client.Projects.ListOrganizationProjectItems(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListOrganizationProjectItems returned error: %v", err)
	}
	if len(items) != 1 || items[0].GetID() != 17 {
		t.Fatalf("Projects.ListOrganizationProjectItems returned %+v", items)
	}

	const methodName = "ListOrganizationProjectItems"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListOrganizationProjectItems(ctx, "\n", 1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListOrganizationProjectItems(ctx, "o", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListOrganizationProjectItems(ctxBypass, "o", 1, &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: Ptr("b"), After: Ptr("a")}}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_AddOrganizationProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		if body != `{"type":"Issue","id":99}`+"\n" { // encoder adds newline
			t.Fatalf("unexpected body: %s", body)
		}
		fmt.Fprint(w, `{"id":99,"node_id":"PVTI_new"}`)
	})

	ctx := t.Context()
	item, _, err := client.Projects.AddOrganizationProjectItem(ctx, "o", 1, &AddProjectItemOptions{Type: "Issue", ID: 99})
	if err != nil {
		t.Fatalf("Projects.AddOrganizationProjectItem returned error: %v", err)
	}
	if item.GetID() != 99 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_AddProjectItemForOrg_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":1}`)
	})
	ctx := t.Context()
	const methodName = "AddOrganizationProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.AddOrganizationProjectItem(ctx, "o", 1, &AddProjectItemOptions{Type: "Issue", ID: 1})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetOrganizationProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":17,"node_id":"PVTI_node"}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{}
	item, _, err := client.Projects.GetOrganizationProjectItem(ctx, "o", 1, 17, opts)
	if err != nil {
		t.Fatalf("GetOrganizationProjectItem error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_GetOrganizationProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":17}`)
	})
	ctx := t.Context()
	const methodName = "GetOrganizationProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetOrganizationProjectItem(ctx, "o", 1, 17, &GetProjectItemOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetOrganizationProjectItem_WithFieldsOption(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// Verify that fields option is properly added as comma-separated URL parameter
		testFormValues(t, r, values{"fields": "123,456,789"})
		fmt.Fprint(w, `{
			"id":17,
			"node_id":"PVTI_node_fields",
			"fields":[
				{"id":123,"name":"Status","data_type":"single_select"},
				{"id":456,"name":"Priority","data_type":"single_select"},
				{"id":789,"name":"Assignee","data_type":"text"}
			]
		}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{
		Fields: []int64{123, 456, 789}, // Request specific field IDs
	}
	item, _, err := client.Projects.GetOrganizationProjectItem(ctx, "o", 1, 17, opts)
	if err != nil {
		t.Fatalf("GetOrganizationProjectItem error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}
	const methodName = "GetOrganizationProjectItemWithFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.GetOrganizationProjectItem(ctx, "\n", 1, 17, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetOrganizationProjectItem(ctx, "o", 1, 17, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateOrganizationProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		if body != `{"archived":true}`+"\n" {
			t.Fatalf("unexpected body: %s", body)
		}
		fmt.Fprint(w, `{"id":17}`)
	})
	archived := true
	ctx := t.Context()
	item, _, err := client.Projects.UpdateOrganizationProjectItem(ctx, "o", 1, 17, &UpdateProjectItemOptions{Archived: &archived})
	if err != nil {
		t.Fatalf("UpdateOrganizationProjectItem error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_UpdateOrganizationProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":17}`)
	})
	archived := true
	ctx := t.Context()
	const methodName = "UpdateProjectItemForOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateOrganizationProjectItem(ctx, "o", 1, 17, &UpdateProjectItemOptions{Archived: &archived})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateOrganizationProjectItem_WithFieldUpdates(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		// Verify the field updates are properly formatted in the request body
		expectedBody := `{"fields":[{"id":123,"value":"Updated text value"},{"id":456,"value":"Done"}]}`
		if body != expectedBody+"\n" {
			t.Fatalf("unexpected body: %s, expected: %s", body, expectedBody)
		}
		fmt.Fprint(w, `{"id":17,"node_id":"PVTI_node_updated"}`)
	})

	ctx := t.Context()
	opts := &UpdateProjectItemOptions{
		Fields: []*UpdateProjectV2Field{
			{ID: 123, Value: "Updated text value"},
			{ID: 456, Value: "Done"},
		},
	}
	item, _, err := client.Projects.UpdateOrganizationProjectItem(ctx, "o", 1, 17, opts)
	if err != nil {
		t.Fatalf("UpdateOrganizationProjectItem error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}

	const methodName = "UpdateOrganizationProjectItemWithFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.UpdateOrganizationProjectItem(ctx, "\n", 1, 17, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateOrganizationProjectItem(ctx, "o", 1, 17, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteOrganizationProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	if _, err := client.Projects.DeleteOrganizationProjectItem(ctx, "o", 1, 17); err != nil {
		t.Fatalf("DeleteOrganizationProjectItem error: %v", err)
	}
}

func TestProjectsService_DeleteOrganizationProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	const methodName = "DeleteOrganizationProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteOrganizationProjectItem(ctx, "o", 1, 17)
	})
}

func TestProjectsService_ListUserProjectItems(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "20", "q": "type:issue"})
		fmt.Fprint(w, `[{"id":7,"node_id":"PVTI_user"}]`)
	})
	ctx := t.Context()
	items, _, err := client.Projects.ListUserProjectItems(ctx, "u", 2, &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{PerPage: Ptr(20)}, Query: Ptr("type:issue")}})
	if err != nil {
		t.Fatalf("ListUserProjectItems error: %v", err)
	}
	if len(items) != 1 || items[0].GetID() != 7 {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestProjectsService_ListUserProjectItems_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})
	ctx := t.Context()
	const methodName = "ListUserProjectItems"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListUserProjectItems(ctx, "u", 2, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListUserProjectItems(ctx, "\n", 2, nil)
		return err
	})
}

func TestProjectsService_AddUserProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		if body != `{"type":"PullRequest","id":123}`+"\n" {
			t.Fatalf("unexpected body: %s", body)
		}
		fmt.Fprint(w, `{"id":123,"node_id":"PVTI_new_user"}`)
	})
	ctx := t.Context()
	item, _, err := client.Projects.AddUserProjectItem(ctx, "u", 2, &AddProjectItemOptions{Type: "PullRequest", ID: 123})
	if err != nil {
		t.Fatalf("AddUserProjectItem error: %v", err)
	}
	if item.GetID() != 123 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_AddUserProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":5}`)
	})
	ctx := t.Context()
	const methodName = "AddUserProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.AddUserProjectItem(ctx, "u", 2, &AddProjectItemOptions{Type: "Issue", ID: 5})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetUserProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":55,"node_id":"PVTI_user_item"}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{}
	item, _, err := client.Projects.GetUserProjectItem(ctx, "u", 2, 55, opts)
	if err != nil {
		t.Fatalf("GetUserProjectItem error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_GetUserProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":55}`)
	})
	ctx := t.Context()
	const methodName = "GetUserProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetUserProjectItem(ctx, "u", 2, 55, &GetProjectItemOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetUserProjectItem_WithFieldsOption(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// Verify that fields option is properly added as comma-separated URL parameter
		testFormValues(t, r, values{"fields": "100,200"})
		fmt.Fprint(w, `{
			"id":55,
			"node_id":"PVTI_user_item_fields",
			"fields":[
				{"id":100,"name":"Status","data_type":"single_select"},
				{"id":200,"name":"Milestone","data_type":"text"}
			]
		}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{
		Fields: []int64{100, 200}, // Request specific field IDs
	}
	item, _, err := client.Projects.GetUserProjectItem(ctx, "u", 2, 55, opts)
	if err != nil {
		t.Fatalf("GetUserProjectItem error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}

	const methodName = "GetUserProjectItemWithFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.GetUserProjectItem(ctx, "\n", 2, 55, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetUserProjectItem(ctx, "u", 2, 55, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateUserProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		if body != `{"archived":false}`+"\n" {
			t.Fatalf("unexpected body: %s", body)
		}
		fmt.Fprint(w, `{"id":55}`)
	})
	archived := false
	ctx := t.Context()
	item, _, err := client.Projects.UpdateUserProjectItem(ctx, "u", 2, 55, &UpdateProjectItemOptions{Archived: &archived})
	if err != nil {
		t.Fatalf("UpdateUserProjectItem error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_UpdateUserProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":55}`)
	})
	archived := false
	ctx := t.Context()
	const methodName = "UpdateUserProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateUserProjectItem(ctx, "u", 2, 55, &UpdateProjectItemOptions{Archived: &archived})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateUserProjectItem_WithFieldUpdates(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		// Verify the field updates are properly formatted in the request body
		expectedBody := `{"fields":[{"id":100,"value":"In Progress"},{"id":200,"value":5}]}`
		if body != expectedBody+"\n" {
			t.Fatalf("unexpected body: %s, expected: %s", body, expectedBody)
		}
		fmt.Fprint(w, `{"id":55,"node_id":"PVTI_user_updated"}`)
	})

	ctx := t.Context()
	opts := &UpdateProjectItemOptions{
		Fields: []*UpdateProjectV2Field{
			{ID: 100, Value: "In Progress"},
			{ID: 200, Value: 5}, // number field
		},
	}
	item, _, err := client.Projects.UpdateUserProjectItem(ctx, "u", 2, 55, opts)
	if err != nil {
		t.Fatalf("UpdateUserProjectItem error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}

	const methodName = "UpdateUserProjectItemWithFields"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.UpdateUserProjectItem(ctx, "\n", 2, 55, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateUserProjectItem(ctx, "u", 2, 55, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteUserProjectItem(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	if _, err := client.Projects.DeleteUserProjectItem(ctx, "u", 2, 55); err != nil {
		t.Fatalf("DeleteUserProjectItem error: %v", err)
	}
}

func TestProjectsService_DeleteUserProjectItem_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	const methodName = "DeleteUserProjectItem"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteUserProjectItem(ctx, "u", 2, 55)
	})
}
