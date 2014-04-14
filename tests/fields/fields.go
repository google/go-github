// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool tests for the JSON mappings in the go-github data types.  It will
// identify fields that are returned by the live GitHub API, but that are not
// currently mapped into a struct field of the relevant go-github type.  This
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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"code.google.com/p/goauth2/oauth"

	"github.com/google/go-github/github"
)

var (
	client *github.Client

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool

	skipURLs = flag.Bool("skip_urls", false, "skip url fields")
)

func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		println("!!! No OAuth token.  Some tests won't run. !!!\n")
		client = github.NewClient(nil)
	} else {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		client = github.NewClient(t.Client())
		auth = true
	}

	//testType("rate_limit", &github.RateLimits{})
	testType("users/octocat", &github.User{})
	testType("orgs/google", &github.Organization{})
	testType("repos/google/go-github", &github.Repository{})
}

// testType fetches the JSON resource at urlStr and compares its keys to the
// struct fields of typ.
//
// TODO: handle resources that are more easily fetched as an array of objects,
// rather than a single object (e.g. a user's public keys).  In this case, we
// should just take the first object in the array, and use that.  In that case,
// should typ also be specified as a slice?
func testType(urlStr string, typ interface{}) error {
	req, err := client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	// start with a json.RawMessage so we can decode multiple ways below
	raw := new(json.RawMessage)
	_, err = client.Do(req, raw)
	if err != nil {
		return err
	}

	// unmarshall directly to a map
	var m1 map[string]interface{}
	err = json.Unmarshal(*raw, &m1)
	if err != nil {
		return err
	}

	// unarmshall to typ first, then re-marshall and unmarshall to a map
	err = json.Unmarshal(*raw, typ)
	if err != nil {
		return err
	}

	byt, err := json.Marshal(typ)
	if err != nil {
		return err
	}

	var m2 map[string]interface{}
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
