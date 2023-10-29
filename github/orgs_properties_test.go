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

func TestOrganizationsService_CreateOrUpdateCustomProperty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/schema/name", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
		"property_name": "name",
		"value_type": "single_select",
		"required": true,
		"default_value": "production",
		"description": "Prod or dev environment",
		"allowed_values":[
		  "production",
		  "development"
		]
	  }`)
	})

	ctx := context.Background()
	property, _, err := client.Organizations.CreateOrUpdateCustomProperty(ctx, "o", "name", &CustomProperty{
		ValueType:     "single_select",
		Required:      Bool(true),
		DefaultValue:  String("production"),
		Description:   String("Prod or dev environment"),
		AllowedValues: []string{"production", "development"},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrUpdateCustomProperty returned error: %v", err)
	}

	want := &CustomProperty{
		PropertyName:  String("name"),
		ValueType:     "single_select",
		Required:      Bool(true),
		DefaultValue:  String("production"),
		Description:   String("Prod or dev environment"),
		AllowedValues: []string{"production", "development"},
	}
	if !cmp.Equal(property, want) {
		t.Errorf("Organizations.CreateOrUpdateCustomProperty returned %+v, want %+v", property, want)
	}

	const methodName = "CreateOrUpdateCustomProperty"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrUpdateCustomProperty(ctx, "o", "name", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetCustomProperty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/schema/name", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		"property_name": "name",
		"value_type": "single_select",
		"required": true,
		"default_value": "production",
		"description": "Prod or dev environment",
		"allowed_values":[
		  "production",
		  "development"
		]
	  }`)
	})

	ctx := context.Background()
	property, _, err := client.Organizations.GetCustomProperty(ctx, "o", "name")
	if err != nil {
		t.Errorf("Organizations.GetCustomProperty returned error: %v", err)
	}

	want := &CustomProperty{
		PropertyName:  String("name"),
		ValueType:     "single_select",
		Required:      Bool(true),
		DefaultValue:  String("production"),
		Description:   String("Prod or dev environment"),
		AllowedValues: []string{"production", "development"},
	}
	if !cmp.Equal(property, want) {
		t.Errorf("Organizations.GetCustomProperty returned %+v, want %+v", property, want)
	}

	const methodName = "GetCustomProperty"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetCustomProperty(ctx, "o", "name")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_RemoveCustomProperty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/schema/name", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.RemoveCustomProperty(ctx, "o", "name")
	if err != nil {
		t.Errorf("Organizations.RemoveCustomProperty returned error: %v", err)
	}

	const methodName = "RemoveCustomProperty"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.RemoveCustomProperty(ctx, "0", "name")
	})
}

func TestOrganizationsService_GetAllCustomProperties(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		{
          "property_name": "name",
          "value_type": "single_select",
          "required": true,
          "default_value": "production",
          "description": "Prod or dev environment",
          "allowed_values":[
            "production",
            "development"
          ]
        },
        {
          "property_name": "service",
          "value_type": "string"
        },
        {
          "property_name": "team",
          "value_type": "string",
          "description": "Team owning the repository"
        }
        ]`)
	})

	ctx := context.Background()
	properties, _, err := client.Organizations.GetAllCustomProperties(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetAllCustomProperties returned error: %v", err)
	}

	want := []*CustomProperty{
		{
			PropertyName:  String("name"),
			ValueType:     "single_select",
			Required:      Bool(true),
			DefaultValue:  String("production"),
			Description:   String("Prod or dev environment"),
			AllowedValues: []string{"production", "development"},
		},
		{
			PropertyName: String("service"),
			ValueType:    "string",
		},
		{
			PropertyName: String("team"),
			ValueType:    "string",
			Description:  String("Team owning the repository"),
		},
	}
	if !cmp.Equal(properties, want) {
		t.Errorf("Organizations.GetAllCustomProperties returned %+v, want %+v", properties, want)
	}

	const methodName = "GetAllCustomProperties"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAllCustomProperties(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
