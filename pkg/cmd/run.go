package cmd

import (
	"errors"
	"fmt"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/xqtr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	runCmdName      = "run"
	shortRunCmdDesc = "Runs the steps of a given job file."
)

// newRunCmdHandler is the actual handler of the `xqtr run` subcommand which calls the core application to execute
// the job yaml files
func newRunCmdHandler(cfg *config.XqtRConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		xqtR := xqtr.NewXqtR(cfg)
		xqtR.Run() // core application that parses job yaml files in order to exec them
	}
}

// newArgsValidatorHandler is a function used to perform some extra validations on the flags and arguments supplied
// to the `xqtr run` cobra command.
func newArgsValidatorHandler(cfg *config.XqtRConfig) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// checks for acceptable log level
		zeroLogLevel := config.ParseLevel(cfg.LogLevel)

		if zeroLogLevel == zerolog.Disabled {
			return errors.New(fmt.Sprintf("Unknown log level: %s", cfg.LogLevel))
		}

		return nil
	}
}

// newRunCmd creates a new `xqtr run` subcommand and add its flags and handlers. The core application `xqtR` is
// called by this run command. See `newRunCmdHandler()` function for details.
func newRunCmd(cfg *config.XqtRConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   runCmdName,
		Short: shortRunCmdDesc,
		Args:  newArgsValidatorHandler(cfg),
		Run:   newRunCmdHandler(cfg),
	}

	// persistent flags -> cascades to subcommands
	cmd.PersistentFlags().StringVarP(&cfg.JobFilePath, "file", "f", "job.yaml", "the job file yaml location with the steps to be executed")
	cmd.PersistentFlags().BoolVar(&cfg.IsDryRun, "dry-run", false, "runs the job's steps without actually executing anything")

	return cmd
}
