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

func TestOrganizationsService_GetOrganizationCustomPropertyValues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/org-properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"property_name": "team",
			"value": "core"
		},
		{
			"property_name": "level",
			"value": "gold"
		}]`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.GetOrganizationCustomPropertyValues(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetOrganizationCustomPropertyValues returned error: %v", err)
	}

	want := []*PropertyValue{
		{PropertyName: Ptr("team"), Value: Ptr("core")},
		{PropertyName: Ptr("level"), Value: Ptr("gold")},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Organizations.GetOrganizationCustomPropertyValues = %+v, want %+v", got, want)
	}

	const methodName = "GetOrganizationCustomPropertyValues"

	testBadOptions(t, methodName, func() error {
		_, _, err := client.Organizations.GetOrganizationCustomPropertyValues(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrganizationCustomPropertyValues(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrUpdateOrgCustomPropertyValues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/org-properties/values", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{}`)
	})

	ctx := t.Context()
	values := []*PropertyValue{
		{PropertyName: Ptr("team"), Value: Ptr("core")},
		{PropertyName: Ptr("level"), Value: Ptr("gold")},
	}

	props := OrganizationCustomPropertyValues{
		Properties: values,
	}
	_, err := client.Organizations.CreateOrUpdateOrgCustomPropertyValues(ctx, "o", props)
	if err != nil {
		t.Errorf("Organizations.CreateOrUpdateOrgCustomPropertyValues returned error: %v", err)
	}

	const methodName = "CreateOrUpdateOrgCustomPropertyValues"
	testBadOptions(t, methodName, func() error {
		_, err := client.Organizations.CreateOrUpdateOrgCustomPropertyValues(ctx, "\n", props)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.CreateOrUpdateOrgCustomPropertyValues(ctx, "o", props)
	})
}
