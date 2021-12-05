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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"reflect"
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
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
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

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.URL, server.Close
}

// openTestFile creates a new file with the given name and content for testing.
// In order to ensure the exact file name, this function will create a new temp
// directory, and create the file in that directory. It is the caller's
// responsibility to remove the directory and its contents when no longer needed.
func openTestFile(name, content string) (file *os.File, dir string, err error) {
	dir, err = ioutil.TempDir("", "go-github")
	if err != nil {
		return nil, dir, err
	}

	file, err = os.OpenFile(path.Join(dir, name), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, dir, err
	}

	fmt.Fprint(file, content)

	// close and re-open the file to keep file.Stat() happy
	file.Close()
	file, err = os.Open(file.Name())
	if err != nil {
		return nil, dir, err
	}

	return file, dir, err
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

	r.ParseForm()
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, want %s", v, j, w)
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
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	resp, err = f()
	if bypass := resp.Request.Context().Value(bypassRateLimitCheck); bypass != nil {
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(code)
	})

	ctx := context.Background()
	_, _, err := client.Repositories.ListHooks(ctx, "o", "r", nil)

	switch e := err.(type) {
	case *ErrorResponse:
	case *RateLimitError:
	case *AbuseRateLimitError:
		if code != e.Response.StatusCode {
			t.Error("Error response does not contain status code")
		}
	default:
		t.Error("Unknown error response type")
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestClient(t *testing.T) {
	c := NewClient(nil)
	c2 := c.Client()
	if c.client == c2 {
		t.Error("Client returned same http.Client, but should be different")
	}
}

func TestNewEnterpriseClient(t *testing.T) {
	baseURL := "https://custom-url/api/v3/"
	uploadURL := "https://custom-upload-url/api/uploads/"
	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), baseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), uploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_addsTrailingSlashToURLs(t *testing.T) {
	baseURL := "https://custom-url/api/v3"
	uploadURL := "https://custom-upload-url/api/uploads"
	formattedBaseURL := baseURL + "/"
	formattedUploadURL := uploadURL + "/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_addsEnterpriseSuffixToURLs(t *testing.T) {
	baseURL := "https://custom-url/"
	uploadURL := "https://custom-upload-url/"
	formattedBaseURL := baseURL + "api/v3/"
	formattedUploadURL := uploadURL + "api/uploads/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_addsEnterpriseSuffixAndTrailingSlashToURLs(t *testing.T) {
	baseURL := "https://custom-url"
	uploadURL := "https://custom-upload-url"
	formattedBaseURL := baseURL + "/api/v3/"
	formattedUploadURL := uploadURL + "/api/uploads/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_badBaseURL(t *testing.T) {
	baseURL := "bogus\nbase\nURL"
	uploadURL := "https://custom-upload-url/api/uploads/"
	if _, err := NewEnterpriseClient(baseURL, uploadURL, nil); err == nil {
		t.Fatal("NewEnterpriseClient returned nil, expected error")
	}
}

func TestNewEnterpriseClient_badUploadURL(t *testing.T) {
	baseURL := "https://custom-url/api/v3/"
	uploadURL := "bogus\nupload\nURL"
	if _, err := NewEnterpriseClient(baseURL, uploadURL, nil); err == nil {
		t.Fatal("NewEnterpriseClient returned nil, expected error")
	}
}

