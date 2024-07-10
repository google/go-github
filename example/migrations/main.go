// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// migrations demonstrates the functionality of
// the user data migration API for the authenticated GitHub
// user and lists all the user migrations.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v63/github"
)

func fetchAllUserMigrations() ([]*github.UserMigration, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken("<GITHUB_AUTH_TOKEN>")

	migrations, _, err := client.Migrations.ListUserMigrations(ctx, &github.ListOptions{Page: 1})
	return migrations, err
}

func main() {
	migrations, err := fetchAllUserMigrations()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	for i, m := range migrations {
		fmt.Printf("%v. %v", i+1, m.GetID())
	}
}
