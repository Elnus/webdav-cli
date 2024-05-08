package cmd

import (
	"context"
	"os"
	"strings"
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

		var readWg sync.WaitGroup
		readWg.Add(2)
		// 列出所有本地文件
		ltMap := make(map[string]webdav.FileInfo)
		go func() {
			for _, v := range readLDir(ctx, vars.localDir, vars.recursive) {
				_, rootPath, _ := strings.Cut(v.Path, vars.localDir)
				if v.IsDir {
					rootPath = rootPath + string(os.PathSeparator)
				}
				ltMap[rootPath] = v
			}
			readWg.Done()
		}()

		// 列出所有远程路径
		rtMap := make(map[string]webdav.FileInfo)
		go func() {
			for _, v := range readRDir(ctx, vars.Client, vars.remoteDir, vars.recursive) {
				_, rootPath, _ := strings.Cut(v.Path, vars.remoteDir)
				rtMap[rootPath] = v
			}
			readWg.Done()
		}()
		readWg.Wait()

		var checkWg sync.WaitGroup
		checkWg.Add(2)
		downList := make(map[string]string)
		upList := make(map[string]string)

		go func() {
			for i, v := range ltMap {
				if value, exists := rtMap[i]; !exists {
					if !v.IsDir {
						upList[(vars.remoteDir + i)] = v.Path
					}
				} else {
					if !v.IsDir && vars.overwrite {
						if v.ModTime.After(value.ModTime) {
							upList[(vars.remoteDir + i)] = v.Path
						}
						if v.ModTime.Before(value.ModTime) {
							downList[v.Path] = (vars.remoteDir + i)
						}
					}
				}
			}
			checkWg.Done()
		}()
		go func() {
			for i, v := range rtMap {
				if _, exists := ltMap[i]; !exists {
					if v.IsDir {
						makeLocalDir(ctx, (vars.localDir + i))
						continue
					}
					if !v.IsDir {
						downList[vars.localDir+i] = v.Path
					}
				}
			}
			checkWg.Done()
		}()
		checkWg.Wait()

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
