package cmd

import (
	"fmt"
	"log"

<<<<<<< HEAD
	"../lib"
=======
	"github.com/sjeanpierre/datadog_synthetics_manager/lib"
>>>>>>> edbe8078564d4d452eb18b9a16ff1d74074d1bca
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
	log.Println("NOT YET IMPLEMENTED")
}
