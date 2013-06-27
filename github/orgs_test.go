// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrganizationsService_List_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	orgs, err := client.Organizations.List("", nil)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []Organization{Organization{ID: 1}, Organization{ID: 2}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}
}

func TestOrganizationsService_List_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{2}
	orgs, err := client.Organizations.List("u", opt)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []Organization{Organization{ID: 1}, Organization{ID: 2}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}
}

func TestOrganizationsService_List_invalidUser(t *testing.T) {
	_, err := client.Organizations.List("%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "login":"l", "url":"u", "avatar_url": "a", "location":"l"}`)
	})

	org, err := client.Organizations.Get("o")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	want := &Organization{ID: 1, Login: "l", URL: "u", AvatarURL: "a", Location: "l"}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Organizations.Get returned %+v, want %+v", org, want)
	}
}

func TestOrganizationsService_Get_invalidOrg(t *testing.T) {
	_, err := client.Organizations.Get("%")
	testURLParseError(t, err)
}

func TestOrganizationsService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &Organization{Login: "l"}

	mux.HandleFunc("/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		v := new(Organization)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	org, err := client.Organizations.Edit("o", input)
	if err != nil {
		t.Errorf("Organizations.Edit returned error: %v", err)
	}

	want := &Organization{ID: 1}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Organizations.Edit returned %+v, want %+v", org, want)
	}
}

func TestOrganizationsService_Edit_invalidOrg(t *testing.T) {
	_, err := client.Organizations.Edit("%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_ListMembers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	members, err := client.Organizations.ListMembers("o")
	if err != nil {
		t.Errorf("Organizations.ListMembers returned error: %v", err)
	}

	want := []User{User{ID: 1}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Organizations.ListMembers returned %+v, want %+v", members, want)
	}
}

func TestOrganizationsService_ListMembers_invalidOrg(t *testing.T) {
	_, err := client.Organizations.ListMembers("%")
	testURLParseError(t, err)
}

func TestOrganizationsService_ListPublicMembers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	members, err := client.Organizations.ListPublicMembers("o")
	if err != nil {
		t.Errorf("Organizations.ListPublicMembers returned error: %v", err)
	}

	want := []User{User{ID: 1}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Organizations.ListPublicMembers returned %+v, want %+v", members, want)
	}
}

func TestOrganizationsService_ListPublicMembers_invalidOrg(t *testing.T) {
	_, err := client.Organizations.ListPublicMembers("%")
	testURLParseError(t, err)
}

