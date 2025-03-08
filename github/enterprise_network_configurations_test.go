// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_ListEnterpriseNetworkConfigurations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc(" /enterprises/e/network-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "3", "per_page": "2"})
		fmt.Fprintf(w, `{
		  "total_count": 2,
		  "network_configurations": [
			{
			  "id": "123456789ABCDEF",
			  "name": "configuration_one",
			  "compute_service": "actions",
			  "network_settings_ids": [
				"23456789ABDCEF1",
				"3456789ABDCEF12"
			  ],
			  "created_on": "2024-04-09T17:30:15Z"
			},
			{
			  "id": "456789ABDCEF123",
			  "name": "configuration_two",
			  "compute_service": "none",
			  "network_settings_ids": [
				"56789ABDCEF1234",
				"6789ABDCEF12345"
			  ],
			  "created_on": "2024-11-02T12:30:30Z"
			}
		  ]
		}`)
	})

	ctx := context.Background()

	opts := &ListOptions{Page: 3, PerPage: 2}
	configurations, _, err := client.Enterprise.ListEnterpriseNetworkConfigurations(ctx, "e", opts)
	if err != nil {
		t.Errorf("Enterprise.ListEnterpriseNetworkConfigurations returned error: %v", err)
	}

	want := &EnterpriseNetworkConfigurations{
		TotalCount: Ptr(int64(2)),
		NetworkConfigurations: []*EnterpriseNetworkConfiguration{
			{
				ID:                 Ptr("123456789ABCDEF"),
				Name:               Ptr("configuration_one"),
				ComputeService:     Ptr(ComputeService("actions")),
				NetworkSettingsIDs: []string{"23456789ABDCEF1", "3456789ABDCEF12"},
				CreatedOn:          &Timestamp{time.Date(2024, 4, 9, 17, 30, 15, 0, time.UTC)},
			},
			{
				ID:                 Ptr("456789ABDCEF123"),
				Name:               Ptr("configuration_two"),
				ComputeService:     Ptr(ComputeService("none")),
				NetworkSettingsIDs: []string{"56789ABDCEF1234", "6789ABDCEF12345"},
				CreatedOn:          &Timestamp{time.Date(2024, 11, 2, 12, 30, 30, 0, time.UTC)},
			},
		},
	}
	if !cmp.Equal(configurations, want) {
		t.Errorf("Enterprise.ListEnterpriseNetworkConfigurations mismatch (-want +got):\n%s", cmp.Diff(want, configurations))
	}

	const methodName = "ListEnterpriseNetworkConfigurations"
	testBadOptions(t, methodName, func() error {
		_, _, err = client.Enterprise.ListEnterpriseNetworkConfigurations(ctx, "\ne", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListEnterpriseNetworkConfigurations(ctx, "e", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateEnterpriseNetworkConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/network-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, `{
		  "id": "123456789ABCDEF",
		  "name": "configuration_one",
		  "compute_service": "actions",
		  "network_settings_ids": [
			"23456789ABDCEF1",
			"3456789ABDCEF12"
		  ],
		  "created_on": "2024-04-09T17:30:15Z"
		}`)
	})

	ctx := context.Background()

	req := EnterpriseNetworkConfigurationRequest{
		Name:           Ptr("configuration_one"),
		ComputeService: Ptr(ComputeService("actions")),
		NetworkSettingsIDs: []string{
			"23456789ABDCEF1",
			"3456789ABDCEF12",
		},
	}
	configuration, _, err := client.Enterprise.CreateEnterpriseNetworkConfiguration(ctx, "e", req)
	if err != nil {
		t.Errorf("Enterprise.CreateEnterpriseNetworkConfiguration returned error: %v", err)
	}

	want := &EnterpriseNetworkConfiguration{
		ID:                 Ptr("123456789ABCDEF"),
		Name:               Ptr("configuration_one"),
		ComputeService:     Ptr(ComputeService("actions")),
		NetworkSettingsIDs: []string{"23456789ABDCEF1", "3456789ABDCEF12"},
		CreatedOn:          &Timestamp{time.Date(2024, 4, 9, 17, 30, 15, 0, time.UTC)},
	}
	if !cmp.Equal(configuration, want) {
		t.Errorf("Enterprise.CreateEnterpriseNetworkConfiguration mismatch (-want +got):\n%s", cmp.Diff(want, configuration))
	}

	const methodName = "CreateEnterpriseNetworkConfiguration"
	testBadOptions(t, methodName, func() error {
		_, _, err = client.Enterprise.CreateEnterpriseNetworkConfiguration(ctx, "\ne", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateEnterpriseNetworkConfiguration(ctx, "e", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetEnterpriseNetworkConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/network-configurations/123456789ABCDEF", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{
		  "id": "123456789ABCDEF",
		  "name": "configuration_one",
		  "compute_service": "actions",
		  "network_settings_ids": [
			"23456789ABDCEF1",
			"3456789ABDCEF12"
		  ],
		  "created_on": "2024-12-10T19:00:15Z"
		}`)
	})

	ctx := context.Background()
	configuration, _, err := client.Enterprise.GetEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF")
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseNetworkConfiguration returned err: %v", err)
	}

	want := &EnterpriseNetworkConfiguration{
		ID:                 Ptr("123456789ABCDEF"),
		Name:               Ptr("configuration_one"),
		ComputeService:     Ptr(ComputeService("actions")),
		NetworkSettingsIDs: []string{"23456789ABDCEF1", "3456789ABDCEF12"},
		CreatedOn:          &Timestamp{time.Date(2024, 12, 10, 19, 00, 15, 0, time.UTC)},
	}
	if !cmp.Equal(configuration, want) {
		t.Errorf("Enterprise.GetEnterpriseNetworkConfiguration mismatch (-want +got):\n%s", cmp.Diff(want, configuration))
	}

	const methodName = "GetEnterpriseNetworkConfiguration"
	testBadOptions(t, methodName, func() error {
		_, _, err = client.Enterprise.GetEnterpriseNetworkConfiguration(ctx, "\ne", "123456789ABCDEF")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateEnterpriseNetworkConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/network-configurations/123456789ABCDEF", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprintf(w, `{
		  "id": "123456789ABCDEF",
		  "name": "updated_configuration_one",
		  "compute_service": "none",
		  "network_settings_ids": [
			"456789ABDCEF123",
			"56789ABDCEF1234"
		  ],
		  "created_on": "2024-12-10T19:00:15Z"
		}`)
	})

	ctx := context.Background()
	req := EnterpriseNetworkConfigurationRequest{
		Name: Ptr("updated_configuration_one"),
		NetworkSettingsIDs: []string{
			"456789ABDCEF123",
			"56789ABDCEF1234",
		},
		ComputeService: Ptr(ComputeService("none")),
	}
	configuration, _, err := client.Enterprise.UpdateEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF", req)
	if err != nil {
		t.Errorf("Enterprise.UpdateEnterpriseNetworkConfiguration returned error %v", err)
	}

	want := &EnterpriseNetworkConfiguration{
		ID:                 Ptr("123456789ABCDEF"),
		Name:               Ptr("updated_configuration_one"),
		ComputeService:     Ptr(ComputeService("none")),
		NetworkSettingsIDs: []string{"456789ABDCEF123", "56789ABDCEF1234"},
		CreatedOn:          &Timestamp{time.Date(2024, 12, 10, 19, 00, 15, 0, time.UTC)},
	}
	if !cmp.Equal(configuration, want) {
		t.Errorf("Enterprise.UpdateEnterpriseNetworkConfiguration mismatch (-want +get)\n%s", cmp.Diff(want, configuration))
	}

	const methodName = "UpdateEnterpriseNetworkConfiguration"
	testBadOptions(t, methodName, func() error {
		_, _, err = client.Enterprise.UpdateEnterpriseNetworkConfiguration(ctx, "\ne", "123456789ABCDEF", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteEnterpriseNetworkConfiguration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/network-configurations/123456789ABCDEF", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.DeleteEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF")
	if err != nil {
		t.Errorf("Enterprise.DeleteEnterpriseNetworkConfiguration returned error %v", err)
	}

	const methodName = "DeleteEnterpriseNetworkConfiguration"
	testBadOptions(t, methodName, func() error {
		_, err = client.Enterprise.DeleteEnterpriseNetworkConfiguration(ctx, "\ne", "123456789ABCDEF")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteEnterpriseNetworkConfiguration(ctx, "e", "123456789ABCDEF")
	})
}

func TestEnterpriseService_GetEnterpriseNetworkSettingsResource(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/network-settings/123456789ABCDEF", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{
			"id": "220F78DACB92BBFBC5E6F22DE1CCF52309D",
			"network_configuration_id": "934E208B3EE0BD60CF5F752C426BFB53562",
			"name": "my_network_settings",
			"subnet_id": "/subscriptions/14839728-3ad9-43ab-bd2b-fa6ad0f75e2a/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet/subnets/my-subnet",
			"region": "germanywestcentral"
		}`)
	})

	ctx := context.Background()
	resource, _, err := client.Enterprise.GetEnterpriseNetworkSettingsResource(ctx, "e", "123456789ABCDEF")
	if err != nil {
		t.Errorf("Enterprise.GetEnterpriseNetworkSettingsResource returned error %v", err)
	}

	want := &EnterpriseNetworkSettingsResource{
		ID:                     Ptr("220F78DACB92BBFBC5E6F22DE1CCF52309D"),
		NetworkConfigurationID: Ptr("934E208B3EE0BD60CF5F752C426BFB53562"),
		Name:                   Ptr("my_network_settings"),
		SubnetID:               Ptr("/subscriptions/14839728-3ad9-43ab-bd2b-fa6ad0f75e2a/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet/subnets/my-subnet"),
		Region:                 Ptr("germanywestcentral"),
	}
	if !cmp.Equal(resource, want) {
		t.Errorf("Enterprise.GetEnterpriseNetworkSettingsResource mistach (-want +got):\n%s", cmp.Diff(want, resource))
	}

	const methodName = "GetEnterpriseNetworkSettingsResource"
	testBadOptions(t, methodName, func() error {
		_, _, err = client.Enterprise.GetEnterpriseNetworkSettingsResource(ctx, "\ne", "123456789ABCDEF")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetEnterpriseNetworkSettingsResource(ctx, "e", "123456789ABCDEF")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
