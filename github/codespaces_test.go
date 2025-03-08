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

func TestCodespacesService_ListInRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/codespaces", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"total_count":2,"codespaces":[{"id":1,"name":"monalisa-octocat-hello-world-g4wpq6h95q","environment_id":"26a7c758-7299-4a73-b978-5a92a7ae98a0","owner":{"login":"octocat"},"billable_owner":{"login":"octocat"},"repository":{"id":1296269},"machine":{"name":"standardLinux","display_name":"4 cores, 8 GB RAM, 64 GB storage","operating_system":"linux","storage_in_bytes":68719476736,"memory_in_bytes":8589934592,"cpus":4},"prebuild":false,"devcontainer_path":".devcontainer/devcontainer.json","created_at":"2021-10-14T00:53:30-06:00","updated_at":"2021-10-14T00:53:32-06:00","last_used_at":"2021-10-14T00:53:30-06:00","state":"Available","url":"https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q","git_status":{"ahead":0,"behind":0,"has_unpushed_changes":false,"has_uncommitted_changes":false,"ref":"main"},"location":"WestUs2","idle_timeout_minutes":60,"web_url":"https://monalisa-octocat-hello-world-g4wpq6h95q.github.dev","machines_url":"https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/machines","start_url":"https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/start","stop_url":"https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/stop","recent_folders":["testfolder1","testfolder2"]},{"id":2}]}`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	codespaces, _, err := client.Codespaces.ListInRepo(ctx, "owner", "repo", opt)
	if err != nil {
		t.Errorf("Codespaces.ListInRepo returned error: %v", err)
	}

	want := &ListCodespaces{TotalCount: Ptr(2), Codespaces: []*Codespace{
		{
			ID:            Ptr(int64(1)),
			Name:          Ptr("monalisa-octocat-hello-world-g4wpq6h95q"),
			EnvironmentID: Ptr("26a7c758-7299-4a73-b978-5a92a7ae98a0"),
			Owner: &User{
				Login: Ptr("octocat"),
			},
			BillableOwner: &User{
				Login: Ptr("octocat"),
			},
			Repository: &Repository{
				ID: Ptr(int64(1296269)),
			},
			Machine: &CodespacesMachine{
				Name:            Ptr("standardLinux"),
				DisplayName:     Ptr("4 cores, 8 GB RAM, 64 GB storage"),
				OperatingSystem: Ptr("linux"),
				StorageInBytes:  Ptr(int64(68719476736)),
				MemoryInBytes:   Ptr(int64(8589934592)),
				CPUs:            Ptr(4),
			},
			Prebuild:         Ptr(false),
			DevcontainerPath: Ptr(".devcontainer/devcontainer.json"),
			CreatedAt:        &Timestamp{time.Date(2021, 10, 14, 0, 53, 30, 0, time.FixedZone("", -6*60*60))},
			UpdatedAt:        &Timestamp{time.Date(2021, 10, 14, 0, 53, 32, 0, time.FixedZone("", -6*60*60))},
			LastUsedAt:       &Timestamp{time.Date(2021, 10, 14, 0, 53, 30, 0, time.FixedZone("", -6*60*60))},
			State:            Ptr("Available"),
			URL:              Ptr("https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q"),
			GitStatus: &CodespacesGitStatus{
				Ahead:                 Ptr(0),
				Behind:                Ptr(0),
				HasUnpushedChanges:    Ptr(false),
				HasUncommittedChanges: Ptr(false),
				Ref:                   Ptr("main"),
			},
			Location:           Ptr("WestUs2"),
			IdleTimeoutMinutes: Ptr(60),
			WebURL:             Ptr("https://monalisa-octocat-hello-world-g4wpq6h95q.github.dev"),
			MachinesURL:        Ptr("https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/machines"),
			StartURL:           Ptr("https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/start"),
			StopURL:            Ptr("https://api.github.com/user/codespaces/monalisa-octocat-hello-world-g4wpq6h95q/stop"),
			RecentFolders: []string{
				"testfolder1",
				"testfolder2",
			},
		},
		{
			ID: Ptr(int64(2)),
		},
	}}
	if !cmp.Equal(codespaces, want) {
		t.Errorf("Codespaces.ListInRepo returned %+v, want %+v", codespaces, want)
	}

	const methodName = "ListInRepo"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.ListInRepo(ctx, "", "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_List(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/codespaces", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":          "1",
			"per_page":      "2",
			"repository_id": "1296269",
		})
		fmt.Fprint(w, `{"total_count":1,"codespaces":[{"id":1, "repository": {"id": 1296269}}]}`)
	})

	opt := &ListCodespacesOptions{ListOptions: ListOptions{Page: 1, PerPage: 2}, RepositoryID: 1296269}
	ctx := context.Background()
	codespaces, _, err := client.Codespaces.List(ctx, opt)
	if err != nil {
		t.Errorf("Codespaces.List returned error: %v", err)
	}

	want := &ListCodespaces{TotalCount: Ptr(1), Codespaces: []*Codespace{
		{
			ID: Ptr(int64(1)),
			Repository: &Repository{
				ID: Ptr(int64(1296269)),
			},
		},
	}}
	if !cmp.Equal(codespaces, want) {
		t.Errorf("Codespaces.ListInRepo returned %+v, want %+v", codespaces, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_CreateInRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/codespaces", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"ref":"main","geo":"WestUs2","machine":"standardLinux","idle_timeout_minutes":60}`+"\n")
		fmt.Fprint(w, `{"id":1, "repository": {"id": 1296269}}`)
	})
	input := &CreateCodespaceOptions{
		Ref:                Ptr("main"),
		Geo:                Ptr("WestUs2"),
		Machine:            Ptr("standardLinux"),
		IdleTimeoutMinutes: Ptr(60),
	}
	ctx := context.Background()
	codespace, _, err := client.Codespaces.CreateInRepo(ctx, "owner", "repo", input)
	if err != nil {
		t.Errorf("Codespaces.CreateInRepo returned error: %v", err)
	}
	want := &Codespace{
		ID: Ptr(int64(1)),
		Repository: &Repository{
			ID: Ptr(int64(1296269)),
		},
	}

	if !cmp.Equal(codespace, want) {
		t.Errorf("Codespaces.CreateInRepo returned %+v, want %+v", codespace, want)
	}

	const methodName = "CreateInRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Codespaces.CreateInRepo(ctx, "\n", "", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.CreateInRepo(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_Start(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/codespaces/codespace_1/start", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1, "repository": {"id": 1296269}}`)
	})
	ctx := context.Background()
	codespace, _, err := client.Codespaces.Start(ctx, "codespace_1")
	if err != nil {
		t.Errorf("Codespaces.Start returned error: %v", err)
	}
	want := &Codespace{
		ID: Ptr(int64(1)),
		Repository: &Repository{
			ID: Ptr(int64(1296269)),
		},
	}

	if !cmp.Equal(codespace, want) {
		t.Errorf("Codespaces.Start returned %+v, want %+v", codespace, want)
	}

	const methodName = "Start"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Codespaces.Start(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.Start(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_Stop(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/codespaces/codespace_1/stop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1, "repository": {"id": 1296269}}`)
	})
	ctx := context.Background()
	codespace, _, err := client.Codespaces.Stop(ctx, "codespace_1")
	if err != nil {
		t.Errorf("Codespaces.Stop returned error: %v", err)
	}
	want := &Codespace{
		ID: Ptr(int64(1)),
		Repository: &Repository{
			ID: Ptr(int64(1296269)),
		},
	}

	if !cmp.Equal(codespace, want) {
		t.Errorf("Codespaces.Stop returned %+v, want %+v", codespace, want)
	}

	const methodName = "Stop"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Codespaces.Stop(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Codespaces.Stop(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodespacesService_Delete(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/codespaces/codespace_1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Codespaces.Delete(ctx, "codespace_1")
	if err != nil {
		t.Errorf("Codespaces.Delete return error: %v", err)
	}

	const methodName = "Delete"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Codespaces.Delete(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Codespaces.Delete(ctx, "codespace_1")
	})
}
