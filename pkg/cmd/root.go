package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	wb "webdav-cli/pkg/webdav"

	"github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

type globalVar struct {
	Client                                        *webdav.Client
	usr, pwd, server, remoteDir, localDir, config string
	overwrite, recursive                          bool
	timeout                                       time.Duration
}

var vars globalVar

var rootCmd = &cobra.Command{
	Use:   "webdav-cli",
	Short: "",
	Long:  "webdav-cli is a cli tools to sync to webdav server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		vars.Client = wb.InitClient(&http.Client{}, vars.server, vars.usr, vars.pwd)
	},
}

func Exec() {
	rootCmd.PersistentFlags().StringVarP(&vars.server, "server", "s", "", "webdav host ip")
	rootCmd.PersistentFlags().StringVarP(&vars.localDir, "local-dir", "l", "", "local directories that need to be synchronized")
	rootCmd.PersistentFlags().StringVarP(&vars.remoteDir, "remote-dir", "r", "", "remote server sync directory")
	// rootCmd.PersistentFlags().StringVarP(&vars.config, "config-file", "c", "", "read config from yaml file")
	rootCmd.PersistentFlags().StringVarP(&vars.usr, "username", "u", "", "username of logon webdav server")
	rootCmd.PersistentFlags().StringVarP(&vars.pwd, "password", "p", "", "password of logon webdav server")
	rootCmd.PersistentFlags().DurationVarP(&vars.timeout, "timeout", "t", 30*time.Second, "timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&vars.overwrite, "overwrite", false, "ignore the files with the same name between localdir and remotedir")
	rootCmd.PersistentFlags().BoolVar(&vars.recursive, "recursive", false, "recursively all directory files")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func checkLocalIsNotExist(ctx context.Context, name string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		_, err := os.Stat(name)
		if err != nil {
			log.Println(fmt.Errorf("Root:Local Item Is Not Exist:%w", err))
		}
		return !os.IsNotExist(err)
	}
}

func makeLocalDir(ctx context.Context, path string) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := os.MkdirAll(path, 0666); err != nil {
			log.Fatal(fmt.Errorf("Root:Make Local Dir Err:%w", err))
		}
	}
}

func checkRemoteIsNotExist(ctx context.Context, name string) bool {
	if _, err := vars.Client.Stat(ctx, name); err != nil {
		log.Printf("Root:%s Is Not Exist,Creating\n", name)
		return true
	}
	return false
}

func makeRemoteDir(ctx context.Context, path string) {
	if err := vars.Client.Mkdir(ctx, path); err != nil {
		log.Fatal(fmt.Errorf("Root:Make Remote Dir Err:%w", err))
	}
}
