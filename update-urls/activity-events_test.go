// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func newActivitiesEventsPipeline() *pipelineSetup {
	return &pipelineSetup{
		baseURL:              "https://developer.github.com/v3/activity/events/",
		endpointsFromWebsite: activityEventsWant,
		filename:             "activity_events.go",
		serviceName:          "ActivityService",
		originalGoSource:     activityEventsGoFileOriginal,
		wantGoSource:         activityEventsGoFileWant,
		wantNumEndpoints:     7,
	}
}

func TestPipeline_ActivityEvents(t *testing.T) {
	ps := newActivitiesEventsPipeline()
	ps.setup(t, false, false)
	ps.validate(t)
}

func TestPipeline_ActivityEvents_FirstStripAllURLs(t *testing.T) {
	ps := newActivitiesEventsPipeline()
	ps.setup(t, true, false)
	ps.validate(t)
}

func TestPipeline_ActivityEvents_FirstDestroyReceivers(t *testing.T) {
	ps := newActivitiesEventsPipeline()
	ps.setup(t, false, true)
	ps.validate(t)
}

func TestPipeline_ActivityEvents_FirstStripAllURLsAndDestroyReceivers(t *testing.T) {
	ps := newActivitiesEventsPipeline()
	ps.setup(t, true, true)
	ps.validate(t)
}

func TestParseWebPageEndpoints_ActivityEvents(t *testing.T) {
	got, want := parseWebPageEndpoints(activityEventsTestWebPage), activityEventsWant
	testWebPageHelper(t, got, want)
}

var activityEventsWant = endpointsByFragmentID{
	"list-public-events": []*Endpoint{
		{urlFormats: []string{"events"}, httpMethod: "GET"},
	},

	"list-repository-events": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/events"}, httpMethod: "GET"},
	},

	"list-public-events-for-a-network-of-repositories": []*Endpoint{
		{urlFormats: []string{"networks/%v/%v/events"}, httpMethod: "GET"},
	},

	"list-events-received-by-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"users/%v/received_events"}, httpMethod: "GET"},
	},

	"list-events-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"users/%v/events"}, httpMethod: "GET"},
	},

	"list-public-events-for-a-user": []*Endpoint{
		{urlFormats: []string{"users/%v/events/public"}, httpMethod: "GET"},
	},

	"list-organization-events-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"users/%v/events/orgs/%v"}, httpMethod: "GET"},
	},

	"list-public-organization-events": []*Endpoint{
		{urlFormats: []string{"orgs/%v/events"}, httpMethod: "GET"},
	},

	"list-public-events-received-by-a-user": []*Endpoint{
		{urlFormats: []string{"users/%v/received_events/public"}, httpMethod: "GET"},
	},
}

