// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Example struct {
	Strings []*string `json:"strings,omitempty"` // want `use \[\]string instead of \[\]\*string`
	Things  []Thing   `json:"things,omitempty"`  // want `use \[\]\*Thing instead of \[\]Thing`
}

type Thing struct {
}

func main() {
	slice1 := []*string{} // want `use \[\]string instead of \[\]\*string`
	_ = slice1
	slice2 := []Thing{} // want `use \[\]\*Thing instead of \[\]Thing`
	_ = slice2
}
