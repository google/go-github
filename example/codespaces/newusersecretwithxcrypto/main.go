// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newusersecretwithxcrypto creates a new secret in GitHub for a given user.
// newusersecretwithxcrypto uses x/crypto/nacl/box instead of sodium.
// It does not depend on any native libraries and is easier to cross-compile for different platforms.
// Quite possibly there is a performance penalty due to this.
//
// newusersecretwithxcrypto takes in one argument for the name of the secret to add, and 2 flags owner, repo.
// If owner/repo are defined then it adds the secret to that repository
// The secret value is pulled from an environment variable based on the secret name.
// To authenticate with GitHub, provide your token via an environment variable GITHUB_AUTH_TOKEN.
//
// To verify the new secret, navigate to GitHub User > Settings > left side options bar > Codespaces > Secrets.
//
// Usage:
//
//	export GITHUB_AUTH_TOKEN=<auth token from github that has secret create rights>
//	export SECRET_VARIABLE=<secret value of the secret variable>
//	go run main.go SECRET_VARIABLE
//
// Example:
//
//	export GITHUB_AUTH_TOKEN=0000000000000000
//	export SECRET_VARIABLE="my-secret"
//	go run main.go SECRET_VARIABLE
package main

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v84/github"
	"golang.org/x/crypto/nacl/box"
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

	if err := addUserSecret(ctx, client, secretName, secretValue, *owner, *repo); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added secret %q to the authenticated user\n", secretName)
}

func getSecretName() (string, error) {
	secretName := flag.Arg(0)
	if secretName == "" {
		return "", errors.New("missing argument secret name")
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

// addUserSecret will add a secret to a GitHub user for use in GitHub Codespaces.
//
// The secretName and secretValue will determine the name of the secret added and it's corresponding value.
//
// The actual transmission of the secret value to GitHub using the api requires that the secret value is encrypted
// using the public key of the target user. This encryption is done using x/crypto/nacl/box.
//
// First, the public key of the user is retrieved. The public key comes base64
// encoded, so it must be decoded prior to use.
//
// Second, the decode key is converted into a fixed size byte array.
//
// Third, the secret value is converted into a slice of bytes.
//
// Fourth, the secret is encrypted with box.SealAnonymous using the user's decoded public key.
//
// Fifth, the encrypted secret is encoded as a base64 string to be used in a github.EncodedSecret type.
//
// Sixth, The other two properties of the github.EncodedSecret type are determined. The name of the secret to be added
// (string not base64), and the KeyID of the public key used to encrypt the secret.
// This can be retrieved via the public key's GetKeyID method.
//
// Seventh, the github.EncodedSecret is passed into the GitHub client.Codespaces.CreateOrUpdateUserSecret method to
// populate the secret in the GitHub user.
//
// Finally, if a repo and owner are passed in, it adds the repo to the user secret.
func addUserSecret(ctx context.Context, client *github.Client, secretName, secretValue, owner, repo string) error {
	publicKey, _, err := client.Codespaces.GetUserPublicKey(ctx)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := client.Codespaces.CreateOrUpdateUserSecret(ctx, encryptedSecret); err != nil {
		return fmt.Errorf("client.Codespaces.CreateOrUpdateUserSecret returned error: %v", err)
	}

	if owner != "" && repo != "" {
		r, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			return fmt.Errorf("client.Repositories.Get returned error: %v", err)
		}
		_, err = client.Codespaces.AddSelectedRepoToUserSecret(ctx, encryptedSecret.Name, r)
		if err != nil {
			return fmt.Errorf("client.Codespaces.AddSelectedRepoToUserSecret returned error: %v", err)
		}
		fmt.Printf("Added secret %q to %v/%v\n", secretName, owner, repo)
	}

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	encryptedBytes, err := box.SealAnonymous([]byte{}, []byte(secretValue), &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
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
