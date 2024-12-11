// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"
	"time"
)

func TestPullRequestThread_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestThread{}, "{}")

	createdAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	updatedAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	reactions := &Reactions{
		TotalCount: Ptr(1),
		PlusOne:    Ptr(1),
		MinusOne:   Ptr(0),
		Laugh:      Ptr(0),
		Confused:   Ptr(0),
		Heart:      Ptr(0),
		Hooray:     Ptr(0),
		Rocket:     Ptr(0),
		Eyes:       Ptr(0),
		URL:        Ptr("u"),
	}
	user := &User{
		Login:       Ptr("ll"),
		ID:          Ptr(int64(123)),
		AvatarURL:   Ptr("a"),
		GravatarID:  Ptr("g"),
		Name:        Ptr("n"),
		Company:     Ptr("c"),
		Blog:        Ptr("b"),
		Location:    Ptr("l"),
		Email:       Ptr("e"),
		Hireable:    Ptr(true),
		PublicRepos: Ptr(1),
		Followers:   Ptr(1),
		Following:   Ptr(1),
		CreatedAt:   &Timestamp{referenceTime},
		URL:         Ptr("u"),
	}
	comment := &PullRequestComment{
		ID:                  Ptr(int64(10)),
		InReplyTo:           Ptr(int64(8)),
		Body:                Ptr("Test comment"),
		Path:                Ptr("file1.txt"),
		DiffHunk:            Ptr("@@ -16,33 +16,40 @@ fmt.Println()"),
		PullRequestReviewID: Ptr(int64(42)),
		Position:            Ptr(1),
		OriginalPosition:    Ptr(4),
		StartLine:           Ptr(2),
		Line:                Ptr(3),
		OriginalLine:        Ptr(2),
		OriginalStartLine:   Ptr(2),
		Side:                Ptr("RIGHT"),
		StartSide:           Ptr("LEFT"),
		CommitID:            Ptr("ab"),
		OriginalCommitID:    Ptr("9c"),
		User:                user,
		Reactions:           reactions,
		CreatedAt:           &createdAt,
		UpdatedAt:           &updatedAt,
		URL:                 Ptr("pullrequestcommentUrl"),
		HTMLURL:             Ptr("pullrequestcommentHTMLUrl"),
		PullRequestURL:      Ptr("pullrequestcommentPullRequestURL"),
	}

	u := &PullRequestThread{
		ID:       Ptr(int64(1)),
		NodeID:   Ptr("nid"),
		Comments: []*PullRequestComment{comment, comment},
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"comments": [
			{
				"id": 10,
				"in_reply_to_id": 8,
				"body": "Test comment",
				"path": "file1.txt",
				"diff_hunk": "@@ -16,33 +16,40 @@ fmt.Println()",
				"pull_request_review_id": 42,
				"position": 1,
				"original_position": 4,
				"start_line": 2,
				"line": 3,
				"original_line": 2,
				"original_start_line": 2,
				"side": "RIGHT",
				"start_side": "LEFT",
				"commit_id": "ab",
				"original_commit_id": "9c",
				"user": {
					"login": "ll",
					"id": 123,
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
				"reactions": {
					"total_count": 1,
					"+1": 1,
					"-1": 0,
					"laugh": 0,
					"confused": 0,
					"heart": 0,
					"hooray": 0,
					"rocket": 0,
					"eyes": 0,
					"url": "u"
				},
				"created_at": "2002-02-10T15:30:00Z",
				"updated_at": "2002-02-10T15:30:00Z",
				"url": "pullrequestcommentUrl",
				"html_url": "pullrequestcommentHTMLUrl",
				"pull_request_url": "pullrequestcommentPullRequestURL"
			},
			{
				"id": 10,
				"in_reply_to_id": 8,
				"body": "Test comment",
				"path": "file1.txt",
				"diff_hunk": "@@ -16,33 +16,40 @@ fmt.Println()",
				"pull_request_review_id": 42,
				"position": 1,
				"original_position": 4,
				"start_line": 2,
				"line": 3,
				"original_line": 2,
				"original_start_line": 2,
				"side": "RIGHT",
				"start_side": "LEFT",
				"commit_id": "ab",
				"original_commit_id": "9c",
				"user": {
					"login": "ll",
					"id": 123,
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
				"reactions": {
					"total_count": 1,
					"+1": 1,
					"-1": 0,
					"laugh": 0,
					"confused": 0,
					"heart": 0,
					"hooray": 0,
					"rocket": 0,
					"eyes": 0,
					"url": "u"
				},
				"created_at": "2002-02-10T15:30:00Z",
				"updated_at": "2002-02-10T15:30:00Z",
				"url": "pullrequestcommentUrl",
				"html_url": "pullrequestcommentHTMLUrl",
				"pull_request_url": "pullrequestcommentPullRequestURL"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
