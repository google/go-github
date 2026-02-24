// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/api-v3"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (client *Client, mux *http.ServeMux, serverURL string) {
	t.Helper()
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// Create a custom transport with isolated connection pool
	transport := &http.Transport{
		// Controls connection reuse - false allows reuse, true forces new connections for each request
		DisableKeepAlives: false,
		// Maximum concurrent connections per host (active + idle)
		MaxConnsPerHost: 10,
		// Maximum idle connections maintained per host for reuse
		MaxIdleConnsPerHost: 5,
		// Maximum total idle connections across all hosts
		MaxIdleConns: 20,
		// How long an idle connection remains in the pool before being closed
		IdleConnTimeout: 20 * time.Second,
	}

	// Create HTTP client with the isolated transport
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	// client is the GitHub client being tested and is
	// configured to use test server.
	client = NewClient(httpClient)

	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	t.Cleanup(server.Close)

	return client, mux, server.URL
}

// openTestFile creates a new file with the given name and content for testing.
// In order to ensure the exact file name, this function will create a new temp
// directory, and create the file in that directory. The file is automatically
// cleaned up after the test.
func openTestFile(t *testing.T, name, content string) *os.File {
	t.Helper()
	fname := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(fname, []byte(content), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { file.Close() })

	return file
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	assertNilError(t, r.ParseForm())
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error to be returned")
	}
	var uerr *url.Error
	if !errors.As(err, &uerr) || uerr.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %v, want %v", got, want)
	}
}

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v any, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	got, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if diff := cmp.Diff(string(w), string(got)); diff != "" {
		t.Errorf("json.Marshal returned:\n%v\nwant:\n%v\ndiff:\n%v", got, w, diff)
	}
}

// Test how bad options are handled. Method f under test should
// return an error.
func testBadOptions(t *testing.T, methodName string, f func() error) {
	t.Helper()
	if methodName == "" {
		t.Error("testBadOptions: must supply method methodName")
	}
	if err := f(); err == nil {
		t.Errorf("bad options %v err = nil, want error", methodName)
	}
}

// Test function under NewRequest failure and then s.client.Do failure.
// Method f should be a regular call that would normally succeed, but
// should return an error when NewRequest or s.client.Do fails.
func testNewRequestAndDoFailure(t *testing.T, methodName string, client *Client, f func() (*Response, error)) {
	testNewRequestAndDoFailureCategory(t, methodName, client, CoreCategory, f)
}

// testNewRequestAndDoFailureCategory works Like testNewRequestAndDoFailure, but allows setting the category.
func testNewRequestAndDoFailureCategory(t *testing.T, methodName string, client *Client, category RateLimitCategory, f func() (*Response, error)) {
	t.Helper()
	if methodName == "" {
		t.Error("testNewRequestAndDoFailure: must supply method methodName")
	}

	client.BaseURL.Path = ""
	resp, err := f()
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' %v resp = %#v, want nil", methodName, resp)
	}
	if err == nil {
		t.Errorf("client.BaseURL.Path='' %v err = nil, want error", methodName)
	}

	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[category].Reset.Time = time.Now().Add(10 * time.Minute)
	resp, err = f()
	if client.DisableRateLimitCheck {
		return
	}
	if bypass := resp.Request.Context().Value(BypassRateLimitCheck); bypass != nil {
		return
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		if resp != nil {
			t.Errorf("rate.Reset.Time > now %v resp = %#v, want StatusCode=%v", methodName, resp.Response, want)
		} else {
			t.Errorf("rate.Reset.Time > now %v resp = nil, want StatusCode=%v", methodName, want)
		}
	}
	if err == nil {
		t.Errorf("rate.Reset.Time > now %v err = nil, want error", methodName)
	}
}

// Test that all error response types contain the status code.
func testErrorResponseForStatusCode(t *testing.T, code int) {
	t.Helper()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(code)
	})

	ctx := t.Context()
	_, _, err := client.Repositories.ListHooks(ctx, "o", "r", nil)

	var abuseErr *AbuseRateLimitError
	switch {
	case errors.As(err, new(*ErrorResponse)):
	case errors.As(err, new(*RateLimitError)):
	case errors.As(err, &abuseErr):
		if code != abuseErr.Response.StatusCode {
			t.Error("Error response does not contain status code")
		}
	default:
		t.Error("Unknown error response type")
	}
}

func assertNoDiff(t *testing.T, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff mismatch (-want +got):\n%v", diff)
	}
}

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertWrite(t *testing.T, w io.Writer, data []byte) {
	t.Helper()
	_, err := w.Write(data)
	assertNilError(t, err)
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, defaultUserAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewClientWithEnvProxy(t *testing.T) {
	t.Parallel()
	client := NewClientWithEnvProxy()
	if got, want := client.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestClient(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	c2 := c.Client()
	if c.client == c2 {
		t.Error("Client returned same http.Client, but should be different")
	}
}

func TestWithAuthToken(t *testing.T) {
	t.Parallel()
	token := "gh_test_token"

	validate := func(t *testing.T, c *http.Client, token string) {
		t.Helper()
		want := token
		if want != "" {
			want = "Bearer " + want
		}
		gotReq := false
		headerVal := ""
		srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			gotReq = true
			headerVal = r.Header.Get("Authorization")
		}))
		_, err := c.Get(srv.URL)
		assertNilError(t, err)
		if !gotReq {
			t.Error("request not sent")
		}
		if headerVal != want {
			t.Errorf("Authorization header is %v, want %v", headerVal, want)
		}
	}

	t.Run("zero-value Client", func(t *testing.T) {
		t.Parallel()
		c := new(Client).WithAuthToken(token)
		validate(t, c.Client(), token)
	})

	t.Run("NewClient", func(t *testing.T) {
		t.Parallel()
		httpClient := &http.Client{}
		client := NewClient(httpClient).WithAuthToken(token)
		validate(t, client.Client(), token)
		// make sure the original client isn't setting auth headers now
		validate(t, httpClient, "")
	})

	t.Run("NewTokenClient", func(t *testing.T) {
		t.Parallel()
		validate(t, NewTokenClient(t.Context(), token).Client(), token)
	})

	t.Run("do not set Authorization when empty token", func(t *testing.T) {
		t.Parallel()
		c := new(Client).WithAuthToken("")

		gotReq := false
		ifAuthorizationSet := false
		srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			gotReq = true
			_, ifAuthorizationSet = r.Header["Authorization"]
		}))
		_, err := c.client.Get(srv.URL)
		assertNilError(t, err)
		if !gotReq {
			t.Error("request not sent")
		}
		if ifAuthorizationSet {
			t.Error("The header 'Authorization' must not be set")
		}
	})
}

