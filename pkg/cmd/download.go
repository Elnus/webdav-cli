package cmd

import (
	"context"

	wb "github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file from webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), checkCountFlags(cmd, "timeout"))
		defer cancel()
		Client.Move(ctx, checkStringFlags(cmd, "remote-dir"), "", &wb.MoveOptions{NoOverwrite: true})
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
