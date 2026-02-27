// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newreposecretwithlibsodium creates a new secret in GitHub for a given owner/repo.
// newreposecretwithlibsodium depends on sodium being installed. Installation instructions for Sodium can be found at this url:
// https://github.com/jedisct1/libsodium
//
// nnewreposecretwithlibsodium has two required flags for owner and repo, and takes in one argument for the name of the secret to add.
// The secret value is pulled from an environment variable based on the secret name.
// To authenticate with GitHub, provide your token via an environment variable GITHUB_AUTH_TOKEN.
//
// To verify the new secret, navigate to GitHub Repository > Settings > left side options bar > Secrets.
//
// Usage:
//
//	export GITHUB_AUTH_TOKEN=<auth token from github that has secret create rights>
//	export SECRET_VARIABLE=<secret value of the secret variable>
//	go run main.go -owner <owner name> -repo <repository name> SECRET_VARIABLE
//
// Example:
//
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
	"github.com/google/go-github/v84/github"
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
		log.Fatal("please provide required flag --owner to specify GitHub user/org owner")
	}

	secretName, err := getSecretName()
	if err != nil {
		log.Fatal(err)
	}

	secretValue, err := getSecretValue(secretName)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	if err := addRepoSecret(ctx, client, *owner, *repo, secretName, secretValue); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added secret %q to the repo %v/%v\n", secretName, *owner, *repo)
}

func getSecretName() (string, error) {
	secretName := flag.Arg(0)
	if secretName == "" {
		return "", fmt.Errorf("missing argument secret name")
	}
	return secretName, nil
}

func getSecretValue(secretName string) (string, error) {
	secretValue := os.Getenv(secretName)
	if secretValue == "" {
		return "", fmt.Errorf("secret value not found under env variable %q", secretName)
	}
	return secretValue, nil
}

// addRepoSecret will add a secret to a GitHub repo for use in GitHub Actions.
//
// Finally, the secretName and secretValue will determine the name of the secret added and it's corresponding value.
//
// The actual transmission of the secret value to GitHub using the api requires that the secret value is encrypted
// using the public key of the target repo. This encryption is done using sodium.
//
// First, the public key of the repo is retrieved. The public key comes base64
// encoded, so it must be decoded prior to use in sodiumlib.
//
// Second, the secret value is converted into a slice of bytes.
//
// Third, the secret is encrypted with sodium.CryptoBoxSeal using the repo's decoded public key.
//
// Fourth, the encrypted secret is encoded as a base64 string to be used in a github.EncodedSecret type.
//
// Fifth, The other two properties of the github.EncodedSecret type are determined. The name of the secret to be added
// (string not base64), and the KeyID of the public key used to encrypt the secret.
// This can be retrieved via the public key's GetKeyID method.
//
// Finally, the github.EncodedSecret is passed into the GitHub client.Actions.CreateOrUpdateRepoSecret method to
// populate the secret in the GitHub repo.
func addRepoSecret(ctx context.Context, client *github.Client, owner, repo, secretName, secretValue string) error {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret); err != nil {
		return fmt.Errorf("client.Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	encryptedBytes, exit := sodium.CryptoBoxSeal([]byte(secretValue), decodedPublicKey)
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
