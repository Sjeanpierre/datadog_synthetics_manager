package cmd

import (
	"../lib"
	"fmt"
	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/spf13/cobra"
	"strings"
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

// ListFunctions are used by Pingdom commands to retrieve single or multiple Pingdom checks
func ListFunctions(cmd *cobra.Command, args []string)  ([]pingdom.CheckResponse,error){
	var checks []pingdom.CheckResponse
	id, err := cmd.Flags().GetInt("ID")
	if err != nil {
		return checks, fmt.Errorf("could not parse ID flag, err %s",err)
	}
	// Handle -i && --ID flag being used in list or pull command
	if id != 0 {
		checks, err = lib.GetPingdomCheck(id)
		if err != nil {
			return checks, err
		}
		return checks,nil
	}
	//Pull list of checks from Pingdom
	checks, err = lib.ListPingdomChecks()
	if err != nil {
		return checks,fmt.Errorf("could not list Pingdom checks, error %s",err)
	}
	//Handle -f && --filter flag being set in pull or Listing
	filter,err := cmd.Flags().GetString("filter")
	if filter != "" {
		var filteredCheckIDs []int
		for _, check := range checks {
			if strings.Contains(check.Name,filter) {
				filteredCheckIDs = append(filteredCheckIDs,check.ID)
			}
		}
		checks = lib.GetPingdomChecks(filteredCheckIDs)
	}
	return checks,nil
}