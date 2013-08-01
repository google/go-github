// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

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
