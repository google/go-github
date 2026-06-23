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
	"log"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/google/go-github/v88/github"
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
		c, err := github.NewClient()
		if err != nil {
			log.Fatalf("Error creating GitHub client: %v", err)
		}
		client = c
	} else {
		c, err := github.NewClient(github.WithAuthToken(token))
		if err != nil {
			log.Fatalf("Error creating GitHub client with token: %v", err)
		}
		client = c
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
	ctx := context.Background()

	req, err := client.NewRequest(ctx, "GET", urlStr, nil)
	if err != nil {
		return err
	}

	raw := new(json.RawMessage)
	_, err = client.Do(req, raw)
	if err != nil {
		return err
	}

	missing, err := missingFields(*raw, typ, *skipURLs)
	if err != nil {
		return err
	}
	for _, m := range missing {
		fmt.Printf("%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), m.key, m.value)
	}

	return nil
}

// missingField is a JSON key returned by the API that has no corresponding
// struct field in the relevant go-github type.
type missingField struct {
	key   string
	value any
}

// missingFields returns, sorted by key, the JSON keys present in the API
// response raw that have no corresponding field in typ. typ may be a pointer to
// a struct or a pointer to a slice of structs (the element type is used, and
// the first element of the response array is inspected).
//
// The expected set of struct keys is derived from the type via reflection
// rather than by re-marshaling a decoded value. Re-marshaling drops keys for
// any field whose value is the zero value because go-github fields use
// ",omitempty"; a field that the API returns as null therefore decodes to a nil
// pointer and disappears on re-marshal, which previously caused those fields to
// be reported as missing even though they are mapped (see issue #576).
func missingFields(raw json.RawMessage, typ any, skipURLs bool) ([]*missingField, error) {
	v := reflect.Indirect(reflect.ValueOf(typ))

	// Decode the raw response into a map of its top-level keys. For a slice
	// response, inspect the first element.
	var m1 map[string]any
	if v.Kind() == reflect.Slice {
		var s []map[string]any
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		if len(s) > 0 {
			m1 = s[0]
		}
	} else {
		if err := json.Unmarshal(raw, &m1); err != nil {
			return nil, err
		}
	}

	t := v.Type()
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	want := jsonFieldNames(t)

	var missing []*missingField
	for k, val := range m1 {
		if skipURLs && strings.HasSuffix(k, "_url") {
			continue
		}
		if _, ok := want[k]; !ok {
			missing = append(missing, &missingField{key: k, value: val})
		}
	}
	slices.SortFunc(missing, func(a, b *missingField) int {
		return strings.Compare(a.key, b.key)
	})
	return missing, nil
}

// jsonFieldNames returns the set of top-level JSON object keys that a value of
// type t marshals to, following encoding/json's rules: it skips unexported
// fields and fields tagged `json:"-"`, uses the tag name (the part before the
// first comma) when present and otherwise the field name, and promotes the
// fields of anonymous (embedded) structs. Unlike re-marshaling a value, this is
// value-independent, so a ",omitempty" field is never dropped.
func jsonFieldNames(t reflect.Type) map[string]struct{} {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	names := map[string]struct{}{}
	if t.Kind() != reflect.Struct {
		return names
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}

		// An embedded struct without an explicit JSON name has its fields
		// promoted to the parent's top-level keys.
		if f.Anonymous && name == "" {
			ft := f.Type
			for ft.Kind() == reflect.Pointer {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				for k := range jsonFieldNames(ft) {
					names[k] = struct{}{}
				}
				continue
			}
		}

		if !f.IsExported() {
			continue
		}
		if name == "" {
			name = f.Name
		}
		names[name] = struct{}{}
	}

	return names
}
