package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file from webdav",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()
		downloadFunc(ctx, vars.localDir, vars.remoteDir)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func downloadFunc(ctx context.Context, ld, rd string) {
	res := readRDir(ctx, vars.Client, rd, vars.recursive)
	for _, v := range res {
		lItemPath := fmt.Sprintf("%s%s", ld, v.Path)
		switch v.IsDir {
		case true:
			if checkLocalIsNotExist(ctx, lItemPath) {
				makeLocalDir(ctx, lItemPath)
			}
		case false:
			if checkLocalIsNotExist(ctx, lItemPath) || vars.overwrite {
				downloadFile(ctx, lItemPath, v.Path)
			}
		}
	}
}

func downloadFile(ctx context.Context, lItemPath, rItemPath string) {
	file, err := vars.Client.Open(ctx, rItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Open Remote File Err:%w", err))
	}
	defer file.Close()

	osFile, err := os.Create(lItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Create local File Err:%w", err))
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, file)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Write Local File Err:%w", err))
	}
}

func readRDir(ctx context.Context, c *webdav.Client, path string, recurse bool) []webdav.FileInfo {
	items, err := c.ReadDir(ctx, path, false)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:List Local Item Err:%w", err))
	}
	if recurse {
		var res []webdav.FileInfo
		res = append(res, items[0])
		for _, v := range items {
			if !v.IsDir {
				res = append(res, v)
				continue
			}
			if v.Path != path {
				res = append(res, readRDir(ctx, c, v.Path, recurse)...)
			}
		}
		return res
	}
	return items
}
