// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAppsService_GetHookConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app/hook/config", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"content_type": "json",
			"insecure_ssl": "0",
			"secret": "********",
			"url": "https://example.com/webhook"
		}`)
	})

	ctx := context.Background()
	config, _, err := client.Apps.GetHookConfig(ctx)
	if err != nil {
		t.Errorf("Apps.GetHookConfig returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: String("json"),
		InsecureSSL: String("0"),
		Secret:      String("********"),
		URL:         String("https://example.com/webhook"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Apps.GetHookConfig returned %+v, want %+v", config, want)
	}

	const methodName = "GetHookConfig"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.GetHookConfig(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_UpdateHookConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &HookConfig{
		ContentType: String("json"),
		InsecureSSL: String("1"),
		Secret:      String("s"),
		URL:         String("u"),
	}

	mux.HandleFunc("/app/hook/config", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"content_type":"json","insecure_ssl":"1","url":"u","secret":"s"}`+"\n")
		fmt.Fprint(w, `{
			"content_type": "json",
			"insecure_ssl": "1",
			"secret": "********",
			"url": "u"
		}`)
	})

	ctx := context.Background()
	config, _, err := client.Apps.UpdateHookConfig(ctx, input)
	if err != nil {
		t.Errorf("Apps.UpdateHookConfig returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: String("json"),
		InsecureSSL: String("1"),
		Secret:      String("********"),
		URL:         String("u"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Apps.UpdateHookConfig returned %+v, want %+v", config, want)
	}

	const methodName = "UpdateHookConfig"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.UpdateHookConfig(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
