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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Import{
		VCS:         String("git"),
		VCSURL:      String("url"),
		VCSUsername: String("u"),
		VCSPassword: String("p"),
	}

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		json.NewDecoder(r.Body).Decode(v)

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
	want := &Import{Status: String("importing")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"status":"complete"}`)
	})

	ctx := context.Background()
	got, _, err := client.Migrations.ImportProgress(ctx, "o", "r")
	if err != nil {
		t.Errorf("ImportProgress returned error: %v", err)
	}
	want := &Import{Status: String("complete")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Import{
		VCS:         String("git"),
		VCSURL:      String("url"),
		VCSUsername: String("u"),
		VCSPassword: String("p"),
	}

	mux.HandleFunc("/repos/o/r/import", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		json.NewDecoder(r.Body).Decode(v)

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
	want := &Import{Status: String("importing")}
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
	client, mux, _, teardown := setup()
	defer teardown()

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
		{ID: Int64(1), Name: String("a")},
		{ID: Int64(2), Name: String("b")},
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &SourceImportAuthor{Name: String("n"), Email: String("e")}

	mux.HandleFunc("/repos/o/r/import/authors/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(SourceImportAuthor)
		json.NewDecoder(r.Body).Decode(v)

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
	want := &SourceImportAuthor{ID: Int64(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Import{UseLFS: String("opt_in")}

	mux.HandleFunc("/repos/o/r/import/lfs", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		json.NewDecoder(r.Body).Decode(v)

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
	want := &Import{Status: String("importing")}
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
	client, mux, _, teardown := setup()
	defer teardown()

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
		{OID: String("a")},
		{OID: String("b")},
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
	client, mux, _, teardown := setup()
	defer teardown()

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
	testJSONMarshal(t, &LargeFile{}, "{}")

	u := &LargeFile{
		RefName: String("rn"),
		Path:    String("p"),
		OID:     String("oid"),
		Size:    Int(1),
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
	testJSONMarshal(t, &SourceImportAuthor{}, "{}")

	u := &SourceImportAuthor{
		ID:         Int64(1),
		RemoteID:   String("rid"),
		RemoteName: String("rn"),
		Email:      String("e"),
		Name:       String("n"),
		URL:        String("url"),
		ImportURL:  String("iurl"),
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
	testJSONMarshal(t, &Import{}, "{}")

	u := &Import{
		VCSURL:          String("vcsurl"),
		VCS:             String("vcs"),
		VCSUsername:     String("vcsusr"),
		VCSPassword:     String("vcspass"),
		TFVCProject:     String("tfvcp"),
		UseLFS:          String("uselfs"),
		HasLargeFiles:   Bool(false),
		LargeFilesSize:  Int(1),
		LargeFilesCount: Int(1),
		Status:          String("status"),
		CommitCount:     Int(1),
		StatusText:      String("statustxt"),
		AuthorsCount:    Int(1),
		Percent:         Int(1),
		PushPercent:     Int(1),
		URL:             String("url"),
		HTMLURL:         String("hurl"),
		AuthorsURL:      String("aurl"),
		RepositoryURL:   String("rurl"),
		Message:         String("msg"),
		FailedStep:      String("fs"),
		HumanName:       String("hn"),
		ProjectChoices:  []*Import{{VCSURL: String("vcsurl")}},
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
