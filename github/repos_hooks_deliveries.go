// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// HookDelivery represents the data that is received from GitHub's Webhook Delivery API
//
// GitHub API docs:
// - https://docs.github.com/en/rest/reference/repos#list-deliveries-for-a-repository-webhook
// - https://docs.github.com/en/rest/reference/repos#get-a-delivery-for-a-repository-webhook
type HookDelivery struct {
	ID             *int64     `json:"id"`
	GUID           *string    `json:"guid"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	Redelivery     *bool      `json:"redelivery"`
	Duration       *float64   `json:"duration"`
	Status         *string    `json:"status"`
	StatusCode     *int       `json:"status_code"`
	Event          *string    `json:"event"`
	Action         *string    `json:"action"`
	InstallationID *string    `json:"installation_id"`
	RepositoryID   *int64     `json:"repository_id"`

	// Request is populated by GetHookDelivery
	Request *HookRequest `json:"request,omitempty"`
	// Response is populated by GetHookDelivery
	Response *HookResponse `json:"response,omitempty"`
}

func (d HookDelivery) String() string {
	return Stringify(d)
}

type HookRequest struct {
	Header     map[string]string `json:"headers,omitempty"`
	RawPayload *json.RawMessage  `json:"payload,omitempty"`
}

func (r HookRequest) String() string {
	return Stringify(r)
}

type HookResponse struct {
	Header     map[string]string `json:"headers,omitempty"`
	RawPayload *json.RawMessage  `json:"payload,omitempty"`
}

func (r HookResponse) String() string {
	return Stringify(r)
}

// ListHookDeliveries lists webhook deliveries for a webhook configured in a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/repos#list-deliveries-for-a-repository-webhook
func (s *RepositoriesService) ListHookDeliveries(ctx context.Context, owner, repo string, id int64, opts *ListCursorOptions) ([]*HookDelivery, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d/deliveries", owner, repo, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deliveries := []*HookDelivery{}
	resp, err := s.client.Do(ctx, req, &deliveries)
	if err != nil {
		return nil, resp, err
	}

	return deliveries, resp, nil
}

// GetHookDelivery returns a delivery for a webhook configured in a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/repos#get-a-delivery-for-a-repository-webhook
func (s *RepositoriesService) GetHookDelivery(ctx context.Context, owner, repo string, hookID, deliveryID int64) (*HookDelivery, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%d/deliveries/%d", owner, repo, hookID, deliveryID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	h := new(HookDelivery)
	resp, err := s.client.Do(ctx, req, h)
	if err != nil {
		return nil, resp, err
	}

	return h, resp, nil
}

func (d *HookDelivery) ParseRequestPayload() (interface{}, error) {
	var payload interface{}
	switch *d.Event {
	case "check_run":
		payload = &CheckRunEvent{}
	case "check_suite":
		payload = &CheckSuiteEvent{}
	case "commit_comment":
		payload = &CommitCommentEvent{}
	case "content_reference":
		payload = &ContentReferenceEvent{}
	case "create":
		payload = &CreateEvent{}
	case "delete":
		payload = &DeleteEvent{}
	case "deploy_ket":
		payload = &DeployKeyEvent{}
	case "deployment":
		payload = &DeploymentEvent{}
	case "deployment_status":
		payload = &DeploymentStatusEvent{}
	case "fork":
		payload = &ForkEvent{}
	case "github_app_authorization":
		payload = &GitHubAppAuthorizationEvent{}
	case "gollum":
		payload = &GollumEvent{}
	case "installation":
		payload = &InstallationEvent{}
	case "installation_repositories":
		payload = &InstallationRepositoriesEvent{}
	case "issue_comment":
		payload = &IssueCommentEvent{}
	case "issues":
		payload = &IssuesEvent{}
	case "label":
		payload = &LabelEvent{}
	case "marketplace_purchase":
		payload = &MarketplacePurchaseEvent{}
	case "member_event":
		payload = &MemberEvent{}
	case "membership_event":
		payload = &MembershipEvent{}
	case "meta":
		payload = &MetaEvent{}
	case "milestone":
		payload = &MilestoneEvent{}
	case "organization":
		payload = &OrganizationEvent{}
	case "org_block":
		payload = &OrgBlockEvent{}
	case "package":
		payload = &PackageEvent{}
	case "page_build":
		payload = &PageBuildEvent{}
	case "ping":
		payload = &PingEvent{}
	case "project":
		payload = &ProjectEvent{}
	case "project_card":
		payload = &ProjectCardEvent{}
	case "project_column":
		payload = &ProjectColumnEvent{}
	case "public":
		payload = &PublicEvent{}
	case "pull_request":
		payload = &PullRequestEvent{}
	case "pull_request_review":
		payload = &PullRequestReviewEvent{}
	case "pull_request_review_comment":
		payload = &PullRequestReviewCommentEvent{}
	case "pull_request_target":
		payload = &PullRequestTargetEvent{}
	case "push":
		payload = &PushEvent{}
	case "release":
		payload = &ReleaseEvent{}
	case "repository":
		payload = &RepositoryEvent{}
	case "repository_dispatch":
		payload = &RepositoryDispatchEvent{}
	case "repository_vulnerability_alert":
		payload = &RepositoryVulnerabilityAlertEvent{}
	case "star":
		payload = &StarEvent{}
	case "status":
		payload = &StatusEvent{}
	case "team":
		payload = &TeamEvent{}
	case "team_add":
		payload = &TeamAddEvent{}
	case "user":
		payload = &UserEvent{}
	case "watch":
		payload = &WatchEvent{}
	case "workflow_dispatch":
		payload = &WorkflowDispatchEvent{}
	case "workflow_run":
		payload = &WorkflowRunEvent{}
	}
	err := json.Unmarshal(*d.Request.RawPayload, &payload)
	return payload, err
}
