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

func TestSecretScanningService_ListCustomPatternsForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state":           "published",
			"push_protection": "enabled",
			"sort":            "created",
			"direction":       "desc",
			"page":            "2",
		})
		fmt.Fprint(w, `[
			{
				"id": 1,
				"name": "Custom pattern",
				"pattern": "[A-Z]{2}-[0-9]{4}",
				"slug": "custom-pattern",
				"state": "published",
				"push_protection_enabled": true,
				"start_delimiter": "\\b",
				"end_delimiter": "\\b",
				"must_match": ["ID-.*"],
				"must_not_match": ["TEST-.*"],
				"custom_pattern_version": "v1",
				"created_at": `+referenceTimeStr+`,
				"updated_at": `+referenceTimeStr+`
			}
		]`)
	})

	ctx := t.Context()
	opts := &SecretScanningCustomPatternListOptions{
		State:          "published",
		PushProtection: "enabled",
		Sort:           "created",
		Direction:      "desc",
		ListOptions:    ListOptions{Page: 2},
	}
	patterns, _, err := client.SecretScanning.ListCustomPatternsForRepo(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListCustomPatternsForRepo returned error: %v", err)
	}

	want := []*SecretScanningCustomPattern{
		{
			ID:                    1,
			Name:                  "Custom pattern",
			Pattern:               "[A-Z]{2}-[0-9]{4}",
			Slug:                  "custom-pattern",
			State:                 "published",
			PushProtectionEnabled: true,
			StartDelimiter:        Ptr(`\b`),
			EndDelimiter:          Ptr(`\b`),
			MustMatch:             []string{"ID-.*"},
			MustNotMatch:          []string{"TEST-.*"},
			CustomPatternVersion:  Ptr("v1"),
			CreatedAt:             &referenceTimestamp,
			UpdatedAt:             &referenceTimestamp,
		},
	}
	if !cmp.Equal(patterns, want) {
		t.Errorf("SecretScanning.ListCustomPatternsForRepo returned %+v, want %+v", patterns, want)
	}

	const methodName = "ListCustomPatternsForRepo"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListCustomPatternsForRepo(ctx, "\n", "\n", &SecretScanningCustomPatternListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListCustomPatternsForRepo(ctx, "o", "r", &SecretScanningCustomPatternListOptions{})
		return resp, err
	})
}

func TestSecretScanningService_CreateCustomPatternsForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningCreateCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternRequest{
			{
				Name:    "Custom pattern",
				Pattern: "[A-Z]{2}-[0-9]{4}",
			},
		},
	}

	mux.HandleFunc("/repos/o/r/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"created_patterns": [
				{
					"id": 1,
					"name": "Custom pattern",
					"pattern": "[A-Z]{2}-[0-9]{4}",
					"slug": "custom-pattern",
					"state": "published",
					"push_protection_enabled": false
				}
			]
		}`)
	})

	ctx := t.Context()
	result, _, err := client.SecretScanning.CreateCustomPatternsForRepo(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("SecretScanning.CreateCustomPatternsForRepo returned error: %v", err)
	}

	want := &SecretScanningCreateCustomPatternsResponse{
		CreatedPatterns: []*SecretScanningCustomPattern{
			{
				ID:                    1,
				Name:                  "Custom pattern",
				Pattern:               "[A-Z]{2}-[0-9]{4}",
				Slug:                  "custom-pattern",
				State:                 "published",
				PushProtectionEnabled: false,
			},
		},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("SecretScanning.CreateCustomPatternsForRepo returned %+v, want %+v", result, want)
	}

	const methodName = "CreateCustomPatternsForRepo"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.CreateCustomPatternsForRepo(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.CreateCustomPatternsForRepo(ctx, "o", "r", input)
		return resp, err
	})
}

func TestSecretScanningService_UpdateCustomPatternForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningUpdateCustomPatternRequest{
		Pattern:              Ptr("[A-Z]{3}-[0-9]{4}"),
		CustomPatternVersion: Ptr("v1"),
	}

	mux.HandleFunc("/repos/o/r/secret-scanning/custom-patterns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"id": 1,
			"name": "Custom pattern",
			"pattern": "[A-Z]{3}-[0-9]{4}",
			"slug": "custom-pattern",
			"state": "published",
			"push_protection_enabled": false,
			"custom_pattern_version": "v2"
		}`)
	})

	ctx := t.Context()
	pattern, _, err := client.SecretScanning.UpdateCustomPatternForRepo(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("SecretScanning.UpdateCustomPatternForRepo returned error: %v", err)
	}

	want := &SecretScanningCustomPattern{
		ID:                    1,
		Name:                  "Custom pattern",
		Pattern:               "[A-Z]{3}-[0-9]{4}",
		Slug:                  "custom-pattern",
		State:                 "published",
		PushProtectionEnabled: false,
		CustomPatternVersion:  Ptr("v2"),
	}
	if !cmp.Equal(pattern, want) {
		t.Errorf("SecretScanning.UpdateCustomPatternForRepo returned %+v, want %+v", pattern, want)
	}

	const methodName = "UpdateCustomPatternForRepo"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdateCustomPatternForRepo(ctx, "\n", "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdateCustomPatternForRepo(ctx, "o", "r", 1, input)
		return resp, err
	})
}

