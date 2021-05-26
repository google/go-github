// Copyright 2013 The go-github AUTHORS. All rights reserved.
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

func TestActivityService_ListEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListEvents(ctx, opt)
	if err != nil {
		t.Errorf("Activities.ListEvents returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListEvents returned %+v, want %+v", events, want)
	}

	const methodName = "ListEvents"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListEvents(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListRepositoryEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListRepositoryEvents(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Activities.ListRepositoryEvents returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListRepositoryEvents returned %+v, want %+v", events, want)
	}

	const methodName = "ListRepositoryEvents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListRepositoryEvents(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListRepositoryEvents(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListRepositoryEvents_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListRepositoryEvents(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestActivityService_ListIssueEventsForRepository(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListIssueEventsForRepository(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Activities.ListIssueEventsForRepository returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListIssueEventsForRepository returned %+v, want %+v", events, want)
	}

	const methodName = "ListIssueEventsForRepository"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListIssueEventsForRepository(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListIssueEventsForRepository(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListIssueEventsForRepository_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListIssueEventsForRepository(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestActivityService_ListEventsForRepoNetwork(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/networks/o/r/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListEventsForRepoNetwork(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Activities.ListEventsForRepoNetwork returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListEventsForRepoNetwork returned %+v, want %+v", events, want)
	}

	const methodName = "ListEventsForRepoNetwork"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListEventsForRepoNetwork(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListEventsForRepoNetwork(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListEventsForRepoNetwork_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListEventsForRepoNetwork(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestActivityService_ListEventsForOrganization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListEventsForOrganization(ctx, "o", opt)
	if err != nil {
		t.Errorf("Activities.ListEventsForOrganization returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListEventsForOrganization returned %+v, want %+v", events, want)
	}

	const methodName = "ListEventsForOrganization"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListEventsForOrganization(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListEventsForOrganization(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListEventsForOrganization_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListEventsForOrganization(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestActivityService_ListEventsPerformedByUser_all(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, "u", false, opt)
	if err != nil {
		t.Errorf("Events.ListPerformedByUser returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Events.ListPerformedByUser returned %+v, want %+v", events, want)
	}

	const methodName = "ListEventsPerformedByUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListEventsPerformedByUser(ctx, "\n", false, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListEventsPerformedByUser(ctx, "u", false, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListEventsPerformedByUser_publicOnly(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/events/public", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	ctx := context.Background()
	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, "u", true, nil)
	if err != nil {
		t.Errorf("Events.ListPerformedByUser returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Events.ListPerformedByUser returned %+v, want %+v", events, want)
	}
}

func TestActivityService_ListEventsPerformedByUser_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListEventsPerformedByUser(ctx, "%", false, nil)
	testURLParseError(t, err)
}

func TestActivityService_ListEventsReceivedByUser_all(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/received_events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListEventsReceivedByUser(ctx, "u", false, opt)
	if err != nil {
		t.Errorf("Events.ListReceivedByUser returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Events.ListReceivedUser returned %+v, want %+v", events, want)
	}

	const methodName = "ListEventsReceivedByUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListEventsReceivedByUser(ctx, "\n", false, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListEventsReceivedByUser(ctx, "u", false, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_ListEventsReceivedByUser_publicOnly(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/received_events/public", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	ctx := context.Background()
	events, _, err := client.Activity.ListEventsReceivedByUser(ctx, "u", true, nil)
	if err != nil {
		t.Errorf("Events.ListReceivedByUser returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Events.ListReceivedByUser returned %+v, want %+v", events, want)
	}
}

func TestActivityService_ListEventsReceivedByUser_invalidUser(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Activity.ListEventsReceivedByUser(ctx, "%", false, nil)
	testURLParseError(t, err)
}

func TestActivityService_ListUserEventsForOrganization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/events/orgs/o", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	events, _, err := client.Activity.ListUserEventsForOrganization(ctx, "o", "u", opt)
	if err != nil {
		t.Errorf("Activities.ListUserEventsForOrganization returned error: %v", err)
	}

	want := []*Event{{ID: String("1")}, {ID: String("2")}}
	if !cmp.Equal(events, want) {
		t.Errorf("Activities.ListUserEventsForOrganization returned %+v, want %+v", events, want)
	}

	const methodName = "ListUserEventsForOrganization"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Activity.ListUserEventsForOrganization(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListUserEventsForOrganization(ctx, "o", "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActivityService_EventParsePayload_typed(t *testing.T) {
	raw := []byte(`{"type": "PushEvent","payload":{"push_id": 1}}`)
	var event *Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("Unmarshal Event returned error: %v", err)
	}

	want := &PushEvent{PushID: Int64(1)}
	got, err := event.ParsePayload()
	if err != nil {
		t.Fatalf("ParsePayload returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Event.ParsePayload returned %+v, want %+v", got, want)
	}
}

// TestEvent_Payload_untyped checks that unrecognized events are parsed to an
// interface{} value (instead of being discarded or throwing an error), for
// forward compatibility with new event types.
func TestActivityService_EventParsePayload_untyped(t *testing.T) {
	raw := []byte(`{"type": "UnrecognizedEvent","payload":{"field": "val"}}`)
	var event *Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("Unmarshal Event returned error: %v", err)
	}

	want := map[string]interface{}{"field": "val"}
	got, err := event.ParsePayload()
	if err != nil {
		t.Fatalf("ParsePayload returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Event.ParsePayload returned %+v, want %+v", got, want)
	}
}

func TestActivityService_EventParsePayload_installation(t *testing.T) {
	raw := []byte(`{"type": "PullRequestEvent","payload":{"installation":{"id":1}}}`)
	var event *Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("Unmarshal Event returned error: %v", err)
	}

	want := &PullRequestEvent{Installation: &Installation{ID: Int64(1)}}
	got, err := event.ParsePayload()
	if err != nil {
		t.Fatalf("ParsePayload returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Event.ParsePayload returned %+v, want %+v", got, want)
	}
}
