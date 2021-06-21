// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"
)

func TestEditChange_Marshal_TitleChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	TitleFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("TitleFrom"),
	}

	u := &EditChange{
		Title: &TitleFrom,
		Body:  nil,
		Base:  nil,
	}

	want := `{
		"title": {
			"from": "TitleFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_BodyChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	BodyFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BodyFrom"),
	}

	u := &EditChange{
		Title: nil,
		Body:  &BodyFrom,
		Base:  nil,
	}

	want := `{
		"body": {
			"from": "BodyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestEditChange_Marshal_BaseChange(t *testing.T) {
	testJSONMarshal(t, &EditChange{}, "{}")

	RefFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BaseRefFrom"),
	}

	SHAFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BaseSHAFrom"),
	}

	Base := struct {
		Ref *struct {
			From *string `json:"from,omitempty"`
		} `json:"ref,omitempty"`
		SHA *struct {
			From *string `json:"from,omitempty"`
		} `json:"sha,omitempty"`
	}{
		Ref: &RefFrom,
		SHA: &SHAFrom,
	}

	u := &EditChange{
		Title: nil,
		Body:  nil,
		Base:  &Base,
	}

	want := `{
		"base": {
			"ref": {
				"from": "BaseRefFrom"
			},
			"sha": {
				"from": "BaseSHAFrom"
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPage_Marshal(t *testing.T) {
	testJSONMarshal(t, &Page{}, "{}")

	u := &Page{
		PageName: String("p"),
		Title:    String("t"),
		Summary:  String("s"),
		Action:   String("a"),
		SHA:      String("s"),
		HTMLURL:  String("h"),
	}

	want := `{
		"page_name": "p",
		"title": "t",
		"summary": "s",
		"action": "a",
		"sha": "s",
		"html_url": "h"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamChange_Marshal_DescriptionChange(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	DescriptionFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("DescriptionFrom"),
	}

	u := &TeamChange{
		Description: &DescriptionFrom,
		Name:        nil,
		Privacy:     nil,
		Repository:  nil,
	}

	want := `{
		"description": {
			"from": "DescriptionFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	NameFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("NameFrom"),
	}

	u := &TeamChange{
		Description: nil,
		Name:        &NameFrom,
		Privacy:     nil,
		Repository:  nil,
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamChange_Marshal_PrivacyChange(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	PrivacyFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("PrivacyFrom"),
	}

	u := &TeamChange{
		Description: nil,
		Name:        nil,
		Privacy:     &PrivacyFrom,
		Repository:  nil,
	}

	want := `{
		"privacy": {
			"from": "PrivacyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamChange_Marshal_RepositoryChange(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	From := struct {
		Admin *bool `json:"admin,omitempty"`
		Pull  *bool `json:"pull,omitempty"`
		Push  *bool `json:"push,omitempty"`
	}{
		Admin: Bool(true),
		Pull:  Bool(true),
		Push:  Bool(false),
	}

	PermissionsFrom := struct {
		From *struct {
			Admin *bool `json:"admin,omitempty"`
			Pull  *bool `json:"pull,omitempty"`
			Push  *bool `json:"push,omitempty"`
		} `json:"from,omitempty"`
	}{
		From: &From,
	}

	Repository := struct {
		Permissions *struct {
			From *struct {
				Admin *bool `json:"admin,omitempty"`
				Pull  *bool `json:"pull,omitempty"`
				Push  *bool `json:"push,omitempty"`
			} `json:"from,omitempty"`
		} `json:"permissions,omitempty"`
	}{
		Permissions: &PermissionsFrom,
	}

	u := &TeamChange{
		Description: nil,
		Name:        nil,
		Privacy:     nil,
		Repository:  &Repository,
	}

	want := `{
		"repository": {
			"permissions": {
				"from": {
					"admin": true,
					"pull": true,
					"push": false
				}
			}
		  }
	}`

	testJSONMarshal(t, u, want)
}
