// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRateLimits_String(t *testing.T) {
	t.Parallel()
	v := RateLimits{
		Core:                      &Rate{},
		Search:                    &Rate{},
		GraphQL:                   &Rate{},
		IntegrationManifest:       &Rate{},
		SourceImport:              &Rate{},
		CodeScanningUpload:        &Rate{},
		ActionsRunnerRegistration: &Rate{},
		SCIM:                      &Rate{},
		DependencySnapshots:       &Rate{},
		CodeSearch:                &Rate{},
		AuditLog:                  &Rate{},
		DependencySBOM:            &Rate{},
	}
	want := `github.RateLimits{Core:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, Search:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, GraphQL:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, IntegrationManifest:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, SourceImport:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, CodeScanningUpload:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, ActionsRunnerRegistration:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, SCIM:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, DependencySnapshots:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, CodeSearch:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, AuditLog:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}, DependencySBOM:github.Rate{Limit:0, Remaining:0, Used:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}, Resource:""}}`
	if got := v.String(); got != want {
		t.Errorf("RateLimits.String = %v, want %v", got, want)
	}
}

func TestRateLimits(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"used":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"used":1,"reset":1372700874},
			"graphql": {"limit":4,"remaining":3,"used":1,"reset":1372700875},
			"integration_manifest": {"limit":5,"remaining":4,"used":1,"reset":1372700876},
			"source_import": {"limit":6,"remaining":5,"used":1,"reset":1372700877},
			"code_scanning_upload": {"limit":7,"remaining":6,"used":1,"reset":1372700878},
			"actions_runner_registration": {"limit":8,"remaining":7,"used":1,"reset":1372700879},
			"scim": {"limit":9,"remaining":8,"used":1,"reset":1372700880},
			"dependency_snapshots": {"limit":10,"remaining":9,"used":1,"reset":1372700881},
			"code_search": {"limit":11,"remaining":10,"used":1,"reset":1372700882},
			"audit_log": {"limit": 12,"remaining":11,"used":1,"reset":1372700883},
			"dependency_sbom": {"limit": 100,"remaining":100,"used":0,"reset":1372700884}
		}}`)
	})

	ctx := t.Context()
	rate, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
		GraphQL: &Rate{
			Limit:     4,
			Remaining: 3,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 55, 0, time.UTC).Local()},
		},
		IntegrationManifest: &Rate{
			Limit:     5,
			Remaining: 4,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 56, 0, time.UTC).Local()},
		},
		SourceImport: &Rate{
			Limit:     6,
			Remaining: 5,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 57, 0, time.UTC).Local()},
		},
		CodeScanningUpload: &Rate{
			Limit:     7,
			Remaining: 6,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 58, 0, time.UTC).Local()},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     8,
			Remaining: 7,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 59, 0, time.UTC).Local()},
		},
		SCIM: &Rate{
			Limit:     9,
			Remaining: 8,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 0, 0, time.UTC).Local()},
		},
		DependencySnapshots: &Rate{
			Limit:     10,
			Remaining: 9,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 1, 0, time.UTC).Local()},
		},
		CodeSearch: &Rate{
			Limit:     11,
			Remaining: 10,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 2, 0, time.UTC).Local()},
		},
		AuditLog: &Rate{
			Limit:     12,
			Remaining: 11,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 3, 0, time.UTC).Local()},
		},
		DependencySBOM: &Rate{
			Limit:     100,
			Remaining: 100,
			Used:      0,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 4, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}
	tests := []struct {
		category RateLimitCategory
		rate     *Rate
	}{
		{
			category: CoreCategory,
			rate:     want.Core,
		},
		{
			category: SearchCategory,
			rate:     want.Search,
		},
		{
			category: GraphqlCategory,
			rate:     want.GraphQL,
		},
		{
			category: IntegrationManifestCategory,
			rate:     want.IntegrationManifest,
		},
		{
			category: SourceImportCategory,
			rate:     want.SourceImport,
		},
		{
			category: CodeScanningUploadCategory,
			rate:     want.CodeScanningUpload,
		},
		{
			category: ActionsRunnerRegistrationCategory,
			rate:     want.ActionsRunnerRegistration,
		},
		{
			category: ScimCategory,
			rate:     want.SCIM,
		},
		{
			category: DependencySnapshotsCategory,
			rate:     want.DependencySnapshots,
		},
		{
			category: CodeSearchCategory,
			rate:     want.CodeSearch,
		},
		{
			category: AuditLogCategory,
			rate:     want.AuditLog,
		},
		{
			category: DependencySBOMCategory,
			rate:     want.DependencySBOM,
		},
	}

	for _, tt := range tests {
		if got, want := client.rateLimits[tt.category], *tt.rate; got != want {
			t.Errorf("client.rateLimits[%v] is %+v, want %+v", tt.category, got, want)
		}
	}
}

func TestRateLimits_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "RateLimits"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.RateLimit.Get(ctx)
		return resp, err
	})
}

func TestRateLimits_overQuota(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	client.rateLimits[CoreCategory] = Rate{
		Limit:     1,
		Remaining: 0,
		Used:      1,
		Reset:     Timestamp{time.Now().Add(time.Hour).Local()},
	}
	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"used":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"used":1,"reset":1372700874},
			"graphql": {"limit":4,"remaining":3,"used":1,"reset":1372700875},
			"integration_manifest": {"limit":5,"remaining":4,"used":1,"reset":1372700876},
			"source_import": {"limit":6,"remaining":5,"used":1,"reset":1372700877},
			"code_scanning_upload": {"limit":7,"remaining":6,"used":1,"reset":1372700878},
			"actions_runner_registration": {"limit":8,"remaining":7,"used":1,"reset":1372700879},
			"scim": {"limit":9,"remaining":8,"used":1,"reset":1372700880},
			"dependency_snapshots": {"limit":10,"remaining":9,"used":1,"reset":1372700881},
			"code_search": {"limit":11,"remaining":10,"used":1,"reset":1372700882},
			"audit_log": {"limit":12,"remaining":11,"used":1,"reset":1372700883},
			"dependency_sbom": {"limit":13,"remaining":12,"used":1,"reset":1372700884}
		}}`)
	})

	ctx := t.Context()
	rate, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
		GraphQL: &Rate{
			Limit:     4,
			Remaining: 3,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 55, 0, time.UTC).Local()},
		},
		IntegrationManifest: &Rate{
			Limit:     5,
			Remaining: 4,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 56, 0, time.UTC).Local()},
		},
		SourceImport: &Rate{
			Limit:     6,
			Remaining: 5,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 57, 0, time.UTC).Local()},
		},
		CodeScanningUpload: &Rate{
			Limit:     7,
			Remaining: 6,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 58, 0, time.UTC).Local()},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     8,
			Remaining: 7,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 59, 0, time.UTC).Local()},
		},
		SCIM: &Rate{
			Limit:     9,
			Remaining: 8,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 0, 0, time.UTC).Local()},
		},
		DependencySnapshots: &Rate{
			Limit:     10,
			Remaining: 9,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 1, 0, time.UTC).Local()},
		},
		CodeSearch: &Rate{
			Limit:     11,
			Remaining: 10,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 2, 0, time.UTC).Local()},
		},
		AuditLog: &Rate{
			Limit:     12,
			Remaining: 11,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 3, 0, time.UTC).Local()},
		},
		DependencySBOM: &Rate{
			Limit:     13,
			Remaining: 12,
			Used:      1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 4, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}

	tests := []struct {
		category RateLimitCategory
		rate     *Rate
	}{
		{
			category: CoreCategory,
			rate:     want.Core,
		},
		{
			category: SearchCategory,
			rate:     want.Search,
		},
		{
			category: GraphqlCategory,
			rate:     want.GraphQL,
		},
		{
			category: IntegrationManifestCategory,
			rate:     want.IntegrationManifest,
		},
		{
			category: SourceImportCategory,
			rate:     want.SourceImport,
		},
		{
			category: CodeScanningUploadCategory,
			rate:     want.CodeScanningUpload,
		},
		{
			category: ActionsRunnerRegistrationCategory,
			rate:     want.ActionsRunnerRegistration,
		},
		{
			category: ScimCategory,
			rate:     want.SCIM,
		},
		{
			category: DependencySnapshotsCategory,
			rate:     want.DependencySnapshots,
		},
		{
			category: CodeSearchCategory,
			rate:     want.CodeSearch,
		},
		{
			category: AuditLogCategory,
			rate:     want.AuditLog,
		},
		{
			category: DependencySBOMCategory,
			rate:     want.DependencySBOM,
		},
	}
	for _, tt := range tests {
		if got, want := client.rateLimits[tt.category], *tt.rate; got != want {
			t.Errorf("client.rateLimits[%v] is %+v, want %+v", tt.category, got, want)
		}
	}
}

