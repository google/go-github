// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/google/go-github/github"
)

func main() {
	client := github.NewClient(nil)

	fmt.Println("Retrieving the go-github README.md file.")
	encodedText, _, err := client.Repositories.GetReadme("google", "go-github", &github.RepositoryContentGetOptions{})

	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
	}
	if encodedText == nil {
		fmt.Println("The returned text is nil. Are you sure it exists?")
	}
	text, err := encodedText.Decode()
	if err != nil {
		fmt.Printf("Decoding failed: %v", err)
	}
	readme := string(text)
	fmt.Printf("Converted readme:\n%v\n", readme)
}
