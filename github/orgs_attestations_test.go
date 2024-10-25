// Copyright 2024 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListAttestations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/attestations/digest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"attestations": [
				{
					"repository_id": 1,
					"bundle": {}
				},
				{
					"repository_id": 2,
					"bundle": {}
				}
			]
		}`)
	})
	ctx := context.Background()
	attestations, _, err := client.Organizations.ListAttestations(ctx, "o", "digest", &ListOptions{})
	if err != nil {
		t.Errorf("Organizations.ListAttestations returned error: %v", err)
	}

	want := &AttestationsResponse{
		Attestations: []*Attestation{
			{
				RepositoryID: 1,
				Bundle:       []byte(`{}`),
			},
			{
				RepositoryID: 2,
				Bundle:       []byte(`{}`),
			},
		},
	}

	if !cmp.Equal(attestations, want) {
		t.Errorf("Organizations.ListAttestations = %+v, want %+v", attestations, want)
	}
}
