// Copyright 2026 The go-github AUTHORS. All rights reserved.
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

	input := &CreateArtifactDeploymentRequest{
		Name:               Ptr("test-n"),
		Digest:             Ptr("sha256:123"),
		Version:            Ptr("v1.0.0"),
		Status:             Ptr("deployed"),
		LogicalEnvironment: Ptr("prod"),
		DeploymentName:     Ptr("dep-1"),
		RuntimeRisks:       []string{"critical-resource", "internet-exposed"},
		GithubRepository:   Ptr("octo-org/octo-repo"),
		Tags: map[string]string{
			"data-access": "sensitive",
		},
	}

	_ = input.String()

	mux.HandleFunc("/orgs/o/artifacts/metadata/deployment-record", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONMarshal(t, input, `{"digest":"sha256:123","name":"test-n","version":"v1.0.0","status":"deployed","logical_environment":"prod","deployment_name":"dep-1","tags":{"data-access":"sensitive"},"runtime_risks":["critical-resource","internet-exposed"],"github_repository":"octo-org/octo-repo"}`)
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

	_ = want.String()
	_ = want.DeploymentRecords[0].String()

	if !cmp.Equal(got, want) {
		t.Errorf("CreateArtifactDeploymentRecord returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_SetClusterDeploymentRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ClusterDeploymentRecordsRequest{
		LogicalEnvironment:  Ptr("prod"),
		PhysicalEnvironment: Ptr("pacific-east"),
		Deployments: []*CreateArtifactDeploymentRequest{
			{
				Name:    Ptr("awesome-image"),
				Version: Ptr("v2.0"),
				Status:  Ptr("deployed"),
			},
		},
	}

	_ = input.String()

	mux.HandleFunc("/orgs/o/artifacts/metadata/deployment-record/cluster/c1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONMarshal(t, input, `{"logical_environment":"prod","physical_environment":"pacific-east","deployments":[{"name":"awesome-image","version":"v2.0","status":"deployed"}]}`)
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

func TestOrganizationsService_CreateArtifactStorageRecord(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &CreateArtifactStorageRequest{
		Name:             Ptr("libfoo"),
		Version:          Ptr("v1.2.3"),
		Path:             Ptr("target/libs"),
		GithubRepository: Ptr("org/repo"),
		RegistryURL:      Ptr("https://reg.example.com"),
		Status:           Ptr("active"),
	}

	_ = input.String()

	mux.HandleFunc("/orgs/o/artifacts/metadata/storage-record", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONMarshal(t, input, `{"name":"libfoo","version":"v1.2.3","path":"target/libs","registry_url":"https://reg.example.com","status":"active","github_repository":"org/repo"}`)
		fmt.Fprint(w, `{"total_count":1,"storage_records":[{"name":"libfoo"}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.CreateArtifactStorageRecord(ctx, "o", input)
	if err != nil {
		t.Errorf("CreateArtifactStorageRecord returned error: %v", err)
	}

	want := &ArtifactStorageResponse{
		TotalCount:     Ptr(1),
		StorageRecords: []*ArtifactStorageRecord{{Name: Ptr("libfoo")}},
	}

	_ = want.String()
	_ = want.StorageRecords[0].String()

	if !cmp.Equal(got, want) {
		t.Errorf("CreateArtifactStorageRecord returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_ListArtifactDeploymentRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/artifacts/sha256:abc/metadata/deployment-records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1,"deployment_records":[{"id":1, "runtime_risks": ["sensitive-data"]}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.ListArtifactDeploymentRecords(ctx, "o", "sha256:abc")
	if err != nil {
		t.Errorf("ListArtifactDeploymentRecords returned error: %v", err)
	}

	want := &ArtifactDeploymentResponse{
		TotalCount: Ptr(1),
		DeploymentRecords: []*ArtifactDeploymentRecord{
			{ID: Ptr(int64(1)), RuntimeRisks: []string{"sensitive-data"}},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListArtifactDeploymentRecords returned %+v, want %+v", got, want)
	}
}

func TestOrganizationsService_ListArtifactStorageRecords(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/artifacts/sha256:abc/metadata/storage-records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1,"storage_records":[{"name":"libfoo"}]}`)
	})

	ctx := t.Context()
	got, _, err := client.Organizations.ListArtifactStorageRecords(ctx, "o", "sha256:abc")
	if err != nil {
		t.Errorf("ListArtifactStorageRecords returned error: %v", err)
	}

	want := &ArtifactStorageResponse{
		TotalCount:     Ptr(1),
		StorageRecords: []*ArtifactStorageRecord{{Name: Ptr("libfoo")}},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListArtifactStorageRecords returned %+v, want %+v", got, want)
	}
}
