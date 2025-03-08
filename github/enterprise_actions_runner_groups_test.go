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

func TestEnterpriseService_ListRunnerGroups(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"runner_groups":[{"id":1,"name":"Default","visibility":"all","default":true,"runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/1/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":true,"selected_workflows":["a","b"]},{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_organizations_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations","runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners","inherited":true,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]},{"id":3,"name":"expensive-hardware","visibility":"private","default":false,"runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/3/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]}]}`)
	})

	opts := &ListEnterpriseRunnerGroupOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	groups, _, err := client.Enterprise.ListRunnerGroups(ctx, "o", opts)
	if err != nil {
		t.Errorf("Enterprise.ListRunnerGroups returned error: %v", err)
	}

	want := &EnterpriseRunnerGroups{
		TotalCount: Ptr(3),
		RunnerGroups: []*EnterpriseRunnerGroup{
			{ID: Ptr(int64(1)), Name: Ptr("Default"), Visibility: Ptr("all"), Default: Ptr(true), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/1/runners"), Inherited: Ptr(false), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(true), SelectedWorkflows: []string{"a", "b"}},
			{ID: Ptr(int64(2)), Name: Ptr("octo-runner-group"), Visibility: Ptr("selected"), Default: Ptr(false), SelectedOrganizationsURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations"), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners"), Inherited: Ptr(true), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(false), SelectedWorkflows: []string{}},
			{ID: Ptr(int64(3)), Name: Ptr("expensive-hardware"), Visibility: Ptr("private"), Default: Ptr(false), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/3/runners"), Inherited: Ptr(false), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(false), SelectedWorkflows: []string{}},
		},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Enterprise.ListRunnerGroups returned %+v, want %+v", groups, want)
	}

	const methodName = "ListRunnerGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRunnerGroups(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListRunnerGroups(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListRunnerGroupsVisibleToOrganization(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2", "visible_to_organization": "github"})
		fmt.Fprint(w, `{"total_count":3,"runner_groups":[{"id":1,"name":"Default","visibility":"all","default":true,"runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/1/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]},{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_organizations_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations","runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners","inherited":true,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]},{"id":3,"name":"expensive-hardware","visibility":"private","default":false,"runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/3/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]}]}`)
	})

	opts := &ListEnterpriseRunnerGroupOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}, VisibleToOrganization: "github"}
	ctx := context.Background()
	groups, _, err := client.Enterprise.ListRunnerGroups(ctx, "o", opts)
	if err != nil {
		t.Errorf("Enterprise.ListRunnerGroups returned error: %v", err)
	}

	want := &EnterpriseRunnerGroups{
		TotalCount: Ptr(3),
		RunnerGroups: []*EnterpriseRunnerGroup{
			{ID: Ptr(int64(1)), Name: Ptr("Default"), Visibility: Ptr("all"), Default: Ptr(true), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/1/runners"), Inherited: Ptr(false), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(false), SelectedWorkflows: []string{}},
			{ID: Ptr(int64(2)), Name: Ptr("octo-runner-group"), Visibility: Ptr("selected"), Default: Ptr(false), SelectedOrganizationsURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations"), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners"), Inherited: Ptr(true), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(false), SelectedWorkflows: []string{}},
			{ID: Ptr(int64(3)), Name: Ptr("expensive-hardware"), Visibility: Ptr("private"), Default: Ptr(false), RunnersURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/3/runners"), Inherited: Ptr(false), AllowsPublicRepositories: Ptr(true), RestrictedToWorkflows: Ptr(false), SelectedWorkflows: []string{}},
		},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Enterprise.ListRunnerGroups returned %+v, want %+v", groups, want)
	}

	const methodName = "ListRunnerGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRunnerGroups(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListRunnerGroups(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_organizations_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations","runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]}`)
	})

	ctx := context.Background()
	group, _, err := client.Enterprise.GetEnterpriseRunnerGroup(ctx, "o", 2)
	if err != nil {
		t.Errorf("Enterprise.GetRunnerGroup returned error: %v", err)
	}

	want := &EnterpriseRunnerGroup{
		ID:                       Ptr(int64(2)),
		Name:                     Ptr("octo-runner-group"),
		Visibility:               Ptr("selected"),
		Default:                  Ptr(false),
		SelectedOrganizationsURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations"),
		RunnersURL:               Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners"),
		Inherited:                Ptr(false),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}

	if !cmp.Equal(group, want) {
		t.Errorf("Enterprise.GetRunnerGroup returned %+v, want %+v", group, want)
	}

	const methodName = "GetRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetEnterpriseRunnerGroup(ctx, "\n", 2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseRunnerGroup(ctx, "o", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.DeleteEnterpriseRunnerGroup(ctx, "o", 2)
	if err != nil {
		t.Errorf("Enterprise.DeleteRunnerGroup returned error: %v", err)
	}

	const methodName = "DeleteRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.DeleteEnterpriseRunnerGroup(ctx, "\n", 2)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteEnterpriseRunnerGroup(ctx, "o", 2)
	})
}

