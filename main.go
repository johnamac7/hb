package main

import (
	"fmt"
	"time"

	"github.com/damianoneill/hbcli/cmd"
)

import (
	_ "github.com/damianoneill/hbcli/cmd/provision"
)


var (
	// Version should be set as an argument in the build
	Version = "NOT SET"
)

func main() {
	t := time.Now()
	Version := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	cmd.Execute(Version)
}
