// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListCampaigns(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := t.Context()

	mux.HandleFunc("/orgs/o/campaigns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "sort": "created", "direction": "desc", "page": "2"})
		fmt.Fprint(w, `[
		{
			"number": 1,
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-02T00:00:00Z",
			"name": "c1",
			"description": "d1",
			"managers": [{"login": "u"}],
			"ends_at": "2026-02-01T00:00:00Z",
			"state": "open",
			"contact_link": "https://example.com/c1",
			"alert_stats": {"open_count": 1, "closed_count": 2, "in_progress_count": 3}
		},
		{
			"number": 2,
			"created_at": "2026-01-03T00:00:00Z",
			"updated_at": "2026-01-04T00:00:00Z",
			"description": "d2",
			"managers": [],
			"ends_at": "2026-02-02T00:00:00Z",
			"state": "closed",
			"contact_link": "https://example.com/c2"
		}]`)
	})

	opts := &ListCampaignsOptions{State: "open", Sort: "created", Direction: "desc", ListOptions: ListOptions{Page: 2}}
	campaigns, _, err := client.Organizations.ListCampaigns(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListCampaigns returned error: %v", err)
	}

	want := []*Campaign{
		{
			Number:      1,
			CreatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
			UpdatedAt:   Timestamp{time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)},
			Name:        new("c1"),
			Description: "d1",
			Managers:    []*User{{Login: new("u")}},
			EndsAt:      Timestamp{time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)},
			State:       "open",
			ContactLink: "https://example.com/c1",
			AlertStats:  &CampaignAlertStats{OpenCount: 1, ClosedCount: 2, InProgressCount: 3},
		},
		{
			Number:      2,
			CreatedAt:   Timestamp{time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC)},
			UpdatedAt:   Timestamp{time.Date(2026, time.January, 4, 0, 0, 0, 0, time.UTC)},
			Description: "d2",
			Managers:    []*User{},
			EndsAt:      Timestamp{time.Date(2026, time.February, 2, 0, 0, 0, 0, time.UTC)},
			State:       "closed",
			ContactLink: "https://example.com/c2",
		},
	}
	if !cmp.Equal(campaigns, want) {
		t.Errorf("Organizations.ListCampaigns returned %+v, want %+v", campaigns, want)
	}

	const methodName = "ListCampaigns"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCampaigns(ctx, "\n", opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListCampaigns(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateCampaign(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := t.Context()

	input := CreateCampaignRequest{
		Name:        "c1",
		Description: "d1",
		Managers:    []string{"u"},
		EndsAt:      Timestamp{time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)},
		ContactLink: new("https://example.com/c1"),
		CodeScanningAlerts: []*CampaignCodeScanningAlert{
			{RepositoryID: 1, AlertNumbers: []int64{1, 2}},
		},
		GenerateIssues: new(true),
	}

	mux.HandleFunc("/orgs/o/campaigns", func(w http.ResponseWriter, r *http.Request) {
		var v *CreateCampaignRequest
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

		testMethod(t, r, "POST")
		if !cmp.Equal(*v, input) {
			t.Errorf("Request body = %+v, want %+v", *v, input)
		}
		fmt.Fprint(w, `{
			"number": 1,
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-01T00:00:00Z",
			"name": "c1",
			"description": "d1",
			"managers": [{"login": "u"}],
			"ends_at": "2026-02-01T00:00:00Z",
			"state": "open",
			"contact_link": "https://example.com/c1"
		}`)
	})

	campaign, _, err := client.Organizations.CreateCampaign(ctx, "o", input)
	if err != nil {
		t.Errorf("Organizations.CreateCampaign returned error: %v", err)
	}

	want := &Campaign{
		Number:      1,
		CreatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		UpdatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Name:        new("c1"),
		Description: "d1",
		Managers:    []*User{{Login: new("u")}},
		EndsAt:      Timestamp{time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)},
		State:       "open",
		ContactLink: "https://example.com/c1",
	}
	if !cmp.Equal(campaign, want) {
		t.Errorf("Organizations.CreateCampaign returned %+v, want %+v", campaign, want)
	}

	const methodName = "CreateCampaign"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.CreateCampaign(ctx, "\n", input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateCampaign(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetCampaign(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := t.Context()

	mux.HandleFunc("/orgs/o/campaigns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"number": 1,
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-01T00:00:00Z",
			"description": "d1",
			"managers": [{"login": "u"}],
			"ends_at": "2026-02-01T00:00:00Z",
			"state": "open",
			"contact_link": "https://example.com/c1"
		}`)
	})

	campaign, _, err := client.Organizations.GetCampaign(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.GetCampaign returned error: %v", err)
	}

	want := &Campaign{
		Number:      1,
		CreatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		UpdatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Description: "d1",
		Managers:    []*User{{Login: new("u")}},
		EndsAt:      Timestamp{time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)},
		State:       "open",
		ContactLink: "https://example.com/c1",
	}
	if !cmp.Equal(campaign, want) {
		t.Errorf("Organizations.GetCampaign returned %+v, want %+v", campaign, want)
	}

	const methodName = "GetCampaign"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetCampaign(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetCampaign(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateCampaign(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := t.Context()

	input := UpdateCampaignRequest{
		Name:        new("c2"),
		Description: new("d2"),
		State:       new("closed"),
	}

	mux.HandleFunc("/orgs/o/campaigns/1", func(w http.ResponseWriter, r *http.Request) {
		var v *UpdateCampaignRequest
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(*v, input) {
			t.Errorf("Request body = %+v, want %+v", *v, input)
		}
		fmt.Fprint(w, `{
			"number": 1,
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-05T00:00:00Z",
			"name": "c2",
			"description": "d2",
			"managers": [{"login": "u"}],
			"ends_at": "2026-02-01T00:00:00Z",
			"closed_at": "2026-01-05T00:00:00Z",
			"state": "closed",
			"contact_link": "https://example.com/c1"
		}`)
	})

	campaign, _, err := client.Organizations.UpdateCampaign(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.UpdateCampaign returned error: %v", err)
	}

	want := &Campaign{
		Number:      1,
		CreatedAt:   Timestamp{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		UpdatedAt:   Timestamp{time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC)},
		Name:        new("c2"),
		Description: "d2",
		Managers:    []*User{{Login: new("u")}},
		EndsAt:      Timestamp{time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)},
		ClosedAt:    &Timestamp{time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC)},
		State:       "closed",
		ContactLink: "https://example.com/c1",
	}
	if !cmp.Equal(campaign, want) {
		t.Errorf("Organizations.UpdateCampaign returned %+v, want %+v", campaign, want)
	}

	const methodName = "UpdateCampaign"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.UpdateCampaign(ctx, "\n", 1, input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateCampaign(ctx, "o", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeleteCampaign(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	ctx := t.Context()

	mux.HandleFunc("/orgs/o/campaigns/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Organizations.DeleteCampaign(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.DeleteCampaign returned error: %v", err)
	}

	const methodName = "DeleteCampaign"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.DeleteCampaign(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.DeleteCampaign(ctx, "o", 1)
	})
}
