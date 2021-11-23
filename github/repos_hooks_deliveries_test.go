// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListHookDeliveries(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/deliveries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"cursor": "v1_12077215967"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListCursorOptions{Cursor: "v1_12077215967"}

	ctx := context.Background()
	hooks, _, err := client.Repositories.ListHookDeliveries(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Repositories.ListHookDeliveries returned error: %v", err)
	}

	want := []*HookDelivery{{ID: Int64(1)}, {ID: Int64(2)}}
	if d := cmp.Diff(hooks, want); d != "" {
		t.Errorf("Repositories.ListHooks want (-), got (+):\n%s", d)
	}

	const methodName = "ListHookDeliveries"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListHookDeliveries(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListHookDeliveries(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListHookDeliveries_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.ListHookDeliveries(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestRepositoriesService_GetHookDelivery(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/deliveries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Repositories.GetHookDelivery(ctx, "o", "r", 1, 1)
	if err != nil {
		t.Errorf("Repositories.GetHookDelivery returned error: %v", err)
	}

	want := &HookDelivery{ID: Int64(1)}
	if !cmp.Equal(hook, want) {
		t.Errorf("Repositories.GetHookDelivery returned %+v, want %+v", hook, want)
	}

	const methodName = "GetHookDelivery"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetHookDelivery(ctx, "\n", "\n", -1, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetHookDelivery(ctx, "o", "r", 1, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetHookDelivery_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Repositories.GetHookDelivery(ctx, "%", "%", 1, 1)
	testURLParseError(t, err)
}

var hookDeliveryPayloadTypeToStruct = map[string]interface{}{
	"check_run":                      &CheckRunEvent{},
	"check_suite":                    &CheckSuiteEvent{},
	"commit_comment":                 &CommitCommentEvent{},
	"content_reference":              &ContentReferenceEvent{},
	"create":                         &CreateEvent{},
	"delete":                         &DeleteEvent{},
	"deploy_key":                     &DeployKeyEvent{},
	"deployment":                     &DeploymentEvent{},
	"deployment_status":              &DeploymentStatusEvent{},
	"fork":                           &ForkEvent{},
	"github_app_authorization":       &GitHubAppAuthorizationEvent{},
	"gollum":                         &GollumEvent{},
	"installation":                   &InstallationEvent{},
	"installation_repositories":      &InstallationRepositoriesEvent{},
	"issue_comment":                  &IssueCommentEvent{},
	"issues":                         &IssuesEvent{},
	"label":                          &LabelEvent{},
	"marketplace_purchase":           &MarketplacePurchaseEvent{},
	"member":                         &MemberEvent{},
	"membership":                     &MembershipEvent{},
	"meta":                           &MetaEvent{},
	"milestone":                      &MilestoneEvent{},
	"organization":                   &OrganizationEvent{},
	"org_block":                      &OrgBlockEvent{},
	"package":                        &PackageEvent{},
	"page_build":                     &PageBuildEvent{},
	"ping":                           &PingEvent{},
	"project":                        &ProjectEvent{},
	"project_card":                   &ProjectCardEvent{},
	"project_column":                 &ProjectColumnEvent{},
	"public":                         &PublicEvent{},
	"pull_request":                   &PullRequestEvent{},
	"pull_request_review":            &PullRequestReviewEvent{},
	"pull_request_review_comment":    &PullRequestReviewCommentEvent{},
	"pull_request_target":            &PullRequestTargetEvent{},
	"push":                           &PushEvent{},
	"release":                        &ReleaseEvent{},
	"repository":                     &RepositoryEvent{},
	"repository_dispatch":            &RepositoryDispatchEvent{},
	"repository_vulnerability_alert": &RepositoryVulnerabilityAlertEvent{},
	"star":                           &StarEvent{},
	"status":                         &StatusEvent{},
	"team":                           &TeamEvent{},
	"team_add":                       &TeamAddEvent{},
	"user":                           &UserEvent{},
	"watch":                          &WatchEvent{},
	"workflow_dispatch":              &WorkflowDispatchEvent{},
	"workflow_job":                   &WorkflowJobEvent{},
	"workflow_run":                   &WorkflowRunEvent{},
}

func TestHookDelivery_ParsePayload(t *testing.T) {
	for evt, obj := range hookDeliveryPayloadTypeToStruct {
		t.Run(evt, func(t *testing.T) {
			bs, err := json.Marshal(obj)
			if err != nil {
				t.Fatal(err)
			}

			p := json.RawMessage(bs)

			d := &HookDelivery{
				Event: String(evt),
				Request: &HookRequest{
					RawPayload: &p,
				},
			}

			got, err := d.ParseRequestPayload()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(obj, got) {
				t.Errorf("want %T %v, got %T %v", obj, obj, got, got)
			}
		})
	}
}

func TestHookDelivery_ParsePayload_invalidEvent(t *testing.T) {
	p := json.RawMessage(nil)

	d := &HookDelivery{
		Event: String("some_invalid_event"),
		Request: &HookRequest{
			RawPayload: &p,
		},
	}

	_, err := d.ParseRequestPayload()
	if err == nil || err.Error() != `unsupported event type "some_invalid_event"` {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHookDelivery_ParsePayload_invalidPayload(t *testing.T) {
	p := json.RawMessage([]byte(`{"check_run":{"id":"invalid"}}`))

	d := &HookDelivery{
		Event: String("check_run"),
		Request: &HookRequest{
			RawPayload: &p,
		},
	}

	_, err := d.ParseRequestPayload()
	if err == nil || err.Error() != "json: cannot unmarshal string into Go struct field CheckRun.check_run.id of type int64" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHookRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &HookRequest{}, "{}")

	header := make(map[string]string)
	header["key"] = "value"

	jsonMsg, _ := json.Marshal(&header)

	r := &HookRequest{
		Headers:    header,
		RawPayload: (*json.RawMessage)(&jsonMsg),
	}

	want := `{
		"headers": {
			"key": "value"
		},
		"payload": {
			"key": "value"
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestHookResponse_Marshal(t *testing.T) {
	testJSONMarshal(t, &HookResponse{}, "{}")

	header := make(map[string]string)
	header["key"] = "value"

	jsonMsg, _ := json.Marshal(&header)

	r := &HookResponse{
		Headers:    header,
		RawPayload: (*json.RawMessage)(&jsonMsg),
	}

	want := `{
		"headers": {
			"key": "value"
		},
		"payload": {
			"key": "value"
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestHookDelivery_Marshal(t *testing.T) {
	testJSONMarshal(t, &HookDelivery{}, "{}")

	header := make(map[string]string)
	header["key"] = "value"

	jsonMsg, _ := json.Marshal(&header)

	r := &HookDelivery{
		ID:             Int64(1),
		GUID:           String("guid"),
		DeliveredAt:    &Timestamp{referenceTime},
		Redelivery:     Bool(true),
		Duration:       Float64(1),
		Status:         String("guid"),
		StatusCode:     Int(1),
		Event:          String("guid"),
		Action:         String("guid"),
		InstallationID: String("guid"),
		RepositoryID:   Int64(1),
		Request: &HookRequest{
			Headers:    header,
			RawPayload: (*json.RawMessage)(&jsonMsg),
		},
		Response: &HookResponse{
			Headers:    header,
			RawPayload: (*json.RawMessage)(&jsonMsg),
		},
	}

	want := `{
		"id": 1,
		"guid": "guid",
		"delivered_at": ` + referenceTimeStr + `,
		"redelivery": true,
		"duration": 1,
		"status": "guid",
		"status_code": 1,
		"event": "guid",
		"action": "guid",
		"installation_id": "guid",
		"repository_id": 1,
		"request": {
			"headers": {
				"key": "value"
			},
			"payload": {
				"key": "value"
			}
		},
		"response": {
			"headers": {
				"key": "value"
			},
			"payload": {
				"key": "value"
			}
		}
	}`

	testJSONMarshal(t, r, want)
}
