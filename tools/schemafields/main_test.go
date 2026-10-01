// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

// TestRepoRelative checks that the summary names a file the way the findings do, with forward
// slashes, on every platform. A Windows build is where this is worth a test of its own: there,
// filepath.Rel separates with a backslash, and the summary would name the exceptions file
// differently from the findings printed above it.
func TestRepoRelative(t *testing.T) {
	t.Parallel()
	const repo = "/checkout"
	assertEqual(t, "tools/schemafields/exceptions.txt",
		repoRelative(repo, "/checkout/tools/schemafields/exceptions.txt"))

	// A path outside the checkout keeps the spelling it was given.
	assertEqual(t, "/elsewhere/exceptions.txt", repoRelative(repo, "/elsewhere/exceptions.txt"))
}
