// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"testing"
	"time"
)

func TestStringify(t *testing.T) {
	t.Parallel()
	var nilPointer *string

	var tests = []struct {
		in  interface{}
		out string
	}{
		// basic types
		{"foo", `"foo"`},
		{123, `123`},
		{1.5, `1.5`},
		{false, `false`},
		{
			[]string{"a", "b"},
			`["a" "b"]`,
		},
		{
			struct {
				A []string
			}{nil},
			// nil slice is skipped
			`{}`,
		},
		{
			struct {
				A string
			}{"foo"},
			// structs not of a named type get no prefix
			`{A:"foo"}`,
		},

		// pointers
		{nilPointer, `<nil>`},
		{String("foo"), `"foo"`},
		{Int(123), `123`},
		{Bool(false), `false`},
		{
			[]*string{String("a"), String("b")},
			`["a" "b"]`,
		},

		// actual GitHub structs
		{
			Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)},
			`github.Timestamp{2006-01-02 15:04:05 +0000 UTC}`,
		},
		{
			&Timestamp{time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)},
			`github.Timestamp{2006-01-02 15:04:05 +0000 UTC}`,
		},
		{
			User{ID: Int64(123), Name: String("n")},
			`github.User{ID:123, Name:"n"}`,
		},
		{
			Repository{Owner: &User{ID: Int64(123)}},
			`github.Repository{Owner:github.User{ID:123}}`,
		},
	}

	for i, tt := range tests {
		s := Stringify(tt.in)
		if s != tt.out {
			t.Errorf("%d. Stringify(%q) => %q, want %q", i, tt.in, s, tt.out)
		}
	}
}

// Directly test the String() methods on various GitHub types. We don't do an
// exhaustive test of all the various field types, since TestStringify() above
// takes care of that. Rather, we just make sure that Stringify() is being
// used to build the strings, which we do by verifying that pointers are
// stringified as their underlying value.
func TestString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		in  interface{}
		out string
	}{
		{CodeResult{Name: String("n")}, `github.CodeResult{Name:"n"}`},
		{CommitAuthor{Name: String("n")}, `github.CommitAuthor{Name:"n"}`},
		{CommitFile{SHA: String("s")}, `github.CommitFile{SHA:"s"}`},
		{CommitStats{Total: Int(1)}, `github.CommitStats{Total:1}`},
		{CommitsComparison{TotalCommits: Int(1)}, `github.CommitsComparison{TotalCommits:1}`},
		{Commit{SHA: String("s")}, `github.Commit{SHA:"s"}`},
		{Event{ID: String("1")}, `github.Event{ID:"1"}`},
		{GistComment{ID: Int64(1)}, `github.GistComment{ID:1}`},
		{GistFile{Size: Int(1)}, `github.GistFile{Size:1}`},
		{Gist{ID: String("1")}, `github.Gist{ID:"1"}`},
		{GitObject{SHA: String("s")}, `github.GitObject{SHA:"s"}`},
		{Gitignore{Name: String("n")}, `github.Gitignore{Name:"n"}`},
		{Hook{ID: Int64(1)}, `github.Hook{ID:1}`},
		{IssueComment{ID: Int64(1)}, `github.IssueComment{ID:1}`},
		{Issue{Number: Int(1)}, `github.Issue{Number:1}`},
		{Key{ID: Int64(1)}, `github.Key{ID:1}`},
		{Label{ID: Int64(1), Name: String("l")}, `github.Label{ID:1, Name:"l"}`},
		{Organization{ID: Int64(1)}, `github.Organization{ID:1}`},
		{PullRequestComment{ID: Int64(1)}, `github.PullRequestComment{ID:1}`},
		{PullRequest{Number: Int(1)}, `github.PullRequest{Number:1}`},
		{PullRequestReview{ID: Int64(1)}, `github.PullRequestReview{ID:1}`},
		{DraftReviewComment{Position: Int(1)}, `github.DraftReviewComment{Position:1}`},
		{PullRequestReviewRequest{Body: String("r")}, `github.PullRequestReviewRequest{Body:"r"}`},
		{PullRequestReviewDismissalRequest{Message: String("r")}, `github.PullRequestReviewDismissalRequest{Message:"r"}`},
		{HeadCommit{SHA: String("s")}, `github.HeadCommit{SHA:"s"}`},
		{PushEvent{PushID: Int64(1)}, `github.PushEvent{PushID:1}`},
		{Reference{Ref: String("r")}, `github.Reference{Ref:"r"}`},
		{ReleaseAsset{ID: Int64(1)}, `github.ReleaseAsset{ID:1}`},
		{RepoStatus{ID: Int64(1)}, `github.RepoStatus{ID:1}`},
		{RepositoryComment{ID: Int64(1)}, `github.RepositoryComment{ID:1}`},
		{RepositoryCommit{SHA: String("s")}, `github.RepositoryCommit{SHA:"s"}`},
		{RepositoryContent{Name: String("n")}, `github.RepositoryContent{Name:"n"}`},
		{RepositoryRelease{ID: Int64(1)}, `github.RepositoryRelease{ID:1}`},
		{Repository{ID: Int64(1)}, `github.Repository{ID:1}`},
		{Team{ID: Int64(1)}, `github.Team{ID:1}`},
		{TreeEntry{SHA: String("s")}, `github.TreeEntry{SHA:"s"}`},
		{Tree{SHA: String("s")}, `github.Tree{SHA:"s"}`},
		{User{ID: Int64(1)}, `github.User{ID:1}`},
		{WebHookAuthor{Name: String("n")}, `github.CommitAuthor{Name:"n"}`},
		{WebHookCommit{ID: String("1")}, `github.HeadCommit{ID:"1"}`},
		{WebHookPayload{Ref: String("r")}, `github.PushEvent{Ref:"r"}`},
	}

	for i, tt := range tests {
		s := tt.in.(fmt.Stringer).String()
		if s != tt.out {
			t.Errorf("%d. String() => %q, want %q", i, tt.in, tt.out)
		}
	}
}