func TestEnterpriseService_CreateRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_organizations_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations","runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]}`)
	})

	ctx := context.Background()
	req := CreateEnterpriseRunnerGroupRequest{
		Name:                     Ptr("octo-runner-group"),
		Visibility:               Ptr("selected"),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}
	group, _, err := client.Enterprise.CreateEnterpriseRunnerGroup(ctx, "o", req)
	if err != nil {
		t.Errorf("Enterprise.CreateRunnerGroup returned error: %v", err)
	}

	want := &EnterpriseRunnerGroup{
		ID:                       Ptr(int64(2)),
		Name:                     Ptr("octo-runner-group"),
		Visibility:               Ptr("selected"),
		Default:                  Ptr(false),
		SelectedOrganizationsURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations"),
		RunnersURL:               Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners"),
		Inherited:                Ptr(false),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}

	if !cmp.Equal(group, want) {
		t.Errorf("Enterprise.CreateRunnerGroup returned %+v, want %+v", group, want)
	}

	const methodName = "CreateRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.CreateEnterpriseRunnerGroup(ctx, "\n", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseRunnerGroup(ctx, "o", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":2,"name":"octo-runner-group","visibility":"selected","default":false,"selected_organizations_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations","runners_url":"https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners","inherited":false,"allows_public_repositories":true,"restricted_to_workflows":false,"selected_workflows":[]}`)
	})

	ctx := context.Background()
	req := UpdateEnterpriseRunnerGroupRequest{
		Name:                     Ptr("octo-runner-group"),
		Visibility:               Ptr("selected"),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}
	group, _, err := client.Enterprise.UpdateEnterpriseRunnerGroup(ctx, "o", 2, req)
	if err != nil {
		t.Errorf("Enterprise.UpdateRunnerGroup returned error: %v", err)
	}

	want := &EnterpriseRunnerGroup{
		ID:                       Ptr(int64(2)),
		Name:                     Ptr("octo-runner-group"),
		Visibility:               Ptr("selected"),
		Default:                  Ptr(false),
		SelectedOrganizationsURL: Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/organizations"),
		RunnersURL:               Ptr("https://api.github.com/enterprises/octo-enterprise/actions/runner_groups/2/runners"),
		Inherited:                Ptr(false),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}

	if !cmp.Equal(group, want) {
		t.Errorf("Enterprise.UpdateRunnerGroup returned %+v, want %+v", group, want)
	}

	const methodName = "UpdateRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.UpdateEnterpriseRunnerGroup(ctx, "\n", 2, req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateEnterpriseRunnerGroup(ctx, "o", 2, req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListOrganizationAccessRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "1", "page": "1"})
		fmt.Fprint(w, `{"total_count": 1, "organizations": [{"id": 43, "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", "name": "Hello-World", "login": "octocat"}]}`)
	})

	ctx := context.Background()
	opts := &ListOptions{Page: 1, PerPage: 1}
	groups, _, err := client.Enterprise.ListOrganizationAccessRunnerGroup(ctx, "o", 2, opts)
	if err != nil {
		t.Errorf("Enterprise.ListOrganizationAccessRunnerGroup returned error: %v", err)
	}

	want := &ListOrganizations{
		TotalCount: Ptr(1),
		Organizations: []*Organization{
			{ID: Ptr(int64(43)), NodeID: Ptr("MDEwOlJlcG9zaXRvcnkxMjk2MjY5"), Name: Ptr("Hello-World"), Login: Ptr("octocat")},
		},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Enterprise.ListOrganizationAccessRunnerGroup returned %+v, want %+v", groups, want)
	}

	const methodName = "ListOrganizationAccessRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListOrganizationAccessRunnerGroup(ctx, "\n", 2, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListOrganizationAccessRunnerGroup(ctx, "o", 2, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_SetOrganizationAccessRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	req := SetOrgAccessRunnerGroupRequest{
		SelectedOrganizationIDs: []int64{
			1,
			2,
		},
	}

	ctx := context.Background()
	_, err := client.Enterprise.SetOrganizationAccessRunnerGroup(ctx, "o", 2, req)
	if err != nil {
		t.Errorf("Enterprise.SetOrganizationAccessRunnerGroup returned error: %v", err)
	}

	const methodName = "SetRepositoryAccessRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.SetOrganizationAccessRunnerGroup(ctx, "\n", 2, req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.SetOrganizationAccessRunnerGroup(ctx, "o", 2, req)
	})
}

