// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseAppsService_ListRepositoriesForOrgInstallation(t *testing.T) {
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1"})
		fmt.Fprint(w, `{"total_count":1, "repositories":[{"id":1}]}`)
	})

	ctx := context.Background()
	repos, _, err := client.EnterpriseApps.ListRepositoriesForOrgInstallation(ctx, "e", "o", 1, &ListOptions{Page: 1})
	if err != nil {
		t.Errorf("EnterpriseApps.ListRepositoriesForOrgInstallation returned error: %v", err)
	}

	want := &ListRepositories{TotalCount: Ptr(1), Repositories: []*Repository{{ID: Ptr(int64(1))}}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("EnterpriseApps.ListRepositoriesForOrgInstallation returned diff (-want +got):\n%s", diff)
	}

	const methodName = "ListRepositoriesForOrgInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.EnterpriseApps.ListRepositoriesForOrgInstallation(ctx, "\n", "\n", -1, &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.EnterpriseApps.ListRepositoriesForOrgInstallation(ctx, "e", "o", 1, &ListOptions{})
		return resp, err
	})
}

func TestEnterpriseAppsService_ToggleInstallationRepositories(t *testing.T) {
	client, mux, _ := setup(t)

	input := &EnterpriseInstallationRepositoriesToggleOptions{
		RepositorySelection:   String("selected"),
		SelectedRepositoryIDs: []int64{1, 2},
	}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"repository_selection":"selected","selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `{"total_count":2, "repositories":[{"id":1},{"id":2}]}`)
	})

	ctx := context.Background()
	repos, _, err := client.EnterpriseApps.ToggleInstallationRepositories(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("EnterpriseApps.ToggleInstallationRepositories returned error: %v", err)
	}

	want := &ListRepositories{TotalCount: Ptr(2), Repositories: []*Repository{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("EnterpriseApps.ToggleInstallationRepositories returned diff (-want +got):\n%s", diff)
	}

	const methodName = "ToggleInstallationRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.EnterpriseApps.ToggleInstallationRepositories(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.EnterpriseApps.ToggleInstallationRepositories(ctx, "e", "o", 1, input)
		return resp, err
	})
}

func TestEnterpriseAppsService_AddRepositoriesToInstallation(t *testing.T) {
	client, mux, _ := setup(t)

	input := &EnterpriseInstallationRepositoriesOptions{SelectedRepositoryIDs: []int64{1, 2}}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `{"total_count":2, "repositories":[{"id":1},{"id":2}]}`)
	})

	ctx := context.Background()
	repos, _, err := client.EnterpriseApps.AddRepositoriesToInstallation(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("EnterpriseApps.AddRepositoriesToInstallation returned error: %v", err)
	}

	want := &ListRepositories{TotalCount: Ptr(2), Repositories: []*Repository{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("EnterpriseApps.AddRepositoriesToInstallation returned diff (-want +got):\n%s", diff)
	}

	const methodName = "AddRepositoriesToInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.EnterpriseApps.AddRepositoriesToInstallation(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.EnterpriseApps.AddRepositoriesToInstallation(ctx, "e", "o", 1, input)
		return resp, err
	})
}

func TestEnterpriseAppsService_RemoveRepositoriesFromInstallation(t *testing.T) {
	client, mux, _ := setup(t)

	input := &EnterpriseInstallationRepositoriesOptions{SelectedRepositoryIDs: []int64{1, 2}}

	mux.HandleFunc("/enterprises/e/apps/organizations/o/installations/1/repositories/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"selected_repository_ids":[1,2]}`+"\n")
		fmt.Fprint(w, `{"total_count":2, "repositories":[{"id":1},{"id":2}]}`)
	})

	ctx := context.Background()
	repos, _, err := client.EnterpriseApps.RemoveRepositoriesFromInstallation(ctx, "e", "o", 1, input)
	if err != nil {
		t.Errorf("EnterpriseApps.RemoveRepositoriesFromInstallation returned error: %v", err)
	}

	want := &ListRepositories{TotalCount: Ptr(2), Repositories: []*Repository{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}}
	if diff := cmp.Diff(repos, want); diff != "" {
		t.Errorf("EnterpriseApps.RemoveRepositoriesFromInstallation returned diff (-want +got):\n%s", diff)
	}

	const methodName = "RemoveRepositoriesFromInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.EnterpriseApps.RemoveRepositoriesFromInstallation(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.EnterpriseApps.RemoveRepositoriesFromInstallation(ctx, "e", "o", 1, input)
		return resp, err
	})
}
