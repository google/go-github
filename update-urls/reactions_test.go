// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"strings"
	"testing"
)

func newReactionsPipeline() *pipelineSetup {
	return &pipelineSetup{
		baseURL:              "https://docs.github.com/en/rest/reactions/",
		endpointsFromWebsite: reactionsWant,
		filename:             "reactions.go",
		serviceName:          "ReactionsService",
		originalGoSource:     strings.ReplaceAll(reactionsGoFileOriginal, "\r", ""),
		wantGoSource:         strings.ReplaceAll(reactionsGoFileWant, "\r", ""),
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
	got, err := parseWebPageEndpoints(reactionsTestWebPage)
	if err != nil {
		t.Fatal(err)
	}
	testWebPageHelper(t, got, reactionsWant)
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

//go:embed testdata/reactions.html
var reactionsTestWebPage string

//go:embed testdata/reactions-original.go
var reactionsGoFileOriginal string

//go:embed testdata/reactions-want.go
var reactionsGoFileWant string
