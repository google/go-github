// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTeamsService_ListTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	teams, _, err := client.Teams.ListTeams(ctx, "o", opt)
	if err != nil {
		t.Errorf("Teams.ListTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !cmp.Equal(teams, want) {
		t.Errorf("Teams.ListTeams returned %+v, want %+v", teams, want)
	}

	const methodName = "ListTeams"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeams(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeams(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListTeams_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.ListTeams(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestTeamsService_GetTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "description": "d", "url":"u", "slug": "s", "permission":"p", "ldap_dn":"cn=n,ou=groups,dc=example,dc=com", "parent":null}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.GetTeamByID(ctx, 1, 1)
	if err != nil {
		t.Errorf("Teams.GetTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"), LDAPDN: String("cn=n,ou=groups,dc=example,dc=com")}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.GetTeamByID returned %+v, want %+v", team, want)
	}

	const methodName = "GetTeamByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamByID(ctx, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamByID(ctx, 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_GetTeamByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	team, resp, err := client.Teams.GetTeamByID(ctx, 1, 2)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamByID returned status %d, want %d", got, want)
	}
	if team != nil {
		t.Errorf("Teams.GetTeamByID returned %+v, want nil", team)
	}
}

func TestTeamsService_GetTeamBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "description": "d", "url":"u", "slug": "s", "permission":"p", "ldap_dn":"cn=n,ou=groups,dc=example,dc=com", "parent":null}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.GetTeamBySlug(ctx, "o", "s")
	if err != nil {
		t.Errorf("Teams.GetTeamBySlug returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"), LDAPDN: String("cn=n,ou=groups,dc=example,dc=com")}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.GetTeamBySlug returned %+v, want %+v", team, want)
	}

	const methodName = "GetTeamBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.GetTeamBySlug(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.GetTeamBySlug(ctx, "o", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_GetTeamBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.GetTeamBySlug(ctx, "%", "s")
	testURLParseError(t, err)
}

func TestTeamsService_GetTeamBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	team, resp, err := client.Teams.GetTeamBySlug(ctx, "o", "s")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.GetTeamBySlug returned status %d, want %d", got, want)
	}
	if team != nil {
		t.Errorf("Teams.GetTeamBySlug returned %+v, want nil", team)
	}
}

func TestTeamsService_CreateTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed"), RepoNames: []string{"r"}}

	mux.HandleFunc("/orgs/o/teams", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.CreateTeam(ctx, "o", input)
	if err != nil {
		t.Errorf("Teams.CreateTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.CreateTeam returned %+v, want %+v", team, want)
	}

	const methodName = "CreateTeam"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.CreateTeam(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.CreateTeam(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_CreateTeam_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.CreateTeam(ctx, "%", NewTeam{})
	testURLParseError(t, err)
}

func TestTeamsService_EditTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed")}

	mux.HandleFunc("/organizations/1/team/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.EditTeamByID(ctx, 1, 1, input, false)
	if err != nil {
		t.Errorf("Teams.EditTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.EditTeamByID returned %+v, want %+v", team, want)
	}

	const methodName = "EditTeamByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.EditTeamByID(ctx, -1, -1, input, false)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.EditTeamByID(ctx, 1, 1, input, false)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_EditTeamByID_RemoveParent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed")}
	var body string

	mux.HandleFunc("/organizations/1/team/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read body: %v", err)
		}
		body = string(buf)
		json.NewDecoder(bytes.NewBuffer(buf)).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.EditTeamByID(ctx, 1, 1, input, true)
	if err != nil {
		t.Errorf("Teams.EditTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.EditTeamByID returned %+v, want %+v", team, want)
	}

	if want := `{"name":"n","parent_team_id":null,"privacy":"closed"}` + "\n"; body != want {
		t.Errorf("Teams.EditTeamByID body = %+v, want %+v", body, want)
	}
}

func TestTeamsService_EditTeamBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed")}

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.EditTeamBySlug(ctx, "o", "s", input, false)
	if err != nil {
		t.Errorf("Teams.EditTeamBySlug returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.EditTeamBySlug returned %+v, want %+v", team, want)
	}

	const methodName = "EditTeamBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.EditTeamBySlug(ctx, "\n", "\n", input, false)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.EditTeamBySlug(ctx, "o", "s", input, false)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_EditTeamBySlug_RemoveParent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed")}
	var body string

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read body: %v", err)
		}
		body = string(buf)
		json.NewDecoder(bytes.NewBuffer(buf)).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	team, _, err := client.Teams.EditTeamBySlug(ctx, "o", "s", input, true)
	if err != nil {
		t.Errorf("Teams.EditTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !cmp.Equal(team, want) {
		t.Errorf("Teams.EditTeam returned %+v, want %+v", team, want)
	}

	if want := `{"name":"n","parent_team_id":null,"privacy":"closed"}` + "\n"; body != want {
		t.Errorf("Teams.EditTeam body = %+v, want %+v", body, want)
	}
}