func TestNewEnterpriseClient_URLHasExistingAPIPrefix_AddTrailingSlash(t *testing.T) {
	baseURL := "https://api.custom-url"
	uploadURL := "https://api.custom-upload-url"
	formattedBaseURL := baseURL + "/"
	formattedUploadURL := uploadURL + "/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_URLHasExistingAPIPrefixAndTrailingSlash(t *testing.T) {
	baseURL := "https://api.custom-url/"
	uploadURL := "https://api.custom-upload-url/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), baseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), uploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_URLHasAPISubdomain_AddTrailingSlash(t *testing.T) {
	baseURL := "https://catalog.api.custom-url"
	uploadURL := "https://catalog.api.custom-upload-url"
	formattedBaseURL := baseURL + "/"
	formattedUploadURL := uploadURL + "/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_URLHasAPISubdomainAndTrailingSlash(t *testing.T) {
	baseURL := "https://catalog.api.custom-url/"
	uploadURL := "https://catalog.api.custom-upload-url/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), baseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), uploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_URLIsNotAProperAPISubdomain_addsEnterpriseSuffixAndSlash(t *testing.T) {
	baseURL := "https://cloud-api.custom-url"
	uploadURL := "https://cloud-api.custom-upload-url"
	formattedBaseURL := baseURL + "/api/v3/"
	formattedUploadURL := uploadURL + "/api/uploads/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

func TestNewEnterpriseClient_URLIsNotAProperAPISubdomain_addsEnterpriseSuffix(t *testing.T) {
	baseURL := "https://cloud-api.custom-url/"
	uploadURL := "https://cloud-api.custom-upload-url/"
	formattedBaseURL := baseURL + "api/v3/"
	formattedUploadURL := uploadURL + "api/uploads/"

	c, err := NewEnterpriseClient(baseURL, uploadURL, nil)
	if err != nil {
		t.Fatalf("NewEnterpriseClient returned unexpected error: %v", err)
	}

	if got, want := c.BaseURL.String(), formattedBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UploadURL.String(), formattedUploadURL; got != want {
		t.Errorf("NewClient UploadURL is %v, want %v", got, want)
	}
}

// Ensure that length of Client.rateLimits is the same as number of fields in RateLimits struct.
func TestClient_rateLimits(t *testing.T) {
	if got, want := len(Client{}.rateLimits), reflect.TypeOf(RateLimits{}).NumField(); got != want {
		t.Errorf("len(Client{}.rateLimits) is %v, want %v", got, want)
	}
}

func TestRateLimits_String(t *testing.T) {
	v := RateLimits{
		Core:   &Rate{},
		Search: &Rate{},
	}
	want := `github.RateLimits{Core:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}, Search:github.Rate{Limit:0, Remaining:0, Reset:github.Timestamp{0001-01-01 00:00:00 +0000 UTC}}}`
	if got := v.String(); got != want {
		t.Errorf("RateLimits.String = %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &User{Login: String("l")}, `{"login":"l"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil)
	if _, err := c.NewRequest("BOGUS\nMETHOD", ".", nil); err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}
}

// ensure that no User-Agent header is set if the client's UserAgent is empty.
// This caused a problem with Google's internal http client.
func TestNewRequest_emptyUserAgent(t *testing.T) {
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
	c := NewClient(nil)
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
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
		if _, err := c.NewRequest(http.MethodGet, "test", nil); test.wantError && err == nil {
			t.Fatalf("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewRequest returned unexpected error: %v.", err)
		}
	}
}

func TestNewUploadRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewUploadRequest(":", nil, 0, "")
	testURLParseError(t, err)
}

func TestNewUploadRequest_errorForNoTrailingSlash(t *testing.T) {
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
			t.Fatalf("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewUploadRequest returned unexpected error: %v.", err)
		}
	}
}

