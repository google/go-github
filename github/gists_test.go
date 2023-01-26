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
	testJSONMarshal(t, &Gist{}, "{}")

	createdAt := time.Date(2010, time.February, 10, 10, 10, 0, 0, time.UTC)
	updatedAt := time.Date(2010, time.February, 10, 10, 10, 0, 0, time.UTC)

	u := &Gist{
		ID:          String("i"),
		Description: String("description"),
		Public:      Bool(true),
		Owner: &User{
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
		},
		Files: map[GistFilename]GistFile{
			"gistfile.py": {
				Size:     Int(167),
				Filename: String("gistfile.py"),
				Language: String("Python"),
				Type:     String("application/x-python"),
				RawURL:   String("raw-url"),
				Content:  String("c"),
			},
		},
		Comments:   Int(1),
		HTMLURL:    String("html-url"),
		GitPullURL: String("gitpull-url"),
		GitPushURL: String("gitpush-url"),
		CreatedAt:  &Timestamp{createdAt},
		UpdatedAt:  &Timestamp{updatedAt},
		NodeID:     String("node"),
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
	testJSONMarshal(t, &GistCommit{}, "{}")

	u := &GistCommit{
		URL:     String("u"),
		Version: String("v"),
		User: &User{
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
		},
		ChangeStatus: &CommitStats{
			Additions: Int(1),
			Deletions: Int(1),
			Total:     Int(2),
		},
		CommittedAt: &Timestamp{referenceTime},
		NodeID:      String("node"),
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
	testJSONMarshal(t, &GistFork{}, "{}")

	u := &GistFork{
		URL: String("u"),
		User: &User{
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
		},
		ID:        String("id"),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		NodeID:    String("node"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	want := []*Gist{{ID: String("1")}}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	ctx := context.Background()
	gists, _, err := client.Gists.List(ctx, "", nil)
	if err != nil {
		t.Errorf("Gists.List returned error: %v", err)
	}

	want := []*Gist{{ID: String("1")}}
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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.List(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	want := []*Gist{{ID: String("1")}}
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	want := []*Gist{{ID: String("1")}}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Get(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Get returned error: %v", err)
	}

	want := &Gist{ID: String("1")}
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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.Get(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_GetRevision(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.GetRevision(ctx, "1", "s")
	if err != nil {
		t.Errorf("Gists.Get returned error: %v", err)
	}

	want := &Gist{ID: String("1")}
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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.GetRevision(ctx, "%", "%")
	testURLParseError(t, err)
}

func TestGistsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Gist{
		Description: String("Gist description"),
		Public:      Bool(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Content: String("Gist file content")},
		},
	}

	mux.HandleFunc("/gists", func(w http.ResponseWriter, r *http.Request) {
		v := new(Gist)
		json.NewDecoder(r.Body).Decode(v)

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
		ID:          String("1"),
		Description: String("Gist description"),
		Public:      Bool(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Filename: String("test.txt")},
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Gist{
		Description: String("New description"),
		Files: map[GistFilename]GistFile{
			"new.txt": {Content: String("new file content")},
		},
	}

	mux.HandleFunc("/gists/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Gist)
		json.NewDecoder(r.Body).Decode(v)

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
		ID:          String("1"),
		Description: String("new description"),
		Public:      Bool(false),
		Files: map[GistFilename]GistFile{
			"test.txt": {Filename: String("test.txt")},
			"new.txt":  {Filename: String("new.txt")},
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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.Edit(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_ListCommits(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
		URL:         String("https://api.github.com/gists/1/1"),
		Version:     String("1"),
		User:        &User{ID: Int64(1)},
		CommittedAt: &Timestamp{time.Date(2010, time.January, 1, 00, 00, 00, 0, time.UTC)},
		ChangeStatus: &CommitStats{
			Additions: Int(180),
			Deletions: Int(0),
			Total:     Int(180),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.ListCommits(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestGistsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Gists.Delete(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Star(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Gists.Star(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Unstar(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Gists.Unstar(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_IsStarred_hasStar(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Gists.IsStarred(ctx, "%")
	testURLParseError(t, err)
}

func TestGistsService_Fork(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/gists/1/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id": "2"}`)
	})

	ctx := context.Background()
	gist, _, err := client.Gists.Fork(ctx, "1")
	if err != nil {
		t.Errorf("Gists.Fork returned error: %v", err)
	}

	want := &Gist{ID: String("2")}
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
	client, mux, _, teardown := setup()
	defer teardown()

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
		URL:       String("https://api.github.com/gists/1"),
		ID:        String("1"),
		User:      &User{ID: Int64(1)},
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
	client, mux, _, teardown := setup()
	defer teardown()

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
	testJSONMarshal(t, &GistFile{}, "{}")

	u := &GistFile{
		Size:     Int(1),
		Filename: String("fn"),
		Language: String("lan"),
		Type:     String("type"),
		RawURL:   String("rurl"),
		Content:  String("con"),
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
