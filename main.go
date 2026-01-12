package main

import (
	"github.com/xyue92/gitai/cmd"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {
	// Set version in cmd package
	cmd.SetVersion(Version)
	cmd.Execute()
}