func TestTeamsService_DeleteTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Teams.DeleteTeamByID(ctx, 1, 1)
	if err != nil {
		t.Errorf("Teams.DeleteTeamByID returned error: %v", err)
	}

	const methodName = "DeleteTeamByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.DeleteTeamByID(ctx, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.DeleteTeamByID(ctx, 1, 1)
	})
}

func TestTeamsService_DeleteTeamBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Teams.DeleteTeamBySlug(ctx, "o", "s")
	if err != nil {
		t.Errorf("Teams.DeleteTeamBySlug returned error: %v", err)
	}

	const methodName = "DeleteTeamBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.DeleteTeamBySlug(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.DeleteTeamBySlug(ctx, "o", "s")
	})
}

func TestTeamsService_ListChildTeamsByParentID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	teams, _, err := client.Teams.ListChildTeamsByParentID(ctx, 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListChildTeamsByParentID returned error: %v", err)
	}

	want := []*Team{{ID: Int64(2)}}
	if !cmp.Equal(teams, want) {
		t.Errorf("Teams.ListChildTeamsByParentID returned %+v, want %+v", teams, want)
	}

	const methodName = "ListChildTeamsByParentID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListChildTeamsByParentID(ctx, -1, -2, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListChildTeamsByParentID(ctx, 1, 2, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListChildTeamsByParentSlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	teams, _, err := client.Teams.ListChildTeamsByParentSlug(ctx, "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListChildTeamsByParentSlug returned error: %v", err)
	}

	want := []*Team{{ID: Int64(2)}}
	if !cmp.Equal(teams, want) {
		t.Errorf("Teams.ListChildTeamsByParentSlug returned %+v, want %+v", teams, want)
	}

	const methodName = "ListChildTeamsByParentSlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListChildTeamsByParentSlug(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListChildTeamsByParentSlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListTeamReposByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeTopicsPreview}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	members, _, err := client.Teams.ListTeamReposByID(ctx, 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.ListTeamReposByID returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !cmp.Equal(members, want) {
		t.Errorf("Teams.ListTeamReposByID returned %+v, want %+v", members, want)
	}

	const methodName = "ListTeamReposByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamReposByID(ctx, -1, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamReposByID(ctx, 1, 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListTeamReposBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeTopicsPreview}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	members, _, err := client.Teams.ListTeamReposBySlug(ctx, "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamReposBySlug returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !cmp.Equal(members, want) {
		t.Errorf("Teams.ListTeamReposBySlug returned %+v, want %+v", members, want)
	}

	const methodName = "ListTeamReposBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamReposBySlug(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamReposBySlug(ctx, "o", "s", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_IsTeamRepoByID_true(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeOrgPermissionRepo}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	repo, _, err := client.Teams.IsTeamRepoByID(ctx, 1, 1, "owner", "repo")
	if err != nil {
		t.Errorf("Teams.IsTeamRepoByID returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !cmp.Equal(repo, want) {
		t.Errorf("Teams.IsTeamRepoByID returned %+v, want %+v", repo, want)
	}

	const methodName = "IsTeamRepoByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.IsTeamRepoByID(ctx, -1, -1, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.IsTeamRepoByID(ctx, 1, 1, "owner", "repo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_IsTeamRepoBySlug_true(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeOrgPermissionRepo}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	repo, _, err := client.Teams.IsTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.IsTeamRepoBySlug returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !cmp.Equal(repo, want) {
		t.Errorf("Teams.IsTeamRepoBySlug returned %+v, want %+v", repo, want)
	}

	const methodName = "IsTeamRepoBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.IsTeamRepoBySlug(ctx, "\n", "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.IsTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_IsTeamRepoByID_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	repo, resp, err := client.Teams.IsTeamRepoByID(ctx, 1, 1, "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.IsTeamRepoByID returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepoByID returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepoBySlug_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	repo, resp, err := client.Teams.IsTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.IsTeamRepoByID returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepoByID returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepoByID_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	ctx := context.Background()
	repo, resp, err := client.Teams.IsTeamRepoByID(ctx, 1, 1, "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("Teams.IsTeamRepoByID returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepoByID returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepoBySlug_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	ctx := context.Background()
	repo, resp, err := client.Teams.IsTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("Teams.IsTeamRepoBySlug returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepoBySlug returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.IsTeamRepoByID(ctx, 1, 1, "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_IsTeamRepoBySlug_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Teams.IsTeamRepoBySlug(ctx, "o", "s", "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_AddTeamRepoByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamRepoOptions{Permission: "admin"}

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamRepoOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoByID(ctx, 1, 1, "owner", "repo", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamRepoByID returned error: %v", err)
	}

	const methodName = "AddTeamRepoByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.AddTeamRepoByID(ctx, 1, 1, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.AddTeamRepoByID(ctx, 1, 1, "owner", "repo", opt)
	})
}

func TestTeamsService_AddTeamRepoBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamRepoOptions{Permission: "admin"}

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamRepoOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoBySlug(ctx, "org", "slug", "owner", "repo", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamRepoBySlug returned error: %v", err)
	}

	const methodName = "AddTeamRepoBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.AddTeamRepoBySlug(ctx, "\n", "\n", "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.AddTeamRepoBySlug(ctx, "org", "slug", "owner", "repo", opt)
	})
}

