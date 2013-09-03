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

func TestRepositoriesService_List_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	repos, _, err := client.Repositories.List("", nil)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}, {ID: Int(2)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"type":      "owner",
			"sort":      "created",
			"direction": "asc",
			"page":      "2",
		})

		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListOptions{"owner", "created", "asc", ListOptions{Page: 2}}
	repos, _, err := client.Repositories.List("u", opt)
	if err != nil {
		t.Errorf("Repositories.List returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.List returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_List_invalidUser(t *testing.T) {
	_, _, err := client.Repositories.List("%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListByOrg(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"type": "forks",
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListByOrgOptions{"forks", ListOptions{Page: 2}}
	repos, _, err := client.Repositories.ListByOrg("o", opt)
	if err != nil {
		t.Errorf("Repositories.ListByOrg returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.ListByOrg returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_ListByOrg_invalidOrg(t *testing.T) {
	_, _, err := client.Repositories.ListByOrg("%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since":    "1",
			"page":     "2",
			"per_page": "3",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &RepositoryListAllOptions{1, ListOptions{2, 3}}
	repos, _, err := client.Repositories.ListAll(opt)
	if err != nil {
		t.Errorf("Repositories.ListAll returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repositories.ListAll returned %+v, want %+v", repos, want)
	}
}

func TestRepositoriesService_Create_user(t *testing.T) {
	setup()
	defer teardown()

	input := &Repository{Name: String("n")}

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Repositories.Create("", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: Int(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Create_org(t *testing.T) {
	setup()
	defer teardown()

	input := &Repository{Name: String("n")}

	mux.HandleFunc("/orgs/o/repos", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Repositories.Create("o", input)
	if err != nil {
		t.Errorf("Repositories.Create returned error: %v", err)
	}

	want := &Repository{ID: Int(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Create returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Create_invalidOrg(t *testing.T) {
	_, _, err := client.Repositories.Create("%", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"}}`)
	})

	repo, _, err := client.Repositories.Get("o", "r")
	if err != nil {
		t.Errorf("Repositories.Get returned error: %v", err)
	}

	want := &Repository{ID: Int(1), Name: String("n"), Description: String("d"), Owner: &User{Login: String("l")}}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Get returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Edit(t *testing.T) {
	setup()
	defer teardown()

	i := true
	input := &Repository{HasIssues: &i}

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(Repository)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Repositories.Edit("o", "r", input)
	if err != nil {
		t.Errorf("Repositories.Edit returned error: %v", err)
	}

	want := &Repository{ID: Int(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Edit returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Get_invalidOwner(t *testing.T) {
	_, _, err := client.Repositories.Get("%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_Edit_invalidOwner(t *testing.T) {
	_, _, err := client.Repositories.Edit("%", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListLanguages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/languages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"go":1}`)
	})

	languages, _, err := client.Repositories.ListLanguages("o", "r")
	if err != nil {
		t.Errorf("Repositories.ListLanguages returned error: %v", err)
	}

	want := map[string]int{"go": 1}
	if !reflect.DeepEqual(languages, want) {
		t.Errorf("Repositories.ListLanguages returned %+v, want %+v", languages, want)
	}
}

func TestRepositoriesService_ListLanguages_invalidOwner(t *testing.T) {
	_, _, err := client.Repositories.ListLanguages("%", "%")
	testURLParseError(t, err)
}