func TestWithEnterpriseURLs(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		baseURL       string
		wantBaseURL   string
		uploadURL     string
		wantUploadURL string
		wantErr       string
	}{
		{
			name:          "does not modify properly formed URLs",
			baseURL:       "https://custom-url/api/v3/",
			wantBaseURL:   "https://custom-url/api/v3/",
			uploadURL:     "https://custom-upload-url/api/uploads/",
			wantUploadURL: "https://custom-upload-url/api/uploads/",
		},
		{
			name:          "adds trailing slash",
			baseURL:       "https://custom-url/api/v3",
			wantBaseURL:   "https://custom-url/api/v3/",
			uploadURL:     "https://custom-upload-url/api/uploads",
			wantUploadURL: "https://custom-upload-url/api/uploads/",
		},
		{
			name:          "adds enterprise suffix",
			baseURL:       "https://custom-url/",
			wantBaseURL:   "https://custom-url/api/v3/",
			uploadURL:     "https://custom-upload-url/",
			wantUploadURL: "https://custom-upload-url/api/uploads/",
		},
		{
			name:          "adds enterprise suffix and trailing slash",
			baseURL:       "https://custom-url",
			wantBaseURL:   "https://custom-url/api/v3/",
			uploadURL:     "https://custom-upload-url",
			wantUploadURL: "https://custom-upload-url/api/uploads/",
		},
		{
			name:      "bad base URL",
			baseURL:   "bogus\nbase\nURL",
			uploadURL: "https://custom-upload-url/api/uploads/",
			wantErr:   `invalid control character in URL`,
		},
		{
			name:      "bad upload URL",
			baseURL:   "https://custom-url/api/v3/",
			uploadURL: "bogus\nupload\nURL",
			wantErr:   `invalid control character in URL`,
		},
		{
			name:          "URL has existing API prefix, adds trailing slash",
			baseURL:       "https://api.custom-url",
			wantBaseURL:   "https://api.custom-url/",
			uploadURL:     "https://api.custom-upload-url",
			wantUploadURL: "https://api.custom-upload-url/",
		},
		{
			name:          "URL has existing API prefix and trailing slash",
			baseURL:       "https://api.custom-url/",
			wantBaseURL:   "https://api.custom-url/",
			uploadURL:     "https://api.custom-upload-url/",
			wantUploadURL: "https://api.custom-upload-url/",
		},
		{
			name:          "URL has API subdomain, adds trailing slash",
			baseURL:       "https://catalog.api.custom-url",
			wantBaseURL:   "https://catalog.api.custom-url/",
			uploadURL:     "https://catalog.api.custom-upload-url",
			wantUploadURL: "https://catalog.api.custom-upload-url/",
		},
		{
			name:          "URL has API subdomain and trailing slash",
			baseURL:       "https://catalog.api.custom-url/",
			wantBaseURL:   "https://catalog.api.custom-url/",
			uploadURL:     "https://catalog.api.custom-upload-url/",
			wantUploadURL: "https://catalog.api.custom-upload-url/",
		},
		{
			name:          "URL is not a proper API subdomain, adds enterprise suffix and slash",
			baseURL:       "https://cloud-api.custom-url",
			wantBaseURL:   "https://cloud-api.custom-url/api/v3/",
			uploadURL:     "https://cloud-api.custom-upload-url",
			wantUploadURL: "https://cloud-api.custom-upload-url/api/uploads/",
		},
		{
			name:          "URL is not a proper API subdomain, adds enterprise suffix",
			baseURL:       "https://cloud-api.custom-url/",
			wantBaseURL:   "https://cloud-api.custom-url/api/v3/",
			uploadURL:     "https://cloud-api.custom-upload-url/",
			wantUploadURL: "https://cloud-api.custom-upload-url/api/uploads/",
		},
		{
			name:          "URL has uploads subdomain, does not modify",
			baseURL:       "https://api.custom-url/",
			wantBaseURL:   "https://api.custom-url/",
			uploadURL:     "https://uploads.custom-upload-url/",
			wantUploadURL: "https://uploads.custom-upload-url/",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validate := func(c *Client, err error) {
				t.Helper()
				if test.wantErr != "" {
					if err == nil || !strings.Contains(err.Error(), test.wantErr) {
						t.Fatalf("error does not contain expected string %q: %v", test.wantErr, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("got unexpected error: %v", err)
				}
				if c.BaseURL.String() != test.wantBaseURL {
					t.Errorf("BaseURL is %v, want %v", c.BaseURL, test.wantBaseURL)
				}
				if c.UploadURL.String() != test.wantUploadURL {
					t.Errorf("UploadURL is %v, want %v", c.UploadURL, test.wantUploadURL)
				}
			}
			validate(NewClient(nil).WithEnterpriseURLs(test.baseURL, test.uploadURL))
			validate(new(Client).WithEnterpriseURLs(test.baseURL, test.uploadURL))
			validate(NewEnterpriseClient(test.baseURL, test.uploadURL, nil))
		})
	}
}

// Ensure that length of Client.rateLimits is the same as number of fields in RateLimits struct.
func TestClient_rateLimits(t *testing.T) {
	t.Parallel()
	if got, want := len(Client{}.rateLimits), reflect.TypeFor[RateLimits]().NumField(); got != want {
		t.Errorf("len(Client{}.rateLimits) is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &User{Login: Ptr("l")}, `{"login":"l"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	userAgent := req.Header.Get("User-Agent")

	// test that default user-agent is attached to the request
	if got, want := userAgent, c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}

	if !strings.Contains(userAgent, Version) {
		t.Errorf("NewRequest() User-Agent should contain %v, found %v", Version, userAgent)
	}

	apiVersion := req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, defaultAPIVersion; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}

	req, _ = c.NewRequest("GET", inURL, inBody, WithVersion("2022-11-29"))
	apiVersion = req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, "2022-11-29"; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)

	type T struct {
		F func()
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
}

func TestNewRequest_badURL(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_badMethod(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	if _, err := c.NewRequest("BOGUS\nMETHOD", ".", nil); err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}
}

// ensure that no User-Agent header is set if the client's UserAgent is empty.
// This caused a problem with Google's internal http client.
func TestNewRequest_emptyUserAgent(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	c.UserAgent = ""
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if _, ok := req.Header["User-Agent"]; ok {
		t.Fatal("constructed request contains unexpected User-Agent header")
	}
}

// If a nil body is passed to github.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request
// body that is an empty string versus one that is not set at all. However in
// certain cases, intermediate systems may treat these differently resulting in
// subtle errors.
func TestNewRequest_emptyBody(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatal("constructed request contains a non-nil Body")
	}
}

func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api/v3", wantError: true},
		{rawurl: "https://example.com/api/v3/", wantError: false},
	}
	c := NewClient(nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewRequest("GET", "test", nil); test.wantError && err == nil {
			t.Fatal("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewRequest returned unexpected error: %v.", err)
		}
	}
}

func TestNewFormRequest(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	form := url.Values{}
	form.Add("login", "l")
	inBody, outBody := strings.NewReader(form.Encode()), "login=l"
	req, _ := c.NewFormRequest(inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewFormRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was form encoded
	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewFormRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewFormRequest() User-Agent is %v, want %v", got, want)
	}

	apiVersion := req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, defaultAPIVersion; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}

	req, _ = c.NewFormRequest(inURL, inBody, WithVersion("2022-11-29"))
	apiVersion = req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, "2022-11-29"; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}
}

func TestNewFormRequest_badURL(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	_, err := c.NewFormRequest(":", nil)
	testURLParseError(t, err)
}

func TestNewFormRequest_emptyUserAgent(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	c.UserAgent = ""
	req, err := c.NewFormRequest(".", nil)
	if err != nil {
		t.Fatalf("NewFormRequest returned unexpected error: %v", err)
	}
	if _, ok := req.Header["User-Agent"]; ok {
		t.Fatal("constructed request contains unexpected User-Agent header")
	}
}

func TestNewFormRequest_emptyBody(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	req, err := c.NewFormRequest(".", nil)
	if err != nil {
		t.Fatalf("NewFormRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatal("constructed request contains a non-nil Body")
	}
}

func TestNewFormRequest_errorForNoTrailingSlash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		rawURL    string
		wantError bool
	}{
		{rawURL: "https://example.com/api/v3", wantError: true},
		{rawURL: "https://example.com/api/v3/", wantError: false},
	}
	c := NewClient(nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawURL)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewFormRequest("test", nil); test.wantError && err == nil {
			t.Fatal("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewFormRequest returned unexpected error: %v.", err)
		}
	}
}

func TestNewUploadRequest_WithVersion(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	req, _ := c.NewUploadRequest("https://example.com/", nil, 0, "")

	apiVersion := req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, defaultAPIVersion; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}

	req, _ = c.NewUploadRequest("https://example.com/", nil, 0, "", WithVersion("2022-11-29"))
	apiVersion = req.Header.Get(headerAPIVersion)
	if got, want := apiVersion, "2022-11-29"; got != want {
		t.Errorf("NewRequest() %v header is %v, want %v", headerAPIVersion, got, want)
	}
}

func TestNewUploadRequest_badURL(t *testing.T) {
	t.Parallel()
	c := NewClient(nil)
	_, err := c.NewUploadRequest(":", nil, 0, "")
	testURLParseError(t, err)

	const methodName = "NewUploadRequest"
	testBadOptions(t, methodName, func() (err error) {
		_, err = c.NewUploadRequest("\n", nil, -1, "\n")
		return err
	})
}

func TestNewUploadRequest_errorForNoTrailingSlash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api/uploads", wantError: true},
		{rawurl: "https://example.com/api/uploads/", wantError: false},
	}
	c := NewClient(nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.UploadURL = u
		if _, err = c.NewUploadRequest("test", nil, 0, ""); test.wantError && err == nil {
			t.Fatal("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewUploadRequest returned unexpected error: %v.", err)
		}
	}
}

