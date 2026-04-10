// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/google/go-github/v84/github"

func main() {
	// Literal argument cannot be addressed.
	_ = github.Ptr("a.txt")

	const file = "a.txt"
	_ = github.Ptr(file)

	for range []int{1, 2} {
		_ = github.Ptr("a")
	}

	_ = github.Ptr(getOptions().Mode)
}

func getOptions() struct {
	Mode string
} {
	return struct {
		Mode string
	}{Mode: "gfm"}
}

func unqualifiedIdentifierArgumentDoesNotWarn() {
	_ = Ptr("x")
}

func Bool(v bool) *bool { return Ptr(v) }

func Int(v int) *int { return Ptr(v) }

func Int64(v int64) *int64 { return Ptr(v) }

func String(v string) *string { return Ptr(v) }

func Ptr[T any](v T) *T {
	return &v
}
