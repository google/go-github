// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRequiredReviewer_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var testCases = map[string]struct {
		data      []byte
		wantRule  []*RequiredReviewer
		wantError bool
	}{
		"User Reviewer": {
			data:      []byte(`[{"type": "User", "reviewer": {"id": 1,"login": "octocat"}}]`),
			wantRule:  []*RequiredReviewer{{Type: String("User"), Reviewer: &User{ID: Int64(1), Login: String("octocat")}}},
			wantError: false,
		},
		"Team Reviewer": {
			data:      []byte(`[{"type": "Team", "reviewer": {"id": 1, "name": "Justice League"}}]`),
			wantRule:  []*RequiredReviewer{{Type: String("Team"), Reviewer: &Team{ID: Int64(1), Name: String("Justice League")}}},
			wantError: false,
		},
		"Both Types Reviewer": {
			data:      []byte(`[{"type": "User", "reviewer": {"id": 1,"login": "octocat"}},{"type": "Team", "reviewer": {"id": 1, "name": "Justice League"}}]`),
			wantRule:  []*RequiredReviewer{{Type: String("User"), Reviewer: &User{ID: Int64(1), Login: String("octocat")}}, {Type: String("Team"), Reviewer: &Team{ID: Int64(1), Name: String("Justice League")}}},
			wantError: false,
		},
		"Empty JSON Object": {
			data:      []byte(`[]`),
			wantRule:  []*RequiredReviewer{},
			wantError: false,
		},
		"Bad JSON Object": {
			data:      []byte(`[badjson: 1]`),
			wantRule:  []*RequiredReviewer{},
			wantError: true,
		},
		"Wrong Type Type in Reviewer Object": {
			data:      []byte(`[{"type": 1, "reviewer": {"id": 1}}]`),
			wantRule:  []*RequiredReviewer{{Type: nil, Reviewer: nil}},
			wantError: true,
		},
		"Wrong ID Type in User Object": {
			data:      []byte(`[{"type": "User", "reviewer": {"id": "string"}}]`),
			wantRule:  []*RequiredReviewer{{Type: String("User"), Reviewer: nil}},
			wantError: true,
		},
		"Wrong ID Type in Team Object": {
			data:      []byte(`[{"type": "Team", "reviewer": {"id": "string"}}]`),
			wantRule:  []*RequiredReviewer{{Type: String("Team"), Reviewer: nil}},
			wantError: true,
		},
		"Wrong Type of Reviewer": {
			data:      []byte(`[{"type": "Cat", "reviewer": {"id": 1,"login": "octocat"}}]`),
			wantRule:  []*RequiredReviewer{{Type: nil, Reviewer: nil}},
			wantError: true,
		},
	}

	for name, test := range testCases {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			rule := []*RequiredReviewer{}
			err := json.Unmarshal(test.data, &rule)
			if err != nil && !test.wantError {
				t.Errorf("RequiredReviewer.UnmarshalJSON returned an error when we expected nil")
			}
			if err == nil && test.wantError {
				t.Errorf("RequiredReviewer.UnmarshalJSON returned no error when we expected one")
			}
			if !cmp.Equal(test.wantRule, rule) {
				t.Errorf("RequiredReviewer.UnmarshalJSON expected rule %+v, got %+v", test.wantRule, rule)
			}
		})
	}
}

func TestCreateUpdateEnvironment_MarshalJSON(t *testing.T) {
	t.Parallel()
	cu := &CreateUpdateEnvironment{}

	got, err := cu.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON: %v", err)
	}

	want := `{"wait_timer":0,"reviewers":null,"can_admins_bypass":true,"deployment_branch_policy":null}`
	if string(got) != want {
		t.Errorf("MarshalJSON = %s, want %v", got, want)
	}
}

