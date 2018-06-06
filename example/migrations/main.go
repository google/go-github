// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// migrations command demostrate the functionality of
// user data migration for authenticated Github user
// and list all the user migrations.
package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func FetchAllUserMigrations() ([]*github.UserMigration, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "b453839afb5b7871b9dabfe8184b314747f6523a"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user_migrations, _, err := client.Migrations.ListUserMigrations(ctx)

	return user_migrations, err
}

func main() {
	user_migrations, err := FetchAllUserMigrations()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	for i, user_migration := range user_migrations {
		fmt.Printf("%v. %v", i+1, user_migration.GetID())
	}
}
