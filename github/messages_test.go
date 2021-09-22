// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMessageMAC_BadHashTypePrefix(t *testing.T) {
	const signature = "bogus1=1234567"
	if _, _, err := messageMAC(signature); err == nil {
		t.Fatal("messageMAC returned nil; wanted error")
	}
}

func TestValidatePayload(t *testing.T) {
	const defaultBody = `{"yo":true}` // All tests below use the default request body and signature.
	const defaultSignature = "sha1=126f2c800419c60137ce748d7672e77b65cf16d6"
	secretKey := []byte("0123456789abcdef")
	tests := []struct {
		signature       string
		signatureHeader string
		eventID         string
		event           string
		wantEventID     string
		wantEvent       string
		wantPayload     string
	}{
		// The following tests generate expected errors:
		{},                         // Missing signature
		{signature: "yo"},          // Missing signature prefix
		{signature: "sha1=yo"},     // Signature not hex string
		{signature: "sha1=012345"}, // Invalid signature
		// The following tests expect err=nil:
		{
			signature:   defaultSignature,
			eventID:     "dead-beef",
			event:       "ping",
			wantEventID: "dead-beef",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   defaultSignature,
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   "sha256=b1f8020f5b4cd42042f807dd939015c4a418bc1ff7f604dd55b0a19b5d953d9b",
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:       "sha256=b1f8020f5b4cd42042f807dd939015c4a418bc1ff7f604dd55b0a19b5d953d9b",
			signatureHeader: SHA256SignatureHeader,
			event:           "ping",
			wantEvent:       "ping",
			wantPayload:     defaultBody,
		},
		{
			signature:   "sha512=8456767023c1195682e182a23b3f5d19150ecea598fde8cb85918f7281b16079471b1329f92b912c4d8bd7455cb159777db8f29608b20c7c87323ba65ae62e1f",
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
	}

	for _, test := range tests {
		buf := bytes.NewBufferString(defaultBody)
		req, err := http.NewRequest("GET", "http://localhost/event", buf)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		if test.signature != "" {
			if test.signatureHeader != "" {
				req.Header.Set(test.signatureHeader, test.signature)
			} else {
				req.Header.Set(SHA1SignatureHeader, test.signature)
			}
		}
		req.Header.Set("Content-Type", "application/json")

		got, err := ValidatePayload(req, secretKey)
		if err != nil {
			if test.wantPayload != "" {
				t.Errorf("ValidatePayload(%#v): err = %v, want nil", test, err)
			}
			continue
		}
		if string(got) != test.wantPayload {
			t.Errorf("ValidatePayload = %q, want %q", got, test.wantPayload)
		}
	}
}

func TestValidatePayload_FormGet(t *testing.T) {
	payload := `{"yo":true}`
	signature := "sha1=3374ef144403e8035423b23b02e2c9d7a4c50368"
	secretKey := []byte("0123456789abcdef")

	form := url.Values{}
	form.Add("payload", payload)
	req, err := http.NewRequest("POST", "http://localhost/event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(SHA1SignatureHeader, signature)

	got, err := ValidatePayload(req, secretKey)
	if err != nil {
		t.Errorf("ValidatePayload(%#v): err = %v, want nil", payload, err)
	}
	if string(got) != payload {
		t.Errorf("ValidatePayload = %q, want %q", got, payload)
	}

	// check that if payload is invalid we get error
	req.Header.Set(SHA1SignatureHeader, "invalid signature")
	if _, err = ValidatePayload(req, []byte{0}); err == nil {
		t.Error("ValidatePayload = nil, want err")
	}
}

func TestValidatePayload_FormPost(t *testing.T) {
	payload := `{"yo":true}`
	signature := "sha1=3374ef144403e8035423b23b02e2c9d7a4c50368"
	secretKey := []byte("0123456789abcdef")

	form := url.Values{}
	form.Set("payload", payload)
	buf := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequest("POST", "http://localhost/event", buf)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(SHA1SignatureHeader, signature)

	got, err := ValidatePayload(req, secretKey)
	if err != nil {
		t.Errorf("ValidatePayload(%#v): err = %v, want nil", payload, err)
	}
	if string(got) != payload {
		t.Errorf("ValidatePayload = %q, want %q", got, payload)
	}

	// check that if payload is invalid we get error
	req.Header.Set(SHA1SignatureHeader, "invalid signature")
	if _, err = ValidatePayload(req, []byte{0}); err == nil {
		t.Error("ValidatePayload = nil, want err")
	}
}

func TestValidatePayload_InvalidContentType(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost/event", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "invalid content type")
	if _, err = ValidatePayload(req, nil); err == nil {
		t.Error("ValidatePayload = nil, want err")
	}
}

func TestValidatePayload_NoSecretKey(t *testing.T) {
	payload := `{"yo":true}`

	form := url.Values{}
	form.Set("payload", payload)
	buf := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequest("POST", "http://localhost/event", buf)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	got, err := ValidatePayload(req, nil)
	if err != nil {
		t.Errorf("ValidatePayload(%#v): err = %v, want nil", payload, err)
	}
	if string(got) != payload {
		t.Errorf("ValidatePayload = %q, want %q", got, payload)
	}
}

// badReader satisfies io.Reader but always returns an error.
type badReader struct{}

func (b *badReader) Read(p []byte) (int, error) {
	return 0, errors.New("bad reader")
}

func (b *badReader) Close() error { return errors.New("bad reader") }

func TestValidatePayload_BadRequestBody(t *testing.T) {
	tests := []struct {
		contentType string
	}{
		{contentType: "application/json"},
		{contentType: "application/x-www-form-urlencoded"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			req := &http.Request{
				Header: http.Header{"Content-Type": []string{tt.contentType}},
				Body:   &badReader{},
			}
			if _, err := ValidatePayload(req, nil); err == nil {
				t.Fatal("ValidatePayload returned nil; want error")
			}
		})
	}
}

func TestValidatePayload_InvalidContentTypeParams(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost/event", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=")
	if _, err = ValidatePayload(req, nil); err == nil {
		t.Error("ValidatePayload = nil, want err")
	}
}