var activityEventsTestWebPage = `<!DOCTYPE html>
<html lang="en" prefix="og: http://ogp.me/ns#">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <meta http-equiv="Content-Language" content="en-us" />
  <meta http-equiv="imagetoolbar" content="false" />
  <meta name="MSSmartTagsPreventParsing" content="true" />
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Events | GitHub Developer Guide</title>
  <meta property="og:url" content="https://developer.github.com/v3/activity/events/" />
<meta property="og:site_name" content="GitHub Developer" />
<meta property="og:title" content="Events" />
<meta property="og:description" content="Get started with one of our guides, or jump straight into the API documentation." />
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
<meta property="twitter:title" content="Events" />
<meta property="twitter:description" content="Get started with one of our guides, or jump straight into the API documentation." />
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
    <a class="site-header-logo mt-1" href="/"><img src="/assets/images/github-developer-logo.svg" alt="GitHub Developer"></a>
    <nav class="site-header-nav" aria-label="Main Navigation">
      <div class="dropdown" id="api-docs-dropdown">
          <div class="dropdown-button" tabIndex="0" aria-haspopup="true" aria-expanded="false"><span class="dropdown-button-link" tabIndex="-1">Docs</span> <div class="dropdown-caret"></div></div>
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
          <div class="dropdown-button" tabIndex="0" aria-haspopup="true" aria-expanded="false"><span class="dropdown-button-link" tabIndex="-1">Versions</span> <div class="dropdown-caret"></div></div>
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
          <input type="text" id="searchfield" class="form-control" placeholder="Search…"
            autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false" />
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
    <option value="/v3/apps/marketplace/">GitHub Marketplace</option></h3>
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
    <option value="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review Comment</option>
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
<a id="events" class="anchor" href="#events" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Events</h1>

<p>This is a read-only API to the GitHub events. These events power the various activity streams on the site. An events API for repository issues is also available. For more information, see the "<a href="/v3/issues/events/">Issue Events API</a>."</p>

<ul id="markdown-toc">
<li><a href="#list-public-events" id="markdown-toc-list-public-events">List public events</a></li>
<li><a href="#list-repository-events" id="markdown-toc-list-repository-events">List repository events</a></li>
<li><a href="#list-public-events-for-a-network-of-repositories" id="markdown-toc-list-public-events-for-a-network-of-repositories">List public events for a network of repositories</a></li>
<li><a href="#list-public-organization-events" id="markdown-toc-list-public-organization-events">List public organization events</a></li>
<li><a href="#list-events-received-by-the-authenticated-user" id="markdown-toc-list-events-received-by-the-authenticated-user">List events received by the authenticated user</a></li>
<li><a href="#list-public-events-received-by-a-user" id="markdown-toc-list-public-events-received-by-a-user">List public events received by a user</a></li>
<li><a href="#list-events-for-the-authenticated-user" id="markdown-toc-list-events-for-the-authenticated-user">List events for the authenticated user</a></li>
<li><a href="#list-public-events-for-a-user" id="markdown-toc-list-public-events-for-a-user">List public events for a user</a></li>
<li><a href="#list-organization-events-for-the-authenticated-user" id="markdown-toc-list-organization-events-for-the-authenticated-user">List organization events for the authenticated user</a></li>
</ul>

<p>Events are optimized for polling with the "ETag" header.  If no new events have been triggered, you will see a "304 Not Modified" response, and your current rate limit will be untouched.  There is also an "X-Poll-Interval" header that specifies how often (in seconds) you are allowed to poll.  In times of high
server load, the time may increase.  Please obey the header.</p>

<pre class="command-line">
<span class="command">curl -I https://api.github.com/users/tater/events</span>
<span class="output">HTTP/1.1 200 OK</span>
<span class="output">X-Poll-Interval: 60</span>
<span class="output">ETag: "a18c3bded88eb5dbb5c849a489412bf3"</span>
<span class="comment"># The quotes around the ETag value are important</span>
<span class="command">curl -I https://api.github.com/users/tater/events \</span>
<span class="command">   -H 'If-None-Match: "a18c3bded88eb5dbb5c849a489412bf3"'</span>
<span class="output">HTTP/1.1 304 Not Modified</span>
<span class="output">X-Poll-Interval: 60</span>
</pre>

<p>Events support <a href="/v3/#pagination">pagination</a>, however the <code>per_page</code> option is unsupported. The fixed page size is 30 items. Fetching up to ten pages is supported, for a total of 300 events.</p>

<p>Only events created within the past 90 days will be included in timelines. Events older than 90 days will not be included (even if the total number of events in the timeline is less than 300).</p>

<p>All Events have the same response format:</p>

<pre class="highlight highlight-headers"><code>Status: 200 OK
Link: &lt;https://api.github.com/resource?page=2&gt;; rel="next",
      &lt;https://api.github.com/resource?page=5&gt;; rel="last"
</code></pre>


<pre class="highlight highlight-json"><code><span class="p">[</span><span class="w">
  </span><span class="p">{</span><span class="w">
    </span><span class="nt">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"Event"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"public"</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="p">,</span><span class="w">
    </span><span class="nt">"payload"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"repo"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">3</span><span class="p">,</span><span class="w">
      </span><span class="nt">"name"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat/Hello-World"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/repos/octocat/Hello-World"</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"actor"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"octocat"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/users/octocat"</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"org"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
      </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
      </span><span class="nt">"login"</span><span class="p">:</span><span class="w"> </span><span class="s2">"github"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"gravatar_id"</span><span class="p">:</span><span class="w"> </span><span class="s2">""</span><span class="p">,</span><span class="w">
      </span><span class="nt">"url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://api.github.com/orgs/github"</span><span class="p">,</span><span class="w">
      </span><span class="nt">"avatar_url"</span><span class="p">:</span><span class="w"> </span><span class="s2">"https://github.com/images/error/octocat_happy.gif"</span><span class="w">
    </span><span class="p">},</span><span class="w">
    </span><span class="nt">"created_at"</span><span class="p">:</span><span class="w"> </span><span class="s2">"2011-09-06T17:26:27Z"</span><span class="p">,</span><span class="w">
    </span><span class="nt">"id"</span><span class="p">:</span><span class="w"> </span><span class="s2">"12345"</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">]</span><span class="w">
</span></code></pre>


<h2>
<a id="list-public-events" class="anchor" href="#list-public-events" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List public events<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<p>We delay the public events feed by five minutes, which means the most recent event returned by the public events API actually occurred at least five minutes ago.</p>

<pre><code>GET /events
</code></pre>

<h2>
<a id="list-repository-events" class="anchor" href="#list-repository-events" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List repository events<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<pre><code>GET /repos/:owner/:repo/events
</code></pre>

<h2>
<a id="list-public-events-for-a-network-of-repositories" class="anchor" href="#list-public-events-for-a-network-of-repositories" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List public events for a network of repositories<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<pre><code>GET /networks/:owner/:repo/events
</code></pre>

<h2>
<a id="list-public-organization-events" class="anchor" href="#list-public-organization-events" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List public organization events<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<pre><code>GET /orgs/:org/events
</code></pre>

<h2>
<a id="list-events-received-by-the-authenticated-user" class="anchor" href="#list-events-received-by-the-authenticated-user" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List events received by the authenticated user<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<p>These are events that you've received by watching repos and following users. If you are authenticated as the given user, you will see private events. Otherwise, you'll only see public events.</p>

<pre><code>GET /users/:username/received_events
</code></pre>

<h2>
<a id="list-public-events-received-by-a-user" class="anchor" href="#list-public-events-received-by-a-user" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List public events received by a user<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<pre><code>GET /users/:username/received_events/public
</code></pre>

<h2>
<a id="list-events-for-the-authenticated-user" class="anchor" href="#list-events-for-the-authenticated-user" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List events for the authenticated user<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<p>If you are authenticated as the given user, you will see your private events. Otherwise, you'll only see public events.</p>

<pre><code>GET /users/:username/events
</code></pre>

<h2>
<a id="list-public-events-for-a-user" class="anchor" href="#list-public-events-for-a-user" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List public events for a user<a href="/apps/" class="tooltip-link github-apps-marker octicon octicon-info" title="Enabled for GitHub Apps"></a>
</h2>

<pre><code>GET /users/:username/events/public
</code></pre>

<h2>
<a id="list-organization-events-for-the-authenticated-user" class="anchor" href="#list-organization-events-for-the-authenticated-user" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>List organization events for the authenticated user</h2>

<p>This is the user's organization dashboard. You must be authenticated as the user to view this.</p>

<pre><code>GET /users/:username/events/orgs/:org
</code></pre>
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
    <option value="/v3/apps/marketplace/">GitHub Marketplace</option></h3>
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
    <option value="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review Comment</option>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/">Overview</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/activity/">Activity</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/checks/">Checks</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/checks/runs/">Check Runs</a></li>
          <li><a href="/v3/checks/suites/">Check Suites</a></li>
        </ul>
      </li>

      
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/gists/">Gists</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/gists/comments/">Comments</a></li>
        </ul>
      </li>
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/git/">Git Data</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/git/blobs/">Blobs</a></li>
          <li><a href="/v3/git/commits/">Commits</a></li>
          <li><a href="/v3/git/refs/">References</a></li>
          <li><a href="/v3/git/tags/">Tags</a></li>
          <li><a href="/v3/git/trees/">Trees</a></li>
        </ul>
      </li>
      
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/actions/">GitHub Actions</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/apps/">GitHub Apps</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/interactions/">Interactions</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/interactions/orgs/">Organization</a></li>
          <li><a href="/v3/interactions/repos/">Repository</a></li>
        </ul>
      </li>
      

      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/issues/">Issues</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/migrations/">Migrations</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/migrations/orgs/">Organization</a></li>
          <li><a href="/v3/migrations/source_imports/">Source Imports</a></li>
          <li><a href="/v3/migrations/users/">User</a></li>
        </ul>
      </li>
      

      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/misc/">Miscellaneous</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/orgs/">Organizations</a></h3>
        <ul class="js-guides">
          
          <li><a href="/v3/orgs/blocking/">Blocking Users</a></li>
          
          <li><a href="/v3/orgs/members/">Members</a></li>
          <li><a href="/v3/orgs/outside_collaborators/">Outside Collaborators</a></li>
          <li><a href="/v3/orgs/hooks/">Webhooks</a></li>
        </ul>
      </li>
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/projects/">Projects</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/projects/cards/">Cards</a></li>
          <li><a href="/v3/projects/collaborators/">Collaborators</a></li>
          <li><a href="/v3/projects/columns/">Columns</a></li>
        </ul>
      </li>
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/pulls/">Pull Requests</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/pulls/reviews/">Reviews</a></li>
          <li><a href="/v3/pulls/comments/">Review Comments</a></li>
          <li><a href="/v3/pulls/review_requests/">Review Requests</a></li>
        </ul>
      </li>
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/reactions/">Reactions</a></h3>
        <ul class="js-guides">
          <li><a href="/v3/reactions/#list-reactions-for-a-commit-comment">Commit Comment</a></li>
          <li><a href="/v3/reactions/#list-reactions-for-an-issue">Issue</a></li>
          <li><a href="/v3/reactions/#list-reactions-for-an-issue-comment">Issue Comment</a></li>
          <li><a href="/v3/reactions/#list-reactions-for-a-pull-request-review-comment">Pull Request Review Comment</a></li>
          <li><a href="/v3/reactions/#list-reactions-for-a-team-discussion">Team Discussion</a></li>
          <li><a href="/v3/reactions/#list-reactions-for-a-team-discussion-comment">Team Discussion Comment</a></li>
        </ul>
      </li>
      <li class="js-topic">
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/repos/">Repositories</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/search/">Search</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/teams/">Teams</a></h3>
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
        <h3><a href="#" class="js-expand-btn collapsed arrow-btn" data-proofer-ignore></a><a href="/v3/users/">Users</a></h3>
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
  
  <div class="sidebar-module api-status py-4 text-center"><a href="https://status.github.com" class="unknown">API Status</a></div>

  
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
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-3769691-37', 'github.com');
  ga('send', 'pageview');
</script>

</body>
</html>
`

