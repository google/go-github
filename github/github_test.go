// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
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
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

// openTestFile creates a new file with the given name and content for testing.
// In order to ensure the exact file name, this function will create a new temp
// directory, and create the file in that directory.  It is the caller's
// responsibility to remove the directy and its contents when no longer needed.
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
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if value := r.Header.Get(header); want != value {
		t.Errorf("Header %s = %s, want: %s", header, value, want)
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

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read body")
	}
	str := string(b)
	if want != str {
		t.Errorf("Body = %s, want: %s", str, want)
	}
}

// Helper function to test that a value is marshalled to JSON as expected.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %v", v)
	}

	w := new(bytes.Buffer)
	err = json.Compact(w, []byte(want))
	if err != nil {
		t.Errorf("String is not valid json: %s", want)
	}

	if w.String() != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, want %s", v, j, w)
	}

	// now go the other direction and make sure things unmarshal as expected
	u := reflect.ValueOf(v).Interface()
	if err := json.Unmarshal([]byte(want), u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v", want)
	}

	if !reflect.DeepEqual(v, u) {
		t.Errorf("json.Unmarshal(%q) returned %s, want %s", want, u, v)
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
	inBody, outBody := &User{Login: String("l")}, `{"login":"l"}`+"\n"
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
	if want, got := 1, response.FirstPage; want != got {
		t.Errorf("response.FirstPage: %v, want %v", want, got)
	}
	if want, got := 2, response.PrevPage; want != got {
		t.Errorf("response.PrevPage: %v, want %v", want, got)
	}
	if want, got := 4, response.NextPage; want != got {
		t.Errorf("response.NextPage: %v, want %v", want, got)
	}
	if want, got := 5, response.LastPage; want != got {
		t.Errorf("response.LastPage: %v, want %v", want, got)
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
	if want, got := 0, response.FirstPage; want != got {
		t.Errorf("response.FirstPage: %v, want %v", want, got)
	}
	if want, got := 0, response.PrevPage; want != got {
		t.Errorf("response.PrevPage: %v, want %v", want, got)
	}
	if want, got := 0, response.NextPage; want != got {
		t.Errorf("response.NextPage: %v, want %v", want, got)
	}
	if want, got := 0, response.LastPage; want != got {
		t.Errorf("response.LastPage: %v, want %v", want, got)
	}

	// more invalid URLs
	r = http.Response{
		Header: http.Header{
			"Link": {`<https://api.github.com/%?page=2>; rel="first"`},
		},
	}

	response = newResponse(&r)
	if want, got := 0, response.FirstPage; want != got {
		t.Errorf("response.FirstPage: %v, want %v", want, got)
	}
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
		w.Header().Add(headerRateReset, "1372700873")
	})

	var want int

	if want = 0; client.Rate.Limit != want {
		t.Errorf("Client rate limit = %v, want %v", client.Rate.Limit, want)
	}
	if want = 0; client.Rate.Limit != want {
		t.Errorf("Client rate remaining = %v, got %v", client.Rate.Remaining, want)
	}
	if !client.Rate.Reset.IsZero() {
		t.Errorf("Client rate reset not initialized to zero value")
	}

	req, _ := client.NewRequest("GET", "/", nil)
	client.Do(req, nil)

	if want = 60; client.Rate.Limit != want {
		t.Errorf("Client rate limit = %v, want %v", client.Rate.Limit, want)
	}
	if want = 59; client.Rate.Remaining != want {
		t.Errorf("Client rate remaining = %v, want %v", client.Rate.Remaining, want)
	}
	reset := time.Date(2013, 7, 1, 17, 47, 53, 0, time.UTC)
	if client.Rate.Reset.UTC() != reset {
		t.Errorf("Client rate reset = %v, want %v", client.Rate.Reset, reset)
	}
}

func TestDo_rateLimit_errorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerRateLimit, "60")
		w.Header().Add(headerRateRemaining, "59")
		w.Header().Add(headerRateReset, "1372700873")
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
	reset := time.Date(2013, 7, 1, 17, 47, 53, 0, time.UTC)
	if client.Rate.Reset.UTC() != reset {
		t.Errorf("Client rate reset = %v, want %v", client.Rate.Reset, reset)
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
		Errors:   []Error{{Resource: "r", Field: "f", Code: "c"}},
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
		//fmt.Fprint(w, `{"resources":{"core": {"limit":2,"remaining":1,"reset":1372700873}}}`)
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874}
		}}`)
	})

	rate, _, err := client.RateLimit()
	if err != nil {
		t.Errorf("Rate limit returned error: %v", err)
	}

	want := &Rate{
		Limit:     2,
		Remaining: 1,
		Reset:     Timestamp{time.Date(2013, 7, 1, 17, 47, 53, 0, time.UTC).Local()},
	}
	if !reflect.DeepEqual(rate, want) {
		t.Errorf("RateLimit returned %+v, want %+v", rate, want)
	}
}

func TestRateLimits(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"resources":{
			"core": {"limit":2,"remaining":1,"reset":1372700873},
			"search": {"limit":3,"remaining":2,"reset":1372700874}
		}}`)
	})

	rate, _, err := client.RateLimits()
	if err != nil {
		t.Errorf("RateLimits returned error: %v", err)
	}

	want := &RateLimits{
		Core: &Rate{
			Limit:     2,
			Remaining: 1,
			Reset:     Timestamp{time.Date(2013, 7, 1, 17, 47, 53, 0, time.UTC).Local()},
		},
		Search: &Rate{
			Limit:     3,
			Remaining: 2,
			Reset:     Timestamp{time.Date(2013, 7, 1, 17, 47, 54, 0, time.UTC).Local()},
		},
	}
	if !reflect.DeepEqual(rate, want) {
		t.Errorf("RateLimits returned %+v, want %+v", rate, want)
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
