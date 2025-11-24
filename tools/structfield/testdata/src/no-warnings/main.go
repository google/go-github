// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"time"
)

type Example struct {
	GithubThing string  `json:"github_thing"`      // Should not be flagged
	ID          *string `json:"id,omitempty"`      // Should not be flagged
	Strings     *string `json:"strings,omitempty"` // Should not be flagged
	Ref         *string `json:"$ref,omitempty"`    // Should not be flagged
	Query       string  `json:"q"`                 // Should not be flagged due to exception
}

type JSONFieldType struct {
	ID             *string           `json:"id,omitempty"`              // Should not be flagged
	HookAttributes map[string]string `json:"hook_attributes,omitempty"` // Should not be flagged
	Inputs         json.RawMessage   `json:"inputs,omitempty"`          // Should not be flagged
	Exception      string            `json:"exception,omitempty"`       // Should not be flagged due to exception
	Value          any               `json:"value,omitempty"`

	Page    *string    `url:"page,omitempty"`         // Should not be flagged
	PerPage *int       `url:"per_page,omitempty"`     // Should not be flagged
	Labels  []string   `url:"labels,omitempty,comma"` // Should not be flagged
	Since   *time.Time `url:"since,omitempty"`        // Should not be flagged
	Fields  []int64    `url:"fields,omitempty,comma"` // Should not be flagged
}
