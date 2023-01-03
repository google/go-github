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

func TestActionsService_GetOrganizationOIDCSubjectClaimCustomizationTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"include_claim_keys":["repo","context"]}`)
	})

	ctx := context.Background()
	template, _, err := client.Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate returned error: %v", err)
	}

	want := &OIDCSubjectClaimCustomizationTemplate{IncludeClaimKeys: []string{"repo", "context"}}
	if !cmp.Equal(template, want) {
		t.Errorf("Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate returned %+v, want %+v", template, want)
	}

	const methodName = "GetOrganizationOIDCSubjectClaimCustomizationTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepositoryOIDCSubjectClaimCustomizationTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"use_default":false,"include_claim_keys":["repo","context"]}`)
	})

	ctx := context.Background()
	template, _, err := client.Actions.GetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetRepositoryOIDCSubjectClaimCustomizationTemplate returned error: %v", err)
	}

	want := &OIDCSubjectClaimCustomizationTemplate{UseDefault: Bool(false), IncludeClaimKeys: []string{"repo", "context"}}
	if !cmp.Equal(template, want) {
		t.Errorf("Actions.GetOrganizationOIDCSubjectClaimCustomizationTemplate returned %+v, want %+v", template, want)
	}

	const methodName = "GetRepositoryOIDCSubjectClaimCustomizationTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetOrganizationOIDCSubjectClaimCustomizationTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"include_claim_keys":["repo","context"]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &OIDCSubjectClaimCustomizationTemplate{
		IncludeClaimKeys: []string{"repo", "context"},
	}
	ctx := context.Background()
	_, err := client.Actions.SetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.SetOrganizationOIDCSubjectClaimCustomizationTemplate returned error: %v", err)
	}

	const methodName = "SetOrganizationOIDCSubjectClaimCustomizationTemplate"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetOrganizationOIDCSubjectClaimCustomizationTemplate(ctx, "o", input)
	})
}

func TestActionsService_SetRepositoryOIDCSubjectClaimCustomizationTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/oidc/customization/sub", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"use_default":false,"include_claim_keys":["repo","context"]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &OIDCSubjectClaimCustomizationTemplate{
		UseDefault:       Bool(false),
		IncludeClaimKeys: []string{"repo", "context"},
	}
	ctx := context.Background()
	_, err := client.Actions.SetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.SetRepositoryOIDCSubjectClaimCustomizationTemplate returned error: %v", err)
	}

	const methodName = "SetRepositoryOIDCSubjectClaimCustomizationTemplate"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetRepositoryOIDCSubjectClaimCustomizationTemplate(ctx, "o", "r", input)
	})
}
