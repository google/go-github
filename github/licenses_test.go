// Copyright 2013 The go-github AUTHORS. All rights reserved.
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

func TestRepositoryLicense_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoryLicense{}, "{}")

	rl := &RepositoryLicense{
		Name:        String("n"),
		Path:        String("p"),
		SHA:         String("s"),
		Size:        Int(1),
		URL:         String("u"),
		HTMLURL:     String("h"),
		GitURL:      String("g"),
		DownloadURL: String("d"),
		Type:        String("t"),
		Content:     String("c"),
		Encoding:    String("e"),
		License: &License{
			Key:            String("k"),
			Name:           String("n"),
			URL:            String("u"),
			SPDXID:         String("s"),
			HTMLURL:        String("h"),
			Featured:       Bool(true),
			Description:    String("d"),
			Implementation: String("i"),
			Permissions:    &[]string{"p"},
			Conditions:     &[]string{"c"},
			Limitations:    &[]string{"l"},
			Body:           String("b"),
		},
	}
	want := `{
		"name": "n",
		"path": "p",
		"sha": "s",
		"size": 1,
		"url": "u",
		"html_url": "h",
		"git_url": "g",
		"download_url": "d",
		"type": "t",
		"content": "c",
		"encoding": "e",
		"license": {
			"key": "k",
			"name": "n",
			"url": "u",
			"spdx_id": "s",
			"html_url": "h",
			"featured": true,
			"description": "d",
			"implementation": "i",
			"permissions": ["p"],
			"conditions": ["c"],
			"limitations": ["l"],
			"body": "b"
		}
	}`
	testJSONMarshal(t, rl, want)
}

func TestLicense_Marshal(t *testing.T) {
	testJSONMarshal(t, &License{}, "{}")

	l := &License{
		Key:            String("k"),
		Name:           String("n"),
		URL:            String("u"),
		SPDXID:         String("s"),
		HTMLURL:        String("h"),
		Featured:       Bool(true),
		Description:    String("d"),
		Implementation: String("i"),
		Permissions:    &[]string{"p"},
		Conditions:     &[]string{"c"},
		Limitations:    &[]string{"l"},
		Body:           String("b"),
	}
	want := `{
		"key": "k",
		"name": "n",
		"url": "u",
		"spdx_id": "s",
		"html_url": "h",
		"featured": true,
		"description": "d",
		"implementation": "i",
		"permissions": ["p"],
		"conditions": ["c"],
		"limitations": ["l"],
		"body": "b"
	}`
	testJSONMarshal(t, l, want)
}

func TestLicensesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/licenses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"key":"mit","name":"MIT","spdx_id":"MIT","url":"https://api.github.com/licenses/mit","featured":true}]`)
	})

	ctx := context.Background()
	licenses, _, err := client.Licenses.List(ctx)
	if err != nil {
		t.Errorf("Licenses.List returned error: %v", err)
	}

	want := []*License{{
		Key:      String("mit"),
		Name:     String("MIT"),
		SPDXID:   String("MIT"),
		URL:      String("https://api.github.com/licenses/mit"),
		Featured: Bool(true),
	}}
	if !cmp.Equal(licenses, want) {
		t.Errorf("Licenses.List returned %+v, want %+v", licenses, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Licenses.List(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestLicensesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/licenses/mit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key":"mit","name":"MIT"}`)
	})

	ctx := context.Background()
	license, _, err := client.Licenses.Get(ctx, "mit")
	if err != nil {
		t.Errorf("Licenses.Get returned error: %v", err)
	}

	want := &License{Key: String("mit"), Name: String("MIT")}
	if !cmp.Equal(license, want) {
		t.Errorf("Licenses.Get returned %+v, want %+v", license, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Licenses.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Licenses.Get(ctx, "mit")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestLicensesService_Get_invalidTemplate(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Licenses.Get(ctx, "%")
	testURLParseError(t, err)
}
