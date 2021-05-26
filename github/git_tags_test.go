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

	"github.com/google/go-cmp/cmp"
)

func TestGitService_GetTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/tags/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"tag": "t"}`)
	})

	ctx := context.Background()
	tag, _, err := client.Git.GetTag(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetTag returned error: %v", err)
	}

	want := &Tag{Tag: String("t")}
	if !cmp.Equal(tag, want) {
		t.Errorf("Git.GetTag returned %+v, want %+v", tag, want)
	}

	const methodName = "GetTag"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetTag(ctx, "\n", "\n", "\n")
		return err
	})
}

func TestGitService_CreateTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &createTagRequest{Tag: String("t"), Object: String("s")}

	mux.HandleFunc("/repos/o/r/git/tags", func(w http.ResponseWriter, r *http.Request) {
		v := new(createTagRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"tag": "t"}`)
	})

	ctx := context.Background()
	tag, _, err := client.Git.CreateTag(ctx, "o", "r", &Tag{
		Tag:    input.Tag,
		Object: &GitObject{SHA: input.Object},
	})
	if err != nil {
		t.Errorf("Git.CreateTag returned error: %v", err)
	}

	want := &Tag{Tag: String("t")}
	if !cmp.Equal(tag, want) {
		t.Errorf("Git.GetTag returned %+v, want %+v", tag, want)
	}

	const methodName = "CreateTag"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateTag(ctx, "\n", "\n", &Tag{
			Tag:    input.Tag,
			Object: &GitObject{SHA: input.Object},
		})
		return err
	})
}
