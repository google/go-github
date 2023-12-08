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

func TestOrganizationsService_CreateOrUpdateCustomProperties(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"properties":[{"property_name":"name","value_type":"single_select","required":true},{"property_name":"service","value_type":"string"}]}`+"\n")
		fmt.Fprint(w, `[
		{
          "property_name": "name",
          "value_type": "single_select",
          "required": true
        },
        {
          "property_name": "service",
          "value_type": "string"
        }
        ]`)
	})

	ctx := context.Background()
	properties, _, err := client.Organizations.CreateOrUpdateCustomProperties(ctx, "o", []*CustomProperty{
		{
			PropertyName: String("name"),
			ValueType:    "single_select",
			Required:     Bool(true),
		},
		{
			PropertyName: String("service"),
			ValueType:    "string",
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrUpdateCustomProperties returned error: %v", err)
	}

	want := []*CustomProperty{
		{
			PropertyName: String("name"),
			ValueType:    "single_select",
			Required:     Bool(true),
		},
		{
			PropertyName: String("service"),
			ValueType:    "string",
		},
	}

	if !cmp.Equal(properties, want) {
		t.Errorf("Organizations.CreateOrUpdateCustomProperties returned %+v, want %+v", properties, want)
	}

	const methodName = "CreateOrUpdateCustomProperties"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrUpdateCustomProperties(ctx, "o", nil)
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

func TestOrganizationsService_ListCustomPropertyValues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1", "per_page": "100"})
		fmt.Fprint(w, `[{
		"repository_id": 1296269,
		"repository_name": "Hello-World",
		"repository_full_name": "octocat/Hello-World",	
		"properties": [
		{
          "property_name": "environment",
          "value": "production"
        },
        {
          "property_name": "service",
          "value": "web"
        }
		]
        }]`)
	})

	ctx := context.Background()
	repoPropertyValues, _, err := client.Organizations.ListCustomPropertyValues(ctx, "o", &ListOptions{
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		t.Errorf("Organizations.ListCustomPropertyValues returned error: %v", err)
	}

	want := []*RepoCustomPropertyValue{
		{
			RepositoryID:       1296269,
			RepositoryName:     "Hello-World",
			RepositoryFullName: "octocat/Hello-World",
			Properties: []*CustomPropertyValue{
				{
					PropertyName: "environment",
					Value:        String("production"),
				},
				{
					PropertyName: "service",
					Value:        String("web"),
				},
			},
		},
	}

	if !cmp.Equal(repoPropertyValues, want) {
		t.Errorf("Organizations.ListCustomPropertyValues returned %+v, want %+v", repoPropertyValues, want)
	}

	const methodName = "ListCustomPropertyValues"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCustomPropertyValues(ctx, "\n", &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListCustomPropertyValues(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrUpdateRepoCustomPropertyValues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"repository_names":["repo"],"properties":[{"property_name":"service","value":"string"}]}`+"\n")
	})

	ctx := context.Background()
	_, err := client.Organizations.CreateOrUpdateRepoCustomPropertyValues(ctx, "o", []string{"repo"}, []*CustomPropertyValue{
		{
			PropertyName: "service",
			Value:        String("string"),
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrUpdateCustomPropertyValuesForRepos returned error: %v", err)
	}

	const methodName = "CreateOrUpdateCustomPropertyValuesForRepos"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.CreateOrUpdateRepoCustomPropertyValues(ctx, "o", nil, nil)
	})
}