func TestResponse_populatePageValues(t *testing.T) {
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?page=1>; rel="first",` +
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
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?since=1>; rel="first",` +
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
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?since=2021-12-04T10%3A43%3A42Z&page=1>; rel="first",` +
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
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?after=a1b2c3&before=>; rel="next",` +
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
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?page=1>,` +
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
	r := http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/?since=1>,` +
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
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	ctx := context.Background()
	client.Do(ctx, req, body)

	want := &foo{"a"}
	if !cmp.Equal(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_nilContext(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(nil, req, nil)

	if !errors.Is(err, errNonNilContext) {
		t.Errorf("Expected context must be non-nil error")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

// Test handling of an error caused by the internal http client's Do()
// function. A redirect loop is pretty unlikely to occur within the GitHub
// API, but does allow us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusFound)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

// Test that an error caused by the internal http client's Do() function
// does not leak the client secret.
func TestDo_sanitizeURL(t *testing.T) {
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
	ctx := context.Background()
	_, err = unauthedClient.Do(ctx, req, nil)
	if err == nil {
		t.Fatal("Expected error to be returned.")
	}
	if strings.Contains(err.Error(), "client_secret=secret") {
		t.Errorf("Do error contains secret, should be redacted:\n%q", err)
	}
}

func TestDo_rateLimit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "60")
		w.Header().Set(headerRateRemaining, "59")
		w.Header().Set(headerRateReset, "1372700873")
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
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
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if resp.Rate.Reset.UTC() != reset {
		t.Errorf("Client rate reset = %v, want %v", resp.Rate.Reset, reset)
	}
}

// ensure rate limit is still parsed, even for error responses
func TestDo_rateLimit_errorResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "60")
		w.Header().Set(headerRateRemaining, "59")
		w.Header().Set(headerRateReset, "1372700873")
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if _, ok := err.(*RateLimitError); ok {
		t.Errorf("Did not expect a *RateLimitError error; got %#v.", err)
	}
	if got, want := resp.Rate.Limit, 60; got != want {
		t.Errorf("Client rate limit = %v, want %v", got, want)
	}
	if got, want := resp.Rate.Remaining, 59; got != want {
		t.Errorf("Client rate remaining = %v, want %v", got, want)
	}
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if resp.Rate.Reset.UTC() != reset {
		t.Errorf("Client rate reset = %v, want %v", resp.Rate.Reset, reset)
	}
}

// Ensure *RateLimitError is returned when API rate limit is exceeded.
func TestDo_rateLimit_rateLimitError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "60")
		w.Header().Set(headerRateRemaining, "0")
		w.Header().Set(headerRateReset, "1372700873")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("Expected a *RateLimitError error; got %#v.", err)
	}
	if got, want := rateLimitErr.Rate.Limit, 60; got != want {
		t.Errorf("rateLimitErr rate limit = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Remaining, 0; got != want {
		t.Errorf("rateLimitErr rate remaining = %v, want %v", got, want)
	}
	reset := time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC)
	if rateLimitErr.Rate.Reset.UTC() != reset {
		t.Errorf("rateLimitErr rate reset = %v, want %v", rateLimitErr.Rate.Reset.UTC(), reset)
	}
}

// Ensure a network call is not made when it's known that API rate limit is still exceeded.
func TestDo_rateLimit_noNetworkCall(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reset := time.Now().UTC().Add(time.Minute).Round(time.Second) // Rate reset is a minute from now, with 1 second precision.

	mux.HandleFunc("/first", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "60")
		w.Header().Set(headerRateRemaining, "0")
		w.Header().Set(headerRateReset, fmt.Sprint(reset.Unix()))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	madeNetworkCall := false
	mux.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
		madeNetworkCall = true
	})

	// First request is made, and it makes the client aware of rate reset time being in the future.
	req, _ := client.NewRequest("GET", "first", nil)
	ctx := context.Background()
	client.Do(ctx, req, nil)

	// Second request should not cause a network call to be made, since client can predict a rate limit error.
	req, _ = client.NewRequest("GET", "second", nil)
	_, err := client.Do(ctx, req, nil)

	if madeNetworkCall {
		t.Fatal("Network call was made, even though rate limit is known to still be exceeded.")
	}

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("Expected a *RateLimitError error; got %#v.", err)
	}
	if got, want := rateLimitErr.Rate.Limit, 60; got != want {
		t.Errorf("rateLimitErr rate limit = %v, want %v", got, want)
	}
	if got, want := rateLimitErr.Rate.Remaining, 0; got != want {
		t.Errorf("rateLimitErr rate remaining = %v, want %v", got, want)
	}
	if rateLimitErr.Rate.Reset.UTC() != reset {
		t.Errorf("rateLimitErr rate reset = %v, want %v", rateLimitErr.Rate.Reset.UTC(), reset)
	}
}

