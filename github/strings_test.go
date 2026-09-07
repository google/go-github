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

	tests := []struct {
		in  any
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
		{new("foo"), `"foo"`},
		{new(123), `123`},
		{new(false), `false`},
		{
			//nolint:sliceofpointers
			[]*string{new("a"), new("b")},
			`["a" "b"]`,
		},

		// actual GitHub structs
		{
			Timestamp{time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)},
			`github.Timestamp{2006-01-02 15:04:05 +0000 UTC}`,
		},
		{
			&Timestamp{time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)},
			`github.Timestamp{2006-01-02 15:04:05 +0000 UTC}`,
		},
		{
			User{ID: new(int64(123)), Name: new("n")},
			`github.User{ID:123, Name:"n"}`,
		},
		{
			Repository{Owner: &User{ID: new(int64(123))}},
			`github.Repository{Owner:github.User{ID:123}}`,
		},
	}

	for i, tt := range tests {
		s := Stringify(tt.in)
		if s != tt.out {
			t.Errorf("%v. Stringify(%q) => %q, want %q", i, tt.in, s, tt.out)
		}
	}
}

func TestStringify_Primitives(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  any
		out string
	}{
		// Bool
		{true, "true"},
		{false, "false"},

		// Int variants
		{int(1), "1"},
		{int8(2), "2"},
		{int16(3), "3"},
		{int32(4), "4"},
		{int64(5), "5"},

		// Uint variants
		{uint(6), "6"},
		{uint8(7), "7"},
		{uint16(8), "8"},
		{uint32(9), "9"},
		{uint64(10), "10"},
		{uintptr(11), "11"},

		// Float variants (Precision Correctness)
		{float32(1.1), "1.1"},
		{float64(1.1), "1.1"},
		{float32(1.0000001), "1.0000001"},
		{float64(1.000000000000001), "1.000000000000001"},

		// Boundary Cases
		{int8(-128), "-128"},
		{int8(127), "127"},
		{uint64(18446744073709551615), "18446744073709551615"},

		// String Optimization
		{"hello", `"hello"`},
		{"", `""`},
	}

	for i, tt := range tests {
		s := Stringify(tt.in)
		if s != tt.out {
			t.Errorf("%v. Stringify(%T) => %q, want %q", i, tt.in, s, tt.out)
		}
	}
}

