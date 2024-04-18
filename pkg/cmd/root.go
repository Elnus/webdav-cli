package cmd

import (
	"log"
	"net/http"
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
	rootCmd.PersistentFlags().StringVarP(&vars.password, "password", "p", "", "username of logon webdav server")
	rootCmd.PersistentFlags().Bool("ignore-samename-file", false, "ignore the files with the same name between localdir and remotedir")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func checkStringFlags(cmd *cobra.Command, arg string) string {
	res, err := cmd.Flags().GetString(arg)
	if err != nil {
		log.Println(err)
		return ""
	}
	return res
}

func checkBoolFlags(cmd *cobra.Command, arg string) bool {
	res, err := cmd.Flags().GetBool(arg)
	if err != nil {
		log.Println(err)
		return false
	}
	return res
}
