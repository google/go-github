// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	for key, want := range values {
		if v := r.FormValue(key); v != want {
			t.Errorf("Request parameter %v = %v, want %v", key, v, want)
		}
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &User{Login: "l"}, `{"login":"l"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, want %v", inBody, string(body), outBody)
	}

	// test that default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[int]interface{}
	}
	_, err := c.NewRequest("GET", "/", &T{})

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

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.  A redirect loop is pretty unlikely to occur within the GitHub
// API, but does allow us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_rateLimit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerRateLimit, "60")
		w.Header().Add(headerRateRemaining, "59")
	})

	var want int

	if want = 0; client.Rate.Limit != want {
		t.Errorf("Client rate limit = %v, want %v", client.Rate.Limit, want)
	}
	if want = 0; client.Rate.Limit != want {
		t.Errorf("Client rate remaining, got %v", client.Rate.Remaining, want)
	}

	req, _ := client.NewRequest("GET", "/", nil)
	client.Do(req, nil)

	if want = 60; client.Rate.Limit != want {
		t.Errorf("Client rate limit = %v, want %v", client.Rate.Limit, want)
	}
	if want = 59; client.Rate.Remaining != want {
		t.Errorf("Client rate remaining = %v, want %v", client.Rate.Remaining, want)
	}
}

func TestDo_rateLimit_errorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerRateLimit, "60")
		w.Header().Add(headerRateRemaining, "59")
		http.Error(w, "Bad Request", 400)
	})

	var want int

	req, _ := client.NewRequest("GET", "/", nil)
	client.Do(req, nil)

	if want = 60; client.Rate.Limit != want {
		t.Errorf("Client rate limit = %v, want %v", client.Rate.Limit, want)
	}
	if want = 59; client.Rate.Remaining != want {
		t.Errorf("Client rate remaining = %v, want %v", client.Rate.Remaining, want)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m", 
			"errors": [{"resource": "r", "field": "f", "code": "c"}]}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
		Errors:   []Error{Error{Resource: "r", Field: "f", Code: "c"}},
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

// ensure that we properly handle API errors that do not contain a response
// body
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
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
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

func TestRateLimit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"rate":{"limit":2,"remaining":1}}`)
	})

	rate, err := client.RateLimit()
	if err != nil {
		t.Errorf("Rate limit returned error: %v", err)
	}

	want := &Rate{Limit: 2, Remaining: 1}
	if !reflect.DeepEqual(rate, want) {
		t.Errorf("RateLimit returned %+v, want %+v", rate, want)
	}
}

func TestUnauthenticatedRateLimitedTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var v, want string
		q := r.URL.Query()
		if v, want = q.Get("client_id"), "id"; v != want {
			t.Errorf("OAuth Client ID = %v, want %v", v, want)
		}
		if v, want = q.Get("client_secret"), "secret"; v != want {
			t.Errorf("OAuth Client Secret = %v, want %v", v, want)
		}
	})

	tp := &UnauthenticatedRateLimitedTransport{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	unauthedClient := NewClient(tp.Client())
	unauthedClient.BaseURL = client.BaseURL
	req, _ := unauthedClient.NewRequest("GET", "/", nil)
	unauthedClient.Do(req, nil)
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
