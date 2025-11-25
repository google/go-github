// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structfield

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRun(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()
	plugin, _ := New(map[string]any{
		"allowed-tag-names": []any{
			"JSONFieldName.Query",
			"URLFieldName.Query",
		},
		"allowed-tag-types": []any{
			"JSONFieldType.Exception",
			"URLFieldType.Exception",
		},
	})
	analyzers, _ := plugin.BuildAnalyzers()
	analysistest.Run(t, testdata, analyzers[0], "has-warnings", "no-warnings")
}
