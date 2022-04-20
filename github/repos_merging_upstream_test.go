// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestRepositoriesService_MergeUpstream(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &RepositoryMergeUpstreamRequest{
		Branch: String("b"),
	}

	mux.HandleFunc("/repos/o/r/merge-upstream", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryMergeUpstreamRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"merge_type":"m"}`)
	})

	ctx := context.Background()
	result, _, err := client.Repositories.MergeUpstream(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.MergeUpstream returned error: %v", err)
	}

	want := &MergeUpstreamResult{MergeType: String("m")}
	if !cmp.Equal(result, want) {
		t.Errorf("Repositories.MergeUpstream returned %+v, want %+v", result, want)
	}

	const methodName = "MergeUpstream"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.MergeUpstream(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.MergeUpstream(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
