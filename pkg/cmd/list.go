package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files from directory",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		remoteDir := checkStringFlags(cmd, "remote-dir")
		ignore := checkBoolFlags(cmd, "ignore-samename-file")
		res, err := Client.ReadDir(context.Background(), remoteDir, ignore)
		if err != nil {
			log.Panicln(err)
			return
		}
		for _, v := range res {
			fmt.Println("----------------------------")
			fmt.Printf("Path:%s  |  ", v.Path)
			fmt.Printf("ModTime:%v  |  ", v.ModTime)
			fmt.Printf("IsDir:%v\n", v.IsDir)
		}
	},
}

func init() {
	listCmd.Flags().Bool("recursive", false, "Recursively list all directory files")
	rootCmd.AddCommand(listCmd)
}
