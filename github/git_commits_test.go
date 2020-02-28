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
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/openpgp"
)

func TestCommit_Marshal(t *testing.T) {
	testJSONMarshal(t, &Commit{}, "{}")

	u := &Commit{
		SHA: String("s"),
		Author: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		Committer: &CommitAuthor{
			Date:  &referenceTime,
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		Message: String("m"),
		Tree: &Tree{
			SHA: String("s"),
			Entries: []*TreeEntry{{
				SHA:     String("s"),
				Path:    String("p"),
				Mode:    String("m"),
				Type:    String("t"),
				Size:    Int(1),
				Content: String("c"),
				URL:     String("u"),
			}},
			Truncated: Bool(false),
		},
		Parents: nil,
		Stats: &CommitStats{
			Additions: Int(1),
			Deletions: Int(1),
			Total:     Int(1),
		},
		HTMLURL: String("h"),
		URL:     String("u"),
		Verification: &SignatureVerification{
			Verified:  Bool(false),
			Reason:    String("r"),
			Signature: String("s"),
			Payload:   String("p"),
		},
		NodeID:       String("n"),
		CommentCount: Int(1),
		SigningKey:   &openpgp.Entity{},
	}

	want := `{
		"sha": "s",
		"author": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"committer": {
			"date": ` + referenceTimeStr + `,
			"name": "n",
			"email": "e",
			"username": "u"
		},
		"message": "m",
		"tree": {
			"sha": "s",
			"tree": [
				{
					"sha": "s",
					"path": "p",
					"mode": "m",
					"type": "t",
					"size": 1,
					"content": "c",
					"url": "u"
				}
			],
			"truncated": false
		},
		"stats": {
			"additions": 1,
			"deletions": 1,
			"total": 1
		},
		"html_url": "h",
		"url": "u",
		"verification": {
			"verified": false,
			"reason": "r",
			"signature": "s",
			"payload": "p"
		},
		"node_id": "n",
		"comment_count": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestGitService_GetCommit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"sha":"s","message":"Commit Message.","author":{"name":"n"}}`)
	})

	commit, _, err := client.Git.GetCommit(context.Background(), "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("s"), Message: String("Commit Message."), Author: &CommitAuthor{Name: String("n")}}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.GetCommit returned %+v, want %+v", commit, want)
	}
}

func TestGitService_GetCommit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Git.GetCommit(context.Background(), "%", "%", "%")
	testURLParseError(t, err)
}

func TestGitService_CreateCommit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Commit{
		Message: String("Commit Message."),
		Tree:    &Tree{SHA: String("t")},
		Parents: []*Commit{{SHA: String("p")}},
	}

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")

		want := &createCommit{
			Message: input.Message,
			Tree:    String("t"),
			Parents: []string{"p"},
		}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"s"}`)
	})

	commit, _, err := client.Git.CreateCommit(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("s")}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}
}

func TestGitService_CreateSignedCommit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	signature := "----- BEGIN PGP SIGNATURE -----\n\naaaa\naaaa\n----- END PGP SIGNATURE -----"

	input := &Commit{
		Message: String("Commit Message."),
		Tree:    &Tree{SHA: String("t")},
		Parents: []*Commit{{SHA: String("p")}},
		Verification: &SignatureVerification{
			Signature: String(signature),
		},
	}

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")

		want := &createCommit{
			Message:   input.Message,
			Tree:      String("t"),
			Parents:   []string{"p"},
			Signature: String(signature),
		}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"commitSha"}`)
	})

	commit, _, err := client.Git.CreateCommit(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("commitSha")}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}
}
func TestGitService_CreateSignedCommitWithInvalidParams(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	input := &Commit{
		SigningKey: &openpgp.Entity{},
	}

	_, _, err := client.Git.CreateCommit(context.Background(), "o", "r", input)
	if err == nil {
		t.Errorf("Expected error to be returned because invalid params was passed")
	}
}

