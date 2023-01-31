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
	testJSONMarshal(t, &PullRequestThread{}, "{}")

	createdAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	updatedAt := Timestamp{time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)}
	reactions := &Reactions{
		TotalCount: Int(1),
		PlusOne:    Int(1),
		MinusOne:   Int(0),
		Laugh:      Int(0),
		Confused:   Int(0),
		Heart:      Int(0),
		Hooray:     Int(0),
		Rocket:     Int(0),
		Eyes:       Int(0),
		URL:        String("u"),
	}
	user := &User{
		Login:       String("ll"),
		ID:          Int64(123),
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
	}
	comment := &PullRequestComment{
		ID:                  Int64(10),
		InReplyTo:           Int64(8),
		Body:                String("Test comment"),
		Path:                String("file1.txt"),
		DiffHunk:            String("@@ -16,33 +16,40 @@ fmt.Println()"),
		PullRequestReviewID: Int64(42),
		Position:            Int(1),
		OriginalPosition:    Int(4),
		StartLine:           Int(2),
		Line:                Int(3),
		OriginalLine:        Int(2),
		OriginalStartLine:   Int(2),
		Side:                String("RIGHT"),
		StartSide:           String("LEFT"),
		CommitID:            String("ab"),
		OriginalCommitID:    String("9c"),
		User:                user,
		Reactions:           reactions,
		CreatedAt:           &createdAt,
		UpdatedAt:           &updatedAt,
		URL:                 String("pullrequestcommentUrl"),
		HTMLURL:             String("pullrequestcommentHTMLUrl"),
		PullRequestURL:      String("pullrequestcommentPullRequestURL"),
	}

	u := &PullRequestThread{
		ID:       Int64(1),
		NodeID:   String("nid"),
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