func TestTeamsService_AddTeamRepoByID_noAccess(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoByID(ctx, 1, 1, "owner", "repo", nil)
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestTeamsService_AddTeamRepoBySlug_noAccess(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoBySlug(ctx, "org", "slug", "owner", "repo", nil)
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestTeamsService_AddTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoByID(ctx, 1, 1, "%", "r", nil)
	testURLParseError(t, err)
}

func TestTeamsService_AddTeamRepoBySlug_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Teams.AddTeamRepoBySlug(ctx, "o", "s", "%", "r", nil)
	testURLParseError(t, err)
}

func TestTeamsService_RemoveTeamRepoByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamRepoByID(ctx, 1, 1, "owner", "repo")
	if err != nil {
		t.Errorf("Teams.RemoveTeamRepoByID returned error: %v", err)
	}

	const methodName = "RemoveTeamRepoByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamRepoByID(ctx, -1, -1, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamRepoByID(ctx, 1, 1, "owner", "repo")
	})
}

func TestTeamsService_RemoveTeamRepoBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.RemoveTeamRepoBySlug returned error: %v", err)
	}

	const methodName = "RemoveTeamRepoBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamRepoBySlug(ctx, "\n", "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamRepoBySlug(ctx, "org", "slug", "owner", "repo")
	})
}

func TestTeamsService_RemoveTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamRepoByID(ctx, 1, 1, "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_RemoveTeamRepoBySlug_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamRepoBySlug(ctx, "o", "s", "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_ListUserTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1}
	ctx := context.Background()
	teams, _, err := client.Teams.ListUserTeams(ctx, opt)
	if err != nil {
		t.Errorf("Teams.ListUserTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !cmp.Equal(teams, want) {
		t.Errorf("Teams.ListUserTeams returned %+v, want %+v", teams, want)
	}

	const methodName = "ListUserTeams"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListUserTeams(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListProjectsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/organizations/1/team/1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	projects, _, err := client.Teams.ListTeamProjectsByID(ctx, 1, 1)
	if err != nil {
		t.Errorf("Teams.ListTeamProjectsByID returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !cmp.Equal(projects, want) {
		t.Errorf("Teams.ListTeamProjectsByID returned %+v, want %+v", projects, want)
	}

	const methodName = "ListTeamProjectsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamProjectsByID(ctx, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamProjectsByID(ctx, 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListProjectsBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	projects, _, err := client.Teams.ListTeamProjectsBySlug(ctx, "o", "s")
	if err != nil {
		t.Errorf("Teams.ListTeamProjectsBySlug returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !cmp.Equal(projects, want) {
		t.Errorf("Teams.ListTeamProjectsBySlug returned %+v, want %+v", projects, want)
	}

	const methodName = "ListTeamProjectsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListTeamProjectsBySlug(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListTeamProjectsBySlug(ctx, "o", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ReviewProjectsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/organizations/1/team/1/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Teams.ReviewTeamProjectsByID(ctx, 1, 1, 1)
	if err != nil {
		t.Errorf("Teams.ReviewTeamProjectsByID returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Teams.ReviewTeamProjectsByID returned %+v, want %+v", project, want)
	}

	const methodName = "ReviewTeamProjectsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ReviewTeamProjectsByID(ctx, -1, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ReviewTeamProjectsByID(ctx, 1, 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ReviewProjectsBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	project, _, err := client.Teams.ReviewTeamProjectsBySlug(ctx, "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.ReviewTeamProjectsBySlug returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !cmp.Equal(project, want) {
		t.Errorf("Teams.ReviewTeamProjectsBySlug returned %+v, want %+v", project, want)
	}

	const methodName = "ReviewTeamProjectsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ReviewTeamProjectsBySlug(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ReviewTeamProjectsBySlug(ctx, "o", "s", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_AddTeamProjectByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamProjectOptions{
		Permission: String("admin"),
	}

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/organizations/1/team/1/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))

		v := &TeamProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamProjectByID(ctx, 1, 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.AddTeamProjectByID returned error: %v", err)
	}

	const methodName = "AddTeamProjectByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.AddTeamProjectByID(ctx, -1, -1, -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.AddTeamProjectByID(ctx, 1, 1, 1, opt)
	})
}

