// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestProjectsService_ListProjectsForOrg(t *testing.T) {
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

	opts := &ListProjectsOptions{Query: "alpha", ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: "2", Before: "1"}}
	ctx := t.Context()
	projects, _, err := client.Projects.ListProjectsForOrg(ctx, "o", opts)
	if err != nil {
		t.Fatalf("Projects.ListProjectsForOrg returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 1 || projects[0].GetTitle() != "T1" {
		t.Fatalf("Projects.ListProjectsForOrg returned %+v", projects)
	}

	const methodName = "ListProjectsForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectsForOrg(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass := context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectsForOrg(ctxBypass, "o", &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_GetProjectForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"title":"OrgProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	project, _, err := client.Projects.GetProjectForOrg(ctx, "o", 1)
	if err != nil {
		t.Fatalf("Projects.GetProjectForOrg returned error: %v", err)
	}
	if project.GetID() != 1 || project.GetTitle() != "OrgProj" {
		t.Fatalf("Projects.GetProjectForOrg returned %+v", project)
	}

	const methodName = "GetProjectForOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectForOrg(ctx, "o", 1)
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

	opts := &ListProjectsOptions{Query: "beta", ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "1", After: "2", PerPage: 2}}
	ctx := t.Context()
	var ctxBypass context.Context
	projects, _, err := client.Projects.ListProjectsForUser(ctx, "u", opts)
	if err != nil {
		t.Fatalf("Projects.ListProjectsForUser returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 2 || projects[0].GetTitle() != "UProj" {
		t.Fatalf("Projects.ListProjectsForUser returned %+v", projects)
	}

	const methodName = "ListProjectsForUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectsForUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectsForUser(ctx, "u", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass = context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectsForUser(ctxBypass, "u", &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_GetProjectForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":2,"title":"UProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	project, _, err := client.Projects.GetProjectForUser(ctx, "u", 2)
	if err != nil {
		t.Fatalf("Projects.GetProjectForUser returned error: %v", err)
	}
	if project.GetID() != 2 || project.GetTitle() != "UProj" {
		t.Fatalf("Projects.GetProjectForUser returned %+v", project)
	}

	const methodName = "GetProjectForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectForUser(ctx, "u", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListProjectsForOrg_pagination(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// First page returns a Link header with rel="next" containing an after cursor (after=cursor2)
	mux.HandleFunc("/orgs/o/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		after := q.Get("after")
		before := q.Get("before")
		if after == "" && before == "" {
			// first request
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2?after=cursor2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":1,"title":"P1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "cursor2" {
			// second request simulates a previous link
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2?before=cursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":2,"title":"P2","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		// unexpected state
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})

	ctx := t.Context()
	first, resp, err := client.Projects.ListProjectsForOrg(ctx, "o", nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].GetID() != 1 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "cursor2" {
		t.Fatalf("expected resp.After=cursor2 got %q", resp.After)
	}

	// Use resp.After as opts.After for next page
	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: resp.After}}
	second, resp2, err := client.Projects.ListProjectsForOrg(ctx, "o", opts)
	if err != nil {
		t.Fatalf("second page error: %v", err)
	}
	if len(second) != 1 || second[0].GetID() != 2 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "cursor2" {
		t.Fatalf("expected resp2.Before=cursor2 got %q", resp2.Before)
	}
}

