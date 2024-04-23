package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync localDir and RemoteDir",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		// 列出所有本地文件
		localItems, err := webdav.LocalFileSystem("/").ReadDir(ctx, vars.localDir, vars.recursive)
		if err != nil {
			log.Fatal(err)
		}
		ltMap := make(map[string]time.Time)
		for _, v := range localItems {
			if v.IsDir {
				v.Path = v.Path + "/"
			}
			p := strings.Replace(v.Path, vars.localDir, "/", 1)
			ltMap[p] = v.ModTime
		}

		// 列出所有远程路径
		remoteItems := newReadDir(ctx, vars.Client, vars.remoteDir, vars.recursive)
		rtMap := make(map[string]time.Time)
		for _, v := range remoteItems {
			p := strings.Replace(v.Path, vars.remoteDir, "/", 1)
			rtMap[p] = v.ModTime
		}

		for _, v := range remoteItems {
			value, ok := ltMap[v.Path]
			fmt.Println(v.Path, ok)
			if ok && !v.IsDir {
				if rtMap[v.Path].After(value) {
					fmt.Println("本地版本太久，下载", v.Path)
					downloadFile(ctx, vars.localDir, v.Path)
				}
				if rtMap[v.Path].Before(value) {
					fmt.Println("远程版本太久，上传", v.Path)
					// uploadFunc(ctx, v.Path, vars.remoteDir)
				}
			}
			if !ok && !v.IsDir {
				fmt.Println("本地没有，下载", v.Path)
				downloadFile(ctx, vars.localDir, v.Path)
			}
		}

		// for _, v := range localItems {
		// 	_, ok := rtMap[v.Path]
		// 	if !ok {
		// 		// uploadFunc(ctx, v.Path, vars.remoteDir)
		// 	}
		// }
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
