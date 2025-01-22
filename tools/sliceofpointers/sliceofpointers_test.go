package sliceofpointers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRun(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()
	plugin, _ := New(nil)
	analyzers, _ := plugin.BuildAnalyzers()
	analysistest.Run(t, testdata, analyzers[0], "has-warnings", "no-warnings")
}
