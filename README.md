# go-github #

go-github is Go library for accessing the [GitHub API][].

**Documentation:** <http://godoc.org/github.com/google/go-github/github>  
**Build Status:** [![Build Status](https://travis-ci.org/google/go-github.png?branch=master)](https://travis-ci.org/google/go-github)

go-github is tested using Go version 1.1.  There are have been
[reports][issue-9] of JSON marshalling errors using Go 1.0.3, at least on
Windows 7.

[issue-9]: https://github.com/google/go-github/issues/9

## Usage ##

```go
import "github.com/google/go-github/github"
```

Construct a new GitHub client, then use the various services on the client to
access different parts of the GitHub API.  For example, to list all
organizations for user "willnorris":

```go
client := github.NewClient(nil)
orgs, err := client.Organizations.List("willnorris", nil)
```

Some API methods have optional parameters that can be passed.  For example,
to list repositories for org "github", sorted by the time they were last
updated:

```go
client := github.NewClient(nil)
opt := &github.RepositoryListByOrgOptions{Sort: "updated"}
repos, err := client.Repositories.ListByOrg("github", opt)
```

The go-github library does not directly handle authentication.  Instead, when
creating a new client, pass an `http.Client` that can handle authentication for
you.  The easiest and recommended way to do this is using the [goauth2][]
library, but you can always use any other library that provides an
`http.Client`.  If you have an OAuth2 access token (for example, a [personal
API token][]), you can use it with the goauth2 library like the following:

```go
t := &oauth.Transport{
  Token: &oauth.Token{AccessToken: "... your access token ..."},
}

client := github.NewClient(t.Client())

// list all repositories for the authenticated user
repos, err := client.Repositories.List("", nil)
```

See the [goauth2 docs][] for complete instructions on using that library.

Also note that when using an authenticated Client, all calls made by the client
will include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

[GitHub API]: http://developer.github.com/v3/
[goauth2]: https://code.google.com/p/goauth2/
[goauth2 docs]: http://godoc.org/code.google.com/p/goauth2/oauth
[personal API token]: https://github.com/blog/1509-personal-api-tokens


## Roadmap ##

This library is being initially developed for an internal application at
Google, so API methods will likely be implemented in the order that they are
needed by that application.  You can track the status of implementation in
[this Google spreadsheet][].  Eventually, I would like to cover the entire
GitHub API, so contributions are of course [always welcome][].  The calling
pattern is pretty well established, so adding new methods is relatively
straightforward.

[this Google spreadsheet]: https://docs.google.com/spreadsheet/ccc?key=0ApoVX4GOiXr-dGNKN1pObFh6ek1DR2FKUjBNZ1FmaEE&usp=sharing
[always welcome]: CONTRIBUTING.md


## License ##

This library is distributed under the BSD-style license found in the LICENSE
file.
