// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_ListAssignees(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	assignees, _, err := client.Issues.ListAssignees("o", "r")
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []User{User{ID: 1}}
	if !reflect.DeepEqual(assignees, want) {
		t.Errorf("Issues.ListAssignees returned %+v, want %+v", assignees, want)
	}
}

func TestIssuesService_ListAssignees_invalidOwner(t *testing.T) {
	_, _, err := client.Issues.ListAssignees("%", "r")
	testURLParseError(t, err)
}

func TestIssuesService_CheckAssignee_true(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	assignee, _, err := client.Issues.CheckAssignee("o", "r", "u")
	if err != nil {
		t.Errorf("Issues.CheckAssignee returned error: %v", err)
	}
	if want := true; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_false(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	assignee, _, err := client.Issues.CheckAssignee("o", "r", "u")
	if err != nil {
		t.Errorf("Issues.CheckAssignee returned error: %v", err)
	}
	if want := false; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	assignee, _, err := client.Issues.CheckAssignee("o", "r", "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_invalidOwner(t *testing.T) {
	_, _, err := client.Issues.CheckAssignee("%", "r", "u")
	testURLParseError(t, err)
}
