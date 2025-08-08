package cmd

import (
	"github.com/puppetma4ster/koyane-framework/internal/core/analyzer"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	all     bool
	general bool
	content bool
)

var analyzeCmd = &cobra.Command{
	Use:   output.AnalyzeHelpTexts["use"],
	Short: output.AnalyzeHelpTexts["short"],
	Long:  output.AnalyzeHelpTexts["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]

		if general {
			content := analyzer.NewContentDummy()
			general, err := analyzer.NewGeneralAnalyzer(inputPath)
			if err != nil {
				return
			}
			printer := output.NewAnalyzePrinter(general, content)
			printer.PrintAllGeneralInfo()
			printer.FlushGeneral()
		}
		if content {
			general := analyzer.NewGeneralDummy()
			content, err := analyzer.NewAnalyzerContent(inputPath, true, true, true,
				false, true, true, false)
			if err != nil {
				return
			}
			printer := output.NewAnalyzePrinter(general, content)
			printer.PrintAllContentInfo()
			printer.FlushContent()

		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().BoolVarP(&all, "all", "a", false, output.AnalyzeHelpTexts["all"])
	analyzeCmd.Flags().BoolVarP(&general, "general", "g", false, output.AnalyzeHelpTexts["generate"])
	analyzeCmd.Flags().BoolVarP(&content, "content", "c", false, output.AnalyzeHelpTexts["content"])

}
