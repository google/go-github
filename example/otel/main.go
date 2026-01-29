// Copyright 2026 The go-github Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v82/github"
	"github.com/google/go-github/v82/otel"
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

	// Create GitHub client
	client := github.NewClient(httpClient)

	// Make a request (Get Rate Limits is public and cheap)
	limits, resp, err := client.RateLimits(context.Background())
	if err != nil {
		log.Printf("Error fetching rate limits: %v", err)
	} else {
		fmt.Printf("Core Rate Limit: %d/%d (Resets at %v)\n", 
			limits.GetCore().Remaining, 
			limits.GetCore().Limit, 
			limits.GetCore().Reset)
	}
    
    // Check if we captured attributes in response
    if resp != nil {
        fmt.Printf("Response Status: %s\n", resp.Status)
    }
}
