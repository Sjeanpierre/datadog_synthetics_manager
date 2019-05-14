package cmd

import (
	"../lib"
	"fmt"
	"github.com/spf13/cobra"
)

// ListCmd represents the listing of checks from Pingdom
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists pingdom checks",
	Run: func(cmd *cobra.Command, args []string) {
		checks, err := ListFunctions(cmd, args)
		if err != nil {
			fmt.Println(err)
		}
		js, err := cmd.Flags().GetBool("json")
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
	ListCmd.Flags().BoolP("json", "j", false, "return data in json format")
	ListCmd.Flags().IntP("ID", "i", 0, "Pingdom check ID to lookup")
	ListCmd.Flags().StringP("filter", "f", "", "Part of check name to match against")
}
