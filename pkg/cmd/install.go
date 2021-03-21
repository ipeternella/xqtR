package cmd

import (
	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newRunCmd creates a new run subcommand and add its flags
func newRunCmd(cfg *config.XqtRConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Installs all software found in an Minstall file.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Running...")
		},
	}

	cmd.PersistentFlags().StringVarP(&cfg.JobFilePath, "file", "f", "minstall.yml", "The install file path")
	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "Logging verbosity")
	cmd.PersistentFlags().BoolVar(&cfg.IsDryRun, "dry-run", false, "Runs the program without installing or downloading anything")

	return cmd
}
