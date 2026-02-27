// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The basicauth command demonstrates using the github.BasicAuthTransport,
// including handling two-factor authentication. This won't currently work for
// accounts that use SMS to receive one-time passwords.
//
// Deprecation Notice: GitHub will discontinue password authentication to the API.
// You must now authenticate to the GitHub API with an API token, such as an OAuth access token,
// GitHub App installation access token, or personal access token, depending on what you need to do with the token.
// Password authentication to the API will be removed on November 13, 2020.
// See the tokenauth example for details.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v84/github"
	"golang.org/x/term"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("GitHub Password: ")
	password, _ := term.ReadPassword(int(os.Stdin.Fd()))

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(string(password)),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")

	// Is this a two-factor auth error? If so, prompt for OTP and try again.
	if errors.As(err, new(*github.TwoFactorAuthError)) {
		fmt.Print("\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		user, _, err = client.Users.Get(ctx, "")
	}

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	fmt.Printf("\n%v\n", github.Stringify(user))
}
