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
	Short: "コマンドのバージョンを表示します",
	Long:  "コマンドのバージョンを表示します",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dummy Data Generator v0.0.1")
	},
}
