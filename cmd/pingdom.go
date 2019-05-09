package cmd

import (
	"github.com/spf13/cobra"
)

// PingdomCmd represents the Pingdom  root command, does not perform any actions
var PingdomCmd = &cobra.Command{
	Use:   "pingdom",
	Short: "checks Pingdom",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(PingdomCmd)
}