package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version indication",
	Long:  `Version indication`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version used: %v\n", Version())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