func TestEnterpriseService_AddOrganizationAccessRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/organizations/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	ctx := context.Background()
	_, err := client.Enterprise.AddOrganizationAccessRunnerGroup(ctx, "o", 2, 42)
	if err != nil {
		t.Errorf("Enterprise.AddOrganizationAccessRunnerGroup returned error: %v", err)
	}

	const methodName = "AddOrganizationAccessRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.AddOrganizationAccessRunnerGroup(ctx, "\n", 2, 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.AddOrganizationAccessRunnerGroup(ctx, "o", 2, 42)
	})
}

func TestEnterpriseService_RemoveOrganizationAccessRunnerGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/organizations/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.RemoveOrganizationAccessRunnerGroup(ctx, "o", 2, 42)
	if err != nil {
		t.Errorf("Enterprise.RemoveOrganizationAccessRunnerGroup returned error: %v", err)
	}

	const methodName = "RemoveOrganizationAccessRunnerGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.RemoveOrganizationAccessRunnerGroup(ctx, "\n", 2, 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.RemoveOrganizationAccessRunnerGroup(ctx, "o", 2, 42)
	})
}

func TestEnterpriseService_ListEnterpriseRunnerGroupRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":2,"runners":[{"id":23,"name":"MBP","os":"macos","status":"online"},{"id":24,"name":"iMac","os":"macos","status":"offline"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	runners, _, err := client.Enterprise.ListRunnerGroupRunners(ctx, "o", 2, opts)
	if err != nil {
		t.Errorf("Enterprise.ListEnterpriseRunnerGroupRunners returned error: %v", err)
	}

	want := &Runners{
		TotalCount: 2,
		Runners: []*Runner{
			{ID: Ptr(int64(23)), Name: Ptr("MBP"), OS: Ptr("macos"), Status: Ptr("online")},
			{ID: Ptr(int64(24)), Name: Ptr("iMac"), OS: Ptr("macos"), Status: Ptr("offline")},
		},
	}
	if !cmp.Equal(runners, want) {
		t.Errorf("Enterprise.ListEnterpriseRunnerGroupRunners returned %+v, want %+v", runners, want)
	}

	const methodName = "ListEnterpriseRunnerGroupRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRunnerGroupRunners(ctx, "\n", 2, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListRunnerGroupRunners(ctx, "o", 2, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_SetEnterpriseRunnerGroupRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	req := SetRunnerGroupRunnersRequest{
		Runners: []int64{
			1,
			2,
		},
	}

	ctx := context.Background()
	_, err := client.Enterprise.SetRunnerGroupRunners(ctx, "o", 2, req)
	if err != nil {
		t.Errorf("Enterprise.SetEnterpriseRunnerGroupRunners returned error: %v", err)
	}

	const methodName = "SetEnterpriseRunnerGroupRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.SetRunnerGroupRunners(ctx, "\n", 2, req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.SetRunnerGroupRunners(ctx, "o", 2, req)
	})
}

func TestEnterpriseService_AddEnterpriseRunnerGroupRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/runners/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	ctx := context.Background()
	_, err := client.Enterprise.AddRunnerGroupRunners(ctx, "o", 2, 42)
	if err != nil {
		t.Errorf("Enterprise.AddEnterpriseRunnerGroupRunners returned error: %v", err)
	}

	const methodName = "AddEnterpriseRunnerGroupRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.AddRunnerGroupRunners(ctx, "\n", 2, 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.AddRunnerGroupRunners(ctx, "o", 2, 42)
	})
}

