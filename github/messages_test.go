// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"net/http"
	"testing"
)

func TestValidatePayload(t *testing.T) {
	const defaultBody = `{"yo":true}` // All tests below use the default request body and signature.
	const defaultSignature = "sha1=126f2c800419c60137ce748d7672e77b65cf16d6"
	secretKey := []byte("0123456789abcdef")
	tests := []struct {
		signature   string
		eventID     string
		event       string
		wantEventID string
		wantEvent   string
		wantPayload string
	}{
		// The following tests generate expected errors:
		{},                         // Missing signature
		{signature: "yo"},          // Missing signature prefix
		{signature: "sha1=yo"},     // Signature not hex string
		{signature: "sha1=012345"}, // Invalid signature
		// The following tests expect err=nil:
		{
			signature:   defaultSignature,
			eventID:     "dead-beef",
			event:       "ping",
			wantEventID: "dead-beef",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   defaultSignature,
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   "sha256=b1f8020f5b4cd42042f807dd939015c4a418bc1ff7f604dd55b0a19b5d953d9b",
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   "sha512=8456767023c1195682e182a23b3f5d19150ecea598fde8cb85918f7281b16079471b1329f92b912c4d8bd7455cb159777db8f29608b20c7c87323ba65ae62e1f",
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
	}

	for _, test := range tests {
		buf := bytes.NewBufferString(defaultBody)
		req, err := http.NewRequest("GET", "http://localhost/event", buf)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		if test.signature != "" {
			req.Header.Set(signatureHeader, test.signature)
		}

		got, err := ValidatePayload(req, secretKey)
		if err != nil {
			if test.wantPayload != "" {
				t.Errorf("ValidatePayload(%#v): err = %v, want nil", test, err)
			}
			continue
		}
		if string(got) != test.wantPayload {
			t.Errorf("ValidatePayload = %q, want %q", got, test.wantPayload)
		}
	}
}
