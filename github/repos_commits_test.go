// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoriesService_ListCommits(t *testing.T) {
	setup()
	defer teardown()

	// given
	mux.HandleFunc("/repos/o/r/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `[
		  {
		    "sha": "s",
		    "commit": { "message": "m" },
		    "author": { "login": "l" },
		    "committer": { "login": "l" },
		    "parents": [ { "sha": "s" } ]
		  }
		]`)
	})

	wantAuthor := &User{
		Login: String("l"),
	}

	want := []RepositoryCommit{
		{
			SHA: String("s"),
			Commit: &Commit{
				Message: String("m"),
			},
			Author:    wantAuthor,
			Committer: wantAuthor,
			Parents: []Commit{
				{
					SHA: String("s"),
				},
			},
		},
	}

	// when
	commits, _, err := client.Repositories.ListCommits("o", "r", nil)

	// then
	if err != nil {
		t.Errorf("Repositories.ListCommits returned error: %v", err)
	}

	if !reflect.DeepEqual(commits, want) {
		t.Errorf("Repositories.ListCommits returned \n%+v, want \n%+v", commits, want)
	}
}

func TestRepositoriesService_GetCommit(t *testing.T) {
	setup()
	defer teardown()

	// given
	mux.HandleFunc("/repos/o/r/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e", func(w http.ResponseWriter, r *http.Request) {
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
		      "raw_url": "r",
		      "blob_url": "b",
		      "patch": "p"
		    }
		  ]
		}`)
	})

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
		Parents: []Commit{
			{
				SHA: String("s"),
			},
		},
		Stats: &CommitStats{
			Additions: Int(104),
			Deletions: Int(4),
			Total:     Int(108),
		},
		Files: []CommitFile{
			{
				Filename:  String("f"),
				Additions: Int(10),
				Deletions: Int(2),
				Changes:   Int(12),
				Status:    String("s"),
				Patch:     String("p"),
			},
		},
	}

	// when
	commit, _, err := client.Repositories.GetCommit("o", "r", "6dcb09b5b57875f334f61aebed695e2e4193db5e")

	// then
	if err != nil {
		t.Errorf("Repositories.GetCommit returned error: %v", err)
	}

	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Repositories.GetCommit returned \n%+v, want \n%+v", commit, want)
	}
}

func TestRepositoriesService_CompareCommits(t *testing.T) {
	setup()
	defer teardown()

	// given
	base := "6dcb09b5b57875f334f61aebed695e2e4193db5e"
	head := "0328041d1152db8ae77652d1618a02e57f745f17"
	url := fmt.Sprintf("/repos/o/r/compare/%v...%v", base, head)

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
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
		    "author": { "login": "n" },
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
		  "files": [ { "filename": "f" } ]
		}`)
	})

	wantUser := &User{
		Login: String("l"),
	}

	wantAuthor := &CommitAuthor{
		Name: String("n"),
	}

	want := &CommitsComparison{
		Status:       String("s"),
		AheadBy:      Int(1),
		BehindBy:     Int(2),
		TotalCommits: Int(1),
		BaseCommit: &RepositoryCommit{
			Commit: &Commit{
				Author: wantAuthor,
			},
			Author:    wantUser,
			Committer: wantUser,
			Message:   String("m"),
		},
		Commits: []RepositoryCommit{
			{
				SHA: String("s"),
			},
		},
		Files: []CommitFile{
			{
				Filename: String("f"),
			},
		},
	}

	// when
	got, _, err := client.Repositories.CompareCommits("o", "r", base, head)

	// then
	if err != nil {
		t.Errorf("Repositories.CompareCommits returned error: %v", err)
	}

	if reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.CompareCommits 	returned \n%+v, want \n%+v", got.Files[0].SHA, want.Files[0].SHA)
	}
}
