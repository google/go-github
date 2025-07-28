// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The social_accounts command demonstrates how to use the social accounts API.
// It lists social accounts for a given username.
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/google/go-github/v74/github"
)

// Fetch social accounts for a user.
func fetchSocialAccounts(username string) ([]*github.SocialAccount, error) {
	client := github.NewClient(nil)
	accounts, _, err := client.Users.ListUserSocialAccounts(context.Background(), username, nil)
	return accounts, err
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	accounts, err := fetchSocialAccounts(username)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Printf("No social accounts found for %v\n", username)
		return
	}

	fmt.Printf("Social accounts for %s:\n", username)
	for i, account := range accounts {
		fmt.Printf("%v. Provider: %v, URL: %v\n",
			i+1,
			*account.Provider,
			*account.URL)
	}
}
