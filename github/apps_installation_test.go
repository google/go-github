// Copyright 2016 The go-github AUTHORS. All rights reserved.
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

func TestAppsService_ListRepos(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/installation/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"repositories": [{"id":1}]}`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	repositories, _, err := client.Apps.ListRepos(ctx, opt)
	if err != nil {
		t.Errorf("Apps.ListRepos returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(repositories, want) {
		t.Errorf("Apps.ListRepos returned %+v, want %+v", repositories, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Apps.ListRepos(ctx, nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListRepos = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListRepos resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListRepos err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Apps.ListRepos(ctx, nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListRepos = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListRepos resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListRepos err = nil, want error")
	}
}

func TestAppsService_ListUserRepos(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/installations/1/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"repositories": [{"id":1}]}`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	repositories, _, err := client.Apps.ListUserRepos(ctx, 1, opt)
	if err != nil {
		t.Errorf("Apps.ListUserRepos returned error: %v", err)
	}

	want := []*Repository{{ID: Int64(1)}}
	if !reflect.DeepEqual(repositories, want) {
		t.Errorf("Apps.ListUserRepos returned %+v, want %+v", repositories, want)
	}

	// Test addOptions failure
	_, _, err = client.Apps.ListUserRepos(ctx, -1, &ListOptions{})
	if err == nil {
		t.Error("bad options ListUserRepos err = nil, want error")
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Apps.ListUserRepos(ctx, 1, nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListUserRepos = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListUserRepos resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListUserRepos err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Apps.ListUserRepos(ctx, 1, nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListUserRepos = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListUserRepos resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListUserRepos err = nil, want error")
	}
}

func TestAppsService_AddRepository(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/installations/1/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1,"name":"n","description":"d","owner":{"login":"l"},"license":{"key":"mit"}}`)
	})

	ctx := context.Background()
	repo, _, err := client.Apps.AddRepository(ctx, 1, 1)
	if err != nil {
		t.Errorf("Apps.AddRepository returned error: %v", err)
	}

	want := &Repository{ID: Int64(1), Name: String("n"), Description: String("d"), Owner: &User{Login: String("l")}, License: &License{Key: String("mit")}}
	if !reflect.DeepEqual(repo, want) {
		t.Errorf("AddRepository returned %+v, want %+v", repo, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Apps.AddRepository(ctx, 1, 1)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' AddRepository = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' AddRepository resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' AddRepository err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Apps.AddRepository(ctx, 1, 1)
	if got != nil {
		t.Errorf("rate.Reset.Time > now AddRepository = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now AddRepository resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now AddRepository err = nil, want error")
	}
}

func TestAppsService_RemoveRepository(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/installations/1/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Apps.RemoveRepository(ctx, 1, 1)
	if err != nil {
		t.Errorf("Apps.RemoveRepository returned error: %v", err)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	resp, err := client.Apps.RemoveRepository(ctx, 1, 1)
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' RemoveRepository resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' RemoveRepository err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	resp, err = client.Apps.RemoveRepository(ctx, 1, 1)
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now RemoveRepository resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now RemoveRepository err = nil, want error")
	}
}

func TestAppsService_RevokeInstallationToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/installation/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Apps.RevokeInstallationToken(ctx)
	if err != nil {
		t.Errorf("Apps.RevokeInstallationToken returned error: %v", err)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	resp, err := client.Apps.RevokeInstallationToken(ctx)
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' RevokeInstallationToken resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' RevokeInstallationToken err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	resp, err = client.Apps.RevokeInstallationToken(ctx)
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now RevokeInstallationToken resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now RevokeInstallationToken err = nil, want error")
	}
}
