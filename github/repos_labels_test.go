// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CreateTestRepos() (goodRepo *Repository, badRepo *Repository) {
	user := "o"
	invalidUser := "\n"
	repoName := "r"
	invalidName := "%"
	return &Repository{
			Owner: &User{
				Login: &user,
			},
			Name: &repoName,
		}, &Repository{
			Owner: &User{
				Login: &invalidUser,
			},
			Name: &invalidName,
		}
}

func TestRepoService_ListLabels(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name": "a"},{"name": "b"}]`)
	})
	goodRepo, badRepo := CreateTestRepos()
	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	labels, _, err := client.Repositories.ListLabels(ctx, goodRepo, opt)
	if err != nil {
		t.Errorf("Repo.ListLabels returned error: %v", err)
	}

	want := []*Label{{Name: String("a")}, {Name: String("b")}}
	if !cmp.Equal(labels, want) {
		t.Errorf("Repo.ListLabels returned %+v, want %+v", labels, want)
	}

	const methodName = "ListLabels"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListLabels(ctx, badRepo, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListLabels(ctx, goodRepo, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepoService_ListLabels_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, badRepo := CreateTestRepos()
	_, _, err := client.Repositories.ListLabels(ctx, badRepo, nil)
	testURLParseError(t, err)
}

func TestRepoService_GetLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "name": "n", "color": "c", "description": "d"}`)
	})
	goodRepo, badRepo := CreateTestRepos()

	ctx := context.Background()
	label, _, err := client.Repositories.GetLabel(ctx, goodRepo, "n")
	if err != nil {
		t.Errorf("Repo.GetLabel returned error: %v", err)
	}

	want := &Label{URL: String("u"), Name: String("n"), Color: String("c"), Description: String("d")}
	if !cmp.Equal(label, want) {
		t.Errorf("Repo.GetLabel returned %+v, want %+v", label, want)
	}

	const methodName = "GetLabel"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetLabel(ctx, badRepo, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetLabel(ctx, goodRepo, "n")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepoService_GetLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, badRepo := CreateTestRepos()
	_, _, err := client.Repositories.GetLabel(ctx, badRepo, "%")
	testURLParseError(t, err)
}

func TestRepoService_CreateLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Label{Name: String("n")}

	mux.HandleFunc("/repos/o/r/labels", func(w http.ResponseWriter, r *http.Request) {
		v := new(Label)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"url":"u"}`)
	})

	ctx := context.Background()
	goodRepo, badRepo := CreateTestRepos()
	label, _, err := client.Repositories.CreateLabel(ctx, goodRepo, input)
	if err != nil {
		t.Errorf("Repo.CreateLabel returned error: %v", err)
	}

	want := &Label{URL: String("u")}
	if !cmp.Equal(label, want) {
		t.Errorf("Repo.CreateLabel returned %+v, want %+v", label, want)
	}

	const methodName = "CreateLabel"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateLabel(ctx, badRepo, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateLabel(ctx, goodRepo, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepoService_CreateLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, badRepo := CreateTestRepos()
	_, _, err := client.Repositories.CreateLabel(ctx, badRepo, nil)
	testURLParseError(t, err)
}

func TestRepoService_EditLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Label{Name: String("z")}

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		v := new(Label)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"url":"u"}`)
	})

	ctx := context.Background()
	goodRepo, badRepo := CreateTestRepos()
	label, _, err := client.Repositories.EditLabel(ctx, goodRepo, "n", input)
	if err != nil {
		t.Errorf("Repo.EditLabel returned error: %v", err)
	}

	want := &Label{URL: String("u")}
	if !cmp.Equal(label, want) {
		t.Errorf("Repo.EditLabel returned %+v, want %+v", label, want)
	}

	const methodName = "EditLabel"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditLabel(ctx, badRepo, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditLabel(ctx, goodRepo, "n", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepoService_EditLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, badRepo := CreateTestRepos()
	_, _, err := client.Repositories.EditLabel(ctx, badRepo, "%", nil)
	testURLParseError(t, err)
}

func TestRepoService_DeleteLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	goodRepo, badRepo := CreateTestRepos()
	_, err := client.Repositories.DeleteLabel(ctx, goodRepo, "n")
	if err != nil {
		t.Errorf("Repo.DeleteLabel returned error: %v", err)
	}

	const methodName = "DeleteLabel"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteLabel(ctx, badRepo, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteLabel(ctx, goodRepo, "n")
	})
}

func TestRepoService_DeleteLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, badRepo := CreateTestRepos()
	_, err := client.Repositories.DeleteLabel(ctx, badRepo, "%")
	testURLParseError(t, err)
}
