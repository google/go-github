// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool tests for the JSON mappings in the go-github data types. It will
// identify fields that are returned by the live GitHub API, but that are not
// currently mapped into a struct field of the relevant go-github type. This
// helps to ensure that all relevant data returned by the API is being made
// accessible, particularly new fields that are periodically (and sometimes
// quietly) added to the API over time.
//
// These tests simply aid in identifying which fields aren't being mapped; it
// is not necessarily true that every one of them should always be mapped.
// Some fields may be undocumented for a reason, either because they aren't
// actually used yet or should not be relied upon.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/google/go-github/v84/github"
)

var (
	client *github.Client

	skipURLs = flag.Bool("skip_urls", false, "skip url fields")
)

func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		fmt.Print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = github.NewClient(nil)
	} else {
		client = github.NewClient(nil).WithAuthToken(token)
	}

	for _, tt := range []struct {
		url string
		typ any
	}{
		{"users/octocat", &github.User{}},
		{"user", &github.User{}},
		{"users/willnorris/keys", &[]github.Key{}},
		{"orgs/google-test", &github.Organization{}},
		{"repos/google/go-github", &github.Repository{}},
		{"repos/google/go-github/issues/1", &github.Issue{}},
		{"/gists/9257657", &github.Gist{}},
	} {
		err := testType(tt.url, tt.typ)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

// testType fetches the JSON resource at urlStr and compares its keys to the
// struct fields of typ.
func testType(urlStr string, typ any) error {
	slice := reflect.Indirect(reflect.ValueOf(typ)).Kind() == reflect.Slice

	req, err := client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	// start with a json.RawMessage so we can decode multiple ways below
	raw := new(json.RawMessage)
	_, err = client.Do(context.Background(), req, raw)
	if err != nil {
		return err
	}

	// unmarshal directly to a map
	var m1 map[string]any
	if slice {
		var s []map[string]any
		err = json.Unmarshal(*raw, &s)
		if err != nil {
			return err
		}
		m1 = s[0]
	} else {
		err = json.Unmarshal(*raw, &m1)
		if err != nil {
			return err
		}
	}

	// unmarshal to typ first, then re-marshal and unmarshal to a map
	err = json.Unmarshal(*raw, typ)
	if err != nil {
		return err
	}

	var byt []byte
	if slice {
		// use first item in slice
		v := reflect.Indirect(reflect.ValueOf(typ))
		byt, err = json.Marshal(v.Index(0).Interface())
		if err != nil {
			return err
		}
	} else {
		byt, err = json.Marshal(typ)
		if err != nil {
			return err
		}
	}

	var m2 map[string]any
	err = json.Unmarshal(byt, &m2)
	if err != nil {
		return err
	}

	// now compare the two maps
	for k, v := range m1 {
		if *skipURLs && strings.HasSuffix(k, "_url") {
			continue
		}
		if _, ok := m2[k]; !ok {
			fmt.Printf("%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), k, v)
		}
	}

	return nil
}
