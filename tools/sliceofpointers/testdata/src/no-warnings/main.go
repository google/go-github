// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Example struct {
	Strings []string `json:"strings,omitempty"` // Should not be flagged
	Things  []*Thing `json:"things,omitempty"`  // Should not be flagged
}

type Thing struct {
}

func main() {
	slice1 := []string{} // Should not be flagged
	_ = slice1
	slice2 := []*Thing{} // Should not be flagged
	_ = slice2
}
