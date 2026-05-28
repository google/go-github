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

func TestCodeQualityService_GetSetup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/code-quality/setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{
			"state": "configured",
			"languages": ["javascript-typescript", "python"],
			"runner_type": "standard",
			"runner_label": null,
			"updated_at": "2026-01-01T00:00:00Z",
			"schedule": "weekly"
		}`)
	})

	ctx := t.Context()
	cfg, _, err := client.CodeQuality.GetSetup(ctx, "o", "r")
	if err != nil {
		t.Fatalf("CodeQuality.GetSetup returned error: %v", err)
	}

	want := &CodeQualitySetupConfiguration{
		State:      Ptr("configured"),
		Languages:  []string{"javascript-typescript", "python"},
		RunnerType: Ptr("standard"),
		UpdatedAt:  &Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Schedule:   Ptr("weekly"),
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Errorf("CodeQuality.GetSetup mismatch (-want +got):\n%v", diff)
	}

	const methodName = "GetSetup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeQuality.GetSetup(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeQuality.GetSetup(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeQualityService_UpdateSetup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &UpdateCodeQualitySetupOptions{
		State:     Ptr("configured"),
		Languages: []string{"javascript-typescript", "python", "ruby"},
	}

	mux.HandleFunc("/repos/o/r/code-quality/setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"run_id": 42,
			"run_url": "https://api.github.com/repos/octocat/hello-world/actions/runs/42"
		}`)
	})

	ctx := t.Context()
	result, _, err := client.CodeQuality.UpdateSetup(ctx, "o", "r", input)
	if err != nil {
		t.Fatalf("CodeQuality.UpdateSetup returned error: %v", err)
	}

	want := &UpdateCodeQualitySetupResponse{
		RunID:  Ptr(int64(42)),
		RunURL: Ptr("https://api.github.com/repos/octocat/hello-world/actions/runs/42"),
	}
	if diff := cmp.Diff(want, result); diff != "" {
		t.Errorf("CodeQuality.UpdateSetup mismatch (-want +got):\n%v", diff)
	}

	const methodName = "UpdateSetup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodeQuality.UpdateSetup(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodeQuality.UpdateSetup(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeQualityService_UpdateSetup_withRunnerLabel(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &UpdateCodeQualitySetupOptions{
		State:       Ptr("configured"),
		RunnerType:  Ptr("labeled"),
		RunnerLabel: Ptr("my-runner"),
		Languages:   []string{"go", "python"},
	}

	mux.HandleFunc("/repos/o/r/code-quality/setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"run_id": 99,
			"run_url": "https://api.github.com/repos/octocat/hello-world/actions/runs/99"
		}`)
	})

	ctx := t.Context()
	result, _, err := client.CodeQuality.UpdateSetup(ctx, "o", "r", input)
	if err != nil {
		t.Fatalf("CodeQuality.UpdateSetup returned error: %v", err)
	}

	want := &UpdateCodeQualitySetupResponse{
		RunID:  Ptr(int64(99)),
		RunURL: Ptr("https://api.github.com/repos/octocat/hello-world/actions/runs/99"),
	}
	if diff := cmp.Diff(want, result); diff != "" {
		t.Errorf("CodeQuality.UpdateSetup mismatch (-want +got):\n%v", diff)
	}
}

func TestCodeQualityService_UpdateSetup_notConfigured(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &UpdateCodeQualitySetupOptions{
		State: Ptr("not-configured"),
	}

	mux.HandleFunc("/repos/o/r/code-quality/setup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{}`)
	})

	ctx := t.Context()
	_, _, err := client.CodeQuality.UpdateSetup(ctx, "o", "r", input)
	if err != nil {
		t.Fatalf("CodeQuality.UpdateSetup returned error: %v", err)
	}
}

func TestCodeQualityService_GetSetup_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.CodeQuality.GetSetup(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestCodeQualityService_UpdateSetup_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.CodeQuality.UpdateSetup(ctx, "%", "r", nil)
	testURLParseError(t, err)
}
