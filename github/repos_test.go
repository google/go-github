// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"encoding/json"
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
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
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
		var v, want string
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if v, want = r.FormValue("type"), "owner"; v != want {
			t.Errorf("Request type parameter = %v, want %v", v, want)
		}
		if v, want = r.FormValue("sort"), "created"; v != want {
			t.Errorf("Request sort parameter = %v, want %v", v, want)
		}
		if v, want = r.FormValue("direction"), "asc"; v != want {
			t.Errorf("Request direction parameter = %v, want %v", v, want)
		}
		if v, want = r.FormValue("page"), "2"; v != want {
			t.Errorf("Request page parameter = %v, want %v", v, want)
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
		var v, want string
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if v, want = r.FormValue("type"), "forks"; v != want {
			t.Errorf("Request type parameter = %v, want %v", v, want)
		}
		if v, want = r.FormValue("page"), "2"; v != want {
			t.Errorf("Request page parameter = %v, want %v", v, want)
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

func TestRepositoriesService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories", func(w http.ResponseWriter, r *http.Request) {
		var v, want string
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if v, want = r.FormValue("since"), "1"; v != want {
			t.Errorf("Request since parameter = %v, want %v", v, want)
		}
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListAllOptions{1}
	repos, err := client.Repositories.ListAll(opt)
	if err != nil {
		t.Errorf("Repositories.ListAll returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.ListAll returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_Create_user(t *testing.T) {
	setup()
	defer teardown()

	input := &Repository{Name: "n"}

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	repo, err := client.Repositories.Create("", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: 1}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Create_org(t *testing.T) {
	setup()
	defer teardown()

	input := &Repository{Name: "n"}

	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	repo, err := client.Repositories.Create("o", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: 1}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Create_invalidOrg(t *testing.T) {
	_, err := client.Repositories.Create("%", nil)
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
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
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

func TestRepositoriesService_ListForks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		var v, want string
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if v, want = r.FormValue("sort"), "newest"; v != want {
			t.Errorf("Request type parameter = %v, want %v", v, want)
		}
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &RepositoryListForksOptions{Sort: "newest"}
	repos, err := client.Repositories.ListForks("o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListForks returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}, Repository{ID: 2}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.ListForks returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_ListForks_invalidOwner(t *testing.T) {
	_, err := client.Repositories.ListForks("%", "r", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestRepositoriesService_CreateFork(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		var v, want string
		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		if v, want = r.FormValue("organization"), "o"; v != want {
			t.Errorf("Request type parameter = %v, want %v", v, want)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	opt := &RepositoryCreateForkOptions{Organization: "o"}
	repo, err := client.Repositories.CreateFork("o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.CreateFork returned error: %v", err)
	}

	want := &Repository{ID: 1}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.CreateFork returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_CreateFork_invalidOwner(t *testing.T) {
	_, err := client.Repositories.CreateFork("%", "r", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestRepositoriesService_ListStatuses(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
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
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestRepositoriesService_CreateStatus(t *testing.T) {
	setup()
	defer teardown()

	input := &RepoStatus{State: "s", TargetURL: "t", Description: "d"}

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepoStatus)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
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
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}
