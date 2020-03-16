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

func TestActionsService_ListRunnerApplicationDownloads(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners/downloads", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"os":"osx","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-osx-x64-2.164.0.tar.gz","filename":"actions-runner-osx-x64-2.164.0.tar.gz"},{"os":"linux","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-x64-2.164.0.tar.gz","filename":"actions-runner-linux-x64-2.164.0.tar.gz"},{"os": "linux","architecture":"arm","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm-2.164.0.tar.gz","filename":"actions-runner-linux-arm-2.164.0.tar.gz"},{"os":"win","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-win-x64-2.164.0.zip","filename":"actions-runner-win-x64-2.164.0.zip"},{"os":"linux","architecture":"arm64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm64-2.164.0.tar.gz","filename":"actions-runner-linux-arm64-2.164.0.tar.gz"}]`)
	})

	downloads, _, err := client.Actions.ListRunnerApplicationDownloads(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Actions.ListRunnerApplicationDownloads returned error: %v", err)
	}

	want := []*RunnerApplicationDownload{
		{OS: String("osx"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-osx-x64-2.164.0.tar.gz"), Filename: String("actions-runner-osx-x64-2.164.0.tar.gz")},
		{OS: String("linux"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-x64-2.164.0.tar.gz"), Filename: String("actions-runner-linux-x64-2.164.0.tar.gz")},
		{OS: String("linux"), Architecture: String("arm"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm-2.164.0.tar.gz"), Filename: String("actions-runner-linux-arm-2.164.0.tar.gz")},
		{OS: String("win"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-win-x64-2.164.0.zip"), Filename: String("actions-runner-win-x64-2.164.0.zip")},
		{OS: String("linux"), Architecture: String("arm64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm64-2.164.0.tar.gz"), Filename: String("actions-runner-linux-arm64-2.164.0.tar.gz")},
	}
	if !reflect.DeepEqual(downloads, want) {
		t.Errorf("Actions.ListRunnerApplicationDownloads returned %+v, want %+v", downloads, want)
	}
}

func TestActionsService_CreateRegistrationToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners/registration-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"LLBF3JGZDX3P5PMEXLND6TS6FCWO6","expires_at":"2020-01-22T12:13:35.123Z"}`)
	})

	token, _, err := client.Actions.CreateRegistrationToken(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Actions.CreateRegistrationToken returned error: %v", err)
	}

	want := &RegistrationToken{Token: String("LLBF3JGZDX3P5PMEXLND6TS6FCWO6"),
		ExpiresAt: &Timestamp{time.Date(2020, time.January, 22, 12, 13, 35,
			123000000, time.UTC)}}
	if !reflect.DeepEqual(token, want) {
		t.Errorf("Actions.CreateRegistrationToken returned %+v, want %+v", token, want)
	}
}

func TestActionsService_ListRunners(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `[{"id":23,"name":"MBP","os":"macos","status":"online"},{"id":24,"name":"iMac","os":"macos","status":"offline"}]`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	runners, _, err := client.Actions.ListRunners(context.Background(), "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRunners returned error: %v", err)
	}

	want := []*Runner{
		{ID: Int64(23), Name: String("MBP"), OS: String("macos"), Status: String("online")},
		{ID: Int64(24), Name: String("iMac"), OS: String("macos"), Status: String("offline")},
	}
	if !reflect.DeepEqual(runners, want) {
		t.Errorf("Actions.ListRunners returned %+v, want %+v", runners, want)
	}
}

func TestActionsService_GetRunner(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners/23", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":23,"name":"MBP","os":"macos","status":"online"}`)
	})

	runner, _, err := client.Actions.GetRunner(context.Background(), "o", "r", 23)
	if err != nil {
		t.Errorf("Actions.GetRunner returned error: %v", err)
	}

	want := &Runner{
		ID:     Int64(23),
		Name:   String("MBP"),
		OS:     String("macos"),
		Status: String("online"),
	}
	if !reflect.DeepEqual(runner, want) {
		t.Errorf("Actions.GetRunner returned %+v, want %+v", runner, want)
	}
}

func TestActionsService_CreateRemoveToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners/remove-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"AABF3JGZDX3P5PMEXLND6TS6FCWO6","expires_at":"2020-01-29T12:13:35.123Z"}`)
	})

	token, _, err := client.Actions.CreateRemoveToken(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Actions.CreateRemoveToken returned error: %v", err)
	}

	want := &RemoveToken{Token: String("AABF3JGZDX3P5PMEXLND6TS6FCWO6"), ExpiresAt: &Timestamp{time.Date(2020, time.January, 29, 12, 13, 35, 123000000, time.UTC)}}
	if !reflect.DeepEqual(token, want) {
		t.Errorf("Actions.CreateRemoveToken returned %+v, want %+v", token, want)
	}
}

func TestActionsService_RemoveRunner(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runners/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Actions.RemoveRunner(context.Background(), "o", "r", 21)
	if err != nil {
		t.Errorf("Actions.RemoveRunner returned error: %v", err)
	}
}
