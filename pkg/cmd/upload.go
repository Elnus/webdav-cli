package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload file to webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("is upload")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
