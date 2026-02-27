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
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Payload did not panic but should have")
		}
	}()

	name := "UserEvent"
	body := json.RawMessage("[") // bogus JSON
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestPayload_NoPanic(t *testing.T) {
	t.Parallel()
	name := "UserEvent"
	body := json.RawMessage("{}")
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestEmptyEvent_NoPanic(t *testing.T) {
	t.Parallel()
	e := &Event{}
	if _, err := e.ParsePayload(); err == nil {
		t.Error("ParsePayload unexpectedly succeeded on empty event")
	}

	e = nil
	if _, err := e.ParsePayload(); err == nil {
		t.Error("ParsePayload unexpectedly succeeded on nil event")
	}
}

func TestEvent_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Event{}, "{}")

	l := make(map[string]any)
	l["key"] = "value"

	jsonMsg, _ := json.Marshal(&l)

	u := &Event{
		Type:       Ptr("t"),
		Public:     Ptr(false),
		RawPayload: (*json.RawMessage)(&jsonMsg),
		Repo: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("s"),
			Name: Ptr("n"),
		},
		Actor: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Org: &Organization{
			BillingEmail:                         Ptr("be"),
			Blog:                                 Ptr("b"),
			Company:                              Ptr("c"),
			Email:                                Ptr("e"),
			TwitterUsername:                      Ptr("tu"),
			Location:                             Ptr("loc"),
			Name:                                 Ptr("n"),
			Description:                          Ptr("d"),
			IsVerified:                           Ptr(true),
			HasOrganizationProjects:              Ptr(true),
			HasRepositoryProjects:                Ptr(true),
			DefaultRepoPermission:                Ptr("drp"),
			MembersCanCreateRepos:                Ptr(true),
			MembersCanCreateInternalRepos:        Ptr(true),
			MembersCanCreatePrivateRepos:         Ptr(true),
			MembersCanCreatePublicRepos:          Ptr(false),
			MembersAllowedRepositoryCreationType: Ptr("marct"),
			MembersCanCreatePages:                Ptr(true),
			MembersCanCreatePublicPages:          Ptr(false),
			MembersCanCreatePrivatePages:         Ptr(true),
		},
		CreatedAt: &Timestamp{referenceTime},
		ID:        Ptr("id"),
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

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}
