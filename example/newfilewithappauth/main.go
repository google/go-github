// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newfilewithappauth demonstrates the functionality of GitHub's app authentication
// methods by fetching an installation access token as a GitHub App and then
// using that token to authenticate with the GitHub API.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v93/github"
	"github.com/jferrl/go-githubauth"
)

func main() {
	gitHost := "https://git.api.com"

	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub App Client ID: ")
	clientID, _ := r.ReadString('\n')

	fmt.Print("Path to GitHub App Private Key (.pem): ")
	pemPath, _ := r.ReadString('\n')
	privateKey, err := os.ReadFile(strings.TrimSpace(pemPath))
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(strings.TrimSpace(clientID), privateKey)
	if err != nil {
		log.Fatalf("failed to create app token source: %v", err)
	}

	// create git client with app transport
	client, err := github.NewClient(github.WithTokenSource(appTokenSource), github.WithEnterpriseURLs(gitHost, gitHost), github.WithTimeout(30*time.Second))
	if err != nil {
		log.Fatalf("failed to create git client for app: %v", err)
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		log.Fatalf("failed to list installations: %v", err)
	}

	// capture our installationId for our app
	// we need this for the access token
	var installID int64
	for _, val := range installations {
		installID = val.GetID()
		break
	}

	if installID == 0 {
		log.Fatal("no installation found")
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installID, appTokenSource, githubauth.WithEnterpriseURL(gitHost))

	apiClient, err := github.NewClient(github.WithTokenSource(installationTokenSource), github.WithEnterpriseURLs(gitHost, gitHost), github.WithTimeout(30*time.Second))
	if err != nil {
		log.Fatalf("failed to create new git client with token: %v", err)
	}

	_, resp, err := apiClient.Repositories.CreateFile(
		context.Background(),
		"repoOwner",
		"sample-repo",
		"example/foo.txt",
		&github.RepositoryContentFileOptions{
			Content: []byte("foo"),
			Message: new("sample commit"),
			SHA:     nil,
		},
	)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}

	log.Printf("file written status code: %v", resp.StatusCode)
}
