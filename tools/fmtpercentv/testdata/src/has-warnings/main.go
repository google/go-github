// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	_ = fmt.Sprintf("some/%d/url", 1)          // want `use %v instead of %d`
	_ = fmt.Sprintf("some/%s/url", "yo")       // want `use %v instead of %s`
	_ = fmt.Sprintf("some/%s/%d/url", "yo", 1) // want `use %v instead of %s and %d`
	_ = fmt.Sprintf("some/%d/%s/url", 1, "yo") // want `use %v instead of %s and %d`
	fmt.Printf("some %d", 1)                   // want `use %v instead of %d`
	fmt.Printf("some %s", "thing")             // want `use %v instead of %s`
}
