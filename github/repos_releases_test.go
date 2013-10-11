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

func TestRepositoriesService_ListReleases(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/u/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mimeReleasePreview)
		fmt.Fprint(w, `[{"id":1}]`)
	})

	releases, _, err := client.Repositories.ListReleases("o", "u")
	if err != nil {
		t.Errorf("Repositories.ListReleases returned error: %v", err)
	}
	want := []RepositoryRelease{{ID: Int(1)}}
	if !reflect.DeepEqual(releases, want) {
		t.Errorf("Repositories.ListReleases returned %+v, want %+v", releases, want)
	}
}

func TestRepositoriesService_GetRelease(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/u/releases/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mimeReleasePreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	release, resp, err := client.Repositories.GetRelease("o", "u", 1)
	if err != nil {
		t.Errorf("Repositories.GetRelease returned error: %v\n%v", err, resp.Body)
	}

	want := &RepositoryRelease{ID: Int(1)}
	if !reflect.DeepEqual(release, want) {
		t.Errorf("Repositories.GetRelease returned %+v, want %+v", release, want)
	}
}

func TestRepositoriesService_CreateRelease(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryRelease{Name: String("v1.0")}

	mux.HandleFunc("/repos/o/u/releases", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryRelease)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mimeReleasePreview)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	release, _, err := client.Repositories.CreateRelease("o", "u", input)
	if err != nil {
		t.Errorf("Repositories.CreateRelease returned error: %v", err)
	}

	want := &RepositoryRelease{ID: Int(1)}
	if !reflect.DeepEqual(release, want) {
		t.Errorf("Repositories.CreateRelease returned %+v, want %+v", release, want)
	}
}

func TestRepositoriesService_EditRelease(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryRelease{Name: String("n")}

	mux.HandleFunc("/repos/o/u/releases/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryRelease)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mimeReleasePreview)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	release, _, err := client.Repositories.EditRelease("o", "u", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditRelease returned error: %v", err)
	}
	want := &RepositoryRelease{ID: Int(1)}
	if !reflect.DeepEqual(release, want) {
		t.Errorf("Repositories.EditRelease returned = %+v, want %+v", release, want)
	}
}

func TestRepositoriesService_DeleteRelease(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/u/releases/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mimeReleasePreview)
	})

	_, err := client.Repositories.DeleteRelease("o", "u", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteRelease returned error: %v", err)
	}
}
