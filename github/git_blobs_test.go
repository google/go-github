// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGitService_GetBlob(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/git/blobs/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			  "sha": "s",
			  "content": "blob content"
			}`)
	})

	ctx := context.Background()
	blob, _, err := client.Git.GetBlob(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetBlob returned error: %v", err)
	}

	want := Blob{
		SHA:     Ptr("s"),
		Content: Ptr("blob content"),
	}

	if !cmp.Equal(*blob, want) {
		t.Errorf("Blob.Get returned %+v, want %+v", *blob, want)
	}

	const methodName = "GetBlob"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetBlob(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetBlob(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_GetBlob_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Git.GetBlob(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestGitService_GetBlobRaw(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/git/blobs/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "application/vnd.github.v3.raw")

		fmt.Fprint(w, `raw contents here`)
	})

	ctx := context.Background()
	blob, _, err := client.Git.GetBlobRaw(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetBlobRaw returned error: %v", err)
	}

	want := []byte("raw contents here")
	if !bytes.Equal(blob, want) {
		t.Errorf("GetBlobRaw returned %q, want %q", blob, want)
	}

	const methodName = "GetBlobRaw"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetBlobRaw(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetBlobRaw(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateBlob(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Blob{
		SHA:      Ptr("s"),
		Content:  Ptr("blob content"),
		Encoding: Ptr("utf-8"),
		Size:     Ptr(12),
	}

	mux.HandleFunc("/repos/o/r/git/blobs", func(w http.ResponseWriter, r *http.Request) {
		v := new(Blob)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")

		want := input
		if !cmp.Equal(v, want) {
			t.Errorf("Git.CreateBlob request body: %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{
		 "sha": "s",
		 "content": "blob content",
		 "encoding": "utf-8",
		 "size": 12
		}`)
	})

	ctx := context.Background()
	blob, _, err := client.Git.CreateBlob(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateBlob returned error: %v", err)
	}

	want := input

	if !cmp.Equal(*blob, *want) {
		t.Errorf("Git.CreateBlob returned %+v, want %+v", *blob, *want)
	}

	const methodName = "CreateBlob"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateBlob(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateBlob(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateBlob_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Git.CreateBlob(ctx, "%", "%", &Blob{})
	testURLParseError(t, err)
}

func TestBlob_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Blob{}, "{}")

	u := &Blob{
		Content:  Ptr("content"),
		Encoding: Ptr("encoding"),
		SHA:      Ptr("sha"),
		Size:     Ptr(1),
		URL:      Ptr("url"),
		NodeID:   Ptr("nid"),
	}

	want := `{
		"content": "content",
		"encoding": "encoding",
		"sha": "sha",
		"size": 1,
		"url": "url",
		"node_id": "nid"
	}`

	testJSONMarshal(t, u, want)
}
