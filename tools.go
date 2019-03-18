// +build tools

package tools

import (
	// used by continuous integration scripts
	_ "github.com/mattn/goveralls"
	_ "golang.org/x/tools/cmd/cover"
)
