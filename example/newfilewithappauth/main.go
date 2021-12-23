/*
Copyright 2018 The go-github AUTHORS. All rights reserved.

Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

newfilewithappauth demonstrates the functionality of github's app authentication
methods by fetching an installation access token and reauthenticating to github
with OAuth configurations
 */
package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v41/github"
	"github.com/bradleyfalzon/ghinstallation"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	const gitHost = "https://git.api.com"

	//get private key data
	//download pem from github app
	privatePem, err := ioutil.ReadFile("path/to/pem")
	if err != nil {
		panic(err)
	}

	// Wrap the shared transport for use with the
	// app ID and app API
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, 10, privatePem)
	if err != nil {
		fmt.Printf("faild to create app transport: %v\n", err)
	}
	//your enterprise git URI
	itr.BaseURL = gitHost

	//create git client with app transport
	client, err := github.NewEnterpriseClient(
		gitHost,
		gitHost,
		&http.Client{
		Transport: itr,
		Timeout:   time.Second * 30,
	})
	if err != nil {
		fmt.Printf("faild to create git client for app: %v\n", err)
	}

	//list installations
	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		fmt.Printf("failed to list installations: %v\n", err)
	}

	//capture our installationId for our app
	//we need this for the access token
	var installID int64
	for _, install := range installations {
		installID = install.GetID()
	}

	//get installation token
	token, _, err := client.Apps.CreateInstallationToken(
		context.Background(),
		installID,
		&github.InstallationTokenOptions{})
	if err != nil {
		fmt.Printf("failed to create installation token: %v\n", err)
	}

	//create OAuth client with access token to interact with files
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.GetToken()},
	)
	oAuthClient := oauth2.NewClient(context.Background(), ts)

	//create new git hub client with accessToken
	apiClient, err := github.NewEnterpriseClient(gitHost, gitHost, oAuthClient)
	if err != nil {
		fmt.Printf("failed to create new git client with token: %v\n", err)
	}

	// create new file
	_, resp, err := apiClient.Repositories.CreateFile(
		context.Background(),
		"repoOwner",
		"sample-repo",
		"example/foo.txt",
		&github.RepositoryContentFileOptions{
			Content: []byte("foo"),
			Message: github.String("sample commit"),
			SHA:     nil,
		})
	if err != nil {
		fmt.Printf("failed to create new file: %v\n", err)
	}

	fmt.Printf("file written: %v", resp.StatusCode)

}
