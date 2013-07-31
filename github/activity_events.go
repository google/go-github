// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Event represents a GitHub event.
type Event struct {
	Type       string          `json:"type,omitempty"`
	Public     bool            `json:"public"`
	RawPayload json.RawMessage `json:"payload,omitempty"`
	Repo       *Repository     `json:"repo,omitempty"`
	Actor      *User           `json:"actor,omitempty"`
	Org        *Organization   `json:"org,omitempty"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	ID         string          `json:"id,omitempty"`
}

// Payload returns the parsed event payload. For recognized event types
// (PushEvent), a value of the corresponding struct type will be returned.
func (e *Event) Payload() (payload interface{}) {
	switch e.Type {
	case "PushEvent":
		payload = &PushEvent{}
	}
	if err := json.Unmarshal(e.RawPayload, &payload); err != nil {
		panic(err.Error())
	}
	return payload
}

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: http://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	PushID  int               `json:"push_id,omitempty"`
	Head    string            `json:"head,omitempty"`
	Ref     string            `json:"ref,omitempty"`
	Size    int               `json:"ref,omitempty"`
	Commits []PushEventCommit `json:"commits,omitempty"`
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	SHA      string        `json:"sha,omitempty"`
	Message  string        `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      string        `json:"url,omitempty"`
	Distinct bool          `json:"distinct"`
}

// ListEventsPerformedByUser lists the events performed by a user. If publicOnly is
// true, only public events will be returned.
//
// GitHub API docs: http://developer.github.com/v3/activity/events/#list-events-performed-by-a-user
func (s *ActivityService) ListEventsPerformedByUser(user string, publicOnly bool, opt *ListOptions) ([]Event, *Response, error) {
	var u string
	if publicOnly {
		u = fmt.Sprintf("users/%v/events/public", user)
	} else {
		u = fmt.Sprintf("users/%v/events", user)
	}

	if opt != nil {
		params := url.Values{
			"page": []string{strconv.Itoa(opt.Page)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	events := new([]Event)
	resp, err := s.client.Do(req, events)
	return *events, resp, err
}
