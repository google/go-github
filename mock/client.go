package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

type RequestMatch = string

// Users
var RequestMatchUsersGet RequestMatch = "github.(*UsersService).Get"

// Orgs
var RequestMatchOrganizationsList RequestMatch = "github.(*OrganizationsService).List"

func MatchIncomingRequest(r *http.Request, rm RequestMatch) bool {
	pc := make([]uintptr, 100)
	n := runtime.Callers(0, pc)
	if n == 0 {
		return false
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	for {
		frame, more := frames.Next()

		if strings.Contains(frame.File, "go-github") &&
			strings.Contains(frame.Function, "github.") &&
			strings.Contains(frame.Function, "Service") {
			splitFuncName := strings.Split(frame.Function, "/")
			methodCall := splitFuncName[len(splitFuncName)-1]

			if methodCall == rm {
				return true
			}
		}

		if !more {
			return false
		}
	}
}

type MockRoundTripper struct {
	RequestMocks map[RequestMatch][][]byte
}

// RoundTrip implements http.RoundTripper interface
func (mrt *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	for requestMatch, respBodies := range mrt.RequestMocks {
		if MatchIncomingRequest(r, requestMatch) {
			if len(respBodies) == 0 {
				fmt.Printf(
					"no more available mocked responses for endpoit %s\n",
					r.URL.Path,
				)

				fmt.Println("please add the required RequestMatch to the MockHttpClient. Eg.")
				fmt.Println(`
				mockedHttpClient := NewMockHttpClient(
					WithRequestMatch(
						RequestMatchUsersGet,
						MustMarshall(github.User{
							Name: github.String("foobar"),
						}),
					),
					WithRequestMatch(
						RequestMatchOrganizationsList,
						MustMarshall([]github.Organization{
							{
								Name: github.String("foobar123"),
							},
						}),
					),
				)
				`)

				panic(nil)
			}

			resp := respBodies[0]

			defer func(mrt *MockRoundTripper, rm RequestMatch) {
				mrt.RequestMocks[rm] = mrt.RequestMocks[rm][1:]
			}(mrt, requestMatch)

			re := bytes.NewReader(resp)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(re),
			}, nil
		}
	}

	return nil, fmt.Errorf(
		"couldn find a mock request that matches the request sent to: %s",
		r.URL.Path,
	)

}

var _ http.RoundTripper = &MockRoundTripper{}

type MockHttpClientOption func(*MockRoundTripper)

func WithRequestMatch(
	rm RequestMatch,
	marshalled []byte,
) MockHttpClientOption {
	return func(mrt *MockRoundTripper) {
		if _, found := mrt.RequestMocks[rm]; !found {
			mrt.RequestMocks[rm] = make([][]byte, 0)
		}

		mrt.RequestMocks[rm] = append(
			mrt.RequestMocks[rm],
			marshalled,
		)
	}
}

func NewMockHttpClient(options ...MockHttpClientOption) *http.Client {
	rt := &MockRoundTripper{
		RequestMocks: make(map[RequestMatch][][]byte),
	}

	for _, o := range options {
		o(rt)
	}

	return &http.Client{
		Transport: rt,
	}
}

func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)

	if err == nil {
		return b
	}

	panic(err)
}
