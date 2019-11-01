package main

import (
	"github.com/damianoneill/hb/cmd"

	_ "github.com/damianoneill/hb/cmd/provision"
)

var (
	// Version should be set as an argument in the build
	Version = "1.0.0"
)

func main() {
	cmd.Execute(Version)
}
