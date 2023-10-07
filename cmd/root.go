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

var cfgFile string

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "試験データの設定ファイルを指定します")
	rootCmd.PersistentFlags().BoolP("help", "", false, "コマンドのヘルプを表示します")
	rootCmd.PersistentFlags().StringP("host", "h", "localhost", "DBのホスト名を指定します")
	rootCmd.PersistentFlags().StringP("dbuser", "u", "root", "DBのユーザ名を指定します")
	rootCmd.PersistentFlags().StringP("port", "P", "5432", "DBのポート番号を指定します")
	rootCmd.PersistentFlags().StringP("password", "p", "password", "DBのログインパスワードを指定します")
	rootCmd.PersistentFlags().StringP("database", "d", "mydb", "データベース名を指定します")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("dbuser", rootCmd.PersistentFlags().Lookup("dbuser"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	viper.SetDefault("host", "localhost")
	viper.SetDefault("dbuser", "root")
	viper.SetDefault("port", "5432")
	viper.SetDefault("password", "password")
	viper.SetDefault("database", "mydb")

	// rootCmd.MarkFlagRequired("config")

	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
