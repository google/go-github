# go-github #

[![go-github release (latest SemVer)](https://img.shields.io/github/v/release/google/go-github?sort=semver)](https://github.com/google/go-github/releases)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/google/go-github/v63/github)
[![Test Status](https://github.com/google/go-github/workflows/tests/badge.svg)](https://github.com/google/go-github/actions?query=workflow%3Atests)
[![Test Coverage](https://codecov.io/gh/google/go-github/branch/master/graph/badge.svg)](https://codecov.io/gh/google/go-github)
[![Discuss at go-github@googlegroups.com](https://img.shields.io/badge/discuss-go--github%40googlegroups.com-blue.svg)](https://groups.google.com/group/go-github)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/796/badge)](https://bestpractices.coreinfrastructure.org/projects/796)

go-github is a Go client library for accessing the [GitHub API v3][].

Currently, **go-github requires Go version 1.13 or greater**.  go-github tracks
[Go's version support policy][support-policy].  We do our best not to break
older versions of Go if we don't have to, but due to tooling constraints, we
don't always test older versions.

[support-policy]: https://golang.org/doc/devel/release.html#policy

If you're interested in using the [GraphQL API v4][], the recommended library is
[shurcooL/githubv4][].

## Installation ##

go-github is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/google/go-github/v63
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/google/go-github/v63/github"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/google/go-github/v63@master
```

## Usage ##

```go
import "github.com/google/go-github/v63/github"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
import "github.com/google/go-github/github" // with go modules disabled
```

Construct a new GitHub client, then use the various services on the client to
access different parts of the GitHub API. For example:

```go
client := github.NewClient(nil)

// list all organizations for user "willnorris"
orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
```

Some API methods have optional parameters that can be passed. For example:

```go
client := github.NewClient(nil)

// list public repositories for org "github"
opt := &github.RepositoryListByOrgOptions{Type: "public"}
repos, _, err := client.Repositories.ListByOrg(context.Background(), "github", opt)
```

The services of a client divide the API into logical chunks and correspond to
the structure of the GitHub API documentation at
https://docs.github.com/en/rest .

NOTE: Using the [context](https://godoc.org/context) package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

For more sample code snippets, head over to the
[example](https://github.com/google/go-github/tree/master/example) directory.

### Authentication ###

Use the `WithAuthToken` method to configure your client to authenticate using an
OAuth token (for example, a [personal access token][]). This is what is needed
for a majority of use cases aside from GitHub Apps.

```go
client := github.NewClient(nil).WithAuthToken("... your access token ...")
```

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

For API methods that require HTTP Basic Authentication, use the
[`BasicAuthTransport`](https://godoc.org/github.com/google/go-github/github#BasicAuthTransport).

#### As a GitHub App ####

GitHub Apps authentication can be provided by different pkgs like [ghinstallation](https://github.com/bradleyfalzon/ghinstallation) or [go-githubauth](https://github.com/jferrl/go-githubauth).

> **Note**: Most endpoints (ex. [`GET /rate_limit`]) require access token authentication
> while a few others (ex. [`GET /app/hook/deliveries`]) require [JWT] authentication.

[`GET /rate_limit`]: https://docs.github.com/en/rest/rate-limit#get-rate-limit-status-for-the-authenticated-user
[`GET /app/hook/deliveries`]: https://docs.github.com/en/rest/apps/webhooks#list-deliveries-for-an-app-webhook
[JWT]: https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-a-github-app

`ghinstallation` provides `Transport`, which implements `http.RoundTripper` to provide authentication as an installation for GitHub Apps.

Here is an example of how to authenticate as a GitHub App using the `ghinstallation` package:

```go
import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v63/github"
)

func main() {
	// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 1, 99, "2016-10-19.private-key.pem")

	// Or for endpoints that require JWT authentication
	// itr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 1, "2016-10-19.private-key.pem")

	if err != nil {
		// Handle error.
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})

	// Use client...
}
```

`go-githubauth` implements a set of `oauth2.TokenSource` to be used with `oauth2.Client`. An `oauth2.Client` can be injected into the `github.Client` to authenticate requests.

Other example using `go-githubauth`:

```go
package main

import (
 "context"
 "fmt"
 "os"
 "strconv"

 "github.com/google/go-github/v63/github"
 "github.com/jferrl/go-githubauth"
 "golang.org/x/oauth2"
)

func main() {
 privateKey := []byte(os.Getenv("GITHUB_APP_PRIVATE_KEY"))

 appTokenSource, err := githubauth.NewApplicationTokenSource(1112, privateKey)
 if err != nil {
  fmt.Println("Error creating application token source:", err)
  return
 }

 installationTokenSource := githubauth.NewInstallationTokenSource(1113, appTokenSource)

 // oauth2.NewClient uses oauth2.ReuseTokenSource to reuse the token until it expires.
 // The token will be automatically refreshed when it expires.
 // InstallationTokenSource has the mechanism to refresh the token when it expires.
 httpClient := oauth2.NewClient(context.Background(), installationTokenSource)

 client := github.NewClient(httpClient)
}
```

*Note*: In order to interact with certain APIs, for example writing a file to a repo, one must generate an installation token
using the installation ID of the GitHub app and authenticate with the OAuth method mentioned above. See the examples.

### Rate Limiting ###

GitHub imposes a rate limit on all API clients. Unauthenticated clients are
limited to 60 requests per hour, while authenticated clients can make up to
5,000 requests per hour. The Search API has a custom rate limit. Unauthenticated
clients are limited to 10 requests per minute, while authenticated clients
can make up to 30 requests per minute. To receive the higher rate limit when
making calls that are not issued on behalf of a user,
use `UnauthenticatedRateLimitedTransport`.

The returned `Response.Rate` value contains the rate limit information
from the most recent API call. If a recent enough response isn't
available, you can use `RateLimits` to fetch the most up-to-date rate
limit data for the client.

To detect an API rate limit error, you can check if its type is `*github.RateLimitError`:

```go
repos, _, err := client.Repositories.List(ctx, "", nil)
if _, ok := err.(*github.RateLimitError); ok {
	log.Println("hit rate limit")
}
```

Learn more about GitHub rate limiting at
https://docs.github.com/en/rest/rate-limit .

In addition to these rate limits, GitHub imposes a secondary rate limit on all API clients.
This rate limit prevents clients from making too many concurrent requests.

To detect an API secondary rate limit error, you can check if its type is `*github.AbuseRateLimitError`:

```go
repos, _, err := client.Repositories.List(ctx, "", nil)
if _, ok := err.(*github.AbuseRateLimitError); ok {
	log.Println("hit secondary rate limit")
}
```

Alternatively, you can block until the rate limit is reset by using the `context.WithValue` method:

```go
repos, _, err := client.Repositories.List(context.WithValue(ctx, github.SleepUntilPrimaryRateLimitResetWhenRateLimited, true), "", nil)
```

You can use [go-github-ratelimit](https://github.com/gofri/go-github-ratelimit) to handle
secondary rate limit sleep-and-retry for you.

Learn more about GitHub secondary rate limiting at
https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28#about-secondary-rate-limits .

### Accepted Status ###

Some endpoints may return a 202 Accepted status code, meaning that the
information required is not yet ready and was scheduled to be gathered on
the GitHub side. Methods known to behave like this are documented specifying
this behavior.

To detect this condition of error, you can check if its type is
`*github.AcceptedError`:

```go
stats, _, err := client.Repositories.ListContributorsStats(ctx, org, repo)
if _, ok := err.(*github.AcceptedError); ok {
	log.Println("scheduled on GitHub side")
}
```

### Conditional Requests ###

The GitHub API has good support for conditional requests which will help
prevent you from burning through your rate limit, as well as help speed up your
application. `go-github` does not handle conditional requests directly, but is
instead designed to work with a caching `http.Transport`. We recommend using
https://github.com/gregjones/httpcache for that. For example:

```go
import "github.com/gregjones/httpcache"

	client := github.NewClient(
		httpcache.NewMemoryCacheTransport().Client()
    ).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
```

Learn more about GitHub conditional requests at
https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api?apiVersion=2022-11-28#use-conditional-requests-if-appropriate

### Creating and Updating Resources ###

All structs for GitHub resources use pointer values for all non-repeated fields.
This allows distinguishing between unset fields and those set to a zero-value.
Helper functions have been provided to easily create these pointers for string,
bool, and int values. For example:

```go
// create a new private repository named "foo"
repo := &github.Repository{
	Name:    github.String("foo"),
	Private: github.Bool(true),
}
client.Repositories.Create(ctx, "", repo)
```

Users who have worked with protocol buffers should find this pattern familiar.

### Pagination ###

All requests for resource collections (repos, pull requests, issues, etc.)
support pagination. Pagination options are described in the
`github.ListOptions` struct and passed to the list methods directly or as an
embedded type of a more specific list options struct (for example
`github.PullRequestListOptions`). Pages information is available via the
`github.Response` struct.

```go
client := github.NewClient(nil)

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
```

### Webhooks ###

`go-github` provides structs for almost all [GitHub webhook events][] as well as functions to validate them and unmarshal JSON payloads from `http.Request` structs.

```go
func (s *GitHubEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, s.webhookSecretKey)
	if err != nil { ... }
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil { ... }
	switch event := event.(type) {
	case *github.CommitCommentEvent:
		processCommitCommentEvent(event)
	case *github.CreateEvent:
		processCreateEvent(event)
	...
	}
}
```

Furthermore, there are libraries like [cbrgm/githubevents][] that build upon the example above and provide functions to subscribe callbacks to specific events.

For complete usage of go-github, see the full [package docs][].

[GitHub API v3]: https://docs.github.com/en/rest
[personal access token]: https://github.com/blog/1509-personal-api-tokens
[package docs]: https://pkg.go.dev/github.com/google/go-github/v63/github
[GraphQL API v4]: https://developer.github.com/v4/
[shurcooL/githubv4]: https://github.com/shurcooL/githubv4
[GitHub webhook events]: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads
[cbrgm/githubevents]: https://github.com/cbrgm/githubevents

### Testing code that uses `go-github`

The repo [migueleliasweb/go-github-mock](https://github.com/migueleliasweb/go-github-mock) provides a way to mock responses. Check the repo for more details.

### Integration Tests ###

You can run integration tests from the `test` directory. See the integration tests [README](test/README.md).

## Contributing ##
I would like to cover the entire GitHub API and contributions are of course always welcome. The
calling pattern is pretty well established, so adding new methods is relatively
straightforward. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## Versioning ##

In general, go-github follows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. For self-contained libraries, the
application of semantic versioning is relatively straightforward and generally
understood. But because go-github is a client library for the GitHub API, which
itself changes behavior, and because we are typically pretty aggressive about
implementing preview features of the GitHub API, we've adopted the following
versioning policy:

* We increment the **major version** with any incompatible change to
	non-preview functionality, including changes to the exported Go API surface
	or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to
	functionality, as well as any changes to preview functionality in the GitHub
	API. GitHub makes no guarantee about the stability of preview functionality,
	so neither do we consider it a stable part of the go-github API.
* We increment the **patch version** with any backwards-compatible bug fixes.

Preview functionality may take the form of entire methods or simply additional
data returned from an otherwise non-preview method. Refer to the GitHub API
documentation for details on preview functionality.

### Calendar Versioning ###

As of 2022-11-28, GitHub [has announced](https://github.blog/2022-11-28-to-infinity-and-beyond-enabling-the-future-of-githubs-rest-api-with-api-versioning/)
that they are starting to version their v3 API based on "calendar-versioning".

In practice, our goal is to make per-method version overrides (at
least in the core library) rare and temporary.

Our understanding of the GitHub docs is that they will be revving the
entire API to each new date-based version, even if only a few methods
have breaking changes. Other methods will accept the new version with
their existing functionality. So when a new date-based version of the
GitHub API is released, we (the repo maintainers) plan to:

* update each method that had breaking changes, overriding their
  per-method API version header. This may happen in one or multiple
  commits and PRs, and is all done in the main branch.

* once all of the methods with breaking changes have been updated,
  have a final commit that bumps the default API version, and remove
  all of the per-method overrides. That would now get a major version
  bump when the next go-github release is made.

### Version Compatibility Table ###

The following table identifies which version of the GitHub API is
supported by this (and past) versions of this repo (go-github).
Versions prior to 48.2.0 are not listed.

| go-github Version | GitHub v3 API Version |
| ----------------- | --------------------- |
| 63.0.0            | 2022-11-28            |
| 62.0.0            | 2022-11-28            |
| 61.0.0            | 2022-11-28            |
| 60.0.0            | 2022-11-28            |
| 59.0.0            | 2022-11-28            |
| 58.0.0            | 2022-11-28            |
| 57.0.0            | 2022-11-28            |
| 56.0.0            | 2022-11-28            |
| 55.0.0            | 2022-11-28            |
| 54.0.0            | 2022-11-28            |
| 53.2.0            | 2022-11-28            |
| 53.1.0            | 2022-11-28            |
| 53.0.0            | 2022-11-28            |
| 52.0.0            | 2022-11-28            |
| 51.0.0            | 2022-11-28            |
| 50.2.0            | 2022-11-28            |
| 50.1.0            | 2022-11-28            |
| 50.0.0            | 2022-11-28            |
| 49.1.0            | 2022-11-28            |
| 49.0.0            | 2022-11-28            |
| 48.2.0            | 2022-11-28            |

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.
