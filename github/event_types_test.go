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

func TestHeadCommit_Marshal(t *testing.T) {
	testJSONMarshal(t, &HeadCommit{}, "{}")

	u := &HeadCommit{
		Message: String("m"),
		Author: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		URL:       String("u"),
		Distinct:  Bool(true),
		SHA:       String("s"),
		ID:        String("id"),
		TreeID:    String("tid"),
		Timestamp: &Timestamp{referenceTime},
		Committer: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		Added:    []string{"a"},
		Removed:  []string{"r"},
		Modified: []string{"m"},
	}

	want := `{
		"message": "m",
		"author": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"url": "u",
		"distinct": true,
		"sha": "s",
		"id": "id",
		"tree_id": "tid",
		"timestamp": ` + referenceTimeStr + `,
		"committer": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"added": [
			"a"
		],
		"removed":  [
			"r"
		],
		"modified":  [
			"m"
		]
	}`

	testJSONMarshal(t, u, want)
}
