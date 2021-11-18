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
	"testing"

	"github.com/google/go-cmp/cmp"
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
	ctx := context.Background()
	members, _, err := client.Teams.ListTeamMembersByID(ctx, 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersByID returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !cmp.Equal(members, want) {
		t.Errorf("Teams.ListTeamMembersByID returned %+v, want %+v", members, want)
	}

	const methodName = "ListTeamMembersByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamMembersByID(ctx, -1, -2, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamMembersByID(ctx, 1, 2, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	members, resp, err := client.Teams.ListTeamMembersByID(ctx, 1, 2, opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.ListTeamMembersByID returned status %d, want %d", got, want)
	}
	if members != nil {
		t.Errorf("Teams.ListTeamMembersByID returned %+v, want nil", members)
	}

	const methodName = "ListTeamMembersByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamMembersByID(ctx, 1, 2, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamMembersByID(ctx, 1, 2, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	members, _, err := client.Teams.ListTeamMembersBySlug(ctx, "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembersBySlug returned error: %v", err)
	}

	want := []*User{{ID: Int64(1)}}
	if !cmp.Equal(members, want) {
		t.Errorf("Teams.ListTeamMembersBySlug returned %+v, want %+v", members, want)
	}

	const methodName = "ListTeamMembersBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamMembersBySlug(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamMembersBySlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	members, resp, err := client.Teams.ListTeamMembersBySlug(ctx, "o", "s", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.ListTeamMembersBySlug returned status %d, want %d", got, want)
	}
	if members != nil {
		t.Errorf("Teams.ListTeamMembersBySlug returned %+v, want nil", members)
	}

	const methodName = "ListTeamMembersBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamMembersBySlug(ctx, "o", "s", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamMembersBySlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__ListTeamMembersBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.ListTeamMembersBySlug(ctx, "%", "s", nil)
	testURLParseError(t, err)
}

