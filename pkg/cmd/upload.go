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
		uploadFunc(ctx, vars.localDir, vars.remoteDir)
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

func uploadFunc(ctx context.Context, ld, rd string) {
	items, err := webdav.LocalFileSystem("/").ReadDir(ctx, ld, vars.recursive)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:List Remote Item Err:%w", err))
	}
	basePath := rd + splitStr(ld)
	for _, v := range items {
		_, subPath, _ := strings.Cut(v.Path, splitStr(basePath))
		switch v.IsDir {
		case true:
			dirPath := basePath + subPath + "/"
			if !checkRemoteIsNotExist(ctx, dirPath) {
				continue
			}
			makeRemoteDir(ctx, dirPath)
		case false:
			filePath := basePath + subPath
			if checkRemoteIsNotExist(ctx, filePath) || vars.overwrite {
				uploadFile(ctx, filePath, v.Path)
			}
		}
	}
}

func uploadFile(ctx context.Context, path, name string) {
	localFile, err := webdav.LocalFileSystem("/").Open(ctx, name)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Open Local File Err:%w", err))
	}
	data, err := io.ReadAll(localFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Read Local File Err:%w", err))
	}
	r, err := vars.Client.Create(ctx, path)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Create Remote File Err:%w", err))
	}
	_, err = r.Write(data)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Write Remote File Err:%w", err))
	}
	defer r.Close()
}