var activityEventsGoFileOriginal = `// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListEvents drinks from the firehose of all public events across GitHub.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events
func (s *ActivityService) ListEvents(ctx context.Context, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := addOptions("events", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListRepositoryEvents lists events for a repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-repository-events
func (s *ActivityService) ListRepositoryEvents(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/events", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// Note that ActivityService.ListIssueEventsForRepository was moved to:
// IssuesService.ListRepositoryEvents.

// ListEventsForRepoNetwork lists public events for a network of repositories.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-for-a-network-of-repositories
func (s *ActivityService) ListEventsForRepoNetwork(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("networks/%v/%v/events", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsForOrganization lists public events for an organization.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-for-an-organization
func (s *ActivityService) ListEventsForOrganization(ctx context.Context, org string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("orgs/%v/events", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsPerformedByUser lists the events performed by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-events-for-the-authenticated-user
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-for-a-user
func (s *ActivityService) ListEventsPerformedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	var u string
	if publicOnly {
		u = fmt.Sprintf("users/%v/events/public", user)
	} else {
		u = fmt.Sprintf("users/%v/events", user)
	}
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsReceivedByUser lists the events received by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-events-received-by-the-authenticated-user
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-received-by-a-user
func (s *ActivityService) ListEventsReceivedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	var u string
	if publicOnly {
		u = fmt.Sprintf("users/%v/received_events/public", user)
	} else {
		u = fmt.Sprintf("users/%v/received_events", user)
	}
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListUserEventsForOrganization provides the user’s organization dashboard. You
// must be authenticated as the user to view this.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-events-for-an-organization
func (s *ActivityService) ListUserEventsForOrganization(ctx context.Context, org, user string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("users/%v/events/orgs/%v", user, org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}
`