func TestResponse_populatePageValues(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?page=1>; rel="first",` +
					` <https://api.github.com/?page=2>; rel="prev",` +
					` <https://api.github.com/?page=4>; rel="next",` +
					` <https://api.github.com/?page=5>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 1; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 2; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 4; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 5; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; want != got {
		t.Errorf("response.NextPageToken: %v, want %v", got, want)
	}
}

func TestResponse_populateSinceValues(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?since=1>; rel="first",` +
					` <https://api.github.com/?since=2>; rel="prev",` +
					` <https://api.github.com/?since=4>; rel="next",` +
					` <https://api.github.com/?since=5>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 1; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 2; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 4; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 5; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; want != got {
		t.Errorf("response.NextPageToken: %v, want %v", got, want)
	}
}

func TestResponse_SinceWithPage(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?since=2021-12-04T10%3A43%3A42Z&page=1>; rel="first",` +
					` <https://api.github.com/?since=2021-12-04T10%3A43%3A42Z&page=2>; rel="prev",` +
					` <https://api.github.com/?since=2021-12-04T10%3A43%3A42Z&page=4>; rel="next",` +
					` <https://api.github.com/?since=2021-12-04T10%3A43%3A42Z&page=5>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 1; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 2; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 4; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 5; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; want != got {
		t.Errorf("response.NextPageToken: %v, want %v", got, want)
	}
}

func TestResponse_cursorPagination(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Status": {"200 OK"},
			"Link":   {`<https://api.github.com/resource?per_page=2&page=url-encoded-next-page-token>; rel="next"`},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 0; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 0; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 0; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, "url-encoded-next-page-token"; want != got {
		t.Errorf("response.NextPageToken: %v, want %v", got, want)
	}

	// cursor-based pagination with "cursor" param
	r = http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?cursor=v1_12345678>; rel="next"`,
			},
		},
	}

	response = newResponse(&r)
	if got, want := response.Cursor, "v1_12345678"; got != want {
		t.Errorf("response.Cursor: %v, want %v", got, want)
	}
}

func TestResponse_beforeAfterPagination(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?after=a1b2c3&before=>; rel="next",` +
					` <https://api.github.com/?after=&before=>; rel="first",` +
					` <https://api.github.com/?after=&before=d4e5f6>; rel="prev",`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.Before, "d4e5f6"; got != want {
		t.Errorf("response.Before: %v, want %v", got, want)
	}
	if got, want := response.After, "a1b2c3"; got != want {
		t.Errorf("response.After: %v, want %v", got, want)
	}
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 0; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 0; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 0; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; want != got {
		t.Errorf("response.NextPageToken: %v, want %v", got, want)
	}
}

func TestResponse_populatePageValues_invalid(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?page=1>,` +
					`<https://api.github.com/?page=abc>; rel="first",` +
					`https://api.github.com/?page=2; rel="prev",` +
					`<https://api.github.com/>; rel="next",` +
					`<https://api.github.com/?page=>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 0; got != want {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 0; got != want {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 0; got != want {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}

	// more invalid URLs
	r = http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/%?page=2>; rel="first"`},
		},
	}

	response = newResponse(&r)
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
}

func TestResponse_populateSinceValues_invalid(t *testing.T) {
	t.Parallel()
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/?since=1>,` +
					`<https://api.github.com/?since=abc>; rel="first",` +
					`https://api.github.com/?since=2; rel="prev",` +
					`<https://api.github.com/>; rel="next",` +
					`<https://api.github.com/?since=>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 0; got != want {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 0; got != want {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 0; got != want {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}

	// more invalid URLs
	r = http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/%?since=2>; rel="first"`},
		},
	}

	response = newResponse(&r)
	if got, want := response.FirstPage, 0; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	ctx := t.Context()
	_, err := client.Do(ctx, req, body)
	assertNilError(t, err)

	want := &foo{"a"}
	if !cmp.Equal(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_nilContext(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(nil, req, nil)

	if !errors.Is(err, errNonNilContext) {
		t.Error("Expected context must be non-nil error")
	}
}

func TestDo_httpError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %v status code.", resp.StatusCode)
	}
}

// Test handling of an error caused by the internal http client's Do()
// function. A redirect loop is pretty unlikely to occur within the GitHub
// API, but does allow us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusFound)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if !errors.As(err, new(*url.Error)) {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_preservesResponseInHTTPError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{
			"message": "Resource not found",
			"documentation_url": "https://docs.github.com/rest/reference/repos#get-a-repository"
		}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	var resp *Response
	var data any
	resp, err := client.Do(t.Context(), req, &data)

	if err == nil {
		t.Fatal("Expected error response")
	}

	// Verify error type and access to status code
	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("Expected *ErrorResponse error, got %T", err)
	}

	// Verify status code is accessible from both Response and ErrorResponse
	if resp == nil {
		t.Fatal("Expected response to be returned even with error")
	}
	if got, want := resp.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Response status = %v, want %v", got, want)
	}
	if got, want := errResp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Error response status = %v, want %v", got, want)
	}

	// Verify error contains proper message
	if !strings.Contains(errResp.Message, "Resource not found") {
		t.Errorf("Error message = %q, want to contain 'Resource not found'", errResp.Message)
	}
}

// Test that an error caused by the internal http client's Do() function
// does not leak the client secret.
func TestDo_sanitizeURL(t *testing.T) {
	t.Parallel()
	tp := &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	unauthedClient := NewClient(tp.Client())
	unauthedClient.BaseURL = &url.URL{Scheme: "http", Host: "127.0.0.1:0", Path: "/"} // Use port 0 on purpose to trigger a dial TCP error, expect to get "dial tcp 127.0.0.1:0: connect: can't assign requested address".
	req, err := unauthedClient.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	ctx := t.Context()
	_, err = unauthedClient.Do(ctx, req, nil)
	if err == nil {
		t.Fatal("Expected error to be returned.")
	}
	if strings.Contains(err.Error(), "client_secret=secret") {
		t.Errorf("Do error contains secret, should be redacted:\n%q", err)
	}
}

func TestDo_rateLimit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "59")
		w.Header().Set(HeaderRateUsed, "1")
		w.Header().Set(HeaderRateReset, "1372700873")
		w.Header().Set(HeaderRateResource, "core")
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	if got, want := resp.Rate.Limit, 60; got != want {
		t.Errorf("Client rate limit = %v, want %v", got, want)
	}
	if got, want := resp.Rate.Remaining, 59; got != want {
		t.Errorf("Client rate remaining = %v, want %v", got, want)
	}
	if got, want := resp.Rate.Used, 1; got != want {
		t.Errorf("Client rate used = %v, want %v", got, want)
	}
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if !resp.Rate.Reset.UTC().Equal(reset) {
		t.Errorf("Client rate reset = %v, want %v", resp.Rate.Reset.UTC(), reset)
	}
	if got, want := resp.Rate.Resource, "core"; got != want {
		t.Errorf("Client rate resource = %v, want %v", got, want)
	}
}

func TestDo_rateLimitCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method   string
		url      string
		category RateLimitCategory
	}{
		{
			method:   "GET",
			url:      "/",
			category: CoreCategory,
		},
		{
			method:   "GET",
			url:      "/search/issues?q=rate",
			category: SearchCategory,
		},
		{
			method:   "GET",
			url:      "/graphql",
			category: GraphqlCategory,
		},
		{
			method:   "POST",
			url:      "/app-manifests/code/conversions",
			category: IntegrationManifestCategory,
		},
		{
			method:   "GET",
			url:      "/app-manifests/code/conversions",
			category: CoreCategory, // only POST requests are in the integration manifest category
		},
		{
			method:   "PUT",
			url:      "/repos/google/go-github/import",
			category: SourceImportCategory,
		},
		{
			method:   "GET",
			url:      "/repos/google/go-github/import",
			category: CoreCategory, // only PUT requests are in the source import category
		},
		{
			method:   "POST",
			url:      "/repos/google/go-github/code-scanning/sarifs",
			category: CodeScanningUploadCategory,
		},
		{
			method:   "GET",
			url:      "/scim/v2/organizations/ORG/Users",
			category: ScimCategory,
		},
		{
			method:   "POST",
			url:      "/repos/google/go-github/dependency-graph/snapshots",
			category: DependencySnapshotsCategory,
		},
		{
			method:   "GET",
			url:      "/search/code?q=rate",
			category: CodeSearchCategory,
		},
		{
			method:   "GET",
			url:      "/orgs/google/audit-log",
			category: AuditLogCategory,
		},
		{
			method:   "GET",
			url:      "/repos/google/go-github/dependency-graph/sbom",
			category: DependencySBOMCategory,
		},
		// missing a check for actionsRunnerRegistrationCategory: API not found
	}

	for _, tt := range tests {
		if got, want := GetRateLimitCategory(tt.method, tt.url), tt.category; got != want {
			t.Errorf("expecting category %v, found %v", got, want)
		}
	}
}

