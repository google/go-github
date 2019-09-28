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
	checkSuite, _, err := client.Checks.GetCheckSuite(context.Background(), "o", "r", 1)
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
	if !reflect.DeepEqual(checkSuite, want) {
		t.Errorf("Checks.GetCheckSuite return %+v, want %+v", checkSuite, want)
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

	checkRunAnnotations, _, err := client.Checks.ListCheckRunAnnotations(context.Background(), "o", "r", 1, &ListOptions{Page: 1})
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

	if !reflect.DeepEqual(checkRunAnnotations, want) {
		t.Errorf("Checks.ListCheckRunAnnotations returned %+v, want %+v", checkRunAnnotations, want)
	}
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
                        "head_sha":"deadbeef",
			"status": "completed",
			"conclusion": "neutral",
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": "2018-05-04T01:14:52Z",
                        "output":{"title": "Mighty test report", "summary":"There are 0 failures, 2 warnings and 1 notice", "text":"You may have misspelled some words."}}`)
	})
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	updateCheckRunOpt := UpdateCheckRunOptions{
		Name:        "testUpdateCheckRun",
		HeadSHA:     String("deadbeef"),
		Status:      String("completed"),
		CompletedAt: &Timestamp{startedAt},
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String("There are 0 failures, 2 warnings and 1 notice"),
			Text:    String("You may have misspelled some words."),
		},
	}

	checkRun, _, err := client.Checks.UpdateCheckRun(context.Background(), "o", "r", 1, updateCheckRunOpt)
	if err != nil {
		t.Errorf("Checks.UpdateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:          Int64(1),
		Status:      String("completed"),
		StartedAt:   &Timestamp{startedAt},
		CompletedAt: &Timestamp{startedAt},
		Conclusion:  String("neutral"),
		HeadSHA:     String("deadbeef"),
		Name:        String("testUpdateCheckRun"),
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String("There are 0 failures, 2 warnings and 1 notice"),
			Text:    String("You may have misspelled some words."),
		},
	}
	if !reflect.DeepEqual(checkRun, want) {
		t.Errorf("Checks.UpdateCheckRun return %+v, want %+v", checkRun, want)
	}
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
	checkRuns, _, err := client.Checks.ListCheckRunsForRef(context.Background(), "o", "r", "master", opt)
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
		}},
	}

	if !reflect.DeepEqual(checkRuns, want) {
		t.Errorf("Checks.ListCheckRunsForRef returned %+v, want %+v", checkRuns, want)
	}
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
	checkRuns, _, err := client.Checks.ListCheckRunsCheckSuite(context.Background(), "o", "r", 1, opt)
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

	if !reflect.DeepEqual(checkRuns, want) {
		t.Errorf("Checks.ListCheckRunsCheckSuite returned %+v, want %+v", checkRuns, want)
	}
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
	checkSuites, _, err := client.Checks.ListCheckSuitesForRef(context.Background(), "o", "r", "master", opt)
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

	if !reflect.DeepEqual(checkSuites, want) {
		t.Errorf("Checks.ListCheckSuitesForRef returned %+v, want %+v", checkSuites, want)
	}
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
	prefResults, _, err := client.Checks.SetCheckSuitePreferences(context.Background(), "o", "r", opt)
	if err != nil {
		t.Errorf("Checks.SetCheckSuitePreferences return error: %v", err)
	}

	p := &PreferenceList{
		AutoTriggerChecks: a,
	}
	want := &CheckSuitePreferenceResults{
		Preferences: p,
	}

	if !reflect.DeepEqual(prefResults, want) {
		t.Errorf("Checks.SetCheckSuitePreferences return %+v, want %+v", prefResults, want)
	}
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

	checkSuite, _, err := client.Checks.CreateCheckSuite(context.Background(), "o", "r", checkSuiteOpt)
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
	if !reflect.DeepEqual(checkSuite, want) {
		t.Errorf("Checks.CreateCheckSuite return %+v, want %+v", checkSuite, want)
	}
}

func TestChecksService_ReRequestCheckSuite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/check-suites/1/rerequest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeCheckRunsPreview)
		w.WriteHeader(http.StatusCreated)
	})
	resp, err := client.Checks.ReRequestCheckSuite(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.ReRequestCheckSuite return error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("Checks.ReRequestCheckSuite = %v, want %v", got, want)
	}
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
