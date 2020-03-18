// Copyright 2019 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The simple command demonstrates a simple functionality which
// prompts the user for a GitHub topic and lists all the public
// organization memberships of the specified topic.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v30/github"
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

	for _, topic := range topics.Topics {
		fmt.Println(*topic.Name)
	}
}