// Ensure rate limit is still parsed, even for error responses.
func TestDo_rateLimit_errorResponse(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "59")
		w.Header().Set(HeaderRateUsed, "1")
		w.Header().Set(HeaderRateReset, "1372700873")
		w.Header().Set(HeaderRateResource, "core")
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if errors.As(err, new(*RateLimitError)) {
		t.Errorf("Did not expect a *RateLimitError error; got %#v.", err)
	}
	if got, want := resp.Rate.Limit, 60; got != want {
		t.Errorf("Client rate limit = %v, want %v", got, want)
	}
	if got, want := resp.Rate.Remaining, 59; got != want {
		t.Errorf("Client rate remaining = %v, want %v", got, want)
	}
	if got, want := resp.Rate.Used, 1; got != want {
		t.Errorf("Client rate used = %v, want %v", got, want)
	}
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if !resp.Rate.Reset.UTC().Equal(reset) {
		t.Errorf("Client rate reset = %v, want %v", resp.Rate.Reset, reset)
	}
	if got, want := resp.Rate.Resource, "core"; got != want {
		t.Errorf("Client rate resource = %v, want %v", got, want)
	}
}

// Ensure *RateLimitError is returned when API rate limit is exceeded.
func TestDo_rateLimit_rateLimitError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "0")
		w.Header().Set(HeaderRateUsed, "60")
		w.Header().Set(HeaderRateReset, "1372700873")
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("Expected a *RateLimitError error; got %#v.", err)
	}
	if got, want := rateLimitErr.Rate.Limit, 60; got != want {
		t.Errorf("rateLimitErr rate limit = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Remaining, 0; got != want {
		t.Errorf("rateLimitErr rate remaining = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Used, 60; got != want {
		t.Errorf("rateLimitErr rate used = %v, want %v", got, want)
	}
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if !rateLimitErr.Rate.Reset.UTC().Equal(reset) {
		t.Errorf("rateLimitErr rate reset = %v, want %v", rateLimitErr.Rate.Reset.UTC(), reset)
	}
	if got, want := rateLimitErr.Rate.Resource, "core"; got != want {
		t.Errorf("rateLimitErr rate resource = %v, want %v", got, want)
	}
}

// Ensure a network call is not made when it's known that API rate limit is still exceeded.
func TestDo_rateLimit_noNetworkCall(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Minute).Round(time.Second) // Rate reset is a minute from now, with 1 second precision.

	mux.HandleFunc("/first", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "0")
		w.Header().Set(HeaderRateUsed, "60")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	madeNetworkCall := false
	mux.HandleFunc("/second", func(http.ResponseWriter, *http.Request) {
		madeNetworkCall = true
	})

	// First request is made, and it makes the client aware of rate reset time being in the future.
	req, _ := client.NewRequest("GET", "first", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}

	// Second request should not cause a network call to be made, since client can predict a rate limit error.
	req, _ = client.NewRequest("GET", "second", nil)
	_, err = client.Do(ctx, req, nil)

	if madeNetworkCall {
		t.Fatal("Network call was made, even though rate limit is known to still be exceeded.")
	}

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("Expected a *RateLimitError error; got %#v.", err)
	}
	if got, want := rateLimitErr.Rate.Limit, 60; got != want {
		t.Errorf("rateLimitErr rate limit = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Remaining, 0; got != want {
		t.Errorf("rateLimitErr rate remaining = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Used, 60; got != want {
		t.Errorf("rateLimitErr rate used = %v, want %v", got, want)
	}
	if !rateLimitErr.Rate.Reset.UTC().Equal(reset) {
		t.Errorf("rateLimitErr rate reset = %v, want %v", rateLimitErr.Rate.Reset.UTC(), reset)
	}
	if got, want := rateLimitErr.Rate.Resource, "core"; got != want {
		t.Errorf("rateLimitErr rate resource = %v, want %v", got, want)
	}
}

// Ignore rate limit headers if the response was served from cache.
func TestDo_rateLimit_ignoredFromCache(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Minute).Round(time.Second) // Rate reset is a minute from now, with 1 second precision.

	// By adding the X-From-Cache header we pretend this is served from a cache.
	mux.HandleFunc("/first", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-From-Cache", "1")
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "0")
		w.Header().Set(HeaderRateUsed, "60")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	madeNetworkCall := false
	mux.HandleFunc("/second", func(http.ResponseWriter, *http.Request) {
		madeNetworkCall = true
	})

	// First request is made so afterwards we can check the returned rate limit headers were ignored.
	req, _ := client.NewRequest("GET", "first", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}

	// Second request should not by hindered by rate limits.
	req, _ = client.NewRequest("GET", "second", nil)
	_, err = client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Second request failed, even though the rate limits from the cache should've been ignored: %v", err)
	}
	if !madeNetworkCall {
		t.Fatal("Network call was not made, even though the rate limits from the cache should've been ignored")
	}
}

// Ensure sleeps until the rate limit is reset when the client is rate limited.
func TestDo_rateLimit_sleepUntilResponseResetLimit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Second)

	firstRequest := true
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		if firstRequest {
			firstRequest = false
			w.Header().Set(HeaderRateLimit, "60")
			w.Header().Set(HeaderRateRemaining, "0")
			w.Header().Set(HeaderRateUsed, "60")
			w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Unix()))
			w.Header().Set(HeaderRateResource, "core")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
			return
		}
		w.Header().Set(HeaderRateLimit, "5000")
		w.Header().Set(HeaderRateRemaining, "5000")
		w.Header().Set(HeaderRateUsed, "0")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Add(time.Hour).Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(context.WithValue(ctx, SleepUntilPrimaryRateLimitResetWhenRateLimited, true), req, nil)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("Response status code = %v, want %v", got, want)
	}
}

// Ensure tries to sleep until the rate limit is reset when the client is rate limited, but only once.
func TestDo_rateLimit_sleepUntilResponseResetLimitRetryOnce(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Second)

	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "0")
		w.Header().Set(HeaderRateUsed, "60")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(context.WithValue(ctx, SleepUntilPrimaryRateLimitResetWhenRateLimited, true), req, nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if got, want := requestCount, 2; got != want {
		t.Errorf("Expected 2 requests, got %v", got)
	}
}

// Ensure a network call is not made when it's known that API rate limit is still exceeded.
func TestDo_rateLimit_sleepUntilClientResetLimit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Second)
	client.rateLimits[CoreCategory] = Rate{Limit: 5000, Remaining: 0, Reset: Timestamp{reset}}
	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "5000")
		w.Header().Set(HeaderRateRemaining, "5000")
		w.Header().Set(HeaderRateUsed, "0")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Add(time.Hour).Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	})
	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(context.WithValue(ctx, SleepUntilPrimaryRateLimitResetWhenRateLimited, true), req, nil)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("Response status code = %v, want %v", got, want)
	}
	if got, want := requestCount, 1; got != want {
		t.Errorf("Expected 1 request, got %v", got)
	}
}

// Ensure sleep is aborted when the context is cancelled.
func TestDo_rateLimit_abortSleepContextCancelled(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// We use a 1 minute reset time to ensure the sleep is not completed.
	reset := time.Now().UTC().Add(time.Minute)
	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "60")
		w.Header().Set(HeaderRateRemaining, "0")
		w.Header().Set(HeaderRateUsed, "60")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	_, err := client.Do(context.WithValue(ctx, SleepUntilPrimaryRateLimitResetWhenRateLimited, true), req, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("Expected context deadline exceeded error.")
	}
	if got, want := requestCount, 1; got != want {
		t.Errorf("Expected 1 requests, got %v", got)
	}
}

