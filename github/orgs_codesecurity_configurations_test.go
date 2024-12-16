// Copyright 2024 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrganizationsService_GetCodeSecurityConfigurations(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-security/configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		{
			"id":1,
			"name":"config1",
			"code_scanning_default_setup": "enabled"
		},
		{
			"id":2,
			"name":"config2",
			"private_vulnerability_reporting": "enabled"
		}]`)
	})

	configurations, _, err := client.Organizations.GetCodeSecurityConfigurations(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetOrganizationCodeSecurityConfigurations returned error: %v", err)
	}

	want := []*CodeSecurityConfiguration{
		{ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled")},
		{ID: Ptr(int64(2)), Name: Ptr("config2"), PrivateVulnerabilityReporting: Ptr("enabled")},
	}
	if !reflect.DeepEqual(configurations, want) {
		t.Errorf("Organizations.GetCodeSecurityConfigurations returned %+v, want %+v", configurations, want)
	}
	const methodName = "GetCodeSecurityConfigurations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetCodeSecurityConfigurations(ctx, "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetCodeSecurityConfigurations(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := context.Background()

	mux.HandleFunc("/orgs/o/code-security/configurations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id":1,
			"name":"config1",
			"code_scanning_default_setup": "enabled"
		}`)
	})

	configuration, _, err := client.Organizations.GetCodeSecurityConfiguration(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.GetCodeSecurityConfiguration returned error: %v", err)
	}

	want := &CodeSecurityConfiguration{ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled")}
	if !reflect.DeepEqual(configuration, want) {
		t.Errorf("Organizations.GetCodeSecurityConfiguration returned %+v, want %+v", configuration, want)
	}

	const methodName = "GetCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetCodeSecurityConfiguration(ctx, "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetCodeSecurityConfiguration(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := context.Background()

	input := &CodeSecurityConfiguration{
		Name:                     Ptr("config1"),
		CodeScanningDefaultSetup: Ptr("enabled"),
	}

	mux.HandleFunc("/orgs/o/code-security/configurations", func(w http.ResponseWriter, r *http.Request) {
		v := new(CodeSecurityConfiguration)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Errorf("Organizations.CreateCodeSecurityConfiguration request body decode failed: %v", err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Organizations.CreateCodeSecurityConfiguration request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
			"id":1,
			"name":"config1",
			"code_scanning_default_setup": "enabled"
		}`)
	})

	configuration, _, err := client.Organizations.CreateCodeSecurityConfiguration(ctx, "o", input)
	if err != nil {
		t.Errorf("Organizations.CreateCodeSecurityConfiguration returned error: %v", err)
	}

	want := &CodeSecurityConfiguration{ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled")}
	if !reflect.DeepEqual(configuration, want) {
		t.Errorf("Organizations.CreateCodeSecurityConfiguration returned %+v, want %+v", configuration, want)
	}

	const methodName = "CreateCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.CreateCodeSecurityConfiguration(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateCodeSecurityConfiguration(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetDefaultCodeSecurityConfigurations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := context.Background()

	mux.HandleFunc("/orgs/o/code-security/configurations/defaults", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		{
			"id":1,
			"name":"config1",
			"code_scanning_default_setup": "enabled"
		},
		{
			"id":2,
			"name":"config2",
			"private_vulnerability_reporting": "enabled"
		}]`)
	})

	configurations, _, err := client.Organizations.GetDefaultCodeSecurityConfigurations(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetDefaultCodeSecurityConfigurations returned error: %v", err)
	}

	want := []*CodeSecurityConfiguration{
		{ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled")},
		{ID: Ptr(int64(2)), Name: Ptr("config2"), PrivateVulnerabilityReporting: Ptr("enabled")},
	}
	if !reflect.DeepEqual(configurations, want) {
		t.Errorf("Organizations.GetDefaultCodeSecurityConfigurations returned %+v, want %+v", configurations, want)
	}

	const methodName = "GetDefaultCodeSecurityConfigurations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetDefaultCodeSecurityConfigurations(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetDefaultCodeSecurityConfigurations(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DetachCodeSecurityConfigurationsFromRepositories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := context.Background()

	mux.HandleFunc("/orgs/o/code-security/configurations/detach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Organizations.DetachCodeSecurityConfigurationsFromRepositories(ctx, "o", []int64{1})
	if err != nil {
		t.Errorf("Organizations.DetachCodeSecurityConfigurationsFromRepositories returned error: %v", err)
	}

	want := http.StatusNoContent
	if resp.StatusCode != want {
		t.Errorf("Organizations.DetachCodeSecurityConfigurationsFromRepositories returned status %d, want %d", resp.StatusCode, want)
	}

	const methodName = "DetachCodeSecurityConfigurationsFromRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DetachCodeSecurityConfigurationsFromRepositories(ctx, "\n", []int64{1})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.Organizations.DetachCodeSecurityConfigurationsFromRepositories(ctx, "o", []int64{1})
		return resp, err
	})
}