func TestRepositoriesService_ListEnvironments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/environments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1, "environments":[{"id":1}, {"id": 2}]}`)
	})

	opt := &EnvironmentListOptions{
		ListOptions: ListOptions{
			Page:    2,
			PerPage: 2,
		},
	}
	ctx := context.Background()
	environments, _, err := client.Repositories.ListEnvironments(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListEnvironments returned error: %v", err)
	}
	want := &EnvResponse{TotalCount: Int(1), Environments: []*Environment{{ID: Int64(1)}, {ID: Int64(2)}}}
	if !cmp.Equal(environments, want) {
		t.Errorf("Repositories.ListEnvironments returned %+v, want %+v", environments, want)
	}

	const methodName = "ListEnvironments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListEnvironments(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListEnvironments(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetEnvironment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/environments/e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": 1,"name": "staging", "deployment_branch_policy": {"protected_branches": true,	"custom_branch_policies": false}, "can_admins_bypass": false}`)
	})

	ctx := context.Background()
	release, resp, err := client.Repositories.GetEnvironment(ctx, "o", "r", "e")
	if err != nil {
		t.Errorf("Repositories.GetEnvironment returned error: %v\n%v", err, resp.Body)
	}

	want := &Environment{ID: Int64(1), Name: String("staging"), DeploymentBranchPolicy: &BranchPolicy{ProtectedBranches: Bool(true), CustomBranchPolicies: Bool(false)}, CanAdminsBypass: Bool(false)}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.GetEnvironment returned %+v, want %+v", release, want)
	}

	const methodName = "GetEnvironment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetEnvironment(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetEnvironment(ctx, "o", "r", "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateEnvironment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &CreateUpdateEnvironment{
		WaitTimer: Int(30),
	}

	mux.HandleFunc("/repos/o/r/environments/e", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreateUpdateEnvironment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		want := &CreateUpdateEnvironment{WaitTimer: Int(30), CanAdminsBypass: Bool(true)}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"id": 1, "name": "staging",	"protection_rules": [{"id": 1, "type": "wait_timer", "wait_timer": 30}]}`)
	})

	ctx := context.Background()
	release, _, err := client.Repositories.CreateUpdateEnvironment(ctx, "o", "r", "e", input)
	if err != nil {
		t.Errorf("Repositories.CreateUpdateEnvironment returned error: %v", err)
	}

	want := &Environment{ID: Int64(1), Name: String("staging"), ProtectionRules: []*ProtectionRule{{ID: Int64(1), Type: String("wait_timer"), WaitTimer: Int(30)}}}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.CreateUpdateEnvironment returned %+v, want %+v", release, want)
	}

	const methodName = "CreateUpdateEnvironment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateUpdateEnvironment(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateUpdateEnvironment(ctx, "o", "r", "e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateEnvironment_noEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &CreateUpdateEnvironment{}
	callCount := 0

	mux.HandleFunc("/repos/o/r/environments/e", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreateUpdateEnvironment)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if callCount == 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			callCount++
		} else {
			want := &CreateUpdateEnvironment{}
			if !cmp.Equal(v, want) {
				t.Errorf("Request body = %+v, want %+v", v, want)
			}
			fmt.Fprint(w, `{"id": 1, "name": "staging",	"protection_rules": []}`)
		}
	})

	ctx := context.Background()
	release, _, err := client.Repositories.CreateUpdateEnvironment(ctx, "o", "r", "e", input)
	if err != nil {
		t.Errorf("Repositories.CreateUpdateEnvironment returned error: %v", err)
	}

	want := &Environment{ID: Int64(1), Name: String("staging"), ProtectionRules: []*ProtectionRule{}}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.CreateUpdateEnvironment returned %+v, want %+v", release, want)
	}
}

func TestRepositoriesService_createNewEnvNoEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &CreateUpdateEnvironment{
		DeploymentBranchPolicy: &BranchPolicy{
			ProtectedBranches:    Bool(true),
			CustomBranchPolicies: Bool(false),
		},
	}

	mux.HandleFunc("/repos/o/r/environments/e", func(w http.ResponseWriter, r *http.Request) {
		v := new(createUpdateEnvironmentNoEnterprise)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		want := &createUpdateEnvironmentNoEnterprise{
			DeploymentBranchPolicy: &BranchPolicy{
				ProtectedBranches:    Bool(true),
				CustomBranchPolicies: Bool(false),
			},
		}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"id": 1, "name": "staging",	"protection_rules": [{"id": 1, "node_id": "id", "type": "branch_policy"}], "deployment_branch_policy": {"protected_branches": true, "custom_branch_policies": false}}`)
	})

	ctx := context.Background()
	release, _, err := client.Repositories.createNewEnvNoEnterprise(ctx, "repos/o/r/environments/e", input)
	if err != nil {
		t.Errorf("Repositories.createNewEnvNoEnterprise returned error: %v", err)
	}

	want := &Environment{
		ID:   Int64(1),
		Name: String("staging"),
		ProtectionRules: []*ProtectionRule{
			{
				ID:     Int64(1),
				NodeID: String("id"),
				Type:   String("branch_policy"),
			},
		},
		DeploymentBranchPolicy: &BranchPolicy{
			ProtectedBranches:    Bool(true),
			CustomBranchPolicies: Bool(false),
		},
	}
	if !cmp.Equal(release, want) {
		t.Errorf("Repositories.createNewEnvNoEnterprise returned %+v, want %+v", release, want)
	}

	const methodName = "createNewEnvNoEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.createNewEnvNoEnterprise(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.createNewEnvNoEnterprise(ctx, "repos/o/r/environments/e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteEnvironment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/environments/e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteEnvironment(ctx, "o", "r", "e")
	if err != nil {
		t.Errorf("Repositories.DeleteEnvironment returned error: %v", err)
	}

	const methodName = "DeleteEnvironment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteEnvironment(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteEnvironment(ctx, "o", "r", "e")
	})
}

