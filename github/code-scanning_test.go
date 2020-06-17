// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
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
	if !reflect.DeepEqual(id, want) {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}

	// Test: HTMLURL is nil
	a = &Alert{}
	id = a.ID()
	want = 0
	if !reflect.DeepEqual(id, want) {
		t.Errorf("Alert.ID error returned %+v, want %+v", id, want)
	}

	// Test: ID can't be parsed as an int
	a = &Alert{
		HTMLURL: String("https://github.com/o/r/security/code-scanning/bad88"),
	}
	id = a.ID()
	want = 0
	if !reflect.DeepEqual(id, want) {
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
				"tool":"CodeQL",
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
				"tool":"CodeQL",
				"created_at":"2020-05-06T12:00:00Z",
				"open":true,
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/88",
				"html_url":"https://github.com/o/r/security/code-scanning/88"
				}]`)
	})

	opts := &AlertListOptions{State: "open", Ref: "heads/master"}
	alerts, _, err := client.CodeScanning.ListAlertsForRepo(context.Background(), "o", "r", opts)
	if err != nil {
		t.Errorf("CodeScanning.ListAlertsForRepo returned error: %v", err)
	}

	date := Timestamp{time.Date(2020, time.May, 06, 12, 00, 00, 0, time.UTC)}
	want := []*Alert{
		{
			RuleID:          String("js/trivial-conditional"),
			RuleSeverity:    String("warning"),
			RuleDescription: String("Useless conditional"),
			Tool:            String("CodeQL"),
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
			Tool:            String("CodeQL"),
			CreatedAt:       &date,
			Open:            Bool(true),
			ClosedBy:        nil,
			ClosedAt:        nil,
			URL:             String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
			HTMLURL:         String("https://github.com/o/r/security/code-scanning/88"),
		},
	}
	if !reflect.DeepEqual(alerts, want) {
		t.Errorf("CodeScanning.ListAlertsForRepo returned %+v, want %+v", alerts, want)
	}
}

func TestActionsService_GetAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/code-scanning/alerts/88", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"rule_id":"js/useless-expression",
				"rule_severity":"warning",
				"rule_description":"Expression has no effect",
				"tool":"CodeQL",
				"created_at":"2019-01-02T15:04:05Z",
				"open":true,
				"closed_by":null,
				"closed_at":null,
				"url":"https://api.github.com/repos/o/r/code-scanning/alerts/88",
				"html_url":"https://github.com/o/r/security/code-scanning/88"}`)
	})

	alert, _, err := client.CodeScanning.GetAlert(context.Background(), "o", "r", 88)
	if err != nil {
		t.Errorf("CodeScanning.GetAlert returned error: %v", err)
	}

	date := Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}
	want := &Alert{
		RuleID:          String("js/useless-expression"),
		RuleSeverity:    String("warning"),
		RuleDescription: String("Expression has no effect"),
		Tool:            String("CodeQL"),
		CreatedAt:       &date,
		Open:            Bool(true),
		ClosedBy:        nil,
		ClosedAt:        nil,
		URL:             String("https://api.github.com/repos/o/r/code-scanning/alerts/88"),
		HTMLURL:         String("https://github.com/o/r/security/code-scanning/88"),
	}
	if !reflect.DeepEqual(alert, want) {
		t.Errorf("CodeScanning.GetAlert returned %+v, want %+v", alert, want)
	}
}
