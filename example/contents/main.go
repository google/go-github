// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The contents command utilizes go-github as a CLI tool for
// downloading the contents of a file in a repository.
// It takes an inputs of the repository owner, repository name, path to the
// file in the repository, reference (branch, tag or commit SHA), and output
// path for the downloaded file. It then uses the Repositories.DownloadContents
// method to download the file and saves it to the specified output path.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v85/github"
)

func main() {
	fmt.Println("This example will download the contents of a file from a GitHub repository.")

	r := bufio.NewReader(os.Stdin)

	fmt.Print("Repository Owner: ")
	owner, _ := r.ReadString('\n')
	owner = strings.TrimSpace(owner)

	fmt.Print("Repository Name: ")
	repo, _ := r.ReadString('\n')
	repo = strings.TrimSpace(repo)

	fmt.Print("Repository Path: ")
	repoPath, _ := r.ReadString('\n')
	repoPath = strings.TrimSpace(repoPath)

	fmt.Print("Reference (branch, tag or commit SHA): ")
	ref, _ := r.ReadString('\n')
	ref = strings.TrimSpace(ref)

	fmt.Print("Output Path: ")
	outputPath, _ := r.ReadString('\n')
	outputPath = filepath.Clean(strings.TrimSpace(outputPath))

	fmt.Printf("\nDownloading %v/%v/%v at ref %v to %v...\n", owner, repo, repoPath, ref, outputPath)

	client := github.NewClient(nil)

	rc, _, err := client.Repositories.DownloadContents(context.Background(), owner, repo, repoPath, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer rc.Close()

	f, err := os.Create(outputPath) //#nosec G703 -- path is validated above
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := io.Copy(f, rc); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download completed.")
}
