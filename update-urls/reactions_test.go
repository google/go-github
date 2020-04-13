// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func newReactionsPipeline() *pipelineSetup {
	return &pipelineSetup{
		baseURL:              "https://developer.github.com/v3/reactions/",
		endpointsFromWebsite: reactionsWant,
		filename:             "reactions.go",
		serviceName:          "ReactionsService",
		originalGoSource:     reactionsGoFileOriginal,
		wantGoSource:         reactionsGoFileWant,
		wantNumEndpoints:     25,
	}
}

func TestPipeline_Reactions(t *testing.T) {
	ps := newReactionsPipeline()
	ps.setup(t, false, false)
	ps.validate(t)
}

func TestPipeline_Reactions_FirstStripAllURLs(t *testing.T) {
	ps := newReactionsPipeline()
	ps.setup(t, true, false)
	ps.validate(t)
}

func TestPipeline_Reactions_FirstDestroyReceivers(t *testing.T) {
	ps := newReactionsPipeline()
	ps.setup(t, false, true)
	ps.validate(t)
}

func TestPipeline_Reactions_FirstStripAllURLsAndDestroyReceivers(t *testing.T) {
	ps := newReactionsPipeline()
	ps.setup(t, true, true)
	ps.validate(t)
}

func TestParseWebPageEndpoints_Reactions(t *testing.T) {
	got, want := parseWebPageEndpoints(reactionsTestWebPage), reactionsWant
	testWebPageHelper(t, got, want)
}

