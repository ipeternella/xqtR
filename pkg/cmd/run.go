package cmd

import (
	"errors"
	"fmt"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/startup"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	runCmdName      = "run"
	shortRunCmdDesc = "Runs the steps of a given job file."
)

// newRunCmd creates a new run subcommand and add its flags
func newRunCmd(cfg *config.XqtRConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   runCmdName,
		Short: shortRunCmdDesc,
		Args:  newArgsValidatorHandler(cfg),
		Run:   newRunCmdHandler(cfg),
	}

	// persistent flags -> cascades to subcommands
	cmd.PersistentFlags().StringVarP(&cfg.JobFilePath, "file", "f", "job.yml", "the job file yaml location with the steps to be executed")
	cmd.PersistentFlags().BoolVar(&cfg.IsDryRun, "dry-run", false, "runs the job's steps without actually executing anything")

	return cmd
}

func newRunCmdHandler(cfg *config.XqtRConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		startup.Boot(cfg)
		log.Info().Msg("Running...")
	}
}

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
