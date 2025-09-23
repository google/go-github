// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ClassroomService handles communication with the GitHub Classroom related
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/classroom/classroom
type ClassroomService service

// Classroom represents a GitHub Classroom.
type Classroom struct {
	ID           *int64        `json:"id,omitempty"`
	Name         *string       `json:"name,omitempty"`
	Archived     *bool         `json:"archived,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
	URL          *string       `json:"url,omitempty"`
}

func (c Classroom) String() string {
	return Stringify(c)
}

// ClassroomAssignment represents a GitHub Classroom assignment.
type ClassroomAssignment struct {
	ID                          *int64      `json:"id,omitempty"`
	PublicRepo                  *bool       `json:"public_repo,omitempty"`
	Title                       *string     `json:"title,omitempty"`
	Type                        *string     `json:"type,omitempty"`
	InviteLink                  *string     `json:"invite_link,omitempty"`
	InvitationsEnabled          *bool       `json:"invitations_enabled,omitempty"`
	Slug                        *string     `json:"slug,omitempty"`
	StudentsAreRepoAdmins       *bool       `json:"students_are_repo_admins,omitempty"`
	FeedbackPullRequestsEnabled *bool       `json:"feedback_pull_requests_enabled,omitempty"`
	MaxTeams                    *int        `json:"max_teams,omitempty"`
	MaxMembers                  *int        `json:"max_members,omitempty"`
	Editor                      *string     `json:"editor,omitempty"`
	Accepted                    *int        `json:"accepted,omitempty"`
	Submitted                   *int        `json:"submitted,omitempty"`
	Passing                     *int        `json:"passing,omitempty"`
	Language                    *string     `json:"language,omitempty"`
	Deadline                    *Timestamp  `json:"deadline,omitempty"`
	StarterCodeRepository       *Repository `json:"starter_code_repository,omitempty"`
	Classroom                   *Classroom  `json:"classroom,omitempty"`
}

func (a ClassroomAssignment) String() string {
	return Stringify(a)
}

// GetAssignment gets a GitHub Classroom assignment. Assignment will only be
// returned if the current user is an administrator of the GitHub Classroom
// for the assignment.
//
// GitHub API docs: https://docs.github.com/rest/classroom/classroom#get-an-assignment
//
//meta:operation GET /assignments/{assignment_id}
func (s *ClassroomService) GetAssignment(ctx context.Context, assignmentID int64) (*ClassroomAssignment, *Response, error) {
	u := fmt.Sprintf("assignments/%v", assignmentID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	assignment := new(ClassroomAssignment)
	resp, err := s.client.Do(ctx, req, assignment)
	if err != nil {
		return nil, resp, err
	}

	return assignment, resp, nil
}

// GetClassroom gets a GitHub Classroom for the current user. Classroom will only be
// returned if the current user is an administrator of the GitHub Classroom.
//
// GitHub API docs: https://docs.github.com/rest/classroom/classroom#get-a-classroom
//
//meta:operation GET /classrooms/{classroom_id}
func (s *ClassroomService) GetClassroom(ctx context.Context, classroomID int64) (*Classroom, *Response, error) {
	u := fmt.Sprintf("classrooms/%v", classroomID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	classroom := new(Classroom)
	resp, err := s.client.Do(ctx, req, classroom)
	if err != nil {
		return nil, resp, err
	}

	return classroom, resp, nil
}

// ListClassrooms lists GitHub Classrooms for the current user. Classrooms will only be
// returned if the current user is an administrator of one or more GitHub Classrooms.
//
// GitHub API docs: https://docs.github.com/rest/classroom/classroom#list-classrooms
//
//meta:operation GET /classrooms
func (s *ClassroomService) ListClassrooms(ctx context.Context, opts *ListOptions) ([]*Classroom, *Response, error) {
	u, err := addOptions("classrooms", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var classrooms []*Classroom
	resp, err := s.client.Do(ctx, req, &classrooms)
	if err != nil {
		return nil, resp, err
	}

	return classrooms, resp, nil
}

// ListClassroomAssignments lists GitHub Classroom assignments for a classroom. Assignments will only be
// returned if the current user is an administrator of the GitHub Classroom.
//
// GitHub API docs: https://docs.github.com/rest/classroom/classroom#list-assignments-for-a-classroom
//
//meta:operation GET /classrooms/{classroom_id}/assignments
func (s *ClassroomService) ListClassroomAssignments(ctx context.Context, classroomID int64, opts *ListOptions) ([]*ClassroomAssignment, *Response, error) {
	u := fmt.Sprintf("classrooms/%v/assignments", classroomID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var assignments []*ClassroomAssignment
	resp, err := s.client.Do(ctx, req, &assignments)
	if err != nil {
		return nil, resp, err
	}

	return assignments, resp, nil
}
