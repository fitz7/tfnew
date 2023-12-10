//go:build tools
// +build tools

package tools

import (
	_ "github.com/go-semantic-release/semantic-release/v2/cmd/semantic-release"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/spf13/cobra-cli"
	_ "gotest.tools/gotestsum"
)
