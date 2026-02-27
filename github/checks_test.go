// Copyright 2018 The go-github AUTHORS. All rights reserved.
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

func TestChecksService_GetCheckRun(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	checkRun, _, err := client.Checks.GetCheckRun(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.GetCheckRun return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	completeAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")

	want := &CheckRun{
		ID:          Ptr(int64(1)),
		Status:      Ptr("completed"),
		Conclusion:  Ptr("neutral"),
		StartedAt:   &Timestamp{startedAt},
		CompletedAt: &Timestamp{completeAt},
		Name:        Ptr("testCheckRun"),
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
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	checkSuite, _, err := client.Checks.GetCheckSuite(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.GetCheckSuite return error: %v", err)
	}
	want := &CheckSuite{
		ID:         Ptr(int64(1)),
		HeadBranch: Ptr("master"),
		HeadSHA:    Ptr("deadbeef"),
		AfterSHA:   Ptr("deadbeefa"),
		BeforeSHA:  Ptr("deadbeefb"),
		Status:     Ptr("completed"),
		Conclusion: Ptr("neutral"),
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		Status:    Ptr("in_progress"),
		StartedAt: &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   Ptr("Mighty test report"),
			Summary: Ptr(""),
			Text:    Ptr(""),
		},
	}

	ctx := t.Context()
	checkRun, _, err := client.Checks.CreateCheckRun(ctx, "o", "r", checkRunOpt)
	if err != nil {
		t.Errorf("Checks.CreateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:        Ptr(int64(1)),
		Status:    Ptr("in_progress"),
		StartedAt: &Timestamp{startedAt},
		HeadSHA:   Ptr("deadbeef"),
		Name:      Ptr("testCreateCheckRun"),
		Output: &CheckRunOutput{
			Title:   Ptr("Mighty test report"),
			Summary: Ptr(""),
			Text:    Ptr(""),
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
	t.Parallel()
	client, mux, _ := setup(t)

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

	ctx := t.Context()
	checkRunAnnotations, _, err := client.Checks.ListCheckRunAnnotations(ctx, "o", "r", 1, &ListOptions{Page: 1})
	if err != nil {
		t.Errorf("Checks.ListCheckRunAnnotations return error: %v", err)
	}

	want := []*CheckRunAnnotation{{
		Path:            Ptr("README.md"),
		StartLine:       Ptr(2),
		EndLine:         Ptr(2),
		StartColumn:     Ptr(1),
		EndColumn:       Ptr(5),
		AnnotationLevel: Ptr("warning"),
		Message:         Ptr("Check your spelling for 'banaas'."),
		Title:           Ptr("Spell check"),
		RawDetails:      Ptr("Do you mean 'bananas' or 'banana'?"),
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		Status:      Ptr("completed"),
		CompletedAt: &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   Ptr("Mighty test report"),
			Summary: Ptr("There are 0 failures, 2 warnings and 1 notice"),
			Text:    Ptr("You may have misspelled some words."),
		},
	}

	ctx := t.Context()
	checkRun, _, err := client.Checks.UpdateCheckRun(ctx, "o", "r", 1, updateCheckRunOpt)
	if err != nil {
		t.Errorf("Checks.UpdateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:          Ptr(int64(1)),
		Status:      Ptr("completed"),
		StartedAt:   &Timestamp{startedAt},
		CompletedAt: &Timestamp{startedAt},
		Conclusion:  Ptr("neutral"),
		Name:        Ptr("testUpdateCheckRun"),
		Output: &CheckRunOutput{
			Title:   Ptr("Mighty test report"),
			Summary: Ptr("There are 0 failures, 2 warnings and 1 notice"),
			Text:    Ptr("You may have misspelled some words."),
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		CheckName:   Ptr("testing"),
		Status:      Ptr("completed"),
		Filter:      Ptr("all"),
		AppID:       Ptr(int64(1)),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := t.Context()
	checkRuns, _, err := client.Checks.ListCheckRunsForRef(ctx, "o", "r", "master", opt)
	if err != nil {
		t.Errorf("Checks.ListCheckRunsForRef return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	want := &ListCheckRunsResults{
		Total: Ptr(1),
		CheckRuns: []*CheckRun{{
			ID:          Ptr(int64(1)),
			Status:      Ptr("completed"),
			StartedAt:   &Timestamp{startedAt},
			CompletedAt: &Timestamp{startedAt},
			Conclusion:  Ptr("neutral"),
			HeadSHA:     Ptr("deadbeef"),
			App:         &App{ID: Ptr(int64(1))},
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		CheckName:   Ptr("testing"),
		Status:      Ptr("completed"),
		Filter:      Ptr("all"),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := t.Context()
	checkRuns, _, err := client.Checks.ListCheckRunsCheckSuite(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Checks.ListCheckRunsCheckSuite return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	want := &ListCheckRunsResults{
		Total: Ptr(1),
		CheckRuns: []*CheckRun{{
			ID:          Ptr(int64(1)),
			Status:      Ptr("completed"),
			StartedAt:   &Timestamp{startedAt},
			CompletedAt: &Timestamp{startedAt},
			Conclusion:  Ptr("neutral"),
			HeadSHA:     Ptr("deadbeef"),
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		CheckName:   Ptr("testing"),
		AppID:       Ptr(int64(2)),
		ListOptions: ListOptions{Page: 1},
	}
	ctx := t.Context()
	checkSuites, _, err := client.Checks.ListCheckSuitesForRef(ctx, "o", "r", "master", opt)
	if err != nil {
		t.Errorf("Checks.ListCheckSuitesForRef return error: %v", err)
	}
	want := &ListCheckSuiteResults{
		Total: Ptr(1),
		CheckSuites: []*CheckSuite{{
			ID:         Ptr(int64(1)),
			Status:     Ptr("completed"),
			Conclusion: Ptr("neutral"),
			HeadSHA:    Ptr("deadbeef"),
			HeadBranch: Ptr("master"),
			BeforeSHA:  Ptr("deadbeefb"),
			AfterSHA:   Ptr("deadbeefa"),
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/check-suites/preferences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		testBody(t, r, `{"auto_trigger_checks":[{"app_id":2,"setting":false}]}`+"\n")
		fmt.Fprint(w, `{"preferences":{"auto_trigger_checks":[{"app_id": 2,"setting": false}]}}`)
	})
	a := []*AutoTriggerCheck{{
		AppID:   Ptr(int64(2)),
		Setting: Ptr(false),
	}}
	opt := CheckSuitePreferenceOptions{AutoTriggerChecks: a}
	ctx := t.Context()
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
	t.Parallel()
	client, mux, _ := setup(t)

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
		HeadBranch: Ptr("master"),
	}

	ctx := t.Context()
	checkSuite, _, err := client.Checks.CreateCheckSuite(ctx, "o", "r", checkSuiteOpt)
	if err != nil {
		t.Errorf("Checks.CreateCheckSuite return error: %v", err)
	}

	want := &CheckSuite{
		ID:         Ptr(int64(2)),
		Status:     Ptr("completed"),
		HeadSHA:    Ptr("deadbeef"),
		HeadBranch: Ptr("master"),
		Conclusion: Ptr("neutral"),
		BeforeSHA:  Ptr("deadbeefb"),
		AfterSHA:   Ptr("deadbeefa"),
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/check-suites/1/rerequest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		w.WriteHeader(http.StatusCreated)
	})
	ctx := t.Context()
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
	t.Parallel()
	testJSONMarshal(t, &CheckRun{}, "{}")

	now := time.Now()
	ts := now.Format(time.RFC3339Nano)

	c := CheckRun{
		ID:          Ptr(int64(1)),
		NodeID:      Ptr("n"),
		HeadSHA:     Ptr("h"),
		ExternalID:  Ptr("1"),
		URL:         Ptr("u"),
		HTMLURL:     Ptr("u"),
		DetailsURL:  Ptr("u"),
		Status:      Ptr("s"),
		Conclusion:  Ptr("c"),
		StartedAt:   &Timestamp{Time: now},
		CompletedAt: &Timestamp{Time: now},
		Output: &CheckRunOutput{
			Annotations: []*CheckRunAnnotation{
				{
					AnnotationLevel: Ptr("a"),
					EndLine:         Ptr(1),
					Message:         Ptr("m"),
					Path:            Ptr("p"),
					RawDetails:      Ptr("r"),
					StartLine:       Ptr(1),
					Title:           Ptr("t"),
				},
			},
			AnnotationsCount: Ptr(1),
			AnnotationsURL:   Ptr("a"),
			Images: []*CheckRunImage{
				{
					Alt:      Ptr("a"),
					ImageURL: Ptr("i"),
					Caption:  Ptr("c"),
				},
			},
			Title:   Ptr("t"),
			Summary: Ptr("s"),
			Text:    Ptr("t"),
		},
		Name: Ptr("n"),
		CheckSuite: &CheckSuite{
			ID: Ptr(int64(1)),
		},
		App: &App{
			ID:     Ptr(int64(1)),
			NodeID: Ptr("n"),
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Name:        Ptr("n"),
			Description: Ptr("d"),
			HTMLURL:     Ptr("h"),
			ExternalURL: Ptr("u"),
			CreatedAt:   &Timestamp{now},
			UpdatedAt:   &Timestamp{now},
		},
		PullRequests: []*PullRequest{
			{
				URL:    Ptr("u"),
				ID:     Ptr(int64(1)),
				Number: Ptr(1),
				Head: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("s"),
						Name: Ptr("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("u"),
						Name: Ptr("n"),
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
		"started_at": "%v",
		"completed_at": "%v",
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
			"created_at": "%v",
			"updated_at": "%v"
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
	t.Parallel()
	testJSONMarshal(t, &CheckSuite{}, "{}")

	now := time.Now()
	ts := now.Format(time.RFC3339Nano)

	c := CheckSuite{
		ID:         Ptr(int64(1)),
		NodeID:     Ptr("n"),
		HeadBranch: Ptr("h"),
		HeadSHA:    Ptr("h"),
		URL:        Ptr("u"),
		BeforeSHA:  Ptr("b"),
		AfterSHA:   Ptr("a"),
		Status:     Ptr("s"),
		Conclusion: Ptr("c"),
		App: &App{
			ID:     Ptr(int64(1)),
			NodeID: Ptr("n"),
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Name:        Ptr("n"),
			Description: Ptr("d"),
			HTMLURL:     Ptr("h"),
			ExternalURL: Ptr("u"),
			CreatedAt:   &Timestamp{now},
			UpdatedAt:   &Timestamp{now},
		},
		Repository: &Repository{
			ID: Ptr(int64(1)),
		},
		PullRequests: []*PullRequest{
			{
				URL:    Ptr("u"),
				ID:     Ptr(int64(1)),
				Number: Ptr(1),
				Head: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("s"),
						Name: Ptr("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("u"),
						Name: Ptr("n"),
					},
				},
			},
		},
		HeadCommit: &Commit{
			SHA: Ptr("s"),
		},
		LatestCheckRunsCount: Ptr(int64(1)),
		Rerequestable:        Ptr(true),
		RunsRerequestable:    Ptr(true),
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
				"created_at": "%v",
				"updated_at": "%v"
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
		},
		"latest_check_runs_count": 1,
		"rerequestable": true,
		"runs_rerequestable": true
		}`, ts, ts)

	testJSONMarshal(t, &c, w)
}

func TestCheckRunAnnotation_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CheckRunAnnotation{}, "{}")

	u := &CheckRunAnnotation{
		Path:            Ptr("p"),
		StartLine:       Ptr(1),
		EndLine:         Ptr(1),
		StartColumn:     Ptr(1),
		EndColumn:       Ptr(1),
		AnnotationLevel: Ptr("al"),
		Message:         Ptr("m"),
		Title:           Ptr("t"),
		RawDetails:      Ptr("rd"),
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
	t.Parallel()
	testJSONMarshal(t, &CheckRunImage{}, "{}")

	u := &CheckRunImage{
		Alt:      Ptr("a"),
		ImageURL: Ptr("i"),
		Caption:  Ptr("c"),
	}

	want := `{
		"alt": "a",
		"image_url": "i",
		"caption": "c"
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckRunAction_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CheckRunAction{}, `{
		"label": "",
		"description": "",
		"identifier": ""
	}`)

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
	t.Parallel()
	testJSONMarshal(t, &AutoTriggerCheck{}, "{}")

	u := &AutoTriggerCheck{
		AppID:   Ptr(int64(1)),
		Setting: Ptr(false),
	}

	want := `{
		"app_id": 1,
		"setting": false
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateCheckSuiteOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CreateCheckSuiteOptions{}, `{"head_sha": ""}`)

	u := &CreateCheckSuiteOptions{
		HeadSHA:    "hsha",
		HeadBranch: Ptr("hb"),
	}

	want := `{
		"head_sha": "hsha",
		"head_branch": "hb"
	}`

	testJSONMarshal(t, u, want)
}

func TestCheckRunOutput_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CheckRunOutput{}, "{}")

	u := &CheckRunOutput{
		Title:            Ptr("ti"),
		Summary:          Ptr("s"),
		Text:             Ptr("t"),
		AnnotationsCount: Ptr(1),
		AnnotationsURL:   Ptr("au"),
		Annotations: []*CheckRunAnnotation{
			{
				Path:            Ptr("p"),
				StartLine:       Ptr(1),
				EndLine:         Ptr(1),
				StartColumn:     Ptr(1),
				EndColumn:       Ptr(1),
				AnnotationLevel: Ptr("al"),
				Message:         Ptr("m"),
				Title:           Ptr("t"),
				RawDetails:      Ptr("rd"),
			},
		},
		Images: []*CheckRunImage{
			{
				Alt:      Ptr("a"),
				ImageURL: Ptr("i"),
				Caption:  Ptr("c"),
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
	t.Parallel()
	testJSONMarshal(t, &CreateCheckRunOptions{}, `{
		"name": "",
		"head_sha": ""
	}`)

	u := &CreateCheckRunOptions{
		Name:        "n",
		HeadSHA:     "hsha",
		DetailsURL:  Ptr("durl"),
		ExternalID:  Ptr("eid"),
		Status:      Ptr("s"),
		Conclusion:  Ptr("c"),
		StartedAt:   &Timestamp{referenceTime},
		CompletedAt: &Timestamp{referenceTime},
		Output: &CheckRunOutput{
			Title:            Ptr("ti"),
			Summary:          Ptr("s"),
			Text:             Ptr("t"),
			AnnotationsCount: Ptr(1),
			AnnotationsURL:   Ptr("au"),
			Annotations: []*CheckRunAnnotation{
				{
					Path:            Ptr("p"),
					StartLine:       Ptr(1),
					EndLine:         Ptr(1),
					StartColumn:     Ptr(1),
					EndColumn:       Ptr(1),
					AnnotationLevel: Ptr("al"),
					Message:         Ptr("m"),
					Title:           Ptr("t"),
					RawDetails:      Ptr("rd"),
				},
			},
			Images: []*CheckRunImage{
				{
					Alt:      Ptr("a"),
					ImageURL: Ptr("i"),
					Caption:  Ptr("c"),
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
	t.Parallel()
	testJSONMarshal(t, &UpdateCheckRunOptions{}, `{"name": ""}`)

	u := &UpdateCheckRunOptions{
		Name:        "n",
		DetailsURL:  Ptr("durl"),
		ExternalID:  Ptr("eid"),
		Status:      Ptr("s"),
		Conclusion:  Ptr("c"),
		CompletedAt: &Timestamp{referenceTime},
		Output: &CheckRunOutput{
			Title:            Ptr("ti"),
			Summary:          Ptr("s"),
			Text:             Ptr("t"),
			AnnotationsCount: Ptr(1),
			AnnotationsURL:   Ptr("au"),
			Annotations: []*CheckRunAnnotation{
				{
					Path:            Ptr("p"),
					StartLine:       Ptr(1),
					EndLine:         Ptr(1),
					StartColumn:     Ptr(1),
					EndColumn:       Ptr(1),
					AnnotationLevel: Ptr("al"),
					Message:         Ptr("m"),
					Title:           Ptr("t"),
					RawDetails:      Ptr("rd"),
				},
			},
			Images: []*CheckRunImage{
				{
					Alt:      Ptr("a"),
					ImageURL: Ptr("i"),
					Caption:  Ptr("c"),
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
	t.Parallel()
	testJSONMarshal(t, &ListCheckRunsResults{}, "{}")

	l := &ListCheckRunsResults{
		Total: Ptr(1),
		CheckRuns: []*CheckRun{
			{
				ID:          Ptr(int64(1)),
				NodeID:      Ptr("n"),
				HeadSHA:     Ptr("h"),
				ExternalID:  Ptr("1"),
				URL:         Ptr("u"),
				HTMLURL:     Ptr("u"),
				DetailsURL:  Ptr("u"),
				Status:      Ptr("s"),
				Conclusion:  Ptr("c"),
				StartedAt:   &Timestamp{referenceTime},
				CompletedAt: &Timestamp{referenceTime},
				Output: &CheckRunOutput{
					Annotations: []*CheckRunAnnotation{
						{
							AnnotationLevel: Ptr("a"),
							EndLine:         Ptr(1),
							Message:         Ptr("m"),
							Path:            Ptr("p"),
							RawDetails:      Ptr("r"),
							StartLine:       Ptr(1),
							Title:           Ptr("t"),
						},
					},
					AnnotationsCount: Ptr(1),
					AnnotationsURL:   Ptr("a"),
					Images: []*CheckRunImage{
						{
							Alt:      Ptr("a"),
							ImageURL: Ptr("i"),
							Caption:  Ptr("c"),
						},
					},
					Title:   Ptr("t"),
					Summary: Ptr("s"),
					Text:    Ptr("t"),
				},
				Name: Ptr("n"),
				CheckSuite: &CheckSuite{
					ID: Ptr(int64(1)),
				},
				App: &App{
					ID:     Ptr(int64(1)),
					NodeID: Ptr("n"),
					Owner: &User{
						Login:     Ptr("l"),
						ID:        Ptr(int64(1)),
						NodeID:    Ptr("n"),
						URL:       Ptr("u"),
						ReposURL:  Ptr("r"),
						EventsURL: Ptr("e"),
						AvatarURL: Ptr("a"),
					},
					Name:        Ptr("n"),
					Description: Ptr("d"),
					HTMLURL:     Ptr("h"),
					ExternalURL: Ptr("u"),
					CreatedAt:   &Timestamp{referenceTime},
					UpdatedAt:   &Timestamp{referenceTime},
				},
				PullRequests: []*PullRequest{
					{
						URL:    Ptr("u"),
						ID:     Ptr(int64(1)),
						Number: Ptr(1),
						Head: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("s"),
								Name: Ptr("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("u"),
								Name: Ptr("n"),
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
	t.Parallel()
	testJSONMarshal(t, &ListCheckSuiteResults{}, "{}")

	l := &ListCheckSuiteResults{
		Total: Ptr(1),
		CheckSuites: []*CheckSuite{
			{
				ID:         Ptr(int64(1)),
				NodeID:     Ptr("n"),
				HeadBranch: Ptr("h"),
				HeadSHA:    Ptr("h"),
				URL:        Ptr("u"),
				BeforeSHA:  Ptr("b"),
				AfterSHA:   Ptr("a"),
				Status:     Ptr("s"),
				Conclusion: Ptr("c"),
				App: &App{
					ID:     Ptr(int64(1)),
					NodeID: Ptr("n"),
					Owner: &User{
						Login:     Ptr("l"),
						ID:        Ptr(int64(1)),
						NodeID:    Ptr("n"),
						URL:       Ptr("u"),
						ReposURL:  Ptr("r"),
						EventsURL: Ptr("e"),
						AvatarURL: Ptr("a"),
					},
					Name:        Ptr("n"),
					Description: Ptr("d"),
					HTMLURL:     Ptr("h"),
					ExternalURL: Ptr("u"),
					CreatedAt:   &Timestamp{referenceTime},
					UpdatedAt:   &Timestamp{referenceTime},
				},
				Repository: &Repository{
					ID: Ptr(int64(1)),
				},
				PullRequests: []*PullRequest{
					{
						URL:    Ptr("u"),
						ID:     Ptr(int64(1)),
						Number: Ptr(1),
						Head: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("s"),
								Name: Ptr("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("u"),
								Name: Ptr("n"),
							},
						},
					},
				},
				HeadCommit: &Commit{
					SHA: Ptr("s"),
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
	t.Parallel()
	testJSONMarshal(t, &CheckSuitePreferenceOptions{}, "{}")

	u := &CheckSuitePreferenceOptions{
		AutoTriggerChecks: []*AutoTriggerCheck{
			{
				AppID:   Ptr(int64(1)),
				Setting: Ptr(false),
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
	t.Parallel()
	testJSONMarshal(t, &PreferenceList{}, "{}")

	u := &PreferenceList{
		AutoTriggerChecks: []*AutoTriggerCheck{
			{
				AppID:   Ptr(int64(1)),
				Setting: Ptr(false),
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
	t.Parallel()
	testJSONMarshal(t, &CheckSuitePreferenceResults{}, "{}")

	u := &CheckSuitePreferenceResults{
		Preferences: &PreferenceList{
			AutoTriggerChecks: []*AutoTriggerCheck{
				{
					AppID:   Ptr(int64(1)),
					Setting: Ptr(false),
				},
			},
		},
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("u"),
			Name: Ptr("n"),
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/check-runs/1/rerequest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		w.WriteHeader(http.StatusCreated)
	})
	ctx := t.Context()
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
