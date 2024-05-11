package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "make webdav-cli server dir",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()
		err := vars.Client.Mkdir(ctx, vars.remoteDir)
		if err != nil {
			log.Printf("create dir %s failed\n", vars.remoteDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}