func TestSecretScanningService_DeleteCustomPatternsForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningDeleteCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternToDelete{
			{PatternID: 1},
		},
	}

	mux.HandleFunc("/repos/o/r/secret-scanning/custom-patterns", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testJSONBody(t, r, input)
	})

	ctx := t.Context()
	_, err := client.SecretScanning.DeleteCustomPatternsForRepo(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("SecretScanning.DeleteCustomPatternsForRepo returned error: %v", err)
	}

	const methodName = "DeleteCustomPatternsForRepo"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecretScanning.DeleteCustomPatternsForRepo(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SecretScanning.DeleteCustomPatternsForRepo(ctx, "o", "r", input)
	})
}

func TestSecretScanningService_ListCustomPatternsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state":           "published",
			"push_protection": "enabled",
			"sort":            "created",
			"direction":       "desc",
			"page":            "2",
		})
		fmt.Fprint(w, `[
			{
				"id": 1,
				"name": "Custom pattern",
				"pattern": "[A-Z]{2}-[0-9]{4}",
				"slug": "custom-pattern",
				"state": "published",
				"push_protection_enabled": true,
				"start_delimiter": "\\b",
				"end_delimiter": "\\b",
				"must_match": ["ID-.*"],
				"must_not_match": ["TEST-.*"],
				"custom_pattern_version": "v1",
				"created_at": `+referenceTimeStr+`,
				"updated_at": `+referenceTimeStr+`
			}
		]`)
	})

	ctx := t.Context()
	opts := &SecretScanningCustomPatternListOptions{
		State:          "published",
		PushProtection: "enabled",
		Sort:           "created",
		Direction:      "desc",
		ListOptions:    ListOptions{Page: 2},
	}
	patterns, _, err := client.SecretScanning.ListCustomPatternsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListCustomPatternsForOrg returned error: %v", err)
	}

	want := []*SecretScanningCustomPattern{
		{
			ID:                    1,
			Name:                  "Custom pattern",
			Pattern:               "[A-Z]{2}-[0-9]{4}",
			Slug:                  "custom-pattern",
			State:                 "published",
			PushProtectionEnabled: true,
			StartDelimiter:        Ptr(`\b`),
			EndDelimiter:          Ptr(`\b`),
			MustMatch:             []string{"ID-.*"},
			MustNotMatch:          []string{"TEST-.*"},
			CustomPatternVersion:  Ptr("v1"),
			CreatedAt:             &referenceTimestamp,
			UpdatedAt:             &referenceTimestamp,
		},
	}
	if !cmp.Equal(patterns, want) {
		t.Errorf("SecretScanning.ListCustomPatternsForOrg returned %+v, want %+v", patterns, want)
	}

	const methodName = "ListCustomPatternsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListCustomPatternsForOrg(ctx, "\n", &SecretScanningCustomPatternListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListCustomPatternsForOrg(ctx, "o", &SecretScanningCustomPatternListOptions{})
		return resp, err
	})
}

