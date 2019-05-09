package cmd

import (
	"../lib"
	"fmt"
	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/spf13/cobra"
	"log"
)



// ListCmd represents the listing of checks from Pingdom
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists pingdom checks",
	Run: func(cmd *cobra.Command, args []string) {
		var checks []pingdom.CheckResponse
		var err error
		id, err := cmd.Flags().GetString("ID")
		if id != "" {
			checks, err = lib.GetPingdomCheck(id)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			checks, err = lib.ListPingdomChecks()
			if err != nil {
				log.Println(err)
				return
			}
		}
		js,err := cmd.Flags().GetBool("json")
		if err == nil && js {
			lib.CheckJson(checks)
			return
		}

		for _, check := range checks {
			fmt.Printf("%d | %s | %s\n", check.ID, check.Name, check.Status)
		}
	},
}



func init() {
	PingdomCmd.AddCommand(ListCmd)
	ListCmd.Flags().BoolP("json","j",false,"return data in json format")
	ListCmd.Flags().StringP("ID","i","","Pingdom check ID to lookup")
}