func TestStringify_BufferPool(t *testing.T) {
	t.Parallel()
	// Verify that concurrent usage of Stringify is safe and doesn't corrupt buffers.
	// While we can't easily verify reuse without exposing internal metrics,
	// we can verify correctness under load which implies proper Reset() handling.
	const goroutines = 10
	const iterations = 100

	errCh := make(chan error, goroutines)

	for range goroutines {
		go func() {
			for range iterations {
				// Use a mix of types to exercise different code paths
				s1 := Stringify(123)
				if s1 != "123" {
					errCh <- fmt.Errorf("got %q, want %q", s1, "123")
					return
				}

				s2 := Stringify("test")
				if s2 != `"test"` {
					errCh <- fmt.Errorf("got %q, want %q", s2, `"test"`)
					return
				}
			}
			errCh <- nil
		}()
	}

	for range goroutines {
		if err := <-errCh; err != nil {
			t.Error(err)
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
	tests := []struct {
		in  any
		out string
	}{
		{CodeResult{Name: new("n")}, `github.CodeResult{Name:"n"}`},
		{CommitAuthor{Name: new("n")}, `github.CommitAuthor{Name:"n"}`},
		{CommitFile{SHA: new("s")}, `github.CommitFile{SHA:"s"}`},
		{CommitStats{Total: new(1)}, `github.CommitStats{Total:1}`},
		{CommitsComparison{TotalCommits: new(1)}, `github.CommitsComparison{TotalCommits:1}`},
		{Commit{SHA: new("s")}, `github.Commit{SHA:"s"}`},
		{Event{ID: new("1")}, `github.Event{ID:"1"}`},
		{GistComment{ID: new(int64(1))}, `github.GistComment{ID:1}`},
		{GistFile{Size: new(1)}, `github.GistFile{Size:1}`},
		{Gist{ID: new("1")}, `github.Gist{ID:"1"}`},
		{GitObject{SHA: new("s")}, `github.GitObject{SHA:"s"}`},
		{Gitignore{Name: new("n")}, `github.Gitignore{Name:"n"}`},
		{Hook{ID: new(int64(1))}, `github.Hook{ID:1}`},
		{IssueComment{ID: new(int64(1))}, `github.IssueComment{ID:1}`},
		{Issue{Number: new(1)}, `github.Issue{Number:1}`},
		{SubIssue{ID: new(int64(1))}, `github.SubIssue{ID:1}`},
		{Key{ID: new(int64(1))}, `github.Key{ID:1}`},
		{Label{ID: 1, Name: "l"}, `github.Label{ID:1, URL:"", Name:"l", Color:"", Default:false, NodeID:""}`},
		{Organization{ID: new(int64(1))}, `github.Organization{ID:1}`},
		{PullRequestComment{ID: new(int64(1))}, `github.PullRequestComment{ID:1}`},
		{PullRequest{Number: new(1)}, `github.PullRequest{Number:1}`},
		{PullRequestReview{ID: new(int64(1))}, `github.PullRequestReview{ID:1}`},
		{DraftReviewComment{Position: new(1)}, `github.DraftReviewComment{Position:1}`},
		{PullRequestReviewRequest{Body: new("r")}, `github.PullRequestReviewRequest{Body:"r"}`},
		{PullRequestDismissReviewRequest{Message: "r"}, `github.PullRequestDismissReviewRequest{Message:"r"}`},
		{HeadCommit{SHA: new("s")}, `github.HeadCommit{SHA:"s"}`},
		{PushEvent{PushID: new(int64(1))}, `github.PushEvent{PushID:1}`},
		{Reference{Ref: new("r")}, `github.Reference{Ref:"r"}`},
		{ReleaseAsset{ID: new(int64(1))}, `github.ReleaseAsset{ID:1}`},
		{RepoStatus{ID: new(int64(1))}, `github.RepoStatus{ID:1}`},
		{RepositoryComment{ID: new(int64(1))}, `github.RepositoryComment{ID:1}`},
		{RepositoryCommit{SHA: new("s")}, `github.RepositoryCommit{SHA:"s"}`},
		{RepositoryContent{Name: new("n")}, `github.RepositoryContent{Name:"n"}`},
		{RepositoryRelease{ID: 1}, `github.RepositoryRelease{TagName:"", TargetCommitish:"", Draft:false, Prerelease:false, ID:1, CreatedAt:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, URL:"", HTMLURL:"", AssetsURL:"", UploadURL:"", NodeID:""}`},
		{Repository{ID: new(int64(1))}, `github.Repository{ID:1}`},
		{Team{ID: new(int64(1))}, `github.Team{ID:1}`},
		{TreeEntry{SHA: new("s")}, `github.TreeEntry{SHA:"s"}`},
		{Tree{SHA: new("s")}, `github.Tree{SHA:"s"}`},
		{User{ID: new(int64(1))}, `github.User{ID:1}`},
		{WebHookAuthor{Name: new("n")}, `github.CommitAuthor{Name:"n"}`},
		{WebHookCommit{ID: new("1")}, `github.HeadCommit{ID:"1"}`},
		{WebHookPayload{Ref: new("r")}, `github.PushEvent{Ref:"r"}`},
	}

	for i, tt := range tests {
		s := tt.in.(fmt.Stringer).String()
		if s != tt.out {
			t.Errorf("%v. String() => %q, want %q", i, tt.in, tt.out)
		}
	}
}

func TestStringify_Floats(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  any
		out string
	}{
		{float32(1.1), "1.1"},
		{float64(1.1), "1.1"},
		{float32(1.0000001), "1.0000001"},
		{struct{ F float32 }{1.1}, "{F:1.1}"},
	}

	for i, tt := range tests {
		s := Stringify(tt.in)
		if s != tt.out {
			t.Errorf("%v. Stringify(%v) = %q, want %q", i, tt.in, s, tt.out)
		}
	}
}
