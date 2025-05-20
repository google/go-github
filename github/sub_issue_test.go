// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestSubIssuesService_Add(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &SubIssueRequest{SubIssueID: 42}

	mux.HandleFunc("/repos/o/r/issues/1/sub_issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(SubIssueRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":42, "number":1}`)
	})

	ctx := context.Background()
	got, _, err := client.SubIssue.Add(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("SubIssues.Add returned error: %v", err)
	}

	want := &SubIssue{Number: Ptr(1), ID: Int64(42)}
	if !cmp.Equal(got, want) {
		t.Errorf("SubIssues.Add = %+v, want %+v", got, want)
	}
}

func TestSubIssuesService_ListByIssue(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/1/sub_issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3)

		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	opt := &IssueListOptions{}
	issues, _, err := client.SubIssue.ListByIssue(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("SubIssues.ListByIssue returned error: %v", err)
	}

	want := []*SubIssue{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(issues, want) {
		t.Errorf("SubIssues.ListByIssue = %+v, want %+v", issues, want)
	}
}

func TestSubIssuesService_Remove(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &SubIssueRequest{SubIssueID: 42}

	mux.HandleFunc("/repos/o/r/issues/1/sub_issues", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Accept", mediaTypeV3)

		v := new(SubIssueRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "DELETE")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":42, "number":1}`)
	})

	ctx := context.Background()
	got, _, err := client.SubIssue.Remove(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("SubIssues.Remove returned error: %v", err)
	}

	want := &SubIssue{ID: Int64(42), Number: Ptr(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("SubIssues.Remove = %+v, want %+v", got, want)
	}
}

func TestSubIssuesService_Reprioritize(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &SubIssueRequest{SubIssueID: 42, AfterID: 5}

	mux.HandleFunc("/repos/o/r/issues/1/sub_issues/priority", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Accept", mediaTypeV3)

		v := new(SubIssueRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":42, "number":1}`)
	})

	ctx := context.Background()
	got, _, err := client.SubIssue.Reprioritize(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("SubIssues.Reprioritize returned error: %v", err)
	}

	want := &SubIssue{ID: Int64(42), Number: Ptr(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("SubIssues.Reprioritize = %+v, want %+v", got, want)
	}
}
