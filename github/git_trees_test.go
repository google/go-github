// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshalJSON_withNilContentAndSHA(t *testing.T) {
	te := &TreeEntry{
		Path: String("path"),
		Mode: String("mode"),
		Type: String("type"),
		Size: Int(1),
		URL:  String("url"),
	}

	got, err := te.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON: %v", err)
	}

	want := `{"sha":null,"path":"path","mode":"mode","type":"type"}`
	if string(got) != want {
		t.Errorf("MarshalJSON = %s, want %v", got, want)
	}
}

func TestGitService_GetTree(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/trees/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			  "sha": "s",
			  "tree": [ { "type": "blob" } ],
			  "truncated": true
			}`)
	})

	ctx := context.Background()
	tree, _, err := client.Git.GetTree(ctx, "o", "r", "s", true)
	if err != nil {
		t.Errorf("Git.GetTree returned error: %v", err)
	}

	want := Tree{
		SHA: String("s"),
		Entries: []*TreeEntry{
			{
				Type: String("blob"),
			},
		},
		Truncated: Bool(true),
	}
	if !cmp.Equal(*tree, want) {
		t.Errorf("Tree.Get returned %+v, want %+v", *tree, want)
	}

	const methodName = "GetTree"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetTree(ctx, "\n", "\n", "\n", true)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetTree(ctx, "o", "r", "s", true)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_GetTree_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Git.GetTree(ctx, "%", "%", "%", false)
	testURLParseError(t, err)
}

func TestGitService_CreateTree(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []*TreeEntry{
		{
			Path: String("file.rb"),
			Mode: String("100644"),
			Type: String("blob"),
			SHA:  String("7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"),
		},
	}

	mux.HandleFunc("/repos/o/r/git/trees", func(w http.ResponseWriter, r *http.Request) {
		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unable to read body: %v", err)
		}

		testMethod(t, r, "POST")

		want := []byte(`{"base_tree":"b","tree":[{"sha":"7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b","path":"file.rb","mode":"100644","type":"blob"}]}` + "\n")
		if !bytes.Equal(got, want) {
			t.Errorf("Git.CreateTree request body: %s, want %s", got, want)
		}

		fmt.Fprint(w, `{
		  "sha": "cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		  "tree": [
		    {
		      "path": "file.rb",
		      "mode": "100644",
		      "type": "blob",
		      "size": 132,
		      "sha": "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"
		    }
		  ]
		}`)
	})

	ctx := context.Background()
	tree, _, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Git.CreateTree returned error: %v", err)
	}

	want := Tree{
		String("cd8274d15fa3ae2ab983129fb037999f264ba9a7"),
		[]*TreeEntry{
			{
				Path: String("file.rb"),
				Mode: String("100644"),
				Type: String("blob"),
				Size: Int(132),
				SHA:  String("7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"),
			},
		},
		nil,
	}

	if !cmp.Equal(*tree, want) {
		t.Errorf("Git.CreateTree returned %+v, want %+v", *tree, want)
	}

	const methodName = "CreateTree"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateTree(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateTree_Content(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []*TreeEntry{
		{
			Path:    String("content.md"),
			Mode:    String("100644"),
			Content: String("file content"),
		},
	}

	mux.HandleFunc("/repos/o/r/git/trees", func(w http.ResponseWriter, r *http.Request) {
		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unable to read body: %v", err)
		}

		testMethod(t, r, "POST")

		want := []byte(`{"base_tree":"b","tree":[{"path":"content.md","mode":"100644","content":"file content"}]}` + "\n")
		if !bytes.Equal(got, want) {
			t.Errorf("Git.CreateTree request body: %s, want %s", got, want)
		}

		fmt.Fprint(w, `{
		  "sha": "5c6780ad2c68743383b740fd1dab6f6a33202b11",
		  "url": "https://api.github.com/repos/o/r/git/trees/5c6780ad2c68743383b740fd1dab6f6a33202b11",
		  "tree": [
		    {
			  "mode": "100644",
			  "type": "blob",
			  "sha":  "aad8feacf6f8063150476a7b2bd9770f2794c08b",
			  "path": "content.md",
			  "size": 12,
			  "url": "https://api.github.com/repos/o/r/git/blobs/aad8feacf6f8063150476a7b2bd9770f2794c08b"
		    }
		  ]
		}`)
	})

	ctx := context.Background()
	tree, _, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Git.CreateTree returned error: %v", err)
	}

	want := Tree{
		String("5c6780ad2c68743383b740fd1dab6f6a33202b11"),
		[]*TreeEntry{
			{
				Path: String("content.md"),
				Mode: String("100644"),
				Type: String("blob"),
				Size: Int(12),
				SHA:  String("aad8feacf6f8063150476a7b2bd9770f2794c08b"),
				URL:  String("https://api.github.com/repos/o/r/git/blobs/aad8feacf6f8063150476a7b2bd9770f2794c08b"),
			},
		},
		nil,
	}

	if !cmp.Equal(*tree, want) {
		t.Errorf("Git.CreateTree returned %+v, want %+v", *tree, want)
	}

	const methodName = "CreateTree"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateTree(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateTree_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []*TreeEntry{
		{
			Path: String("content.md"),
			Mode: String("100644"),
		},
	}

	mux.HandleFunc("/repos/o/r/git/trees", func(w http.ResponseWriter, r *http.Request) {
		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unable to read body: %v", err)
		}

		testMethod(t, r, "POST")

		want := []byte(`{"base_tree":"b","tree":[{"sha":null,"path":"content.md","mode":"100644"}]}` + "\n")
		if !bytes.Equal(got, want) {
			t.Errorf("Git.CreateTree request body: %s, want %s", got, want)
		}

		fmt.Fprint(w, `{
		  "sha": "5c6780ad2c68743383b740fd1dab6f6a33202b11",
		  "url": "https://api.github.com/repos/o/r/git/trees/5c6780ad2c68743383b740fd1dab6f6a33202b11",
		  "tree": [
		    {
			  "mode": "100644",
			  "type": "blob",
			  "sha":  null,
			  "path": "content.md",
			  "size": 12,
			  "url": "https://api.github.com/repos/o/r/git/blobs/aad8feacf6f8063150476a7b2bd9770f2794c08b"
		    }
		  ]
		}`)
	})

	ctx := context.Background()
	tree, _, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
	if err != nil {
		t.Errorf("Git.CreateTree returned error: %v", err)
	}

	want := Tree{
		String("5c6780ad2c68743383b740fd1dab6f6a33202b11"),
		[]*TreeEntry{
			{
				Path: String("content.md"),
				Mode: String("100644"),
				Type: String("blob"),
				Size: Int(12),
				SHA:  nil,
				URL:  String("https://api.github.com/repos/o/r/git/blobs/aad8feacf6f8063150476a7b2bd9770f2794c08b"),
			},
		},
		nil,
	}

	if !cmp.Equal(*tree, want) {
		t.Errorf("Git.CreateTree returned %+v, want %+v", *tree, want)
	}

	const methodName = "CreateTree"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateTree(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateTree(ctx, "o", "r", "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateTree_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Git.CreateTree(ctx, "%", "%", "", nil)
	testURLParseError(t, err)
}

func TestTree_Marshal(t *testing.T) {
	testJSONMarshal(t, &Tree{}, "{}")

	u := &Tree{
		SHA: String("sha"),
		Entries: []*TreeEntry{
			{
				SHA:     String("sha"),
				Path:    String("path"),
				Mode:    String("mode"),
				Type:    String("type"),
				Size:    Int(1),
				Content: String("content"),
				URL:     String("url"),
			},
		},
		Truncated: Bool(false),
	}

	want := `{
		"sha": "sha",
		"tree": [
			{
				"sha": "sha",
				"path": "path",
				"mode": "mode",
				"type": "type",
				"size": 1,
				"content": "content",
				"url": "url"
			}
		],
		"truncated": false
	}`

	testJSONMarshal(t, u, want)
}

func TestTreeEntry_Marshal(t *testing.T) {
	testJSONMarshal(t, &TreeEntry{}, "{}")

	u := &TreeEntry{
		SHA:     String("sha"),
		Path:    String("path"),
		Mode:    String("mode"),
		Type:    String("type"),
		Size:    Int(1),
		Content: String("content"),
		URL:     String("url"),
	}

	want := `{
		"sha": "sha",
		"path": "path",
		"mode": "mode",
		"type": "type",
		"size": 1,
		"content": "content",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestTreeEntryWithFileDelete_Marshal(t *testing.T) {
	testJSONMarshal(t, &treeEntryWithFileDelete{}, "{}")

	u := &treeEntryWithFileDelete{
		SHA:     String("sha"),
		Path:    String("path"),
		Mode:    String("mode"),
		Type:    String("type"),
		Size:    Int(1),
		Content: String("content"),
		URL:     String("url"),
	}

	want := `{
		"sha": "sha",
		"path": "path",
		"mode": "mode",
		"type": "type",
		"size": 1,
		"content": "content",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateTree_Marshal(t *testing.T) {
	testJSONMarshal(t, &createTree{}, "{}")

	u := &createTree{
		BaseTree: "bt",
		Entries:  []interface{}{"e"},
	}

	want := `{
		"base_tree": "bt",
		"tree": ["e"]
	}`

	testJSONMarshal(t, u, want)
}
