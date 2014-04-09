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

	testType("rate_limit", github.Rate{})
	testType("users/octocat", github.User{})
	testType("orgs/google", github.Organization{})
	testType("repos/google/go-github", github.Repository{})
}

func checkAuth(name string) bool {
	if !auth {
		fmt.Printf("No auth - skipping portions of %v\n", name)
	}
	return auth
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

	// I'm thinking we might want to unmarshall the response both as a
	// map[string]interface{} as well as typ, though I'm not 100% sure.
	// That's why we unmarshal to json.RawMessage first here.
	raw := new(json.RawMessage)
	_, err = client.Do(req, raw)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	err = json.Unmarshal(*raw, &m)
	if err != nil {
		return err
	}

	fields := jsonFields(typ)

	for k, v := range m {
		if *skipURLs && strings.HasSuffix(k, "_url") {
			continue
		}
		if _, ok := fields[k]; !ok {
			fmt.Printf("%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), k, v)
		}
	}

	return nil
}

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// jsonFields returns a map of JSON fields that have an explicit mapping to a
// field in v.  The fields will be in the returned map's keys; the map values
// for those keys is currently undefined.
func jsonFields(v interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	typ := reflect.TypeOf(v)
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" { // unexported
			continue
		}

		tag := sf.Tag.Get("json")
		if tag == "-" {
			continue
		}

		name, opts := parseTag(tag)
		fields[name] = opts
	}
	return fields
}
