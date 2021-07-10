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

	u := &EditChange{
		Title: &EditTitle{
			From: String("TitleFrom"),
		},
		Body: nil,
		Base: nil,
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

	u := &EditChange{
		Title: nil,
		Body: &EditBody{
			From: String("BodyFrom"),
		},
		Base: nil,
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

	Base := EditBase{
		Ref: &EditRef{
			From: String("BaseRefFrom"),
		},
		SHA: &EditSHA{
			From: String("BaseSHAFrom"),
		},
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

func TestProjectChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: &ProjectName{From: String("NameFrom")},
		Body: nil,
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectChange_Marshal_BodyChange(t *testing.T) {
	testJSONMarshal(t, &ProjectChange{}, "{}")

	u := &ProjectChange{
		Name: nil,
		Body: &ProjectBody{From: String("BodyFrom")},
	}

	want := `{
		"body": {
			"from": "BodyFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardChange_Marshal_NoteChange(t *testing.T) {
	testJSONMarshal(t, &ProjectCardChange{}, "{}")

	u := &ProjectCardChange{
		Note: &ProjectCardNote{From: String("NoteFrom")},
	}

	want := `{
		"note": {
			"from": "NoteFrom"
		  }
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnChange{}, "{}")

	u := &ProjectColumnChange{
		Name: &ProjectColumnName{From: String("NameFrom")},
	}

	want := `{
		"name": {
			"from": "NameFrom"
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

func TestTeamChange_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamChange{}, "{}")

	u := &TeamChange{
		Description: &TeamDescription{
			From: String("DescriptionFrom"),
		},
		Name: &TeamName{
			From: String("NameFrom"),
		},
		Privacy: &TeamPrivacy{
			From: String("PrivacyFrom"),
		},
		Repository: &TeamRepository{
			Permissions: &TeamPermissions{
				From: &TeamPermissionsFrom{
					Admin: Bool(false),
					Pull:  Bool(false),
					Push:  Bool(false),
				},
			},
		},
	}

	want := `{
		"description": {
			"from": "DescriptionFrom"
		},
		"name": {
			"from": "NameFrom"
		},
		"privacy": {
			"from": "PrivacyFrom"
		},
		"repository": {
			"permissions": {
				"from": {
					"admin": false,
					"pull": false,
					"push": false
				}
			}
		}
	}`

	testJSONMarshal(t, u, want)
}
