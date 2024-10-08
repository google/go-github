// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMigrationService_StartMigration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		w.WriteHeader(http.StatusCreated)
		assertWrite(t, w, migrationJSON)
	})

	opt := &MigrationOptions{
		LockRepositories:   true,
		ExcludeAttachments: false,
	}
	ctx := context.Background()
	got, _, err := client.Migrations.StartMigration(ctx, "o", []string{"r"}, opt)
	if err != nil {
		t.Errorf("StartMigration returned error: %v", err)
	}
	if want := wantMigration; !cmp.Equal(got, want) {
		t.Errorf("StartMigration = %+v, want %+v", got, want)
	}

	const methodName = "StartMigration"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.StartMigration(ctx, "\n", []string{"\n"}, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.StartMigration(ctx, "o", []string{"r"}, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_ListMigrations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(fmt.Sprintf("[%s]", migrationJSON)))
	})

	ctx := context.Background()
	got, _, err := client.Migrations.ListMigrations(ctx, "o", &ListOptions{Page: 1, PerPage: 2})
	if err != nil {
		t.Errorf("ListMigrations returned error: %v", err)
	}
	if want := []*Migration{wantMigration}; !cmp.Equal(got, want) {
		t.Errorf("ListMigrations = %+v, want %+v", got, want)
	}

	const methodName = "ListMigrations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.ListMigrations(ctx, "\n", &ListOptions{Page: 1, PerPage: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.ListMigrations(ctx, "o", &ListOptions{Page: 1, PerPage: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_MigrationStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, migrationJSON)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.MigrationStatus(ctx, "o", 1)
	if err != nil {
		t.Errorf("MigrationStatus returned error: %v", err)
	}
	if want := wantMigration; !cmp.Equal(got, want) {
		t.Errorf("MigrationStatus = %+v, want %+v", got, want)
	}

	const methodName = "MigrationStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.MigrationStatus(ctx, "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.MigrationStatus(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_MigrationArchiveURL(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations/1/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		http.Redirect(w, r, "/yo", http.StatusFound)
	})
	mux.HandleFunc("/yo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte("0123456789abcdef"))
	})

	ctx := context.Background()
	got, err := client.Migrations.MigrationArchiveURL(ctx, "o", 1)
	if err != nil {
		t.Errorf("MigrationStatus returned error: %v", err)
	}
	if want := "/yo"; !strings.HasSuffix(got, want) {
		t.Errorf("MigrationArchiveURL = %+v, want %+v", got, want)
	}

	const methodName = "MigrationArchiveURL"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Migrations.MigrationArchiveURL(ctx, "\n", -1)
		return err
	})
}

func TestMigrationService_DeleteMigration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations/1/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Migrations.DeleteMigration(ctx, "o", 1); err != nil {
		t.Errorf("DeleteMigration returned error: %v", err)
	}

	const methodName = "DeleteMigration"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Migrations.DeleteMigration(ctx, "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Migrations.DeleteMigration(ctx, "o", 1)
	})
}

func TestMigrationService_UnlockRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/migrations/1/repos/r/lock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeMigrationsPreview)

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Migrations.UnlockRepo(ctx, "o", 1, "r"); err != nil {
		t.Errorf("UnlockRepo returned error: %v", err)
	}

	const methodName = "UnlockRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Migrations.UnlockRepo(ctx, "\n", -1, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Migrations.UnlockRepo(ctx, "o", 1, "r")
	})
}

var migrationJSON = []byte(`{
  "id": 79,
  "guid": "0b989ba4-242f-11e5-81e1-c7b6966d2516",
  "state": "pending",
  "lock_repositories": true,
  "exclude_attachments": false,
  "url": "https://api.github.com/orgs/octo-org/migrations/79",
  "created_at": "2015-07-06T15:33:38-07:00",
  "updated_at": "2015-07-06T15:33:38-07:00",
  "repositories": [
    {
      "id": 1296269,
      "name": "Hello-World",
      "full_name": "octocat/Hello-World",
      "description": "This your first repo!"
    }
  ]
}`)

var wantMigration = &Migration{
	ID:                 Int64(79),
	GUID:               String("0b989ba4-242f-11e5-81e1-c7b6966d2516"),
	State:              String("pending"),
	LockRepositories:   Bool(true),
	ExcludeAttachments: Bool(false),
	URL:                String("https://api.github.com/orgs/octo-org/migrations/79"),
	CreatedAt:          String("2015-07-06T15:33:38-07:00"),
	UpdatedAt:          String("2015-07-06T15:33:38-07:00"),
	Repositories: []*Repository{
		{
			ID:          Int64(1296269),
			Name:        String("Hello-World"),
			FullName:    String("octocat/Hello-World"),
			Description: String("This your first repo!"),
		},
	},
}

func TestMigration_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Migration{}, "{}")

	u := &Migration{
		ID:                 Int64(1),
		GUID:               String("guid"),
		State:              String("state"),
		LockRepositories:   Bool(false),
		ExcludeAttachments: Bool(false),
		URL:                String("url"),
		CreatedAt:          String("ca"),
		UpdatedAt:          String("ua"),
		Repositories:       []*Repository{{ID: Int64(1)}},
	}

	want := `{
		"id": 1,
		"guid": "guid",
		"state": "state",
		"lock_repositories": false,
		"exclude_attachments": false,
		"url": "url",
		"created_at": "ca",
		"updated_at": "ua",
		"repositories": [
			{
				"id": 1
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestStartMigration_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &startMigration{}, "{}")

	u := &startMigration{
		Repositories:       []string{"r"},
		LockRepositories:   Bool(false),
		ExcludeAttachments: Bool(false),
	}

	want := `{
		"repositories": [
			"r"
		],
		"lock_repositories": false,
		"exclude_attachments": false
	}`

	testJSONMarshal(t, u, want)
}
