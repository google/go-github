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

func TestActivityService_ListStarred_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/starred", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	repos, _, err := client.Activity.ListStarred("", nil)
	if err != nil {
		t.Errorf("Activity.ListStarred returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Activity.ListStarred returned %+v, want %+v", repos, want)
	}
}

func TestActivityService_ListStarred_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/starred", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"sort":      "created",
			"direction": "asc",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &ActivityListStarredOptions{"created", "asc", ListOptions{Page: 2}}
	repos, _, err := client.Activity.ListStarred("u", opt)
	if err != nil {
		t.Errorf("Activity.ListStarred returned error: %v", err)
	}

	want := []Repository{{ID: Int(2)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Activity.ListStarred returned %+v, want %+v", repos, want)
	}
}

func TestActivityService_ListStarred_invalidUser(t *testing.T) {
	_, _, err := client.Activity.ListStarred("%", nil)
	testURLParseError(t, err)
}
