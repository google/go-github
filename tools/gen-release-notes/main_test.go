// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/compare-vXX.html
var compareVXXHTML string

//go:embed testdata/release-notes-vXY.txt
var releaseNotes string

func TestGenReleaseNotes(t *testing.T) {
	t.Parallel()
	text := strings.ReplaceAll(compareVXXHTML, "\r\n", "\n")
	got := genReleaseNotes(text)
	got = strings.ReplaceAll(got, "\r\n", "\n")
	want := strings.ReplaceAll(releaseNotes, "\r\n", "\n")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(got)
		t.Errorf("genReleaseNotes mismatch (-want +got):\n%v", diff)
	}
}

func TestSplitIntoPRs(t *testing.T) {
	t.Parallel()

	text := strings.ReplaceAll(compareVXXHTML, "\r\n", "\n")
	text = text[191600:]

	got := splitIntoPRs(text)
	want := []string{
		`* Bump go-github from v75 to v76 in /scrape (#3783)`,
		`* Add custom jsonfieldname linter to ensure Go field name matches JSON tag name (#3757)`,
		`* chore: Fix typo in comment (#3786)`,
		`* feat: Add support for private registries endpoints (#3785)`,
		"* Only set `Authorization` when `token` is available (#3789)",
		`* test: Ensure Authorization is not set with empty token (#3790)`,
		`* Fix spelling issues (#3792)`,
		"* refactor!: Remove pointer from required field of CreateStatus API (#3794)\n  BREAKING CHANGE: `RepositoriesService.CreateStatus` now takes value for `status`, not pointer.",
		`* Add test cases for JSON resource marshaling - SCIM (#3798)`,
		`* fix: Org/Enterprise UpdateRepositoryRulesetClearBypassActor sends empty array (#3796)`,
		"* feat!: Address post-merge enterprise billing cost center review (#3805)\n  BREAKING CHANGES: Various `EnterpriseService` structs have been renamed for consistency.",
		`* feat!: Add support for project items CRUD and project fields read operations (#3793)`,
	}

	if len(got) != len(want) {
		t.Log(strings.Join(got, "\n"))
		t.Fatalf("splitIntoPRs = %v lines, want %v", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("splitIntoPRs[%v] =\n%v\n, want \n%v", i, got[i], want[i])
		}
	}
}

func TestMatchDivs(t *testing.T) {
	t.Parallel()

	text := `<div class="flex-auto min-width-0 js-details-container Details">
    <p class="mb-1" >
        <a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b6248e6f6aec019e75ba2c8e189bfe89f36b7d01">Bump go-github from v75 to v76 in /scrape (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3514191578" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3783" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3783/hovercard" href="https://github.com/google/go-github/pull/3783">#3783</a><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b6248e6f6aec019e75ba2c8e189bfe89f36b7d01">)</a>

    </p>


    <div class="d-flex flex-items-center mt-1" >

<div class="AvatarStack flex-self-start  " >
  <div class="AvatarStack-body" >
      <a class="avatar avatar-user" style="width:20px;height:20px;" data-test-selector="commits-avatar-stack-avatar-link" data-hovercard-type="user" data-hovercard-url="/users/gmlewis/hovercard" data-octo-click="hovercard-link-click" data-octo-dimensions="link_type:self" href="/gmlewis">
        <img data-test-selector="commits-avatar-stack-avatar-image" src="https://avatars.githubusercontent.com/u/6598971?s=40&amp;v=4" width="20" height="20" alt="@gmlewis" class=" avatar-user" />
</a>  </div>
</div>

      <div class="f6 color-fg-muted min-width-0">
            <a class="commit-author user-mention" title="View all commits by gmlewis" data-hovercard-type="user" data-hovercard-url="/users/gmlewis/hovercard" data-octo-click="hovercard-link-click" data-octo-dimensions="link_type:self" href="/google/go-github/commits?author=gmlewis">gmlewis</a>

  authored
  <relative-time datetime="2025-10-14T14:29:25Z" class="no-wrap">Oct 14, 2025</relative-time>

      </div>
      <div class="ml-1">

<batch-deferred-content class="d-inline-block" data-url="/google/go-github/commits/checks-statuses-rollups">
    <input type="hidden" name="oid" value="b6248e6f6aec019e75ba2c8e189bfe89f36b7d01" data-targets="batch-deferred-content.inputs" autocomplete="off" />
    <input type="hidden" name="dropdown_direction" value="e" data-targets="batch-deferred-content.inputs" autocomplete="off" />
    <input type="hidden" name="disable_live_updates" value="false" data-targets="batch-deferred-content.inputs" autocomplete="off" />



  <div class="commit-build-statuses">
    <span class="Skeleton d-inline-block" style="width:12px; height:12px;"></span>
  </div>

</batch-deferred-content>
      </div>
    </div>
  </div>

  <div class="d-flex flex-shrink-0" >
<div data-view-component="true" class="Overlay Overlay--size-auto">
      <div data-view-component="true" class="Overlay-body Overlay-body--paddingNone">          <action-list>
  <div data-view-component="true">
    <div data-view-component="true" class="ButtonGroup mr-2 d-none d-sm-flex">
    <div>
<div aria-live="polite" aria-atomic="true" class="sr-only" data-clipboard-copy-feedback></div>
<div>
<div data-view-component="true" class="ButtonGroup d-none d-sm-flex">
    <div>
<div data-view-component="true" class="TimelineItem TimelineItem--condensed pt-2 pb-2">
  <div data-view-component="true" class="TimelineItem-badge"><svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-git-commit">
</svg></div>
  <div data-view-component="true" class="TimelineItem-body">      <h3 class="f5 text-normal color-fg-muted" >Commits on Oct 17, 2025</h3>
`
	want := `* Bump go-github from v75 to v76 in /scrape (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3514191578" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3783" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3783/hovercard" href="https://github.com/google/go-github/pull/3783">#3783</a><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b6248e6f6aec019e75ba2c8e189bfe89f36b7d01">)</a>

    </p>


    <div class="d-flex flex-items-center mt-1" >

`

	got := matchDivs(text)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("matchDivs mismatch (-want +got):\n%v", diff)
	}
}

func TestGetTagSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		text          string
		wantTagSeq    []string
		wantInnerText []string
	}{
		{
			name: "simple case",
			text: `* Bump go-github from v75 to v76 in /scrape (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3514191578" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3783" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3783/hovercard" href="https://github.com/google/go-github/pull/3783">#3783</a><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b6248e6f6aec019e75ba2c8e189bfe89f36b7d01">)</a>

    </p>


    <div class="d-flex flex-items-center mt-1" >`,
			wantTagSeq:    []string{"/a", "a", "/a", "a", "/a", "/p", "div"},
			wantInnerText: []string{"#3783", ")"},
		},
		{
			name: "has ellipses",
			text: `* Add custom jsonfieldname linter to ensure Go field name matches JSON …</a>

        <span class="hidden-text-expander inline">
            <button aria-expanded="false" aria-label="Commit message body" type="button" data-view-component="true" class="ellipsis-expander js-details-target btn">    &hellip;
</button>        </span>
    </p>

      <div class="my-2 Details-content--hidden" ><pre class="text-small ws-pre-wrap">…tag name (<a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3485416274" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3757" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3757/hovercard" href="https://github.com/google/go-github/pull/3757">#3757</a>)</pre>`,
			wantTagSeq:    []string{"/a", "span", "button", "/button", "/span", "/p", "div", "pre", "a", "/a", "/pre"},
			wantInnerText: []string{"…tag name (", "#3757", ")"},
		},
		{
			name: "has code blocks",
			text: `* Only set</a> <code><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b755d64c365491b0e36528802c4b9e89570ae50e">Authorization</a></code> <a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b755d64c365491b0e36528802c4b9e89570ae50e">when</a> <code><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b755d64c365491b0e36528802c4b9e89570ae50e">token</a></code> <a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b755d64c365491b0e36528802c4b9e89570ae50e">is available (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3545203207" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3789" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3789/hovercard" href="https://github.com/google/go-github/pull/3789">#3789</a><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/b755d64c365491b0e36528802c4b9e89570ae50e">)</a>

    </p>


    <div class="d-flex flex-items-center mt-1" >`,
			wantTagSeq:    []string{"/a", "a", "/a", "a", "/a", "a", "/a", "a", "/a", "a", "/a", "a", "/a", "/p", "div"},
			wantInnerText: []string{" `", "Authorization", "` ", "when", " `", "token", "` ", "is available (", "#3789", ")"},
		},
		{
			name: "breaking change with backticks",
			text: `* refactor!: Remove pointer from required field of CreateStatus API (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3560597146" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3794" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3794/hovercard" href="https://github.com/google/go-github/pull/3794">#3794</a>

        <span class="hidden-text-expander inline">
            <button aria-expanded="false" aria-label="Commit message body" type="button" data-view-component="true" class="ellipsis-expander js-details-target btn">    &hellip;
</button>        </span>
    </p>

      <div class="my-2 Details-content--hidden" ><pre class="text-small ws-pre-wrap"><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3560597146" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3794" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3794/hovercard" href="https://github.com/google/go-github/pull/3794"></a>)

BREAKING CHANGE: ` + "`" + `RepositoriesService.CreateStatus` + "`" + ` now takes value for ` + "`" + `status` + "`" + `, not pointer.</pre>`,
			wantTagSeq:    []string{"/a", "a", "/a", "span", "button", "/button", "/span", "/p", "div", "pre", "a", "/a", "/pre"},
			wantInnerText: []string{"#3794", ")\n\nBREAKING CHANGE: `RepositoriesService.CreateStatus` now takes value for `status`, not pointer."},
		},
		{
			name: "more ellipses",
			text: `* fix: Org/Enterprise UpdateRepositoryRulesetClearBypassActor sends emp…</a>

        <span class="hidden-text-expander inline">
            <button aria-expanded="false" aria-label="Commit message body" type="button" data-view-component="true" class="ellipsis-expander js-details-target btn">    &hellip;
</button>        </span>
    </p>

      <div class="my-2 Details-content--hidden" ><pre class="text-small ws-pre-wrap">…ty array (<a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3571786078" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3796" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3796/hovercard" href="https://github.com/google/go-github/pull/3796">#3796</a>)</pre>"`,
			wantTagSeq:    []string{"/a", "span", "button", "/button", "/span", "/p", "div", "pre", "a", "/a", "/pre"},
			wantInnerText: []string{"…ty array (", "#3796", ")"},
		},
		{
			name: "another breaking change",
			text: `* feat!: Add support for project items CRUD and project fields read ope…</a>

        <span class="hidden-text-expander inline">
            <button aria-expanded="false" aria-label="Commit message body" type="button" data-view-component="true" class="ellipsis-expander js-details-target btn">    &hellip;
</button>        </span>
    </p>

      <div class="my-2 Details-content--hidden" ><pre class="text-small ws-pre-wrap">…rations (<a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3558743250" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3793" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3793/hovercard" href="https://github.com/google/go-github/pull/3793">#3793</a>)</pre>`,
			wantTagSeq:    []string{"/a", "span", "button", "/button", "/span", "/p", "div", "pre", "a", "/a", "/pre"},
			wantInnerText: []string{"…rations (", "#3793", ")"},
		},
		{
			name: "bug: missing newline",
			text: `* feat!: Address post-merge enterprise billing cost center review (</a><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="3591077447" data-permission-text="Title is private" data-url="https://github.com/google/go-github/issues/3805" data-hovercard-type="pull_request" data-hovercard-url="/google/go-github/pull/3805/hovercard" href="https://github.com/google/go-github/pull/3805">#3805</a><a class="Link--primary text-bold js-navigation-open markdown-title" href="/google/go-github/commit/1b0a91c5a79dae9c0a014a93f3e447398ec53fa2">)</a>

        <span class="hidden-text-expander inline">
            <button aria-expanded="false" aria-label="Commit message body" type="button" data-view-component="true" class="ellipsis-expander js-details-target btn">    &hellip;
</button>        </span>
    </p>

      <div class="my-2 Details-content--hidden" ><pre class="text-small ws-pre-wrap">BREAKING CHANGES: Various ` + "`" + `EnterpriseService` + "`" + ` structs have been renamed for consistency.</pre></div>

    <div class="d-flex flex-items-center mt-1" >

`,
			wantTagSeq:    []string{"/a", "a", "/a", "a", "/a", "span", "button", "/button", "/span", "/p", "div", "pre", "/pre", "/div", "div"},
			wantInnerText: []string{"#3805", ")", "\n\nBREAKING CHANGES: Various `EnterpriseService` structs have been renamed for consistency."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotTagSeq, gotInnerText := getTagSequence(tt.text)
			if diff := cmp.Diff(tt.wantTagSeq, gotTagSeq); diff != "" {
				t.Errorf("gotTagSeq=\n%#v,\n wantTagSeq=\n%#v", gotTagSeq, tt.wantTagSeq)
			}
			if diff := cmp.Diff(tt.wantInnerText, gotInnerText); diff != "" {
				t.Errorf("gotInnerText=\n%#v,\n wantInnerText=\n%#v", gotInnerText, tt.wantInnerText)
			}
		})
	}
}
