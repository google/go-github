// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_GetOrganizationCustomPropertySchema(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"properties": [{
				"property_name": "team",
				"value_type": "string",
				"description": "Team name"
			}]
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetOrganizationCustomPropertySchema(ctx, "e")
	if err != nil {
		t.Fatalf("Enterprise.GetOrganizationCustomPropertySchema returned error: %v", err)
	}

	want := &EnterpriseCustomPropertySchema{
		Properties: []*CustomProperty{
			{
				PropertyName: Ptr("team"),
				ValueType:    "string",
				Description:  Ptr("Team name"),
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetOrganizationCustomPropertySchema = %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationCustomPropertySchema"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetOrganizationCustomPropertySchema(ctx, "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetOrganizationCustomPropertySchema(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateOrUpdateEnterpriseCustomPropertySchema(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{}`)
	})

	ctx := t.Context()
	schema := EnterpriseCustomPropertySchema{
		Properties: []*CustomProperty{{PropertyName: Ptr("team")}},
	}
	_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomPropertySchema(ctx, "e", schema)
	if err != nil {
		t.Errorf("Enterprise.CreateOrUpdateEnterpriseCustomPropertySchema returned error: %v", err)
	}

	const methodName = "CreateOrUpdateEnterpriseCustomPropertySchema"
	testBadOptions(t, methodName, func() error {
		_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomPropertySchema(ctx, "\n", schema)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.CreateOrUpdateEnterpriseCustomPropertySchema(ctx, "e", schema)
	})
}

func TestEnterpriseService_GetEnterpriseCustomProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema/prop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"property_name": "team",
			"value_type": "string",
			"description": "Team name"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetEnterpriseCustomProperty(ctx, "e", "prop")
	if err != nil {
		t.Fatalf("Enterprise.GetEnterpriseCustomProperty returned error: %v", err)
	}

	want := &CustomProperty{
		PropertyName: Ptr("team"),
		ValueType:    "string",
		Description:  Ptr("Team name"),
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetEnterpriseCustomProperty = %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseCustomProperty"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetEnterpriseCustomProperty(ctx, "\n", "prop")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseCustomProperty(ctx, "e", "prop")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateOrUpdateEnterpriseCustomProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema/prop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{}`)
	})

	ctx := t.Context()
	property := CustomProperty{PropertyName: Ptr("team")}
	_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomProperty(ctx, "e", "prop", property)
	if err != nil {
		t.Errorf("Enterprise.CreateOrUpdateEnterpriseCustomProperty returned error: %v", err)
	}

	const methodName = "CreateOrUpdateEnterpriseCustomProperty"
	testBadOptions(t, methodName, func() error {
		_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomProperty(ctx, "\n", "prop", property)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.CreateOrUpdateEnterpriseCustomProperty(ctx, "e", "prop", property)
	})
}

func TestEnterpriseService_DeleteEnterpriseCustomProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema/prop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Enterprise.DeleteEnterpriseCustomProperty(ctx, "e", "prop")
	if err != nil {
		t.Errorf("Enterprise.DeleteEnterpriseCustomProperty returned error: %v", err)
	}

	const methodName = "DeleteEnterpriseCustomProperty"
	testBadOptions(t, methodName, func() error {
		_, err := client.Enterprise.DeleteEnterpriseCustomProperty(ctx, "\n", "prop")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteEnterpriseCustomProperty(ctx, "e", "prop")
	})
}

func TestEnterpriseService_GetEnterpriseCustomPropertyValues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"organization_id": 1,
			"organization_login": "org1",
			"properties": [
				{"property_name": "team", "value": "core"}
			]
		}]`)
	})

	ctx := t.Context()
	opts := &ListOptions{Page: 1, PerPage: 10}
	got, _, err := client.Enterprise.GetEnterpriseCustomPropertyValues(ctx, "e", opts)
	if err != nil {
		t.Fatalf("Enterprise.GetEnterpriseCustomPropertyValues returned error: %v", err)
	}

	want := []*EnterpriseCustomPropertiesValues{
		{
			OrganizationID:    Ptr(int64(1)),
			OrganizationLogin: Ptr("org1"),
			Properties: []*CustomPropertyValue{
				{PropertyName: "team", Value: "core"},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetEnterpriseCustomPropertyValues = %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseCustomPropertyValues"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetEnterpriseCustomPropertyValues(ctx, "\n", opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseCustomPropertyValues(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateOrUpdateEnterpriseCustomPropertyValues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{}`)
	})

	ctx := t.Context()
	values := []*CustomPropertyValue{{PropertyName: "team", Value: Ptr("core")}}
	orgs := []string{"org1"}

	opts := EnterpriseCustomPropertyValuesRequest{
		OrganizationLogin: orgs,
		Properties:        values,
	}
	_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomPropertyValues(ctx, "e", opts)
	if err != nil {
		t.Errorf("Enterprise.CreateOrUpdateEnterpriseCustomPropertyValues returned error: %v", err)
	}

	const methodName = "CreateOrUpdateEnterpriseCustomPropertyValues"
	testBadOptions(t, methodName, func() error {
		_, err := client.Enterprise.CreateOrUpdateEnterpriseCustomPropertyValues(ctx, "\n", opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.CreateOrUpdateEnterpriseCustomPropertyValues(ctx, "e", opts)
	})
}
