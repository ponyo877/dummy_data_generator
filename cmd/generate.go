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
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate dummy data",
	Long:  "generate dummy data",
	Run: func(_ *cobra.Command, _ []string) {
		engine := viper.GetString("engine")
		client, err := database.NewDatabaseClient(engine)
		if err != nil {
			log.Fatalf("failed to create database client: %v\n", err)
		}
		repository := repository.NewGenerateRepository(client)
		service := generator.NewService(repository)
		if err := service.Generate(); err != nil {
			log.Fatalf("failed to generate dummy data: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
