//go:build tools
// +build tools

package tools

/* These are various tools we put in here so they exist in vnedor direcotry (i think) */

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/segmentio/golines"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
