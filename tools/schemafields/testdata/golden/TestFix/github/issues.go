// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// IssuesService is a stand-in for the real service.
type IssuesService struct{}

// CreateCommentRequest creates a comment. Every schema agrees that Body is required.
type CreateCommentRequest struct {
	// Body is required.
	Body string `json:"body"`

	// Extra is not in the schema at all.
	Extra *string `json:"extra,omitempty"`
}

// UpdateCommentRequest updates a comment. The schema leaves Body optional.
type UpdateCommentRequest struct {
	// Body is optional.
	Body *string `json:"body,omitempty"`
}

// DualRequest is the body of two operations that disagree about Body.
type DualRequest struct {
	// Body is required when creating and optional when updating.
	Body string `json:"body"`
}

// UntrackedRequest has no operation annotation, so it cannot be mapped to a schema.
type UntrackedRequest struct {
	// Name is a field of an unchecked struct.
	Name string `json:"name"`
}

//meta:operation POST /repos/{owner}/{repo}/issues/{number}/comments
func (s *IssuesService) CreateComment(body CreateCommentRequest) error { return nil }

//meta:operation PATCH /repos/{owner}/{repo}/issues/comments/{comment_id}
func (s *IssuesService) UpdateComment(body UpdateCommentRequest) error { return nil }

//meta:operation POST /repos/{owner}/{repo}/issues/{number}/dual
func (s *IssuesService) CreateDual(body DualRequest) error { return nil }

//meta:operation PATCH /repos/{owner}/{repo}/issues/{number}/dual
func (s *IssuesService) UpdateDual(body DualRequest) error { return nil }

//meta:operation DELETE /repos/{owner}/{repo}/issues/{number}/dual
func (s *IssuesService) DeleteDual(body DualRequest) error { return nil }

// Untracked takes a body whose method has no operation annotation.
func (s *IssuesService) Untracked(body UntrackedRequest) error { return nil }

// Converted takes a body that paramcheck has not converted to a value yet.
func (s *IssuesService) Converted(body *CreateCommentRequest) error { return nil }
