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

func TestEditTitle_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditTitle{}, "{}")

	u := &EditTitle{
		From: String("EditTitleFrom"),
	}

	want := `{
		"from": "EditTitleFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBody_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditBody{}, "{}")

	u := &EditBody{
		From: String("EditBodyFrom"),
	}

	want := `{
		"from": "EditBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditBase_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditBase{}, "{}")

	u := &EditBase{
		Ref: &EditRef{
			From: String("EditRefFrom"),
		},
		SHA: &EditSHA{
			From: String("EditSHAFrom"),
		},
	}

	want := `{
		"ref": {
			"from": "EditRefFrom"
		},
		"sha": {
			"from": "EditSHAFrom"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestEditRef_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditRef{}, "{}")

	u := &EditRef{
		From: String("EditRefFrom"),
	}

	want := `{
		"from": "EditRefFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestEditSHA_Marshal(t *testing.T) {
	testJSONMarshal(t, &EditSHA{}, "{}")

	u := &EditSHA{
		From: String("EditSHAFrom"),
	}

	want := `{
		"from": "EditSHAFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectName_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectName{}, "{}")

	u := &ProjectName{
		From: String("ProjectNameFrom"),
	}

	want := `{
		"from": "ProjectNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectBody_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectBody{}, "{}")

	u := &ProjectBody{
		From: String("ProjectBodyFrom"),
	}

	want := `{
		"from": "ProjectBodyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectCardNote_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectCardNote{}, "{}")

	u := &ProjectCardNote{
		From: String("ProjectCardNoteFrom"),
	}

	want := `{
		"from": "ProjectCardNoteFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestProjectColumnName_Marshal(t *testing.T) {
	testJSONMarshal(t, &ProjectColumnName{}, "{}")

	u := &ProjectColumnName{
		From: String("ProjectColumnNameFrom"),
	}

	want := `{
		"from": "ProjectColumnNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamDescription_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamDescription{}, "{}")

	u := &TeamDescription{
		From: String("TeamDescriptionFrom"),
	}

	want := `{
		"from": "TeamDescriptionFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamName_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamName{}, "{}")

	u := &TeamName{
		From: String("TeamNameFrom"),
	}

	want := `{
		"from": "TeamNameFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPrivacy_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPrivacy{}, "{}")

	u := &TeamPrivacy{
		From: String("TeamPrivacyFrom"),
	}

	want := `{
		"from": "TeamPrivacyFrom"
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamRepository_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamRepository{}, "{}")

	u := &TeamRepository{
		Permissions: &TeamPermissions{
			From: &TeamPermissionsFrom{
				Admin: Bool(true),
				Pull:  Bool(true),
				Push:  Bool(true),
			},
		},
	}

	want := `{
		"permissions": {
			"from": {
				"admin": true,
				"pull": true,
				"push": true
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPermissions_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPermissions{}, "{}")

	u := &TeamPermissions{
		From: &TeamPermissionsFrom{
			Admin: Bool(true),
			Pull:  Bool(true),
			Push:  Bool(true),
		},
	}

	want := `{
		"from": {
			"admin": true,
			"pull": true,
			"push": true
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestTeamPermissionsFrom_Marshal(t *testing.T) {
	testJSONMarshal(t, &TeamPermissionsFrom{}, "{}")

	u := &TeamPermissionsFrom{
		Admin: Bool(true),
		Pull:  Bool(true),
		Push:  Bool(true),
	}

	want := `{
		"admin": true,
		"pull": true,
		"push": true
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoryVulnerabilityAlert_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryVulnerabilityAlert{}, "{}")

	u := &RepositoryVulnerabilityAlert{
		ID:                  Int64(1),
		AffectedRange:       String("ar"),
		AffectedPackageName: String("apn"),
		ExternalReference:   String("er"),
		ExternalIdentifier:  String("ei"),
		FixedIn:             String("fi"),
		Dismisser: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		DismissReason: String("dr"),
		DismissedAt:   &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"affected_range": "ar",
		"affected_package_name": "apn",
		"external_reference": "er",
		"external_identifier": "ei",
		"fixed_in": "fi",
		"dismisser": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"dismiss_reason": "dr",
		"dismissed_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
