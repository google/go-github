// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This embed-interface example is a copy of the "simple" example
// and its purpose is to demonstrate how embedding an interface
// in a struct makes it easy to mock one or more methods.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v34/github"
)

type myClient struct {
	client *github.Client
}

// Fetch all the public organizations' membership of a user.
//
func (m *myClient) fetchOrganizations(ctx context.Context, username string) ([]*github.Organization, error) {
	orgs, _, err := m.client.Organizations.List(ctx, username, nil)
	return orgs, err
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	c := &myClient{
		client: github.NewClient(nil),
	}

	ctx := context.Background()
	organizations, err := c.fetchOrganizations(ctx, username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, organization := range organizations {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}
}
