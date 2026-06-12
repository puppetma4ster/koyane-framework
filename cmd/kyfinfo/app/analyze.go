package app

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/analyzer"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	allArg         bool
	generalArg     bool
	contentArg     bool
	statsArg       bool
	saveFileArg    string
	withDuplicates bool
)

var analyzeCmd = &cobra.Command{
	Use:   output.AnalyzeHelpTexts["use"],
	Short: output.AnalyzeHelpTexts["short"],
	Long:  output.AnalyzeHelpTexts["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]

		stop := make(chan struct{})
		go output.Spinner("Analyze File", stop)

		if allArg || generalArg && contentArg && statsArg {
			general, err := analyzer.NewGeneralAnalyzer(inputPath)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
			isSaved, err := analyzer.IsSavedObj(general)
			var content *analyzer.AnalyzerContent
			if isSaved { // checking if wordlist is not changed and already analyzed
				general, content, err = analyzer.LoadFromYamlAll(*general) // loading wordlist analytics
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			} else { // when file is not analyzed yet
				content, err = analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true, true, true, true, true)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			}
			printer := output.NewAnalyzePrinter(general, content)
			printer.PrintAllGeneralInfo()
			printer.PrintAllContentInfo()
			printer.PrintAllStatsInfo(withDuplicates)

			err = analyzer.SaveToYamlAll(general, content) // saveing the values from the wordlist
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
			printer.FlushGeneral()
			printer.FlushContent()
			printer.FlushStats()
		} else if statsArg && contentArg {
			general := analyzer.NewGeneralDummy()
			content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true,
				true, true, true, true)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}

			printer := output.NewAnalyzePrinter(general, content)

			printer.PrintAllContentInfo()
			printer.PrintAllStatsInfo(withDuplicates)

			printer.FlushContent()
			printer.FlushStats()
		} else {
			if generalArg {
				content := analyzer.NewContentDummy()
				general, err := analyzer.NewGeneralAnalyzer(inputPath)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}

				printer := output.NewAnalyzePrinter(general, content)

				printer.PrintAllGeneralInfo()

				printer.FlushGeneral()
			}
			if contentArg {
				general := analyzer.NewGeneralDummy()
				content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true,
					false, true, true, false)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				printer := output.NewAnalyzePrinter(general, content)

				printer.PrintAllContentInfo()

				printer.FlushContent()

			}
			if statsArg {
				general := analyzer.NewGeneralDummy()
				content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, false, false,
					true, false, true, true)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				printer := output.NewAnalyzePrinter(general, content)
				printer.PrintAllStatsInfo(withDuplicates)
				printer.FlushStats()
			}
		}
		close(stop)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := analyzeCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	analyzeCmd.Flags().BoolVarP(&allArg, "all", "a", false, output.AnalyzeHelpTexts["all"])
	analyzeCmd.Flags().BoolVarP(&generalArg, "general", "g", false, output.AnalyzeHelpTexts["generate"])
	analyzeCmd.Flags().BoolVarP(&contentArg, "content", "c", false, output.AnalyzeHelpTexts["content"])
	analyzeCmd.Flags().BoolVarP(&statsArg, "stats", "s", false, output.AnalyzeHelpTexts["stats"])
	analyzeCmd.Flags().BoolVarP(&withDuplicates, "duplicates", "d", false, output.AnalyzeHelpTexts["duplicates"])
	analyzeCmd.Flags().StringVarP(&saveFileArg, "save-file", "O", "", output.AnalyzeHelpTexts["saveFile"])

}
