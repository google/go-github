// Copyright 2020 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_CreateRegistrationToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/actions/runners/registration-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"LLBF3JGZDX3P5PMEXLND6TS6FCWO6","expires_at":"2020-01-22T12:13:35.123Z"}`)
	})

	ctx := context.Background()
	token, _, err := client.Enterprise.CreateRegistrationToken(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.CreateRegistrationToken returned error: %v", err)
	}

	want := &RegistrationToken{Token: String("LLBF3JGZDX3P5PMEXLND6TS6FCWO6"),
		ExpiresAt: &Timestamp{time.Date(2020, time.January, 22, 12, 13, 35,
			123000000, time.UTC)}}
	if !reflect.DeepEqual(token, want) {
		t.Errorf("Enterprise.CreateRegistrationToken returned %+v, want %+v", token, want)
	}

	// Test addOptions failure
	_, _, err = client.Enterprise.CreateRegistrationToken(ctx, "\n")
	if err == nil {
		t.Error("bad options CreateRegistrationToken err = nil, want error")
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Enterprise.CreateRegistrationToken(ctx, "e")
	if got != nil {
		t.Errorf("client.BaseURL.Path='' CreateRegistrationToken = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' CreateRegistrationToken resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' CreateRegistrationToken err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Enterprise.CreateRegistrationToken(ctx, "e")
	if got != nil {
		t.Errorf("rate.Reset.Time > now CreateRegistrationToken = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now CreateRegistrationToken resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now CreateRegistrationToken err = nil, want error")
	}
}

func TestEnterpriseService_ListRunners(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/actions/runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":2,"runners":[{"id":23,"name":"MBP","os":"macos","status":"online"},{"id":24,"name":"iMac","os":"macos","status":"offline"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	runners, _, err := client.Enterprise.ListRunners(ctx, "e", opts)
	if err != nil {
		t.Errorf("Enterprise.ListRunners returned error: %v", err)
	}

	want := &Runners{
		TotalCount: 2,
		Runners: []*Runner{
			{ID: Int64(23), Name: String("MBP"), OS: String("macos"), Status: String("online")},
			{ID: Int64(24), Name: String("iMac"), OS: String("macos"), Status: String("offline")},
		},
	}
	if !reflect.DeepEqual(runners, want) {
		t.Errorf("Actions.ListRunners returned %+v, want %+v", runners, want)
	}

	// Test addOptions failure
	_, _, err = client.Enterprise.ListRunners(ctx, "\n", &ListOptions{})
	if err == nil {
		t.Error("bad options ListRunners err = nil, want error")
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Enterprise.ListRunners(ctx, "e", nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListRunners = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListRunners resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListRunners err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Enterprise.ListRunners(ctx, "e", nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListRunners = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListRunners resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListRunners err = nil, want error")
	}
}
