// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestIssuesService_ListIssueEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectCardDetailsPreview)
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	events, _, err := client.Issues.ListIssueEvents(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListIssueEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Int64(1)}}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("Issues.ListIssueEvents returned %+v, want %+v", events, want)
	}

	// Test addOptions failure
	_, _, err = client.Issues.ListIssueEvents(ctx, "\n", "\n", -1, &ListOptions{})
	if err == nil {
		t.Error("bad options ListIssueEvents err = nil, want error")
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Issues.ListIssueEvents(ctx, "o", "r", 1, nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListIssueEvents = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListIssueEvents resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListIssueEvents err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Issues.ListIssueEvents(ctx, "o", "r", 1, nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListIssueEvents = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListIssueEvents resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListIssueEvents err = nil, want error")
	}
}

func TestIssuesService_ListRepositoryEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	events, _, err := client.Issues.ListRepositoryEvents(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListRepositoryEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Int64(1)}}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("Issues.ListRepositoryEvents returned %+v, want %+v", events, want)
	}

	// Test addOptions failure
	_, _, err = client.Issues.ListRepositoryEvents(ctx, "\n", "\n", &ListOptions{})
	if err == nil {
		t.Error("bad options ListRepositoryEvents err = nil, want error")
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Issues.ListRepositoryEvents(ctx, "o", "r", nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListRepositoryEvents = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListRepositoryEvents resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListRepositoryEvents err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Issues.ListRepositoryEvents(ctx, "o", "r", nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListRepositoryEvents = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListRepositoryEvents resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListRepositoryEvents err = nil, want error")
	}
}

func TestIssuesService_GetEvent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/events/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	event, _, err := client.Issues.GetEvent(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetEvent returned error: %v", err)
	}

	want := &IssueEvent{ID: Int64(1)}
	if !reflect.DeepEqual(event, want) {
		t.Errorf("Issues.GetEvent returned %+v, want %+v", event, want)
	}

	// Test addOptions failure
	_, _, err = client.Issues.GetEvent(ctx, "\n", "\n", -1)
	if err == nil {
		t.Error("bad options GetEvent err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err := client.Issues.GetEvent(ctx, "o", "r", 1)
	if got != nil {
		t.Errorf("rate.Reset.Time > now GetEvent = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now GetEvent resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now GetEvent err = nil, want error")
	}
}
