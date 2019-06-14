// Copyright 2016 The go-github AUTHORS. All rights reserved.
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
	"time"
)

func TestAdminService_UpdateUserLDAPMapping(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &UserLDAPMapping{
		LDAPDN: String("uid=asdf,ou=users,dc=github,dc=com"),
	}

	mux.HandleFunc("/admin/ldap/users/u/mapping", func(w http.ResponseWriter, r *http.Request) {
		v := new(UserLDAPMapping)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1,"ldap_dn":"uid=asdf,ou=users,dc=github,dc=com"}`)
	})

	ctx := context.Background()
	mapping, _, err := client.Admin.UpdateUserLDAPMapping(ctx, "u", input)
	if err != nil {
		t.Errorf("Admin.UpdateUserLDAPMapping returned error: %v", err)
	}

	want := &UserLDAPMapping{
		ID:     Int64(1),
		LDAPDN: String("uid=asdf,ou=users,dc=github,dc=com"),
	}
	if !reflect.DeepEqual(mapping, want) {
		t.Errorf("Admin.UpdateUserLDAPMapping returned %+v, want %+v", mapping, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Admin.UpdateUserLDAPMapping(ctx, "u", input)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' UpdateUserLDAPMapping = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' UpdateUserLDAPMapping resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' UpdateUserLDAPMapping err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Admin.UpdateUserLDAPMapping(ctx, "u", input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now UpdateUserLDAPMapping = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now UpdateUserLDAPMapping resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now UpdateUserLDAPMapping err = nil, want error")
	}
}

func TestAdminService_UpdateTeamLDAPMapping(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &TeamLDAPMapping{
		LDAPDN: String("cn=Enterprise Ops,ou=teams,dc=github,dc=com"),
	}

	mux.HandleFunc("/admin/ldap/teams/1/mapping", func(w http.ResponseWriter, r *http.Request) {
		v := new(TeamLDAPMapping)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":1,"ldap_dn":"cn=Enterprise Ops,ou=teams,dc=github,dc=com"}`)
	})

	ctx := context.Background()
	mapping, _, err := client.Admin.UpdateTeamLDAPMapping(ctx, 1, input)
	if err != nil {
		t.Errorf("Admin.UpdateTeamLDAPMapping returned error: %v", err)
	}

	want := &TeamLDAPMapping{
		ID:     Int64(1),
		LDAPDN: String("cn=Enterprise Ops,ou=teams,dc=github,dc=com"),
	}
	if !reflect.DeepEqual(mapping, want) {
		t.Errorf("Admin.UpdateTeamLDAPMapping returned %+v, want %+v", mapping, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Admin.UpdateTeamLDAPMapping(ctx, 1, input)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' UpdateTeamLDAPMapping = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' UpdateTeamLDAPMapping resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' UpdateTeamLDAPMapping err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Admin.UpdateTeamLDAPMapping(ctx, 1, input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now UpdateTeamLDAPMapping = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now UpdateTeamLDAPMapping resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now UpdateTeamLDAPMapping err = nil, want error")
	}
}

func TestAdminService_TeamLDAPMapping_String(t *testing.T) {
	v := &TeamLDAPMapping{
		ID:              Int64(1),
		LDAPDN:          String("a"),
		URL:             String("b"),
		Name:            String("c"),
		Slug:            String("d"),
		Description:     String("e"),
		Privacy:         String("f"),
		Permission:      String("g"),
		MembersURL:      String("h"),
		RepositoriesURL: String("i"),
	}

	want := `github.TeamLDAPMapping{ID:1, LDAPDN:"a", URL:"b", Name:"c", Slug:"d", Description:"e", Privacy:"f", Permission:"g", MembersURL:"h", RepositoriesURL:"i"}`
	if got := v.String(); got != want {
		t.Errorf("TeamLDAPMapping.String = `%v`, want `%v`", got, want)
	}
}

func TestAdminService_UserLDAPMapping_String(t *testing.T) {
	v := &UserLDAPMapping{
		ID:                Int64(1),
		LDAPDN:            String("a"),
		Login:             String("b"),
		AvatarURL:         String("c"),
		GravatarID:        String("d"),
		Type:              String("e"),
		SiteAdmin:         Bool(true),
		URL:               String("f"),
		EventsURL:         String("g"),
		FollowingURL:      String("h"),
		FollowersURL:      String("i"),
		GistsURL:          String("j"),
		OrganizationsURL:  String("k"),
		ReceivedEventsURL: String("l"),
		ReposURL:          String("m"),
		StarredURL:        String("n"),
		SubscriptionsURL:  String("o"),
	}

	want := `github.UserLDAPMapping{ID:1, LDAPDN:"a", Login:"b", AvatarURL:"c", GravatarID:"d", Type:"e", SiteAdmin:true, URL:"f", EventsURL:"g", FollowingURL:"h", FollowersURL:"i", GistsURL:"j", OrganizationsURL:"k", ReceivedEventsURL:"l", ReposURL:"m", StarredURL:"n", SubscriptionsURL:"o"}`
	if got := v.String(); got != want {
		t.Errorf("UserLDAPMapping.String = `%v`, want `%v`", got, want)
	}
}
