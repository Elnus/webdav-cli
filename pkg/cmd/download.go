package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file from webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		remoteDir := checkStringFlags(cmd, "remote-dir")
		localDir := checkStringFlags(cmd, "local-dir")
		ctx, cancel := context.WithTimeout(context.Background(), checkCountFlags(cmd, "timeout"))
		defer cancel()

		res, err := Client.ReadDir(ctx, remoteDir, checkBoolFlags(cmd, "recursive"))
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range res {
			path := fmt.Sprintf("%s%s", localDir, v.Path)
			switch v.IsDir {
			case true:
				if checkIsNotExist(path) {
					if os.Mkdir(path, 0644) != nil {
						log.Fatal(err)
					}
				}
			case false:
				if checkIsNotExist(path) || checkBoolFlags(cmd, "overwrite") {
					downloadFile(ctx, path, v.Path)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
