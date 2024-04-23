package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

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
	res := newReadDir(ctx, vars.Client, rd, vars.recursive)
	for _, v := range res {
		path := fmt.Sprintf("%s%s", ld, v.Path)
		switch v.IsDir {
		case true:
			if checkLocalIsNotExist(ctx, path) {
				makeLocalDir(ctx, path)
			}
		case false:
			if checkLocalIsNotExist(ctx, path) || vars.overwrite {
				downloadFile(ctx, path, v.Path)
			}
		}
	}
}

func downloadFile(ctx context.Context, path, name string) {
	file, err := vars.Client.Open(ctx, name)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Open Remote File Err:%w", err))
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Read Remote File Err:%w", err))
	}
	osFile, err := webdav.LocalFileSystem("/").Create(ctx, path)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Create local File Err:%w", err))
	}
	_, err = osFile.Write(data)
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:Write Local File Err:%w", err))
	}
	defer osFile.Close()
}

func newReadDir(ctx context.Context, c *webdav.Client, path string, recurse bool) []webdav.FileInfo {
	var res []webdav.FileInfo
	items, err := c.ReadDir(ctx, path, false)
	res = append(res, items[0])
	if err != nil {
		log.Fatal(fmt.Errorf("DownLoad:List Local Item Err:%w", err))
	}
	if recurse {
		for _, v := range items {
			if !v.IsDir {
				res = append(res, v)
				continue
			}
			if v.IsDir && v.Path != path {
				res = append(res, newReadDir(ctx, c, v.Path, recurse)...)
			}
		}
		//fmt.Printf("res:%v\n", res)
		return res
	}
	return items
}
