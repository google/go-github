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
				"created_at":"2020-05-06T12:00:00Z",
				"open":true,
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
				"created_at":"2020-05-06T12:00:00Z",
				"open":true,
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
			CreatedAt:       &date,
			Open:            Bool(true),
			ClosedBy:        nil,
			ClosedAt:        nil,
			URL:             String("https://api.github.com/repos/o/r/code-scanning/alerts/25"),
			HTMLURL:         String("https://github.com/o/r/security/code-scanning/25"),
		},
		{
			RuleID:          String("js/useless-expression"),
			RuleSeverity:    String("warning"),
			RuleDescription: String("Expression has no effect"),
			Tool:            &Tool{Name: String("CodeQL"), GUID: nil, Version: String("1.4.0")},
			CreatedAt:       &date,
			Open:            Bool(true),
			ClosedBy:        nil,
			ClosedAt:        nil,
			URL:             String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
			HTMLURL:         String("https://github.com/o/r/security/code-scanning/88"),
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
				"created_at":"2019-01-02T15:04:05Z",
				"open":true,
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
		CreatedAt:       &date,
		Open:            Bool(true),
		ClosedBy:        nil,
		ClosedAt:        nil,
		URL:             String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
		HTMLURL:         String("https://github.com/o/r/security/code-scanning/88"),
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
		Open:      Bool(false),
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
		"open": false,
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
