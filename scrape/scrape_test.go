package scrape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

// setup a test HTTP server along with a scrape.Client that is configured to
// talk to that test server. Tests should register handlers on the mux which
// provide mock responses for the GitHub pages being tested.
func setup() (client *Client, mux *http.ServeMux, cleanup func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient(nil)
	client.baseURL, _ = url.Parse(server.URL + "/")

	return client, mux, server.Close
}

func copyTestFile(t *testing.T, w io.Writer, filename string) {
	t.Helper()
	f, err := os.Open("testdata/" + filename)
	if err != nil {
		t.Errorf("unable to open test file: %v", err)
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
