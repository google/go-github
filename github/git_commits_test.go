// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"errors"
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
	t.Parallel()
	testJSONMarshal(t, &Commit{}, "{}")

	u := &Commit{
		SHA: Ptr("s"),
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("n"),
			Email: Ptr("e"),
			Login: Ptr("u"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("n"),
			Email: Ptr("e"),
			Login: Ptr("u"),
		},
		Message: Ptr("m"),
		Tree: &Tree{
			SHA: Ptr("s"),
			Entries: []*TreeEntry{{
				SHA:     Ptr("s"),
				Path:    Ptr("p"),
				Mode:    Ptr("m"),
				Type:    Ptr("t"),
				Size:    Ptr(1),
				Content: Ptr("c"),
				URL:     Ptr("u"),
			}},
			Truncated: Ptr(false),
		},
		Parents: nil,
		HTMLURL: Ptr("h"),
		URL:     Ptr("u"),
		Verification: &SignatureVerification{
			Verified:  Ptr(false),
			Reason:    Ptr("r"),
			Signature: Ptr("s"),
			Payload:   Ptr("p"),
		},
		NodeID:       Ptr("n"),
		CommentCount: Ptr(1),
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/git/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"sha":"s","message":"Commit Message.","author":{"name":"n"}}`)
	})

	ctx := context.Background()
	commit, _, err := client.Git.GetCommit(ctx, "o", "r", "s")
	if err != nil {
		t.Errorf("Git.GetCommit returned error: %v", err)
	}

	want := &Commit{SHA: Ptr("s"), Message: Ptr("Commit Message."), Author: &CommitAuthor{Name: Ptr("n")}}
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
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Git.GetCommit(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestGitService_CreateCommit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Commit{
		Message: Ptr("Commit Message."),
		Tree:    &Tree{SHA: Ptr("t")},
		Parents: []*Commit{{SHA: Ptr("p")}},
	}

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")

		want := &createCommit{
			Message: input.Message,
			Tree:    Ptr("t"),
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

	want := &Commit{SHA: Ptr("s")}
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
	t.Parallel()
	client, mux, _ := setup(t)

	signature := "----- BEGIN PGP SIGNATURE -----\n\naaaa\naaaa\n----- END PGP SIGNATURE -----"

	input := &Commit{
		Message: Ptr("Commit Message."),
		Tree:    &Tree{SHA: Ptr("t")},
		Parents: []*Commit{{SHA: Ptr("p")}},
		Verification: &SignatureVerification{
			Signature: Ptr(signature),
		},
	}

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")

		want := &createCommit{
			Message:   input.Message,
			Tree:      Ptr("t"),
			Parents:   []string{"p"},
			Signature: Ptr(signature),
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

	want := &Commit{SHA: Ptr("commitSha")}
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
	t.Parallel()
	client, _, _ := setup(t)

	input := &Commit{}

	ctx := context.Background()
	opts := CreateCommitOptions{Signer: uncalledSigner(t)}
	_, _, err := client.Git.CreateCommit(ctx, "o", "r", input, &opts)
	if err == nil {
		t.Errorf("Expected error to be returned because invalid params were passed")
	}
}

func TestGitService_CreateCommitWithNilCommit(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Git.CreateCommit(ctx, "o", "r", nil, nil)
	if err == nil {
		t.Errorf("Expected error to be returned because commit=nil")
	}
}

func TestGitService_CreateCommit_WithSigner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	signature := "my voice is my password"
	date := time.Date(2017, time.May, 4, 0, 3, 43, 0, time.FixedZone("CEST", 2*3600))
	author := CommitAuthor{
		Name:  Ptr("go-github"),
		Email: Ptr("go-github@github.com"),
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
		Message: Ptr("Commit Message."),
		Tree:    &Tree{SHA: Ptr("t")},
		Parents: []*Commit{{SHA: Ptr("p")}},
		Author:  &author,
	}
	wantBody := createCommit{
		Message:   input.Message,
		Tree:      Ptr("t"),
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
	wantCommit := &Commit{SHA: Ptr(sha)}
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
	t.Parallel()
	a := &createCommit{
		Message: Ptr("Commit Message."),
		Tree:    Ptr("t"),
		Parents: []string{"p"},
	}

	_, err := createSignature(nil, a)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_nilCommit(t *testing.T) {
	t.Parallel()
	_, err := createSignature(uncalledSigner(t), nil)

	if err == nil {
		t.Errorf("Expected error to be returned because no author was passed")
	}
}

func TestGitService_createSignature_signerError(t *testing.T) {
	t.Parallel()
	a := &createCommit{
		Message: Ptr("Commit Message."),
		Tree:    Ptr("t"),
		Parents: []string{"p"},
		Author:  &CommitAuthor{Name: Ptr("go-github")},
	}

	signer := mockSigner(t, "", errors.New("signer error"), "")
	_, err := createSignature(signer, a)

	if err == nil {
		t.Errorf("Expected error to be returned because signer returned an error")
	}
}

func TestGitService_createSignatureMessage_nilCommit(t *testing.T) {
	t.Parallel()
	_, err := createSignatureMessage(nil)
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_nilMessage(t *testing.T) {
	t.Parallel()
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	_, err := createSignatureMessage(&createCommit{
		Message: nil,
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  Ptr("go-github"),
			Email: Ptr("go-github@github.com"),
			Date:  &Timestamp{date},
		},
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_emptyMessage(t *testing.T) {
	t.Parallel()
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")
	emptyString := ""
	_, err := createSignatureMessage(&createCommit{
		Message: &emptyString,
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  Ptr("go-github"),
			Email: Ptr("go-github@github.com"),
			Date:  &Timestamp{date},
		},
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_nilAuthor(t *testing.T) {
	t.Parallel()
	_, err := createSignatureMessage(&createCommit{
		Message: Ptr("Commit Message."),
		Parents: []string{"p"},
		Author:  nil,
	})
	if err == nil {
		t.Errorf("Expected error to be returned due to nil key")
	}
}

func TestGitService_createSignatureMessage_withoutTree(t *testing.T) {
	t.Parallel()
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	msg, _ := createSignatureMessage(&createCommit{
		Message: Ptr("Commit Message."),
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  Ptr("go-github"),
			Email: Ptr("go-github@github.com"),
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
	t.Parallel()
	date, _ := time.Parse("Mon Jan 02 15:04:05 2006 -0700", "Thu May 04 00:03:43 2017 +0200")

	msg, _ := createSignatureMessage(&createCommit{
		Message: Ptr("Commit Message."),
		Parents: []string{"p"},
		Author: &CommitAuthor{
			Name:  Ptr("go-github"),
			Email: Ptr("go-github@github.com"),
			Date:  &Timestamp{date},
		},
		Committer: &CommitAuthor{
			Name:  Ptr("foo"),
			Email: Ptr("foo@bar.com"),
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
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Git.CreateCommit(ctx, "%", "%", &Commit{}, nil)
	testURLParseError(t, err)
}

func TestSignatureVerification_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SignatureVerification{}, "{}")

	u := &SignatureVerification{
		Verified:  Ptr(true),
		Reason:    Ptr("reason"),
		Signature: Ptr("sign"),
		Payload:   Ptr("payload"),
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
	t.Parallel()
	testJSONMarshal(t, &CommitAuthor{}, "{}")

	u := &CommitAuthor{
		Date:  &Timestamp{referenceTime},
		Name:  Ptr("name"),
		Email: Ptr("email"),
		Login: Ptr("login"),
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
	t.Parallel()
	testJSONMarshal(t, &createCommit{}, "{}")

	u := &createCommit{
		Author: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("name"),
			Email: Ptr("email"),
			Login: Ptr("login"),
		},
		Committer: &CommitAuthor{
			Date:  &Timestamp{referenceTime},
			Name:  Ptr("name"),
			Email: Ptr("email"),
			Login: Ptr("login"),
		},
		Message:   Ptr("message"),
		Tree:      Ptr("tree"),
		Parents:   []string{"p"},
		Signature: Ptr("sign"),
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
