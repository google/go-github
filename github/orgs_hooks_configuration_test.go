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

func TestOrganizationsService_GetHookConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/hooks/1/config", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"content_type": "json", "insecure_ssl": "0", "secret": "********", "url": "https://example.com/webhook"}`)
	})

	ctx := context.Background()
	config, _, err := client.Organizations.GetHookConfiguration(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.GetHookConfiguration returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: Ptr("json"),
		InsecureSSL: Ptr("0"),
		Secret:      Ptr("********"),
		URL:         Ptr("https://example.com/webhook"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Organizations.GetHookConfiguration returned %+v, want %+v", config, want)
	}

	const methodName = "GetHookConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetHookConfiguration(ctx, "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetHookConfiguration(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetHookConfiguration_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.GetHookConfiguration(ctx, "%", 1)
	testURLParseError(t, err)
}

func TestOrganizationsService_EditHookConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &HookConfig{}

	mux.HandleFunc("/orgs/o/hooks/1/config", func(w http.ResponseWriter, r *http.Request) {
		v := new(HookConfig)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"content_type": "json", "insecure_ssl": "0", "secret": "********", "url": "https://example.com/webhook"}`)
	})

	ctx := context.Background()
	config, _, err := client.Organizations.EditHookConfiguration(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.EditHookConfiguration returned error: %v", err)
	}

	want := &HookConfig{
		ContentType: Ptr("json"),
		InsecureSSL: Ptr("0"),
		Secret:      Ptr("********"),
		URL:         Ptr("https://example.com/webhook"),
	}
	if !cmp.Equal(config, want) {
		t.Errorf("Organizations.EditHookConfiguration returned %+v, want %+v", config, want)
	}

	const methodName = "EditHookConfiguration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.EditHookConfiguration(ctx, "\n", -1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.EditHookConfiguration(ctx, "o", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_EditHookConfiguration_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Organizations.EditHookConfiguration(ctx, "%", 1, nil)
	testURLParseError(t, err)
}
