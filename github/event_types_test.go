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

func TestPushEventRepository_Marshal(t *testing.T) {
	testJSONMarshal(t, &PushEventRepository{}, "{}")

	u := &PushEventRepository{
		ID:       Int64(1),
		NodeID:   String("nid"),
		Name:     String("n"),
		FullName: String("fn"),
		Owner: &User{
			Login:       String("l"),
			ID:          Int64(1),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		},
		Private:         Bool(true),
		Description:     String("d"),
		Fork:            Bool(true),
		CreatedAt:       &Timestamp{referenceTime},
		PushedAt:        &Timestamp{referenceTime},
		UpdatedAt:       &Timestamp{referenceTime},
		Homepage:        String("h"),
		PullsURL:        String("p"),
		Size:            Int(1),
		StargazersCount: Int(1),
		WatchersCount:   Int(1),
		Language:        String("l"),
		HasIssues:       Bool(true),
		HasDownloads:    Bool(true),
		HasWiki:         Bool(true),
		HasPages:        Bool(true),
		ForksCount:      Int(1),
		Archived:        Bool(true),
		Disabled:        Bool(true),
		OpenIssuesCount: Int(1),
		DefaultBranch:   String("d"),
		MasterBranch:    String("m"),
		Organization:    String("o"),
		URL:             String("u"),
		ArchiveURL:      String("a"),
		HTMLURL:         String("h"),
		StatusesURL:     String("s"),
		GitURL:          String("g"),
		SSHURL:          String("s"),
		CloneURL:        String("c"),
		SVNURL:          String("s"),
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"name": "n",
		"full_name": "fn",
		"owner": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"private": true,
		"description": "d",
		"fork": true,
		"created_at": ` + referenceTimeStr + `,
		"pushed_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"homepage": "h",
		"pulls_url": "p",
		"size": 1,
		"stargazers_count": 1,
		"watchers_count": 1,
		"language": "l",
		"has_issues": true,
		"has_downloads": true,
		"has_wiki": true,
		"has_pages": true,
		"forks_count": 1,
		"archived": true,
		"disabled": true,
		"open_issues_count": 1,
		"default_branch": "d",
		"master_branch": "m",
		"organization": "o",
		"url": "u",
		"archive_url": "a",
		"html_url": "h",
		"statuses_url": "s",
		"git_url": "g",
		"ssh_url": "s",
		"clone_url": "c",
		"svn_url": "s"
	}`

	testJSONMarshal(t, u, want)
}
