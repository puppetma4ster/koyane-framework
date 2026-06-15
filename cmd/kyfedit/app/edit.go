package app

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/editor"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	inputFilePath     string
	outputFilePath    string
	sortArg           bool
	filterMaskArg     string
	filterRangeArg    string
	filterRegExArg    string
	europeanArg       bool
	deleteOriginalArg bool
	stdout            bool
	quiet             bool
)
var editCmd = &cobra.Command{
	Use:   output.GenerateEditHelpTexts["use"],
	Short: output.GenerateEditHelpTexts["short"],
	Long:  output.GenerateEditHelpTexts["long"],
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// print banner
		if !quiet && !stdout {
			output.PrintBanner()

		}
		stop := make(chan struct{})
		if !stdout {
			go output.Spinner("Edit File", stop)
		}
		err := editor.EditWordlist(inputFilePath, outputFilePath, filterRangeArg, filterMaskArg, filterRegExArg)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		close(stop)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := editCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {

	editCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "input file path")

	editCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "output file path")

	editCmd.Flags().BoolVar(&stdout, "stdout", false, "sort words")
	editCmd.Flags().BoolVarP(&sortArg, "sort", "s", false, output.GenerateEditHelpTexts["sort"])
	editCmd.Flags().StringVarP(&filterMaskArg, "filter-mask", "m", "", output.GenerateEditHelpTexts["removeMask"])
	editCmd.Flags().StringVarP(&filterRangeArg, "filter-range", "r", "", output.GenerateEditHelpTexts["removeRange"])
	editCmd.Flags().StringVarP(&filterRegExArg, "filter-chars", "x", "", output.GenerateEditHelpTexts["removeChars"])
	editCmd.Flags().BoolVar(&europeanArg, "european", false, output.GenerateEditHelpTexts["european"])
	editCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, output.GenerateEditHelpTexts["quiet"])
	editCmd.Flags().BoolVarP(&deleteOriginalArg, "delete", "d", false, output.GenerateEditHelpTexts["delete"])

}