// Ensure sleep is aborted when the context is cancelled on initial request.
func TestDo_rateLimit_abortSleepContextCancelledClientLimit(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(time.Minute)
	client.rateLimits[CoreCategory] = Rate{Limit: 5000, Remaining: 0, Reset: Timestamp{reset}}
	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "5000")
		w.Header().Set(HeaderRateRemaining, "5000")
		w.Header().Set(HeaderRateUsed, "0")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Add(time.Hour).Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	})
	req, _ := client.NewRequest("GET", ".", nil)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	_, err := client.Do(context.WithValue(ctx, SleepUntilPrimaryRateLimitResetWhenRateLimited, true), req, nil)
	var rateLimitError *RateLimitError
	if !errors.As(err, &rateLimitError) {
		t.Fatalf("Expected a *rateLimitError error; got %#v.", err)
	}
	if got, wantSuffix := rateLimitError.Message, "Context cancelled while waiting for rate limit to reset until"; !strings.HasPrefix(got, wantSuffix) {
		t.Errorf("Expected request to be prevented because context cancellation, got: %v.", got)
	}
	if got, want := requestCount, 0; got != want {
		t.Errorf("Expected 1 requests, got %v", got)
	}
}

// Ensure *AbuseRateLimitError is returned when the response indicates that
// the client has triggered an abuse detection mechanism.
func TestDo_rateLimit_abuseRateLimitError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		// When the abuse rate limit error is of the "temporarily blocked from content creation" type,
		// there is no "Retry-After" header.
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism and have been temporarily blocked from content creation. Please retry your request again later.",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var abuseRateLimitErr *AbuseRateLimitError
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if got, want := abuseRateLimitErr.RetryAfter, (*time.Duration)(nil); got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

// Ensure *AbuseRateLimitError is returned when the response indicates that
// the client has triggered an abuse detection mechanism on GitHub Enterprise.
func TestDo_rateLimit_abuseRateLimitErrorEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		// When the abuse rate limit error is of the "temporarily blocked from content creation" type,
		// there is no "Retry-After" header.
		// This response returns a documentation url like the one returned for GitHub Enterprise, this
		// url changes between versions but follows roughly the same format.
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism and have been temporarily blocked from content creation. Please retry your request again later.",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var abuseRateLimitErr *AbuseRateLimitError
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if got, want := abuseRateLimitErr.RetryAfter, (*time.Duration)(nil); got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

// Ensure *AbuseRateLimitError.RetryAfter is parsed correctly for the Retry-After header.
func TestDo_rateLimit_abuseRateLimitError_retryAfter(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set(headerRetryAfter, "123") // Retry after value of 123 seconds.
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism ...",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var abuseRateLimitErr *AbuseRateLimitError
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatal("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	if got, want := *abuseRateLimitErr.RetryAfter, 123*time.Second; got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}

	// expect prevention of a following request
	if _, err = client.Do(ctx, req, nil); err == nil {
		t.Error("Expected error to be returned.")
	}
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatal("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	// the saved duration might be a bit smaller than Retry-After because the duration is calculated from the expected end-of-cooldown time
	if got, want := *abuseRateLimitErr.RetryAfter, 123*time.Second; want-got > 1*time.Second {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
	if got, wantSuffix := abuseRateLimitErr.Message, "not making remote request."; !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("Expected request to be prevented because of secondary rate limit, got: %v.", got)
	}
}

// Ensure *AbuseRateLimitError.RetryAfter is parsed correctly for the x-ratelimit-reset header.
func TestDo_rateLimit_abuseRateLimitError_xRateLimitReset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	// x-ratelimit-reset value of 123 seconds into the future.
	blockUntil := time.Now().Add(time.Duration(123) * time.Second).Unix()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set(HeaderRateReset, strconv.Itoa(int(blockUntil)))
		w.Header().Set(HeaderRateRemaining, "1") // set remaining to a value > 0 to distinct from a primary rate limit
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism ...",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var abuseRateLimitErr *AbuseRateLimitError
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatal("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	// the retry after value might be a bit smaller than the original duration because the duration is calculated from the expected end-of-cooldown time
	if got, want := *abuseRateLimitErr.RetryAfter, 123*time.Second; want-got > 1*time.Second {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}

	// expect prevention of a following request
	if _, err = client.Do(ctx, req, nil); err == nil {
		t.Error("Expected error to be returned.")
	}
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatal("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	// the saved duration might be a bit smaller than Retry-After because the duration is calculated from the expected end-of-cooldown time
	if got, want := *abuseRateLimitErr.RetryAfter, 123*time.Second; want-got > 1*time.Second {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
	if got, wantSuffix := abuseRateLimitErr.Message, "not making remote request."; !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("Expected request to be prevented because of secondary rate limit, got: %v.", got)
	}
}

// Ensure *AbuseRateLimitError.RetryAfter respect a max duration if specified.
func TestDo_rateLimit_abuseRateLimitError_maxDuration(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	// specify a max retry after duration of 1 min
	client.MaxSecondaryRateLimitRetryAfterDuration = 60 * time.Second

	// x-ratelimit-reset value of 1h into the future, to make sure we are way over the max wait time duration.
	blockUntil := time.Now().Add(1 * time.Hour).Unix()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set(HeaderRateReset, strconv.Itoa(int(blockUntil)))
		w.Header().Set(HeaderRateRemaining, "1") // set remaining to a value > 0 to distinct from a primary rate limit
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism ...",
   "documentation_url": "https://docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	var abuseRateLimitErr *AbuseRateLimitError
	if !errors.As(err, &abuseRateLimitErr) {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatal("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	// check that the retry after is set to be the max allowed duration
	if got, want := *abuseRateLimitErr.RetryAfter, client.MaxSecondaryRateLimitRetryAfterDuration; got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

// Make network call if client has disabled the rate limit check.
func TestDo_rateLimit_disableRateLimitCheck(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	client.DisableRateLimitCheck = true

	reset := time.Now().UTC().Add(60 * time.Second)
	client.rateLimits[CoreCategory] = Rate{Limit: 5000, Remaining: 0, Reset: Timestamp{reset}}
	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "5000")
		w.Header().Set(HeaderRateRemaining, "5000")
		w.Header().Set(HeaderRateUsed, "0")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Add(time.Hour).Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	})
	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("Response status code = %v, want %v", got, want)
	}
	if got, want := requestCount, 1; got != want {
		t.Errorf("Expected 1 request, got %v", got)
	}
	if got, want := client.rateLimits[CoreCategory].Remaining, 0; got != want {
		t.Errorf("Expected 0 requests remaining, got %v", got)
	}
}

// Make network call if client has bypassed the rate limit check.
func TestDo_rateLimit_bypassRateLimitCheck(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	reset := time.Now().UTC().Add(60 * time.Second)
	client.rateLimits[CoreCategory] = Rate{Limit: 5000, Remaining: 0, Reset: Timestamp{reset}}
	requestCount := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.Header().Set(HeaderRateLimit, "5000")
		w.Header().Set(HeaderRateRemaining, "5000")
		w.Header().Set(HeaderRateUsed, "0")
		w.Header().Set(HeaderRateReset, fmt.Sprint(reset.Add(time.Hour).Unix()))
		w.Header().Set(HeaderRateResource, "core")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	})
	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	resp, err := client.Do(context.WithValue(ctx, BypassRateLimitCheck, true), req, nil)
	if err != nil {
		t.Errorf("Do returned unexpected error: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("Response status code = %v, want %v", got, want)
	}
	if got, want := requestCount, 1; got != want {
		t.Errorf("Expected 1 request, got %v", got)
	}
	if got, want := client.rateLimits[CoreCategory].Remaining, 5000; got != want {
		t.Errorf("Expected 5000 requests remaining, got %v", got)
	}
}

func TestDo_noContent(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := client.Do(ctx, req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestBareDoUntilFound_redirectLoop(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusMovedPermanently)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, _, err := client.bareDoUntilFound(ctx, req, 1)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if !errors.As(err, new(*RedirectionError)) {
		t.Errorf("Expected a Redirection error; got %#v.", err)
	}
}

func TestBareDoUntilFound_UnexpectedRedirection(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusSeeOther)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, _, err := client.bareDoUntilFound(ctx, req, 1)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if !errors.As(err, new(*RedirectionError)) {
		t.Errorf("Expected a Redirection error; got %#v.", err)
	}
}

func TestSanitizeURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"/?a=b", "/?a=b"},
		{"/?a=b&client_secret=secret", "/?a=b&client_secret=REDACTED"},
		{"/?a=b&client_id=id&client_secret=secret", "/?a=b&client_id=id&client_secret=REDACTED"},
	}

	for _, tt := range tests {
		inURL, _ := url.Parse(tt.in)
		want, _ := url.Parse(tt.want)

		if got := sanitizeURL(inURL); !cmp.Equal(got, want) {
			t.Errorf("sanitizeURL(%v) returned %v, want %v", tt.in, got, want)
		}
	}
}

