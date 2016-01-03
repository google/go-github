// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// These tests call the live GitHub API, and therefore require a little more
// setup to run.  See https://github.com/google/go-github/tree/master/tests/integration
// for more information

package tests

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"math/rand"
)

var (
	client *github.Client

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool
)

func init() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		print("!!! No OAuth token.  Some tests won't run. !!!\n\n")
		client = github.NewClient(nil)
	} else {
		tc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		client = github.NewClient(tc)
		auth = true
	}
}

func checkAuth(name string) bool {
	if !auth {
		fmt.Printf("No auth - skipping portions of %v\n", name)
	}
	return auth
}

func createRandomTestRepository(owner string, autoinit bool) (*github.Repository, error) {
	// create random repo name that does not currently exist
	var repoName string
	for {
		repoName = fmt.Sprintf("test-%d", rand.Int())
		_, resp, err := client.Repositories.Get(owner, repoName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				// found a non-existant repo, perfect
				break
			}

			return nil, err
		}
	}

	// create the repository
	repo, _, err := client.Repositories.Create("", &github.Repository{Name: github.String(repoName), AutoInit:github.Bool(autoinit)})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
