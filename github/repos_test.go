// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRepositoriesService_List_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method = %v, want %v", r.Method, "GET")
		}
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	repos, err := client.Repositories.List("", nil)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}, Repository{ID: 2}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		var v string
		if r.Method != "GET" {
			t.Errorf("Request method = %v, want %v", r.Method, "GET")
		}
		if v = r.FormValue("type"); v != "owner" {
			t.Errorf("Request type parameter = %v, want %v", v, "owner")
		}
		if v = r.FormValue("sort"); v != "created" {
			t.Errorf("Request sort parameter = %v, want %v", v, "created")
		}
		if v = r.FormValue("direction"); v != "asc" {
			t.Errorf("Request direction parameter = %v, want %v", v, "created")
		}
		if v = r.FormValue("page"); v != "2" {
			t.Errorf("Request page parameter = %v, want %v", v, "2")
		}

		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListOptions{"owner", "created", "asc", 2}
	repos, err := client.Repositories.List("u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_invalidUser(t *testing.T) {
	_, err := client.Repositories.List("%", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestRepositoriesService_ListByOrg(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method = %v, want %v", r.Method, "GET")
		}
		v := r.FormValue("type")
		if v != "forks" {
			t.Errorf("Request type parameter = %v, want %v", v, "forks")
		}
		v = r.FormValue("page")
		if v != "2" {
			t.Errorf("Request page parameter = %v, want %v", v, "2")
		}
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListByOrgOptions{"forks", 2}
	repos, err := client.Repositories.ListByOrg("o", opt)
	if err != nil {
		t.Errorf("Repositories.ListByOrg returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.ListByOrg returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_ListByOrg_invalidOrg(t *testing.T) {
	_, err := client.Repositories.ListByOrg("%", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestRepositoriesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method = %v, want %v", r.Method, "GET")
		}
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"}}`)
	})

	repo, err := client.Repositories.Get("o", "r")
	if err != nil {
		t.Errorf("Repositories.Get returned error: %v", err)
	}

	want := &Repository{ID: 1, Name: "n", Description: "d", Owner: &User{Login: "l"}}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Get returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Get_invalidOwner(t *testing.T) {
	_, err := client.Repositories.Get("%", "r")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}
