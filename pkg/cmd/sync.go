package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"

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
			rootPath := unifiedPath(vars.localDir, v.Path)
			if v.IsDir {
				rootPath = rootPath + "/"
			}
			// fmt.Println("lt:", v.Path)
			ltMap[rootPath] = v
		}

		// 列出所有远程路径
		remoteItems := newReadDir(ctx, vars.Client, vars.remoteDir, vars.recursive)
		rtMap := make(map[string]webdav.FileInfo)
		for _, v := range remoteItems {
			rootPath := unifiedPath(vars.remoteDir, v.Path)
			fmt.Println("rt:", v.Path)
			rtMap[rootPath] = v
		}
		for i := range rtMap {
			fmt.Println(i)
		}

		var checkWg sync.WaitGroup
		checkWg.Add(2)
		downList := make(map[string]string)
		upList := make(map[string]string)

		go func() {
			for i, v := range ltMap {
				if value, exists := rtMap[i]; !exists {
					if !v.IsDir {
						//fmt.Println("上传本地文件:", v.Path)
						upList[(vars.remoteDir + i)] = v.Path
						//uploadFile(ctx, (vars.remoteDir + i), v.Path)
					}
				} else {
					// 遍历出本地和远程都有的item
					fmt.Println("本地和远程都有的item:", i, v.IsDir)
					if !v.IsDir && vars.overwrite {
						if v.ModTime.After(value.ModTime) {
							//fmt.Println("远程版本太旧，上传", v.Path)
							upList[(vars.remoteDir + i)] = v.Path
							//uploadFile(ctx, (vars.remoteDir + i), v.Path)
						}
						if v.ModTime.Before(value.ModTime) {
							//fmt.Println("本地版本太旧，下载", v.Path)
							downList[v.Path] = (vars.remoteDir + i)
							//downloadFile(ctx, v.Path, (vars.remoteDir + i))
						}
					}
				}
			}
			checkWg.Done()
		}()
		// 遍历出本地有但远程没有的item
		go func() {
			// 遍历出远程有但本地没有的item
			for i, v := range rtMap {
				if _, exists := ltMap[i]; !exists {
					// 遍历出远程有但本地没有的item
					fmt.Printf("远程有但本地没有的item,i:%v  v.IsDir:%v  v.Patt:%v\n", i, v.IsDir, v.Path)
					if v.IsDir {
						fmt.Println("创建本地文件夹:", (vars.localDir + i))
						makeLocalDir(ctx, (vars.localDir + i))
						continue
					}
					if !v.IsDir {
						downList[v.Path] = (vars.remoteDir + i)
						//downloadFile(ctx, (vars.localDir + i), v.Path)
					}
				}
			}
			checkWg.Done()
		}()
		checkWg.Wait()

		fmt.Println("Check Is Done")
		var actWg sync.WaitGroup
		actWg.Add(2)
		go func() {
			for i, v := range downList {
				downloadFile(ctx, i, v)
			}
			actWg.Done()
		}()
		go func() {
			for i, v := range upList {
				uploadFile(ctx, i, v)
			}
			actWg.Done()
		}()
		actWg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
