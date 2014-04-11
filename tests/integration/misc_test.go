// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import "testing"

func TestEmojis(t *testing.T) {
	emoji, _, err := client.ListEmojis()
	if err != nil {
		t.Fatalf("ListEmojis returned error: %v", err)
	}

	if len(emoji) == 0 {
		t.Errorf("ListEmojis returned no emojis")
	}

	if _, ok := emoji["+1"]; !ok {
		t.Errorf("ListEmojis missing '+1' emoji")
	}
}

func TestAPIMeta(t *testing.T) {
	meta, _, err := client.APIMeta()
	if err != nil {
		t.Fatalf("APIMeta returned error: %v", err)
	}

	if len(meta.Hooks) == 0 {
		t.Errorf("APIMeta returned no hook addresses")
	}

	if len(meta.Git) == 0 {
		t.Errorf("APIMeta returned no git addresses")
	}

	if !*meta.VerifiablePasswordAuthentication {
		t.Errorf("APIMeta VerifiablePasswordAuthentication is false")
	}
}
