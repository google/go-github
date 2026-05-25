// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDependabotService_ListRepoAlerts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open"})
		fmt.Fprint(w, `[{"number":1,"state":"open"},{"number":42,"state":"fixed"}]`)
	})

	opts := &ListAlertsOptions{State: Ptr("open")}
	ctx := t.Context()
	alerts, _, err := client.Dependabot.ListRepoAlerts(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Dependabot.ListRepoAlerts returned error: %v", err)
	}

	want := []*DependabotAlert{
		{Number: Ptr(1), State: Ptr("open")},
		{Number: Ptr(42), State: Ptr("fixed")},
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("Dependabot.ListRepoAlerts returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListRepoAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListRepoAlerts(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListRepoAlerts(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_GetRepoAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependabot/alerts/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":42,"state":"fixed"}`)
	})

	ctx := t.Context()
	alert, _, err := client.Dependabot.GetRepoAlert(ctx, "o", "r", 42)
	if err != nil {
		t.Errorf("Dependabot.GetRepoAlert returned error: %v", err)
	}

	want := &DependabotAlert{
		Number: Ptr(42),
		State:  Ptr("fixed"),
	}
	if !cmp.Equal(alert, want) {
		t.Errorf("Dependabot.GetRepoAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "GetRepoAlert"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetRepoAlert(ctx, "\n", "\n", 0)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetRepoAlert(ctx, "o", "r", 42)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_ListOrgAlerts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/dependabot/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open"})
		fmt.Fprint(w, `[{"number":1,"state":"open"},{"number":42,"state":"fixed"}]`)
	})

	opts := &ListAlertsOptions{State: Ptr("open")}
	ctx := t.Context()
	alerts, _, err := client.Dependabot.ListOrgAlerts(ctx, "o", opts)
	if err != nil {
		t.Errorf("Dependabot.ListOrgAlerts returned error: %v", err)
	}

	want := []*DependabotAlert{
		{Number: Ptr(1), State: Ptr("open")},
		{Number: Ptr(42), State: Ptr("fixed")},
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("Dependabot.ListOrgAlerts returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListOrgAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListOrgAlerts(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListOrgAlerts(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_UpdateAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	state := Ptr("dismissed")
	dismissedReason := Ptr("no_bandwidth")
	dismissedComment := Ptr("no time to fix this")

	alertState := &DependabotAlertState{State: *state, DismissedReason: dismissedReason, DismissedComment: dismissedComment}

	mux.HandleFunc("/repos/o/r/dependabot/alerts/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"number":42,"state":"dismissed","dismissed_reason":"no_bandwidth","dismissed_comment":"no time to fix this"}`)
	})

	ctx := t.Context()
	alert, _, err := client.Dependabot.UpdateAlert(ctx, "o", "r", 42, alertState)
	if err != nil {
		t.Errorf("Dependabot.UpdateAlert returned error: %v", err)
	}

	want := &DependabotAlert{
		Number:           Ptr(42),
		State:            Ptr("dismissed"),
		DismissedReason:  Ptr("no_bandwidth"),
		DismissedComment: Ptr("no time to fix this"),
	}
	if !cmp.Equal(alert, want) {
		t.Errorf("Dependabot.UpdateAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "UpdateAlert"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.UpdateAlert(ctx, "\n", "\n", 0, alertState)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.UpdateAlert(ctx, "o", "r", 42, alertState)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependency_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Dependency{}, "{}")

	h := &Dependency{
		Package: &VulnerabilityPackage{
			Ecosystem: Ptr("pip"),
			Name:      Ptr("django"),
		},
		ManifestPath: Ptr("path/to/requirements.txt"),
		Scope:        Ptr("runtime"),
	}

	want := `{
		"package": {
        "ecosystem": "pip",
        "name": "django"
      },
      "manifest_path": "path/to/requirements.txt",
      "scope": "runtime"
	}`

	testJSONMarshal(t, h, want)
}

func TestAdvisoryCVSS_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AdvisoryCVSS{}, "{}")

	h := &AdvisoryCVSS{
		Score:        Ptr(7.5),
		VectorString: Ptr("CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N"),
	}

	want := `{
		"vector_string": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N",
        "score": 7.5
	}`

	testJSONMarshal(t, h, want)
}

func TestAdvisoryCWEs_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AdvisoryCWEs{}, "{}")

	h := &AdvisoryCWEs{
		CWEID: Ptr("CWE-200"),
		Name:  Ptr("Exposure of Sensitive Information to an Unauthorized Actor"),
	}

	want := `{
		"cwe_id": "CWE-200",
		"name": "Exposure of Sensitive Information to an Unauthorized Actor"
	}`

	testJSONMarshal(t, h, want)
}

func TestAdvisoryEPSS_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AdvisoryEPSS{}, `{"percentage": 0, "percentile": 0}`)

	h := &AdvisoryEPSS{
		Percentage: 0.05,
		Percentile: 0.5,
	}

	want := `{
		"percentage": 0.05,
		"percentile": 0.5
	}`

	testJSONMarshal(t, h, want)
}

