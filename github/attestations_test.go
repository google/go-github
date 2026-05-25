// Copyright 2024 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
)

func TestAttestation_Marshal(t *testing.T) {
	testJSONMarshalOnly(t, &Attestation{}, `{"bundle": null, "repository_id": 0}`)

	u := &Attestation{
		Bundle:       json.RawMessage(`{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json"}`),
		RepositoryID: 1,
	}

	want := `{
		"bundle": {
			"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json"
		},
		"repository_id": 1
	}`

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}

func TestAttestationsResponse_Marshal(t *testing.T) {
	testJSONMarshal(t, &AttestationsResponse{}, `{"attestations": null}`)

	u := &AttestationsResponse{
		Attestations: []*Attestation{
			{
				Bundle:       json.RawMessage(`{"verificationMaterial":{"certificate":{"rawBytes":"abc"}}}`),
				RepositoryID: 1,
			},
			{
				Bundle:       json.RawMessage(`{"dsseEnvelope":{"payload":"def"}}`),
				RepositoryID: 2,
			},
		},
	}

	want := `{
		"attestations": [
			{
				"bundle": {
					"verificationMaterial": {
						"certificate": {
							"rawBytes": "abc"
						}
					}
				},
				"repository_id": 1
			},
			{
				"bundle": {
					"dsseEnvelope": {
						"payload": "def"
					}
				},
				"repository_id": 2
			}
		]
	}`

	testJSONMarshal(t, u, want, cmpJSONRawMessageComparator())
}