func TestValidatePayload_ValidContentTypeParams(t *testing.T) {
	var requestBody = `{"yo":true}`
	buf := bytes.NewBufferString(requestBody)

	req, err := http.NewRequest("POST", "http://localhost/event", buf)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	_, err = ValidatePayload(req, nil)
	if err != nil {
		t.Error("ValidatePayload = nil, want err")
	}
}

func TestParseWebHook(t *testing.T) {
	tests := []struct {
		payload     interface{}
		messageType string
	}{
		{
			payload:     &BranchProtectionRuleEvent{},
			messageType: "branch_protection_rule",
		},
		{
			payload:     &CheckRunEvent{},
			messageType: "check_run",
		},
		{
			payload:     &CheckSuiteEvent{},
			messageType: "check_suite",
		},
		{
			payload:     &CommitCommentEvent{},
			messageType: "commit_comment",
		},
		{
			payload:     &ContentReferenceEvent{},
			messageType: "content_reference",
		},
		{
			payload:     &CreateEvent{},
			messageType: "create",
		},
		{
			payload:     &DeleteEvent{},
			messageType: "delete",
		},
		{
			payload:     &DeployKeyEvent{},
			messageType: "deploy_key",
		},
		{
			payload:     &DeploymentEvent{},
			messageType: "deployment",
		},

		{
			payload:     &DeploymentStatusEvent{},
			messageType: "deployment_status",
		},
		{
			payload:     &ForkEvent{},
			messageType: "fork",
		},
		{
			payload:     &GitHubAppAuthorizationEvent{},
			messageType: "github_app_authorization",
		},
		{
			payload:     &GollumEvent{},
			messageType: "gollum",
		},
		{
			payload:     &InstallationEvent{},
			messageType: "installation",
		},
		{
			payload:     &InstallationRepositoriesEvent{},
			messageType: "installation_repositories",
		},
		{
			payload:     &IssueCommentEvent{},
			messageType: "issue_comment",
		},
		{
			payload:     &IssuesEvent{},
			messageType: "issues",
		},
		{
			payload:     &LabelEvent{},
			messageType: "label",
		},
		{
			payload:     &MarketplacePurchaseEvent{},
			messageType: "marketplace_purchase",
		},
		{
			payload:     &MemberEvent{},
			messageType: "member",
		},
		{
			payload:     &MembershipEvent{},
			messageType: "membership",
		},
		{
			payload:     &MetaEvent{},
			messageType: "meta",
		},
		{
			payload:     &MilestoneEvent{},
			messageType: "milestone",
		},
		{
			payload:     &OrganizationEvent{},
			messageType: "organization",
		},
		{
			payload:     &OrgBlockEvent{},
			messageType: "org_block",
		},
		{
			payload:     &PackageEvent{},
			messageType: "package",
		},
		{
			payload:     &PageBuildEvent{},
			messageType: "page_build",
		},
		{
			payload:     &PingEvent{},
			messageType: "ping",
		},
		{
			payload:     &ProjectEvent{},
			messageType: "project",
		},
		{
			payload:     &ProjectCardEvent{},
			messageType: "project_card",
		},
		{
			payload:     &ProjectColumnEvent{},
			messageType: "project_column",
		},
		{
			payload:     &PublicEvent{},
			messageType: "public",
		},
		{
			payload:     &PullRequestEvent{},
			messageType: "pull_request",
		},
		{
			payload:     &PullRequestReviewEvent{},
			messageType: "pull_request_review",
		},
		{
			payload:     &PullRequestReviewCommentEvent{},
			messageType: "pull_request_review_comment",
		},
		{
			payload:     &PullRequestTargetEvent{},
			messageType: "pull_request_target",
		},
		{
			payload:     &PushEvent{},
			messageType: "push",
		},
		{
			payload:     &ReleaseEvent{},
			messageType: "release",
		},
		{
			payload:     &RepositoryEvent{},
			messageType: "repository",
		},
		{
			payload:     &RepositoryVulnerabilityAlertEvent{},
			messageType: "repository_vulnerability_alert",
		},
		{
			payload:     &StarEvent{},
			messageType: "star",
		},
		{
			payload:     &StatusEvent{},
			messageType: "status",
		},
		{
			payload:     &TeamEvent{},
			messageType: "team",
		},
		{
			payload:     &TeamAddEvent{},
			messageType: "team_add",
		},
		{
			payload:     &UserEvent{},
			messageType: "user",
		},
		{
			payload:     &WatchEvent{},
			messageType: "watch",
		},
		{
			payload:     &RepositoryDispatchEvent{},
			messageType: "repository_dispatch",
		},
		{
			payload:     &WorkflowDispatchEvent{},
			messageType: "workflow_dispatch",
		},
		{
			payload:     &WorkflowJobEvent{},
			messageType: "workflow_job",
		},
		{
			payload:     &WorkflowRunEvent{},
			messageType: "workflow_run",
		},
	}

	for _, test := range tests {
		p, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Marshal(%#v): %v", test.payload, err)
		}
		got, err := ParseWebHook(test.messageType, p)
		if err != nil {
			t.Fatalf("ParseWebHook: %v", err)
		}
		if want := test.payload; !cmp.Equal(got, want) {
			t.Errorf("ParseWebHook(%#v, %#v) = %#v, want %#v", test.messageType, p, got, want)
		}
	}
}

