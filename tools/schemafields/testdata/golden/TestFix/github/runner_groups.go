// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// RunnerGroupsService is a stand-in for the real service.
type RunnerGroupsService struct{}

// CreateRunnerGroupRequest creates a runner group.
type CreateRunnerGroupRequest struct {
	// Name is required by the schema, but the tag lets it be omitted.
	Name string `json:"name"`

	// Visibility is optional.
	Visibility *string `json:"visibility,omitempty"`

	// SelectedRepositoryIDs is optional and has no omit option.
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitzero"`

	// ResponseOnly is readOnly in the schema, so it is never sent.
	ResponseOnly string `json:"response_only,omitempty"`
}

// UpdateRunnerGroupRequest updates a runner group.
type UpdateRunnerGroupRequest struct {
	// Name is required and not nullable, but the field is a pointer.
	Name string `json:"name"`
}

// NullableRequest has a property that is required and nullable.
type NullableRequest struct {
	// Name may be null, so a pointer is correct even though it is required.
	Name *string `json:"name"`
}

// ValueTypeRequest has an optional property of a value type.
type ValueTypeRequest struct {
	// Count is optional, but a value type is always sent.
	Count *int `json:"count,omitempty"`
}

// StructTypeRequest has an optional property of a struct value type.
type StructTypeRequest struct {
	// Inner is optional and CONTRIBUTING.md asks for omitzero on structs.
	Inner InnerConfig `json:"inner,omitzero"`
}

// InnerConfig is a nested struct.
type InnerConfig struct {
	// Enabled is optional.
	Enabled *bool `json:"enabled,omitempty"`
}

// AllOfRequest takes its schema from an allOf composition.
type AllOfRequest struct {
	// Name is required by the referenced schema.
	Name string `json:"name"`

	// Extra is optional.
	Extra *string `json:"extra,omitempty"`
}

// OneOfRequest takes its schema from a oneOf composition.
type OneOfRequest struct {
	// Kind is required by every variant.
	Kind string `json:"kind"`

	// KindTag is required by every variant, but the tag lets it be omitted.
	KindTag string `json:"kind_tag"`

	// Size is required by only one variant, so it stays optional.
	Size *int `json:"size,omitempty"`
}

// MissingRequest cannot supply a property that the schema requires.
type MissingRequest struct {
	// Name is optional.
	Name *string `json:"name,omitempty"`
}

//meta:operation POST /orgs/{org}/actions/runner-groups
func (s *RunnerGroupsService) CreateRunnerGroup(body CreateRunnerGroupRequest) error { return nil }

//meta:operation PATCH /orgs/{org}/actions/runner-groups/{group_id}
func (s *RunnerGroupsService) UpdateRunnerGroup(body UpdateRunnerGroupRequest) error { return nil }

//meta:operation POST /orgs/{org}/nullable
func (s *RunnerGroupsService) Nullable(body NullableRequest) error { return nil }

//meta:operation POST /orgs/{org}/value-type
func (s *RunnerGroupsService) ValueType(body ValueTypeRequest) error { return nil }

//meta:operation POST /orgs/{org}/struct-type
func (s *RunnerGroupsService) StructType(body StructTypeRequest) error { return nil }

//meta:operation POST /orgs/{org}/allof
func (s *RunnerGroupsService) AllOf(body AllOfRequest) error { return nil }

//meta:operation POST /orgs/{org}/oneof
func (s *RunnerGroupsService) OneOf(body OneOfRequest) error { return nil }

//meta:operation POST /orgs/{org}/missing
func (s *RunnerGroupsService) Missing(body MissingRequest) error { return nil }
