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

		// 列出所有本地文件
		localItems, err := webdav.LocalFileSystem("/").ReadDir(ctx, vars.localDir, vars.recursive)
		if err != nil {
			log.Fatal(err)
		}
		ltMap := make(map[string]webdav.FileInfo)
		for _, v := range localItems {
			var oldStr = vars.localDir
			if str := string(vars.localDir[len(vars.localDir)-1]); str == "/" {
				oldStr = vars.localDir[0 : len(vars.localDir)-1]
			}
			p := strings.Replace(v.Path, oldStr, "", 1)
			if v.IsDir {
				p = p + "/"
			}
			ltMap[p] = v
		}

		// 列出所有远程路径
		remoteItems := newReadDir(ctx, vars.Client, vars.remoteDir, vars.recursive)
		rtMap := make(map[string]webdav.FileInfo)
		for _, v := range remoteItems {
			p := strings.Replace(v.Path, vars.remoteDir, "/", 1)
			rtMap[p] = v

		}

		// for _, v := range remoteItems {
		// 	value, ok := ltMap[v.Path]
		// 	lItemPath := vars.localDir + v.Path
		// 	//fmt.Println(v.Path, ok)
		// 	if ok && !v.IsDir {
		// 		if rtMap[v.Path].After(value) {
		// 			fmt.Println("本地版本太久，下载", v.Path)
		// 			downloadFile(ctx, lItemPath, v.Path)
		// 		}
		// 		if rtMap[v.Path].Before(value) {
		// 			fmt.Println("远程版本太久，上传", v.Path)
		// 			// uploadFunc(ctx, v.Path, vars.remoteDir)
		// 		}
		// 	}
		// 	if !ok {
		// 		if v.IsDir && v.Path != "/" {
		// 			if checkLocalIsNotExist(ctx, lItemPath) {
		// 				fmt.Println("无路径，创建", lItemPath)
		// 				makeLocalDir(ctx, lItemPath)
		// 			}
		// 		}
		// 		if !v.IsDir {
		// 			fmt.Println("本地没有，下载", v.Path)
		// 			downloadFile(ctx, lItemPath, v.Path)
		// 		}
		// 	}
		// }

		// 遍历出本地有但远程没有的item
		for i, v := range ltMap {
			if _, exists := rtMap[i]; !exists {
				// 遍历出本地有但远程没有的item
				fmt.Println("本地有但远程没有的item:", i, v.IsDir)
			} else {
				// 遍历出本地和远程都有的item
				fmt.Println("本地和远程都有的item:", i, v.IsDir)
			}
		}

		// 遍历出远程有但本地没有的item
		for i, v := range rtMap {
			if _, exists := ltMap[i]; !exists {
				// 遍历出远程有但本地没有的item
				fmt.Println("远程有但本地没有的item:", i, v.IsDir)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
