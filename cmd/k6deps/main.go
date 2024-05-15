// Package main contains the main function for k6deps CLI tool.
package main

import "fmt"

//nolint:gochecknoglobals
var (
	appname = "k6pack"
	version = "dev"
)

func main() {
	fmt.Printf("%s version %s\n", appname, version)
}
