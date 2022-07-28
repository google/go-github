package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
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

	fmt.Print("Tag pattern(leave blank if don't want to create new tag protection): ")
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
		tagProtection, _, err := client.Repositories.CreateTagProtection(context.Background(), owner, repo, pattern)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		println()
		fmt.Printf("New tag protection created in github.com/%v/%v\n", owner, repo)
		tp, _ := json.Marshal(tagProtection)
		fmt.Println(string(tp))
	}

	// list all tag protection
	println()
	fmt.Printf("List all tag protection in github.com/%v/%v\n", owner, repo)
	tagProtections, _, err := client.Repositories.ListTagProtection(context.Background(), owner, repo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	results, _ := json.Marshal(tagProtections)
	fmt.Println(string(results))
}