func TestOrganizationsService_UpdateCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	input := &CodeSecurityConfiguration{
		Name:                     Ptr("config1"),
		CodeScanningDefaultSetup: Ptr("enabled"),
	}

	mux.HandleFunc("/orgs/o/code-security/configurations/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(CodeSecurityConfiguration)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Errorf("Organizations.UpdateCodeSecurityConfiguration request body decode failed: %v", err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Organizations.UpdateCodeSecurityConfiguration request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
			"id":1,
			"name":"config1",
			"code_scanning_default_setup": "enabled"
		}`)
	})

	configuration, _, err := client.Organizations.UpdateCodeSecurityConfiguration(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.UpdateCodeSecurityConfiguration returned error: %v", err)
	}

	want := &CodeSecurityConfiguration{ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled")}
	if !reflect.DeepEqual(configuration, want) {
		t.Errorf("Organizations.UpdateCodeSecurityConfiguration returned %+v, want %+v", configuration, want)
	}

	const methodName = "UpdateCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.UpdateCodeSecurityConfiguration(ctx, "\n", -1, input)
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateCodeSecurityConfiguration(ctx, "o", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-security/configurations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Organizations.DeleteCodeSecurityConfiguration(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.DeleteCodeSecurityConfiguration returned error: %v", err)
	}

	want := http.StatusNoContent
	if resp.StatusCode != want {
		t.Errorf("Organizations.DeleteCodeSecurityConfiguration returned status %d, want %d", resp.StatusCode, want)
	}

	const methodName = "DeleteCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DeleteCodeSecurityConfiguration(ctx, "\n", -1)
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.Organizations.DeleteCodeSecurityConfiguration(ctx, "o", 1)
		return resp, err
	})
}

func TestOrganizationsService_AttachCodeSecurityConfigurationsToRepositories(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-security/configurations/1/attach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		type request struct {
			Scope                 string  `json:"scope"`
			SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitempty"`
		}
		v := new(request)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Errorf("Organizations.AttachCodeSecurityConfigurationsToRepositories request body decode failed: %v", err)
		}
		if v.Scope != "selected" {
			t.Errorf("Organizations.AttachCodeSecurityConfigurationsToRepositories request body scope = %s, want selected", v.Scope)
		}
		if !reflect.DeepEqual(v.SelectedRepositoryIDs, []int64{5, 20}) {
			t.Errorf("Organizations.AttachCodeSecurityConfigurationsToRepositories request body selected_repository_ids = %+v, want %+v", v.SelectedRepositoryIDs, []int64{5, 20})
		}
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.Organizations.AttachCodeSecurityConfigurationsToRepositories(ctx, "o", int64(1), "selected", []int64{5, 20})
	if err != nil {
		t.Errorf("Organizations.AttachCodeSecurityConfigurationsToRepositories returned error: %v", err)
	}

	want := http.StatusAccepted
	if resp.StatusCode != want {
		t.Errorf("Organizations.AttachCodeSecurityConfigurationsToRepositories returned status %d, want %d", resp.StatusCode, want)
	}

	const methodName = "AttachCodeSecurityConfigurationsToRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.AttachCodeSecurityConfigurationsToRepositories(ctx, "\n", -1, "", nil)
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.Organizations.AttachCodeSecurityConfigurationsToRepositories(ctx, "o", 1, "selected", []int64{5, 20})
		return resp, err
	})
}

