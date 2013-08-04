// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import

// Label represents a GitHib label on an Issue
"fmt"

type Label struct {
	URL   string `json:"url,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

func (label Label) String() string {
	return fmt.Sprint(label.Name)
}
