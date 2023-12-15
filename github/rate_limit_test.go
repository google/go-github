// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRateLimits_String(t *testing.T) {
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
	}
	want := `github.RateLimits{Core:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, Search:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, GraphQL:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, IntegrationManifest:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, SourceImport:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, CodeScanningUpload:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, ActionsRunnerRegistration:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, SCIM:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, DependencySnapshots:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, CodeSearch:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}}`
	if got := v.String(); got != want {
		t.Errorf("RateLimits.String = %v, want %v", got, want)
	}
}

func TestRateLimits(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874},
			"graphql": {"limit":4,"remaining":3,"reset":1372700875},
			"integration_manifest": {"limit":5,"remaining":4,"reset":1372700876},
			"source_import": {"limit":6,"remaining":5,"reset":1372700877},
			"code_scanning_upload": {"limit":7,"remaining":6,"reset":1372700878},
			"actions_runner_registration": {"limit":8,"remaining":7,"reset":1372700879},
			"scim": {"limit":9,"remaining":8,"reset":1372700880},
			"dependency_snapshots": {"limit":10,"remaining":9,"reset":1372700881},
			"code_search": {"limit":11,"remaining":10,"reset":1372700882}
		}}`)
	})

	ctx := context.Background()
	rate, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
		GraphQL: &Rate{
			Limit:     4,
			Remaining: 3,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 55, 0, time.UTC).Local()},
		},
		IntegrationManifest: &Rate{
			Limit:     5,
			Remaining: 4,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 56, 0, time.UTC).Local()},
		},
		SourceImport: &Rate{
			Limit:     6,
			Remaining: 5,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 57, 0, time.UTC).Local()},
		},
		CodeScanningUpload: &Rate{
			Limit:     7,
			Remaining: 6,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 58, 0, time.UTC).Local()},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     8,
			Remaining: 7,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 59, 0, time.UTC).Local()},
		},
		SCIM: &Rate{
			Limit:     9,
			Remaining: 8,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 00, 0, time.UTC).Local()},
		},
		DependencySnapshots: &Rate{
			Limit:     10,
			Remaining: 9,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 1, 0, time.UTC).Local()},
		},
		CodeSearch: &Rate{
			Limit:     11,
			Remaining: 10,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 2, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}
	tests := []struct {
		category rateLimitCategory
		rate     *Rate
	}{
		{
			category: coreCategory,
			rate:     want.Core,
		},
		{
			category: searchCategory,
			rate:     want.Search,
		},
		{
			category: graphqlCategory,
			rate:     want.GraphQL,
		},
		{
			category: integrationManifestCategory,
			rate:     want.IntegrationManifest,
		},
		{
			category: sourceImportCategory,
			rate:     want.SourceImport,
		},
		{
			category: codeScanningUploadCategory,
			rate:     want.CodeScanningUpload,
		},
		{
			category: actionsRunnerRegistrationCategory,
			rate:     want.ActionsRunnerRegistration,
		},
		{
			category: scimCategory,
			rate:     want.SCIM,
		},
		{
			category: dependencySnapshotsCategory,
			rate:     want.DependencySnapshots,
		},
		{
			category: codeSearchCategory,
			rate:     want.CodeSearch,
		},
	}

	for _, tt := range tests {
		if got, want := client.rateLimits[tt.category], *tt.rate; got != want {
			t.Errorf("client.rateLimits[%v] is %+v, want %+v", tt.category, got, want)
		}
	}
}

func TestRateLimits_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "RateLimits"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.RateLimit.Get(ctx)
		return resp, err
	})
}

func TestRateLimits_overQuota(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	client.rateLimits[coreCategory] = Rate{
		Limit:     1,
		Remaining: 0,
		Reset:     Timestamp{time.Now().Add(time.Hour).Local()},
	}
	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874},
			"graphql": {"limit":4,"remaining":3,"reset":1372700875},
			"integration_manifest": {"limit":5,"remaining":4,"reset":1372700876},
			"source_import": {"limit":6,"remaining":5,"reset":1372700877},
			"code_scanning_upload": {"limit":7,"remaining":6,"reset":1372700878},
			"actions_runner_registration": {"limit":8,"remaining":7,"reset":1372700879},
			"scim": {"limit":9,"remaining":8,"reset":1372700880},
			"dependency_snapshots": {"limit":10,"remaining":9,"reset":1372700881},
			"code_search": {"limit":11,"remaining":10,"reset":1372700882}
		}}`)
	})

	ctx := context.Background()
	rate, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
		GraphQL: &Rate{
			Limit:     4,
			Remaining: 3,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 55, 0, time.UTC).Local()},
		},
		IntegrationManifest: &Rate{
			Limit:     5,
			Remaining: 4,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 56, 0, time.UTC).Local()},
		},
		SourceImport: &Rate{
			Limit:     6,
			Remaining: 5,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 57, 0, time.UTC).Local()},
		},
		CodeScanningUpload: &Rate{
			Limit:     7,
			Remaining: 6,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 58, 0, time.UTC).Local()},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     8,
			Remaining: 7,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 59, 0, time.UTC).Local()},
		},
		SCIM: &Rate{
			Limit:     9,
			Remaining: 8,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 00, 0, time.UTC).Local()},
		},
		DependencySnapshots: &Rate{
			Limit:     10,
			Remaining: 9,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 1, 0, time.UTC).Local()},
		},
		CodeSearch: &Rate{
			Limit:     11,
			Remaining: 10,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 48, 2, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}

	tests := []struct {
		category rateLimitCategory
		rate     *Rate
	}{
		{
			category: coreCategory,
			rate:     want.Core,
		},
		{
			category: searchCategory,
			rate:     want.Search,
		},
		{
			category: graphqlCategory,
			rate:     want.GraphQL,
		},
		{
			category: integrationManifestCategory,
			rate:     want.IntegrationManifest,
		},
		{
			category: sourceImportCategory,
			rate:     want.SourceImport,
		},
		{
			category: codeScanningUploadCategory,
			rate:     want.CodeScanningUpload,
		},
		{
			category: actionsRunnerRegistrationCategory,
			rate:     want.ActionsRunnerRegistration,
		},
		{
			category: scimCategory,
			rate:     want.SCIM,
		},
		{
			category: dependencySnapshotsCategory,
			rate:     want.DependencySnapshots,
		},
		{
			category: codeSearchCategory,
			rate:     want.CodeSearch,
		},
	}
	for _, tt := range tests {
		if got, want := client.rateLimits[tt.category], *tt.rate; got != want {
			t.Errorf("client.rateLimits[%v] is %+v, want %+v", tt.category, got, want)
		}
	}
}

func TestRateLimits_Marshal(t *testing.T) {
	testJSONMarshal(t, &RateLimits{}, "{}")

	u := &RateLimits{
		Core: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		Search: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		GraphQL: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		IntegrationManifest: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		SourceImport: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		CodeScanningUpload: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		ActionsRunnerRegistration: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		SCIM: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		DependencySnapshots: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		CodeSearch: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
	}

	want := `{
		"core": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"search": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"graphql": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"integration_manifest": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"source_import": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"code_scanning_upload": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"actions_runner_registration": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"scim": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"dependency_snapshots": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"code_search": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestRate_Marshal(t *testing.T) {
	testJSONMarshal(t, &Rate{}, "{}")

	u := &Rate{
		Limit:     1,
		Remaining: 1,
		Reset:     Timestamp{referenceTime},
	}

	want := `{
		"limit": 1,
		"remaining": 1,
		"reset": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
