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

func TestTeamsService_GetTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "description": "d", "url":"u", "slug": "s", "permission":"p", "ldap_dn":"cn=n,ou=groups,dc=example,dc=com", "parent":null}`)
	})

	team, _, err := client.Teams.GetTeam(context.Background(), 1)
	if err != nil {
		t.Errorf("Teams.GetTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"), LDAPDN: String("cn=n,ou=groups,dc=example,dc=com")}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.GetTeam returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_GetTeam_nestedTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "description": "d", "url":"u", "slug": "s", "permission":"p",
		"parent": {"id":2, "name":"n", "description": "d", "parent": null}}`)
	})

	team, _, err := client.Teams.GetTeam(context.Background(), 1)
	if err != nil {
		t.Errorf("Teams.GetTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1), Name: String("n"), Description: String("d"), URL: String("u"), Slug: String("s"), Permission: String("p"),
		Parent: &Team{ID: Int64(2), Name: String("n"), Description: String("d")},
	}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.GetTeam returned %+v, want %+v", team, want)
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

func TestTeamsService_EditTeam(t *testing.T) {
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

	team, _, err := client.Teams.EditTeam(context.Background(), "o", "s", input, false)
	if err != nil {
		t.Errorf("Teams.EditTeam returned error: %v", err)
	}

	want := &Team{ID: Int64(1)}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Teams.EditTeam returned %+v, want %+v", team, want)
	}
}

func TestTeamsService_EditTeam_RemoveParent(t *testing.T) {
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

	team, _, err := client.Teams.EditTeam(context.Background(), "o", "s", input, true)
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

func TestTeamsService_DeleteTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Teams.DeleteTeam(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.DeleteTeam returned error: %v", err)
	}
}

func TestTeamsService_ListChildTeams(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	teams, _, err := client.Teams.ListChildTeams(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeams returned error: %v", err)
	}

	want := []*Team{{ID: Int64(2)}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.ListTeams returned %+v, want %+v", teams, want)
	}
}

func TestTeamsService_ListTeamRepos(t *testing.T) {
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
	members, _, err := client.Teams.ListTeamRepos(context.Background(), "o", "s", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamRepos returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Teams.ListTeamRepos returned %+v, want %+v", members, want)
	}
}

func TestTeamsService_IsTeamRepo_true(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		wantAcceptHeaders := []string{mediaTypeOrgPermissionRepo}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	repo, _, err := client.Teams.IsTeamRepo(context.Background(), "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.IsTeamRepo returned error: %v", err)
	}

	want := &Repository{ID: Int64(1)}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("Teams.IsTeamRepo returned %+v, want %+v", repo, want)
	}
}

func TestTeamsService_IsTeamRepo_false(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	repo, resp, err := client.Teams.IsTeamRepo(context.Background(), "org", "slug", "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Teams.IsTeamRepo returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepo returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepo_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	repo, resp, err := client.Teams.IsTeamRepo(context.Background(), "org", "slug", "owner", "repo")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("Teams.IsTeamRepo returned status %d, want %d", got, want)
	}
	if repo != nil {
		t.Errorf("Teams.IsTeamRepo returned %+v, want nil", repo)
	}
}

func TestTeamsService_IsTeamRepo_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Teams.IsTeamRepo(context.Background(), "o", "s", "%", "r")
	testURLParseError(t, err)
}

func TestTeamsService_AddTeamRepo(t *testing.T) {
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

	_, err := client.Teams.AddTeamRepo(context.Background(), "org", "slug", "owner", "repo", opt)
	if err != nil {
		t.Errorf("Teams.AddTeamRepo returned error: %v", err)
	}
}

func TestTeamsService_AddTeamRepo_noAccess(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	_, err := client.Teams.AddTeamRepo(context.Background(), "org", "slug", "owner", "repo", nil)
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestTeamsService_AddTeamRepo_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.AddTeamRepo(context.Background(), "o", "s", "%", "r", nil)
	testURLParseError(t, err)
}

func TestTeamsService_RemoveTeamRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/org/teams/slug/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamRepo(context.Background(), "org", "slug", "owner", "repo")
	if err != nil {
		t.Errorf("Teams.RemoveTeamRepo returned error: %v", err)
	}
}

func TestTeamsService_RemoveTeamRepo_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Teams.RemoveTeamRepo(context.Background(), "o", "s", "%", "r")
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

func TestTeamsService_ListProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `[{"id":1}]`)
	})

	projects, _, err := client.Teams.ListTeamProjects(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.ListTeamProjects returned error: %v", err)
	}

	want := []*Project{{ID: Int64(1)}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Teams.ListTeamProjects returned %+v, want %+v", projects, want)
	}
}

func TestTeamsService_ReviewProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":1}`)
	})

	project, _, err := client.Teams.ReviewTeamProjects(context.Background(), "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.ReviewTeamProjects returned error: %v", err)
	}

	want := &Project{ID: Int64(1)}
	if !reflect.DeepEqual(project, want) {
		t.Errorf("Teams.ReviewTeamProjects returned %+v, want %+v", project, want)
	}
}

func TestTeamsService_AddTeamProject(t *testing.T) {
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

	_, err := client.Teams.AddTeamProject(context.Background(), "o", "s", 1, opt)
	if err != nil {
		t.Errorf("Teams.AddTeamProject returned error: %v", err)
	}
}

func TestTeamsService_RemoveTeamProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantAcceptHeaders := []string{mediaTypeProjectsPreview}
	mux.HandleFunc("/orgs/o/teams/s/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Teams.RemoveTeamProject(context.Background(), "o", "s", 1)
	if err != nil {
		t.Errorf("Teams.RemoveTeamProject returned error: %v", err)
	}
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

func TestTeamsService_ListIDPGroupsForTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"groups": [{"group_id": "1",  "group_name": "n", "group_description": "d"}]}`)
	})

	groups, _, err := client.Teams.ListIDPGroupsForTeam(context.Background(), "o", "s")
	if err != nil {
		t.Errorf("Teams.ListIDPGroupsForTeam returned error: %v", err)
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
		t.Errorf("Teams.ListIDPGroupsForTeam returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnections(t *testing.T) {
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

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnections(context.Background(), "o", "s", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnections returned error: %v", err)
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
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnections returned %+v. want %+v", groups, want)
	}
}

func TestTeamsService_CreateOrUpdateIDPGroupConnections_empty(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams/s/team-sync/group-mappings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"groups": []}`)
	})

	input := IDPGroupList{
		Groups: []*IDPGroup{},
	}

	groups, _, err := client.Teams.CreateOrUpdateIDPGroupConnections(context.Background(), "o", "s", input)
	if err != nil {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnections returned error: %v", err)
	}

	want := &IDPGroupList{
		Groups: []*IDPGroup{},
	}
	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Teams.CreateOrUpdateIDPGroupConnections returned %+v. want %+v", groups, want)
	}
}
