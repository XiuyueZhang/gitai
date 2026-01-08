package main

import (
	"github.com/xyue92/gitai/cmd"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {
	cmd.Execute()
}