func TestRepoEnvironment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EnvResponse{}, "{}")

	repoEnv := &EnvResponse{
		TotalCount: Int(1),
		Environments: []*Environment{
			{
				Owner:           String("me"),
				Repo:            String("se"),
				EnvironmentName: String("dev"),
				WaitTimer:       Int(123),
				Reviewers: []*EnvReviewers{
					{
						Type: String("main"),
						ID:   Int64(1),
					},
					{
						Type: String("rev"),
						ID:   Int64(2),
					},
				},
				DeploymentBranchPolicy: &BranchPolicy{
					ProtectedBranches:    Bool(false),
					CustomBranchPolicies: Bool(false),
				},
				ID:        Int64(2),
				NodeID:    String("star"),
				Name:      String("eg"),
				URL:       String("https://hey.in"),
				HTMLURL:   String("htmlurl"),
				CreatedAt: &Timestamp{referenceTime},
				UpdatedAt: &Timestamp{referenceTime},
				ProtectionRules: []*ProtectionRule{
					{
						ID:        Int64(21),
						NodeID:    String("mnb"),
						Type:      String("ewq"),
						WaitTimer: Int(9090),
					},
				},
			},
		},
	}

	want := `{
		"total_count":1,
		"environments":[
		   {
			  "owner":"me",
			  "repo":"se",
			  "environment_name":"dev",
			  "wait_timer":123,
			  "reviewers":[
				 {
					"type":"main",
					"id":1
				 },
				 {
					"type":"rev",
					"id":2
				 }
			  ],
			  "deployment_branch_policy":{
				 "protected_branches":false,
				 "custom_branch_policies":false
			  },
			  "id":2,
			  "node_id":"star",
			  "name":"eg",
			  "url":"https://hey.in",
			  "html_url":"htmlurl",
			  "created_at":` + referenceTimeStr + `,
			  "updated_at":` + referenceTimeStr + `,
			  "protection_rules":[
				 {
					"id":21,
					"node_id":"mnb",
					"type":"ewq",
					"wait_timer":9090
				 }
			  ]
		   }
		]
	 }`

	testJSONMarshal(t, repoEnv, want)
}

func TestEnvReviewers_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EnvReviewers{}, "{}")

	repoEnv := &EnvReviewers{
		Type: String("main"),
		ID:   Int64(1),
	}

	want := `{
		"type":"main",
		"id":1
	}`

	testJSONMarshal(t, repoEnv, want)
}

func TestEnvironment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Environment{}, "{}")

	repoEnv := &Environment{
		Owner:           String("o"),
		Repo:            String("r"),
		EnvironmentName: String("e"),
		WaitTimer:       Int(123),
		Reviewers: []*EnvReviewers{
			{
				Type: String("main"),
				ID:   Int64(1),
			},
			{
				Type: String("rev"),
				ID:   Int64(2),
			},
		},
		DeploymentBranchPolicy: &BranchPolicy{
			ProtectedBranches:    Bool(false),
			CustomBranchPolicies: Bool(false),
		},
		ID:        Int64(2),
		NodeID:    String("star"),
		Name:      String("eg"),
		URL:       String("https://hey.in"),
		HTMLURL:   String("htmlurl"),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		ProtectionRules: []*ProtectionRule{
			{
				ID:        Int64(21),
				NodeID:    String("mnb"),
				Type:      String("ewq"),
				WaitTimer: Int(9090),
			},
		},
	}

	want := `{
		"owner":"o",
		"repo":"r",
		"environment_name":"e",
		"wait_timer":123,
		"reviewers":[
			{
				"type":"main",
				"id":1
			},
			{
				"type":"rev",
				"id":2
			}
		],
		"deployment_branch_policy":{
			"protected_branches":false,
			"custom_branch_policies":false
		},
		"id":2,
		"node_id":"star",
		"name":"eg",
		"url":"https://hey.in",
		"html_url":"htmlurl",
		"created_at":` + referenceTimeStr + `,
		"updated_at":` + referenceTimeStr + `,
		"protection_rules":[
			{
				"id":21,
				"node_id":"mnb",
				"type":"ewq",
				"wait_timer":9090
			}
		]
	}`

	testJSONMarshal(t, repoEnv, want)
}

func TestBranchPolicy_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &BranchPolicy{}, "{}")

	bp := &BranchPolicy{
		ProtectedBranches:    Bool(false),
		CustomBranchPolicies: Bool(false),
	}

	want := `{
		"protected_branches": false,
		"custom_branch_policies": false
	}`

	testJSONMarshal(t, bp, want)
}