// Ensure *AbuseRateLimitError is returned when the response indicates that
// the client has triggered an abuse detection mechanism.
func TestDo_rateLimit_abuseRateLimitError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		// When the abuse rate limit error is of the "temporarily blocked from content creation" type,
		// there is no "Retry-After" header.
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism and have been temporarily blocked from content creation. Please retry your request again later.",
   "documentation_url": "https://docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	abuseRateLimitErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if got, want := abuseRateLimitErr.RetryAfter, (*time.Duration)(nil); got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

// Ensure *AbuseRateLimitError is returned when the response indicates that
// the client has triggered an abuse detection mechanism on GitHub Enterprise.
func TestDo_rateLimit_abuseRateLimitErrorEnterprise(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		// When the abuse rate limit error is of the "temporarily blocked from content creation" type,
		// there is no "Retry-After" header.
		// This response returns a documentation url like the one returned for GitHub Enterprise, this
		// url changes between versions but follows roughly the same format.
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism and have been temporarily blocked from content creation. Please retry your request again later.",
   "documentation_url": "https://docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	abuseRateLimitErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if got, want := abuseRateLimitErr.RetryAfter, (*time.Duration)(nil); got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

// Ensure *AbuseRateLimitError.RetryAfter is parsed correctly.
func TestDo_rateLimit_abuseRateLimitError_retryAfter(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Retry-After", "123") // Retry after value of 123 seconds.
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{
   "message": "You have triggered an abuse detection mechanism ...",
   "documentation_url": "https://docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"
}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	abuseRateLimitErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected a *AbuseRateLimitError error; got %#v.", err)
	}
	if abuseRateLimitErr.RetryAfter == nil {
		t.Fatalf("abuseRateLimitErr RetryAfter is nil, expected not-nil")
	}
	if got, want := *abuseRateLimitErr.RetryAfter, 123*time.Second; got != want {
		t.Errorf("abuseRateLimitErr RetryAfter = %v, want %v", got, want)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestSanitizeURL(t *testing.T) {
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
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"errors": [{"resource": "r", "field": "f", "code": "c"}],
			"block": {"reason": "dmca", "created_at": "2016-03-17T15:39:46Z"}}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
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
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusForbidden,
		Header:     http.Header{},
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "url"}`)),
	}
	res.Header.Set(headerRateLimit, "60")
	res.Header.Set(headerRateRemaining, "0")
	res.Header.Set(headerRateReset, "243424")

	err := CheckResponse(res).(*RateLimitError)

	if err == nil {
		t.Errorf("Expected error response.")
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
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusForbidden,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"documentation_url": "docs.github.com/en/free-pro-team@latest/rest/overview/resources-in-the-rest-api#abuse-rate-limits"}`)),
	}
	err := CheckResponse(res).(*AbuseRateLimitError)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &AbuseRateLimitError{
		Response: res,
		Message:  "m",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCompareHttpResponse(t *testing.T) {
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
			v := compareHTTPResponse(tc.h1, tc.h2)
			if tc.expected != v {
				t.Errorf("Expected %t, got %t for (%#v, %#v)", tc.expected, v, tc.h1, tc.h2)
			}
		})
	}
}

func TestErrorResponse_Is(t *testing.T) {
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
			otherError: errors.New("Github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.wantSame != err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", err, tc.otherError)
			}
		})
	}
}

func TestRateLimitError_Is(t *testing.T) {
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
			otherError: errors.New("Github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.wantSame != tc.err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", tc.err, tc.otherError)
			}
		})
	}
}

func TestAbuseRateLimitError_Is(t *testing.T) {
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
			otherError: errors.New("Github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.wantSame != tc.err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", tc.err, tc.otherError)
			}
		})
	}
}

