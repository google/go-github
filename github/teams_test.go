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
	"reflect"
	"strings"
	"testing"
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
	teams, _, err := client.Teams.ListTeams(context.Background(), "o", opt)
	if err != nil {
		t.Errorf("Teams.ListTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.ListTeams returned %+v, want %+v", teams, want)
	}
}

func TestTeamsService_ListTeams_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.ListTeams(context.Background(), "%", nil)
	testURLParseError(t, err)
}

func TestTeamsService_GetTeamByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "description": "d", "url":"u", "slug": "s", "permission":"p", "ldap_dn":"cn=n,ou=groups,dc=example,dc=com", "parent":null}`)
	})

	team, _, err := client.Teams.GetTeamByID(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Teams.GetTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"), LDAPDN: String("cn=n,ou=groups,dc=example,dc=com")}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.GetTeamByID returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_GetTeamByID_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	team, resp, err := client.Teams.GetTeamByID(context.Background(), 1, 2)
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

	team, _, err := client.Teams.GetTeamBySlug(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.GetTeamBySlug returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"), LDAPDN: String("cn=n,ou=groups,dc=example,dc=com")}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.GetTeamBySlug returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_GetTeamBySlug_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.GetTeamBySlug(context.Background(), "%", "s")
	testURLParseError(t, err)
}

func TestTeamsService_GetTeamBySlug_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	team, resp, err := client.Teams.GetTeamBySlug(context.Background(), "o", "s")
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
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, _, err := client.Teams.CreateTeam(context.Background(), "o", input)
	if err != nil {
		t.Errorf("Teams.CreateTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.CreateTeam returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_CreateTeam_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.CreateTeam(context.Background(), "%", NewTeam{})
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
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, _, err := client.Teams.EditTeamByID(context.Background(), 1, 1, input, false)
	if err != nil {
		t.Errorf("Teams.EditTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.EditTeamByID returned %+v, want %+v", team, want)
	}
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
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, _, err := client.Teams.EditTeamByID(context.Background(), 1, 1, input, true)
	if err != nil {
		t.Errorf("Teams.EditTeamByID returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.EditTeamByID returned %+v, want %+v", team, want)
	}

	if want := `{"name":"n","parent_team_id":null,"privacy":"closed"}` + "\n"; body != want {
		t.Errorf("Teams.EditTeamByID body = %+v, want %+v", body, want)
	}
}

func TestTeamsService_EditTeamByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := NewTeam{Name: "n", Privacy: String("closed")}

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewTeam)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, _, err := client.Teams.EditTeamByName(context.Background(), "o", "s", input, false)
	if err != nil {
		t.Errorf("Teams.EditTeamByName returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.EditTeamByName returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_EditTeamByName_RemoveParent(t *testing.T) {
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
		if !reflect.DeepEqual(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, _, err := client.Teams.EditTeamByName(context.Background(), "o", "s", input, true)
	if err != nil {
		t.Errorf("Teams.EditTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
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

	_, err := client.Teams.DeleteTeamByID(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Teams.DeleteTeamByID returned error: %v", err)
	}
}

func TestTeamsService_DeleteTeamByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Teams.DeleteTeamByName(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.DeleteTeamByName returned error: %v", err)
	}
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
	teams, _, err := client.Teams.ListChildTeamsByParentID(context.Background(), 1, 2, opt)
	if err != nil {
		t.Errorf("Teams.ListChildTeamsByParentID returned error: %v", err)
	}

	want := []*Team{{ID: Int64(2)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.ListChildTeamsByParentID returned %+v, want %+v", teams, want)
	}
}

func TestTeamsService_ListChildTeamsByParentName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	teams, _, err := client.Teams.ListChildTeamsByParentName(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListChildTeamsByParentName returned error: %v", err)
	}

	want := []*Team{{ID: Int64(2)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.ListChildTeamsByParentName returned %+v, want %+v", teams, want)
	}
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
	members, _, err := client.Teams.ListTeamReposByID(context.Background(), 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.ListTeamReposByID returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamReposByID returned %+v, want %+v", members, want)
	}
}

func TestTeamsService_ListTeamReposByName(t *testing.T) {
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
	members, _, err := client.Teams.ListTeamReposByName(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamReposByName returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamReposByName returned %+v, want %+v", members, want)
	}
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

	repo, _, err := client.Teams.IsTeamRepoByID(context.Background(), 1, 1, "owner", "repo")
	if err != nil {
		t.Errorf("Teams.IsTeamRepoByID returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Teams.IsTeamRepoByID returned %+v, want %+v", repo, want)
	}
}

func TestTeamsService_IsTeamRepoByName_true(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeOrgPermissionRepo}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Teams.IsTeamRepoByName(context.Background(), "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.IsTeamRepoByName returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Teams.IsTeamRepoByName returned %+v, want %+v", repo, want)
	}
}

func TestTeamsService_IsTeamRepoByID_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	repo, resp, err := client.Teams.IsTeamRepoByID(context.Background(), 1, 1, "owner", "repo")
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

func TestTeamsService_IsTeamRepoByName_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	repo, resp, err := client.Teams.IsTeamRepoByName(context.Background(), "org", "slug", "owner", "repo")
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

	repo, resp, err := client.Teams.IsTeamRepoByID(context.Background(), 1, 1, "owner", "repo")
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

func TestTeamsService_IsTeamRepoByName_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	repo, resp, err := client.Teams.IsTeamRepoByName(context.Background(), "org", "slug", "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("Teams.IsTeamRepoByName returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepoByName returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.IsTeamRepoByID(context.Background(), 1, 1, "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_IsTeamRepoByName_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.IsTeamRepoByName(context.Background(), "o", "s", "%", "r")
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
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.AddTeamRepoByID(context.Background(), 1, 1, "owner", "repo", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamRepoByID returned error: %v", err)
	}
}

func TestTeamsService_AddTeamRepoByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	opt := &TeamAddTeamRepoOptions{Permission: "admin"}

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamAddTeamRepoOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.AddTeamRepoByName(context.Background(), "org", "slug", "owner", "repo", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamRepoByName returned error: %v", err)
	}
}

func TestTeamsService_AddTeamRepoByID_noAccess(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	_, err := client.Teams.AddTeamRepoByID(context.Background(), 1, 1, "owner", "repo", nil)
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestTeamsService_AddTeamRepoByName_noAccess(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	_, err := client.Teams.AddTeamRepoByName(context.Background(), "org", "slug", "owner", "repo", nil)
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestTeamsService_AddTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.AddTeamRepoByID(context.Background(), 1, 1, "%", "r", nil)
	testURLParseError(t, err)
}

func TestTeamsService_AddTeamRepoByName_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.AddTeamRepoByName(context.Background(), "o", "s", "%", "r", nil)
	testURLParseError(t, err)
}

func TestTeamsService_RemoveTeamRepoByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/1/team/1/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamRepoByID(context.Background(), 1, 1, "owner", "repo")
	if err != nil {
		t.Errorf("Teams.RemoveTeamRepoByID returned error: %v", err)
	}
}

func TestTeamsService_RemoveTeamRepoByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamRepoByName(context.Background(), "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.RemoveTeamRepoByName returned error: %v", err)
	}
}

func TestTeamsService_RemoveTeamRepoByID_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.RemoveTeamRepoByID(context.Background(), 1, 1, "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_RemoveTeamRepoByName_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.RemoveTeamRepoByName(context.Background(), "o", "s", "%", "r")
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
	teams, _, err := client.Teams.ListUserTeams(context.Background(), opt)
	if err != nil {
		t.Errorf("Teams.ListUserTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(1)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.ListUserTeams returned %+v, want %+v", teams, want)
	}
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

	projects, _, err := client.Teams.ListTeamProjectsByID(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Teams.ListTeamProjectsByID returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Teams.ListTeamProjectsByID returned %+v, want %+v", projects, want)
	}
}

func TestTeamsService_ListProjectsByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1}]`)
	})

	projects, _, err := client.Teams.ListTeamProjectsByName(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.ListTeamProjectsByName returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Teams.ListTeamProjectsByName returned %+v, want %+v", projects, want)
	}
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

	project, _, err := client.Teams.ReviewTeamProjectsByID(context.Background(), 1, 1, 1)
	if err != nil {
		t.Errorf("Teams.ReviewTeamProjectsByID returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Teams.ReviewTeamProjectsByID returned %+v, want %+v", project, want)
	}
}

func TestTeamsService_ReviewProjectsByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Teams.ReviewTeamProjectsByName(context.Background(), "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.ReviewTeamProjectsByName returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Teams.ReviewTeamProjectsByName returned %+v, want %+v", project, want)
	}
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
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.AddTeamProjectByID(context.Background(), 1, 1, 1, opt)
	if err != nil {
		t.Errorf("Teams.AddTeamProjectByID returned error: %v", err)
	}
}

func TestTeamsService_AddTeamProjectByName(t *testing.T) {
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
		if !reflect.DeepEqual(v, opt) {
			t.Errorf("Request body = %+v, want %+v", v, opt)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.AddTeamProjectByName(context.Background(), "o", "s", 1, opt)
	if err != nil {
		t.Errorf("Teams.AddTeamProjectByName returned error: %v", err)
	}
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

	_, err := client.Teams.RemoveTeamProjectByID(context.Background(), 1, 1, 1)
	if err != nil {
		t.Errorf("Teams.RemoveTeamProjectByID returned error: %v", err)
	}
}

func TestTeamsService_RemoveTeamProjectByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamProjectByName(context.Background(), "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.RemoveTeamProjectByName returned error: %v", err)
	}
}
