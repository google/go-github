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
	"reflect"
	"testing"
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
	projects, _, err := client.Users.ListProjects(context.Background(), "u", opt)
	if err != nil {
		t.Errorf("Users.ListProjects returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Users.ListProjects returned %+v, want %+v", projects, want)
	}
}

func TestUsersService_CreateProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &CreateUserProjectOptions{Name: "Project Name", Body: String("Project body.")}

	mux.HandleFunc("/users/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)

		v := &CreateUserProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Users.CreateProject(context.Background(), input)
	if err != nil {
		t.Errorf("Users.CreateProject returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Users.CreateProject returned %+v, want %+v", project, want)
	}
}
