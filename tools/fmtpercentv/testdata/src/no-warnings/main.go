// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	_ = fmt.Sprintf("some/%v/url", 1) // Should not be flagged
	fmt.Printf("some %v", 1)          // Should not be flagged
	fmt.Printf("some %v", "thing")    // Should not be flagged
}
