// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGist_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Gist{}, "{}")

	createdAt := time.Date(2010, time.February, 10, 10, 10, 0, 0, time.UTC)
	updatedAt := time.Date(2010, time.February, 10, 10, 10, 0, 0, time.UTC)

	u := &Gist{
		ID:          Ptr("i"),
		Description: Ptr("description"),
		Public:      Ptr(true),
		Owner: &User{
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
		},
		Files: map[GistFilename]GistFile{
			"gistfile.py": {
				Size:     Ptr(167),
				Filename: Ptr("gistfile.py"),
				Language: Ptr("Python"),
				Type:     Ptr("application/x-python"),
				RawURL:   Ptr("raw-url"),
				Content:  Ptr("c"),
			},
		},
		Comments:   Ptr(1),
		HTMLURL:    Ptr("html-url"),
		GitPullURL: Ptr("gitpull-url"),
		GitPushURL: Ptr("gitpush-url"),
		CreatedAt:  &Timestamp{createdAt},
		UpdatedAt:  &Timestamp{updatedAt},
		NodeID:     Ptr("node"),
	}

	want := `{
		"id": "i",
		"description": "description",
		"public": true,
		"owner": {
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
		"files": {
			"gistfile.py": {
				"size": 167,
				"filename": "gistfile.py",
				"language": "Python",
				"type": "application/x-python",
				"raw_url": "raw-url",
				"content": "c"
			}
		},
		"comments": 1,
		"html_url": "html-url",
		"git_pull_url": "gitpull-url",
		"git_push_url": "gitpush-url",
		"created_at": "2010-02-10T10:10:00Z",
		"updated_at": "2010-02-10T10:10:00Z",
		"node_id": "node"
	}`

	testJSONMarshal(t, u, want)
}

func TestGistCommit_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GistCommit{}, "{}")

	u := &GistCommit{
		URL:     Ptr("u"),
		Version: Ptr("v"),
		User: &User{
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
		},
		ChangeStatus: &CommitStats{
			Additions: Ptr(1),
			Deletions: Ptr(1),
			Total:     Ptr(2),
		},
		CommittedAt: &Timestamp{referenceTime},
		NodeID:      Ptr("node"),
	}

	want := `{
		"url": "u",
		"version": "v",
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
		"change_status": {
			"additions": 1,
			"deletions": 1,
			"total": 2
		},
		"committed_at": ` + referenceTimeStr + `,
		"node_id": "node"
	}`

	testJSONMarshal(t, u, want)
}

func TestGistFork_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GistFork{}, "{}")

	u := &GistFork{
		URL: Ptr("u"),
		User: &User{
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
		},
		ID:        Ptr("id"),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		NodeID:    Ptr("node"),
	}

	want := `{
		"url": "u",
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
		"id": "id",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"node_id": "node"
	}`

	testJSONMarshal(t, u, want)
}

