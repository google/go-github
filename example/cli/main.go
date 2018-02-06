package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		fmt.Println("Unauthorized: No token present")
	}
}
