// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListRuleSuites(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/rule-suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"ref":               "refs/heads/main",
			"repository_name":   "r",
			"time_period":       "week",
			"actor_name":        "octocat",
			"rule_suite_result": "fail",
			"evaluate_status":   "evaluate",
			"page":              "2",
			"per_page":          "30",
		})
		fmt.Fprint(w, `[{
			"id": 21,
			"actor_id": 12,
			"actor_name": "octocat",
			"before_sha": "893f768e172fb1bc9c5d6f3dd48557e45f14e01d",
			"after_sha": "dedd88641a362b6b4ea872da4847d6131a164d01",
			"ref": "refs/heads/i-see-everything",
			"repository_id": 404,
			"repository_name": "r",
			"pushed_at": `+referenceTimeStr+`,
			"result": "bypass",
			"evaluation_result": "fail"
		}]`)
	})

	ctx := t.Context()
	opts := &ListRuleSuitesOptions{
		Ref:             "refs/heads/main",
		RepositoryName:  "r",
		TimePeriod:      "week",
		ActorName:       "octocat",
		RuleSuiteResult: "fail",
		EvaluateStatus:  "evaluate",
		ListOptions:     ListOptions{Page: 2, PerPage: 30},
	}
	ruleSuites, _, err := client.Organizations.ListRuleSuites(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListRuleSuites returned error: %v", err)
	}

	want := []*RuleSuite{{
		ID:               new(int64(21)),
		ActorID:          new(int64(12)),
		ActorName:        new("octocat"),
		BeforeSHA:        new("893f768e172fb1bc9c5d6f3dd48557e45f14e01d"),
		AfterSHA:         new("dedd88641a362b6b4ea872da4847d6131a164d01"),
		Ref:              new("refs/heads/i-see-everything"),
		RepositoryID:     new(int64(404)),
		RepositoryName:   new("r"),
		PushedAt:         &referenceTimestamp,
		Result:           new("bypass"),
		EvaluationResult: new("fail"),
	}}
	if !cmp.Equal(ruleSuites, want) {
		t.Errorf("Organizations.ListRuleSuites returned %+v, want %+v", ruleSuites, want)
	}

	const methodName = "ListRuleSuites"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListRuleSuites(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListRuleSuites(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetRuleSuite(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/rule-suites/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"include_metadata": "true"})
		fmt.Fprint(w, `{
			"id": 21,
			"actor_id": 12,
			"actor_name": "octocat",
			"before_sha": "893f768e172fb1bc9c5d6f3dd48557e45f14e01d",
			"after_sha": "dedd88641a362b6b4ea872da4847d6131a164d01",
			"ref": "refs/heads/i-see-everything",
			"repository_id": 404,
			"repository_name": "r",
			"pushed_at": `+referenceTimeStr+`,
			"result": "bypass",
			"evaluation_result": "fail",
			"rule_evaluations": [
				{
					"rule_source": {
						"type": "ruleset",
						"id": 2,
						"name": "Author email must be a GitHub email address"
					},
					"enforcement": "active",
					"result": "pass",
					"rule_type": "commit_author_email_pattern",
					"details": "Commit author email is valid"
				},
				{
					"rule_source": {"type": "ruleset", "id": 3, "name": "Require reviews"},
					"enforcement": "active",
					"result": "fail",
					"rule_type": "pull_request",
					"metadata": {
						"pull_request": {
							"id": 100,
							"number": 7,
							"user": {"id": 1, "login": "octocat", "type": "User"},
							"reviews": [{"id": 200, "user": {"id": 2, "login": "hubot", "type": "User"}, "state": "APPROVED"}],
							"proof_of_presence": {
								"state": "pending",
								"code": "not_verified",
								"base_sha": "",
								"head_sha": "",
								"required_level": "",
								"verified_level": "",
								"verified_at": "",
								"verifier": null
							}
						}
					}
				},
				{
					"rule_source": {"type": "ruleset", "id": 4, "name": "Require CI"},
					"enforcement": "active",
					"result": "pass",
					"rule_type": "required_status_checks",
					"metadata": {
						"checks": [{
							"id": 300,
							"context": "ci/build",
							"state": "success",
							"type": "check_run",
							"app": {"id": 15368, "slug": "github-actions", "name": "GitHub Actions"}
						}]
					}
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &GetRuleSuiteOptions{IncludeMetadata: true}
	ruleSuite, _, err := client.Organizations.GetRuleSuite(ctx, "o", 21, opts)
	if err != nil {
		t.Errorf("Organizations.GetRuleSuite returned error: %v", err)
	}

	want := &RuleSuite{
		ID:               new(int64(21)),
		ActorID:          new(int64(12)),
		ActorName:        new("octocat"),
		BeforeSHA:        new("893f768e172fb1bc9c5d6f3dd48557e45f14e01d"),
		AfterSHA:         new("dedd88641a362b6b4ea872da4847d6131a164d01"),
		Ref:              new("refs/heads/i-see-everything"),
		RepositoryID:     new(int64(404)),
		RepositoryName:   new("r"),
		PushedAt:         &referenceTimestamp,
		Result:           new("bypass"),
		EvaluationResult: new("fail"),
		RuleEvaluations: []*RuleEvaluation{
			{
				RuleSource: &RuleEvaluationSource{
					Type: new("ruleset"),
					ID:   new(int64(2)),
					Name: new("Author email must be a GitHub email address"),
				},
				Enforcement: new("active"),
				Result:      new("pass"),
				RuleType:    new("commit_author_email_pattern"),
				Details:     new("Commit author email is valid"),
			},
			{
				RuleSource:  &RuleEvaluationSource{Type: new("ruleset"), ID: new(int64(3)), Name: new("Require reviews")},
				Enforcement: new("active"),
				Result:      new("fail"),
				RuleType:    new("pull_request"),
				Metadata: &RuleEvaluationMetadata{
					PullRequest: &RuleEvaluationPullRequest{
						ID:     new(int64(100)),
						Number: new(7),
						User:   &User{ID: new(int64(1)), Login: new("octocat"), Type: new("User")},
						Reviews: []*PullRequestReview{{
							ID:    new(int64(200)),
							User:  &User{ID: new(int64(2)), Login: new("hubot"), Type: new("User")},
							State: new("APPROVED"),
						}},
						ProofOfPresence: &RuleEvaluationProofOfPresence{
							State:         new("pending"),
							Code:          new("not_verified"),
							BaseSHA:       new(""),
							HeadSHA:       new(""),
							RequiredLevel: new(""),
							VerifiedLevel: new(""),
							VerifiedAt:    new(""),
						},
					},
				},
			},
			{
				RuleSource:  &RuleEvaluationSource{Type: new("ruleset"), ID: new(int64(4)), Name: new("Require CI")},
				Enforcement: new("active"),
				Result:      new("pass"),
				RuleType:    new("required_status_checks"),
				Metadata: &RuleEvaluationMetadata{
					Checks: []*RuleEvaluationStatusCheck{{
						ID:      new(int64(300)),
						Context: new("ci/build"),
						State:   new("success"),
						Type:    new("check_run"),
						App:     &App{ID: new(int64(15368)), Slug: new("github-actions"), Name: new("GitHub Actions")},
					}},
				},
			},
		},
	}
	if !cmp.Equal(ruleSuite, want) {
		t.Errorf("Organizations.GetRuleSuite returned %+v, want %+v", ruleSuite, want)
	}

	const methodName = "GetRuleSuite"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetRuleSuite(ctx, "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetRuleSuite(ctx, "o", 21, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
