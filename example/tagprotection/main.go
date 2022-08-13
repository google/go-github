// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tagprotection command demonstrates the functionality that
// prompts the user for GitHub owner, repo, tag protection pattern and token,
// then creates a new tag protection if the user entered a pattern at the prompt.
// Otherwise, it will just list all existing tag protections.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/google/go-github/v45/github"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

func main() {
	// read github owner, repo, token from standard input
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Org/User name: ")
	owner, _ := r.ReadString('\n')
	owner = strings.TrimSpace(owner)

	fmt.Print("GitHub repo name: ")
	repo, _ := r.ReadString('\n')
	repo = strings.TrimSpace(repo)

	fmt.Print("Tag pattern(leave blank to not create new tag protection): ")
	pattern, _ := r.ReadString('\n')
	pattern = strings.TrimSpace(pattern)

	fmt.Print("GitHub Token: ")
	byteToken, _ := terminal.ReadPassword(int(syscall.Stdin))
	println()
	token := string(byteToken)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// create new tag protection
	if pattern != "" {
		tagProtection, _, err := client.Repositories.CreateTagProtection(ctx, owner, repo, pattern)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		println()
		fmt.Printf("New tag protection created in github.com/%v/%v\n", owner, repo)
		tp, _ := json.Marshal(tagProtection)
		fmt.Println(string(tp))
	}

	// list all tag protection
	println()
	fmt.Printf("List all tag protection in github.com/%v/%v\n", owner, repo)
	tagProtections, _, err := client.Repositories.ListTagProtection(ctx, owner, repo)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	results, _ := json.Marshal(tagProtections)
	fmt.Println(string(results))
}
