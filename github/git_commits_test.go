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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mockSigner(t *testing.T, signature string, emitErr error, wantMessage string) MessageSignerFunc {
	return func(w io.Writer, r io.Reader) error {
		t.Helper()
		message, err := io.ReadAll(r)
		assertNilError(t, err)
		if wantMessage != "" && string(message) != wantMessage {
			t.Errorf("MessageSignerFunc got %q, want %q", string(message), wantMessage)
		}
		assertWrite(t, w, []byte(signature))
		return emitErr
	}
}

func uncalledSigner(t *testing.T) MessageSignerFunc {
	return func(w io.Writer, r io.Reader) error {
		t.Error("MessageSignerFunc should not be called")
		return nil
	}
}

func TestCommit_Marshal(t *testing.T) {
	testJSONMarshal(t, &Commit{}, "{}")

	u := &Commit{
		SHA: String("s"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("n"),
			Email: String("e"),
			Login: String("u"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
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

	ctx := context.Background()
	commit, _, err := client.Git.GetCommit(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("s"), Message: String("Commit Message."), Author: &CommitAuthor{Name: String("n")}}
	if !cmp.Equal(commit, want) {
		t.Errorf("Git.GetCommit returned %+v, want %+v", commit, want)
	}

	const methodName = "GetCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetCommit(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetCommit(ctx, "o", "r", "s")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_GetCommit_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Git.GetCommit(ctx, "%", "%", "%")
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
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")

		want := &createCommit{
			Message: input.Message,
			Tree:    String("t"),
			Parents: []string{"p"},
		}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"s"}`)
	})

	ctx := context.Background()
	commit, _, err := client.Git.CreateCommit(ctx, "o", "r", input, nil)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("s")}
	if !cmp.Equal(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}

	const methodName = "CreateCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateCommit(ctx, "\n", "\n", input, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateCommit(ctx, "o", "r", input, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
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
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")

		want := &createCommit{
			Message:   input.Message,
			Tree:      String("t"),
			Parents:   []string{"p"},
			Signature: String(signature),
		}
		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"commitSha"}`)
	})

	ctx := context.Background()
	commit, _, err := client.Git.CreateCommit(ctx, "o", "r", input, nil)
	if err != nil {
		t.Errorf("Git.CreateCommit returned error: %v", err)
	}

	want := &Commit{SHA: String("commitSha")}
	if !cmp.Equal(commit, want) {
		t.Errorf("Git.CreateCommit returned %+v, want %+v", commit, want)
	}

	const methodName = "CreateCommit"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateCommit(ctx, "\n", "\n", input, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateCommit(ctx, "o", "r", input, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateSignedCommitWithInvalidParams(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	input := &Commit{}

	ctx := context.Background()
	opts := CreateCommitOptions{Signer: uncalledSigner(t)}
	_, _, err := client.Git.CreateCommit(ctx, "o", "r", input, &opts)
	if err == nil {
		t.Errorf("Expected error to be returned because invalid params were passed")
	}
}

func TestGitService_CreateCommitWithNilCommit(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Git.CreateCommit(ctx, "o", "r", nil, nil)
	if err == nil {
		t.Errorf("Expected error to be returned because commit=nil")
	}
}

func TestGitService_CreateCommit_WithSigner(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	signature := "my voice is my password"
	date := time.Date(2017, time.May, 4, 0, 3, 43, 0, time.FixedZone("CEST", 2*3600))
	author := CommitAuthor{
		Name:  String("go-github"),
		Email: String("go-github@github.com"),
		Date:  &Timestamp{date},
	}
	wantMessage := `tree t
parent p
author go-github <go-github@github.com> 1493849023 +0200
committer go-github <go-github@github.com> 1493849023 +0200

Commit Message.`
	sha := "commitSha"
	input := &Commit{
		SHA:     &sha,
		Message: String("Commit Message."),
		Tree:    &Tree{SHA: String("t")},
		Parents: []*Commit{{SHA: String("p")}},
		Author:  &author,
	}
	wantBody := createCommit{
		Message:   input.Message,
		Tree:      String("t"),
		Parents:   []string{"p"},
		Author:    &author,
		Signature: &signature,
	}
	var gotBody createCommit
	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		assertNilError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		testMethod(t, r, "POST")
		fmt.Fprintf(w, `{"sha":"%s"}`, sha)
	})
	ctx := context.Background()
	wantCommit := &Commit{SHA: String(sha)}
	opts := CreateCommitOptions{Signer: mockSigner(t, signature, nil, wantMessage)}
	commit, _, err := client.Git.CreateCommit(ctx, "o", "r", input, &opts)
	assertNilError(t, err)
	if cmp.Diff(gotBody, wantBody) != "" {
		t.Errorf("Request body = %+v, want %+v\n%s", gotBody, wantBody, cmp.Diff(gotBody, wantBody))
	}
	if cmp.Diff(commit, wantCommit) != "" {
		t.Errorf("Git.CreateCommit returned %+v, want %+v\n%s", commit, wantCommit, cmp.Diff(commit, wantCommit))
	}
}

func TestGitService_createSignature_nilSigner(t *testing.T) {
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
	_, err := createSignature(uncalledSigner(t), nil)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_signerError(t *testing.T) {
	a := &createCommit{
		Message: String("Commit Message."),
		Tree:    String("t"),
		Parents: []string{"p"},
		Author:  &CommitAuthor{Name: String("go-github")},
	}

	signer := mockSigner(t, "", fmt.Errorf("signer error"), "")
	_, err := createSignature(signer, a)

	if err == nil {
		t.Errorf("Expected error to be returned because signer returned an error")
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
			Date:  &Timestamp{date},
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
			Date:  &Timestamp{date},
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
			Date:  &Timestamp{date},
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
			Date:  &Timestamp{date},
		},
		Committer: &CommitAuthor{
			Name:  String("foo"),
			Email: String("foo@bar.com"),
			Date:  &Timestamp{date},
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

	ctx := context.Background()
	_, _, err := client.Git.CreateCommit(ctx, "%", "%", &Commit{}, nil)
	testURLParseError(t, err)
}

func TestSignatureVerification_Marshal(t *testing.T) {
	testJSONMarshal(t, &SignatureVerification{}, "{}")

	u := &SignatureVerification{
		Verified:  Bool(true),
		Reason:    String("reason"),
		Signature: String("sign"),
		Payload:   String("payload"),
	}

	want := `{
		"verified": true,
		"reason": "reason",
		"signature": "sign",
		"payload": "payload"
	}`

	testJSONMarshal(t, u, want)
}

func TestCommitAuthor_Marshal(t *testing.T) {
	testJSONMarshal(t, &CommitAuthor{}, "{}")

	u := &CommitAuthor{
		Date:  &Timestamp{referenceTime},
		Name:  String("name"),
		Email: String("email"),
		Login: String("login"),
	}

	want := `{
		"date": ` + referenceTimeStr + `,
		"name": "name",
		"email": "email",
		"username": "login"
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateCommit_Marshal(t *testing.T) {
	testJSONMarshal(t, &createCommit{}, "{}")

	u := &createCommit{
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  String("name"),
			Email: String("email"),
			Login: String("login"),
		},
		Message:   String("message"),
		Tree:      String("tree"),
		Parents:   []string{"p"},
		Signature: String("sign"),
	}

	want := `{
		"author": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		},
		"committer": {
			"date": ` + referenceTimeStr + `,
			"name": "name",
			"email": "email",
			"username": "login"
		},
		"message": "message",
		"tree": "tree",
		"parents": [
			"p"
		],
		"signature": "sign"
	}`

	testJSONMarshal(t, u, want)
}
