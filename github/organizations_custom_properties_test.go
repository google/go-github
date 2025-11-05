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

func TestEnterpriseService_GetEnterpriseCustomPropertySchema(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/org-properties/schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"properties": [{
				"name": "team",
				"value_type": "string",
				"description": "Team name"
			}]
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetEnterpriseCustomPropertySchema(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseCustomPropertySchema returned error: %v", err)
	}

	want := &EnterpriseCustomPropertySchema{
		Properties: []*Property{
			{
				Name:        "team",
				ValueType:   "string",
				Description: Ptr("Team name"),
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetEnterpriseCustomPropertySchema = %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseCustomPropertySchema"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetEnterpriseCustomPropertySchema(ctx, "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseCustomPropertySchema(ctx, "e")
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
		Properties: []*Property{{Name: "team"}},
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
			"name": "team",
			"value_type": "string",
			"description": "Team name"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetEnterpriseCustomProperty(ctx, "e", "prop")
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseCustomProperty returned error: %v", err)
	}

	want := &Property{
		Name:        "team",
		ValueType:   "string",
		Description: Ptr("Team name"),
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
	property := Property{Name: "team"}
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
			"properties": [{"property_name": "team", "value": "core"}]
		}]`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetEnterpriseCustomPropertyValues(ctx, "e", nil)
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseCustomPropertyValues returned error: %v", err)
	}

	want := []*CustomPropertiesValues{
		{
			OrganizationID:    Ptr(int64(1)),
			OrganizationLogin: Ptr("org1"),
			Properties: []*PropertyValue{
				{PropertyName: Ptr("team"), Value: Ptr("core")},
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetEnterpriseCustomPropertyValues = %+v, want %+v", got, want)
	}

	const methodName = "GetEnterpriseCustomPropertyValues"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetEnterpriseCustomPropertyValues(ctx, "\n", nil)
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
	values := []*PropertyValue{{PropertyName: Ptr("team"), Value: Ptr("core")}}
	orgs := []string{"org1"}

	opts := EnterpriseCustomPropertyValues{
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
