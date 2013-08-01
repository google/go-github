// Copyright 2013 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_ListStatuses(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	statuses, err := client.Repositories.ListStatuses("o", "r", "r")
	if err != nil {
		t.Errorf("Repositories.ListStatuses returned error: %v", err)
	}

	want := []RepoStatus{RepoStatus{ID: 1}}
	if !reflect.DeepEqual(statuses, want) {
		t.Errorf("Repositories.ListStatuses returned %+v, want %+v", statuses, want)
	}
}

func TestRepositoriesService_ListStatuses_invalidOwner(t *testing.T) {
	_, err := client.Repositories.ListStatuses("%", "r", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateStatus(t *testing.T) {
	setup()
	defer teardown()

	input := &RepoStatus{State: "s", TargetURL: "t", Description: "d"}

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepoStatus)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	status, err := client.Repositories.CreateStatus("o", "r", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateStatus returned error: %v", err)
	}

	want := &RepoStatus{ID: 1}
	if !reflect.DeepEqual(status, want) {
		t.Errorf("Repositories.CreateStatus returned %+v, want %+v", status, want)
	}
}

func TestRepositoriesService_CreateStatus_invalidOwner(t *testing.T) {
	_, err := client.Repositories.CreateStatus("%", "r", "r", nil)
	testURLParseError(t, err)
}
