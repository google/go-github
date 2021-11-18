// Copyright 2019 The go-github AUTHORS. All rights reserved.
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

func TestUsersService_ListProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"state": "open", "page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ProjectListOptions{State: "open", ListOptions: ListOptions{Page: 2}}
	ctx := context.Background()
	projects, _, err := client.Users.ListProjects(ctx, "u", opt)
	if err != nil {
		t.Errorf("Users.ListProjects returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !cmp.Equal(projects, want) {
		t.Errorf("Users.ListProjects returned %+v, want %+v", projects, want)
	}

	const methodName = "ListProjects"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Users.ListProjects(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListProjects(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_CreateProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &CreateUserProjectOptions{Name: "Project Name", Body: String("Project body.")}

	mux.HandleFunc("/user/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &CreateUserProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Users.CreateProject(ctx, input)
	if err != nil {
		t.Errorf("Users.CreateProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Users.CreateProject returned %+v, want %+v", project, want)
	}

	const methodName = "CreateProject"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.CreateProject(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCreateUserProjectOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &CreateUserProjectOptions{}, `{}`)

	c := CreateUserProjectOptions{
		Name: "SomeProject",
		Body: String("SomeProjectBody"),
	}

	want := `{
			"name": "SomeProject",
			"body": "SomeProjectBody"
		}`

	testJSONMarshal(t, c, want)
}
