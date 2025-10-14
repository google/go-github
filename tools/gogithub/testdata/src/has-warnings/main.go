// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "time"

type Example struct {
	Strings []*string `json:"strings,omitempty"` // want `use \[\]string instead of \[\]\*string`
	Things  []Thing   `json:"things,omitempty"`  // want `use \[\]\*Thing instead of \[\]Thing`
}

type Thing struct {
}

type OmitemptyExample struct {
	Name   string    `json:"name,omitempty"`   // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
	Age    int       `json:"age,omitempty"`    // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
	Count  int64     `json:"count,omitempty"`  // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
	Active bool      `json:"active,omitempty"` // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
	Thing  Thing     `json:"thing,omitempty"`  // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
	Time   time.Time `json:"time,omitempty"`   // want `using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted`
}

func main() {
	slice1 := []*string{} // want `use \[\]string instead of \[\]\*string`
	_ = slice1
	slice2 := []Thing{} // want `use \[\]\*Thing instead of \[\]Thing`
	_ = slice2
}
