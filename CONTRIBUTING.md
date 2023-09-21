# How to contribute #

We'd love to accept your patches and contributions to this project. There are
a just a few small guidelines you need to follow.


## Contributor License Agreement ##

Contributions to any Google project must be accompanied by a Contributor
License Agreement. This is not a copyright **assignment**, it simply gives
Google permission to use and redistribute your contributions as part of the
project. Head over to <https://cla.developers.google.com/> to see your current
agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.


## Reporting issues ##

Bugs, feature requests, and development-related questions should be directed to
our [GitHub issue tracker](https://github.com/google/go-github/issues).  If
reporting a bug, please try and provide as much context as possible such as
your operating system, Go version, and anything else that might be relevant to
the bug.  For feature requests, please explain what you're trying to do, and
how the requested feature would help you do that.

Security related bugs can either be reported in the issue tracker, or if they
are more sensitive, emailed to <opensource@google.com>.

## Submitting a patch ##

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
   are able to be formatted properly by various git tools.

6. Finally, push the commits to your fork and submit a [pull request][].
   **NOTE:** Please do not use force-push on PRs in this repo, as it makes it
   more difficult for reviewers to see what has changed since the last code
   review. We always perform "squash and merge" actions on PRs in this repo, so it doesn't
   matter how many commits your PR has, as they will end up being a single commit after merging.
   This is done to make a much cleaner `git log` history and helps to find regressions in the code
   using existing tools such as `git bisect`.

[forking]: https://help.github.com/articles/fork-a-repo
[well-formed commit messages]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
[pull request]: https://help.github.com/articles/creating-a-pull-request
[monitored by codecov.io]: https://codecov.io/gh/google/go-github

## Scripts ##

The `script` directory has shell scripts that help with common development
tasks.

**script/fmt.sh** formats all go code in the repository.

**script/generate.sh** runs code generators and `go mod tidy` on all modules. With
`--check` it checks that the generated files are current.

**script/lint.sh** runs linters on the project and checks generated files are
current.

**script/test.sh** runs tests on all modules.

## Other notes on code organization ##

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


## Maintainer's Guide ##

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
