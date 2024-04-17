package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webdav-cli",
	Short: "",
	Long:  "webdav-cli is a cli tools to sync to webdav server",
}

type globalVar struct {
	server    string
	localDir  string
	remoteDir string
}

var vars globalVar

func Exec() {
	rootCmd.PersistentFlags().StringVarP(&vars.server, "server", "s", "", "webdav host ip")
	rootCmd.PersistentFlags().StringVarP(&vars.localDir, "local-dir", "t", "", "Local directories that need to be synchronized")
	rootCmd.PersistentFlags().StringVarP(&vars.remoteDir, "remote-dir", "l", "", "webdav server sync directory")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		return
	}
}
