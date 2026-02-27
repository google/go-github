// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The commitpr command utilizes go-github as a CLI tool for
// pushing files to a branch and creating a pull request from it.
// It takes an auth token as an environment variable and creates
// the commit and the PR under the account affiliated with that token.
//
// The purpose of this example is to show how to use refs, trees and commits to
// create commits and pull requests.
//
// Note, if you want to push a single file, you probably prefer to use the
// content API. An example is available here:
// https://pkg.go.dev/github.com/google/go-github/v84/github#example-RepositoriesService-CreateFile
//
// Note, for this to work at least 1 commit is needed, so you if you use this
// after creating a repository you might want to make sure you set `AutoInit` to
// `true`.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/google/go-github/v84/github"
)

var (
	sourceOwner   = flag.String("source-owner", "", "Name of the owner (user or org) of the repo to create the commit in.")
	sourceRepo    = flag.String("source-repo", "", "Name of repo to create the commit in.")
	commitMessage = flag.String("commit-message", "", "Content of the commit message.")
	commitBranch  = flag.String("commit-branch", "", "Name of branch to create the commit in. If it does not already exists, it will be created using the `base-branch` parameter")
	repoBranch    = flag.String("repo-branch", "", "Name of the repository where the changes in the pull request were made. This field is required for cross-repository pull requests if both repositories are owned by the same organization")
	baseBranch    = flag.String("base-branch", "master", "Name of branch to create the `commit-branch` from.")
	prRepoOwner   = flag.String("merge-repo-owner", "", "Name of the owner (user or org) of the repo to create the PR against. If not specified, the value of the `-source-owner` flag will be used.")
	prRepo        = flag.String("merge-repo", "", "Name of repo to create the PR against. If not specified, the value of the `-source-repo` flag will be used.")
	prBranch      = flag.String("merge-branch", "master", "Name of branch to create the PR against (the one you want to merge your branch in via the PR).")
	prSubject     = flag.String("pr-title", "", "Title of the pull request. If not specified, no pull request will be created.")
	prDescription = flag.String("pr-text", "", "Text to put in the description of the pull request.")
	sourceFiles   = flag.String("files", "", `Comma-separated list of files to commit and their location.
The local file is separated by its target location by a semi-colon.
If the file should be in the same location with the same name, you can just put the file name and omit the repetition.
Example: README.md,main.go:github/examples/commitpr/main.go`)
	authorName  = flag.String("author-name", "", "Name of the author of the commit.")
	authorEmail = flag.String("author-email", "", "Email of the author of the commit.")
	privateKey  = flag.String("private-key", "", "Path to the private key to use to sign the commit.")
)

var (
	client *github.Client
	ctx    = context.Background()
)

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, *sourceOwner, *sourceRepo, branchRef(*commitBranch)); err == nil {
		return ref, nil
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	if *commitBranch == *baseBranch {
		return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if *baseBranch == "" {
		return nil, errors.New("the `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, *sourceOwner, *sourceRepo, branchRef(*baseBranch)); err != nil {
		return nil, err
	}
	newRef := github.CreateRef{Ref: branchRef(*commitBranch), SHA: *baseRef.Object.SHA}
	ref, _, err = client.Git.CreateRef(ctx, *sourceOwner, *sourceRepo, newRef)
	return ref, err
}

// branchRef generates the fully qualified git reference for the given branch name.
func branchRef(name string) string {
	return "refs/heads/" + name
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for fileArg := range strings.SplitSeq(*sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.Ptr(file), Type: github.Ptr("blob"), Content: github.Ptr(string(content)), Mode: github.Ptr("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{Time: date}, Name: authorName, Email: authorEmail}
	commit := github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	opts := github.CreateCommitOptions{}
	if *privateKey != "" {
		armoredBlock, e := os.ReadFile(*privateKey)
		if e != nil {
			return e
		}
		keyring, e := openpgp.ReadArmoredKeyRing(bytes.NewReader(armoredBlock))
		if e != nil {
			return e
		}
		if len(keyring) != 1 {
			return errors.New("expected exactly one key in the keyring")
		}
		key := keyring[0]
		opts.Signer = github.MessageSignerFunc(func(w io.Writer, r io.Reader) error {
			return openpgp.ArmoredDetachSign(w, key, r, nil)
		})
	}
	newCommit, _, err := client.Git.CreateCommit(ctx, *sourceOwner, *sourceRepo, commit, &opts)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, *sourceOwner, *sourceRepo, *ref.Ref, github.UpdateRef{
		SHA:   *newCommit.SHA,
		Force: github.Ptr(false),
	})
	return err
}

// createPR creates a pull request. Based on: https://pkg.go.dev/github.com/google/go-github/v84/github#example-PullRequestsService-Create
func createPR() (err error) {
	if *prSubject == "" {
		return errors.New("missing `-pr-title` flag; skipping PR creation")
	}

	if *prRepoOwner != "" && *prRepoOwner != *sourceOwner {
		*commitBranch = fmt.Sprintf("%v:%v", *sourceOwner, *commitBranch)
	} else {
		prRepoOwner = sourceOwner
	}

	if *prRepo == "" {
		prRepo = sourceRepo
	}

	newPR := &github.NewPullRequest{
		Title:               prSubject,
		Head:                commitBranch,
		HeadRepo:            repoBranch,
		Base:                prBranch,
		Body:                prDescription,
		MaintainerCanModify: github.Ptr(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, *prRepoOwner, *prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %v\n", pr.GetHTMLURL())
	return nil
}

func main() {
	flag.Parse()
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if *sourceOwner == "" || *sourceRepo == "" || *commitBranch == "" || *sourceFiles == "" || *authorName == "" || *authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}
	client = github.NewClient(nil).WithAuthToken(token)

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %v\n", err)
	}
	if ref == nil {
		log.Fatal("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %v\n", err)
	}

	if err := pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %v\n", err)
	}

	if err := createPR(); err != nil {
		log.Fatalf("Error while creating the pull request: %v", err)
	}
}