func TestEnterpriseService_RemoveEnterpriseRunnerGroupRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runner-groups/2/runners/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.RemoveRunnerGroupRunners(ctx, "o", 2, 42)
	if err != nil {
		t.Errorf("Enterprise.RemoveEnterpriseRunnerGroupRunners returned error: %v", err)
	}

	const methodName = "RemoveEnterpriseRunnerGroupRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.RemoveRunnerGroupRunners(ctx, "\n", 2, 42)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.RemoveRunnerGroupRunners(ctx, "o", 2, 42)
	})
}

func TestEnterpriseRunnerGroup_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EnterpriseRunnerGroup{}, "{}")

	u := &EnterpriseRunnerGroup{
		ID:                       Ptr(int64(1)),
		Name:                     Ptr("n"),
		Visibility:               Ptr("v"),
		Default:                  Ptr(true),
		SelectedOrganizationsURL: Ptr("s"),
		RunnersURL:               Ptr("r"),
		Inherited:                Ptr(true),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}

	want := `{
		"id": 1,
		"name": "n",
		"visibility": "v",
		"default": true,
		"selected_organizations_url": "s",
		"runners_url": "r",
		"inherited": true,
		"allows_public_repositories": true,
		"restricted_to_workflows": false,
		"selected_workflows": []
	}`

	testJSONMarshal(t, u, want)
}

func TestEnterpriseRunnerGroups_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &EnterpriseRunnerGroups{}, "{}")

	u := &EnterpriseRunnerGroups{
		TotalCount: Ptr(1),
		RunnerGroups: []*EnterpriseRunnerGroup{
			{
				ID:                       Ptr(int64(1)),
				Name:                     Ptr("n"),
				Visibility:               Ptr("v"),
				Default:                  Ptr(true),
				SelectedOrganizationsURL: Ptr("s"),
				RunnersURL:               Ptr("r"),
				Inherited:                Ptr(true),
				AllowsPublicRepositories: Ptr(true),
				RestrictedToWorkflows:    Ptr(false),
				SelectedWorkflows:        []string{},
			},
		},
	}

	want := `{
		"total_count": 1,
		"runner_groups": [{
			"id": 1,
			"name": "n",
			"visibility": "v",
			"default": true,
			"selected_organizations_url": "s",
			"runners_url": "r",
			"inherited": true,
			"allows_public_repositories": true,
			"restricted_to_workflows": false,
			"selected_workflows": []
		}]
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateEnterpriseRunnerGroupRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CreateEnterpriseRunnerGroupRequest{}, "{}")

	u := &CreateEnterpriseRunnerGroupRequest{
		Name:                     Ptr("n"),
		Visibility:               Ptr("v"),
		SelectedOrganizationIDs:  []int64{1},
		Runners:                  []int64{1},
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(true),
		SelectedWorkflows:        []string{"a", "b"},
	}

	want := `{
		"name": "n",
		"visibility": "v",
		"selected_organization_ids": [1],
		"runners": [1],
		"allows_public_repositories": true,
		"restricted_to_workflows": true,
		"selected_workflows": ["a","b"]
	}`

	testJSONMarshal(t, u, want)
}

func TestUpdateEnterpriseRunnerGroupRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UpdateEnterpriseRunnerGroupRequest{}, "{}")

	u := &UpdateEnterpriseRunnerGroupRequest{
		Name:                     Ptr("n"),
		Visibility:               Ptr("v"),
		AllowsPublicRepositories: Ptr(true),
		RestrictedToWorkflows:    Ptr(false),
		SelectedWorkflows:        []string{},
	}

	want := `{
		"name": "n",
		"visibility": "v",
		"allows_public_repositories": true,
		"restricted_to_workflows": false,
		"selected_workflows": []
	}`

	testJSONMarshal(t, u, want)
}

func TestSetOrgAccessRunnerGroupRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SetOrgAccessRunnerGroupRequest{}, "{}")

	u := &SetOrgAccessRunnerGroupRequest{
		SelectedOrganizationIDs: []int64{1},
	}

	want := `{
		"selected_organization_ids": [1]
	}`

	testJSONMarshal(t, u, want)
}
