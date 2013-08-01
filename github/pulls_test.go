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

func TestPullRequestsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state": "closed",
			"head":  "h",
			"base":  "b",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &PullRequestListOptions{"closed", "h", "b"}
	pulls, _, err := client.PullRequests.List("o", "r", opt)

	if err != nil {
		t.Errorf("PullRequests.List returned error: %v", err)
	}

	want := []PullRequest{{Number: Int(1)}}
	if !reflect.DeepEqual(pulls, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", pulls, want)
	}
}

func TestPullRequestsService_List_invalidOwner(t *testing.T) {
	_, _, err := client.PullRequests.List("%", "r", nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":1}`)
	})

	pull, _, err := client.PullRequests.Get("o", "r", 1)

	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{Number: Int(1)}
	if !reflect.DeepEqual(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Get_invalidOwner(t *testing.T) {
	_, _, err := client.PullRequests.Get("%", "r", 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &PullRequest{Title: String("t")}

	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	pull, _, err := client.PullRequests.Create("o", "r", input)
	if err != nil {
		t.Errorf("PullRequests.Create returned error: %v", err)
	}

	want := &PullRequest{Number: Int(1)}
	if !reflect.DeepEqual(pull, want) {
		t.Errorf("PullRequests.Create returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Create_invalidOwner(t *testing.T) {
	_, _, err := client.PullRequests.Create("%", "r", nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &PullRequest{Title: String("t")}

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	pull, _, err := client.PullRequests.Edit("o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.Edit returned error: %v", err)
	}

	want := &PullRequest{Number: Int(1)}
	if !reflect.DeepEqual(pull, want) {
		t.Errorf("PullRequests.Edit returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Edit_invalidOwner(t *testing.T) {
	_, _, err := client.PullRequests.Edit("%", "r", 1, nil)
	testURLParseError(t, err)
}