func TestTeamsService__GetTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	ctx := context.Background()
	membership, _, err := client.Teams.GetTeamMembershipByID(ctx, 1, 2, "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !cmp.Equal(membership, want) {
		t.Errorf("Teams.GetTeamMembershipByID returned %+v, want %+v", membership, want)
	}

	const methodName = "GetTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamMembershipByID(ctx, -1, -2, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamMembershipByID(ctx, 1, 2, "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__GetTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	membership, resp, err := client.Teams.GetTeamMembershipByID(ctx, 1, 2, "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamMembershipByID returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.GetTeamMembershipByID returned %+v, want nil", membership)
	}

	const methodName = "GetTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamMembershipByID(ctx, 1, 2, "u")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamMembershipByID(ctx, 1, 2, "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__GetTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "state":"active"}`)
	})

	ctx := context.Background()
	membership, _, err := client.Teams.GetTeamMembershipBySlug(ctx, "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.GetTeamMembershipBySlug returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("active")}
	if !cmp.Equal(membership, want) {
		t.Errorf("Teams.GetTeamMembershipBySlug returned %+v, want %+v", membership, want)
	}

	const methodName = "GetTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamMembershipBySlug(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamMembershipBySlug(ctx, "o", "s", "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__GetTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	membership, resp, err := client.Teams.GetTeamMembershipBySlug(ctx, "o", "s", "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamMembershipBySlug returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.GetTeamMembershipBySlug returned %+v, want nil", membership)
	}

	const methodName = "GetTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamMembershipBySlug(ctx, "o", "s", "u")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamMembershipBySlug(ctx, "o", "s", "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__GetTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.GetTeamMembershipBySlug(ctx, "%s", "s", "u")
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
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	ctx := context.Background()
	membership, _, err := client.Teams.AddTeamMembershipByID(ctx, 1, 2, "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipByID returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !cmp.Equal(membership, want) {
		t.Errorf("Teams.AddTeamMembershipByID returned %+v, want %+v", membership, want)
	}

	const methodName = "AddTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.AddTeamMembershipByID(ctx, -1, -2, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.AddTeamMembershipByID(ctx, 1, 2, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__AddTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	membership, resp, err := client.Teams.AddTeamMembershipByID(ctx, 1, 2, "u", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.AddTeamMembershipByID returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.AddTeamMembershipByID returned %+v, want nil", membership)
	}

	const methodName = "AddTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.AddTeamMembershipByID(ctx, 1, 2, "u", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.AddTeamMembershipByID(ctx, 1, 2, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__AddTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		fmt.Fprint(w, `{"url":"u", "state":"pending"}`)
	})

	ctx := context.Background()
	membership, _, err := client.Teams.AddTeamMembershipBySlug(ctx, "o", "s", "u", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamMembershipBySlug returned error: %v", err)
	}

	want := &Membership{URL: String("u"), State: String("pending")}
	if !cmp.Equal(membership, want) {
		t.Errorf("Teams.AddTeamMembershipBySlug returned %+v, want %+v", membership, want)
	}

	const methodName = "AddTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.AddTeamMembershipBySlug(ctx, "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.AddTeamMembershipBySlug(ctx, "o", "s", "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__AddTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamMembershipOptions{Role: "maintainer"}

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamMembershipOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	membership, resp, err := client.Teams.AddTeamMembershipBySlug(ctx, "o", "s", "u", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.AddTeamMembershipBySlug returned status %d, want %d", got, want)
	}
	if membership != nil {
		t.Errorf("Teams.AddTeamMembershipBySlug returned %+v, want nil", membership)
	}

	const methodName = "AddTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.AddTeamMembershipBySlug(ctx, "o", "s", "u", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.AddTeamMembershipBySlug(ctx, "o", "s", "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__AddTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.AddTeamMembershipBySlug(ctx, "%", "s", "u", nil)
	testURLParseError(t, err)
}

func TestTeamsService__RemoveTeamMembershipByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamMembershipByID(ctx, 1, 2, "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned error: %v", err)
	}

	const methodName = "RemoveTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamMembershipByID(ctx, -1, -2, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamMembershipByID(ctx, 1, 2, "u")
	})
}

func TestTeamsService__RemoveTeamMembershipByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Teams.RemoveTeamMembershipByID(ctx, 1, 2, "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}

	const methodName = "RemoveTeamMembershipByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamMembershipByID(ctx, 1, 2, "u")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamMembershipByID(ctx, 1, 2, "u")
	})
}

func TestTeamsService__RemoveTeamMembershipBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamMembershipBySlug(ctx, "o", "s", "u")
	if err != nil {
		t.Errorf("Teams.RemoveTeamMembershipBySlug returned error: %v", err)
	}

	const methodName = "RemoveTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamMembershipBySlug(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamMembershipBySlug(ctx, "o", "s", "u")
	})
}

func TestTeamsService__RemoveTeamMembershipBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/memberships/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Teams.RemoveTeamMembershipBySlug(ctx, "o", "s", "u")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipBySlug returned status %d, want %d", got, want)
	}

	const methodName = "RemoveTeamMembershipBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamMembershipBySlug(ctx, "o", "s", "u")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamMembershipBySlug(ctx, "o", "s", "u")
	})
}

func TestTeamsService__RemoveTeamMembershipBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamMembershipBySlug(ctx, "%", "s", "u")
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
	ctx := context.Background()
	invitations, _, err := client.Teams.ListPendingTeamInvitationsByID(ctx, 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned error: %v", err)
	}

	want := []*Invitation{{ID: Int64(1)}}
	if !cmp.Equal(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned %+v, want %+v", invitations, want)
	}

	const methodName = "ListPendingTeamInvitationsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListPendingTeamInvitationsByID(ctx, -1, -2, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListPendingTeamInvitationsByID(ctx, 1, 2, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	invitations, resp, err := client.Teams.ListPendingTeamInvitationsByID(ctx, 1, 2, opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}
	if invitations != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned %+v, want nil", invitations)
	}

	const methodName = "ListPendingTeamInvitationsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListPendingTeamInvitationsByID(ctx, 1, 2, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListPendingTeamInvitationsByID(ctx, 1, 2, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	invitations, _, err := client.Teams.ListPendingTeamInvitationsBySlug(ctx, "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned error: %v", err)
	}

	want := []*Invitation{{ID: Int64(1)}}
	if !cmp.Equal(invitations, want) {
		t.Errorf("Teams.ListPendingTeamInvitationsByID returned %+v, want %+v", invitations, want)
	}

	const methodName = "ListPendingTeamInvitationsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListPendingTeamInvitationsBySlug(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListPendingTeamInvitationsBySlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
	ctx := context.Background()
	invitations, resp, err := client.Teams.ListPendingTeamInvitationsBySlug(ctx, "o", "s", opt)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.RemoveTeamMembershipByID returned status %d, want %d", got, want)
	}
	if invitations != nil {
		t.Errorf("Teams.RemoveTeamMembershipByID returned %+v, want nil", invitations)
	}

	const methodName = "ListPendingTeamInvitationsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListPendingTeamInvitationsBySlug(ctx, "o", "s", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListPendingTeamInvitationsBySlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService__ListPendingTeamInvitationsBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.ListPendingTeamInvitationsBySlug(ctx, "%", "s", nil)
	testURLParseError(t, err)
}

func TestTeamAddTeamMembershipOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamAddTeamMembershipOptions{}, "{}")

	u := &TeamAddTeamMembershipOptions{
		Role: "role",
	}

	want := `{
		"role": "role"
	}`

	testJSONMarshal(t, u, want)
}
