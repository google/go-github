// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func Test_normalizedOpName(t *testing.T) {
	for _, td := range []struct {
		name string
		want string
	}{
		{name: "", want: ""},
		{name: "get /foo/{id}", want: "GET /foo/*"},
		{name: "get foo", want: "GET /foo"},
	} {
		t.Run(td.name, func(t *testing.T) {
			got := normalizedOpName(td.name)
			if got != td.want {
				t.Errorf("normalizedOpName() = %v, want %v", got, td.want)
			}
		})
	}
}
