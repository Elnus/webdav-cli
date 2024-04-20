package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		items, err := webdav.LocalFileSystem("/").ReadDir(ctx, vars.localDir, vars.recursive)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range items {
			switch v.IsDir {
			case true:
				fmt.Println("this is a fold")
			case false:
				localFile, err := webdav.LocalFileSystem("/").Open(ctx, v.Path)
				if err != nil {
					log.Fatal(err)
				}
				data, err := io.ReadAll(localFile)
				if err != nil {
					log.Fatal(err)
				}
				_, err = vars.Client.Stat(ctx, v.Path)
				if err != nil {
					log.Fatal(err)
				}
				r, err := vars.Client.Create(ctx, v.Path)
				if err != nil {
					log.Fatal(err)
				}
				defer r.Close()
				r.Write(data)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
