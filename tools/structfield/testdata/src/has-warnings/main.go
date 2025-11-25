// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type JSONFieldName struct {
	GitHubThing      string  `json:"github_thing"`               // want `change Go field name "GitHubThing" to "GithubThing" for tag "github_thing" in struct "JSONFieldName"`
	Id               *string `json:"id,omitempty"`               // want `change Go field name "Id" to "ID" for tag "id" in struct "JSONFieldName"`
	strings          *string `json:"strings,omitempty"`          // want `change Go field name "strings" to "Strings" for tag "strings" in struct "JSONFieldName"`
	camelcaseexample *int    `json:"camelCaseExample,omitempty"` // want `change Go field name "camelcaseexample" to "CamelCaseExample" for tag "camelCaseExample" in struct "JSONFieldName"`
	DollarRef        string  `json:"$ref"`                       // want `change Go field name "DollarRef" to "Ref" for tag "\$ref" in struct "JSONFieldName"`
}

type JSONFieldType struct {
	String                         string             `json:"string,omitempty"`                              // want `change the "String" field type to "\*string" in the struct "JSONFieldType" because its tag uses "omitempty"`
	SliceOfStringPointers          []*string          `json:"slice_of_string_pointers,omitempty"`            // want `change the "SliceOfStringPointers" field type to "\[\]string" in the struct "JSONFieldType"`
	PointerToSliceOfStrings        *[]string          `json:"pointer_to_slice_of_strings,omitempty"`         // want `change the "PointerToSliceOfStrings" field type to "\[\]string" in the struct "JSONFieldType"`
	SliceOfStructs                 []Struct           `json:"slice_of_structs,omitempty"`                    // want `change the "SliceOfStructs" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	PointerToSliceOfStructs        *[]Struct          `json:"pointer_to_slice_of_structs,omitempty"`         // want `change the "PointerToSliceOfStructs" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	PointerToSliceOfPointerStructs *[]*Struct         `json:"pointer_to_slice_of_pointer_structs,omitempty"` // want `change the "PointerToSliceOfPointerStructs" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	PointerToMap                   *map[string]string `json:"pointer_to_map,omitempty"`                      // want `change the "PointerToMap" field type to "map\[string\]string" in the struct "JSONFieldType"`
	SliceOfInts                    []*int             `json:"slice_of_ints,omitempty"`                       // want `change the "SliceOfInts" field type to "\[\]int" in the struct "JSONFieldType"`
}

type Struct struct{}

type URLFieldName struct {
	GitHubThing string `url:"github_thing"` // want `change Go field name "GitHubThing" to "GithubThing" for tag "github_thing" in struct "URLFieldName"`
}

type URLFieldType struct {
	Page          string `url:"page,omitempty"`          // want `change the "Page" field type to "\*string" in the struct "URLFieldType" because its tag uses "omitempty"`
	PerPage       int    `url:"per_page,omitempty"`      // want `change the "PerPage" field type to "\*int" in the struct "URLFieldType" because its tag uses "omitempty"`
	Participating bool   `url:"participating,omitempty"` // want `change the "Participating" field type to "\*bool" in the struct "URLFieldType" because its tag uses "omitempty"`
}
