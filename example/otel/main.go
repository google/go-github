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

	"github.com/google/go-github/v84/github"
	"github.com/google/go-github/v84/otel"
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

	// Configure HTTP client with OTel transport
	httpClient := &http.Client{
		Transport: otel.NewTransport(
			http.DefaultTransport,
			otel.WithTracerProvider(tp),
		),
	}

	client := github.NewClient(httpClient)

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

	// Check if we captured attributes in response
	if resp != nil {
		fmt.Printf("Response Status: %v\n", resp.Status)
	}
}
