// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEventsService_ListPerformedByUser_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/events", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	events, err := client.Events.ListPerformedByUser("u", false, nil)
	if err != nil {
		t.Errorf("Events.ListPerformedByUser returned error: %v", err)
	}

	want := []Event{Event{ID: "1"}, Event{ID: "2"}}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("Events.ListPerformedByUser returned %+v, want %+v", events, want)
	}
}

func TestEventsService_ListPerformedByUser_publicOnly(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/events/public", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	})

	events, err := client.Events.ListPerformedByUser("u", true, nil)
	if err != nil {
		t.Errorf("Events.ListPerformedByUser returned error: %v", err)
	}

	want := []Event{Event{ID: "1"}, Event{ID: "2"}}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("Events.ListPerformedByUser returned %+v, want %+v", events, want)
	}
}

func TestEvent_Payload_typed(t *testing.T) {
	raw := []byte(`{"type": "PushEvent","payload":{"push_id": 1}}`)
	var event *Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("Unmarshal Event returned error: %v", err)
	}

	want := &PushEvent{PushID: 1}
	if !reflect.DeepEqual(event.Payload(), want) {
		t.Errorf("Event Payload returned %+v, want %+v", event.Payload(), want)
	}
}

// TestEvent_Payload_untyped checks that unrecognized events are parsed to an
// interface{} value (instead of being discarded or throwing an error), for
// forward compatibility with new event types.
func TestEvent_Payload_untyped(t *testing.T) {
	raw := []byte(`{"type": "UnrecognizedEvent","payload":{"field": "val"}}`)
	var event *Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("Unmarshal Event returned error: %v", err)
	}

	want := map[string]interface{}{"field": "val"}
	if !reflect.DeepEqual(event.Payload(), want) {
		t.Errorf("Event Payload returned %+v, want %+v", event.Payload(), want)
	}
}
