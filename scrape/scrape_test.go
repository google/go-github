package scrape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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

func copyTestFile(w io.Writer, filename string) error {
	f, err := os.Open("testdata/" + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}
