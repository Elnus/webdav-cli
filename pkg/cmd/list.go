package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files from directory",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		res := newReadDir(ctx, vars.Client, vars.remoteDir, vars.recursive)

		fmt.Println("----------------------------")
		for _, v := range res {
			fmt.Printf("Path:%s  |  ", v.Path)
			fmt.Printf("ModTime:%v  |  ", v.ModTime)
			fmt.Printf("IsDir:%v\n", v.IsDir)
			fmt.Println("----------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
