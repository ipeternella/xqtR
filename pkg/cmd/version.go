package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version ",
		Short: "Prints xqtR's current version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(versionText, version)
		},
	}

	return versionCmd
}
