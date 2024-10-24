// Copyright 2024 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
)

type Attestation struct {
	// The attestation's Sigstore Bundle.
	// Refer to the sigstore bundle specification for more info:
	// https://github.com/sigstore/protobuf-specs/blob/main/protos/sigstore_bundle.proto
	Bundle       *json.RawMessage `json:"bundle"`
	RepositoryID *int             `json:"repository_id"`
}

type AttestationsResponse struct {
	Attestations []*Attestation `json:"attestations"`
}
