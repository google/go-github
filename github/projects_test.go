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

	mux.HandleFunc("/users/u/projectsV2/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3,"title":"UserProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	project, _, err := client.Projects.GetProjectForUser(ctx, "u", 3)
	if err != nil {
		t.Fatalf("Projects.GetProjectForUser returned error: %v", err)
	}
	if project.GetID() != 3 || project.GetTitle() != "UserProj" {
		t.Fatalf("Projects.GetProjectForUser returned %+v", project)
	}

	const methodName = "GetProjectForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectForUser(ctx, "u", 3)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListProjectFieldsForOrg(t *testing.T) {
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
			"dataType": "single_select",
			"url": "https://api.github.com/projects/1/fields/field1",
			"options": [
				{"id": 1, "name": "Todo", "color": "blue", "description": "Tasks to be done"},
				{"id": 2, "name": "In Progress", "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		},
		{
			"id": 2,
			"node_id": "node_2",
			"name": "Priority",
			"dataType": "text",
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
		t.Fatalf("Projects.ListProjectFieldsForOrg returned %d fields, want 2", len(fields))
	}
	if fields[0].ID == nil || *fields[0].ID != 1 || fields[1].ID == nil || *fields[1].ID != 2 {
		t.Fatalf("unexpected field IDs: %+v", fields)
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
	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectFieldsForOrg(ctxBypass, "o", 1, &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_ListProjectFieldsForUser(t *testing.T) {
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
			"dataType": "single_select",
			"url": "https://api.github.com/projects/1/fields/field1",
			"options": [
				{"id": 1, "name": "Todo", "color": "blue", "description": "Tasks to be done"},
				{"id": 2, "name": "In Progress", "color": "yellow"}
			],
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		},
		{
			"id": 2,
			"node_id": "node_2",
			"name": "Priority",
			"dataType": "text",
			"url": "https://api.github.com/projects/1/fields/field2",
			"created_at": "2011-01-02T15:04:05Z",
			"updated_at": "2012-01-02T15:04:05Z"
		}
		]`)
	})

	opts := &ListProjectsOptions{Query: "text", ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: "2", Before: "1"}}
	ctx := t.Context()
	fields, _, err := client.Projects.ListProjectFieldsForUser(ctx, "u", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListProjectFieldsForUser returned error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("Projects.ListProjectFieldsForUser returned %d fields, want 2", len(fields))
	}
	if fields[0].ID == nil || *fields[0].ID != 1 || fields[1].ID == nil || *fields[1].ID != 2 {
		t.Fatalf("unexpected field IDs: %+v", fields)
	}

	const methodName = "ListProjectFieldsForUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectFieldsForUser(ctx, "\n", 1, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectFieldsForUser(ctx, "u", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectFieldsForUser(ctxBypass, "u", 1, &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_ListProjectsForUser_pagination(t *testing.T) {
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

func TestProjectsService_ListProjectsForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})
	ctx := t.Context()
	const methodName = "ListProjectsForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectsForUser(ctx, "u", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	// bad options (bad username) should error
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectsForUser(ctx, "\n", nil)
		return err
	})
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
	if len(first) != 1 || first[0].ID == nil || *first[0].ID != 1 {
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
	if len(second) != 1 || second[0].ID == nil || *second[0].ID != 2 {
		t.Fatalf("unexpected second page %+v", second)
	}
	if resp2.Before != "cursor2" {
		t.Fatalf("expected resp2.Before=cursor2 got %q", resp2.Before)
	}
}

func TestProjectsService_ListProjectsForOrg_pagination(t *testing.T) {
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
	first, resp, err := client.Projects.ListProjectsForOrg(ctx, "o", nil)
	if err != nil {
		t.Fatalf("first page error: %v", err)
	}
	if len(first) != 1 || first[0].GetID() != 20 {
		t.Fatalf("unexpected first page %+v", first)
	}
	if resp.After != "ocursor2" {
		t.Fatalf("expected resp.After=ocursor2 got %q", resp.After)
	}

	opts := &ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: resp.After}}
	second, resp2, err := client.Projects.ListProjectsForOrg(ctx, "o", opts)
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
		ID:       Ptr(int64(2)),
		NodeID:   "node_1",
		Name:     "Status",
		DataType: "single_select",
		URL:      "https://api.github.com/projects/1/fields/field1",
		Options: []*ProjectV2FieldOption{
			{
				ID:          Ptr(int64(1)),
				Name:        "Todo",
				Color:       "blue",
				Description: "Tasks to be done",
			},
		},
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
	}

	want := `{
        "id": 2,
        "node_id": "node_1",
        "name": "Status",
        "dataType": "single_select",
        "url": "https://api.github.com/projects/1/fields/field1",
        "options": [
            {
                "id": 1,
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

func TestProjectsService_ListProjectItemsForOrg(t *testing.T) {
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

	opts := &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{After: "2", Before: "1", PerPage: 50}, Query: "status:open"}, Fields: []int64{10, 11}}
	ctx := t.Context()
	items, _, err := client.Projects.ListProjectItemsForOrg(ctx, "o", 1, opts)
	if err != nil {
		t.Fatalf("Projects.ListProjectItemsForOrg returned error: %v", err)
	}
	if len(items) != 1 || items[0].GetID() != 17 {
		t.Fatalf("Projects.ListProjectItemsForOrg returned %+v", items)
	}

	const methodName = "ListProjectItemsForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectItemsForOrg(ctx, "\n", 1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectItemsForOrg(ctx, "o", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	ctxBypass := context.WithValue(ctx, BypassRateLimitCheck, true)
	if _, _, err = client.Projects.ListProjectItemsForOrg(ctxBypass, "o", 1, &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{Before: "b", After: "a"}}}); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestProjectsService_AddProjectItemForOrg(t *testing.T) {
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
	item, _, err := client.Projects.AddProjectItemForOrg(ctx, "o", 1, &AddProjectItemOptions{Type: "Issue", ID: 99})
	if err != nil {
		t.Fatalf("Projects.AddProjectItemForOrg returned error: %v", err)
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
	const methodName = "AddProjectItemForOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.AddProjectItemForOrg(ctx, "o", 1, &AddProjectItemOptions{Type: "Issue", ID: 1})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetProjectItemForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":17,"node_id":"PVTI_node"}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{}
	item, _, err := client.Projects.GetProjectItemForOrg(ctx, "o", 1, 17, opts)
	if err != nil {
		t.Fatalf("GetProjectItemForOrg error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_GetProjectItemForOrg_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":17}`)
	})
	ctx := t.Context()
	const methodName = "GetProjectItemForOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectItemForOrg(ctx, "o", 1, 17, &GetProjectItemOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateProjectItemForOrg(t *testing.T) {
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
	item, _, err := client.Projects.UpdateProjectItemForOrg(ctx, "o", 1, 17, &UpdateProjectItemOptions{Archived: &archived})
	if err != nil {
		t.Fatalf("UpdateProjectItemForOrg error: %v", err)
	}
	if item.GetID() != 17 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_UpdateProjectItemForOrg_error(t *testing.T) {
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
		got, resp, err := client.Projects.UpdateProjectItemForOrg(ctx, "o", 1, 17, &UpdateProjectItemOptions{Archived: &archived})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteProjectItemForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	if _, err := client.Projects.DeleteProjectItemForOrg(ctx, "o", 1, 17); err != nil {
		t.Fatalf("DeleteProjectItemForOrg error: %v", err)
	}
}

func TestProjectsService_DeleteProjectItemForOrg_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/orgs/o/projectsV2/1/items/17", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	const methodName = "DeleteProjectItemForOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteProjectItemForOrg(ctx, "o", 1, 17)
	})
}

func TestProjectsService_ListProjectItemsForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "20", "q": "type:issue"})
		fmt.Fprint(w, `[{"id":7,"node_id":"PVTI_user"}]`)
	})
	ctx := t.Context()
	items, _, err := client.Projects.ListProjectItemsForUser(ctx, "u", 2, &ListProjectItemsOptions{ListProjectsOptions: ListProjectsOptions{ListProjectsPaginationOptions: ListProjectsPaginationOptions{PerPage: 20}, Query: "type:issue"}})
	if err != nil {
		t.Fatalf("ListProjectItemsForUser error: %v", err)
	}
	if len(items) != 1 || items[0].GetID() != 7 {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestProjectsService_ListProjectItemsForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[]`)
	})
	ctx := t.Context()
	const methodName = "ListProjectItemsForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListProjectItemsForUser(ctx, "u", 2, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListProjectItemsForUser(ctx, "\n", 2, nil)
		return err
	})
}

