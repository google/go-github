// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/go-github/v84/github"
)

func main() {
	file, content := getFileAndContent()
	opts := struct {
		Mode string
	}{Mode: "gfm"}

	_ = github.Ptr(string(content))

	_ = github.Ptr(file) // want `replace github.Ptr\(file\) with &file`

	other := "b.txt"
	_ = github.Ptr(other)     // want `replace github.Ptr\(other\) with &other`
	_ = github.Ptr(opts.Mode) // want `replace github.Ptr\(opts.Mode\) with &opts.Mode`

	for _, loopFile := range []string{"x", "y"} {
		_ = github.Ptr(loopFile) // want `replace github.Ptr\(loopFile\) with &loopFile`
	}

	name := "before"
	_ = github.Ptr(name) // want `replace github.Ptr\(name\) with &name`
	name = "after"
	_ = github.Ptr(name) // want `replace github.Ptr\(name\) with &name`

	i := 1
	_ = Ptr(i)         // want `replace github.Ptr\(i\) with &i`
	_ = Ptr(opts.Mode) // want `replace github.Ptr\(opts.Mode\) with &opts.Mode`
}

func getFileAndContent() (string, []byte) {
	return "", nil
}

func Ptr[T any](v T) *T {
	return &v
}