func TestGitService_CreateSignedCommitWithKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	s := strings.NewReader(testGPGKey)
	keyring, err := openpgp.ReadArmoredKeyRing(s)
	if err != nil {
		t.Errorf("Error reading keyring: %+v", err)
	}

	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")
	author := CommitAuthor{
		Name:  String("go-github"),
		Email: String("go-github@github.com"),
		Date:  &date,
	}
	input := &Commit{
		Message:    String("Commit Message."),
		Tree:       &Tree{SHA: String("t")},
		Parents:    []*Commit{{SHA: String("p")}},
		SigningKey: keyring[0],
		Author:     &author,
	}

	messageReader := strings.NewReader(`tree t
parent p
author go-github <go-github@github.com> 1493849023 +0200
committer go-github <go-github@github.com> 1493849023 +0200

Commit Message.`)

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")

		want := &createCommit{
			Message: input.Message,
			Tree:    String("t"),
			Parents: []string{"p"},
			Author:  &author,
		}

		sigReader := strings.NewReader(*v.Signature)
		signer, err := openpgp.CheckArmoredDetachedSignature(keyring, messageReader, sigReader)
		if err != nil {
			t.Errorf("Error verifying signature: %+v", err)
		}
		if signer.Identities["go-github <go-github@github.com>"].Name != "go-github <go-github@github.com>" {
			t.Errorf("Signer is incorrect. got: %+v, want %+v", signer.Identities["go-github <go-github@github.com>"].Name, "go-github <go-github@github.com>")
		}
		// Nullify Signature since we checked it above
		v.Signature = nil
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"commitSha"}`)
	})

	commit, _, err := client.Git.CreateCommit(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("commitSha")}
	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}
}

func TestGitService_createSignature_nilSigningKey(t *testing.T) {
	a := &createCommit{
		Message: String("Commit Message."),
		Tree:    String("t"),
		Parents: []string{"p"},
	}

	_, err := createSignature(nil, a)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_nilCommit(t *testing.T) {
	_, err := createSignature(&openpgp.Entity{}, nil)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_noAuthor(t *testing.T) {
	a := &createCommit{
		Message: String("Commit Message."),
		Tree:    String("t"),
		Parents: []string{"p"},
	}

	_, err := createSignature(&openpgp.Entity{}, a)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_invalidKey(t *testing.T) {
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	_, err := createSignature(&openpgp.Entity{}, &createCommit{
		Message: String("Commit Message."),
		Tree:    String("t"),
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  String("go-github"),
			Email: String("go-github@github.com"),
			Date:  &date,
		},
	})

	if err == nil {
		t.Errorf("Expected error to be returned due to invalid key")
	}
}

func TestGitService_createSignatureMessage_nilCommit(t *testing.T) {
	_, err := createSignatureMessage(nil)
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_nilMessage(t *testing.T) {
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	_, err := createSignatureMessage(&createCommit{
		Message: nil,
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  String("go-github"),
			Email: String("go-github@github.com"),
			Date:  &date,
		},
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_emptyMessage(t *testing.T) {
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")
	emptyString := ""
	_, err := createSignatureMessage(&createCommit{
		Message: &emptyString,
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  String("go-github"),
			Email: String("go-github@github.com"),
			Date:  &date,
		},
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_nilAuthor(t *testing.T) {
	_, err := createSignatureMessage(&createCommit{
		Message: String("Commit Message."),
		Parents: []string{"p"},
		Author:  nil,
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_withoutTree(t *testing.T) {
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	msg, _ := createSignatureMessage(&createCommit{
		Message: String("Commit Message."),
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  String("go-github"),
			Email: String("go-github@github.com"),
			Date:  &date,
		},
	})
	expected := `parent p
author go-github <go-github@github.com> 1493849023 +0200
committer go-github <go-github@github.com> 1493849023 +0200

Commit Message.`
	if msg != expected {
		t.Errorf("Returned message incorrect. returned %s, want %s", msg, expected)
	}
}

func TestGitService_createSignatureMessage_withoutCommitter(t *testing.T) {
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	msg, _ := createSignatureMessage(&createCommit{
		Message: String("Commit Message."),
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  String("go-github"),
			Email: String("go-github@github.com"),
			Date:  &date,
		},
		Committer: &CommitAuthor{
			Name:  String("foo"),
			Email: String("foo@bar.com"),
			Date:  &date,
		},
	})
	expected := `parent p
author go-github <go-github@github.com> 1493849023 +0200
committer foo <foo@bar.com> 1493849023 +0200

Commit Message.`
	if msg != expected {
		t.Errorf("Returned message incorrect. returned %s, want %s", msg, expected)
	}
}

func TestGitService_CreateCommit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Git.CreateCommit(context.Background(), "%", "%", &Commit{})
	testURLParseError(t, err)
}

const testGPGKey = `
-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=VFSn
-----END PGP PRIVATE KEY BLOCK-----`
