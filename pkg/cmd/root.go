// Package cmd configures all the cmds required by the xqtR app. The run command, created by newRunCmd(),
// creates the main xqtR which is used to parse yaml job files and execute them.
package cmd

import (
	"fmt"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/spf13/cobra"
)

const (
	appName      = "xqtr"
	shortAppDesc = "xqtR is a tool that lets you run jobs on your machine."
	longAppDesc  = "xqtR (short for executor) is a tool that parses jobs files in order to execute them to perform some tasks."
	versionText  = "xqtR current version: %s\n"
	version      = "1.0.0"
)

// NewXqtRCmd creates the `xqtr` command and its nested subcommands.
func NewXqtRCmd() *cobra.Command {
	xqtRConfig := config.NewXqtRConfigWithDefaults() // Main cfg of the command-line app

	// main cmd
	var rootCmd = &cobra.Command{
		Use:     appName,
		Short:   shortAppDesc,
		Long:    longAppDesc,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // main command just runs help
		},
	}
	rootCmd.SetVersionTemplate(fmt.Sprintf(versionText, version))

	// logging global flag
	rootCmd.PersistentFlags().StringVar(&xqtRConfig.LogLevel, "log-level", "info", "logging verbosity: ('debug', 'info', 'warn', 'error', 'fatal')")

	// cli commands wire up
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newRunCmd(&xqtRConfig)) // creates xqtR's core app instance to parse yaml jobs and exec them!

	return rootCmd
}