func TestOrganizationsService_CheckMembership(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	member, err := client.Organizations.CheckMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.CheckMembership returned error: %v", err)
	}
	if want := true; member != want {
		t.Errorf("Organizations.CheckMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 404 response is interpreted as "false" and not an error
func TestOrganizationsService_CheckMembership_notMember(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	member, err := client.Organizations.CheckMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.CheckMembership returned error: %+v", err)
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 400 response is interpreted as an actual error, and not simply
// as "false" like the above case of a 404
func TestOrganizationsService_CheckMembership_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	member, err := client.Organizations.CheckMembership("o", "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckMembership returned %+v, want %+v", member, want)
	}
}

func TestOrganizationsService_CheckMembership_invalidOrg(t *testing.T) {
	_, err := client.Organizations.CheckMembership("%", "u")
	testURLParseError(t, err)
}

func TestOrganizationsService_CheckPublicMembership(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	member, err := client.Organizations.CheckPublicMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.CheckPublicMembership returned error: %v", err)
	}
	if want := true; member != want {
		t.Errorf("Organizations.CheckPublicMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 404 response is interpreted as "false" and not an error
func TestOrganizationsService_CheckPublicMembership_notMember(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	member, err := client.Organizations.CheckPublicMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.CheckPublicMembership returned error: %v", err)
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckPublicMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 400 response is interpreted as an actual error, and not simply
// as "false" like the above case of a 404
func TestOrganizationsService_CheckPublicMembership_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	member, err := client.Organizations.CheckPublicMembership("o", "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckPublicMembership returned %+v, want %+v", member, want)
	}
}

func TestOrganizationsService_CheckPublicMembership_invalidOrg(t *testing.T) {
	_, err := client.Organizations.CheckPublicMembership("%", "u")
	testURLParseError(t, err)
}

func TestOrganizationsService_RemoveMember(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Organizations.RemoveMember("o", "u")
	if err != nil {
		t.Errorf("Organizations.RemoveMember returned error: %v", err)
	}
}

func TestOrganizationsService_RemoveMember_invalidOrg(t *testing.T) {
	err := client.Organizations.RemoveMember("%", "u")
	testURLParseError(t, err)
}

func TestOrganizationsService_ListTeams(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	teams, err := client.Organizations.ListTeams("o")
	if err != nil {
		t.Errorf("Organizations.ListTeams returned error: %v", err)
	}

	want := []Team{Team{ID: 1}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Organizations.ListTeams returned %+v, want %+v", teams, want)
	}
}

func TestOrganizationsService_ListTeams_invalidOrg(t *testing.T) {
	_, err := client.Organizations.ListTeams("%")
	testURLParseError(t, err)
}

func TestOrganizationsService_GetTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "name":"n", "url":"u", "slug": "s", "permission":"p"}`)
	})

	team, err := client.Organizations.GetTeam(1)
	if err != nil {
		t.Errorf("Organizations.GetTeam returned error: %v", err)
	}

	want := &Team{ID: 1, Name: "n", URL: "u", Slug: "s", Permission: "p"}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Organizations.GetTeam returned %+v, want %+v", team, want)
	}
}

func TestOrganizationsService_CreateTeam(t *testing.T) {
	setup()
	defer teardown()

	input := &Team{Name: "n"}

	mux.HandleFunc("/orgs/o/teams", func(w http.ResponseWriter, r *http.Request) {
		v := new(Team)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, err := client.Organizations.CreateTeam("o", input)
	if err != nil {
		t.Errorf("Organizations.CreateTeam returned error: %v", err)
	}

	want := &Team{ID: 1}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Organizations.CreateTeam returned %+v, want %+v", team, want)
	}
}

func TestOrganizationsService_CreateTeam_invalidOrg(t *testing.T) {
	_, err := client.Organizations.CreateTeam("%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_EditTeam(t *testing.T) {
	setup()
	defer teardown()

	input := &Team{Name: "n"}

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Team)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	team, err := client.Organizations.EditTeam(1, input)
	if err != nil {
		t.Errorf("Organizations.EditTeam returned error: %v", err)
	}

	want := &Team{ID: 1}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Organizations.EditTeam returned %+v, want %+v", team, want)
	}
}

func TestOrganizationsService_DeleteTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Organizations.DeleteTeam(1)
	if err != nil {
		t.Errorf("Organizations.DeleteTeam returned error: %v", err)
	}
}

func TestOrganizationsService_ListTeamMembers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	members, err := client.Organizations.ListTeamMembers(1)
	if err != nil {
		t.Errorf("Organizations.ListTeamMembers returned error: %v", err)
	}

	want := []User{User{ID: 1}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Organizations.ListTeamMembers returned %+v, want %+v", members, want)
	}
}

func TestOrganizationsService_CheckTeamMembership_true(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	member, err := client.Organizations.CheckTeamMembership(1, "u")
	if err != nil {
		t.Errorf("Organizations.CheckTeamMembership returned error: %v", err)
	}
	if want := true; member != want {
		t.Errorf("Organizations.CheckTeamMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 404 response is interpreted as "false" and not an error
func TestOrganizationsService_CheckTeamMembership_false(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	member, err := client.Organizations.CheckTeamMembership(1, "u")
	if err != nil {
		t.Errorf("Organizations.CheckTeamMembership returned error: %+v", err)
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckTeamMembership returned %+v, want %+v", member, want)
	}
}

// ensure that a 400 response is interpreted as an actual error, and not simply
// as "false" like the above case of a 404
func TestOrganizationsService_CheckTeamMembership_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	member, err := client.Organizations.CheckTeamMembership(1, "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; member != want {
		t.Errorf("Organizations.CheckTeamMembership returned %+v, want %+v", member, want)
	}
}

func TestOrganizationsService_CheckMembership_invalidUser(t *testing.T) {
	_, err := client.Organizations.CheckTeamMembership(1, "%")
	testURLParseError(t, err)
}

func TestOrganizationsService_AddTeamMember(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.AddTeamMember(1, "u")
	if err != nil {
		t.Errorf("Organizations.AddTeamMember returned error: %v", err)
	}
}

func TestOrganizationsService_AddTeamMember_invalidUser(t *testing.T) {
	err := client.Organizations.AddTeamMember(1, "%")
	testURLParseError(t, err)
}

func TestOrganizationsService_RemoveTeamMember(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.RemoveTeamMember(1, "u")
	if err != nil {
		t.Errorf("Organizations.RemoveTeamMember returned error: %v", err)
	}
}

func TestOrganizationsService_RemoveTeamMember_invalidUser(t *testing.T) {
	err := client.Organizations.RemoveTeamMember(1, "%")
	testURLParseError(t, err)
}

func TestOrganizationsService_PublicizeMembership(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.PublicizeMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.PublicizeMembership returned error: %v", err)
	}
}

func TestOrganizationsService_PublicizeMembership_invalidOrg(t *testing.T) {
	err := client.Organizations.PublicizeMembership("%", "u")
	testURLParseError(t, err)
}

func TestOrganizationsService_ConcealMembership(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/public_members/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.ConcealMembership("o", "u")
	if err != nil {
		t.Errorf("Organizations.ConcealMembership returned error: %v", err)
	}
}

func TestOrganizationsService_ConcealMembership_invalidOrg(t *testing.T) {
	err := client.Organizations.ConcealMembership("%", "u")
	testURLParseError(t, err)
}

func TestOrganizationsService_ListTeamRepos(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	members, err := client.Organizations.ListTeamRepos(1)
	if err != nil {
		t.Errorf("Organizations.ListTeamRepos returned error: %v", err)
	}

	want := []Repository{Repository{ID: 1}}
	if !reflect.DeepEqual(members, want) {
		t.Errorf("Organizations.ListTeamRepos returned %+v, want %+v", members, want)
	}
}

func TestOrganizationsService_CheckTeamRepo_true(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	managed, err := client.Organizations.CheckTeamRepo(1, "o", "r")
	if err != nil {
		t.Errorf("Organizations.CheckTeamRepo returned error: %v", err)
	}
	if want := true; managed != want {
		t.Errorf("Organizations.CheckTeamRepo returned %+v, want %+v", managed, want)
	}
}

func TestOrganizationsService_CheckTeamRepo_false(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	managed, err := client.Organizations.CheckTeamRepo(1, "o", "r")
	if err != nil {
		t.Errorf("Organizations.CheckTeamRepo returned error: %v", err)
	}
	if want := false; managed != want {
		t.Errorf("Organizations.CheckTeamRepo returned %+v, want %+v", managed, want)
	}
}

func TestOrganizationsService_CheckTeamRepo_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	managed, err := client.Organizations.CheckTeamRepo(1, "o", "r")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; managed != want {
		t.Errorf("Organizations.CheckTeamRepo returned %+v, want %+v", managed, want)
	}
}

func TestOrganizationsService_CheckTeamRepo_invalidOwner(t *testing.T) {
	_, err := client.Organizations.CheckTeamRepo(1, "%", "r")
	testURLParseError(t, err)
}

func TestOrganizationsService_AddTeamRepo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.AddTeamRepo(1, "o", "r")
	if err != nil {
		t.Errorf("Organizations.AddTeamRepo returned error: %v", err)
	}
}

func TestOrganizationsService_AddTeamRepo_noAccess(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(422)
	})

	err := client.Organizations.AddTeamRepo(1, "o", "r")
	if err == nil {
		t.Errorf("Expcted error to be returned")
	}
}

func TestOrganizationsService_AddTeamRepo_invalidOwner(t *testing.T) {
	err := client.Organizations.AddTeamRepo(1, "%", "r")
	testURLParseError(t, err)
}

func TestOrganizationsService_RemoveTeamRepo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Organizations.RemoveTeamRepo(1, "o", "r")
	if err != nil {
		t.Errorf("Organizations.RemoveTeamRepo returned error: %v", err)
	}
}

func TestOrganizationsService_RemoveTeamRepo_invalidOwner(t *testing.T) {
	err := client.Organizations.RemoveTeamRepo(1, "%", "r")
	testURLParseError(t, err)
}
