package cmd

import (
	"errors"
	"fmt"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newRunCmd creates a new run subcommand and add its flags
func newRunCmd(cfg *config.XqtRConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Installs all software found in an Minstall file.",
		Args: func(cmd *cobra.Command, args []string) error {
			// checks for acceptable log level
			_, okLogLevel := config.StrToLogLevelMapping[cfg.LogLevel]

			if !okLogLevel {
				return errors.New(fmt.Sprintf("Unknown log level: %s", cfg.LogLevel))
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Running...")
		},
	}

	// persistent flags -> cascades to subcommands
	cmd.PersistentFlags().StringVarP(&cfg.JobFilePath, "file", "f", "job.yml", "The install file path")
	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "Logging verbosity: ('debug', 'info', 'warn', 'error', 'fatal')")
	cmd.PersistentFlags().BoolVar(&cfg.IsDryRun, "dry-run", false, "Runs the program without installing or downloading anything")

	return cmd
}
