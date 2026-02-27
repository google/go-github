// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The uploadreleaseassetfromrelease example demonstrates how to upload
// a release asset using the UploadReleaseAssetFromRelease helper.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v84/github"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_AUTH_TOKEN not set")
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	owner := "OWNER"
	repo := "REPO"
	releaseID := int64(1)

	// Fetch the release (UploadURL is populated by the API)
	release, _, err := client.Repositories.GetRelease(ctx, owner, repo, releaseID)
	if err != nil {
		log.Fatalf("GetRelease failed: %v", err)
	}

	// Asset content
	data := []byte("Hello from go-github!\n")
	reader := bytes.NewReader(data)
	size := int64(len(data))

	opts := &github.UploadOptions{
		Name:  "example.txt",
		Label: "Example asset",
	}

	asset, _, err := client.Repositories.UploadReleaseAssetFromRelease(
		ctx,
		release,
		opts,
		reader,
		size,
	)
	if err != nil {
		log.Fatalf("UploadReleaseAssetFromRelease failed: %v", err)
	}

	fmt.Printf("Uploaded asset ID: %v\n", asset.GetID())
}
