package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

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
		basePath := vars.remoteDir + splitStr(vars.localDir)
		for _, v := range items {
			_, subPath, _ := strings.Cut(v.Path, splitStr(basePath))
			switch v.IsDir {
			case true:
				dirPath := basePath + subPath + "/"
				if !checkRemoteIsNotExist(ctx, dirPath) {
					continue
				}
				if err := vars.Client.Mkdir(ctx, dirPath); err != nil {
					log.Fatal(fmt.Errorf("make remote dir err:%s", err))
				}
			case false:
				filePath := basePath + subPath
				if checkRemoteIsNotExist(ctx, filePath) || vars.overwrite {
					uploadFile(ctx, filePath, v.Path)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func splitStr(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	str := strings.Split(path, "/")
	return str[len(str)-2]
}

func uploadFile(ctx context.Context, path, name string) {
	localFile, err := webdav.LocalFileSystem("/").Open(ctx, name)
	if err != nil {
		log.Fatal(fmt.Errorf("open local file err:%w", err))
	}
	data, err := io.ReadAll(localFile)
	if err != nil {
		log.Fatal(fmt.Errorf("read local file err:%w", err))
	}
	r, err := vars.Client.Create(ctx, path)
	if err != nil {
		log.Fatal(fmt.Errorf("create remote file err:%w", err))
	}
	_, err = r.Write(data)
	if err != nil {
		log.Fatal(fmt.Errorf("write remote file err:%w", err))
	}
	defer r.Close()
}
