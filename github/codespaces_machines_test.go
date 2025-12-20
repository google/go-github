// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestCodespacesService_ListMachinesTypesForRepository(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/codespaces/machines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"total_count": 1,
			"machines": [
				{
					"name": "standardLinux",
					"display_name": "4 cores, 8 GB RAM, 64 GB storage",
					"operating_system": "linux",
					"storage_in_bytes": 68719476736,
					"memory_in_bytes": 17179869184,
					"cpus": 4,
					"prebuild_availability": "ready"
				}
			]
		}`)
	})

	ctx := t.Context()
	opt := &ListMachinesOptions{
		Ref:      Ptr("main"),
		Location: Ptr("WestUs2"),
		ClientIP: Ptr("1.2.3.4"),
	}

	got, _, err := client.Codespaces.ListMachinesTypesForRepository(
		ctx,
		"owner",
		"repo",
		opt,
	)
	if err != nil {
		t.Fatalf("Codespaces.ListMachinesTypesForRepository returned error: %v", err)
	}

	want := &CodespaceMachines{
		TotalCount: 1,
		Machines: []*Machine{
			{
				Name:                 "standardLinux",
				DisplayName:          "4 cores, 8 GB RAM, 64 GB storage",
				OperatingSystem:      "linux",
				StorageInBytes:       68719476736,
				MemoryInBytes:        17179869184,
				CPUs:                 4,
				PrebuildAvailability: "ready",
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Codespaces.ListMachinesTypesForRepository returned %+v, want %+v", got, want)
	}

	const methodName = "ListMachinesTypesForRepository"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.ListMachinesTypesForRepository(ctx, "/n", "/n", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_ListMachinesTypesForCodespace(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/codespaces/codespace_1/machines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"total_count": 1,
			"machines": [
				{
					"name": "standardLinux",
					"display_name": "4 cores, 8 GB RAM, 64 GB storage",
					"operating_system": "linux",
					"storage_in_bytes": 68719476736,
					"memory_in_bytes": 17179869184,
					"cpus": 4,
					"prebuild_availability": "ready"
				}
			]
		}`)
	})

	ctx := t.Context()
	got, _, err := client.Codespaces.ListMachinesTypesForCodespace(ctx, "codespace_1")
	if err != nil {
		t.Fatalf("Codespaces.ListMachinesTypesForCodespace returned error: %v", err)
	}

	want := &CodespaceMachines{
		TotalCount: 1,
		Machines: []*Machine{
			{
				Name:                 "standardLinux",
				DisplayName:          "4 cores, 8 GB RAM, 64 GB storage",
				OperatingSystem:      "linux",
				StorageInBytes:       68719476736,
				MemoryInBytes:        17179869184,
				CPUs:                 4,
				PrebuildAvailability: "ready",
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Codespaces.ListMachinesTypesForCodespace returned %+v, want %+v", got, want)
	}

	const methodName = "ListMachinesTypesForCodespace"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.ListMachinesTypesForCodespace(ctx, "/n")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
