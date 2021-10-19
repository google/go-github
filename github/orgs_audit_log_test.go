// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// `Use` of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationService_GetAuditLog(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/audit-log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `[
		{
        "active": true,
        "workflow_id": 123456,
        "head_branch": "master",
        "org": "o",
        "trigger_id": null,
        "repo": "o/blue-crayon-1",
        "created_at": 1615077308538,
        "head_sha": "5acdeadbeef64d1a62388e901e5cdc9358644b37",
        "conclusion": "success",
        "actor": "testactor",
        "completed_at": "2021-03-07T00:35:08.000Z",
        "@timestamp": 1615077308538,
        "name": "Code scanning - action",
        "action": "workflows.completed_workflow_run",
        "started_at": "2021-03-07T00:33:04.000Z",
        "event": "schedule",
        "workflow_run_id": 628312345,
        "_document_id": "beeZYapIUe-wKg5-beadb33",
        "config": {
            "content_type": "json",
            "insecure_ssl": "0",
            "url": "https://example.com/deadbeef-new-hook"
         },
        "events": ["code_scanning_alert"]
		}]`)
	})
	ctx := context.Background()
	getOpts := GetAuditLogOptions{
		Include: String("all"),
		Phrase:  String("action:workflows"),
		Order:   String("asc"),
	}

	auditEntries, resp, err := client.Organizations.GetAuditLog(ctx, "o", &getOpts)
	if err != nil {
		t.Errorf("Organizations.GetAuditLog returned error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2021-03-07T00:33:04.000Z")
	completedAt, _ := time.Parse(time.RFC3339, "2021-03-07T00:35:08.000Z")
	timestamp := time.Unix(0, 1615077308538*1e6)

	want := []*AuditEntry{
		{
			Timestamp:     &Timestamp{timestamp},
			DocumentID:    String("beeZYapIUe-wKg5-beadb33"),
			Action:        String("workflows.completed_workflow_run"),
			Actor:         String("testactor"),
			Active:        Bool(true),
			CompletedAt:   &Timestamp{completedAt},
			Conclusion:    String("success"),
			CreatedAt:     &Timestamp{timestamp},
			Event:         String("schedule"),
			HeadBranch:    String("master"),
			HeadSHA:       String("5acdeadbeef64d1a62388e901e5cdc9358644b37"),
			Name:          String("Code scanning - action"),
			Org:           String("o"),
			Repo:          String("o/blue-crayon-1"),
			StartedAt:     &Timestamp{startedAt},
			WorkflowID:    Int64(123456),
			WorkflowRunID: Int64(628312345),
			Events:        []string{"code_scanning_alert"},
			Config: &HookConfig{
				ContentType: String("json"),
				InsecureSSL: String("0"),
				URL:         String("https://example.com/deadbeef-new-hook"),
			},
		},
	}

	if !cmp.Equal(auditEntries, want) {
		t.Errorf("Organizations.GetAuditLog return \ngot: %+v,\nwant:%+v", auditEntries, want)
	}

	// assert query string has lower case params
	requestedQuery := resp.Request.URL.RawQuery
	if !strings.Contains(requestedQuery, "phrase") {
		t.Errorf("Organizations.GetAuditLog query string \ngot: %+v,\nwant:%+v", requestedQuery, "phrase")
	}

	const methodName = "GetAuditLog"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetAuditLog(ctx, "\n", &getOpts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAuditLog(ctx, "o", &GetAuditLogOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGetAuditLogOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &GetAuditLogOptions{}, "{}")

	u := &GetAuditLogOptions{
		Phrase:  String("p"),
		Include: String("i"),
		Order:   String("o"),
		ListCursorOptions: ListCursorOptions{
			Page:    "p",
			PerPage: 1,
			After:   "a",
			Before:  "b",
		},
	}

	want := `{
		"phrase": "p",
		"include": "i",
		"order": "o",
		"Page": "p",
		"PerPage": 1,
		"After": "a",
		"Before": "b"
	}`

	testJSONMarshal(t, u, want)
}

func TestHookConfig_Marshal(t *testing.T) {
	testJSONMarshal(t, &HookConfig{}, "{}")

	u := &HookConfig{
		ContentType: String("ct"),
		InsecureSSL: String("ct"),
		URL:         String("url"),
	}

	want := `{
		"content_type": "ct",
		"insecure_ssl": "ct",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuditEntry_Marshal(t *testing.T) {
	testJSONMarshal(t, &AuditEntry{}, "{}")

	u := &AuditEntry{
		Action:                String("a"),
		Active:                Bool(false),
		ActiveWas:             Bool(false),
		Actor:                 String("ac"),
		BlockedUser:           String("bu"),
		Business:              String("b"),
		CancelledAt:           &Timestamp{referenceTime},
		CompletedAt:           &Timestamp{referenceTime},
		Conclusion:            String("c"),
		Config:                &HookConfig{URL: String("s")},
		ConfigWas:             &HookConfig{URL: String("s")},
		ContentType:           String("ct"),
		CreatedAt:             &Timestamp{referenceTime},
		DeployKeyFingerprint:  String("dkf"),
		DocumentID:            String("did"),
		Emoji:                 String("e"),
		EnvironmentName:       String("en"),
		Event:                 String("e"),
		Events:                []string{"s"},
		EventsWere:            []string{"s"},
		Explanation:           String("e"),
		Fingerprint:           String("f"),
		HeadBranch:            String("hb"),
		HeadSHA:               String("hsha"),
		HookID:                Int64(1),
		IsHostedRunner:        Bool(false),
		JobName:               String("jn"),
		LimitedAvailability:   Bool(false),
		Message:               String("m"),
		Name:                  String("n"),
		OldUser:               String("ou"),
		OpenSSHPublicKey:      String("osshpk"),
		Org:                   String("o"),
		PreviousVisibility:    String("pv"),
		ReadOnly:              String("ro"),
		Repo:                  String("r"),
		Repository:            String("repo"),
		RepositoryPublic:      Bool(false),
		RunnerGroupID:         Int64(1),
		RunnerGroupName:       String("rgn"),
		RunnerID:              Int64(1),
		RunnerLabels:          []string{"s"},
		RunnerName:            String("rn"),
		SecretsPassed:         []string{"s"},
		SourceVersion:         String("sv"),
		StartedAt:             &Timestamp{referenceTime},
		TargetLogin:           String("tl"),
		TargetVersion:         String("tv"),
		Team:                  String("t"),
		Timestamp:             &Timestamp{referenceTime},
		TransportProtocolName: String("tpn"),
		TransportProtocol:     Int(1),
		TriggerID:             Int64(1),
		User:                  String("u"),
		Visibility:            String("v"),
		WorkflowID:            Int64(1),
		WorkflowRunID:         Int64(1),
	}

	want := `{
		"action": "a",
		"active": false,
		"active_was": false,
		"actor": "ac",
		"blocked_user": "bu",
		"business": "b",
		"cancelled_at": ` + referenceTimeStr + `,
		"completed_at": ` + referenceTimeStr + `,
		"conclusion": "c",
		"config": {
			"url": "s"
		},
		"config_was": {
			"url": "s"
		},
		"content_type": "ct",
		"created_at": ` + referenceTimeStr + `,
		"deploy_key_fingerprint": "dkf",
		"_document_id": "did",
		"emoji": "e",
		"environment_name": "en",
		"event": "e",
		"events": [
			"s"
		],
		"events_were": [
			"s"
		],
		"explanation": "e",
		"fingerprint": "f",
		"head_branch": "hb",
		"head_sha": "hsha",
		"hook_id": 1,
		"is_hosted_runner": false,
		"job_name": "jn",
		"limited_availability": false,
		"message": "m",
		"name": "n",
		"old_user": "ou",
		"openssh_public_key": "osshpk",
		"org": "o",
		"previous_visibility": "pv",
		"read_only": "ro",
		"repo": "r",
		"repository": "repo",
		"repository_public": false,
		"runner_group_id": 1,
		"runner_group_name": "rgn",
		"runner_id": 1,
		"runner_labels": [
			"s"
		],
		"runner_name": "rn",
		"secrets_passed": [
			"s"
		],
		"source_version": "sv",
		"started_at": ` + referenceTimeStr + `,
		"target_login": "tl",
		"target_version": "tv",
		"team": "t",
		"@timestamp": ` + referenceTimeStr + `,
		"transport_protocol_name": "tpn",
		"transport_protocol": 1,
		"trigger_id": 1,
		"user": "u",
		"visibility": "v",
		"workflow_id": 1,
		"workflow_run_id": 1
	}`

	testJSONMarshal(t, u, want)
}
