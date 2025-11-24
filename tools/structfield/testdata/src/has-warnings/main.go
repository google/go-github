// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Example struct {
	GitHubThing      string  `json:"github_thing"`               // want `change Go field name "GitHubThing" to "GithubThing" for tag "github_thing" in struct "Example"`
	Id               *string `json:"id,omitempty"`               // want `change Go field name "Id" to "ID" for tag "id" in struct "Example"`
	strings          *string `json:"strings,omitempty"`          // want `change Go field name "strings" to "Strings" for tag "strings" in struct "Example"`
	camelcaseexample *int    `json:"camelCaseExample,omitempty"` // want `change Go field name "camelcaseexample" to "CamelCaseExample" for tag "camelCaseExample" in struct "Example"`
	DollarRef        string  `json:"$ref"`                       // want `change Go field name "DollarRef" to "Ref" for tag "\$ref" in struct "Example"`
}

type JSONFieldType struct {
	ID string `json:"id,omitempty"` // want `change the \"ID\" field type to "\*string" in the struct "JSONFieldType" because its tag uses "omitempty"`

	Page          string `url:"page,omitempty"`          // want `change the "Page" field type to "\*string" in the struct "JSONFieldType" because its tag uses "omitempty"`
	PerPage       int    `url:"per_page,omitempty"`      // want `change the "PerPage" field type to "\*int" in the struct "JSONFieldType" because its tag uses "omitempty"`
	Participating bool   `url:"participating,omitempty"` // want `change the "Participating" field type to "\*bool" in the struct "JSONFieldType" because its tag uses "omitempty"`
}