func TestDependabotSecurityAdvisory_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DependabotSecurityAdvisory{}, "{}")

	publishedAt, _ := time.Parse(time.RFC3339, "2018-10-03T21:13:54Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2022-04-26T18:35:37Z")

	h := &DependabotSecurityAdvisory{
		GHSAID:      Ptr("GHSA-rf4j-j272-fj86"),
		CVEID:       Ptr("CVE-2018-6188"),
		Summary:     Ptr("Django allows remote attackers to obtain potentially sensitive information by leveraging data exposure from the confirm_login_allowed() method, as demonstrated by discovering whether a user account is inactive"),
		Description: Ptr("django.contrib.auth.forms.AuthenticationForm in Django 2.0 before 2.0.2, and 1.11.8 and 1.11.9, allows remote attackers to obtain potentially sensitive information by leveraging data exposure from the confirm_login_allowed() method, as demonstrated by discovering whether a user account is inactive."),
		Vulnerabilities: []*AdvisoryVulnerability{
			{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("pip"),
					Name:      Ptr("django"),
				},
				Severity:               Ptr("high"),
				VulnerableVersionRange: Ptr(">= 2.0.0, < 2.0.2"),
				FirstPatchedVersion:    &FirstPatchedVersion{Identifier: Ptr("2.0.2")},
			},
			{
				Package: &VulnerabilityPackage{
					Ecosystem: Ptr("pip"),
					Name:      Ptr("django"),
				},
				Severity:               Ptr("high"),
				VulnerableVersionRange: Ptr(">= 1.11.8, < 1.11.10"),
				FirstPatchedVersion:    &FirstPatchedVersion{Identifier: Ptr("1.11.10")},
			},
		},
		Severity: Ptr("high"),
		CVSS: &AdvisoryCVSS{
			Score:        Ptr(7.5),
			VectorString: Ptr("CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N"),
		},
		CWEs: []*AdvisoryCWEs{
			{
				CWEID: Ptr("CWE-200"),
				Name:  Ptr("Exposure of Sensitive Information to an Unauthorized Actor"),
			},
		},
		EPSS: &AdvisoryEPSS{
			Percentage: 0.05,
			Percentile: 0.5,
		},
		Identifiers: []*AdvisoryIdentifier{
			{
				Type:  Ptr("GHSA"),
				Value: Ptr("GHSA-rf4j-j272-fj86"),
			},
			{
				Type:  Ptr("CVE"),
				Value: Ptr("CVE-2018-6188"),
			},
		},
		References: []*AdvisoryReference{
			{
				URL: Ptr("https://example.com/vuln/detail/CVE-2018-6188"),
			},
			{
				URL: Ptr("https://github.com/advisories/GHSA-rf4j-j272-fj86"),
			},
			{
				URL: Ptr("https://example.com/3559-1/"),
			},
			{
				URL: Ptr("https://example.com/weblog/2018/feb/01/security-releases/"),
			},
			{
				URL: Ptr("https://example.com/id/1040422"),
			},
		},
		PublishedAt: &Timestamp{publishedAt},
		UpdatedAt:   &Timestamp{updatedAt},
		WithdrawnAt: nil,
	}

	want := `{
	  "ghsa_id": "GHSA-rf4j-j272-fj86",
      "cve_id": "CVE-2018-6188",
      "summary": "Django allows remote attackers to obtain potentially sensitive information by leveraging data exposure from the confirm_login_allowed() method, as demonstrated by discovering whether a user account is inactive",
      "description": "django.contrib.auth.forms.AuthenticationForm in Django 2.0 before 2.0.2, and 1.11.8 and 1.11.9, allows remote attackers to obtain potentially sensitive information by leveraging data exposure from the confirm_login_allowed() method, as demonstrated by discovering whether a user account is inactive.",
      "vulnerabilities": [
        {
          "package": {
            "ecosystem": "pip",
            "name": "django"
          },
          "severity": "high",
          "vulnerable_version_range": ">= 2.0.0, < 2.0.2",
          "first_patched_version": {
            "identifier": "2.0.2"
          }
        },
        {
          "package": {
            "ecosystem": "pip",
            "name": "django"
          },
          "severity": "high",
          "vulnerable_version_range": ">= 1.11.8, < 1.11.10",
          "first_patched_version": {
            "identifier": "1.11.10"
          }
        }
      ],
      "severity": "high",
      "cvss": {
        "vector_string": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N",
        "score": 7.5
      },
      "cwes": [
        {
          "cwe_id": "CWE-200",
          "name": "Exposure of Sensitive Information to an Unauthorized Actor"
        }
      ],
      "epss": {
        "percentage": 0.05,
        "percentile": 0.5
      },
      "identifiers": [
        {
          "type": "GHSA",
          "value": "GHSA-rf4j-j272-fj86"
        },
        {
          "type": "CVE",
          "value": "CVE-2018-6188"
        }
      ],
      "references": [
        {
          "url": "https://example.com/vuln/detail/CVE-2018-6188"
        },
        {
          "url": "https://github.com/advisories/GHSA-rf4j-j272-fj86"
        },
        {
          "url": "https://example.com/3559-1/"
        },
        {
          "url": "https://example.com/weblog/2018/feb/01/security-releases/"
        },
        {
          "url": "https://example.com/id/1040422"
        }
      ],
      "published_at": "2018-10-03T21:13:54Z",
      "updated_at": "2022-04-26T18:35:37Z"
	}`

	testJSONMarshal(t, h, want)
}

