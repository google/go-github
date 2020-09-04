// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newreposecret creates a new secret in GitHub for a given owner/repo.
// It has two required flags for owner and repo, and takes in one argument for the name of the secret to add.
// Provide the value of the secret you want to add with an environment variable of the same name.
// To authenticate with GitHub provide it via an environment variable GITHUB_AUTH_TOKEN.
//
// To verify the new secret, navigate to GitHub Repository > Settings > left side options bar > Secrets.
//
// Usage:
//	export GITHUB_AUTH_TOKEN=<auth token from github that has secret create rights>
//	export SECRET_VARIABLE=<secret value of the secret variable>
//	go run main.go -owner <owner name> -repo <repository name> SECRET_VARIABLE
//
// Example:
//	export GITHUB_AUTH_TOKEN=0000000000000000
//	export SECRET_VARIABLE="my-secret"
//	go run main.go -owner google -repo go-github SECRET_VARIABLE
package main

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var (
	repo  = flag.String("repo", "", "The repo that the secret should be added to, ex. go-github")
	owner = flag.String("owner", "", "The owner of there repo this should be added to, ex. google")
)

func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}

	if *repo == "" {
		log.Fatal("please provide required flag --repo to specify GitHub repository ")
	}

	if *owner == "" {
		log.Fatal("please provide required flag --owner to speify GitHub user/org owner")
	}

	secretName, err := getSecretName()
	if err != nil {
		log.Fatal(err)
	}

	secretValue, err := getSecretValue(secretName)
	if err != nil {
		log.Fatal(err)
	}

	ctx, client, err := githubAuth(token)
	if err != nil {
		log.Fatalf("unable to authorize using env GITHUB_AUTH_TOKEN: %v", err)
	}

	err = addRepoSecret(ctx, client, *owner, *repo, secretName, secretValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added secret %q to the repo %v/%v\n", secretName, *owner, *repo)
}

func getSecretName() (string, error) {
	secretName := flag.Arg(0)
	if secretName == "" {
		err := fmt.Errorf("missing argument secret name")
		return "", err
	}
	return secretName, nil
}

func getSecretValue(secretName string) (string, error) {
	secretValue := os.Getenv(secretName)
	if secretValue == "" {
		err := fmt.Errorf("secret value not found under env variable %q", secretName)
		return "", err
	}
	return secretValue, nil
}

// githubAuth returns a GitHub client and context
// It reads an environment variable GITHUB_AUTH_TOKEN
// for a GitHub API token with secret read/write permissions
// cannot be the default token from GitHub actions
func githubAuth(token string) (context.Context, *github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

// addRepoSecret will add a secret to a GitHub repo for use in GitHub Actions
//
// To communicate with GitHub addRepoSecret requires a client and context.
// To determine what repository to add this secret to it will use the owner and repo combination
// Finally the secretName and secretValue will determine the name of the secret added and it's corresponding value
//
// The actual transmission of the secret value to GitHub using the api requires that the secret value is encrypted
// using the public key of the target repo. This encryption must be done using sodium. This function has a hard
// dependency on sodium being installed where this is run. Sodium can be installed from this url:
// https://formulae.brew.sh/formula/libsodium
//
// First step of the upload process by this function is to get the public key of the repo. The public key comes base64
// encoded, so it must be decoded prior to use in sodiumlib.
//
// Second the secret value starts as a string type, but must be converted into a slice of bytes.
// Third the decoded public key of the repo is used to encrypt the secret with sodium.CryptoBoxSeal resulting
// in a secret encrypted as bytes
//
// Finally the secret bytes need to be encoded as a base64 string and used in a github.EncodedSecret type
// that is passed into the GitHub client.Actions.CreateOrUpdateRepoSecret method to populate the secret in GitHub.
//
// The other two properties of the github.EncodedSecret type are the name of the secret to be added (string not base64)
// and the KeyID of the public key used to encrypt the secret, which is gettable from the public key's GetKeyID method.
//
// Finally it passes in the github.EncodedSecret object to CreateOrUpdateRepoSecret which creates or updates the secret
// in GitHub
func addRepoSecret(ctx context.Context, client *github.Client, owner string, repo, secretName string, secretValue string) error {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret)
	if err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err))
	}

	secretBytes := []byte(secretValue)
	encryptedBytes, exit := sodium.CryptoBoxSeal(secretBytes, decodedPublicKey)
	if exit != 0 {
		return nil, errors.New("sodium.CryptoBoxSeal exited with non zero exit code")
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}
