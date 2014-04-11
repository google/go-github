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
