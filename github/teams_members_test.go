// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamsService__ListTeamMembersByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, _, err := client.Teams.ListTeamMembersByID(context.Background(), 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersByID returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamMembersByID returned %+v, want %+v", members, want)
	}
}

func TestTeamsService__ListTeamMembersByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		w.WriteHeader(http.StatusNotFound)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, resp, err := client.Teams.ListTeamMembersByID(context.Background(), 1, 2, opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.ListTeamMembersByID returned status %d, want %d", got, want)
	}
	if members != nil {
		t.Errorf("Teams.ListTeamMembersByID returned %+v, want nil", members)
	}
}

func TestTeamsService__ListTeamMembersBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, _, err := client.Teams.ListTeamMembersBySlug(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersBySlug returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamMembersBySlug returned %+v, want %+v", members, want)
	}
}

func TestTeamsService__ListTeamMembersBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"role": "member", "page": "2"})
		w.WriteHeader(http.StatusNotFound)
	})

	opt := &TeamListTeamMembersOptions{Role: "member", ListOptions: ListOptions{Page: 2}}
	members, resp, err := client.Teams.ListTeamMembersBySlug(context.Background(), "o", "s", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.ListTeamMembersBySlug returned status %d, want %d", got, want)
	}
	if members != nil {
		t.Errorf("Teams.ListTeamMembersBySlug returned %+v, want nil", members)
	}
}

func TestTeamsService__ListTeamMembersBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.ListTeamMembersBySlug(context.Background(), "%", "s", nil)
	testURLParseError(t, err)
}

func TestTeamsService__GetTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	membership, _, err := client.Teams.GetTeamMembershipByID(context.Background(), 1, 2, "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.GetTeamMembershipByID returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__GetTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	membership, resp, err := client.Teams.GetTeamMembershipByID(context.Background(), 1, 2, "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamMembershipByID returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.GetTeamMembershipByID returned %+v, want nil", membership)
	}
}

func TestTeamsService__GetTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	membership, _, err := client.Teams.GetTeamMembershipBySlug(context.Background(), "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipBySlug returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.GetTeamMembershipBySlug returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__GetTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	membership, resp, err := client.Teams.GetTeamMembershipBySlug(context.Background(), "o", "s", "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamMembershipBySlug returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.GetTeamMembershipBySlug returned %+v, want nil", membership)
	}
}

func TestTeamsService__GetTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.GetTeamMembershipBySlug(context.Background(), "%s", "s", "u")
	testURLParseError(t, err)
}

func TestTeamsService__AddTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	membership, _, err := client.Teams.AddTeamMembershipByID(context.Background(), 1, 2, "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.AddTeamMembershipByID returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__AddTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNotFound)
	})

	membership, resp, err := client.Teams.AddTeamMembershipByID(context.Background(), 1, 2, "u", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.AddTeamMembershipByID returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.AddTeamMembershipByID returned %+v, want nil", membership)
	}
}

func TestTeamsService__AddTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	membership, _, err := client.Teams.AddTeamMembershipBySlug(context.Background(), "o", "s", "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipBySlug returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !reflect.DeepEqual(membership, want) {
		t.Errorf("Teams.AddTeamMembershipBySlug returned %+v, want %+v", membership, want)
	}
}

func TestTeamsService__AddTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNotFound)
	})

	membership, resp, err := client.Teams.AddTeamMembershipBySlug(context.Background(), "o", "s", "u", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.AddTeamMembershipBySlug returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.AddTeamMembershipBySlug returned %+v, want nil", membership)
	}
}

func TestTeamsService__AddTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.AddTeamMembershipBySlug(context.Background(), "%", "s", "u", nil)
	testURLParseError(t, err)
}

func TestTeamsService__RemoveTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamMembershipByID(context.Background(), 1, 2, "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned error: %v", err)
	}
}

func TestTeamsService__RemoveTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	resp, err := client.Teams.RemoveTeamMembershipByID(context.Background(), 1, 2, "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}
}

func TestTeamsService__RemoveTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamMembershipBySlug(context.Background(), "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipBySlug returned error: %v", err)
	}
}

func TestTeamsService__RemoveTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	resp, err := client.Teams.RemoveTeamMembershipBySlug(context.Background(), "o", "s", "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipBySlug returned status %d, want %d", got, want)
	}
}

func TestTeamsService__RemoveTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.RemoveTeamMembershipBySlug(context.Background(), "%", "s", "u")
	testURLParseError(t, err)
}

func TestTeamsService__ListPendingTeamInvitationsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	invitations, _, err := client.Teams.ListPendingTeamInvitationsByID(context.Background(), 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned error: %v", err)
	}

	want := []*Invitation{{ID: Int64(1)}}
	if !reflect.DeepEqual(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned %+v, want %+v", invitations, want)
	}
}

func TestTeamsService__ListPendingTeamInvitationsByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		w.WriteHeader(http.StatusNotFound)
	})

	opt := &ListOptions{Page: 2}
	invitations, resp, err := client.Teams.ListPendingTeamInvitationsByID(context.Background(), 1, 2, opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}
	if invitations != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned %+v, want nil", invitations)
	}
}

func TestTeamsService__ListPendingTeamInvitationsBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	invitations, _, err := client.Teams.ListPendingTeamInvitationsBySlug(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned error: %v", err)
	}

	want := []*Invitation{{ID: Int64(1)}}
	if !reflect.DeepEqual(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned %+v, want %+v", invitations, want)
	}
}

func TestTeamsService__ListPendingTeamInvitationsBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/invitations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		w.WriteHeader(http.StatusNotFound)
	})

	opt := &ListOptions{Page: 2}
	invitations, resp, err := client.Teams.ListPendingTeamInvitationsBySlug(context.Background(), "o", "s", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}
	if invitations != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned %+v, want nil", invitations)
	}
}

func TestTeamsService__ListPendingTeamInvitationsBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.ListPendingTeamInvitationsBySlug(context.Background(), "%", "s", nil)
	testURLParseError(t, err)
}
