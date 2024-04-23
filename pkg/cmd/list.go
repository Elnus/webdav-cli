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
		ctx, cancel := context.WithTimeout(context.Background(), vars.timeout)
		defer cancel()

		res, err := vars.Client.ReadDir(ctx, vars.remoteDir, vars.recursive)
		if err != nil {
			log.Fatal(fmt.Errorf("List:List Remote Item Err:%w", err))
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
