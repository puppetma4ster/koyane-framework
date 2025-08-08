/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   output.GenerateRootHelpTexts["use"],
	Short: output.GenerateRootHelpTexts["short"],
	Long:  output.GenerateRootHelpTexts["long"],
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help() // printing help when no command is specified
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) { // is always executed regardless of command / flag
		output.PrintStatus("statusRoot", "generateTemp") //Temp path management
		err := utils.CreateTempDir()
		if err != nil { // creates temp folder to /tmp/koyane_framework_tmp
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.koyane-framework.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	figure.NewFigure("KYF", "slant", true).Print()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
