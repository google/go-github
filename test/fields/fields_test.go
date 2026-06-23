// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/google/go-github/v88/github"
)

func TestMissingFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		typ      any
		skipURLs bool
		want     []string
	}{
		{
			// Regression test for #576: a mapped field whose API value is null
			// must not be reported as missing. Previously the struct was
			// re-marshaled and ",omitempty" dropped the nil milestone key.
			name: "null-valued mapped field is not reported",
			raw:  `{"number": 1, "milestone": null, "comments": 2}`,
			typ:  &github.Issue{},
			want: nil,
		},
		{
			name: "genuinely unmapped field is reported",
			raw:  `{"milestone": null, "this_field_does_not_exist": 123}`,
			typ:  &github.Issue{},
			want: []string{"this_field_does_not_exist"},
		},
		{
			name:     "url fields are skipped when requested",
			raw:      `{"bogus_url": "x", "another_bogus": 1}`,
			typ:      &github.Issue{},
			skipURLs: true,
			want:     []string{"another_bogus"},
		},
		{
			name: "slice response inspects the first element",
			raw:  `[{"totally_bogus": true}]`,
			typ:  &[]github.Key{},
			want: []string{"totally_bogus"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := missingFields(json.RawMessage(tt.raw), tt.typ, tt.skipURLs)
			if err != nil {
				t.Fatalf("missingFields returned error: %v", err)
			}

			var keys []string
			for _, m := range got {
				keys = append(keys, m.key)
			}
			if !slices.Equal(keys, tt.want) {
				t.Errorf("missingFields keys = %v, want %v", keys, tt.want)
			}
		})
	}
}
