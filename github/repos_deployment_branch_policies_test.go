// Copyright 2023 The go-github AUTHORS. All rights reserved.
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
)

func TestRepositoriesService_ListDeploymentBranchPolicies(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/environments/e/deployment-branch-policies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"total_count":2, "branch_policies":[{"id":1}, {"id": 2}]}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.ListDeploymentBranchPolicies(ctx, "o", "r", "e")
	if err != nil {
		t.Errorf("Repositories.ListDeploymentBranchPolicies returned error: %v", err)
	}

	want := &DeploymentBranchPolicyResponse{
		BranchPolicies: []*DeploymentBranchPolicy{
			{ID: Int64(1)},
			{ID: Int64(2)},
		},
		TotalCount: Int(2),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.ListDeploymentBranchPolicies = %+v, want %+v", got, want)
	}

	const methodName = "ListDeploymentBranchPolicies"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListDeploymentBranchPolicies(ctx, "o", "r", "e")
		if got != nil {
			t.Errorf("got non-nil Repositories.ListDeploymentBranchPolicies response: %+v", got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetDeploymentBranchPolicy(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/environments/e/deployment-branch-policies/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetDeploymentBranchPolicy(ctx, "o", "r", "e", 1)
	if err != nil {
		t.Errorf("Repositories.GetDeploymentBranchPolicy returned error: %v", err)
	}

	want := &DeploymentBranchPolicy{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.GetDeploymentBranchPolicy = %+v, want %+v", got, want)
	}

	const methodName = "GetDeploymentBranchPolicy"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetDeploymentBranchPolicy(ctx, "o", "r", "e", 1)
		if got != nil {
			t.Errorf("got non-nil Repositories.GetDeploymentBranchPolicy response: %+v", got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateDeploymentBranchPolicy(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/environments/e/deployment-branch-policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.CreateDeploymentBranchPolicy(ctx, "o", "r", "e", &DeploymentBranchPolicyRequest{Name: String("n")})
	if err != nil {
		t.Errorf("Repositories.CreateDeploymentBranchPolicy returned error: %v", err)
	}

	want := &DeploymentBranchPolicy{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.CreateDeploymentBranchPolicy = %+v, want %+v", got, want)
	}

	const methodName = "CreateDeploymentBranchPolicy"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateDeploymentBranchPolicy(ctx, "o", "r", "e", &DeploymentBranchPolicyRequest{Name: String("n")})
		if got != nil {
			t.Errorf("got non-nil Repositories.CreateDeploymentBranchPolicy response: %+v", got)
		}
		return resp, err
	})
}

func TestRepositoriesService_UpdateDeploymentBranchPolicy(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/environments/e/deployment-branch-policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.UpdateDeploymentBranchPolicy(ctx, "o", "r", "e", 1, &DeploymentBranchPolicyRequest{Name: String("n")})
	if err != nil {
		t.Errorf("Repositories.UpdateDeploymentBranchPolicy returned error: %v", err)
	}

	want := &DeploymentBranchPolicy{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.UpdateDeploymentBranchPolicy = %+v, want %+v", got, want)
	}

	const methodName = "UpdateDeploymentBranchPolicy"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.UpdateDeploymentBranchPolicy(ctx, "o", "r", "e", 1, &DeploymentBranchPolicyRequest{Name: String("n")})
		if got != nil {
			t.Errorf("got non-nil Repositories.UpdateDeploymentBranchPolicy response: %+v", got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteDeploymentBranchPolicy(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/environments/e/deployment-branch-policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Repositories.DeleteDeploymentBranchPolicy(ctx, "o", "r", "e", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteDeploymentBranchPolicy returned error: %v", err)
	}

	const methodName = "DeleteDeploymentBranchPolicy"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteDeploymentBranchPolicy(ctx, "o", "r", "e", 1)
	})
}
