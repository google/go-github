// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestListGlobalSecurityAdvisories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/advisories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"cve_id": "CVE-xoxo-1234"})

		fmt.Fprint(w, `[{
				"id": 1,
				"ghsa_id": "GHSA-xoxo-1234-xoxo",
				"cve_id": "CVE-xoxo-1234",
				"url": "https://api.github.com/advisories/GHSA-xoxo-1234-xoxo",
				"html_url": "https://github.com/advisories/GHSA-xoxo-1234-xoxo",
				"repository_advisory_url": "https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo",
				"summary": "Heartbleed security advisory",
				"description": "This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information.",
				"type": "reviewed",
				"severity": "high",
				"source_code_location": "https://github.com/project/a-package",
				"identifiers": [
					{
						"type": "GHSA",
						"value": "GHSA-xoxo-1234-xoxo"
					},
					{
						"type": "CVE",
						"value": "CVE-xoxo-1234"
					}
				],
				"references": ["https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"],
				"published_at": "1996-06-20T00:00:00Z",
				"updated_at": "1996-06-20T00:00:00Z",
				"github_reviewed_at": "1996-06-20T00:00:00Z",
				"nvd_published_at": "1996-06-20T00:00:00Z",
				"withdrawn_at": null,
				"vulnerabilities": [
					{
						"package": {
							"ecosystem": "npm",
							"name": "a-package"
						},
						"first_patched_version": "1.0.3",
						"vulnerable_version_range": "<=1.0.2",
						"vulnerable_functions": ["a_function"]
					}
				],
				"cvss": {
					"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H",
					"score": 7.6
				},
				"cwes": [
					{
						"cwe_id": "CWE-400",
						"name": "Uncontrolled Resource Consumption"
					}
				],
				"credits": [
					{
						"user": {
							"login": "user",
							"id": 1,
							"node_id": "12=",
							"avatar_url": "a",
							"gravatar_id": "",
							"url": "a",
							"html_url": "b",
							"followers_url": "b",
							"following_url": "c",
							"gists_url": "d",
							"starred_url": "e",
							"subscriptions_url": "f",
							"organizations_url": "g",
							"repos_url": "h",
							"events_url": "i",
							"received_events_url": "j",
							"type": "User",
							"site_admin": false
						},
						"type": "analyst"
					}
				]
			}
		]`)
	})

	ctx := context.Background()
	opts := &ListGlobalSecurityAdvisoriesOptions{CVEID: String("CVE-xoxo-1234")}

	advisories, _, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, opts)
	if err != nil {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*GlobalSecurityAdvisory{
		{
			ID: Int64(1),
			SecurityAdvisory: SecurityAdvisory{
				GHSAID:      String("GHSA-xoxo-1234-xoxo"),
				CVEID:       String("CVE-xoxo-1234"),
				URL:         String("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
				HTMLURL:     String("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
				Severity:    String("high"),
				Summary:     String("Heartbleed security advisory"),
				Description: String("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
				Identifiers: []*AdvisoryIdentifier{
					{
						Type:  String("GHSA"),
						Value: String("GHSA-xoxo-1234-xoxo"),
					},
					{
						Type:  String("CVE"),
						Value: String("CVE-xoxo-1234"),
					},
				},
				PublishedAt: &date,
				UpdatedAt:   &date,
				WithdrawnAt: nil,
				CVSS: &AdvisoryCVSS{
					VectorString: String("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
					Score:        Float64(7.6),
				},
				CWEs: []*AdvisoryCWEs{
					{
						CWEID: String("CWE-400"),
						Name:  String("Uncontrolled Resource Consumption"),
					},
				},
			},
			References: []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
			Vulnerabilities: []*Vulnerabilities{
				{
					Package: &VulnerabilityPackage{
						Ecosystem: String("npm"),
						Name:      String("a-package"),
					},
					FirstPatchedVersion:    String("1.0.3"),
					VulnerableVersionRange: String("<=1.0.2"),
					VulnerableFunctions:    []string{"a_function"},
				},
			},
			RepositoryAdvisoryURL: String("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
			Type:                  String("reviewed"),
			SourceCodeLocation:    String("https://github.com/project/a-package"),
			GitHubReviewedAt:      &date,
			NVDPublishedAt:        &date,
			Credits: []*Credit{
				{
					User: &User{
						Login:             String("user"),
						ID:                Int64(1),
						NodeID:            String("12="),
						AvatarURL:         String("a"),
						GravatarID:        String(""),
						URL:               String("a"),
						HTMLURL:           String("b"),
						FollowersURL:      String("b"),
						FollowingURL:      String("c"),
						GistsURL:          String("d"),
						StarredURL:        String("e"),
						SubscriptionsURL:  String("f"),
						OrganizationsURL:  String("g"),
						ReposURL:          String("h"),
						EventsURL:         String("i"),
						ReceivedEventsURL: String("j"),
						Type:              String("User"),
						SiteAdmin:         Bool(false),
					},
					Type: String("analyst"),
				},
			},
		},
	}

	if !cmp.Equal(advisories, want) {
		t.Errorf("SecurityAdvisories.ListGlobalSecurityAdvisories %+v, want %+v", advisories, want)
	}

	const methodName = "ListGlobalSecurityAdvisories"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecurityAdvisories.ListGlobalSecurityAdvisories(ctx, nil)
		return resp, err
	})
}

func TestGetGlobalSecurityAdvisories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/advisories/GHSA-xoxo-1234-xoxo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"id": 1,
			"ghsa_id": "GHSA-xoxo-1234-xoxo",
			"cve_id": "CVE-xoxo-1234",
			"url": "https://api.github.com/advisories/GHSA-xoxo-1234-xoxo",
			"html_url": "https://github.com/advisories/GHSA-xoxo-1234-xoxo",
			"repository_advisory_url": "https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo",
			"summary": "Heartbleed security advisory",
			"description": "This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information.",
			"type": "reviewed",
			"severity": "high",
			"source_code_location": "https://github.com/project/a-package",
			"identifiers": [
				{
					"type": "GHSA",
					"value": "GHSA-xoxo-1234-xoxo"
				},
				{
					"type": "CVE",
					"value": "CVE-xoxo-1234"
				}
			],
			"references": ["https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"],
			"published_at": "1996-06-20T00:00:00Z",
			"updated_at": "1996-06-20T00:00:00Z",
			"github_reviewed_at": "1996-06-20T00:00:00Z",
			"nvd_published_at": "1996-06-20T00:00:00Z",
			"withdrawn_at": null,
			"vulnerabilities": [
				{
					"package": {
						"ecosystem": "npm",
						"name": "a-package"
					},
					"first_patched_version": "1.0.3",
					"vulnerable_version_range": "<=1.0.2",
					"vulnerable_functions": ["a_function"]
				}
			],
			"cvss": {
				"vector_string": "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H",
				"score": 7.6
			},
			"cwes": [
				{
					"cwe_id": "CWE-400",
					"name": "Uncontrolled Resource Consumption"
				}
			],
			"credits": [
				{
					"user": {
						"login": "user",
						"id": 1,
						"node_id": "12=",
						"avatar_url": "a",
						"gravatar_id": "",
						"url": "a",
						"html_url": "b",
						"followers_url": "b",
						"following_url": "c",
						"gists_url": "d",
						"starred_url": "e",
						"subscriptions_url": "f",
						"organizations_url": "g",
						"repos_url": "h",
						"events_url": "i",
						"received_events_url": "j",
						"type": "User",
						"site_admin": false
					},
					"type": "analyst"
				}
			]
		}`)
	})

	ctx := context.Background()
	advisory, _, err := client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "GHSA-xoxo-1234-xoxo")
	if err != nil {
		t.Errorf("SecurityAdvisories.GetGlobalSecurityAdvisories returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := &GlobalSecurityAdvisory{
		ID: Int64(1),
		SecurityAdvisory: SecurityAdvisory{
			GHSAID:      String("GHSA-xoxo-1234-xoxo"),
			CVEID:       String("CVE-xoxo-1234"),
			URL:         String("https://api.github.com/advisories/GHSA-xoxo-1234-xoxo"),
			HTMLURL:     String("https://github.com/advisories/GHSA-xoxo-1234-xoxo"),
			Severity:    String("high"),
			Summary:     String("Heartbleed security advisory"),
			Description: String("This bug allows an attacker to read portions of the affected server’s memory, potentially disclosing sensitive information."),
			Identifiers: []*AdvisoryIdentifier{
				{
					Type:  String("GHSA"),
					Value: String("GHSA-xoxo-1234-xoxo"),
				},
				{
					Type:  String("CVE"),
					Value: String("CVE-xoxo-1234"),
				},
			},
			PublishedAt: &date,
			UpdatedAt:   &date,
			WithdrawnAt: nil,
			CVSS: &AdvisoryCVSS{
				VectorString: String("CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:C/C:H/I:H/A:H"),
				Score:        Float64(7.6),
			},
			CWEs: []*AdvisoryCWEs{
				{
					CWEID: String("CWE-400"),
					Name:  String("Uncontrolled Resource Consumption"),
				},
			},
		},
		RepositoryAdvisoryURL: String("https://api.github.com/repos/project/a-package/security-advisories/GHSA-xoxo-1234-xoxo"),
		Type:                  String("reviewed"),
		SourceCodeLocation:    String("https://github.com/project/a-package"),
		References:            []string{"https://nvd.nist.gov/vuln/detail/CVE-xoxo-1234"},
		GitHubReviewedAt:      &date,
		NVDPublishedAt:        &date,

		Vulnerabilities: []*Vulnerabilities{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: String("npm"),
					Name:      String("a-package"),
				},
				FirstPatchedVersion:    String("1.0.3"),
				VulnerableVersionRange: String("<=1.0.2"),
				VulnerableFunctions:    []string{"a_function"},
			},
		},
		Credits: []*Credit{
			{
				User: &User{
					Login:             String("user"),
					ID:                Int64(1),
					NodeID:            String("12="),
					AvatarURL:         String("a"),
					GravatarID:        String(""),
					URL:               String("a"),
					HTMLURL:           String("b"),
					FollowersURL:      String("b"),
					FollowingURL:      String("c"),
					GistsURL:          String("d"),
					StarredURL:        String("e"),
					SubscriptionsURL:  String("f"),
					OrganizationsURL:  String("g"),
					ReposURL:          String("h"),
					EventsURL:         String("i"),
					ReceivedEventsURL: String("j"),
					Type:              String("User"),
					SiteAdmin:         Bool(false),
				},
				Type: String("analyst"),
			},
		},
	}

	if !cmp.Equal(advisory, want) {
		t.Errorf("SecurityAdvisories.GetGlobalSecurityAdvisories %+v, want %+v", advisory, want)
	}

	const methodName = "GetGlobalSecurityAdvisories"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "CVE-\n-1234")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.SecurityAdvisories.GetGlobalSecurityAdvisories(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
