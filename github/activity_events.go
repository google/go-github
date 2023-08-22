// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// ListEvents drinks from the firehose of all public events across GitHub.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-public-events
func (s *ActivityService) ListEvents(ctx context.Context, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := newURLString("events")
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListRepositoryEvents lists events for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-repository-events
func (s *ActivityService) ListRepositoryEvents(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := newURLString("repos/%v/%v/events", owner, repo)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListIssueEventsForRepository lists issue events for a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/issues/events#list-issue-events-for-a-repository
func (s *ActivityService) ListIssueEventsForRepository(ctx context.Context, owner, repo string, opts *ListOptions) ([]*IssueEvent, *Response, error) {
	u, err := newURLString("repos/%v/%v/issues/events", owner, repo)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*IssueEvent
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsForRepoNetwork lists public events for a network of repositories.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-public-events-for-a-network-of-repositories
func (s *ActivityService) ListEventsForRepoNetwork(ctx context.Context, owner, repo string, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := newURLString("networks/%v/%v/events", owner, repo)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsForOrganization lists public events for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-public-organization-events
func (s *ActivityService) ListEventsForOrganization(ctx context.Context, org string, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := newURLString("orgs/%v/events", org)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsPerformedByUser lists the events performed by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-events-for-the-authenticated-user
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-public-events-for-a-user
func (s *ActivityService) ListEventsPerformedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	urlFormat := "users/%v/events"
	if publicOnly {
		urlFormat += "/public"
	}
	u, err := newURLString(urlFormat, user)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListEventsReceivedByUser lists the events received by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-events-received-by-the-authenticated-user
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-public-events-received-by-a-user
func (s *ActivityService) ListEventsReceivedByUser(ctx context.Context, user string, publicOnly bool, opts *ListOptions) ([]*Event, *Response, error) {
	urlFormat := "users/%v/received_events"
	if publicOnly {
		urlFormat += "/public"
	}
	u, err := newURLString(urlFormat, user)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// ListUserEventsForOrganization provides the user’s organization dashboard. You
// must be authenticated as the user to view this.
//
// GitHub API docs: https://docs.github.com/en/rest/activity/events#list-organization-events-for-the-authenticated-user
func (s *ActivityService) ListUserEventsForOrganization(ctx context.Context, org, user string, opts *ListOptions) ([]*Event, *Response, error) {
	u, err := newURLString("users/%v/events/orgs/%v", user, org)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}