func TestTeamsService_AddTeamProjectBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamProjectOptions{
		Permission: String("admin"),
	}

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))

		v := &TeamProjectOptions{}
		json.NewDecoder(r.Body).Decode(v)
		if !cmp.Equal(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.AddTeamProjectBySlug(ctx, "o", "s", 1, opt)
	if err != nil {
		t.Errorf("Teams.AddTeamProjectBySlug returned error: %v", err)
	}

	const methodName = "AddTeamProjectBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.AddTeamProjectBySlug(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.AddTeamProjectBySlug(ctx, "o", "s", 1, opt)
	})
}

func TestTeamsService_RemoveTeamProjectByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/organizations/1/team/1/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamProjectByID(ctx, 1, 1, 1)
	if err != nil {
		t.Errorf("Teams.RemoveTeamProjectByID returned error: %v", err)
	}

	const methodName = "RemoveTeamProjectByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamProjectByID(ctx, -1, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamProjectByID(ctx, 1, 1, 1)
	})
}

func TestTeamsService_RemoveTeamProjectBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Teams.RemoveTeamProjectBySlug(ctx, "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.RemoveTeamProjectBySlug returned error: %v", err)
	}

	const methodName = "RemoveTeamProjectBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Teams.RemoveTeamProjectBySlug(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Teams.RemoveTeamProjectBySlug(ctx, "o", "s", 1)
	})
}

