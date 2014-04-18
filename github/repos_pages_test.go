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

func TestRepositoriesService_GetPages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","cname":"c","custom_404":false}`)
	})

	page, _, err := client.Repositories.GetPages("o", "r")
	if err != nil {
		t.Errorf("Repositories.GetPages returned error: %v", err)
	}

	c := false
	want := &Pages{URL: String("u"), Status: String("s"), CNAME: String("c"), Custom404: &c}
	if !reflect.DeepEqual(page, want) {
		t.Errorf("Repositories.GetPages returned %+v, want %+v", page, want)
	}
}

func TestRepositoriesService_ListPagesBuilds(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"url":"u","status":"s","commit":"c"}]`)
	})

	pages, _, err := client.Repositories.ListPagesBuilds("o", "r")
	if err != nil {
		t.Errorf("Repositories.ListPagesBuilds returned error: %v", err)
	}

	want := []PagesBuild{{URL: String("u"), Status: String("s"), Commit: String("c")}}
	if !reflect.DeepEqual(pages, want) {
		t.Errorf("Repositories.ListPagesBuilds returned %+v, want %+v", pages, want)
	}
}

func TestRepositoriesService_ListLatestPagesBuilds(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pages/builds/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u","status":"s","commit":"c"}`)
	})

	build, _, err := client.Repositories.ListLatestPagesBuilds("o", "r")
	if err != nil {
		t.Errorf("Repositories.ListLatestPagesBuilds returned error: %v", err)
	}

	want := &PagesBuild{URL: String("u"), Status: String("s"), Commit: String("c")}
	if !reflect.DeepEqual(build, want) {
		t.Errorf("Repositories.ListLatestPagesBuilds returned %+v, want %+v", build, want)
	}
}
