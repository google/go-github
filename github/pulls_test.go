// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPullRequestsService_List(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state":     "closed",
			"head":      "h",
			"base":      "b",
			"sort":      "created",
			"direction": "desc",
			"page":      "2",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opts := &PullRequestListOptions{"closed", "h", "b", "created", "desc", ListOptions{Page: 2}}
	ctx := t.Context()
	pulls, _, err := client.PullRequests.List(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("PullRequests.List returned error: %v", err)
	}

	want := []*PullRequest{{Number: new(1)}}
	if !cmp.Equal(pulls, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", pulls, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.List(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.List(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListPullRequestsWithCommit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/commits/sha/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeListPullsOrBranchesForCommitPreview)
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opts := &ListOptions{Page: 2}
	ctx := t.Context()
	pulls, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, "o", "r", "sha", opts)
	if err != nil {
		t.Errorf("PullRequests.ListPullRequestsWithCommit returned error: %v", err)
	}

	want := []*PullRequest{{Number: new(1)}}
	if !cmp.Equal(pulls, want) {
		t.Errorf("PullRequests.ListPullRequestsWithCommit returned %+v, want %+v", pulls, want)
	}

	const methodName = "ListPullRequestsWithCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListPullRequestsWithCommit(ctx, "\n", "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListPullRequestsWithCommit(ctx, "o", "r", "sha", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_List_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.PullRequests.List(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_Get(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := t.Context()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{Number: new(1)}
	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.Get(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.Get(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_GetRaw_diff(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	const rawStr = "@@diff content"

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3Diff)
		fmt.Fprint(w, rawStr)
	})

	ctx := t.Context()
	got, _, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{Diff})
	if err != nil {
		t.Fatalf("PullRequests.GetRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("PullRequests.GetRaw returned %v want %v", got, want)
	}

	const methodName = "GetRaw"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.GetRaw(ctx, "\n", "\n", -1, RawOptions{Diff})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{Diff})
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_GetRaw_patch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	const rawStr = "@@patch content"

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3Patch)
		fmt.Fprint(w, rawStr)
	})

	ctx := t.Context()
	got, _, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{Patch})
	if err != nil {
		t.Fatalf("PullRequests.GetRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("PullRequests.GetRaw returned %v want %v", got, want)
	}
}

func TestPullRequestsService_GetRaw_invalid(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{100})
	if err == nil {
		t.Fatal("PullRequests.GetRaw should return error")
	}
	if !strings.Contains(err.Error(), "unsupported raw type") {
		t.Error("PullRequests.GetRaw should return unsupported raw type error")
	}
}

