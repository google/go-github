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

func TestGitService_GetBlob(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/blobs/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{"sha": "t"}`)
	})

	blob, _, err := client.Git.GetBlob("o", "r", "s")

	if err != nil {
		t.Errorf("Git.GetTag returned error: %v", err)
	}

	want := &Blob{SHA: "t"}
	if !reflect.DeepEqual(blob, want) {
		t.Errorf("Git.GetTag returned %+v, want %+v", blob, want)
	}
}

func TestGitService_CreateBlob(t *testing.T) {
	setup()
	defer teardown()

	input := &Blob{Content: "t", Encoding: "s"}

	mux.HandleFunc("/repos/o/r/git/blobs", func(w http.ResponseWriter, r *http.Request) {
		v := new(Blob)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"sha": "t"}`)
	})

	blob, _, err := client.Git.CreateBlob("o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateBlob returned error: %v", err)
	}

	want := &Blob{SHA: "t"}
	if !reflect.DeepEqual(blob, want) {
		t.Errorf("Git.GetBlob returned %+v, want %+v", blob, want)
	}
}