// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The auditlogstream command demonstrates managing enterprise audit log
// streams for Azure Blob Storage using the go-github library.
//
// Usage — create (github.com):
//
//	export GITHUB_AUTH_TOKEN=<your token>
//	go run main.go create \
//	  -enterprise=my-enterprise \
//	  -container=my-container \
//	  -sas-url=<plain-text-sas-url>
//
// Usage — create (GitHub Enterprise Server):
//
//	export GITHUB_AUTH_TOKEN=<your token>
//	go run main.go create \
//	  -base-url=https://github.example.com/api/v3/ \
//	  -enterprise=my-enterprise \
//	  -container=my-container \
//	  -sas-url=<plain-text-sas-url>
//
// Usage — delete:
//
//	export GITHUB_AUTH_TOKEN=<your token>
//	go run main.go delete \
//	  -base-url=https://github.example.com/api/v3/ \
//	  -enterprise=my-enterprise \
//	  -stream-id=42
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v83/github"
	"golang.org/x/crypto/nacl/box"
)

// encryptSecret encrypts a plain-text secret using libsodium's sealed box
// (crypto_box_seal), which is what GitHub's API expects for encrypted credentials.
func encryptSecret(publicKeyB64, secret string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return "", fmt.Errorf("decoding public key: %w", err)
	}
	if len(publicKeyBytes) != 32 {
		return "", fmt.Errorf("public key must be 32 bytes, got %d", len(publicKeyBytes))
	}
	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)

	encrypted, err := box.SealAnonymous(nil, []byte(secret), &publicKey, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("encrypting secret: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <create|delete> [flags]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		runCreate(os.Args[2:])
	case "delete":
		runDelete(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q. Must be one of: create, delete\n", os.Args[1])
		os.Exit(1)
	}
}

func runCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	baseURL := fs.String("base-url", "https://api.github.com/", "GitHub API base URL. For GitHub Enterprise Server use https://HOSTNAME/api/v3/.")
	enterprise := fs.String("enterprise", "", "Name of the GitHub enterprise slug (required).")
	container := fs.String("container", "", "Azure Blob Storage container name (required).")
	sasURL := fs.String("sas-url", "", "Plain-text Azure SAS URL to encrypt and submit (required).")
	enabled := fs.Bool("enabled", true, "Whether the stream should be enabled immediately.")
	fs.Parse(args)

	token := requireEnv("GITHUB_AUTH_TOKEN")
	requireFlag("enterprise", *enterprise)
	requireFlag("container", *container)
	requireFlag("sas-url", *sasURL)

	ctx := context.Background()
	client := newClient(token, *baseURL)

	// Step 1: Fetch the enterprise's public streaming key.
	streamKey, _, err := client.Enterprise.GetAuditLogStreamKey(ctx, *enterprise)
	if err != nil {
		log.Fatalf("Error fetching audit log stream key: %v", err)
	}
	fmt.Printf("Retrieved stream key ID: %s\n", streamKey.GetKeyID())

	// Step 2: Encrypt the SAS URL using the public key (sealed box / crypto_box_seal).
	encryptedSASURL, err := encryptSecret(streamKey.GetPublicKey(), *sasURL)
	if err != nil {
		log.Fatalf("Error encrypting SAS URL: %v", err)
	}
	fmt.Println("SAS URL encrypted successfully.")

	// Step 3: Create the audit log stream.
	config := github.NewAzureBlobStreamConfig(*enabled, &github.AzureBlobConfig{
		KeyID:           streamKey.KeyID,
		Container:       github.Ptr(*container),
		EncryptedSASURL: github.Ptr(encryptedSASURL),
	})

	stream, _, err := client.Enterprise.CreateAuditLogStream(ctx, *enterprise, config)
	if err != nil {
		log.Fatalf("Error creating audit log stream: %v", err)
	}

	fmt.Printf("Successfully created audit log stream:\n")
	fmt.Printf("  ID:         %d\n", stream.GetID())
	fmt.Printf("  Type:       %s\n", stream.GetStreamType())
	fmt.Printf("  Enabled:    %v\n", stream.GetEnabled())
	fmt.Printf("  Created at: %v\n", stream.GetCreatedAt())
}

func runDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	baseURL := fs.String("base-url", "https://api.github.com/", "GitHub API base URL. For GitHub Enterprise Server use https://HOSTNAME/api/v3/.")
	enterprise := fs.String("enterprise", "", "Name of the GitHub enterprise slug (required).")
	streamID := fs.Int64("stream-id", 0, "ID of the audit log stream to delete (required).")
	fs.Parse(args)

	token := requireEnv("GITHUB_AUTH_TOKEN")
	requireFlag("enterprise", *enterprise)
	if *streamID == 0 {
		log.Fatal("flag -stream-id is required")
	}

	ctx := context.Background()
	client := newClient(token, *baseURL)

	_, err := client.Enterprise.DeleteAuditLogStream(ctx, *enterprise, *streamID)
	if err != nil {
		log.Fatalf("Error deleting audit log stream: %v", err)
	}

	fmt.Printf("Successfully deleted audit log stream %d.\n", *streamID)
}

func newClient(token, baseURL string) *github.Client {
	client, err := github.NewClient(nil).WithAuthToken(token).WithEnterpriseURLs(baseURL, baseURL)
	if err != nil {
		log.Fatalf("Error creating GitHub client: %v", err)
	}
	return client
}

func requireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("environment variable %s is not set", name)
	}
	return val
}

func requireFlag(name, val string) {
	if val == "" {
		log.Fatalf("flag -%s is required", name)
	}
}
