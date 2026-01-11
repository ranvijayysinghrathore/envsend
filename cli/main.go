package main

import (
	"github.com/ranvijayysinghrathore/envsend/cli/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cmd.Execute()
}
