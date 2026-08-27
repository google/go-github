// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// fixedResponseTransport serves a canned JSON payload without touching the
// network so that BenchmarkDo measures the real request/decode path alone.
type fixedResponseTransport struct {
	payload []byte
}

func (t *fixedResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(t.payload)),
		Request:    req,
	}, nil
}

func BenchmarkDo(b *testing.B) {
	for _, sizeKB := range []int{1, 500} {
		b.Run(fmt.Sprintf("%vKB", sizeKB), func(b *testing.B) {
			payload, err := json.Marshal(map[string]string{"body": strings.Repeat("a", sizeKB*1024)})
			if err != nil {
				b.Fatalf("json.Marshal returned error: %v", err)
			}
			client, err := NewClient(WithHTTPClient(&http.Client{
				Transport: &fixedResponseTransport{payload: payload},
			}))
			if err != nil {
				b.Fatalf("NewClient returned error: %v", err)
			}
			req, err := client.NewRequest(b.Context(), "GET", ".", nil)
			if err != nil {
				b.Fatalf("NewRequest returned error: %v", err)
			}

			b.ReportAllocs()
			for b.Loop() {
				var v map[string]string
				if _, err := client.Do(req, &v); err != nil {
					b.Fatalf("Do returned error: %v", err)
				}
			}
		})
	}
}