func TestOrganizationsService_SetDefaultCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-security/configurations/1/defaults", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprintf(w, `
		{
			"default_for_new_repos": "all",
			"configuration":
				{
					"id": 1,
					"name": "config1",
					"code_scanning_default_setup": "enabled"
				}
		}`)
	})
	got, resp, err := client.Organizations.SetDefaultCodeSecurityConfiguration(ctx, "o", 1, "all")
	if err != nil {
		t.Errorf("Organizations.SetDefaultCodeSecurityConfiguration returned error: %v", err)
	}
	wantStatus := http.StatusOK
	if resp.StatusCode != wantStatus {
		t.Errorf("Organizations.SetDefaultCodeSecurityConfiguration returned status %d, want %d", resp.StatusCode, wantStatus)
	}
	want := &CodeSecurityConfigurationWithDefaultForNewRepos{
		DefaultForNewRepos: Ptr("all"),
		Configuration: &CodeSecurityConfiguration{
			ID: Ptr(int64(1)), Name: Ptr("config1"), CodeScanningDefaultSetup: Ptr("enabled"),
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Organizations.SetDefaultCodeSecurityConfiguration returned %+v, want %+v", got, want)
	}

	const methodName = "SetDefaultCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.SetDefaultCodeSecurityConfiguration(ctx, "\n", -1, "")
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.SetDefaultCodeSecurityConfiguration(ctx, "o", 1, "all")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetRepositoriesForCodeSecurityConfiguration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-security/configurations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		{
			"id":8,
			"name":"repo8"
		},
		{
			"id":42,
			"name":"repo42"
		}]`)
	})

	repositories, _, err := client.Organizations.GetRepositoriesForCodeSecurityConfiguration(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.GetRepositoriesForCodeSecurityConfiguration returned error: %v", err)
	}

	want := []*Repository{
		{ID: Ptr(int64(8)), Name: Ptr("repo8")},
		{ID: Ptr(int64(42)), Name: Ptr("repo42")},
	}
	if !reflect.DeepEqual(repositories, want) {
		t.Errorf("Organizations.GetRepositoriesForCodeSecurityConfiguration returned %+v, want %+v", repositories, want)
	}

	const methodName = "GetRepositoriesForCodeSecurityConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetRepositoriesForCodeSecurityConfiguration(ctx, "\n", -1)
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetRepositoriesForCodeSecurityConfiguration(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetCodeSecurityConfigurationForRepository(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo8/code-security-configuration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		    "state": "attached",
			"configuration": {
				"id":42,
				"name":"config42",
				"code_scanning_default_setup": "enabled"
			}
		}`)
	})

	rc, _, err := client.Organizations.GetCodeSecurityConfigurationForRepository(ctx, "o", "repo8")
	if err != nil {
		t.Errorf("Organizations.GetCodeSecurityConfigurationForRepository returned error: %v", err)
	}
	c := &CodeSecurityConfiguration{ID: Ptr(int64(42)), Name: Ptr("config42"), CodeScanningDefaultSetup: Ptr("enabled")}
	want := &RepositoryCodeSecurityConfiguration{
		State:         Ptr("attached"),
		Configuration: c,
	}
	if !reflect.DeepEqual(rc, want) {
		t.Errorf("Organizations.GetCodeSecurityConfigurationForRepository returned %+v, want %+v", rc, want)
	}

	const methodName = "GetCodeSecurityConfigurationForRepository"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetCodeSecurityConfigurationForRepository(ctx, "\n", "\n")
		return
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetCodeSecurityConfigurationForRepository(ctx, "o", "repo8")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
