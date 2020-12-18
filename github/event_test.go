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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Payload did not panic but should have")
		}
	}()

	name := "UserEvent"
	body := json.RawMessage("[") // bogus JSON
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}

func TestPayload_NoPanic(t *testing.T) {
	name := "UserEvent"
	body := json.RawMessage("{}")
	e := &Event{Type: &name, RawPayload: &body}
	e.Payload()
}