func TestParseWebHook_BadMessageType(t *testing.T) {
	if _, err := ParseWebHook("bogus message type", []byte("{}")); err == nil {
		t.Fatal("ParseWebHook returned nil; wanted error")
	}
}

func TestValidatePayloadFromBody_UnableToParseBody(t *testing.T) {
	if _, err := ValidatePayloadFromBody("application/x-www-form-urlencoded", bytes.NewReader([]byte(`%`)), "sha1=", []byte{}); err == nil {
		t.Errorf("ValidatePayloadFromBody returned nil; wanted error")
	}
}

func TestValidatePayloadFromBody_UnsupportedContentType(t *testing.T) {
	if _, err := ValidatePayloadFromBody("invalid", bytes.NewReader([]byte(`{}`)), "sha1=", []byte{}); err == nil {
		t.Errorf("ValidatePayloadFromBody returned nil; wanted error")
	}
}

func TestDeliveryID(t *testing.T) {
	id := "8970a780-244e-11e7-91ca-da3aabcb9793"
	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		t.Fatalf("DeliveryID: %v", err)
	}
	req.Header.Set("X-Github-Delivery", id)

	got := DeliveryID(req)
	if got != id {
		t.Errorf("DeliveryID(%#v) = %q, want %q", req, got, id)
	}
}

func TestWebHookType(t *testing.T) {
	want := "yo"
	req := &http.Request{
		Header: http.Header{EventTypeHeader: []string{want}},
	}
	if got := WebHookType(req); got != want {
		t.Errorf("WebHookType = %q, want %q", got, want)
	}
}
