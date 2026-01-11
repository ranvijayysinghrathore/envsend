package main

import (
	"github.com/yourusername/envsend/cli/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cmd.Execute()
}
