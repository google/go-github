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
	t.Parallel()
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
	t.Parallel()
	client, mux, _ := setup(t)

	expectedSarifID := &SarifID{
		ID:  String("testid"),
		URL: String("https://example.com/testurl"),
	}

	mux.HandleFunc("/repos/o/r/code-scanning/sarifs", func(w http.ResponseWriter, r *http.Request) {
		v := new(SarifAnalysis)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))
		testMethod(t, r, "POST")
		want := &SarifAnalysis{CommitSHA: String("abc"), Ref: String("ref/head/main"), Sarif: String("abc"), CheckoutURI: String("uri"), StartedAt: &Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)}, ToolName: String("codeql-cli")}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		w.WriteHeader(http.StatusAccepted)
		respBody, _ := json.Marshal(expectedSarifID)
		_, _ = w.Write(respBody)
	})

	ctx := context.Background()
	sarifAnalysis := &SarifAnalysis{CommitSHA: String("abc"), Ref: String("ref/head/main"), Sarif: String("abc"), CheckoutURI: String("uri"), StartedAt: &Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)}, ToolName: String("codeql-cli")}
	respSarifID, _, err := client.CodeScanning.UploadSarif(ctx, "o", "r", sarifAnalysis)
	if err != nil {
		t.Errorf("CodeScanning.UploadSarif returned error: %v", err)
	}
	if !cmp.Equal(expectedSarifID, respSarifID) {
		t.Errorf("Sarif response = %+v, want %+v", respSarifID, expectedSarifID)
	}

	const methodName = "UploadSarif"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.UploadSarif(ctx, "\n", "\n", sarifAnalysis)
		return err
	})

	testNewRequestAndDoFailureCategory(t, methodName, client, CodeScanningUploadCategory, func() (*Response, error) {
		_, resp, err := client.CodeScanning.UploadSarif(ctx, "o", "r", sarifAnalysis)
		return resp, err
	})
}

