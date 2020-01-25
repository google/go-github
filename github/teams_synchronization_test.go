// Copyright 2018 The go-github AUTHORS. All rights reserved.
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

func TestTeamsService_ListIDPGroupsInOrganization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/team-sync/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	opt := &ListOptions{Page: 2}
	groups, _, err := client.Teams.ListIDPGroupsInOrganization(context.Background(), "o", opt)
	if err != nil {
		t.Errorf("Teams.ListIDPGroupsInOrganization returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.ListIDPGroupsInOrganization returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_ListIDPGroupsForTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	groups, _, err := client.Teams.ListIDPGroupsForTeamByID(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Teams.ListIDPGroupsForTeamByID returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.ListIDPGroupsForTeamByID returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_ListIDPGroupsForTeamByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	groups, _, err := client.Teams.ListIDPGroupsForTeamByName(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.ListIDPGroupsForTeamByName returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.ListIDPGroupsForTeamByName returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByID(context.Background(), 1, 1, input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByName(context.Background(), "o", "s", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByName returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{
			{
				GroupID:          String("1"),
				GroupName:        String("n"),
				GroupDescription: String("d"),
			},
		},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByName returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsByID_empty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": []}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{},
	}

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByID(context.Background(), 1, 1, input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsByName_empty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": []}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{},
	}

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByName(context.Background(), "o", "s", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByName returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByName returned %+v. want %+v", groups, want)
	}
}
