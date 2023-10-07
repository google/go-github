// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSecurityAdvisoriesService_RequestCVE(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_ok/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_accepted/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_ok")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	_, err = client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_accepted")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	const methodName = "RequestCVE"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecurityAdvisories.RequestCVE(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id")
		if err == nil {
			t.Errorf("testNewRequestAndDoFailure %v should have return err", methodName)
		}
		return resp, err
	})
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_BadRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_UnmarshalError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err == nil {
		t.Errorf("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal number into Go struct field SecurityAdvisory.ghsa_id of type string") {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[
			{
				"ghsa_id": "GHSA-abcd-1234-efgh",
   				"cve_id": "CVE-2050-00000"
 			}
		]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg return status %d, want %d", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: String("GHSA-abcd-1234-efgh"),
			CVEID:  String("CVE-2050-00000"),
		},
	}
	if !cmp.Equal(advisories, want) {
		t.Errorf("ListRepositorySecurityAdvisoriesForOrg returned %+v, want %+v", advisories, want)
	}

	methodName := "ListRepositorySecurityAdvisoriesForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "\n", &ListRepositorySecurityAdvisoriesOptions{
			Sort: "\n",
		})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_BadRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		http.Error(w, "Bad Request", 400)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		query := r.URL.Query()
		if query.Get("state") != "draft" {
			t.Errorf("ListRepositorySecurityAdvisories returned %+v, want %+v", query.Get("state"), "draft")
		}

		http.NotFound(w, r)
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", &ListRepositorySecurityAdvisoriesOptions{
		State: "draft",
	})
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories_UnmarshalError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[{"ghsa_id": 12334354}]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected unmarshal error")
	} else if !strings.Contains(err.Error(), "json: cannot unmarshal number into Go struct field SecurityAdvisory.ghsa_id of type string") {
		t.Errorf("ListRepositorySecurityAdvisories returned unexpected error: %v", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}
	if advisories != nil {
		t.Errorf("ListRepositorySecurityAdvisories return %+v, want nil", advisories)
	}
}

func TestSecurityAdvisoriesService_ListRepositorySecurityAdvisories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`[
			{
				"ghsa_id": "GHSA-abcd-1234-efgh",
   				"cve_id": "CVE-2050-00000"
 			}
		]`))
	})

	ctx := context.Background()
	advisories, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("ListRepositorySecurityAdvisories returned error: %v, want nil", err)
	}
	if got, want := resp.Response.StatusCode, http.StatusOK; got != want {
		t.Errorf("ListRepositorySecurityAdvisories return status %d, want %d", got, want)
	}

	want := []*SecurityAdvisory{
		{
			GHSAID: String("GHSA-abcd-1234-efgh"),
			CVEID:  String("CVE-2050-00000"),
		},
	}
	if !cmp.Equal(advisories, want) {
		t.Errorf("ListRepositorySecurityAdvisories returned %+v, want %+v", advisories, want)
	}

	methodName := "ListRepositorySecurityAdvisories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "\n", "\n", &ListRepositorySecurityAdvisoriesOptions{
			Sort: "\n",
		})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.ListRepositorySecurityAdvisories(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
