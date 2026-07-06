package cmd

import (
	"github.com/spf13/cobra"
)

func activateCloud() {
	
}

var offCmd = &cobra.Command{
	Use:   "on",
	Short: "turn on self hosting",
	Long:  "turn on self hosting",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	offCmd.Flags().String("image", "", "docker image name to run")
}
