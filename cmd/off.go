package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func activateCloud() {

}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "turn off self hosting",
	Long:  "turn off self hosting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("off command")
	},
}

func init() {

}
