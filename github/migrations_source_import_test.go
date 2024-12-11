// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMigrationService_StartImport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Import{
		VCS:         Ptr("git"),
		VCSURL:      Ptr("url"),
		VCSUsername: Ptr("u"),
		VCSPassword: Ptr("p"),
	}

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"status":"importing"}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.StartImport(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("StartImport returned error: %v", err)
	}
	want := &Import{Status: Ptr("importing")}
	if !cmp.Equal(got, want) {
		t.Errorf("StartImport = %+v, want %+v", got, want)
	}

	const methodName = "StartImport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.StartImport(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.StartImport(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_ImportProgress(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"status":"complete"}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.ImportProgress(ctx, "o", "r")
	if err != nil {
		t.Errorf("ImportProgress returned error: %v", err)
	}
	want := &Import{Status: Ptr("complete")}
	if !cmp.Equal(got, want) {
		t.Errorf("ImportProgress = %+v, want %+v", got, want)
	}

	const methodName = "ImportProgress"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.ImportProgress(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.ImportProgress(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_UpdateImport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Import{
		VCS:         Ptr("git"),
		VCSURL:      Ptr("url"),
		VCSUsername: Ptr("u"),
		VCSPassword: Ptr("p"),
	}

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"status":"importing"}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.UpdateImport(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("UpdateImport returned error: %v", err)
	}
	want := &Import{Status: Ptr("importing")}
	if !cmp.Equal(got, want) {
		t.Errorf("UpdateImport = %+v, want %+v", got, want)
	}

	const methodName = "UpdateImport"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.UpdateImport(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.UpdateImport(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_CommitAuthors(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import/authors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1,"name":"a"},{"id":2,"name":"b"}]`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.CommitAuthors(ctx, "o", "r")
	if err != nil {
		t.Errorf("CommitAuthors returned error: %v", err)
	}
	want := []*SourceImportAuthor{
		{ID: Ptr(int64(1)), Name: Ptr("a")},
		{ID: Ptr(int64(2)), Name: Ptr("b")},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("CommitAuthors = %+v, want %+v", got, want)
	}

	const methodName = "CommitAuthors"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.CommitAuthors(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.CommitAuthors(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_MapCommitAuthor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &SourceImportAuthor{Name: Ptr("n"), Email: Ptr("e")}

	mux.HandleFunc("/repos/o/r/import/authors/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(SourceImportAuthor)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id": 1}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.MapCommitAuthor(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("MapCommitAuthor returned error: %v", err)
	}
	want := &SourceImportAuthor{ID: Ptr(int64(1))}
	if !cmp.Equal(got, want) {
		t.Errorf("MapCommitAuthor = %+v, want %+v", got, want)
	}

	const methodName = "MapCommitAuthor"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.MapCommitAuthor(ctx, "\n", "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.MapCommitAuthor(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_SetLFSPreference(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &Import{UseLFS: Ptr("opt_in")}

	mux.HandleFunc("/repos/o/r/import/lfs", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"status":"importing"}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.SetLFSPreference(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("SetLFSPreference returned error: %v", err)
	}
	want := &Import{Status: Ptr("importing")}
	if !cmp.Equal(got, want) {
		t.Errorf("SetLFSPreference = %+v, want %+v", got, want)
	}

	const methodName = "SetLFSPreference"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.SetLFSPreference(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.SetLFSPreference(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_LargeFiles(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import/large_files", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"oid":"a"},{"oid":"b"}]`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.LargeFiles(ctx, "o", "r")
	if err != nil {
		t.Errorf("LargeFiles returned error: %v", err)
	}
	want := []*LargeFile{
		{OID: Ptr("a")},
		{OID: Ptr("b")},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("LargeFiles = %+v, want %+v", got, want)
	}

	const methodName = "LargeFiles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Migrations.LargeFiles(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Migrations.LargeFiles(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMigrationService_CancelImport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Migrations.CancelImport(ctx, "o", "r")
	if err != nil {
		t.Errorf("CancelImport returned error: %v", err)
	}

	const methodName = "CancelImport"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Migrations.CancelImport(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Migrations.CancelImport(ctx, "o", "r")
	})
}

func TestLargeFile_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &LargeFile{}, "{}")

	u := &LargeFile{
		RefName: Ptr("rn"),
		Path:    Ptr("p"),
		OID:     Ptr("oid"),
		Size:    Ptr(1),
	}

	want := `{
		"ref_name": "rn",
		"path": "p",
		"oid": "oid",
		"size": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestSourceImportAuthor_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SourceImportAuthor{}, "{}")

	u := &SourceImportAuthor{
		ID:         Ptr(int64(1)),
		RemoteID:   Ptr("rid"),
		RemoteName: Ptr("rn"),
		Email:      Ptr("e"),
		Name:       Ptr("n"),
		URL:        Ptr("url"),
		ImportURL:  Ptr("iurl"),
	}

	want := `{
		"id": 1,
		"remote_id": "rid",
		"remote_name": "rn",
		"email": "e",
		"name": "n",
		"url": "url",
		"import_url": "iurl"
	}`

	testJSONMarshal(t, u, want)
}

func TestImport_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Import{}, "{}")

	u := &Import{
		VCSURL:          Ptr("vcsurl"),
		VCS:             Ptr("vcs"),
		VCSUsername:     Ptr("vcsusr"),
		VCSPassword:     Ptr("vcspass"),
		TFVCProject:     Ptr("tfvcp"),
		UseLFS:          Ptr("uselfs"),
		HasLargeFiles:   Ptr(false),
		LargeFilesSize:  Ptr(1),
		LargeFilesCount: Ptr(1),
		Status:          Ptr("status"),
		CommitCount:     Ptr(1),
		StatusText:      Ptr("statustxt"),
		AuthorsCount:    Ptr(1),
		Percent:         Ptr(1),
		PushPercent:     Ptr(1),
		URL:             Ptr("url"),
		HTMLURL:         Ptr("hurl"),
		AuthorsURL:      Ptr("aurl"),
		RepositoryURL:   Ptr("rurl"),
		Message:         Ptr("msg"),
		FailedStep:      Ptr("fs"),
		HumanName:       Ptr("hn"),
		ProjectChoices:  []*Import{{VCSURL: Ptr("vcsurl")}},
	}

	want := `{
		"vcs_url": "vcsurl",
		"vcs": "vcs",
		"vcs_username": "vcsusr",
		"vcs_password": "vcspass",
		"tfvc_project": "tfvcp",
		"use_lfs": "uselfs",
		"has_large_files": false,
		"large_files_size": 1,
		"large_files_count": 1,
		"status": "status",
		"commit_count": 1,
		"status_text": "statustxt",
		"authors_count": 1,
		"percent": 1,
		"push_percent": 1,
		"url": "url",
		"html_url": "hurl",
		"authors_url": "aurl",
		"repository_url": "rurl",
		"message": "msg",
		"failed_step": "fs",
		"human_name": "hn",
		"project_choices": [
			{
				"vcs_url": "vcsurl"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
