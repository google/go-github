// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tokenauth command demonstrates using a Personal Access Token (PAT) to
// authenticate with GitHub.
// You can test out a GitHub Personal Access Token using this simple example.
// You can generate them here: https://github.com/settings/tokens
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v84/github"
	"golang.org/x/term"
)

func main() {
	fmt.Print("GitHub Token: ")
	token, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(string(token))

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	fmt.Printf("\n%v\n", github.Stringify(user))
}
