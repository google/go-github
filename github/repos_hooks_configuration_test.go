// Copyright 2023 The go-github AUTHORS. All rights reserved.
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

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetHookConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/hooks/1/config", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"content_type": "json", "insecure_ssl": "0", "secret": "********", "url": "https://example.com/webhook"}`)
	})

	ctx := context.Background()
	config, _, err := client.Repositories.GetHookConfiguration(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.GetHookConfiguration returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: Ptr("json"),
		InsecureSSL: Ptr("0"),
		Secret:      Ptr("********"),
		URL:         Ptr("https://example.com/webhook"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Repositories.GetHookConfiguration returned %+v, want %+v", config, want)
	}

	const methodName = "GetHookConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetHookConfiguration(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetHookConfiguration(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetHookConfiguration_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.GetHookConfiguration(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestRepositoriesService_EditHookConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &HookConfig{}

	mux.HandleFunc("/repos/o/r/hooks/1/config", func(w http.ResponseWriter, r *http.Request) {
		v := new(HookConfig)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"content_type": "json", "insecure_ssl": "0", "secret": "********", "url": "https://example.com/webhook"}`)
	})

	ctx := context.Background()
	config, _, err := client.Repositories.EditHookConfiguration(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.EditHookConfiguration returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: Ptr("json"),
		InsecureSSL: Ptr("0"),
		Secret:      Ptr("********"),
		URL:         Ptr("https://example.com/webhook"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Repositories.EditHookConfiguration returned %+v, want %+v", config, want)
	}

	const methodName = "EditHookConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditHookConfiguration(ctx, "\n", "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditHookConfiguration(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EditHookConfiguration_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Repositories.EditHookConfiguration(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}
