// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_Alert_ID(t *testing.T) {
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

func TestActionsService_ListAlertsForRepo(t *testing.T) {
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

func TestActionsService_GetAlert(t *testing.T) {
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
	testJSONMarshal(t, &Tool{}, "{}")

	u := &Message{
		Text: String("text"),
	}

	want := `{
		"text": "text"
	}`

	testJSONMarshal(t, u, want)
}