func TestPullRequestsService_Get_links(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"number":1,
			"_links":{
				"self":{"href":"https://api.github.com/repos/octocat/Hello-World/pulls/1347"},
				"html":{"href":"https://github.com/octocat/Hello-World/pull/1347"},
				"issue":{"href":"https://api.github.com/repos/octocat/Hello-World/issues/1347"},
				"comments":{"href":"https://api.github.com/repos/octocat/Hello-World/issues/1347/comments"},
				"review_comments":{"href":"https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments"},
				"review_comment":{"href":"https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"},
				"commits":{"href":"https://api.github.com/repos/octocat/Hello-World/pulls/1347/commits"},
				"statuses":{"href":"https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e"}
				}
			}`)
	})

	ctx := t.Context()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number: new(1),
		Links: &PRLinks{
			Self: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/pulls/1347"),
			}, HTML: &PRLink{
				HRef: new("https://github.com/octocat/Hello-World/pull/1347"),
			}, Issue: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/issues/1347"),
			}, Comments: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/issues/1347/comments"),
			}, ReviewComments: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments"),
			}, ReviewComment: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"),
			}, Commits: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/pulls/1347/commits"),
			}, Statuses: &PRLink{
				HRef: new("https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e"),
			},
		},
	}
	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Get_headAndBase(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":1,"head":{"ref":"r2","repo":{"id":2}},"base":{"ref":"r1","repo":{"id":1}}}`)
	})

	ctx := t.Context()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number: new(1),
		Head: &PullRequestBranch{
			Ref:  new("r2"),
			Repo: &Repository{ID: new(int64(2))},
		},
		Base: &PullRequestBranch{
			Ref:  new("r1"),
			Repo: &Repository{ID: new(int64(1))},
		},
	}
	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Get_urlFields(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":1,
			"url": "https://api.github.com/repos/octocat/Hello-World/pulls/1347",
			"html_url": "https://github.com/octocat/Hello-World/pull/1347",
			"issue_url": "https://api.github.com/repos/octocat/Hello-World/issues/1347",
			"statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"diff_url": "https://github.com/octocat/Hello-World/pull/1347.diff",
			"patch_url": "https://github.com/octocat/Hello-World/pull/1347.patch",
			"review_comments_url": "https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments",
			"review_comment_url": "https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"}`)
	})

	ctx := t.Context()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number:            new(1),
		URL:               new("https://api.github.com/repos/octocat/Hello-World/pulls/1347"),
		HTMLURL:           new("https://github.com/octocat/Hello-World/pull/1347"),
		IssueURL:          new("https://api.github.com/repos/octocat/Hello-World/issues/1347"),
		StatusesURL:       new("https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e"),
		DiffURL:           new("https://github.com/octocat/Hello-World/pull/1347.diff"),
		PatchURL:          new("https://github.com/octocat/Hello-World/pull/1347.patch"),
		ReviewCommentsURL: new("https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments"),
		ReviewCommentURL:  new("https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"),
	}

	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Get_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.PullRequests.Get(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := CreatePullRequest{Title: new("t")}

	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)

		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := t.Context()
	pull, _, err := client.PullRequests.Create(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("PullRequests.Create returned error: %v", err)
	}

	want := &PullRequest{Number: new(1)}
	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Create returned %+v, want %+v", pull, want)
	}

	const methodName = "Create"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.Create(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.Create(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_Create_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.PullRequests.Create(ctx, "%", "r", CreatePullRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_UpdateBranch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/update-branch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeUpdatePullRequestBranchPreview)
		fmt.Fprint(w, `
			{
			  "message": "Updating pull request branch.",
			  "url": "https://github.com/repos/o/r/pulls/1"
			}`)
	})

	opts := &PullRequestBranchUpdateOptions{
		ExpectedHeadSHA: new("s"),
	}

	ctx := t.Context()
	pull, _, err := client.PullRequests.UpdateBranch(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.UpdateBranch returned error: %v", err)
	}

	want := &PullRequestBranchUpdateResponse{
		Message: new("Updating pull request branch."),
		URL:     new("https://github.com/repos/o/r/pulls/1"),
	}

	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.UpdateBranch returned %+v, want %+v", pull, want)
	}

	const methodName = "UpdateBranch"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.UpdateBranch(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.UpdateBranch(ctx, "o", "r", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_Edit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	tests := []struct {
		input        *PullRequest
		sendResponse string
		want         *PullRequest
		wantUpdate   *pullRequestUpdate
	}{
		{
			input:        &PullRequest{Title: new("t")},
			sendResponse: `{"number":1}`,
			want:         &PullRequest{Number: new(1)},
			wantUpdate: &pullRequestUpdate{
				Title: new("t"),
			},
		},
		{
			// base update
			input:        &PullRequest{Base: &PullRequestBranch{Ref: new("master")}},
			sendResponse: `{"number":1,"base":{"ref":"master"}}`,
			want: &PullRequest{
				Number: new(1),
				Base:   &PullRequestBranch{Ref: new("master")},
			},
			wantUpdate: &pullRequestUpdate{
				Base: new("master"),
			},
		},
	}

	for i, tt := range tests {
		madeRequest := false
		mux.HandleFunc(fmt.Sprintf("/repos/o/r/pulls/%v", i), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			testJSONBody(t, r, tt.wantUpdate)
			_, err := io.WriteString(w, tt.sendResponse)
			assertNilError(t, err)
			madeRequest = true
		})

		ctx := t.Context()
		pull, _, err := client.PullRequests.Edit(ctx, "o", "r", i, tt.input)
		if err != nil {
			t.Errorf("%v: PullRequests.Edit returned error: %v", i, err)
		}

		if !cmp.Equal(pull, tt.want) {
			t.Errorf("%v: PullRequests.Edit returned %+v, want %+v", i, pull, tt.want)
		}

		if !madeRequest {
			t.Errorf("%v: PullRequest.Edit did not make the expected request", i)
		}

		const methodName = "Edit"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.PullRequests.Edit(ctx, "\n", "\n", -i, tt.input)
			return err
		})
	}
	testNewRequestAndDoFailure(t, "Edit", client, func() (*Response, error) {
		got, resp, err := client.PullRequests.Edit(t.Context(), "o", "r", 1, &PullRequest{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", "Edit", got)
		}
		return resp, err
	})
}

func TestPullRequestsService_Edit_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.PullRequests.Edit(ctx, "%", "r", 1, &PullRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_ListCommits(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `
			[
			  {
			    "sha": "3",
			    "parents": [
			      {
			        "sha": "2"
			      }
			    ]
			  },
			  {
			    "sha": "2",
			    "parents": [
			      {
			        "sha": "1"
			      }
			    ]
			  }
			]`)
	})

	opts := &ListOptions{Page: 2}
	ctx := t.Context()
	commits, _, err := client.PullRequests.ListCommits(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.ListCommits returned error: %v", err)
	}

	want := []*RepositoryCommit{
		{
			SHA: new("3"),
			Parents: []*Commit{
				{
					SHA: new("2"),
				},
			},
		},
		{
			SHA: new("2"),
			Parents: []*Commit{
				{
					SHA: new("1"),
				},
			},
		},
	}
	if !cmp.Equal(commits, want) {
		t.Errorf("PullRequests.ListCommits returned %+v, want %+v", commits, want)
	}

	const methodName = "ListCommits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListCommits(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListCommits(ctx, "o", "r", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_ListFiles(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/files", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `
			[
			  {
			    "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			    "filename": "file1.txt",
			    "status": "added",
			    "additions": 103,
			    "deletions": 21,
			    "changes": 124,
			    "patch": "@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"
			  },
			  {
			    "sha": "f61aebed695e2e4193db5e6dcb09b5b57875f334",
			    "filename": "file2.txt",
			    "status": "modified",
			    "additions": 5,
			    "deletions": 3,
			    "changes": 103,
			    "patch": "@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"
			  }
			]`)
	})

	opts := &ListOptions{Page: 2}
	ctx := t.Context()
	commitFiles, _, err := client.PullRequests.ListFiles(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.ListFiles returned error: %v", err)
	}

	want := []*CommitFile{
		{
			SHA:       new("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
			Filename:  new("file1.txt"),
			Additions: new(103),
			Deletions: new(21),
			Changes:   new(124),
			Status:    new("added"),
			Patch:     new("@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"),
		},
		{
			SHA:       new("f61aebed695e2e4193db5e6dcb09b5b57875f334"),
			Filename:  new("file2.txt"),
			Additions: new(5),
			Deletions: new(3),
			Changes:   new(103),
			Status:    new("modified"),
			Patch:     new("@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"),
		},
	}

	if !cmp.Equal(commitFiles, want) {
		t.Errorf("PullRequests.ListFiles returned %+v, want %+v", commitFiles, want)
	}

	const methodName = "ListFiles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListFiles(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListFiles(ctx, "o", "r", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_IsMerged(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	isMerged, _, err := client.PullRequests.IsMerged(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.IsMerged returned error: %v", err)
	}

	want := true
	if !cmp.Equal(isMerged, want) {
		t.Errorf("PullRequests.IsMerged returned %+v, want %+v", isMerged, want)
	}

	const methodName = "IsMerged"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.IsMerged(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.IsMerged(ctx, "o", "r", 1)
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestPullRequestsService_Merge(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `
			{
			  "sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			  "merged": true,
			  "message": "Pull Request successfully merged"
			}`)
	})

	options := &PullRequestOptions{MergeMethod: "rebase"}
	ctx := t.Context()
	merge, _, err := client.PullRequests.Merge(ctx, "o", "r", 1, "merging pull request", options)
	if err != nil {
		t.Errorf("PullRequests.Merge returned error: %v", err)
	}

	want := &PullRequestMergeResult{
		SHA:     new("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
		Merged:  new(true),
		Message: new("Pull Request successfully merged"),
	}
	if !cmp.Equal(merge, want) {
		t.Errorf("PullRequests.Merge returned %+v, want %+v", merge, want)
	}

	const methodName = "Merge"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.Merge(ctx, "\n", "\n", -1, "\n", options)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.Merge(ctx, "o", "r", 1, "merging pull request", options)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

// Test that different merge options produce expected PUT requests. See issue https://github.com/google/go-github/issues/500.
func TestPullRequestsService_Merge_options(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	tests := []struct {
		options *PullRequestOptions
		want    pullRequestMergeRequest
	}{
		{
			options: nil,
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
			},
		},
		{
			options: &PullRequestOptions{},
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
			},
		},
		{
			options: &PullRequestOptions{MergeMethod: "rebase"},
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
				MergeMethod:   "rebase",
			},
		},
		{
			options: &PullRequestOptions{SHA: "6dcb09b5b57875f334f61aebed695e2e4193db5e"},
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
				SHA:           "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			},
		},
		{
			options: &PullRequestOptions{
				CommitTitle: "Extra detail",
				SHA:         "6dcb09b5b57875f334f61aebed695e2e4193db5e",
				MergeMethod: "squash",
			},
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
				SHA:           "6dcb09b5b57875f334f61aebed695e2e4193db5e",
				CommitTitle:   "Extra detail",
				MergeMethod:   "squash",
			},
		},
		{
			options: &PullRequestOptions{
				DontDefaultIfBlank: true,
			},
			want: pullRequestMergeRequest{
				CommitMessage: new("merging pull request"),
			},
		},
	}

	for i, test := range tests {
		madeRequest := false
		mux.HandleFunc(fmt.Sprintf("/repos/o/r/pulls/%v/merge", i), func(_ http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testJSONBody(t, r, test.want)
			madeRequest = true
		})
		ctx := t.Context()
		_, _, _ = client.PullRequests.Merge(ctx, "o", "r", i, "merging pull request", test.options)
		if !madeRequest {
			t.Errorf("%v: PullRequests.Merge(%#v): expected request was not made", i, test.options)
		}
	}
}

func TestPullRequestsService_Merge_Blank_Message(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	madeRequest := false
	want := &pullRequestMergeRequest{}
	mux.HandleFunc("/repos/o/r/pulls/1/merge", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, want)
		madeRequest = true
	})

	ctx := t.Context()
	_, _, _ = client.PullRequests.Merge(ctx, "o", "r", 1, "", nil)
	if !madeRequest {
		t.Error("TestPullRequestsService_Merge_Blank_Message #1 did not make request")
	}

	madeRequest = false
	opts := PullRequestOptions{
		DontDefaultIfBlank: true,
	}
	want = &pullRequestMergeRequest{
		CommitMessage: new(""),
	}
	_, _, _ = client.PullRequests.Merge(ctx, "o", "r", 1, "", &opts)
	if !madeRequest {
		t.Error("TestPullRequestsService_Merge_Blank_Message #2 did not make request")
	}
}
