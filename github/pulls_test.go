// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
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
	ctx := context.Background()
	pulls, _, err := client.PullRequests.List(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("PullRequests.List returned error: %v", err)
	}

	want := []*PullRequest{{Number: Ptr(1)}}
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
	ctx := context.Background()
	pulls, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, "o", "r", "sha", opts)
	if err != nil {
		t.Errorf("PullRequests.ListPullRequestsWithCommit returned error: %v", err)
	}

	want := []*PullRequest{{Number: Ptr(1)}}
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

	ctx := context.Background()
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

	ctx := context.Background()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{Number: Ptr(1)}
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

	ctx := context.Background()
	got, _, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{Diff})
	if err != nil {
		t.Fatalf("PullRequests.GetRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("PullRequests.GetRaw returned %s want %s", got, want)
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

	ctx := context.Background()
	got, _, err := client.PullRequests.GetRaw(ctx, "o", "r", 1, RawOptions{Patch})
	if err != nil {
		t.Fatalf("PullRequests.GetRaw returned error: %v", err)
	}
	want := rawStr
	if got != want {
		t.Errorf("PullRequests.GetRaw returned %s want %s", got, want)
	}
}

func TestPullRequestsService_GetRaw_invalid(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
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

	ctx := context.Background()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number: Ptr(1),
		Links: &PRLinks{
			Self: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/1347"),
			}, HTML: &PRLink{
				HRef: Ptr("https://github.com/octocat/Hello-World/pull/1347"),
			}, Issue: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/issues/1347"),
			}, Comments: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/issues/1347/comments"),
			}, ReviewComments: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments"),
			}, ReviewComment: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"),
			}, Commits: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/1347/commits"),
			}, Statuses: &PRLink{
				HRef: Ptr("https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e"),
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

	ctx := context.Background()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number: Ptr(1),
		Head: &PullRequestBranch{
			Ref:  Ptr("r2"),
			Repo: &Repository{ID: Ptr(int64(2))},
		},
		Base: &PullRequestBranch{
			Ref:  Ptr("r1"),
			Repo: &Repository{ID: Ptr(int64(1))},
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

	ctx := context.Background()
	pull, _, err := client.PullRequests.Get(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &PullRequest{
		Number:            Ptr(1),
		URL:               Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/1347"),
		HTMLURL:           Ptr("https://github.com/octocat/Hello-World/pull/1347"),
		IssueURL:          Ptr("https://api.github.com/repos/octocat/Hello-World/issues/1347"),
		StatusesURL:       Ptr("https://api.github.com/repos/octocat/Hello-World/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e"),
		DiffURL:           Ptr("https://github.com/octocat/Hello-World/pull/1347.diff"),
		PatchURL:          Ptr("https://github.com/octocat/Hello-World/pull/1347.patch"),
		ReviewCommentsURL: Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/1347/comments"),
		ReviewCommentURL:  Ptr("https://api.github.com/repos/octocat/Hello-World/pulls/comments{/number}"),
	}

	if !cmp.Equal(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_Get_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.PullRequests.Get(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &NewPullRequest{Title: Ptr("t")}

	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewPullRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	ctx := context.Background()
	pull, _, err := client.PullRequests.Create(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("PullRequests.Create returned error: %v", err)
	}

	want := &PullRequest{Number: Ptr(1)}
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

	ctx := context.Background()
	_, _, err := client.PullRequests.Create(ctx, "%", "r", nil)
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
		ExpectedHeadSHA: Ptr("s"),
	}

	ctx := context.Background()
	pull, _, err := client.PullRequests.UpdateBranch(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.UpdateBranch returned error: %v", err)
	}

	want := &PullRequestBranchUpdateResponse{
		Message: Ptr("Updating pull request branch."),
		URL:     Ptr("https://github.com/repos/o/r/pulls/1"),
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

		wantUpdate string
		want       *PullRequest
	}{
		{
			input:        &PullRequest{Title: Ptr("t")},
			sendResponse: `{"number":1}`,
			wantUpdate:   `{"title":"t"}`,
			want:         &PullRequest{Number: Ptr(1)},
		},
		{
			// base update
			input:        &PullRequest{Base: &PullRequestBranch{Ref: Ptr("master")}},
			sendResponse: `{"number":1,"base":{"ref":"master"}}`,
			wantUpdate:   `{"base":"master"}`,
			want: &PullRequest{
				Number: Ptr(1),
				Base:   &PullRequestBranch{Ref: Ptr("master")},
			},
		},
	}

	for i, tt := range tests {
		madeRequest := false
		mux.HandleFunc(fmt.Sprintf("/repos/o/r/pulls/%v", i), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			testBody(t, r, tt.wantUpdate+"\n")
			_, err := io.WriteString(w, tt.sendResponse)
			assertNilError(t, err)
			madeRequest = true
		})

		ctx := context.Background()
		pull, _, err := client.PullRequests.Edit(ctx, "o", "r", i, tt.input)
		if err != nil {
			t.Errorf("%d: PullRequests.Edit returned error: %v", i, err)
		}

		if !cmp.Equal(pull, tt.want) {
			t.Errorf("%d: PullRequests.Edit returned %+v, want %+v", i, pull, tt.want)
		}

		if !madeRequest {
			t.Errorf("%d: PullRequest.Edit did not make the expected request", i)
		}

		const methodName = "Edit"
		testBadOptions(t, methodName, func() (err error) {
			_, _, err = client.PullRequests.Edit(ctx, "\n", "\n", -i, tt.input)
			return err
		})
	}
}

func TestPullRequestsService_Edit_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
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
	ctx := context.Background()
	commits, _, err := client.PullRequests.ListCommits(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.ListCommits returned error: %v", err)
	}

	want := []*RepositoryCommit{
		{
			SHA: Ptr("3"),
			Parents: []*Commit{
				{
					SHA: Ptr("2"),
				},
			},
		},
		{
			SHA: Ptr("2"),
			Parents: []*Commit{
				{
					SHA: Ptr("1"),
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
	ctx := context.Background()
	commitFiles, _, err := client.PullRequests.ListFiles(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("PullRequests.ListFiles returned error: %v", err)
	}

	want := []*CommitFile{
		{
			SHA:       Ptr("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
			Filename:  Ptr("file1.txt"),
			Additions: Ptr(103),
			Deletions: Ptr(21),
			Changes:   Ptr(124),
			Status:    Ptr("added"),
			Patch:     Ptr("@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"),
		},
		{
			SHA:       Ptr("f61aebed695e2e4193db5e6dcb09b5b57875f334"),
			Filename:  Ptr("file2.txt"),
			Additions: Ptr(5),
			Deletions: Ptr(3),
			Changes:   Ptr(103),
			Status:    Ptr("modified"),
			Patch:     Ptr("@@ -132,7 +132,7 @@ module Test @@ -1000,7 +1000,7 @@ module Test"),
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

	ctx := context.Background()
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
	ctx := context.Background()
	merge, _, err := client.PullRequests.Merge(ctx, "o", "r", 1, "merging pull request", options)
	if err != nil {
		t.Errorf("PullRequests.Merge returned error: %v", err)
	}

	want := &PullRequestMergeResult{
		SHA:     Ptr("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
		Merged:  Ptr(true),
		Message: Ptr("Pull Request successfully merged"),
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
		options  *PullRequestOptions
		wantBody string
	}{
		{
			options:  nil,
			wantBody: `{"commit_message":"merging pull request"}`,
		},
		{
			options:  &PullRequestOptions{},
			wantBody: `{"commit_message":"merging pull request"}`,
		},
		{
			options:  &PullRequestOptions{MergeMethod: "rebase"},
			wantBody: `{"commit_message":"merging pull request","merge_method":"rebase"}`,
		},
		{
			options:  &PullRequestOptions{SHA: "6dcb09b5b57875f334f61aebed695e2e4193db5e"},
			wantBody: `{"commit_message":"merging pull request","sha":"6dcb09b5b57875f334f61aebed695e2e4193db5e"}`,
		},
		{
			options: &PullRequestOptions{
				CommitTitle: "Extra detail",
				SHA:         "6dcb09b5b57875f334f61aebed695e2e4193db5e",
				MergeMethod: "squash",
			},
			wantBody: `{"commit_message":"merging pull request","commit_title":"Extra detail","merge_method":"squash","sha":"6dcb09b5b57875f334f61aebed695e2e4193db5e"}`,
		},
		{
			options: &PullRequestOptions{
				DontDefaultIfBlank: true,
			},
			wantBody: `{"commit_message":"merging pull request"}`,
		},
	}

	for i, test := range tests {
		madeRequest := false
		mux.HandleFunc(fmt.Sprintf("/repos/o/r/pulls/%d/merge", i), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testBody(t, r, test.wantBody+"\n")
			madeRequest = true
		})
		ctx := context.Background()
		_, _, _ = client.PullRequests.Merge(ctx, "o", "r", i, "merging pull request", test.options)
		if !madeRequest {
			t.Errorf("%d: PullRequests.Merge(%#v): expected request was not made", i, test.options)
		}
	}
}

func TestPullRequestsService_Merge_Blank_Message(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	madeRequest := false
	expectedBody := ""
	mux.HandleFunc("/repos/o/r/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, expectedBody+"\n")
		madeRequest = true
	})

	ctx := context.Background()
	expectedBody = `{}`
	_, _, _ = client.PullRequests.Merge(ctx, "o", "r", 1, "", nil)
	if !madeRequest {
		t.Error("TestPullRequestsService_Merge_Blank_Message #1 did not make request")
	}

	madeRequest = false
	opts := PullRequestOptions{
		DontDefaultIfBlank: true,
	}
	expectedBody = `{"commit_message":""}`
	_, _, _ = client.PullRequests.Merge(ctx, "o", "r", 1, "", &opts)
	if !madeRequest {
		t.Error("TestPullRequestsService_Merge_Blank_Message #2 did not make request")
	}
}

func TestPullRequestMergeRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &pullRequestMergeRequest{}, "{}")

	u := &pullRequestMergeRequest{
		CommitMessage: Ptr("cm"),
		CommitTitle:   "ct",
		MergeMethod:   "mm",
		SHA:           "sha",
	}

	want := `{
		"commit_message": "cm",
		"commit_title": "ct",
		"merge_method": "mm",
		"sha": "sha"
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestMergeResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestMergeResult{}, "{}")

	u := &PullRequestMergeResult{
		SHA:     Ptr("sha"),
		Merged:  Ptr(false),
		Message: Ptr("msg"),
	}

	want := `{
		"sha": "sha",
		"merged": false,
		"message": "msg"
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestUpdate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &pullRequestUpdate{}, "{}")

	u := &pullRequestUpdate{
		Title:               Ptr("title"),
		Body:                Ptr("body"),
		State:               Ptr("state"),
		Base:                Ptr("base"),
		MaintainerCanModify: Ptr(false),
	}

	want := `{
		"title": "title",
		"body": "body",
		"state": "state",
		"base": "base",
		"maintainer_can_modify": false
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestBranchUpdateResponse_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestBranchUpdateResponse{}, "{}")

	u := &PullRequestBranchUpdateResponse{
		Message: Ptr("message"),
		URL:     Ptr("url"),
	}

	want := `{
		"message": "message",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestBranchUpdateOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestBranchUpdateOptions{}, "{}")

	u := &PullRequestBranchUpdateOptions{
		ExpectedHeadSHA: Ptr("eh"),
	}

	want := `{
		"expected_head_sha": "eh"
	}`

	testJSONMarshal(t, u, want)
}

func TestNewPullRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &NewPullRequest{}, "{}")

	u := &NewPullRequest{
		Title:               Ptr("eh"),
		Head:                Ptr("eh"),
		HeadRepo:            Ptr("eh"),
		Base:                Ptr("eh"),
		Body:                Ptr("eh"),
		Issue:               Ptr(1),
		MaintainerCanModify: Ptr(false),
		Draft:               Ptr(false),
	}

	want := `{
		"title": "eh",
		"head": "eh",
		"head_repo": "eh",
		"base": "eh",
		"body": "eh",
		"issue": 1,
		"maintainer_can_modify": false,
		"draft": false
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestBranch_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestBranch{}, "{}")

	u := &PullRequestBranch{
		Label: Ptr("label"),
		Ref:   Ptr("ref"),
		SHA:   Ptr("sha"),
		Repo:  &Repository{ID: Ptr(int64(1))},
		User: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
	}

	want := `{
		"label": "label",
		"ref": "ref",
		"sha": "sha",
		"repo": {
			"id": 1
		},
		"user": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPRLink_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PRLink{}, "{}")

	u := &PRLink{
		HRef: Ptr("href"),
	}

	want := `{
		"href": "href"
	}`

	testJSONMarshal(t, u, want)
}

func TestPRLinks_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PRLinks{}, "{}")

	u := &PRLinks{
		Self: &PRLink{
			HRef: Ptr("href"),
		},
		HTML: &PRLink{
			HRef: Ptr("href"),
		},
		Issue: &PRLink{
			HRef: Ptr("href"),
		},
		Comments: &PRLink{
			HRef: Ptr("href"),
		},
		ReviewComments: &PRLink{
			HRef: Ptr("href"),
		},
		ReviewComment: &PRLink{
			HRef: Ptr("href"),
		},
		Commits: &PRLink{
			HRef: Ptr("href"),
		},
		Statuses: &PRLink{
			HRef: Ptr("href"),
		},
	}

	want := `{
		"self": {
			"href": "href"
		},
		"html": {
			"href": "href"
		},
		"issue": {
			"href": "href"
		},
		"comments": {
			"href": "href"
		},
		"review_comments": {
			"href": "href"
		},
		"review_comment": {
			"href": "href"
		},
		"commits": {
			"href": "href"
		},
		"statuses": {
			"href": "href"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestAutoMerge_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequestAutoMerge{}, "{}")

	u := &PullRequestAutoMerge{
		EnabledBy: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		MergeMethod:   Ptr("mm"),
		CommitTitle:   Ptr("ct"),
		CommitMessage: Ptr("cm"),
	}

	want := `{
		"enabled_by": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"merge_method": "mm",
		"commit_title": "ct",
		"commit_message": "cm"
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PullRequest{}, "{}")

	u := &PullRequest{
		ID:        Ptr(int64(1)),
		Number:    Ptr(1),
		State:     Ptr("state"),
		Locked:    Ptr(false),
		Title:     Ptr("title"),
		Body:      Ptr("body"),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		ClosedAt:  &Timestamp{referenceTime},
		MergedAt:  &Timestamp{referenceTime},
		Labels:    []*Label{{ID: Ptr(int64(1))}},
		User: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Draft:          Ptr(false),
		Merged:         Ptr(false),
		Mergeable:      Ptr(false),
		MergeableState: Ptr("ms"),
		MergedBy: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		MergeCommitSHA:    Ptr("mcs"),
		Rebaseable:        Ptr(false),
		Comments:          Ptr(1),
		Commits:           Ptr(1),
		Additions:         Ptr(1),
		Deletions:         Ptr(1),
		ChangedFiles:      Ptr(1),
		URL:               Ptr("url"),
		HTMLURL:           Ptr("hurl"),
		IssueURL:          Ptr("iurl"),
		StatusesURL:       Ptr("surl"),
		DiffURL:           Ptr("durl"),
		PatchURL:          Ptr("purl"),
		CommitsURL:        Ptr("curl"),
		CommentsURL:       Ptr("comurl"),
		ReviewCommentsURL: Ptr("rcurls"),
		ReviewCommentURL:  Ptr("rcurl"),
		ReviewComments:    Ptr(1),
		Assignee: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Assignees: []*User{
			{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
		},
		Milestone:           &Milestone{ID: Ptr(int64(1))},
		MaintainerCanModify: Ptr(true),
		AuthorAssociation:   Ptr("aa"),
		NodeID:              Ptr("nid"),
		RequestedReviewers: []*User{
			{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
		},
		AutoMerge: &PullRequestAutoMerge{
			EnabledBy: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			MergeMethod:   Ptr("mm"),
			CommitTitle:   Ptr("ct"),
			CommitMessage: Ptr("cm"),
		},
		RequestedTeams: []*Team{{ID: Ptr(int64(1))}},
		Links: &PRLinks{
			Self: &PRLink{
				HRef: Ptr("href"),
			},
			HTML: &PRLink{
				HRef: Ptr("href"),
			},
			Issue: &PRLink{
				HRef: Ptr("href"),
			},
			Comments: &PRLink{
				HRef: Ptr("href"),
			},
			ReviewComments: &PRLink{
				HRef: Ptr("href"),
			},
			ReviewComment: &PRLink{
				HRef: Ptr("href"),
			},
			Commits: &PRLink{
				HRef: Ptr("href"),
			},
			Statuses: &PRLink{
				HRef: Ptr("href"),
			},
		},
		Head: &PullRequestBranch{
			Ref:  Ptr("r2"),
			Repo: &Repository{ID: Ptr(int64(2))},
		},
		Base: &PullRequestBranch{
			Ref:  Ptr("r2"),
			Repo: &Repository{ID: Ptr(int64(2))},
		},
		ActiveLockReason: Ptr("alr"),
	}

	want := `{
		"id": 1,
		"number": 1,
		"state": "state",
		"locked": false,
		"title": "title",
		"body": "body",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"closed_at": ` + referenceTimeStr + `,
		"merged_at": ` + referenceTimeStr + `,
		"labels": [
			{
				"id": 1
			}
		],
		"user": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"draft": false,
		"merged": false,
		"mergeable": false,
		"mergeable_state": "ms",
		"merged_by": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"merge_commit_sha": "mcs",
		"rebaseable": false,
		"comments": 1,
		"commits": 1,
		"additions": 1,
		"deletions": 1,
		"changed_files": 1,
		"url": "url",
		"html_url": "hurl",
		"issue_url": "iurl",
		"statuses_url": "surl",
		"diff_url": "durl",
		"patch_url": "purl",
		"commits_url": "curl",
		"comments_url": "comurl",
		"review_comments_url": "rcurls",
		"review_comment_url": "rcurl",
		"review_comments": 1,
		"assignee": {
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
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"assignees": [
			{
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
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			}
		],
		"milestone": {
			"id": 1
		},
		"maintainer_can_modify": true,
		"author_association": "aa",
		"node_id": "nid",
		"requested_reviewers": [
			{
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
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			}
		],
		"auto_merge": {
			"enabled_by": {
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
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			},
			"merge_method": "mm",
			"commit_title": "ct",
			"commit_message": "cm"
		},
		"requested_teams": [
			{
				"id": 1
			}
		],
		"_links": {
			"self": {
				"href": "href"
			},
			"html": {
				"href": "href"
			},
			"issue": {
				"href": "href"
			},
			"comments": {
				"href": "href"
			},
			"review_comments": {
				"href": "href"
			},
			"review_comment": {
				"href": "href"
			},
			"commits": {
				"href": "href"
			},
			"statuses": {
				"href": "href"
			}
		},
		"head": {
			"ref": "r2",
			"repo": {
				"id": 2
			}
		},
		"base": {
			"ref": "r2",
			"repo": {
				"id": 2
			}
		},
		"active_lock_reason": "alr"
	}`

	testJSONMarshal(t, u, want)
}
