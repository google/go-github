// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListHostedRunners(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"total_count": 2,
			"runners": [
				{
					"id": 5,
					"name": "My hosted ubuntu runner",
					"runner_group_id": 2,
					"platform": "linux-x64",
					"image": {
						"id": "ubuntu-20.04",
						"size": 86
					},
					"machine_size_details": {
						"id": "4-core",
						"cpu_cores": 4,
						"memory_gb": 16,
						"storage_gb": 150
					},
					"status": "Ready",
					"maximum_runners": 10,
					"public_ip_enabled": true,
					"public_ips": [
						{
							"enabled": true,
							"prefix": "20.80.208.150",
							"length": 31
						}
					],
					"last_active_on": "2022-10-09T23:39:01Z"
				},
				{
					"id": 7,
					"name": "My hosted Windows runner",
					"runner_group_id": 2,
					"platform": "win-x64",
					"image": {
						"id": "windows-latest",
						"size": 256
					},
					"machine_size_details": {
						"id": "8-core",
						"cpu_cores": 8,
						"memory_gb": 32,
						"storage_gb": 300
					},
					"status": "Ready",
					"maximum_runners": 20,
					"public_ip_enabled": false,
					"public_ips": [],
					"last_active_on": "2023-04-26T15:23:37Z"
				}
			]
		}`)
	})
	opts := &ListOptions{Page: 1, PerPage: 1}
	ctx := context.Background()
	hostedRunners, _, err := client.Actions.ListHostedRunners(ctx, "o", opts)
	if err != nil {
		t.Errorf("Actions.ListHostedRunners returned error: %v", err)
	}

	want := &HostedRunners{
		TotalCount: 2,
		Runners: []*HostedRunner{
			{
				ID:            Ptr(int64(5)),
				Name:          Ptr("My hosted ubuntu runner"),
				RunnerGroupID: Ptr(int64(2)),
				Platform:      Ptr("linux-x64"),
				Image: &HostedRunnerImageDetail{
					ID:   Ptr("ubuntu-20.04"),
					Size: Ptr(86),
				},
				MachineSizeDetails: &HostedRunnerMachineSpec{
					ID:        "4-core",
					CPUCores:  4,
					MemoryGB:  16,
					StorageGB: 150,
				},
				Status:          Ptr("Ready"),
				MaximumRunners:  Ptr(int64(10)),
				PublicIPEnabled: Ptr(true),
				PublicIPs: []*HostedRunnerPublicIP{
					{
						Enabled: true,
						Prefix:  "20.80.208.150",
						Length:  31,
					},
				},
				LastActiveOn: Ptr("2022-10-09T23:39:01Z"),
			},
			{
				ID:            Ptr(int64(7)),
				Name:          Ptr("My hosted Windows runner"),
				RunnerGroupID: Ptr(int64(2)),
				Platform:      Ptr("win-x64"),
				Image: &HostedRunnerImageDetail{
					ID:   Ptr("windows-latest"),
					Size: Ptr(256),
				},
				MachineSizeDetails: &HostedRunnerMachineSpec{
					ID:        "8-core",
					CPUCores:  8,
					MemoryGB:  32,
					StorageGB: 300,
				},
				Status:          Ptr("Ready"),
				MaximumRunners:  Ptr(int64(20)),
				PublicIPEnabled: Ptr(false),
				PublicIPs:       []*HostedRunnerPublicIP{},
				LastActiveOn:    Ptr("2023-04-26T15:23:37Z"),
			},
		},
	}
	if !cmp.Equal(hostedRunners, want) {
		t.Errorf("Actions.ListHostedRunners returned %+v, want %+v", hostedRunners, want)
	}

	const methodName = "ListHostedRunners"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListHostedRunners(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListHostedRunners(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateHostedRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 5,
			"name": "My hosted ubuntu runner",
			"runner_group_id": 2,
			"platform": "linux-x64",
			"image": {
				"id": "ubuntu-20.04",
				"size": 86
			},
			"machine_size_details": {
				"id": "4-core",
				"cpu_cores": 4,
				"memory_gb": 16,
				"storage_gb": 150
			},
			"status": "Ready",
			"maximum_runners": 10,
			"public_ip_enabled": true,
			"public_ips": [
				{
					"enabled": true,
					"prefix": "20.80.208.150",
					"length": 31
				}
			],
			"last_active_on": "2022-10-09T23:39:01Z"
			}`)
	})

	ctx := context.Background()
	req := &CreateHostedRunnerRequest{
		Name: "My Hosted runner",
		Image: HostedRunnerImage{
			ID:      "ubuntu-latest",
			Source:  "github",
			Version: "latest",
		},
		RunnerGroupID:  1,
		Size:           "4-core",
		MaximumRunners: 50,
		EnableStaticIP: false,
	}
	hostedRunner, _, err := client.Actions.CreateHostedRunner(ctx, "o", req)
	if err != nil {
		t.Errorf("Actions.CreateHostedRunner returned error: %v", err)
	}

	want := &HostedRunner{
		ID:            Ptr(int64(5)),
		Name:          Ptr("My hosted ubuntu runner"),
		RunnerGroupID: Ptr(int64(2)),
		Platform:      Ptr("linux-x64"),
		Image: &HostedRunnerImageDetail{
			ID:   Ptr("ubuntu-20.04"),
			Size: Ptr(86),
		},
		MachineSizeDetails: &HostedRunnerMachineSpec{
			ID:        "4-core",
			CPUCores:  4,
			MemoryGB:  16,
			StorageGB: 150,
		},
		Status:          Ptr("Ready"),
		MaximumRunners:  Ptr(int64(10)),
		PublicIPEnabled: Ptr(true),
		PublicIPs: []*HostedRunnerPublicIP{
			{
				Enabled: true,
				Prefix:  "20.80.208.150",
				Length:  31,
			},
		},
		LastActiveOn: Ptr("2022-10-09T23:39:01Z"),
	}

	if !cmp.Equal(hostedRunner, want) {
		t.Errorf("Actions.CreateHostedRunner returned %+v, want %+v", hostedRunner, want)
	}

	const methodName = "CreateHostedRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.CreateHostedRunner(ctx, "\n", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.CreateHostedRunner(ctx, "o", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunnerGithubOwnedImages(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/images/github-owned", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"total_count": 1,
			"images": [
				{
					"id": "ubuntu-20.04",
					"platform": "linux-x64",
					"size_gb": 86,
					"display_name": "20.04",
					"source": "github"
				}
			]
			}`)
	})

	ctx := context.Background()
	hostedRunnerImages, _, err := client.Actions.GetHostedRunnerGithubOwnedImages(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetHostedRunnerGithubOwnedImages returned error: %v", err)
	}

	want := &HostedRunnerImages{
		TotalCount: 1,
		Images: []*HostedRunnerImageSpecs{
			{
				ID:          "ubuntu-20.04",
				Platform:    "linux-x64",
				SizeGB:      86,
				DisplayName: "20.04",
				Source:      "github",
			},
		},
	}

	if !cmp.Equal(hostedRunnerImages, want) {
		t.Errorf("Actions.GetHostedRunnerGithubOwnedImages returned %+v, want %+v", hostedRunnerImages, want)
	}

	const methodName = "GetHostedRunnerGithubOwnedImages"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunnerGithubOwnedImages(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunnerGithubOwnedImages(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunnerPartnerImages(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/images/partner", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"total_count": 1,
			"images": [
				{
					"id": "ubuntu-20.04",
					"platform": "linux-x64",
					"size_gb": 86,
					"display_name": "20.04",
					"source": "partner"
				}
			]
			}`)
	})

	ctx := context.Background()
	hostedRunnerImages, _, err := client.Actions.GetHostedRunnerPartnerImages(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetHostedRunnerPartnerImages returned error: %v", err)
	}

	want := &HostedRunnerImages{
		TotalCount: 1,
		Images: []*HostedRunnerImageSpecs{
			{
				ID:          "ubuntu-20.04",
				Platform:    "linux-x64",
				SizeGB:      86,
				DisplayName: "20.04",
				Source:      "partner",
			},
		},
	}

	if !cmp.Equal(hostedRunnerImages, want) {
		t.Errorf("Actions.GetHostedRunnerPartnerImages returned %+v, want %+v", hostedRunnerImages, want)
	}

	const methodName = "GetHostedRunnerPartnerImages"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunnerPartnerImages(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunnerPartnerImages(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunnerLimits(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"public_ips": {
				"current_usage": 17,
				"maximum": 50
			}
		}`)
	})

	ctx := context.Background()
	publicIPLimits, _, err := client.Actions.GetHostedRunnerLimits(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetPartnerImages returned error: %v", err)
	}

	want := &HostedRunnerPublicIPLimits{
		PublicIPs: &PublicIPUsage{
			CurrentUsage: 17,
			Maximum:      50,
		},
	}

	if !cmp.Equal(publicIPLimits, want) {
		t.Errorf("Actions.GetHostedRunnerLimits returned %+v, want %+v", publicIPLimits, want)
	}

	const methodName = "GetHostedRunnerLimits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunnerLimits(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunnerLimits(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunnerMachineSpecs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/machine-sizes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"total_count": 1,
			"machine_specs": [
				{
					"id": "4-core",
 	 				"cpu_cores": 4,
  					"memory_gb": 16,
  					"storage_gb": 150
				}
			]
			}`)
	})

	ctx := context.Background()
	machineSpecs, _, err := client.Actions.GetHostedRunnerMachineSpecs(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetHostedRunnerMachineSpecs returned error: %v", err)
	}
	want := &HostedRunnerMachineSpecs{
		TotalCount: 1,
		MachineSpecs: []*HostedRunnerMachineSpec{
			{
				ID:        "4-core",
				CPUCores:  4,
				MemoryGB:  16,
				StorageGB: 150,
			},
		},
	}

	if !cmp.Equal(machineSpecs, want) {
		t.Errorf("Actions.GetHostedRunnerMachineSpecs returned %+v, want %+v", machineSpecs, want)
	}

	const methodName = "GetHostedRunnerMachineSpecs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunnerMachineSpecs(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunnerMachineSpecs(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunnerPlatforms(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/platforms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"total_count": 1,
  			"platforms": [
    			"linux-x64",
    			"win-x64"
  			]
		}`)
	})

	ctx := context.Background()
	platforms, _, err := client.Actions.GetHostedRunnerPlatforms(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetHostedRunnerPlatforms returned error: %v", err)
	}
	want := &HostedRunnerPlatforms{
		TotalCount: 1,
		Platforms: []string{
			"linux-x64",
			"win-x64",
		},
	}

	if !cmp.Equal(platforms, want) {
		t.Errorf("Actions.GetHostedRunnerPlatforms returned %+v, want %+v", platforms, want)
	}

	const methodName = "GetHostedRunnerPlatforms"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunnerPlatforms(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunnerPlatforms(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetHostedRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/23", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 5,
			"name": "My hosted ubuntu runner",
			"runner_group_id": 2,
			"platform": "linux-x64",
			"image": {
				"id": "ubuntu-20.04",
				"size": 86
			},
			"machine_size_details": {
				"id": "4-core",
				"cpu_cores": 4,
				"memory_gb": 16,
				"storage_gb": 150
			},
			"status": "Ready",
			"maximum_runners": 10,
			"public_ip_enabled": true,
			"public_ips": [
				{
				"enabled": true,
				"prefix": "20.80.208.150",
				"length": 31
				}
			],
			"last_active_on": "2022-10-09T23:39:01Z"
		}`)
	})

	ctx := context.Background()
	hostedRunner, _, err := client.Actions.GetHostedRunner(ctx, "o", 23)
	if err != nil {
		t.Errorf("Actions.GetHostedRunner returned error: %v", err)
	}

	want := &HostedRunner{
		ID:            Ptr(int64(5)),
		Name:          Ptr("My hosted ubuntu runner"),
		RunnerGroupID: Ptr(int64(2)),
		Platform:      Ptr("linux-x64"),
		Image: &HostedRunnerImageDetail{
			ID:   Ptr("ubuntu-20.04"),
			Size: Ptr(86),
		},
		MachineSizeDetails: &HostedRunnerMachineSpec{
			ID:        "4-core",
			CPUCores:  4,
			MemoryGB:  16,
			StorageGB: 150,
		},
		Status:          Ptr("Ready"),
		MaximumRunners:  Ptr(int64(10)),
		PublicIPEnabled: Ptr(true),
		PublicIPs: []*HostedRunnerPublicIP{
			{
				Enabled: true,
				Prefix:  "20.80.208.150",
				Length:  31,
			},
		},
		LastActiveOn: Ptr("2022-10-09T23:39:01Z"),
	}

	if !cmp.Equal(hostedRunner, want) {
		t.Errorf("Actions.GetHostedRunner returned %+v, want %+v", hostedRunner, want)
	}

	const methodName = "GetHostedRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetHostedRunner(ctx, "\n", 23)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetHostedRunner(ctx, "o", 23)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_UpdateHostedRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/23", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{
			"id": 5,
			"name": "My hosted ubuntu runner",
			"runner_group_id": 2,
			"platform": "linux-x64",
			"image": {
				"id": "ubuntu-20.04",
				"size": 86
			},
			"machine_size_details": {
				"id": "4-core",
				"cpu_cores": 4,
				"memory_gb": 16,
				"storage_gb": 150
			},
			"status": "Ready",
			"maximum_runners": 10,
			"public_ip_enabled": true,
			"public_ips": [
				{
					"enabled": true,
					"prefix": "20.80.208.150",
					"length": 31
				}
			],
			"last_active_on": "2022-10-09T23:39:01Z"
			}`)
	})

	ctx := context.Background()
	req := UpdateHostedRunnerRequest{
		Name:           "My larger runner",
		RunnerGroupID:  1,
		MaximumRunners: 50,
		EnableStaticIP: false,
		ImageVersion:   "1.0.0",
	}
	hostedRunner, _, err := client.Actions.UpdateHostedRunner(ctx, "o", 23, req)
	if err != nil {
		t.Errorf("Actions.UpdateHostedRunner returned error: %v", err)
	}

	want := &HostedRunner{
		ID:            Ptr(int64(5)),
		Name:          Ptr("My hosted ubuntu runner"),
		RunnerGroupID: Ptr(int64(2)),
		Platform:      Ptr("linux-x64"),
		Image: &HostedRunnerImageDetail{
			ID:   Ptr("ubuntu-20.04"),
			Size: Ptr(86),
		},
		MachineSizeDetails: &HostedRunnerMachineSpec{
			ID:        "4-core",
			CPUCores:  4,
			MemoryGB:  16,
			StorageGB: 150,
		},
		Status:          Ptr("Ready"),
		MaximumRunners:  Ptr(int64(10)),
		PublicIPEnabled: Ptr(true),
		PublicIPs: []*HostedRunnerPublicIP{
			{
				Enabled: true,
				Prefix:  "20.80.208.150",
				Length:  31,
			},
		},
		LastActiveOn: Ptr("2022-10-09T23:39:01Z"),
	}

	if !cmp.Equal(hostedRunner, want) {
		t.Errorf("Actions.UpdateHostedRunner returned %+v, want %+v", hostedRunner, want)
	}

	const methodName = "UpdateHostedRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.UpdateHostedRunner(ctx, "\n", 23, req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.UpdateHostedRunner(ctx, "o", 23, req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_DeleteHostedRunner(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/hosted-runners/23", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"id": 5,
			"name": "My hosted ubuntu runner",
			"runner_group_id": 2,
			"platform": "linux-x64",
			"image": {
				"id": "ubuntu-20.04",
				"size": 86
			},
			"machine_size_details": {
				"id": "4-core",
				"cpu_cores": 4,
				"memory_gb": 16,
				"storage_gb": 150
			},
			"status": "Ready",
			"maximum_runners": 10,
			"public_ip_enabled": true,
			"public_ips": [
				{
				"enabled": true,
				"prefix": "20.80.208.150",
				"length": 31
				}
			],
			"last_active_on": "2022-10-09T23:39:01Z"
		}`)
	})

	ctx := context.Background()
	hostedRunner, _, err := client.Actions.DeleteHostedRunner(ctx, "o", 23)
	if err != nil {
		t.Errorf("Actions.GetHostedRunner returned error: %v", err)
	}

	want := &HostedRunner{
		ID:            Ptr(int64(5)),
		Name:          Ptr("My hosted ubuntu runner"),
		RunnerGroupID: Ptr(int64(2)),
		Platform:      Ptr("linux-x64"),
		Image: &HostedRunnerImageDetail{
			ID:   Ptr("ubuntu-20.04"),
			Size: Ptr(86),
		},
		MachineSizeDetails: &HostedRunnerMachineSpec{
			ID:        "4-core",
			CPUCores:  4,
			MemoryGB:  16,
			StorageGB: 150,
		},
		Status:          Ptr("Ready"),
		MaximumRunners:  Ptr(int64(10)),
		PublicIPEnabled: Ptr(true),
		PublicIPs: []*HostedRunnerPublicIP{
			{
				Enabled: true,
				Prefix:  "20.80.208.150",
				Length:  31,
			},
		},
		LastActiveOn: Ptr("2022-10-09T23:39:01Z"),
	}

	if !cmp.Equal(hostedRunner, want) {
		t.Errorf("Actions.DeleteHostedRunner returned %+v, want %+v", hostedRunner, want)
	}

	const methodName = "DeleteHostedRunner"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.DeleteHostedRunner(ctx, "\n", 23)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.DeleteHostedRunner(ctx, "o", 23)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
