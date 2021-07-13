// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "testing"

func TestInteractionRestriction_Marshal(t *testing.T) {
	testJSONMarshal(t, &InteractionRestriction{}, "{}")

	u := &InteractionRestriction{
		Limit:     String("limit"),
		Origin:    String("origin"),
		ExpiresAt: &Timestamp{referenceTime},
	}

	want := `{
		"limit": "limit",
		"origin": "origin",
		"expires_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}
