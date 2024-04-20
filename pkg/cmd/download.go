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
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		res, err := vars.Client.ReadDir(ctx, vars.remoteDir, vars.recursive)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range res {
			path := fmt.Sprintf("%s%s", vars.localDir, v.Path)
			switch v.IsDir {
			case true:
				if checkIsNotExist(path) {
					if os.Mkdir(path, 0644) != nil {
						log.Fatal(err)
					}
				}
			case false:
				if checkIsNotExist(path) || vars.overwrite {
					downloadFile(ctx, path, v.Path)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
