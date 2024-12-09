// Copyright 2023 The go-github AUTHORS. All rights reserved.
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

func TestActionsService_GetOrgOIDCSubjectClaimCustomTemplate(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"include_claim_keys":["repo","context"]}`)
	})

	ctx := context.Background()
	template, _, err := client.Actions.GetOrgOIDCSubjectClaimCustomTemplate(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetOrgOIDCSubjectClaimCustomTemplate returned error: %v", err)
	}

	want := &OIDCSubjectClaimCustomTemplate{IncludeClaimKeys: []string{"repo", "context"}}
	if !cmp.Equal(template, want) {
		t.Errorf("Actions.GetOrgOIDCSubjectClaimCustomTemplate returned %+v, want %+v", template, want)
	}

	const methodName = "GetOrgOIDCSubjectClaimCustomTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrgOIDCSubjectClaimCustomTemplate(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrgOIDCSubjectClaimCustomTemplate(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepoOIDCSubjectClaimCustomTemplate(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"use_default":false,"include_claim_keys":["repo","context"]}`)
	})

	ctx := context.Background()
	template, _, err := client.Actions.GetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepoOIDCSubjectClaimCustomTemplate returned error: %v", err)
	}

	want := &OIDCSubjectClaimCustomTemplate{UseDefault: Ptr(false), IncludeClaimKeys: []string{"repo", "context"}}
	if !cmp.Equal(template, want) {
		t.Errorf("Actions.GetOrgOIDCSubjectClaimCustomTemplate returned %+v, want %+v", template, want)
	}

	const methodName = "GetRepoOIDCSubjectClaimCustomTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepoOIDCSubjectClaimCustomTemplate(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetOrgOIDCSubjectClaimCustomTemplate(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"include_claim_keys":["repo","context"]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &OIDCSubjectClaimCustomTemplate{
		IncludeClaimKeys: []string{"repo", "context"},
	}
	ctx := context.Background()
	_, err := client.Actions.SetOrgOIDCSubjectClaimCustomTemplate(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.SetOrgOIDCSubjectClaimCustomTemplate returned error: %v", err)
	}

	const methodName = "SetOrgOIDCSubjectClaimCustomTemplate"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetOrgOIDCSubjectClaimCustomTemplate(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetOrgOIDCSubjectClaimCustomTemplate(ctx, "o", input)
	})
}

func TestActionsService_SetRepoOIDCSubjectClaimCustomTemplate(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"use_default":false,"include_claim_keys":["repo","context"]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &OIDCSubjectClaimCustomTemplate{
		UseDefault:       Ptr(false),
		IncludeClaimKeys: []string{"repo", "context"},
	}
	ctx := context.Background()
	_, err := client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.SetRepoOIDCSubjectClaimCustomTemplate returned error: %v", err)
	}

	const methodName = "SetRepoOIDCSubjectClaimCustomTemplate"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r", input)
	})
}

func TestActionService_SetRepoOIDCSubjectClaimCustomTemplateToDefault(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"use_default":true}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &OIDCSubjectClaimCustomTemplate{
		UseDefault: Ptr(true),
	}
	ctx := context.Background()
	_, err := client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.SetRepoOIDCSubjectClaimCustomTemplate returned error: %v", err)
	}

	const methodName = "SetRepoOIDCSubjectClaimCustomTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetRepoOIDCSubjectClaimCustomTemplate(ctx, "o", "r", input)
	})
}

func TestOIDCSubjectClaimCustomTemplate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &OIDCSubjectClaimCustomTemplate{}, "{}")

	u := &OIDCSubjectClaimCustomTemplate{
		UseDefault:       Ptr(false),
		IncludeClaimKeys: []string{"s"},
	}

	want := `{
		"use_default": false,
		"include_claim_keys": [
			"s"
		]
	}`

	testJSONMarshal(t, u, want)
}
