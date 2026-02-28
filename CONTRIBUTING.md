# How to contribute

We'd love to accept your patches and contributions to this project. There are
a just a few small guidelines you need to follow.

## Contributor License Agreement

Contributions to any Google project must be accompanied by a Contributor
License Agreement. This is not a copyright **assignment**, it simply gives
Google permission to use and redistribute your contributions as part of the
project. Head over to <https://cla.developers.google.com/> to see your current
agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.

## Reporting issues

Bugs, feature requests, and development-related questions should be directed to
our [GitHub issue tracker](https://github.com/google/go-github/issues).  If
reporting a bug, please try and provide as much context as possible such as
your operating system, Go version, and anything else that might be relevant to
the bug.  For feature requests, please explain what you're trying to do, and
how the requested feature would help you do that.

Security related bugs can either be reported in the issue tracker, or if they
are more sensitive, emailed to <opensource@google.com>.

## Reviewing PRs

In addition to writing code, community projects also require community
contributions in other ways; one of these is reviewing code contributions. If
you are willing to review PRs please open a PR to add your GitHub username to
the [REVIEWERS](./REVIEWERS) file. By adding your GitHub username to the list
of reviewers you are giving contributors permission to request a review for a
PR that has already been approved by a maintainer. If you are asked to review a
PR and either do not have the time or do not think you are able to you should
feel comfortable politely saying no.

If at any time you would like to remove your permission to be contacted for a
review you can open a PR to remove your name from the [REVIEWERS](./REVIEWERS)
file.

## Submitting a patch

1. It's generally best to start by opening a new issue describing the bug or
   feature you're intending to fix. Even if you think it's relatively minor,
   it's helpful to know what people are working on. Mention in the initial issue
   that you are planning to work on that bug or feature so that it can be
   assigned to you.

2. Follow the normal process of [forking][] the project, and set up a new branch
   to work in. It's important that each group of changes be done in separate
   branches in order to ensure that a pull request only includes the commits
   related to that bug or feature.

3. Any significant changes should almost always be accompanied by tests. The
   project already has good test coverage, so look at some of the existing tests
   if you're unsure how to go about it. Coverage is [monitored by codecov.io][],
   which flags pull requests that decrease test coverage. This doesn't
   necessarily mean that PRs with decreased coverage won't be merged. Sometimes
   a decrease in coverage makes sense, but if your PR is flagged, you should
   either add tests to cover those lines or add a PR comment explaining the
   untested lines.

4. Run `script/fmt.sh`, `script/test.sh` and `script/lint.sh` to format your code and
   check that it passes all tests and linters. `script/lint.sh` may also tell you
   that generated files need to be updated. If so, run `script/generate.sh` to
   update them.

5. Do your best to have [well-formed commit messages][] for each change. This
   provides consistency throughout the project, and ensures that commit messages
   are able to be formatted properly by various git tools. See next section for
   more details.

6. Finally, push the commits to your fork and submit a [pull request][].
   **NOTE:** Please do not use force-push on PRs in this repo, as it makes it
   more difficult for reviewers to see what has changed since the last code
   review. We always perform "squash and merge" actions on PRs in this repo, so it doesn't
   matter how many commits your PR has, as they will end up being a single commit after merging.
   This is done to make a much cleaner `git log` history and helps to find regressions in the code
   using existing tools such as `git bisect`.

   - If your PR needs additional reviews you can request one of the
     [REVIEWERS][] takes a look by mentioning them in a PR comment.

[forking]: https://help.github.com/articles/fork-a-repo
[well-formed commit messages]: https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
[pull request]: https://help.github.com/articles/creating-a-pull-request
[monitored by codecov.io]: https://codecov.io/gh/google/go-github
[REVIEWERS]: ./REVIEWERS

### Use proper commit messages and PR titles

Effective Git commit messages and subject lines hold immense significance in comprehending alterations and enhancing the code's maintainability.

Always commit the changes to your fork and push them to the corresponding original repo by sending a Pull Request (PR). Follow the best practices for writing commit messages/PR titles.

1. Limit the subject line to 50 characters
2. Capitalize the subject line
3. Do not end the subject line with a period
4. Use the imperative mood in the subject line. A properly formed Git commit subject line should always be able to complete the following sentence:
   If applied, this commit will `<your subject line here>`

(This above advice can be found all over the internet, but was copied from [here](https://learn-ballerina.github.io/best_practices/use_proper_titles.html).)

5. You may optionally prefix the PR title with the type of PR it is, in lower case,
   followed by a colon. For example, `feat:`, `chore:`, `fix:`, `docs:`, etc.
   For breaking API changes, add an exclamation point.
   For example, `feat!:`, `chore!:`, `fix!:`, etc.

### Windows users

Use Git Bash as a terminal or WSL instead of PowerShell.

To avoid [issues][] with a few linters and formatters within golangci-lint,
make sure you check out files only with LF endings:

```sh
git config core.autocrlf false
git config core.eol lf
```

To convert an existing cloned repo from CRLF to LF, use the following commands:

```sh
git config core.autocrlf false
git rm --cached -r .
git reset --hard HEAD
```

[issues]: https://github.com/golangci/golangci-lint/discussions/5840

## Tips

Although we have not (yet) banned AI-driven contributions to this repo (as many
other open source projects have), we encourage you to read and honor the following
tips (which are frequently ignored by AI-driven PRs):

* Always review your own PRs using the same user interface that your actual
  reviewers will be using with a critical eye, attempting to anticipate what your
  reviewers will call out, _before_ asking anyone to review your PR.
* Come up with a short and appropriate PR title.
* Come up with a short, well-written, and appropriate PR description (we don't
  need hundreds of lines of text here - keep it short and to-the-point so a
  reviewer can _quickly_ determine what the PR is all about).
* If a PR involves bug fixes, it should certainly include a unit test (or tests)
  that demonstrates the bug - without the PR changes, the new unit test would
  fail, but with the included PR changes, the new test(s) pass.
* When possible, try to make smaller, focused PRs (which are easier to review
  and easier for others to understand).

## Code Comments

Every exported method and type needs to have code comments that follow
[Go Doc Comments][]. A typical method's comments will look like this:

```go
// Get fetches a repository.
//
// GitHub API docs: https://docs.github.com/rest/repos/repos#get-a-repository
//
//meta:operation GET /repos/{owner}/{repo}
func (s *RepositoriesService) Get(ctx context.Context, owner, repo string) (*Repository, *Response, error) {
    u := fmt.Sprintf("repos/%v/%v", owner, repo)
    req, err := s.client.NewRequest("GET", u, nil)
    ...
}
```
And the returned type `Repository` will have comments like this:

```go
// Repository represents a GitHub repository.
type Repository struct {
    ID     *int64  `json:"id,omitempty"`
    NodeID *string `json:"node_id,omitempty"`
    Owner  *User   `json:"owner,omitempty"`
    ...
}
```

The first line is the name of the method followed by a short description. This
could also be a longer description if needed, but there is no need to repeat any
details that are documented in GitHub's documentation because users are expected
to follow the documentation links to learn more.

After the description comes a link to the GitHub API documentation. This is
added or fixed automatically when you run `script/generate.sh`, so you won't
need to set this yourself.

Finally, the `//meta:operation` comment is a directive to the code generator
that maps the method to the corresponding OpenAPI operation. Once again, there
can be multiple directives for methods that call multiple
endpoints. `script/generate.sh` will normalize these directives for you, so if
you are adding a new method you can use the pattern from the `u := fmt.Sprintf`
line instead of looking up what the url parameters are called in the OpenAPI
description.

[Go Doc Comments]: https://go.dev/doc/comment

## Metadata

GitHub publishes [OpenAPI descriptions of their API][]. We use these
descriptions to keep documentation links up to date and to keep track of which
methods call which endpoints via the `//meta:operation` comments described
above. GitHub's descriptions are far too large to keep in this repository or to
pull down every time we generate code, so we keep only the metadata we need
in `openapi_operations.yaml`.

### openapi_operations.yaml

Most contributors won't need to interact with `openapi_operations.yaml`, but it
may be useful to know what it is. Its sections are:

- `openapi_operations` - is the metadata that comes from GitHub's OpenAPI
  descriptions. It is generated by `script/metadata.sh update-openapi` and
  should not be edited by hand. In the rare case where it needs to be
  overridden, use the `operation_overrides` section instead.

  An operation consists of `name`, `documentation_url`,
  and `openapi_files`. `openapi_files` is the list of files where the operation
  is described. In order or preference, values can be "api.github.com.json" for
  operations available on the free plan, "ghec.json" for operations available on
  GitHub Enterprise Cloud or "ghes-<version>.json" for operations available on
  GitHub Enterprise Server. When an operation is described in multiple ghes
  files, only the most recent version is included. `documentation_url` is the
  URL that should be linked from godoc. It is the documentation link found in
  the first file listed in `openapi_files`.

- `openapi_commit` - is the git commit that `script/metadata.sh update-openapi`
  saw when it last updated `openapi_operations`. It is not necessarily the most
  recent commit seen because `update-openapi` doesn't update the file when
  there are no changes to `openapi_operations`.

- `operations` - contains manually added metadata that is not in GitHub's
  OpenAPI descriptions. There are only a few of these. Some have
  documentation_urls that point to relevant GitHub documentation that is not in
  the OpenAPI descriptions. Others have no documentation_url and result in a
  note in the generated code that the documentation is missing.

- `operation_overrides` - is where we override the documentation_url for
  operations where the link in the OpenAPI descriptions is wrong.

Please note that if your PR unit tests are failing due to an out-of-date
`openapi_operations.yaml` file, simply ask the maintainer(s) of this repo
to update it for you so that your PR doesn't need to include changes
to this auto-generated file.

### tools/metadata

The `tools/metadata` package is a command-line tool for working with metadata.
In a typical workflow, you won't use it directly, but you will use it indirectly
through `script/generate.sh` and `script/lint.sh`.

Its subcommands are:

- `update-openapi` - updates `openapi_operations.yaml` with the latest
  information from GitHub's OpenAPI descriptions. With `--validate` it will
  validate that the descriptions are correct as of the commit
  in `openapi_commit`. `update-openapi --validate` is called
  by `script/lint.sh`.

- `update-go` - updates Go files with documentation URLs and formats comments.
  It is used by `script/generate.sh`.

- `format` - formats white space in `openapi_operations.yaml` and sorts its
  arrays. It is used by `script/fmt.sh`.

- `unused` - lists operations from `openapi_operations.yaml` that are not mapped
  from any methods.

[OpenAPI descriptions of their API]: https://github.com/github/rest-api-description

## Scripts

The `script` directory has shell scripts that help with common development
tasks.

**script/fmt.sh** formats all Go code in the repository.

**script/generate.sh** runs code generators and `go mod tidy` on all modules. With
`--check` it checks that the generated files are current.

**script/lint.sh** runs linters on the project and checks generated files are
current.

**script/metadata.sh** runs `tools/metadata`. See the [Metadata](#metadata)
section for more information.

**script/test.sh** runs tests on all modules.

## Other notes on code organization

Currently, everything is defined in the main `github` package, with API methods
broken into separate service objects. These services map directly to how
the [GitHub API documentation][] is organized, so use that as your guide for
where to put new methods.

Code is organized in files also based pretty closely on the GitHub API
documentation, following the format `{service}_{api}.go`. For example, methods
defined at <https://docs.github.com/en/rest/webhooks/repos> live in
[repos_hooks.go][].

[GitHub API documentation]: https://docs.github.com/en/rest
[repos_hooks.go]: https://github.com/google/go-github/blob/master/github/repos_hooks.go

## Maintainer's Guide

(These notes are mostly only for people merging in pull requests.)

**Verify CLAs.** CLAs must be on file for the pull request submitter and commit
author(s). Google's CLA verification system should handle this automatically
and will set commit statuses as appropriate. If there's ever any question about
a pull request, ask [willnorris](https://github.com/willnorris).

**Always try to maintain a clean, linear git history.** With very few
exceptions, running `git log` should not show a bunch of branching and merging.

Never use the GitHub "merge" button, since it always creates a merge commit.
Instead, check out the pull request locally ([these git aliases
help][git-aliases]), then cherry-pick or rebase them onto master. If there are
small cleanup commits, especially as a result of addressing code review
comments, these should almost always be squashed down to a single commit. Don't
bother squashing commits that really deserve to be separate though. If needed,
feel free to amend additional small changes to the code or commit message that
aren't worth going through code review for.

If you made any changes like squashing commits, rebasing onto master, etc, then
GitHub won't recognize that this is the same commit in order to mark the pull
request as "merged". So instead, amend the commit message to include a line
"Fixes #0", referencing the pull request number. This would be in addition to
any other "Fixes" lines for closing related issues. If you forget to do this,
you can also leave a comment on the pull request [like this][rebase-comment].
If you made any other changes, it's worth noting that as well, [like
this][modified-comment].

[git-aliases]: https://github.com/willnorris/dotfiles/blob/d640d010c23b1116bdb3d4dc12088ed26120d87d/git/.gitconfig#L13-L15
[rebase-comment]: https://github.com/google/go-github/pull/277#issuecomment-183035491
[modified-comment]: https://github.com/google/go-github/pull/280#issuecomment-184859046

**When creating a release, don't forget to update the `Version` constant in `github.go`.** This is used to
send the version in the `User-Agent` header to identify clients to the GitHub API.