func TestCodeScanningService_GetSARIF(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/sarifs/abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"processing_status": "s",
			"analyses_url": "u"
		}`)
	})

	ctx := context.Background()
	sarifUpload, _, err := client.CodeScanning.GetSARIF(ctx, "o", "r", "abc")
	if err != nil {
		t.Errorf("CodeScanning.GetSARIF returned error: %v", err)
	}

	want := &SARIFUpload{
		ProcessingStatus: String("s"),
		AnalysesURL:      String("u"),
	}
	if !cmp.Equal(sarifUpload, want) {
		t.Errorf("CodeScanning.GetSARIF returned %+v, want %+v", sarifUpload, want)
	}

	const methodName = "GetSARIF"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.GetSARIF(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.GetSARIF(ctx, "o", "r", "abc")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_ListAlertsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "ref": "heads/master", "severity": "warning", "tool_name": "CodeQL"})
		fmt.Fprint(w, `[{
				"repository": {
					"id": 1,
					"name": "n",
					"url": "url"
				},
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

	opts := &AlertListOptions{State: "open", Ref: "heads/master", Severity: "warning", ToolName: "CodeQL"}
	ctx := context.Background()
	alerts, _, err := client.CodeScanning.ListAlertsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAlertsForOrg returned error: %v", err)
	}

	date := Timestamp{time.Date(2020, time.May, 06, 12, 00, 00, 0, time.UTC)}
	want := []*Alert{
		{
			Repository: &Repository{
				ID:   Int64(1),
				URL:  String("url"),
				Name: String("n"),
			},
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
		t.Errorf("CodeScanning.ListAlertsForOrg returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListAlertsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListAlertsForOrg(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_ListAlertsForOrgLisCursorOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/code-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "ref": "heads/master", "severity": "warning", "tool_name": "CodeQL", "per_page": "1", "before": "deadbeefb", "after": "deadbeefa"})
		fmt.Fprint(w, `[{
				"repository": {
					"id": 1,
					"name": "n",
					"url": "url"
				},
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
				}]`)
	})

	opts := &AlertListOptions{State: "open", Ref: "heads/master", Severity: "warning", ToolName: "CodeQL", ListCursorOptions: ListCursorOptions{PerPage: 1, Before: "deadbeefb", After: "deadbeefa"}}
	ctx := context.Background()
	alerts, _, err := client.CodeScanning.ListAlertsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAlertsForOrg returned error: %v", err)
	}

	date := Timestamp{time.Date(2020, time.May, 06, 12, 00, 00, 0, time.UTC)}
	want := []*Alert{
		{
			Repository: &Repository{
				ID:   Int64(1),
				URL:  String("url"),
				Name: String("n"),
			},
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
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("CodeScanning.ListAlertsForOrg returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListAlertsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListAlertsForOrg(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_ListAlertsForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "ref": "heads/master", "severity": "warning", "tool_name": "CodeQL"})
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

	opts := &AlertListOptions{State: "open", Ref: "heads/master", Severity: "warning", ToolName: "CodeQL"}
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

func TestCodeScanningService_UpdateAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/alerts/88", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
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
					"state": "dismissed",
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
				"state":"dismissed",
				"dismissed_reason": "false positive",
				"dismissed_comment": "This alert is not actually correct as sanitizer is used",
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/88",
				"html_url":"https://github.com/o/r/security/code-scanning/88"}`)
	})

	ctx := context.Background()
	dismissedComment := String("This alert is not actually correct as sanitizer is used")
	dismissedReason := String("false positive")
	state := String("dismissed")
	stateInfo := &CodeScanningAlertState{State: *state, DismissedReason: dismissedReason, DismissedComment: dismissedComment}
	alert, _, err := client.CodeScanning.UpdateAlert(ctx, "o", "r", 88, stateInfo)
	if err != nil {
		t.Errorf("CodeScanning.UpdateAlert returned error: %v", err)
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
		CreatedAt:        &date,
		State:            state,
		DismissedReason:  dismissedReason,
		DismissedComment: dismissedComment,
		ClosedBy:         nil,
		ClosedAt:         nil,
		URL:              String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
		HTMLURL:          String("https://github.com/o/r/security/code-scanning/88"),
		MostRecentInstance: &MostRecentInstance{
			Ref:       String("refs/heads/main"),
			State:     String("dismissed"),
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
		t.Errorf("CodeScanning.UpdateAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "UpdateAlert"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.UpdateAlert(ctx, "\n", "\n", -88, stateInfo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.UpdateAlert(ctx, "o", "r", 88, stateInfo)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_ListAlertInstances(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/alerts/88/instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "ref": "refs/heads/main",
			  "analysis_key": ".github/workflows/codeql-analysis.yml:analyze",
			  "environment": "",
			  "category": ".github/workflows/codeql-analysis.yml:analyze",
			  "state": "open",
			  "fixed_at": null,
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
			}
		  ]`)
	})

	opts := &AlertInstancesListOptions{Ref: "heads/main", ListOptions: ListOptions{Page: 1}}
	ctx := context.Background()
	instances, _, err := client.CodeScanning.ListAlertInstances(ctx, "o", "r", 88, opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAlertInstances returned error: %v", err)
	}

	want := []*MostRecentInstance{
		{
			Ref:         String("refs/heads/main"),
			AnalysisKey: String(".github/workflows/codeql-analysis.yml:analyze"),
			Category:    String(".github/workflows/codeql-analysis.yml:analyze"),
			Environment: String(""),
			State:       String("open"),
			CommitSHA:   String("abcdefg12345"),
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
	if !cmp.Equal(instances, want) {
		t.Errorf("CodeScanning.ListAlertInstances returned %+v, want %+v", instances, want)
	}

	const methodName = "ListAlertInstances"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListAlertInstances(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListAlertInstances(ctx, "o", "r", 88, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_GetAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/alerts/88", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"rule_id":"js/useless-expression",
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
			"html_url":"https://github.com/o/r/security/code-scanning/88"
		}`)
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	t.Parallel()
	client, mux, _ := setup(t)

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

func TestCodeScanningService_DeleteAnalysis(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/analyses/40", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"next_analysis_url": "a",
			"confirm_delete_url": "b"
		}`)
	})

	ctx := context.Background()
	analysis, _, err := client.CodeScanning.DeleteAnalysis(ctx, "o", "r", 40)
	if err != nil {
		t.Errorf("CodeScanning.DeleteAnalysis returned error: %v", err)
	}

	want := &DeleteAnalysis{
		NextAnalysisURL:  String("a"),
		ConfirmDeleteURL: String("b"),
	}
	if !cmp.Equal(analysis, want) {
		t.Errorf("CodeScanning.DeleteAnalysis returned %+v, want %+v", analysis, want)
	}

	const methodName = "DeleteAnalysis"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.DeleteAnalysis(ctx, "\n", "\n", -123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.DeleteAnalysis(ctx, "o", "r", 40)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_ListCodeQLDatabases(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/codeql/databases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"id": 1,
				"name": "name",
				"language": "language",
				"uploader": {
					"login": "a",
					"id": 1,
					"node_id": "b",
					"avatar_url": "c",
					"gravatar_id": "d",
					"url": "e",
					"html_url": "f",
					"followers_url": "g",
					"following_url": "h",
					"gists_url": "i",
					"starred_url": "j",
					"subscriptions_url": "k",
					"organizations_url": "l",
					"repos_url": "m",
					"events_url": "n",
					"received_events_url": "o",
					"type": "p",
					"site_admin": false
				},
				"content_type": "r",
				"size": 1024,
				"created_at": "2021-01-13T11:55:49Z",
				"updated_at": "2021-01-13T11:55:49Z",
				"url": "s"
			}
		]`)
	})

	ctx := context.Background()
	databases, _, err := client.CodeScanning.ListCodeQLDatabases(ctx, "o", "r")
	if err != nil {
		t.Errorf("CodeScanning.ListCodeQLDatabases returned error: %v", err)
	}

	date := &Timestamp{time.Date(2021, time.January, 13, 11, 55, 49, 0, time.UTC)}
	want := []*CodeQLDatabase{
		{
			ID:       Int64(1),
			Name:     String("name"),
			Language: String("language"),
			Uploader: &User{
				Login:             String("a"),
				ID:                Int64(1),
				NodeID:            String("b"),
				AvatarURL:         String("c"),
				GravatarID:        String("d"),
				URL:               String("e"),
				HTMLURL:           String("f"),
				FollowersURL:      String("g"),
				FollowingURL:      String("h"),
				GistsURL:          String("i"),
				StarredURL:        String("j"),
				SubscriptionsURL:  String("k"),
				OrganizationsURL:  String("l"),
				ReposURL:          String("m"),
				EventsURL:         String("n"),
				ReceivedEventsURL: String("o"),
				Type:              String("p"),
				SiteAdmin:         Bool(false),
			},
			ContentType: String("r"),
			Size:        Int64(1024),
			CreatedAt:   date,
			UpdatedAt:   date,
			URL:         String("s"),
		},
	}

	if !cmp.Equal(databases, want) {
		t.Errorf("CodeScanning.ListCodeQLDatabases returned %+v, want %+v", databases, want)
	}

	const methodName = "ListCodeQLDatabases"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.ListCodeQLDatabases(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.ListCodeQLDatabases(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_GetCodeQLDatabase(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/codeql/databases/lang", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 1,
			"name": "name",
			"language": "language",
			"uploader": {
				"login": "a",
				"id": 1,
				"node_id": "b",
				"avatar_url": "c",
				"gravatar_id": "d",
				"url": "e",
				"html_url": "f",
				"followers_url": "g",
				"following_url": "h",
				"gists_url": "i",
				"starred_url": "j",
				"subscriptions_url": "k",
				"organizations_url": "l",
				"repos_url": "m",
				"events_url": "n",
				"received_events_url": "o",
				"type": "p",
				"site_admin": false
			},
			"content_type": "r",
			"size": 1024,
			"created_at": "2021-01-13T11:55:49Z",
			"updated_at": "2021-01-13T11:55:49Z",
			"url": "s"
		}`)
	})

	ctx := context.Background()
	database, _, err := client.CodeScanning.GetCodeQLDatabase(ctx, "o", "r", "lang")
	if err != nil {
		t.Errorf("CodeScanning.GetCodeQLDatabase returned error: %v", err)
	}

	date := &Timestamp{time.Date(2021, time.January, 13, 11, 55, 49, 0, time.UTC)}
	want := &CodeQLDatabase{
		ID:       Int64(1),
		Name:     String("name"),
		Language: String("language"),
		Uploader: &User{
			Login:             String("a"),
			ID:                Int64(1),
			NodeID:            String("b"),
			AvatarURL:         String("c"),
			GravatarID:        String("d"),
			URL:               String("e"),
			HTMLURL:           String("f"),
			FollowersURL:      String("g"),
			FollowingURL:      String("h"),
			GistsURL:          String("i"),
			StarredURL:        String("j"),
			SubscriptionsURL:  String("k"),
			OrganizationsURL:  String("l"),
			ReposURL:          String("m"),
			EventsURL:         String("n"),
			ReceivedEventsURL: String("o"),
			Type:              String("p"),
			SiteAdmin:         Bool(false),
		},
		ContentType: String("r"),
		Size:        Int64(1024),
		CreatedAt:   date,
		UpdatedAt:   date,
		URL:         String("s"),
	}

	if !cmp.Equal(database, want) {
		t.Errorf("CodeScanning.GetCodeQLDatabase returned %+v, want %+v", database, want)
	}

	const methodName = "GetCodeQLDatabase"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.GetCodeQLDatabase(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.GetCodeQLDatabase(ctx, "o", "r", "lang")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_GetDefaultSetupConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/default-setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, err := fmt.Fprint(w, `{
		"state": "configured",
		"languages": [
			"javascript",
			"javascript-typescript",
			"typescript"
		],
		"query_suite": "default",
		"updated_at": "2006-01-02T15:04:05Z"
		}`)
		if err != nil {
			t.Fatal(err)
		}
	})

	ctx := context.Background()
	cfg, _, err := client.CodeScanning.GetDefaultSetupConfiguration(ctx, "o", "r")
	if err != nil {
		t.Errorf("CodeScanning.GetDefaultSetupConfiguration returned error: %v", err)
	}

	date := &Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)}
	want := &DefaultSetupConfiguration{
		State:      String("configured"),
		Languages:  []string{"javascript", "javascript-typescript", "typescript"},
		QuerySuite: String("default"),
		UpdatedAt:  date,
	}
	if !cmp.Equal(cfg, want) {
		t.Errorf("CodeScanning.GetDefaultSetupConfiguration returned %+v, want %+v", cfg, want)
	}

	const methodName = "GetDefaultSetupConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.GetDefaultSetupConfiguration(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.GetDefaultSetupConfiguration(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeScanningService_UpdateDefaultSetupConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-scanning/default-setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		_, err := fmt.Fprint(w, `{
		"run_id": 5301214200,
		"run_url": "https://api.github.com/repos/o/r/actions/runs/5301214200"
		}`)
		if err != nil {
			t.Fatal(err)
		}
	})

	ctx := context.Background()
	options := &UpdateDefaultSetupConfigurationOptions{
		State:      "configured",
		Languages:  []string{"go"},
		QuerySuite: String("default"),
	}
	got, _, err := client.CodeScanning.UpdateDefaultSetupConfiguration(ctx, "o", "r", options)
	if err != nil {
		t.Errorf("CodeScanning.UpdateDefaultSetupConfiguration returned error: %v", err)
	}

	want := &UpdateDefaultSetupConfigurationResponse{
		RunID:  Int64(5301214200),
		RunURL: String("https://api.github.com/repos/o/r/actions/runs/5301214200"),
	}
	if !cmp.Equal(got, want) {
		t.Errorf("CodeScanning.UpdateDefaultSetupConfiguration returned %+v, want %+v", got, want)
	}

	const methodName = "UpdateDefaultSetupConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeScanning.UpdateDefaultSetupConfiguration(ctx, "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeScanning.UpdateDefaultSetupConfiguration(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
