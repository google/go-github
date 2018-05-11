// Copyright 2018 The go-github AUTHORS. All rights reserved.
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

func TestChecksService_GetCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-runs/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		fmt.Fprint(w, `{
			"id": 1,
                        "name":"testCheckRun",
			"status": "completed",
			"conclusion": "neutral",
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": "2018-05-04T01:14:52Z"}`)
	})
	checkRun, _, err := client.Checks.GetCheckRun(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.GetCheckRun return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	completeAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")

	want := &CheckRun{
		ID:          Int64(1),
		Status:      String("completed"),
		Conclusion:  String("neutral"),
		StartedAt:   &Timestamp{startedAt},
		CompletedAt: &Timestamp{completeAt},
		Name:        String("testCheckRun"),
	}
	if !reflect.DeepEqual(checkRun, want) {
		t.Errorf("Checks.GetCheckRun return %+v, want %+v", checkRun, want)
	}
}

func TestChecksService_CreateCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		fmt.Fprint(w, `{
			"id": 1,
                        "name":"testCreateCheckRun",
                        "head_sha":"deadbeef",
			"status": "in_progress",
			"conclusion": null,
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": null,
                        "output":{"title": "Mighty test report", "summary":"", "text":""}}`)
	})
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	checkRunOpt := CreateCheckRunOptions{
		HeadBranch: "master",
		Name:       "testCreateCheckRun",
		HeadSHA:    "deadbeef",
		Status:     String("in_progress"),
		StartedAt:  &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String(""),
			Text:    String(""),
		},
	}

	checkRun, _, err := client.Checks.CreateCheckRun(context.Background(), "o", "r", checkRunOpt)
	if err != nil {
		t.Errorf("Checks.CreateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:        Int64(1),
		Status:    String("in_progress"),
		StartedAt: &Timestamp{startedAt},
		HeadSHA:   String("deadbeef"),
		Name:      String("testCreateCheckRun"),
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String(""),
			Text:    String(""),
		},
	}
	if !reflect.DeepEqual(checkRun, want) {
		t.Errorf("Checks.CreateCheckRun return %+v, want %+v", checkRun, want)
	}
}

func TestChecksService_ListCheckRunAnnotations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-runs/1/annotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testFormValues(t, r, values{
			"page": "1",
		})
		fmt.Fprint(w, `[{
		                           "filename": "README.md",
		                           "blob_href": "https://github.com/octocat/Hello-World/blob/837db83be4137ca555d9a5598d0a1ea2987ecfee/README.md",
		                           "start_line": 2,
		                           "end_line": 2,
		                           "warning_level": "warning",
		                           "message": "Check your spelling for 'banaas'.",
                                           "title": "Spell check",
		                           "raw_details": "Do you mean 'bananas' or 'banana'?"}]`,
		)
	})

	checkRunAnnotations, _, err := client.Checks.ListCheckRunAnnotations(context.Background(), "o", "r", 1, &ListOptions{Page: 1})
	if err != nil {
		t.Errorf("Checks.ListCheckRunAnnotations return error: %v", err)
	}

	want := []*CheckAnnotation{{
		FileName:     String("README.md"),
		BlobHRef:     String("https://github.com/octocat/Hello-World/blob/837db83be4137ca555d9a5598d0a1ea2987ecfee/README.md"),
		StartLine:    Int(2),
		EndLine:      Int(2),
		WarningLevel: String("warning"),
		Message:      String("Check your spelling for 'banaas'."),
		RawDetails:   String("Do you mean 'bananas' or 'banana'?"),
		Title:        String("Spell check"),
	}}

	if !reflect.DeepEqual(checkRunAnnotations, want) {
		t.Errorf("Checks.ListCheckRunAnnotations returned %+v, want %+v", checkRunAnnotations, want)
	}
}
