package cmd

import (
	"context"
	"os"
	"runtime"
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
			defer readWg.Done()
			for _, v := range readLDir(ctx, vars.localDir, vars.recursive) {
				_, rootPath, _ := strings.Cut(v.Path, vars.localDir)
				if v.IsDir {
					rootPath = rootPath + string(os.PathSeparator)
				}
				ltMap[rootPath] = v
			}
		}()

		// 列出所有远程路径
		rtMap := make(map[string]webdav.FileInfo)
		go func() {
			defer readWg.Done()
			for _, v := range readRDir(ctx, vars.Client, vars.remoteDir, vars.recursive) {
				_, rootPath, _ := strings.Cut(v.Path, vars.remoteDir)
				rtMap[rootPath] = v
			}
		}()
		readWg.Wait()

		var checkWg sync.WaitGroup
		checkWg.Add(2)
		downloadChan := make(chan [2]string, 10)
		uploadChan := make(chan [2]string, 10)

		go func() {
			defer checkWg.Done()
			for i, v := range ltMap {
				if value, exists := rtMap[i]; !exists {
					if !v.IsDir {
						uploadChan <- [2]string{vars.remoteDir + i, v.Path}
					}
				} else {
					if !v.IsDir && vars.overwrite {
						if v.ModTime.After(value.ModTime) {
							uploadChan <- [2]string{vars.remoteDir + i, v.Path}
						}
						if v.ModTime.Before(value.ModTime) {
							downloadChan <- [2]string{v.Path, vars.remoteDir + i}
						}
					}
				}
			}
		}()
		go func() {
			defer checkWg.Done()
			for i, v := range rtMap {
				if _, exists := ltMap[i]; !exists {
					if v.IsDir {
						makeLocalDir(ctx, (vars.localDir + i))
					} else {
						downloadChan <- [2]string{vars.localDir + i, v.Path}
					}
				}
			}
		}()
		checkWg.Wait()
		close(uploadChan)
		close(downloadChan)

		var actWg sync.WaitGroup
		actWg.Add(2 * runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				defer actWg.Done()
				for {
					if v, ok := <-downloadChan; ok {
						downloadFile(ctx, v[0], v[1])
					} else {
						break
					}
				}
			}()
		}
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				defer actWg.Done()
				for {
					if v, ok := <-uploadChan; ok {
						uploadFile(ctx, v[0], v[1])
					} else {
						break
					}
				}
			}()
		}
		actWg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
