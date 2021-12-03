// Copyright 2020 The go-github AUTHORS. All rights reserved.
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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCodeScanningService_Alert_ID(t *testing.T) {
	// Test: nil Alert ID == 0
	var a *Alert
	id := a.ID()
	var want int64
	if id != want {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}

	// Test: Valid HTMLURL
	a = &Alert{
		HTMLURL: String("https://github.com/o/r/security/code-scanning/88"),
	}
	id = a.ID()
	want = 88
	if !cmp.Equal(id, want) {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}

	// Test: HTMLURL is nil
	a = &Alert{}
	id = a.ID()
	want = 0
	if !cmp.Equal(id, want) {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}

	// Test: ID can't be parsed as an int
	a = &Alert{
		HTMLURL: String("https://github.com/o/r/security/code-scanning/bad88"),
	}
	id = a.ID()
	want = 0
	if !cmp.Equal(id, want) {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}
}

func TestCodeScanningService_UploadSarif(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/sarifs", func(w http.ResponseWriter, r *http.Request) {
		v := new(SarifAnalysis)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")
		want := &SarifAnalysis{CommitSHA: String("abc"), Ref: String("ref/head/main"), Sarif: String("abc"), CheckoutURI: String("uri"), StartedAt: &Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)}, ToolName: String("codeql-cli")}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"commit_sha":"abc","ref":"ref/head/main","sarif":"abc"}`)
	})

	ctx := context.Background()
	sarifAnalysis := &SarifAnalysis{CommitSHA: String("abc"), Ref: String("ref/head/main"), Sarif: String("abc"), CheckoutURI: String("uri"), StartedAt: &Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)}, ToolName: String("codeql-cli")}
	_, _, err := client.CodeScanning.UploadSarif(ctx, "o", "r", sarifAnalysis)
	if err != nil {
		t.Errorf("CodeScanning.UploadSarif returned error: %v", err)
	}

	const methodName = "UploadSarif"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.UploadSarif(ctx, "\n", "\n", sarifAnalysis)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.CodeScanning.UploadSarif(ctx, "o", "r", sarifAnalysis)
		return resp, err
	})
}

func TestCodeScanningService_ListAlertsForRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "ref": "heads/master"})
		fmt.Fprint(w, `[{
				"rule_id":"js/trivial-conditional",
				"rule_severity":"warning",
				"rule_description":"Useless conditional",
				"tool": {
					"name": "CodeQL",
					"guid": null,
					"version": "1.4.0"
				},
				"rule": {
					"id": "js/trivial-conditional",
					"severity": "warning",
					"description": "Useless conditional",
					"name": "js/trivial-conditional",
					"full_description": "Expression has no effect",
					"help": "Expression has no effect"
				},
				"most_recent_instance": {
					"ref": "refs/heads/main",
					"state": "open",
					"commit_sha": "abcdefg12345",
					"message": {
						"text": "This path depends on a user-provided value."
					},
					"location": {
						"path": "spec-main/api-session-spec.ts",
						"start_line": 917,
						"end_line": 917,
						"start_column": 7,
						"end_column": 18
					},
					"classifications": [
						"test"
					]
				},
				"created_at":"2020-05-06T12:00:00Z",
				"state":"open",
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/25",
				"html_url":"https://github.com/o/r/security/code-scanning/25"
				},
				{
				"rule_id":"js/useless-expression",
				"rule_severity":"warning",
				"rule_description":"Expression has no effect",
				"tool": {
					"name": "CodeQL",
					"guid": null,
					"version": "1.4.0"
				},
				"rule": {
					"id": "js/useless-expression",
					"severity": "warning",
					"description": "Expression has no effect",
					"name": "js/useless-expression",
					"full_description": "Expression has no effect",
					"help": "Expression has no effect"
				},
				"most_recent_instance": {
					"ref": "refs/heads/main",
					"state": "open",
					"commit_sha": "abcdefg12345",
					"message": {
						"text": "This path depends on a user-provided value."
					},
					"location": {
						"path": "spec-main/api-session-spec.ts",
						"start_line": 917,
						"end_line": 917,
						"start_column": 7,
						"end_column": 18
					},
					"classifications": [
						"test"
					]
				},
				"created_at":"2020-05-06T12:00:00Z",
				"state":"open",
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/88",
				"html_url":"https://github.com/o/r/security/code-scanning/88"
				}]`)
	})

	opts := &AlertListOptions{State: "open", Ref: "heads/master"}
	ctx := context.Background()
	alerts, _, err := client.CodeScanning.ListAlertsForRepo(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAlertsForRepo returned error: %v", err)
	}

	date := Timestamp{time.Date(2020, time.May, 06, 12, 00, 00, 0, time.UTC)}
	want := []*Alert{
		{
			RuleID:          String("js/trivial-conditional"),
			RuleSeverity:    String("warning"),
			RuleDescription: String("Useless conditional"),
			Tool:            &Tool{Name: String("CodeQL"), GUID: nil, Version: String("1.4.0")},
			Rule: &Rule{
				ID:              String("js/trivial-conditional"),
				Severity:        String("warning"),
				Description:     String("Useless conditional"),
				Name:            String("js/trivial-conditional"),
				FullDescription: String("Expression has no effect"),
				Help:            String("Expression has no effect"),
			},
			CreatedAt: &date,
			State:     String("open"),
			ClosedBy:  nil,
			ClosedAt:  nil,
			URL:       String("https://api.github.com/repos/o/r/code-scanning/alerts/25"),
			HTMLURL:   String("https://github.com/o/r/security/code-scanning/25"),
			MostRecentInstance: &MostRecentInstance{
				Ref:       String("refs/heads/main"),
				State:     String("open"),
				CommitSHA: String("abcdefg12345"),
				Message: &Message{
					Text: String("This path depends on a user-provided value."),
				},
				Location: &Location{
					Path:        String("spec-main/api-session-spec.ts"),
					StartLine:   Int(917),
					EndLine:     Int(917),
					StartColumn: Int(7),
					EndColumn:   Int(18),
				},
				Classifications: []string{"test"},
			},
		},
		{
			RuleID:          String("js/useless-expression"),
			RuleSeverity:    String("warning"),
			RuleDescription: String("Expression has no effect"),
			Tool:            &Tool{Name: String("CodeQL"), GUID: nil, Version: String("1.4.0")},
			Rule: &Rule{
				ID:              String("js/useless-expression"),
				Severity:        String("warning"),
				Description:     String("Expression has no effect"),
				Name:            String("js/useless-expression"),
				FullDescription: String("Expression has no effect"),
				Help:            String("Expression has no effect"),
			},
			CreatedAt: &date,
			State:     String("open"),
			ClosedBy:  nil,
			ClosedAt:  nil,
			URL:       String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
			HTMLURL:   String("https://github.com/o/r/security/code-scanning/88"),
			MostRecentInstance: &MostRecentInstance{
				Ref:       String("refs/heads/main"),
				State:     String("open"),
				CommitSHA: String("abcdefg12345"),
				Message: &Message{
					Text: String("This path depends on a user-provided value."),
				},
				Location: &Location{
					Path:        String("spec-main/api-session-spec.ts"),
					StartLine:   Int(917),
					EndLine:     Int(917),
					StartColumn: Int(7),
					EndColumn:   Int(18),
				},
				Classifications: []string{"test"},
			},
		},
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("CodeScanning.ListAlertsForRepo returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListAlertsForRepo(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListAlertsForRepo(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_GetAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/alerts/88", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"rule_id":"js/useless-expression",
				"rule_severity":"warning",
				"rule_description":"Expression has no effect",
				"tool": {
					"name": "CodeQL",
					"guid": null,
					"version": "1.4.0"
				},
				"rule": {
					"id": "useless expression",
					"severity": "warning",
					"description": "Expression has no effect",
					"name": "useless expression",
					"full_description": "Expression has no effect",
					"help": "Expression has no effect"
				},
				"most_recent_instance": {
					"ref": "refs/heads/main",
					"state": "open",
					"commit_sha": "abcdefg12345",
					"message": {
						"text": "This path depends on a user-provided value."
					},
					"location": {
						"path": "spec-main/api-session-spec.ts",
						"start_line": 917,
						"end_line": 917,
						"start_column": 7,
						"end_column": 18
					},
					"classifications": [
						"test"
					]
				},      
				"created_at":"2019-01-02T15:04:05Z",
				"state":"open",
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/88",
				"html_url":"https://github.com/o/r/security/code-scanning/88"}`)
	})

	ctx := context.Background()
	alert, _, err := client.CodeScanning.GetAlert(ctx, "o", "r", 88)
	if err != nil {
		t.Errorf("CodeScanning.GetAlert returned error: %v", err)
	}

	date := Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}
	want := &Alert{
		RuleID:          String("js/useless-expression"),
		RuleSeverity:    String("warning"),
		RuleDescription: String("Expression has no effect"),
		Tool:            &Tool{Name: String("CodeQL"), GUID: nil, Version: String("1.4.0")},
		Rule: &Rule{
			ID:              String("useless expression"),
			Severity:        String("warning"),
			Description:     String("Expression has no effect"),
			Name:            String("useless expression"),
			FullDescription: String("Expression has no effect"),
			Help:            String("Expression has no effect"),
		},
		CreatedAt: &date,
		State:     String("open"),
		ClosedBy:  nil,
		ClosedAt:  nil,
		URL:       String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
		HTMLURL:   String("https://github.com/o/r/security/code-scanning/88"),
		MostRecentInstance: &MostRecentInstance{
			Ref:       String("refs/heads/main"),
			State:     String("open"),
			CommitSHA: String("abcdefg12345"),
			Message: &Message{
				Text: String("This path depends on a user-provided value."),
			},
			Location: &Location{
				Path:        String("spec-main/api-session-spec.ts"),
				StartLine:   Int(917),
				EndLine:     Int(917),
				StartColumn: Int(7),
				EndColumn:   Int(18),
			},
			Classifications: []string{"test"},
		},
	}
	if !cmp.Equal(alert, want) {
		t.Errorf("CodeScanning.GetAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "GetAlert"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.GetAlert(ctx, "\n", "\n", -88)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.GetAlert(ctx, "o", "r", 88)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAlert_Marshal(t *testing.T) {
	testJSONMarshal(t, &Alert{}, "{}")

	u := &Alert{
		RuleID:          String("rid"),
		RuleSeverity:    String("rs"),
		RuleDescription: String("rd"),
		Tool: &Tool{
			Name:    String("n"),
			GUID:    String("g"),
			Version: String("v"),
		},
		CreatedAt: &Timestamp{referenceTime},
		State:     String("fixed"),
		ClosedBy: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		ClosedAt: &Timestamp{referenceTime},
		URL:      String("url"),
		HTMLURL:  String("hurl"),
	}

	want := `{
		"rule_id": "rid",
		"rule_severity": "rs",
		"rule_description": "rd",
		"tool": {
			"name": "n",
			"guid": "g",
			"version": "v"
		},
		"created_at": ` + referenceTimeStr + `,
		"state": "fixed",
		"closed_by": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"closed_at": ` + referenceTimeStr + `,
		"url": "url",
		"html_url": "hurl"
	}`

	testJSONMarshal(t, u, want)
}

func TestLocation_Marshal(t *testing.T) {
	testJSONMarshal(t, &Location{}, "{}")

	u := &Location{
		Path:        String("path"),
		StartLine:   Int(1),
		EndLine:     Int(2),
		StartColumn: Int(3),
		EndColumn:   Int(4),
	}

	want := `{
		"path": "path",
		"start_line": 1,
		"end_line": 2,
		"start_column": 3,
		"end_column": 4
	}`

	testJSONMarshal(t, u, want)
}

func TestRule_Marshal(t *testing.T) {
	testJSONMarshal(t, &Rule{}, "{}")

	u := &Rule{
		ID:                    String("1"),
		Severity:              String("3"),
		Description:           String("description"),
		Name:                  String("first"),
		SecuritySeverityLevel: String("2"),
		FullDescription:       String("summary"),
		Tags:                  []string{"tag1", "tag2"},
		Help:                  String("Help Text"),
	}

	want := `{
		"id":                      "1",
		"severity":                "3",
		"description":             "description",
		"name":                    "first",
		"security_severity_level": "2",
		"full_description":        "summary",
		"tags":                    ["tag1", "tag2"],
		"help":                    "Help Text"
	}`

	testJSONMarshal(t, u, want)
}

func TestTool_Marshal(t *testing.T) {
	testJSONMarshal(t, &Tool{}, "{}")

	u := &Tool{
		Name:    String("name"),
		GUID:    String("guid"),
		Version: String("ver"),
	}

	want := `{
		"name": "name",
		"guid": "guid",
		"version": "ver"
	}`

	testJSONMarshal(t, u, want)
}

func TestMessage_Marshal(t *testing.T) {
	testJSONMarshal(t, &Message{}, "{}")

	u := &Message{
		Text: String("text"),
	}

	want := `{
		"text": "text"
	}`

	testJSONMarshal(t, u, want)
}

func TestCodeScanningService_ListAnalysesForRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/analyses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"sarif_id": "8981cd8e-b078-4ac3-a3be-1dad7dbd0b582", "ref": "heads/master"})
		fmt.Fprint(w, `[
			  {
				"ref": "refs/heads/main",
				"commit_sha": "d99612c3e1f2970085cfbaeadf8f010ef69bad83",
				"analysis_key": ".github/workflows/codeql-analysis.yml:analyze",
				"environment": "{\"language\":\"python\"}",
				"error": "",
				"category": ".github/workflows/codeql-analysis.yml:analyze/language:python",
				"created_at": "2020-08-27T15:05:21Z",
				"results_count": 17,
				"rules_count": 49,
				"id": 201,
				"url": "https://api.github.com/repos/o/r/code-scanning/analyses/201",
				"sarif_id": "8981cd8e-b078-4ac3-a3be-1dad7dbd0b582",
				"tool": {
				  "name": "CodeQL",
				  "guid": null,
				  "version": "2.4.0"
				},
				"deletable": true,
				"warning": ""
			  },
			  {
				"ref": "refs/heads/my-branch",
				"commit_sha": "c8cff6510d4d084fb1b4aa13b64b97ca12b07321",
				"analysis_key": ".github/workflows/shiftleft.yml:build",
				"environment": "{}",
				"error": "",
				"category": ".github/workflows/shiftleft.yml:build/",
				"created_at": "2020-08-27T15:05:21Z",
				"results_count": 17,
				"rules_count": 32,
				"id": 200,
				"url": "https://api.github.com/repos/o/r/code-scanning/analyses/200",
				"sarif_id": "8981cd8e-b078-4ac3-a3be-1dad7dbd0b582",
				"tool": {
				  "name": "Python Security ScanningAnalysis",
				  "guid": null,
				  "version": "1.2.0"
				},
				"deletable": true,
				"warning": ""
			  }
			]`)
	})

	opts := &AnalysesListOptions{SarifID: String("8981cd8e-b078-4ac3-a3be-1dad7dbd0b582"), Ref: String("heads/master")}
	ctx := context.Background()
	analyses, _, err := client.CodeScanning.ListAnalysesForRepo(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAnalysesForRepo returned error: %v", err)
	}

	date := &Timestamp{time.Date(2020, time.August, 27, 15, 05, 21, 0, time.UTC)}
	want := []*ScanningAnalysis{
		{
			ID:           Int64(201),
			Ref:          String("refs/heads/main"),
			CommitSHA:    String("d99612c3e1f2970085cfbaeadf8f010ef69bad83"),
			AnalysisKey:  String(".github/workflows/codeql-analysis.yml:analyze"),
			Environment:  String("{\"language\":\"python\"}"),
			Error:        String(""),
			Category:     String(".github/workflows/codeql-analysis.yml:analyze/language:python"),
			CreatedAt:    date,
			ResultsCount: Int(17),
			RulesCount:   Int(49),
			URL:          String("https://api.github.com/repos/o/r/code-scanning/analyses/201"),
			SarifID:      String("8981cd8e-b078-4ac3-a3be-1dad7dbd0b582"),
			Tool: &Tool{
				Name:    String("CodeQL"),
				GUID:    nil,
				Version: String("2.4.0"),
			},
			Deletable: Bool(true),
			Warning:   String(""),
		},
		{
			ID:           Int64(200),
			Ref:          String("refs/heads/my-branch"),
			CommitSHA:    String("c8cff6510d4d084fb1b4aa13b64b97ca12b07321"),
			AnalysisKey:  String(".github/workflows/shiftleft.yml:build"),
			Environment:  String("{}"),
			Error:        String(""),
			Category:     String(".github/workflows/shiftleft.yml:build/"),
			CreatedAt:    date,
			ResultsCount: Int(17),
			RulesCount:   Int(32),
			URL:          String("https://api.github.com/repos/o/r/code-scanning/analyses/200"),
			SarifID:      String("8981cd8e-b078-4ac3-a3be-1dad7dbd0b582"),
			Tool: &Tool{
				Name:    String("Python Security ScanningAnalysis"),
				GUID:    nil,
				Version: String("1.2.0"),
			},
			Deletable: Bool(true),
			Warning:   String(""),
		},
	}
	if !cmp.Equal(analyses, want) {
		t.Errorf("CodeScanning.ListAnalysesForRepo returned %+v, want %+v", analyses, want)
	}

	const methodName = "ListAnalysesForRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListAnalysesForRepo(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListAnalysesForRepo(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_GetAnalysis(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/analyses/3602840", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			  "ref": "refs/heads/main",
			  "commit_sha": "c18c69115654ff0166991962832dc2bd7756e655",
			  "analysis_key": ".github/workflows/codeql-analysis.yml:analyze",
			  "environment": "{\"language\":\"javascript\"}",
			  "error": "",
			  "category": ".github/workflows/codeql-analysis.yml:analyze/language:javascript",
			  "created_at": "2021-01-13T11:55:49Z",
			  "results_count": 3,
			  "rules_count": 67,
			  "id": 3602840,
			  "url": "https://api.github.com/repos/o/r/code-scanning/analyses/201",
			  "sarif_id": "47177e22-5596-11eb-80a1-c1e54ef945c6",
			  "tool": {
				"name": "CodeQL",
				"guid": null,
				"version": "2.4.0"
			  },
			  "deletable": true,
			  "warning": ""
			}`)
	})

	ctx := context.Background()
	analysis, _, err := client.CodeScanning.GetAnalysis(ctx, "o", "r", 3602840)
	if err != nil {
		t.Errorf("CodeScanning.GetAnalysis returned error: %v", err)
	}

	date := &Timestamp{time.Date(2021, time.January, 13, 11, 55, 49, 0, time.UTC)}
	want := &ScanningAnalysis{
		ID:           Int64(3602840),
		Ref:          String("refs/heads/main"),
		CommitSHA:    String("c18c69115654ff0166991962832dc2bd7756e655"),
		AnalysisKey:  String(".github/workflows/codeql-analysis.yml:analyze"),
		Environment:  String("{\"language\":\"javascript\"}"),
		Error:        String(""),
		Category:     String(".github/workflows/codeql-analysis.yml:analyze/language:javascript"),
		CreatedAt:    date,
		ResultsCount: Int(3),
		RulesCount:   Int(67),
		URL:          String("https://api.github.com/repos/o/r/code-scanning/analyses/201"),
		SarifID:      String("47177e22-5596-11eb-80a1-c1e54ef945c6"),
		Tool: &Tool{
			Name:    String("CodeQL"),
			GUID:    nil,
			Version: String("2.4.0"),
		},
		Deletable: Bool(true),
		Warning:   String(""),
	}
	if !cmp.Equal(analysis, want) {
		t.Errorf("CodeScanning.GetAnalysis returned %+v, want %+v", analysis, want)
	}

	const methodName = "GetAnalysis"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.GetAnalysis(ctx, "\n", "\n", -123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.GetAnalysis(ctx, "o", "r", 3602840)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
