// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Example struct {
	GitHubThing      string `json:"github_thing"`               // want `change Go field name "GitHubThing" to "GithubThing" for JSON tag "github_thing" in struct "Example"`
	Id               string `json:"id,omitempty"`               // want `change Go field name "Id" to "ID" for JSON tag "id" in struct "Example"`
	strings          string `json:"strings,omitempty"`          // want `change Go field name "strings" to "Strings" for JSON tag "strings" in struct "Example"`
	camelcaseexample *int   `json:"camelCaseExample,omitempty"` // want `change Go field name "camelcaseexample" to "CamelCaseExample" for JSON tag "camelCaseExample" in struct "Example"`
	DollarRef        string `json:"$ref"`                       // want `change Go field name "DollarRef" to "Ref" for JSON tag "\$ref" in struct "Example"`
}
