// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// CustomProperty represents the organization custom property object.
type CustomProperty struct {
	PropertyName *string `json:"property_name,omitempty"`
	// Possible values for ValueType are: string, single_select
	ValueType     string   `json:"value_type"`
	Required      *bool    `json:"required,omitempty"`
	DefaultValue  *string  `json:"default_value,omitempty"`
	Description   *string  `json:"description,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
}
