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
		var res []webdav.FileInfo
		switch vars.recursive {
		case true:
			res = RecurseReadAndMkdir(ctx, vars.Client, vars.remoteDir, res)
		case false:
			f, err := vars.Client.ReadDir(ctx, vars.remoteDir, vars.recursive)
			res = f
			if err != nil {
				log.Fatal(fmt.Errorf("read remote dir err:%w", err))
			}
		}

		for _, v := range res {
			path := fmt.Sprintf("%s%s", vars.localDir, v.Path)
			switch v.IsDir {
			case true:
				if checkLocalIsNotExist(ctx, path) {
					if err := webdav.LocalFileSystem("/").Mkdir(ctx, path); err != nil {
						log.Fatal(fmt.Errorf("make local dir err:%w", err))
					}
				}
			case false:
				if checkLocalIsNotExist(ctx, path) || vars.overwrite {
					downloadFile(ctx, path, v.Path)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func downloadFile(ctx context.Context, path, name string) {
	file, err := vars.Client.Open(ctx, name)
	if err != nil {
		log.Fatal(fmt.Errorf("open remote file err:%w", err))
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(fmt.Errorf("read remote file err:%w", err))
	}
	fmt.Println("path", path)
	osFile, err := webdav.LocalFileSystem("/").Create(ctx, path)
	if err != nil {
		log.Fatal(fmt.Errorf("create local file err:%w", err))
	}
	_, err = osFile.Write(data)
	if err != nil {
		log.Fatal(fmt.Errorf("write local file err:%w", err))
	}
	defer osFile.Close()
}

func RecurseReadAndMkdir(ctx context.Context, c *webdav.Client, path string, res []webdav.FileInfo) []webdav.FileInfo {
	f, err := c.ReadDir(ctx, path, false)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range f {
		if v.IsDir && v.Path != path {
			res = append(res, v)
			RecurseReadAndMkdir(ctx, c, v.Path, res)
			continue
		}
		res = append(res, v)
	}
	fmt.Println(res)
	return res
}
