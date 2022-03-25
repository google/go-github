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

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetTag(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
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
	inputTag := &Tag{
		Tag:    input.Tag,
		Object: &GitObject{SHA: input.Object},
	}
	tag, _, err := client.Git.CreateTag(ctx, "o", "r", inputTag)
	if err != nil {
		t.Errorf("Git.CreateTag returned error: %v", err)
	}

	want := &Tag{Tag: String("t")}
	if !cmp.Equal(tag, want) {
		t.Errorf("Git.GetTag returned %+v, want %+v", tag, want)
	}

	const methodName = "CreateTag"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateTag(ctx, "\n", "\n", inputTag)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateTag(ctx, "o", "r", inputTag)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestTag_Marshal(t *testing.T) {
	testJSONMarshal(t, &Tag{}, "{}")

	u := &Tag{
		Tag:     String("tag"),
		SHA:     String("sha"),
		URL:     String("url"),
		Message: String("msg"),
		Tagger: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
		},
		Object: &GitObject{
			Type: String("type"),
			SHA:  String("sha"),
			URL:  String("url"),
		},
		Verification: &SignatureVerification{
			Verified:  Bool(true),
			Reason:    String("reason"),
			Signature: String("sign"),
			Payload:   String("payload"),
		},
		NodeID: String("nid"),
	}

	want := `{
		"tag": "tag",
		"sha": "sha",
		"url": "url",
		"message": "msg",
		"tagger": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		},
		"object": {
			"type": "type",
			"sha": "sha",
			"url": "url"
		},
		"verification": {
			"verified": true,
			"reason": "reason",
			"signature": "sign",
			"payload": "payload"
		},
		"node_id": "nid"
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateTagRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &createTagRequest{}, "{}")

	u := &createTagRequest{
		Tag:     String("tag"),
		Message: String("msg"),
		Object:  String("obj"),
		Type:    String("type"),
		Tagger: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
		},
	}

	want := `{
		"tag": "tag",
		"message": "msg",
		"object": "obj",
		"type": "type",
		"tagger": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		}
	}`

	testJSONMarshal(t, u, want)
}