var activityEventsGoFileWant = `// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListEvents drinks from the firehose of all public events across GitHub.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events
func (s *ActivityService) ListEvents(ctx context.Context, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := addOptions("events", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListRepositoryEvents lists events for a repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-repository-events
func (s *ActivityService) ListRepositoryEvents(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/events", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// Note that ActivityService.ListIssueEventsForRepository was moved to:
// IssuesService.ListRepositoryEvents.

// ListEventsForRepoNetwork lists public events for a network of repositories.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-for-a-network-of-repositories
func (s *ActivityService) ListEventsForRepoNetwork(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("networks/%v/%v/events", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsForOrganization lists public events for an organization.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-organization-events
func (s *ActivityService) ListEventsForOrganization(ctx context.Context, org string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("orgs/%v/events", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsPerformedByUser lists the events performed by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-events-for-the-authenticated-user
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-for-a-user
func (s *ActivityService) ListEventsPerformedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	var u string
	if publicOnly {
		u = fmt.Sprintf("users/%v/events/public", user)
	} else {
		u = fmt.Sprintf("users/%v/events", user)
	}
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsReceivedByUser lists the events received by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-events-received-by-the-authenticated-user
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-public-events-received-by-a-user
func (s *ActivityService) ListEventsReceivedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	var u string
	if publicOnly {
		u = fmt.Sprintf("users/%v/received_events/public", user)
	} else {
		u = fmt.Sprintf("users/%v/received_events", user)
	}
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListUserEventsForOrganization provides the user’s organization dashboard. You
// must be authenticated as the user to view this.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/#list-organization-events-for-the-authenticated-user
func (s *ActivityService) ListUserEventsForOrganization(ctx context.Context, org, user string, opts *ListOptions) ([]*Event, *Response, error) {
	u := fmt.Sprintf("users/%v/events/orgs/%v", user, org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}
`
