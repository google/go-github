// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRoundTripWithOptionalFollowRedirect_HonorsHTTPClientTimeout(t *testing.T) {
	t.Parallel()
	const (
		clientTimeout = 10 * time.Millisecond
		transportWait = 100 * time.Millisecond
	)

	client, mux, _ := setup(t)
	client, err := client.Clone(WithTimeout(clientTimeout))
	if err != nil {
		t.Fatalf("Client.Clone returned error: %v", err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(transportWait)
		w.WriteHeader(http.StatusOK)
	})

	_, err = client.roundTripWithOptionalFollowRedirect(t.Context(), ".", 0)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("roundTripWithOptionalFollowRedirect error = %v, want context deadline exceeded", err)
	}
}