func TestCheckResponse(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: io.NopCloser(strings.NewReader(`{"message":"m",
			"errors": [{"resource": "r", "field": "f", "code": "c"}],
			"block": {"reason": "dmca", "created_at": "2016-03-17T15:39:46Z"}}`)),
	}
	var err *ErrorResponse
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
		Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
		Block: &ErrorBlock{
			Reason:    "dmca",
			CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
		},
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_RateLimit(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusForbidden,
		Header:     http.Header{},
		Body: io.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "url"}`)),
	}
	res.Header.Set(HeaderRateLimit, "60")
	res.Header.Set(HeaderRateRemaining, "0")
	res.Header.Set(HeaderRateUsed, "1")
	res.Header.Set(HeaderRateReset, "243424")
	res.Header.Set(HeaderRateResource, "core")

	var err *RateLimitError
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &RateLimitError{
		Rate:     parseRate(res),
		Response: res,
		Message:  "m",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_AbuseRateLimit(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusForbidden,
		Body: io.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "docs.github.com/en/rest/overview/resources-in-the-rest-api#abuse-rate-limits"}`)),
	}
	var err *AbuseRateLimitError
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &AbuseRateLimitError{
		Response: res,
		Message:  "m",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

// TestCheckResponse_RateLimit_TooManyRequests tests that HTTP 429 with
// X-RateLimit-Remaining: 0 is correctly detected as RateLimitError.
// GitHub API can return either 403 or 429 for rate limiting.
// See: https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
func TestCheckResponse_RateLimit_TooManyRequests(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
		Body: io.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "url"}`)),
	}
	res.Header.Set(HeaderRateLimit, "60")
	res.Header.Set(HeaderRateRemaining, "0")
	res.Header.Set(HeaderRateUsed, "60")
	res.Header.Set(HeaderRateReset, "243424")
	res.Header.Set(HeaderRateResource, "core")

	var err *RateLimitError
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &RateLimitError{
		Rate:     parseRate(res),
		Response: res,
		Message:  "m",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

// TestCheckResponse_AbuseRateLimit_TooManyRequests tests that HTTP 429 with
// secondary rate limit documentation_url is correctly detected as AbuseRateLimitError.
// GitHub API can return either 403 or 429 for secondary rate limits.
// See: https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api#about-secondary-rate-limits
func TestCheckResponse_AbuseRateLimit_TooManyRequests(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
		Body: io.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "https://docs.github.com/rest/overview/rate-limits-for-the-rest-api#about-secondary-rate-limits"}`)),
	}
	res.Header.Set(headerRetryAfter, "60")

	var err *AbuseRateLimitError
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Fatal("Expected error response.")
	}

	if err.Response != res {
		t.Errorf("Response = %v, want %v", err.Response, res)
	}
	if err.Message != "m" {
		t.Errorf("Message = %q, want %q", err.Message, "m")
	}
	if err.RetryAfter == nil {
		t.Error("Expected RetryAfter to be set")
	} else if *err.RetryAfter != 60*time.Second {
		t.Errorf("RetryAfter = %v, want %v", *err.RetryAfter, 60*time.Second)
	}
}

func TestCheckResponse_RedirectionError(t *testing.T) {
	t.Parallel()
	urlStr := "/foo/bar"

	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusFound,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(``)),
	}
	res.Header.Set("Location", urlStr)
	var err *RedirectionError
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	wantedURL, parseErr := url.Parse(urlStr)
	if parseErr != nil {
		t.Errorf("Error parsing fixture url: %v", parseErr)
	}

	want := &RedirectionError{
		Response:   res,
		StatusCode: http.StatusFound,
		Location:   wantedURL,
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCompareHttpResponse(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		h1       *http.Response
		h2       *http.Response
		expected bool
	}{
		"both are nil": {
			expected: true,
		},
		"both are non nil - same StatusCode": {
			expected: true,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 200},
		},
		"both are non nil - different StatusCode": {
			expected: false,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 404},
		},
		"one is nil, other is not": {
			expected: false,
			h2:       &http.Response{},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v := compareHTTPResponse(tc.h1, tc.h2)
			if tc.expected != v {
				t.Errorf("Expected %t, got %t for (%#v, %#v)", tc.expected, v, tc.h1, tc.h2)
			}
		})
	}
}

func TestErrorResponse_Is(t *testing.T) {
	t.Parallel()
	err := &ErrorResponse{
		Response: &http.Response{},
		Message:  "m",
		Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
		Block: &ErrorBlock{
			Reason:    "r",
			CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
		},
		DocumentationURL: "https://github.com",
	}
	testcases := map[string]struct {
		wantSame   bool
		otherError error
	}{
		"errors are same": {
			wantSame: true,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Message": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m1",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - DocumentationURL": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://google.com",
			},
		},
		"errors have different values - Response is nil": {
			wantSame: false,
			otherError: &ErrorResponse{
				Errors:  []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message: "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Errors": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r1", Field: "f1", Code: "c1"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Errors have different length": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Block - one is nil, other is not": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response:         &http.Response{},
				Errors:           []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:          "m",
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Block - different Reason": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r1",
					CreatedAt: &Timestamp{time.Date(2016, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Block - different CreatedAt #1": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: nil,
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different values - Block - different CreatedAt #2": {
			wantSame: false,
			otherError: &ErrorResponse{
				Response: &http.Response{},
				Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
				Message:  "m",
				Block: &ErrorBlock{
					Reason:    "r",
					CreatedAt: &Timestamp{time.Date(2017, time.March, 17, 15, 39, 46, 0, time.UTC)},
				},
				DocumentationURL: "https://github.com",
			},
		},
		"errors have different types": {
			wantSame:   false,
			otherError: errors.New("github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.wantSame != err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", err, tc.otherError)
			}
		})
	}
}

func TestRateLimitError_Is(t *testing.T) {
	t.Parallel()
	err := &RateLimitError{
		Response: &http.Response{},
		Message:  "Github",
	}
	testcases := map[string]struct {
		wantSame   bool
		err        *RateLimitError
		otherError error
	}{
		"errors are same": {
			wantSame: true,
			err:      err,
			otherError: &RateLimitError{
				Response: &http.Response{},
				Message:  "Github",
			},
		},
		"errors are same - Response is nil": {
			wantSame: true,
			err: &RateLimitError{
				Message: "Github",
			},
			otherError: &RateLimitError{
				Message: "Github",
			},
		},
		"errors have different values - Rate": {
			wantSame: false,
			err:      err,
			otherError: &RateLimitError{
				Rate:     Rate{Limit: 10},
				Response: &http.Response{},
				Message:  "Gitlab",
			},
		},
		"errors have different values - Response is nil": {
			wantSame: false,
			err:      err,
			otherError: &RateLimitError{
				Message: "Github",
			},
		},
		"errors have different values - StatusCode": {
			wantSame: false,
			err:      err,
			otherError: &RateLimitError{
				Response: &http.Response{StatusCode: 200},
				Message:  "Github",
			},
		},
		"errors have different types": {
			wantSame:   false,
			err:        err,
			otherError: errors.New("github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.wantSame != tc.err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", tc.err, tc.otherError)
			}
		})
	}
}

