// Package main contains CLI documentation generator tool.
package main

import (
	_ "embed"
	"strings"

	"github.com/grafana/clireadme"
	"github.com/grafana/k6deps/cmd"
)

func main() {
	root := cmd.New()
	root.Use = strings.ReplaceAll(root.Use, "deps", "k6deps")
	clireadme.Main(root, 1)
}
