package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		// 对齐路径
		//rootPath := "/"
		// 列出所有本地文件
		localItems, err := webdav.LocalFileSystem("/").ReadDir(ctx, vars.localDir, vars.recursive)
		if err != nil {
			log.Fatal(err)
		}
		ltMap := make(map[string]struct{})
		for _, v := range localItems {
			p := strings.Replace(v.Path, vars.localDir, "/", 1)
			ltMap[p] = struct{}{}
		}
		fmt.Printf("lt文件：%s\n", ltMap)
		// 列出所有远程路径
		remoteItems := newReadDir(ctx, vars.Client, vars.remoteDir, vars.recursive)
		rtMap := make(map[string]struct{})
		for _, v := range remoteItems {
			p := strings.Replace(v.Path, vars.remoteDir, "/", 1)
			rtMap[p] = struct{}{}
		}
		fmt.Printf("rt文件：%s\n", rtMap)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
