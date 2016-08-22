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

func TestRepositoriesService_ListPaths(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("repos/1/2/traffic/popular/paths", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `[{"path":"/google/go-github/pull/393/files",
			 "title":"add unit test for #342",
			 "Counts": 6,
			 "uniques": 5}]`)
	})
	paths, _, err := client.Repositories.ListPaths("1", "2")
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := []*Path{{
		Path:    String("/google/go-github/pull/393/files"),
		Title:   String("add unit test for #342"),
		Count:   Int(6),
		Uniques: Int(5),
	}}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("Repositories.ListPaths returned %+v, want %+v", paths, want)
	}

}
