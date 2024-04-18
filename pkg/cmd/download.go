package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file from webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("is download")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
