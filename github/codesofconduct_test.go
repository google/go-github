// Copyright 2023 The go-github AUTHORS. All rights reserved.
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

func TestCodesOfConductService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `[{
						"key": "key",
						"name": "name",
						"url": "url"}
						]`)
	})

	ctx := context.Background()
	cs, _, err := client.ListCodesOfConduct(ctx)
	assertNilError(t, err)

	want := []*CodeOfConduct{
		{
			Key:  String("key"),
			Name: String("name"),
			URL:  String("url"),
		}}
	if !cmp.Equal(want, cs) {
		t.Errorf("returned %+v, want %+v", cs, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodesOfConduct.List(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodesOfConductService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct/k", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `{
						"key": "key",
						"name": "name",
						"url": "url",
						"body": "body"}`,
		)
	})

	ctx := context.Background()
	coc, _, err := client.GetCodeOfConduct(ctx, "k")
	assertNilError(t, err)

	want := &CodeOfConduct{
		Key:  String("key"),
		Name: String("name"),
		URL:  String("url"),
		Body: String("body"),
	}
	if !cmp.Equal(want, coc) {
		t.Errorf("returned %+v, want %+v", coc, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodesOfConduct.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodesOfConduct.Get(ctx, "k")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeOfConduct_Marshal(t *testing.T) {
	testJSONMarshal(t, &CodeOfConduct{}, "{}")

	a := &CodeOfConduct{
		Name: String("name"),
		Key:  String("key"),
		URL:  String("url"),
		Body: String("body"),
	}

	want := `{
		"name": "name",
		"key": "key",
		"url": "url",
		"body": "body"
	}`

	testJSONMarshal(t, a, want)
}
