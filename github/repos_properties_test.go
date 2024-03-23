// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetAllCustomPropertyValues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		{
          "property_name": "environment",
          "value": "production"
        },
        {
          "property_name": "service",
          "value": "web"
        }
		]`)
	})

	ctx := context.Background()
	customPropertyValues, _, err := client.Repositories.GetAllCustomPropertyValues(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetAllCustomPropertyValues returned error: %v", err)
	}

	want := []*CustomPropertyValue{
		{
			PropertyName: "environment",
			Value:        String("production"),
		},
		{
			PropertyName: "service",
			Value:        String("web"),
		},
	}

	if !cmp.Equal(customPropertyValues, want) {
		t.Errorf("Repositories.GetAllCustomPropertyValues returned %+v, want %+v", customPropertyValues, want)
	}

	const methodName = "GetAllCustomPropertyValues"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetAllCustomPropertyValues(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCreateOrUpdateRepoCustomPropertyValues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	// Mock API endpoint
	mux.HandleFunc("/repos/o/repo/properties/values", func(w http.ResponseWriter, r *http.Request) {
		// Check HTTP method
		if r.Method != http.MethodPatch {
			t.Errorf("unexpected HTTP method: %s, expected: %s", r.Method, http.MethodPatch)
		}
		// Check request body
		var requestBody map[string][]*CustomPropertyValue
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		expectedRequestBody := map[string][]*CustomPropertyValue{
			"properties": {
				{
					PropertyName: "service",
					Value:        String("string"),
				},
			},
		}
		if !reflect.DeepEqual(requestBody, expectedRequestBody) {
			t.Errorf("unexpected request body, got: %v, want: %v", requestBody, expectedRequestBody)
		}

		// Respond with a dummy error response (status code 404)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
	})

	// Call the method
	ctx := context.Background()
	_, err := client.Repositories.CreateOrUpdateRepoCustomPropertyValues(ctx, "o", "repo", []*CustomPropertyValue{
		{
			PropertyName: "service",
			Value:        String("string"),
		},
	})
	if err == nil {
		t.Error("expected an error from CreateOrUpdateRepoCustomPropertyValues, but got nil")
	} else if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain status code 404, got: %v", err)
	}
}
