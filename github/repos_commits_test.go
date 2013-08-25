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
	"time"
)

const timeFormat = "2011-01-14 16:00:49 +0000 UTC"

func TestRepositoriesService_ListCommits(t *testing.T) {
	setup()
	defer teardown()

	// given
	mux.HandleFunc("/repos/o/r/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `[
		  {
		    "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		    "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		    "commit": {
		      "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "author": {
		        "name": "Monalisa Octocat",
		        "email": "support@github.com",
		        "date": "2011-01-14T16:00:49Z"
		      },
		      "committer": {
		        "name": "Monalisa Octocat",
		        "email": "support@github.com",
		        "date": "2011-01-14T16:00:49Z"
		      },
		      "message": "Fix all the bugs",
		      "tree": {
		        "url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		        "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		      }
		    },
		    "author": {
		      "login": "octocat",
		      "id": 1,
		      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		      "gravatar_id": "somehexcode",
		      "url": "https://api.github.com/users/octocat"
		    },
		    "committer": {
		      "login": "octocat",
		      "id": 1,
		      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		      "gravatar_id": "somehexcode",
		      "url": "https://api.github.com/users/octocat"
		    },
		    "parents": [
		      {
		        "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		        "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		      }
		    ]
		  }
		]`)
	})

	wantCommitTime := time.Date(2011, time.January, 14, 16, 0, 49, 0, time.UTC)

	wantCommitAuthor := &CommitAuthor{
		Name:  String("Monalisa Octocat"),
		Email: String("support@github.com"),
		Date:  &wantCommitTime,
	}

	wantRepoCommitAuthor := &RepositoryCommitAuthor{
		Login:      String("octocat"),
		Id:         Int(1),
		AvatarUrl:  String("https://github.com/images/error/octocat_happy.gif"),
		GravatarId: String("somehexcode"),
	}

	want := []RepositoryCommit{
		{
			SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
			Commit: &Commit{
				Author:    wantCommitAuthor,
				Committer: wantCommitAuthor,
				Message:   String("Fix all the bugs"),
				Tree: &Tree{
					SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
				},
			},
			Author:    wantRepoCommitAuthor,
			Committer: wantRepoCommitAuthor,
			Parents: []Commit{
				{
					SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
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
		  "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		  "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		  "commit": {
		    "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		    "author": {
		      "name": "Monalisa Octocat",
		      "email": "support@github.com",
		      "date": "2011-01-14T16:00:49Z"
		    },
		    "committer": {
		      "name": "Monalisa Octocat",
		      "email": "support@github.com",
		      "date": "2011-01-14T16:00:49Z"
		    },
		    "message": "Fix all the bugs",
		    "tree": {
		      "url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		    }
		  },
		  "author": {
		    "login": "octocat",
		    "id": 1,
		    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		    "gravatar_id": "somehexcode",
		    "url": "https://api.github.com/users/octocat"
		  },
		  "committer": {
		    "login": "octocat",
		    "id": 1,
		    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		    "gravatar_id": "somehexcode",
		    "url": "https://api.github.com/users/octocat"
		  },
		  "parents": [
		    {
		      "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		    }
		  ],
		  "stats": {
		    "additions": 104,
		    "deletions": 4,
		    "total": 108
		  },
		  "files": [
		    {
		      "filename": "file1.txt",
		      "additions": 10,
		      "deletions": 2,
		      "changes": 12,
		      "status": "modified",
		      "raw_url": "https://github.com/octocat/Hello-World/raw/7ca483543807a51b6079e54ac4cc392bc29ae284/file1.txt",
		      "blob_url": "https://github.com/octocat/Hello-World/blob/7ca483543807a51b6079e54ac4cc392bc29ae284/file1.txt",
		      "patch": "@@ -29,7 +29,7 @@\n....."
		    }
		  ]
		}`)
	})

	wantCommitTime := time.Date(2011, time.January, 14, 16, 0, 49, 0, time.UTC)

	wantCommitAuthor := &CommitAuthor{
		Name:  String("Monalisa Octocat"),
		Email: String("support@github.com"),
		Date:  &wantCommitTime,
	}

	wantTree := &Tree{
		SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
	}

	wantRepoCommitAuthor := &RepositoryCommitAuthor{
		Login:      String("octocat"),
		Id:         Int(1),
		AvatarUrl:  String("https://github.com/images/error/octocat_happy.gif"),
		GravatarId: String("somehexcode"),
	}

	want := &RepositoryCommit{
		SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
		Commit: &Commit{
			Author:    wantCommitAuthor,
			Committer: wantCommitAuthor,
			Message:   String("Fix all the bugs"),
			Tree:      wantTree,
		},
		Author:    wantRepoCommitAuthor,
		Committer: wantRepoCommitAuthor,
		Parents: []Commit{
			{
				SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
			},
		},
		Stats: &CommitStats{
			Additions: Int(104),
			Deletions: Int(4),
			Total:     Int(108),
		},
		Files: []CommitFile{
			{
				Filename:  String("file1.txt"),
				Additions: Int(10),
				Deletions: Int(2),
				Changes:   Int(12),
				Status:    String("modified"),
				Patch:     String("@@ -29,7 +29,7 @@\n....."),
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
		t.Errorf("Repositories.GetCommit returned \n%+v, want \n%+v", commit.Commit.SHA, want.Commit.SHA)
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
		  "url": "https://api.github.com/repos/octocat/Hello-World/compare/master...topic",
		  "html_url": "https://github.com/octocat/Hello-World/compare/master...topic",
		  "permalink_url": "https://github.com/octocat/Hello-World/compare/octocat:bbcd538c8e72b8c175046e27cc8f907076331401...octocat:0328041d1152db8ae77652d1618a02e57f745f17",
		  "diff_url": "https://github.com/octocat/Hello-World/compare/master...topic.diff",
		  "patch_url": "https://github.com/octocat/Hello-World/compare/master...topic.patch",

		  "base_commit": {
		    "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		    "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		    "commit": {
		      "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "author": {
		        "name": "Monalisa Octocat",
		        "email": "support@github.com",
		        "date": "2011-01-14T16:00:49Z"
		      },
		      "committer": {
		        "name": "Monalisa Octocat",
		        "email": "support@github.com",
		        "date": "2011-01-14T16:00:49Z"
		      },
		      "message": "Fix all the bugs",
		      "tree": {
		        "url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		        "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		      }
		    },
		    "author": {
		      "login": "octocat",
		      "id": 1,
		      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		      "gravatar_id": "someauthorhexcode",
		      "url": "https://api.github.com/users/octocat"
		    },
		    "committer": {
		      "login": "octocat",
		      "id": 1,
		      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		      "gravatar_id": "somecommitterhexcode",
		      "url": "https://api.github.com/users/octocat"
		    },
		    "parents": [
		      {
		        "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		        "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		      }
		    ]
		  },

		  "status": "behind",
		  "ahead_by": 1,
		  "behind_by": 2,
		  "total_commits": 1,

		  "commits": [
		    {
		      "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "commit": {
		        "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		        "author": {
		          "name": "Monalisa Octocat",
		          "email": "support@github.com",
		          "date": "2011-01-14T16:00:49Z"
		        },
		        "committer": {
		          "name": "Monalisa Octocat",
		          "email": "support@github.com",
		          "date": "2011-01-14T16:00:49Z"
		        },
		        "message": "Fix all the bugs",
		        "tree": {
		          "url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		          "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		        }
		      },
		      "author": {
		        "login": "octocat",
		        "id": 1,
		        "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		        "gravatar_id": "somehexcode",
		        "url": "https://api.github.com/users/octocat"
		      },
		      "committer": {
		        "login": "octocat",
		        "id": 1,
		        "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		        "gravatar_id": "somehexcode",
		        "url": "https://api.github.com/users/octocat"
		      },
		      "parents": [
		        {
		          "url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		          "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		        }
		      ]
		    }
		  ],
		  "files": [
		    {
		      "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		      "filename": "example.txt",
		      "status": "added",
		      "additions": 103,
		      "deletions": 21,
		      "changes": 124,
		      "blob_url": "https://github.com/octocat/Hello-World/blob/6dcb09b5b57875f334f61aebed695e2e4193db5e/file1.txt",
		      "raw_url": "https://github.com/octocat/Hello-World/raw/6dcb09b5b57875f334f61aebed695e2e4193db5e/file1.txt",
		      "patch": "@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"
		    }
		  ]
		}`)
	})

	wantCommitTime := time.Date(2011, time.January, 14, 16, 0, 49, 0, time.UTC)

	wantCommitAuthor := &CommitAuthor{
		Name:  String("Monalisa Octocat"),
		Email: String("support@github.com"),
		Date:  &wantCommitTime,
	}

	wantTree := &Tree{
		SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
	}

	wantRepoCommitAuthor := &RepositoryCommitAuthor{
		Login:      String("octocat"),
		Id:         Int(1),
		AvatarUrl:  String("https://github.com/images/error/octocat_happy.gif"),
		GravatarId: String("somehexcode"),
	}

	author := &CommitAuthor{
		Date:  &wantCommitTime,
		Name:  String("Monalisa Octocat"),
		Email: String("support@github.com"),
	}

	sha := String("6dcb09b5b57875f334f61aebed695e2e4193db5e")

	want := &CommitsComparison{
		Status:       String("behind"),
		AheadBy:      Int(1),
		BehindBy:     Int(2),
		TotalCommits: Int(1),
		BaseCommit: &BaseComparedCommit{
			Commit: &Commit{
				SHA:       sha,
				Author:    author,
				Committer: author,
			},
			Author:    wantRepoCommitAuthor,
			Committer: wantRepoCommitAuthor,
			Message:   String("Fix all the bugs"),
		},
		Commits: []RepositoryCommit{
			{
				SHA: sha,
				Commit: &Commit{
					SHA:       String(""),
					Author:    wantCommitAuthor,
					Committer: wantCommitAuthor,
					Message:   String("Fix all the bugs"),
					Tree:      wantTree,
					Parents:   nil,
				},
				Author:    wantRepoCommitAuthor,
				Committer: wantRepoCommitAuthor,
				Parents: []Commit{
					{
						SHA: String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
					},
				},
			},
		},
		Files: []CommitFile{
			{
				Filename:  String("file1.txt"),
				Additions: Int(103),
				Deletions: Int(21),
				Changes:   Int(124),
				Status:    String("added"),
				Patch:     String("@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"),
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
