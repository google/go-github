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
	"strings"
	"testing"
	"time"
)

func TestOrganizationService_GetAuditLog(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/audit-log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `[
	{
		"@timestamp": 1615077308538,
		"_document_id": "beeZYapIUe-wKg5-beadb33",
		"action": "workflows.completed_workflow_run",
		"active": true,
		"actor": "testactor",
		"actor_ip": "10.0.0.1",
		"actor_location": {
			"country_code": "US"
		},
		"cancelled_at": "2021-03-07T00:35:08.000Z",
		"completed_at": "2021-03-07T00:35:08.000Z",
		"conclusion": "success",
		"config": {
			"content_type": "json",
			"insecure_ssl": "0",
			"url": "https://example.com/deadbeef-new-hook"
		},
		"created_at": 1615077308538,
		"event": "schedule",
		"events": ["code_scanning_alert"],
		"hashed_token": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		"head_branch": "master",
		"head_sha": "5acdeadbeef64d1a62388e901e5cdc9358644b37",
		"job_workflow_ref": "testorg/testrepo/.github/workflows/testjob.yml@refs/pull/1/merge",
		"name": "Code scanning - action",
		"oauth_application_id": 1,
		"old_permission": "read",
		"org": "o",
		"org_id": 1,
		"overridden_codes": [
			"review_policy_not_satisfied"
		],
		"permission": "admin",
		"pull_request_id": 1,
		"pull_request_title": "a pr title",
		"pull_request_url": "https://github.com/testorg/testrepo/pull/1",
		"reasons": [
			{
				"code": "a code",
				"message": "a message"
			}
		],
		"programmatic_access_type": "GitHub App server-to-server token",
		"referrer": "a referrer",
		"repo": "o/blue-crayon-1",
		"run_attempt": 1,
		"run_number": 1,
		"started_at": "2021-03-07T00:33:04.000Z",
		"token_id": 1,
		"token_scopes": "gist,repo:read",
		"topic": "cp1-iad.ingest.github.actions.v0.WorkflowUpdate",
		"trigger_id": null,
		"user_agent": "a user agent",
		"workflow_id": 123456,
		"workflow_run_id": 628312345
	}]`)
	})
	ctx := context.Background()
	getOpts := GetAuditLogOptions{
		Include: Ptr("all"),
		Phrase:  Ptr("action:workflows"),
		Order:   Ptr("asc"),
	}

	auditEntries, resp, err := client.Organizations.GetAuditLog(ctx, "o", &getOpts)
	if err != nil {
		t.Errorf("Organizations.GetAuditLog returned error: %v", err)
	}
	timestamp := time.Unix(0, 1615077308538*1e6)

	want := []*AuditEntry{
		{
			Timestamp:  &Timestamp{timestamp},
			DocumentID: Ptr("beeZYapIUe-wKg5-beadb33"),
			Action:     Ptr("workflows.completed_workflow_run"),
			Actor:      Ptr("testactor"),
			ActorLocation: &ActorLocation{
				CountryCode: Ptr("US"),
			},
			CreatedAt:   &Timestamp{timestamp},
			HashedToken: Ptr("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			Org:         Ptr("o"),
			OrgID:       Ptr(int64(1)),
			TokenID:     Ptr(int64(1)),
			TokenScopes: Ptr("gist,repo:read"),
			AdditionalFields: map[string]any{
				"actor_ip":                 "10.0.0.1",
				"active":                   true,
				"cancelled_at":             "2021-03-07T00:35:08.000Z",
				"completed_at":             "2021-03-07T00:35:08.000Z",
				"conclusion":               "success",
				"event":                    "schedule",
				"head_branch":              "master",
				"head_sha":                 "5acdeadbeef64d1a62388e901e5cdc9358644b37",
				"job_workflow_ref":         "testorg/testrepo/.github/workflows/testjob.yml@refs/pull/1/merge",
				"name":                     "Code scanning - action",
				"oauth_application_id":     float64(1),
				"old_permission":           "read",
				"overridden_codes":         []any{"review_policy_not_satisfied"},
				"permission":               "admin",
				"programmatic_access_type": "GitHub App server-to-server token",
				"pull_request_id":          float64(1),
				"pull_request_title":       "a pr title",
				"pull_request_url":         "https://github.com/testorg/testrepo/pull/1",
				"reasons": []any{map[string]any{
					"code":    "a code",
					"message": "a message",
				}},
				"referrer":        "a referrer",
				"repo":            "o/blue-crayon-1",
				"run_attempt":     float64(1),
				"run_number":      float64(1),
				"started_at":      "2021-03-07T00:33:04.000Z",
				"topic":           "cp1-iad.ingest.github.actions.v0.WorkflowUpdate",
				"user_agent":      "a user agent",
				"workflow_id":     float64(123456),
				"workflow_run_id": float64(628312345),
				"events":          []any{"code_scanning_alert"},
				"config": map[string]any{
					"content_type": "json",
					"insecure_ssl": "0",
					"url":          "https://example.com/deadbeef-new-hook",
				},
			},
		},
	}

	assertNoDiff(t, want, auditEntries)

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

	testNewRequestAndDoFailureCategory(t, methodName, client, AuditLogCategory, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAuditLog(ctx, "o", &GetAuditLogOptions{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGetAuditLogOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &GetAuditLogOptions{}, `{
		"After": "",
		"Before": "",
		"Cursor": "",
		"First": 0,
		"Include": null,
		"Last": 0,
		"Order": null,
		"Page": "",
		"PerPage": 0,
		"Phrase": null
	}`)

	u := &GetAuditLogOptions{
		Phrase:  Ptr("p"),
		Include: Ptr("i"),
		Order:   Ptr("o"),
		ListCursorOptions: ListCursorOptions{
			Page:    "p",
			PerPage: 1,
			After:   "a",
			Before:  "b",
		},
	}

	want := `{
		"Phrase": "p",
		"Include": "i",
		"Order": "o",
		"Page": "p",
		"PerPage": 1,
		"After": "a",
		"Before": "b",
		"Cursor": "",
    	"First": 0,
		"Last": 0
	}`

	testJSONMarshal(t, u, want)
}

func TestHookConfig_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &HookConfig{}, "{}")

	u := &HookConfig{
		ContentType: Ptr("ct"),
		InsecureSSL: Ptr("ct"),
		URL:         Ptr("url"),
	}

	want := `{
		"content_type": "ct",
		"insecure_ssl": "ct",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestAuditEntry_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AuditEntry{}, "{}")

	u := &AuditEntry{
		Action:                   Ptr("a"),
		Actor:                    Ptr("ac"),
		ActorLocation:            &ActorLocation{CountryCode: Ptr("alcc")},
		Business:                 Ptr("b"),
		CreatedAt:                &Timestamp{referenceTime},
		DocumentID:               Ptr("did"),
		ExternalIdentityNameID:   Ptr("ein"),
		ExternalIdentityUsername: Ptr("eiu"),
		HashedToken:              Ptr("ht"),
		Org:                      Ptr("o"),
		OrgID:                    Ptr(int64(1)),
		Timestamp:                &Timestamp{referenceTime},
		TokenID:                  Ptr(int64(1)),
		TokenScopes:              Ptr("ts"),
		User:                     Ptr("u"),
		Data: map[string]any{
			"old_name":  "on",
			"old_login": "ol",
		},
		AdditionalFields: map[string]any{
			"active":       false,
			"active_was":   false,
			"actor_ip":     "aip",
			"blocked_user": "bu",
			"cancelled_at": "2021-03-07T00:35:08.000Z",
			"completed_at": "2021-03-07T00:35:08.000Z",
			"conclusion":   "c",
			"config": map[string]any{
				"url": "s",
			},
			"config_was": map[string]any{
				"url": "s",
			},
			"content_type":             "ct",
			"deploy_key_fingerprint":   "dkf",
			"emoji":                    "e",
			"environment_name":         "en",
			"event":                    "e",
			"events":                   []any{"s"},
			"events_were":              []any{"s"},
			"explanation":              "e",
			"fingerprint":              "f",
			"head_branch":              "hb",
			"head_sha":                 "hsha",
			"hook_id":                  float64(1),
			"is_hosted_runner":         false,
			"job_name":                 "jn",
			"limited_availability":     false,
			"message":                  "m",
			"name":                     "n",
			"old_permission":           "op",
			"old_user":                 "ou",
			"openssh_public_key":       "osshpk",
			"permission":               "p",
			"previous_visibility":      "pv",
			"programmatic_access_type": "pat",
			"pull_request_id":          float64(1),
			"pull_request_title":       "prt",
			"pull_request_url":         "pru",
			"read_only":                "ro",
			"reasons": []any{
				map[string]any{
					"code":    "c",
					"message": "m",
				},
			},
			"referrer":                "a referrer",
			"repo":                    "r",
			"repository":              "repo",
			"repository_public":       false,
			"run_attempt":             1,
			"runner_group_id":         1,
			"runner_group_name":       "rgn",
			"runner_id":               1,
			"runner_labels":           []any{"s"},
			"runner_name":             "rn",
			"secrets_passed":          []any{"s"},
			"source_version":          "sv",
			"started_at":              "2006-01-02T15:04:05Z",
			"target_login":            "tl",
			"target_version":          "tv",
			"team":                    "t",
			"topic":                   "cp1-iad.ingest.github.actions.v0.WorkflowUpdate",
			"transport_protocol":      1,
			"transport_protocol_name": "tpn",
			"trigger_id":              1,
			"user_agent":              "ua",
			"visibility":              "v",
			"workflow_id":             1,
			"workflow_run_id":         1,
		},
	}

	want := `{
		"action": "a",
		"active": false,
		"active_was": false,
		"actor": "ac",
		"actor_ip": "aip",
		"actor_location": {
			"country_code": "alcc"
		},
		"blocked_user": "bu",
		"business": "b",
		"cancelled_at": "2021-03-07T00:35:08.000Z",
		"completed_at": "2021-03-07T00:35:08.000Z",
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
		"external_identity_nameid": "ein",
		"external_identity_username": "eiu",
		"fingerprint": "f",
		"hashed_token": "ht",
		"head_branch": "hb",
		"head_sha": "hsha",
		"hook_id": 1,
		"is_hosted_runner": false,
		"job_name": "jn",
		"limited_availability": false,
		"message": "m",
		"name": "n",
		"old_permission": "op",
		"old_user": "ou",
		"openssh_public_key": "osshpk",
		"org": "o",
		"org_id": 1,
		"permission": "p",
		"previous_visibility": "pv",
		"programmatic_access_type": "pat",
		"pull_request_id": 1,
		"pull_request_title": "prt",
		"pull_request_url": "pru",
		"reasons": [{
			"code": "c",
			"message": "m"
		}],
		"referrer": "a referrer",
		"read_only": "ro",
		"repo": "r",
		"repository": "repo",
		"repository_public": false,
		"run_attempt": 1,
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
		"token_id": 1,
		"token_scopes": "ts",
		"topic": "cp1-iad.ingest.github.actions.v0.WorkflowUpdate",
		"transport_protocol_name": "tpn",
		"transport_protocol": 1,
		"trigger_id": 1,
		"user": "u",
		"user_agent": "ua",
		"visibility": "v",
		"workflow_id": 1,
		"workflow_run_id": 1,
		"data": {
			"old_name": "on",
			"old_login": "ol"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestAuditEntry_Unmarshal(t *testing.T) {
	t.Parallel()

	// Test case 1: JSON with both defined fields and additional fields
	jsonData := `{
		"action": "org.update_member",
		"actor": "testuser",
		"actor_location": {
			"country_code": "US"
		},
		"created_at": "2021-03-07T00:35:08.000Z",
		"org": "testorg",
		"org_id": 12345,
		"user": "memberuser",
		"user_id": 67890,
		"custom_field": "custom_value",
		"another_field": 42,
		"nested_field": {
			"key": "value"
		},
		"array_field": ["item1", "item2"]
	}`

	var entry AuditEntry
	err := json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Errorf("Error unmarshaling JSON: %v", err)
	}

	// Check defined fields
	if *entry.Action != "org.update_member" {
		t.Errorf("Action = %v, want %v", *entry.Action, "org.update_member")
	}
	if *entry.Actor != "testuser" {
		t.Errorf("Actor = %v, want %v", *entry.Actor, "testuser")
	}
	if *entry.ActorLocation.CountryCode != "US" {
		t.Errorf("ActorLocation.CountryCode = %v, want %v", *entry.ActorLocation.CountryCode, "US")
	}
	if *entry.Org != "testorg" {
		t.Errorf("Org = %v, want %v", *entry.Org, "testorg")
	}
	if *entry.OrgID != 12345 {
		t.Errorf("OrgID = %v, want %v", *entry.OrgID, 12345)
	}
	if *entry.User != "memberuser" {
		t.Errorf("User = %v, want %v", *entry.User, "memberuser")
	}
	if *entry.UserID != 67890 {
		t.Errorf("UserID = %v, want %v", *entry.UserID, 67890)
	}

	// Check additional fields
	if entry.AdditionalFields["custom_field"] != "custom_value" {
		t.Errorf("AdditionalFields[\"custom_field\"] = %v, want %v", entry.AdditionalFields["custom_field"], "custom_value")
	}
	if entry.AdditionalFields["another_field"] != float64(42) {
		t.Errorf("AdditionalFields[\"another_field\"] = %v, want %v", entry.AdditionalFields["another_field"], float64(42))
	}

	// Check nested fields
	nestedField, ok := entry.AdditionalFields["nested_field"].(map[string]interface{})
	if !ok {
		t.Errorf("AdditionalFields[\"nested_field\"] is not a map")
	} else if nestedField["key"] != "value" {
		t.Errorf("AdditionalFields[\"nested_field\"][\"key\"] = %v, want %v", nestedField["key"], "value")
	}

	// Check array fields
	arrayField, ok := entry.AdditionalFields["array_field"].([]interface{})
	if !ok {
		t.Errorf("AdditionalFields[\"array_field\"] is not an array")
	} else {
		if len(arrayField) != 2 {
			t.Errorf("len(AdditionalFields[\"array_field\"]) = %v, want %v", len(arrayField), 2)
		}
		if arrayField[0] != "item1" {
			t.Errorf("AdditionalFields[\"array_field\"][0] = %v, want %v", arrayField[0], "item1")
		}
		if arrayField[1] != "item2" {
			t.Errorf("AdditionalFields[\"array_field\"][1] = %v, want %v", arrayField[1], "item2")
		}
	}

	// Test case 2: Empty JSON
	err = json.Unmarshal([]byte("{}"), &entry)
	if err != nil {
		t.Errorf("Error unmarshaling empty JSON: %v", err)
	}
	if entry.AdditionalFields != nil {
		t.Errorf("AdditionalFields = %v, want nil", entry.AdditionalFields)
	}

	// Test case 3: JSON with only additional fields
	jsonData = `{
		"custom_field": "custom_value",
		"another_field": 42
	}`

	err = json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Errorf("Error unmarshaling JSON with only additional fields: %v", err)
	}
	if entry.AdditionalFields["custom_field"] != "custom_value" {
		t.Errorf("AdditionalFields[\"custom_field\"] = %v, want %v", entry.AdditionalFields["custom_field"], "custom_value")
	}
	if entry.AdditionalFields["another_field"] != float64(42) {
		t.Errorf("AdditionalFields[\"another_field\"] = %v, want %v", entry.AdditionalFields["another_field"], float64(42))
	}

	// Test case 4: Test that nil values in AdditionalFields are removed
	jsonData = `{
		"action": "org.update_member",
		"null_field": null
	}`

	err = json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Errorf("Error unmarshaling JSON with null field: %v", err)
	}
	if _, exists := entry.AdditionalFields["null_field"]; exists {
		t.Errorf("AdditionalFields contains null_field, but it should have been removed")
	}
}
