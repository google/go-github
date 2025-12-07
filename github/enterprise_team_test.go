// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseService_ListTeams(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": 1,
			"url": "https://example.com/team1",
			"member_url": "https://example.com/members",
			"name": "Team One",
			"html_url": "https://example.com/html",
			"slug": "team-one",
			"created_at": "2020-01-01T00:00:00Z",
			"updated_at": "2020-01-02T00:00:00Z",
			"group_id": "99"
		}]`)
	})

	ctx := t.Context()
	opts := &ListOptions{Page: 1, PerPage: 10}
	got, _, err := client.Enterprise.ListTeams(ctx, "e", opts)
	if err != nil {
		t.Fatalf("Enterprise.ListTeams returned error: %v", err)
	}

	want := []*EnterpriseTeam{
		{
			ID:        1,
			URL:       "https://example.com/team1",
			MemberURL: "https://example.com/members",
			Name:      "Team One",
			HTMLURL:   "https://example.com/html",
			Slug:      "team-one",
			GroupID:   "99",
			CreatedAt: Timestamp{Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			UpdatedAt: Timestamp{Time: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.ListTeams = %+v, want %+v", got, want)
	}

	const methodName = "ListTeams"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.ListTeams(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListTeams(ctx, "e", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateTeam(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := EnterpriseTeamCreateOrUpdateRequest{
		Name: "New Team",
	}

	mux.HandleFunc("/enterprises/e/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"New Team"}`+"\n")
		fmt.Fprint(w, `{
			"id": 10,
			"name": "New Team",
			"slug": "new-team",
			"url": "https://example.com/team"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.CreateTeam(ctx, "e", input)
	if err != nil {
		t.Fatalf("Enterprise.CreateTeam returned error: %v", err)
	}

	want := &EnterpriseTeam{
		ID:   10,
		Name: "New Team",
		Slug: "new-team",
		URL:  "https://example.com/team",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.CreateTeam = %+v, want %+v", got, want)
	}

	const methodName = "CreateTeam"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateTeam(ctx, "e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetTeam(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/teams/t1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 2,
			"name": "Team One",
			"slug": "t1"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.GetTeam(ctx, "e", "t1")
	if err != nil {
		t.Fatalf("Enterprise.GetTeam returned error: %v", err)
	}

	want := &EnterpriseTeam{
		ID:   2,
		Name: "Team One",
		Slug: "t1",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.GetTeam = %+v, want %+v", got, want)
	}

	const methodName = "GetTeam"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.GetTeam(ctx, "\n", "t1")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetTeam(ctx, "e", "t1")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateTeam(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := EnterpriseTeamCreateOrUpdateRequest{
		Name: "Updated Team",
	}

	mux.HandleFunc("/enterprises/e/teams/t1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"name":"Updated Team"}`+"\n")
		fmt.Fprint(w, `{
			"id": 3,
			"name": "Updated Team",
			"slug": "t1"
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Enterprise.UpdateTeam(ctx, "e", "t1", input)
	if err != nil {
		t.Fatalf("Enterprise.UpdateTeam returned error: %v", err)
	}

	want := &EnterpriseTeam{
		ID:   3,
		Name: "Updated Team",
		Slug: "t1",
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Enterprise.UpdateTeam = %+v, want %+v", got, want)
	}

	const methodName = "UpdateTeam"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Enterprise.UpdateTeam(ctx, "\n", "t1", input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateTeam(ctx, "e", "t1", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteTeam(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/teams/t1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Enterprise.DeleteTeam(ctx, "e", "t1")
	if err != nil {
		t.Fatalf("Enterprise.DeleteTeam returned error: %v", err)
	}

	const methodName = "DeleteTeam"
	testBadOptions(t, methodName, func() error {
		_, err := client.Enterprise.DeleteTeam(ctx, "\n", "t1")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteTeam(ctx, "e", "t1")
	})
}
