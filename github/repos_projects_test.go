// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ProjectListOptions{ListOptions: ListOptions{Page: 2}}
	ctx := context.Background()
	projects, _, err := client.Repositories.ListProjects(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListProjects returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !cmp.Equal(projects, want) {
		t.Errorf("Repositories.ListProjects returned %+v, want %+v", projects, want)
	}

	const methodName = "ListProjects"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListProjects(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListProjects(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ProjectOptions{Name: String("Project Name"), Body: String("Project body.")}

	mux.HandleFunc("/repos/o/r/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &ProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Repositories.CreateProject(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Repositories.CreateProject returned %+v, want %+v", project, want)
	}

	const methodName = "CreateProject"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateProject(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateProject(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