func TestAcceptedError_Is(t *testing.T) {
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
			otherError: errors.New("Github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.wantSame != err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", err, tc.otherError)
			}
		})
	}
}

// ensure that we properly handle API errors that do not contain a response body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_unexpectedErrorStructure(t *testing.T) {
	httpBody := `{"message":"m", "errors": ["error 1"]}`
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(httpBody)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
		Errors:   []Error{{Message: "error 1"}},
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
	data, err2 := ioutil.ReadAll(err.Response.Body)
	if err2 != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got := string(data); got != httpBody {
		t.Errorf("ErrorResponse.Response.Body = %q, want %q", got, httpBody)
	}
}

func TestParseBooleanResponse_true(t *testing.T) {
	result, err := parseBoolResponse(nil)
	if err != nil {
		t.Errorf("parseBoolResponse returned error: %+v", err)
	}

	if want := true; result != want {
		t.Errorf("parseBoolResponse returned %+v, want: %+v", result, want)
	}
}

func TestParseBooleanResponse_false(t *testing.T) {
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
	v := &ErrorResponse{Response: &http.Response{StatusCode: http.StatusBadRequest}}
	result, err := parseBoolResponse(v)

	if err == nil {
		t.Errorf("Expected error to be returned.")
	}

	if want := false; result != want {
		t.Errorf("parseBoolResponse returned %+v, want: %+v", result, want)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func TestError_Error(t *testing.T) {
	err := Error{}
	if err.Error() == "" {
		t.Errorf("Expected non-empty Error.Error()")
	}
}

func TestRateLimits(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874}
		}}`)
	})

	ctx := context.Background()
	rate, _, err := client.RateLimits(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}

	if got, want := client.rateLimits[coreCategory], *want.Core; got != want {
		t.Errorf("client.rateLimits[coreCategory] is %+v, want %+v", got, want)
	}
	if got, want := client.rateLimits[searchCategory], *want.Search; got != want {
		t.Errorf("client.rateLimits[searchCategory] is %+v, want %+v", got, want)
	}
}

func TestRateLimits_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "RateLimits"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.RateLimits(ctx)
		return resp, err
	})
}

func TestRateLimits_overQuota(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	client.rateLimits[coreCategory] = Rate{
		Limit:     1,
		Remaining: 0,
		Reset:     Timestamp{time.Now().Add(time.Hour).Local()},
	}
	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874}
		}}`)
	})

	ctx := context.Background()
	rate, _, err := client.RateLimits(ctx)
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Reset:     Timestamp{time.Date(2013, time.July, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
	}
	if !cmp.Equal(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
	}

	if got, want := client.rateLimits[coreCategory], *want.Core; got != want {
		t.Errorf("client.rateLimits[coreCategory] is %+v, want %+v", got, want)
	}
	if got, want := client.rateLimits[searchCategory], *want.Search; got != want {
		t.Errorf("client.rateLimits[searchCategory] is %+v, want %+v", got, want)
	}
}

func TestSetCredentialsAsHeaders(t *testing.T) {
	req := new(http.Request)
	id, secret := "id", "secret"
	modifiedRequest := setCredentialsAsHeaders(req, id, secret)

	actualID, actualSecret, ok := modifiedRequest.BasicAuth()
	if !ok {
		t.Errorf("request does not contain basic credentials")
	}

	if actualID != id {
		t.Errorf("id is %s, want %s", actualID, id)
	}

	if actualSecret != secret {
		t.Errorf("secret is %s, want %s", actualSecret, secret)
	}
}