func TestProjectsService_AddProjectItemForUser(t *testing.T) {
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
	item, _, err := client.Projects.AddProjectItemForUser(ctx, "u", 2, &AddProjectItemOptions{Type: "PullRequest", ID: 123})
	if err != nil {
		t.Fatalf("AddProjectItemForUser error: %v", err)
	}
	if item.GetID() != 123 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_AddProjectItemForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":5}`)
	})
	ctx := t.Context()
	const methodName = "AddProjectItemForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.AddProjectItemForUser(ctx, "u", 2, &AddProjectItemOptions{Type: "Issue", ID: 5})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetProjectItemForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":55,"node_id":"PVTI_user_item"}`)
	})
	ctx := t.Context()
	opts := &GetProjectItemOptions{}
	item, _, err := client.Projects.GetProjectItemForUser(ctx, "u", 2, 55, opts)
	if err != nil {
		t.Fatalf("GetProjectItemForUser error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_GetProjectItemForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":55}`)
	})
	ctx := t.Context()
	const methodName = "GetProjectItemForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetProjectItemForUser(ctx, "u", 2, 55, &GetProjectItemOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_UpdateProjectItemForUser(t *testing.T) {
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
	item, _, err := client.Projects.UpdateProjectItemForUser(ctx, "u", 2, 55, &UpdateProjectItemOptions{Archived: &archived})
	if err != nil {
		t.Fatalf("UpdateProjectItemForUser error: %v", err)
	}
	if item.GetID() != 55 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestProjectsService_UpdateProjectItemForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":55}`)
	})
	archived := false
	ctx := t.Context()
	const methodName = "UpdateProjectItemForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.UpdateProjectItemForUser(ctx, "u", 2, 55, &UpdateProjectItemOptions{Archived: &archived})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_DeleteProjectItemForUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	if _, err := client.Projects.DeleteProjectItemForUser(ctx, "u", 2, 55); err != nil {
		t.Fatalf("DeleteProjectItemForUser error: %v", err)
	}
}

func TestProjectsService_DeleteProjectItemForUser_error(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/users/u/projectsV2/2/items/55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	const methodName = "DeleteProjectItemForUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Projects.DeleteProjectItemForUser(ctx, "u", 2, 55)
	})
}
