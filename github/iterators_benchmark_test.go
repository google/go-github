// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func setupBenchmark(b *testing.B) (*Client, *http.ServeMux) {
	b.Helper()
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	server := httptest.NewServer(apiHandler)

	client := NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	b.Cleanup(server.Close)
	return client, mux
}

func BenchmarkIterators(b *testing.B) {
	client, mux := setupBenchmark(b)

	mux.HandleFunc("/users/u/repos", func(w http.ResponseWriter, r *http.Request) {
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

	ctx := b.Context()
	opts := &RepositoryListByUserOptions{
		ListOptions: ListOptions{PerPage: 1},
	}

	b.Run("Manual", func(b *testing.B) {
		for b.Loop() {
			count := 0
			curOpts := *opts
			for {
				repos, resp, err := client.Repositories.ListByUser(ctx, "u", &curOpts)
				if err != nil {
					b.Fatalf("ListByUser returned error: %v", err)
				}
				for range repos {
					count++
				}
				if resp.NextPage == 0 {
					break
				}
				curOpts.Page = resp.NextPage
			}
			if count != 2 {
				b.Fatalf("Expected 2 items, got %d", count)
			}
		}
	})

	b.Run("Iterator", func(b *testing.B) {
		for b.Loop() {
			count := 0
			for _, err := range client.Repositories.ListByUserIter(ctx, "u", opts) {
				if err != nil {
					b.Fatalf("ListByUserIter returned error: %v", err)
				}
				count++
			}
			if count != 2 {
				b.Fatalf("Expected 2 items, got %d", count)
			}
		}
	})
}