func TestUnauthenticatedRateLimitedTransport(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	clientID, clientSecret := "id", "secret"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id, secret, ok := r.BasicAuth()
		if !ok {
			t.Errorf("request does not contain basic auth credentials")
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
	ctx := context.Background()
	unauthedClient.Do(ctx, req, nil)
}

func TestUnauthenticatedRateLimitedTransport_missingFields(t *testing.T) {
	// missing ClientID
	tp := &UnauthenticatedRateLimitedTransport{
		ClientSecret: "secret",
	}
	_, err := tp.RoundTrip(nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	// missing ClientSecret
	tp = &UnauthenticatedRateLimitedTransport{
		ClientID: "id",
	}
	_, err = tp.RoundTrip(nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestUnauthenticatedRateLimitedTransport_transport(t *testing.T) {
	// default transport
	tp := &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	if tp.transport() != http.DefaultTransport {
		t.Errorf("Expected http.DefaultTransport to be used.")
	}

	// custom transport
	tp = &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
		Transport:    &http.Transport{},
	}
	if tp.transport() == http.DefaultTransport {
		t.Errorf("Expected custom transport to be used.")
	}
}

func TestBasicAuthTransport(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	username, password, otp := "u", "p", "123456"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			t.Errorf("request does not contain basic auth credentials")
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
	ctx := context.Background()
	basicAuthClient.Do(ctx, req, nil)
}

func TestBasicAuthTransport_transport(t *testing.T) {
	// default transport
	tp := &BasicAuthTransport{}
	if tp.transport() != http.DefaultTransport {
		t.Errorf("Expected http.DefaultTransport to be used.")
	}

	// custom transport
	tp = &BasicAuthTransport{
		Transport: &http.Transport{},
	}
	if tp.transport() == http.DefaultTransport {
		t.Errorf("Expected custom transport to be used.")
	}
}

func TestFormatRateReset(t *testing.T) {
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
	issue := &Issue{User: nil}
	got := issue.GetUser().GetPlan().GetName()
	want := ""
	if got != want {
		t.Errorf("Issues.Get.GetUser().GetPlan().GetName() returned %+v, want %+v", got, want)
	}
}

func TestTwoFactorAuthError(t *testing.T) {
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
	a := &AcceptedError{}
	if got, want := a.Error(), "try again later"; !strings.Contains(got, want) {
		t.Errorf("AcceptedError = %q, want %q", got, want)
	}
}

func TestAbuseRateLimitError(t *testing.T) {
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

func TestAddOptions_QueryValues(t *testing.T) {
	if _, err := addOptions("yo", ""); err == nil {
		t.Error("addOptions err = nil, want error")
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest("GET", "test-url", nil)
	if err != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %v", err)
	}

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll returned error: %v", err)
	}
	if string(got) != expectedBody {
		t.Fatalf("Expected %q, got %q", expectedBody, string(got))
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("resp.Body.Close() returned error: %v", err)
	}
}

// roundTripperFunc creates a mock RoundTripper (transport)
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestErrorResponse_Marshal(t *testing.T) {
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

func TestRate_Marshal(t *testing.T) {
	testJSONMarshal(t, &Rate{}, "{}")

	u := &Rate{
		Limit:     1,
		Remaining: 1,
		Reset:     Timestamp{referenceTime},
	}

	want := `{
		"limit": 1,
		"remaining": 1,
		"reset": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestRateLimits_Marshal(t *testing.T) {
	testJSONMarshal(t, &RateLimits{}, "{}")

	u := &RateLimits{
		Core: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
		Search: &Rate{
			Limit:     1,
			Remaining: 1,
			Reset:     Timestamp{referenceTime},
		},
	}

	want := `{
		"core": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		},
		"search": {
			"limit": 1,
			"remaining": 1,
			"reset": ` + referenceTimeStr + `
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestParseTokenExpiration(t *testing.T) {
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
	}

	for _, tt := range tests {
		res := &http.Response{
			Request: &http.Request{},
			Header:  http.Header{},
		}

		res.Header.Set(headerTokenExpiration, tt.header)
		exp := parseTokenExpiration(res)
		if !exp.Equal(tt.want) {
			t.Errorf("parseTokenExpiration returned %#v, want %#v", exp, tt.want)
		}
	}
}
