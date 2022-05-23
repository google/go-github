// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListCustomRepoRoles(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/custom_roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "custom_roles": [{ "id": 1, "name": "Developer"}]}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListCustomRepoRoles(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.ListCustomRepoRoles returned error: %v", err)
	}

	want := &OrganizationCustomRepoRoles{TotalCount: Int(1), CustomRepoRoles: []*CustomRepoRoles{{ID: Int64(1), Name: String("Developer")}}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListCustomRepoRoles returned %+v, want %+v", apps, want)
	}

	const methodName = "ListCustomRepoRoles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCustomRepoRoles(ctx, "\no")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListCustomRepoRoles(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
