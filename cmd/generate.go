/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ponyo877/dummy_data_generator/internal/database"
	"github.com/ponyo877/dummy_data_generator/internal/repository"
	"github.com/ponyo877/dummy_data_generator/internal/usecase/generator"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "試験データを生成します",
	Long:  "試験データを生成します",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := database.PostgresClient()
		if err != nil {
			log.Fatalf("DBクライアントの作成に失敗しました: %v\n", err)
		}
		repository := repository.NewGenerateRepository(client)
		service := generator.NewService(repository)
		if err := service.Generate(); err != nil {
			log.Fatalf("試験データの作成に失敗しました: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}