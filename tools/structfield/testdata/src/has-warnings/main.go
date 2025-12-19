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

	Count                              int                `json:"count,omitzero"`                                    // want `the "Count" field in struct "JSONFieldType" uses "omitzero" with a primitive type; remove "omitzero", as it is only allowed with structs, maps, and slices`
	Size                               *int               `json:"size,omitzero"`                                     // want `the "Size" field in struct "JSONFieldType" uses "omitzero" with a primitive type; remove "omitzero" and use only "omitempty" for pointer primitive types`
	PointerToSliceOfStringsZero        *[]string          `json:"pointer_to_slice_of_strings_zero,omitzero"`         // want `change the "PointerToSliceOfStringsZero" field type to "\[\]string" in the struct "JSONFieldType"`
	PointerToSliceOfStructsZero        *[]Struct          `json:"pointer_to_slice_of_structs_zero,omitzero"`         // want `change the "PointerToSliceOfStructsZero" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	PointerToSliceOfPointerStructsZero *[]*Struct         `json:"pointer_to_slice_of_pointer_structs_zero,omitzero"` // want `change the "PointerToSliceOfPointerStructsZero" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	PointerSliceInt                    *[]int             `json:"pointer_slice_int,omitzero"`                        // want `change the "PointerSliceInt" field type to "\[\]int" in the struct "JSONFieldType"`
	AnyZero                            any                `json:"any_zero,omitzero"`                                 // want `the "AnyZero" field in struct "JSONFieldType" uses "omitzero"; remove "omitzero", as it is only allowed with structs, maps, and slices`
	StringPointerSlice                 []*string          `json:"string_pointer_slice,omitzero"`                     // want `change the "StringPointerSlice" field type to "\[\]string" in the struct "JSONFieldType"`
	PointerToMapZero                   *map[string]string `json:"pointer__to_map_zero,omitzero"`                     // want `change the "PointerToMapZero" field type to "map\[string\]string" in the struct "JSONFieldType"`
	SliceOfStructsZero                 []Struct           `json:"slice_of_structs_zero,omitzero"`                    // want `change the "SliceOfStructsZero" field type to "\[\]\*Struct" in the struct "JSONFieldType"`
	StructZero                         Struct             `json:"struct_zero,omitzero"`                              // want `change the "StructZero" field type to "\*Struct" in the struct "JSONFieldType"`

	AnyBoth              any     `json:"any_both,omitempty,omitzero"`                // want `the "AnyBoth" field in struct "JSONFieldType" uses "omitzero"; remove "omitzero", as it is only allowed with structs, maps, and slices`
	NonPointerStructBoth Struct  `json:"non_pointer_struct_both,omitempty,omitzero"` // want `change the "NonPointerStructBoth" field type to "\*Struct" in the struct "JSONFieldType"`
	PointerStringBoth    *string `json:"pointer_string_both,omitempty,omitzero"`     // want `the "PointerStringBoth" field in struct "JSONFieldType" uses "omitzero" with a primitive type; remove "omitzero" and use only "omitempty" for pointer primitive types`
	StructZeroBoth       Struct  `json:"struct_zero_both,omitempty,omitzero"`        // want `change the "StructZeroBoth" field type to "\*Struct" in the struct "JSONFieldType"`
}

type Struct struct{}

type URLFieldName struct {
	GitHubThing string `url:"github_thing"` // want `change Go field name "GitHubThing" to "GithubThing" for tag "github_thing" in struct "URLFieldName"`
}

type URLFieldType struct {
	Page          string `url:"page,omitempty"`          // want `change the "Page" field type to "\*string" in the struct "URLFieldType" because its tag uses "omitempty"`
	PerPage       int    `url:"per_page,omitempty"`      // want `change the "PerPage" field type to "\*int" in the struct "URLFieldType" because its tag uses "omitempty"`
	Participating bool   `url:"participating,omitempty"` // want `change the "Participating" field type to "\*bool" in the struct "URLFieldType" because its tag uses "omitempty"`

	PerPageZeros []int `url:"per_page_zeros,omitzero"`          // want `the "PerPageZeros" field in struct "URLFieldType" uses unsupported omitzero tag for URL tags`
	PerPageBoth  *int  `url:"per_page_both,omitempty,omitzero"` // want `the "PerPageBoth" field in struct "URLFieldType" uses unsupported omitzero tag for URL tags`
}
