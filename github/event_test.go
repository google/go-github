// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
)

func TestPayload_Panic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Payload did not panic but should have")
		}
	}()

	name := "UserEvent"
	body := json.RawMessage("[") // bogus JSON
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestPayload_NoPanic(t *testing.T) {
	t.Parallel()
	name := "UserEvent"
	body := json.RawMessage("{}")
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestEmptyEvent_NoPanic(t *testing.T) {
	t.Parallel()
	e := &Event{}
	if _, err := e.ParsePayload(); err == nil {
		t.Error("ParsePayload unexpectedly succeeded on empty event")
	}

	e = nil
	if _, err := e.ParsePayload(); err == nil {
		t.Error("ParsePayload unexpectedly succeeded on nil event")
	}
}
