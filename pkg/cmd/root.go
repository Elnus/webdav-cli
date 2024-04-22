package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		vars.usr = checkStringFlags(cmd, "username")
		vars.pwd = checkStringFlags(cmd, "password")
		vars.server = checkStringFlags(cmd, "server")
		vars.remoteDir = checkStringFlags(cmd, "remote-dir")
		vars.localDir = checkStringFlags(cmd, "local-dir")
		vars.overwrite = checkBoolFlags(cmd, "overwrite")
		vars.recursive = checkBoolFlags(cmd, "recursive")
		vars.timeout = checkCountFlags(cmd, "timeout")
		vars.Client = wb.InitClient(&http.Client{}, vars.server, vars.usr, vars.pwd)
	},
}

func Exec() {
	rootCmd.PersistentFlags().StringVarP(&vars.server, "server", "s", "", "webdav host ip")
	rootCmd.PersistentFlags().StringVarP(&vars.localDir, "local-dir", "l", "", "local directories that need to be synchronized")
	rootCmd.PersistentFlags().StringVarP(&vars.remoteDir, "remote-dir", "r", "", "remote server sync directory")
	rootCmd.PersistentFlags().StringVarP(&vars.config, "config-file", "c", "", "read config from yaml file")
	rootCmd.PersistentFlags().StringVarP(&vars.usr, "username", "u", "", "username of logon webdav server")
	rootCmd.PersistentFlags().StringVarP(&vars.pwd, "password", "p", "", "password of logon webdav server")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 30*time.Second, "timeout in seconds")
	rootCmd.PersistentFlags().Bool("overwrite", false, "ignore the files with the same name between localdir and remotedir")
	rootCmd.PersistentFlags().Bool("recursive", false, "recursively all directory files")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func checkStringFlags(cmd *cobra.Command, arg string) string {
	res, err := cmd.Flags().GetString(arg)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return res
}

func checkBoolFlags(cmd *cobra.Command, arg string) bool {
	res, err := cmd.Flags().GetBool(arg)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return res
}

func checkCountFlags(cmd *cobra.Command, arg string) time.Duration {
	res, err := cmd.Flags().GetDuration(arg)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func checkLocalIsNotExist(ctx context.Context, name string) bool {
	if _, err := webdav.LocalFileSystem("/").Stat(ctx, name); err != nil {
		log.Println(fmt.Errorf("local stat:%w", err))
		return true
	}
	return false
}

func checkRemoteIsNotExist(ctx context.Context, name string) bool {
	if _, err := vars.Client.Stat(ctx, name); err != nil {
		log.Println(fmt.Errorf("remote stat:%w", err))
		return true
	}
	return false
}
