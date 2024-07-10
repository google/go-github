// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package demo provides an app that shows how to use the github package on
// Google App Engine.
package demo

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v63/github"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := appengine.NewContext(r)
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_AUTH_TOKEN"))

	commits, _, err := client.Repositories.ListCommits(ctx, "google", "go-github", nil)
	if err != nil {
		log.Errorf(ctx, "ListCommits: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, commit := range commits {
		fmt.Fprintln(w, commit.GetHTMLURL())
	}
}
