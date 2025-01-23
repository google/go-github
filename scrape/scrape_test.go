// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scrape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

// setup a test HTTP server along with a scrape.Client that is configured to
// talk to that test server. Tests should register handlers on the mux which
// provide mock responses for the GitHub pages being tested.
func setup(t *testing.T) (client *Client, mux *http.ServeMux) {
	t.Helper()
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient(nil)
	client.baseURL, _ = url.Parse(server.URL + "/")

	t.Cleanup(server.Close)

	return client, mux
}

func copyTestFile(t *testing.T, w io.Writer, filename string) {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("unable to open test file: %v", err)
	}
	_, err = io.Copy(w, f)
	if err != nil {
		t.Errorf("failure copying test file: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Errorf("failure closing test file: %v", err)
	}
}
