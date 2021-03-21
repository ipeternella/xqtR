package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version ",
		Short: "Prints xqtR's current version.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("xqtr's current version: 1.0.0")
		},
	}

	return versionCmd
}
