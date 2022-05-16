// Copyright 2018 The go-github AUTHORS. All rights reserved.
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
	ctx := context.Background()
	checkRun, _, err := client.Checks.GetCheckRun(ctx, "o", "r", 1)
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
	if !cmp.Equal(checkRun, want) {
		t.Errorf("Checks.GetCheckRun return %+v, want %+v", checkRun, want)
	}

	const methodName = "GetCheckRun"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.GetCheckRun(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.GetCheckRun(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_GetCheckSuite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		fmt.Fprint(w, `{
			"id": 1,
                        "head_branch":"master",
			"head_sha": "deadbeef",
			"conclusion": "neutral",
                        "before": "deadbeefb",
                        "after": "deadbeefa",
			"status": "completed"}`)
	})
	ctx := context.Background()
	checkSuite, _, err := client.Checks.GetCheckSuite(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.GetCheckSuite return error: %v", err)
	}
	want := &CheckSuite{
		ID:         Int64(1),
		HeadBranch: String("master"),
		HeadSHA:    String("deadbeef"),
		AfterSHA:   String("deadbeefa"),
		BeforeSHA:  String("deadbeefb"),
		Status:     String("completed"),
		Conclusion: String("neutral"),
	}
	if !cmp.Equal(checkSuite, want) {
		t.Errorf("Checks.GetCheckSuite return %+v, want %+v", checkSuite, want)
	}

	const methodName = "GetCheckSuite"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.GetCheckSuite(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.GetCheckSuite(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		Name:      "testCreateCheckRun",
		HeadSHA:   "deadbeef",
		Status:    String("in_progress"),
		StartedAt: &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String(""),
			Text:    String(""),
		},
	}

	ctx := context.Background()
	checkRun, _, err := client.Checks.CreateCheckRun(ctx, "o", "r", checkRunOpt)
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
	if !cmp.Equal(checkRun, want) {
		t.Errorf("Checks.CreateCheckRun return %+v, want %+v", checkRun, want)
	}

	const methodName = "CreateCheckRun"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.CreateCheckRun(ctx, "\n", "\n", CreateCheckRunOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.CreateCheckRun(ctx, "o", "r", checkRunOpt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		                           "path": "README.md",
		                           "start_line": 2,
		                           "end_line": 2,
		                           "start_column": 1,
		                           "end_column": 5,									
		                           "annotation_level": "warning",
		                           "message": "Check your spelling for 'banaas'.",
                                           "title": "Spell check",
		                           "raw_details": "Do you mean 'bananas' or 'banana'?"}]`,
		)
	})

	ctx := context.Background()
	checkRunAnnotations, _, err := client.Checks.ListCheckRunAnnotations(ctx, "o", "r", 1, &ListOptions{Page: 1})
	if err != nil {
		t.Errorf("Checks.ListCheckRunAnnotations return error: %v", err)
	}

	want := []*CheckRunAnnotation{{
		Path:            String("README.md"),
		StartLine:       Int(2),
		EndLine:         Int(2),
		StartColumn:     Int(1),
		EndColumn:       Int(5),
		AnnotationLevel: String("warning"),
		Message:         String("Check your spelling for 'banaas'."),
		Title:           String("Spell check"),
		RawDetails:      String("Do you mean 'bananas' or 'banana'?"),
	}}

	if !cmp.Equal(checkRunAnnotations, want) {
		t.Errorf("Checks.ListCheckRunAnnotations returned %+v, want %+v", checkRunAnnotations, want)
	}

	const methodName = "ListCheckRunAnnotations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.ListCheckRunAnnotations(ctx, "\n", "\n", -1, &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.ListCheckRunAnnotations(ctx, "o", "r", 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_UpdateCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-runs/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		fmt.Fprint(w, `{
			"id": 1,
                        "name":"testUpdateCheckRun",
			"status": "completed",            
			"conclusion": "neutral",
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": "2018-05-04T01:14:52Z",
                        "output":{"title": "Mighty test report", "summary":"There are 0 failures, 2 warnings and 1 notice", "text":"You may have misspelled some words."}}`)
	})
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	updateCheckRunOpt := UpdateCheckRunOptions{
		Name:        "testUpdateCheckRun",
		Status:      String("completed"),
		CompletedAt: &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String("There are 0 failures, 2 warnings and 1 notice"),
			Text:    String("You may have misspelled some words."),
		},
	}

	ctx := context.Background()
	checkRun, _, err := client.Checks.UpdateCheckRun(ctx, "o", "r", 1, updateCheckRunOpt)
	if err != nil {
		t.Errorf("Checks.UpdateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:          Int64(1),
		Status:      String("completed"),
		StartedAt:   &Timestamp{startedAt},
		CompletedAt: &Timestamp{startedAt},
		Conclusion:  String("neutral"),
		Name:        String("testUpdateCheckRun"),
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String("There are 0 failures, 2 warnings and 1 notice"),
			Text:    String("You may have misspelled some words."),
		},
	}
	if !cmp.Equal(checkRun, want) {
		t.Errorf("Checks.UpdateCheckRun return %+v, want %+v", checkRun, want)
	}

	const methodName = "UpdateCheckRun"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.UpdateCheckRun(ctx, "\n", "\n", -1, UpdateCheckRunOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.UpdateCheckRun(ctx, "o", "r", 1, updateCheckRunOpt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_ListCheckRunsForRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/commits/master/check-runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testFormValues(t, r, values{
			"check_name": "testing",
			"page":       "1",
			"status":     "completed",
			"filter":     "all",
			"app_id":     "1",
		})
		fmt.Fprint(w, `{"total_count":1,
                                "check_runs": [{
                                    "id": 1,
                                    "head_sha": "deadbeef",
                                    "status": "completed",
                                    "conclusion": "neutral",
                                    "started_at": "2018-05-04T01:14:52Z",
                                    "completed_at": "2018-05-04T01:14:52Z",
                                    "app": {
                                      "id": 1}}]}`,
		)
	})

	opt := &ListCheckRunsOptions{
		CheckName:   String("testing"),
		Status:      String("completed"),
		Filter:      String("all"),
		AppID:       Int64(1),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := context.Background()
	checkRuns, _, err := client.Checks.ListCheckRunsForRef(ctx, "o", "r", "master", opt)
	if err != nil {
		t.Errorf("Checks.ListCheckRunsForRef return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	want := &ListCheckRunsResults{
		Total: Int(1),
		CheckRuns: []*CheckRun{{
			ID:          Int64(1),
			Status:      String("completed"),
			StartedAt:   &Timestamp{startedAt},
			CompletedAt: &Timestamp{startedAt},
			Conclusion:  String("neutral"),
			HeadSHA:     String("deadbeef"),
			App:         &App{ID: Int64(1)},
		}},
	}

	if !cmp.Equal(checkRuns, want) {
		t.Errorf("Checks.ListCheckRunsForRef returned %+v, want %+v", checkRuns, want)
	}

	const methodName = "ListCheckRunsForRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.ListCheckRunsForRef(ctx, "\n", "\n", "\n", &ListCheckRunsOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.ListCheckRunsForRef(ctx, "o", "r", "master", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_ListCheckRunsCheckSuite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites/1/check-runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testFormValues(t, r, values{
			"check_name": "testing",
			"page":       "1",
			"status":     "completed",
			"filter":     "all",
		})
		fmt.Fprint(w, `{"total_count":1,
                                "check_runs": [{
                                    "id": 1,
                                    "head_sha": "deadbeef",
                                    "status": "completed",
                                    "conclusion": "neutral",
                                    "started_at": "2018-05-04T01:14:52Z",
                                    "completed_at": "2018-05-04T01:14:52Z"}]}`,
		)
	})

	opt := &ListCheckRunsOptions{
		CheckName:   String("testing"),
		Status:      String("completed"),
		Filter:      String("all"),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := context.Background()
	checkRuns, _, err := client.Checks.ListCheckRunsCheckSuite(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Checks.ListCheckRunsCheckSuite return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	want := &ListCheckRunsResults{
		Total: Int(1),
		CheckRuns: []*CheckRun{{
			ID:          Int64(1),
			Status:      String("completed"),
			StartedAt:   &Timestamp{startedAt},
			CompletedAt: &Timestamp{startedAt},
			Conclusion:  String("neutral"),
			HeadSHA:     String("deadbeef"),
		}},
	}

	if !cmp.Equal(checkRuns, want) {
		t.Errorf("Checks.ListCheckRunsCheckSuite returned %+v, want %+v", checkRuns, want)
	}

	const methodName = "ListCheckRunsCheckSuite"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.ListCheckRunsCheckSuite(ctx, "\n", "\n", -1, &ListCheckRunsOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.ListCheckRunsCheckSuite(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_ListCheckSuiteForRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/commits/master/check-suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testFormValues(t, r, values{
			"check_name": "testing",
			"page":       "1",
			"app_id":     "2",
		})
		fmt.Fprint(w, `{"total_count":1,
                                "check_suites": [{
                                    "id": 1,
                                    "head_sha": "deadbeef",
                                    "head_branch": "master",
                                    "status": "completed",
                                    "conclusion": "neutral",
                                    "before": "deadbeefb",
                                    "after": "deadbeefa"}]}`,
		)
	})

	opt := &ListCheckSuiteOptions{
		CheckName:   String("testing"),
		AppID:       Int(2),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := context.Background()
	checkSuites, _, err := client.Checks.ListCheckSuitesForRef(ctx, "o", "r", "master", opt)
	if err != nil {
		t.Errorf("Checks.ListCheckSuitesForRef return error: %v", err)
	}
	want := &ListCheckSuiteResults{
		Total: Int(1),
		CheckSuites: []*CheckSuite{{
			ID:         Int64(1),
			Status:     String("completed"),
			Conclusion: String("neutral"),
			HeadSHA:    String("deadbeef"),
			HeadBranch: String("master"),
			BeforeSHA:  String("deadbeefb"),
			AfterSHA:   String("deadbeefa"),
		}},
	}

	if !cmp.Equal(checkSuites, want) {
		t.Errorf("Checks.ListCheckSuitesForRef returned %+v, want %+v", checkSuites, want)
	}

	const methodName = "ListCheckSuitesForRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.ListCheckSuitesForRef(ctx, "\n", "\n", "\n", &ListCheckSuiteOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.ListCheckSuitesForRef(ctx, "o", "r", "master", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_SetCheckSuitePreferences(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites/preferences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testBody(t, r, `{"auto_trigger_checks":[{"app_id":2,"setting":false}]}`+"\n")
		fmt.Fprint(w, `{"preferences":{"auto_trigger_checks":[{"app_id": 2,"setting": false}]}}`)
	})
	a := []*AutoTriggerCheck{{
		AppID:   Int64(2),
		Setting: Bool(false),
	}}
	opt := CheckSuitePreferenceOptions{AutoTriggerChecks: a}
	ctx := context.Background()
	prefResults, _, err := client.Checks.SetCheckSuitePreferences(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Checks.SetCheckSuitePreferences return error: %v", err)
	}

	p := &PreferenceList{
		AutoTriggerChecks: a,
	}
	want := &CheckSuitePreferenceResults{
		Preferences: p,
	}

	if !cmp.Equal(prefResults, want) {
		t.Errorf("Checks.SetCheckSuitePreferences return %+v, want %+v", prefResults, want)
	}

	const methodName = "SetCheckSuitePreferences"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.SetCheckSuitePreferences(ctx, "\n", "\n", CheckSuitePreferenceOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.SetCheckSuitePreferences(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_CreateCheckSuite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		fmt.Fprint(w, `{
			"id": 2,
                        "head_branch":"master",
                        "head_sha":"deadbeef",
			"status": "completed",
			"conclusion": "neutral",
                        "before": "deadbeefb",
                        "after": "deadbeefa"}`)
	})

	checkSuiteOpt := CreateCheckSuiteOptions{
		HeadSHA:    "deadbeef",
		HeadBranch: String("master"),
	}

	ctx := context.Background()
	checkSuite, _, err := client.Checks.CreateCheckSuite(ctx, "o", "r", checkSuiteOpt)
	if err != nil {
		t.Errorf("Checks.CreateCheckSuite return error: %v", err)
	}

	want := &CheckSuite{
		ID:         Int64(2),
		Status:     String("completed"),
		HeadSHA:    String("deadbeef"),
		HeadBranch: String("master"),
		Conclusion: String("neutral"),
		BeforeSHA:  String("deadbeefb"),
		AfterSHA:   String("deadbeefa"),
	}
	if !cmp.Equal(checkSuite, want) {
		t.Errorf("Checks.CreateCheckSuite return %+v, want %+v", checkSuite, want)
	}

	const methodName = "CreateCheckSuite"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Checks.CreateCheckSuite(ctx, "\n", "\n", CreateCheckSuiteOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Checks.CreateCheckSuite(ctx, "o", "r", checkSuiteOpt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestChecksService_ReRequestCheckSuite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites/1/rerequest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		w.WriteHeader(http.StatusCreated)
	})
	ctx := context.Background()
	resp, err := client.Checks.ReRequestCheckSuite(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.ReRequestCheckSuite return error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("Checks.ReRequestCheckSuite = %v, want %v", got, want)
	}

	const methodName = "ReRequestCheckSuite"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Checks.ReRequestCheckSuite(ctx, "\n", "\n", 1)
		return err
	})
}

func Test_CheckRunMarshal(t *testing.T) {
	testJSONMarshal(t, &CheckRun{}, "{}")

	now := time.Now()
	ts := now.Format(time.RFC3339Nano)

	c := CheckRun{
		ID:          Int64(1),
		NodeID:      String("n"),
		HeadSHA:     String("h"),
		ExternalID:  String("1"),
		URL:         String("u"),
		HTMLURL:     String("u"),
		DetailsURL:  String("u"),
		Status:      String("s"),
		Conclusion:  String("c"),
		StartedAt:   &Timestamp{Time: now},
		CompletedAt: &Timestamp{Time: now},
		Output: &CheckRunOutput{
			Annotations: []*CheckRunAnnotation{
				{
					AnnotationLevel: String("a"),
					EndLine:         Int(1),
					Message:         String("m"),
					Path:            String("p"),
					RawDetails:      String("r"),
					StartLine:       Int(1),
					Title:           String("t"),
				},
			},
			AnnotationsCount: Int(1),
			AnnotationsURL:   String("a"),
			Images: []*CheckRunImage{
				{
					Alt:      String("a"),
					ImageURL: String("i"),
					Caption:  String("c"),
				},
			},
			Title:   String("t"),
			Summary: String("s"),
			Text:    String("t"),
		},
		Name: String("n"),
		CheckSuite: &CheckSuite{
			ID: Int64(1),
		},
		App: &App{
			ID:     Int64(1),
			NodeID: String("n"),
			Owner: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Name:        String("n"),
			Description: String("d"),
			HTMLURL:     String("h"),
			ExternalURL: String("u"),
			CreatedAt:   &Timestamp{now},
			UpdatedAt:   &Timestamp{now},
		},
		PullRequests: []*PullRequest{
			{
				URL:    String("u"),
				ID:     Int64(1),
				Number: Int(1),
				Head: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("s"),
						Name: String("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("u"),
						Name: String("n"),
					},
				},
			},
		},
	}
	w := fmt.Sprintf(`{
		"id": 1,
		"node_id": "n",
		"head_sha": "h",
		"external_id": "1",
		"url": "u",
		"html_url": "u",
		"details_url": "u",
		"status": "s",
		"conclusion": "c",
		"started_at": "%s",
		"completed_at": "%s",
		"output": {
			"title": "t",
			"summary": "s",
			"text": "t",
			"annotations_count": 1,
			"annotations_url": "a",
			"annotations": [
				{
					"path": "p",
					"start_line": 1,
					"end_line": 1,
					"annotation_level": "a",
					"message": "m",
					"title": "t",
					"raw_details": "r"
				}
			],
			"images": [
				{
					"alt": "a",
					"image_url": "i",
					"caption": "c"
				}
			]
		},
		"name": "n",
		"check_suite": {
			"id": 1
		},
		"app": {
			"id": 1,
			"node_id": "n",
			"owner": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"name": "n",
			"description": "d",
			"external_url": "u",
			"html_url": "h",
			"created_at": "%s",
			"updated_at": "%s"
		},
		"pull_requests": [
			{
				"id": 1,
				"number": 1,
				"url": "u",
				"head": {
					"ref": "r",
					"sha": "s",
					"repo": {
						"id": 1,
						"name": "n",
						"url": "s"
					}
				},
				"base": {
					"ref": "r",
					"sha": "s",
					"repo": {
						"id": 1,
						"name": "n",
						"url": "u"
					}
				}
			}
		]
	  }`, ts, ts, ts, ts)

	testJSONMarshal(t, &c, w)
}

func Test_CheckSuiteMarshal(t *testing.T) {
	testJSONMarshal(t, &CheckSuite{}, "{}")

	now := time.Now()
	ts := now.Format(time.RFC3339Nano)

	c := CheckSuite{
		ID:         Int64(1),
		NodeID:     String("n"),
		HeadBranch: String("h"),
		HeadSHA:    String("h"),
		URL:        String("u"),
		BeforeSHA:  String("b"),
		AfterSHA:   String("a"),
		Status:     String("s"),
		Conclusion: String("c"),
		App: &App{
			ID:     Int64(1),
			NodeID: String("n"),
			Owner: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Name:        String("n"),
			Description: String("d"),
			HTMLURL:     String("h"),
			ExternalURL: String("u"),
			CreatedAt:   &Timestamp{now},
			UpdatedAt:   &Timestamp{now},
		},
		Repository: &Repository{
			ID: Int64(1),
		},
		PullRequests: []*PullRequest{
			{
				URL:    String("u"),
				ID:     Int64(1),
				Number: Int(1),
				Head: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("s"),
						Name: String("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("u"),
						Name: String("n"),
					},
				},
			},
		},
		HeadCommit: &Commit{
			SHA: String("s"),
		},
	}

	w := fmt.Sprintf(`{
			"id": 1,
			"node_id": "n",
			"head_branch": "h",
			"head_sha": "h",
			"url": "u",
			"before": "b",
			"after": "a",
			"status": "s",
			"conclusion": "c",
			"app": {
				"id": 1,
				"node_id": "n",
				"owner": {
					"login": "l",
					"id": 1,
					"node_id": "n",
					"avatar_url": "a",
					"url": "u",
					"events_url": "e",
					"repos_url": "r"
				},
				"name": "n",
				"description": "d",
				"external_url": "u",
				"html_url": "h",
				"created_at": "%s",
				"updated_at": "%s"
			},
			"repository": {
				"id": 1
			},
			"pull_requests": [
			{
				"id": 1,
				"number": 1,
				"url": "u",
				"head": {
					"ref": "r",
					"sha": "s",
					"repo": {
						"id": 1,
						"name": "n",
						"url": "s"
					}
				},
				"base": {
					"ref": "r",
					"sha": "s",
					"repo": {
						"id": 1,
						"name": "n",
						"url": "u"
					}
				}
			}
		],
		"head_commit": {
			"sha": "s"
		}
		}`, ts, ts)

	testJSONMarshal(t, &c, w)
}

func TestCheckRunAnnotation_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckRunAnnotation{}, "{}")

	u := &CheckRunAnnotation{
		Path:            String("p"),
		StartLine:       Int(1),
		EndLine:         Int(1),
		StartColumn:     Int(1),
		EndColumn:       Int(1),
		AnnotationLevel: String("al"),
		Message:         String("m"),
		Title:           String("t"),
		RawDetails:      String("rd"),
	}

	want := `{
		"path": "p",
		"start_line": 1,
		"end_line": 1,
		"start_column": 1,
		"end_column": 1,
		"annotation_level": "al",
		"message": "m",
		"title": "t",
		"raw_details": "rd"
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckRunImage_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckRunImage{}, "{}")

	u := &CheckRunImage{
		Alt:      String("a"),
		ImageURL: String("i"),
		Caption:  String("c"),
	}

	want := `{
		"alt": "a",
		"image_url": "i",
		"caption": "c"
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckRunAction_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckRunAction{}, "{}")

	u := &CheckRunAction{
		Label:       "l",
		Description: "d",
		Identifier:  "i",
	}

	want := `{
		"label": "l",
		"description": "d",
		"identifier": "i"
	}`

	testJSONMarshal(t, u, want)
}

func TestAutoTriggerCheck_Marshal(t *testing.T) {
	testJSONMarshal(t, &AutoTriggerCheck{}, "{}")

	u := &AutoTriggerCheck{
		AppID:   Int64(1),
		Setting: Bool(false),
	}

	want := `{
		"app_id": 1,
		"setting": false
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateCheckSuiteOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &CreateCheckSuiteOptions{}, "{}")

	u := &CreateCheckSuiteOptions{
		HeadSHA:    "hsha",
		HeadBranch: String("hb"),
	}

	want := `{
		"head_sha": "hsha",
		"head_branch": "hb"
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckRunOutput_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckRunOutput{}, "{}")

	u := &CheckRunOutput{
		Title:            String("ti"),
		Summary:          String("s"),
		Text:             String("t"),
		AnnotationsCount: Int(1),
		AnnotationsURL:   String("au"),
		Annotations: []*CheckRunAnnotation{
			{
				Path:            String("p"),
				StartLine:       Int(1),
				EndLine:         Int(1),
				StartColumn:     Int(1),
				EndColumn:       Int(1),
				AnnotationLevel: String("al"),
				Message:         String("m"),
				Title:           String("t"),
				RawDetails:      String("rd"),
			},
		},
		Images: []*CheckRunImage{
			{
				Alt:      String("a"),
				ImageURL: String("i"),
				Caption:  String("c"),
			},
		},
	}

	want := `{
		"title": "ti",
		"summary": "s",
		"text": "t",
		"annotations_count": 1,
		"annotations_url": "au",
		"annotations": [
			{
				"path": "p",
				"start_line": 1,
				"end_line": 1,
				"start_column": 1,
				"end_column": 1,
				"annotation_level": "al",
				"message": "m",
				"title": "t",
				"raw_details": "rd"
			}
		],
		"images": [
			{
				"alt": "a",
				"image_url": "i",
				"caption": "c"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateCheckRunOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &CreateCheckRunOptions{}, "{}")

	u := &CreateCheckRunOptions{
		Name:        "n",
		HeadSHA:     "hsha",
		DetailsURL:  String("durl"),
		ExternalID:  String("eid"),
		Status:      String("s"),
		Conclusion:  String("c"),
		StartedAt:   &Timestamp{referenceTime},
		CompletedAt: &Timestamp{referenceTime},
		Output: &CheckRunOutput{
			Title:            String("ti"),
			Summary:          String("s"),
			Text:             String("t"),
			AnnotationsCount: Int(1),
			AnnotationsURL:   String("au"),
			Annotations: []*CheckRunAnnotation{
				{
					Path:            String("p"),
					StartLine:       Int(1),
					EndLine:         Int(1),
					StartColumn:     Int(1),
					EndColumn:       Int(1),
					AnnotationLevel: String("al"),
					Message:         String("m"),
					Title:           String("t"),
					RawDetails:      String("rd"),
				},
			},
			Images: []*CheckRunImage{
				{
					Alt:      String("a"),
					ImageURL: String("i"),
					Caption:  String("c"),
				},
			},
		},
		Actions: []*CheckRunAction{
			{
				Label:       "l",
				Description: "d",
				Identifier:  "i",
			},
		},
	}

	want := `{
		"name": "n",
		"head_sha": "hsha",
		"details_url": "durl",
		"external_id": "eid",
		"status": "s",
		"conclusion": "c",
		"started_at": ` + referenceTimeStr + `,
		"completed_at": ` + referenceTimeStr + `,
		"output": {
			"title": "ti",
			"summary": "s",
			"text": "t",
			"annotations_count": 1,
			"annotations_url": "au",
			"annotations": [
				{
					"path": "p",
					"start_line": 1,
					"end_line": 1,
					"start_column": 1,
					"end_column": 1,
					"annotation_level": "al",
					"message": "m",
					"title": "t",
					"raw_details": "rd"
				}
			],
			"images": [
				{
					"alt": "a",
					"image_url": "i",
					"caption": "c"
				}
			]
		},
		"actions": [
			{
				"label": "l",
				"description": "d",
				"identifier": "i"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestUpdateCheckRunOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &UpdateCheckRunOptions{}, "{}")

	u := &UpdateCheckRunOptions{
		Name:        "n",
		DetailsURL:  String("durl"),
		ExternalID:  String("eid"),
		Status:      String("s"),
		Conclusion:  String("c"),
		CompletedAt: &Timestamp{referenceTime},
		Output: &CheckRunOutput{
			Title:            String("ti"),
			Summary:          String("s"),
			Text:             String("t"),
			AnnotationsCount: Int(1),
			AnnotationsURL:   String("au"),
			Annotations: []*CheckRunAnnotation{
				{
					Path:            String("p"),
					StartLine:       Int(1),
					EndLine:         Int(1),
					StartColumn:     Int(1),
					EndColumn:       Int(1),
					AnnotationLevel: String("al"),
					Message:         String("m"),
					Title:           String("t"),
					RawDetails:      String("rd"),
				},
			},
			Images: []*CheckRunImage{
				{
					Alt:      String("a"),
					ImageURL: String("i"),
					Caption:  String("c"),
				},
			},
		},
		Actions: []*CheckRunAction{
			{
				Label:       "l",
				Description: "d",
				Identifier:  "i",
			},
		},
	}

	want := `{
		"name": "n",
		"details_url": "durl",
		"external_id": "eid",
		"status": "s",
		"conclusion": "c",
		"completed_at": ` + referenceTimeStr + `,
		"output": {
			"title": "ti",
			"summary": "s",
			"text": "t",
			"annotations_count": 1,
			"annotations_url": "au",
			"annotations": [
				{
					"path": "p",
					"start_line": 1,
					"end_line": 1,
					"start_column": 1,
					"end_column": 1,
					"annotation_level": "al",
					"message": "m",
					"title": "t",
					"raw_details": "rd"
				}
			],
			"images": [
				{
					"alt": "a",
					"image_url": "i",
					"caption": "c"
				}
			]
		},
		"actions": [
			{
				"label": "l",
				"description": "d",
				"identifier": "i"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestListCheckRunsResults_Marshal(t *testing.T) {
	testJSONMarshal(t, &ListCheckRunsResults{}, "{}")

	l := &ListCheckRunsResults{
		Total: Int(1),
		CheckRuns: []*CheckRun{
			{
				ID:          Int64(1),
				NodeID:      String("n"),
				HeadSHA:     String("h"),
				ExternalID:  String("1"),
				URL:         String("u"),
				HTMLURL:     String("u"),
				DetailsURL:  String("u"),
				Status:      String("s"),
				Conclusion:  String("c"),
				StartedAt:   &Timestamp{referenceTime},
				CompletedAt: &Timestamp{referenceTime},
				Output: &CheckRunOutput{
					Annotations: []*CheckRunAnnotation{
						{
							AnnotationLevel: String("a"),
							EndLine:         Int(1),
							Message:         String("m"),
							Path:            String("p"),
							RawDetails:      String("r"),
							StartLine:       Int(1),
							Title:           String("t"),
						},
					},
					AnnotationsCount: Int(1),
					AnnotationsURL:   String("a"),
					Images: []*CheckRunImage{
						{
							Alt:      String("a"),
							ImageURL: String("i"),
							Caption:  String("c"),
						},
					},
					Title:   String("t"),
					Summary: String("s"),
					Text:    String("t"),
				},
				Name: String("n"),
				CheckSuite: &CheckSuite{
					ID: Int64(1),
				},
				App: &App{
					ID:     Int64(1),
					NodeID: String("n"),
					Owner: &User{
						Login:     String("l"),
						ID:        Int64(1),
						NodeID:    String("n"),
						URL:       String("u"),
						ReposURL:  String("r"),
						EventsURL: String("e"),
						AvatarURL: String("a"),
					},
					Name:        String("n"),
					Description: String("d"),
					HTMLURL:     String("h"),
					ExternalURL: String("u"),
					CreatedAt:   &Timestamp{referenceTime},
					UpdatedAt:   &Timestamp{referenceTime},
				},
				PullRequests: []*PullRequest{
					{
						URL:    String("u"),
						ID:     Int64(1),
						Number: Int(1),
						Head: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("s"),
								Name: String("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("u"),
								Name: String("n"),
							},
						},
					},
				},
			},
		},
	}

	w := `{
		"total_count": 1,
		"check_runs": [
			{
				"id": 1,
				"node_id": "n",
				"head_sha": "h",
				"external_id": "1",
				"url": "u",
				"html_url": "u",
				"details_url": "u",
				"status": "s",
				"conclusion": "c",
				"started_at": ` + referenceTimeStr + `,
				"completed_at": ` + referenceTimeStr + `,
				"output": {
					"title": "t",
					"summary": "s",
					"text": "t",
					"annotations_count": 1,
					"annotations_url": "a",
					"annotations": [
						{
							"path": "p",
							"start_line": 1,
							"end_line": 1,
							"annotation_level": "a",
							"message": "m",
							"title": "t",
							"raw_details": "r"
						}
					],
					"images": [
						{
							"alt": "a",
							"image_url": "i",
							"caption": "c"
						}
					]
				},
				"name": "n",
				"check_suite": {
					"id": 1
				},
				"app": {
					"id": 1,
					"node_id": "n",
					"owner": {
						"login": "l",
						"id": 1,
						"node_id": "n",
						"avatar_url": "a",
						"url": "u",
						"events_url": "e",
						"repos_url": "r"
					},
					"name": "n",
					"description": "d",
					"external_url": "u",
					"html_url": "h",
					"created_at": ` + referenceTimeStr + `,
					"updated_at": ` + referenceTimeStr + `
				},
				"pull_requests": [
					{
						"id": 1,
						"number": 1,
						"url": "u",
						"head": {
							"ref": "r",
							"sha": "s",
							"repo": {
								"id": 1,
								"name": "n",
								"url": "s"
							}
						},
						"base": {
							"ref": "r",
							"sha": "s",
							"repo": {
								"id": 1,
								"name": "n",
								"url": "u"
							}
						}
					}
				]
			}
		]
	}`

	testJSONMarshal(t, &l, w)
}

func TestListCheckSuiteResults_Marshal(t *testing.T) {
	testJSONMarshal(t, &ListCheckSuiteResults{}, "{}")

	l := &ListCheckSuiteResults{
		Total: Int(1),
		CheckSuites: []*CheckSuite{
			{
				ID:         Int64(1),
				NodeID:     String("n"),
				HeadBranch: String("h"),
				HeadSHA:    String("h"),
				URL:        String("u"),
				BeforeSHA:  String("b"),
				AfterSHA:   String("a"),
				Status:     String("s"),
				Conclusion: String("c"),
				App: &App{
					ID:     Int64(1),
					NodeID: String("n"),
					Owner: &User{
						Login:     String("l"),
						ID:        Int64(1),
						NodeID:    String("n"),
						URL:       String("u"),
						ReposURL:  String("r"),
						EventsURL: String("e"),
						AvatarURL: String("a"),
					},
					Name:        String("n"),
					Description: String("d"),
					HTMLURL:     String("h"),
					ExternalURL: String("u"),
					CreatedAt:   &Timestamp{referenceTime},
					UpdatedAt:   &Timestamp{referenceTime},
				},
				Repository: &Repository{
					ID: Int64(1),
				},
				PullRequests: []*PullRequest{
					{
						URL:    String("u"),
						ID:     Int64(1),
						Number: Int(1),
						Head: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("s"),
								Name: String("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("u"),
								Name: String("n"),
							},
						},
					},
				},
				HeadCommit: &Commit{
					SHA: String("s"),
				},
			},
		},
	}

	w := `{
		"total_count": 1,
		"check_suites": [
			{
				"id": 1,
				"node_id": "n",
				"head_branch": "h",
				"head_sha": "h",
				"url": "u",
				"before": "b",
				"after": "a",
				"status": "s",
				"conclusion": "c",
				"app": {
					"id": 1,
					"node_id": "n",
					"owner": {
						"login": "l",
						"id": 1,
						"node_id": "n",
						"avatar_url": "a",
						"url": "u",
						"events_url": "e",
						"repos_url": "r"
					},
					"name": "n",
					"description": "d",
					"external_url": "u",
					"html_url": "h",
					"created_at": ` + referenceTimeStr + `,
					"updated_at": ` + referenceTimeStr + `
				},
				"repository": {
					"id": 1
				},
				"pull_requests": [
				{
					"id": 1,
					"number": 1,
					"url": "u",
					"head": {
						"ref": "r",
						"sha": "s",
						"repo": {
							"id": 1,
							"name": "n",
							"url": "s"
						}
					},
					"base": {
						"ref": "r",
						"sha": "s",
						"repo": {
							"id": 1,
							"name": "n",
							"url": "u"
						}
					}
				}
			],
			"head_commit": {
				"sha": "s"
			}
			}
		]
	}`

	testJSONMarshal(t, &l, w)
}

func TestCheckSuitePreferenceOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckSuitePreferenceOptions{}, "{}")

	u := &CheckSuitePreferenceOptions{
		AutoTriggerChecks: []*AutoTriggerCheck{
			{
				AppID:   Int64(1),
				Setting: Bool(false),
			},
		},
	}

	want := `{
		"auto_trigger_checks": [
			{
				"app_id": 1,
				"setting": false
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestPreferenceList_Marshal(t *testing.T) {
	testJSONMarshal(t, &PreferenceList{}, "{}")

	u := &PreferenceList{
		AutoTriggerChecks: []*AutoTriggerCheck{
			{
				AppID:   Int64(1),
				Setting: Bool(false),
			},
		},
	}

	want := `{
		"auto_trigger_checks": [
			{
				"app_id": 1,
				"setting": false
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckSuitePreferenceResults_Marshal(t *testing.T) {
	testJSONMarshal(t, &CheckSuitePreferenceResults{}, "{}")

	u := &CheckSuitePreferenceResults{
		Preferences: &PreferenceList{
			AutoTriggerChecks: []*AutoTriggerCheck{
				{
					AppID:   Int64(1),
					Setting: Bool(false),
				},
			},
		},
		Repository: &Repository{
			ID:   Int64(1),
			URL:  String("u"),
			Name: String("n"),
		},
	}

	want := `{
		"preferences": {
			"auto_trigger_checks": [
				{
					"app_id": 1,
					"setting": false
				}
			]
		},
		"repository": {
			"id":1,
			"name":"n",
			"url":"u"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestChecksService_ReRequestCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-runs/1/rerequest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		w.WriteHeader(http.StatusCreated)
	})
	ctx := context.Background()
	resp, err := client.Checks.ReRequestCheckRun(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.ReRequestCheckRun return error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("Checks.ReRequestCheckRun = %v, want %v", got, want)
	}

	const methodName = "ReRequestCheckRun"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Checks.ReRequestCheckRun(ctx, "\n", "\n", 1)
		return err
	})
}
