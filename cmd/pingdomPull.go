package cmd

import (
	"../lib"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pulls Pingdom check into a file",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checks, err := lib.GetPingdomCheck(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		err = lib.CheckDownload(checks)
		if err != nil {
			fmt.Println("encountered an error downloading pingdom checks",err)
		}
	},
}

func init() {
	PingdomCmd.AddCommand(PullCmd)
}