var reactionsWant = endpointsByFragmentID{
	"list-reactions-for-a-commit-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/comments/%v/reactions"}, httpMethod: "GET"}},

	"delete-a-commit-comment-reaction": []*Endpoint{
		{urlFormats: []string{"repositories/%v/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"repos/%v/%v/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"create-reaction-for-an-issue": []*Endpoint{{urlFormats: []string{"repos/%v/%v/issues/%v/reactions"}, httpMethod: "POST"}},

	"delete-an-issue-reaction": []*Endpoint{
		{urlFormats: []string{"repositories/%v/issues/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"repos/%v/%v/issues/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"create-reaction-for-a-pull-request-review-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/pulls/comments/%v/reactions"}, httpMethod: "POST"}},

	"list-reactions-for-a-team-discussion": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/reactions"}, httpMethod: "GET"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/reactions"}, httpMethod: "GET"},
	},

	"delete-a-reaction-legacy": []*Endpoint{{urlFormats: []string{"reactions/%v"}, httpMethod: "DELETE"}},

	"list-reactions-for-a-team-discussion-comment-legacy": []*Endpoint{{urlFormats: []string{"teams/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "GET"}},

	"delete-an-issue-comment-reaction": []*Endpoint{
		{urlFormats: []string{"repositories/%v/issues/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"repos/%v/%v/issues/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"list-reactions-for-a-pull-request-review-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/pulls/comments/%v/reactions"}, httpMethod: "GET"}},

	"create-reaction-for-a-team-discussion-legacy": []*Endpoint{{urlFormats: []string{"teams/%v/discussions/%v/reactions"}, httpMethod: "POST"}},

	"create-reaction-for-a-team-discussion-comment-legacy": []*Endpoint{{urlFormats: []string{"teams/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "POST"}},

	"create-reaction-for-a-commit-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/comments/%v/reactions"}, httpMethod: "POST"}},

	"list-reactions-for-an-issue": []*Endpoint{{urlFormats: []string{"repos/%v/%v/issues/%v/reactions"}, httpMethod: "GET"}},

	"create-reaction-for-an-issue-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/issues/comments/%v/reactions"}, httpMethod: "POST"}},

	"create-reaction-for-a-team-discussion": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/reactions"}, httpMethod: "POST"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/reactions"}, httpMethod: "POST"},
	},

	"delete-team-discussion-reaction": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"create-reaction-for-a-team-discussion-comment": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "POST"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "POST"},
	},

	"list-reactions-for-an-issue-comment": []*Endpoint{{urlFormats: []string{"repos/%v/%v/issues/comments/%v/reactions"}, httpMethod: "GET"}},

	"delete-a-pull-request-comment-reaction": []*Endpoint{
		{urlFormats: []string{"repositories/%v/pulls/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"repos/%v/%v/pulls/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"list-reactions-for-a-team-discussion-comment": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "GET"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/comments/%v/reactions"}, httpMethod: "GET"},
	},

	"delete-team-discussion-comment-reaction": []*Endpoint{
		{urlFormats: []string{"organizations/%v/team/%v/discussions/%v/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
		{urlFormats: []string{"orgs/%v/teams/%v/discussions/%v/comments/%v/reactions/%v"}, httpMethod: "DELETE"},
	},

	"list-reactions-for-a-team-discussion-legacy": []*Endpoint{{urlFormats: []string{"teams/%v/discussions/%v/reactions"}, httpMethod: "GET"}},
}

var reactionsTestWebPage = `<!DOCTYPE html>
<html lang="en" prefix="og: http://ogp.me/ns#">

<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <meta http-equiv="Content-Language" content="en-us" />
  <meta http-equiv="imagetoolbar" content="false" />
  <meta name="MSSmartTagsPreventParsing" content="true" />
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Reactions | GitHub Developer Guide</title>
  <meta property="og:url" content="https://developer.github.com/v3/reactions/" />
  <meta property="og:site_name" content="GitHub Developer" />
  <meta property="og:title" content="Reactions" />
  <meta property="og:description"
    content="Get started with one of our guides, or jump straight into the API documentation." />
  <meta property="og:type" content="website" />
  <meta property="og:author" content="https://www.facebook.com/GitHub" />
  <meta property="og:image" content="https://og.github.com/logo/github-logo@1200x1200.png" />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="1200" />
  <meta property="og:image" content="https://og.github.com/mark/github-mark@1200x630.png" />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta property="og:image" content="https://og.github.com/octocat/github-octocat@1200x630.png" />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta property="twitter:card" content="summary_large_image" />
  <meta property="twitter:site" content="@github" />
  <meta property="twitter:site:id" content="13334762" />
  <meta property="twitter:creator" content="@githubapi" />
  <meta property="twitter:creator:id" content="@539153822" />
  <meta property="twitter:title" content="Reactions" />
  <meta property="twitter:description"
    content="Get started with one of our guides, or jump straight into the API documentation." />
  <meta property="twitter:image:src" content="https://og.github.com/logo/github-logo@1200x1200.png" />
  <meta property="twitter:image:width" content="1200" />
  <meta property="twitter:image:height" content="1200" />


  <link rel="alternate" type="application/atom+xml" title="API Changes" href="/changes.atom" />

  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:400,300,400italic,500">
  <link href="/assets/stylesheets/application.css" rel="stylesheet" type="text/css" />
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
  <script src="/assets/javascripts/application.js" type="text/javascript"></script>

</head>


<body class="api ">


  <header class="site-header">
    <div id="header" class="container">
      <a class="site-header-logo mt-1" href="/"><img src="/assets/images/github-developer-logo.svg"
          alt="GitHub Developer"></a>
      <nav class="site-header-nav" aria-label="Main Navigation">
        <div class="dropdown" id="api-docs-dropdown">
          <div class="dropdown-button" tabIndex="0" aria-haspopup="true" aria-expanded="false"><span
              class="dropdown-button-link" tabIndex="-1">Docs</span>
            <div class="dropdown-caret"></div>
          </div>
          <div id="js-dropdown-menu" class="dropdown-menu">
            <a class="dropdown-menu-item" href="/apps/">Apps</a>
            <a class="dropdown-menu-item" href="/actions/">GitHub Actions</a>
            <a class="dropdown-menu-item" href="/marketplace/">GitHub Marketplace</a>
            <a class="dropdown-menu-item" href="/webhooks/">Webhooks</a>
            <a class="dropdown-menu-item" href="/partnerships/">Partnerships</a>
            <a class="dropdown-menu-item" href="/v3/">REST API v3</a>
            <a class="dropdown-menu-item" href="/v4/">GraphQL API v4</a>
          </div>
        </div>
        <a class="site-header-nav-item" href="/changes/">Blog</a>
        <a class="site-header-nav-item" href="/forum/">Forum</a>
        <div class="dropdown" id="versions-dropdown">
          <div class="dropdown-button" tabIndex="0" aria-haspopup="true" aria-expanded="false"><span
              class="dropdown-button-link" tabIndex="-1">Versions</span>
            <div class="dropdown-caret"></div>
          </div>
          <div id="js-dropdown-menu" class="dropdown-menu">
            <a class="dropdown-menu-item" data-proofer-ignore>GitHub.com</a>

            <a class="dropdown-menu-item" data-proofer-ignore>GitHub Enterprise Server 2.20</a>

            <a class="dropdown-menu-item" data-proofer-ignore>GitHub Enterprise Server 2.19</a>

            <a class="dropdown-menu-item" data-proofer-ignore>GitHub Enterprise Server 2.18</a>

            <a class="dropdown-menu-item" data-proofer-ignore>GitHub Enterprise Server 2.17</a>

          </div>
        </div>
        <form accept-charset="UTF-8" action="#" class="site-header-nav-item site-header-search">
          <label id="search-container" class="mb-0">
            <input type="text" id="searchfield" class="form-control" placeholder="Search…" autocomplete="off"
              autocorrect="off" autocapitalize="off" spellcheck="false" />
            <div class="cancel-search"></div>
            <ul id="search-results"></ul>
          </label>
        </form>
      </nav>
    </div>
  </header>


  <div class="container">

    <div class="sub-nav mt-5 mb-3">
      <h2><a href="/v3/">REST API v3</a></h2>
      <ul>

        <li><a href="/v3/" class="active">Reference</a></li>
        <li><a href="/v3/guides/">Guides</a></li>
        <li><a href="/v3/libraries/">Libraries</a></li>

      </ul>

    </div>

    <div class="nav-select nav-select-top hide-lg hide-xl">
      <select class="form-select mt-4" onchange="if (this.value) window.location.href=this.value">
        <option value="">Navigate the docs…</option>

        <optgroup label="Overview">
          <option value="/v3/">API Overview</a></h3>
          <option value="/v3/media/">Media Types</option>
          <option value="/v3/oauth_authorizations/">OAuth Authorizations API</option>
          <option value="/v3/auth/">Other Authentication Methods</option>
          <option value="/v3/troubleshooting/">Troubleshooting</option>
          <option value="/v3/previews/">API Previews</option>
          <option value="/v3/versions/">Versions</option>
        </optgroup>

        <optgroup label="Activity">
          <option value="/v3/activity/">Activity overview</a></h3>
          <option value="/v3/activity/events/">Events</option>
          <option value="/v3/activity/events/types/">Event Types &amp; Payloads</option>
          <option value="/v3/activity/feeds/">Feeds</option>
          <option value="/v3/activity/notifications/">Notifications</option>
          <option value="/v3/activity/starring/">Starring</option>
          <option value="/v3/activity/watching/">Watching</option>
        </optgroup>

        <optgroup label="Checks">
          <option value="/v3/checks/">Checks</a></h3>
          <option value="/v3/checks/runs/">Check Runs</option>
          <option value="/v3/checks/suites/">Check Suites</option>
        </optgroup>



        <optgroup label="Gists">
          <option value="/v3/gists/">Gists overview</a></h3>
          <option value="/v3/gists/comments/">Comments</option>
        </optgroup>

        <optgroup label="Git Data">
          <option value="/v3/git/">Git Data overview</a></h3>
          <option value="/v3/git/blobs/">Blobs</option>
          <option value="/v3/git/commits/">Commits</option>
          <option value="/v3/git/refs/">References</option>
          <option value="/v3/git/tags/">Tags</option>
          <option value="/v3/git/trees/">Trees</option>
        </optgroup>


        <optgroup label="GitHub Actions">
          <option value="/v3/actions/">GitHub Actions Overview</a></h3>
          <option value="/v3/actions/artifacts/">Artifacts</option>
          <option value="/v3/actions/secrets/">Secrets</option>
          <option value="/v3/actions/self_hosted_runners/">Self-hosted runners</option>
          <option value="/v3/actions/workflows/">Workflows</option>
          <option value="/v3/actions/workflow_jobs/">Workflow jobs</option>
          <option value="/v3/actions/workflow_runs/">Workflow runs</option>
        </optgroup>


        <optgroup label="GitHub Apps">
          <option value="/v3/apps/">GitHub Apps overview</a></h3>

          <option value="/v3/apps/oauth_applications/">OAuth Applications API</option>

          <option value="/v3/apps/installations/">Installations</option>
          <option value="/v3/apps/permissions/">Permissions</option>
          <option value="/v3/apps/available-endpoints/">Available Endpoints</option>
        </optgroup>


        <optgroup label="Marketplace">
          <option value="/v3/apps/marketplace/">GitHub Marketplace</option>
          </h3>
        </optgroup>



        <optgroup label="Interactions">
          <option value="/v3/interactions/">Interactions</a></h3>
          <option value="/v3/interactions/orgs/">Organization</option>
          <option value="/v3/interactions/repos/">Repository</option>
        </optgroup>


        <optgroup label="Issues">
          <option value="/v3/issues/">Issues overview</a></h3>
          <option value="/v3/issues/assignees/">Assignees</option>
          <option value="/v3/issues/comments/">Comments</option>
          <option value="/v3/issues/events/">Events</option>
          <option value="/v3/issues/labels/">Labels</option>
          <option value="/v3/issues/milestones/">Milestones</option>
          <option value="/v3/issues/timeline/">Timeline</option>
        </optgroup>


        <optgroup label="Migrations">
          <option value="/v3/migrations/">Migrations overview</a></h3>
          <option value="/v3/migrations/orgs/">Organization</option>
          <option value="/v3/migrations/source_imports/">Source Imports</option>
          <option value="/v3/migrations/users/">User</option>
        </optgroup>


        <optgroup label="Miscellaneous">
          <option value="/v3/misc/">Miscellaneous overview</a></h3>
          <option value="/v3/codes_of_conduct/">Codes of Conduct</option>
          <option value="/v3/emojis/">Emojis</option>
          <option value="/v3/gitignore/">Gitignore</option>
          <option value="/v3/licenses/">Licenses</option>
          <option value="/v3/markdown/">Markdown</option>
          <option value="/v3/meta/">Meta</option>

          <option value="/v3/rate_limit/">Rate Limit</option>

        </optgroup>

        <optgroup label="Organizations">
          <option value="/v3/orgs/">Organizations overview</a></h3>

          <option value="/v3/orgs/blocking/">Blocking Users &#40;Organizations&#41;</option>

          <option value="/v3/orgs/members/">Members</option>
          <option value="/v3/orgs/outside_collaborators/">Outside Collaborators</option>
          <option value="/v3/orgs/hooks/">Webhooks</option>
        </optgroup>

        <optgroup label="Projects">
          <option value="/v3/projects/">Projects overview</a></h3>
          <option value="/v3/projects/cards/">Cards</option>
          <option value="/v3/projects/collaborators/">Collaborators</option>
          <option value="/v3/projects/columns/">Columns</option>
        </optgroup>

        <optgroup label="Pull Requests">
          <option value="/v3/pulls/">Pull Requests overview</a></h3>
          <option value="/v3/pulls/reviews/">Reviews</option>
          <option value="/v3/pulls/comments/">Review Comments</option>
          <option value="/v3/pulls/review_requests/">Review Requests</option>
        </optgroup>

        <optgroup label="Reactions">
          <option value="/v3/reactions/">Reactions overview</a></h3>
          <option value="/v3/reactions/#list-reactions-for-a-commit-comment">Commit Comment</option>
          <option value="/v3/reactions/#list-reactions-for-an-issue">Issue</option>
          <option value="/v3/reactions/#list-reactions-for-an-issue-comment">Issue Comment</option>
          <option value="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review Comment
          </option>
          <option value="/v3/reactions/#list-reactions-for-a-team-discussion">Team Discussion</option>
          <option value="/v3/reactions/#list-reactions-for-a-team-discussion-comment">Team Discussion Comment</option>
        </optgroup>

        <optgroup label="Repositories">
          <option value="/v3/repos/">Repositories overview</a></h3>
          <option value="/v3/repos/branches/">Branches</option>
          <option value="/v3/repos/collaborators/">Collaborators</option>
          <option value="/v3/repos/comments/">Comments</option>
          <option value="/v3/repos/commits/">Commits</option>

          <option value="/v3/repos/community/">Community</option>

          <option value="/v3/repos/contents/">Contents</option>
          <option value="/v3/repos/keys/">Deploy Keys</option>
          <option value="/v3/repos/deployments/">Deployments</option>
          <option value="/v3/repos/downloads/">Downloads</option>
          <option value="/v3/repos/forks/">Forks</option>
          <option value="/v3/repos/invitations/">Invitations</option>
          <option value="/v3/repos/merging/">Merging</option>
          <option value="/v3/repos/pages/">Pages</option>
          <option value="/v3/repos/releases/">Releases</option>
          <option value="/v3/repos/statistics/">Statistics</option>
          <option value="/v3/repos/statuses/">Statuses</option>

          <option value="/v3/repos/traffic/">Traffic</option>

          <option value="/v3/repos/hooks/">Webhooks</option>
        </optgroup>

        <optgroup label="Search">
          <option value="/v3/search/">Search overview</a></h3>
          <option value="/v3/search/#search-repositories">Repositories</option>
          <option value="/v3/search/#search-code">Code</option>
          <option value="/v3/search/#search-commits">Commits</option>
          <option value="/v3/search/#search-issues-and-pull-requests">Issues</option>
          <option value="/v3/search/#search-users">Users</option>
          <option value="/v3/search/#search-topics">Topics</option>
          <option value="/v3/search/#text-match-metadata">Text match metadata</<option>
          <option value="/v3/search/legacy/">Legacy search</option>
        </optgroup>

        <optgroup label="Teams">
          <option value="/v3/teams/">Teams</a></h3>
          <option value="/v3/teams/discussions/">Discussions</option>
          <option value="/v3/teams/discussion_comments/">Discussion comments</option>
          <option value="/v3/teams/members/">Members</option>
          <option value="/v3/teams/team_sync/">Team sunchronization</option>
        </optgroup>


        <optgroup label="SCIM">
          <option value="/v3/scim/">SCIM</a></h3>
        </optgroup>


        <optgroup label="Users">
          <option value="/v3/users/">Users overview</a></h3>

          <option value="/v3/users/blocking/">Blocking Users</option>

          <option value="/v3/users/emails/">Emails</option>
          <option value="/v3/users/followers/">Followers</option>
          <option value="/v3/users/keys/">Git SSH Keys</option>
          <option value="/v3/users/gpg_keys/">GPG Keys</option>
        </optgroup>
      </select>

    </div>

    <section class="main pt-5 clearfix gut-3" role="main">
      <div class="content col-md-8">
        <h1>
          <a id="reactions" class="anchor" href="#reactions" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Reactions</h1>

        <ul id="markdown-toc">
          <li><a href="#reaction-types" id="markdown-toc-reaction-types">Reaction types</a></li>
          <li><a href="#list-reactions-for-a-commit-comment" id="markdown-toc-list-reactions-for-a-commit-comment">List
              reactions for a commit comment</a></li>
          <li><a href="#create-reaction-for-a-commit-comment"
              id="markdown-toc-create-reaction-for-a-commit-comment">Create reaction for a commit comment</a></li>
          <li><a href="#delete-a-commit-comment-reaction" id="markdown-toc-delete-a-commit-comment-reaction">Delete a
              commit comment reaction</a></li>
          <li><a href="#list-reactions-for-an-issue" id="markdown-toc-list-reactions-for-an-issue">List reactions for an
              issue</a></li>
          <li><a href="#create-reaction-for-an-issue" id="markdown-toc-create-reaction-for-an-issue">Create reaction for
              an issue</a></li>
          <li><a href="#delete-an-issue-reaction" id="markdown-toc-delete-an-issue-reaction">Delete an issue
              reaction</a></li>
          <li><a href="#list-reactions-for-an-issue-comment" id="markdown-toc-list-reactions-for-an-issue-comment">List
              reactions for an issue comment</a></li>
          <li><a href="#create-reaction-for-an-issue-comment"
              id="markdown-toc-create-reaction-for-an-issue-comment">Create reaction for an issue comment</a></li>
          <li><a href="#delete-an-issue-comment-reaction" id="markdown-toc-delete-an-issue-comment-reaction">Delete an
              issue comment reaction</a></li>
          <li><a href="#list-reactions-for-a-pull-request-review-comment"
              id="markdown-toc-list-reactions-for-a-pull-request-review-comment">List reactions for a pull request
              review comment</a></li>
          <li><a href="#create-reaction-for-a-pull-request-review-comment"
              id="markdown-toc-create-reaction-for-a-pull-request-review-comment">Create reaction for a pull request
              review comment</a></li>
          <li><a href="#delete-a-pull-request-comment-reaction"
              id="markdown-toc-delete-a-pull-request-comment-reaction">Delete a pull request comment reaction</a></li>
          <li><a href="#list-reactions-for-a-team-discussion"
              id="markdown-toc-list-reactions-for-a-team-discussion">List reactions for a team discussion</a></li>
          <li><a href="#create-reaction-for-a-team-discussion"
              id="markdown-toc-create-reaction-for-a-team-discussion">Create reaction for a team discussion</a></li>
          <li><a href="#delete-team-discussion-reaction" id="markdown-toc-delete-team-discussion-reaction">Delete team
              discussion reaction</a></li>
          <li><a href="#list-reactions-for-a-team-discussion-comment"
              id="markdown-toc-list-reactions-for-a-team-discussion-comment">List reactions for a team discussion
              comment</a></li>
          <li><a href="#create-reaction-for-a-team-discussion-comment"
              id="markdown-toc-create-reaction-for-a-team-discussion-comment">Create reaction for a team discussion
              comment</a></li>
          <li><a href="#delete-team-discussion-comment-reaction"
              id="markdown-toc-delete-team-discussion-comment-reaction">Delete team discussion comment reaction</a></li>
          <li><a href="#delete-a-reaction-legacy" id="markdown-toc-delete-a-reaction-legacy">Delete a reaction
              (Legacy)</a></li>
          <li><a href="#list-reactions-for-a-team-discussion-legacy"
              id="markdown-toc-list-reactions-for-a-team-discussion-legacy">List reactions for a team discussion
              (Legacy)</a></li>
          <li><a href="#create-reaction-for-a-team-discussion-legacy"
              id="markdown-toc-create-reaction-for-a-team-discussion-legacy">Create reaction for a team discussion
              (Legacy)</a></li>
          <li><a href="#list-reactions-for-a-team-discussion-comment-legacy"
              id="markdown-toc-list-reactions-for-a-team-discussion-comment-legacy">List reactions for a team discussion
              comment (Legacy)</a></li>
          <li><a href="#create-reaction-for-a-team-discussion-comment-legacy"
              id="markdown-toc-create-reaction-for-a-team-discussion-comment-legacy">Create reaction for a team
              discussion comment (Legacy)</a></li>
        </ul>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <h2>
          <a id="reaction-types" class="anchor" href="#reaction-types" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Reaction types</h2>

        <p>When creating a reaction, the allowed values for the <code>content</code> parameter are as follows (with the
          corresponding emoji for reference):</p>

        <table>
          <thead>
            <tr>
              <th>content</th>
              <th>emoji</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>+1</code></td>
              <td><img class="emoji" title=":+1:" alt=":+1:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f44d.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>-1</code></td>
              <td><img class="emoji" title=":-1:" alt=":-1:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f44e.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>laugh</code></td>
              <td><img class="emoji" title=":smile:" alt=":smile:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f604.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>confused</code></td>
              <td><img class="emoji" title=":confused:" alt=":confused:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f615.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>heart</code></td>
              <td><img class="emoji" title=":heart:" alt=":heart:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/2764.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>hooray</code></td>
              <td><img class="emoji" title=":tada:" alt=":tada:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f389.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>rocket</code></td>
              <td><img class="emoji" title=":rocket:" alt=":rocket:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f680.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
            <tr>
              <td><code>eyes</code></td>
              <td><img class="emoji" title=":eyes:" alt=":eyes:"
                  src="https://github.githubassets.com/images/icons/emoji/unicode/1f440.png" height="20" width="20"
                  align="absmiddle"></td>
            </tr>
          </tbody>
        </table>

        <h2>
          <a id="list-reactions-for-a-commit-comment" class="anchor" href="#list-reactions-for-a-commit-comment"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List reactions for a
          commit comment<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/repos/comments/">commit comment</a>.</p>

        <pre><code>GET /repos/:owner/:repo/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters" class="anchor" href="#parameters" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a commit comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response" class="anchor" href="#response" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-commit-comment" class="anchor" href="#create-reaction-for-a-commit-comment"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Create reaction for a
          commit comment<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/repos/comments/">commit comment</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this commit comment.</p>

        <pre><code>POST /repos/:owner/:repo/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters-1" class="anchor" href="#parameters-1" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                commit comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example" class="anchor" href="#example" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-1" class="anchor" href="#response-1" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-a-commit-comment-reaction" class="anchor" href="#delete-a-commit-comment-reaction"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Delete a commit comment
          reaction<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a repository by <code>repository_id</code> using the route
            <code>DELETE /repositories/:repository_id/comments/:comment_id/reactions/:reaction_id</code>.</p>

        </div>

        <p>Delete a reaction to a <a href="/v3/repos/comments/">commit comment</a>.</p>

        <pre><code>DELETE /repos/:owner/:repo/comments/:comment_id/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-2" class="anchor" href="#response-2" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-an-issue" class="anchor" href="#list-reactions-for-an-issue"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List reactions for an
          issue<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to an <a href="/v3/issues/">issue</a>.</p>

        <pre><code>GET /repos/:owner/:repo/issues/:issue_number/reactions
</code></pre>

        <h3>
          <a id="parameters-2" class="anchor" href="#parameters-2" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to an issue.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-3" class="anchor" href="#response-3" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-an-issue" class="anchor" href="#create-reaction-for-an-issue"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Create reaction for an
          issue</h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to an <a href="/v3/issues/">issue</a>. A response with a <code>Status: 200 OK</code> means
          that you already added the reaction type to this issue.</p>

        <pre><code>POST /repos/:owner/:repo/issues/:issue_number/reactions
</code></pre>

        <h3>
          <a id="parameters-3" class="anchor" href="#parameters-3" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                issue.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-1" class="anchor" href="#example-1" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-4" class="anchor" href="#response-4" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-an-issue-reaction" class="anchor" href="#delete-an-issue-reaction" aria-hidden="true"><span
              aria-hidden="true" class="octicon octicon-link"></span></a>Delete an issue reaction<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a repository by <code>repository_id</code> using the route
            <code>DELETE /repositories/:repository_id/issues/:issue_number/reactions/:reaction_id</code>.</p>

        </div>

        <p>Delete a reaction to an <a href="/v3/issues/">issue</a>.</p>

        <pre><code>DELETE /repos/:owner/:repo/issues/:issue_number/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-5" class="anchor" href="#response-5" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-an-issue-comment" class="anchor" href="#list-reactions-for-an-issue-comment"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List reactions for an
          issue comment<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to an <a href="/v3/issues/comments/">issue comment</a>.</p>

        <pre><code>GET /repos/:owner/:repo/issues/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters-4" class="anchor" href="#parameters-4" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to an issue comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-6" class="anchor" href="#response-6" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-an-issue-comment" class="anchor" href="#create-reaction-for-an-issue-comment"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Create reaction for an
          issue comment<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to an <a href="/v3/issues/comments/">issue comment</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this issue comment.</p>

        <pre><code>POST /repos/:owner/:repo/issues/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters-5" class="anchor" href="#parameters-5" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                issue comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-2" class="anchor" href="#example-2" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-7" class="anchor" href="#response-7" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-an-issue-comment-reaction" class="anchor" href="#delete-an-issue-comment-reaction"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Delete an issue comment
          reaction<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a repository by <code>repository_id</code> using the route
            <code>DELETE delete /repositories/:repository_id/issues/comments/:comment_id/reactions/:reaction_id</code>.
          </p>

        </div>

        <p>Delete a reaction to an <a href="/v3/issues/comments/">issue comment</a>.</p>

        <pre><code>DELETE /repos/:owner/:repo/issues/comments/:comment_id/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-8" class="anchor" href="#response-8" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-a-pull-request-review-comment" class="anchor"
            href="#list-reactions-for-a-pull-request-review-comment" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>List reactions for a pull request review comment<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/pulls/comments/">pull request review comment</a>.</p>

        <pre><code>GET /repos/:owner/:repo/pulls/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters-6" class="anchor" href="#parameters-6" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a pull request review comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-9" class="anchor" href="#response-9" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-pull-request-review-comment" class="anchor"
            href="#create-reaction-for-a-pull-request-review-comment" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Create reaction for a pull request review comment<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/pulls/comments/">pull request review comment</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this pull request review
          comment.</p>

        <pre><code>POST /repos/:owner/:repo/pulls/comments/:comment_id/reactions
</code></pre>

        <h3>
          <a id="parameters-7" class="anchor" href="#parameters-7" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                pull request review comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-3" class="anchor" href="#example-3" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-10" class="anchor" href="#response-10" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-a-pull-request-comment-reaction" class="anchor" href="#delete-a-pull-request-comment-reaction"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Delete a pull request
          comment reaction<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a repository by <code>repository_id</code> using the route
            <code>DELETE /repositories/:repository_id/pulls/comments/:comment_id/reactions/:reaction_id.</code></p>

        </div>

        <p>Delete a reaction to a <a href="/v3/pulls/comments/">pull request review comment</a>.</p>

        <pre><code>DELETE /repos/:owner/:repo/pulls/comments/:comment_id/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-11" class="anchor" href="#response-11" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-a-team-discussion" class="anchor" href="#list-reactions-for-a-team-discussion"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List reactions for a
          team discussion<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/teams/discussions/">team discussion</a>. OAuth access tokens require the
          <code>read:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team by <code>org_id</code> and <code>team_id</code> using
            the route <code>GET /organizations/:org_id/team/:team_id/discussions/:discussion_number/reactions</code>.
          </p>

        </div>

        <pre><code>GET /orgs/:org/teams/:team_slug/discussions/:discussion_number/reactions
</code></pre>

        <h3>
          <a id="parameters-8" class="anchor" href="#parameters-8" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a team discussion.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-12" class="anchor" href="#response-12" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-team-discussion" class="anchor" href="#create-reaction-for-a-team-discussion"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Create reaction for a
          team discussion</h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/teams/discussions/">team discussion</a>. OAuth access tokens require the
          <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this team discussion.</p>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team by <code>org_id</code> and <code>team_id</code> using
            the route <code>POST /organizations/:org_id/team/:team_id/discussions/:discussion_number/reactions</code>.
          </p>

        </div>

        <pre><code>POST /orgs/:org/teams/:team_slug/discussions/:discussion_number/reactions
</code></pre>

        <h3>
          <a id="parameters-9" class="anchor" href="#parameters-9" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                team discussion.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-4" class="anchor" href="#example-4" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-13" class="anchor" href="#response-13" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-team-discussion-reaction" class="anchor" href="#delete-team-discussion-reaction"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Delete team discussion
          reaction<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team or organization with <code>team_id</code> and
            <code>org_id</code> using the route
            <code>DELETE /organizations/:org_id/team/:team_id/discussions/:discussion_number/reactions/:reaction_id</code>.
          </p>

        </div>

        <p>Delete a reaction to a <a href="/v3/teams/discussions/">team discussion</a>. OAuth access tokens require the
          <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <pre><code>DELETE /orgs/:org/teams/:team_slug/discussions/:discussion_number/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-14" class="anchor" href="#response-14" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-a-team-discussion-comment" class="anchor"
            href="#list-reactions-for-a-team-discussion-comment" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>List reactions for a team discussion comment<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/teams/discussion_comments/">team discussion comment</a>. OAuth access
          tokens require the <code>read:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team by <code>org_id</code> and <code>team_id</code> using
            the route
            <code>GET /organizations/:org_id/team/:team_id/discussions/:discussion_number/comments/:comment_number/reactions</code>.
          </p>

        </div>

        <pre><code>GET /orgs/:org/teams/:team_slug/discussions/:discussion_number/comments/:comment_number/reactions
</code></pre>

        <h3>
          <a id="parameters-10" class="anchor" href="#parameters-10" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a team discussion comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-15" class="anchor" href="#response-15" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-team-discussion-comment" class="anchor"
            href="#create-reaction-for-a-team-discussion-comment" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Create reaction for a team discussion comment<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/teams/discussion_comments/">team discussion comment</a>. OAuth access
          tokens require the <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this team discussion comment.
        </p>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team by <code>org_id</code> and <code>team_id</code> using
            the route
            <code>POST /organizations/:org_id/team/:team_id/discussions/:discussion_number/comments/:comment_number/reactions</code>.
          </p>

        </div>

        <pre><code>POST /orgs/:org/teams/:team_slug/discussions/:discussion_number/comments/:comment_number/reactions
</code></pre>

        <h3>
          <a id="parameters-11" class="anchor" href="#parameters-11" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                team discussion comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-5" class="anchor" href="#example-5" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-16" class="anchor" href="#response-16" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="delete-team-discussion-comment-reaction" class="anchor" href="#delete-team-discussion-comment-reaction"
            aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Delete team discussion
          comment reaction<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> You can also specify a team or organization with <code>team_id</code> and
            <code>org_id</code> using the route
            <code>DELETE /organizations/:org_id/team/:team_id/discussions/:discussion_number/comments/:comment_number/reactions/:reaction_id</code>.
          </p>

        </div>

        <p>Delete a reaction to a <a href="/v3/teams/discussion_comments/">team discussion comment</a>. OAuth access
          tokens require the <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <pre><code>DELETE /orgs/:org/teams/:team_slug/discussions/:discussion_number/comments/:comment_number/reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-17" class="anchor" href="#response-17" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="delete-a-reaction-legacy" class="anchor" href="#delete-a-reaction-legacy" aria-hidden="true"><span
              aria-hidden="true" class="octicon octicon-link"></span></a>Delete a reaction (Legacy)<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert warning">

          <p><strong>Deprecation Notice:</strong> This endpoint route is deprecated and will be removed from the
            Reactions API. We recommend migrating your existing code to use the new delete reactions endpoints. For more
            information, see this <a href="/changes/2020-02-26-new-delete-reactions-endpoints/">blog post</a>.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>OAuth access tokens require the <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>, when deleting a <a
            href="/v3/teams/discussions/">team discussion</a> or <a href="/v3/teams/discussion_comments/">team
            discussion comment</a>.</p>

        <pre><code>DELETE /reactions/:reaction_id
</code></pre>

        <h3>
          <a id="response-18" class="anchor" href="#response-18" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 204 No Content
</code></pre>


        <h2>
          <a id="list-reactions-for-a-team-discussion-legacy" class="anchor"
            href="#list-reactions-for-a-team-discussion-legacy" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>List reactions for a team discussion (Legacy)<a href="/apps/"
            class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert warning">

          <p><strong>Deprecation Notice:</strong> This endpoint route is deprecated and will be removed from the Teams
            API. We recommend migrating your existing code to use the new <a
              href="/v3/reactions/#list-reactions-for-a-team-discussion"><code>List reactions for a team discussion</code></a>
            endpoint.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/teams/discussions/">team discussion</a>. OAuth access tokens require the
          <code>read:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <pre><code>GET /teams/:team_id/discussions/:discussion_number/reactions
</code></pre>

        <h3>
          <a id="parameters-12" class="anchor" href="#parameters-12" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a team discussion.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-19" class="anchor" href="#response-19" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-team-discussion-legacy" class="anchor"
            href="#create-reaction-for-a-team-discussion-legacy" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Create reaction for a team discussion (Legacy)</h2>

        <div class="alert warning">

          <p><strong>Deprecation Notice:</strong> This endpoint route is deprecated and will be removed from the Teams
            API. We recommend migrating your existing code to use the new <a
              href="/v3/reactions/#create-reaction-for-a-team-discussion"><code>Create reaction for a team discussion</code></a>
            endpoint.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/teams/discussions/">team discussion</a>. OAuth access tokens require the
          <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this team discussion.</p>

        <pre><code>POST /teams/:team_id/discussions/:discussion_number/reactions
</code></pre>

        <h3>
          <a id="parameters-13" class="anchor" href="#parameters-13" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                team discussion.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-6" class="anchor" href="#example-6" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-20" class="anchor" href="#response-20" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h2>
          <a id="list-reactions-for-a-team-discussion-comment-legacy" class="anchor"
            href="#list-reactions-for-a-team-discussion-comment-legacy" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>List reactions for a team discussion comment (Legacy)<a
            href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert warning">

          <p><strong>Deprecation Notice:</strong> This endpoint route is deprecated and will be removed from the Teams
            API. We recommend migrating your existing code to use the new <a
              href="/v3/reactions/#list-reactions-for-a-team-discussion-comment"><code>List reactions for a team discussion comment</code></a>
            endpoint.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>List the reactions to a <a href="/v3/teams/discussion_comments/">team discussion comment</a>. OAuth access
          tokens require the <code>read:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>.</p>

        <pre><code>GET /teams/:team_id/discussions/:discussion_number/comments/:comment_number/reactions
</code></pre>

        <h3>
          <a id="parameters-14" class="anchor" href="#parameters-14" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>Returns a single <a href="/v3/reactions/#reaction-types">reaction type</a>. Omit this parameter to
                list all reactions to a team discussion comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="response-21" class="anchor" href="#response-21" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


        <h2>
          <a id="create-reaction-for-a-team-discussion-comment-legacy" class="anchor"
            href="#create-reaction-for-a-team-discussion-comment-legacy" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Create reaction for a team discussion comment (Legacy)<a
            href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info"
            title="Enabled for GitHub Apps"></a>
        </h2>

        <div class="alert warning">

          <p><strong>Deprecation Notice:</strong> This endpoint route is deprecated and will be removed from the Teams
            API. We recommend migrating your existing code to use the new <a
              href="/v3/reactions/#create-reaction-for-a-team-discussion-comment"><code>Create reaction for a team discussion comment</code></a>
            endpoint.</p>

        </div>

        <div class="alert note">

          <p><strong>Note:</strong> APIs for managing reactions are currently available for developers to preview. See
            the <a href="/changes/2016-05-12-reactions-api-preview">blog post</a> for full details. To access the API
            during the preview period, you must provide a custom <a href="/v3/media">media type</a> in the
            <code>Accept</code> header:</p>

          <pre><code>  application/vnd.github.squirrel-girl-preview+json
</code></pre>

        </div>

        <div class="alert warning">

          <p><strong>Warning:</strong> The API may change without advance notice during the preview period. Preview
            features are not supported for production use. If you experience any issues, contact <a
              href="https://github.com/contact">GitHub Support</a> or <a href="https://premium.githubsupport.com">GitHub
              Premium Support</a>.</p>

        </div>

        <p>Create a reaction to a <a href="/v3/teams/discussion_comments/">team discussion comment</a>. OAuth access
          tokens require the <code>write:discussion</code> <a
            href="/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/">scope</a>. A response with a
          <code>Status: 200 OK</code> means that you already added the reaction type to this team discussion comment.
        </p>

        <pre><code>POST /teams/:team_id/discussions/:discussion_number/comments/:comment_number/reactions
</code></pre>

        <h3>
          <a id="parameters-15" class="anchor" href="#parameters-15" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Parameters</h3>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>content</code></td>
              <td><code>string</code></td>
              <td>
                <strong>Required</strong>. The <a href="/v3/reactions/#reaction-types">reaction type</a> to add to the
                team discussion comment.</td>
            </tr>
          </tbody>
        </table>

        <h3>
          <a id="example-7" class="anchor" href="#example-7" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Example</h3>

        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>


        <h3>
          <a id="response-22" class="anchor" href="#response-22" aria-hidden="true"><span aria-hidden="true"
              class="octicon octicon-link"></span></a>Response</h3>

        <pre class="highlight highlight-headers"><code>Status: 201 Created
</code></pre>


        <pre class="highlight highlight-json"><code><span class="p">{</span><span class="w">
  </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDg6UmVhY3Rpb24x"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"user"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
    </span><span class="nt">"node_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"MDQ6VXNlcjE="</span><span class="p">,</span><span class="w">
    </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
    </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"html_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/octocat"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"followers_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/followers"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"following_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/following{/other_user}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"gists_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/gists{/gist_id}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"starred_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/starred{/owner}{/repo}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"subscriptions_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/subscriptions"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"organizations_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/orgs"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"repos_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/repos"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/events{/privacy}"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"received_events_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat/received_events"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"User"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"site_admin"</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
  </span><span class="p">},</span><span class="w">
  </span><span class="nt">"content"</span><span class="p">:</span><span class="w"> </span><span class="s2">"heart"</span><span class="p">,</span><span class="w">
  </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2016-05-20T20:09:31Z"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span></code></pre>

      </div>

      <div class="nav-select nav-select-bottom hide-lg hide-xl">
        <select class="form-select mt-4" onchange="if (this.value) window.location.href=this.value">
          <option value="">Navigate the docs…</option>

          <optgroup label="Overview">
            <option value="/v3/">API Overview</a></h3>
            <option value="/v3/media/">Media Types</option>
            <option value="/v3/oauth_authorizations/">OAuth Authorizations API</option>
            <option value="/v3/auth/">Other Authentication Methods</option>
            <option value="/v3/troubleshooting/">Troubleshooting</option>
            <option value="/v3/previews/">API Previews</option>
            <option value="/v3/versions/">Versions</option>
          </optgroup>

          <optgroup label="Activity">
            <option value="/v3/activity/">Activity overview</a></h3>
            <option value="/v3/activity/events/">Events</option>
            <option value="/v3/activity/events/types/">Event Types &amp; Payloads</option>
            <option value="/v3/activity/feeds/">Feeds</option>
            <option value="/v3/activity/notifications/">Notifications</option>
            <option value="/v3/activity/starring/">Starring</option>
            <option value="/v3/activity/watching/">Watching</option>
          </optgroup>

          <optgroup label="Checks">
            <option value="/v3/checks/">Checks</a></h3>
            <option value="/v3/checks/runs/">Check Runs</option>
            <option value="/v3/checks/suites/">Check Suites</option>
          </optgroup>



          <optgroup label="Gists">
            <option value="/v3/gists/">Gists overview</a></h3>
            <option value="/v3/gists/comments/">Comments</option>
          </optgroup>

          <optgroup label="Git Data">
            <option value="/v3/git/">Git Data overview</a></h3>
            <option value="/v3/git/blobs/">Blobs</option>
            <option value="/v3/git/commits/">Commits</option>
            <option value="/v3/git/refs/">References</option>
            <option value="/v3/git/tags/">Tags</option>
            <option value="/v3/git/trees/">Trees</option>
          </optgroup>


          <optgroup label="GitHub Actions">
            <option value="/v3/actions/">GitHub Actions Overview</a></h3>
            <option value="/v3/actions/artifacts/">Artifacts</option>
            <option value="/v3/actions/secrets/">Secrets</option>
            <option value="/v3/actions/self_hosted_runners/">Self-hosted runners</option>
            <option value="/v3/actions/workflows/">Workflows</option>
            <option value="/v3/actions/workflow_jobs/">Workflow jobs</option>
            <option value="/v3/actions/workflow_runs/">Workflow runs</option>
          </optgroup>


          <optgroup label="GitHub Apps">
            <option value="/v3/apps/">GitHub Apps overview</a></h3>

            <option value="/v3/apps/oauth_applications/">OAuth Applications API</option>

            <option value="/v3/apps/installations/">Installations</option>
            <option value="/v3/apps/permissions/">Permissions</option>
            <option value="/v3/apps/available-endpoints/">Available Endpoints</option>
          </optgroup>


          <optgroup label="Marketplace">
            <option value="/v3/apps/marketplace/">GitHub Marketplace</option>
            </h3>
          </optgroup>



          <optgroup label="Interactions">
            <option value="/v3/interactions/">Interactions</a></h3>
            <option value="/v3/interactions/orgs/">Organization</option>
            <option value="/v3/interactions/repos/">Repository</option>
          </optgroup>


          <optgroup label="Issues">
            <option value="/v3/issues/">Issues overview</a></h3>
            <option value="/v3/issues/assignees/">Assignees</option>
            <option value="/v3/issues/comments/">Comments</option>
            <option value="/v3/issues/events/">Events</option>
            <option value="/v3/issues/labels/">Labels</option>
            <option value="/v3/issues/milestones/">Milestones</option>
            <option value="/v3/issues/timeline/">Timeline</option>
          </optgroup>


          <optgroup label="Migrations">
            <option value="/v3/migrations/">Migrations overview</a></h3>
            <option value="/v3/migrations/orgs/">Organization</option>
            <option value="/v3/migrations/source_imports/">Source Imports</option>
            <option value="/v3/migrations/users/">User</option>
          </optgroup>


          <optgroup label="Miscellaneous">
            <option value="/v3/misc/">Miscellaneous overview</a></h3>
            <option value="/v3/codes_of_conduct/">Codes of Conduct</option>
            <option value="/v3/emojis/">Emojis</option>
            <option value="/v3/gitignore/">Gitignore</option>
            <option value="/v3/licenses/">Licenses</option>
            <option value="/v3/markdown/">Markdown</option>
            <option value="/v3/meta/">Meta</option>

            <option value="/v3/rate_limit/">Rate Limit</option>

          </optgroup>

          <optgroup label="Organizations">
            <option value="/v3/orgs/">Organizations overview</a></h3>

            <option value="/v3/orgs/blocking/">Blocking Users &#40;Organizations&#41;</option>

            <option value="/v3/orgs/members/">Members</option>
            <option value="/v3/orgs/outside_collaborators/">Outside Collaborators</option>
            <option value="/v3/orgs/hooks/">Webhooks</option>
          </optgroup>

          <optgroup label="Projects">
            <option value="/v3/projects/">Projects overview</a></h3>
            <option value="/v3/projects/cards/">Cards</option>
            <option value="/v3/projects/collaborators/">Collaborators</option>
            <option value="/v3/projects/columns/">Columns</option>
          </optgroup>

          <optgroup label="Pull Requests">
            <option value="/v3/pulls/">Pull Requests overview</a></h3>
            <option value="/v3/pulls/reviews/">Reviews</option>
            <option value="/v3/pulls/comments/">Review Comments</option>
            <option value="/v3/pulls/review_requests/">Review Requests</option>
          </optgroup>

          <optgroup label="Reactions">
            <option value="/v3/reactions/">Reactions overview</a></h3>
            <option value="/v3/reactions/#list-reactions-for-a-commit-comment">Commit Comment</option>
            <option value="/v3/reactions/#list-reactions-for-an-issue">Issue</option>
            <option value="/v3/reactions/#list-reactions-for-an-issue-comment">Issue Comment</option>
            <option value="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review Comment
            </option>
            <option value="/v3/reactions/#list-reactions-for-a-team-discussion">Team Discussion</option>
            <option value="/v3/reactions/#list-reactions-for-a-team-discussion-comment">Team Discussion Comment</option>
          </optgroup>

          <optgroup label="Repositories">
            <option value="/v3/repos/">Repositories overview</a></h3>
            <option value="/v3/repos/branches/">Branches</option>
            <option value="/v3/repos/collaborators/">Collaborators</option>
            <option value="/v3/repos/comments/">Comments</option>
            <option value="/v3/repos/commits/">Commits</option>

            <option value="/v3/repos/community/">Community</option>

            <option value="/v3/repos/contents/">Contents</option>
            <option value="/v3/repos/keys/">Deploy Keys</option>
            <option value="/v3/repos/deployments/">Deployments</option>
            <option value="/v3/repos/downloads/">Downloads</option>
            <option value="/v3/repos/forks/">Forks</option>
            <option value="/v3/repos/invitations/">Invitations</option>
            <option value="/v3/repos/merging/">Merging</option>
            <option value="/v3/repos/pages/">Pages</option>
            <option value="/v3/repos/releases/">Releases</option>
            <option value="/v3/repos/statistics/">Statistics</option>
            <option value="/v3/repos/statuses/">Statuses</option>

            <option value="/v3/repos/traffic/">Traffic</option>

            <option value="/v3/repos/hooks/">Webhooks</option>
          </optgroup>

          <optgroup label="Search">
            <option value="/v3/search/">Search overview</a></h3>
            <option value="/v3/search/#search-repositories">Repositories</option>
            <option value="/v3/search/#search-code">Code</option>
            <option value="/v3/search/#search-commits">Commits</option>
            <option value="/v3/search/#search-issues-and-pull-requests">Issues</option>
            <option value="/v3/search/#search-users">Users</option>
            <option value="/v3/search/#search-topics">Topics</option>
            <option value="/v3/search/#text-match-metadata">Text match metadata</<option>
            <option value="/v3/search/legacy/">Legacy search</option>
          </optgroup>

          <optgroup label="Teams">
            <option value="/v3/teams/">Teams</a></h3>
            <option value="/v3/teams/discussions/">Discussions</option>
            <option value="/v3/teams/discussion_comments/">Discussion comments</option>
            <option value="/v3/teams/members/">Members</option>
            <option value="/v3/teams/team_sync/">Team sunchronization</option>
          </optgroup>


          <optgroup label="SCIM">
            <option value="/v3/scim/">SCIM</a></h3>
          </optgroup>


          <optgroup label="Users">
            <option value="/v3/users/">Users overview</a></h3>

            <option value="/v3/users/blocking/">Blocking Users</option>

            <option value="/v3/users/emails/">Emails</option>
            <option value="/v3/users/followers/">Followers</option>
            <option value="/v3/users/keys/">Git SSH Keys</option>
            <option value="/v3/users/gpg_keys/">GPG Keys</option>
          </optgroup>
        </select>

      </div>

      <div id="js-sidebar" class="sidebar-shell col-md-4 hide-sm hide-md d-md-block d-lg-block d-xl-block">
        <div class="js-toggle-list sidebar-module sidebar-menu expandable">
          <ul>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/">Overview</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/media/">Media Types</a></li>
                <li><a href="/v3/oauth_authorizations/">OAuth Authorizations API</a></li>
                <li><a href="/v3/auth/">Other Authentication Methods</a></li>
                <li><a href="/v3/troubleshooting/">Troubleshooting</a></li>
                <li><a href="/v3/previews/">API Previews</a></li>
                <li><a href="/v3/versions/">Versions</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/activity/">Activity</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/activity/events/">Events</a></li>
                <li><a href="/v3/activity/events/types/">Event Types &amp; Payloads</a></li>
                <li><a href="/v3/activity/feeds/">Feeds</a></li>
                <li><a href="/v3/activity/notifications/">Notifications</a></li>
                <li><a href="/v3/activity/starring/">Starring</a></li>
                <li><a href="/v3/activity/watching/">Watching</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/checks/">Checks</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/checks/runs/">Check Runs</a></li>
                <li><a href="/v3/checks/suites/">Check Suites</a></li>
              </ul>
            </li>


            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/gists/">Gists</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/gists/comments/">Comments</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/git/">Git
                  Data</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/git/blobs/">Blobs</a></li>
                <li><a href="/v3/git/commits/">Commits</a></li>
                <li><a href="/v3/git/refs/">References</a></li>
                <li><a href="/v3/git/tags/">Tags</a></li>
                <li><a href="/v3/git/trees/">Trees</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/actions/">GitHub Actions</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/actions/artifacts/">Artifacts</a></li>
                <li><a href="/v3/actions/secrets/">Secrets</a></li>
                <li><a href="/v3/actions/self_hosted_runners/">Self-hosted runners</a></li>
                <li><a href="/v3/actions/workflows/">Workflows</a></li>
                <li><a href="/v3/actions/workflow_jobs/">Workflow jobs</a></li>
                <li><a href="/v3/actions/workflow_runs/">Workflow runs</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/apps/">GitHub Apps</a></h3>
              <ul class="js-guides">

                <li><a href="/v3/apps/oauth_applications/">OAuth Applications API</a></li>

                <li><a href="/v3/apps/installations/">Installations</a></li>
                <li><a href="/v3/apps/permissions/">Permissions</a></li>
                <li><a href="/v3/apps/available-endpoints/">Available Endpoints</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3 class="standalone">
                <a href="/v3/apps/marketplace/">GitHub Marketplace</a>
              </h3>
            </li>


            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/interactions/">Interactions</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/interactions/orgs/">Organization</a></li>
                <li><a href="/v3/interactions/repos/">Repository</a></li>
              </ul>
            </li>


            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/issues/">Issues</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/issues/assignees/">Assignees</a></li>
                <li><a href="/v3/issues/comments/">Comments</a></li>
                <li><a href="/v3/issues/events/">Events</a></li>
                <li><a href="/v3/issues/labels/">Labels</a></li>
                <li><a href="/v3/issues/milestones/">Milestones</a></li>
                <li><a href="/v3/issues/timeline/">Timeline</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/migrations/">Migrations</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/migrations/orgs/">Organization</a></li>
                <li><a href="/v3/migrations/source_imports/">Source Imports</a></li>
                <li><a href="/v3/migrations/users/">User</a></li>
              </ul>
            </li>


            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/misc/">Miscellaneous</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/codes_of_conduct/">Codes of Conduct</a></li>
                <li><a href="/v3/emojis/">Emojis</a></li>
                <li><a href="/v3/gitignore/">Gitignore</a></li>
                <li><a href="/v3/licenses/">Licenses</a></li>
                <li><a href="/v3/markdown/">Markdown</a></li>
                <li><a href="/v3/meta/">Meta</a></li>
                <li><a href="/v3/rate_limit/">Rate Limit</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/orgs/">Organizations</a></h3>
              <ul class="js-guides">

                <li><a href="/v3/orgs/blocking/">Blocking Users</a></li>

                <li><a href="/v3/orgs/members/">Members</a></li>
                <li><a href="/v3/orgs/outside_collaborators/">Outside Collaborators</a></li>
                <li><a href="/v3/orgs/hooks/">Webhooks</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/projects/">Projects</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/projects/cards/">Cards</a></li>
                <li><a href="/v3/projects/collaborators/">Collaborators</a></li>
                <li><a href="/v3/projects/columns/">Columns</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/pulls/">Pull Requests</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/pulls/reviews/">Reviews</a></li>
                <li><a href="/v3/pulls/comments/">Review Comments</a></li>
                <li><a href="/v3/pulls/review_requests/">Review Requests</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/reactions/">Reactions</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/reactions/#list-reactions-for-a-commit-comment">Commit Comment</a></li>
                <li><a href="/v3/reactions/#list-reactions-for-an-issue">Issue</a></li>
                <li><a href="/v3/reactions/#list-reactions-for-an-issue-comment">Issue Comment</a></li>
                <li><a href="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review
                    Comment</a></li>
                <li><a href="/v3/reactions/#list-reactions-for-a-team-discussion">Team Discussion</a></li>
                <li><a href="/v3/reactions/#list-reactions-for-a-team-discussion-comment">Team Discussion Comment</a>
                </li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/repos/">Repositories</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/repos/branches/">Branches</a></li>
                <li><a href="/v3/repos/collaborators/">Collaborators</a></li>
                <li><a href="/v3/repos/comments/">Comments</a></li>
                <li><a href="/v3/repos/commits/">Commits</a></li>

                <li><a href="/v3/repos/community/">Community</a></li>

                <li><a href="/v3/repos/contents/">Contents</a></li>
                <li><a href="/v3/repos/keys/">Deploy Keys</a></li>
                <li><a href="/v3/repos/deployments/">Deployments</a></li>
                <li><a href="/v3/repos/downloads/">Downloads</a></li>
                <li><a href="/v3/repos/forks/">Forks</a></li>
                <li><a href="/v3/repos/invitations/">Invitations</a></li>
                <li><a href="/v3/repos/merging/">Merging</a></li>
                <li><a href="/v3/repos/pages/">Pages</a></li>
                <li><a href="/v3/repos/releases/">Releases</a></li>
                <li><a href="/v3/repos/statistics/">Statistics</a></li>
                <li><a href="/v3/repos/statuses/">Statuses</a></li>

                <li><a href="/v3/repos/traffic/">Traffic</a></li>

                <li><a href="/v3/repos/hooks/">Webhooks</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/search/">Search</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/search/#search-repositories">Repositories</a></li>
                <li><a href="/v3/search/#search-code">Code</a></li>
                <li><a href="/v3/search/#search-commits">Commits</a></li>
                <li><a href="/v3/search/#search-issues-and-pull-requests">Issues</a></li>
                <li><a href="/v3/search/#search-users">Users</a></li>
                <li><a href="/v3/search/#search-topics">Topics</a></li>
                <li><a href="/v3/search/#text-match-metadata">Text match metadata</a></li>
                <li><a href="/v3/search/legacy/">Legacy search</a></li>
              </ul>
            </li>
            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/teams/">Teams</a></h3>
              <ul class="js-guides">
                <li><a href="/v3/teams/discussions/">Discussions</a></li>
                <li><a href="/v3/teams/discussion_comments/">Discussion comments</a></li>
                <li><a href="/v3/teams/members/">Members</a></li>
                <li><a href="/v3/teams/team_sync/">Team synchronization</a></li>
              </ul>
            </li>

            <li class="js-topic">
              <h3 class="standalone">
                <a href="/v3/scim/">SCIM</a>
              </h3>
            </li>

            <li class="js-topic">
              <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a
                  href="/v3/users/">Users</a></h3>
              <ul class="js-guides">

                <li><a href="/v3/users/blocking/">Blocking Users</a></li>

                <li><a href="/v3/users/emails/">Emails</a></li>
                <li><a href="/v3/users/followers/">Followers</a></li>
                <li><a href="/v3/users/keys/">Git SSH Keys</a></li>
                <li><a href="/v3/users/gpg_keys/">GPG Keys</a></li>
              </ul>
            </li>
          </ul>
        </div> <!-- /sidebar-module -->

        <div class="sidebar-module api-status py-4 text-center"><a href="https://status.github.com" class="unknown">API
            Status</a></div>


      </div><!-- /sidebar-shell -->

    </section>
  </div>

  <footer class="footer text-muted">
    <div class="container">
      <span class="footer-legal">&copy; 2020 GitHub Inc. All rights reserved.</span>
      <span class="mega-octicon octicon-mark-github footer-mark"></span>
      <nav class="footer-nav">
        <a class="footer-nav-item" href="https://help.github.com/articles/github-terms-of-service">Terms of service</a>
        <a class="footer-nav-item" href="https://github.com/site/privacy">Privacy</a>
        <a class="footer-nav-item" href="https://github.com/security">Security</a>
        <a class="footer-nav-item" href="https://github.com/support">Support</a>
      </nav>
    </div>
  </footer>

  <script>
    (function (i, s, o, g, r, a, m) {
      i['GoogleAnalyticsObject'] = r; i[r] = i[r] || function () {
        (i[r].q = i[r].q || []).push(arguments)
      }, i[r].l = 1 * new Date(); a = s.createElement(o),
        m = s.getElementsByTagName(o)[0]; a.async = 1; a.src = g; m.parentNode.insertBefore(a, m)
    })(window, document, 'script', '//www.google-analytics.com/analytics.js', 'ga');

    ga('create', 'UA-3769691-37', 'github.com');
    ga('send', 'pageview');
  </script>

</body>

</html>
`

var reactionsGoFileOriginal = `// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
)

// ReactionsService provides access to the reactions-related functions in the
// GitHub API.
type ReactionsService service

// Reaction represents a GitHub reaction.
type Reaction struct {
	// ID is the Reaction ID.
	ID     *int64  ` + "`" + `json:"id,omitempty"` + "`" + `
	User   *User   ` + "`" + `json:"user,omitempty"` + "`" + `
	NodeID *string ` + "`" + `json:"node_id,omitempty"` + "`" + `
	// Content is the type of reaction.
	// Possible values are:
	//     "+1", "-1", "laugh", "confused", "heart", "hooray".
	Content *string ` + "`" + `json:"content,omitempty"` + "`" + `
}

// Reactions represents a summary of GitHub reactions.
type Reactions struct {
	TotalCount *int    ` + "`" + `json:"total_count,omitempty"` + "`" + `
	PlusOne    *int    ` + "`" + `json:"+1,omitempty"` + "`" + `
	MinusOne   *int    ` + "`" + `json:"-1,omitempty"` + "`" + `
	Laugh      *int    ` + "`" + `json:"laugh,omitempty"` + "`" + `
	Confused   *int    ` + "`" + `json:"confused,omitempty"` + "`" + `
	Heart      *int    ` + "`" + `json:"heart,omitempty"` + "`" + `
	Hooray     *int    ` + "`" + `json:"hooray,omitempty"` + "`" + `
	URL        *string ` + "`" + `json:"url,omitempty"` + "`" + `
}

func (r Reaction) String() string {
	return Stringify(r)
}

// ListCommentReactionOptions specifies the optional parameters to the
// ReactionsService.ListCommentReactions method.
type ListCommentReactionOptions struct {
	// Content restricts the returned comment reactions to only those with the given type.
	// Omit this parameter to list all reactions to a commit comment.
	// Possible values are: "+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", or "eyes".
	Content string ` + "`" + `url:"content,omitempty"` + "`" + `

	ListOptions
}

// ListCommentReactions lists the reactions for a commit comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-commit-comment
func (s *ReactionsService) ListCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListCommentReactionOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateCommentReaction creates a reaction for a commit comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-commit-comment
func (s *ReactionsService) CreateCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteCommentReaction deletes the reaction for a commit comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-commit-comment-reaction
func (s *ReactionsService) DeleteCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, u)
}

// DeleteCommentReactionByID deletes the reaction for a commit comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-commit-comment-reaction
func (s *ReactionsService) DeleteCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	u := fmt.Sprintf("repositories/%v/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, u)
}

// ListIssueReactions lists the reactions for an issue.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-an-issue
func (s *ReactionsService) ListIssueReactions(ctx context.Context, owner, repo string, number int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/reactions", owner, repo, number)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateIssueReaction creates a reaction for an issue.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-an-issue
func (s *ReactionsService) CreateIssueReaction(ctx context.Context, owner, repo string, number int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/reactions", owner, repo, number)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteIssueReaction deletes the reaction to an issue.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-reaction
func (s *ReactionsService) DeleteIssueReaction(ctx context.Context, owner, repo string, issueNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/issues/%v/reactions/%v", owner, repo, issueNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteIssueReactionByID deletes the reaction to an issue by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-reaction
func (s *ReactionsService) DeleteIssueReactionByID(ctx context.Context, repoID, issueNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/issues/%v/reactions/%v", repoID, issueNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListIssueCommentReactions lists the reactions for an issue comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-an-issue-comment
func (s *ReactionsService) ListIssueCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateIssueCommentReaction creates a reaction for an issue comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-an-issue-comment
func (s *ReactionsService) CreateIssueCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteIssueCommentReaction deletes the reaction to an issue comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-comment-reaction
func (s *ReactionsService) DeleteIssueCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteIssueCommentReactionByID deletes the reaction to an issue comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-comment-reaction
func (s *ReactionsService) DeleteIssueCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/issues/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListPullRequestCommentReactions lists the reactions for a pull request review comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-an-issue-comment
func (s *ReactionsService) ListPullRequestCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreatePullRequestCommentReaction creates a reaction for a pull request review comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-an-issue-comment
func (s *ReactionsService) CreatePullRequestCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeletePullRequestCommentReaction deletes the reaction to a pull request review comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-pull-request-comment-reaction
func (s *ReactionsService) DeletePullRequestCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeletePullRequestCommentReactionByID deletes the reaction to a pull request review comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-pull-request-comment-reaction
func (s *ReactionsService) DeletePullRequestCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/pulls/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListTeamDiscussionReactions lists the reactions for a team discussion.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-team-discussion
func (s *ReactionsService) ListTeamDiscussionReactions(ctx context.Context, teamID int64, discussionNumber int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/reactions", teamID, discussionNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateTeamDiscussionReaction creates a reaction for a team discussion.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-team-discussion
func (s *ReactionsService) CreateTeamDiscussionReaction(ctx context.Context, teamID int64, discussionNumber int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/reactions", teamID, discussionNumber)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteTeamDiscussionReaction deletes the reaction to a team discussion.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-reaction
func (s *ReactionsService) DeleteTeamDiscussionReaction(ctx context.Context, org, teamSlug string, discussionNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/teams/%v/discussions/%v/reactions/%v", org, teamSlug, discussionNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteTeamDiscussionReactionByOrgIDAndTeamID deletes the reaction to a team discussion by organization ID and team ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-reaction
func (s *ReactionsService) DeleteTeamDiscussionReactionByOrgIDAndTeamID(ctx context.Context, orgID, teamID, discussionNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("organizations/%v/team/%v/discussions/%v/reactions/%v", orgID, teamID, discussionNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListTeamDiscussionCommentReactions lists the reactions for a team discussion comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-team-discussion-comment
func (s *ReactionsService) ListTeamDiscussionCommentReactions(ctx context.Context, teamID int64, discussionNumber, commentNumber int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/comments/%v/reactions", teamID, discussionNumber, commentNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, nil, err
	}
	return m, resp, nil
}

// CreateTeamDiscussionCommentReaction creates a reaction for a team discussion comment.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-team-discussion-comment
func (s *ReactionsService) CreateTeamDiscussionCommentReaction(ctx context.Context, teamID int64, discussionNumber, commentNumber int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/comments/%v/reactions", teamID, discussionNumber, commentNumber)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteTeamDiscussionCommentReaction deletes the reaction to a team discussion comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-comment-reaction
func (s *ReactionsService) DeleteTeamDiscussionCommentReaction(ctx context.Context, org, teamSlug string, discussionNumber, commentNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/teams/%v/discussions/%v/comments/%v/reactions/%v", org, teamSlug, discussionNumber, commentNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteTeamDiscussionCommentReactionByOrgIDAndTeamID deletes the reaction to a team discussion comment by organization ID and team ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-comment-reaction
func (s *ReactionsService) DeleteTeamDiscussionCommentReactionByOrgIDAndTeamID(ctx context.Context, orgID, teamID, discussionNumber, commentNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("organizations/%v/team/%v/discussions/%v/comments/%v/reactions/%v", orgID, teamID, discussionNumber, commentNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

func (s *ReactionsService) deleteReaction(ctx context.Context, url string) (*Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	return s.client.Do(ctx, req, nil)
}
`

var reactionsGoFileWant = `// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
)

// ReactionsService provides access to the reactions-related functions in the
// GitHub API.
type ReactionsService service

// Reaction represents a GitHub reaction.
type Reaction struct {
	// ID is the Reaction ID.
	ID     *int64  ` + "`" + `json:"id,omitempty"` + "`" + `
	User   *User   ` + "`" + `json:"user,omitempty"` + "`" + `
	NodeID *string ` + "`" + `json:"node_id,omitempty"` + "`" + `
	// Content is the type of reaction.
	// Possible values are:
	//     "+1", "-1", "laugh", "confused", "heart", "hooray".
	Content *string ` + "`" + `json:"content,omitempty"` + "`" + `
}

// Reactions represents a summary of GitHub reactions.
type Reactions struct {
	TotalCount *int    ` + "`" + `json:"total_count,omitempty"` + "`" + `
	PlusOne    *int    ` + "`" + `json:"+1,omitempty"` + "`" + `
	MinusOne   *int    ` + "`" + `json:"-1,omitempty"` + "`" + `
	Laugh      *int    ` + "`" + `json:"laugh,omitempty"` + "`" + `
	Confused   *int    ` + "`" + `json:"confused,omitempty"` + "`" + `
	Heart      *int    ` + "`" + `json:"heart,omitempty"` + "`" + `
	Hooray     *int    ` + "`" + `json:"hooray,omitempty"` + "`" + `
	URL        *string ` + "`" + `json:"url,omitempty"` + "`" + `
}

func (r Reaction) String() string {
	return Stringify(r)
}

// ListCommentReactionOptions specifies the optional parameters to the
// ReactionsService.ListCommentReactions method.
type ListCommentReactionOptions struct {
	// Content restricts the returned comment reactions to only those with the given type.
	// Omit this parameter to list all reactions to a commit comment.
	// Possible values are: "+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", or "eyes".
	Content string ` + "`" + `url:"content,omitempty"` + "`" + `

	ListOptions
}

// ListCommentReactions lists the reactions for a commit comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-commit-comment
func (s *ReactionsService) ListCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListCommentReactionOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateCommentReaction creates a reaction for a commit comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-commit-comment
func (s *ReactionsService) CreateCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteCommentReaction deletes the reaction for a commit comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-commit-comment-reaction
func (s *ReactionsService) DeleteCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, u)
}

// DeleteCommentReactionByID deletes the reaction for a commit comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-commit-comment-reaction
func (s *ReactionsService) DeleteCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	u := fmt.Sprintf("repositories/%v/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, u)
}

// ListIssueReactions lists the reactions for an issue.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-an-issue
func (s *ReactionsService) ListIssueReactions(ctx context.Context, owner, repo string, number int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/reactions", owner, repo, number)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateIssueReaction creates a reaction for an issue.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-an-issue
func (s *ReactionsService) CreateIssueReaction(ctx context.Context, owner, repo string, number int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/reactions", owner, repo, number)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteIssueReaction deletes the reaction to an issue.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-reaction
func (s *ReactionsService) DeleteIssueReaction(ctx context.Context, owner, repo string, issueNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/issues/%v/reactions/%v", owner, repo, issueNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteIssueReactionByID deletes the reaction to an issue by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-reaction
func (s *ReactionsService) DeleteIssueReactionByID(ctx context.Context, repoID, issueNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/issues/%v/reactions/%v", repoID, issueNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListIssueCommentReactions lists the reactions for an issue comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-an-issue-comment
func (s *ReactionsService) ListIssueCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateIssueCommentReaction creates a reaction for an issue comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-an-issue-comment
func (s *ReactionsService) CreateIssueCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteIssueCommentReaction deletes the reaction to an issue comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-comment-reaction
func (s *ReactionsService) DeleteIssueCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/issues/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteIssueCommentReactionByID deletes the reaction to an issue comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-an-issue-comment-reaction
func (s *ReactionsService) DeleteIssueCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/issues/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListPullRequestCommentReactions lists the reactions for a pull request review comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-pull-request-review-comment
func (s *ReactionsService) ListPullRequestCommentReactions(ctx context.Context, owner, repo string, id int64, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreatePullRequestCommentReaction creates a reaction for a pull request review comment.
// Note that if you have already created a reaction of type content, the
// previously created reaction will be returned with Status: 200 OK.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-pull-request-review-comment
func (s *ReactionsService) CreatePullRequestCommentReaction(ctx context.Context, owner, repo string, id int64, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions", owner, repo, id)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeletePullRequestCommentReaction deletes the reaction to a pull request review comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-pull-request-comment-reaction
func (s *ReactionsService) DeletePullRequestCommentReaction(ctx context.Context, owner, repo string, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repos/%v/%v/pulls/comments/%v/reactions/%v", owner, repo, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeletePullRequestCommentReactionByID deletes the reaction to a pull request review comment by repository ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-a-pull-request-comment-reaction
func (s *ReactionsService) DeletePullRequestCommentReactionByID(ctx context.Context, repoID, commentID, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("repositories/%v/pulls/comments/%v/reactions/%v", repoID, commentID, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListTeamDiscussionReactions lists the reactions for a team discussion.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-team-discussion-legacy
func (s *ReactionsService) ListTeamDiscussionReactions(ctx context.Context, teamID int64, discussionNumber int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/reactions", teamID, discussionNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateTeamDiscussionReaction creates a reaction for a team discussion.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-team-discussion-legacy
func (s *ReactionsService) CreateTeamDiscussionReaction(ctx context.Context, teamID int64, discussionNumber int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/reactions", teamID, discussionNumber)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteTeamDiscussionReaction deletes the reaction to a team discussion.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-reaction
func (s *ReactionsService) DeleteTeamDiscussionReaction(ctx context.Context, org, teamSlug string, discussionNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/teams/%v/discussions/%v/reactions/%v", org, teamSlug, discussionNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteTeamDiscussionReactionByOrgIDAndTeamID deletes the reaction to a team discussion by organization ID and team ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-reaction
func (s *ReactionsService) DeleteTeamDiscussionReactionByOrgIDAndTeamID(ctx context.Context, orgID, teamID, discussionNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("organizations/%v/team/%v/discussions/%v/reactions/%v", orgID, teamID, discussionNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// ListTeamDiscussionCommentReactions lists the reactions for a team discussion comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#list-reactions-for-a-team-discussion-comment-legacy
func (s *ReactionsService) ListTeamDiscussionCommentReactions(ctx context.Context, teamID int64, discussionNumber, commentNumber int, opts *ListOptions) ([]*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/comments/%v/reactions", teamID, discussionNumber, commentNumber)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	var m []*Reaction
	resp, err := s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, nil, err
	}
	return m, resp, nil
}

// CreateTeamDiscussionCommentReaction creates a reaction for a team discussion comment.
// The content should have one of the following values: "+1", "-1", "laugh", "confused", "heart", "hooray".
//
// GitHub API docs: https://developer.github.com/v3/reactions/#create-reaction-for-a-team-discussion-comment-legacy
func (s *ReactionsService) CreateTeamDiscussionCommentReaction(ctx context.Context, teamID int64, discussionNumber, commentNumber int, content string) (*Reaction, *Response, error) {
	u := fmt.Sprintf("teams/%v/discussions/%v/comments/%v/reactions", teamID, discussionNumber, commentNumber)

	body := &Reaction{Content: String(content)}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeReactionsPreview)

	m := &Reaction{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteTeamDiscussionCommentReaction deletes the reaction to a team discussion comment.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-comment-reaction
func (s *ReactionsService) DeleteTeamDiscussionCommentReaction(ctx context.Context, org, teamSlug string, discussionNumber, commentNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/teams/%v/discussions/%v/comments/%v/reactions/%v", org, teamSlug, discussionNumber, commentNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

// DeleteTeamDiscussionCommentReactionByOrgIDAndTeamID deletes the reaction to a team discussion comment by organization ID and team ID.
//
// GitHub API docs: https://developer.github.com/v3/reactions/#delete-team-discussion-comment-reaction
func (s *ReactionsService) DeleteTeamDiscussionCommentReactionByOrgIDAndTeamID(ctx context.Context, orgID, teamID, discussionNumber, commentNumber int, reactionID int64) (*Response, error) {
	url := fmt.Sprintf("organizations/%v/team/%v/discussions/%v/comments/%v/reactions/%v", orgID, teamID, discussionNumber, commentNumber, reactionID)

	return s.deleteReaction(ctx, url)
}

func (s *ReactionsService) deleteReaction(ctx context.Context, url string) (*Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	req.Header.Set("Accept", mediaTypeReactionsPreview)

	return s.client.Do(ctx, req, nil)
}
`
