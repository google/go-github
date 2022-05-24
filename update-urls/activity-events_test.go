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

func newActivitiesEventsPipeline() *pipelineSetup {
	return &pipelineSetup{
		baseURL:              "https://docs.github.com/en/rest/activity/events/",
		endpointsFromWebsite: activityEventsWant,
		filename:             "activity_events.go",
		serviceName:          "ActivityService",
		originalGoSource:     strings.ReplaceAll(activityEventsGoFileOriginal, "\r", ""),
		wantGoSource:         strings.ReplaceAll(activityEventsGoFileWant, "\r", ""),
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
	got, err := parseWebPageEndpoints(activityEventsTestWebPage)
	if err != nil {
		t.Fatal(err)
	}
	testWebPageHelper(t, got, activityEventsWant)
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

	// Updated docs - consolidated into single page.

	"delete-a-thread-subscription": []*Endpoint{
		{urlFormats: []string{"notifications/threads/%v/subscription"}, httpMethod: "DELETE"},
	},

	"mark-notifications-as-read": []*Endpoint{
		{urlFormats: []string{"notifications"}, httpMethod: "PUT"},
	},

	"set-a-thread-subscription": []*Endpoint{
		{urlFormats: []string{"notifications/threads/%v/subscription"}, httpMethod: "PUT"},
	},

	"delete-a-repository-subscription": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/subscription"}, httpMethod: "DELETE"},
	},

	"star-a-repository-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"user/starred/%v/%v"}, httpMethod: "PUT"},
	},

	"list-repositories-starred-by-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"user/starred"}, httpMethod: "GET"},
	},

	"list-watchers": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/subscribers"}, httpMethod: "GET"},
	},

	"get-feeds": []*Endpoint{
		{urlFormats: []string{"feeds"}, httpMethod: "GET"},
	},

	"get-a-thread": []*Endpoint{
		{urlFormats: []string{"notifications/threads/%v"}, httpMethod: "GET"},
	},

	"mark-a-thread-as-read": []*Endpoint{
		{urlFormats: []string{"notifications/threads/%v"}, httpMethod: "PATCH"},
	},

	"list-stargazers": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/stargazers"}, httpMethod: "GET"},
	},

	"list-repositories-watched-by-a-user": []*Endpoint{
		{urlFormats: []string{"users/%v/subscriptions"}, httpMethod: "GET"},
	},

	"list-repository-notifications-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/notifications"}, httpMethod: "GET"},
	},

	"mark-repository-notifications-as-read": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/notifications"}, httpMethod: "PUT"},
	},

	"check-if-a-repository-is-starred-by-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"user/starred/%v/%v"}, httpMethod: "GET"},
	},

	"list-notifications-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"notifications"}, httpMethod: "GET"},
	},

	"get-a-thread-subscription-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"notifications/threads/%v/subscription"}, httpMethod: "GET"},
	},

	"unstar-a-repository-for-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"user/starred/%v/%v"}, httpMethod: "DELETE"},
	},

	"list-repositories-watched-by-the-authenticated-user": []*Endpoint{
		{urlFormats: []string{"user/subscriptions"}, httpMethod: "GET"},
	},

	"get-a-repository-subscription": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/subscription"}, httpMethod: "GET"},
	},

	"set-a-repository-subscription": []*Endpoint{
		{urlFormats: []string{"repos/%v/%v/subscription"}, httpMethod: "PUT"},
	},

	"list-repositories-starred-by-a-user": []*Endpoint{
		{urlFormats: []string{"users/%v/starred"}, httpMethod: "GET"},
	},
}

//go:embed testdata/activity-events.html
var activityEventsTestWebPage string

//go:embed testdata/activity_events-original.go
var activityEventsGoFileOriginal string

//go:embed testdata/activity_events-want.go
var activityEventsGoFileWant string
