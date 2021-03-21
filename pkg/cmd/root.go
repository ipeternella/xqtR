package cmd

import (
	"github.com/spf13/cobra"
)

// NewXqtRCmd creates the `xqtr` command and its nested subcommands.
func NewXqtRCmd() *cobra.Command {

	// main cmd
	var rootCmd = &cobra.Command{
		Use:     "xqtr",
		Short:   "xqtr is a tool that lets you run jobs on your machine.",
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // main command just runs help
		},
	}
	rootCmd.SetVersionTemplate("xqtr version is 1.0.0\n")

	// cli commands wire up
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
