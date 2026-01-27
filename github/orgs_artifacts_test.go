// Copyright 2024 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_CreateArtifactDeploymentRecord(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ArtifactDeploymentRecord{Name: Ptr("test-n")}

	mux.HandleFunc("/orgs/o/artifacts/metadata/deployment-record", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONMarshal(t, input, `{"name":"test-n"}`)
		fmt.Fprint(w, `{"total_count":1,"deployment_records":[{"id":1}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.CreateArtifactDeploymentRecord(ctx, "o", input)
	if err != nil {
		t.Errorf("CreateArtifactDeploymentRecord returned error: %v", err)
	}

	want := &ArtifactDeploymentResponse{
		TotalCount:        Ptr(1),
		DeploymentRecords: []*ArtifactDeploymentRecord{{ID: Ptr(int64(1))}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("CreateArtifactDeploymentRecord returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_SetClusterDeploymentRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ArtifactDeploymentRecord{Name: Ptr("cluster-deploy")}

	mux.HandleFunc("/orgs/o/artifacts/metadata/deployment-record/cluster/c1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"total_count":1,"deployment_records":[{"id":2}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.SetClusterDeploymentRecords(ctx, "o", "c1", input)
	if err != nil {
		t.Errorf("SetClusterDeploymentRecords returned error: %v", err)
	}

	want := &ArtifactDeploymentResponse{
		TotalCount:        Ptr(1),
		DeploymentRecords: []*ArtifactDeploymentRecord{{ID: Ptr(int64(2))}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("SetClusterDeploymentRecords returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_ListArtifactDeploymentRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/artifacts/d/metadata/deployment-records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `{"total_count":1,"deployment_records":[{"id":1}]}`)
	})

	ctx := t.Context()
	opts := &ListOptions{Page: 2}
	got, _, err := client.Organizations.ListArtifactDeploymentRecords(ctx, "o", "d", opts)
	if err != nil {
		t.Errorf("ListArtifactDeploymentRecords returned error: %v", err)
	}

	want := &ArtifactDeploymentResponse{
		TotalCount:        Ptr(1),
		DeploymentRecords: []*ArtifactDeploymentRecord{{ID: Ptr(int64(1))}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListArtifactDeploymentRecords returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_CreateArtifactStorageRecord(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ArtifactStorageRecord{Name: Ptr("s-test")}

	mux.HandleFunc("/orgs/o/artifacts/metadata/storage-record", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"total_count":1,"storage_records":[{"name":"s-test"}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.CreateArtifactStorageRecord(ctx, "o", input)
	if err != nil {
		t.Errorf("CreateArtifactStorageRecord returned error: %v", err)
	}

	want := &ArtifactStorageResponse{
		TotalCount:     Ptr(1),
		StorageRecords: []*ArtifactStorageRecord{{Name: Ptr("s-test")}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("CreateArtifactStorageRecord returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_ListArtifactStorageRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/artifacts/d/metadata/storage-records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1,"storage_records":[{"name":"s-test"}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.ListArtifactStorageRecords(ctx, "o", "d", nil)
	if err != nil {
		t.Errorf("ListArtifactStorageRecords returned error: %v", err)
	}

	want := &ArtifactStorageResponse{
		TotalCount:     Ptr(1),
		StorageRecords: []*ArtifactStorageRecord{{Name: Ptr("s-test")}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListArtifactStorageRecords returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_ArtifactMetadata_InvalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)
	ctx := t.Context()

	_, _, err := client.Organizations.CreateArtifactDeploymentRecord(ctx, "%", nil)
	testURLParseError(t, err)

	_, _, err = client.Organizations.SetClusterDeploymentRecords(ctx, "%", "c", nil)
	testURLParseError(t, err)

	_, _, err = client.Organizations.CreateArtifactStorageRecord(ctx, "%", nil)
	testURLParseError(t, err)

	_, _, err = client.Organizations.ListArtifactDeploymentRecords(ctx, "%", "d", nil)
	testURLParseError(t, err)

	_, _, err = client.Organizations.ListArtifactDeploymentRecords(ctx, "%", "d", &ListOptions{})
	testURLParseError(t, err)

	_, _, err = client.Organizations.ListArtifactStorageRecords(ctx, "%", "d", nil)
	testURLParseError(t, err)

	_, _, err = client.Organizations.ListArtifactStorageRecords(ctx, "%", "d", &ListOptions{})
	testURLParseError(t, err)
}

func TestArtifactDeploymentRecord_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ArtifactDeploymentRecord{}, "{}")

	u := &ArtifactDeploymentRecord{
		ID:     Ptr(int64(1)),
		Name:   Ptr("n"),
		Status: Ptr("s"),
	}

	want := `{
		"id": 1,
		"name": "n",
		"status": "s"
	}`

	testJSONMarshal(t, u, want)
}

func TestArtifactMetadata_String(t *testing.T) {
	t.Parallel()

	r1 := ArtifactDeploymentRecord{Name: Ptr("n")}
	_ = r1.String()

	r2 := ArtifactDeploymentResponse{TotalCount: Ptr(1)}
	_ = r2.String()

	r3 := ArtifactStorageRecord{Name: Ptr("n")}
	_ = r3.String()

	r4 := ArtifactStorageResponse{TotalCount: Ptr(1)}
	_ = r4.String()
}