func TestSecretScanningService_CreateCustomPatternsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningCreateCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternRequest{
			{
				Name:    "Custom pattern",
				Pattern: "[A-Z]{2}-[0-9]{4}",
			},
		},
	}

	mux.HandleFunc("/orgs/o/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"created_patterns": [
				{
					"id": 1,
					"name": "Custom pattern",
					"pattern": "[A-Z]{2}-[0-9]{4}",
					"slug": "custom-pattern",
					"state": "published",
					"push_protection_enabled": false
				}
			]
		}`)
	})

	ctx := t.Context()
	result, _, err := client.SecretScanning.CreateCustomPatternsForOrg(ctx, "o", input)
	if err != nil {
		t.Errorf("SecretScanning.CreateCustomPatternsForOrg returned error: %v", err)
	}

	want := &SecretScanningCreateCustomPatternsResponse{
		CreatedPatterns: []*SecretScanningCustomPattern{
			{
				ID:                    1,
				Name:                  "Custom pattern",
				Pattern:               "[A-Z]{2}-[0-9]{4}",
				Slug:                  "custom-pattern",
				State:                 "published",
				PushProtectionEnabled: false,
			},
		},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("SecretScanning.CreateCustomPatternsForOrg returned %+v, want %+v", result, want)
	}

	const methodName = "CreateCustomPatternsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.CreateCustomPatternsForOrg(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.CreateCustomPatternsForOrg(ctx, "o", input)
		return resp, err
	})
}

func TestSecretScanningService_UpdateCustomPatternForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningUpdateCustomPatternRequest{
		Pattern:              Ptr("[A-Z]{3}-[0-9]{4}"),
		CustomPatternVersion: Ptr("v1"),
	}

	mux.HandleFunc("/orgs/o/secret-scanning/custom-patterns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"id": 1,
			"name": "Custom pattern",
			"pattern": "[A-Z]{3}-[0-9]{4}",
			"slug": "custom-pattern",
			"state": "published",
			"push_protection_enabled": false,
			"custom_pattern_version": "v2"
		}`)
	})

	ctx := t.Context()
	pattern, _, err := client.SecretScanning.UpdateCustomPatternForOrg(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("SecretScanning.UpdateCustomPatternForOrg returned error: %v", err)
	}

	want := &SecretScanningCustomPattern{
		ID:                    1,
		Name:                  "Custom pattern",
		Pattern:               "[A-Z]{3}-[0-9]{4}",
		Slug:                  "custom-pattern",
		State:                 "published",
		PushProtectionEnabled: false,
		CustomPatternVersion:  Ptr("v2"),
	}
	if !cmp.Equal(pattern, want) {
		t.Errorf("SecretScanning.UpdateCustomPatternForOrg returned %+v, want %+v", pattern, want)
	}

	const methodName = "UpdateCustomPatternForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdateCustomPatternForOrg(ctx, "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdateCustomPatternForOrg(ctx, "o", 1, input)
		return resp, err
	})
}

func TestSecretScanningService_DeleteCustomPatternsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningDeleteCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternToDelete{
			{PatternID: 1},
		},
	}

	mux.HandleFunc("/orgs/o/secret-scanning/custom-patterns", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testJSONBody(t, r, input)
	})

	ctx := t.Context()
	_, err := client.SecretScanning.DeleteCustomPatternsForOrg(ctx, "o", input)
	if err != nil {
		t.Errorf("SecretScanning.DeleteCustomPatternsForOrg returned error: %v", err)
	}

	const methodName = "DeleteCustomPatternsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecretScanning.DeleteCustomPatternsForOrg(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SecretScanning.DeleteCustomPatternsForOrg(ctx, "o", input)
	})
}

func TestSecretScanningService_ListCustomPatternsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state":           "published",
			"push_protection": "enabled",
			"sort":            "created",
			"direction":       "desc",
			"page":            "2",
		})
		fmt.Fprint(w, `[
			{
				"id": 1,
				"name": "Custom pattern",
				"pattern": "[A-Z]{2}-[0-9]{4}",
				"slug": "custom-pattern",
				"state": "published",
				"push_protection_enabled": true,
				"start_delimiter": "\\b",
				"end_delimiter": "\\b",
				"must_match": ["ID-.*"],
				"must_not_match": ["TEST-.*"],
				"custom_pattern_version": "v1",
				"created_at": `+referenceTimeStr+`,
				"updated_at": `+referenceTimeStr+`
			}
		]`)
	})

	ctx := t.Context()
	opts := &SecretScanningCustomPatternListOptions{
		State:          "published",
		PushProtection: "enabled",
		Sort:           "created",
		Direction:      "desc",
		ListOptions:    ListOptions{Page: 2},
	}
	patterns, _, err := client.SecretScanning.ListCustomPatternsForEnterprise(ctx, "e", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListCustomPatternsForEnterprise returned error: %v", err)
	}

	want := []*SecretScanningCustomPattern{
		{
			ID:                    1,
			Name:                  "Custom pattern",
			Pattern:               "[A-Z]{2}-[0-9]{4}",
			Slug:                  "custom-pattern",
			State:                 "published",
			PushProtectionEnabled: true,
			StartDelimiter:        Ptr(`\b`),
			EndDelimiter:          Ptr(`\b`),
			MustMatch:             []string{"ID-.*"},
			MustNotMatch:          []string{"TEST-.*"},
			CustomPatternVersion:  Ptr("v1"),
			CreatedAt:             &referenceTimestamp,
			UpdatedAt:             &referenceTimestamp,
		},
	}
	if !cmp.Equal(patterns, want) {
		t.Errorf("SecretScanning.ListCustomPatternsForEnterprise returned %+v, want %+v", patterns, want)
	}

	const methodName = "ListCustomPatternsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListCustomPatternsForEnterprise(ctx, "\n", &SecretScanningCustomPatternListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListCustomPatternsForEnterprise(ctx, "e", &SecretScanningCustomPatternListOptions{})
		return resp, err
	})
}

