// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIterators_Table(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mock      func(w http.ResponseWriter, r *http.Request)
		opts      *RepositoryListByUserOptions
		wantCount int
		wantErr   bool
	}{
		{
			name: "single page",
			mock: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, `[{"id":1}]`)
			},
			wantCount: 1,
		},
		{
			name: "multi page",
			mock: func(w http.ResponseWriter, r *http.Request) {
				page := r.URL.Query().Get("page")
				if page == "" || page == "1" {
					w.Header().Set("Link", `<http://localhost/api-v3/users/u/repos?page=2>; rel="next"`)
					fmt.Fprint(w, `[{"id":1}]`)
				} else {
					fmt.Fprint(w, `[{"id":2}]`)
				}
			},
			wantCount: 2,
		},
		{
			name: "error on first page",
			mock: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
				tt.mock(w, r)
			})

			ctx := t.Context()
			count := 0
			var iterErr error
			for _, err := range client.Repositories.ListByUserIter(ctx, "u", tt.opts) {
				if err != nil {
					iterErr = err
					break // Stop iteration on error
				}
				count++
			}

			if tt.wantErr {
				if iterErr == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if iterErr != nil {
					t.Errorf("unexpected error: %v", iterErr)
				}
				if count != tt.wantCount {
					t.Errorf("got %v items, want %v", count, tt.wantCount)
				}
			}
		})
	}
}

func TestIterators_Safety(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "1" || page == "" {
			w.Header().Set("Link", `<http://localhost/api-v3/users/u/repos?page=2>; rel="next"`)
		}
		fmt.Fprint(w, `[]`)
	})

	opts := &RepositoryListByUserOptions{
		ListOptions: ListOptions{Page: 1},
	}
	originalPage := opts.Page

	ctx := t.Context()
	// Run iterator. Even if no items, it should finish when pagination ends.
	for _, err := range client.Repositories.ListByUserIter(ctx, "u", opts) {
		if err != nil {
			break
		}
	}

	if opts.Page != originalPage {
		t.Errorf("original opts.Page mutated! got %v, want %v", opts.Page, originalPage)
	}
}
