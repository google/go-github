// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v82/github"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type mockTransport struct {
	Response *http.Response
	Err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	// Return valid response with injected headers
	if m.Response != nil {
		return m.Response, nil
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header)}, nil
}

func TestNewTransport_Defaults(t *testing.T) {
	transport := NewTransport(nil)
	if transport.Base != http.DefaultTransport {
		t.Error("NewTransport(nil) should result in http.DefaultTransport")
	}
	if transport.Tracer == nil {
		t.Error("NewTransport(nil) should set default Tracer")
	}
	if transport.Meter == nil {
		t.Error("NewTransport(nil) should set default Meter")
	}
}

func TestRoundTrip_Spans(t *testing.T) {
	// Setup Trace Provider
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))

	// Setup Headers
	headers := http.Header{}
	headers.Set("X-Ratelimit-Limit", "5000")
	headers.Set("X-Ratelimit-Remaining", "4999")
	headers.Set(github.HeaderRateReset, "1372700873")
	headers.Set("X-Github-Request-Id", "1234-5678")
	headers.Set(github.HeaderRateResource, "core")

	mockResp := &http.Response{
		StatusCode: 200,
		Header:     headers,
	}

	transport := NewTransport(
		&mockTransport{Response: mockResp},
		WithTracerProvider(tp),
	)

	req := httptest.NewRequest("GET", "https://api.github.com/users/google", nil)
	_, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("Expected 1 span, got %d", len(spans))
	}
	span := spans[0]

	// Verify Name
	if span.Name != "github/GET" {
		t.Errorf("Expected span name 'github/GET', got '%s'", span.Name)
	}

	// Verify Attributes
	attrs := make(map[attribute.Key]attribute.Value)
	for _, a := range span.Attributes {
		attrs[a.Key] = a.Value
	}

	expectedStringAttrs := map[attribute.Key]string{
		"http.method":                "GET",
		"http.url":                   "https://api.github.com/users/google",
		"http.host":                  "api.github.com",
		"github.rate_limit.reset":    "1372700873",
		"github.request_id":          "1234-5678",
		"github.rate_limit.resource": "core",
	}

	for k, v := range expectedStringAttrs {
		if got, ok := attrs[k]; !ok || got.AsString() != v {
			t.Errorf("Expected attr '%s' = '%s', got '%v'", k, v, got)
		}
	}

	expectedIntAttrs := map[attribute.Key]int64{
		"http.status_code":            200,
		"github.rate_limit.limit":     5000,
		"github.rate_limit.remaining": 4999,
	}

	for k, v := range expectedIntAttrs {
		if got, ok := attrs[k]; !ok || got.AsInt64() != v {
			t.Errorf("Expected attr '%s' = '%d', got '%v'", k, v, got)
		}
	}
}

func TestRoundTrip_Error(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))

	mockErr := errors.New("network failure")
	transport := NewTransport(
		&mockTransport{Err: mockErr},
		WithTracerProvider(tp),
	)

	req := httptest.NewRequest("POST", "https://api.github.com/repos/new", nil)
	_, err := transport.RoundTrip(req)

	if !errors.Is(err, mockErr) {
		t.Errorf("Expected error '%v', got '%v'", mockErr, err)
	}

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("Expected 1 span, got %d", len(spans))
	}
	span := spans[0]

	if span.Status.Code != codes.Error {
		t.Errorf("Expected span status Error, got %v", span.Status.Code)
	}
	if span.Status.Description != "network failure" {
		t.Errorf("Expected span description 'network failure', got '%s'", span.Status.Description)
	}
}

func TestRoundTrip_HTTPError(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))

	mockResp := &http.Response{
		StatusCode: 404,
		Header:     make(http.Header),
	}
	transport := NewTransport(
		&mockTransport{Response: mockResp},
		WithTracerProvider(tp),
	)

	req := httptest.NewRequest("DELETE", "https://api.github.com/user", nil)
	_, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	spans := exporter.GetSpans()
	span := spans[0]

	if span.Status.Code != codes.Error {
		t.Errorf("Expected span status Error, got %v", span.Status.Code)
	}
	if span.Status.Description != "HTTP 404" {
		t.Errorf("Expected span description 'HTTP 404', got '%s'", span.Status.Description)
	}
}

func TestWithMeterProvider(t *testing.T) {
	meter := otel.GetMeterProvider()
	transport := NewTransport(nil, WithMeterProvider(meter))
	if transport.Meter == nil {
		t.Error("WithMeterProvider failed to set Meter")
	}
}
