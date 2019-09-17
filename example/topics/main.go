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

	"github.com/google/go-github/v28/github"
)

// Fetch all the public organizations' membership of a user.
//
func FetchTopics(topic string) (*github.TopicsSearchResult, error) {
	client := github.NewClient(nil)
	topics, _, err := client.Search.Topics(context.Background(), topic, nil)
	return topics, err
}

func main() {
	var topic string
	fmt.Print("Enter GitHub topic: ")
	fmt.Scanf("%s", &topic)

	topics, err := FetchTopics(topic)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(topics)
}