func TestGistsService_List_specifiedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	since := "2013-01-01T00:00:00Z"

	mux.HandleFunc("/users/u/gists", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": since,
		})
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	opt := &GistListOptions{Since: time.Date(2013, time.January, 1, 0, 0, 0, 0, time.UTC)}
	ctx := context.Background()
	gists, _, err := client.Gists.List(ctx, "u", opt)
	if err != nil {
		t.Errorf("Gists.List returned error: %v", err)
	}

	want := []*Gist{{ID: Ptr("1")}}
	if !cmp.Equal(gists, want) {
		t.Errorf("Gists.List returned %+v, want %+v", gists, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.List(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.List(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_List_authenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	ctx := context.Background()
	gists, _, err := client.Gists.List(ctx, "", nil)
	if err != nil {
		t.Errorf("Gists.List returned error: %v", err)
	}

	want := []*Gist{{ID: Ptr("1")}}
	if !cmp.Equal(gists, want) {
		t.Errorf("Gists.List returned %+v, want %+v", gists, want)
	}

	const methodName = "List"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.List(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.List(ctx, "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_List_invalidUser(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.List(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_ListAll(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	since := "2013-01-01T00:00:00Z"

	mux.HandleFunc("/gists/public", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": since,
		})
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	opt := &GistListOptions{Since: time.Date(2013, time.January, 1, 0, 0, 0, 0, time.UTC)}
	ctx := context.Background()
	gists, _, err := client.Gists.ListAll(ctx, opt)
	if err != nil {
		t.Errorf("Gists.ListAll returned error: %v", err)
	}

	want := []*Gist{{ID: Ptr("1")}}
	if !cmp.Equal(gists, want) {
		t.Errorf("Gists.ListAll returned %+v, want %+v", gists, want)
	}

	const methodName = "ListAll"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListAll(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListStarred(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	since := "2013-01-01T00:00:00Z"

	mux.HandleFunc("/gists/starred", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": since,
		})
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	opt := &GistListOptions{Since: time.Date(2013, time.January, 1, 0, 0, 0, 0, time.UTC)}
	ctx := context.Background()
	gists, _, err := client.Gists.ListStarred(ctx, opt)
	if err != nil {
		t.Errorf("Gists.ListStarred returned error: %v", err)
	}

	want := []*Gist{{ID: Ptr("1")}}
	if !cmp.Equal(gists, want) {
		t.Errorf("Gists.ListStarred returned %+v, want %+v", gists, want)
	}

	const methodName = "ListStarred"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListStarred(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_Get(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Get(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Get returned error: %v", err)
	}

	want := &Gist{ID: Ptr("1")}
	if !cmp.Equal(gist, want) {
		t.Errorf("Gists.Get returned %+v, want %+v", gist, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.Get(ctx, "1")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_Get_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.Get(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_GetRevision(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.GetRevision(ctx, "1", "s")
	if err != nil {
		t.Errorf("Gists.Get returned error: %v", err)
	}

	want := &Gist{ID: Ptr("1")}
	if !cmp.Equal(gist, want) {
		t.Errorf("Gists.Get returned %+v, want %+v", gist, want)
	}

	const methodName = "GetRevision"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.GetRevision(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.GetRevision(ctx, "1", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_GetRevision_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.GetRevision(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestGistsService_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Gist{
		Description: Ptr("Gist description"),
		Public:      Ptr(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Content: Ptr("Gist file content")},
		},
	}

	mux.HandleFunc("/gists", func(w http.ResponseWriter, r *http.Request) {
		v := new(Gist)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "1",
				"description": "Gist description",
				"public": false,
				"files": {
					"test.txt": {
						"filename": "test.txt"
					}
				}
			}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Create(ctx, input)
	if err != nil {
		t.Errorf("Gists.Create returned error: %v", err)
	}

	want := &Gist{
		ID:          Ptr("1"),
		Description: Ptr("Gist description"),
		Public:      Ptr(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Filename: Ptr("test.txt")},
		},
	}
	if !cmp.Equal(gist, want) {
		t.Errorf("Gists.Create returned %+v, want %+v", gist, want)
	}

	const methodName = "Create"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.Create(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_Edit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Gist{
		Description: Ptr("New description"),
		Files: map[GistFilename]GistFile{
			"new.txt": {Content: Ptr("new file content")},
		},
	}

	mux.HandleFunc("/gists/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Gist)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "1",
				"description": "new description",
				"public": false,
				"files": {
					"test.txt": {
						"filename": "test.txt"
					},
					"new.txt": {
						"filename": "new.txt"
					}
				}
			}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Edit(ctx, "1", input)
	if err != nil {
		t.Errorf("Gists.Edit returned error: %v", err)
	}

	want := &Gist{
		ID:          Ptr("1"),
		Description: Ptr("new description"),
		Public:      Ptr(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Filename: Ptr("test.txt")},
			"new.txt":  {Filename: Ptr("new.txt")},
		},
	}
	if !cmp.Equal(gist, want) {
		t.Errorf("Gists.Edit returned %+v, want %+v", gist, want)
	}

	const methodName = "Edit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.Edit(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.Edit(ctx, "1", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_Edit_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.Edit(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_ListCommits(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, nil)
		fmt.Fprint(w, `
		  [
		    {
		      "url": "https://api.github.com/gists/1/1",
		      "version": "1",
		      "user": {
		        "id": 1
		      },
		      "change_status": {
		        "deletions": 0,
		        "additions": 180,
		        "total": 180
		      },
		      "committed_at": "2010-01-01T00:00:00Z"
		    }
		  ]
		`)
	})

	ctx := context.Background()
	gistCommits, _, err := client.Gists.ListCommits(ctx, "1", nil)
	if err != nil {
		t.Errorf("Gists.ListCommits returned error: %v", err)
	}

	want := []*GistCommit{{
		URL:         Ptr("https://api.github.com/gists/1/1"),
		Version:     Ptr("1"),
		User:        &User{ID: Ptr(int64(1))},
		CommittedAt: &Timestamp{time.Date(2010, time.January, 1, 00, 00, 00, 0, time.UTC)},
		ChangeStatus: &CommitStats{
			Additions: Ptr(180),
			Deletions: Ptr(0),
			Total:     Ptr(180),
		}}}

	if !cmp.Equal(gistCommits, want) {
		t.Errorf("Gists.ListCommits returned %+v, want %+v", gistCommits, want)
	}

	const methodName = "ListCommits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.ListCommits(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListCommits(ctx, "1", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListCommits_withOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, err := client.Gists.ListCommits(ctx, "1", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Gists.ListCommits returned error: %v", err)
	}

	const methodName = "ListCommits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.ListCommits(ctx, "\n", &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListCommits(ctx, "1", &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListCommits_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.ListCommits(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_Delete(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Gists.Delete(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Delete returned error: %v", err)
	}

	const methodName = "Delete"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Gists.Delete(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Gists.Delete(ctx, "1")
	})
}

func TestGistsService_Delete_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Gists.Delete(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Star(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/star", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	ctx := context.Background()
	_, err := client.Gists.Star(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Star returned error: %v", err)
	}

	const methodName = "Star"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Gists.Star(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Gists.Star(ctx, "1")
	})
}

func TestGistsService_Star_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Gists.Star(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Unstar(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/star", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Gists.Unstar(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Unstar returned error: %v", err)
	}

	const methodName = "Unstar"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Gists.Unstar(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Gists.Unstar(ctx, "1")
	})
}

func TestGistsService_Unstar_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Gists.Unstar(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_IsStarred_hasStar(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/star", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	star, _, err := client.Gists.IsStarred(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Starred returned error: %v", err)
	}
	if want := true; star != want {
		t.Errorf("Gists.Starred returned %+v, want %+v", star, want)
	}

	const methodName = "IsStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.IsStarred(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.IsStarred(ctx, "1")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_IsStarred_noStar(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/star", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	star, _, err := client.Gists.IsStarred(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Starred returned error: %v", err)
	}
	if want := false; star != want {
		t.Errorf("Gists.Starred returned %+v, want %+v", star, want)
	}

	const methodName = "IsStarred"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.IsStarred(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.IsStarred(ctx, "1")
		if got {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want false", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_IsStarred_invalidID(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Gists.IsStarred(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Fork(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id": "2"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Fork(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Fork returned error: %v", err)
	}

	want := &Gist{ID: Ptr("2")}
	if !cmp.Equal(gist, want) {
		t.Errorf("Gists.Fork returned %+v, want %+v", gist, want)
	}

	const methodName = "Fork"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.Fork(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.Fork(ctx, "1")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListForks(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, nil)
		fmt.Fprint(w, `
		  [
		    {"url": "https://api.github.com/gists/1",
		     "user": {"id": 1},
		     "id": "1",
		     "created_at": "2010-01-01T00:00:00Z",
		     "updated_at": "2013-01-01T00:00:00Z"
		    }
		  ]
		`)
	})

	ctx := context.Background()
	gistForks, _, err := client.Gists.ListForks(ctx, "1", nil)
	if err != nil {
		t.Errorf("Gists.ListForks returned error: %v", err)
	}

	want := []*GistFork{{
		URL:       Ptr("https://api.github.com/gists/1"),
		ID:        Ptr("1"),
		User:      &User{ID: Ptr(int64(1))},
		CreatedAt: &Timestamp{time.Date(2010, time.January, 1, 00, 00, 00, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2013, time.January, 1, 00, 00, 00, 0, time.UTC)}}}

	if !cmp.Equal(gistForks, want) {
		t.Errorf("Gists.ListForks returned %+v, want %+v", gistForks, want)
	}

	const methodName = "ListForks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.ListForks(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListForks(ctx, "1", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistsService_ListForks_withOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/gists/1/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	gistForks, _, err := client.Gists.ListForks(ctx, "1", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Gists.ListForks returned error: %v", err)
	}

	want := []*GistFork{}
	if !cmp.Equal(gistForks, want) {
		t.Errorf("Gists.ListForks returned %+v, want %+v", gistForks, want)
	}

	const methodName = "ListForks"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Gists.ListForks(ctx, "\n", &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Gists.ListForks(ctx, "1", &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGistFile_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GistFile{}, "{}")

	u := &GistFile{
		Size:     Ptr(1),
		Filename: Ptr("fn"),
		Language: Ptr("lan"),
		Type:     Ptr("type"),
		RawURL:   Ptr("rurl"),
		Content:  Ptr("con"),
	}

	want := `{
		"size": 1,
		"filename": "fn",
		"language": "lan",
		"type": "type",
		"raw_url": "rurl",
		"content": "con"
	}`

	testJSONMarshal(t, u, want)
}
