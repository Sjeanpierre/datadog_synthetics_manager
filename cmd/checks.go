package cmd

import (
	"../lib"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var checksCmd = &cobra.Command{
	Use:   "checks",
	Short: "List Synthetic checks in Account",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			lib.GetSyntheticsCheck(args[0])
		} else {
			lib.ListSyntheticsChecks()
		}
	},
}

func init() {
	rootCmd.AddCommand(checksCmd)
}
