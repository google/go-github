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

func TestProjectChange_Marshal_NameChange(t *testing.T) {
	testJSONMarshal(t, &ProjectChange{}, "{}")

	NameFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("NameFrom"),
	}

	u := &ProjectChange{
		Name: &NameFrom,
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

	BodyFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("BodyFrom"),
	}

	u := &ProjectChange{
		Name: nil,
		Body: &BodyFrom,
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

	NoteFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("NoteFrom"),
	}

	u := &ProjectCardChange{
		Note: &NoteFrom,
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

	NameFrom := struct {
		From *string `json:"from,omitempty"`
	}{
		From: String("NameFrom"),
	}

	u := &ProjectColumnChange{
		Name: &NameFrom,
	}

	want := `{
		"name": {
			"from": "NameFrom"
		  }
	}`

	testJSONMarshal(t, u, want)

}
