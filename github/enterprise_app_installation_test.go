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

func TestEnterpriseService_ListInstallableEnterpriseOrgs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/installable_organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1, "login":"org1"}]`)
	})

	ctx := t.Context()
	opts := &ListOptions{Page: 1, PerPage: 10}
	got, _, err := client.Enterprise.ListInstallableEnterpriseOrgs(ctx, "e", opts)
	if err != nil {
		t.Fatalf("Enterprise.ListInstallableEnterpriseOrgs returned error: %v", err)
	}

	want := []*InstallableOrganization{
		{ID: Ptr(int64(1)), Login: Ptr("org1")},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.ListInstallableEnterpriseOrgs = %+v, want %+v", got, want)
	}

	const methodName = "ListInstallableEnterpriseOrgs"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.ListInstallableEnterpriseOrgs(ctx, "\n", opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListInstallableEnterpriseOrgs(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListEnterpriseOrgAccessibleRepositories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/installable_organizations/org1/accessible_repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":10, "name":"repo1", "full_name":"org1/repo1"}]`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	repos, _, err := client.Enterprise.ListEnterpriseOrgAccessibleRepositories(ctx, "e", "org1", opts)
	if err != nil {
		t.Errorf("Enterprise.ListEnterpriseOrgAccessibleRepositories returned error: %v", err)
	}

	want := []*AccessibleRepository{
		{ID: Ptr(int64(10)), Name: Ptr("repo1"), FullName: Ptr("org1/repo1")},
	}

	if !cmp.Equal(repos, want) {
		t.Errorf("Enterprise.ListEnterpriseOrgAccessibleRepositories returned %+v, want %+v", repos, want)
	}

	const methodName = "ListEnterpriseOrgAccessibleRepositories"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListEnterpriseOrgAccessibleRepositories(ctx, "\n", "org1", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListEnterpriseOrgAccessibleRepositories(ctx, "e", "org1", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListEnterpriseAppInstallationForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/organizations/org1/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `[{"id":99}]`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	installations, _, err := client.Enterprise.ListEnterpriseAppInstallationForOrg(ctx, "e", "org1", opts)
	if err != nil {
		t.Errorf("ListEnterpriseAppInstallationForOrg returned error: %v", err)
	}

	want := &[]Installation{
		{ID: Ptr(int64(99))},
	}

	if !cmp.Equal(installations, want) {
		t.Errorf("ListEnterpriseAppInstallationForOrg returned %+v, want %+v", installations, want)
	}

	const methodName = "ListEnterpriseAppInstallationForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListEnterpriseAppInstallationForOrg(ctx, "\n", "org1", &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListEnterpriseAppInstallationForOrg(ctx, "e", "org1", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_InstallEnterpriseOrgApp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/organizations/org1/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"client_id":"cid","repository_selection":"selected","repository":["r1","r2"]}`+"\n")
		fmt.Fprint(w, `{"id":555}`)
	})

	req := AppInstallationRequest{
		ClientID:            "cid",
		RepositorySelection: "selected",
		Repository:          []string{"r1", "r2"},
	}

	ctx := t.Context()
	installation, _, err := client.Enterprise.InstallEnterpriseOrgApp(ctx, "e", "org1", req)
	if err != nil {
		t.Errorf("InstallEnterpriseOrgApp returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(555))}

	if !cmp.Equal(installation, want) {
		t.Errorf("InstallEnterpriseOrgApp returned %+v, want %+v", installation, want)
	}

	const methodName = "InstallEnterpriseOrgApp"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.InstallEnterpriseOrgApp(ctx, "e", "org1", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UninstallEnterpriseOrgApp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/apps/organizations/org1/installations/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	resp, err := client.Enterprise.UninstallEnterpriseOrgApp(ctx, "e", "org1", 123)
	if err != nil {
		t.Errorf("UninstallEnterpriseOrgApp returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("UninstallEnterpriseOrgApp returned status %v, want %v", resp.StatusCode, http.StatusNoContent)
	}

	const methodName = "UninstallEnterpriseOrgApp"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.UninstallEnterpriseOrgApp(ctx, "e", "org1", 123)
	})
}
