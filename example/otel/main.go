// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates how to use the otel transport to instrument
// the go-github client with OpenTelemetry tracing.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/otel/v89"
	"github.com/google/go-github/v89/github"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	// Initialize stdout exporter to see traces in console
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Configure OTel transport
	t := otel.NewTransport(
		http.DefaultTransport,
		otel.WithTracerProvider(tp),
	)

	client, err := github.NewClient(github.WithTransport(t))
	if err != nil {
		log.Fatalf("Error creating GitHub client: %v", err)
	}

	// Make a request (Get Rate Limits is public and cheap)
	limits, resp, err := client.RateLimit.Get(context.Background())
	if err != nil {
		log.Printf("Error fetching rate limits: %v", err)
	} else {
		fmt.Printf("Core Rate Limit: %v/%v (Resets at %v)\n",
			limits.GetCore().Remaining,
			limits.GetCore().Limit,
			limits.GetCore().Reset)
	}

	// Print the HTTP response status when available.
	if resp != nil {
		fmt.Printf("Response Status: %v\n", resp.Status)
	}
}
