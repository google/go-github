// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ReviewPersonalAccessTokenRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := ReviewPersonalAccessTokenRequestOptions{
		Action: "a",
		Reason: String("r"),
	}

	mux.HandleFunc("/orgs/o/personal-access-token-requests/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(ReviewPersonalAccessTokenRequestOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPost)
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	res, err := client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned error: %v", err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned %v, want %v", res.StatusCode, http.StatusNoContent)
	}

	const methodName = "ReviewPersonalAccessTokenRequest"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "\n", 0, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", 1, input)
	})
}

func TestReviewPersonalAccessTokenRequestOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ReviewPersonalAccessTokenRequestOptions{}, "{}")

	u := &ReviewPersonalAccessTokenRequestOptions{
		Action: "a",
		Reason: String("r"),
	}

	want := `{
		"action": "a",
		"reason": "r"
	}`

	testJSONMarshal(t, u, want)
}
