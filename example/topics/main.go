// Copyright 2019 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The simple command demonstrates the functionality that
// prompts the user for a GitHub topic and lists all the entities
// that are related to the specified topic or subject.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v84/github"
)

// Fetch and lists all the public topics associated with the specified GitHub topic.
func fetchTopics(topic string) (*github.TopicsSearchResult, error) {
	client := github.NewClient(nil)
	topics, _, err := client.Search.Topics(context.Background(), topic, nil)
	return topics, err
}

func main() {
	var topic string
	fmt.Print("Enter GitHub topic: ")
	fmt.Scanf("%s", &topic)

	topics, err := fetchTopics(topic)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, topic := range topics.Topics {
		fmt.Println(*topic.Name)
	}
}
