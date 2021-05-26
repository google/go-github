// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListCommits(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	// given
	mux.HandleFunc("/repos/o/r/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r,
			values{
				"sha":    "s",
				"path":   "p",
				"author": "a",
				"since":  "2013-08-01T00:00:00Z",
				"until":  "2013-09-03T00:00:00Z",
			})
		fmt.Fprintf(w, `[{"sha": "s"}]`)
	})

	opt := &CommitsListOptions{
		SHA:    "s",
		Path:   "p",
		Author: "a",
		Since:  time.Date(2013, time.August, 1, 0, 0, 0, 0, time.UTC),
		Until:  time.Date(2013, time.September, 3, 0, 0, 0, 0, time.UTC),
	}
	ctx := context.Background()
	commits, _, err := client.Repositories.ListCommits(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListCommits returned error: %v", err)
	}

	want := []*RepositoryCommit{{SHA: String("s")}}
	if !cmp.Equal(commits, want) {
		t.Errorf("Repositories.ListCommits returned %+v, want %+v", commits, want)
	}

	const methodName = "ListCommits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCommits(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCommits(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetCommit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{
		  "sha": "s",
		  "commit": { "message": "m" },
		  "author": { "login": "l" },
		  "committer": { "login": "l" },
		  "parents": [ { "sha": "s" } ],
		  "stats": { "additions": 104, "deletions": 4, "total": 108 },
		  "files": [
		    {
		      "filename": "f",
		      "additions": 10,
		      "deletions": 2,
		      "changes": 12,
		      "status": "s",
		      "patch": "p",
		      "blob_url": "b",
		      "raw_url": "r",
		      "contents_url": "c"
		    }
		  ]
		}`)
	})

	ctx := context.Background()
	commit, _, err := client.Repositories.GetCommit(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Repositories.GetCommit returned error: %v", err)
	}

	want := &RepositoryCommit{
		SHA: String("s"),
		Commit: &Commit{
			Message: String("m"),
		},
		Author: &User{
			Login: String("l"),
		},
		Committer: &User{
			Login: String("l"),
		},
		Parents: []*Commit{
			{
				SHA: String("s"),
			},
		},
		Stats: &CommitStats{
			Additions: Int(104),
			Deletions: Int(4),
			Total:     Int(108),
		},
		Files: []*CommitFile{
			{
				Filename:    String("f"),
				Additions:   Int(10),
				Deletions:   Int(2),
				Changes:     Int(12),
				Status:      String("s"),
				Patch:       String("p"),
				BlobURL:     String("b"),
				RawURL:      String("r"),
				ContentsURL: String("c"),
			},
		},
	}
	if !cmp.Equal(commit, want) {
		t.Errorf("Repositories.GetCommit returned \n%+v, want \n%+v", commit, want)
	}

	const methodName = "GetCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCommit(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCommit(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetCommitRaw_diff(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	const rawStr = "@@diff content"

	mux.HandleFunc("/repos/o/r/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3Diff)
		fmt.Fprint(w, rawStr)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommitRaw(ctx, "o", "r", "s", RawOptions{Type: Diff})
	if err != nil {
		t.Fatalf("Repositories.GetCommitRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("Repositories.GetCommitRaw returned %s want %s", got, want)
	}

	const methodName = "GetCommitRaw"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCommitRaw(ctx, "\n", "\n", "\n", RawOptions{Type: Diff})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCommitRaw(ctx, "o", "r", "s", RawOptions{Type: Diff})
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want ''", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetCommitRaw_patch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	const rawStr = "@@patch content"

	mux.HandleFunc("/repos/o/r/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3Patch)
		fmt.Fprint(w, rawStr)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommitRaw(ctx, "o", "r", "s", RawOptions{Type: Patch})
	if err != nil {
		t.Fatalf("Repositories.GetCommitRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("Repositories.GetCommitRaw returned %s want %s", got, want)
	}
}

func TestRepositoriesService_GetCommitRaw_invalid(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.GetCommitRaw(ctx, "o", "r", "s", RawOptions{100})
	if err == nil {
		t.Fatal("Repositories.GetCommitRaw should return error")
	}
	if !strings.Contains(err.Error(), "unsupported raw type") {
		t.Error("Repositories.GetCommitRaw should return unsupported raw type error")
	}
}

func TestRepositoriesService_GetCommitSHA1(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	const sha1 = "01234abcde"

	mux.HandleFunc("/repos/o/r/commits/master", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)

		fmt.Fprintf(w, sha1)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommitSHA1(ctx, "o", "r", "master", "")
	if err != nil {
		t.Errorf("Repositories.GetCommitSHA1 returned error: %v", err)
	}

	want := sha1
	if got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}

	mux.HandleFunc("/repos/o/r/commits/tag", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)
		testHeader(t, r, "If-None-Match", `"`+sha1+`"`)

		w.WriteHeader(http.StatusNotModified)
	})

	got, _, err = client.Repositories.GetCommitSHA1(ctx, "o", "r", "tag", sha1)
	if err == nil {
		t.Errorf("Expected HTTP 304 response")
	}

	want = ""
	if got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}

	const methodName = "GetCommitSHA1"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCommitSHA1(ctx, "\n", "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCommitSHA1(ctx, "o", "r", "master", "")
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want ''", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_NonAlphabetCharacter_GetCommitSHA1(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	const sha1 = "01234abcde"

	mux.HandleFunc("/repos/o/r/commits/master%20hash", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)

		fmt.Fprintf(w, sha1)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommitSHA1(ctx, "o", "r", "master%20hash", "")
	if err != nil {
		t.Errorf("Repositories.GetCommitSHA1 returned error: %v", err)
	}

	if want := sha1; got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}

	mux.HandleFunc("/repos/o/r/commits/tag", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)
		testHeader(t, r, "If-None-Match", `"`+sha1+`"`)

		w.WriteHeader(http.StatusNotModified)
	})

	got, _, err = client.Repositories.GetCommitSHA1(ctx, "o", "r", "tag", sha1)
	if err == nil {
		t.Errorf("Expected HTTP 304 response")
	}

	if want := ""; got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}
}

func TestRepositoriesService_TrailingPercent_GetCommitSHA1(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	const sha1 = "01234abcde"

	mux.HandleFunc("/repos/o/r/commits/comm%", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)

		fmt.Fprintf(w, sha1)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommitSHA1(ctx, "o", "r", "comm%", "")
	if err != nil {
		t.Errorf("Repositories.GetCommitSHA1 returned error: %v", err)
	}

	if want := sha1; got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}

	mux.HandleFunc("/repos/o/r/commits/tag", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3SHA)
		testHeader(t, r, "If-None-Match", `"`+sha1+`"`)

		w.WriteHeader(http.StatusNotModified)
	})

	got, _, err = client.Repositories.GetCommitSHA1(ctx, "o", "r", "tag", sha1)
	if err == nil {
		t.Errorf("Expected HTTP 304 response")
	}

	if want := ""; got != want {
		t.Errorf("Repositories.GetCommitSHA1 = %v, want %v", got, want)
	}
}

func TestRepositoriesService_CompareCommits(t *testing.T) {
	testCases := []struct {
		base string
		head string
	}{
		{base: "b", head: "h"},
		{base: "123base", head: "head123"},
		{base: "`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+123base", head: "head123`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+"},
	}

	for _, sample := range testCases {
		client, mux, _, teardown := setup()

		base := sample.base
		head := sample.head
		escapedBase := url.QueryEscape(base)
		escapedHead := url.QueryEscape(head)

		pattern := fmt.Sprintf("/repos/o/r/compare/%v...%v", base, head)

		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprintf(w, `{
		  "base_commit": {
		    "sha": "s",
		    "commit": {
		      "author": { "name": "n" },
		      "committer": { "name": "n" },
		      "message": "m",
		      "tree": { "sha": "t" }
		    },
		    "author": { "login": "l" },
		    "committer": { "login": "l" },
		    "parents": [ { "sha": "s" } ]
		  },
		  "status": "s",
		  "ahead_by": 1,
		  "behind_by": 2,
		  "total_commits": 1,
		  "commits": [
		    {
		      "sha": "s",
		      "commit": { "author": { "name": "n" } },
		      "author": { "login": "l" },
		      "committer": { "login": "l" },
		      "parents": [ { "sha": "s" } ]
		    }
		  ],
		  "files": [ { "filename": "f" } ],
		  "html_url":      "https://github.com/o/r/compare/%[1]v...%[2]v",
		  "permalink_url": "https://github.com/o/r/compare/o:bbcd538c8e72b8c175046e27cc8f907076331401...o:0328041d1152db8ae77652d1618a02e57f745f17",
		  "diff_url":      "https://github.com/o/r/compare/%[1]v...%[2]v.diff",
		  "patch_url":     "https://github.com/o/r/compare/%[1]v...%[2]v.patch",
		  "url":           "https://api.github.com/repos/o/r/compare/%[1]v...%[2]v"
		}`, escapedBase, escapedHead)
		})

		ctx := context.Background()
		got, _, err := client.Repositories.CompareCommits(ctx, "o", "r", base, head)
		if err != nil {
			t.Errorf("Repositories.CompareCommits returned error: %v", err)
		}

		want := &CommitsComparison{
			BaseCommit: &RepositoryCommit{
				SHA: String("s"),
				Commit: &Commit{
					Author:    &CommitAuthor{Name: String("n")},
					Committer: &CommitAuthor{Name: String("n")},
					Message:   String("m"),
					Tree:      &Tree{SHA: String("t")},
				},
				Author:    &User{Login: String("l")},
				Committer: &User{Login: String("l")},
				Parents: []*Commit{
					{
						SHA: String("s"),
					},
				},
			},
			Status:       String("s"),
			AheadBy:      Int(1),
			BehindBy:     Int(2),
			TotalCommits: Int(1),
			Commits: []*RepositoryCommit{
				{
					SHA: String("s"),
					Commit: &Commit{
						Author: &CommitAuthor{Name: String("n")},
					},
					Author:    &User{Login: String("l")},
					Committer: &User{Login: String("l")},
					Parents: []*Commit{
						{
							SHA: String("s"),
						},
					},
				},
			},
			Files: []*CommitFile{
				{
					Filename: String("f"),
				},
			},
			HTMLURL:      String(fmt.Sprintf("https://github.com/o/r/compare/%v...%v", escapedBase, escapedHead)),
			PermalinkURL: String("https://github.com/o/r/compare/o:bbcd538c8e72b8c175046e27cc8f907076331401...o:0328041d1152db8ae77652d1618a02e57f745f17"),
			DiffURL:      String(fmt.Sprintf("https://github.com/o/r/compare/%v...%v.diff", escapedBase, escapedHead)),
			PatchURL:     String(fmt.Sprintf("https://github.com/o/r/compare/%v...%v.patch", escapedBase, escapedHead)),
			URL:          String(fmt.Sprintf("https://api.github.com/repos/o/r/compare/%v...%v", escapedBase, escapedHead)),
		}

		if !cmp.Equal(got, want) {
			t.Errorf("Repositories.CompareCommits returned \n%+v, want \n%+v", got, want)
		}

		const methodName = "CompareCommits"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.Repositories.CompareCommits(ctx, "\n", "\n", "\n", "\n")
			return err
		})

		testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
			got, resp, err := client.Repositories.CompareCommits(ctx, "o", "r", base, head)
			if got != nil {
				t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
			}
			return resp, err
		})

		teardown()
	}
}

func TestRepositoriesService_CompareCommitsRaw_diff(t *testing.T) {
	testCases := []struct {
		base string
		head string
	}{
		{base: "b", head: "h"},
		{base: "123base", head: "head123"},
		{base: "`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+123base", head: "head123`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+"},
	}

	for _, sample := range testCases {
		client, mux, _, teardown := setup()

		base := sample.base
		head := sample.head
		pattern := fmt.Sprintf("/repos/o/r/compare/%v...%v", base, head)
		const rawStr = "@@diff content"

		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "Accept", mediaTypeV3Diff)
			fmt.Fprint(w, rawStr)
		})

		ctx := context.Background()
		got, _, err := client.Repositories.CompareCommitsRaw(ctx, "o", "r", base, head, RawOptions{Type: Diff})
		if err != nil {
			t.Fatalf("Repositories.GetCommitRaw returned error: %v", err)
		}
		want := rawStr
		if got != want {
			t.Errorf("Repositories.GetCommitRaw returned %s want %s", got, want)
		}

		const methodName = "CompareCommitsRaw"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.Repositories.CompareCommitsRaw(ctx, "\n", "\n", "\n", "\n", RawOptions{Type: Diff})
			return err
		})

		testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
			got, resp, err := client.Repositories.CompareCommitsRaw(ctx, "o", "r", base, head, RawOptions{Type: Diff})
			if got != "" {
				t.Errorf("testNewRequestAndDoFailure %v = %#v, want ''", methodName, got)
			}
			return resp, err
		})

		teardown()
	}
}

func TestRepositoriesService_CompareCommitsRaw_patch(t *testing.T) {
	testCases := []struct {
		base string
		head string
	}{
		{base: "b", head: "h"},
		{base: "123base", head: "head123"},
		{base: "`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+123base", head: "head123`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+"},
	}

	for _, sample := range testCases {
		client, mux, _, teardown := setup()

		base := sample.base
		head := sample.head
		pattern := fmt.Sprintf("/repos/o/r/compare/%v...%v", base, head)
		const rawStr = "@@patch content"

		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "Accept", mediaTypeV3Patch)
			fmt.Fprint(w, rawStr)
		})

		ctx := context.Background()
		got, _, err := client.Repositories.CompareCommitsRaw(ctx, "o", "r", base, head, RawOptions{Type: Patch})
		if err != nil {
			t.Fatalf("Repositories.GetCommitRaw returned error: %v", err)
		}
		want := rawStr
		if got != want {
			t.Errorf("Repositories.GetCommitRaw returned %s want %s", got, want)
		}

		teardown()
	}
}

func TestRepositoriesService_CompareCommitsRaw_invalid(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		base string
		head string
	}{
		{base: "b", head: "h"},
		{base: "123base", head: "head123"},
		{base: "`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+123base", head: "head123`~!@#$%^&*()_+-=[]\\{}|;':\",./<>?/*-+"},
	}

	for _, sample := range testCases {
		client, _, _, teardown := setup()
		_, _, err := client.Repositories.CompareCommitsRaw(ctx, "o", "r", sample.base, sample.head, RawOptions{100})
		if err == nil {
			t.Fatal("Repositories.GetCommitRaw should return error")
		}
		if !strings.Contains(err.Error(), "unsupported raw type") {
			t.Error("Repositories.GetCommitRaw should return unsupported raw type error")
		}
		teardown()
	}
}

func TestRepositoriesService_ListBranchesHeadCommit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/commits/s/branches-where-head", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `[{"name": "b","commit":{"sha":"2e90302801c870f17b6152327d9b9a03c8eca0e2","url":"https://api.github.com/repos/google/go-github/commits/2e90302801c870f17b6152327d9b9a03c8eca0e2"},"protected":true}]`)
	})

	ctx := context.Background()
	branches, _, err := client.Repositories.ListBranchesHeadCommit(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Repositories.ListBranchesHeadCommit returned error: %v", err)
	}

	want := []*BranchCommit{
		{
			Name: String("b"),
			Commit: &Commit{
				SHA: String("2e90302801c870f17b6152327d9b9a03c8eca0e2"),
				URL: String("https://api.github.com/repos/google/go-github/commits/2e90302801c870f17b6152327d9b9a03c8eca0e2"),
			},
			Protected: Bool(true),
		},
	}
	if !cmp.Equal(branches, want) {
		t.Errorf("Repositories.ListBranchesHeadCommit returned %+v, want %+v", branches, want)
	}

	const methodName = "ListBranchesHeadCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListBranchesHeadCommit(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListBranchesHeadCommit(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
