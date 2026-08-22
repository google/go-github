// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// legacyDecodeResponse simulates the behavior before Symmetrical Pooling
// (io.ReadAll -> json.Unmarshal).
func legacyDecodeResponse(resp *http.Response, v any) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		return json.Unmarshal(data, v)
	}
	return nil
}

// pooledDecodeResponse simulates the new behavior with Symmetrical Pooling
// (requestBufferPool -> ReadFrom -> json.Unmarshal).
func pooledDecodeResponse(resp *http.Response, v any) error {
	respBuf := requestBufferPool.Get().(*bytes.Buffer)
	defer func() {
		respBuf.Reset()
		requestBufferPool.Put(respBuf)
	}()

	_, err := respBuf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if respBuf.Len() > 0 {
		b := respBuf.Bytes()
		return json.Unmarshal(b, v)
	}
	return nil
}

type dummyReadCloser struct {
	io.Reader
}

func (d *dummyReadCloser) Close() error { return nil }

func BenchmarkDecodeResponse_Legacy(b *testing.B) {
	payload, _ := json.Marshal(map[string]string{"title": "benchmark_test", "body": strings.Repeat("a", 1024*500)}) // 500KB JSON

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		resp := &http.Response{
			Body: &dummyReadCloser{Reader: bytes.NewReader(payload)},
		}
		var v map[string]string
		b.StartTimer()

		_ = legacyDecodeResponse(resp, &v)
	}
}

func BenchmarkDecodeResponse_Pooled(b *testing.B) {
	payload, _ := json.Marshal(map[string]string{"title": "benchmark_test", "body": strings.Repeat("a", 1024*500)}) // 500KB JSON

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		resp := &http.Response{
			Body: &dummyReadCloser{Reader: bytes.NewReader(payload)},
		}
		var v map[string]string
		b.StartTimer()

		_ = pooledDecodeResponse(resp, &v)
	}
}