func TestProjectsService_ListProjectsForUser_pagination(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		after := q.Get("after")
		before := q.Get("before")
		if after == "" && before == "" { // first page
			w.Header().Set("Link", "<http://example.org/users/u/projectsV2?after=ucursor2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":10,"title":"UP1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "ucursor2" { // second page provides prev
			w.Header().Set("Link", "<http://example.org/users/u/projectsV2?before=ucursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":11,"title":"UP2","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})

	ctx := t.Context()
	first, resp, err := client.Projects.ListProjectsForUser(ctx, "u", nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].GetID() != 10 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "ucursor2" {
		t.Fatalf("expected resp.After=ucursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: resp.After}}
	second, resp2, err := client.Projects.ListProjectsForUser(ctx, "u", opts)
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

func TestProjectsService_ListProjectFieldsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if q.Get("before") == "b" && q.Get("after") == "a" {
			fmt.Fprint(w, `[]`)
			return
		}
		testFormValues(t, r, values{"q": "text", "after": "2", "before": "1"})
		fmt.Fprint(w, `[
			{
				"id": 1,
				"node_id": "node_1",
				"name": "Status",
				"data_type": "single_select",
				"project_url": "https://api.github.com/projects/1",
				"url": "https://api.github.com/projects/1/fields/field1",
				"options": [
					{
						"id": "option1",
						"name": "Todo",
						"color": "blue",
						"description": "Tasks to be done"
					},
					{
						"id": "option2",
						"name": "In Progress",
						"color": "yellow"
					}
				],
				"created_at": "2011-01-02T15:04:05Z",
				"updated_at": "2012-01-02T15:04:05Z"
			},
			{
				"id": 2,
				"node_id": "node_2",
				"name": "Priority",
				"data_type": "text",
				"project_url": "https://api.github.com/projects/1",
				"url": "https://api.github.com/projects/1/fields/field2",
				"created_at": "2011-01-02T15:04:05Z",
				"updated_at": "2012-01-02T15:04:05Z"
			}
		]`)
	})

	opts := &ListProjectsOptions{Query: "text", ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: "2", Before: "1"}}
	ctx := t.Context()
	fields, _, err := client.Projects.ListProjectFieldsForOrg(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListProjectFieldsForOrg returned error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Projects.ListProjectFieldsForOrg returned %v fields, want 2", len(fields))
	}

	field1 := fields[0]
	if field1.ID != 1 || field1.Name != "Status" || field1.DataType != "single_select" {
		t.Errorf("First field: got ID=%v, Name=%v, DataType=%v; want 1, Status, single_select", field1.ID, field1.Name, field1.DataType)
	}
	if len(field1.Options) != 2 {
		t.Errorf("First field options: got %v, want 2", len(field1.Options))
	} else {
		name0, name1 := field1.Options[0].Name, field1.Options[1].Name
		if field1.Options[0].Name != "Todo" || field1.Options[1].Name != "In Progress" {
			t.Errorf("First field option names: got %q, %q; want Todo, In Progress", name0, name1)
		}
	}

	// Validate second field (without options)
	field2 := fields[1]
	if field2.ID != 2 || field2.Name != "Priority" || field2.DataType != "text" {
		t.Errorf("Second field: got ID=%v, Name=%v, DataType=%v; want 2, Priority, text", field2.ID, field2.Name, field2.DataType)
	}
	if len(field2.Options) != 0 {
		t.Errorf("Second field options: got %v, want 0", len(field2.Options))
	}

	const methodName = "ListProjectFieldsForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectFieldsForOrg(ctx, "\n", 1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectFieldsForOrg(ctx, "o", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass := context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectFieldsForOrg(ctxBypass, "o", 1, &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_ListProjectFieldsForOrg_pagination(t *testing.T) {
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
			fmt.Fprint(w, `[{"id":1,"name":"Status","dataType":"single_select","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		if after == "cursor2" {
			// second request simulates a previous link
			w.Header().Set("Link", "<http://example.org/orgs/o/projectsV2/1/fields?before=cursor2>; rel=\"prev\"")
			fmt.Fprint(w, `[{"id":2,"name":"Priority","dataType":"text","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
			return
		}
		// unexpected state
		http.Error(w, "unexpected query", http.StatusBadRequest)
	})

	ctx := t.Context()
	first, resp, err := client.Projects.ListProjectFieldsForOrg(ctx, "o", 1, nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].ID != 1 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "cursor2" {
		t.Fatalf("expected resp.After=cursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: resp.After}}
	second, resp2, err := client.Projects.ListProjectFieldsForOrg(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("second page error: %v", err)
	}
	if len(second) != 1 || second[0].ID != 2 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "cursor2" {
		t.Fatalf("expected resp2.Before=cursor2 got %q", resp2.Before)
	}
}

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

func TestProjectV2Field_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2Field{}, "{}")       // empty struct
	testJSONMarshal(t, &ProjectV2FieldOption{}, "{}") // option struct still individually testable

	optVal := &ProjectV2FieldOption{Color: "blue", Description: "Tasks to be done", ID: "option1", Name: "Todo"}

	field := &ProjectV2Field{
		ID:         1,
		NodeID:     Ptr("node_1"),
		Name:       "Status",
		DataType:   "single_select",
		URL:        Ptr("https://api.github.com/projects/1/fields/field1"),
		ProjectURL: "https://api.github.com/projects/1",
		Options:    []*ProjectV2FieldOption{optVal},
		CreatedAt:  referenceTime,
		UpdatedAt:  referenceTime,
	}

	want := `{
			"id": 1,
	        "node_id": "node_1",
	        "name": "Status",
	        "data_type": "single_select",
	        "url": "https://api.github.com/projects/1/fields/field1",
	        "project_url": "https://api.github.com/projects/1",
	        "options": [
	            {
	                "id": "option1",
	                "name": "Todo",
	                "color": "blue",
	                "description": "Tasks to be done"
	            }
	        ],
	        "created_at": ` + referenceTimeStr + `,
	        "updated_at": ` + referenceTimeStr + `
	    }`

	testJSONMarshal(t, field, want)
}
