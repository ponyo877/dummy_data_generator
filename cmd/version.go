package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show command version",
	Long:  "show command version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Dummy Data Generator v0.0.1")
	},
}
