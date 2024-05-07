package cmd

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Get Abs Path Err:%w", err))
	}
	str := strings.Split(absPath+string(os.PathSeparator), string(os.PathSeparator))
	return str[len(str)-2]
}

func uploadFunc(ctx context.Context, ld, rd string) {
	items, err := readLDir(ctx, ld, vars.recursive)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:List Remote Item Err:%w", err))
	}
	basePath := rd + splitStr(ld)
	for _, v := range items {
		_, subPath, _ := strings.Cut(v.Path, splitStr(basePath))
		switch v.IsDir {
		case true:
			dirPath := basePath + subPath + string(os.PathSeparator)
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

func uploadFile(ctx context.Context, rItemPath, lItemPath string) {
	localFile, err := os.Open(lItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Open Local File Err:%w", err))
	}
	data, err := io.ReadAll(localFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Read Local File Err:%w", err))
	}
	r, err := vars.Client.Create(ctx, rItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Create Remote File Err:%w", err))
	}
	_, err = r.Write(data)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Write Remote File Err:%w", err))
	}
	defer r.Close()
}

func readLDir(ctx context.Context, path string, recurse bool) ([]webdav.FileInfo, error) {
	var fi []webdav.FileInfo
	var depth int = 0
	if recurse {
		depth = 63
	}
	rootDepth := strings.Count(path, string(os.PathSeparator))
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if strings.Count(path, string(os.PathSeparator)) > depth+rootDepth {
			return filepath.SkipDir
		}
		fi = append(fi, webdav.FileInfo{Path: path, ModTime: info.ModTime(), IsDir: info.IsDir()})
		return err
	})
	return fi, err
}
