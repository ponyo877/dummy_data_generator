/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringSliceP("config", "c", []string{"config.yaml"}, "dummy data config file, multi case: -c \"cfg_*.yaml\" or -c cfg_1.yaml,cfg_2.yaml")
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	rootCmd.PersistentFlags().StringP("engine", "e", "postgres", "database engine, support postgres and mysql")
	rootCmd.PersistentFlags().StringP("host", "h", "127.0.0.1", "database server host or socket directory")
	rootCmd.PersistentFlags().StringP("dbuser", "u", "root", "database user name")
	rootCmd.PersistentFlags().StringP("port", "P", "5432", "database server port")
	rootCmd.PersistentFlags().StringP("password", "p", "password", "database password to use when connecting to serve")
	rootCmd.PersistentFlags().StringP("database", "d", "mydb", "the database to use")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("engine", rootCmd.PersistentFlags().Lookup("engine"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("dbuser", rootCmd.PersistentFlags().Lookup("dbuser"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	viper.SetDefault("engine", "postgres")
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("dbuser", "root")
	viper.SetDefault("port", "5432")
	viper.SetDefault("password", "password")
	viper.SetDefault("database", "mydb")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
