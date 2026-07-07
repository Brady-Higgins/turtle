package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "turtle",
	Short: "The Turtle CLI",
	Long:  "The Turtle CLI is hybrid cloud tool",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
}
