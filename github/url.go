package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// parseURL parses the given string as a URL and ensures that the path ends
// with a slash. If the input string is empty, it returns an error. If the URL
// cannot be parsed, it returns the parsing error.
func parseURL(s string) (*url.URL, error) {
	if s == "" {
		return nil, errors.New("url cannot be empty")
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u, nil
}