func TestAbuseRateLimitError_Is(t *testing.T) {
	t.Parallel()
	t1 := 1 * time.Second
	t2 := 2 * time.Second
	err := &AbuseRateLimitError{
		Response:   &http.Response{},
		Message:    "Github",
		RetryAfter: &t1,
	}
	testcases := map[string]struct {
		wantSame   bool
		err        *AbuseRateLimitError
		otherError error
	}{
		"errors are same": {
			wantSame: true,
			err:      err,
			otherError: &AbuseRateLimitError{
				Response:   &http.Response{},
				Message:    "Github",
				RetryAfter: &t1,
			},
		},
		"errors are same - Response is nil": {
			wantSame: true,
			err: &AbuseRateLimitError{
				Message:    "Github",
				RetryAfter: &t1,
			},
			otherError: &AbuseRateLimitError{
				Message:    "Github",
				RetryAfter: &t1,
			},
		},
		"errors have different values - Message": {
			wantSame: false,
			err:      err,
			otherError: &AbuseRateLimitError{
				Response:   &http.Response{},
				Message:    "Gitlab",
				RetryAfter: nil,
			},
		},
		"errors have different values - RetryAfter": {
			wantSame: false,
			err:      err,
			otherError: &AbuseRateLimitError{
				Response:   &http.Response{},
				Message:    "Github",
				RetryAfter: &t2,
			},
		},
		"errors have different values - Response is nil": {
			wantSame: false,
			err:      err,
			otherError: &AbuseRateLimitError{
				Message:    "Github",
				RetryAfter: &t1,
			},
		},
		"errors have different values - StatusCode": {
			wantSame: false,
			err:      err,
			otherError: &AbuseRateLimitError{
				Response:   &http.Response{StatusCode: 200},
				Message:    "Github",
				RetryAfter: &t1,
			},
		},
		"errors have different types": {
			wantSame:   false,
			err:        err,
			otherError: errors.New("github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.wantSame != tc.err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", tc.err, tc.otherError)
			}
		})
	}
}

func TestAcceptedError_Is(t *testing.T) {
	t.Parallel()
	err := &AcceptedError{Raw: []byte("Github")}
	testcases := map[string]struct {
		wantSame   bool
		otherError error
	}{
		"errors are same": {
			wantSame:   true,
			otherError: &AcceptedError{Raw: []byte("Github")},
		},
		"errors have different values": {
			wantSame:   false,
			otherError: &AcceptedError{Raw: []byte("Gitlab")},
		},
		"errors have different types": {
			wantSame:   false,
			otherError: errors.New("github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.wantSame != err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", err, tc.otherError)
			}
		})
	}
}

// Ensure that we properly handle API errors that do not contain a response body.
func TestCheckResponse_noBody(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	var err *ErrorResponse
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_unexpectedErrorStructure(t *testing.T) {
	t.Parallel()
	httpBody := `{"message":"m", "errors": ["error 1"]}`
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(httpBody)),
	}
	var err *ErrorResponse
	errors.As(CheckResponse(res), &err)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
		Errors:   []Error{{Message: "error 1"}},
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
	data, err2 := io.ReadAll(err.Response.Body)
	if err2 != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got := string(data); got != httpBody {
		t.Errorf("ErrorResponse.Response.Body = %q, want %q", got, httpBody)
	}
}

func TestParseBooleanResponse_true(t *testing.T) {
	t.Parallel()
	result, err := parseBoolResponse(nil)
	if err != nil {
		t.Errorf("parseBoolResponse returned error: %+v", err)
	}

	if want := true; result != want {
		t.Errorf("parseBoolResponse returned %+v, want: %+v", result, want)
	}
}

func TestParseBooleanResponse_false(t *testing.T) {
	t.Parallel()
	v := &ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	result, err := parseBoolResponse(v)
	if err != nil {
		t.Errorf("parseBoolResponse returned error: %+v", err)
	}

	if want := false; result != want {
		t.Errorf("parseBoolResponse returned %+v, want: %+v", result, want)
	}
}

func TestParseBooleanResponse_error(t *testing.T) {
	t.Parallel()
	v := &ErrorResponse{Response: &http.Response{StatusCode: http.StatusBadRequest}}
	result, err := parseBoolResponse(v)

	if err == nil {
		t.Error("Expected error to be returned.")
	}

	if want := false; result != want {
		t.Errorf("parseBoolResponse returned %+v, want: %+v", result, want)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	t.Parallel()
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Error("Expected non-empty ErrorResponse.Error()")
	}

	// dont panic if request is nil
	res = &http.Response{}
	err = ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Error("Expected non-empty ErrorResponse.Error()")
	}

	// dont panic if response is nil
	err = ErrorResponse{Message: "m"}
	if err.Error() == "" {
		t.Error("Expected non-empty ErrorResponse.Error()")
	}
}

func TestError_Error(t *testing.T) {
	t.Parallel()
	err := Error{}
	if err.Error() == "" {
		t.Error("Expected non-empty Error.Error()")
	}
}

func TestSetCredentialsAsHeaders(t *testing.T) {
	t.Parallel()
	req := new(http.Request)
	id, secret := "id", "secret"
	modifiedRequest := setCredentialsAsHeaders(req, id, secret)

	actualID, actualSecret, ok := modifiedRequest.BasicAuth()
	if !ok {
		t.Error("request does not contain basic credentials")
	}

	if actualID != id {
		t.Errorf("id is %v, want %v", actualID, id)
	}

	if actualSecret != secret {
		t.Errorf("secret is %v, want %v", actualSecret, secret)
	}
}

func TestUnauthenticatedRateLimitedTransport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	clientID, clientSecret := "id", "secret"
	mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		id, secret, ok := r.BasicAuth()
		if !ok {
			t.Error("request does not contain basic auth credentials")
		}
		if id != clientID {
			t.Errorf("request contained basic auth username %q, want %q", id, clientID)
		}
		if secret != clientSecret {
			t.Errorf("request contained basic auth password %q, want %q", secret, clientSecret)
		}
	})

	tp := &UnauthenticatedRateLimitedTransport{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	unauthedClient := NewClient(tp.Client())
	unauthedClient.BaseURL = client.BaseURL
	req, _ := unauthedClient.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := unauthedClient.Do(ctx, req, nil)
	assertNilError(t, err)
}

func TestUnauthenticatedRateLimitedTransport_missingFields(t *testing.T) {
	t.Parallel()
	// missing ClientID
	tp := &UnauthenticatedRateLimitedTransport{
		ClientSecret: "secret",
	}
	_, err := tp.RoundTrip(nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}

	// missing ClientSecret
	tp = &UnauthenticatedRateLimitedTransport{
		ClientID: "id",
	}
	_, err = tp.RoundTrip(nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestUnauthenticatedRateLimitedTransport_transport(t *testing.T) {
	t.Parallel()
	// default transport
	tp := &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	if tp.transport() != http.DefaultTransport {
		t.Error("Expected http.DefaultTransport to be used.")
	}

	// custom transport
	tp = &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
		Transport:    &http.Transport{},
	}
	if tp.transport() == http.DefaultTransport {
		t.Error("Expected custom transport to be used.")
	}
}

func TestBasicAuthTransport(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	username, password, otp := "u", "p", "123456"

	mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			t.Error("request does not contain basic auth credentials")
		}
		if u != username {
			t.Errorf("request contained basic auth username %q, want %q", u, username)
		}
		if p != password {
			t.Errorf("request contained basic auth password %q, want %q", p, password)
		}
		if got, want := r.Header.Get(headerOTP), otp; got != want {
			t.Errorf("request contained OTP %q, want %q", got, want)
		}
	})

	tp := &BasicAuthTransport{
		Username: username,
		Password: password,
		OTP:      otp,
	}
	basicAuthClient := NewClient(tp.Client())
	basicAuthClient.BaseURL = client.BaseURL
	req, _ := basicAuthClient.NewRequest("GET", ".", nil)
	ctx := t.Context()
	_, err := basicAuthClient.Do(ctx, req, nil)
	assertNilError(t, err)
}

func TestBasicAuthTransport_transport(t *testing.T) {
	t.Parallel()
	// default transport
	tp := &BasicAuthTransport{}
	if tp.transport() != http.DefaultTransport {
		t.Error("Expected http.DefaultTransport to be used.")
	}

	// custom transport
	tp = &BasicAuthTransport{
		Transport: &http.Transport{},
	}
	if tp.transport() == http.DefaultTransport {
		t.Error("Expected custom transport to be used.")
	}
}

