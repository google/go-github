// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAPIMeta_Marshal(t *testing.T) {
	testJSONMarshal(t, &APIMeta{}, "{}")

	a := &APIMeta{
		Hooks:                            []string{"h"},
		Git:                              []string{"g"},
		VerifiablePasswordAuthentication: Bool(true),
		Pages:                            []string{"p"},
		Importer:                         []string{"i"},
		Actions:                          []string{"a"},
		Dependabot:                       []string{"d"},
		SSHKeyFingerprints:               map[string]string{"a": "f"},
		SSHKeys:                          []string{"k"},
		API:                              []string{"a"},
		Web:                              []string{"w"},
	}
	want := `{
		"hooks":["h"],
		"git":["g"],
		"verifiable_password_authentication":true,
		"pages":["p"],
		"importer":["i"],
		"actions":["a"],
		"dependabot":["d"],
		"ssh_key_fingerprints":{"a":"f"},
		"ssh_keys":["k"],
		"api":["a"],
		"web":["w"]
	}`

	testJSONMarshal(t, a, want)
}

func TestMetaService_APIMeta(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"web":["w"],"api":["a"],"hooks":["h"], "git":["g"], "pages":["p"], "importer":["i"], "actions":["a"], "dependabot":["d"], "verifiable_password_authentication": true}`)
	})

	ctx := context.Background()
	meta, _, err := client.APIMeta(ctx)
	if err != nil {
		t.Errorf("APIMeta returned error: %v", err)
	}

	want := &APIMeta{
		Hooks:      []string{"h"},
		Git:        []string{"g"},
		Pages:      []string{"p"},
		Importer:   []string{"i"},
		Actions:    []string{"a"},
		Dependabot: []string{"d"},
		API:        []string{"a"},
		Web:        []string{"w"},

		VerifiablePasswordAuthentication: Bool(true),
	}
	if !cmp.Equal(want, meta) {
		t.Errorf("APIMeta returned %+v, want %+v", meta, want)
	}

	const methodName = "APIMeta"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.APIMeta(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMetaService_Octocat(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := "input"
	output := "sample text"

	mux.HandleFunc("/octocat", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"s": input})
		w.Header().Set("Content-Type", "application/octocat-stream")
		fmt.Fprint(w, output)
	})

	ctx := context.Background()
	got, _, err := client.Octocat(ctx, input)
	if err != nil {
		t.Errorf("Octocat returned error: %v", err)
	}

	if want := output; got != want {
		t.Errorf("Octocat returned %+v, want %+v", got, want)
	}

	const methodName = "Octocat"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Octocat(ctx, input)
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMetaService_Zen(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	output := "sample text"

	mux.HandleFunc("/zen", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		fmt.Fprint(w, output)
	})

	ctx := context.Background()
	got, _, err := client.Zen(ctx)
	if err != nil {
		t.Errorf("Zen returned error: %v", err)
	}

	if want := output; got != want {
		t.Errorf("Zen returned %+v, want %+v", got, want)
	}

	const methodName = "Zen"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Zen(ctx)
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