func TestSecretScanningService_CreateCustomPatternsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningCreateCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternRequest{
			{
				Name:    "Custom pattern",
				Pattern: "[A-Z]{2}-[0-9]{4}",
			},
		},
	}

	mux.HandleFunc("/enterprises/e/secret-scanning/custom-patterns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"created_patterns": [
				{
					"id": 1,
					"name": "Custom pattern",
					"pattern": "[A-Z]{2}-[0-9]{4}",
					"slug": "custom-pattern",
					"state": "published",
					"push_protection_enabled": false
				}
			]
		}`)
	})

	ctx := t.Context()
	result, _, err := client.SecretScanning.CreateCustomPatternsForEnterprise(ctx, "e", input)
	if err != nil {
		t.Errorf("SecretScanning.CreateCustomPatternsForEnterprise returned error: %v", err)
	}

	want := &SecretScanningCreateCustomPatternsResponse{
		CreatedPatterns: []*SecretScanningCustomPattern{
			{
				ID:                    1,
				Name:                  "Custom pattern",
				Pattern:               "[A-Z]{2}-[0-9]{4}",
				Slug:                  "custom-pattern",
				State:                 "published",
				PushProtectionEnabled: false,
			},
		},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("SecretScanning.CreateCustomPatternsForEnterprise returned %+v, want %+v", result, want)
	}

	const methodName = "CreateCustomPatternsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.CreateCustomPatternsForEnterprise(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.CreateCustomPatternsForEnterprise(ctx, "e", input)
		return resp, err
	})
}

func TestSecretScanningService_UpdateCustomPatternForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningUpdateCustomPatternRequest{
		Pattern:              Ptr("[A-Z]{3}-[0-9]{4}"),
		CustomPatternVersion: Ptr("v1"),
	}

	mux.HandleFunc("/enterprises/e/secret-scanning/custom-patterns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{
			"id": 1,
			"name": "Custom pattern",
			"pattern": "[A-Z]{3}-[0-9]{4}",
			"slug": "custom-pattern",
			"state": "published",
			"push_protection_enabled": false,
			"custom_pattern_version": "v2"
		}`)
	})

	ctx := t.Context()
	pattern, _, err := client.SecretScanning.UpdateCustomPatternForEnterprise(ctx, "e", 1, input)
	if err != nil {
		t.Errorf("SecretScanning.UpdateCustomPatternForEnterprise returned error: %v", err)
	}

	want := &SecretScanningCustomPattern{
		ID:                    1,
		Name:                  "Custom pattern",
		Pattern:               "[A-Z]{3}-[0-9]{4}",
		Slug:                  "custom-pattern",
		State:                 "published",
		PushProtectionEnabled: false,
		CustomPatternVersion:  Ptr("v2"),
	}
	if !cmp.Equal(pattern, want) {
		t.Errorf("SecretScanning.UpdateCustomPatternForEnterprise returned %+v, want %+v", pattern, want)
	}

	const methodName = "UpdateCustomPatternForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdateCustomPatternForEnterprise(ctx, "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdateCustomPatternForEnterprise(ctx, "e", 1, input)
		return resp, err
	})
}

func TestSecretScanningService_DeleteCustomPatternsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := SecretScanningDeleteCustomPatternsRequest{
		Patterns: []*SecretScanningCustomPatternToDelete{
			{PatternID: 1},
		},
	}

	mux.HandleFunc("/enterprises/e/secret-scanning/custom-patterns", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testJSONBody(t, r, input)
	})

	ctx := t.Context()
	_, err := client.SecretScanning.DeleteCustomPatternsForEnterprise(ctx, "e", input)
	if err != nil {
		t.Errorf("SecretScanning.DeleteCustomPatternsForEnterprise returned error: %v", err)
	}

	const methodName = "DeleteCustomPatternsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecretScanning.DeleteCustomPatternsForEnterprise(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SecretScanning.DeleteCustomPatternsForEnterprise(ctx, "e", input)
	})
}
