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

func TestIterators(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
		// Just consume body
		r.Body.Close()

		page := r.URL.Query().Get("page")
		switch page {
		case "", "1":
			w.Header().Set("Link", `<http://localhost/users/u/repos?page=2>; rel="next"`)
			fmt.Fprint(w, `[{"id":1}]`)
		case "2":
			fmt.Fprint(w, `[{"id":2}]`)
		}
	})

	ctx := t.Context()
	opts := &RepositoryListByUserOptions{
		ListOptions: ListOptions{PerPage: 1},
	}

	var repos []*Repository
	for repo, err := range client.Repositories.ListByUserIter(ctx, "u", opts) {
		if err != nil {
			t.Fatalf("ListByUserIter returned error: %v", err)
		}
		repos = append(repos, repo)
	}

	if len(repos) != 2 {
		t.Errorf("ListByUserIter returned %v repos, want 2", len(repos))
	}
	if repos[0].GetID() != 1 {
		t.Errorf("repo[0].ID = %v, want 1", repos[0].GetID())
	}
	if repos[1].GetID() != 2 {
		t.Errorf("repo[1].ID = %v, want 2", repos[1].GetID())
	}
}