func TestRateLimits_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RateLimits{}, `{
		"core": null,
		"search": null,
		"graphql": null,
		"integration_manifest": null,
		"source_import": null,
		"code_scanning_upload": null,
		"actions_runner_registration": null,
		"scim": null,
		"dependency_snapshots": null,
		"code_search": null,
		"audit_log": null,
		"dependency_sbom": null
	}`)

	u := &RateLimits{
		Core: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		Search: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		GraphQL: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		IntegrationManifest: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		SourceImport: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		CodeScanningUpload: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		SCIM: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		DependencySnapshots: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		CodeSearch: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		AuditLog: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
		DependencySBOM: &Rate{
			Limit:     1,
			Remaining: 1,
			Used:      0,
			Reset:     Timestamp{referenceTime},
		},
	}

	want := `{
		"core": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"search": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"graphql": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"integration_manifest": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"source_import": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"code_scanning_upload": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"actions_runner_registration": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"scim": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"dependency_snapshots": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"code_search": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"audit_log": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		},
		"dependency_sbom": {
			"limit": 1,
			"remaining": 1,
			"used": 0,
			"reset": ` + referenceTimeStr + `
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestRate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Rate{}, `{
		"limit": 0,
		"remaining": 0,
		"used": 0,
		"reset": "0001-01-01T00:00:00Z"
	}`)

	u := &Rate{
		Limit:     1,
		Remaining: 1,
		Used:      0,
		Reset:     Timestamp{referenceTime},
		Resource:  "core",
	}

	want := `{
		"limit": 1,
		"remaining": 1,
		"used": 0,
		"reset": ` + referenceTimeStr + `,
		"resource": "core"
	}`

	testJSONMarshal(t, u, want)
}
