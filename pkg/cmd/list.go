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
		ctx, cancel := context.WithTimeout(context.Background(), checkCountFlags(cmd, "timeout"))
		defer cancel()

		res, err := Client.ReadDir(ctx, checkStringFlags(cmd, "remote-dir"), checkBoolFlags(cmd, "recursive"))
		if err != nil {
			log.Fatal(err)
		}
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
