// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package otel provides OpenTelemetry instrumentation for the go-github client.
package otel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-github/v84/github"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// instrumentationName is the name of this instrumentation package.
	instrumentationName = "github.com/google/go-github/v84/otel"
)

// Transport is an http.RoundTripper that instrument requests with OpenTelemetry.
type Transport struct {
	Base   http.RoundTripper
	Tracer trace.Tracer
	Meter  metric.Meter
}

// NewTransport creates a new OpenTelemetry transport.
func NewTransport(base http.RoundTripper, opts ...Option) *Transport {
	if base == nil {
		base = http.DefaultTransport
	}
	t := &Transport{Base: base}
	for _, opt := range opts {
		opt(t)
	}
	if t.Tracer == nil {
		t.Tracer = otel.Tracer(instrumentationName)
	}
	if t.Meter == nil {
		t.Meter = otel.Meter(instrumentationName)
	}
	return t
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	spanName := fmt.Sprintf("github/%v", req.Method)
	// Start Span
	ctx, span := t.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// Inject Attributes
	span.SetAttributes(
		attribute.String("http.method", req.Method),
		attribute.String("http.url", req.URL.String()),
		attribute.String("http.host", req.URL.Host),
	)

	// Inject Propagation Headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Execute Request
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Capture response attributes
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	// Capture GitHub Specifics
	if limit := resp.Header.Get(github.HeaderRateLimit); limit != "" {
		if v, err := strconv.Atoi(limit); err == nil {
			span.SetAttributes(attribute.Int("github.rate_limit.limit", v))
		}
	}
	if remaining := resp.Header.Get(github.HeaderRateRemaining); remaining != "" {
		if v, err := strconv.Atoi(remaining); err == nil {
			span.SetAttributes(attribute.Int("github.rate_limit.remaining", v))
		}
	}
	if reset := resp.Header.Get(github.HeaderRateReset); reset != "" {
		span.SetAttributes(attribute.String("github.rate_limit.reset", reset))
	}
	if reqID := resp.Header.Get(github.HeaderRequestID); reqID != "" {
		span.SetAttributes(attribute.String("github.request_id", reqID))
	}
	if resource := resp.Header.Get(github.HeaderRateResource); resource != "" {
		span.SetAttributes(attribute.String("github.rate_limit.resource", resource))
	}

	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %v", resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "OK")
	}

	return resp, nil
}

// Option applies configuration to Transport.
type Option func(*Transport)

// WithTracerProvider configures the TracerProvider.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(t *Transport) {
		t.Tracer = tp.Tracer(instrumentationName)
	}
}

// WithMeterProvider configures the MeterProvider.
func WithMeterProvider(mp metric.MeterProvider) Option {
	return func(t *Transport) {
		t.Meter = mp.Meter(instrumentationName)
	}
}
