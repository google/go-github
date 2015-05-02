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
)

var (
	client *github.Client

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool
)

// tokenSource is an oauth2.TokenSource which returns a static access token
type tokenSource struct {
	token *oauth2.Token
}

// Token implements the oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func init() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		print("!!! No OAuth token.  Some tests won't run. !!!\n\n")
		client = github.NewClient(nil)
	} else {
		tc := oauth2.NewClient(oauth2.NoContext, &tokenSource{
			&oauth2.Token{AccessToken: token},
		})
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
