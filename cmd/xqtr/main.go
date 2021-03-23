// Package main is, naturally, the main-entry point of the application which just creates a new
// `xqtR` command instance. The core `xqtR` application (which parses yaml jobs and executes them)
// is created and called only by the `xqtr run` subcommand as the `xqtr` root command only runs help.
package main

import (
	"github.com/IgooorGP/xqtR/pkg/cmd"
)

func main() {
	xqtR := cmd.NewXqtRCmd()

	xqtR.Execute()
}
