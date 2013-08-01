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

func TestGitService_GetCommit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"sha":"s","message":"m","author":{"name":"n"}}`)
	})

	commit, err := client.Git.GetCommit("o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetCommit returned error: %v", err)
	}

	want := &Commit{SHA: "s", Message: "m", Author: &CommitAuthor{Name: "n"}}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.GetCommit returned %+v, want %+v", commit, want)
	}
}

func TestGitService_CreateCommit(t *testing.T) {
	setup()
	defer teardown()

	input := &Commit{Message: "m", Tree: &Tree{SHA: "t"}}

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(Commit)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"sha":"s"}`)
	})

	commit, err := client.Git.CreateCommit("o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: "s"}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}
}
