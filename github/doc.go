// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github provides a client for using the GitHub API.

Usage:

	import "github.com/google/go-github/v93/github"

Construct a new GitHub client using [NewClient], then use the various services on the client to
access different parts of the GitHub API. For example:

	client, err := github.NewClient()
	if err != nil {
		// Handle error.
	}

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(ctx, "willnorris", nil)

Some API methods have optional parameters that can be passed. For example:

	client, err := github.NewClient()
	if err != nil {
		// Handle error.
	}

	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(ctx, "github", opt)

The services of a client divide the API into logical chunks and correspond to
the structure of the GitHub API documentation at
https://docs.github.com/rest?apiVersion=2022-11-28.

NOTE: Using the [context] package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then [context.Background]
can be used as a starting point.

For more sample code snippets, head over to the https://github.com/google/go-github/tree/master/example directory.

# Authentication

Use [WithAuthToken] to configure your client to authenticate using an OAuth token
(for example, a personal access token). This is what is needed for a majority of use cases
aside from GitHub Apps.

	client, err := github.NewClient(github.WithAuthToken("... your access token ..."))
	if err != nil {
		// Handle error.
	}

Note that when using an authenticated [Client], all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

For API methods that require HTTP Basic Authentication, use the
[WithBasicAuth] option.

GitHub Apps authentication can be provided by the
[github.com/jferrl/go-githubauth](https://github.com/jferrl/go-githubauth) package.
Most endpoints (for example, `GET /rate_limit`) require access token
authentication, while a few others (for example, `GET /app/hook/deliveries`)
require [JWT authentication](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-a-github-app).

`go-githubauth` provides functions that create `oauth2.TokenSource`
implementations. Configure one for the client with [WithTokenSource].
`NewApplicationTokenSource` accepts either the App's client ID (a `string`,
which GitHub now recommends) or its legacy app ID (an `int64`).

For example, create an application token source and then an installation token
source to authenticate a client for an installation:

	privateKey := []byte(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	clientID := os.Getenv("GITHUB_APP_CLIENT_ID")
	installationID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_INSTALLATION_ID"), 10, 64)
	if err != nil {
		// Handle error.
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(clientID, privateKey)
	if err != nil {
		// Handle error.
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installationID, appTokenSource)
	client, err := github.NewClient(github.WithTokenSource(installationTokenSource))
	if err != nil {
		// Handle error.
	}

Some APIs, such as writing a file to a repository, require an installation
token. Use the GitHub App's installation ID to generate one, then authenticate
with the OAuth token source as shown above.

# Rate Limiting

GitHub imposes a rate limit on all API clients. Unauthenticated clients are
limited to 60 requests per hour, while authenticated clients can make up to
5,000 requests per hour. The Search API has a custom rate limit. Unauthenticated
clients are limited to 10 requests per minute, while authenticated clients
can make up to 30 requests per minute. To receive the higher rate limit when
making calls that are not issued on behalf of a user, use [github.WithBasicAuth].

The returned [Response].[Rate] value contains the rate limit information
from the most recent API call. If a recent enough response isn't
available, you can use RateLimits to fetch the most up-to-date rate
limit data for the client.

To detect an API rate limit error, you can check if its type is *[RateLimitError].
For secondary rate limits, you can check if its type is *[AbuseRateLimitError]:

	repos, _, err := client.Repositories.List(ctx, "", nil)
	if errors.As(err, new(*github.RateLimitError)) {
		log.Println("hit rate limit")
	}
	if errors.As(err, new(*github.AbuseRateLimitError)) {
		log.Println("hit secondary rate limit")
	}

Learn more about GitHub rate limiting at
https://docs.github.com/rest/rate-limit?apiVersion=2022-11-28.

# Accepted Status

Some endpoints may return a 202 Accepted status code, meaning that the
information required is not yet ready and was scheduled to be gathered on
the GitHub side. Methods known to behave like this are documented specifying
this behavior.

To detect this condition of error, you can check if its type is
*[AcceptedError]:

	stats, _, err := client.Repositories.ListContributorsStats(ctx, org, repo)
	if errors.As(err, new(*github.AcceptedError)) {
		log.Println("scheduled on GitHub side")
	}

# Conditional Requests

The GitHub REST API has good support for conditional HTTP requests
via the ETag header which will help prevent you from burning through your
rate limit, as well as help speed up your application. go-github does not
handle conditional requests directly, but is instead designed to work with a
caching [http.Transport].

Typically, an RFC 9111 compliant HTTP cache such as https://github.com/bartventer/httpcache
is recommended. Alternatively, the https://github.com/bored-engineer/github-conditional-http-transport
package relies on (undocumented) GitHub specific cache logic and is
recommended when making requests using short-lived credentials such as a
GitHub App installation token.

Learn more about GitHub conditional requests at
https://docs.github.com/rest/using-the-rest-api/best-practices-for-using-the-rest-api?apiVersion=2022-11-28#use-conditional-requests-if-appropriate.

# Creating and Updating Resources

Many resource and request structs use pointers for optional scalar fields. This
allows distinguishing between unset fields and those set to a zero-value.
Required fields such as [CreateRepositoryRequest.Name] use value types.
Use the new builtin to easily create these pointers for string,
bool, and int values. For example:

	// create a new private repository named "foo"
	repo := github.CreateRepositoryRequest{
		Name:    "foo",
		Private: new(true),
	}
	client.Repositories.Create(ctx, "", repo)

Users who have worked with protocol buffers should find this pattern familiar.

# Pagination

All requests for resource collections (repos, pull requests, issues, etc.)
support pagination. Pagination options are described in the
[ListOptions] struct and passed to the list methods directly or as an
embedded type of a more specific list options struct (for example
[PullRequestListOptions]). Pages information is available via the
[Response] struct.

	client, err := github.NewClient()
	if err != nil {
		// Handle error.
	}

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, "github", opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
*/
package github
