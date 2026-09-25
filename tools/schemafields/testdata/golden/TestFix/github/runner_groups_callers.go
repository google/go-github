// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// Callers builds the requests whose field types -fix changes, so that the compiler reports
// the literals that have to be repaired with them.
func Callers() []any {
	return []any{
		// The field is named by the literal.
		CreateRunnerGroupRequest{Name: "a", Visibility: new("all")},

		// The field is named by the literal, and the change only unwraps the type.
		UpdateRunnerGroupRequest{Name: "b"},

		// The literal takes its type from the slice that holds it.
		[]*UpdateRunnerGroupRequest{{Name: "c"}},
	}
}
