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
		for _, v := range readLDir(ctx, vars.localDir, vars.recursive) {
			uploadFunc(ctx, vars.localDir, vars.remoteDir, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func uploadFunc(ctx context.Context, ld, rd string, v webdav.FileInfo) {
	_, subPath, _ := strings.Cut(v.Path, ld)
	switch v.IsDir {
	case true:
		rItemPath := rd + subPath + string(os.PathSeparator)
		if checkRemoteIsNotExist(ctx, rItemPath) {
			makeRemoteDir(ctx, rItemPath)
		}
	case false:
		rItemPath := rd + subPath
		if checkRemoteIsNotExist(ctx, rItemPath) || vars.overwrite {
			uploadFile(ctx, rItemPath, v.Path)
		}
	}
}

func uploadFile(ctx context.Context, rItemPath, lItemPath string) {
	localFile, err := os.Open(lItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Open Local File Err:%w", err))
	}

	r, err := vars.Client.Create(ctx, rItemPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Create Remote File Err:%w", err))
	}
	defer r.Close()

	_, err = io.Copy(r, localFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:Write Remote File Err:%w", err))
	}
}

func readLDir(ctx context.Context, path string, recurse bool) []webdav.FileInfo {
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
	if err != nil {
		log.Fatal(fmt.Errorf("Upload:List Remote Item Err:%w", err))
	}
	return fi
}
