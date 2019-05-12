package cmd

import (
	"../lib"
	"fmt"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pulls Pingdom check into a file",
	Run: func(cmd *cobra.Command, args []string) {
		checks, err := ListFunctions(cmd, args)
		err = lib.DownloadCheck(checks)
		if err != nil {
			fmt.Println("encountered an error downloading pingdom checks", err)
		}
		fmt.Printf("Pull complete, %d checks were downloaded\n",len(checks))
	},
}

func init() {
	PingdomCmd.AddCommand(PullCmd)
	PullCmd.Flags().IntP("ID", "i", 0, "Pingdom check ID to download")
	PullCmd.Flags().StringP("filter", "f", "", "Part of check name to match against for download")
}