func TestFormatRateReset(t *testing.T) {
	t.Parallel()
	d := 120*time.Minute + 12*time.Second
	got := formatRateReset(d)
	want := "[rate reset in 120m12s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 14*time.Minute + 2*time.Second
	got = formatRateReset(d)
	want = "[rate reset in 14m02s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 2*time.Minute + 2*time.Second
	got = formatRateReset(d)
	want = "[rate reset in 2m02s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 12 * time.Second
	got = formatRateReset(d)
	want = "[rate reset in 12s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = -1 * (2*time.Hour + 2*time.Second)
	got = formatRateReset(d)
	want = "[rate limit was reset 120m02s ago]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}
}

func TestNestedStructAccessorNoPanic(t *testing.T) {
	t.Parallel()
	issue := &Issue{User: nil}
	got := issue.GetUser().GetPlan().GetName()
	want := ""
	if got != want {
		t.Errorf("Issues.Get.GetUser().GetPlan().GetName() returned %+v, want %+v", got, want)
	}
}

func TestTwoFactorAuthError(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	e := &TwoFactorAuthError{
		Response: &http.Response{
			Request:    &http.Request{Method: "PUT", URL: u},
			StatusCode: http.StatusTooManyRequests,
		},
		Message: "<msg>",
	}
	if got, want := e.Error(), "PUT https://example.com: 429 <msg> []"; got != want {
		t.Errorf("TwoFactorAuthError = %q, want %q", got, want)
	}
}

func TestRateLimitError(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	r := &RateLimitError{
		Response: &http.Response{
			Request:    &http.Request{Method: "PUT", URL: u},
			StatusCode: http.StatusTooManyRequests,
		},
		Message: "<msg>",
	}
	if got, want := r.Error(), "PUT https://example.com: 429 <msg> [rate limit was reset"; !strings.Contains(got, want) {
		t.Errorf("RateLimitError = %q, want %q", got, want)
	}
}

func TestAcceptedError(t *testing.T) {
	t.Parallel()
	a := &AcceptedError{}
	if got, want := a.Error(), "try again later"; !strings.Contains(got, want) {
		t.Errorf("AcceptedError = %q, want %q", got, want)
	}
}

func TestAbuseRateLimitError(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	r := &AbuseRateLimitError{
		Response: &http.Response{
			Request:    &http.Request{Method: "PUT", URL: u},
			StatusCode: http.StatusTooManyRequests,
		},
		Message: "<msg>",
	}
	if got, want := r.Error(), "PUT https://example.com: 429 <msg>"; got != want {
		t.Errorf("AbuseRateLimitError = %q, want %q", got, want)
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, expectedBody)
	})

	ctx := t.Context()
	req, err := client.NewRequest("GET", "test-url", nil)
	if err != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %v", err)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}
	if string(got) != expectedBody {
		t.Fatalf("Expected %q, got %q", expectedBody, string(got))
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("resp.Body.Close() returned error: %v", err)
	}
}

func TestErrorResponse_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ErrorResponse{}, "{}")

	u := &ErrorResponse{
		Message: "msg",
		Errors: []Error{
			{
				Resource: "res",
				Field:    "f",
				Code:     "c",
				Message:  "msg",
			},
		},
		Block: &ErrorBlock{
			Reason:    "reason",
			CreatedAt: &Timestamp{referenceTime},
		},
		DocumentationURL: "doc",
	}

	want := `{
		"message": "msg",
		"errors": [
			{
				"resource": "res",
				"field": "f",
				"code": "c",
				"message": "msg"
			}
		],
		"block": {
			"reason": "reason",
			"created_at": ` + referenceTimeStr + `
		},
		"documentation_url": "doc"
	}`

	testJSONMarshal(t, u, want)
}

func TestErrorBlock_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ErrorBlock{}, "{}")

	u := &ErrorBlock{
		Reason:    "reason",
		CreatedAt: &Timestamp{referenceTime},
	}

	want := `{
		"reason": "reason",
		"created_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestRateLimitError_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RateLimitError{}, "{}")

	u := &RateLimitError{
		Rate: Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		Message: "msg",
	}

	want := `{
		"Rate": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"message": "msg"
	}`

	testJSONMarshal(t, u, want)
}

func TestAbuseRateLimitError_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &AbuseRateLimitError{}, "{}")

	u := &AbuseRateLimitError{
		Message: "msg",
	}

	want := `{
		"message": "msg"
	}`

	testJSONMarshal(t, u, want)
}

func TestError_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Error{}, "{}")

	u := &Error{
		Resource: "res",
		Field:    "field",
		Code:     "code",
		Message:  "msg",
	}

	want := `{
		"resource": "res",
		"field": "field",
		"code": "code",
		"message": "msg"
	}`

	testJSONMarshal(t, u, want)
}

func TestParseTokenExpiration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		header string
		want   Timestamp
	}{
		{
			header: "",
			want:   Timestamp{},
		},
		{
			header: "this is a garbage",
			want:   Timestamp{},
		},
		{
			header: "2021-09-03 02:34:04 UTC",
			want:   Timestamp{time.Date(2021, time.September, 3, 2, 34, 4, 0, time.UTC)},
		},
		{
			header: "2021-09-03 14:34:04 UTC",
			want:   Timestamp{time.Date(2021, time.September, 3, 14, 34, 4, 0, time.UTC)},
		},
		// Some tokens include the timezone offset instead of the timezone.
		// https://github.com/google/go-github/issues/2649
		{
			header: "2023-04-26 20:23:26 +0200",
			want:   Timestamp{time.Date(2023, time.April, 26, 18, 23, 26, 0, time.UTC)},
		},
	}

	for _, tt := range tests {
		res := &http.Response{
			Request: &http.Request{},
			Header:  http.Header{},
		}

		res.Header.Set(headerTokenExpiration, tt.header)
		exp := parseTokenExpiration(res)
		if !exp.Equal(tt.want) {
			t.Errorf("parseTokenExpiration of %q\nreturned %#v\n    want %#v", tt.header, exp, tt.want)
		}
	}
}

func TestClientCopy_leak_transport(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		accessToken := r.Header.Get("Authorization")
		_, _ = fmt.Fprintf(w, `{"login": "%v"}`, accessToken)
	}))
	clientPreconfiguredWithURLs, err := NewClient(nil).WithEnterpriseURLs(srv.URL, srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	aliceClient := clientPreconfiguredWithURLs.WithAuthToken("alice")
	bobClient := clientPreconfiguredWithURLs.WithAuthToken("bob")

	alice, _, err := aliceClient.Users.Get(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}

	assertNoDiff(t, "Bearer alice", alice.GetLogin())

	bob, _, err := bobClient.Users.Get(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}

	assertNoDiff(t, "Bearer bob", bob.GetLogin())
}

func TestPtr(t *testing.T) {
	t.Parallel()
	equal := func(t *testing.T, want, got any) {
		t.Helper()
		if !cmp.Equal(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}

	equal(t, true, *Ptr(true))
	equal(t, int(10), *Ptr(int(10)))
	equal(t, int64(-10), *Ptr(int64(-10)))
	equal(t, "str", *Ptr("str"))
}

func TestDeploymentProtectionRuleEvent_GetRunID(t *testing.T) {
	t.Parallel()

	var want int64 = 123456789
	url := "https://api.github.com/repos/dummy-org/dummy-repo/actions/runs/123456789/deployment_protection_rule"

	e := DeploymentProtectionRuleEvent{
		DeploymentCallbackURL: &url,
	}

	got, _ := e.GetRunID()
	if got != want {
		t.Errorf("want %#v, got %#v", want, got)
	}

	want = 123456789
	url = "repos/dummy-org/dummy-repo/actions/runs/123456789/deployment_protection_rule"

	e = DeploymentProtectionRuleEvent{
		DeploymentCallbackURL: &url,
	}

	got, _ = e.GetRunID()
	if got != want {
		t.Errorf("want %#v, got %#v", want, got)
	}

	want = -1
	url = "https://api.github.com/repos/dummy-org/dummy-repo/actions/runs/abc123/deployment_protection_rule"
	got, err := e.GetRunID()
	if err == nil {
		t.Error("Expected error to be returned")
	}

	if got != want {
		t.Errorf("want %#v, got %#v", want, got)
	}
}
