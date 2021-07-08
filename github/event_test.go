// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
)

func TestPayload_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Payload did not panic but should have")
		}
	}()

	name := "UserEvent"
	body := json.RawMessage("[") // bogus JSON
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestPayload_NoPanic(t *testing.T) {
	name := "UserEvent"
	body := json.RawMessage("{}")
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &Event{}, "{}")

	l := make(map[string]interface{})
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &Event{
		Type:       String("t"),
		Public:     Bool(false),
		RawPayload: (*json.RawMessage)(&jsonMsg),
		Repo: &Repository{
			ID:   Int64(1),
			URL:  String("s"),
			Name: String("n"),
		},
		Actor: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Org: &Organization{
			BillingEmail:                         String("be"),
			Blog:                                 String("b"),
			Company:                              String("c"),
			Email:                                String("e"),
			TwitterUsername:                      String("tu"),
			Location:                             String("loc"),
			Name:                                 String("n"),
			Description:                          String("d"),
			IsVerified:                           Bool(true),
			HasOrganizationProjects:              Bool(true),
			HasRepositoryProjects:                Bool(true),
			DefaultRepoPermission:                String("drp"),
			MembersCanCreateRepos:                Bool(true),
			MembersCanCreateInternalRepos:        Bool(true),
			MembersCanCreatePrivateRepos:         Bool(true),
			MembersCanCreatePublicRepos:          Bool(false),
			MembersAllowedRepositoryCreationType: String("marct"),
			MembersCanCreatePages:                Bool(true),
			MembersCanCreatePublicPages:          Bool(false),
			MembersCanCreatePrivatePages:         Bool(true),
		},
		CreatedAt: &referenceTime,
		ID:        String("id"),
	}

	want := `{
		"type": "t",
		"public": false,
		"payload": {
			"key": "value"
		},
		"repo": {
			"id": 1,
			"name": "n",
			"url": "s"
		},
		"actor": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"org": {
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "loc",
			"email": "e",
			"twitter_username": "tu",
			"description": "d",
			"billing_email": "be",
			"is_verified": true,
			"has_organization_projects": true,
			"has_repository_projects": true,
			"default_repository_permission": "drp",
			"members_can_create_repositories": true,
			"members_can_create_public_repositories": false,
			"members_can_create_private_repositories": true,
			"members_can_create_internal_repositories": true,
			"members_allowed_repository_creation_type": "marct",
			"members_can_create_pages": true,
			"members_can_create_public_pages": false,
			"members_can_create_private_pages": true
		},
		"created_at": ` + referenceTimeStr + `,
		"id": "id"
	}`

	testJSONMarshal(t, u, want)
}
