package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show webdav-cli current version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("webcli version : %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
