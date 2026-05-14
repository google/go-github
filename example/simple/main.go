// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The simple command demonstrates a simple functionality which
// prompts the user for a GitHub username and lists all the public
// organization memberships of the specified username.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v86/github"
)

// Fetch all the public organizations' membership of a user.
func fetchOrganizations(username string) ([]*github.Organization, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, err
	}
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	organizations, err := fetchOrganizations(username)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i, organization := range organizations {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}
}
