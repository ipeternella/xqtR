package cmd

import (
	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/spf13/cobra"
)

const (
	appName      = "xqtr"
	shortAppDesc = "xqtR is a tool that lets you run jobs on your machine."
	longAppDesc  = "xqtR (short for executor) is a tool that parses jobs files in order to execute them to perform some tasks."
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
	rootCmd.SetVersionTemplate("xqtr version is 1.0.0\n")

	// logging global flag
	rootCmd.PersistentFlags().StringVar(&xqtRConfig.LogLevel, "log-level", "info", "logging verbosity: ('debug', 'info', 'warn', 'error', 'fatal')")

	// cli commands wire up
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newRunCmd(&xqtRConfig)) // main xqtr command

	return rootCmd
}
