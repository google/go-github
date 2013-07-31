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
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"type":      "owner",
			"sort":      "created",
			"direction": "asc",
			"page":      "2",
		})

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
	testURLParseError(t, err)
}

func TestRepositoriesService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1"})
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

		testMethod(t, r, "POST")
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

		testMethod(t, r, "POST")
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
	testURLParseError(t, err)
}

func TestRepositoriesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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

	repo, err := client.Repositories.Edit("o", "r", input)
	if err != nil {
		t.Errorf("Repositories.Edit returned error: %v", err)
	}

	want := &Repository{ID: 1}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Repositories.Edit returned %+v, want %+v", repo, want)
	}
}

func TestRepositoriesService_Get_invalidOwner(t *testing.T) {
	_, err := client.Repositories.Get("%", "r")
	testURLParseError(t, err)
}

func TestRepositoriesService_Edit_invalidOwner(t *testing.T) {
	_, err := client.Repositories.Edit("%", "r", nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_ListForks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"sort": "newest"})
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
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateFork(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"organization": "o"})
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
	testURLParseError(t, err)
}

func TestRepositoriesService_ListStatuses(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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
	testURLParseError(t, err)
}

func TestRepositoriesService_CreateStatus(t *testing.T) {
	setup()
	defer teardown()

	input := &RepoStatus{State: "s", TargetURL: "t", Description: "d"}

	mux.HandleFunc("/repos/o/r/statuses/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepoStatus)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
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
	testURLParseError(t, err)
}

func TestRepositoriesService_ListLanguages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/u/r/languages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"go":1}`)
	})

	languages, err := client.Repositories.ListLanguages("u", "r")
	if err != nil {
		t.Errorf("Repositories.ListLanguages returned error: %v", err)
	}

	want := map[string]int{"go": 1}
	if !reflect.DeepEqual(languages, want) {
		t.Errorf("Repositories.ListLanguages returned %+v, want %+v", languages, want)
	}
}

func TestRepositoriesService_CreateHook(t *testing.T) {
	setup()
	defer teardown()

	input := &Hook{Name: "t"}

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Repositories.CreateHook("o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateHook returned error: %v", err)
	}

	want := &Hook{ID: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Repositories.CreateHook returned %+v, want %+v", hook, want)
	}
}

func TestRepositoriesService_ListHooks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{2}

	hooks, err := client.Repositories.ListHooks("o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListHooks returned error: %v", err)
	}

	want := []Hook{Hook{ID: 1}, Hook{ID: 2}}
	if !reflect.DeepEqual(hooks, want) {
		t.Errorf("Repositories.ListHooks returned %+v, want %+v", hooks, want)
	}
}

func TestRepositoriesService_GetHook(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Repositories.GetHook("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetHook returned error: %v", err)
	}

	want := &Hook{ID: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Repositories.GetHook returned %+v, want %+v", hook, want)
	}
}

func TestRepositoriesService_EditHook(t *testing.T) {
	setup()
	defer teardown()

	input := &Hook{Name: "t"}

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Repositories.EditHook("o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditHook returned error: %v", err)
	}

	want := &Hook{ID: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Repositories.EditHook returned %+v, want %+v", hook, want)
	}
}

func TestRepositoriesService_DeleteHook(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Repositories.DeleteHook("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteHook returned error: %v", err)
	}
}

func TestRepositoriesService_TestHook(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	err := client.Repositories.TestHook("o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.TestHook returned error: %v", err)
	}
}
