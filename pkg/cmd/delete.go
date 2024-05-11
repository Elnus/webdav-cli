package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete webdav server file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		err := vars.Client.RemoveAll(ctx, vars.remoteDir)
		if err != nil {
			log.Printf("delete item failed:%s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
