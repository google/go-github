package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		fmt.Println("Unauthorized: No token present")
		return
	}
	name := os.Args[1:]
	if len(name) == 0 {
		fmt.Println("No name: New repos must be given a name")
		return
	}
	fmt.Println(name)
}
