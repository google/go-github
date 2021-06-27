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
	"testing"

	"github.com/google/go-cmp/cmp"
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
		if !cmp.Equal(v, input) {
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
	if !cmp.Equal(mapping, want) {
		t.Errorf("Admin.UpdateUserLDAPMapping returned %+v, want %+v", mapping, want)
	}

	const methodName = "UpdateUserLDAPMapping"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Admin.UpdateUserLDAPMapping(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Admin.UpdateUserLDAPMapping(ctx, "u", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		if !cmp.Equal(v, input) {
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
	if !cmp.Equal(mapping, want) {
		t.Errorf("Admin.UpdateTeamLDAPMapping returned %+v, want %+v", mapping, want)
	}

	const methodName = "UpdateTeamLDAPMapping"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Admin.UpdateTeamLDAPMapping(ctx, -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Admin.UpdateTeamLDAPMapping(ctx, 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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

func TestTeamLDAPMapping_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamLDAPMapping{}, "{}")

	u := &TeamLDAPMapping{
		ID:              Int64(1),
		LDAPDN:          String("ldapdn"),
		URL:             String("u"),
		Name:            String("n"),
		Slug:            String("s"),
		Description:     String("d"),
		Privacy:         String("p"),
		Permission:      String("per"),
		MembersURL:      String("mu"),
		RepositoriesURL: String("ru"),
	}

	want := `{
		"id": 1,
		"ldap_dn": "ldapdn",
		"url": "u",
		"name": "n",
		"slug": "s",
		"description": "d",
		"privacy": "p",
		"permission": "per",
		"members_url": "mu",
		"repositories_url": "ru"
	}`

	testJSONMarshal(t, u, want)
}

func TestUserLDAPMapping_Marshal(t *testing.T) {
	testJSONMarshal(t, &UserLDAPMapping{}, "{}")

	u := &UserLDAPMapping{
		ID:                Int64(1),
		LDAPDN:            String("ldapdn"),
		Login:             String("l"),
		AvatarURL:         String("au"),
		GravatarID:        String("gi"),
		Type:              String("t"),
		SiteAdmin:         Bool(true),
		URL:               String("u"),
		EventsURL:         String("eu"),
		FollowingURL:      String("fu"),
		FollowersURL:      String("fu"),
		GistsURL:          String("gu"),
		OrganizationsURL:  String("ou"),
		ReceivedEventsURL: String("reu"),
		ReposURL:          String("ru"),
		StarredURL:        String("su"),
		SubscriptionsURL:  String("subu"),
	}

	want := `{
		"id": 1,
		"ldap_dn": "ldapdn",
		"login": "l",
		"avatar_url": "au",
		"gravatar_id": "gi",
		"type": "t",
		"site_admin": true,
		"url": "u",
		"events_url": "eu",
		"following_url": "fu",
		"followers_url": "fu",
		"gists_url": "gu",
		"organizations_url": "ou",
		"received_events_url": "reu",
		"repos_url": "ru",
		"starred_url": "su",
		"subscriptions_url": "subu"
	}`

	testJSONMarshal(t, u, want)
}

func TestEnterprise_Marshal(t *testing.T) {
	testJSONMarshal(t, &Enterprise{}, "{}")

	u := &Enterprise{
		ID:          Int(1),
		Slug:        String("s"),
		Name:        String("n"),
		NodeID:      String("nid"),
		AvatarURL:   String("au"),
		Description: String("d"),
		WebsiteURL:  String("wu"),
		HTMLURL:     String("hu"),
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"slug": "s",
		"name": "n",
		"node_id": "nid",
		"avatar_url": "au",
		"description": "d",
		"website_url": "wu",
		"html_url": "hu",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
