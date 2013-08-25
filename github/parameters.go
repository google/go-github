package github

import (
	"net/url"
	"strings"
	"time"
)

// Appends parameter to given url if it's value is not nil.
// The value will be URL Encoded.
func withOptionalParameter(baseUrl, parameter string, value *string) string {
	if value == nil {
		return baseUrl
	}

	params := url.Values{
		parameter: []string{*value},
	}

	return appendUrlParams(baseUrl, params)
}

func withOptionalTimeParameter(baseUrl, parameter string, value time.Time) string {
	if value.IsZero() {
		return baseUrl
	}

	params := url.Values{}
	if !value.IsZero() {
		params.Add(parameter, value.Format(time.RFC3339)) // todo, should be ISO 8601???? http://localhost:9292/v3/repos/commits/index.html
	}

	return appendUrlParams(baseUrl, params)
}

func appendUrlParams(baseUrl string, params url.Values) string {
	if strings.Contains(baseUrl, "?") {
		return baseUrl + "&" + params.Encode()
	} else {
		return baseUrl + "?" + params.Encode()
	}
}
