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

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "cnt",
	Short: "count target table data",
	Long:  "count target table data",
	Run: func(_ *cobra.Command, _ []string) {
		engine := viper.GetString("engine")
		client, err := database.NewDatabaseClient(engine)
		if err != nil {
			log.Fatalf("failed to create database client: %v\n", err)
		}
		repository := repository.NewGenerateRepository(client)
		service := generator.NewService(repository)
		if err := service.Count(); err != nil {
			log.Fatalf("failed to count target table data: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}
