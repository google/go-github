// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetRulesForBranch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rules/branches/branch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "ruleset_id": 42069,
			  "ruleset_source_type": "Repository",
			  "ruleset_source": "google/a",
			  "type": "creation"
			},
			{
			  "ruleset_id": 42069,
			  "ruleset_source_type": "Organization",
			  "ruleset_source": "google",
			  "type": "update",
			  "parameters": {
			    "update_allows_fetch_and_merge": true
			  }
			}
		]`)
	})

	ctx := context.Background()
	rules, _, err := client.Repositories.GetRulesForBranch(ctx, "o", "repo", "branch")
	if err != nil {
		t.Errorf("Repositories.GetRulesForBranch returned error: %v", err)
	}

	want := &BranchRules{
		Creation: []*BranchRuleMetadata{{RulesetSourceType: RulesetSourceTypeRepository, RulesetSource: "google/a", RulesetID: 42069}},
		Update:   []*UpdateBranchRule{{BranchRuleMetadata: BranchRuleMetadata{RulesetSourceType: RulesetSourceTypeOrganization, RulesetSource: "google", RulesetID: 42069}, Parameters: UpdateRuleParameters{UpdateAllowsFetchAndMerge: true}}},
	}

	if !cmp.Equal(rules, want) {
		t.Errorf("Repositories.GetRulesForBranch returned %+v, want %+v", rules, want)
	}

	const methodName = "GetRulesForBranch"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRulesForBranch(ctx, "o", "repo", "branch")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetAllRulesets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `[
			{
			  "id": 42,
			  "name": "ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "active",
			  "created_at": %[1]s,
			  "updated_at": %[1]s
			},
			{
			  "id": 314,
			  "name": "Another ruleset",
			  "source_type": "Repository",
			  "source": "o/repo",
			  "enforcement": "active",
			  "created_at": %[1]s,
			  "updated_at": %[1]s
			}
		]`, referenceTimeStr)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.GetAllRulesets(ctx, "o", "repo", false)
	if err != nil {
		t.Errorf("Repositories.GetAllRulesets returned error: %v", err)
	}

	want := []*RepositoryRuleset{
		{
			ID:          Ptr(int64(42)),
			Name:        "ruleset",
			SourceType:  Ptr(RulesetSourceTypeRepository),
			Source:      "o/repo",
			Enforcement: RulesetEnforcementActive,
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
		{
			ID:          Ptr(int64(314)),
			Name:        "Another ruleset",
			SourceType:  Ptr(RulesetSourceTypeRepository),
			Source:      "o/repo",
			Enforcement: RulesetEnforcementActive,
			CreatedAt:   &Timestamp{referenceTime},
			UpdatedAt:   &Timestamp{referenceTime},
		},
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.GetAllRulesets returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "GetAllRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAllRulesets(ctx, "o", "repo", false)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "active"
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.CreateRuleset(ctx, "o", "repo", RepositoryRuleset{
		Name:        "ruleset",
		Enforcement: RulesetEnforcementActive,
	})
	if err != nil {
		t.Errorf("Repositories.CreateRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr(RulesetSourceTypeRepository),
		Source:      "o/repo",
		Enforcement: RulesetEnforcementActive,
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.CreateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "CreateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRuleset(ctx, "o", "repo", RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateRulesetWithPushRules(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "active",
			"target": "push",
			"rules": [
				{
					"type": "file_path_restriction",
					"parameters": {
						"restricted_file_paths": ["/a/file"]
					}
				},
				{
					"type": "max_file_path_length",
					"parameters": {
						"max_file_path_length": 255
					}
				},
				{
					"type": "file_extension_restriction",
					"parameters": {
						"restricted_file_extensions": [".exe"]
					}
				},
				{
					"type": "max_file_size",
					"parameters": {
						"max_file_size": 1024
					}
				}
			]
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.CreateRuleset(ctx, "o", "repo", RepositoryRuleset{
		Name:        "ruleset",
		Enforcement: RulesetEnforcementActive,
	})
	if err != nil {
		t.Errorf("Repositories.CreateRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr(RulesetSourceTypeRepository),
		Source:      "o/repo",
		Target:      Ptr(RulesetTargetPush),
		Enforcement: RulesetEnforcementActive,
		Rules: &RepositoryRulesetRules{
			FilePathRestriction:      &FilePathRestrictionRuleParameters{RestrictedFilePaths: []string{"/a/file"}},
			MaxFilePathLength:        &MaxFilePathLengthRuleParameters{MaxFilePathLength: 255},
			FileExtensionRestriction: &FileExtensionRestrictionRuleParameters{RestrictedFileExtensions: []string{".exe"}},
			MaxFileSize:              &MaxFileSizeRuleParameters{MaxFileSize: 1024},
		},
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.CreateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "CreateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateRuleset(ctx, "o", "repo", RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.GetRuleset(ctx, "o", "repo", 42, true)
	if err != nil {
		t.Errorf("Repositories.GetRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: RulesetEnforcementActive,
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
	}
	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.GetRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "GetRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetRuleset(ctx, "o", "repo", 42, true)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "active"
		}`)
	})

	ctx := context.Background()
	ruleSet, _, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, RepositoryRuleset{
		Name:        "ruleset",
		Enforcement: RulesetEnforcementActive,
	})
	if err != nil {
		t.Errorf("Repositories.UpdateRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr(RulesetSourceTypeRepository),
		Source:      "o/repo",
		Enforcement: "active",
	}

	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.UpdateRuleset returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "UpdateRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRuleset(ctx, "o", "repo", 42, RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateRulesetClearBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "active"
			"conditions": {
				"ref_name": {
					"include": [
						"refs/heads/main",
						"refs/heads/master"
					],
					"exclude": [
						"refs/heads/dev*"
					]
				}
			},
			"rules": [
			  {
					"type": "creation"
			  }
			]
		}`)
	})

	ctx := context.Background()

	_, err := client.Repositories.UpdateRulesetClearBypassActor(ctx, "o", "repo", 42)
	if err != nil {
		t.Errorf("Repositories.UpdateRulesetClearBypassActor returned error: %v \n", err)
	}

	const methodName = "UpdateRulesetClearBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.UpdateRulesetClearBypassActor(ctx, "o", "repo", 42)
	})
}

func TestRepositoriesService_UpdateRulesetNoBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	rs := RepositoryRuleset{
		Name:        "ruleset",
		Source:      "o/repo",
		Enforcement: RulesetEnforcementActive,
	}

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 42,
			"name": "ruleset",
			"source_type": "Repository",
			"source": "o/repo",
			"enforcement": "active"
		}`)
	})

	ctx := context.Background()

	ruleSet, _, err := client.Repositories.UpdateRulesetNoBypassActor(ctx, "o", "repo", 42, rs)
	if err != nil {
		t.Errorf("Repositories.UpdateRulesetNoBypassActor returned error: %v \n", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(42)),
		Name:        "ruleset",
		SourceType:  Ptr(RulesetSourceTypeRepository),
		Source:      "o/repo",
		Enforcement: RulesetEnforcementActive,
	}

	if !cmp.Equal(ruleSet, want) {
		t.Errorf("Repositories.UpdateRulesetNoBypassActor returned %+v, want %+v", ruleSet, want)
	}

	const methodName = "UpdateRulesetNoBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateRulesetNoBypassActor(ctx, "o", "repo", 42, RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/repo/rulesets/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteRuleset(ctx, "o", "repo", 42)
	if err != nil {
		t.Errorf("Repositories.DeleteRuleset returned error: %v", err)
	}

	const methodName = "DeleteRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteRuleset(ctx, "o", "repo", 42)
	})
}
