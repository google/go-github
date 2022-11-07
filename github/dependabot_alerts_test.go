// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestDependabotService_ListRepoAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/dependabot/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open"})
		fmt.Fprint(w, `[{"number":1,"state":"open"},{"number":42,"state":"fixed"}]`)
	})

	opts := &ListAlertsOptions{State: String("open")}
	ctx := context.Background()
	alerts, _, err := client.Dependabot.ListRepoAlerts(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Dependabot.ListRepoAlerts returned error: %v", err)
	}

	want := []*DependabotAlert{
		{Number: Int(1), State: String("open")},
		{Number: Int(42), State: String("fixed")},
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("Dependabot.ListRepoAlerts returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListRepoAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListRepoAlerts(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListRepoAlerts(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_GetRepoAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/dependabot/alerts/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":42,"state":"fixed"}`)
	})

	ctx := context.Background()
	alert, _, err := client.Dependabot.GetRepoAlert(ctx, "o", "r", 42)
	if err != nil {
		t.Errorf("Dependabot.GetRepoAlert returned error: %v", err)
	}

	want := &DependabotAlert{
		Number: Int(42),
		State:  String("fixed"),
	}
	if !cmp.Equal(alert, want) {
		t.Errorf("Dependabot.GetRepoAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "GetRepoAlert"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.GetRepoAlert(ctx, "\n", "\n", 0)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.GetRepoAlert(ctx, "o", "r", 42)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependabotService_ListOrgAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/dependabot/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open"})
		fmt.Fprint(w, `[{"number":1,"state":"open"},{"number":42,"state":"fixed"}]`)
	})

	opts := &ListAlertsOptions{State: String("open")}
	ctx := context.Background()
	alerts, _, err := client.Dependabot.ListOrgAlerts(ctx, "o", opts)
	if err != nil {
		t.Errorf("Dependabot.ListOrgAlerts returned error: %v", err)
	}

	want := []*DependabotAlert{
		{Number: Int(1), State: String("open")},
		{Number: Int(42), State: String("fixed")},
	}
	if !cmp.Equal(alerts, want) {
		t.Errorf("Dependabot.ListOrgAlerts returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListOrgAlerts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Dependabot.ListOrgAlerts(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Dependabot.ListOrgAlerts(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
