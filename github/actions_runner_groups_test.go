// Copyright 2020 The go-github AUTHORS. All rights reserved.
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

func TestActionsService_ListOrganizationRunnerGroups(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/runner-groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"runner_groups":[{"id":1,"name":"Default","visibility":"all","default":true,"runners_url":"https://api.github.com/orgs/octo-org/actions/runner_groups/1/runners","inherited":false,"allows_public_repositories":true},{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/runner_groups/2/repositories","runners_url":"https://api.github.com/orgs/octo-org/actions/runner_groups/2/runners","inherited":true,"allows_public_repositories":true},{"id":3,"name":"expensive-hardware","visibility":"private","default":false,"runners_url":"https://api.github.com/orgs/octo-org/actions/runner_groups/3/runners","inherited":false,"allows_public_repositories":true}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	groups, _, err := client.Actions.ListOrganizationRunnerGroups(ctx, "o", opts)
	if err != nil {
		t.Errorf("Actions.ListRunners returned error: %v", err)
	}

	want := &RunnerGroups{
		TotalCount: 3,
		RunnerGroups: []*RunnerGroup{
			{ID: Int64(1), Name: String("Default"), Visibility: String("all"), Default: Bool(true), RunnersURL: String("https://api.github.com/orgs/octo-org/actions/runner_groups/1/runners"), Inherited: Bool(false), AllowsPublicRepositories: Bool(true)},
			{ID: Int64(2), Name: String("octo-runner-group"), Visibility: String("selected"), Default: Bool(false), SelectedRepositoriesURL: String("https://api.github.com/orgs/octo-org/actions/runner_groups/2/repositories"), RunnersURL: String("https://api.github.com/orgs/octo-org/actions/runner_groups/2/runners"), Inherited: Bool(true), AllowsPublicRepositories: Bool(true)},
			{ID: Int64(3), Name: String("expensive-hardware"), Visibility: String("private"), Default: Bool(false), RunnersURL: String("https://api.github.com/orgs/octo-org/actions/runner_groups/3/runners"), Inherited: Bool(false), AllowsPublicRepositories: Bool(true)},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Actions.ListOrganizationRunnerGroups returned %+v, want %+v", groups, want)
	}

	const methodName = "ListOrganizationRunnerGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListOrganizationRunnerGroups(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListOrganizationRunnerGroups(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
