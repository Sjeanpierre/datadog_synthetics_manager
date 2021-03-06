package cmd

import (
	"fmt"
	"log"

	"../lib"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local config to Datadog",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		singleCheck(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

func singleCheck(file string) {
	fc := lib.ReadFile(file)
	synth, err := lib.YAMLtoSynth(fc)
	if err != nil {
		log.Printf("encoutered error parsing file %s, error %s", file, err)
		return
	}
	if synth.PublicId != nil {
		//we are doing an update here
		check, err := lib.UpdateSyntheticsTest(synth)
		if err != nil {
			fmt.Printf("encountered issue updating Synthetics check, error %s", err)
		}
		fmt.Printf("%s | %s | %s\n", *check.PublicId, *check.Name, *check.Status)
		return
	}
	//we are creating a new check here
	check, err := lib.CreateSyntheticsTest(synth)
	if err != nil {
		fmt.Printf("encountered issue creating Synthetics check, error %s", err)
		return
	}
	fmt.Println("New Synthetic created from",file)
	fmt.Printf("%s | %s | %s\n", *check.PublicId, *check.Name, *check.Status)
}
