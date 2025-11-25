// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"time"
)

type JSONFieldName struct {
	GithubThing string  `json:"github_thing"`
	ID          *string `json:"id,omitempty"`
	Strings     *string `json:"strings,omitempty"`
	Ref         *string `json:"$ref,omitempty"`
	Query       string  `json:"q"`
}

type JSONFieldType struct {
	WithoutTag string

	ID                    *string           `json:"id,omitempty"`
	HookAttributes        map[string]string `json:"hook_attributes,omitempty"`
	Inputs                json.RawMessage   `json:"inputs,omitempty"`
	Exception             string            `json:"exception,omitempty"`
	Value                 any               `json:"value,omitempty"`
	SliceOfPointerStructs []*Struct         `json:"slice_of_pointer_structs,omitempty"`
}

type URLFieldName struct {
	ID    *string `url:"id,omitempty"`
	Query string  `url:"q"`
}

type URLFieldType struct {
	Page      *string    `url:"page,omitempty"`
	PerPage   *int       `url:"per_page,omitempty"`
	Labels    []string   `url:"labels,omitempty,comma"`
	Since     *time.Time `url:"since,omitempty"`
	Fields    []int64    `url:"fields,omitempty,comma"`
	Exception string     `url:"exception,omitempty"`
}

type Struct struct{}
