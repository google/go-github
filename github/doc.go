// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github provides a client for using the GitHub API.

Access different parts of the GitHub API using the various services on a GitHub
Client:

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List("willnorris", nil)

Set optional parameters for an API method by passing an Options object.

	// list recently updated repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Sort: "updated"}
	repos, _, err := client.Repositories.ListByOrg("github", opt)

Make authenticated API calls by constructing a GitHub client using an OAuth
capable http.Client:

	import "code.google.com/p/goauth2/oauth"

	// simple OAuth transport if you already have an access token;
	// see goauth2 library for full usage
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "..."},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(nil)

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

The full GitHub API is documented at http://developer.github.com/v3/.
*/
package github
