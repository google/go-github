// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

func main() {
	file, content := getFileAndContent()
	opts := struct {
		Mode string
	}{Mode: "gfm"}

	_ = Ptr(string(content))

	_ = Ptr(file) // want `replace github.Ptr\(file\) with &file`

	other := "b.txt"
	_ = Ptr(other)     // want `replace github.Ptr\(other\) with &other`
	_ = Ptr(opts.Mode) // want `replace github.Ptr\(opts.Mode\) with &opts.Mode`

	for _, loopFile := range []string{"x", "y"} {
		_ = Ptr(loopFile) // want `replace github.Ptr\(loopFile\) with &loopFile`
	}

	name := "before"
	_ = Ptr(name) // want `replace github.Ptr\(name\) with &name`
	name = "after"
	_ = Ptr(name) // want `replace github.Ptr\(name\) with &name`

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