func TestDependabotAlert_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DependabotAlert{}, "{}")

	h := &DependabotAlert{
		Number: Ptr(42),
		State:  Ptr("dismissed"),
		Dependency: &Dependency{
			Package: &VulnerabilityPackage{
				Ecosystem: Ptr("npm"),
				Name:      Ptr("minimist"),
			},
			ManifestPath: Ptr("package-lock.json"),
			Scope:        Ptr("runtime"),
		},
		SecurityAdvisory: &DependabotSecurityAdvisory{
			GHSAID:   Ptr("GHSA-vh95-rmgr-6w4m"),
			CVEID:    Ptr("CVE-2020-7598"),
			Summary:  Ptr("Prototype pollution in minimist"),
			Severity: Ptr("high"),
			EPSS: &AdvisoryEPSS{
				Percentage: 0.05,
				Percentile: 0.5,
			},
		},
		SecurityVulnerability: &AdvisoryVulnerability{
			Package: &VulnerabilityPackage{
				Ecosystem: Ptr("npm"),
				Name:      Ptr("minimist"),
			},
			Severity:               Ptr("high"),
			VulnerableVersionRange: Ptr("< 1.2.3"),
			FirstPatchedVersion:    &FirstPatchedVersion{Identifier: Ptr("1.2.3")},
			PatchedVersions:        Ptr(">= 1.2.3"),
			VulnerableFunctions:    []string{"parse"},
		},
		URL:              Ptr("https://api.github.com/repos/o/r/dependabot/alerts/42"),
		HTMLURL:          Ptr("https://github.com/o/r/security/dependabot/42"),
		CreatedAt:        &Timestamp{referenceTime},
		UpdatedAt:        &Timestamp{referenceTime},
		DismissedAt:      &Timestamp{referenceTime},
		DismissedBy:      &User{Login: Ptr("octocat"), ID: Ptr(int64(1))},
		DismissedReason:  Ptr("tolerable_risk"),
		DismissedComment: Ptr("risk accepted"),
		FixedAt:          &Timestamp{referenceTime},
		AutoDismissedAt:  &Timestamp{referenceTime},
		Repository: &Repository{
			Owner:    &User{Login: Ptr("o")},
			Name:     Ptr("r"),
			FullName: Ptr("o/r"),
			Private:  Ptr(false),
		},
	}

	want := `{
		"number": 42,
		"state": "dismissed",
		"dependency": {
			"package": {
				"ecosystem": "npm",
				"name": "minimist"
			},
			"manifest_path": "package-lock.json",
			"scope": "runtime"
		},
		"security_advisory": {
			"ghsa_id": "GHSA-vh95-rmgr-6w4m",
			"cve_id": "CVE-2020-7598",
			"summary": "Prototype pollution in minimist",
			"severity": "high",
			"epss": {
				"percentage": 0.05,
				"percentile": 0.5
			}
		},
		"security_vulnerability": {
			"package": {
				"ecosystem": "npm",
				"name": "minimist"
			},
			"severity": "high",
			"vulnerable_version_range": "< 1.2.3",
			"first_patched_version": {
				"identifier": "1.2.3"
			},
			"patched_versions": ">= 1.2.3",
			"vulnerable_functions": [
				"parse"
			]
		},
		"url": "https://api.github.com/repos/o/r/dependabot/alerts/42",
		"html_url": "https://github.com/o/r/security/dependabot/42",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"dismissed_at": ` + referenceTimeStr + `,
		"dismissed_by": {
			"login": "octocat",
			"id": 1
		},
		"dismissed_reason": "tolerable_risk",
		"dismissed_comment": "risk accepted",
		"fixed_at": ` + referenceTimeStr + `,
		"auto_dismissed_at": ` + referenceTimeStr + `,
		"repository": {
			"owner": {
				"login": "o"
			},
			"name": "r",
			"full_name": "o/r",
			"private": false
		}
	}`

	testJSONMarshal(t, h, want)
}

func TestDependabotAlertState_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DependabotAlertState{}, `{"state": ""}`)

	h := &DependabotAlertState{
		State:            "dismissed",
		DismissedReason:  Ptr("no_bandwidth"),
		DismissedComment: Ptr("no time to fix this"),
	}

	want := `{
		"state": "dismissed",
		"dismissed_reason": "no_bandwidth",
		"dismissed_comment": "no time to fix this"
	}`

	testJSONMarshal(t, h, want)
}
