package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func activateCloud() {

}

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "switch to cloud hosting",
	Long:  "switch to cloud hosting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cloud command")
	},
}

func init() {

}