func TestTeamsService_ListIDPGroupsInOrganization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/team-sync/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "url-encoded-next-page-token",
		})
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	opt := &ListCursorOptions{Page: "url-encoded-next-page-token"}
	ctx := context.Background()
	groups, _, err := client.Teams.ListIDPGroupsInOrganization(ctx, "o", opt)
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
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.ListIDPGroupsInOrganization returned %+v. want %+v", groups, want)
	}

	const methodName = "ListIDPGroupsInOrganization"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListIDPGroupsInOrganization(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListIDPGroupsInOrganization(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListIDPGroupsForTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	ctx := context.Background()
	groups, _, err := client.Teams.ListIDPGroupsForTeamByID(ctx, 1, 1)
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
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.ListIDPGroupsForTeamByID returned %+v. want %+v", groups, want)
	}

	const methodName = "ListIDPGroupsForTeamByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListIDPGroupsForTeamByID(ctx, -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListIDPGroupsForTeamByID(ctx, 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_ListIDPGroupsForTeamBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/slug/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	ctx := context.Background()
	groups, _, err := client.Teams.ListIDPGroupsForTeamBySlug(ctx, "o", "slug")
	if err != nil {
		t.Errorf("Teams.ListIDPGroupsForTeamBySlug returned error: %v", err)
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
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.ListIDPGroupsForTeamBySlug returned %+v. want %+v", groups, want)
	}

	const methodName = "ListIDPGroupsForTeamBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.ListIDPGroupsForTeamBySlug(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.ListIDPGroupsForTeamBySlug(ctx, "o", "slug")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByID(ctx, 1, 1, input)
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
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned %+v. want %+v", groups, want)
	}

	const methodName = "CreateOrUpdateIDPGroupConnectionsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.CreateOrUpdateIDPGroupConnectionsByID(ctx, -1, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByID(ctx, 1, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsBySlug(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/slug/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.Background()
	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsBySlug(ctx, "o", "slug", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsBySlug returned error: %v", err)
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
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsBySlug returned %+v. want %+v", groups, want)
	}

	const methodName = "CreateOrUpdateIDPGroupConnectionsBySlug"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Teams.CreateOrUpdateIDPGroupConnectionsBySlug(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Teams.CreateOrUpdateIDPGroupConnectionsBySlug(ctx, "o", "slug", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

	ctx := context.Background()
	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsByID(ctx, 1, 1, input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsByID returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnectionsBySlug_empty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/slug/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": []}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{},
	}

	ctx := context.Background()
	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnectionsBySlug(ctx, "o", "slug", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsBySlug returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnectionsBySlug returned %+v. want %+v", groups, want)
	}
}

func TestNewTeam_Marshal(t *testing.T) {
	testJSONMarshal(t, &NewTeam{}, "{}")

	u := &NewTeam{
		Name:         "n",
		Description:  String("d"),
		Maintainers:  []string{"m1", "m2"},
		RepoNames:    []string{"repo1", "repo2"},
		ParentTeamID: Int64(1),
		Permission:   String("perm"),
		Privacy:      String("p"),
		LDAPDN:       String("l"),
	}

	want := `{
		"name":           "n",
		"description":    "d",
		"maintainers":    ["m1", "m2"],
		"repo_names":     ["repo1", "repo2"],
		"parent_team_id": 1,
		"permission":     "perm",
		"privacy":        "p",
		"ldap_dn":        "l"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeams_Marshal(t *testing.T) {
	testJSONMarshal(t, &Team{}, "{}")

	u := &Team{
		ID:              Int64(1),
		NodeID:          String("n"),
		Name:            String("n"),
		Description:     String("d"),
		URL:             String("u"),
		Slug:            String("s"),
		Permission:      String("p"),
		Privacy:         String("p"),
		MembersCount:    Int(1),
		ReposCount:      Int(1),
		MembersURL:      String("m"),
		RepositoriesURL: String("r"),
		Organization: &Organization{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			AvatarURL: String("a"),
			HTMLURL:   String("h"),
			Name:      String("n"),
			Company:   String("c"),
			Blog:      String("b"),
			Location:  String("l"),
			Email:     String("e"),
		},
		Parent: &Team{
			ID:           Int64(1),
			NodeID:       String("n"),
			Name:         String("n"),
			Description:  String("d"),
			URL:          String("u"),
			Slug:         String("s"),
			Permission:   String("p"),
			Privacy:      String("p"),
			MembersCount: Int(1),
			ReposCount:   Int(1),
		},
		LDAPDN: String("l"),
	}

	want := `{
		"id": 1,
		"node_id": "n",
		"name": "n",
		"description": "d",
		"url": "u",
		"slug": "s",
		"permission": "p",
		"privacy": "p",
		"members_count": 1,
		"repos_count": 1,
		"members_url": "m",
		"repositories_url": "r",
		"organization": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"html_url": "h",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e"
		},
		"parent": {
			"id": 1,
			"node_id": "n",
			"name": "n",
			"description": "d",
			"url": "u",
			"slug": "s",
			"permission": "p",
			"privacy": "p",
			"members_count": 1,
			"repos_count": 1
		},
		"ldap_dn": "l"	
	}`

	testJSONMarshal(t, u, want)
}

func TestInvitation_Marshal(t *testing.T) {
	testJSONMarshal(t, &Invitation{}, "{}")

	u := &Invitation{
		ID:                Int64(1),
		NodeID:            String("test node"),
		Login:             String("login123"),
		Email:             String("go@github.com"),
		Role:              String("developer"),
		CreatedAt:         &referenceTime,
		TeamCount:         Int(99),
		InvitationTeamURL: String("url"),
	}

	want := `{
		"id": 1,
		"node_id": "test node",
		"login":"login123",
		"email":"go@github.com",
		"role":"developer",
		"created_at":` + referenceTimeStr + `,
		"team_count":99,
		"invitation_team_url":"url"
	}`

	testJSONMarshal(t, u, want)
}

func TestIDPGroup_Marshal(t *testing.T) {
	testJSONMarshal(t, &IDPGroup{}, "{}")

	u := &IDPGroup{
		GroupID:          String("abc1"),
		GroupName:        String("test group"),
		GroupDescription: String("test group descripation"),
	}

	want := `{
		"group_id": "abc1",
		"group_name": "test group",
		"group_description":"test group descripation"
	}`

	testJSONMarshal(t, u, want)
}
