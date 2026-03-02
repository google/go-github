// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The auditlogstream command demonstrates managing enterprise audit log
// streams for Azure Blob Storage using the go-github library.
//
// The GitHub API base URL is read from the GITHUB_API_URL environment
// variable. When running inside a GitHub Actions workflow this is set
// automatically.
//
// Usage — create:
//
//	export GITHUB_AUTH_TOKEN=<your token>
//	export GITHUB_API_URL=https://api.<domain>.ghe.com/ or https://domain/api/v3/
//	go run main.go create \
//	  -enterprise=my-enterprise \
//	  -container=my-container \
//	  -sas-url=<plain-text-sas-url>
//
// Usage — delete:
//
//	export GITHUB_AUTH_TOKEN=<your token>
//	export GITHUB_API_URL=https://api.<domain>.ghe.com/ or https://domain/api/v3/
//	go run main.go delete \
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

	"github.com/google/go-github/v84/github"
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
		return "", fmt.Errorf("public key must be 32 bytes, got %v", len(publicKeyBytes))
	}
	publicKey := [32]byte(publicKeyBytes)

	encrypted, err := box.SealAnonymous(nil, []byte(secret), &publicKey, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("encrypting secret: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <create|delete> [flags]\n", os.Args[0])
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

// newFlagSet creates a FlagSet with the common -enterprise flag pre-registered.
func newFlagSet(name string) (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	enterprise := fs.String("enterprise", "", "Enterprise slug (required).")
	return fs, enterprise
}

// parseAndInit parses the FlagSet, validates the enterprise flag, reads
// environment variables, and returns a ready-to-use context, client, and
// enterprise slug.
func parseAndInit(fs *flag.FlagSet, enterprise *string, args []string) (context.Context, *github.Client, string) {
	if err := fs.Parse(args); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	requireFlag("enterprise", *enterprise)

	token := requireEnv("GITHUB_AUTH_TOKEN")
	apiURL := requireEnv("GITHUB_API_URL")

	return context.Background(), newClient(token, apiURL), *enterprise
}

func runCreate(args []string) {
	fs, enterprise := newFlagSet("create")
	container := fs.String("container", "", "Azure Blob Storage container name (required).")
	sasURL := fs.String("sas-url", "", "Plain-text Azure SAS URL to encrypt and submit (required).")
	enabled := fs.Bool("enabled", true, "Whether the stream should be enabled immediately.")

	ctx, client, ent := parseAndInit(fs, enterprise, args)
	requireFlag("container", *container)
	requireFlag("sas-url", *sasURL)

	streamKey, _, err := client.Enterprise.GetAuditLogStreamKey(ctx, ent)
	if err != nil {
		log.Fatalf("Error fetching audit log stream key: %v", err)
	}
	fmt.Printf("Retrieved stream key ID: %v\n", streamKey.KeyID)

	encryptedSASURL, err := encryptSecret(streamKey.Key, *sasURL)
	if err != nil {
		log.Fatalf("Error encrypting SAS URL: %v", err)
	}
	fmt.Println("SAS URL encrypted successfully.")

	config := github.NewAzureBlobStreamConfig(*enabled, &github.AzureBlobConfig{
		KeyID:           streamKey.KeyID,
		Container:       *container,
		EncryptedSASURL: encryptedSASURL,
	})

	stream, _, err := client.Enterprise.CreateAuditLogStream(ctx, ent, *config)
	if err != nil {
		log.Fatalf("Error creating audit log stream: %v", err)
	}

	fmt.Println("Successfully created audit log stream:")
	fmt.Printf("  ID:         %v\n", stream.ID)
	fmt.Printf("  Type:       %v\n", stream.StreamType)
	fmt.Printf("  Enabled:    %v\n", stream.Enabled)
	fmt.Printf("  Created at: %v\n", stream.CreatedAt)
}

func runDelete(args []string) {
	fs, enterprise := newFlagSet("delete")
	streamID := fs.Int64("stream-id", 0, "ID of the audit log stream to delete (required).")

	ctx, client, ent := parseAndInit(fs, enterprise, args)
	requireIntFlag("stream-id", *streamID)

	_, err := client.Enterprise.DeleteAuditLogStream(ctx, ent, *streamID)
	if err != nil {
		log.Fatalf("Error deleting audit log stream: %v", err)
	}

	fmt.Printf("Successfully deleted audit log stream %v.\n", *streamID)
}

func newClient(token, apiURL string) *github.Client {
	client, err := github.NewClient(nil).WithAuthToken(token).WithEnterpriseURLs(apiURL, apiURL)
	if err != nil {
		log.Fatalf("Error creating GitHub client: %v", err)
	}
	return client
}

func requireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("environment variable %v is not set", name)
	}
	return val
}

func requireFlag(name, val string) {
	if val == "" {
		log.Fatalf("flag -%v is required", name)
	}
}

func requireIntFlag(name string, val int64) {
	if val == 0 {
		log.Fatalf("flag -%v is required", name)
	}
}
