package github

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

type Transport struct {
	Username string
	Password string
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(r)
}

func TestListFeeds(t *testing.T) {
	username := os.Getenv("GITHUB_USERNAME")
	password := os.Getenv("GITHUB_PASSWORD")
	if username != "" && password != "" {
		mux.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")

		})
		client := &http.Client{Transport: &Transport{username, password}}
		cl := NewClient(client)
		res, _, err := cl.Feeds.ListFeeds()
		if err != nil {
			t.Errorf("Feeds.ListFeeds returned error: %v", err)
		}
		item := fmt.Sprintf("https://github.com/%s", username)
		if *res.PublicURL != item {
			t.Errorf("Feeds.ListFeeds returned %+v want %s", *res.PublicURL, item)
		}
	}
}
