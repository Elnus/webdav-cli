package cmd

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	wb "webdav-cli/pkg/webdav"

	"github.com/emersion/go-webdav"
	"github.com/spf13/cobra"
)

var (
	Client *webdav.Client
)

var rootCmd = &cobra.Command{
	Use:   "webdav-cli",
	Short: "",
	Long:  "webdav-cli is a cli tools to sync to webdav server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		usr := checkStringFlags(cmd, "username")
		pwd := checkStringFlags(cmd, "password")
		ep := checkStringFlags(cmd, "server")
		Client = wb.InitClient(&http.Client{}, ep, usr, pwd)
	},
}

type globalVar struct {
	server    string
	localDir  string
	remoteDir string
	config    string
	username  string
	password  string
}

var vars globalVar

func Exec() {
	rootCmd.PersistentFlags().StringVarP(&vars.server, "server", "s", "", "webdav host ip")
	rootCmd.PersistentFlags().StringVarP(&vars.localDir, "local-dir", "l", "", "local directories that need to be synchronized")
	rootCmd.PersistentFlags().StringVarP(&vars.remoteDir, "remote-dir", "r", "", "remote server sync directory")
	rootCmd.PersistentFlags().StringVarP(&vars.config, "config-file", "c", "", "read config from yaml file")
	rootCmd.PersistentFlags().StringVarP(&vars.username, "username", "u", "", "username of logon webdav server")
	rootCmd.PersistentFlags().StringVarP(&vars.password, "password", "p", "", "password of logon webdav server")
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

func checkIsNotExist(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

func downloadFile(ctx context.Context, path, name string) {
	file, err := Client.Open(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	osFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = osFile.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
