// Copyright 2026 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListOrganizationRuleSuites(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/rule-suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"id": 101,
				"actor_id": 12,
				"actor_name": "octocat",
				"before_sha": "abc123",
				"after_sha": "def456",
				"ref": "refs/heads/main",
				"repository_id": 1,
				"repository_name": "repo",
				"pushed_at": "2023-07-06T08:43:03Z",
				"result": "bypass"
			}
		]`)
	})

	ctx := t.Context()
	suites, _, err := client.Organizations.ListOrganizationRuleSuites(ctx, "o", nil)
	if err != nil {
		t.Errorf("Organizations.ListOrganizationRuleSuites returned error: %v", err)
	}

	want := []*RuleSuite{{
		ID:             Ptr(int64(101)),
		ActorID:        Ptr(int64(12)),
		ActorName:      Ptr("octocat"),
		BeforeSHA:      Ptr("abc123"),
		AfterSHA:       Ptr("def456"),
		Ref:            Ptr("refs/heads/main"),
		RepositoryID:   Ptr(int64(1)),
		RepositoryName: Ptr("repo"),
		PushedAt:       Ptr(Timestamp{time.Date(2023, time.July, 6, 8, 43, 3, 0, time.UTC)}),
		Result:         Ptr("bypass"),
	}}

	if !cmp.Equal(suites, want) {
		t.Errorf("Organizations.ListOrganizationRuleSuites returned %+v, want %+v", suites, want)
	}

	const methodName = "ListOrganizationRuleSuites"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListOrganizationRuleSuites(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_ListOrganizationRuleSuites_ListOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/rule-suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":              "2",
			"per_page":          "35",
			"ref":               "refs/heads/main",
			"repository_name":   "repo2",
			"time_period":       "day",
			"actor_name":        "alice",
			"rule_suite_result": "pass",
			"evaluate_status":   "active",
		})
		fmt.Fprint(w, `[
			{
				"id": 201,
				"actor_id": 13,
				"actor_name": "alice",
				"before_sha": "aaa111",
				"after_sha": "bbb222",
				"ref": "refs/heads/feature",
				"repository_id": 2,
				"repository_name": "repo2",
				"pushed_at": "2023-07-07T08:43:03Z",
				"result": "pass"
			}
		]`)
	})

	opts := &ListRuleSuitesOptions{
		Ref:             "refs/heads/main",
		RepositoryName:  "repo2",
		TimePeriod:      "day",
		ActorName:       "alice",
		RuleSuiteResult: "pass",
		EvaluateStatus:  "active",
		Page:            2,
		PerPage:         35,
	}
	ctx := t.Context()
	suites, _, err := client.Organizations.ListOrganizationRuleSuites(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListOrganizationRuleSuites returned error: %v", err)
	}

	want := []*RuleSuite{{
		ID:             Ptr(int64(201)),
		ActorID:        Ptr(int64(13)),
		ActorName:      Ptr("alice"),
		BeforeSHA:      Ptr("aaa111"),
		AfterSHA:       Ptr("bbb222"),
		Ref:            Ptr("refs/heads/feature"),
		RepositoryID:   Ptr(int64(2)),
		RepositoryName: Ptr("repo2"),
		PushedAt:       Ptr(Timestamp{time.Date(2023, time.July, 7, 8, 43, 3, 0, time.UTC)}),
		Result:         Ptr("pass"),
	}}

	if !cmp.Equal(suites, want) {
		t.Errorf("Organizations.ListOrganizationRuleSuites returned %+v, want %+v", suites, want)
	}
}

func TestOrganizationsService_ListOrganizationRuleSuites_addOptionsError(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	_, _, err := client.Organizations.ListOrganizationRuleSuites(t.Context(), "\u007F", &ListRuleSuitesOptions{})
	if err == nil {
		t.Fatal("Organizations.ListOrganizationRuleSuites returned nil error, want non-nil")
	}
}

func TestOrganizationsService_GetOrganizationRuleSuite(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/rule-suites/101", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 101,
			"actor_id": 12,
			"actor_name": "octocat",
			"before_sha": "abc123",
			"after_sha": "def456",
			"ref": "refs/heads/main",
			"repository_id": 1,
			"repository_name": "repo",
			"pushed_at": "2023-07-06T08:43:03Z",
			"result": "pass",
			"evaluation_result": "fail",
			"rule_evaluations": [
				{
					"rule_source": {"type": "ruleset", "id": 3, "name": "Evaluate commit message pattern"},
					"enforcement": "evaluate",
					"result": "fail",
					"rule_type": "commit_message_pattern",
					"details": "pattern did not match"
				}
			]
		}`)
	})

	ctx := t.Context()
	suite, _, err := client.Organizations.GetOrganizationRuleSuite(ctx, "o", 101)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRuleSuite returned error: %v", err)
	}

	want := &RuleSuite{
		ID:               Ptr(int64(101)),
		ActorID:          Ptr(int64(12)),
		ActorName:        Ptr("octocat"),
		BeforeSHA:        Ptr("abc123"),
		AfterSHA:         Ptr("def456"),
		Ref:              Ptr("refs/heads/main"),
		RepositoryID:     Ptr(int64(1)),
		RepositoryName:   Ptr("repo"),
		PushedAt:         Ptr(Timestamp{time.Date(2023, time.July, 6, 8, 43, 3, 0, time.UTC)}),
		Result:           Ptr("pass"),
		EvaluationResult: Ptr("fail"),
		RuleEvaluations: []*RuleEvaluation{{
			RuleSource:  &RuleEvaluationSource{Type: Ptr("ruleset"), ID: Ptr(int64(3)), Name: Ptr("Evaluate commit message pattern")},
			Enforcement: Ptr("evaluate"),
			Result:      Ptr("fail"),
			RuleType:    Ptr("commit_message_pattern"),
			Details:     Ptr("pattern did not match"),
		}},
	}

	if !cmp.Equal(suite, want) {
		t.Errorf("Organizations.GetOrganizationRuleSuite returned %+v, want %+v", suite, want)
	}

	const methodName = "GetOrganizationRuleSuite"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrganizationRuleSuite(ctx, "o", 101)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
