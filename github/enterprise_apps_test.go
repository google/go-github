// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseService_ListRepositoriesForOrgAppInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := t.Context()
	repos, _, err := client.Enterprise.ListRepositoriesForOrgAppInstallation(ctx, "e", "o", 1, &ListOptions{Page: 1})
	if err != nil {
		t.Errorf("Enterprise.ListRepositoriesForOrgAppInstallation returned error: %v", err)
	}

	want := []*AccessibleRepository{{ID: 1}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("Enterprise.ListRepositoriesForOrgAppInstallation returned diff (-want +got):\n%v", diff)
	}

	const methodName = "ListRepositoriesForOrgAppInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRepositoriesForOrgAppInstallation(ctx, "\n", "\n", -1, &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.Enterprise.ListRepositoriesForOrgAppInstallation(ctx, "e", "o", 1, &ListOptions{})
		return resp, err
	})
}

func TestEnterpriseService_UpdateAppInstallationRepositories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := UpdateAppInstallationRepositoriesOptions{
		RepositorySelection:   Ptr("selected"),
		SelectedRepositoryIDs: []int64{1, 2},
	}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"repository_selection":"selected","selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `{"id":1, "repository_selection":"selected"}`)
	})

	ctx := t.Context()
	inst, _, err := client.Enterprise.UpdateAppInstallationRepositories(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("Enterprise.UpdateAppInstallationRepositories returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), RepositorySelection: Ptr("selected")}
	if diff := cmp.Diff(inst, want); diff != "" {
		t.Errorf("Enterprise.UpdateAppInstallationRepositories returned diff (-want +got):\n%v", diff)
	}

	const methodName = "UpdateAppInstallationRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.UpdateAppInstallationRepositories(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.Enterprise.UpdateAppInstallationRepositories(ctx, "e", "o", 1, input)
		return resp, err
	})
}

func TestEnterpriseService_AddRepositoriesToAppInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := AppInstallationRepositoriesOptions{SelectedRepositoryIDs: []int64{1, 2}}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := t.Context()
	repos, _, err := client.Enterprise.AddRepositoriesToAppInstallation(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("Enterprise.AddRepositoriesToAppInstallation returned error: %v", err)
	}

	want := []*AccessibleRepository{{ID: 1}, {ID: 2}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("Enterprise.AddRepositoriesToAppInstallation returned diff (-want +got):\n%v", diff)
	}

	const methodName = "AddRepositoriesToAppInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.AddRepositoriesToAppInstallation(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.Enterprise.AddRepositoriesToAppInstallation(ctx, "e", "o", 1, input)
		return resp, err
	})
}

func TestEnterpriseService_RemoveRepositoriesFromAppInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := AppInstallationRepositoriesOptions{SelectedRepositoryIDs: []int64{1, 2}}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := t.Context()
	repos, _, err := client.Enterprise.RemoveRepositoriesFromAppInstallation(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("Enterprise.RemoveRepositoriesFromAppInstallation returned error: %v", err)
	}

	want := []*AccessibleRepository{{ID: 1}, {ID: 2}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("Enterprise.RemoveRepositoriesFromAppInstallation returned diff (-want +got):\n%v", diff)
	}

	const methodName = "RemoveRepositoriesFromAppInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.RemoveRepositoriesFromAppInstallation(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.Enterprise.RemoveRepositoriesFromAppInstallation(ctx, "e", "o", 1, input)
		return resp, err
	})
}
