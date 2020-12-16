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
	"reflect"
	"testing"
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
		if !reflect.DeepEqual(v, input) {
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("StartImport = %+v, want %+v", got, want)
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ImportProgress = %+v, want %+v", got, want)
	}
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
		if !reflect.DeepEqual(v, input) {
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateImport = %+v, want %+v", got, want)
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CommitAuthors = %+v, want %+v", got, want)
	}
}

func TestMigrationService_MapCommitAuthor(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &SourceImportAuthor{Name: String("n"), Email: String("e")}

	mux.HandleFunc("/repos/o/r/import/authors/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(SourceImportAuthor)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapCommitAuthor = %+v, want %+v", got, want)
	}
}

func TestMigrationService_SetLFSPreference(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Import{UseLFS: String("opt_in")}

	mux.HandleFunc("/repos/o/r/import/lfs", func(w http.ResponseWriter, r *http.Request) {
		v := new(Import)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetLFSPreference = %+v, want %+v", got, want)
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LargeFiles = %+v, want %+v", got, want)
	}
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
}
