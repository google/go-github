// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "encoding/json"

type Example struct {
	Strings []string `json:"strings,omitempty"` // Should not be flagged
	Things  []*Thing `json:"things,omitempty"`  // Should not be flagged
}

type Thing struct {
}

type OmitemptyExample struct {
	// Fields with omitempty using pointer types - should NOT be flagged
	Name   *string `json:"name,omitempty"`   // Should not be flagged
	Age    *int    `json:"age,omitempty"`    // Should not be flagged
	Count  *int64  `json:"count,omitempty"`  // Should not be flagged
	Active *bool   `json:"active,omitempty"` // Should not be flagged
	Thing  *Thing  `json:"thing,omitempty"`  // Should not be flagged

	// Reference types that can be nil - should NOT be flagged
	Map        map[string]any  `json:"map,omitempty"`         // Should not be flagged (map can be nil)
	RawMessage json.RawMessage `json:"raw_message,omitempty"` // Should not be flagged ([]byte can be nil)
	Interface  any             `json:"interface,omitempty"`   // Should not be flagged (interface can be nil)
	Slice      []string        `json:"slice,omitempty"`       // Should not be flagged (slice can be nil)

	// Fields without omitempty using value types - should NOT be flagged
	Name2 string `json:"name2"` // Should not be flagged
	Age2  int    `json:"age2"`  // Should not be flagged

	// Fields without json tag - should NOT be flagged
	Internal string // Should not be flagged
}

func main() {
	slice1 := []string{} // Should not be flagged
	_ = slice1
	slice2 := []*Thing{} // Should not be flagged
	_ = slice2
}
