/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atc",
	Short: "Atcoder command line tool",
	// The number of version is set dynamically later.
	Version: "v0.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setVersion() {
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
}

func init() {
	setVersion()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.atc.toml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	rootCmd.SetVersionTemplate(version)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		config, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".atc" (without extension).
		viper.AddConfigPath(config)
		viper.SetConfigType("toml")
		viper.SetConfigName(".atc")
	}

	viper.AutomaticEnv() // read in environment variables that match
}
