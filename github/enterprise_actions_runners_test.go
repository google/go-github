// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseService_GenerateEnterpriseJITConfig(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &GenerateJITConfigRequest{Name: "test", RunnerGroupID: 1, Labels: []string{"one", "two"}}

	mux.HandleFunc("/enterprises/o/actions/runners/generate-jitconfig", func(w http.ResponseWriter, r *http.Request) {
		v := new(GenerateJITConfigRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Errorf("Request body decode failed: %v", err)
		}

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"encoded_jit_config":"foo"}`)
	})

	ctx := context.Background()
	jitConfig, _, err := client.Enterprise.GenerateEnterpriseJITConfig(ctx, "o", input)
	if err != nil {
		t.Errorf("Enterprise.GenerateEnterpriseJITConfig returned error: %v", err)
	}

	want := &JITRunnerConfig{EncodedJITConfig: String("foo")}
	if !cmp.Equal(jitConfig, want) {
		t.Errorf("Enterprise.GenerateEnterpriseJITConfig returned %+v, want %+v", jitConfig, want)
	}

	const methodName = "GenerateEnterpriseJITConfig"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GenerateEnterpriseJITConfig(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GenerateEnterpriseJITConfig(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateRegistrationToken(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	if !cmp.Equal(token, want) {
		t.Errorf("Enterprise.CreateRegistrationToken returned %+v, want %+v", token, want)
	}

	const methodName = "CreateRegistrationToken"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.CreateRegistrationToken(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateRegistrationToken(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "MBP", "per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":1,"runners":[{"id":23,"name":"MBP","os":"macos","status":"online"}]}`)
	})

	opts := &ListRunnersOptions{
		Name:        String("MBP"),
		ListOptions: ListOptions{Page: 2, PerPage: 2},
	}
	ctx := context.Background()
	runners, _, err := client.Enterprise.ListRunners(ctx, "e", opts)
	if err != nil {
		t.Errorf("Enterprise.ListRunners returned error: %v", err)
	}

	want := &Runners{
		TotalCount: 1,
		Runners: []*Runner{
			{ID: Int64(23), Name: String("MBP"), OS: String("macos"), Status: String("online")},
		},
	}
	if !cmp.Equal(runners, want) {
		t.Errorf("Actions.ListRunners returned %+v, want %+v", runners, want)
	}

	const methodName = "ListRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRunners(ctx, "\n", &ListRunnersOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListRunners(ctx, "e", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/runners/23", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":23,"name":"MBP","os":"macos","status":"online"}`)
	})

	ctx := context.Background()
	runner, _, err := client.Enterprise.GetRunner(ctx, "e", 23)
	if err != nil {
		t.Errorf("Enterprise.GetRunner returned error: %v", err)
	}

	want := &Runner{
		ID:     Int64(23),
		Name:   String("MBP"),
		OS:     String("macos"),
		Status: String("online"),
	}
	if !cmp.Equal(runner, want) {
		t.Errorf("Enterprise.GetRunner returned %+v, want %+v", runner, want)
	}

	const methodName = "GetRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetRunner(ctx, "\n", 23)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetRunner(ctx, "e", 23)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_RemoveRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runners/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.RemoveRunner(ctx, "o", 21)
	if err != nil {
		t.Errorf("Actions.RemoveRunner returned error: %v", err)
	}

	const methodName = "RemoveRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.RemoveRunner(ctx, "\n", 21)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.RemoveRunner(ctx, "o", 21)
	})
}

func TestEnterpriseService_ListRunnerApplicationDownloads(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/o/actions/runners/downloads", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"os":"osx","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-osx-x64-2.164.0.tar.gz","filename":"actions-runner-osx-x64-2.164.0.tar.gz"},{"os":"linux","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-x64-2.164.0.tar.gz","filename":"actions-runner-linux-x64-2.164.0.tar.gz"},{"os": "linux","architecture":"arm","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm-2.164.0.tar.gz","filename":"actions-runner-linux-arm-2.164.0.tar.gz"},{"os":"win","architecture":"x64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-win-x64-2.164.0.zip","filename":"actions-runner-win-x64-2.164.0.zip"},{"os":"linux","architecture":"arm64","download_url":"https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm64-2.164.0.tar.gz","filename":"actions-runner-linux-arm64-2.164.0.tar.gz"}]`)
	})

	ctx := context.Background()
	downloads, _, err := client.Enterprise.ListRunnerApplicationDownloads(ctx, "o")
	if err != nil {
		t.Errorf("Enterprise.ListRunnerApplicationDownloads returned error: %v", err)
	}

	want := []*RunnerApplicationDownload{
		{OS: String("osx"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-osx-x64-2.164.0.tar.gz"), Filename: String("actions-runner-osx-x64-2.164.0.tar.gz")},
		{OS: String("linux"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-x64-2.164.0.tar.gz"), Filename: String("actions-runner-linux-x64-2.164.0.tar.gz")},
		{OS: String("linux"), Architecture: String("arm"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm-2.164.0.tar.gz"), Filename: String("actions-runner-linux-arm-2.164.0.tar.gz")},
		{OS: String("win"), Architecture: String("x64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-win-x64-2.164.0.zip"), Filename: String("actions-runner-win-x64-2.164.0.zip")},
		{OS: String("linux"), Architecture: String("arm64"), DownloadURL: String("https://github.com/actions/runner/releases/download/v2.164.0/actions-runner-linux-arm64-2.164.0.tar.gz"), Filename: String("actions-runner-linux-arm64-2.164.0.tar.gz")},
	}
	if !cmp.Equal(downloads, want) {
		t.Errorf("Enterprise.ListRunnerApplicationDownloads returned %+v, want %+v", downloads, want)
	}

	const methodName = "ListRunnerApplicationDownloads"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListRunnerApplicationDownloads(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListRunnerApplicationDownloads(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
